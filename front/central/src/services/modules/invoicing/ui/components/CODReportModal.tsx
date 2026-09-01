'use client';

import { useState, useCallback } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { DateRangePicker } from '@/shared/ui/date-range-picker';
import { getCODReportAction } from '../../infra/actions';
import type { CODReport } from '../../domain/types';
import { CODReportInvoicesModal } from './CODReportInvoicesModal';

interface CODReportModalProps {
  isOpen: boolean;
  onClose: () => void;
  businessId: number | null;
}

function formatCurrency(value: number): string {
  return new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(value);
}

export function CODReportModal({ isOpen, onClose, businessId }: CODReportModalProps) {
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [report, setReport] = useState<CODReport | null>(null);
  const [drillDown, setDrillDown] = useState<{ accountNumber: string; isCOD: boolean } | null>(null);

  const handleRangeChange = useCallback((from: string | undefined, to: string | undefined) => {
    setDateFrom(from ?? '');
    setDateTo(to ?? '');
  }, []);

  const handleGenerate = useCallback(async () => {
    if (!businessId) {
      setError('Selecciona un negocio antes de generar el reporte');
      return;
    }
    setLoading(true);
    setError('');
    try {
      const data = await getCODReportAction(businessId, dateFrom || undefined, dateTo || undefined);
      setReport(data);
    } catch (err: any) {
      setError('Error al generar el reporte: ' + err.message);
    } finally {
      setLoading(false);
    }
  }, [businessId, dateFrom, dateTo]);

  const handleNewReport = () => {
    setReport(null);
    setDateFrom('');
    setDateTo('');
    setError('');
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />

      <div className="relative bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-4xl flex flex-col max-h-[92vh]">
        <div className="flex items-center justify-between px-6 pt-6 pb-5 border-b border-gray-100 flex-shrink-0">
          <div>
            <h2 className="text-xl font-bold text-gray-900 dark:text-white">Reporte de facturación contra entrega</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
              Cuántas facturas se emitieron y a qué cuenta bancaria se registró cada recibo de caja
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 dark:text-gray-300 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>

        {!report && !loading && (
          <div className="px-6 pt-5 pb-2 flex-shrink-0">
            <div className="flex items-end gap-4">
              <div className="flex-1">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1.5">
                  Rango de fechas{' '}
                  <span className="text-gray-400 font-normal">(vacío = todo el histórico)</span>
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
                onClick={handleGenerate}
                className="px-5 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg text-sm font-semibold transition-colors whitespace-nowrap h-[38px]"
              >
                Generar reporte
              </button>
            </div>
            {error && <p className="mt-2 text-sm text-red-600">{error}</p>}
            <div className="h-80" />
          </div>
        )}

        {loading && (
          <div className="flex flex-col items-center justify-center py-20 text-gray-500 dark:text-gray-400 flex-shrink-0">
            <div className="w-12 h-12 border-4 border-purple-200 border-t-purple-600 rounded-full animate-spin mb-4" />
            <p className="text-sm font-medium">Generando reporte...</p>
          </div>
        )}

        {report && !loading && (
          <div className="flex-1 overflow-y-auto px-6 py-5 min-h-0">
            <p className="text-xs text-gray-400 mb-4">Período: {report.period_label}</p>

            <div className="grid grid-cols-3 gap-3 mb-5">
              <div className="rounded-xl border border-gray-200 dark:border-gray-700 px-4 py-3">
                <p className="text-xs text-gray-500 dark:text-gray-400">Total facturas</p>
                <p className="text-xl font-bold text-gray-900 dark:text-white">{report.total_invoices}</p>
                <p className="text-xs text-gray-400 mt-0.5">{formatCurrency(report.total_amount)}</p>
              </div>
              <div className="rounded-xl border border-amber-200 bg-amber-50 dark:bg-amber-900/20 px-4 py-3">
                <p className="text-xs text-amber-700 dark:text-amber-300">Contra entrega</p>
                <p className="text-xl font-bold text-amber-800 dark:text-amber-200">{report.cod_count}</p>
                <p className="text-xs text-amber-600 dark:text-amber-400 mt-0.5">{formatCurrency(report.cod_amount)}</p>
              </div>
              <div className="rounded-xl border border-green-200 bg-green-50 dark:bg-green-900/20 px-4 py-3">
                <p className="text-xs text-green-700 dark:text-green-300">Pagadas por adelantado</p>
                <p className="text-xl font-bold text-green-800 dark:text-green-200">{report.non_cod_count}</p>
                <p className="text-xs text-green-600 dark:text-green-400 mt-0.5">{formatCurrency(report.non_cod_amount)}</p>
              </div>
            </div>

            <h3 className="text-sm font-bold text-gray-700 dark:text-gray-200 mb-2">Desglose por cuenta bancaria</h3>
            {report.by_account.length === 0 ? (
              <div className="text-center py-10 text-gray-400 text-sm">Sin facturas en el período seleccionado</div>
            ) : (
              <div className="overflow-x-auto rounded-xl border border-gray-200 dark:border-gray-700">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="bg-gradient-to-r from-purple-600 to-purple-700 text-white text-xs uppercase tracking-wide">
                      <th className="px-4 py-3 text-left rounded-tl-xl">Cuenta bancaria</th>
                      <th className="px-4 py-3 text-center">Tipo</th>
                      <th className="px-4 py-3 text-right">Facturas</th>
                      <th className="px-4 py-3 text-right rounded-tr-xl">Total</th>
                    </tr>
                  </thead>
                  <tbody>
                    {report.by_account.map((row, idx) => (
                      <tr
                        key={idx}
                        onClick={() => setDrillDown({ accountNumber: row.account_number, isCOD: row.is_cod })}
                        className="border-b border-gray-100 dark:border-gray-700 cursor-pointer hover:bg-purple-50 dark:hover:bg-purple-900/20 transition-colors"
                      >
                        <td className="px-4 py-3 font-mono font-medium text-gray-700 dark:text-gray-200">{row.account_number}</td>
                        <td className="px-4 py-3 text-center">
                          <span
                            className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium ${
                              row.is_cod
                                ? 'bg-amber-100 text-amber-800'
                                : 'bg-green-100 text-green-800'
                            }`}
                          >
                            {row.is_cod ? 'Contra entrega' : 'Contado'}
                          </span>
                        </td>
                        <td className="px-4 py-3 text-right text-gray-700 dark:text-gray-200">{row.count}</td>
                        <td className="px-4 py-3 text-right font-semibold text-gray-900 dark:text-white">{formatCurrency(row.amount)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

            <div className="mt-4 flex justify-end">
              <button
                onClick={handleNewReport}
                className="text-sm text-purple-600 hover:text-purple-800 font-medium transition-colors"
              >
                Nuevo reporte
              </button>
            </div>
          </div>
        )}
      </div>

      <CODReportInvoicesModal
        isOpen={drillDown !== null}
        onClose={() => setDrillDown(null)}
        businessId={businessId}
        accountNumber={drillDown?.accountNumber ?? ''}
        isCOD={drillDown?.isCOD ?? false}
        startDate={dateFrom || undefined}
        endDate={dateTo || undefined}
      />
    </div>
  );
}
