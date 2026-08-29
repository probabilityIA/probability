'use client';

import { useState, useTransition } from 'react';
import { useRouter } from 'next/navigation';
import { ConfirmModal } from '@/shared/ui';
import { Concept, ConceptKind, EntriesListResponse, Entry } from '../../domain/types';
import { deleteAccountingEntryAction } from '../../infra/actions';
import { formatCOP, formatEntryDate, kindLabel } from '../format';
import { useToast } from '@/shared/providers/toast-provider';
import { EntryFormModal } from './EntryFormModal';

interface EntriesFilters {
    page: number;
    page_size: number;
    from: string;
    to: string;
    concept_id: number | null;
    kind: ConceptKind | null;
}

interface EntriesViewProps {
    entries: EntriesListResponse;
    concepts: Concept[];
    filters: EntriesFilters;
    error: string | null;
}

export function EntriesView({ entries, concepts, filters, error }: EntriesViewProps) {
    const router = useRouter();
    const { showToast } = useToast();
    const [fromValue, setFromValue] = useState(filters.from);
    const [toValue, setToValue] = useState(filters.to);
    const [expandedId, setExpandedId] = useState<number | null>(null);
    const [showCreate, setShowCreate] = useState(false);
    const [entryToDelete, setEntryToDelete] = useState<Entry | null>(null);
    const [isPending, startTransition] = useTransition();

    const navigate = (next: Partial<EntriesFilters>) => {
        const merged = { ...filters, from: fromValue, to: toValue, ...next };
        const params = new URLSearchParams();
        params.set('page', String(merged.page));
        params.set('page_size', String(merged.page_size));
        if (merged.from) params.set('from', merged.from);
        if (merged.to) params.set('to', merged.to);
        if (merged.concept_id) params.set('concept_id', String(merged.concept_id));
        if (merged.kind) params.set('kind', merged.kind);
        startTransition(() => {
            router.replace(`/accounting/movimientos?${params.toString()}`);
        });
    };

    const handleDelete = async () => {
        if (!entryToDelete) return;
        const result = await deleteAccountingEntryAction(entryToDelete.id);
        if (result.success) {
            showToast('Movimiento eliminado', 'success');
            router.refresh();
        } else {
            showToast(result.error, 'error');
        }
        setEntryToDelete(null);
    };

    const totalPages = entries.total_pages || 0;

    return (
        <div className="space-y-4">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Movimientos</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        {entries.total} movimientos en el periodo
                    </p>
                </div>
                <button
                    onClick={() => setShowCreate(true)}
                    className="inline-flex items-center gap-2 px-4 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors"
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    Nuevo movimiento
                </button>
            </div>

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm p-4">
                <div className="flex flex-wrap items-end gap-3">
                    <div>
                        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Desde</label>
                        <input
                            type="date"
                            value={fromValue}
                            onChange={(e) => setFromValue(e.target.value)}
                            className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Hasta</label>
                        <input
                            type="date"
                            value={toValue}
                            onChange={(e) => setToValue(e.target.value)}
                            className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Concepto</label>
                        <select
                            value={filters.concept_id ?? 0}
                            onChange={(e) => navigate({ concept_id: Number(e.target.value) || null, page: 1 })}
                            className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        >
                            <option value={0}>Todos</option>
                            {concepts.map((c) => (
                                <option key={c.id} value={c.id}>{c.name}</option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Tipo</label>
                        <select
                            value={filters.kind ?? ''}
                            onChange={(e) => navigate({ kind: (e.target.value || null) as ConceptKind | null, page: 1 })}
                            className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        >
                            <option value="">Todos</option>
                            <option value="INCOME">Ingreso</option>
                            <option value="EXPENSE">Gasto</option>
                        </select>
                    </div>
                    <button
                        onClick={() => navigate({ page: 1 })}
                        disabled={isPending}
                        className="px-4 py-2 text-sm font-semibold rounded-lg bg-gray-900 dark:bg-gray-700 text-white hover:bg-gray-800 dark:hover:bg-gray-600 transition-colors disabled:opacity-50"
                    >
                        {isPending ? 'Cargando...' : 'Filtrar'}
                    </button>
                </div>
            </div>

            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 text-sm rounded-lg px-4 py-3">
                    {error}
                </div>
            )}

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-sm">
                        <thead>
                            <tr className="text-left text-xs text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-900/40">
                                <th className="px-4 py-2 font-medium">Fecha</th>
                                <th className="px-4 py-2 font-medium">Concepto</th>
                                <th className="px-4 py-2 font-medium">Negocio</th>
                                <th className="px-4 py-2 font-medium">Tipo</th>
                                <th className="px-4 py-2 font-medium">Origen</th>
                                <th className="px-4 py-2 font-medium text-right">Monto</th>
                                <th className="px-4 py-2 font-medium text-right">Impuestos</th>
                                <th className="px-4 py-2 font-medium text-right">Acciones</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                            {entries.data.length === 0 && (
                                <tr>
                                    <td colSpan={8} className="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
                                        No hay movimientos con los filtros seleccionados
                                    </td>
                                </tr>
                            )}
                            {entries.data.map((entry) => {
                                const hasTaxes = (entry.tax_detail?.length ?? 0) > 0;
                                const expanded = expandedId === entry.id;
                                return [
                                    <tr key={entry.id} className="text-gray-800 dark:text-gray-200">
                                        <td className="px-4 py-2 whitespace-nowrap">{formatEntryDate(entry.entry_date)}</td>
                                        <td className="px-4 py-2">
                                            <span className="font-medium">{entry.concept_name}</span>
                                            {entry.description && (
                                                <p className="text-xs text-gray-500 dark:text-gray-400 max-w-xs truncate" title={entry.description}>
                                                    {entry.description}
                                                </p>
                                            )}
                                        </td>
                                        <td className="px-4 py-2">{entry.business_name || '-'}</td>
                                        <td className="px-4 py-2">
                                            <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold ${entry.kind === 'INCOME' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300'}`}>
                                                {kindLabel(entry.kind)}
                                            </span>
                                        </td>
                                        <td className="px-4 py-2">
                                            <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold ${entry.is_automatic ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300' : 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300'}`}>
                                                {entry.is_automatic ? 'Automático' : 'Manual'}
                                            </span>
                                        </td>
                                        <td className="px-4 py-2 text-right font-medium whitespace-nowrap">{formatCOP(entry.amount)}</td>
                                        <td className="px-4 py-2 text-right whitespace-nowrap">
                                            {hasTaxes ? (
                                                <button
                                                    onClick={() => setExpandedId(expanded ? null : entry.id)}
                                                    className="inline-flex items-center gap-1 text-purple-600 dark:text-purple-400 hover:underline"
                                                >
                                                    {formatCOP(entry.tax_total)}
                                                    <svg className={`w-3.5 h-3.5 transition-transform ${expanded ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                                                    </svg>
                                                </button>
                                            ) : (
                                                <span className="text-gray-400">{formatCOP(entry.tax_total)}</span>
                                            )}
                                        </td>
                                        <td className="px-4 py-2 text-right">
                                            {!entry.is_automatic && (
                                                <button
                                                    onClick={() => setEntryToDelete(entry)}
                                                    className="p-1.5 rounded-lg text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 transition-colors"
                                                    title="Eliminar movimiento"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                                    </svg>
                                                </button>
                                            )}
                                        </td>
                                    </tr>,
                                    expanded && hasTaxes ? (
                                        <tr key={`${entry.id}-taxes`} className="bg-gray-50 dark:bg-gray-900/40">
                                            <td colSpan={8} className="px-6 py-3">
                                                <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase mb-2">Detalle de impuestos</p>
                                                <div className="flex flex-wrap gap-4">
                                                    {(entry.tax_detail || []).map((tax) => (
                                                        <div key={tax.tax_code} className="text-xs text-gray-700 dark:text-gray-300">
                                                            <span className="font-medium">{tax.tax_name}</span>
                                                            <span className="text-gray-400"> ({tax.rate_percent}%)</span>
                                                            <span className="ml-1 font-semibold">{formatCOP(tax.amount)}</span>
                                                        </div>
                                                    ))}
                                                </div>
                                            </td>
                                        </tr>
                                    ) : null,
                                ];
                            })}
                        </tbody>
                    </table>
                </div>

                <div className="flex flex-col sm:flex-row items-center justify-between gap-3 px-4 py-3 border-t border-gray-200 dark:border-gray-700">
                    <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
                        <span>Por página:</span>
                        <select
                            value={filters.page_size}
                            onChange={(e) => navigate({ page_size: Number(e.target.value), page: 1 })}
                            className="px-2 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                        >
                            {[10, 25, 50, 100].map((size) => (
                                <option key={size} value={size}>{size}</option>
                            ))}
                        </select>
                    </div>
                    <div className="flex items-center gap-2">
                        <button
                            onClick={() => navigate({ page: filters.page - 1 })}
                            disabled={filters.page <= 1 || isPending}
                            className="px-3 py-1.5 text-sm rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 transition-colors"
                        >
                            Anterior
                        </button>
                        <span className="text-sm text-gray-600 dark:text-gray-300">
                            Pagina {filters.page} de {Math.max(totalPages, 1)}
                        </span>
                        <button
                            onClick={() => navigate({ page: filters.page + 1 })}
                            disabled={filters.page >= totalPages || isPending}
                            className="px-3 py-1.5 text-sm rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 transition-colors"
                        >
                            Siguiente
                        </button>
                    </div>
                </div>
            </div>

            <EntryFormModal
                isOpen={showCreate}
                onClose={() => setShowCreate(false)}
                onCreated={() => {
                    setShowCreate(false);
                    router.refresh();
                }}
                concepts={concepts}
            />

            <ConfirmModal
                isOpen={entryToDelete !== null}
                onClose={() => setEntryToDelete(null)}
                onConfirm={handleDelete}
                title="Eliminar movimiento"
                message={`Se eliminara el movimiento manual de ${entryToDelete ? formatCOP(entryToDelete.amount) : ''} (${entryToDelete?.concept_name || ''}). Esta acción no se puede deshacer.`}
                confirmText="Eliminar"
                type="danger"
            />
        </div>
    );
}
