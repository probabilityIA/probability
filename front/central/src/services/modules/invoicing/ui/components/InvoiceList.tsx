/**
 * Componente para listar facturas con paginación y detalle en modal
 */

'use client';

import { useState, useEffect, useCallback, useMemo, forwardRef, useImperativeHandle } from 'react';
import { EyeIcon } from '@heroicons/react/24/outline';
import { Table } from '@/shared/ui/table';
import { Badge } from '@/shared/ui/badge';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { useToast } from '@/shared/providers/toast-provider';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { BulkCreateInvoiceModal } from './BulkCreateInvoiceModal';
import { InvoiceDetailModal } from './InvoiceDetailPanel';
import { InvoiceComparisonModal } from './InvoiceComparisonModal';
import {
  getInvoicesAction,
  cancelInvoiceAction,
  requestInvoiceComparisonAction,
} from '../../infra/actions';
import { useInvoiceSSE } from '../hooks/useInvoiceSSE';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import type { Invoice, InvoiceFilters, CompareResponseData } from '../../domain/types';

interface InvoiceListProps {
  businessId: number;
  onOpenBulkModal?: () => void;
}

const PAGE_SIZE_DEFAULT = 20;

export const InvoiceList = forwardRef(function InvoiceList(
  { businessId }: InvoiceListProps,
  ref
) {
  const { showToast } = useToast();
  const { isSuperAdmin } = usePermissions();
  const { businesses } = useBusinessesSimple();
  const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
  const [filters, setFilters] = useState<InvoiceFilters>({});
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(PAGE_SIZE_DEFAULT);
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [showBulkModal, setShowBulkModal] = useState(false);
  const [showDetailModal, setShowDetailModal] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);
  const [showCompareModal, setShowCompareModal] = useState(false);
  const [compareLoading, setCompareLoading] = useState(false);
  const [compareCorrelationId, setCompareCorrelationId] = useState<string | null>(null);
  const [compareData, setCompareData] = useState<CompareResponseData | null>(null);

  useImperativeHandle(ref, () => ({
    openBulkModal: () => setShowBulkModal(true),
  }));

  useEffect(() => {
    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/4dbf3696-4a46-47a8-86ba-70e3d0546d6b', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        sessionId: 'debug-session',
        runId: 'invoice-table-ui-v1',
        hypothesisId: 'H1',
        location: 'InvoiceList.tsx:mount',
        message: 'InvoiceList mounted',
        data: { businessIdIsZero: businessId === 0 },
        timestamp: Date.now(),
      }),
    }).catch(() => { });
    // #endregion
  }, [businessId]);

  // SSE: Escuchar eventos en tiempo real
  useInvoiceSSE({
    businessId,
    onInvoiceCreated: (data) => {
      showToast(
        `Factura ${data.invoice_number || ''} creada exitosamente para ${data.customer_name || 'orden'}`,
        'success'
      );
      loadInvoices(currentPage, pageSize);
    },
    onInvoiceFailed: (data) => {
      showToast(
        `Error al crear factura: ${data.error_message || 'Error desconocido'}`,
        'error'
      );
      loadInvoices(currentPage, pageSize);
    },
    onInvoiceCancelled: () => {
      loadInvoices(currentPage, pageSize);
    },
    onCompareReady: (data) => {
      if (!compareCorrelationId || data.correlation_id === compareCorrelationId) {
        setCompareData(data);
        setCompareLoading(false);
      }
    },
  });

  const loadInvoices = useCallback(async (page: number, size: number) => {
    try {
      setLoading(true);
      const effectiveBusinessId = isSuperAdmin
        ? (selectedBusinessId ?? undefined)
        : businessId;
      const finalFilters: InvoiceFilters = effectiveBusinessId
        ? { ...filters, business_id: effectiveBusinessId, page, page_size: size }
        : { ...filters, page, page_size: size };

      const response = await getInvoicesAction(finalFilters);

      // AGRUPAR: Mostrar solo UNA factura por orden
      // Prioridad: 1) Facturas con status != failed, 2) La más reciente
      const grouped = (response.data || []).reduce((acc, invoice) => {
        const existing = acc.get(invoice.order_id);

        if (!existing) {
          // Si no existe, agregar
          acc.set(invoice.order_id, invoice);
        } else {
          // Si existe, reemplazar solo si:
          // - La actual NO es failed y la existente SÍ es failed, O
          // - Ambas son failed/non-failed y la actual es más reciente
          const currentIsFailed = invoice.status === 'failed';
          const existingIsFailed = existing.status === 'failed';

          if ((!currentIsFailed && existingIsFailed) ||
            (currentIsFailed === existingIsFailed &&
              new Date(invoice.created_at) > new Date(existing.created_at))) {
            acc.set(invoice.order_id, invoice);
          }
        }

        return acc;
      }, new Map<string, Invoice>());

      const uniqueInvoices = Array.from(grouped.values());

      setInvoices(uniqueInvoices);
      setTotalCount(uniqueInvoices.length); // Actualizar total con facturas únicas
    } catch (error: any) {
      showToast('Error al cargar facturas: ' + error.message, 'error');
      setInvoices([]);
      setTotalCount(0);
    } finally {
      setLoading(false);
    }
  }, [businessId, filters, showToast, isSuperAdmin, selectedBusinessId]);

  useEffect(() => {
    setCurrentPage(1);
    loadInvoices(1, pageSize);
  }, [businessId, JSON.stringify(filters), selectedBusinessId]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    loadInvoices(page, pageSize);
  };

  const handlePageSizeChange = (size: number) => {
    setPageSize(size);
    setCurrentPage(1);
    loadInvoices(1, size);
  };

  const handleCancelInvoice = async () => {
    if (!selectedInvoice) return;
    try {
      setActionLoading(true);
      await cancelInvoiceAction(selectedInvoice.id);
      showToast('Factura cancelada exitosamente', 'success');
      setShowCancelModal(false);
      loadInvoices(currentPage, pageSize);
    } catch (error: any) {
      showToast('Error al cancelar factura: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

  const handleRowDoubleClick = (invoice: Invoice) => {
    setSelectedInvoice(invoice);
    setShowDetailModal(true);
  };

  const handleRequestComparison = useCallback(async (dateFrom: string, dateTo: string) => {
    if (isSuperAdmin && !selectedBusinessId) {
      showToast('Selecciona un negocio antes de comparar', 'error');
      return;
    }
    setCompareLoading(true);
    setCompareData(null);
    setCompareCorrelationId(null);
    try {
      const result = await requestInvoiceComparisonAction(
        dateFrom,
        dateTo,
        isSuperAdmin ? selectedBusinessId ?? undefined : undefined
      );
      setCompareCorrelationId(result.correlation_id);
    } catch (error: any) {
      showToast('Error al iniciar comparación: ' + error.message, 'error');
      setCompareLoading(false);
    }
  }, [showToast, isSuperAdmin, selectedBusinessId]);

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'secondary' }> = {
      issued: { label: 'Emitida', type: 'success' },
      pending: { label: 'Pendiente', type: 'warning' },
      cancelled: { label: 'Cancelada', type: 'error' },
      failed: { label: 'Fallida', type: 'error' },
    };
    const config = statusConfig[status] || { label: status, type: 'secondary' as const };
    return <Badge type={config.type}>{config.label}</Badge>;
  };

  // ===== DynamicFilters =====
  const availableFilters: FilterOption[] = useMemo(() => [
    {
      key: 'invoice_number',
      label: 'Nº Factura',
      type: 'text',
      placeholder: 'Buscar por número de factura...',
    },
    {
      key: 'order_number',
      label: 'Nº Orden',
      type: 'text',
      placeholder: 'Buscar por #orden o ID Shopify...',
    },
    {
      key: 'customer_name',
      label: 'Cliente',
      type: 'text',
      placeholder: 'Buscar por nombre de cliente...',
    },
    {
      key: 'status',
      label: 'Estado',
      type: 'select',
      options: [
        { value: 'issued', label: 'Emitida' },
        { value: 'pending', label: 'Pendiente' },
        { value: 'cancelled', label: 'Cancelada' },
        { value: 'failed', label: 'Fallida' },
      ],
    },
    {
      key: 'created_at',
      label: 'Rango de fechas',
      type: 'date-range',
    },
  ], []);

  const activeFilters: ActiveFilter[] = useMemo(() => {
    const active: ActiveFilter[] = [];

    if (filters.invoice_number) {
      active.push({ key: 'invoice_number', label: 'Nº Factura', value: filters.invoice_number, type: 'text' });
    }
    if (filters.order_number) {
      active.push({ key: 'order_number', label: 'Nº Orden', value: filters.order_number, type: 'text' });
    }
    if (filters.customer_name) {
      active.push({ key: 'customer_name', label: 'Cliente', value: filters.customer_name, type: 'text' });
    }
    if (filters.status) {
      const statusLabels: Record<string, string> = { issued: 'Emitida', pending: 'Pendiente', cancelled: 'Cancelada', failed: 'Fallida' };
      active.push({ key: 'status', label: 'Estado', value: statusLabels[filters.status] || filters.status, type: 'select' });
    }
    if (filters.start_date || filters.end_date) {
      active.push({
        key: 'created_at',
        label: 'Rango de fechas',
        value: { start: filters.start_date, end: filters.end_date },
        type: 'date-range',
      });
    }

    return active;
  }, [filters]);

  const handleAddFilter = useCallback((filterKey: string, value: any) => {
    setFilters((prev) => {
      const newFilters = { ...prev };
      if (filterKey === 'created_at' && typeof value === 'object') {
        newFilters.start_date = value.start;
        newFilters.end_date = value.end;
      } else {
        (newFilters as any)[filterKey] = value;
      }
      return newFilters;
    });
    setCurrentPage(1);
  }, []);

  const handleRemoveFilter = useCallback((filterKey: string) => {
    setFilters((prev) => {
      const newFilters = { ...prev };
      if (filterKey === 'created_at') {
        delete newFilters.start_date;
        delete newFilters.end_date;
      } else {
        delete (newFilters as any)[filterKey];
      }
      return newFilters;
    });
    setCurrentPage(1);
  }, []);

  const totalPages = Math.max(1, Math.ceil(totalCount / pageSize));

  const columns = [
    {
      key: 'invoice_number',
      label: 'Factura',
      render: (_: unknown, invoice: Invoice) => (
        <div>
          <div className="font-medium">{invoice.invoice_number || 'Sin número'}</div>
          <div className="text-xs text-gray-500">ID: {invoice.id}</div>
        </div>
      ),
    },
    {
      key: 'order_id',
      label: 'Orden',
      render: (_: unknown, invoice: Invoice) => (
        <div className="text-sm text-gray-600 font-mono">
          {invoice.order_id.substring(0, 8)}...
        </div>
      ),
    },
    {
      key: 'customer_name',
      label: 'Cliente',
      render: (_: unknown, invoice: Invoice) => (
        <div>
          <div className="font-medium">{invoice.customer_name || '-'}</div>
          {invoice.customer_email && (
            <div className="text-xs text-gray-500">{invoice.customer_email}</div>
          )}
        </div>
      ),
    },
    {
      key: 'total_amount',
      label: 'Total',
      render: (_: unknown, invoice: Invoice) => (
        <div className="font-semibold">
          {new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: invoice.currency || 'COP',
          }).format(invoice.total_amount)}
        </div>
      ),
    },
    {
      key: 'provider',
      label: 'Proveedor',
      render: (_: unknown, invoice: Invoice) => (
        invoice.provider_logo_url ? (
          <img
            src={invoice.provider_logo_url}
            alt={invoice.provider_name || 'Proveedor'}
            title={invoice.provider_name || 'Proveedor'}
            className="h-6 w-auto object-contain"
          />
        ) : (
          <span className="text-xs text-gray-400">{invoice.provider_name || '—'}</span>
        )
      ),
    },
    {
      key: 'status',
      label: 'Estado',
      render: (_: unknown, invoice: Invoice) => getStatusBadge(invoice.status),
    },
    {
      key: 'created_at',
      label: 'Fecha',
      render: (_: unknown, invoice: Invoice) => (
        <div className="text-sm text-gray-600">
          {new Date(invoice.created_at).toLocaleString('es-CO', {
            day: '2-digit',
            month: 'short',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
          })}
        </div>
      ),
    },
    {
      key: 'actions',
      label: '',
      width: '50px',
      align: 'right' as const,
      render: (_: unknown, invoice: Invoice) => (
        <button
          className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded transition-colors"
          title="Ver detalle"
          onClick={(e) => {
            e.stopPropagation();
            handleRowDoubleClick(invoice);
          }}
        >
          <EyeIcon className="w-5 h-5" />
        </button>
      ),
    },
  ];

  return (
    <>
      {/* Business Selector — solo Super Admin */}
      {isSuperAdmin && businesses.length > 0 && (
        <div className="mb-4 bg-blue-50 border border-blue-200 rounded-lg p-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Seleccionar Negocio (Super Admin)
          </label>
          <select
            value={selectedBusinessId?.toString() ?? ''}
            onChange={(e) => {
              const val = e.target.value;
              setSelectedBusinessId(val ? Number(val) : null);
              setCurrentPage(1);
            }}
            className="w-full max-w-xs px-3 py-2 border border-gray-300 rounded-md text-sm"
          >
            <option value="">Todos los negocios</option>
            {businesses.map((b) => (
              <option key={b.id} value={b.id}>{b.name} (ID: {b.id})</option>
            ))}
          </select>
        </div>
      )}

      {/* Filtros + botón crear */}
      <div className="invoiceFilters flex items-center mb-4">
        <div className="flex-1 min-w-0">
          <DynamicFilters
            availableFilters={availableFilters}
            activeFilters={activeFilters}
            onAddFilter={handleAddFilter}
            onRemoveFilter={handleRemoveFilter}
            className="!p-0 !border-0 !shadow-none !rounded-none"
          />
        </div>
        <button
          onClick={() => {
            setShowCompareModal(true);
            setCompareData(null);
            setCompareLoading(false);
            setCompareCorrelationId(null);
          }}
          className="ml-3 flex-shrink-0 p-2 rounded-full bg-white border-2 border-[#7c3aed] text-[#7c3aed] hover:shadow-lg transition-all duration-200 hover:scale-110"
          title="Comparar con proveedor"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
          </svg>
        </button>
        <button
          onClick={() => setShowBulkModal(true)}
          className="ml-1 flex-shrink-0 p-2 rounded-full bg-white border-2 border-[#7c3aed] text-[#7c3aed] hover:shadow-lg transition-all duration-200 hover:scale-110"
          title="Crear Facturas"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>

      {/* Tabla con paginación */}
      <div className="invoiceTable">
        <Table
          data={invoices}
          columns={columns}
          loading={loading}
          emptyMessage="No hay facturas para mostrar"
          keyExtractor={(invoice: Invoice) => invoice.id}
          onRowDoubleClick={handleRowDoubleClick}
          pagination={{
            currentPage,
            totalPages,
            totalItems: totalCount,
            itemsPerPage: pageSize,
            onPageChange: handlePageChange,
            onItemsPerPageChange: handlePageSizeChange,
            showItemsPerPageSelector: true,
            itemsPerPageOptions: [10, 20, 50],
          }}
        />

        <style jsx>{`
          /* Botones morados en filtros */
          :global(.invoiceFilters) :global(button.btn-primary),
          :global(.invoiceFilters) :global(.btn-primary) {
            background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%) !important;
            border-color: #7c3aed !important;
          }
          :global(.invoiceFilters) :global(button.btn-primary:hover),
          :global(.invoiceFilters) :global(.btn-primary:hover) {
            filter: brightness(1.08);
          }
          /* Chips morados */
          :global(.invoiceFilters) :global(.bg-blue-50) {
            background-color: rgba(124, 58, 237, 0.08) !important;
            color: #6d28d9 !important;
          }
          :global(.invoiceFilters) :global(.text-blue-800) {
            color: #6d28d9 !important;
          }

          /* CTA morado (solo esta pantalla) */
          .invoiceHeader :global(.btn-primary) {
            background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%);
          }
          .invoiceHeader :global(.btn-primary:hover) {
            filter: brightness(1.05);
          }

          /* Tabla más “card-like” fila por fila (solo Facturas Electrónicas) */
          .invoiceTable :global(.table) {
            border-collapse: separate;
            border-spacing: 0 10px; /* separación entre filas */
            background: transparent;
          }

          /* Quitar el borde del contenedor global de Table SOLO aquí */
          .invoiceTable :global(div.overflow-hidden.w-full.rounded-lg.border.border-gray-200.bg-white) {
            border: none !important;
            background: transparent !important;
          }

          .invoiceTable :global(.table th) {
            background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%);
            color: #fff;
            position: sticky;
            top: 0;
            z-index: 1;
          }

          /* Header más llamativo + bordes redondeados */
          .invoiceTable :global(.table thead th) {
            padding-top: 10px;
            padding-bottom: 10px;
            font-size: 0.75rem; /* más pequeño */
            font-weight: 800;
            letter-spacing: 0.06em;
            text-transform: uppercase;
            box-shadow: 0 10px 25px rgba(124, 58, 237, 0.18);
          }

          .invoiceTable :global(.table thead th:first-child) {
            border-top-left-radius: 14px;
            border-bottom-left-radius: 14px;
          }

          .invoiceTable :global(.table thead th:last-child) {
            border-top-right-radius: 14px;
            border-bottom-right-radius: 14px;
          }

          .invoiceTable :global(.table tbody tr) {
            background: rgba(255, 255, 255, 0.95);
            box-shadow: 0 1px 0 rgba(17, 24, 39, 0.04);
            transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
          }

          /* Zebra suave en morado */
          .invoiceTable :global(.table tbody tr:nth-child(even)) {
            background: rgba(124, 58, 237, 0.03);
          }

          .invoiceTable :global(.table tbody tr:hover) {
            background: rgba(124, 58, 237, 0.06);
            box-shadow: 0 10px 25px rgba(17, 24, 39, 0.08);
            transform: translateY(-1px);
          }

          .invoiceTable :global(.table td) {
            border-top: none;
          }

          /* Redondeo de cada fila */
          .invoiceTable :global(.table tbody td:first-child) {
            border-top-left-radius: 12px;
            border-bottom-left-radius: 12px;
          }
          .invoiceTable :global(.table tbody td:last-child) {
            border-top-right-radius: 12px;
            border-bottom-right-radius: 12px;
          }

          /* Acciones: focus consistente */
          .invoiceTable :global(a),
          .invoiceTable :global(button) {
            outline-color: rgba(124, 58, 237, 0.35);
          }
        `}</style>
      </div>

      {/* Modal de detalle de factura */}
      <InvoiceDetailModal
        invoice={selectedInvoice}
        isOpen={showDetailModal}
        onClose={() => setShowDetailModal(false)}
        onCancel={(invoice) => {
          setShowDetailModal(false);
          setSelectedInvoice(invoice);
          setShowCancelModal(true);
        }}
        onRefresh={() => loadInvoices(currentPage, pageSize)}
        businessId={businessId}
      />

      {/* Modal de confirmación de cancelación */}
      <ConfirmModal
        isOpen={showCancelModal}
        onClose={() => setShowCancelModal(false)}
        onConfirm={handleCancelInvoice}
        title="Cancelar Factura"
        message={`¿Estás seguro de que deseas cancelar la factura ${selectedInvoice?.invoice_number}?`}
        confirmText="Sí, cancelar"
        cancelText="No, volver"
        type="danger"
      />

      {/* Modal de creación masiva */}
      <BulkCreateInvoiceModal
        isOpen={showBulkModal}
        onClose={() => setShowBulkModal(false)}
        onSuccess={() => loadInvoices(currentPage, pageSize)}
        businessId={businessId}
      />

      {/* Modal de auditoría comparativa */}
      <InvoiceComparisonModal
        isOpen={showCompareModal}
        onClose={() => setShowCompareModal(false)}
        loading={compareLoading}
        compareData={compareData}
        onRequestComparison={handleRequestComparison}
      />
    </>
  );
});
