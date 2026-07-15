'use client';

import { useCallback, useEffect, useState } from 'react';
import { CheckCircle2, RefreshCw, ShieldCheck, AlertCircle, Package } from 'lucide-react';
import { getSelectableOrdersAction, createDraftCutAction } from '../../infra/actions';
import { CodOrder } from '../../domain/types';
import { formatMoney, formatDateOnly, formatDateTime, carrierLabel } from './helpers';

interface Props {
    isOpen: boolean;
    onClose: () => void;
    onConfirmed: (msg: string) => void;
    periodStart: string;
    periodEnd: string;
    businessId?: number | null;
}

export function CutSelectionModal({ isOpen, onClose, onConfirmed, periodStart, periodEnd, businessId }: Props) {
    const [orders, setOrders] = useState<CodOrder[]>([]);
    const [selected, setSelected] = useState<Set<string>>(new Set());
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        setError(null);
        const res = await getSelectableOrdersAction(periodStart, periodEnd, businessId || undefined);
        if (res.success) {
            const list = (res.data || []) as CodOrder[];
            setOrders(list);
            setSelected(new Set(list.map(o => o.order_id)));
        } else {
            setError((res as any).message || 'Error al cargar las ordenes de la semana');
            setOrders([]);
        }
        setLoading(false);
    }, [periodStart, periodEnd, businessId]);

    useEffect(() => {
        if (isOpen) load();
    }, [isOpen, load]);

    if (!isOpen) return null;

    const toggle = (id: string) => {
        setSelected(prev => {
            const next = new Set(prev);
            if (next.has(id)) next.delete(id); else next.add(id);
            return next;
        });
    };

    const allSelected = orders.length > 0 && selected.size === orders.length;
    const toggleAll = () => {
        setSelected(allSelected ? new Set() : new Set(orders.map(o => o.order_id)));
    };

    const selectedTotal = orders
        .filter(o => selected.has(o.order_id))
        .reduce((sum, o) => sum + (o.cod_total || 0), 0);
    const currency = orders[0]?.currency;

    const doConfirm = async () => {
        if (selected.size === 0) return;
        setSubmitting(true);
        setError(null);
        const res = await createDraftCutAction(periodStart, periodEnd, Array.from(selected), businessId || undefined);
        if (res.success) {
            onConfirmed('Borrador de corte creado. Revisa y confirmalo para consignar.');
            onClose();
        } else {
            setError((res as any).message || 'Error al crear el borrador');
        }
        setSubmitting(false);
    };

    return (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-4xl w-full max-h-[85vh] flex flex-col">
                <div className="px-5 py-4 border-b border-gray-100 dark:border-gray-700">
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
                        <ShieldCheck size={18} className="text-emerald-600" /> Marcar corte de pago
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mt-1">
                        Semana <strong>{formatDateOnly(periodStart)} - {formatDateOnly(periodEnd)}</strong>.
                        Elige las ordenes que iras a consignar. Se crea un borrador que luego confirmas
                        (o cancelas). Las que dejes sin marcar quedan disponibles para otro corte.
                    </p>
                </div>

                <div className="flex-1 overflow-auto px-5 py-3">
                    {error && (
                        <div className="mb-3 p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm flex items-center gap-2">
                            <AlertCircle size={15} /> {error}
                        </div>
                    )}

                    {loading ? (
                        <div className="flex items-center justify-center py-12 text-gray-400">
                            <RefreshCw size={18} className="animate-spin mr-2" /> Cargando ordenes...
                        </div>
                    ) : orders.length === 0 ? (
                        <div className="text-center py-12 text-gray-400 text-sm">
                            <Package size={28} className="mx-auto mb-2 opacity-50" />
                            No hay ordenes entregadas por pagar en esta semana.
                        </div>
                    ) : (
                        <table className="w-full text-sm">
                            <thead>
                                <tr className="text-[11px] uppercase text-gray-500 dark:text-gray-400 border-b border-gray-100 dark:border-gray-700">
                                    <th className="py-2 pr-2 w-8">
                                        <input type="checkbox" checked={allSelected} onChange={toggleAll} className="accent-emerald-600 cursor-pointer" />
                                    </th>
                                    <th className="text-left py-2 font-semibold">Orden</th>
                                    <th className="text-left py-2 font-semibold">Cliente</th>
                                    <th className="text-left py-2 font-semibold">Transportadora</th>
                                    <th className="text-right py-2 font-semibold">COD orden</th>
                                    <th className="text-right py-2 font-semibold">Cargo carrier</th>
                                    <th className="text-right py-2 font-semibold">Total cliente</th>
                                    <th className="text-left py-2 font-semibold">Entregado</th>
                                </tr>
                            </thead>
                            <tbody>
                                {orders.map(o => {
                                    const checked = selected.has(o.order_id);
                                    return (
                                        <tr
                                            key={o.order_id}
                                            onClick={() => toggle(o.order_id)}
                                            className={`border-b border-gray-50 dark:border-gray-700/50 cursor-pointer ${checked ? 'bg-emerald-50/60 dark:bg-emerald-900/20' : 'hover:bg-gray-50 dark:hover:bg-gray-700/30'}`}
                                        >
                                            <td className="py-2 pr-2">
                                                <input type="checkbox" checked={checked} onChange={() => toggle(o.order_id)} onClick={e => e.stopPropagation()} className="accent-emerald-600 cursor-pointer" />
                                            </td>
                                            <td className="py-2 font-mono text-xs font-bold text-[#6d28d9] whitespace-nowrap">#{o.order_number || o.order_id.slice(0, 8)}</td>
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

                <div className="px-5 py-4 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between gap-3">
                    <div className="text-sm">
                        <span className="text-gray-500 dark:text-gray-400">Seleccionadas: </span>
                        <span className="font-bold text-gray-800 dark:text-gray-200">{selected.size}</span>
                        <span className="text-gray-400 mx-2">|</span>
                        <span className="text-gray-500 dark:text-gray-400">Total: </span>
                        <span className="font-bold text-emerald-600">{formatMoney(selectedTotal, currency)}</span>
                    </div>
                    <div className="flex items-center gap-2">
                        <button
                            onClick={onClose}
                            disabled={submitting}
                            className="px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                        >
                            Cancelar
                        </button>
                        <button
                            onClick={doConfirm}
                            disabled={submitting || selected.size === 0}
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
