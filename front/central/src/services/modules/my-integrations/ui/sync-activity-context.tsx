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
    channelNoSku: number;
    skuChanged: number;
    skuTypo: number;
    matchRules?: ProductMatchRuleSummary[];
}

export type SyncResult = InventoryResult | ProductsResult | { kind: 'error'; message: string };

export type DetailGroup =
    | 'both'
    | 'not_associated'
    | 'only_probability'
    | 'only_channel'
    | 'channel_no_sku'
    | 'sku_changed'
    | 'sku_typo'
    | 'sku_spacing'
    | 'updated'
    | 'skipped'
    | 'failed';

export interface SyncDetailItem {
    sku: string;
    label: string;
    tone: 'ok' | 'warn' | 'error';
    group: DetailGroup;
    matchedBy?: string;
    matchedValue?: string;
    parentRef?: string;
    parentLabel?: string;
    variantLabel?: string;
}

export interface ProductMatchRuleSummary {
    probability: string;
    channel: string;
}

export type ProductActionKey = 'associate' | 'createInChannel' | 'createInProbability' | 'updateInProbability' | 'createBothSides';

export interface ProductActionResult {
    ok: boolean;
    message: string;
    pending?: boolean;
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
    runInventoryOne: (integrationId: number, skus?: string[]) => void;
    runCurrent: () => void;
    runInventory: () => void;
    runProducts: () => void;
    reset: () => void;
}

const SyncActivityContext = createContext<SyncActivityValue | null>(null);

const SYNC_TIMEOUT_MS = 90 * 1000;

const RUN_CLOCK_SKEW_MS = 60 * 1000;

const APPLY_SETTLE_MS = 4000;

const APPLY_TIMEOUT_MS = 10 * 60 * 1000;

const MAX_PENDING_RUNS = 20;

const MAX_PENDING_EVENTS_PER_RUN = 100;

interface PendingEvent {
    eventType: string;
    data: any;
}

interface FailedItem {
    sku?: string;
    error?: string;
}

export function buildApplyMessage(
    total: number,
    created: number,
    updated: number,
    failed: number,
    failedItems?: unknown,
): string {
    const parts: string[] = [];
    if (created > 0) parts.push(`${created} creados`);
    if (updated > 0) parts.push(`${updated} actualizados`);
    if (failed > 0) parts.push(`${failed} con error`);
    if (parts.length === 0) return total === 0 ? 'No habia productos por aplicar' : 'Sin cambios';

    let message = `${parts.join(', ')} de ${total}`;
    const items = Array.isArray(failedItems) ? (failedItems as FailedItem[]) : [];
    const firstError = items.find(item => item?.error)?.error;
    if (firstError) {
        const sample = items[0]?.sku ? `${items[0].sku}: ${firstError}` : firstError;
        message += ` — ${sample}`;
        if (items.length > 1) message += ` (y ${items.length - 1} mas)`;
    }
    return message;
}

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
    const applyCompletion = useRef<Map<number, () => void>>(new Map());

    const eligible = useMemo(
        () => integrations.filter(i => i.is_active && getSyncProvider(i.integration_type_id)),
        [integrations],
    );

    const patchNode = useCallback((id: number, state: SyncNodeState) => {
        setNodes(prev => (prev[id] === state ? prev : { ...prev, [id]: state }));
    }, []);

    const pendingEvents = useRef<Map<string, PendingEvent[]>>(new Map());

    const applyEvent = useCallback((id: number, eventType: string, data: any) => {
        try {
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
            } else if (eventType.endsWith('.product.sync.started')) {
                patchNode(id, 'active');
                setProgress(prev => ({ ...prev, [id]: { processed: 0, total: Number(data.total) || 0 } }));
            } else if (eventType.endsWith('.product.sync.progress')) {
                setProgress(prev => ({
                    ...prev,
                    [id]: { processed: Number(data.processed) || 0, total: Number(data.total) || prev[id]?.total || 0 },
                }));
            } else if (eventType.endsWith('.product.sync.completed')) {
                const created = Number(data.created) || 0;
                const updated = Number(data.updated) || 0;
                const failed = Number(data.failed) || 0;
                const total = Number(data.total) || created + updated + failed;
                setActionResult(prev => ({
                    ...prev,
                    [id]: {
                        ok: failed === 0 && total > 0,
                        message: buildApplyMessage(total, created, updated, failed, data.failed_items),
                    },
                }));
                patchNode(id, failed > 0 ? 'error' : 'done');
                applyCompletion.current.get(id)?.();
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
                            channelNoSku: Number(data.channel_no_sku) || 0,
                            skuChanged: Number(data.sku_changed) || 0,
                            skuTypo: Number(data.sku_typo) || 0,
                            matchRules: Array.isArray(data.match_rules) ? (data.match_rules as ProductMatchRuleSummary[]) : undefined,
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

    const drainPending = useCallback((correlationID: string) => {
        const buffered = pendingEvents.current.get(correlationID);
        if (!buffered) return;
        pendingEvents.current.delete(correlationID);
        const id = corrToIntegration.current.get(correlationID);
        if (id === undefined) return;
        for (const item of buffered) applyEvent(id, item.eventType, item.data);
    }, [applyEvent]);

    const bufferEvent = useCallback((correlationID: string, eventType: string, data: any) => {
        const buffered = pendingEvents.current.get(correlationID) || [];
        if (buffered.length >= MAX_PENDING_EVENTS_PER_RUN) return;
        buffered.push({ eventType, data });
        pendingEvents.current.set(correlationID, buffered);
        while (pendingEvents.current.size > MAX_PENDING_RUNS) {
            const oldest = pendingEvents.current.keys().next().value;
            if (oldest === undefined) break;
            pendingEvents.current.delete(oldest);
        }
    }, []);

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data);
            const eventType: string = parsed.type || parsed.metadata?.event_type || '';
            const data = parsed.data;
            const correlationID: string = data?.correlation_id ? String(data.correlation_id) : '';
            if (correlationID === '') return;
            const id = corrToIntegration.current.get(correlationID);
            if (id === undefined) {
                bufferEvent(correlationID, eventType, data);
                return;
            }
            applyEvent(id, eventType, data);
        } catch {
            return;
        }
    }, [applyEvent, bufferEvent]);

    const hydrateFromLastRun = useCallback(async (
        integrationId: number,
        kind: SyncRunKind,
        launchedAt: number,
    ): Promise<boolean> => {
        try {
            const rows = await listSyncRunsAction(businessId ?? undefined);
            const row = rows.find(r => r.integration_id === integrationId && r.kind === kind);
            if (!row?.finished_at) return false;
            const finishedAt = new Date(row.finished_at).getTime();
            if (!Number.isFinite(finishedAt) || finishedAt < launchedAt - RUN_CLOCK_SKEW_MS) return false;

            if (kind === 'products') {
                setResults(prev => ({
                    ...prev,
                    [integrationId]: {
                        kind: 'products',
                        matched: row.matched || 0,
                        notAssociated: row.not_associated || 0,
                        onlyInProbability: row.only_in_probability || 0,
                        onlyInChannel: row.only_in_channel || 0,
                        channelNoSku: row.channel_no_sku || 0,
                        skuChanged: row.sku_changed || 0,
                        skuTypo: row.sku_typo || 0,
                    },
                }));
            } else {
                const total = row.total || 0;
                setProgress(prev => ({ ...prev, [integrationId]: { processed: total, total } }));
                setResults(prev => ({
                    ...prev,
                    [integrationId]: {
                        kind: 'inventory',
                        total,
                        updated: row.updated || 0,
                        unchanged: row.unchanged || 0,
                        skipped: row.skipped || 0,
                        failed: row.failed || 0,
                    },
                }));
            }
            patchNode(integrationId, (row.failed || 0) > 0 ? 'error' : 'done');
            return true;
        } catch {
            return false;
        }
    }, [businessId, patchNode]);

    useSSE({
        businessId: businessId ?? 0,
        eventTypes: GLOBAL_INVENTORY_EVENT_TYPES,
        onMessage: handleMessage,
        enabled: true,
    });

    const reset = useCallback(() => {
        pendingEvents.current.clear();
        setMode('idle');
        setRunning(false);
        setNodes({});
        setProgress({});
        setResults({});
        setDetails({});
        setActionResult({});
    }, []);

    const syncInventoryOne = useCallback(async (integration: Integration, skus?: string[]) => {
        const provider = getSyncProvider(integration.integration_type_id);
        if (!provider) return;

        patchNode(integration.id, 'active');
        const launchedAt = Date.now();
        let result: SyncStartResult | null = null;
        try {
            result = await provider.syncInventory(integration.id, businessId ?? undefined, skus) as SyncStartResult;
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

        const correlationID = result.correlation_id;
        corrToIntegration.current.set(correlationID, integration.id);

        await new Promise<void>(resolve => {
            const timer = window.setTimeout(async () => {
                completion.current.delete(integration.id);
                const recovered = await hydrateFromLastRun(integration.id, 'inventory', launchedAt);
                if (!recovered) {
                    patchNode(integration.id, 'error');
                    setResults(prev => ({
                        ...prev,
                        [integration.id]: { kind: 'error', message: 'Continua en segundo plano' },
                    }));
                }
                resolve();
            }, SYNC_TIMEOUT_MS);
            completion.current.set(integration.id, () => {
                window.clearTimeout(timer);
                completion.current.delete(integration.id);
                resolve();
            });
            drainPending(correlationID);
        });
    }, [businessId, patchNode, drainPending, hydrateFromLastRun]);

    const runInventoryOne = useCallback(async (integrationId: number, skus?: string[]) => {
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

        await syncInventoryOne(integration, skus);

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

        const launchedAt = Date.now();
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

        const correlationID = start.correlation_id;
        corrToIntegration.current.set(correlationID, integration.id);

        await new Promise<void>(resolve => {
            const timer = window.setTimeout(async () => {
                completion.current.delete(integration.id);
                const recovered = await hydrateFromLastRun(integration.id, 'products', launchedAt);
                if (!recovered) {
                    patchNode(integration.id, 'error');
                    setResults(prev => ({
                        ...prev,
                        [integration.id]: { kind: 'error', message: 'Continua en segundo plano' },
                    }));
                }
                resolve();
            }, SYNC_TIMEOUT_MS);
            completion.current.set(integration.id, () => {
                window.clearTimeout(timer);
                completion.current.delete(integration.id);
                resolve();
            });
            drainPending(correlationID);
        });
    }, [businessId, patchNode, drainPending, hydrateFromLastRun]);

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
                ? (provider.associateProducts
                    ? (id: number, bid?: number, only?: string[]) => provider.associateProducts!(id, bid, only)
                    : undefined)
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

                const correlationID = typeof res?.correlation_id === 'string' ? res.correlation_id : '';
                if (correlationID === '') continue;

                corrToIntegration.current.set(correlationID, integrationId);
                setActionResult(prev => ({ ...prev, [integrationId]: { ok: true, pending: true, message: 'Aplicando...' } }));
                await new Promise<void>(resolve => {
                    const timer = window.setTimeout(() => {
                        applyCompletion.current.delete(integrationId);
                        setActionResult(prev => ({
                            ...prev,
                            [integrationId]: { ok: true, pending: true, message: 'Sigue aplicandose en segundo plano' },
                        }));
                        resolve();
                    }, APPLY_TIMEOUT_MS);
                    applyCompletion.current.set(integrationId, () => {
                        window.clearTimeout(timer);
                        applyCompletion.current.delete(integrationId);
                        resolve();
                    });
                    drainPending(correlationID);
                });
            }
            setActionResult(prev => (prev[integrationId]?.message === 'Aplicando...' || !prev[integrationId]
                ? { ...prev, [integrationId]: { ok, message } }
                : prev));
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
    }, [eligible, businessId, actionBusy, reconcileOne, loadLastRuns, drainPending]);

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
