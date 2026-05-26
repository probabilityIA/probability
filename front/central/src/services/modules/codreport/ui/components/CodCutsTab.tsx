'use client';

import { useCallback, useEffect, useState } from 'react';
import {
    Calendar, CheckCircle2, ChevronDown, ChevronUp, RefreshCw, AlertCircle, ShieldCheck, Clock,
} from 'lucide-react';
import { getCodCutsAction, confirmCodCutAction } from '../../infra/actions';
import { PaymentCut } from '../../domain/types';
import { formatMoney, formatDate, formatDateOnly, carrierLabel } from './helpers';

interface Props {
    businessId?: number | null;
    isAdmin: boolean;
}

function periodLabel(start: string, end: string): string {
    return `${formatDateOnly(start)} - ${formatDateOnly(end)}`;
}

export default function CodCutsTab({ businessId, isAdmin }: Props) {
    const [cuts, setCuts] = useState<PaymentCut[]>([]);
    const [canConfirm, setCanConfirm] = useState(false);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [expanded, setExpanded] = useState<string | null>(null);
    const [confirmTarget, setConfirmTarget] = useState<PaymentCut | null>(null);
    const [confirming, setConfirming] = useState(false);
    const [feedback, setFeedback] = useState<{ ok: boolean; msg: string } | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        setError(null);
        const res = await getCodCutsAction(businessId || undefined);
        if (res.success) {
            setCuts(res.data || []);
            setCanConfirm(res.can_confirm);
        } else {
            setError((res as any).message || 'Error al cargar los cortes');
            setCuts([]);
        }
        setLoading(false);
    }, [businessId]);

    useEffect(() => { load(); }, [load]);

    const doConfirm = async () => {
        if (!confirmTarget) return;
        setConfirming(true);
        setFeedback(null);
        const res = await confirmCodCutAction(
            confirmTarget.period_start.slice(0, 10),
            confirmTarget.period_end.slice(0, 10),
            businessId || undefined,
        );
        if (res.success) {
            setFeedback({ ok: true, msg: 'Corte de pago confirmado exitosamente' });
            setConfirmTarget(null);
            await load();
        } else {
            setFeedback({ ok: false, msg: (res as any).message || 'Error al confirmar el corte' });
        }
        setConfirming(false);
        setTimeout(() => setFeedback(null), 3500);
    };

    return (
        <div className="space-y-3">
            <div className="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800 rounded-lg px-3 py-2 border border-gray-200 dark:border-gray-700">
                <Calendar size={14} className="text-purple-600 shrink-0" />
                <span>
                    Cada corte agrupa el recaudo de una semana (lunes a domingo). El administrador confirma el corte
                    para cerrar el pago de esa semana.
                    {!isAdmin && ' Solo ves las semanas ya confirmadas.'}
                </span>
                <div className="flex-1" />
                <button
                    onClick={load}
                    disabled={loading}
                    className="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 disabled:opacity-50"
                >
                    <RefreshCw size={13} className={loading ? 'animate-spin' : ''} />
                </button>
            </div>

            {feedback && (
                <div className={`px-3 py-2 rounded-md text-sm ${feedback.ok ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-red-50 text-red-700 border border-red-200'}`}>
                    {feedback.msg}
                </div>
            )}

            {error && (
                <div className="p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm flex items-center gap-2">
                    <AlertCircle size={15} /> {error}
                </div>
            )}

            {loading && (
                <div className="flex items-center justify-center py-16 text-gray-400">
                    <RefreshCw size={18} className="animate-spin mr-2" /> Cargando cortes...
                </div>
            )}

            {!loading && cuts.length === 0 && !error && (
                <div className="text-center py-16 text-gray-400 text-sm">
                    <Calendar size={30} className="mx-auto mb-2 opacity-50" />
                    {isAdmin ? 'No hay cortes de pago todavia.' : 'No hay semanas confirmadas todavia.'}
                </div>
            )}

            {!loading && cuts.map(cut => {
                const key = cut.period_start.slice(0, 10);
                const isOpen = expanded === key;
                const isConfirmed = cut.status === 'confirmed';
                return (
                    <div
                        key={key}
                        className={`rounded-xl border ${isConfirmed ? 'border-emerald-200 dark:border-emerald-800' : 'border-gray-200 dark:border-gray-700'} bg-white dark:bg-gray-800 overflow-hidden`}
                    >
                        <div className="flex items-center gap-3 px-4 py-3 flex-wrap">
                            <div className="flex items-center gap-2">
                                <Calendar size={15} className="text-purple-600" />
                                <span className="text-sm font-bold text-gray-900 dark:text-white">
                                    {periodLabel(cut.period_start, cut.period_end)}
                                </span>
                            </div>
                            {isConfirmed ? (
                                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold bg-emerald-100 text-emerald-700">
                                    <CheckCircle2 size={12} /> Confirmado
                                </span>
                            ) : (
                                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold bg-amber-100 text-amber-700">
                                    <Clock size={12} /> Pendiente
                                </span>
                            )}
                            <div className="flex-1" />
                            <div className="flex items-center gap-4 text-xs">
                                <div className="text-right">
                                    <div className="text-[10px] uppercase text-gray-400 font-bold">Recaudado</div>
                                    <div className="font-bold text-emerald-600">{formatMoney(cut.total_collected)}</div>
                                </div>
                                <div className="text-right">
                                    <div className="text-[10px] uppercase text-gray-400 font-bold">Ordenes</div>
                                    <div className="font-bold text-gray-700 dark:text-gray-200">{cut.orders_count}</div>
                                </div>
                            </div>
                            <button
                                onClick={() => setExpanded(isOpen ? null : key)}
                                className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
                            >
                                {isOpen ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
                            </button>
                            {canConfirm && !isConfirmed && (
                                <button
                                    onClick={() => setConfirmTarget(cut)}
                                    className="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold rounded-md inline-flex items-center gap-1.5"
                                >
                                    <ShieldCheck size={13} /> Confirmar corte
                                </button>
                            )}
                        </div>

                        {isConfirmed && cut.confirmed_by_name && (
                            <div className="px-4 pb-2 text-[11px] text-gray-500 dark:text-gray-400">
                                Confirmado por {cut.confirmed_by_name}
                                {cut.confirmed_at ? ` el ${formatDate(cut.confirmed_at)}` : ''}
                            </div>
                        )}

                        {isOpen && (
                            <div className="border-t border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/40 px-4 py-3">
                                {cut.by_carrier.length === 0 ? (
                                    <div className="text-xs text-gray-400">Sin recaudo en esta semana.</div>
                                ) : (
                                    <table className="w-full text-xs">
                                        <thead>
                                            <tr className="text-[10px] uppercase text-gray-400">
                                                <th className="text-left py-1 font-semibold">Transportadora</th>
                                                <th className="text-right py-1 font-semibold">Ordenes</th>
                                                <th className="text-right py-1 font-semibold">Recaudado</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {cut.by_carrier.map(c => (
                                                <tr key={c.carrier} className="border-t border-gray-100 dark:border-gray-700/50">
                                                    <td className="py-1.5 font-medium text-gray-800 dark:text-gray-200">{carrierLabel(c.carrier)}</td>
                                                    <td className="py-1.5 text-right text-gray-600 dark:text-gray-300">{c.orders_count}</td>
                                                    <td className="py-1.5 text-right text-emerald-600 font-semibold">{formatMoney(c.total_collected)}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                )}
                            </div>
                        )}
                    </div>
                );
            })}

            {confirmTarget && (
                <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-md w-full p-5">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-1 flex items-center gap-2">
                            <ShieldCheck size={18} className="text-emerald-600" /> Confirmar corte de pago
                        </h3>
                        <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                            Vas a confirmar el recaudo de la semana{' '}
                            <strong>{periodLabel(confirmTarget.period_start, confirmTarget.period_end)}</strong>.
                            Una vez confirmado, el negocio podra ver esta semana como cerrada.
                        </p>
                        <div className="bg-gray-50 dark:bg-gray-900/40 rounded-lg p-3 text-sm space-y-1 mb-4">
                            <div className="flex justify-between"><span className="text-gray-700 dark:text-gray-200 font-semibold">Recaudado</span><span className="font-bold text-emerald-600">{formatMoney(confirmTarget.total_collected)}</span></div>
                        </div>
                        <div className="flex justify-end gap-2">
                            <button
                                onClick={() => setConfirmTarget(null)}
                                disabled={confirming}
                                className="px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700"
                            >
                                Cancelar
                            </button>
                            <button
                                onClick={doConfirm}
                                disabled={confirming}
                                className="px-3 py-2 text-sm rounded-md bg-emerald-600 hover:bg-emerald-700 text-white font-semibold inline-flex items-center gap-1.5 disabled:opacity-50"
                            >
                                {confirming ? <RefreshCw size={14} className="animate-spin" /> : <CheckCircle2 size={14} />}
                                Confirmar corte
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
