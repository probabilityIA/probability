'use client';

import { useState, useEffect, useMemo } from 'react';
import { XMarkIcon, MagnifyingGlassIcon } from '@heroicons/react/24/outline';
import { CookieStorage } from '@/shared/utils/cookie-storage';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import {
  getInvoiceableOrdersAction,
  createBulkInvoicesAction,
} from '../../infra/actions';
import type { InvoiceableOrder } from '../../domain/types';

interface Props {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

interface Business {
  id: number;
  name: string;
}

export function BulkCreateInvoiceModal({ isOpen, onClose, onSuccess }: Props) {
  const [orders, setOrders] = useState<InvoiceableOrder[]>([]);
  const [selectedOrderIds, setSelectedOrderIds] = useState<string[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasLoadedOnce, setHasLoadedOnce] = useState(false);

  // Super admin filters
  const [isSuperAdmin, setIsSuperAdmin] = useState(false);
  const [businesses, setBusinesses] = useState<Business[]>([]);
  const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
  const [loadingBusinesses, setLoadingBusinesses] = useState(false);

  useEffect(() => {
    if (isOpen && !hasLoadedOnce) {
      checkIfSuperAdmin();
      loadOrders();
      setHasLoadedOnce(true);
    } else if (!isOpen) {
      // Reset state cuando se cierra
      setOrders([]);
      setSelectedOrderIds([]);
      setSearchQuery('');
      setError(null);
      setHasLoadedOnce(false);
      setSelectedBusinessId(null);
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
    setSelectedOrderIds((prev) =>
      prev.includes(orderId)
        ? prev.filter((id) => id !== orderId)
        : [...prev, orderId]
    );
  };

  const handleToggleAll = () => {
    if (selectedOrderIds.length === filteredOrders.length) {
      setSelectedOrderIds([]);
    } else {
      setSelectedOrderIds(filteredOrders.map((order) => order.id));
    }
  };

  const handleSubmit = async () => {
    if (selectedOrderIds.length === 0) {
      alert('Selecciona al menos una orden');
      return;
    }

    setSubmitting(true);
    try {
      const result = await createBulkInvoicesAction({
        order_ids: selectedOrderIds,
      });

      if (result.created > 0) {
        const message = result.failed === 0
          ? `✓ ${result.created} factura(s) creada(s) exitosamente`
          : `✓ ${result.created} creada(s), ${result.failed} fallida(s)`;

        alert(message);
        onSuccess();
        onClose();
      } else {
        alert(`No se crearon facturas. ${result.failed} orden(es) fallaron`);
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Error al crear facturas');
    } finally {
      setSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-center justify-center p-4">
        {/* Backdrop */}
        <div
          className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
          onClick={onClose}
        />

        {/* Modal */}
        <div className="relative bg-white rounded-lg shadow-xl max-w-3xl w-full max-h-[90vh] flex flex-col">
          {/* Header */}
          <div className="flex items-center justify-between p-6 border-b">
            <h2 className="text-2xl font-bold">Crear Facturas desde Órdenes</h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600"
            >
              <XMarkIcon className="w-6 h-6" />
            </button>
          </div>

          {/* Content */}
          <div className="flex-1 overflow-y-auto p-6">
            {loading ? (
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
                    onClick={loadOrders}
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

                {/* Filtro de Business (solo Super Admin) */}
                {isSuperAdmin && (
                  <div className="mb-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Filtrar por Negocio
                    </label>
                    <select
                      value={selectedBusinessId ?? ''}
                      onChange={(e) => {
                        const value = e.target.value;
                        handleBusinessFilterChange(value === '' ? null : parseInt(value));
                      }}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      disabled={loadingBusinesses}
                    >
                      <option value="">Todos los negocios</option>
                      {businesses.map((business) => (
                        <option key={business.id} value={business.id}>
                          {business.name}
                        </option>
                      ))}
                    </select>
                    {loadingBusinesses && (
                      <p className="text-xs text-gray-500 mt-1">Cargando negocios...</p>
                    )}
                  </div>
                )}

                {/* Búsqueda */}
                <div className="mb-4 relative">
                  <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
                  <input
                    type="text"
                    placeholder="Buscar orden por número o cliente..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
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

                {/* Orders List */}
                <div className="border border-gray-200 rounded max-h-96 overflow-y-auto">
                  {filteredOrders.map((order) => (
                    <div
                      key={order.id}
                      onClick={() => handleToggleOrder(order.id)}
                      className={`p-4 border-b last:border-b-0 cursor-pointer hover:bg-gray-50 transition-colors ${
                        selectedOrderIds.includes(order.id)
                          ? 'bg-blue-50 border-blue-200'
                          : ''
                      }`}
                    >
                      <div className="flex items-start gap-3">
                        <input
                          type="checkbox"
                          checked={selectedOrderIds.includes(order.id)}
                          onChange={() => {}}
                          className="mt-1 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                        />
                        <div className="flex-1">
                          <div className="flex items-center justify-between mb-1">
                            <span className="font-medium text-gray-900">
                              {order.order_number}
                            </span>
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
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="mt-4 text-sm text-gray-600">
                  {selectedOrderIds.length === 0
                    ? 'Ninguna orden seleccionada'
                    : `${selectedOrderIds.length} orden(es) seleccionada(s)`}
                </div>
              </>
            )}
          </div>

          {/* Footer */}
          <div className="flex justify-end gap-3 p-6 border-t">
            <button
              onClick={onClose}
              disabled={submitting}
              className="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
            >
              Cancelar
            </button>
            <button
              onClick={handleSubmit}
              disabled={submitting || selectedOrderIds.length === 0 || loading}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {submitting
                ? 'Creando facturas...'
                : `Crear ${selectedOrderIds.length} Factura(s)`}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
