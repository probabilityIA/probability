'use client';

import { useState, useCallback, useMemo } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { DateRangePicker } from '@/shared/ui/date-range-picker';
import type { CompareResponseData } from '../../domain/types';

interface SyncCancellationsModalProps {
  isOpen: boolean;
  onClose: () => void;
  loading: boolean;
  syncData: CompareResponseData | null;
  onRequestSync: (dateFrom: string, dateTo: string) => void;
}

const PAGE_SIZE = 25;

function formatCurrency(value: string | number | undefined): string {
  if (value === undefined || value === null || value === '') return '—';
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(num)) return String(value);
  return new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(num);
}

export function SyncCancellationsModal({
  isOpen,
  onClose,
  loading,
  syncData,
  onRequestSync,
}: SyncCancellationsModalProps) {
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');
  const [validationError, setValidationError] = useState('');
  const [page, setPage] = useState(1);

  const handleRangeChange = useCallback((from: string | undefined, to: string | undefined) => {
    setDateFrom(from ?? '');
    setDateTo(to ?? '');
    setValidationError('');
  }, []);

  const handleStart = useCallback(() => {
    setValidationError('');
    if (!dateFrom || !dateTo) {
      setValidationError('Selecciona las fechas de inicio y fin');
      return;
    }
    const diffDays = Math.ceil(
      (new Date(dateTo).getTime() - new Date(dateFrom).getTime()) / 86400000
    );
    if (diffDays < 0) {
      setValidationError('La fecha inicial debe ser anterior a la final');
      return;
    }
    if (diffDays > 30) {
      setValidationError('El rango máximo es de 30 días');
      return;
    }
    setPage(1);
    onRequestSync(dateFrom, dateTo);
  }, [dateFrom, dateTo, onRequestSync]);

  const handleNewSync = () => {
    setDateFrom('');
    setDateTo('');
    setValidationError('');
    setPage(1);
  };

  const results = syncData?.results ?? [];
  const totalPages = Math.max(1, Math.ceil(results.length / PAGE_SIZE));
  const pageRows = useMemo(
    () => results.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE),
    [results, page]
  );

  if (!isOpen) return null;

  const annulled = syncData?.summary.annulled_in_provider ?? 0;
  const released = syncData?.summary.released ?? 0;
  const reviewed = (syncData?.summary.matched ?? 0) + annulled;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />

      <div className="relative bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-5xl flex flex-col max-h-[92vh]">
        <div className="flex items-center justify-between px-6 pt-6 pb-5 border-b border-gray-100 flex-shrink-0">
          <div>
            <h2 className="text-xl font-bold text-gray-900 dark:text-white">Sincronizar facturas anuladas</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
              Detecta facturas anuladas en el proveedor, las cancela en el sistema y libera sus órdenes para re-facturar
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 dark:text-gray-300 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>

        {!syncData && !loading && (
          <div className="px-6 pt-5 pb-2 flex-shrink-0">
            <div className="flex items-end gap-4">
              <div className="flex-1">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1.5">
                  Rango de fechas{' '}
                  <span className="text-gray-400 font-normal">(máximo 30 días)</span>
                </label>
                <DateRangePicker
                  startDate={dateFrom}
                  endDate={dateTo}
                  onChange={handleRangeChange}
                  placeholder="Seleccionar rango de fechas"
                  className="w-full"
                />
              </div>
              <button
                onClick={handleStart}
                disabled={!dateFrom || !dateTo}
                className="px-5 py-2 bg-purple-600 hover:bg-purple-700 disabled:bg-purple-300 disabled:cursor-not-allowed text-white rounded-lg text-sm font-semibold transition-colors whitespace-nowrap h-[38px]"
              >
                Sincronizar
              </button>
            </div>
            {validationError && <p className="mt-2 text-sm text-red-600">{validationError}</p>}
            <div className="h-80" />
          </div>
        )}

        {loading && !syncData && (
          <div className="flex flex-col items-center justify-center py-20 text-gray-500 dark:text-gray-400 flex-shrink-0">
            <div className="w-12 h-12 border-4 border-purple-200 border-t-purple-600 rounded-full animate-spin mb-4" />
            <p className="text-sm font-medium">Consultando el proveedor y liberando órdenes…</p>
            <p className="text-xs mt-1 text-gray-400">Esto puede tardar unos segundos</p>
          </div>
        )}

        {syncData && (
          <div className="flex-1 overflow-y-auto px-6 py-5 min-h-0">
            <div className="flex gap-3 mb-5 flex-wrap items-center">
              <span className="flex items-center gap-2 px-4 py-2 bg-gray-50 border border-gray-200 rounded-full text-sm font-medium text-gray-700">
                🔍 {reviewed} revisadas
              </span>
              <span className="flex items-center gap-2 px-4 py-2 bg-red-50 border border-red-200 rounded-full text-sm font-medium text-red-800">
                🚫 {annulled} anuladas en proveedor
              </span>
              <span className="flex items-center gap-2 px-4 py-2 bg-green-50 border border-green-200 rounded-full text-sm font-medium text-green-800">
                ✅ {released} órdenes liberadas
              </span>
              <span className="ml-auto text-xs text-gray-400">
                {syncData.date_from} → {syncData.date_to}
              </span>
            </div>

            {results.length === 0 ? (
              <div className="text-center py-14 text-gray-400 text-sm">
                No se encontraron facturas anuladas en el rango seleccionado
              </div>
            ) : (
              <>
                <div className="overflow-x-auto rounded-xl border border-gray-200 dark:border-gray-700">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="bg-gradient-to-r from-purple-600 to-purple-700 text-white text-xs uppercase tracking-wide">
                        <th className="px-4 py-3 text-left rounded-tl-xl">Nro. Factura</th>
                        <th className="px-4 py-3 text-left">Fecha Factura</th>
                        <th className="px-4 py-3 text-left">Cliente</th>
                        <th className="px-4 py-3 text-right">Total Proveedor</th>
                        <th className="px-4 py-3 text-center rounded-tr-xl">Resultado</th>
                      </tr>
                    </thead>
                    <tbody>
                      {pageRows.map((row, idx) => (
                        <tr key={idx} className="border-b bg-red-50 border-red-100">
                          <td className="px-4 py-3 font-mono font-medium">
                            {row.prefix ? `${row.prefix}-` : ''}{row.invoice_number || '—'}
                          </td>
                          <td className="px-4 py-3 text-gray-600 dark:text-gray-300">{row.document_date || '—'}</td>
                          <td className="px-4 py-3 text-gray-600 dark:text-gray-300">
                            {row.customer_name || '—'}
                            {row.customer_nit ? <span className="text-gray-400"> · {row.customer_nit}</span> : null}
                          </td>
                          <td className="px-4 py-3 text-right font-semibold">{formatCurrency(row.provider_total)}</td>
                          <td className="px-4 py-3 text-center">
                            {row.released ? (
                              <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                                ✅ Liberada
                              </span>
                            ) : (
                              <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-700">
                                Sin cambio
                              </span>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {totalPages > 1 && (
                  <div className="mt-4 flex items-center justify-center gap-3 text-sm">
                    <button
                      onClick={() => setPage(p => Math.max(1, p - 1))}
                      disabled={page <= 1}
                      className="px-3 py-1.5 rounded-lg border border-gray-200 disabled:opacity-40 hover:bg-gray-50"
                    >
                      Anterior
                    </button>
                    <span className="text-gray-500">
                      Página {page} de {totalPages}
                    </span>
                    <button
                      onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                      disabled={page >= totalPages}
                      className="px-3 py-1.5 rounded-lg border border-gray-200 disabled:opacity-40 hover:bg-gray-50"
                    >
                      Siguiente
                    </button>
                  </div>
                )}
              </>
            )}

            <div className="mt-4 flex justify-end">
              <button
                onClick={handleNewSync}
                className="text-sm text-purple-600 hover:text-purple-800 font-medium transition-colors"
              >
                Nueva sincronización
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
