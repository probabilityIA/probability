'use client';

import { useState, useTransition } from 'react';
import { useRouter } from 'next/navigation';
import { AccountingReport } from '../../domain/types';
import { syncAccountingAction } from '../../infra/actions';
import { formatCOP, kindLabel } from '../format';
import { useToast } from '@/shared/providers/toast-provider';

interface ReportViewProps {
    report: AccountingReport | null;
    error: string | null;
    from: string;
    to: string;
}

export function ReportView({ report, error, from, to }: ReportViewProps) {
    const router = useRouter();
    const { showToast } = useToast();
    const [fromValue, setFromValue] = useState(from);
    const [toValue, setToValue] = useState(to);
    const [syncing, setSyncing] = useState(false);
    const [isPending, startTransition] = useTransition();

    const applyRange = () => {
        if (!fromValue || !toValue) return;
        startTransition(() => {
            router.replace(`/accounting?from=${fromValue}&to=${toValue}`);
        });
    };

    const handleSync = async () => {
        setSyncing(true);
        const result = await syncAccountingAction();
        setSyncing(false);
        if (result.success) {
            showToast(`Sincronizacion completada: ${result.data.created} movimientos creados`, 'success');
            router.refresh();
        } else {
            showToast(result.error, 'error');
        }
    };

    const totals = report?.totals;
    const cards = [
        { label: 'Ganancia real', value: totals?.real_income ?? 0, accent: 'text-emerald-600 dark:text-emerald-400' },
        { label: 'Entradas de caja', value: totals?.cash_in ?? 0, accent: 'text-blue-600 dark:text-blue-400' },
        { label: 'Gastos', value: totals?.real_expense ?? 0, accent: 'text-red-600 dark:text-red-400' },
        { label: 'Salidas de caja', value: totals?.cash_out ?? 0, accent: 'text-orange-600 dark:text-orange-400' },
        { label: 'Impuestos', value: totals?.tax_total ?? 0, accent: 'text-amber-600 dark:text-amber-400' },
        { label: 'Neto', value: totals?.net ?? 0, accent: (totals?.net ?? 0) >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400' },
    ];

    const byConcept = report?.by_concept || [];
    const byTax = report?.by_tax || [];

    return (
        <div className="space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Reporte contable</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        Contabilidad interna de la plataforma
                    </p>
                </div>
                <div className="flex flex-wrap items-end gap-2">
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
                    <button
                        onClick={applyRange}
                        disabled={isPending}
                        className="px-4 py-2 text-sm font-semibold rounded-lg bg-gray-900 dark:bg-gray-700 text-white hover:bg-gray-800 dark:hover:bg-gray-600 transition-colors disabled:opacity-50"
                    >
                        {isPending ? 'Cargando...' : 'Aplicar'}
                    </button>
                    <button
                        onClick={handleSync}
                        disabled={syncing}
                        className="inline-flex items-center gap-2 px-4 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors disabled:opacity-50"
                    >
                        <svg className={`w-4 h-4 ${syncing ? 'animate-spin' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                        {syncing ? 'Sincronizando...' : 'Sincronizar ahora'}
                    </button>
                </div>
            </div>

            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 text-sm rounded-lg px-4 py-3">
                    {error}
                </div>
            )}

            <div className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-6 gap-3">
                {cards.map((card) => (
                    <div key={card.label} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm p-4">
                        <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">{card.label}</p>
                        <p className={`mt-2 text-lg font-bold ${card.accent}`}>{formatCOP(card.value)}</p>
                    </div>
                ))}
            </div>

            <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
                <div className="xl:col-span-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm overflow-hidden">
                    <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
                        <h2 className="text-sm font-semibold text-gray-900 dark:text-white">Por concepto</h2>
                    </div>
                    <div className="overflow-x-auto">
                        <table className="w-full text-sm">
                            <thead>
                                <tr className="text-left text-xs text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-900/40">
                                    <th className="px-4 py-2 font-medium">Concepto</th>
                                    <th className="px-4 py-2 font-medium">Tipo</th>
                                    <th className="px-4 py-2 font-medium text-right">Movs.</th>
                                    <th className="px-4 py-2 font-medium text-right">Monto</th>
                                    <th className="px-4 py-2 font-medium text-right">Impuestos</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                {byConcept.length === 0 && (
                                    <tr>
                                        <td colSpan={5} className="px-4 py-6 text-center text-gray-500 dark:text-gray-400">
                                            Sin movimientos en el periodo
                                        </td>
                                    </tr>
                                )}
                                {byConcept.map((row) => (
                                    <tr key={row.concept_id} className="text-gray-800 dark:text-gray-200">
                                        <td className="px-4 py-2">
                                            <span className="font-medium">{row.name}</span>
                                            <span className="ml-2 text-xs text-gray-400">{row.code}</span>
                                            {row.kind === 'INCOME' && !row.is_real_income && (
                                                <span className="ml-2 text-[10px] font-semibold text-blue-600 dark:text-blue-400 uppercase">Solo caja</span>
                                            )}
                                        </td>
                                        <td className="px-4 py-2">
                                            <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold ${row.kind === 'INCOME' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300'}`}>
                                                {kindLabel(row.kind)}
                                            </span>
                                        </td>
                                        <td className="px-4 py-2 text-right">{row.entries_count}</td>
                                        <td className="px-4 py-2 text-right font-medium">{formatCOP(row.amount)}</td>
                                        <td className="px-4 py-2 text-right">{formatCOP(row.tax_total)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm overflow-hidden">
                    <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
                        <h2 className="text-sm font-semibold text-gray-900 dark:text-white">Por impuesto</h2>
                    </div>
                    <div className="overflow-x-auto">
                        <table className="w-full text-sm">
                            <thead>
                                <tr className="text-left text-xs text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-900/40">
                                    <th className="px-4 py-2 font-medium">Impuesto</th>
                                    <th className="px-4 py-2 font-medium text-right">Monto</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                {byTax.length === 0 && (
                                    <tr>
                                        <td colSpan={2} className="px-4 py-6 text-center text-gray-500 dark:text-gray-400">
                                            Sin impuestos en el periodo
                                        </td>
                                    </tr>
                                )}
                                {byTax.map((row) => (
                                    <tr key={row.tax_code} className="text-gray-800 dark:text-gray-200">
                                        <td className="px-4 py-2">
                                            <span className="font-medium">{row.tax_name}</span>
                                            <span className="ml-2 text-xs text-gray-400">{row.tax_code}</span>
                                        </td>
                                        <td className="px-4 py-2 text-right font-medium">{formatCOP(row.amount)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    );
}
