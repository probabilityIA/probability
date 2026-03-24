/**
 * Modal de detalle de factura con historial de sincronización
 */

'use client';

import { useState, useEffect, useCallback } from 'react';
import { XMarkIcon, ClipboardDocumentIcon, ClipboardDocumentCheckIcon } from '@heroicons/react/24/outline';
import { Button } from '@/shared/ui/button';
import { Badge } from '@/shared/ui/badge';
import { Spinner } from '@/shared/ui/spinner';
import { useToast } from '@/shared/providers/toast-provider';
import {
  getInvoiceSyncLogsAction,
  cancelRetryAction,
  enableRetryAction,
  retryInvoiceAction,
  deletePendingInvoiceAction,
} from '../../infra/actions';
import { useInvoiceSSE } from '../hooks/useInvoiceSSE';
import type { Invoice, SyncLog, InvoiceSSEEventData } from '../../domain/types';

interface InvoiceDetailModalProps {
  invoice: Invoice | null;
  isOpen: boolean;
  onClose: () => void;
  onCancel: (invoice: Invoice) => void;
  onRefresh: () => void;
  onDelete?: () => void;
  businessId: number;
}

export function InvoiceDetailModal({
  invoice,
  isOpen,
  onClose,
  onCancel,
  onRefresh,
  onDelete,
  businessId,
}: InvoiceDetailModalProps) {
  const { showToast } = useToast();
  const [syncLogs, setSyncLogs] = useState<SyncLog[]>([]);
  const [loadingLogs, setLoadingLogs] = useState(true);
  const [cancellingRetry, setCancellingRetry] = useState(false);
  const [retrying, setRetrying] = useState(false);
  const [retryProgress, setRetryProgress] = useState(0);
  const [retryResult, setRetryResult] = useState<'success' | 'failed' | 'pending_validation' | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [copiedField, setCopiedField] = useState<string | null>(null);

  const copyToClipboard = (text: string, fieldId: string) => {
    navigator.clipboard.writeText(text);
    setCopiedField(fieldId);
    setTimeout(() => setCopiedField(null), 2000);
  };

  const CopyButton = ({ text, fieldId }: { text: string; fieldId: string }) => {
    const isCopied = copiedField === fieldId;
    return (
      <button
        onClick={() => copyToClipboard(text, fieldId)}
        className="inline-flex items-center p-0.5 text-gray-400 hover:text-gray-600 dark:text-gray-300 transition-colors"
        title="Copiar"
      >
        {isCopied ? (
          <ClipboardDocumentCheckIcon className="w-3.5 h-3.5 text-green-500" />
        ) : (
          <ClipboardDocumentIcon className="w-3.5 h-3.5" />
        )}
      </button>
    );
  };

  // SSE para escuchar resultado del retry en tiempo real
  // NOTA: No mostrar toasts aquí - InvoiceList ya los muestra (evitar duplicados)
  const handleInvoiceCreated = useCallback((data: InvoiceSSEEventData) => {
    if (!invoice || !retrying) return;
    if (data.invoice_id === invoice.id || data.order_id === invoice.order_id) {
      setRetryProgress(100);
      setRetryResult('success');
      setRetrying(false);
      loadSyncLogs();
      onRefresh();
    }
  }, [invoice, retrying]);

  const handleInvoiceFailed = useCallback((data: InvoiceSSEEventData) => {
    if (!invoice || !retrying) return;
    if (data.invoice_id === invoice.id || data.order_id === invoice.order_id) {
      setRetryProgress(100);
      setRetryResult('failed');
      setRetrying(false);
      loadSyncLogs();
      onRefresh();
    }
  }, [invoice, retrying]);

  const handleInvoicePendingValidation = useCallback((data: InvoiceSSEEventData) => {
    if (!invoice || !retrying) return;
    if (data.invoice_id === invoice.id || data.order_id === invoice.order_id) {
      setRetryProgress(100);
      setRetryResult('pending_validation');
      setRetrying(false);
      loadSyncLogs();
      onRefresh();
    }
  }, [invoice, retrying]);

  useInvoiceSSE({
    businessId,
    onInvoiceCreated: handleInvoiceCreated,
    onInvoiceFailed: handleInvoiceFailed,
    onInvoicePendingValidation: handleInvoicePendingValidation,
  });

  useEffect(() => {
    if (isOpen && invoice) {
      loadSyncLogs();
      setRetrying(false);
      setRetryProgress(0);
      setRetryResult(null);
    } else {
      setSyncLogs([]);
    }
  }, [isOpen, invoice?.id]);

  // Progreso simulado mientras espera SSE
  useEffect(() => {
    if (!retrying) return;
    setRetryProgress(5);
    const interval = setInterval(() => {
      setRetryProgress(prev => {
        if (prev >= 85) { clearInterval(interval); return 85; }
        return prev + Math.random() * 10;
      });
    }, 500);
    return () => clearInterval(interval);
  }, [retrying]);

  const loadSyncLogs = async () => {
    if (!invoice) return;
    try {
      setLoadingLogs(true);
      const logs = await getInvoiceSyncLogsAction(invoice.id);
      setSyncLogs(logs);
    } catch {
      setSyncLogs([]);
    } finally {
      setLoadingLogs(false);
    }
  };

  const handleRetry = async () => {
    if (!invoice) return;
    try {
      setRetrying(true);
      setRetryProgress(0);
      setRetryResult(null);
      await retryInvoiceAction(invoice.id);
      // No cerrar ni mostrar éxito aquí - SSE lo hará cuando llegue el resultado
    } catch (error: any) {
      setRetrying(false);
      setRetryProgress(0);
      showToast('Error al reintentar: ' + error.message, 'error');
    }
  };

  const handleToggleAutoRetry = async () => {
    if (!invoice) return;
    try {
      setCancellingRetry(true);
      if (autoRetriesEnabled) {
        await cancelRetryAction(invoice.id);
        showToast('Reintentos automáticos deshabilitados', 'success');
      } else {
        await enableRetryAction(invoice.id);
        showToast('Reintentos automáticos habilitados', 'success');
      }
      loadSyncLogs();
      onRefresh();
    } catch (error: any) {
      showToast('Error: ' + error.message, 'error');
    } finally {
      setCancellingRetry(false);
    }
  };

  const handleDelete = async () => {
    if (!invoice) return;
    if (!confirm('¿Estás seguro de eliminar esta factura? Esta acción no se puede deshacer.')) return;
    try {
      setDeleting(true);
      await deletePendingInvoiceAction(invoice.id);
      showToast('Factura eliminada exitosamente', 'success');
      onClose();
      if (onDelete) onDelete();
      else onRefresh();
    } catch (error: any) {
      showToast('Error al eliminar: ' + error.message, 'error');
    } finally {
      setDeleting(false);
    }
  };

  const hasPendingRetries = syncLogs.some(
    log => (log.status === 'failed' || log.status === 'pending') && log.next_retry_at
  );

  // Detectar si los reintentos automáticos están cancelados
  const hasCancelledRetries = syncLogs.some(
    log => log.status === 'cancelled'
  );

  // Calcular estado de reintentos desde el último sync log
  const lastLog = syncLogs.length > 0 ? syncLogs[0] : null;
  // Para pending (check_status): sin límite práctico. Para failed (retry): límite normal.
  const maxRetriesReached = invoice?.status === 'pending'
    ? false
    : (lastLog ? lastLog.retry_count >= lastLog.max_retries : false);
  const retriesUsed = lastLog ? lastLog.retry_count : 0;
  const maxRetries = lastLog ? lastLog.max_retries : 3;

  // Estado del toggle: reintentos activos o deshabilitados
  const autoRetriesEnabled = hasPendingRetries;
  const autoRetriesDisabled = hasCancelledRetries && !hasPendingRetries;

  // Puede eliminarse: pending + 3+ intentos de consulta (query)
  const queryAttempts = syncLogs.filter(log => log.operation_type === 'query').length;
  const canDelete = invoice?.status === 'pending' && queryAttempts >= 3;

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString('es-CO', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  const getStatusBadge = (status: string) => {
    const config: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'secondary' | 'primary' }> = {
      issued: { label: 'Emitida', type: 'success' },
      pending: { label: 'Pendiente', type: 'warning' },
      failed: { label: 'Fallida', type: 'error' },
      cancelled: { label: 'Cancelada', type: 'secondary' },
    };
    const c = config[status] || { label: status, type: 'secondary' as const };
    return <Badge type={c.type}>{c.label}</Badge>;
  };

  const getSyncStatusBadge = (status: string) => {
    const config: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'secondary' | 'primary' }> = {
      success: { label: 'Exitoso', type: 'success' },
      processing: { label: 'Procesando', type: 'primary' },
      pending: { label: 'Pendiente', type: 'warning' },
      failed: { label: 'Fallido', type: 'error' },
      cancelled: { label: 'Cancelado', type: 'secondary' },
    };
    const c = config[status] || { label: status, type: 'secondary' as const };
    return <Badge type={c.type}>{c.label}</Badge>;
  };

  const getTriggerLabel = (trigger: string) => {
    const labels: Record<string, string> = {
      auto: 'Automático',
      manual: 'Manual',
      retry_job: 'Reintento',
    };
    return labels[trigger] || trigger;
  };

  if (!isOpen || !invoice) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-center justify-center p-4">
        {/* Backdrop - fondo claro opaco, se ve la plataforma detrás */}
        <div
          className="fixed inset-0 bg-white dark:bg-gray-800/60 backdrop-blur-sm transition-opacity"
          onClick={onClose}
        />

        {/* Modal */}
        <div className="relative bg-white dark:bg-gray-800 rounded-lg shadow-2xl border border-gray-200 dark:border-gray-700 max-w-2xl w-full max-h-[85vh] flex flex-col">
          {/* Header */}
          <div className="flex items-center justify-between px-6 py-4 border-b">
            <div className="flex items-center gap-3">
              <h2 className="text-lg font-bold">
                Factura {invoice.invoice_number || `#${invoice.id}`}
              </h2>
              {getStatusBadge(invoice.status)}
            </div>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 dark:text-gray-300"
            >
              <XMarkIcon className="w-6 h-6" />
            </button>
          </div>

          {/* Content */}
          <div className="flex-1 overflow-y-auto p-6">
            {/* Info de la factura */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Orden</p>
                <p className="font-mono text-sm mt-1 flex items-center gap-1">
                  {invoice.order_id}
                  <CopyButton text={invoice.order_id} fieldId="order_id" />
                </p>
              </div>
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Cliente</p>
                <p className="font-medium text-sm mt-1">{invoice.customer_name || '-'}</p>
              </div>
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Total</p>
                <p className="font-semibold text-sm mt-1">
                  {new Intl.NumberFormat('es-CO', {
                    style: 'currency',
                    currency: invoice.currency || 'COP',
                  }).format(invoice.total_amount)}
                </p>
              </div>
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Creada</p>
                <p className="text-sm mt-1">{formatDate(invoice.created_at)}</p>
              </div>
            </div>

            {/* Datos de la factura emitida */}
            {invoice.status === 'issued' && (invoice.cufe || invoice.pdf_url || invoice.xml_url || invoice.invoice_url) && (
              <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
                <p className="text-xs text-green-600 uppercase tracking-wide font-semibold mb-3">Datos de Factura Electrónica</p>
                <div className="space-y-2">
                  {invoice.cufe && (
                    <div className="flex items-start gap-2">
                      <span className="text-xs text-gray-500 dark:text-gray-400 w-12 shrink-0 pt-0.5">CUFE</span>
                      <span className="text-xs font-mono text-gray-700 dark:text-gray-200 break-all flex-1">{invoice.cufe}</span>
                      <CopyButton text={invoice.cufe} fieldId="cufe" />
                    </div>
                  )}
                  <div className="flex flex-wrap gap-2 mt-2">
                    {invoice.invoice_url && (
                      <a href={invoice.invoice_url} target="_blank" rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 px-3 py-1.5 bg-white dark:bg-gray-800 border border-green-300 rounded-md text-xs font-medium text-green-700 hover:bg-green-100 transition-colors">
                        Ver Factura
                      </a>
                    )}
                    {invoice.pdf_url && (
                      <a href={invoice.pdf_url} target="_blank" rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 px-3 py-1.5 bg-white dark:bg-gray-800 border border-green-300 rounded-md text-xs font-medium text-green-700 hover:bg-green-100 transition-colors">
                        Descargar PDF
                      </a>
                    )}
                    {invoice.xml_url && (
                      <a href={invoice.xml_url} target="_blank" rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 px-3 py-1.5 bg-white dark:bg-gray-800 border border-green-300 rounded-md text-xs font-medium text-green-700 hover:bg-green-100 transition-colors">
                        Descargar XML
                      </a>
                    )}
                  </div>
                </div>
              </div>
            )}

            {/* Documento completo del proveedor (Softpymes) */}
            {invoice.provider_response && (
              <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
                <p className="text-xs text-blue-600 uppercase tracking-wide font-semibold mb-3">
                  Vista Previa de Factura
                </p>
                <details className="group" open>
                  <summary className="text-xs text-gray-700 dark:text-gray-200 cursor-pointer hover:text-blue-600 font-medium flex items-center gap-2">
                    <span>Detalles del documento</span>
                    <span className="text-gray-400 group-open:rotate-180 transition-transform">▼</span>
                  </summary>
                  <div className="mt-3 space-y-3">
                    {/* Info básica del documento */}
                    <div className="grid grid-cols-2 gap-3">
                      {invoice.provider_response.documentNumber && (
                        <div>
                          <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">Número de Documento</p>
                          <p className="text-xs font-mono font-medium text-gray-700 dark:text-gray-200">
                            {invoice.provider_response.documentNumber}
                          </p>
                        </div>
                      )}
                      {invoice.provider_response.documentDate && (
                        <div>
                          <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">Fecha</p>
                          <p className="text-xs font-medium text-gray-700 dark:text-gray-200">
                            {invoice.provider_response.documentDate}
                          </p>
                        </div>
                      )}
                      {invoice.provider_response.customerName && (
                        <div>
                          <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">Cliente</p>
                          <p className="text-xs font-medium text-gray-700 dark:text-gray-200">
                            {invoice.provider_response.customerName}
                          </p>
                        </div>
                      )}
                      {invoice.provider_response.customerIdentification && (
                        <div>
                          <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">Identificación</p>
                          <p className="text-xs font-mono font-medium text-gray-700 dark:text-gray-200">
                            {invoice.provider_response.customerIdentification}
                          </p>
                        </div>
                      )}
                    </div>

                    {/* Totales */}
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-3 p-3 bg-white dark:bg-gray-800/60 rounded">
                      {invoice.provider_response.total && (
                        <div>
                          <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">Total</p>
                          <p className="text-xs font-semibold text-gray-900 dark:text-white">
                            ${invoice.provider_response.total}
                          </p>
                        </div>
                      )}
                      {invoice.provider_response.totalIva && (
                        <div>
                          <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">IVA</p>
                          <p className="text-xs font-medium text-gray-700 dark:text-gray-200">
                            ${invoice.provider_response.totalIva}
                          </p>
                        </div>
                      )}
                      {invoice.provider_response.totalDiscount && (
                        <div>
                          <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">Descuento</p>
                          <p className="text-xs font-medium text-gray-700 dark:text-gray-200">
                            ${invoice.provider_response.totalDiscount}
                          </p>
                        </div>
                      )}
                      {invoice.provider_response.totalWithholdingTax && (
                        <div>
                          <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">Retención</p>
                          <p className="text-xs font-medium text-gray-700 dark:text-gray-200">
                            ${invoice.provider_response.totalWithholdingTax}
                          </p>
                        </div>
                      )}
                    </div>

                    {/* Items del documento */}
                    {invoice.provider_response.details && invoice.provider_response.details.length > 0 && (
                      <div>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mb-2">Items ({invoice.provider_response.details.length})</p>
                        <div className="space-y-2">
                          {invoice.provider_response.details.map((detail: any, idx: number) => (
                            <div key={idx} className="p-2 bg-white dark:bg-gray-800/80 rounded text-xs border border-gray-200 dark:border-gray-700">
                              <div className="flex justify-between items-start mb-1">
                                <span className="font-medium text-gray-900 dark:text-white">{detail.itemName || detail.itemCode}</span>
                                <span className="font-semibold text-gray-900 dark:text-white">${detail.value}</span>
                              </div>
                              <div className="flex gap-3 text-gray-600 dark:text-gray-300">
                                <span>Cant: {detail.quantity}</span>
                                {detail.iva && <span>IVA: ${detail.iva}</span>}
                                {detail.discount && detail.discount !== '0' && detail.discount !== '0.00' && (
                                  <span>Desc: ${detail.discount}</span>
                                )}
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Información de envío */}
                    {invoice.provider_response.shipInformation && (
                      <div className="p-3 bg-white dark:bg-gray-800/60 rounded">
                        <p className="text-xs text-gray-500 dark:text-gray-400 mb-2 font-medium">Información de Envío</p>
                        <div className="space-y-1 text-xs text-gray-700 dark:text-gray-200">
                          {invoice.provider_response.shipInformation.shipAddress && (
                            <p>📍 {invoice.provider_response.shipInformation.shipAddress}</p>
                          )}
                          {invoice.provider_response.shipInformation.shipCity && (
                            <p>🏙️ {invoice.provider_response.shipInformation.shipCity}, {invoice.provider_response.shipInformation.shipDepartment}</p>
                          )}
                          {invoice.provider_response.shipInformation.shipPhone && (
                            <p>📞 {invoice.provider_response.shipInformation.shipPhone}</p>
                          )}
                        </div>
                      </div>
                    )}

                    {/* Recibo de Caja */}
                    {invoice.provider_response.cash_receipt && (
                      <div className={`p-3 rounded ${invoice.provider_response.cash_receipt.status === 'success' ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'}`}>
                        <p className={`text-xs uppercase tracking-wide font-medium mb-2 ${invoice.provider_response.cash_receipt.status === 'success' ? 'text-green-600' : 'text-red-600'}`}>
                          Recibo de Caja {invoice.provider_response.cash_receipt.status === 'success' ? '— Registrado' : '— Error'}
                        </p>
                        {invoice.provider_response.cash_receipt.status === 'success' ? (
                          <div className="grid grid-cols-2 gap-2 text-xs">
                            <div>
                              <span className="text-gray-500 dark:text-gray-400">Medio de pago: </span>
                              <span className="font-medium text-gray-700 dark:text-gray-200">{invoice.provider_response.cash_receipt.payment_type}</span>
                            </div>
                            <div>
                              <span className="text-gray-500 dark:text-gray-400">Monto: </span>
                              <span className="font-medium text-gray-700 dark:text-gray-200">${Number(invoice.provider_response.cash_receipt.amount).toLocaleString()}</span>
                            </div>
                            {invoice.provider_response.cash_receipt.message && (
                              <div className="col-span-2">
                                <span className="text-gray-500 dark:text-gray-400">Respuesta: </span>
                                <span className="font-medium text-gray-700 dark:text-gray-200">{invoice.provider_response.cash_receipt.message}</span>
                              </div>
                            )}
                          </div>
                        ) : (
                          <p className="text-xs text-red-700">{invoice.provider_response.cash_receipt.error}</p>
                        )}
                      </div>
                    )}

                    {/* JSON completo (colapsado) */}
                    <details className="mt-3">
                      <summary className="text-xs text-gray-500 dark:text-gray-400 cursor-pointer hover:text-gray-700 dark:text-gray-200 flex items-center gap-1">
                        <span>Ver JSON completo</span>
                        <CopyButton
                          text={JSON.stringify(invoice.provider_response, null, 2)}
                          fieldId="full-document-json"
                        />
                      </summary>
                      <pre className="mt-2 text-xs bg-white dark:bg-gray-800/80 rounded p-3 overflow-x-auto max-h-64 border border-gray-200 dark:border-gray-700 font-mono text-gray-700 dark:text-gray-200">
                        {JSON.stringify(invoice.provider_response, null, 2)}
                      </pre>
                    </details>
                  </div>
                </details>
              </div>
            )}

            {/* Error message si existe */}
            {invoice.error_message && (
              <div className="mb-6 p-3 bg-red-50 border border-red-200 rounded-lg">
                <p className="text-xs text-red-500 uppercase tracking-wide mb-1">Error</p>
                <p className="text-sm text-red-700 font-mono break-all">{invoice.error_message}</p>
              </div>
            )}

            {/* Barra de progreso del retry */}
            {(retrying || retryResult) && (
              <div className="mb-6 p-4 bg-gray-50 border border-gray-200 dark:border-gray-700 rounded-lg">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-200">
                    {retryResult === 'success'
                      ? 'Factura emitida exitosamente'
                      : retryResult === 'failed'
                        ? (invoice.status === 'pending' ? 'Consulta fallida' : 'Reintento fallido')
                        : retryResult === 'pending_validation'
                          ? 'DIAN aún validando — se consultará de nuevo automáticamente'
                          : (invoice.status === 'pending' ? 'Consultando estado DIAN...' : 'Reintentando emisión...')}
                  </span>
                  <span className="text-sm text-gray-500 dark:text-gray-400">
                    {Math.round(retryProgress)}%
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2.5">
                  <div
                    className={`h-2.5 rounded-full transition-all duration-300 ${
                      retryResult === 'success'
                        ? 'bg-green-500'
                        : retryResult === 'failed'
                          ? 'bg-red-500'
                          : retryResult === 'pending_validation'
                            ? 'bg-amber-500'
                            : 'bg-blue-600'
                    }`}
                    style={{ width: `${Math.min(retryProgress, 100)}%` }}
                  />
                </div>
              </div>
            )}

            {/* Acciones */}
            <div className="flex gap-2 mb-6 pb-6 border-b border-gray-200 dark:border-gray-700">
              {invoice.status === 'failed' && (
                <Button
                  variant="primary"
                  size="sm"
                  onClick={handleRetry}
                  disabled={retrying || maxRetriesReached}
                >
                  {retrying ? 'Reintentando...' : 'Reintentar'}
                </Button>
              )}
              {invoice.status === 'pending' && (
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={handleRetry}
                  disabled={retrying || maxRetriesReached}
                >
                  {retrying ? 'Consultando...' : 'Consultar Estado DIAN'}
                </Button>
              )}
              {(autoRetriesEnabled || autoRetriesDisabled) && (
                <Button
                  variant={autoRetriesEnabled ? 'danger' : 'secondary'}
                  size="sm"
                  onClick={handleToggleAutoRetry}
                  disabled={cancellingRetry}
                >
                  {cancellingRetry
                    ? (autoRetriesEnabled ? 'Deshabilitando...' : 'Habilitando...')
                    : autoRetriesEnabled
                      ? 'Deshabilitar Reintentos'
                      : 'Habilitar Reintentos'}
                </Button>
              )}
              {invoice.status === 'issued' && (
                <Button
                  variant="danger"
                  size="sm"
                  disabled
                  title="Funcionalidad en desarrollo"
                >
                  Cancelar Factura
                </Button>
              )}
              {canDelete && (
                <Button
                  variant="danger"
                  size="sm"
                  onClick={handleDelete}
                  disabled={deleting}
                >
                  {deleting ? 'Eliminando...' : 'Eliminar Factura'}
                </Button>
              )}
            </div>

            {/* Historial de sincronización */}
            <div>
              <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">
                Historial de Sincronización
              </h4>

              {loadingLogs ? (
                <div className="flex justify-center py-6">
                  <Spinner />
                </div>
              ) : syncLogs.length === 0 ? (
                <p className="text-sm text-gray-500 dark:text-gray-400 py-4 text-center">
                  Sin registros de sincronización
                </p>
              ) : (
                <div className="space-y-3">
                  {syncLogs.map((log) => (
                    <div
                      key={log.id}
                      className={`border rounded-lg p-4 ${
                        log.status === 'success'
                          ? 'border-green-200 bg-green-50'
                          : log.status === 'failed'
                            ? 'border-red-200 bg-red-50'
                            : log.status === 'cancelled'
                              ? 'border-gray-200 dark:border-gray-700 bg-gray-50'
                              : 'border-yellow-200 bg-yellow-50'
                      }`}
                    >
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center gap-3">
                          {getSyncStatusBadge(log.status)}
                          <span className="text-xs text-gray-500 dark:text-gray-400">
                            {getTriggerLabel(log.triggered_by)}
                          </span>
                          {log.duration_ms && (
                            <span className="text-xs text-gray-400">
                              {log.duration_ms}ms
                            </span>
                          )}
                        </div>
                        <span className="text-xs text-gray-500 dark:text-gray-400">
                          {formatDate(log.created_at)}
                        </span>
                      </div>

                      {/* Info de reintentos */}
                      <div className="flex items-center gap-4 text-xs text-gray-600 dark:text-gray-300">
                        <span>Intento {log.retry_count + 1} de {log.max_retries}</span>
                        {log.next_retry_at && (log.status === 'failed' || log.status === 'pending') && (
                          <span className="text-orange-600">
                            Próximo reintento: {formatDate(log.next_retry_at)}
                          </span>
                        )}
                      </div>

                      {/* Error message */}
                      {log.error_message && (
                        <div className="mt-2 p-2 bg-white dark:bg-gray-800/60 rounded text-xs text-red-700 font-mono break-all">
                          {log.error_message}
                        </div>
                      )}

                      {/* Request/Response audit data */}
                      {(log.request_payload || log.response_body) && (
                        <details className="mt-2">
                          <summary className="text-xs text-gray-500 dark:text-gray-400 cursor-pointer hover:text-gray-700 dark:text-gray-200">
                            Ver request/response
                          </summary>
                          <div className="mt-1 space-y-1">
                            {log.request_url && (
                              <div className="text-xs font-mono text-gray-600 dark:text-gray-300">
                                URL: {log.request_url}
                              </div>
                            )}
                            {log.response_status != null && log.response_status > 0 && (
                              <div className="text-xs font-mono text-gray-600 dark:text-gray-300">
                                Status: {log.response_status}
                              </div>
                            )}
                            {log.request_payload && (
                              <div>
                                <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5 flex items-center gap-1">
                                  Request:
                                  <CopyButton
                                    text={JSON.stringify(log.request_payload, null, 2)}
                                    fieldId={`req-${log.id}`}
                                  />
                                </p>
                                <pre className="text-xs bg-white dark:bg-gray-800/60 rounded p-2 overflow-x-auto max-h-32">
                                  {JSON.stringify(log.request_payload, null, 2)}
                                </pre>
                              </div>
                            )}
                            {log.response_body && (
                              <div>
                                <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5 flex items-center gap-1">
                                  Response:
                                  <CopyButton
                                    text={JSON.stringify(log.response_body, null, 2)}
                                    fieldId={`res-${log.id}`}
                                  />
                                </p>
                                <pre className="text-xs bg-white dark:bg-gray-800/60 rounded p-2 overflow-x-auto max-h-32">
                                  {JSON.stringify(log.response_body, null, 2)}
                                </pre>
                              </div>
                            )}
                          </div>
                        </details>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
