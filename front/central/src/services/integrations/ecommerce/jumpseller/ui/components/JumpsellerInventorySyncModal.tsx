'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { X, CheckCircle2, Loader2, AlertCircle, RefreshCw, ArrowUpFromLine } from 'lucide-react';
import { useSSE } from '@/shared/hooks/use-sse';
import { syncJumpsellerInventoryAction } from '../../infra/actions';

interface JumpsellerInventorySyncModalProps {
    isOpen: boolean;
    onClose: () => void;
    integrationId: number;
    businessId: number | null;
    onCompleted?: () => void;
}

const INVENTORY_EVENT_TYPES = [
    'jumpseller.inventory.sync.started',
    'jumpseller.inventory.sync.item',
    'jumpseller.inventory.sync.progress',
    'jumpseller.inventory.sync.completed',
];

type Phase = 'idle' | 'starting' | 'running' | 'done' | 'error';

interface SyncItem {
    sku: string;
    quantity: number;
    action: 'updated' | 'failed';
}

export function JumpsellerInventorySyncModal({ isOpen, onClose, integrationId, businessId, onCompleted }: JumpsellerInventorySyncModalProps) {
    const [phase, setPhase] = useState<Phase>('idle');
    const [total, setTotal] = useState(0);
    const [processed, setProcessed] = useState(0);
    const [updated, setUpdated] = useState(0);
    const [failed, setFailed] = useState(0);
    const [items, setItems] = useState<SyncItem[]>([]);
    const [failedSkus, setFailedSkus] = useState<string[]>([]);
    const [failedHidden, setFailedHidden] = useState(0);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);

    const correlationRef = useRef<string | null>(null);

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data);
            const eventType = parsed.type || parsed.metadata?.event_type;
            const data = parsed.data;
            if (!data) return;
            const corr = correlationRef.current;
            if (!corr || data.correlation_id !== corr) return;

            if (eventType === 'jumpseller.inventory.sync.started') {
                setTotal(Number(data.total) || 0);
                setPhase('running');
            } else if (eventType === 'jumpseller.inventory.sync.item') {
                setItems((prev) => [...prev, {
                    sku: String(data.sku || ''),
                    quantity: Number(data.quantity) || 0,
                    action: data.action === 'failed' ? 'failed' : 'updated',
                }]);
            } else if (eventType === 'jumpseller.inventory.sync.progress') {
                setProcessed(Number(data.processed) || 0);
                setTotal(Number(data.total) || 0);
                setUpdated(Number(data.updated) || 0);
                setFailed(Number(data.failed) || 0);
            } else if (eventType === 'jumpseller.inventory.sync.completed') {
                setProcessed(Number(data.total) || 0);
                setTotal(Number(data.total) || 0);
                setUpdated(Number(data.updated) || 0);
                setFailed(Number(data.failed) || 0);
                setFailedSkus(Array.isArray(data.failed_skus) ? data.failed_skus.map(String) : []);
                setFailedHidden(Number(data.failed_hidden) || 0);
                setPhase('done');
                onCompleted?.();
            }
        } catch {
            return;
        }
    }, [onCompleted]);

    useSSE({
        businessId: businessId ?? 0,
        eventTypes: INVENTORY_EVENT_TYPES,
        onMessage: handleMessage,
        enabled: isOpen,
    });

    useEffect(() => {
        if (!isOpen) {
            setPhase('idle');
            setTotal(0);
            setProcessed(0);
            setUpdated(0);
            setFailed(0);
            setItems([]);
            setFailedSkus([]);
            setFailedHidden(0);
            setErrorMessage(null);
            correlationRef.current = null;
            return;
        }

        let cancelled = false;
        const run = async () => {
            setPhase('starting');
            const result: any = await syncJumpsellerInventoryAction(integrationId, businessId ?? undefined);
            if (cancelled) return;
            if (!result?.success || !result?.correlation_id) {
                setErrorMessage(result?.message || 'No se pudo iniciar la sincronizacion');
                setPhase('error');
                return;
            }
            correlationRef.current = result.correlation_id;
            setPhase('running');
        };
        run();

        return () => { cancelled = true; };
    }, [isOpen, integrationId, businessId]);

    if (!isOpen) return null;

    const progressPct = total > 0 ? Math.min(100, Math.round((processed / total) * 100)) : phase === 'done' ? 100 : 0;
    const busy = phase === 'starting' || phase === 'running';

    const badges: Array<{ label: string; value: number; cls: string }> = [
        { label: 'Actualizados', value: updated, cls: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' },
        { label: 'Fallidos', value: failed, cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300' },
    ];

    return (
        <div className="fixed inset-0 z-[1100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl w-full max-w-lg flex flex-col overflow-hidden border border-gray-200 dark:border-gray-700 max-h-[90vh]">
                <div className="flex items-center justify-between px-6 py-5 border-b border-gray-100 dark:border-gray-800 bg-gradient-to-r from-emerald-50 to-teal-50 dark:from-emerald-950/40 dark:to-teal-950/40">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-emerald-500/10 dark:bg-emerald-400/10 flex items-center justify-center">
                            <RefreshCw size={18} className={`text-emerald-600 dark:text-emerald-400 ${busy ? 'animate-spin' : ''}`} />
                        </div>
                        <div>
                            <h2 className="text-lg font-bold text-gray-900 dark:text-white">Sincronizacion de Inventario</h2>
                            <p className="text-xs text-gray-500 dark:text-gray-400 flex items-center gap-1">
                                <ArrowUpFromLine size={12} /> Probability &rarr; Jumpseller
                                {phase === 'starting' && ' - Iniciando...'}
                                {phase === 'running' && ` - ${processed} de ${total}`}
                                {phase === 'done' && ` - Completado (${processed})`}
                                {phase === 'error' && ' - Error'}
                            </p>
                        </div>
                    </div>
                    <button
                        onClick={onClose}
                        disabled={busy}
                        className="p-2 rounded-lg hover:bg-white/50 dark:hover:bg-white/10 text-gray-500 dark:text-gray-400 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                    >
                        <X size={18} />
                    </button>
                </div>

                {phase === 'error' && (
                    <div className="p-6 flex flex-col items-center text-center gap-3">
                        <AlertCircle size={48} className="text-red-500" />
                        <p className="text-gray-700 dark:text-gray-300 font-medium">{errorMessage}</p>
                        <button
                            onClick={onClose}
                            className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 rounded-lg text-sm font-semibold text-gray-800 dark:text-gray-200 transition-colors"
                        >
                            Cerrar
                        </button>
                    </div>
                )}

                {(phase === 'starting' || phase === 'running' || phase === 'done') && (
                    <div className="px-6 py-5 overflow-y-auto">
                        <div className="flex items-center justify-between mb-2">
                            <span className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Progreso</span>
                            <span className="text-sm font-bold text-gray-900 dark:text-white tabular-nums">{progressPct}%</span>
                        </div>
                        <div className="h-2 bg-gray-100 dark:bg-gray-800 rounded-full overflow-hidden">
                            <div
                                className="h-full rounded-full bg-gradient-to-r from-emerald-500 to-teal-500 transition-all duration-300 ease-out"
                                style={{ width: `${progressPct}%` }}
                            />
                        </div>

                        <div className="grid grid-cols-2 gap-2 mt-4">
                            {badges.map((b) => (
                                <div key={b.label} className={`flex items-center justify-between rounded-lg px-3 py-2 text-xs font-semibold ${b.cls}`}>
                                    <span>{b.label}</span>
                                    <span className="tabular-nums">{b.value}</span>
                                </div>
                            ))}
                        </div>

                        {items.length > 0 && (
                            <div className="mt-4">
                                <style>{`@keyframes jsInvRowIn{from{opacity:0;transform:translateY(4px)}to{opacity:1;transform:none}}`}</style>
                                <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1.5">Detalle por producto</p>
                                <div className="max-h-56 overflow-y-auto rounded-lg border border-gray-100 dark:border-gray-800 divide-y divide-gray-100 dark:divide-gray-800">
                                    {items.slice(-20).reverse().map((it, i) => (
                                        <div key={`${it.sku}-${i}`} style={{ animation: 'jsInvRowIn 0.25s ease' }} className="flex items-center justify-between gap-2 px-2.5 py-1.5 text-[11px]">
                                            <span className="text-gray-500 dark:text-gray-400 font-mono truncate">{it.sku || '(sin sku)'}</span>
                                            <div className="flex items-center gap-2 flex-shrink-0">
                                                <span className="tabular-nums text-gray-700 dark:text-gray-200">{it.quantity} u</span>
                                                <ActionBadge action={it.action} />
                                            </div>
                                        </div>
                                    ))}
                                </div>
                                {items.length > 20 && (
                                    <p className="mt-1 text-[11px] text-gray-400">Mostrando los ultimos 20 de {items.length}.</p>
                                )}
                            </div>
                        )}

                        {failed > 0 && failedSkus.length > 0 && (
                            <div className="mt-4 rounded-lg border border-amber-200 dark:border-amber-800 bg-amber-50/40 dark:bg-amber-900/10 p-3">
                                <p className="text-xs font-semibold text-amber-800 dark:text-amber-300">SKUs que fallaron</p>
                                <p className="text-[11px] text-gray-400 mt-0.5">Puedes reintentarlos volviendo a ejecutar la operacion.</p>
                                <div className="mt-2 max-h-32 overflow-y-auto rounded-md bg-white dark:bg-gray-800/60 border border-amber-100 dark:border-amber-900/40 divide-y divide-gray-100 dark:divide-gray-700">
                                    {failedSkus.map((sku, i) => (
                                        <div key={i} className="px-2.5 py-1.5 font-mono text-[11px] text-gray-600 dark:text-gray-300 truncate">{sku || '(sin sku)'}</div>
                                    ))}
                                </div>
                                {failedHidden > 0 && (
                                    <p className="mt-2 text-[11px] font-medium text-amber-700 dark:text-amber-400">y {failedHidden} mas...</p>
                                )}
                            </div>
                        )}

                        {phase === 'starting' && (
                            <div className="flex items-center justify-center py-6 gap-2 text-gray-500 dark:text-gray-400">
                                <Loader2 size={18} className="animate-spin" />
                                <span className="text-sm">Consultando inventario en Probability...</span>
                            </div>
                        )}

                        {phase === 'done' && (
                            <div className="flex justify-end mt-5">
                                <button
                                    onClick={onClose}
                                    className="px-5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-sm font-semibold transition-colors flex items-center gap-2"
                                >
                                    <CheckCircle2 size={16} /> Listo
                                </button>
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}

function ActionBadge({ action }: { action: SyncItem['action'] }) {
    const map = {
        updated: { label: 'Actualizado', cls: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' },
        failed: { label: 'Fallido', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300' },
    };
    const { label, cls } = map[action];
    return <span className={`px-1.5 py-0.5 rounded font-semibold ${cls}`}>{label}</span>;
}
