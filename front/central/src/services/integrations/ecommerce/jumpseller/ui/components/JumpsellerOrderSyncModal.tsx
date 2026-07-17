'use client';

import { useCallback, useEffect, useState } from 'react';
import { X, CheckCircle2, Loader2, AlertCircle, ShoppingBag, ArrowDownToLine } from 'lucide-react';
import { useSSE } from '@/shared/hooks/use-sse';
import { syncOrdersAction } from '@/services/integrations/core/infra/actions';
import { SyncOrdersParams } from '@/services/integrations/core/domain/types';

interface JumpsellerOrderSyncModalProps {
    isOpen: boolean;
    onClose: () => void;
    integrationId: number;
    businessId: number | null;
    createdAtMin?: string;
    createdAtMax?: string;
    onCompleted?: () => void;
}

const ORDER_EVENT_TYPES = [
    'jumpseller.orders.sync.started',
    'jumpseller.orders.sync.item',
    'jumpseller.orders.sync.progress',
    'jumpseller.orders.sync.completed',
    'jumpseller.orders.sync.failed',
];

type Phase = 'idle' | 'starting' | 'running' | 'done' | 'error';

interface OrderItem {
    orderNumber: string;
    externalId: string;
    customerName: string;
    total: number;
    currency: string;
    status: string;
    originalStatus: string;
    items: number;
    action: 'imported' | 'failed';
}

const formatAmount = (value: number, currency: string) => {
    const amount = Number.isFinite(value) ? value.toLocaleString('es-CO', { maximumFractionDigits: 2 }) : '0';
    return currency ? `${amount} ${currency}` : amount;
};

export function JumpsellerOrderSyncModal({ isOpen, onClose, integrationId, businessId, createdAtMin, createdAtMax, onCompleted }: JumpsellerOrderSyncModalProps) {
    const [phase, setPhase] = useState<Phase>('idle');
    const [page, setPage] = useState(0);
    const [processed, setProcessed] = useState(0);
    const [imported, setImported] = useState(0);
    const [failed, setFailed] = useState(0);
    const [duration, setDuration] = useState<string | null>(null);
    const [orders, setOrders] = useState<OrderItem[]>([]);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data);
            const eventType: string = parsed.type || parsed.metadata?.event_type || '';
            const data = parsed.data;
            if (!data) return;
            if (!eventType.startsWith('jumpseller.orders.sync.')) return;

            if (eventType === 'jumpseller.orders.sync.started') {
                setPhase('running');
            } else if (eventType === 'jumpseller.orders.sync.item') {
                setOrders((prev) => [...prev, {
                    orderNumber: String(data.order_number || ''),
                    externalId: String(data.external_id || ''),
                    customerName: String(data.customer_name || ''),
                    total: Number(data.total) || 0,
                    currency: String(data.currency || ''),
                    status: String(data.status || ''),
                    originalStatus: String(data.original_status || ''),
                    items: Number(data.items) || 0,
                    action: data.action === 'failed' ? 'failed' : 'imported',
                }]);
            } else if (eventType === 'jumpseller.orders.sync.progress') {
                setPage(Number(data.page) || 0);
                setProcessed(Number(data.processed) || 0);
                setImported(Number(data.imported) || 0);
                setFailed(Number(data.failed) || 0);
            } else if (eventType === 'jumpseller.orders.sync.completed') {
                const fetched = Number(data.total_fetched) || 0;
                setImported(Number(data.imported) || 0);
                setFailed(Number(data.failed) || 0);
                setProcessed((Number(data.imported) || 0) + (Number(data.failed) || 0) || fetched);
                setDuration(data.duration ? String(data.duration) : null);
                setPhase('done');
                onCompleted?.();
            } else if (eventType === 'jumpseller.orders.sync.failed') {
                setErrorMessage(String(data.error || 'La sincronizacion de ordenes fallo'));
                setPhase('error');
            }
        } catch {
            return;
        }
    }, [onCompleted]);

    useSSE({
        businessId: businessId ?? 0,
        integrationId,
        eventTypes: ORDER_EVENT_TYPES,
        onMessage: handleMessage,
        enabled: isOpen,
    });

    useEffect(() => {
        if (!isOpen) {
            setPhase('idle');
            setPage(0);
            setProcessed(0);
            setImported(0);
            setFailed(0);
            setDuration(null);
            setOrders([]);
            setErrorMessage(null);
            return;
        }

        let cancelled = false;
        const run = async () => {
            setPhase('starting');
            const params: SyncOrdersParams = {};
            if (createdAtMin) params.created_at_min = createdAtMin;
            if (createdAtMax) params.created_at_max = createdAtMax;
            const result: any = await syncOrdersAction(integrationId, Object.keys(params).length > 0 ? params : undefined);
            if (cancelled) return;
            if (!result?.success) {
                setErrorMessage(result?.message || 'No se pudo iniciar la sincronizacion de ordenes');
                setPhase('error');
                return;
            }
            setPhase('running');
        };
        run();

        return () => { cancelled = true; };
    }, [isOpen, integrationId, createdAtMin, createdAtMax]);

    if (!isOpen) return null;

    const busy = phase === 'starting' || phase === 'running';

    const badges: Array<{ label: string; value: number; cls: string }> = [
        { label: 'Importadas', value: imported, cls: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' },
        { label: 'Fallidas', value: failed, cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300' },
    ];

    return (
        <div className="fixed inset-0 z-[1100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl w-full max-w-lg flex flex-col overflow-hidden border border-gray-200 dark:border-gray-700 max-h-[90vh]">
                <div className="flex items-center justify-between px-6 py-5 border-b border-gray-100 dark:border-gray-800 bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-950/40 dark:to-indigo-950/40">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-blue-500/10 dark:bg-blue-400/10 flex items-center justify-center">
                            <ShoppingBag size={18} className={`text-blue-600 dark:text-blue-400 ${busy ? 'animate-pulse' : ''}`} />
                        </div>
                        <div>
                            <h2 className="text-lg font-bold text-gray-900 dark:text-white">Sincronizacion de Ordenes</h2>
                            <p className="text-xs text-gray-500 dark:text-gray-400 flex items-center gap-1">
                                <ArrowDownToLine size={12} /> Jumpseller &rarr; Probability
                                {phase === 'starting' && ' - Iniciando...'}
                                {phase === 'running' && ` - ${processed} procesadas${page > 0 ? ` (pagina ${page})` : ''}`}
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
                        {(createdAtMin || createdAtMax) && (
                            <p className="text-[11px] text-gray-400 mb-3">
                                Periodo: {createdAtMin || 'inicio'} a {createdAtMax || 'hoy'}
                            </p>
                        )}

                        <style>{`@keyframes jsOrdRowIn{from{opacity:0;transform:translateY(4px)}to{opacity:1;transform:none}}@keyframes jsOrdBar{0%{transform:translateX(-100%)}100%{transform:translateX(400%)}}`}</style>

                        <div className="flex items-center justify-between mb-2">
                            <span className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Progreso</span>
                            <span className="text-sm font-bold text-gray-900 dark:text-white tabular-nums">
                                {phase === 'done' ? `${processed} ordenes` : `${processed} procesadas`}
                            </span>
                        </div>
                        <div className="h-2 bg-gray-100 dark:bg-gray-800 rounded-full overflow-hidden">
                            {phase === 'done' ? (
                                <div className="h-full w-full rounded-full bg-gradient-to-r from-blue-500 to-indigo-500" />
                            ) : (
                                <div className="h-full w-1/4 rounded-full bg-gradient-to-r from-blue-400 to-indigo-400" style={{ animation: 'jsOrdBar 1.2s ease-in-out infinite' }} />
                            )}
                        </div>
                        <p className="mt-1 text-[11px] text-gray-400">
                            {phase === 'done'
                                ? `Sincronizacion finalizada${duration ? ` en ${duration}` : ''}.`
                                : 'Jumpseller entrega las ordenes por paginas: el total se conoce al finalizar.'}
                        </p>

                        <div className="grid grid-cols-2 gap-2 mt-4">
                            {badges.map((b) => (
                                <div key={b.label} className={`flex items-center justify-between rounded-lg px-3 py-2 text-xs font-semibold ${b.cls}`}>
                                    <span>{b.label}</span>
                                    <span className="tabular-nums">{b.value}</span>
                                </div>
                            ))}
                        </div>

                        {orders.length > 0 && (
                            <div className="mt-4">
                                <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1.5">Detalle por orden</p>
                                <div className="max-h-64 overflow-y-auto rounded-lg border border-gray-100 dark:border-gray-800 divide-y divide-gray-100 dark:divide-gray-800">
                                    {orders.slice(-20).reverse().map((o, i) => (
                                        <div key={`${o.externalId || o.orderNumber}-${i}`} style={{ animation: 'jsOrdRowIn 0.25s ease' }} className="px-2.5 py-2 text-[11px]">
                                            <div className="flex items-center justify-between gap-2">
                                                <span className="font-mono font-semibold text-gray-700 dark:text-gray-200 truncate">#{o.orderNumber || '(sin numero)'}</span>
                                                <div className="flex items-center gap-2 flex-shrink-0">
                                                    <span className="tabular-nums text-gray-700 dark:text-gray-200">{formatAmount(o.total, o.currency)}</span>
                                                    <ActionBadge action={o.action} />
                                                </div>
                                            </div>
                                            <div className="mt-0.5 flex items-center justify-between gap-2">
                                                <span className="text-gray-500 dark:text-gray-400 truncate">{o.customerName || '(sin cliente)'}</span>
                                                {(o.originalStatus || o.status) && (
                                                    <span className="flex items-center gap-1 flex-shrink-0 text-gray-400">
                                                        <span className="font-mono">{o.originalStatus || '-'}</span>
                                                        <span>&rarr;</span>
                                                        <span className="font-mono font-semibold text-gray-600 dark:text-gray-300">{o.status || '-'}</span>
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                                {orders.length > 20 && (
                                    <p className="mt-1 text-[11px] text-gray-400">Mostrando las ultimas 20 de {orders.length}.</p>
                                )}
                            </div>
                        )}

                        {phase === 'starting' && (
                            <div className="flex items-center justify-center py-6 gap-2 text-gray-500 dark:text-gray-400">
                                <Loader2 size={18} className="animate-spin" />
                                <span className="text-sm">Consultando ordenes en Jumpseller...</span>
                            </div>
                        )}

                        {phase === 'done' && (
                            <div className="flex justify-end mt-5">
                                <button
                                    onClick={onClose}
                                    className="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-semibold transition-colors flex items-center gap-2"
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

function ActionBadge({ action }: { action: OrderItem['action'] }) {
    const map = {
        imported: { label: 'Importada', cls: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' },
        failed: { label: 'Fallida', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300' },
    };
    const { label, cls } = map[action];
    return <span className={`px-1.5 py-0.5 rounded font-semibold ${cls}`}>{label}</span>;
}
