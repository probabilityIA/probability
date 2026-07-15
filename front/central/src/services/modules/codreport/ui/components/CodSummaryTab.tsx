'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { RefreshCw, AlertCircle } from 'lucide-react';
import { getCodSummaryAction } from '../../infra/actions';
import { CodSummary, ReportFilters } from '../../domain/types';
import { formatMoney, formatMoneyShort, carrierLabel } from './helpers';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';

interface Props {
    filters: ReportFilters;
}

const CARD = 'bg-white dark:bg-gray-800 border border-[#ececf0] dark:border-gray-700 rounded-2xl';
const DONUT_COLORS = ['#7c3aed', '#0891b2', '#e11d48', '#ea580c', '#16a34a', '#0ea5e9', '#db2777', '#f59e0b', '#4f46e5', '#14b8a6'];

const BUCKETS: { key: string; label: string; word: string }[] = [
    { key: 'day', label: 'Dias', word: 'dia' },
    { key: 'week', label: 'Semanas', word: 'semana' },
    { key: 'month', label: 'Meses', word: 'mes' },
    { key: 'quarter', label: 'Trimestres', word: 'trimestre' },
    { key: 'semester', label: 'Semestres', word: 'semestre' },
    { key: 'year', label: 'Anios', word: 'anio' },
];

function defaultBucket(range?: string): string {
    switch (range) {
        case '3months': return 'week';
        default: return 'day';
    }
}

interface KpiProps { accent: string; label: string; value: string; sub: string; }
function KpiCard({ accent, label, value, sub }: KpiProps) {
    return (
        <div className={`${CARD} relative overflow-hidden p-[18px]`}>
            <div className="absolute left-0 top-0 bottom-0 w-1" style={{ background: accent }} />
            <div className="flex items-center gap-2 text-[12px] font-bold uppercase tracking-wide" style={{ color: accent }}>
                <span className="w-2 h-2 rounded-full" style={{ background: accent }} />
                {label}
            </div>
            <div className="text-[27px] font-extrabold tracking-tight mt-2.5 text-gray-900 dark:text-white">{value}</div>
            <div className="text-[12.5px] text-gray-400 dark:text-gray-500 font-medium mt-1">{sub}</div>
        </div>
    );
}

export default function CodSummaryTab({ filters }: Props) {
    const [data, setData] = useState<CodSummary | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [userBucket, setUserBucket] = useState<string | null>(null);

    const bucket = userBucket ?? defaultBucket(filters.range);
    const bucketWord = BUCKETS.find(b => b.key === bucket)?.word ?? 'dia';

    const rangeKey = `${filters.range}|${filters.start_date ?? ''}|${filters.end_date ?? ''}`;
    const prevRangeKey = useRef(rangeKey);
    useEffect(() => {
        if (prevRangeKey.current !== rangeKey) {
            prevRangeKey.current = rangeKey;
            setUserBucket(null);
        }
    }, [rangeKey]);

    const load = useCallback(async () => {
        setLoading(true);
        setError(null);
        const res = await getCodSummaryAction(filters, bucket);
        if (res.success && res.data) setData(res.data);
        else { setError((res as any).message || 'Error al cargar el resumen'); setData(null); }
        setLoading(false);
    }, [filters, bucket]);

    useEffect(() => { load(); }, [load]);

    if (loading && !data) {
        return (
            <div className="flex items-center justify-center py-20 text-gray-400">
                <RefreshCw size={20} className="animate-spin mr-2" /> Cargando recaudo...
            </div>
        );
    }
    if (error && !data) {
        return (
            <div className="m-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm flex items-center gap-2">
                <AlertCircle size={16} /> {error || 'Sin datos'}
            </div>
        );
    }
    if (!data) return null;

    const currency = 'COP';
    const funnel = [
        { label: 'En curso', total: data.en_curso_total, orders: data.en_curso_orders, color: '#0ea5b7' },
        { label: 'Entregado', total: data.entregado_total, orders: data.entregado_orders, color: '#16a34a' },
        { label: 'Por pagar', total: data.total_pending, orders: data.orders_pending, color: '#d97706' },
        { label: 'Recaudado', total: data.total_collected, orders: data.orders_collected, color: '#7c3aed' },
    ];
    const funnelMax = Math.max(1, ...funnel.map(f => f.total));

    const detail = [...(data.carrier_detail || [])];
    const donutTotal = detail.reduce((a, c) => a + c.total, 0);
    let acc = 0;
    const segs: string[] = [];
    const donutCarriers = detail.map((c, i) => {
        const color = DONUT_COLORS[i % DONUT_COLORS.length];
        const pct = donutTotal > 0 ? (c.total / donutTotal) * 100 : 0;
        segs.push(`${color} ${acc}% ${acc + pct}%`);
        acc += pct;
        return { ...c, color, pct };
    });
    const donutBg = segs.length ? `conic-gradient(${segs.join(',')})` : '#eee';

    const history = data.history || [];
    const histMax = Math.max(1, ...history.map(h => h.entregado + h.en_curso));

    return (
        <div className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-3.5">
                <KpiCard accent="#0ea5b7" label="En curso" value={formatMoney(data.en_curso_total, currency)} sub={`${data.en_curso_orders} en transito`} />
                <KpiCard accent="#16a34a" label="Entregado" value={formatMoney(data.entregado_total, currency)} sub={`${data.entregado_orders} entregadas`} />
                <KpiCard accent="#d97706" label="Por pagar" value={formatMoney(data.total_pending, currency)} sub={`${data.orders_pending} por consignar`} />
                <KpiCard accent="#7c3aed" label="Recaudado" value={formatMoney(data.total_collected, currency)} sub={`${data.orders_collected} pagadas al cliente`} />
            </div>

            <div className={`${CARD} p-5`}>
                <div className="text-[14px] font-bold text-gray-900 dark:text-white">Flujo del dinero</div>
                <div className="text-[12px] text-gray-400 font-medium mb-4">Del recaudo en transito hasta el pago al cliente</div>
                <div className="flex flex-col gap-3">
                    {funnel.map(f => (
                        <div key={f.label} className="flex items-center gap-3.5">
                            <div className="w-[120px] shrink-0 text-[13px] font-semibold text-gray-600 dark:text-gray-300">{f.label}</div>
                            <div className="flex-1 h-[34px] bg-[#f3f3f6] dark:bg-gray-700 rounded-[9px] overflow-hidden relative">
                                <div className="h-full rounded-[9px] transition-all" style={{ width: `${Math.max((f.total / funnelMax) * 100, f.total > 0 ? 3 : 0)}%`, background: f.color }} />
                            </div>
                            <div className="w-[140px] shrink-0 text-right text-[14px] font-extrabold tabular-nums text-gray-900 dark:text-white">{formatMoney(f.total, currency)}</div>
                            <div className="w-[80px] shrink-0 text-right text-[12px] font-semibold text-gray-400">{f.orders} {f.orders === 1 ? 'orden' : 'ordenes'}</div>
                        </div>
                    ))}
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-[1.35fr_1fr] gap-4">
                <div className={`${CARD} p-5 min-w-0`}>
                    <div className="flex items-start justify-between gap-3 flex-wrap">
                        <div>
                            <div className="text-[14px] font-bold text-gray-900 dark:text-white">Historico de recaudo</div>
                            <div className="text-[12px] text-gray-400 font-medium">Total COD por {bucketWord}</div>
                        </div>
                        <div className="flex items-center gap-3 flex-wrap">
                            <div className="flex gap-3.5">
                                <div className="flex items-center gap-1.5 text-[11.5px] text-gray-400 font-semibold"><span className="w-2.5 h-2.5 rounded-sm" style={{ background: '#16a34a' }} />Entregado</div>
                                <div className="flex items-center gap-1.5 text-[11.5px] text-gray-400 font-semibold"><span className="w-2.5 h-2.5 rounded-sm" style={{ background: '#0ea5b7' }} />En curso</div>
                            </div>
                            <div className="inline-flex rounded-lg border border-[#ececf0] dark:border-gray-700 p-0.5 bg-[#fafafb] dark:bg-gray-900/40">
                                {BUCKETS.map(b => (
                                    <button
                                        key={b.key}
                                        onClick={() => setUserBucket(b.key)}
                                        className={`px-2.5 py-1 text-[11px] font-semibold rounded-md transition ${bucket === b.key ? 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white shadow-sm' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'}`}
                                    >
                                        {b.label}
                                    </button>
                                ))}
                            </div>
                        </div>
                    </div>
                    {history.length === 0 ? (
                        <div className="h-[200px] mt-5 flex items-center justify-center text-gray-400 text-sm">Sin movimientos en el periodo.</div>
                    ) : (
                        <div className={`flex items-end gap-4 h-[200px] mt-5 px-1.5 overflow-x-auto transition-opacity ${loading ? 'opacity-50' : ''}`}>
                            {history.map((d, i) => {
                                const tot = d.entregado + d.en_curso;
                                return (
                                    <div key={i} className="flex-1 min-w-[42px] flex flex-col items-center gap-2 h-full justify-end">
                                        <div className="text-[11px] font-bold text-gray-600 dark:text-gray-300 tabular-nums">{formatMoneyShort(tot)}</div>
                                        <div className="w-full max-w-[60px] flex flex-col justify-end rounded-lg overflow-hidden" style={{ height: `${(tot / histMax) * 100}%`, minHeight: tot > 0 ? 6 : 0 }}>
                                            <div style={{ height: `${tot > 0 ? (d.entregado / tot) * 100 : 0}%`, background: '#16a34a' }} />
                                            <div style={{ height: `${tot > 0 ? (d.en_curso / tot) * 100 : 0}%`, background: '#0ea5b7' }} />
                                        </div>
                                        <div className="text-[11px] text-gray-400 font-semibold whitespace-nowrap">{d.label}</div>
                                    </div>
                                );
                            })}
                        </div>
                    )}
                </div>

                <div className={`${CARD} p-5`}>
                    <div className="text-[14px] font-bold text-gray-900 dark:text-white">Recaudo por transportadora</div>
                    <div className="text-[12px] text-gray-400 font-medium mb-4">Participacion sobre total COD</div>
                    {donutTotal === 0 ? (
                        <div className="py-10 text-center text-gray-400 text-sm">Sin movimientos en el periodo.</div>
                    ) : (
                        <div className="flex items-center gap-5">
                            <div className="relative w-[150px] h-[150px] shrink-0">
                                <div className="absolute inset-0 rounded-full" style={{ background: donutBg }} />
                                <div className="absolute inset-[34px] bg-white dark:bg-gray-800 rounded-full flex flex-col items-center justify-center">
                                    <div className="text-[10px] text-gray-400 font-bold uppercase tracking-wide">Total</div>
                                    <div className="text-[15px] font-extrabold text-gray-900 dark:text-white">{formatMoneyShort(donutTotal)}</div>
                                </div>
                            </div>
                            <div className="flex-1 flex flex-col gap-3">
                                {donutCarriers.map(c => (
                                    <div key={c.carrier} className="flex items-center gap-2.5">
                                        <span className="w-3 h-3 rounded shrink-0" style={{ background: c.color }} />
                                        <div className="flex-1 min-w-0">
                                            <div className="text-[13px] font-semibold text-gray-800 dark:text-gray-200 truncate">{carrierLabel(c.carrier)}</div>
                                            <div className="text-[11.5px] text-gray-400 font-medium">{formatMoney(c.total, currency)}</div>
                                        </div>
                                        <div className="text-[14px] font-extrabold text-gray-600 dark:text-gray-300">{Math.round(c.pct)}%</div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            </div>

            <div className={`${CARD} overflow-hidden`}>
                <div className="px-5 pt-4 pb-3 text-[14px] font-bold text-gray-900 dark:text-white">Detalle por transportadora</div>
                <div className="overflow-x-auto">
                    <table className="w-full border-collapse min-w-[760px]">
                        <thead>
                            <tr className="bg-[#fafafb] dark:bg-gray-900/40">
                                <th className="text-left px-5 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">Transportadora</th>
                                <th className="text-right px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">Ordenes</th>
                                <th className="text-right px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">En curso</th>
                                <th className="text-right px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">Entregado</th>
                                <th className="text-right px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">Cargo carrier</th>
                                <th className="text-right px-5 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">Ticket promedio</th>
                            </tr>
                        </thead>
                        <tbody>
                            {detail.length === 0 && (
                                <tr><td colSpan={6} className="text-center py-10 text-gray-400 text-sm">Sin transportadoras en el periodo.</td></tr>
                            )}
                            {detail.map((d, i) => {
                                const logo = getCarrierLogo(d.carrier);
                                const color = DONUT_COLORS[i % DONUT_COLORS.length];
                                const ticket = d.orders > 0 ? d.total / d.orders : 0;
                                return (
                                    <tr key={d.carrier} style={{ borderTop: '1px solid #f2f2f5' }}>
                                        <td className="px-5 py-3.5">
                                            <div className="flex items-center gap-2.5">
                                                {logo ? (
                                                    <span className="inline-flex items-center justify-center h-7 w-10 rounded border border-gray-200 dark:border-gray-600 bg-white p-0.5 shrink-0"><img src={logo} alt={carrierLabel(d.carrier)} className="max-h-full max-w-full object-contain" /></span>
                                                ) : (
                                                    <span className="w-6 h-6 rounded-md text-white inline-flex items-center justify-center text-[12px] font-extrabold shrink-0" style={{ background: color }}>{(d.carrier || '?').charAt(0).toUpperCase()}</span>
                                                )}
                                                <span className="text-[13.5px] font-semibold text-gray-800 dark:text-gray-200">{carrierLabel(d.carrier)}</span>
                                            </div>
                                        </td>
                                        <td className="px-3 py-3.5 text-right text-[13px] text-gray-500 dark:text-gray-400 tabular-nums">{d.orders}</td>
                                        <td className="px-3 py-3.5 text-right text-[13px] font-semibold text-[#0e8a99] tabular-nums whitespace-nowrap">{formatMoney(d.en_curso, currency)}</td>
                                        <td className="px-3 py-3.5 text-right text-[13px] font-semibold text-[#15803d] tabular-nums whitespace-nowrap">{formatMoney(d.entregado, currency)}</td>
                                        <td className="px-3 py-3.5 text-right text-[13px] text-[#c2410c] tabular-nums whitespace-nowrap">{formatMoney(d.cargo, currency)}</td>
                                        <td className="px-5 py-3.5 text-right text-[13.5px] font-bold text-gray-900 dark:text-white tabular-nums whitespace-nowrap">{formatMoney(ticket, currency)}</td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}
