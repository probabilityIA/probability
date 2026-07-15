'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { X, CheckCircle2, Loader2, AlertCircle, RefreshCw, ArrowDownToLine, ArrowRightLeft, Link2 } from 'lucide-react';
import { useSSE } from '@/shared/hooks/use-sse';
import { reconcileSiigoProductsAction, applySiigoProductsAction } from '../../infra/actions';

interface SiigoProductSyncModalProps {
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
    matchedNotAssociated: Brief[];
    onlyInProbability: Brief[];
    onlyInSiigo: Brief[];
    probabilityNoSku: number;
    siigoNoSku: number;
}

const PRODUCT_EVENT_TYPES = [
    'siigo.product.sync.started',
    'siigo.product.sync.progress',
    'siigo.product.sync.completed',
];

type Phase = 'analyzing' | 'diff' | 'running' | 'done' | 'error';

export function SiigoProductSyncModal({ isOpen, onClose, integrationId, businessId, onCompleted }: SiigoProductSyncModalProps) {
    const [phase, setPhase] = useState<Phase>('analyzing');
    const [diff, setDiff] = useState<Diff | null>(null);
    const [selected, setSelected] = useState<Set<string>>(new Set());
    const [actionLabel, setActionLabel] = useState('Procesando');
    const [total, setTotal] = useState(0);
    const [processed, setProcessed] = useState(0);
    const [created, setCreated] = useState(0);
    const [failed, setFailed] = useState(0);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);

    const correlationRef = useRef<string | null>(null);

    const analyze = useCallback(async () => {
        setPhase('analyzing');
        setErrorMessage(null);
        setSelected(new Set());
        const res: any = await reconcileSiigoProductsAction(integrationId, businessId ?? undefined);
        if (!res?.success) {
            setErrorMessage(res?.message || 'No se pudo analizar los productos');
            setPhase('error');
            return;
        }
        setDiff({
            matched: Number(res.matched) || 0,
            matchedNotAssociated: res.matched_not_associated || [],
            onlyInProbability: res.only_in_probability || [],
            onlyInSiigo: res.only_in_siigo || [],
            probabilityNoSku: Number(res.probability_no_sku) || 0,
            siigoNoSku: Number(res.siigo_no_sku) || 0,
        });
        setPhase('diff');
    }, [integrationId, businessId]);

    useEffect(() => {
        if (!isOpen) {
            setPhase('analyzing');
            setDiff(null);
            setSelected(new Set());
            setTotal(0);
            setProcessed(0);
            setCreated(0);
            setFailed(0);
            setErrorMessage(null);
            correlationRef.current = null;
            return;
        }
        analyze();
    }, [isOpen, analyze]);

    const runApply = async (skus: string[] | undefined, label: string, count: number) => {
        setActionLabel(label);
        setPhase('running');
        setTotal(count);
        setProcessed(0);
        setCreated(0);
        setFailed(0);
        correlationRef.current = null;
        const res: any = await applySiigoProductsAction(integrationId, businessId ?? undefined, skus);
        if (!res?.success || !res?.correlation_id) {
            setErrorMessage(res?.message || 'No se pudo iniciar la operacion');
            setPhase('error');
            return;
        }
        correlationRef.current = res.correlation_id;
    };

    const handleCreate = () => runApply(undefined, 'Creando en Probability', diff?.onlyInSiigo.length || 0);
    const handleAssociateAll = () => {
        const skus = (diff?.matchedNotAssociated || []).map((p) => p.sku);
        runApply(skus, 'Asociando al canal', skus.length);
    };
    const handleAssociateSelected = () => {
        const skus = Array.from(selected);
        runApply(skus, 'Asociando al canal', skus.length);
    };

    const toggleSelected = (sku: string) => {
        setSelected((prev) => {
            const next = new Set(prev);
            if (next.has(sku)) next.delete(sku);
            else next.add(sku);
            return next;
        });
    };

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data);
            const eventType = parsed.type || parsed.metadata?.event_type;
            const data = parsed.data;
            if (!data) return;
            const corr = correlationRef.current;
            if (!corr || data.correlation_id !== corr) return;

            if (eventType === 'siigo.product.sync.started') {
                setTotal(Number(data.total) || 0);
            } else if (eventType === 'siigo.product.sync.progress') {
                setProcessed(Number(data.processed) || 0);
                setCreated(Number(data.created) || 0);
                setFailed(Number(data.failed) || 0);
            } else if (eventType === 'siigo.product.sync.completed') {
                setProcessed(Number(data.total) || 0);
                setTotal(Number(data.total) || 0);
                setCreated(Number(data.created) || 0);
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
    const nothingToDo = diff && diff.onlyInSiigo.length === 0 && diff.matchedNotAssociated.length === 0;

    return (
        <div className="fixed inset-0 z-[1100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl w-full max-w-lg flex flex-col overflow-hidden border border-gray-200 dark:border-gray-700 max-h-[90vh]">
                <div className="flex items-center justify-between px-6 py-5 border-b border-gray-100 dark:border-gray-800 bg-gradient-to-r from-emerald-50 to-green-50 dark:from-emerald-950/40 dark:to-green-950/40">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-emerald-500/10 dark:bg-emerald-400/10 flex items-center justify-center">
                            <ArrowRightLeft size={18} className="text-emerald-600 dark:text-emerald-400" />
                        </div>
                        <div>
                            <h2 className="text-lg font-bold text-gray-900 dark:text-white">Sincronizacion de Productos</h2>
                            <p className="text-xs text-gray-500 dark:text-gray-400">Siigo &rarr; Probability</p>
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
                            <button onClick={analyze} className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-sm font-semibold transition-colors">Reintentar</button>
                        </div>
                    )}

                    {phase === 'diff' && diff && (
                        <div className="space-y-4">
                            <div className="flex items-center gap-2 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 px-3 py-2">
                                <CheckCircle2 size={16} className="text-emerald-600 dark:text-emerald-400" />
                                <span className="text-sm text-emerald-800 dark:text-emerald-300"><strong>{diff.matched}</strong> productos ya asociados a este canal</span>
                            </div>

                            {nothingToDo && (
                                <div className="text-center py-6 text-gray-600 dark:text-gray-300">
                                    <CheckCircle2 size={40} className="mx-auto text-emerald-500 mb-2" />
                                    <p className="font-semibold">Todo al dia</p>
                                    <p className="text-xs text-gray-400 mt-1">No hay productos por crear ni por asociar a este canal.</p>
                                </div>
                            )}

                            {diff.matchedNotAssociated.length > 0 && (
                                <div className="rounded-lg border border-amber-200 dark:border-amber-800 bg-amber-50/40 dark:bg-amber-900/10 p-3">
                                    <div className="flex items-start justify-between gap-3">
                                        <div>
                                            <p className="text-sm font-semibold text-gray-900 dark:text-gray-100">{diff.matchedNotAssociated.length} producto{diff.matchedNotAssociated.length !== 1 ? 's' : ''} existe{diff.matchedNotAssociated.length !== 1 ? 'n' : ''} en Probability pero no asociado{diff.matchedNotAssociated.length !== 1 ? 's' : ''} a este canal</p>
                                            <p className="text-[11px] text-gray-400 mt-0.5">Asocialos para que el canal los reconozca como propios.</p>
                                        </div>
                                        <button onClick={handleAssociateAll} className="inline-flex items-center gap-1.5 whitespace-nowrap rounded-lg bg-amber-600 hover:bg-amber-700 px-3 py-1.5 text-xs font-semibold text-white transition-colors">
                                            <Link2 size={14} /> Asociar todos
                                        </button>
                                    </div>
                                    <SelectableProductList items={diff.matchedNotAssociated} selected={selected} onToggle={toggleSelected} />
                                    <div className="mt-2 flex items-center justify-between">
                                        <span className="text-[11px] text-gray-400">{selected.size} seleccionado{selected.size !== 1 ? 's' : ''}</span>
                                        <button onClick={handleAssociateSelected} disabled={selected.size === 0} className="inline-flex items-center gap-1.5 rounded-lg border border-amber-300 dark:border-amber-700 px-3 py-1.5 text-xs font-semibold text-amber-700 dark:text-amber-300 disabled:opacity-40 disabled:cursor-not-allowed transition-colors">
                                            <Link2 size={14} /> Asociar seleccionados
                                        </button>
                                    </div>
                                </div>
                            )}

                            {diff.onlyInSiigo.length > 0 && (
                                <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
                                    <div className="flex items-start justify-between gap-3">
                                        <div>
                                            <p className="text-sm font-semibold text-gray-900 dark:text-gray-100">En Siigo hay {diff.onlyInSiigo.length} producto{diff.onlyInSiigo.length !== 1 ? 's' : ''} que no estan en Probability</p>
                                            <p className="text-[11px] text-gray-400 mt-0.5">Se crearan en Probability y quedaran asociados a este canal.</p>
                                        </div>
                                        <button onClick={handleCreate} className="inline-flex items-center gap-1.5 whitespace-nowrap rounded-lg bg-emerald-600 hover:bg-emerald-700 px-3 py-1.5 text-xs font-semibold text-white transition-colors">
                                            <ArrowDownToLine size={14} /> Crear en Probability
                                        </button>
                                    </div>
                                    <ProductList items={diff.onlyInSiigo} />
                                </div>
                            )}

                            {diff.onlyInProbability.length > 0 && (
                                <div className="rounded-lg border border-gray-100 dark:border-gray-800 bg-gray-50/60 dark:bg-gray-800/40 p-3">
                                    <p className="text-sm font-semibold text-gray-700 dark:text-gray-300">En Probability hay {diff.onlyInProbability.length} producto{diff.onlyInProbability.length !== 1 ? 's' : ''} que no estan en Siigo</p>
                                    <p className="text-[11px] text-gray-400 mt-0.5">Informativo. La integracion Siigo es de solo lectura: estos no se crean en Siigo.</p>
                                    <ProductList items={diff.onlyInProbability} />
                                </div>
                            )}

                            {(diff.probabilityNoSku > 0 || diff.siigoNoSku > 0) && (
                                <p className="text-[11px] text-amber-600 dark:text-amber-400">
                                    Sin SKU (no se pueden cruzar): {diff.probabilityNoSku} en Probability, {diff.siigoNoSku} en Siigo.
                                </p>
                            )}
                        </div>
                    )}

                    {(phase === 'running' || phase === 'done') && (
                        <div>
                            <div className="flex items-center gap-2 mb-3 text-sm font-medium text-gray-700 dark:text-gray-200">
                                <ArrowDownToLine size={16} className="text-emerald-600" />
                                {actionLabel}
                                {phase === 'running' && <Loader2 size={14} className="animate-spin text-gray-400" />}
                            </div>
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Progreso</span>
                                <span className="text-sm font-bold text-gray-900 dark:text-white tabular-nums">{processed} / {total}</span>
                            </div>
                            <div className="h-2 bg-gray-100 dark:bg-gray-800 rounded-full overflow-hidden">
                                <div className="h-full rounded-full bg-gradient-to-r from-emerald-500 to-green-500 transition-all duration-300" style={{ width: `${progressPct}%` }} />
                            </div>
                            <div className="grid grid-cols-2 gap-2 mt-4">
                                <div className="flex items-center justify-between rounded-lg px-3 py-2 text-xs font-semibold bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
                                    <span>Enviados</span><span className="tabular-nums">{created}</span>
                                </div>
                                <div className="flex items-center justify-between rounded-lg px-3 py-2 text-xs font-semibold bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300">
                                    <span>Fallidos</span><span className="tabular-nums">{failed}</span>
                                </div>
                            </div>

                            {phase === 'done' && (
                                <>
                                    <p className="text-[11px] text-gray-400 mt-4">La operacion corre en segundo plano; puede tardar unos segundos en reflejarse.</p>
                                    <div className="flex justify-between mt-4">
                                        <button onClick={analyze} className="px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg text-sm font-semibold text-gray-700 dark:text-gray-200 transition-colors flex items-center gap-1.5">
                                            <RefreshCw size={14} /> Analizar de nuevo
                                        </button>
                                        <button onClick={onClose} className="px-5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-sm font-semibold transition-colors flex items-center gap-2">
                                            <CheckCircle2 size={16} /> Listo
                                        </button>
                                    </div>
                                </>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
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

function SelectableProductList({ items, selected, onToggle }: { items: Brief[]; selected: Set<string>; onToggle: (sku: string) => void }) {
    if (items.length === 0) return null;
    return (
        <div className="mt-2 max-h-40 overflow-y-auto rounded-md bg-white dark:bg-gray-800/60 border border-amber-100 dark:border-amber-900/40 divide-y divide-gray-100 dark:divide-gray-700">
            {items.slice(0, 200).map((p, i) => (
                <label key={i} className="flex items-center gap-2 px-2.5 py-1.5 text-[11px] cursor-pointer hover:bg-amber-50/50 dark:hover:bg-amber-900/10">
                    <input
                        type="checkbox"
                        checked={selected.has(p.sku)}
                        onChange={() => onToggle(p.sku)}
                        className="h-3.5 w-3.5 rounded border-gray-300 text-amber-600 focus:ring-amber-500"
                    />
                    <span className="text-gray-700 dark:text-gray-200 truncate flex-1">{p.name || '(sin nombre)'}</span>
                    <span className="text-gray-400 font-mono ml-2 flex-shrink-0">{p.sku}</span>
                </label>
            ))}
            {items.length > 200 && <div className="px-2.5 py-1.5 text-[11px] text-gray-400">y {items.length - 200} mas (usa "Asociar todos")...</div>}
        </div>
    );
}
