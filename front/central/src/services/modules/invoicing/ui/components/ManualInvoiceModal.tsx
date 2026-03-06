'use client';

import { useState, useEffect, useMemo, useRef } from 'react';
import { XMarkIcon, MagnifyingGlassIcon, DocumentTextIcon, CheckCircleIcon } from '@heroicons/react/24/outline';
import { CookieStorage } from '@/shared/utils/cookie-storage';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import { getInvoiceableOrdersAction, registerManualInvoiceAction } from '../../infra/actions';
import type { InvoiceableOrder } from '../../domain/types';

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

export function ManualInvoiceModal({ isOpen, onClose, onSuccess, businessId: propBusinessId }: Props) {
  const [invoiceNumber, setInvoiceNumber] = useState('');
  const [orders, setOrders] = useState<InvoiceableOrder[]>([]);
  const [selectedOrder, setSelectedOrder] = useState<InvoiceableOrder | null>(null);
  const [orderSearch, setOrderSearch] = useState('');
  const [showDropdown, setShowDropdown] = useState(false);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  // Super admin
  const [isSuperAdmin, setIsSuperAdmin] = useState(false);
  const [businesses, setBusinesses] = useState<Business[]>([]);
  const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
  const [loadingBusinesses, setLoadingBusinesses] = useState(false);

  const dropdownRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const effectiveBusinessId = propBusinessId ?? selectedBusinessId ?? undefined;
  const superAdminNeedsBusiness = isSuperAdmin && !selectedBusinessId;

  useEffect(() => {
    if (isOpen) {
      const user = CookieStorage.getUser();
      if (user?.is_super_admin) {
        setIsSuperAdmin(true);
        loadBusinesses();
      } else {
        loadOrders();
      }
    } else {
      // Reset state
      setInvoiceNumber('');
      setOrders([]);
      setSelectedOrder(null);
      setOrderSearch('');
      setShowDropdown(false);
      setError(null);
      setSuccess(false);
      setSelectedBusinessId(null);
    }
  }, [isOpen]);

  // Cerrar dropdown al hacer click fuera
  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setShowDropdown(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const loadBusinesses = async () => {
    setLoadingBusinesses(true);
    try {
      const result = await getBusinessesAction();
      setBusinesses(result.data.map((b: any) => ({ id: b.id, name: b.name })));
    } catch {
      console.error('Error loading businesses');
    } finally {
      setLoadingBusinesses(false);
    }
  };

  const loadOrders = async (businessId?: number | null) => {
    setLoading(true);
    setError(null);
    try {
      const result = await getInvoiceableOrdersAction(1, 200, businessId ?? undefined);
      setOrders(result.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar ordenes');
    } finally {
      setLoading(false);
    }
  };

  const handleBusinessChange = (businessId: number | null) => {
    setSelectedBusinessId(businessId);
    setSelectedOrder(null);
    setOrderSearch('');
    if (businessId) {
      loadOrders(businessId);
    } else {
      setOrders([]);
    }
  };

  const filteredOrders = useMemo(() => {
    if (!orderSearch.trim()) return orders;
    const q = orderSearch.toLowerCase();
    return orders.filter(
      o => o.order_number.toLowerCase().includes(q) || o.customer_name.toLowerCase().includes(q)
    );
  }, [orders, orderSearch]);

  const handleSelectOrder = (order: InvoiceableOrder) => {
    setSelectedOrder(order);
    setOrderSearch('');
    setShowDropdown(false);
  };

  const handleSubmit = async () => {
    if (!invoiceNumber.trim()) {
      setError('Ingresa el numero de factura');
      return;
    }
    if (!selectedOrder) {
      setError('Selecciona una orden');
      return;
    }

    setSubmitting(true);
    setError(null);

    const result = await registerManualInvoiceAction(
      selectedOrder.id,
      invoiceNumber.trim(),
      effectiveBusinessId
    );

    setSubmitting(false);

    if (!result.success) {
      setError(result.error);
      return;
    }

    setSuccess(true);
    // Cerrar después de mostrar éxito brevemente
    setTimeout(() => {
      onSuccess();
      onClose();
    }, 1500);
  };

  const formatCurrency = (amount: number, currency: string) => {
    return new Intl.NumberFormat('es-CO', {
      style: 'currency',
      currency: currency || 'COP',
    }).format(amount);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-center justify-center p-4">
        <div
          className="fixed inset-0 backdrop-blur-sm bg-white/10 transition-opacity"
          onClick={onClose}
        />

        <div className="relative bg-white rounded-xl shadow-2xl w-full max-w-lg">
          {/* Header */}
          <div className="flex items-center justify-between p-6 bg-gradient-to-r from-[#7c3aed] to-[#6d28d9] rounded-t-xl">
            <div className="flex items-center gap-3">
              <DocumentTextIcon className="w-7 h-7 text-white" />
              <div>
                <h2 className="text-xl font-bold text-white">Registrar Factura Manual</h2>
                <p className="text-purple-100 text-xs mt-0.5">Asociar una factura externa a una orden</p>
              </div>
            </div>
            <button
              onClick={onClose}
              className="text-white hover:bg-white/20 rounded-full p-2 transition-all duration-200"
            >
              <XMarkIcon className="w-5 h-5" />
            </button>
          </div>

          {/* Content */}
          <div className="p-6 space-y-5">
            {/* Éxito */}
            {success && (
              <div className="flex flex-col items-center py-6">
                <CheckCircleIcon className="w-16 h-16 text-green-500 mb-3" />
                <p className="text-lg font-semibold text-green-700">Factura registrada exitosamente</p>
                <p className="text-sm text-gray-500 mt-1">La orden ya no aparecera como facturable</p>
              </div>
            )}

            {!success && (
              <>
                {/* Selector de negocio (super admin) */}
                {isSuperAdmin && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1.5">
                      Negocio <span className="text-red-500">*</span>
                    </label>
                    <select
                      value={selectedBusinessId ?? ''}
                      onChange={(e) => handleBusinessChange(e.target.value ? Number(e.target.value) : null)}
                      className={`w-full px-3 py-2.5 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[#7c3aed] ${
                        superAdminNeedsBusiness ? 'border-amber-400 bg-amber-50' : 'border-gray-300'
                      }`}
                      disabled={loadingBusinesses}
                    >
                      <option value="">-- Selecciona un negocio --</option>
                      {businesses.map(b => (
                        <option key={b.id} value={b.id}>{b.name}</option>
                      ))}
                    </select>
                  </div>
                )}

                {/* Número de factura */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1.5">
                    Numero de factura <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    value={invoiceNumber}
                    onChange={(e) => setInvoiceNumber(e.target.value)}
                    placeholder="Ej: FE-12345, SETP990012345"
                    className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed]"
                    disabled={superAdminNeedsBusiness}
                  />
                </div>

                {/* Selector de orden con buscador */}
                <div ref={dropdownRef}>
                  <label className="block text-sm font-medium text-gray-700 mb-1.5">
                    Orden a asociar <span className="text-red-500">*</span>
                  </label>

                  {selectedOrder ? (
                    <div className="flex items-center justify-between p-3 bg-purple-50 border border-purple-200 rounded-lg">
                      <div>
                        <span className="font-medium text-gray-900">{selectedOrder.order_number}</span>
                        <span className="text-gray-500 mx-2">-</span>
                        <span className="text-sm text-gray-600">{selectedOrder.customer_name}</span>
                        <span className="text-gray-400 mx-2">|</span>
                        <span className="text-sm font-medium text-gray-700">
                          {formatCurrency(selectedOrder.total_amount, selectedOrder.currency)}
                        </span>
                      </div>
                      <button
                        onClick={() => setSelectedOrder(null)}
                        className="text-gray-400 hover:text-gray-600 p-1"
                        title="Cambiar orden"
                      >
                        <XMarkIcon className="w-4 h-4" />
                      </button>
                    </div>
                  ) : (
                    <div className="relative">
                      <div className="relative">
                        <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <input
                          ref={searchInputRef}
                          type="text"
                          value={orderSearch}
                          onChange={(e) => {
                            setOrderSearch(e.target.value);
                            setShowDropdown(true);
                          }}
                          onFocus={() => setShowDropdown(true)}
                          placeholder={
                            superAdminNeedsBusiness
                              ? 'Selecciona un negocio primero'
                              : loading
                                ? 'Cargando ordenes...'
                                : 'Buscar por numero de orden o cliente...'
                          }
                          className="w-full pl-9 pr-3 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed]"
                          disabled={superAdminNeedsBusiness || loading}
                        />
                      </div>

                      {/* Dropdown de resultados */}
                      {showDropdown && !superAdminNeedsBusiness && (
                        <div className="absolute z-10 mt-1 w-full bg-white border border-gray-200 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                          {loading ? (
                            <div className="p-4 text-center text-gray-500 text-sm">
                              <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-purple-600 mx-auto mb-2" />
                              Cargando ordenes...
                            </div>
                          ) : filteredOrders.length === 0 ? (
                            <div className="p-4 text-center text-gray-500 text-sm">
                              {orders.length === 0
                                ? 'No hay ordenes facturables disponibles'
                                : 'No se encontraron resultados'}
                            </div>
                          ) : (
                            filteredOrders.slice(0, 50).map(order => (
                              <button
                                key={order.id}
                                onClick={() => handleSelectOrder(order)}
                                className="w-full text-left px-4 py-3 hover:bg-purple-50 border-b last:border-b-0 transition-colors"
                              >
                                <div className="flex items-center justify-between">
                                  <div>
                                    <span className="font-medium text-gray-900">{order.order_number}</span>
                                    <span className="text-gray-400 mx-1.5">-</span>
                                    <span className="text-sm text-gray-600">{order.customer_name}</span>
                                  </div>
                                  <span className="text-sm font-medium text-gray-700">
                                    {formatCurrency(order.total_amount, order.currency)}
                                  </span>
                                </div>
                                <p className="text-xs text-gray-400 mt-0.5">
                                  {new Date(order.created_at).toLocaleDateString('es-CO', {
                                    year: 'numeric', month: 'short', day: 'numeric',
                                  })}
                                  {order.is_paid ? ' - Pagado' : ' - Sin pagar'}
                                </p>
                              </button>
                            ))
                          )}
                        </div>
                      )}
                    </div>
                  )}
                </div>

                {/* Error */}
                {error && (
                  <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">
                    {error}
                  </div>
                )}
              </>
            )}
          </div>

          {/* Footer */}
          {!success && (
            <div className="flex justify-end gap-3 p-6 border-t bg-gray-50 rounded-b-xl">
              <button
                onClick={onClose}
                disabled={submitting}
                className="px-5 py-2.5 border-2 border-gray-300 text-gray-700 font-semibold rounded-full hover:bg-gray-100 disabled:opacity-50 transition-all duration-200"
              >
                Cancelar
              </button>
              <button
                onClick={handleSubmit}
                disabled={submitting || !invoiceNumber.trim() || !selectedOrder || superAdminNeedsBusiness}
                className="px-5 py-2.5 bg-gradient-to-r from-[#7c3aed] to-[#6d28d9] text-white font-semibold rounded-full hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 hover:scale-105"
              >
                {submitting ? 'Registrando...' : 'Registrar Factura'}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
