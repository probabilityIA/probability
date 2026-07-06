'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { X, CheckCircle2, Loader2, AlertCircle, RefreshCw, ArrowUpFromLine, ArrowDownToLine, ArrowRightLeft } from 'lucide-react';
import { useSSE } from '@/shared/hooks/use-sse';
import { reconcileWooProductsAction, applyWooProductsAction, syncWooProductsAction } from '../../infra/actions';

interface WooProductSyncModalProps {
    isOpen: boolean;
    onClose: () => void;
    integrationId: number;
    businessId: number | null;
    onCompleted?: () => void;
}

interface Brief {
    sku: string;
    name: string;
}

interface Diff {
    matched: number;
    onlyInProbability: Brief[];
    onlyInWoo: Brief[];
    probabilityNoSku: number;
    woocommerceNoSku: number;
}

const PRODUCT_EVENT_TYPES = [
    'woocommerce.product.sync.started',
    'woocommerce.product.sync.progress',
    'woocommerce.product.sync.item',
    'woocommerce.product.sync.completed',
];

interface SyncItem {
    sku: string;
    name: string;
    quantity: number;
    action: 'created' | 'updated' | 'failed';
}

type Phase = 'analyzing' | 'diff' | 'running' | 'done' | 'error';
type Direction = 'to_woo' | 'to_probability';

export function WooProductSyncModal({ isOpen, onClose, integrationId, businessId, onCompleted }: WooProductSyncModalProps) {
    const [phase, setPhase] = useState<Phase>('analyzing');
    const [diff, setDiff] = useState<Diff | null>(null);
    const [direction, setDirection] = useState<Direction | null>(null);
    const [total, setTotal] = useState(0);
    const [processed, setProcessed] = useState(0);
    const [created, setCreated] = useState(0);
    const [updated, setUpdated] = useState(0);
    const [failed, setFailed] = useState(0);
    const [isFullSync, setIsFullSync] = useState(false);
    const [items, setItems] = useState<SyncItem[]>([]);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);

    const correlationRef = useRef<string | null>(null);

    const analyze = useCallback(async () => {
        setPhase('analyzing');
        setErrorMessage(null);
        const res: any = await reconcileWooProductsAction(integrationId, businessId ?? undefined);
        if (!res?.success) {
            setErrorMessage(res?.message || 'No se pudo analizar los productos');
            setPhase('error');
            return;
        }
        setDiff({
            matched: Number(res.matched) || 0,
            onlyInProbability: res.only_in_probability || [],
            onlyInWoo: res.only_in_woocommerce || [],
            probabilityNoSku: Number(res.probability_no_sku) || 0,
            woocommerceNoSku: Number(res.woocommerce_no_sku) || 0,
        });
        setPhase('diff');
    }, [integrationId, businessId]);

    useEffect(() => {
        if (!isOpen) {
            setPhase('analyzing');
            setDiff(null);
            setDirection(null);
            setTotal(0);
            setProcessed(0);
            setCreated(0);
            setUpdated(0);
            setFailed(0);
            setIsFullSync(false);
            setItems([]);
            setErrorMessage(null);
            correlationRef.current = null;
            return;
        }
        analyze();
    }, [isOpen, analyze]);

    const handleApply = async (dir: Direction) => {
        setDirection(dir);
        setIsFullSync(false);
        setPhase('running');
        setTotal(dir === 'to_woo' ? (diff?.onlyInProbability.length || 0) : (diff?.onlyInWoo.length || 0));
        setProcessed(0);
        setCreated(0);
        setUpdated(0);
        setFailed(0);
        setItems([]);
        correlationRef.current = null;
        const res: any = await applyWooProductsAction(integrationId, dir, businessId ?? undefined);
        if (!res?.success || !res?.correlation_id) {
            setErrorMessage(res?.message || 'No se pudo iniciar la operacion');
            setPhase('error');
            return;
        }
        correlationRef.current = res.correlation_id;
    };

    const handleFullSync = async () => {
        setDirection('to_woo');
        setIsFullSync(true);
        setPhase('running');
        setTotal(diff?.matched || 0);
        setProcessed(0);
        setCreated(0);
        setUpdated(0);
        setFailed(0);
        setItems([]);
        correlationRef.current = null;
        const res: any = await syncWooProductsAction(integrationId, businessId ?? undefined);
        if (!res?.success || !res?.correlation_id) {
            setErrorMessage(res?.message || 'No se pudo iniciar la sincronizacion');
            setPhase('error');
            return;
        }
        correlationRef.current = res.correlation_id;
    };

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data);
            const eventType = parsed.type || parsed.metadata?.event_type;
            const data = parsed.data;
            if (!data) return;
            const corr = correlationRef.current;
            if (!corr || data.correlation_id !== corr) return;

            if (eventType === 'woocommerce.product.sync.started') {
                setTotal(Number(data.total) || 0);
            } else if (eventType === 'woocommerce.product.sync.item') {
                setItems((prev) => [...prev, {
                    sku: String(data.sku || ''),
                    name: String(data.name || ''),
                    quantity: Number(data.quantity) || 0,
                    action: (data.action === 'created' || data.action === 'failed') ? data.action : 'updated',
                }]);
            } else if (eventType === 'woocommerce.product.sync.progress') {
                setProcessed(Number(data.processed) || 0);
                setCreated(Number(data.created) || 0);
                setUpdated(Number(data.updated) || 0);
                setFailed(Number(data.failed) || 0);
            } else if (eventType === 'woocommerce.product.sync.completed') {
                setProcessed(Number(data.total) || 0);
                setTotal(Number(data.total) || 0);
                setCreated(Number(data.created) || 0);
                setUpdated(Number(data.updated) || 0);
                setFailed(Number(data.failed) || 0);
                setPhase('done');
                onCompleted?.();
            }
        } catch {
            return;
        }
    }, [onCompleted]);

    useSSE({
        businessId: businessId ?? 0,
        eventTypes: PRODUCT_EVENT_TYPES,
        onMessage: handleMessage,
        enabled: isOpen && (phase === 'running' || phase === 'done'),
    });

    if (!isOpen) return null;

    const progressPct = total > 0 ? Math.min(100, Math.round((processed / total) * 100)) : phase === 'done' ? 100 : 0;
    const inSync = diff && diff.onlyInProbability.length === 0 && diff.onlyInWoo.length === 0;

    return (
        <div className="fixed inset-0 z-[1100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl w-full max-w-lg flex flex-col overflow-hidden border border-gray-200 dark:border-gray-700 max-h-[90vh]">
                <div className="flex items-center justify-between px-6 py-5 border-b border-gray-100 dark:border-gray-800 bg-gradient-to-r from-violet-50 to-purple-50 dark:from-violet-950/40 dark:to-purple-950/40">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-violet-500/10 dark:bg-violet-400/10 flex items-center justify-center">
                            <ArrowRightLeft size={18} className="text-violet-600 dark:text-violet-400" />
                        </div>
                        <div>
                            <h2 className="text-lg font-bold text-gray-900 dark:text-white">Sincronizacion de Productos</h2>
                            <p className="text-xs text-gray-500 dark:text-gray-400">Probability &harr; WooCommerce</p>
                        </div>
                    </div>
                    <button onClick={onClose} disabled={phase === 'running'} className="p-2 rounded-lg hover:bg-white/50 dark:hover:bg-white/10 text-gray-500 dark:text-gray-400 disabled:opacity-40 disabled:cursor-not-allowed transition-colors">
                        <X size={18} />
                    </button>
                </div>

                <div className="p-6 overflow-y-auto">
                    {phase === 'analyzing' && (
                        <div className="flex items-center justify-center py-10 gap-2 text-gray-500 dark:text-gray-400">
                            <Loader2 size={20} className="animate-spin" />
                            <span className="text-sm">Analizando productos en ambos lados...</span>
                        </div>
                    )}

                    {phase === 'error' && (
                        <div className="flex flex-col items-center text-center gap-3 py-6">
                            <AlertCircle size={48} className="text-red-500" />
                            <p className="text-gray-700 dark:text-gray-300 font-medium">{errorMessage}</p>
                            <button onClick={analyze} className="px-4 py-2 bg-violet-600 hover:bg-violet-700 text-white rounded-lg text-sm font-semibold transition-colors">Reintentar</button>
                        </div>
                    )}

                    {phase === 'diff' && diff && (
                        <div className="space-y-4">
                            <div className="flex items-center gap-2 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 px-3 py-2">
                                <CheckCircle2 size={16} className="text-emerald-600 dark:text-emerald-400" />
                                <span className="text-sm text-emerald-800 dark:text-emerald-300"><strong>{diff.matched}</strong> productos coinciden en ambos lados</span>
                            </div>

                            <div className="rounded-lg border border-violet-200 dark:border-violet-800 bg-violet-50/50 dark:bg-violet-900/10 p-3 flex items-start justify-between gap-3">
                                <div>
                                    <p className="text-sm font-semibold text-gray-900 dark:text-gray-100">Sincronizar stock a WooCommerce</p>
                                    <p className="text-[11px] text-gray-400 mt-0.5">Vincula por SKU los productos que ya existen y actualiza su stock en WooCommerce.</p>
                                </div>
                                <button onClick={handleFullSync} className="inline-flex items-center gap-1.5 whitespace-nowrap rounded-lg bg-violet-600 hover:bg-violet-700 px-3 py-1.5 text-xs font-semibold text-white transition-colors">
                                    <RefreshCw size={14} /> Sincronizar stock
                                </button>
                            </div>

                            {inSync ? (
                                <div className="text-center py-6 text-gray-600 dark:text-gray-300">
                                    <CheckCircle2 size={40} className="mx-auto text-emerald-500 mb-2" />
                                    <p className="font-semibold">Todo sincronizado</p>
                                    <p className="text-xs text-gray-400 mt-1">No hay productos pendientes en ninguno de los dos lados.</p>
                                </div>
                            ) : (
                                <>
                                    {diff.onlyInProbability.length > 0 && (
                                        <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
                                            <div className="flex items-start justify-between gap-3">
                                                <div>
                                                    <p className="text-sm font-semibold text-gray-900 dark:text-gray-100">En Probability hay {diff.onlyInProbability.length} producto{diff.onlyInProbability.length !== 1 ? 's' : ''} que no estan en WooCommerce</p>
                                                    <p className="text-[11px] text-gray-400 mt-0.5">Se crearan en tu tienda WooCommerce (con imagen si tienen).</p>
                                                </div>
                                                <button onClick={() => handleApply('to_woo')} className="inline-flex items-center gap-1.5 whitespace-nowrap rounded-lg bg-violet-600 hover:bg-violet-700 px-3 py-1.5 text-xs font-semibold text-white transition-colors">
                                                    <ArrowUpFromLine size={14} /> Crear en WooCommerce
                                                </button>
                                            </div>
                                            <ProductList items={diff.onlyInProbability} />
                                        </div>
                                    )}

                                    {diff.onlyInWoo.length > 0 && (
                                        <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
                                            <div className="flex items-start justify-between gap-3">
                                                <div>
                                                    <p className="text-sm font-semibold text-gray-900 dark:text-gray-100">En WooCommerce hay {diff.onlyInWoo.length} producto{diff.onlyInWoo.length !== 1 ? 's' : ''} que no estan en Probability</p>
                                                    <p className="text-[11px] text-gray-400 mt-0.5">Se crearan en Probability aplicando tu configuracion de bodegas.</p>
                                                </div>
                                                <button onClick={() => handleApply('to_probability')} className="inline-flex items-center gap-1.5 whitespace-nowrap rounded-lg bg-blue-600 hover:bg-blue-700 px-3 py-1.5 text-xs font-semibold text-white transition-colors">
                                                    <ArrowDownToLine size={14} /> Crear en Probability
                                                </button>
                                            </div>
                                            <ProductList items={diff.onlyInWoo} />
                                        </div>
                                    )}
                                </>
                            )}

                            {(diff.probabilityNoSku > 0 || diff.woocommerceNoSku > 0) && (
                                <p className="text-[11px] text-amber-600 dark:text-amber-400">
                                    Sin SKU (no se pueden cruzar): {diff.probabilityNoSku} en Probability, {diff.woocommerceNoSku} en WooCommerce.
                                </p>
                            )}
                        </div>
                    )}

                    {(phase === 'running' || phase === 'done') && (
                        <div>
                            <div className="flex items-center gap-2 mb-3 text-sm font-medium text-gray-700 dark:text-gray-200">
                                {direction === 'to_woo' ? <ArrowUpFromLine size={16} className="text-violet-600" /> : <ArrowDownToLine size={16} className="text-blue-600" />}
                                {isFullSync ? 'Sincronizando stock con WooCommerce' : direction === 'to_woo' ? 'Creando en WooCommerce' : 'Creando en Probability'}
                                {phase === 'running' && <Loader2 size={14} className="animate-spin text-gray-400" />}
                            </div>
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Progreso</span>
                                <span className="text-sm font-bold text-gray-900 dark:text-white tabular-nums">{processed} / {total}</span>
                            </div>
                            <div className="h-2 bg-gray-100 dark:bg-gray-800 rounded-full overflow-hidden">
                                <div className="h-full rounded-full bg-gradient-to-r from-violet-500 to-purple-500 transition-all duration-300" style={{ width: `${progressPct}%` }} />
                            </div>
                            <div className="grid grid-cols-3 gap-2 mt-4">
                                <div className="flex items-center justify-between rounded-lg px-3 py-2 text-xs font-semibold bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
                                    <span>Creados</span><span className="tabular-nums">{created}</span>
                                </div>
                                <div className="flex items-center justify-between rounded-lg px-3 py-2 text-xs font-semibold bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300">
                                    <span>Mapeados</span><span className="tabular-nums">{updated}</span>
                                </div>
                                <div className="flex items-center justify-between rounded-lg px-3 py-2 text-xs font-semibold bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300">
                                    <span>Fallidos</span><span className="tabular-nums">{failed}</span>
                                </div>
                            </div>

                            {items.length > 0 && (
                                <div className="mt-4">
                                    <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1.5">Detalle por producto</p>
                                    <div className="max-h-48 overflow-y-auto rounded-lg border border-gray-100 dark:border-gray-800 divide-y divide-gray-100 dark:divide-gray-800">
                                        {[...items].reverse().map((it, i) => (
                                            <div key={i} className="flex items-center justify-between gap-2 px-2.5 py-1.5 text-[11px]">
                                                <div className="min-w-0">
                                                    <p className="text-gray-700 dark:text-gray-200 truncate">{it.name || '(sin nombre)'}</p>
                                                    <p className="text-gray-400 font-mono">{it.sku || '(sin sku)'}</p>
                                                </div>
                                                <div className="flex items-center gap-2 flex-shrink-0">
                                                    <span className="tabular-nums text-gray-600 dark:text-gray-300">{it.quantity} u</span>
                                                    <ActionBadge action={it.action} />
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            )}

                            {phase === 'done' && (
                                <div className="flex justify-between mt-5">
                                    <button onClick={analyze} className="px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg text-sm font-semibold text-gray-700 dark:text-gray-200 transition-colors flex items-center gap-1.5">
                                        <RefreshCw size={14} /> Analizar de nuevo
                                    </button>
                                    <button onClick={onClose} className="px-5 py-2 bg-violet-600 hover:bg-violet-700 text-white rounded-lg text-sm font-semibold transition-colors flex items-center gap-2">
                                        <CheckCircle2 size={16} /> Listo
                                    </button>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

function ActionBadge({ action }: { action: SyncItem['action'] }) {
    const map = {
        created: { label: 'Creado', cls: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' },
        updated: { label: 'Actualizado', cls: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300' },
        failed: { label: 'Fallido', cls: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300' },
    };
    const { label, cls } = map[action];
    return <span className={`px-1.5 py-0.5 rounded font-semibold ${cls}`}>{label}</span>;
}

function ProductList({ items }: { items: Brief[] }) {
    if (items.length === 0) return null;
    return (
        <div className="mt-2 max-h-32 overflow-y-auto rounded-md bg-gray-50 dark:bg-gray-800/60 divide-y divide-gray-100 dark:divide-gray-700">
            {items.slice(0, 100).map((p, i) => (
                <div key={i} className="flex items-center justify-between px-2.5 py-1.5 text-[11px]">
                    <span className="text-gray-700 dark:text-gray-200 truncate">{p.name || '(sin nombre)'}</span>
                    <span className="text-gray-400 font-mono ml-2 flex-shrink-0">{p.sku}</span>
                </div>
            ))}
            {items.length > 100 && <div className="px-2.5 py-1.5 text-[11px] text-gray-400">y {items.length - 100} mas...</div>}
        </div>
    );
}
