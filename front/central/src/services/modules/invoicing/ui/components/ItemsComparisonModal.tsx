'use client';

import { XMarkIcon } from '@heroicons/react/24/outline';
import type { ItemCompareResponseData } from '../../domain/types';

interface ItemsComparisonModalProps {
  isOpen: boolean;
  onClose: () => void;
  loading: boolean;
  data: ItemCompareResponseData | null;
  onRequestComparison: () => void;
}

const STATUS_CONFIG = {
  matched: {
    label: 'Coincide',
    icon: '✅',
    rowClass: 'bg-green-50',
    badgeClass: 'bg-green-100 text-green-800',
  },
  provider_only: {
    label: 'Solo en proveedor',
    icon: '🔴',
    rowClass: 'bg-red-50',
    badgeClass: 'bg-red-100 text-red-800',
  },
  system_only: {
    label: 'Solo en sistema',
    icon: '🟡',
    rowClass: 'bg-yellow-50',
    badgeClass: 'bg-yellow-100 text-yellow-800',
  },
} as const;

function formatCurrency(value: number | undefined): string {
  if (value === undefined || value === null) return '—';
  if (value === 0) return '$0';
  return new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(value);
}

export function ItemsComparisonModal({
  isOpen,
  onClose,
  loading,
  data,
  onRequestComparison,
}: ItemsComparisonModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />

      {/* Modal */}
      <div className="relative bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-5xl flex flex-col max-h-[92vh]">
        {/* Header */}
        <div className="flex items-center justify-between px-6 pt-6 pb-5 border-b border-gray-100 flex-shrink-0">
          <div>
            <h2 className="text-xl font-bold text-gray-900">Comparar Productos</h2>
            <p className="text-sm text-gray-500 mt-0.5">
              Compara los productos del sistema contra los registrados en el proveedor de facturación
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>

        {/* Start button */}
        {!data && !loading && (
          <div className="flex flex-col items-center justify-center py-20 px-6">
            <div className="text-6xl mb-4">📦</div>
            <p className="text-gray-600 text-sm mb-6 text-center max-w-md">
              Se compararán todos los productos de tu sistema contra el catálogo del proveedor de facturación, cruzando por código SKU.
            </p>
            <button
              onClick={onRequestComparison}
              className="px-6 py-2.5 bg-purple-600 hover:bg-purple-700 text-white rounded-lg text-sm font-semibold transition-colors"
            >
              Iniciar comparación
            </button>
          </div>
        )}

        {/* Spinner */}
        {loading && !data && (
          <div className="flex flex-col items-center justify-center py-20 text-gray-500 flex-shrink-0">
            <div className="w-12 h-12 border-4 border-purple-200 border-t-purple-600 rounded-full animate-spin mb-4" />
            <p className="text-sm font-medium">Obteniendo catálogo del proveedor…</p>
            <p className="text-xs mt-1 text-gray-400">Esto puede tardar unos segundos</p>
          </div>
        )}

        {/* Results */}
        {data && (
          <div className="flex-1 overflow-y-auto px-6 py-5 min-h-0">
            {/* Summary badges */}
            <div className="flex gap-3 mb-5 flex-wrap items-center">
              <span className="flex items-center gap-2 px-4 py-2 bg-green-50 border border-green-200 rounded-full text-sm font-medium text-green-800">
                ✅ {data.summary.matched} coinciden
              </span>
              <span className="flex items-center gap-2 px-4 py-2 bg-red-50 border border-red-200 rounded-full text-sm font-medium text-red-800">
                🔴 {data.summary.provider_only} solo en proveedor
              </span>
              <span className="flex items-center gap-2 px-4 py-2 bg-yellow-50 border border-yellow-200 rounded-full text-sm font-medium text-yellow-800">
                🟡 {data.summary.system_only} solo en sistema
              </span>
              <span className="ml-auto text-xs text-gray-400">
                Proveedor: {data.summary.total_provider} · Sistema: {data.summary.total_system}
              </span>
            </div>

            {data.results.length === 0 ? (
              <div className="text-center py-14 text-gray-400 text-sm">
                No se encontraron productos para comparar
              </div>
            ) : (
              <div className="overflow-x-auto rounded-xl border border-gray-200 dark:border-gray-700">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="bg-gradient-to-r from-purple-600 to-purple-700 text-white text-xs uppercase tracking-wide">
                      <th className="px-4 py-3 text-left rounded-tl-xl">Código</th>
                      <th className="px-4 py-3 text-left">Nombre Proveedor</th>
                      <th className="px-4 py-3 text-left">Nombre Sistema</th>
                      <th className="px-4 py-3 text-right">Precio Proveedor</th>
                      <th className="px-4 py-3 text-right">Precio Sistema</th>
                      <th className="px-4 py-3 text-right">Diferencia</th>
                      <th className="px-4 py-3 text-center rounded-tr-xl">Estado</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.results.map((row, idx) => {
                      const cfg = STATUS_CONFIG[row.status as keyof typeof STATUS_CONFIG] || STATUS_CONFIG.matched;
                      return (
                        <tr
                          key={idx}
                          className={`border-b ${cfg.rowClass} transition-all hover:brightness-95`}
                        >
                          <td className="px-4 py-3 font-mono font-medium text-gray-800">
                            {row.item_code || '—'}
                          </td>
                          <td className="px-4 py-3 text-gray-700">
                            {row.provider_name || '—'}
                          </td>
                          <td className="px-4 py-3 text-gray-700">
                            {row.system_name || '—'}
                          </td>
                          <td className="px-4 py-3 text-right font-semibold">
                            {row.provider_price ? formatCurrency(row.provider_price) : '—'}
                          </td>
                          <td className="px-4 py-3 text-right font-semibold">
                            {row.system_price ? formatCurrency(row.system_price) : '—'}
                          </td>
                          <td className={`px-4 py-3 text-right font-semibold ${
                            row.price_diff > 0 ? 'text-red-600' : row.price_diff < 0 ? 'text-blue-600' : 'text-gray-400'
                          }`}>
                            {row.status === 'matched' ? (
                              row.price_diff !== 0 ? formatCurrency(row.price_diff) : '—'
                            ) : '—'}
                          </td>
                          <td className="px-4 py-3 text-center">
                            <span
                              className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium ${cfg.badgeClass}`}
                            >
                              {cfg.icon} {cfg.label}
                            </span>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}

            {/* Nueva comparación */}
            <div className="mt-4 flex justify-end">
              <button
                onClick={onRequestComparison}
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
