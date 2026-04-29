'use client';

import { useEffect, useState, useCallback } from 'react';
import { Spinner, Alert } from '@/shared/ui';
import { ProfitReportResponse } from '../../domain/types';
import { shippingProfitReportAction } from '../../infra/actions';
import { getActionError } from '@/shared/utils/action-result';
import ShippingProfitDetailModal from './ShippingProfitDetailModal';

interface Props {
    selectedBusinessId?: number;
}

const fmt = (n: number) => '$ ' + Math.round(n).toLocaleString('es-CO');

const toISO = (d: Date) => {
    const y = d.getFullYear();
    const m = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    return `${y}-${m}-${day}`;
};

type PresetKey = 'today' | 'week' | 'month' | '3months' | 'custom';

function rangeFor(preset: PresetKey): { from: string; to: string } {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const to = toISO(today);
    if (preset === 'today') return { from: to, to };
    if (preset === 'week') {
        const d = new Date(today);
        d.setDate(d.getDate() - 6);
        return { from: toISO(d), to };
    }
    if (preset === 'month') {
        const d = new Date(today.getFullYear(), today.getMonth(), 1);
        return { from: toISO(d), to };
    }
    if (preset === '3months') {
        const d = new Date(today);
        d.setMonth(d.getMonth() - 3);
        return { from: toISO(d), to };
    }
    return { from: to, to };
}

export default function ShippingProfitReport({ selectedBusinessId }: Props) {
    const initial = rangeFor('today');
    const [from, setFrom] = useState<string>(initial.from);
    const [to, setTo] = useState<string>(initial.to);
    const [preset, setPreset] = useState<PresetKey>('today');
    const [data, setData] = useState<ProfitReportResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [detailCarrier, setDetailCarrier] = useState<{ code: string; label: string } | null>(null);

    const applyPreset = (p: PresetKey) => {
        setPreset(p);
        if (p === 'custom') return;
        const r = rangeFor(p);
        setFrom(r.from);
        setTo(r.to);
    };

    const fetchReport = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await shippingProfitReportAction({
                business_id: selectedBusinessId,
                from,
                to,
            });
            setData(r);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar el reporte'));
        } finally {
            setLoading(false);
        }
    }, [selectedBusinessId, from, to]);

    useEffect(() => {
        fetchReport();
    }, [fetchReport]);

    const totalRow = data?.totals;

    return (
        <div className="space-y-4">
            <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700 space-y-3">
                <div className="flex flex-wrap gap-2">
                    {([
                        ['today', 'Hoy'],
                        ['week', 'Semana'],
                        ['month', 'Mes'],
                        ['3months', 'Ultimos 3 meses'],
                        ['custom', 'Personalizado'],
                    ] as [PresetKey, string][]).map(([k, label]) => (
                        <button
                            key={k}
                            onClick={() => applyPreset(k)}
                            className={`px-3 py-1.5 text-sm rounded-md border transition-colors ${
                                preset === k
                                    ? 'bg-purple-600 border-purple-600 text-white'
                                    : 'bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600'
                            }`}
                        >
                            {label}
                        </button>
                    ))}
                </div>
                <div className="flex flex-wrap items-end gap-3">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Desde</label>
                        <input
                            type="date"
                            value={from}
                            onChange={(e) => { setFrom(e.target.value); setPreset('custom'); }}
                            className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-sm text-gray-900 dark:text-white"
                        />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Hasta</label>
                        <input
                            type="date"
                            value={to}
                            onChange={(e) => { setTo(e.target.value); setPreset('custom'); }}
                            className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-sm text-gray-900 dark:text-white"
                        />
                    </div>
                    <button
                        onClick={fetchReport}
                        className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white text-sm font-medium rounded-md"
                    >
                        Aplicar
                    </button>
                </div>
            </div>

            {error && <Alert type="error">{error}</Alert>}

            {totalRow && (
                <div className="grid grid-cols-1 sm:grid-cols-4 gap-3">
                    <SummaryCard label="Guias generadas" value={String(totalRow.shipments)} accent="text-gray-900 dark:text-white" />
                    <SummaryCard label="Cobrado al cliente" value={fmt(totalRow.customer_charge_total)} accent="text-blue-700 dark:text-blue-300" />
                    <SummaryCard label="Costo real carrier" value={fmt(totalRow.carrier_cost_total)} accent="text-orange-700 dark:text-orange-300" />
                    <SummaryCard label="Ganancia" value={fmt(totalRow.profit_total)} accent="text-emerald-700 dark:text-emerald-300" />
                </div>
            )}

            <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-x-auto">
                <table className="w-full text-sm">
                    <thead className="bg-gray-100 dark:bg-gray-700">
                        <tr>
                            <th className="text-left px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Transportadora</th>
                            <th className="text-center px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Guias</th>
                            <th className="text-right px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Cobrado al cliente</th>
                            <th className="text-right px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Costo real carrier</th>
                            <th className="text-right px-4 py-3 font-semibold text-emerald-700 dark:text-emerald-300">Ganancia</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                        {loading ? (
                            <tr>
                                <td colSpan={5} className="text-center py-10"><Spinner size="lg" /></td>
                            </tr>
                        ) : data?.rows && data.rows.length > 0 ? (
                            data.rows.map((r, i) => (
                                <tr
                                    key={`${r.carrier_code}-${i}`}
                                    className="hover:bg-gray-50 dark:hover:bg-gray-700/50 cursor-pointer"
                                    onClick={() => setDetailCarrier({ code: r.carrier_code, label: r.carrier })}
                                >
                                    <td className="px-4 py-3 text-gray-900 dark:text-white font-medium underline decoration-dotted">{r.carrier}</td>
                                    <td className="px-4 py-3 text-center text-gray-700 dark:text-gray-200">{r.shipments}</td>
                                    <td className="px-4 py-3 text-right text-blue-700 dark:text-blue-300">{fmt(r.customer_charge_total)}</td>
                                    <td className="px-4 py-3 text-right text-orange-700 dark:text-orange-300">{fmt(r.carrier_cost_total)}</td>
                                    <td className="px-4 py-3 text-right font-semibold text-emerald-700 dark:text-emerald-300">{fmt(r.profit_total)}</td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={5} className="text-center py-10 text-gray-400">Sin guias generadas en el rango seleccionado</td>
                            </tr>
                        )}
                    </tbody>
                    {totalRow && data?.rows && data.rows.length > 0 && (
                        <tfoot className="bg-gray-50 dark:bg-gray-900/40 font-semibold">
                            <tr>
                                <td className="px-4 py-3 text-gray-900 dark:text-white">Total</td>
                                <td className="px-4 py-3 text-center text-gray-900 dark:text-white">{totalRow.shipments}</td>
                                <td className="px-4 py-3 text-right text-blue-800 dark:text-blue-200">{fmt(totalRow.customer_charge_total)}</td>
                                <td className="px-4 py-3 text-right text-orange-800 dark:text-orange-200">{fmt(totalRow.carrier_cost_total)}</td>
                                <td className="px-4 py-3 text-right text-emerald-800 dark:text-emerald-200">{fmt(totalRow.profit_total)}</td>
                            </tr>
                        </tfoot>
                    )}
                </table>
            </div>

            {detailCarrier && (
                <ShippingProfitDetailModal
                    isOpen={!!detailCarrier}
                    onClose={() => setDetailCarrier(null)}
                    carrier={detailCarrier.code}
                    carrierLabel={detailCarrier.label}
                    from={from}
                    to={to}
                    selectedBusinessId={selectedBusinessId}
                />
            )}
        </div>
    );
}

function SummaryCard({ label, value, accent }: { label: string; value: string; accent: string }) {
    return (
        <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">{label}</p>
            <p className={`text-xl font-bold ${accent}`}>{value}</p>
        </div>
    );
}
