'use client';

import { useCallback, useEffect, useState } from 'react';
import {
    Search, RefreshCw, ChevronLeft, ChevronRight, Package, AlertCircle, CheckCircle2, Clock, Lock, FileCheck2, FileX2,
} from 'lucide-react';
import { getCodOrdersAction } from '../../infra/actions';
import { CodOrder, ReportFilters } from '../../domain/types';
import { formatMoney, formatDate, carrierLabel } from './helpers';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';

interface Props {
    filters: ReportFilters;
}

const STATUS_LABELS: Record<string, { label: string; cls: string }> = {
    pending: { label: 'Pendiente', cls: 'bg-gray-100 text-gray-600' },
    picked_up: { label: 'Recogido', cls: 'bg-blue-100 text-blue-700' },
    in_transit: { label: 'En transito', cls: 'bg-blue-100 text-blue-700' },
    out_for_delivery: { label: 'En reparto', cls: 'bg-indigo-100 text-indigo-700' },
    on_hold: { label: 'En espera', cls: 'bg-amber-100 text-amber-700' },
    delivered: { label: 'Entregado', cls: 'bg-emerald-100 text-emerald-700' },
    failed: { label: 'Fallido', cls: 'bg-red-100 text-red-700' },
    returned: { label: 'Devuelto', cls: 'bg-red-100 text-red-700' },
    cancelled: { label: 'Cancelado', cls: 'bg-gray-100 text-gray-500' },
};

export default function CodOrdersTab({ filters }: Props) {
    const [orders, setOrders] = useState<CodOrder[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const pageSize = 15;
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(0);
    const [collected, setCollected] = useState<'' | 'true' | 'false'>('');
    const [hasGuide, setHasGuide] = useState<'' | 'true' | 'false'>('');
    const [search, setSearch] = useState('');
    const [debounced, setDebounced] = useState('');

    useEffect(() => {
        const t = setTimeout(() => setDebounced(search.trim()), 400);
        return () => clearTimeout(t);
    }, [search]);

    useEffect(() => { setPage(1); }, [filters, collected, hasGuide, debounced]);

    const load = useCallback(async () => {
        setLoading(true);
        setError(null);
        const res = await getCodOrdersAction({
            ...filters,
            page,
            page_size: pageSize,
            collected: collected === '' ? undefined : collected === 'true',
            has_guide: hasGuide === '' ? undefined : hasGuide === 'true',
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
    }, [filters, page, collected, hasGuide, debounced]);

    useEffect(() => { load(); }, [load]);

    return (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700">
            <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center gap-3 flex-wrap">
                <div className="flex items-center gap-2">
                    <Package size={16} className="text-purple-600" />
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Ordenes contra entrega</h3>
                    <span className="text-xs text-gray-500">({total})</span>
                </div>
                <div className="flex-1" />
                <div className="relative">
                    <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input
                        value={search}
                        onChange={e => setSearch(e.target.value)}
                        placeholder="Buscar orden o cliente..."
                        className="pl-9 pr-3 py-1.5 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white w-56"
                    />
                </div>
                <select
                    value={collected}
                    onChange={e => setCollected(e.target.value as any)}
                    className="px-2 py-1.5 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                >
                    <option value="">Todas</option>
                    <option value="true">Recaudadas</option>
                    <option value="false">Por recaudar</option>
                </select>
                <select
                    value={hasGuide}
                    onChange={e => setHasGuide(e.target.value as any)}
                    className="px-2 py-1.5 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                >
                    <option value="">Guia: todas</option>
                    <option value="true">Con guia</option>
                    <option value="false">Sin guia</option>
                </select>
                <button
                    onClick={load}
                    disabled={loading}
                    className="p-1.5 rounded-md border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                    title="Refrescar"
                >
                    <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
                </button>
            </div>

            {error && (
                <div className="m-3 p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm flex items-center gap-2">
                    <AlertCircle size={15} /> {error}
                </div>
            )}

            <div className="overflow-x-auto">
                <table className="w-full text-sm">
                    <thead>
                        <tr className="text-[11px] uppercase text-gray-500 dark:text-gray-400 border-b border-gray-100 dark:border-gray-700">
                            <th className="text-left px-3 py-2 font-semibold">Orden</th>
                            <th className="text-left px-3 py-2 font-semibold">Cliente</th>
                            <th className="text-left px-3 py-2 font-semibold">Transportadora</th>
                            <th className="text-center px-3 py-2 font-semibold">Guia</th>
                            <th className="text-left px-3 py-2 font-semibold">Estado</th>
                            <th className="text-right px-3 py-2 font-semibold">COD orden (prod + envio)</th>
                            <th className="text-right px-3 py-2 font-semibold">Cargo COD carrier</th>
                            <th className="text-right px-3 py-2 font-semibold">Total cliente</th>
                            <th className="text-center px-3 py-2 font-semibold">Recaudo</th>
                            <th className="text-center px-3 py-2 font-semibold">Corte</th>
                            <th className="text-left px-3 py-2 font-semibold">Entregado</th>
                        </tr>
                    </thead>
                    <tbody>
                        {loading && (
                            <tr><td colSpan={11} className="text-center py-10 text-gray-400">
                                <RefreshCw size={18} className="animate-spin inline mr-2" /> Cargando...
                            </td></tr>
                        )}
                        {!loading && orders.length === 0 && !error && (
                            <tr><td colSpan={11} className="text-center py-10 text-gray-400 text-sm">
                                <Package size={28} className="mx-auto mb-2 opacity-50" />
                                No hay ordenes contra entrega en el periodo.
                            </td></tr>
                        )}
                        {!loading && orders.map(o => {
                            const st = STATUS_LABELS[o.status] || { label: o.status, cls: 'bg-gray-100 text-gray-600' };
                            return (
                                <tr key={o.order_id} className="border-b border-gray-50 dark:border-gray-700/50 hover:bg-gray-50 dark:hover:bg-gray-700/30">
                                    <td className="px-3 py-2 font-mono text-xs text-gray-700 dark:text-gray-300">#{o.order_number || o.order_id.slice(0, 8)}</td>
                                    <td className="px-3 py-2 text-gray-900 dark:text-white truncate max-w-[160px]">{o.customer_name || '-'}</td>
                                    <td className="px-3 py-2 text-gray-600 dark:text-gray-300">{carrierLabel(o.carrier)}</td>
                                    <td className="px-3 py-2 text-center">
                                        {o.has_guide ? (
                                            (() => {
                                                const logo = getCarrierLogo(o.carrier);
                                                return logo ? (
                                                    <span className="inline-flex items-center justify-center h-7 w-12 rounded border border-gray-200 dark:border-gray-600 bg-white p-0.5" title={`Guia generada - ${carrierLabel(o.carrier)}`}>
                                                        <img src={logo} alt={carrierLabel(o.carrier)} className="max-h-full max-w-full object-contain" />
                                                    </span>
                                                ) : (
                                                    <span className="inline-flex items-center gap-1 text-emerald-600 text-xs font-semibold" title={`Guia generada - ${carrierLabel(o.carrier)}`}>
                                                        <FileCheck2 size={13} /> Generada
                                                    </span>
                                                );
                                            })()
                                        ) : (
                                            <span className="inline-flex items-center gap-1 text-gray-400 text-xs" title="Sin guia generada">
                                                <FileX2 size={13} /> Sin guia
                                            </span>
                                        )}
                                    </td>
                                    <td className="px-3 py-2">
                                        <span className={`inline-block px-2 py-0.5 rounded-full text-[11px] font-semibold ${st.cls}`}>{st.label}</span>
                                    </td>
                                    <td className="px-3 py-2 text-right font-semibold text-gray-900 dark:text-white">{formatMoney(o.cod_total, o.currency)}</td>
                                    <td className="px-3 py-2 text-right text-amber-700 dark:text-amber-400">{o.cod_carrier_fee > 0 ? formatMoney(o.cod_carrier_fee, o.currency) : '-'}</td>
                                    <td className="px-3 py-2 text-right font-semibold text-blue-700 dark:text-blue-300">{formatMoney(o.cod_total + (o.cod_carrier_fee || 0), o.currency)}</td>
                                    <td className="px-3 py-2 text-center">
                                        {o.collected ? (
                                            <span className="inline-flex items-center gap-1 text-emerald-600 text-xs font-semibold">
                                                <CheckCircle2 size={13} /> Recaudada
                                            </span>
                                        ) : (
                                            <span className="inline-flex items-center gap-1 text-amber-600 text-xs font-semibold">
                                                <Clock size={13} /> Pendiente
                                            </span>
                                        )}
                                    </td>
                                    <td className="px-3 py-2 text-center">
                                        {o.cut_status === 'confirmed' ? (
                                            <span className="inline-flex items-center gap-1 text-emerald-600 text-xs font-semibold">
                                                <Lock size={12} /> Confirmado
                                            </span>
                                        ) : o.collected ? (
                                            <span className="text-gray-400 text-xs">Sin confirmar</span>
                                        ) : (
                                            <span className="text-gray-300 text-xs">-</span>
                                        )}
                                    </td>
                                    <td className="px-3 py-2 text-xs text-gray-500">{o.delivered_at ? formatDate(o.delivered_at) : '-'}</td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            </div>

            {totalPages > 1 && (
                <div className="flex items-center justify-between px-4 py-2 border-t border-gray-200 dark:border-gray-700 text-xs">
                    <span className="text-gray-500">Pagina {page} de {totalPages}</span>
                    <div className="flex items-center gap-1">
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
                </div>
            )}
        </div>
    );
}
