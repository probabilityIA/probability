'use client';

import {
    createContext,
    useCallback,
    useContext,
    useMemo,
    useRef,
    useState,
    type ReactNode,
} from 'react';
import { useSSE } from '@/shared/hooks/use-sse';
import type { Integration } from '@/services/integrations/core/domain/types';
import { getSyncProvider, GLOBAL_INVENTORY_EVENT_TYPES } from './providers';

export type SyncNodeState = 'idle' | 'queued' | 'active' | 'scan' | 'done' | 'error';
export type SyncMode = 'idle' | 'inventory' | 'products';
export type SyncEnvironment = 'inventory' | 'products';

export interface InventoryResult {
    kind: 'inventory';
    total: number;
    updated: number;
    unchanged: number;
    skipped: number;
    failed: number;
}

export interface ProductsResult {
    kind: 'products';
    matched: number;
    notAssociated: number;
    onlyInProbability: number;
    onlyInChannel: number;
}

export type SyncResult = InventoryResult | ProductsResult | { kind: 'error'; message: string };

export interface SyncDetailItem {
    sku: string;
    label: string;
    tone: 'ok' | 'warn' | 'error';
}

interface SyncActivityValue {
    mode: SyncMode;
    running: boolean;
    nodes: Record<number, SyncNodeState>;
    progress: Record<number, { processed: number; total: number }>;
    results: Record<number, SyncResult>;
    details: Record<number, SyncDetailItem[]>;
    environment: SyncEnvironment;
    setEnvironment: (env: SyncEnvironment) => void;
    runCurrent: () => void;
    runInventory: () => void;
    runProducts: () => void;
    reset: () => void;
}

const SyncActivityContext = createContext<SyncActivityValue | null>(null);

const SYNC_TIMEOUT_MS = 6 * 60 * 1000;

interface SyncStartResult {
    success?: boolean;
    correlation_id?: string;
    message?: string;
}

function countOf(value: unknown): number {
    return Array.isArray(value) ? value.length : Number(value) || 0;
}

interface ProviderProps {
    children: ReactNode;
    integrations: Integration[];
    businessId: number | null;
}

export function SyncActivityProvider({ children, integrations, businessId }: ProviderProps) {
    const [mode, setMode] = useState<SyncMode>('idle');
    const [running, setRunning] = useState(false);
    const [nodes, setNodes] = useState<Record<number, SyncNodeState>>({});
    const [progress, setProgress] = useState<Record<number, { processed: number; total: number }>>({});
    const [results, setResults] = useState<Record<number, SyncResult>>({});
    const [details, setDetails] = useState<Record<number, SyncDetailItem[]>>({});
    const [environment, setEnvironment] = useState<SyncEnvironment>('inventory');

    const pushDetail = useCallback((id: number, item: SyncDetailItem) => {
        setDetails(prev => {
            const list = prev[id] || [];
            if (list.length > 400) return prev;
            return { ...prev, [id]: [...list, item] };
        });
    }, []);

    const corrToIntegration = useRef<Map<string, number>>(new Map());
    const completion = useRef<Map<number, () => void>>(new Map());

    const eligible = useMemo(
        () => integrations.filter(i => i.is_active && getSyncProvider(i.integration_type_id)),
        [integrations],
    );

    const patchNode = useCallback((id: number, state: SyncNodeState) => {
        setNodes(prev => (prev[id] === state ? prev : { ...prev, [id]: state }));
    }, []);

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data);
            const eventType: string = parsed.type || parsed.metadata?.event_type || '';
            const data = parsed.data;
            if (!data?.correlation_id) return;
            const id = corrToIntegration.current.get(data.correlation_id);
            if (id === undefined) return;

            if (eventType.endsWith('.inventory.sync.started')) {
                patchNode(id, 'active');
                setProgress(prev => ({ ...prev, [id]: { processed: 0, total: Number(data.total) || 0 } }));
            } else if (eventType.endsWith('.inventory.sync.item')) {
                const action = String(data.action || '');
                const failed = /fail|error/i.test(action);
                pushDetail(id, {
                    sku: String(data.sku || '(sin sku)'),
                    label: failed
                        ? String(data.error || data.message || 'fallo al actualizar')
                        : `${action || 'actualizado'} · ${data.quantity ?? '-'} u.`,
                    tone: failed ? 'error' : /skip|omit|unchanged/i.test(action) ? 'warn' : 'ok',
                });
            } else if (eventType.endsWith('.inventory.sync.progress')) {
                setProgress(prev => ({
                    ...prev,
                    [id]: { processed: Number(data.processed) || 0, total: Number(data.total) || prev[id]?.total || 0 },
                }));
            } else if (eventType.endsWith('.inventory.sync.completed')) {
                const total = Number(data.total) || 0;
                setProgress(prev => ({ ...prev, [id]: { processed: total, total } }));
                setResults(prev => ({
                    ...prev,
                    [id]: {
                        kind: 'inventory',
                        total,
                        updated: Number(data.updated) || 0,
                        unchanged: Number(data.unchanged) || 0,
                        skipped: Number(data.skipped) || 0,
                        failed: Number(data.failed) || 0,
                    },
                }));
                const failedSkus = Array.isArray(data.failed_skus) ? data.failed_skus : [];
                for (const raw of failedSkus) {
                    const sku = typeof raw === 'string' ? raw : String((raw as Record<string, unknown>)?.sku ?? '');
                    const msg = typeof raw === 'string' ? 'fallo al actualizar' : String((raw as Record<string, unknown>)?.error ?? 'fallo al actualizar');
                    pushDetail(id, { sku: sku || '(sin sku)', label: msg, tone: 'error' });
                }
                patchNode(id, 'done');
                completion.current.get(id)?.();
            }
        } catch {
            return;
        }
    }, [patchNode, pushDetail]);

    useSSE({
        businessId: businessId ?? 0,
        eventTypes: GLOBAL_INVENTORY_EVENT_TYPES,
        onMessage: handleMessage,
        enabled: true,
    });

    const reset = useCallback(() => {
        setMode('idle');
        setRunning(false);
        setNodes({});
        setProgress({});
        setResults({});
        setDetails({});
    }, []);

    const runInventory = useCallback(async () => {
        if (running || eligible.length === 0) return;
        setRunning(true);
        setMode('inventory');
        setResults({});
        setProgress({});
        setDetails({});
        corrToIntegration.current.clear();

        const queued: Record<number, SyncNodeState> = {};
        for (const integration of eligible) queued[integration.id] = 'queued';
        setNodes(queued);

        for (const integration of eligible) {
            const provider = getSyncProvider(integration.integration_type_id);
            if (!provider) continue;

            patchNode(integration.id, 'active');
            let result: SyncStartResult | null = null;
            try {
                result = await provider.syncInventory(integration.id, businessId ?? undefined) as SyncStartResult;
            } catch {
                result = null;
            }

            if (!result?.success || !result?.correlation_id) {
                patchNode(integration.id, 'error');
                setResults(prev => ({
                    ...prev,
                    [integration.id]: { kind: 'error', message: result?.message || 'No se pudo iniciar' },
                }));
                continue;
            }

            corrToIntegration.current.set(result.correlation_id, integration.id);

            await new Promise<void>(resolve => {
                const timer = window.setTimeout(() => {
                    patchNode(integration.id, 'error');
                    setResults(prev => ({
                        ...prev,
                        [integration.id]: { kind: 'error', message: 'Continua en segundo plano' },
                    }));
                    completion.current.delete(integration.id);
                    resolve();
                }, SYNC_TIMEOUT_MS);
                completion.current.set(integration.id, () => {
                    window.clearTimeout(timer);
                    completion.current.delete(integration.id);
                    resolve();
                });
            });
        }

        setRunning(false);
        setMode('idle');
    }, [running, eligible, businessId, patchNode]);

    const runProducts = useCallback(async () => {
        if (running || eligible.length === 0) return;
        setRunning(true);
        setMode('products');
        setResults({});
        setProgress({});
        setDetails({});

        const scanning: Record<number, SyncNodeState> = {};
        for (const integration of eligible) scanning[integration.id] = 'scan';
        setNodes(scanning);

        await Promise.all(eligible.map(async integration => {
            const provider = getSyncProvider(integration.integration_type_id);
            if (!provider) return;
            try {
                const res = await provider.reconcileProducts(
                    integration.id,
                    businessId ?? undefined,
                ) as Record<string, unknown>;

                if (res?.success === false) {
                    patchNode(integration.id, 'error');
                    setResults(prev => ({
                        ...prev,
                        [integration.id]: { kind: 'error', message: String(res?.message || 'No se pudo comparar') },
                    }));
                    return;
                }

                const detail: SyncDetailItem[] = [];
                const push = (arr: unknown, tone: SyncDetailItem['tone'], label: string) => {
                    if (!Array.isArray(arr)) return;
                    for (const raw of arr.slice(0, 120)) {
                        const obj = raw as Record<string, unknown>;
                        detail.push({ sku: String(obj?.sku ?? ''), label: `${label}${obj?.name ? ` · ${obj.name}` : ''}`, tone });
                    }
                };
                push(res?.matched_not_associated, 'warn', 'sin asociar');
                push(res?.only_in_probability, 'warn', 'solo en Probability');
                push(res?.[provider.onlyInChannelField], 'warn', 'solo en el canal');
                setDetails(prev => ({ ...prev, [integration.id]: detail }));

                setResults(prev => ({
                    ...prev,
                    [integration.id]: {
                        kind: 'products',
                        matched: Number(res?.matched) || 0,
                        notAssociated: countOf(res?.matched_not_associated),
                        onlyInProbability: countOf(res?.only_in_probability),
                        onlyInChannel: countOf(res?.[provider.onlyInChannelField]),
                    },
                }));
                patchNode(integration.id, 'done');
            } catch {
                patchNode(integration.id, 'error');
                setResults(prev => ({
                    ...prev,
                    [integration.id]: { kind: 'error', message: 'No se pudo comparar' },
                }));
            }
        }));

        setRunning(false);
        setMode('idle');
    }, [running, eligible, businessId, patchNode]);

    const runCurrent = useCallback(() => {
        if (environment === 'products') runProducts();
        else runInventory();
    }, [environment, runInventory, runProducts]);

    const value = useMemo<SyncActivityValue>(() => ({
        mode,
        running,
        nodes,
        progress,
        results,
        details,
        environment,
        setEnvironment,
        runCurrent,
        runInventory,
        runProducts,
        reset,
    }), [mode, running, nodes, progress, results, details, environment, runCurrent, runInventory, runProducts, reset]);

    return <SyncActivityContext.Provider value={value}>{children}</SyncActivityContext.Provider>;
}

const EMPTY: SyncActivityValue = {
    mode: 'idle',
    running: false,
    nodes: {},
    progress: {},
    results: {},
    details: {},
    environment: 'inventory',
    setEnvironment: () => undefined,
    runCurrent: () => undefined,
    runInventory: () => undefined,
    runProducts: () => undefined,
    reset: () => undefined,
};

export function useSyncActivity(): SyncActivityValue {
    return useContext(SyncActivityContext) ?? EMPTY;
}
