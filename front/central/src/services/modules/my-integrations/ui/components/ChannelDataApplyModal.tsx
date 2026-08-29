'use client';

import { useCallback, useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { AlertTriangle, ArrowDownToLine, Loader2, X } from 'lucide-react';
import {
    applyChannelData,
    fetchDataPreview,
    type DataApplyResult,
    type DataMode,
    type DataPreview,
    type DataSummaryCell,
} from '../../infra/repository/channel-data';
import { useSSE } from '@/shared/hooks/use-sse';
import { ACCENT, CARD_BG, CARD_BORDER } from '../panel-theme';

const EVENTOS = [
    'products.channel_data.progress',
    'products.channel_data.completed',
    'products.channel_data.failed',
];

export interface ObjetivoDato {
    field: string;
    label: string;
    cell: DataSummaryCell;
    mode: DataMode;
}

export function fechaCorta(value?: string): string {
    if (!value) return '';
    const d = new Date(value);
    if (Number.isNaN(d.getTime())) return '';
    return d.toLocaleString('es-CO', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' });
}

function recorte(value?: string, largo = 70): string {
    if (!value) return '';
    const limpio = value.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim();
    return limpio.length > largo ? `${limpio.slice(0, largo)}...` : limpio;
}

interface ChannelDataApplyModalProps {
    objetivo: ObjetivoDato;
    businessId: number | null;
    onClose: () => void;
    onApplied: (result: DataApplyResult) => void;
}

export function ChannelDataApplyModal({ objetivo, businessId, onClose, onApplied }: ChannelDataApplyModalProps) {
    const [preview, setPreview] = useState<DataPreview | null>(null);
    const [loading, setLoading] = useState(true);
    const [aplicando, setAplicando] = useState(false);
    const [lote, setLote] = useState<string | null>(null);
    const [avance, setAvance] = useState({ processed: 0, total: 0 });
    const [fallo, setFallo] = useState<string | null>(null);

    useEffect(() => {
        const controller = new AbortController();
        setLoading(true);
        fetchDataPreview(businessId ?? undefined, objetivo.cell.integration_id, objetivo.field, objetivo.mode, controller.signal)
            .then(setPreview)
            .finally(() => {
                if (!controller.signal.aborted) setLoading(false);
            });
        return () => controller.abort();
    }, [businessId, objetivo]);

    const alRecibir = useCallback((event: MessageEvent) => {
        let payload: { type?: string; data?: Record<string, unknown> };
        try {
            payload = JSON.parse(event.data);
        } catch {
            return;
        }
        const data = payload?.data ?? {};
        if (!lote || data.batch_id !== lote) return;

        if (payload.type === 'products.channel_data.progress') {
            setAvance({ processed: Number(data.processed) || 0, total: Number(data.total) || 0 });
            return;
        }
        if (payload.type === 'products.channel_data.failed') {
            setAplicando(false);
            setFallo(String(data.error ?? 'No se pudo completar la actualización'));
            return;
        }
        if (payload.type === 'products.channel_data.completed') {
            setAplicando(false);
            onApplied({
                batch_id: lote,
                field: objetivo.field,
                applied: Number(data.applied) || 0,
                applied_at: new Date().toISOString(),
            });
        }
    }, [lote, objetivo.field, onApplied]);

    useSSE({
        businessId: businessId ?? 0,
        eventTypes: EVENTOS,
        onMessage: alRecibir,
        enabled: Boolean(lote),
    });

    const aplicar = async () => {
        setAplicando(true);
        setFallo(null);
        setAvance({ processed: 0, total: preview?.would_fill ?? 0 });
        const result = await applyChannelData(businessId ?? undefined, objetivo.cell.integration_id, objetivo.field, objetivo.mode);
        if (!result?.batch_id) {
            setAplicando(false);
            setFallo('No se pudo iniciar la actualización');
            return;
        }
        setLote(result.batch_id);
    };

    const porcentaje = avance.total > 0 ? Math.round((avance.processed / avance.total) * 100) : 0;

    const total = (preview?.would_fill ?? 0) + (preview?.would_replace ?? 0);

    return createPortal(
        <div className="fixed inset-0 z-[320] flex items-center justify-center p-4" role="dialog" aria-modal="true">
            <button aria-label="Cancelar" onClick={onClose} className="absolute inset-0 cursor-default bg-gray-900/50 backdrop-blur-[2px]" />

            <div className="relative flex max-h-[85vh] w-full max-w-2xl flex-col overflow-hidden rounded-2xl bg-white shadow-2xl dark:bg-gray-900">
                <div className="flex items-start justify-between gap-3 border-b px-5 py-4" style={{ borderColor: CARD_BORDER }}>
                    <div className="min-w-0">
                        <h3 className="text-[15px] font-bold text-gray-900 dark:text-white">
                            Actualizar {objetivo.label.toLowerCase()} en Probability
                        </h3>
                        <p className="mt-0.5 text-[12px] text-gray-500 dark:text-gray-400">
                            Se escribe en los productos de Probability con lo que dice{' '}
                            <span className="font-semibold">{objetivo.cell.integration_name}</span>. El canal no se modifica.
                            {objetivo.mode === 'fill_empty'
                                ? ' Solo se tocan los productos que hoy tienen ese campo vacio.'
                                : ' Se reemplaza el dato actual, incluso donde ya hay algo escrito.'}
                        </p>
                    </div>
                    <button
                        onClick={onClose}
                        className="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-gray-100 text-gray-500 hover:bg-gray-200 dark:bg-gray-800"
                    >
                        <X size={15} />
                    </button>
                </div>

                <div className="min-h-0 flex-1 space-y-3 overflow-y-auto px-5 py-4">
                    {loading && (
                        <p className="flex items-center justify-center gap-2 py-8 text-[12px] text-gray-500">
                            <Loader2 size={14} className="animate-spin" />
                            Revisando que se veria afectado
                        </p>
                    )}

                    {!loading && preview && (
                        <>
                            <div className="grid grid-cols-2 gap-2">
                                <div className="rounded-xl border px-3 py-2" style={{ borderColor: CARD_BORDER, backgroundColor: CARD_BG }}>
                                    <span className="block text-[10.5px] font-semibold uppercase tracking-wide text-gray-500">Se llenan en Probability</span>
                                    <span className="text-[19px] font-black text-emerald-600">{preview.would_fill.toLocaleString('es-CO')}</span>
                                </div>
                                <div className="rounded-xl border px-3 py-2" style={{ borderColor: CARD_BORDER, backgroundColor: CARD_BG }}>
                                    <span className="block text-[10.5px] font-semibold uppercase tracking-wide text-gray-500">Se reemplazan en Probability</span>
                                    <span className={`text-[19px] font-black ${preview.would_replace > 0 ? 'text-amber-600' : 'text-gray-300'}`}>
                                        {preview.would_replace.toLocaleString('es-CO')}
                                    </span>
                                </div>
                            </div>

                            {preview.conflicts.length > 0 && (
                                <div className="rounded-xl border border-amber-300 bg-amber-50 px-3 py-2.5 dark:border-amber-500/40 dark:bg-amber-900/20">
                                    <p className="flex items-center gap-1.5 text-[12px] font-bold text-amber-800 dark:text-amber-200">
                                        <AlertTriangle size={12} />
                                        Ojo, este dato ya lo escribio alguien mas
                                    </p>
                                    <ul className="mt-1.5 space-y-1">
                                        {preview.conflicts.map(conflict => (
                                            <li key={`${conflict.source}-${conflict.integration_id ?? 0}`} className="text-[11.5px] text-amber-800 dark:text-amber-200">
                                                <span className="font-bold">{conflict.count.toLocaleString('es-CO')}</span>{' '}
                                                {conflict.count === 1 ? 'producto viene' : 'productos vienen'} de{' '}
                                                <span className="font-semibold">
                                                    {conflict.source === 'manual' ? 'una edición manual' : conflict.integration_name}
                                                </span>
                                                {conflict.last_change_at && <span className="text-amber-700/80"> · {fechaCorta(conflict.last_change_at)}</span>}
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            )}

                            {preview.samples.length > 0 && (
                                <div className="overflow-hidden rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                                    <table className="w-full border-collapse text-[11.5px]">
                                        <thead style={{ backgroundColor: CARD_BG }}>
                                            <tr className="text-left text-[10px] font-bold uppercase tracking-wider text-gray-400">
                                                <th className="px-3 py-2">SKU</th>
                                                <th className="px-3 py-2">Hoy en Probability</th>
                                                <th className="px-3 py-2">Quedaria</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {preview.samples.map(sample => (
                                                <tr key={sample.product_id} className="border-t" style={{ borderColor: CARD_BORDER }}>
                                                    <td className="px-3 py-1.5 font-mono font-semibold text-gray-700 dark:text-gray-200">{sample.sku}</td>
                                                    <td className="px-3 py-1.5 text-gray-400">
                                                        {sample.current ? recorte(sample.current) : <span className="italic">vacio</span>}
                                                    </td>
                                                    <td className="px-3 py-1.5 font-semibold text-emerald-700 dark:text-emerald-400">{recorte(sample.incoming)}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            )}

                            <p className="text-[11px] text-gray-500 dark:text-gray-400">
                                Queda registrado quien escribio cada dato en Probability y se puede deshacer completo despues de aplicar.
                            </p>
                        </>
                    )}
                </div>

                {(aplicando || fallo) && (
                    <div className="border-t px-5 py-3" style={{ borderColor: CARD_BORDER }}>
                        {fallo ? (
                            <p className="text-[12px] font-semibold text-red-600">{fallo}</p>
                        ) : (
                            <>
                                <div className="mb-1 flex items-center justify-between text-[11.5px] text-gray-500">
                                    <span>
                                        {avance.total > 0
                                            ? `${avance.processed.toLocaleString('es-CO')} de ${avance.total.toLocaleString('es-CO')}`
                                            : 'Preparando...'}
                                    </span>
                                    <span className="font-bold" style={{ color: ACCENT }}>{porcentaje}%</span>
                                </div>
                                <div className="h-1.5 w-full overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700">
                                    <div
                                        className="h-full rounded-full transition-all duration-300"
                                        style={{ width: `${Math.max(porcentaje, 3)}%`, backgroundColor: ACCENT }}
                                    />
                                </div>
                            </>
                        )}
                    </div>
                )}

                <div className="flex items-center justify-end gap-2 border-t px-5 py-3" style={{ borderColor: CARD_BORDER }}>
                    <button onClick={onClose} className="rounded-lg px-3 py-2 text-[12px] font-semibold text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800">
                        Cancelar
                    </button>
                    <button
                        onClick={aplicar}
                        disabled={loading || aplicando || total === 0}
                        className="inline-flex items-center gap-1.5 rounded-lg px-4 py-2 text-[12px] font-bold text-white transition-opacity disabled:opacity-40"
                        style={{ backgroundColor: ACCENT }}
                    >
                        {aplicando ? <Loader2 size={13} className="animate-spin" /> : <ArrowDownToLine size={13} />}
                        Actualizar {total.toLocaleString('es-CO')} en Probability
                    </button>
                </div>
            </div>
        </div>,
        document.body,
    );
}
