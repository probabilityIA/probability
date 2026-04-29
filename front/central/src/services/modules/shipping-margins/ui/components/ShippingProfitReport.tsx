'use client';

import { useEffect, useState, useCallback } from 'react';
import { Spinner, Alert } from '@/shared/ui';
import { ProfitReportResponse } from '../../domain/types';
import { shippingProfitReportAction } from '../../infra/actions';
import { getActionError } from '@/shared/utils/action-result';

interface Props {
    selectedBusinessId?: number;
}

const fmt = (n: number) => '$ ' + Math.round(n).toLocaleString('es-CO');

function defaultFrom() {
    const d = new Date();
    d.setDate(1);
    return d.toISOString().slice(0, 10);
}

function defaultTo() {
    return new Date().toISOString().slice(0, 10);
}

export default function ShippingProfitReport({ selectedBusinessId }: Props) {
    const [from, setFrom] = useState<string>(defaultFrom());
    const [to, setTo] = useState<string>(defaultTo());
    const [data, setData] = useState<ProfitReportResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

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
            <div className="flex flex-wrap items-end gap-3 bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
                <div>
                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Desde</label>
                    <input
                        type="date"
                        value={from}
                        onChange={(e) => setFrom(e.target.value)}
                        className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-sm text-gray-900 dark:text-white"
                    />
                </div>
                <div>
                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Hasta</label>
                    <input
                        type="date"
                        value={to}
                        onChange={(e) => setTo(e.target.value)}
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
                                <tr key={`${r.carrier_code}-${i}`} className="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                    <td className="px-4 py-3 text-gray-900 dark:text-white font-medium">{r.carrier}</td>
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
