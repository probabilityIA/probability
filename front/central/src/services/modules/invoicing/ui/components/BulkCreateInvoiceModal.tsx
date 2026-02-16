'use client';

import { useState, useEffect, useMemo, useCallback, useRef } from 'react';
import { XMarkIcon, MagnifyingGlassIcon, CheckCircleIcon, XCircleIcon, ExclamationTriangleIcon } from '@heroicons/react/24/outline';
import { CookieStorage } from '@/shared/utils/cookie-storage';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import {
  getInvoiceableOrdersAction,
  createBulkInvoicesAction,
} from '../../infra/actions';
import { useInvoiceSSE } from '../hooks/useInvoiceSSE';
import type { InvoiceableOrder, InvoiceSSEEventData } from '../../domain/types';

interface Props {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  businessId?: number;
}

interface Business {
  id: number;
  name: string;
}

export function BulkCreateInvoiceModal({ isOpen, onClose, onSuccess, businessId: propBusinessId }: Props) {
  const [orders, setOrders] = useState<InvoiceableOrder[]>([]);
  const [selectedOrderIds, setSelectedOrderIds] = useState<string[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasLoadedOnce, setHasLoadedOnce] = useState(false);

  // Bulk job progress tracking (SSE)
  const [bulkProgress, setBulkProgress] = useState<InvoiceSSEEventData | null>(null);
  const [bulkCompleted, setBulkCompleted] = useState(false);
  const submittingTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Individual order processing status
  interface OrderProcessingStatus {
    status: 'pending' | 'processing' | 'success' | 'failed';
    error_message?: string;
    invoice_id?: number;
  }
  const [orderStatuses, setOrderStatuses] = useState<Record<string, OrderProcessingStatus>>({});

  // Super admin filters
  const [isSuperAdmin, setIsSuperAdmin] = useState(false);
  const [businesses, setBusinesses] = useState<Business[]>([]);
  const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
  const [loadingBusinesses, setLoadingBusinesses] = useState(false);
  const [showBusinessAlert, setShowBusinessAlert] = useState(false);

  // Super admin debe seleccionar un negocio antes de operar
  const superAdminNeedsBusiness = isSuperAdmin && !selectedBusinessId;

  // Get businessId for SSE connection
  const currentBusinessId = propBusinessId ?? selectedBusinessId ?? 0;

  // SSE: Real-time progress for bulk invoice jobs
  const handleBulkJobProgress = useCallback((data: InvoiceSSEEventData) => {
    setBulkProgress(data);
  }, []);

  const handleBulkJobCompleted = useCallback((data: InvoiceSSEEventData) => {
    // Clear timeout when SSE event arrives
    if (submittingTimeoutRef.current) {
      clearTimeout(submittingTimeoutRef.current);
      submittingTimeoutRef.current = null;
    }
    setBulkProgress(data);
    setBulkCompleted(true);
    setSubmitting(false);
  }, []);

  // Individual invoice events during bulk job
  const handleInvoiceCreated = useCallback((data: InvoiceSSEEventData) => {
    if (data.order_id && selectedOrderIds.includes(data.order_id)) {
      setOrderStatuses((prev) => ({
        ...prev,
        [data.order_id!]: {
          status: 'success',
          invoice_id: data.invoice_id,
        },
      }));
    }
  }, [selectedOrderIds]);

  const handleInvoiceFailed = useCallback((data: InvoiceSSEEventData) => {
    if (data.order_id && selectedOrderIds.includes(data.order_id)) {
      setOrderStatuses((prev) => ({
        ...prev,
        [data.order_id!]: {
          status: 'failed',
          error_message: data.error_message,
        },
      }));
    }
  }, [selectedOrderIds]);

  useInvoiceSSE({
    businessId: currentBusinessId,
    onBulkJobProgress: handleBulkJobProgress,
    onBulkJobCompleted: handleBulkJobCompleted,
    onInvoiceCreated: handleInvoiceCreated,
    onInvoiceFailed: handleInvoiceFailed,
  });

  useEffect(() => {
    if (isOpen && !hasLoadedOnce) {
      checkIfSuperAdmin();
      // Solo cargar ordenes al inicio si NO es super admin
      // Super admin debe seleccionar un negocio primero
      const user = CookieStorage.getUser();
      if (!user?.is_super_admin) {
        loadOrders();
      }
      setHasLoadedOnce(true);
    } else if (!isOpen) {
      // Reset state cuando se cierra
      setOrders([]);
      setSelectedOrderIds([]);
      setSearchQuery('');
      setError(null);
      setHasLoadedOnce(false);
      setSelectedBusinessId(null);
      setBulkProgress(null);
      setBulkCompleted(false);
      setSubmitting(false);
      setOrderStatuses({});
      setShowBusinessAlert(false);
      // Clear timeout if modal closes
      if (submittingTimeoutRef.current) {
        clearTimeout(submittingTimeoutRef.current);
        submittingTimeoutRef.current = null;
      }
    }
  }, [isOpen, hasLoadedOnce]);

  const checkIfSuperAdmin = async () => {
    const user = CookieStorage.getUser();
    if (user?.is_super_admin) {
      setIsSuperAdmin(true);
      await loadBusinesses();
    }
  };

  const loadBusinesses = async () => {
    setLoadingBusinesses(true);
    try {
      const result = await getBusinessesAction();
      const businessList = result.data.map((b: any) => ({
        id: b.id,
        name: b.name,
      }));
      setBusinesses(businessList);
    } catch (err) {
      console.error('Error loading businesses:', err);
    } finally {
      setLoadingBusinesses(false);
    }
  };

  const loadOrders = async (businessId?: number | null) => {
    setLoading(true);
    setError(null);
    try {
      const result = await getInvoiceableOrdersAction(1, 100, businessId ?? undefined);
      setOrders(result.data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Error al cargar las órdenes facturables';
      setError(errorMessage);
      console.error('Error loading invoiceable orders:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleBusinessFilterChange = (businessId: number | null) => {
    setSelectedBusinessId(businessId);
    setSelectedOrderIds([]); // Reset selection
    loadOrders(businessId);
  };

  const formatCurrency = (amount: number, currency: string) => {
    return new Intl.NumberFormat('es-CO', {
      style: 'currency',
      currency: currency || 'COP',
    }).format(amount);
  };

  // Filtrar órdenes por búsqueda
  const filteredOrders = useMemo(() => {
    if (!searchQuery.trim()) return orders;

    const query = searchQuery.toLowerCase();
    return orders.filter(
      (order) =>
        order.order_number.toLowerCase().includes(query) ||
        order.customer_name.toLowerCase().includes(query)
    );
  }, [orders, searchQuery]);

  const handleToggleOrder = (orderId: string) => {
    if (superAdminNeedsBusiness) {
      setShowBusinessAlert(true);
      return;
    }
    setSelectedOrderIds((prev) =>
      prev.includes(orderId)
        ? prev.filter((id) => id !== orderId)
        : [...prev, orderId]
    );
  };

  const handleToggleAll = () => {
    if (superAdminNeedsBusiness) {
      setShowBusinessAlert(true);
      return;
    }
    if (selectedOrderIds.length === filteredOrders.length) {
      setSelectedOrderIds([]);
    } else {
      setSelectedOrderIds(filteredOrders.map((order) => order.id));
    }
  };

  const handleSubmit = async () => {
    if (superAdminNeedsBusiness) {
      setShowBusinessAlert(true);
      return;
    }
    if (selectedOrderIds.length === 0) {
      alert('Selecciona al menos una orden');
      return;
    }

    setSubmitting(true);
    setBulkProgress(null);
    setBulkCompleted(false);

    // Resetear estados de órdenes a "pending"
    const initialStatuses: Record<string, OrderProcessingStatus> = {};
    selectedOrderIds.forEach((orderId) => {
      initialStatuses[orderId] = { status: 'pending' };
    });
    setOrderStatuses(initialStatuses);

    try {
      await createBulkInvoicesAction({
        order_ids: selectedOrderIds,
        ...(isSuperAdmin && selectedBusinessId ? { business_id: selectedBusinessId } : {}),
      });

      // Async job started - SSE will track progress in real-time
      // Set a fallback timeout in case SSE events don't arrive (connection issues, etc.)
      submittingTimeoutRef.current = setTimeout(() => {
        console.warn('SSE bulk_job.completed event did not arrive within 30 seconds - resetting UI state');
        setSubmitting(false);
        setBulkCompleted(true);
        // Show a message to user that they should check the invoice list
        alert('El proceso de facturación se ha iniciado. Revisa la lista de facturas para ver el resultado.');
      }, 30000); // 30 seconds timeout

      // Don't close modal yet, let the progress bar show
    } catch (err) {
      console.error('Error creating bulk invoices:', err);
      alert(err instanceof Error ? err.message : 'Error al crear facturas');
      setSubmitting(false);
      setBulkProgress(null);
      setBulkCompleted(false);
    }
  };

  const handleCloseAfterCompletion = () => {
    onSuccess();
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-center justify-center p-4">
        {/* Backdrop con blur */}
        <div
          className="fixed inset-0 backdrop-blur-sm bg-white/10 transition-opacity"
          onClick={onClose}
        />

        {/* Modal 95% de la pantalla */}
        <div className="relative bg-white rounded-xl shadow-2xl w-[95vw] h-[95vh] flex flex-col overflow-hidden">
          {/* Header con gradiente morado */}
          <div className="flex items-center justify-between p-8 bg-gradient-to-r from-[#7c3aed] to-[#6d28d9]">
            <div>
              <h2 className="text-3xl font-bold text-white">Crear Facturas</h2>
              <p className="text-purple-100 text-sm mt-1">Selecciona las órdenes para generar facturas electrónicas</p>
            </div>
            <div className="flex items-center gap-4">
              {/* Contador de órdenes seleccionadas */}
              <div className="bg-white/20 backdrop-blur-sm rounded-full px-4 py-2">
                <span className="text-white font-semibold text-lg">
                  {selectedOrderIds.length} seleccionada{selectedOrderIds.length !== 1 ? 's' : ''}
                </span>
              </div>
              <button
                onClick={onClose}
                className="text-white hover:bg-white/20 rounded-full p-2 transition-all duration-200"
              >
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>
          </div>

          {/* Content */}
          <div className="flex-1 overflow-y-auto p-8 bg-white">
            {/* Filtro de Business (solo Super Admin) - siempre visible arriba */}
            {isSuperAdmin && (
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Seleccionar Negocio <span className="text-red-500">*</span>
                </label>
                <select
                  value={selectedBusinessId ?? ''}
                  onChange={(e) => {
                    const value = e.target.value;
                    setShowBusinessAlert(false);
                    handleBusinessFilterChange(value === '' ? null : parseInt(value));
                  }}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    superAdminNeedsBusiness ? 'border-amber-400 bg-amber-50' : 'border-gray-300'
                  }`}
                  disabled={loadingBusinesses}
                >
                  <option value="">-- Selecciona un negocio --</option>
                  {businesses.map((business) => (
                    <option key={business.id} value={business.id}>
                      {business.name}
                    </option>
                  ))}
                </select>
                {loadingBusinesses && (
                  <p className="text-xs text-gray-500 mt-1">Cargando negocios...</p>
                )}
                {superAdminNeedsBusiness && (
                  <div className="mt-3 flex items-start gap-2 p-3 bg-amber-50 border border-amber-300 rounded-lg">
                    <ExclamationTriangleIcon className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
                    <p className="text-sm text-amber-800">
                      Debes seleccionar un negocio antes de poder seleccionar ordenes y facturar.
                    </p>
                  </div>
                )}
              </div>
            )}

            {/* Contenido principal: ordenes */}
            {superAdminNeedsBusiness ? (
              <div className="py-8 text-center text-gray-400">
                <ExclamationTriangleIcon className="w-12 h-12 mx-auto mb-3 text-amber-300" />
                <p className="text-gray-500">Selecciona un negocio para ver las ordenes facturables</p>
              </div>
            ) : loading ? (
              <div className="flex flex-col items-center justify-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
                <p className="mt-4 text-gray-600">Cargando órdenes facturables...</p>
              </div>
            ) : error ? (
              <div className="py-8 text-center">
                <p className="text-red-600 mb-4 font-medium">{error}</p>
                {error.includes('authentication') || error.includes('login') ? (
                  <div className="space-y-3">
                    <p className="text-sm text-gray-600">
                      Tu sesión ha expirado. Por favor, inicia sesión nuevamente.
                    </p>
                    <button
                      onClick={() => {
                        window.location.href = '/login';
                      }}
                      className="px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded"
                    >
                      Ir a Login
                    </button>
                  </div>
                ) : (
                  <button
                    onClick={() => loadOrders(selectedBusinessId)}
                    className="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded"
                  >
                    Reintentar
                  </button>
                )}
              </div>
            ) : orders.length === 0 ? (
              <div className="py-8 text-center text-gray-500">
                <p className="mb-2">No hay órdenes facturables disponibles</p>
                <p className="text-sm">
                  Las órdenes deben estar marcadas como facturables y no tener
                  factura previa
                </p>
              </div>
            ) : (
              <>
                <p className="text-sm text-gray-600 mb-4">
                  Selecciona las órdenes para las cuales deseas crear facturas
                  electrónicas. Cada orden se procesará individualmente.
                </p>

                {/* Búsqueda */}
                <div className="mb-6 relative">
                  <MagnifyingGlassIcon className="absolute left-4 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
                  <input
                    type="text"
                    placeholder="Buscar orden por número o cliente..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="w-full pl-12 pr-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed] transition-all duration-200"
                  />
                </div>

                {/* Select All */}
                <div className="mb-2 flex items-center gap-2 p-3 bg-gray-50 rounded">
                  <input
                    type="checkbox"
                    checked={
                      filteredOrders.length > 0 &&
                      selectedOrderIds.length === filteredOrders.length
                    }
                    onChange={handleToggleAll}
                    className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                  <span className="text-sm font-medium">
                    Seleccionar todas ({filteredOrders.length})
                  </span>
                </div>

                {/* Orders List con scroll suave */}
                <div className="border border-gray-200 rounded-lg max-h-96 overflow-y-auto scroll-smooth">
                  {filteredOrders.map((order) => {
                    const orderStatus = orderStatuses[order.id];

                    // Color de fondo según estado
                    const getBgColor = () => {
                      if (!submitting && !bulkCompleted) {
                        return selectedOrderIds.includes(order.id) ? 'bg-blue-50 border-blue-200' : '';
                      }
                      if (!orderStatus || orderStatus.status === 'pending') return 'bg-gray-50';
                      if (orderStatus.status === 'processing') return 'bg-yellow-50 border-yellow-200';
                      if (orderStatus.status === 'success') return 'bg-green-50 border-green-200';
                      if (orderStatus.status === 'failed') return 'bg-red-50 border-red-200';
                      return '';
                    };

                    // Ícono de estado
                    const StatusIcon = () => {
                      if (!orderStatus) return null;
                      if (orderStatus.status === 'processing') {
                        return <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-yellow-600" />;
                      }
                      if (orderStatus.status === 'success') {
                        return <CheckCircleIcon className="w-5 h-5 text-green-600" />;
                      }
                      if (orderStatus.status === 'failed') {
                        return <XCircleIcon className="w-5 h-5 text-red-600" />;
                      }
                      return null;
                    };

                    return (
                      <div
                        key={order.id}
                        onClick={() => !submitting && !bulkCompleted && handleToggleOrder(order.id)}
                        className={`p-4 border-b last:border-b-0 ${
                          submitting || bulkCompleted ? '' : 'cursor-pointer hover:bg-gray-50'
                        } transition-colors ${getBgColor()}`}
                      >
                        <div className="flex items-start gap-3">
                          <input
                            type="checkbox"
                            checked={selectedOrderIds.includes(order.id)}
                            disabled={submitting || bulkCompleted}
                            onChange={() => {}}
                            className="mt-1 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
                          />
                          <div className="flex-1">
                            <div className="flex items-center justify-between mb-1">
                              <div className="flex items-center gap-2">
                                <span className="font-medium text-gray-900">
                                  {order.order_number}
                                </span>
                                <StatusIcon />
                              </div>
                              <span className="font-semibold text-gray-900">
                                {formatCurrency(order.total_amount, order.currency)}
                              </span>
                            </div>
                            <p className="text-sm text-gray-600">
                              {order.customer_name}
                            </p>
                            <div className="flex items-center gap-2 mt-1">
                              <p className="text-xs text-gray-500">
                                {new Date(order.created_at).toLocaleDateString('es-CO', {
                                  year: 'numeric',
                                  month: 'short',
                                  day: 'numeric',
                                })}
                              </p>
                              {/* Mostrar Business ID para super admin */}
                              {isSuperAdmin && (
                                <>
                                  <span className="text-xs text-gray-400">•</span>
                                  <span className="text-xs text-blue-600 font-medium">
                                    Business #{order.business_id}
                                  </span>
                                </>
                              )}
                            </div>

                            {/* Resultado: siempre visible cuando hay status */}
                            {orderStatus?.status === 'success' && (
                              <div className="mt-2 p-2.5 bg-green-100 border border-green-300 rounded-lg text-xs text-green-800">
                                <div className="flex items-start gap-2">
                                  <CheckCircleIcon className="w-4 h-4 text-green-600 flex-shrink-0 mt-0.5" />
                                  <div>
                                    <p className="font-semibold">Factura creada exitosamente</p>
                                    {orderStatus.invoice_id && (
                                      <p className="mt-0.5 text-green-700">ID: {orderStatus.invoice_id}</p>
                                    )}
                                  </div>
                                </div>
                              </div>
                            )}
                            {orderStatus?.status === 'failed' && (
                              <div className="mt-2 p-2.5 bg-red-100 border border-red-300 rounded-lg text-xs text-red-800">
                                <div className="flex items-start gap-2">
                                  <XCircleIcon className="w-4 h-4 text-red-600 flex-shrink-0 mt-0.5" />
                                  <div className="min-w-0 flex-1">
                                    <p className="font-semibold">Error al crear factura</p>
                                    {orderStatus.error_message && (
                                      <p className="mt-0.5 text-red-700 break-words whitespace-pre-wrap">{orderStatus.error_message}</p>
                                    )}
                                  </div>
                                </div>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>

                <div className="mt-4 text-sm text-gray-600">
                  {selectedOrderIds.length === 0
                    ? 'Ninguna orden seleccionada'
                    : `${selectedOrderIds.length} orden(es) seleccionada(s)`}
                </div>
              </>
            )}
          </div>

          {/* Bulk Progress Bar (SSE real-time) */}
          {(submitting || bulkCompleted) && bulkProgress && (
            <div className="px-6 py-4 border-t bg-gray-50">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-gray-700">
                  {bulkCompleted ? 'Procesamiento completado' : 'Procesando facturas...'}
                </span>
                <span className="text-sm text-gray-500">
                  {bulkProgress.processed ?? 0}/{bulkProgress.total_orders ?? 0}
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-3 overflow-hidden flex">
                {bulkCompleted || (bulkProgress.successful ?? 0) > 0 || (bulkProgress.failed ?? 0) > 0 ? (
                  <>
                    {/* Segmento verde (exitosas) */}
                    {(bulkProgress.successful ?? 0) > 0 && (
                      <div
                        className="h-3 bg-green-500 transition-all duration-300"
                        style={{
                          width: `${((bulkProgress.successful ?? 0) / (bulkProgress.total_orders || 1)) * 100}%`,
                        }}
                      />
                    )}
                    {/* Segmento rojo (fallidas) */}
                    {(bulkProgress.failed ?? 0) > 0 && (
                      <div
                        className="h-3 bg-red-500 transition-all duration-300"
                        style={{
                          width: `${((bulkProgress.failed ?? 0) / (bulkProgress.total_orders || 1)) * 100}%`,
                        }}
                      />
                    )}
                  </>
                ) : (
                  /* Barra azul mientras procesa (sin resultados aún) */
                  <div
                    className="h-3 bg-blue-600 rounded-full transition-all duration-300"
                    style={{ width: `${Math.min(bulkProgress.progress ?? 0, 100)}%` }}
                  />
                )}
              </div>
              <div className="flex gap-4 mt-2 text-xs text-gray-500">
                <span className="text-green-600">
                  Exitosas: {bulkProgress.successful ?? 0}
                </span>
                <span className="text-red-600">
                  Fallidas: {bulkProgress.failed ?? 0}
                </span>
              </div>
            </div>
          )}

          {/* Footer */}
          <div className="flex justify-end gap-3 p-6 border-t bg-gray-50">
            {bulkCompleted ? (
              <button
                onClick={handleCloseAfterCompletion}
                className="px-6 py-2.5 bg-gradient-to-r from-green-600 to-green-700 text-white font-semibold rounded-full hover:shadow-lg transition-all duration-200 hover:scale-105"
              >
                Cerrar
              </button>
            ) : (
              <>
                <button
                  onClick={onClose}
                  disabled={submitting}
                  className="px-6 py-2.5 border-2 border-gray-300 text-gray-700 font-semibold rounded-full hover:bg-gray-100 disabled:opacity-50 transition-all duration-200"
                >
                  Cancelar
                </button>
                <button
                  onClick={handleSubmit}
                  disabled={submitting || selectedOrderIds.length === 0 || loading || superAdminNeedsBusiness}
                  className="px-6 py-2.5 bg-gradient-to-r from-[#7c3aed] to-[#6d28d9] text-white font-semibold rounded-full hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 hover:scale-105"
                >
                  {submitting
                    ? 'Creando facturas...'
                    : `Crear ${selectedOrderIds.length} Factura(s)`}
                </button>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Modal de alerta: Super admin debe seleccionar negocio */}
      {showBusinessAlert && (
        <div className="fixed inset-0 z-[60] flex items-center justify-center">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-sm" onClick={() => setShowBusinessAlert(false)} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-md w-full mx-4 overflow-hidden">
            <div className="bg-amber-500 p-4 flex items-center gap-3">
              <ExclamationTriangleIcon className="w-8 h-8 text-white" />
              <h3 className="text-lg font-bold text-white">Negocio requerido</h3>
            </div>
            <div className="p-6">
              <p className="text-gray-700">
                Como super administrador, debes seleccionar un <strong>negocio</strong> antes de poder seleccionar ordenes y crear facturas.
              </p>
              <p className="text-sm text-gray-500 mt-2">
                Cada factura se asocia a un negocio especifico. Selecciona uno del listado para continuar.
              </p>
            </div>
            <div className="flex justify-end p-4 border-t bg-gray-50">
              <button
                onClick={() => setShowBusinessAlert(false)}
                className="px-5 py-2 bg-amber-500 text-white font-semibold rounded-full hover:bg-amber-600 transition-all duration-200"
              >
                Entendido
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
