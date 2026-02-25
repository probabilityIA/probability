'use client';

import { useState, useCallback } from 'react';
import { XMarkIcon, ChevronDownIcon, ChevronRightIcon } from '@heroicons/react/24/outline';
import { DateRangePicker } from '@/shared/ui/date-range-picker';
import type { CompareResponseData, CompareItemDetail } from '../../domain/types';

interface InvoiceComparisonModalProps {
  isOpen: boolean;
  onClose: () => void;
  loading: boolean;
  compareData: CompareResponseData | null;
  onRequestComparison: (dateFrom: string, dateTo: string) => void;
}

const STATUS_CONFIG = {
  matched: {
    label: 'Coincide',
    icon: '✅',
    rowClass: 'bg-green-50 border-green-100',
    badgeClass: 'bg-green-100 text-green-800',
  },
  provider_only: {
    label: 'Solo en proveedor',
    icon: '🔴',
    rowClass: 'bg-red-50 border-red-100',
    badgeClass: 'bg-red-100 text-red-800',
  },
  system_only: {
    label: 'Solo en sistema',
    icon: '🟡',
    rowClass: 'bg-yellow-50 border-yellow-100',
    badgeClass: 'bg-yellow-100 text-yellow-800',
  },
} as const;

function formatCurrency(value: string | number | undefined): string {
  if (value === undefined || value === null || value === '') return '—';
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(num)) return String(value);
  return new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(num);
}

function ItemsTable({ items, emptyText }: { items?: CompareItemDetail[]; emptyText: string }) {
  if (!items || items.length === 0) {
    return <p className="text-xs text-gray-400 italic">{emptyText}</p>;
  }
  return (
    <table className="w-full text-xs border border-gray-200 rounded">
      <thead className="bg-gray-100">
        <tr>
          <th className="px-2 py-1 text-left font-medium text-gray-600">Código</th>
          <th className="px-2 py-1 text-left font-medium text-gray-600">Nombre</th>
          <th className="px-2 py-1 text-center font-medium text-gray-600">Cant.</th>
          <th className="px-2 py-1 text-right font-medium text-gray-600">Precio Unit.</th>
          <th className="px-2 py-1 text-center font-medium text-gray-600">IVA%</th>
        </tr>
      </thead>
      <tbody>
        {items.map((item, i) => (
          <tr key={i} className="border-t border-gray-100">
            <td className="px-2 py-1 font-mono text-gray-700">{item.item_code || '—'}</td>
            <td className="px-2 py-1 text-gray-700">{item.item_name}</td>
            <td className="px-2 py-1 text-center text-gray-700">{item.quantity}</td>
            <td className="px-2 py-1 text-right text-gray-700">{formatCurrency(item.unit_value)}</td>
            <td className="px-2 py-1 text-center text-gray-700">{item.iva}%</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export function InvoiceComparisonModal({
  isOpen,
  onClose,
  loading,
  compareData,
  onRequestComparison,
}: InvoiceComparisonModalProps) {
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');
  const [validationError, setValidationError] = useState('');
  const [expandedRow, setExpandedRow] = useState<number | null>(null);

  const handleRangeChange = useCallback((from: string | undefined, to: string | undefined) => {
    setDateFrom(from ?? '');
    setDateTo(to ?? '');
    setValidationError('');
  }, []);

  const toggleRow = useCallback((idx: number) => {
    setExpandedRow(prev => (prev === idx ? null : idx));
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
    if (diffDays > 7) {
      setValidationError('El rango máximo es de 7 días');
      return;
    }
    onRequestComparison(dateFrom, dateTo);
  }, [dateFrom, dateTo, onRequestComparison]);

  const handleNewComparison = () => {
    setDateFrom('');
    setDateTo('');
    setValidationError('');
    setExpandedRow(null);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />

      {/* Modal */}
      <div className="relative bg-white rounded-2xl shadow-2xl w-full max-w-6xl flex flex-col max-h-[92vh]">
        {/* Header */}
        <div className="flex items-center justify-between px-6 pt-6 pb-5 border-b border-gray-100 flex-shrink-0">
          <div>
            <h2 className="text-xl font-bold text-gray-900">Auditoría Comparativa de Facturas</h2>
            <p className="text-sm text-gray-500 mt-0.5">
              Compara las facturas del sistema contra las registradas en el proveedor
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>

        {/* ── Form section (overflow visible so the calendar popup floats freely) ── */}
        {!compareData && !loading && (
          <div className="px-6 pt-5 pb-2 flex-shrink-0">
            <div className="flex items-end gap-4">
              {/* DateRangePicker ocupa todo el espacio disponible */}
              <div className="flex-1">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Rango de fechas{' '}
                  <span className="text-gray-400 font-normal">(máximo 7 días)</span>
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
                Iniciar comparación
              </button>
            </div>

            {validationError && (
              <p className="mt-2 text-sm text-red-600">{validationError}</p>
            )}

            {/* Spacer so the calendar popup has room */}
            <div className="h-80" />
          </div>
        )}

        {/* ── Spinner ── */}
        {loading && !compareData && (
          <div className="flex flex-col items-center justify-center py-20 text-gray-500 flex-shrink-0">
            <div className="w-12 h-12 border-4 border-purple-200 border-t-purple-600 rounded-full animate-spin mb-4" />
            <p className="text-sm font-medium">Obteniendo datos del proveedor…</p>
            <p className="text-xs mt-1 text-gray-400">Esto puede tardar unos segundos</p>
          </div>
        )}

        {/* ── Results (scrollable) ── */}
        {compareData && (
          <div className="flex-1 overflow-y-auto px-6 py-5 min-h-0">
            {/* Summary badges */}
            <div className="flex gap-3 mb-5 flex-wrap items-center">
              <span className="flex items-center gap-2 px-4 py-2 bg-green-50 border border-green-200 rounded-full text-sm font-medium text-green-800">
                ✅ {compareData.summary.matched} coinciden
              </span>
              <span className="flex items-center gap-2 px-4 py-2 bg-red-50 border border-red-200 rounded-full text-sm font-medium text-red-800">
                🔴 {compareData.summary.provider_only} solo en proveedor
              </span>
              <span className="flex items-center gap-2 px-4 py-2 bg-yellow-50 border border-yellow-200 rounded-full text-sm font-medium text-yellow-800">
                🟡 {compareData.summary.system_only} solo en sistema
              </span>
              <span className="ml-auto text-xs text-gray-400">
                {compareData.date_from} → {compareData.date_to}
              </span>
            </div>

            {compareData.results.length === 0 ? (
              <div className="text-center py-14 text-gray-400 text-sm">
                No se encontraron facturas en el rango seleccionado
              </div>
            ) : (
              <div className="overflow-x-auto rounded-xl border border-gray-200">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="bg-gradient-to-r from-purple-600 to-purple-700 text-white text-xs uppercase tracking-wide">
                      <th className="px-4 py-3 text-left rounded-tl-xl">Nro. Factura</th>
                      <th className="px-4 py-3 text-left">Fecha Factura</th>
                      <th className="px-4 py-3 text-left">Fecha Orden</th>
                      <th className="px-4 py-3 text-right">Total Proveedor</th>
                      <th className="px-4 py-3 text-right">Total Sistema</th>
                      <th className="px-4 py-3 text-center">Estado</th>
                      <th className="px-4 py-3 text-center rounded-tr-xl w-8">▶</th>
                    </tr>
                  </thead>
                  <tbody>
                    {compareData.results.map((row, idx) => {
                      const cfg = STATUS_CONFIG[row.status];
                      const isExpanded = expandedRow === idx;
                      return (
                        <>
                          <tr
                            key={idx}
                            className={`border-b ${cfg.rowClass} cursor-pointer hover:brightness-95 transition-all`}
                            onClick={() => toggleRow(idx)}
                          >
                            <td className="px-4 py-3 font-mono font-medium">
                              {row.prefix ? `${row.prefix}-` : ''}{row.invoice_number || '—'}
                            </td>
                            <td className="px-4 py-3 text-gray-600">{row.document_date || '—'}</td>
                            <td className="px-4 py-3 text-gray-600">{row.order_created_at || '—'}</td>
                            <td className="px-4 py-3 text-right font-semibold">
                              {row.provider_total ? formatCurrency(row.provider_total) : '—'}
                            </td>
                            <td className="px-4 py-3 text-right font-semibold">
                              {row.system_total !== undefined ? formatCurrency(row.system_total) : '—'}
                            </td>
                            <td className="px-4 py-3 text-center">
                              <span
                                className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium ${cfg.badgeClass}`}
                              >
                                {cfg.icon} {cfg.label}
                              </span>
                            </td>
                            <td className="px-4 py-3 text-center text-gray-400">
                              {isExpanded ? (
                                <ChevronDownIcon className="w-4 h-4 inline" />
                              ) : (
                                <ChevronRightIcon className="w-4 h-4 inline" />
                              )}
                            </td>
                          </tr>
                          {isExpanded && (
                            <tr key={`${idx}-detail`}>
                              <td colSpan={7} className="px-4 py-4 bg-gray-50 border-b border-gray-200">
                                <div className="grid grid-cols-2 gap-6">
                                  {/* Sistema */}
                                  <div>
                                    <h4 className="text-xs font-bold text-gray-600 uppercase mb-2">
                                      📋 Sistema
                                    </h4>
                                    <ItemsTable
                                      items={row.system_items}
                                      emptyText="Sin ítems registrados"
                                    />
                                  </div>
                                  {/* Softpymes */}
                                  <div>
                                    <h4 className="text-xs font-bold text-gray-600 uppercase mb-2">
                                      🔗 Softpymes
                                    </h4>
                                    <ItemsTable
                                      items={row.provider_details}
                                      emptyText="No encontrado en proveedor"
                                    />
                                  </div>
                                </div>
                                {row.customer_name && (
                                  <p className="mt-3 text-xs text-gray-500">
                                    Cliente:{' '}
                                    <span className="font-medium">{row.customer_name}</span>
                                    {row.customer_nit && ` · NIT: ${row.customer_nit}`}
                                  </p>
                                )}
                              </td>
                            </tr>
                          )}
                        </>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}

            {/* Nueva comparación */}
            <div className="mt-4 flex justify-end">
              <button
                onClick={handleNewComparison}
                className="text-sm text-purple-600 hover:text-purple-800 font-medium transition-colors"
              >
                Nueva comparación
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
