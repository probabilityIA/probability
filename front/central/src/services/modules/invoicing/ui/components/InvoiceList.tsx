/**
<<<<<<< HEAD
 * Componente para listar facturas
=======
 * Componente para listar facturas con paginación y detalle en modal
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
 */

'use client';

<<<<<<< HEAD
import { useState, useEffect } from 'react';
import { Table } from '@/shared/ui/table';
import { Button } from '@/shared/ui/button';
import { Badge } from '@/shared/ui/badge';
import { Spinner } from '@/shared/ui/spinner';
import { useToast } from '@/shared/providers/toast-provider';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { getInvoicesAction, cancelInvoiceAction, retryInvoiceAction } from '../../infra/actions';
=======
import { useState, useEffect, useCallback } from 'react';
import { EyeIcon } from '@heroicons/react/24/outline';
import { Table } from '@/shared/ui/table';
import { Button } from '@/shared/ui/button';
import { Badge } from '@/shared/ui/badge';
import { useToast } from '@/shared/providers/toast-provider';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { BulkCreateInvoiceModal } from './BulkCreateInvoiceModal';
import { InvoiceDetailModal } from './InvoiceDetailPanel';
import {
  getInvoicesAction,
  cancelInvoiceAction,
} from '../../infra/actions';
import { useInvoiceSSE } from '../hooks/useInvoiceSSE';
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
import type { Invoice, InvoiceFilters } from '../../domain/types';

interface InvoiceListProps {
  businessId: number;
  filters?: InvoiceFilters;
}

<<<<<<< HEAD
=======
const PAGE_SIZE_DEFAULT = 20;

>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
export function InvoiceList({ businessId, filters = {} }: InvoiceListProps) {
  const { showToast } = useToast();
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);
<<<<<<< HEAD
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  const loadInvoices = async () => {
    try {
      setLoading(true);
      // Super admin (businessId = 0): no filtra por business_id, ve todas las facturas
      // Usuario normal: filtra por su business_id
      const isSuperAdmin = !businessId || businessId === 0;
      const finalFilters = isSuperAdmin
        ? { ...filters }
        : { ...filters, business_id: businessId };

      const response = await getInvoicesAction(finalFilters);
      setInvoices(response.data || []);
    } catch (error: any) {
      showToast('Error al cargar facturas: ' + error.message, 'error');
      setInvoices([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadInvoices();
  }, [businessId, JSON.stringify(filters)]);

  const handleCancelInvoice = async () => {
    if (!selectedInvoice) return;

=======
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(PAGE_SIZE_DEFAULT);
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [showBulkModal, setShowBulkModal] = useState(false);
  const [showDetailModal, setShowDetailModal] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

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
  });

  const loadInvoices = useCallback(async (page: number, size: number) => {
    try {
      setLoading(true);
      const isSuperAdmin = !businessId || businessId === 0;
      const finalFilters: InvoiceFilters = isSuperAdmin
        ? { ...filters, page, page_size: size }
        : { ...filters, business_id: businessId, page, page_size: size };

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
  }, [businessId, filters, showToast]);

  useEffect(() => {
    setCurrentPage(1);
    loadInvoices(1, pageSize);
  }, [businessId, JSON.stringify(filters)]);

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
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
    try {
      setActionLoading(true);
      await cancelInvoiceAction(selectedInvoice.id);
      showToast('Factura cancelada exitosamente', 'success');
      setShowCancelModal(false);
<<<<<<< HEAD
      loadInvoices();
=======
      loadInvoices(currentPage, pageSize);
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
    } catch (error: any) {
      showToast('Error al cancelar factura: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

<<<<<<< HEAD
  const handleRetryInvoice = async (invoice: Invoice) => {
    try {
      setActionLoading(true);
      await retryInvoiceAction(invoice.id);
      showToast('Factura reintentada exitosamente', 'success');
      loadInvoices();
    } catch (error: any) {
      showToast('Error al reintentar factura: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; color: 'green' | 'yellow' | 'red' | 'gray' }> = {
      issued: { label: 'Emitida', color: 'green' },
      pending: { label: 'Pendiente', color: 'yellow' },
      cancelled: { label: 'Cancelada', color: 'red' },
      failed: { label: 'Fallida', color: 'red' },
    };

    const config = statusConfig[status] || { label: status, color: 'gray' as const };
    return <Badge color={config.color}>{config.label}</Badge>;
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center py-12">
        <Spinner />
      </div>
    );
  }
=======
  const handleRowDoubleClick = (invoice: Invoice) => {
    setSelectedInvoice(invoice);
    setShowDetailModal(true);
  };

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

  const totalPages = Math.max(1, Math.ceil(totalCount / pageSize));
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

  const columns = [
    {
      key: 'invoice_number',
<<<<<<< HEAD
      label: 'Número',
      render: (_: unknown, invoice: Invoice) => (
        <div className="font-medium">{invoice.invoice_number || 'N/A'}</div>
=======
      label: 'Factura',
      render: (_: unknown, invoice: Invoice) => (
        <div>
          <div className="font-medium">{invoice.invoice_number || 'Sin número'}</div>
          <div className="text-xs text-gray-500">ID: {invoice.id}</div>
        </div>
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
      ),
    },
    {
      key: 'order_id',
      label: 'Orden',
      render: (_: unknown, invoice: Invoice) => (
<<<<<<< HEAD
        <div className="text-sm text-gray-600">{invoice.order_id}</div>
=======
        <div className="text-sm text-gray-600 font-mono">
          {invoice.order_id.substring(0, 8)}...
        </div>
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
      ),
    },
    {
      key: 'customer_name',
      label: 'Cliente',
      render: (_: unknown, invoice: Invoice) => (
        <div>
<<<<<<< HEAD
          <div className="font-medium">{invoice.customer_name}</div>
=======
          <div className="font-medium">{invoice.customer_name || '-'}</div>
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
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
      key: 'status',
      label: 'Estado',
      render: (_: unknown, invoice: Invoice) => getStatusBadge(invoice.status),
    },
    {
<<<<<<< HEAD
      key: 'issued_at',
      label: 'Fecha',
      render: (_: unknown, invoice: Invoice) => (
        <div className="text-sm text-gray-600">
          {invoice.issued_at
            ? new Date(invoice.issued_at).toLocaleDateString('es-CO')
            : 'Pendiente'}
=======
      key: 'created_at',
      label: 'Fecha',
      render: (_: unknown, invoice: Invoice) => (
        <div className="text-sm text-gray-600">
          {new Date(invoice.created_at).toLocaleDateString('es-CO', {
            day: '2-digit',
            month: 'short',
            year: 'numeric',
          })}
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
        </div>
      ),
    },
    {
      key: 'actions',
<<<<<<< HEAD
      label: 'Acciones',
      render: (_: unknown, invoice: Invoice) => (
        <div className="flex gap-2">
          {invoice.pdf_url && (
            <Button
              variant="secondary"
              size="sm"
              onClick={() => window.open(invoice.pdf_url, '_blank')}
            >
              Ver PDF
            </Button>
          )}
          {invoice.status === 'failed' && (
            <Button
              variant="primary"
              size="sm"
              onClick={() => handleRetryInvoice(invoice)}
              disabled={actionLoading}
            >
              Reintentar
            </Button>
          )}
          {invoice.status === 'issued' && (
            <Button
              variant="danger"
              size="sm"
              onClick={() => {
                setSelectedInvoice(invoice);
                setShowCancelModal(true);
              }}
              disabled={actionLoading}
            >
              Cancelar
            </Button>
          )}
        </div>
=======
      label: '',
      width: '50px',
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
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
      ),
    },
  ];

  return (
    <>
<<<<<<< HEAD
      <Table
        data={invoices}
        columns={columns}
        emptyMessage="No hay facturas para mostrar"
      />

=======
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-2xl font-bold">Facturas Electrónicas</h2>
          <p className="text-sm text-gray-500 mt-1">
            {totalCount} {totalCount === 1 ? 'orden facturada' : 'órdenes facturadas'}
            <span className="ml-2 text-xs text-gray-400">
              (doble clic para ver reintentos y detalle)
            </span>
          </p>
        </div>
        <Button variant="primary" onClick={() => setShowBulkModal(true)}>
          + Crear Facturas desde Órdenes
        </Button>
      </div>

      {/* Tabla con paginación */}
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
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
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
<<<<<<< HEAD
=======

      {/* Modal de creación masiva */}
      <BulkCreateInvoiceModal
        isOpen={showBulkModal}
        onClose={() => setShowBulkModal(false)}
        onSuccess={() => loadInvoices(currentPage, pageSize)}
        businessId={businessId}
      />
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
    </>
  );
}
