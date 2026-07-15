'use client';

import { useCallback, useEffect, useState } from 'react';
import {
    Calendar, CheckCircle2, ChevronDown, ChevronUp, RefreshCw, AlertCircle, ShieldCheck, Clock, Plus, Trash2,
} from 'lucide-react';
import { getCodCutsAction, deleteCodCutAction, confirmCodCutAction } from '../../infra/actions';
import { PaymentCut } from '../../domain/types';
import { formatMoney, formatDateTime, formatDateOnly } from './helpers';
import { CutSelectionModal } from './CutSelectionModal';
import { CutOrdersDetail } from './CutOrdersDetail';

interface Props {
    businessId?: number | null;
    isAdmin: boolean;
}

function periodLabel(start: string, end: string): string {
    return `${formatDateOnly(start)} - ${formatDateOnly(end)}`;
}

function weekBounds(dateStr: string): { start: string; end: string } {
    const [y, m, d] = dateStr.split('-').map(Number);
    const dt = new Date(Date.UTC(y, m - 1, d));
    const diff = (dt.getUTCDay() + 6) % 7;
    const mon = new Date(dt);
    mon.setUTCDate(dt.getUTCDate() - diff);
    const sun = new Date(mon);
    sun.setUTCDate(mon.getUTCDate() + 6);
    return { start: mon.toISOString().slice(0, 10), end: sun.toISOString().slice(0, 10) };
}

function todayStr(): string {
    return new Date().toLocaleDateString('en-CA');
}

export default function CodCutsTab({ businessId, isAdmin }: Props) {
    const [cuts, setCuts] = useState<PaymentCut[]>([]);
    const [canConfirm, setCanConfirm] = useState(false);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [expanded, setExpanded] = useState<string | null>(null);
    const [feedback, setFeedback] = useState<{ ok: boolean; msg: string } | null>(null);
    const [createOpen, setCreateOpen] = useState(false);
    const [createDate, setCreateDate] = useState('');
    const [selectionPeriod, setSelectionPeriod] = useState<{ start: string; end: string } | null>(null);
    const [deleteTarget, setDeleteTarget] = useState<PaymentCut | null>(null);
    const [deleting, setDeleting] = useState(false);
    const [confirmingId, setConfirmingId] = useState<number | null>(null);

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

    const openSelection = (start: string, end: string) => {
        setSelectionPeriod({ start, end });
    };

    const onCutConfirmed = async (msg: string) => {
        setFeedback({ ok: true, msg });
        await load();
        setTimeout(() => setFeedback(null), 3500);
    };

    const continueCreate = () => {
        if (!createDate) return;
        const { start, end } = weekBounds(createDate);
        setCreateOpen(false);
        openSelection(start, end);
    };

    const confirmDraft = async (cut: PaymentCut) => {
        if (!cut.id) return;
        setConfirmingId(cut.id);
        const res = await confirmCodCutAction(cut.id, businessId || undefined);
        if (res.success) {
            setFeedback({ ok: true, msg: 'Corte confirmado y consignado' });
            await load();
        } else {
            setFeedback({ ok: false, msg: (res as any).message || 'Error al confirmar el corte' });
        }
        setConfirmingId(null);
        setTimeout(() => setFeedback(null), 3500);
    };

    const doDelete = async () => {
        if (!deleteTarget?.id) return;
        setDeleting(true);
        const res = await deleteCodCutAction(deleteTarget.id, businessId || undefined);
        if (res.success) {
            setFeedback({ ok: true, msg: 'Corte de pago eliminado' });
            setDeleteTarget(null);
            await load();
        } else {
            setFeedback({ ok: false, msg: (res as any).message || 'Error al eliminar el corte' });
        }
        setDeleting(false);
        setTimeout(() => setFeedback(null), 3500);
    };

    return (
        <div className="space-y-3">
            <div className="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800 rounded-lg px-3 py-2 border border-gray-200 dark:border-gray-700">
                <Calendar size={14} className="text-purple-600 shrink-0" />
                <span>
                    Marca un corte para elegir que ordenes consignar: se crea un borrador que confirmas (cierra el pago)
                    o cancelas (libera las ordenes).
                    {!isAdmin && ' Solo ves los cortes ya confirmados.'}
                </span>
                <div className="flex-1" />
                {canConfirm && (
                    <button
                        onClick={() => { setCreateDate(todayStr()); setCreateOpen(true); }}
                        className="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold rounded-md inline-flex items-center gap-1 shrink-0"
                    >
                        <Plus size={13} /> Marcar corte de pago
                    </button>
                )}
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
                                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold bg-blue-100 text-blue-700">
                                    <Clock size={12} /> Borrador
                                </span>
                            )}
                            <div className="flex-1" />
                            <div className="flex items-center gap-4 text-xs">
                                <div className="text-right">
                                    <div className="text-[10px] uppercase text-gray-400 font-bold">{isConfirmed ? 'Consignado' : 'Por consignar'}</div>
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
                            {canConfirm && !isConfirmed && cut.id ? (
                                <>
                                    <button
                                        onClick={() => confirmDraft(cut)}
                                        disabled={confirmingId === cut.id}
                                        className="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold rounded-md inline-flex items-center gap-1.5 disabled:opacity-50"
                                    >
                                        {confirmingId === cut.id ? <RefreshCw size={13} className="animate-spin" /> : <ShieldCheck size={13} />} Confirmar corte
                                    </button>
                                    <button
                                        onClick={() => setDeleteTarget(cut)}
                                        className="px-3 py-1.5 border border-gray-200 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 text-xs font-semibold rounded-md inline-flex items-center gap-1.5"
                                        title="Cancelar borrador (libera sus ordenes)"
                                    >
                                        <Trash2 size={13} /> Cancelar
                                    </button>
                                </>
                            ) : null}
                            {canConfirm && isConfirmed && cut.id ? (
                                <button
                                    onClick={() => setDeleteTarget(cut)}
                                    className="p-1.5 rounded-md border border-red-200 dark:border-red-800 text-red-600 hover:bg-red-50 dark:hover:bg-red-900/30"
                                    title="Eliminar corte (libera sus ordenes)"
                                >
                                    <Trash2 size={14} />
                                </button>
                            ) : null}
                        </div>

                        {isConfirmed && cut.confirmed_by_name && (
                            <div className="px-4 pb-2 flex items-center gap-2 text-[11px] text-gray-500 dark:text-gray-400">
                                {cut.confirmed_by_avatar ? (
                                    <img src={cut.confirmed_by_avatar} alt={cut.confirmed_by_name} className="w-5 h-5 rounded-full object-cover border border-gray-200 dark:border-gray-600" />
                                ) : (
                                    <span className="w-5 h-5 rounded-full bg-emerald-500 text-white inline-flex items-center justify-center text-[9px] font-bold">
                                        {cut.confirmed_by_name.charAt(0).toUpperCase()}
                                    </span>
                                )}
                                <span>
                                    Confirmado por <b className="text-gray-700 dark:text-gray-200">{cut.confirmed_by_name}</b>
                                    {cut.confirmed_at ? ` el ${formatDateTime(cut.confirmed_at)}` : ''}
                                </span>
                            </div>
                        )}

                        {isOpen && cut.id && (
                            <div className="border-t border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/40 px-4 py-3">
                                {!isConfirmed && (
                                    <div className="text-[11px] text-gray-500 dark:text-gray-400 mb-2">
                                        Ordenes en este borrador. Confirma el corte para consignarlas, o cancelalo para liberarlas.
                                    </div>
                                )}
                                <CutOrdersDetail cutId={cut.id} businessId={businessId} />
                            </div>
                        )}
                    </div>
                );
            })}

            {selectionPeriod && (
                <CutSelectionModal
                    isOpen={!!selectionPeriod}
                    onClose={() => setSelectionPeriod(null)}
                    onConfirmed={onCutConfirmed}
                    periodStart={selectionPeriod.start}
                    periodEnd={selectionPeriod.end}
                    businessId={businessId}
                />
            )}

            {deleteTarget && (
                <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-md w-full p-5">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-1 flex items-center gap-2">
                            <Trash2 size={18} className="text-red-600" /> {deleteTarget.status === 'confirmed' ? 'Eliminar corte de pago' : 'Cancelar borrador'}
                        </h3>
                        <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                            {deleteTarget.status === 'confirmed' ? 'Vas a eliminar el corte confirmado de la semana ' : 'Vas a cancelar el borrador de la semana '}
                            <strong>{periodLabel(deleteTarget.period_start, deleteTarget.period_end)}</strong>.
                            Sus <strong>{deleteTarget.orders_count}</strong> ordenes volveran a quedar disponibles
                            (podras consignarlas en otro corte). Esta accion no se puede deshacer.
                        </p>
                        <div className="flex justify-end gap-2">
                            <button
                                onClick={() => setDeleteTarget(null)}
                                disabled={deleting}
                                className="px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                            >
                                Cancelar
                            </button>
                            <button
                                onClick={doDelete}
                                disabled={deleting}
                                className="px-3 py-2 text-sm rounded-md bg-red-600 hover:bg-red-700 text-white font-semibold inline-flex items-center gap-1.5 disabled:opacity-50"
                            >
                                {deleting ? <RefreshCw size={14} className="animate-spin" /> : <Trash2 size={14} />}
                                {deleteTarget.status === 'confirmed' ? 'Eliminar corte' : 'Cancelar borrador'}
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {createOpen && (
                <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-md w-full p-5">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-1 flex items-center gap-2">
                            <Plus size={18} className="text-emerald-600" /> Marcar corte de pago
                        </h3>
                        <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                            Selecciona cualquier dia de la semana que quieres cerrar. Luego eliges orden por orden
                            cuales se pagaron al cliente (lunes a domingo).
                        </p>
                        <label className="block text-xs font-semibold text-gray-500 dark:text-gray-400 mb-1">
                            Dia de la semana
                        </label>
                        <input
                            type="date"
                            value={createDate}
                            onChange={e => setCreateDate(e.target.value)}
                            className="w-full px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white mb-3"
                        />
                        {createDate && (
                            <div className="bg-gray-50 dark:bg-gray-900/40 rounded-lg p-3 text-sm mb-4 flex items-center gap-2">
                                <Calendar size={14} className="text-purple-600 shrink-0" />
                                <span className="text-gray-700 dark:text-gray-200">
                                    Semana: <strong>{periodLabel(weekBounds(createDate).start, weekBounds(createDate).end)}</strong>
                                </span>
                            </div>
                        )}
                        <div className="flex justify-end gap-2">
                            <button
                                onClick={() => setCreateOpen(false)}
                                className="px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700"
                            >
                                Cancelar
                            </button>
                            <button
                                onClick={continueCreate}
                                disabled={!createDate}
                                className="px-3 py-2 text-sm rounded-md bg-emerald-600 hover:bg-emerald-700 text-white font-semibold inline-flex items-center gap-1.5 disabled:opacity-50"
                            >
                                <CheckCircle2 size={14} />
                                Continuar
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
