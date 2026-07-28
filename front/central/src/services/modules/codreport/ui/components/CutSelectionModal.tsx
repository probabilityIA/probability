'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { CheckCircle2, RefreshCw, ShieldCheck, AlertCircle, Package, Search, X, Calendar } from 'lucide-react';
import { getSelectableOrdersAction, createDraftCutAction } from '../../infra/actions';
import { CodOrder } from '../../domain/types';
import { formatMoney, formatDateTime, carrierLabel } from './helpers';

interface Props {
    isOpen: boolean;
    onClose: () => void;
    onConfirmed: (msg: string) => void;
    periodStart: string;
    periodEnd: string;
    businessId?: number | null;
}

function splitTerms(raw: string): string[] {
    return Array.from(new Set(raw.split(/[\s,;]+/).map(v => v.trim().toUpperCase()).filter(Boolean)));
}

function matchesTerms(order: CodOrder, terms: string[]): boolean {
    if (terms.length === 0) return true;
    const guide = (order.guide_number || '').toUpperCase();
    const number = (order.order_number || '').toUpperCase();
    return terms.some(t => guide.includes(t) || number.includes(t));
}

export function CutSelectionModal({ isOpen, onClose, onConfirmed, periodStart, periodEnd, businessId }: Props) {
    const [orders, setOrders] = useState<CodOrder[]>([]);
    const [picked, setPicked] = useState<Map<string, CodOrder>>(new Map());
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [start, setStart] = useState(periodStart);
    const [end, setEnd] = useState(periodEnd);
    const [guideFilter, setGuideFilter] = useState('');
    const [onlyPicked, setOnlyPicked] = useState(false);
    const firstLoad = useRef(true);

    useEffect(() => {
        if (!isOpen) return;
        setStart(periodStart);
        setEnd(periodEnd);
        setPicked(new Map());
        setGuideFilter('');
        setOnlyPicked(false);
        firstLoad.current = true;
    }, [isOpen, periodStart, periodEnd]);

    const load = useCallback(async () => {
        if (!start || !end || start > end) return;
        setLoading(true);
        setError(null);
        const res = await getSelectableOrdersAction(start, end, businessId || undefined);
        if (res.success) {
            const list = (res.data || []) as CodOrder[];
            setOrders(list);
            if (firstLoad.current) {
                setPicked(new Map(list.map(o => [o.order_id, o])));
                firstLoad.current = false;
            }
        } else {
            setError((res as any).message || 'Error al cargar las ordenes del periodo');
            setOrders([]);
        }
        setLoading(false);
    }, [start, end, businessId]);

    useEffect(() => {
        if (isOpen) load();
    }, [isOpen, load]);

    const terms = useMemo(() => splitTerms(guideFilter), [guideFilter]);

    const visible = useMemo(() => {
        const base = onlyPicked ? Array.from(picked.values()) : orders;
        return base.filter(o => matchesTerms(o, terms));
    }, [orders, picked, onlyPicked, terms]);

    const pickedList = useMemo(() => Array.from(picked.values()), [picked]);
    const pickedTotal = pickedList.reduce((sum, o) => sum + (o.cod_total || 0), 0);
    const currency = pickedList[0]?.currency || orders[0]?.currency;

    if (!isOpen) return null;

    const toggle = (order: CodOrder) => {
        setPicked(prev => {
            const next = new Map(prev);
            if (next.has(order.order_id)) next.delete(order.order_id);
            else next.set(order.order_id, order);
            return next;
        });
    };

    const allVisiblePicked = visible.length > 0 && visible.every(o => picked.has(o.order_id));
    const toggleAllVisible = () => {
        setPicked(prev => {
            const next = new Map(prev);
            if (allVisiblePicked) visible.forEach(o => next.delete(o.order_id));
            else visible.forEach(o => next.set(o.order_id, o));
            return next;
        });
    };

    const doConfirm = async () => {
        if (picked.size === 0) return;
        setSubmitting(true);
        setError(null);
        const res = await createDraftCutAction(start, end, Array.from(picked.keys()), businessId || undefined);
        if (res.success) {
            onConfirmed('Borrador de corte creado. Revisa y confirmalo para consignar.');
            onClose();
        } else {
            setError((res as any).message || 'Error al crear el borrador');
        }
        setSubmitting(false);
    };

    const dateInput = 'px-2 py-1 text-xs rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white';

    return (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-5xl w-full max-h-[88vh] flex flex-col">
                <div className="px-5 py-4 border-b border-gray-100 dark:border-gray-700">
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
                        <ShieldCheck size={18} className="text-emerald-600" /> Marcar corte de pago
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mt-1">
                        Elige las ordenes que iras a consignar. Se crea un borrador que luego confirmas
                        (o cancelas). Las que dejes sin marcar quedan disponibles para otro corte.
                        Lo que marques se conserva aunque cambies el filtro o el periodo.
                    </p>

                    <div className="flex items-center gap-3 mt-3 flex-wrap">
                        <div className="flex items-center gap-2">
                            <Calendar size={14} className="text-purple-600 shrink-0" />
                            <span className="text-xs font-semibold text-gray-500 dark:text-gray-400">Periodo del corte</span>
                            <input type="date" value={start} max={end || undefined} onChange={e => setStart(e.target.value)} className={dateInput} />
                            <span className="text-xs text-gray-400">a</span>
                            <input type="date" value={end} min={start || undefined} onChange={e => setEnd(e.target.value)} className={dateInput} />
                        </div>

                        <div className="relative flex-1 min-w-[220px]">
                            <Search size={14} className="absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" />
                            <input
                                type="text"
                                value={guideFilter}
                                onChange={e => setGuideFilter(e.target.value)}
                                placeholder="Buscar por numero de guia u orden (una o varias)"
                                className="w-full pl-8 pr-8 py-1.5 text-xs rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                            />
                            {guideFilter && (
                                <button
                                    onClick={() => setGuideFilter('')}
                                    className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                                    title="Quitar filtro"
                                >
                                    <X size={14} />
                                </button>
                            )}
                        </div>

                        <label className="flex items-center gap-1.5 text-xs text-gray-600 dark:text-gray-300 cursor-pointer whitespace-nowrap">
                            <input type="checkbox" checked={onlyPicked} onChange={e => setOnlyPicked(e.target.checked)} className="accent-emerald-600 cursor-pointer" />
                            Ver solo seleccionadas
                        </label>
                    </div>
                </div>

                <div className="flex-1 overflow-auto px-5 py-3">
                    {error && (
                        <div className="mb-3 p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm flex items-center gap-2">
                            <AlertCircle size={15} /> {error}
                        </div>
                    )}

                    {start > end && (
                        <div className="mb-3 p-3 bg-amber-50 border border-amber-200 rounded-md text-amber-700 text-sm flex items-center gap-2">
                            <AlertCircle size={15} /> La fecha inicial no puede ser posterior a la final.
                        </div>
                    )}

                    {loading ? (
                        <div className="flex items-center justify-center py-12 text-gray-400">
                            <RefreshCw size={18} className="animate-spin mr-2" /> Cargando ordenes...
                        </div>
                    ) : visible.length === 0 ? (
                        <div className="text-center py-12 text-gray-400 text-sm">
                            <Package size={28} className="mx-auto mb-2 opacity-50" />
                            {terms.length > 0
                                ? 'Ninguna orden coincide con la busqueda en este periodo.'
                                : 'No hay ordenes entregadas por pagar en este periodo.'}
                        </div>
                    ) : (
                        <table className="w-full text-sm">
                            <thead>
                                <tr className="text-[11px] uppercase text-gray-500 dark:text-gray-400 border-b border-gray-100 dark:border-gray-700">
                                    <th className="py-2 pr-2 w-8">
                                        <input type="checkbox" checked={allVisiblePicked} onChange={toggleAllVisible} className="accent-emerald-600 cursor-pointer" />
                                    </th>
                                    <th className="text-left py-2 font-semibold">Orden</th>
                                    <th className="text-left py-2 font-semibold">Numero de guia</th>
                                    <th className="text-left py-2 font-semibold">Cliente</th>
                                    <th className="text-left py-2 font-semibold">Transportadora</th>
                                    <th className="text-right py-2 font-semibold">COD orden</th>
                                    <th className="text-right py-2 font-semibold">Cargo carrier</th>
                                    <th className="text-right py-2 font-semibold">Total cliente</th>
                                    <th className="text-left py-2 font-semibold">Entregado</th>
                                </tr>
                            </thead>
                            <tbody>
                                {visible.map(o => {
                                    const checked = picked.has(o.order_id);
                                    return (
                                        <tr
                                            key={o.order_id}
                                            onClick={() => toggle(o)}
                                            className={`border-b border-gray-50 dark:border-gray-700/50 cursor-pointer ${checked ? 'bg-emerald-50/60 dark:bg-emerald-900/20' : 'hover:bg-gray-50 dark:hover:bg-gray-700/30'}`}
                                        >
                                            <td className="py-2 pr-2">
                                                <input type="checkbox" checked={checked} onChange={() => toggle(o)} onClick={e => e.stopPropagation()} className="accent-emerald-600 cursor-pointer" />
                                            </td>
                                            <td className="py-2 font-mono text-xs font-bold text-[#6d28d9] whitespace-nowrap">#{o.order_number || o.order_id.slice(0, 8)}</td>
                                            <td className="py-2 font-mono text-xs text-gray-600 dark:text-gray-300 whitespace-nowrap">{o.guide_number || '-'}</td>
                                            <td className="py-2">
                                                <div className="text-gray-900 dark:text-white truncate max-w-[150px]">{o.customer_name || '-'}</div>
                                                <div className="text-[10.5px] text-gray-400 whitespace-nowrap">{formatDateTime(o.created_at)}</div>
                                            </td>
                                            <td className="py-2 text-gray-600 dark:text-gray-300 whitespace-nowrap">{carrierLabel(o.carrier)}</td>
                                            <td className="py-2 text-right text-gray-600 dark:text-gray-300 tabular-nums whitespace-nowrap">{formatMoney(o.cod_total, o.currency)}</td>
                                            <td className="py-2 text-right text-[#c2410c] tabular-nums whitespace-nowrap">{o.cod_carrier_fee > 0 ? formatMoney(o.cod_carrier_fee, o.currency) : '-'}</td>
                                            <td className="py-2 text-right font-bold text-gray-900 dark:text-white tabular-nums whitespace-nowrap">{formatMoney(o.cod_total + (o.cod_carrier_fee || 0), o.currency)}</td>
                                            <td className="py-2 text-[11px] text-gray-500 whitespace-nowrap">{o.delivered_at ? formatDateTime(o.delivered_at) : '-'}</td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    )}
                </div>

                <div className="px-5 py-4 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between gap-3 flex-wrap">
                    <div className="text-sm">
                        <span className="text-gray-500 dark:text-gray-400">Seleccionadas: </span>
                        <span className="font-bold text-gray-800 dark:text-gray-200">{picked.size}</span>
                        <span className="text-gray-400 mx-2">|</span>
                        <span className="text-gray-500 dark:text-gray-400">Total: </span>
                        <span className="font-bold text-emerald-600">{formatMoney(pickedTotal, currency)}</span>
                        {(terms.length > 0 || onlyPicked) && (
                            <span className="text-[11px] text-gray-400 ml-2">
                                (mostrando {visible.length} de {orders.length})
                            </span>
                        )}
                    </div>
                    <div className="flex items-center gap-2">
                        {picked.size > 0 && (
                            <button
                                onClick={() => setPicked(new Map())}
                                disabled={submitting}
                                className="px-3 py-2 text-sm rounded-md text-gray-500 hover:text-gray-700 dark:hover:text-gray-200 disabled:opacity-50"
                            >
                                Limpiar seleccion
                            </button>
                        )}
                        <button
                            onClick={onClose}
                            disabled={submitting}
                            className="px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                        >
                            Cancelar
                        </button>
                        <button
                            onClick={doConfirm}
                            disabled={submitting || picked.size === 0}
                            className="px-3 py-2 text-sm rounded-md bg-emerald-600 hover:bg-emerald-700 text-white font-semibold inline-flex items-center gap-1.5 disabled:opacity-50"
                        >
                            {submitting ? <RefreshCw size={14} className="animate-spin" /> : <CheckCircle2 size={14} />}
                            Crear borrador
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
