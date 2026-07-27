'use client';

import {
    createContext,
    useCallback,
    useContext,
    useEffect,
    useMemo,
    useRef,
    useState,
    type ReactNode,
} from 'react';
import { useSSE } from '@/shared/hooks/use-sse';
import type { Integration } from '@/services/integrations/core/domain/types';
import type { SyncRunKind, SyncRunRecord } from '../domain/types';
import { listSyncRunsAction } from '../infra/actions/sync-runs';
import { getSyncProvider, GLOBAL_INVENTORY_EVENT_TYPES } from './providers';

export type SyncNodeState = 'idle' | 'queued' | 'active' | 'scan' | 'done' | 'error';
export type SyncMode = 'idle' | 'inventory' | 'products';
export type SyncEnvironment = 'inventory' | 'products' | 'invoicing';

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

export type DetailGroup =
    | 'both'
    | 'not_associated'
    | 'only_probability'
    | 'only_channel'
    | 'updated'
    | 'skipped'
    | 'failed';

export interface SyncDetailItem {
    sku: string;
    label: string;
    tone: 'ok' | 'warn' | 'error';
    group: DetailGroup;
}

export type ProductActionKey = 'associate' | 'createInChannel' | 'createInProbability' | 'updateInProbability' | 'createBothSides';

export interface ProductActionResult {
    ok: boolean;
    message: string;
}

interface SyncActivityValue {
    mode: SyncMode;
    businessId: number | null;
    running: boolean;
    nodes: Record<number, SyncNodeState>;
    progress: Record<number, { processed: number; total: number }>;
    results: Record<number, SyncResult>;
    details: Record<number, SyncDetailItem[]>;
    environment: SyncEnvironment | null;
    setEnvironment: (env: SyncEnvironment | null) => void;
    canRun: boolean;
    lastRuns: Record<number, Partial<Record<SyncRunKind, SyncRunRecord>>>;
    actionBusy: Record<number, ProductActionKey | null>;
    actionResult: Record<number, ProductActionResult | null>;
    runProductAction: (integrationId: number, action: ProductActionKey, skus?: string[]) => void;
    runInventoryOne: (integrationId: number) => void;
    runCurrent: () => void;
    runInventory: () => void;
    runProducts: () => void;
    reset: () => void;
}

const SyncActivityContext = createContext<SyncActivityValue | null>(null);

const SYNC_TIMEOUT_MS = 6 * 60 * 1000;

const APPLY_SETTLE_MS = 4000;

interface SyncStartResult {
    success?: boolean;
    correlation_id?: string;
    message?: string;
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
    const [environment, setEnvironment] = useState<SyncEnvironment | null>(null);
    const [lastRuns, setLastRuns] = useState<Record<number, Partial<Record<SyncRunKind, SyncRunRecord>>>>({});
    const [actionBusy, setActionBusy] = useState<Record<number, ProductActionKey | null>>({});
    const [actionResult, setActionResult] = useState<Record<number, ProductActionResult | null>>({});

    const loadLastRuns = useCallback(async () => {
        const rows = await listSyncRunsAction(businessId ?? undefined);
        const map: Record<number, Partial<Record<SyncRunKind, SyncRunRecord>>> = {};
        for (const row of rows) {
            map[row.integration_id] = { ...(map[row.integration_id] || {}), [row.kind]: row };
        }
        setLastRuns(map);
    }, [businessId]);

    useEffect(() => {
        loadLastRuns();
    }, [loadLastRuns]);

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
                const skipped = /skip|omit|unchanged/i.test(action);
                pushDetail(id, {
                    sku: String(data.sku || '(sin sku)'),
                    label: failed
                        ? String(data.error || data.message || 'fallo al actualizar')
                        : `${action || 'actualizado'} · ${data.quantity ?? '-'} u.`,
                    tone: failed ? 'error' : skipped ? 'warn' : 'ok',
                    group: failed ? 'failed' : skipped ? 'skipped' : 'updated',
                });
            } else if (eventType.endsWith('.inventory.sync.progress')) {
                setProgress(prev => ({
                    ...prev,
                    [id]: { processed: Number(data.processed) || 0, total: Number(data.total) || prev[id]?.total || 0 },
                }));
            } else if (eventType.endsWith('.product.reconcile.started')) {
                patchNode(id, 'scan');
            } else if (eventType.endsWith('.product.reconcile.completed')) {
                if (data.error) {
                    patchNode(id, 'error');
                    setResults(prev => ({
                        ...prev,
                        [id]: { kind: 'error', message: String(data.error) },
                    }));
                } else {
                    setResults(prev => ({
                        ...prev,
                        [id]: {
                            kind: 'products',
                            matched: Number(data.matched) || 0,
                            notAssociated: Number(data.not_associated) || 0,
                            onlyInProbability: Number(data.only_in_probability) || 0,
                            onlyInChannel: Number(data.only_in_channel) || 0,
                        },
                    }));
                    patchNode(id, 'done');
                }
                completion.current.get(id)?.();
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
                    pushDetail(id, { sku: sku || '(sin sku)', label: msg, tone: 'error', group: 'failed' });
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
        setActionResult({});
    }, []);

    const syncInventoryOne = useCallback(async (integration: Integration) => {
        const provider = getSyncProvider(integration.integration_type_id);
        if (!provider) return;

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
            return;
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
    }, [businessId, patchNode]);

    const runInventoryOne = useCallback(async (integrationId: number) => {
        const integration = eligible.find(i => i.id === integrationId);
        if (!integration || running) return;
        setRunning(true);
        setMode('inventory');
        setResults(prev => {
            const next = { ...prev };
            delete next[integrationId];
            return next;
        });
        setDetails(prev => ({ ...prev, [integrationId]: [] }));
        setProgress(prev => ({ ...prev, [integrationId]: { processed: 0, total: 0 } }));

        await syncInventoryOne(integration);

        setRunning(false);
        setMode('idle');
        loadLastRuns();
    }, [eligible, running, syncInventoryOne, loadLastRuns]);

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
            await syncInventoryOne(integration);
        }

        setRunning(false);
        setMode('idle');
        loadLastRuns();
    }, [running, eligible, syncInventoryOne, loadLastRuns]);

    const reconcileOne = useCallback(async (integration: Integration) => {
        const provider = getSyncProvider(integration.integration_type_id);
        if (!provider) return;

        patchNode(integration.id, 'scan');
        setDetails(prev => ({ ...prev, [integration.id]: [] }));

        let start: SyncStartResult | null = null;
        try {
            start = await provider.reconcileProducts(
                integration.id,
                businessId ?? undefined,
            ) as SyncStartResult;
        } catch {
            start = null;
        }

        if (!start?.success || !start?.correlation_id) {
            patchNode(integration.id, 'error');
            setResults(prev => ({
                ...prev,
                [integration.id]: { kind: 'error', message: start?.message || 'No se pudo comparar' },
            }));
            return;
        }

        corrToIntegration.current.set(start.correlation_id, integration.id);

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
    }, [businessId, patchNode]);

    const runProducts = useCallback(async () => {
        if (running || eligible.length === 0) return;
        setRunning(true);
        setMode('products');
        setResults({});
        setProgress({});
        setDetails({});
        corrToIntegration.current.clear();

        const scanning: Record<number, SyncNodeState> = {};
        for (const integration of eligible) scanning[integration.id] = 'scan';
        setNodes(scanning);

        await Promise.all(eligible.map(integration => reconcileOne(integration)));

        setRunning(false);
        setMode('idle');
        loadLastRuns();
    }, [running, eligible, reconcileOne, loadLastRuns]);

    const runProductAction = useCallback(async (integrationId: number, action: ProductActionKey, skus?: string[]) => {
        const integration = eligible.find(i => i.id === integrationId);
        const provider = integration ? getSyncProvider(integration.integration_type_id) : null;
        if (!integration || !provider || actionBusy[integrationId]) return;

        const steps: ProductActionKey[] = action === 'createBothSides'
            ? ['createInChannel', 'createInProbability']
            : [action];

        const runs = steps
            .map(step => (step === 'associate'
                ? (id: number, bid?: number, only?: string[]) => provider.associateProducts(id, bid, only)
                : step === 'createBothSides'
                    ? undefined
                    : provider.apply[step]))
            .filter(Boolean) as ((id: number, bid?: number, only?: string[]) => Promise<unknown>)[];
        if (runs.length === 0) return;

        setActionBusy(prev => ({ ...prev, [integrationId]: action }));
        setActionResult(prev => ({ ...prev, [integrationId]: null }));
        try {
            let ok = true;
            let message = '';
            for (const run of runs) {
                const res = await run(integrationId, businessId ?? undefined, skus) as Record<string, unknown>;
                if (res?.success === false) {
                    ok = false;
                    message = String(res?.message || 'No se pudo aplicar');
                    break;
                }
                message = String(res?.message || 'Aplicado');
            }
            setActionResult(prev => ({ ...prev, [integrationId]: { ok, message } }));
            if (ok) {
                await new Promise(resolve => setTimeout(resolve, APPLY_SETTLE_MS));
                await reconcileOne(integration);
            }
        } catch {
            setActionResult(prev => ({
                ...prev,
                [integrationId]: { ok: false, message: 'No se pudo aplicar' },
            }));
        } finally {
            setActionBusy(prev => ({ ...prev, [integrationId]: null }));
            loadLastRuns();
        }
    }, [eligible, businessId, actionBusy, reconcileOne, loadLastRuns]);

    const canRun = environment === 'inventory' || environment === 'products';

    const runCurrent = useCallback(() => {
        if (environment === 'products') runProducts();
        else if (environment === 'inventory') runInventory();
    }, [environment, runInventory, runProducts]);

    const value = useMemo<SyncActivityValue>(() => ({
        mode,
        businessId,
        running,
        nodes,
        progress,
        results,
        details,
        environment,
        setEnvironment,
        canRun,
        lastRuns,
        actionBusy,
        actionResult,
        runProductAction,
        runInventoryOne,
        runCurrent,
        runInventory,
        runProducts,
        reset,
    }), [mode, businessId, running, nodes, progress, results, details, environment, canRun, lastRuns, actionBusy, actionResult, runProductAction, runInventoryOne, runCurrent, runInventory, runProducts, reset]);

    return <SyncActivityContext.Provider value={value}>{children}</SyncActivityContext.Provider>;
}

const EMPTY: SyncActivityValue = {
    mode: 'idle',
    businessId: null,
    running: false,
    nodes: {},
    progress: {},
    results: {},
    details: {},
    environment: null,
    setEnvironment: () => undefined,
    canRun: false,
    lastRuns: {},
    actionBusy: {},
    actionResult: {},
    runProductAction: () => undefined,
    runInventoryOne: () => undefined,
    runCurrent: () => undefined,
    runInventory: () => undefined,
    runProducts: () => undefined,
    reset: () => undefined,
};

export function useSyncActivity(): SyncActivityValue {
    return useContext(SyncActivityContext) ?? EMPTY;
}
