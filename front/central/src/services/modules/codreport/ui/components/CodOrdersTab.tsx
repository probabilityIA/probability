'use client';

import { useCallback, useEffect, useState } from 'react';
import {
    RefreshCw, ChevronLeft, ChevronRight, Package, AlertCircle, CheckCircle2, Clock, Lock, FileText, Truck, Ban,
} from 'lucide-react';
import { getCodOrdersAction, getCodSummaryAction, getCarrierConfigsAction } from '../../infra/actions';
import { CodOrder, CodState, CodSummary, ReportFilters } from '../../domain/types';
import { formatMoney, formatDateTime, browserTimeZone, carrierLabel } from './helpers';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { GuidePreviewModal } from './GuidePreviewModal';

interface Props {
    filters: ReportFilters;
}

const CARD = 'bg-white dark:bg-gray-800 border border-[#ececf0] dark:border-gray-700 rounded-2xl';

const STATUS_PILL: Record<string, { bg: string; c: string; label: string }> = {
    pending: { bg: '#f1f5f9', c: '#475569', label: 'Pendiente' },
    picked_up: { bg: '#eef2ff', c: '#4f46e5', label: 'Recogido' },
    in_transit: { bg: '#eef2ff', c: '#4f46e5', label: 'En transito' },
    out_for_delivery: { bg: '#eef2ff', c: '#4338ca', label: 'En reparto' },
    on_hold: { bg: '#fef3c7', c: '#b45309', label: 'En espera' },
    delivered: { bg: '#dcfce7', c: '#15803d', label: 'Entregado' },
    failed: { bg: '#fee2e2', c: '#b91c1c', label: 'Fallido' },
    returned: { bg: '#fee2e2', c: '#b91c1c', label: 'Devuelto' },
    cancelled: { bg: '#f1f5f9', c: '#64748b', label: 'Cancelado' },
};

const CARRIER_CHIP_COLORS = ['#e11d48', '#ea580c', '#0891b2', '#7c3aed', '#0ea5e9', '#16a34a', '#db2777', '#4f46e5', '#f59e0b', '#14b8a6'];
function carrierChipColor(name: string): string {
    let h = 0;
    for (let i = 0; i < name.length; i++) h = (h * 31 + name.charCodeAt(i)) >>> 0;
    return CARRIER_CHIP_COLORS[h % CARRIER_CHIP_COLORS.length];
}

const REC_META: Record<CodState, { c: string; label: string; icon: any; title?: string }> = {
    collected: { c: '#16a34a', label: 'Recaudada', icon: CheckCircle2 },
    pending_payment: { c: '#d97706', label: 'Pendiente de pago', icon: Clock, title: 'Entregada: falta que el administrador la marque como pagada al cliente' },
    in_progress: { c: '#2563eb', label: 'En progreso', icon: Truck, title: 'En curso: no se cuenta como recaudada hasta que se entregue' },
    pending: { c: '#d97706', label: 'Pendiente', icon: Clock },
    not_collectable: { c: '#9a9aa5', label: 'No recaudable', icon: Ban, title: 'No recaudable' },
};

function RecaudoBadge({ state }: { state: CodState }) {
    const m = REC_META[state] || REC_META.not_collectable;
    const Icon = m.icon;
    return (
        <span className="inline-flex items-center gap-1.5 text-[12.5px] font-bold whitespace-nowrap" style={{ color: m.c }} title={m.title}>
            <span className="w-[7px] h-[7px] rounded-full shrink-0" style={{ background: m.c }} />
            <Icon size={13} /> {m.label}
        </span>
    );
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

const ESTADO_TABS: { key: '' | 'false' | 'true'; label: string }[] = [
    { key: '', label: 'Todas' },
    { key: 'false', label: 'Por pagar' },
    { key: 'true', label: 'Pagadas' },
];

export default function CodOrdersTab({ filters }: Props) {
    const [orders, setOrders] = useState<CodOrder[]>([]);
    const [summary, setSummary] = useState<CodSummary | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const pageSize = 15;
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(0);
    const [collected, setCollected] = useState<'' | 'true' | 'false'>('');
    const [hasGuide, setHasGuide] = useState<'' | 'true' | 'false'>('');
    const [status, setStatus] = useState('');
    const [carrier, setCarrier] = useState('');
    const [carriers, setCarriers] = useState<string[]>([]);
    const [search, setSearch] = useState('');
    const [debounced, setDebounced] = useState('');
    const [tz, setTz] = useState('');
    const [guidePreview, setGuidePreview] = useState<CodOrder | null>(null);

    useEffect(() => {
        let cancelled = false;
        getCarrierConfigsAction(filters.business_id).then(res => {
            if (cancelled || !res.success) return;
            const names = (res.data || []).map((c: any) => c.carrier_name).filter(Boolean);
            setCarriers(Array.from(new Set(names)) as string[]);
        });
        return () => { cancelled = true; };
    }, [filters.business_id]);

    const availableFilters: FilterOption[] = [
        { key: 'search', label: 'Orden / Cliente', type: 'text', placeholder: 'Numero de orden o cliente' },
        {
            key: 'status', label: 'Estado', type: 'select',
            options: Object.entries(STATUS_PILL).map(([value, v]) => ({ value, label: v.label })),
        },
        {
            key: 'carrier', label: 'Transportadora', type: 'select',
            options: carriers.map(c => ({ value: c, label: carrierLabel(c) })),
        },
        {
            key: 'has_guide', label: 'Guia', type: 'select',
            options: [{ value: 'true', label: 'Con guia' }, { value: 'false', label: 'Sin guia' }],
        },
    ];

    const activeFilters: ActiveFilter[] = [];
    if (search) activeFilters.push({ key: 'search', label: 'Orden / Cliente', value: search, type: 'text' });
    if (status) activeFilters.push({ key: 'status', label: 'Estado', value: STATUS_PILL[status]?.label || status, type: 'select' });
    if (carrier) activeFilters.push({ key: 'carrier', label: 'Transportadora', value: carrierLabel(carrier), type: 'select' });
    if (hasGuide) activeFilters.push({ key: 'has_guide', label: 'Guia', value: hasGuide === 'true' ? 'Con guia' : 'Sin guia', type: 'select' });

    const onAddFilter = (key: string, value: any) => {
        if (key === 'search') setSearch(String(value));
        else if (key === 'status') setStatus(String(value));
        else if (key === 'carrier') setCarrier(String(value));
        else if (key === 'has_guide') setHasGuide(String(value) as any);
    };
    const onRemoveFilter = (key: string) => {
        if (key === 'search') setSearch('');
        else if (key === 'status') setStatus('');
        else if (key === 'carrier') setCarrier('');
        else if (key === 'has_guide') setHasGuide('');
    };

    useEffect(() => { setTz(browserTimeZone()); }, []);

    useEffect(() => {
        const t = setTimeout(() => setDebounced(search.trim()), 400);
        return () => clearTimeout(t);
    }, [search]);

    useEffect(() => { setPage(1); }, [filters, collected, hasGuide, status, carrier, debounced]);

    useEffect(() => {
        let cancelled = false;
        getCodSummaryAction(filters).then(res => {
            if (!cancelled && res.success) setSummary(res.data);
        });
        return () => { cancelled = true; };
    }, [filters]);

    const load = useCallback(async () => {
        setLoading(true);
        setError(null);
        const res = await getCodOrdersAction({
            ...filters,
            page,
            page_size: pageSize,
            collected: collected === '' ? undefined : collected === 'true',
            has_guide: hasGuide === '' ? undefined : hasGuide === 'true',
            status: status || undefined,
            carrier: carrier || undefined,
            search: debounced || undefined,
        });
        if (res.success) {
            setOrders(res.data || []);
            setTotal(res.total || 0);
            setTotalPages(res.total_pages || 0);
        } else {
            setError((res as any).message || 'Error al cargar las ordenes');
            setOrders([]);
        }
        setLoading(false);
    }, [filters, page, collected, hasGuide, status, carrier, debounced]);

    useEffect(() => { load(); }, [load]);

    const filteredTotal = orders.reduce((a, o) => a + o.cod_total + (o.cod_carrier_fee || 0), 0);
    const currency = orders[0]?.currency;

    return (
        <div className="space-y-5">
            <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-3.5">
                <KpiCard accent="#16a34a" label="Recaudado"
                    value={formatMoney(summary?.total_collected, currency)}
                    sub={`${summary?.orders_collected ?? 0} ordenes pagadas al cliente`} />
                <KpiCard accent="#d97706" label="Por pagar"
                    value={formatMoney(summary?.total_pending, currency)}
                    sub={`${summary?.orders_pending ?? 0} entregadas sin consignar`} />
                <KpiCard accent="#e11d48" label="Descuento carrier"
                    value={formatMoney(summary?.total_discount, currency)}
                    sub="sobre lo recaudado" />
                <KpiCard accent="#7c3aed" label="Neto recibido"
                    value={formatMoney(summary?.total_net, currency)}
                    sub="despues de descuento" />
            </div>

            <div className={`${CARD} overflow-hidden`}>
                <div className="flex items-center gap-3 px-[18px] py-3 border-b border-[#f0f0f3] dark:border-gray-700 flex-wrap">
                    <div className="flex-1 min-w-[200px]">
                        <DynamicFilters
                            availableFilters={availableFilters}
                            activeFilters={activeFilters}
                            onAddFilter={onAddFilter}
                            onRemoveFilter={onRemoveFilter}
                            className="!p-0 !bg-transparent !border-0 !shadow-none !rounded-none"
                        />
                    </div>

                    <span className="text-[12px] font-bold text-[#6d28d9] bg-[#f3ecff] dark:bg-purple-900/40 px-2.5 py-0.5 rounded-full shrink-0">{total}</span>

                    <div className="flex gap-[3px] bg-[#f2f2f5] dark:bg-gray-700 p-[3px] rounded-[11px] shrink-0">
                        {ESTADO_TABS.map(t => {
                            const active = collected === t.key;
                            return (
                                <button
                                    key={t.key}
                                    onClick={() => setCollected(t.key)}
                                    className="px-3.5 py-[7px] rounded-lg text-[13px] font-bold transition-colors"
                                    style={active
                                        ? { background: '#fff', color: '#4c1d95', boxShadow: '0 1px 2px rgba(0,0,0,.08)' }
                                        : { color: '#6a6a76', background: 'transparent' }}
                                >
                                    {t.label}
                                </button>
                            );
                        })}
                    </div>
                </div>

                {error && (
                    <div className="m-3 p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm flex items-center gap-2">
                        <AlertCircle size={15} /> {error}
                    </div>
                )}

                <div className="overflow-x-auto">
                    <table className="w-full border-collapse min-w-[1100px]">
                        <thead>
                            <tr className="bg-[#fafafb] dark:bg-gray-900/40">
                                {['Orden', 'Cliente', 'Transportadora', 'Guia PDF', 'Estado'].map(h => (
                                    <th key={h} className="text-left px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">{h}</th>
                                ))}
                                <th className="text-right px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">COD orden</th>
                                <th className="text-right px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">Cargo carrier</th>
                                <th className="text-right px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">Total cliente</th>
                                <th className="text-left px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider">Recaudo</th>
                                <th className="text-left px-3 py-[11px] text-[10.5px] font-bold text-[#9a9aa5] uppercase tracking-wider" title={tz ? `Zona horaria: ${tz}` : undefined}>Corte / Entregado</th>
                            </tr>
                        </thead>
                        <tbody>
                            {loading && (
                                <tr><td colSpan={10} className="text-center py-12 text-gray-400">
                                    <RefreshCw size={18} className="animate-spin inline mr-2" /> Cargando...
                                </td></tr>
                            )}
                            {!loading && orders.length === 0 && !error && (
                                <tr><td colSpan={10} className="text-center py-12 text-gray-400 text-sm">
                                    <Package size={28} className="mx-auto mb-2 opacity-50" />
                                    No hay ordenes que coincidan con los filtros.
                                </td></tr>
                            )}
                            {!loading && orders.map((o, i) => {
                                const st = STATUS_PILL[o.status] || { bg: '#f1f5f9', c: '#475569', label: o.status };
                                const cc = carrierChipColor(o.carrier || '');
                                const logo = getCarrierLogo(o.carrier);
                                return (
                                    <tr key={o.order_id} className="hover:bg-[#fafafb] dark:hover:bg-gray-700/30" style={{ borderTop: i === 0 ? 'none' : '1px solid #f2f2f5' }}>
                                        <td className="px-3 py-3.5 text-[13px] font-bold text-[#6d28d9] font-mono whitespace-nowrap">#{o.order_number || o.order_id.slice(0, 8)}</td>
                                        <td className="px-3 py-3.5">
                                            <div className="text-[13.5px] font-semibold text-[#26262e] dark:text-white truncate max-w-[170px]">{o.customer_name || '-'}</div>
                                            <div className="text-[11.5px] text-[#a0a0ab] mt-0.5 whitespace-nowrap">{formatDateTime(o.created_at)}</div>
                                        </td>
                                        <td className="px-3 py-3.5">
                                            <div className="flex items-center gap-2">
                                                {logo ? (
                                                    <span className="inline-flex items-center justify-center h-7 w-10 rounded border border-gray-200 dark:border-gray-600 bg-white p-0.5 shrink-0" title={carrierLabel(o.carrier)}>
                                                        <img src={logo} alt={carrierLabel(o.carrier)} className="max-h-full max-w-full object-contain" />
                                                    </span>
                                                ) : (
                                                    <span className="w-[22px] h-[22px] rounded-md text-white inline-flex items-center justify-center text-[11px] font-extrabold shrink-0" style={{ background: cc }}>
                                                        {(o.carrier || '?').charAt(0).toUpperCase()}
                                                    </span>
                                                )}
                                                <span className="text-[13px] text-[#3a3a44] dark:text-gray-300 font-medium whitespace-nowrap">{carrierLabel(o.carrier)}</span>
                                            </div>
                                        </td>
                                        <td className="px-3 py-3.5 text-center">
                                            {o.has_guide && o.shipment_id ? (
                                                <button
                                                    onClick={() => setGuidePreview(o)}
                                                    className="inline-flex items-center gap-1 px-2 py-1 rounded-md text-[11px] font-semibold text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800 hover:bg-emerald-50 dark:hover:bg-emerald-900/30 transition-colors"
                                                    title="Ver y descargar guia PDF"
                                                >
                                                    <FileText size={13} /> Ver PDF
                                                </button>
                                            ) : (
                                                <span className="text-gray-300 text-xs">-</span>
                                            )}
                                        </td>
                                        <td className="px-3 py-3.5">
                                            <span className="inline-flex items-center px-2.5 py-1 rounded-full text-[12px] font-bold whitespace-nowrap" style={{ background: st.bg, color: st.c }}>{st.label}</span>
                                        </td>
                                        <td className="px-3 py-3.5 text-right text-[13px] text-[#4a4a54] dark:text-gray-300 tabular-nums whitespace-nowrap">{formatMoney(o.cod_total, o.currency)}</td>
                                        <td className="px-3 py-3.5 text-right text-[13px] text-[#c2410c] tabular-nums whitespace-nowrap">{o.cod_carrier_fee > 0 ? formatMoney(o.cod_carrier_fee, o.currency) : '-'}</td>
                                        <td className="px-3 py-3.5 text-right text-[13.5px] font-bold text-gray-900 dark:text-white tabular-nums whitespace-nowrap">{formatMoney(o.cod_total + (o.cod_carrier_fee || 0), o.currency)}</td>
                                        <td className="px-3 py-3.5"><RecaudoBadge state={o.cod_state} /></td>
                                        <td className="px-3 py-3.5">
                                            {o.cut_status === 'confirmed' ? (
                                                <div className="text-[12.5px] font-semibold text-emerald-600 inline-flex items-center gap-1"><Lock size={12} /> Confirmado</div>
                                            ) : o.collected ? (
                                                <div className="text-[12.5px] font-semibold text-[#c2410c]">Sin confirmar</div>
                                            ) : (
                                                <div className="text-[12.5px] font-semibold text-[#9a9aa5]">—</div>
                                            )}
                                            <div className="text-[11.5px] text-[#a0a0ab] mt-0.5 whitespace-nowrap">{o.delivered_at ? formatDateTime(o.delivered_at) : '—'}</div>
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>

                <div className="flex items-center justify-between px-[18px] py-3 border-t border-[#f0f0f3] dark:border-gray-700 bg-[#fafafb] dark:bg-gray-900/40 flex-wrap gap-2">
                    <div className="text-[12.5px] text-gray-400 dark:text-gray-500 font-medium">
                        Mostrando <b className="text-[#3a3a44] dark:text-gray-200">{orders.length}</b> de {total} ordenes
                    </div>
                    <div className="flex items-center gap-4">
                        <div className="text-[12.5px] text-gray-400 dark:text-gray-500 font-semibold">
                            Total cliente (pagina): <b className="text-gray-900 dark:text-white">{formatMoney(filteredTotal, currency)}</b>
                        </div>
                        {totalPages > 1 && (
                            <div className="flex items-center gap-1">
                                <span className="text-[12px] text-gray-500 mr-1">Pag {page} / {totalPages}</span>
                                <button
                                    onClick={() => setPage(p => Math.max(1, p - 1))}
                                    disabled={page <= 1}
                                    className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30"
                                >
                                    <ChevronLeft size={15} />
                                </button>
                                <button
                                    onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                                    disabled={page >= totalPages}
                                    className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30"
                                >
                                    <ChevronRight size={15} />
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            <GuidePreviewModal
                isOpen={!!guidePreview}
                onClose={() => setGuidePreview(null)}
                shipmentId={guidePreview?.shipment_id ?? null}
                orderLabel={guidePreview ? `#${guidePreview.order_number || guidePreview.order_id.slice(0, 8)}` : undefined}
                carrierLabel={guidePreview ? carrierLabel(guidePreview.carrier) : undefined}
            />
        </div>
    );
}
