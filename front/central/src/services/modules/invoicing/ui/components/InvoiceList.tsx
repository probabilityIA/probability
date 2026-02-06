/**
 * Componente para listar facturas
 */

'use client';

import { useState, useEffect } from 'react';
import { Table } from '@/shared/ui/table';
import { Button } from '@/shared/ui/button';
import { Badge } from '@/shared/ui/badge';
import { Spinner } from '@/shared/ui/spinner';
import { useToast } from '@/shared/providers/toast-provider';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { BulkCreateInvoiceModal } from './BulkCreateInvoiceModal';
import { getInvoicesAction, cancelInvoiceAction, retryInvoiceAction } from '../../infra/actions';
import type { Invoice, InvoiceFilters } from '../../domain/types';

interface InvoiceListProps {
  businessId: number;
  filters?: InvoiceFilters;
}

export function InvoiceList({ businessId, filters = {} }: InvoiceListProps) {
  const { showToast } = useToast();
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [showBulkModal, setShowBulkModal] = useState(false);
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

    try {
      setActionLoading(true);
      await cancelInvoiceAction(selectedInvoice.id);
      showToast('Factura cancelada exitosamente', 'success');
      setShowCancelModal(false);
      loadInvoices();
    } catch (error: any) {
      showToast('Error al cancelar factura: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

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

  const columns = [
    {
      key: 'invoice_number',
      label: 'Número',
      render: (_: unknown, invoice: Invoice) => (
        <div className="font-medium">{invoice.invoice_number || 'N/A'}</div>
      ),
    },
    {
      key: 'order_id',
      label: 'Orden',
      render: (_: unknown, invoice: Invoice) => (
        <div className="text-sm text-gray-600">{invoice.order_id}</div>
      ),
    },
    {
      key: 'customer_name',
      label: 'Cliente',
      render: (_: unknown, invoice: Invoice) => (
        <div>
          <div className="font-medium">{invoice.customer_name}</div>
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
      key: 'issued_at',
      label: 'Fecha',
      render: (_: unknown, invoice: Invoice) => (
        <div className="text-sm text-gray-600">
          {invoice.issued_at
            ? new Date(invoice.issued_at).toLocaleDateString('es-CO')
            : 'Pendiente'}
        </div>
      ),
    },
    {
      key: 'actions',
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
      ),
    },
  ];

  return (
    <>
      {/* Header con botón de creación masiva */}
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Facturas Electrónicas</h2>
        <Button
          variant="primary"
          onClick={() => setShowBulkModal(true)}
        >
          + Crear Facturas desde Órdenes
        </Button>
      </div>

      <Table
        data={invoices}
        columns={columns}
        emptyMessage="No hay facturas para mostrar"
      />

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

      <BulkCreateInvoiceModal
        isOpen={showBulkModal}
        onClose={() => setShowBulkModal(false)}
        onSuccess={() => {
          // Recargar lista de facturas después de crear
          loadInvoices();
        }}
      />
    </>
  );
}
