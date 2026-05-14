'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { XMarkIcon, MagnifyingGlassIcon, CheckCircleIcon, XCircleIcon, ExclamationTriangleIcon, ClockIcon, FunnelIcon, ArrowsUpDownIcon, ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/24/outline';
import { CookieStorage } from '@/shared/utils/cookie-storage';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import {
  getInvoiceableOrdersAction,
  createBulkInvoicesAction,
} from '../../infra/actions';
import { useInvoiceSSE } from '../hooks/useInvoiceSSE';
import type { InvoiceableOrder, InvoiceSSEEventData, InvoiceableOrdersFilters } from '../../domain/types';

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

interface OrderProcessingStatus {
  status: 'pending' | 'processing' | 'success' | 'failed' | 'pending_validation';
  error_message?: string;
  invoice_id?: number;
  invoice_number?: string;
}

type SortKey = NonNullable<InvoiceableOrdersFilters['sortBy']>;
type SortDir = NonNullable<InvoiceableOrdersFilters['sortOrder']>;

const PAGE_SIZE_OPTIONS = [25, 50, 100, 200];

const SORT_OPTIONS: { value: `${SortKey}:${SortDir}`; label: string }[] = [
  { value: 'created_at:desc',   label: 'Mas recientes' },
  { value: 'created_at:asc',    label: 'Mas antiguas' },
  { value: 'total_amount:desc', label: 'Mayor valor' },
  { value: 'total_amount:asc',  label: 'Menor valor' },
  { value: 'order_number:asc',  label: 'Orden # asc' },
  { value: 'order_number:desc', label: 'Orden # desc' },
  { value: 'customer_name:asc', label: 'Cliente A-Z' },
];

export function BulkCreateInvoiceModal({ isOpen, onClose, onSuccess, businessId: propBusinessId }: Props) {
  const [orders, setOrders] = useState<InvoiceableOrder[]>([]);
  const [total, setTotal] = useState(0);
  const [selectedOrderIds, setSelectedOrderIds] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(50);
  const [sort, setSort] = useState<`${SortKey}:${SortDir}`>('created_at:desc');
  const [showFilters, setShowFilters] = useState(false);
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [orderNumber, setOrderNumber] = useState('');
  const [customerName, setCustomerName] = useState('');

  const [bulkProgress, setBulkProgress] = useState<InvoiceSSEEventData | null>(null);
  const [bulkCompleted, setBulkCompleted] = useState(false);
  const submittingTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [orderStatuses, setOrderStatuses] = useState<Record<string, OrderProcessingStatus>>({});

  const [isSuperAdmin, setIsSuperAdmin] = useState(false);
  const [businesses, setBusinesses] = useState<Business[]>([]);
  const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
  const [loadingBusinesses, setLoadingBusinesses] = useState(false);
  const [showBusinessAlert, setShowBusinessAlert] = useState(false);

  const superAdminNeedsBusiness = isSuperAdmin && !selectedBusinessId;
  const currentBusinessId = propBusinessId ?? selectedBusinessId ?? 0;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  const handleBulkJobProgress = useCallback((data: InvoiceSSEEventData) => setBulkProgress(data), []);
  const handleBulkJobCompleted = useCallback((data: InvoiceSSEEventData) => {
    if (submittingTimeoutRef.current) {
      clearTimeout(submittingTimeoutRef.current);
      submittingTimeoutRef.current = null;
    }
    setBulkProgress(data);
    setBulkCompleted(true);
    setSubmitting(false);
  }, []);
  const handleInvoiceCreated = useCallback((data: InvoiceSSEEventData) => {
    if (data.order_id && selectedOrderIds.has(data.order_id)) {
      setOrderStatuses(prev => ({ ...prev, [data.order_id!]: { status: 'success', invoice_id: data.invoice_id, invoice_number: data.invoice_number } }));
    }
  }, [selectedOrderIds]);
  const handleInvoiceFailed = useCallback((data: InvoiceSSEEventData) => {
    if (data.order_id && selectedOrderIds.has(data.order_id)) {
      setOrderStatuses(prev => ({ ...prev, [data.order_id!]: { status: 'failed', error_message: data.error_message } }));
    }
  }, [selectedOrderIds]);
  const handleInvoicePendingValidation = useCallback((data: InvoiceSSEEventData) => {
    if (data.order_id && selectedOrderIds.has(data.order_id)) {
      setOrderStatuses(prev => ({ ...prev, [data.order_id!]: { status: 'pending_validation', invoice_id: data.invoice_id, invoice_number: data.invoice_number } }));
    }
  }, [selectedOrderIds]);

  useInvoiceSSE({
    businessId: currentBusinessId,
    onBulkJobProgress: handleBulkJobProgress,
    onBulkJobCompleted: handleBulkJobCompleted,
    onInvoiceCreated: handleInvoiceCreated,
    onInvoiceFailed: handleInvoiceFailed,
    onInvoicePendingValidation: handleInvoicePendingValidation,
  });

  const loadOrders = useCallback(async () => {
    if (superAdminNeedsBusiness && !propBusinessId) return;
    setLoading(true);
    setError(null);
    try {
      const [sortBy, sortOrder] = sort.split(':') as [SortKey, SortDir];
      const filters: InvoiceableOrdersFilters = {
        page,
        pageSize,
        businessId: propBusinessId ?? selectedBusinessId ?? undefined,
        startDate: startDate || undefined,
        endDate: endDate || undefined,
        orderNumber: orderNumber || undefined,
        customerName: customerName || search || undefined,
        sortBy,
        sortOrder,
      };
      const result = await getInvoiceableOrdersAction(filters);
      setOrders(result.data);
      setTotal(result.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar las ordenes facturables');
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, sort, propBusinessId, selectedBusinessId, startDate, endDate, orderNumber, customerName, search, superAdminNeedsBusiness]);

  useEffect(() => {
    if (!isOpen) {
      setOrders([]); setSelectedOrderIds(new Set()); setSearch(''); setStartDate(''); setEndDate(''); setOrderNumber(''); setCustomerName('');
      setPage(1); setError(null); setSelectedBusinessId(null); setBulkProgress(null); setBulkCompleted(false); setSubmitting(false); setOrderStatuses({}); setShowBusinessAlert(false);
      if (submittingTimeoutRef.current) { clearTimeout(submittingTimeoutRef.current); submittingTimeoutRef.current = null; }
      return;
    }
    (async () => {
      const user = CookieStorage.getUser();
      if (user?.is_super_admin) {
        setIsSuperAdmin(true);
        setLoadingBusinesses(true);
        try {
          const r = await getBusinessesAction();
          setBusinesses(r.data.map((b: any) => ({ id: b.id, name: b.name })));
        } finally {
          setLoadingBusinesses(false);
        }
      }
    })();
  }, [isOpen]);

  useEffect(() => {
    if (!isOpen) return;
    if (isSuperAdmin && !selectedBusinessId && !propBusinessId) return;
    const t = setTimeout(loadOrders, 300);
    return () => clearTimeout(t);
  }, [isOpen, loadOrders, isSuperAdmin, selectedBusinessId, propBusinessId]);

  const formatCurrency = (amount: number, currency: string) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: currency || 'COP' }).format(amount);

  const translateOrderStatus = (status: string) => {
    const map: Record<string, { label: string; color: string }> = {
      pending: { label: 'Pendiente', color: 'bg-yellow-100 text-yellow-800' },
      confirmed: { label: 'Confirmado', color: 'bg-blue-100 text-blue-800' },
      processing: { label: 'En proceso', color: 'bg-blue-100 text-blue-800' },
      shipped: { label: 'Enviado', color: 'bg-indigo-100 text-indigo-800' },
      delivered: { label: 'Entregado', color: 'bg-green-100 text-green-800' },
      fulfilled: { label: 'Completado', color: 'bg-green-100 text-green-800' },
      partially_fulfilled: { label: 'Parcialmente completado', color: 'bg-orange-100 text-orange-800' },
      cancelled: { label: 'Cancelado', color: 'bg-red-100 text-red-800' },
      refunded: { label: 'Reembolsado', color: 'bg-red-100 text-red-800' },
      new: { label: 'Nuevo', color: 'bg-gray-100 text-gray-700 dark:text-gray-200' },
    };
    return map[status?.toLowerCase()] ?? { label: status || '—', color: 'bg-gray-100 text-gray-600 dark:text-gray-300' };
  };

  const handleToggleOrder = (id: string) => {
    if (superAdminNeedsBusiness) { setShowBusinessAlert(true); return; }
    setSelectedOrderIds(prev => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id); else next.add(id);
      return next;
    });
  };

  const handleSelectPage = () => {
    if (superAdminNeedsBusiness) { setShowBusinessAlert(true); return; }
    setSelectedOrderIds(prev => {
      const next = new Set(prev);
      const pageIds = orders.map(o => o.id);
      const allInPage = pageIds.every(id => next.has(id));
      if (allInPage) pageIds.forEach(id => next.delete(id));
      else pageIds.forEach(id => next.add(id));
      return next;
    });
  };

  const handleSelectAllMatching = async () => {
    if (superAdminNeedsBusiness) { setShowBusinessAlert(true); return; }
    setLoading(true);
    try {
      const [sortBy, sortOrder] = sort.split(':') as [SortKey, SortDir];
      const r = await getInvoiceableOrdersAction({
        page: 1,
        pageSize: 200,
        businessId: propBusinessId ?? selectedBusinessId ?? undefined,
        startDate: startDate || undefined,
        endDate: endDate || undefined,
        orderNumber: orderNumber || undefined,
        customerName: customerName || search || undefined,
        sortBy, sortOrder,
      });
      setSelectedOrderIds(new Set(r.data.map(o => o.id)));
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async () => {
    if (superAdminNeedsBusiness) { setShowBusinessAlert(true); return; }
    if (selectedOrderIds.size === 0) return;
    setSubmitting(true);
    setBulkProgress(null);
    setBulkCompleted(false);
    const initial: Record<string, OrderProcessingStatus> = {};
    selectedOrderIds.forEach(id => { initial[id] = { status: 'pending' }; });
    setOrderStatuses(initial);
    const result = await createBulkInvoicesAction({
      order_ids: Array.from(selectedOrderIds),
      ...(isSuperAdmin && selectedBusinessId ? { business_id: selectedBusinessId } : {}),
    });
    if (!result.success) {
      alert(result.error);
      setSubmitting(false);
      setBulkProgress(null);
      setBulkCompleted(false);
      return;
    }
    submittingTimeoutRef.current = setTimeout(() => {
      setSubmitting(false);
      setBulkCompleted(true);
      alert('El proceso de facturacion se ha iniciado. Revisa la lista de facturas para ver el resultado.');
    }, 30000);
  };

  const handleCloseAfterCompletion = () => { onSuccess(); onClose(); };

  const resetFilters = () => {
    setSearch(''); setStartDate(''); setEndDate(''); setOrderNumber(''); setCustomerName('');
    setPage(1);
  };

  if (!isOpen) return null;

  const pageSelectionState = (() => {
    if (orders.length === 0) return 'empty';
    const allSelected = orders.every(o => selectedOrderIds.has(o.id));
    if (allSelected) return 'all';
    const someSelected = orders.some(o => selectedOrderIds.has(o.id));
    return someSelected ? 'some' : 'none';
  })();

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-center justify-center p-4">
        <div className="fixed inset-0 backdrop-blur-sm bg-white dark:bg-gray-800/10 transition-opacity" onClick={onClose} />

        <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-[95vw] h-[95vh] flex flex-col overflow-hidden">
          <div className="flex items-center justify-between p-6 bg-gradient-to-r from-[#7c3aed] to-[#6d28d9]">
            <div>
              <h2 className="text-2xl font-bold text-white">Crear Facturas</h2>
              <p className="text-purple-100 text-sm mt-1">Selecciona las ordenes para generar facturas electronicas</p>
            </div>
            <div className="flex items-center gap-3">
              <div className="bg-white dark:bg-gray-800/20 backdrop-blur-sm rounded-full px-4 py-2">
                <span className="text-white font-semibold">
                  {selectedOrderIds.size} seleccionada{selectedOrderIds.size !== 1 ? 's' : ''}
                </span>
              </div>
              <button onClick={onClose} className="text-white hover:bg-white dark:bg-gray-800/20 rounded-full p-2 transition-all">
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>
          </div>

          <div className="flex-1 flex flex-col overflow-hidden p-6 bg-white dark:bg-gray-800">
            {isSuperAdmin && (
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                  Seleccionar Negocio <span className="text-red-500">*</span>
                </label>
                <select
                  value={selectedBusinessId ?? ''}
                  onChange={(e) => {
                    const v = e.target.value;
                    setShowBusinessAlert(false);
                    setSelectedBusinessId(v === '' ? null : parseInt(v));
                    setSelectedOrderIds(new Set());
                    setPage(1);
                  }}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${superAdminNeedsBusiness ? 'border-amber-400 bg-amber-50' : 'border-gray-300'}`}
                  disabled={loadingBusinesses}
                >
                  <option value="">-- Selecciona un negocio --</option>
                  {businesses.map(b => <option key={b.id} value={b.id}>{b.name}</option>)}
                </select>
              </div>
            )}

            {superAdminNeedsBusiness ? (
              <div className="py-8 text-center text-gray-400">
                <ExclamationTriangleIcon className="w-12 h-12 mx-auto mb-3 text-amber-300" />
                <p className="text-gray-500 dark:text-gray-400">Selecciona un negocio para ver las ordenes facturables</p>
              </div>
            ) : (
              <div className="flex-1 flex flex-col min-h-0">
                <div className="flex flex-wrap gap-2 items-center mb-3">
                  <div className="relative flex-1 min-w-[260px]">
                    <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <input
                      type="text"
                      placeholder="Buscar por cliente..."
                      value={search}
                      onChange={(e) => { setSearch(e.target.value); setPage(1); }}
                      className="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#7c3aed]"
                    />
                  </div>
                  <button
                    type="button"
                    onClick={() => setShowFilters(s => !s)}
                    className={`flex items-center gap-1.5 px-3 py-2 text-sm rounded-md border transition ${showFilters ? 'bg-violet-50 border-violet-300 text-violet-700' : 'bg-white border-gray-300 text-gray-700 hover:bg-gray-50'}`}
                  >
                    <FunnelIcon className="w-4 h-4" />
                    Filtros
                  </button>
                  <div className="flex items-center gap-1.5 px-3 py-2 text-sm rounded-md border border-gray-300 bg-white">
                    <ArrowsUpDownIcon className="w-4 h-4 text-gray-500" />
                    <select
                      value={sort}
                      onChange={(e) => { setSort(e.target.value as `${SortKey}:${SortDir}`); setPage(1); }}
                      className="bg-transparent outline-none cursor-pointer"
                    >
                      {SORT_OPTIONS.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
                    </select>
                  </div>
                  <select
                    value={pageSize}
                    onChange={(e) => { setPageSize(parseInt(e.target.value)); setPage(1); }}
                    className="px-3 py-2 text-sm border border-gray-300 rounded-md bg-white"
                  >
                    {PAGE_SIZE_OPTIONS.map(n => <option key={n} value={n}>{n} / pag</option>)}
                  </select>
                </div>

                {showFilters && (
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-3 p-3 mb-3 bg-gray-50 rounded-lg border border-gray-200">
                    <div>
                      <label className="block text-xs font-medium text-gray-600 mb-1">Desde</label>
                      <input type="date" value={startDate} onChange={(e) => { setStartDate(e.target.value); setPage(1); }} className="w-full px-2 py-1.5 text-sm border border-gray-300 rounded" />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-gray-600 mb-1">Hasta</label>
                      <input type="date" value={endDate} onChange={(e) => { setEndDate(e.target.value); setPage(1); }} className="w-full px-2 py-1.5 text-sm border border-gray-300 rounded" />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-gray-600 mb-1">N orden</label>
                      <input type="text" value={orderNumber} onChange={(e) => { setOrderNumber(e.target.value); setPage(1); }} placeholder="MYS-0001" className="w-full px-2 py-1.5 text-sm border border-gray-300 rounded" />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-gray-600 mb-1">Cliente</label>
                      <input type="text" value={customerName} onChange={(e) => { setCustomerName(e.target.value); setPage(1); }} placeholder="Juan Perez" className="w-full px-2 py-1.5 text-sm border border-gray-300 rounded" />
                    </div>
                    <div className="md:col-span-4 flex justify-end">
                      <button type="button" onClick={resetFilters} className="px-3 py-1.5 text-xs text-gray-600 hover:text-gray-900 underline">
                        Limpiar filtros
                      </button>
                    </div>
                  </div>
                )}

                <div className="flex items-center gap-3 p-2.5 bg-gray-50 rounded mb-2 flex-wrap">
                  <input
                    type="checkbox"
                    checked={pageSelectionState === 'all'}
                    ref={el => { if (el) el.indeterminate = pageSelectionState === 'some'; }}
                    onChange={handleSelectPage}
                    className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                  <span className="text-sm font-medium">Pagina ({orders.length})</span>
                  <button type="button" onClick={handleSelectAllMatching} disabled={total === 0} className="text-xs text-violet-600 hover:underline disabled:opacity-50">
                    Seleccionar todas las {total} coincidencias
                  </button>
                  <button type="button" onClick={() => setSelectedOrderIds(new Set())} disabled={selectedOrderIds.size === 0} className="text-xs text-gray-500 hover:underline disabled:opacity-30 ml-auto">
                    Limpiar seleccion
                  </button>
                </div>

                <div className="border border-gray-200 dark:border-gray-700 rounded-lg flex-1 min-h-0 overflow-y-auto">
                  {loading ? (
                    <div className="flex items-center justify-center py-12">
                      <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600" />
                    </div>
                  ) : error ? (
                    <div className="py-8 text-center">
                      <p className="text-red-600 mb-4 font-medium">{error}</p>
                      <button onClick={loadOrders} className="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded">Reintentar</button>
                    </div>
                  ) : orders.length === 0 ? (
                    <div className="py-12 text-center text-gray-500 dark:text-gray-400">
                      <p className="mb-2">No hay ordenes facturables que coincidan con los filtros</p>
                    </div>
                  ) : orders.map(order => {
                    const orderStatus = orderStatuses[order.id];
                    const isSelected = selectedOrderIds.has(order.id);
                    const getBgColor = () => {
                      if (!submitting && !bulkCompleted) return isSelected ? 'bg-blue-50 border-blue-200' : '';
                      if (!orderStatus || orderStatus.status === 'pending') return 'bg-gray-50';
                      if (orderStatus.status === 'processing') return 'bg-yellow-50 border-yellow-200';
                      if (orderStatus.status === 'success') return 'bg-green-50 border-green-200';
                      if (orderStatus.status === 'failed') return 'bg-red-50 border-red-200';
                      if (orderStatus.status === 'pending_validation') return 'bg-amber-50 border-amber-200';
                      return '';
                    };
                    const StatusIcon = () => {
                      if (!orderStatus) return null;
                      if (orderStatus.status === 'processing') return <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-yellow-600" />;
                      if (orderStatus.status === 'success') return <CheckCircleIcon className="w-5 h-5 text-green-600" />;
                      if (orderStatus.status === 'failed') return <XCircleIcon className="w-5 h-5 text-red-600" />;
                      if (orderStatus.status === 'pending_validation') return <ClockIcon className="w-5 h-5 text-amber-600" />;
                      return null;
                    };
                    return (
                      <div
                        key={order.id}
                        onClick={() => !submitting && !bulkCompleted && handleToggleOrder(order.id)}
                        className={`p-3 border-b last:border-b-0 ${submitting || bulkCompleted ? '' : 'cursor-pointer hover:bg-gray-50'} transition-colors ${getBgColor()}`}
                      >
                        <div className="flex items-start gap-3">
                          <input type="checkbox" checked={isSelected} disabled={submitting || bulkCompleted} onChange={() => {}} className="mt-1 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50" />
                          <div className="flex-1">
                            <div className="flex items-center justify-between mb-1">
                              <div className="flex items-center gap-2">
                                <span className="font-medium text-gray-900 dark:text-white">{order.order_number}</span>
                                <StatusIcon />
                              </div>
                              <span className="font-semibold text-gray-900 dark:text-white">{formatCurrency(order.total_amount, order.currency)}</span>
                            </div>
                            <p className="text-sm text-gray-600 dark:text-gray-300">{order.customer_name}</p>
                            <div className="flex items-center gap-2 mt-1 flex-wrap">
                              <p className="text-xs text-gray-500 dark:text-gray-400">{new Date(order.created_at).toLocaleDateString('es-CO', { year: 'numeric', month: 'short', day: 'numeric' })}</p>
                              {order.status && (() => {
                                const s = translateOrderStatus(order.status);
                                return (<>
                                  <span className="text-xs text-gray-400">•</span>
                                  <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${s.color}`}>{s.label}</span>
                                </>);
                              })()}
                              <span className="text-xs text-gray-400">•</span>
                              <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${order.is_paid ? 'bg-green-100 text-green-800' : 'bg-orange-100 text-orange-800'}`}>
                                {order.is_paid ? 'Pagado' : 'Sin pagar'}
                              </span>
                              {isSuperAdmin && (<>
                                <span className="text-xs text-gray-400">•</span>
                                <span className="text-xs text-blue-600 font-medium">Business #{order.business_id}</span>
                              </>)}
                            </div>
                            {orderStatus?.status === 'success' && (
                              <div className="mt-2 p-2 bg-green-100 border border-green-300 rounded text-xs text-green-800">
                                <p className="font-semibold">Factura creada</p>
                                {orderStatus.invoice_number && <p className="font-mono">{orderStatus.invoice_number}</p>}
                              </div>
                            )}
                            {orderStatus?.status === 'pending_validation' && (
                              <div className="mt-2 p-2 bg-amber-100 border border-amber-300 rounded text-xs text-amber-800">
                                <p className="font-semibold">Pendiente validacion DIAN</p>
                                {orderStatus.invoice_number && <p className="font-mono">{orderStatus.invoice_number}</p>}
                              </div>
                            )}
                            {orderStatus?.status === 'failed' && (
                              <div className="mt-2 p-2 bg-red-100 border border-red-300 rounded text-xs text-red-800">
                                <p className="font-semibold">Error</p>
                                {orderStatus.error_message && <p className="break-words whitespace-pre-wrap">{orderStatus.error_message}</p>}
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>

                <div className="flex items-center justify-between mt-3 text-sm">
                  <span className="text-gray-600 dark:text-gray-300">
                    {total === 0 ? '0 resultados' : `Pagina ${page} de ${totalPages} - ${total} total`}
                  </span>
                  <div className="flex items-center gap-2">
                    <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page <= 1 || loading} className="p-1.5 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-40">
                      <ChevronLeftIcon className="w-4 h-4" />
                    </button>
                    <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page >= totalPages || loading} className="p-1.5 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-40">
                      <ChevronRightIcon className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            )}
          </div>

          {(submitting || bulkCompleted) && bulkProgress && (
            <div className="px-6 py-3 border-t bg-gray-50">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium">{bulkCompleted ? 'Procesamiento completado' : 'Procesando facturas...'}</span>
                <span className="text-sm text-gray-500">{bulkProgress.processed ?? 0}/{bulkProgress.total_orders ?? 0}</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden flex">
                {(bulkCompleted || (bulkProgress.successful ?? 0) > 0 || (bulkProgress.failed ?? 0) > 0) ? (() => {
                  const total = bulkProgress.total_orders || 1;
                  const pendingValidationCount = Object.values(orderStatuses).filter(s => s.status === 'pending_validation').length;
                  const pureSuccess = Math.max((bulkProgress.successful ?? 0) - pendingValidationCount, 0);
                  return (<>
                    {pureSuccess > 0 && <div className="h-2 bg-green-500" style={{ width: `${(pureSuccess / total) * 100}%` }} />}
                    {pendingValidationCount > 0 && <div className="h-2 bg-amber-400" style={{ width: `${(pendingValidationCount / total) * 100}%` }} />}
                    {(bulkProgress.failed ?? 0) > 0 && <div className="h-2 bg-red-500" style={{ width: `${((bulkProgress.failed ?? 0) / total) * 100}%` }} />}
                  </>);
                })() : (
                  <div className="h-2 bg-blue-600 rounded-full" style={{ width: `${Math.min(bulkProgress.progress ?? 0, 100)}%` }} />
                )}
              </div>
            </div>
          )}

          <div className="flex justify-end gap-3 p-4 border-t bg-gray-50">
            {bulkCompleted ? (
              <button onClick={handleCloseAfterCompletion} className="px-5 py-2 bg-green-600 text-white font-semibold rounded-full hover:shadow-lg transition-all">Cerrar</button>
            ) : (
              <>
                <button onClick={onClose} disabled={submitting} className="px-5 py-2 border border-gray-300 text-gray-700 dark:text-gray-200 font-semibold rounded-full hover:bg-gray-100 disabled:opacity-50">Cancelar</button>
                <button onClick={handleSubmit} disabled={submitting || selectedOrderIds.size === 0 || loading || superAdminNeedsBusiness} className="px-5 py-2 bg-gradient-to-r from-[#7c3aed] to-[#6d28d9] text-white font-semibold rounded-full hover:shadow-lg disabled:opacity-50">
                  {submitting ? 'Creando...' : `Crear ${selectedOrderIds.size} Factura(s)`}
                </button>
              </>
            )}
          </div>
        </div>
      </div>

      {showBusinessAlert && (
        <div className="fixed inset-0 z-[60] flex items-center justify-center">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-sm" onClick={() => setShowBusinessAlert(false)} />
          <div className="relative bg-white dark:bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full mx-4 overflow-hidden">
            <div className="bg-amber-500 p-4 flex items-center gap-3">
              <ExclamationTriangleIcon className="w-8 h-8 text-white" />
              <h3 className="text-lg font-bold text-white">Negocio requerido</h3>
            </div>
            <div className="p-6">
              <p className="text-gray-700 dark:text-gray-200">Como super administrador, debes seleccionar un <strong>negocio</strong> antes de continuar.</p>
            </div>
            <div className="flex justify-end p-4 border-t bg-gray-50">
              <button onClick={() => setShowBusinessAlert(false)} className="px-5 py-2 bg-amber-500 text-white font-semibold rounded-full hover:bg-amber-600">Entendido</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
