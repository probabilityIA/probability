'use client';

import { useCallback, useEffect, useState } from 'react';
import {
    BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid, Cell, Legend,
} from 'recharts';
import { TrendingUp, Clock, Percent, Wallet, RefreshCw, AlertCircle, Truck } from 'lucide-react';
import { getCodSummaryAction } from '../../infra/actions';
import { CodSummary, ReportFilters } from '../../domain/types';
import { formatMoney, formatMoneyShort, carrierLabel, carrierColor } from './helpers';

interface Props {
    filters: ReportFilters;
}

function Kpi({ icon, label, value, sub, tone }: {
    icon: React.ReactNode; label: string; value: string; sub?: string; tone: string;
}) {
    return (
        <div className={`rounded-xl border p-4 ${tone}`}>
            <div className="flex items-center gap-1.5 text-[11px] uppercase font-bold opacity-70">
                {icon} {label}
            </div>
            <div className="text-2xl font-extrabold mt-1.5">{value}</div>
            {sub && <div className="text-xs mt-0.5 opacity-70">{sub}</div>}
        </div>
    );
}

export default function CodSummaryTab({ filters }: Props) {
    const [data, setData] = useState<CodSummary | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        setError(null);
        const res = await getCodSummaryAction(filters);
        if (res.success && res.data) {
            setData(res.data);
        } else {
            setError((res as any).message || 'Error al cargar el resumen');
            setData(null);
        }
        setLoading(false);
    }, [filters]);

    useEffect(() => { load(); }, [load]);

    if (loading) {
        return (
            <div className="flex items-center justify-center py-20 text-gray-400">
                <RefreshCw size={20} className="animate-spin mr-2" /> Cargando recaudo...
            </div>
        );
    }
    if (error || !data) {
        return (
            <div className="m-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm flex items-center gap-2">
                <AlertCircle size={16} /> {error || 'Sin datos'}
            </div>
        );
    }

    const monthlyData = data.monthly.map(m => ({
        label: m.label,
        Recaudado: Math.round(m.collected),
        Neto: Math.round(m.net),
    }));
    const carrierData = data.by_carrier.map(c => ({
        label: carrierLabel(c.carrier),
        Recaudado: Math.round(c.total_collected),
    }));

    return (
        <div className="space-y-4 p-1">
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
                <Kpi
                    icon={<TrendingUp size={12} />}
                    label="Recaudado"
                    value={formatMoney(data.total_collected)}
                    sub={`${data.orders_collected} ordenes entregadas`}
                    tone="bg-emerald-50 border-emerald-200 text-emerald-800 dark:bg-emerald-900/20 dark:border-emerald-800 dark:text-emerald-200"
                />
                <Kpi
                    icon={<Clock size={12} />}
                    label="Pendiente"
                    value={formatMoney(data.total_pending)}
                    sub={`${data.orders_pending} ordenes en transito`}
                    tone="bg-amber-50 border-amber-200 text-amber-800 dark:bg-amber-900/20 dark:border-amber-800 dark:text-amber-200"
                />
                <Kpi
                    icon={<Percent size={12} />}
                    label="Descuento transportadoras"
                    value={formatMoney(data.total_discount)}
                    sub="Comision sobre el recaudo"
                    tone="bg-red-50 border-red-200 text-red-800 dark:bg-red-900/20 dark:border-red-800 dark:text-red-200"
                />
                <Kpi
                    icon={<Wallet size={12} />}
                    label="Neto a recibir"
                    value={formatMoney(data.total_net)}
                    sub="Recaudado menos descuentos"
                    tone="bg-purple-50 border-purple-200 text-purple-800 dark:bg-purple-900/20 dark:border-purple-800 dark:text-purple-200"
                />
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-4">
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-3">
                        Historico de recaudo (ultimos meses)
                    </h3>
                    {monthlyData.length === 0 ? (
                        <div className="text-xs text-gray-400 py-10 text-center">Sin datos historicos.</div>
                    ) : (
                        <ResponsiveContainer width="100%" height={240}>
                            <BarChart data={monthlyData}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                                <XAxis dataKey="label" tick={{ fontSize: 11 }} />
                                <YAxis tick={{ fontSize: 11 }} tickFormatter={v => formatMoneyShort(v)} width={44} />
                                <Tooltip formatter={(v: any) => formatMoney(Number(v))} />
                                <Legend wrapperStyle={{ fontSize: 12 }} />
                                <Bar dataKey="Recaudado" fill="#10b981" radius={[4, 4, 0, 0]} />
                                <Bar dataKey="Neto" fill="#7c3aed" radius={[4, 4, 0, 0]} />
                            </BarChart>
                        </ResponsiveContainer>
                    )}
                </div>

                <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-4">
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-3">
                        Recaudo por transportadora
                    </h3>
                    {carrierData.length === 0 ? (
                        <div className="text-xs text-gray-400 py-10 text-center">Sin recaudo en el periodo.</div>
                    ) : (
                        <ResponsiveContainer width="100%" height={240}>
                            <BarChart data={carrierData} layout="vertical">
                                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                                <XAxis type="number" tick={{ fontSize: 11 }} tickFormatter={v => formatMoneyShort(v)} />
                                <YAxis type="category" dataKey="label" tick={{ fontSize: 11 }} width={92} />
                                <Tooltip formatter={(v: any) => formatMoney(Number(v))} />
                                <Bar dataKey="Recaudado" radius={[0, 4, 4, 0]}>
                                    {carrierData.map((_, i) => (
                                        <Cell key={i} fill={carrierColor(i)} />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    )}
                </div>
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
                <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center gap-2">
                    <Truck size={15} className="text-purple-600" />
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Detalle por transportadora</h3>
                </div>
                {data.by_carrier.length === 0 ? (
                    <div className="text-xs text-gray-400 py-8 text-center">Sin transportadoras con recaudo.</div>
                ) : (
                    <div className="overflow-x-auto">
                        <table className="w-full text-sm">
                            <thead>
                                <tr className="text-[11px] uppercase text-gray-500 dark:text-gray-400 border-b border-gray-100 dark:border-gray-700">
                                    <th className="text-left px-4 py-2 font-semibold">Transportadora</th>
                                    <th className="text-right px-4 py-2 font-semibold">Ordenes</th>
                                    <th className="text-right px-4 py-2 font-semibold">Recaudado</th>
                                    <th className="text-right px-4 py-2 font-semibold">% Desc.</th>
                                    <th className="text-right px-4 py-2 font-semibold">Descuento</th>
                                    <th className="text-right px-4 py-2 font-semibold">Neto</th>
                                </tr>
                            </thead>
                            <tbody>
                                {data.by_carrier.map((c, i) => (
                                    <tr key={c.carrier} className="border-b border-gray-50 dark:border-gray-700/50">
                                        <td className="px-4 py-2">
                                            <span className="inline-flex items-center gap-2 font-medium text-gray-900 dark:text-white">
                                                <span className="w-2.5 h-2.5 rounded-full" style={{ background: carrierColor(i) }} />
                                                {carrierLabel(c.carrier)}
                                            </span>
                                        </td>
                                        <td className="px-4 py-2 text-right text-gray-600 dark:text-gray-300">{c.orders_count}</td>
                                        <td className="px-4 py-2 text-right font-semibold text-emerald-600 dark:text-emerald-400">{formatMoney(c.total_collected)}</td>
                                        <td className="px-4 py-2 text-right text-gray-600 dark:text-gray-300">{c.discount_pct.toFixed(1)}%</td>
                                        <td className="px-4 py-2 text-right text-red-600 dark:text-red-400">{formatMoney(c.total_discount)}</td>
                                        <td className="px-4 py-2 text-right font-semibold text-purple-700 dark:text-purple-300">{formatMoney(c.total_net)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}
