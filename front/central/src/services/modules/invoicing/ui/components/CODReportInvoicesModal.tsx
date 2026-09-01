'use client';

import { useEffect, useState } from 'react';
import { XMarkIcon, ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/24/outline';
import { getCODReportInvoicesAction } from '../../infra/actions';
import type { Invoice } from '../../domain/types';

interface CODReportInvoicesModalProps {
  isOpen: boolean;
  onClose: () => void;
  businessId: number | null;
  accountNumber: string;
  isCOD: boolean;
  startDate?: string;
  endDate?: string;
}

const PAGE_SIZE = 10;

function formatCurrency(value: number): string {
  return new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(value);
}

function formatDate(value?: string): string {
  if (!value) return '-';
  return new Date(value).toLocaleDateString('es-CO', { year: 'numeric', month: 'short', day: 'numeric' });
}

export function CODReportInvoicesModal({
  isOpen,
  onClose,
  businessId,
  accountNumber,
  isCOD,
  startDate,
  endDate,
}: CODReportInvoicesModalProps) {
  const [page, setPage] = useState(1);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!isOpen) return;
    setPage(1);
  }, [isOpen, accountNumber, isCOD]);

  useEffect(() => {
    if (!isOpen || !businessId) return;
    let cancelled = false;
    setLoading(true);
    setError('');
    getCODReportInvoicesAction(businessId, accountNumber, isCOD, startDate, endDate, page, PAGE_SIZE)
      .then((res) => {
        if (cancelled) return;
        setInvoices(res.data);
        setTotal(res.total);
      })
      .catch((err: any) => {
        if (cancelled) return;
        setError('Error al cargar las facturas: ' + err.message);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [isOpen, businessId, accountNumber, isCOD, startDate, endDate, page]);

  if (!isOpen) return null;

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />

      <div className="relative bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-3xl flex flex-col max-h-[85vh]">
        <div className="flex items-center justify-between px-6 pt-6 pb-4 border-b border-gray-100 flex-shrink-0">
          <div>
            <h2 className="text-lg font-bold text-gray-900 dark:text-white">Facturas de la cuenta {accountNumber}</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
              {isCOD ? 'Contra entrega' : 'Contado'} &middot; {total} factura{total === 1 ? '' : 's'}
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 dark:text-gray-300 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto px-6 py-4 min-h-[300px]">
          {loading && (
            <div className="flex flex-col items-center justify-center py-16 text-gray-500 dark:text-gray-400">
              <div className="w-10 h-10 border-4 border-purple-200 border-t-purple-600 rounded-full animate-spin mb-3" />
              <p className="text-sm font-medium">Cargando facturas...</p>
            </div>
          )}

          {!loading && error && <p className="text-sm text-red-600">{error}</p>}

          {!loading && !error && invoices.length === 0 && (
            <div className="text-center py-16 text-gray-400 text-sm">Sin facturas para mostrar</div>
          )}

          {!loading && !error && invoices.length > 0 && (
            <div className="overflow-x-auto rounded-xl border border-gray-200 dark:border-gray-700">
              <table className="w-full text-sm">
                <thead>
                  <tr className="bg-gray-50 dark:bg-gray-700 text-gray-500 dark:text-gray-300 text-xs uppercase tracking-wide">
                    <th className="px-3 py-2 text-left">Factura</th>
                    <th className="px-3 py-2 text-left">Orden</th>
                    <th className="px-3 py-2 text-left">Cliente</th>
                    <th className="px-3 py-2 text-left">Fecha</th>
                    <th className="px-3 py-2 text-right">Total</th>
                  </tr>
                </thead>
                <tbody>
                  {invoices.map((inv) => (
                    <tr key={inv.id} className="border-b border-gray-100 dark:border-gray-700">
                      <td className="px-3 py-2 font-mono text-gray-700 dark:text-gray-200">{inv.invoice_number || '-'}</td>
                      <td className="px-3 py-2 font-mono text-gray-500 dark:text-gray-400">{inv.order_number || '-'}</td>
                      <td className="px-3 py-2 text-gray-700 dark:text-gray-200 truncate max-w-[160px]">{inv.customer_name || '-'}</td>
                      <td className="px-3 py-2 text-gray-500 dark:text-gray-400 whitespace-nowrap">
                        {formatDate(inv.issued_at || inv.created_at)}
                      </td>
                      <td className="px-3 py-2 text-right font-semibold text-gray-900 dark:text-white whitespace-nowrap">
                        {formatCurrency(inv.total_amount)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>

        <div className="flex items-center justify-between px-6 py-4 border-t border-gray-100 flex-shrink-0">
          <span className="text-xs text-gray-400">
            Página {page} de {totalPages}
          </span>
          <div className="flex gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page <= 1 || loading}
              className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 text-gray-500 dark:text-gray-300 disabled:opacity-40 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
              <ChevronLeftIcon className="w-4 h-4" />
            </button>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages || loading}
              className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 text-gray-500 dark:text-gray-300 disabled:opacity-40 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
              <ChevronRightIcon className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
