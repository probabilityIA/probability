'use client';

import { useCallback, useMemo, useRef, useState } from 'react';
import { RefreshCw, CheckCircle2, AlertCircle, Loader2 } from 'lucide-react';
import { useSSE } from '@/shared/hooks/use-sse';
import type { Integration } from '@/services/integrations/core/domain/types';
import { getSyncProvider, GLOBAL_INVENTORY_EVENT_TYPES } from '../providers';

interface GlobalInventoryPanelProps {
    integrations: Integration[];
    businessId: number | null;
    onCompleted?: () => void;
}

type RowStatus = 'idle' | 'queued' | 'starting' | 'running' | 'done' | 'error' | 'background';

interface RowState {
    status: RowStatus;
    total: number;
    processed: number;
    updated: number;
    unchanged: number;
    skipped: number;
    failed: number;
    message?: string;
}

const EMPTY_ROW: RowState = { status: 'idle', total: 0, processed: 0, updated: 0, unchanged: 0, skipped: 0, failed: 0 };

const SYNC_TIMEOUT_MS = 6 * 60 * 1000;

export function GlobalInventoryPanel({ integrations, businessId, onCompleted }: GlobalInventoryPanelProps) {
    const [selected, setSelected] = useState<Set<number>>(() => new Set(integrations.map(i => i.id)));
    const [rows, setRows] = useState<Record<number, RowState>>({});
    const [running, setRunning] = useState(false);

    const corrToIntegrationRef = useRef<Map<string, number>>(new Map());
    const completionRef = useRef<Map<number, () => void>>(new Map());

    const patchRow = useCallback((integrationId: number, patch: Partial<RowState>) => {
        setRows(prev => ({ ...prev, [integrationId]: { ...(prev[integrationId] || EMPTY_ROW), ...patch } }));
    }, []);

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data);
            const eventType: string = parsed.type || parsed.metadata?.event_type || '';
            const data = parsed.data;
            if (!data?.correlation_id) return;
            const integrationId = corrToIntegrationRef.current.get(data.correlation_id);
            if (integrationId === undefined) return;

            if (eventType.endsWith('.inventory.sync.started')) {
                patchRow(integrationId, { status: 'running', total: Number(data.total) || 0 });
            } else if (eventType.endsWith('.inventory.sync.progress')) {
                patchRow(integrationId, {
                    status: 'running',
                    processed: Number(data.processed) || 0,
                    updated: Number(data.updated) || 0,
                    unchanged: Number(data.unchanged) || 0,
                    skipped: Number(data.skipped) || 0,
                    failed: Number(data.failed) || 0,
                });
            } else if (eventType.endsWith('.inventory.sync.completed')) {
                patchRow(integrationId, {
                    status: 'done',
                    total: Number(data.total) || 0,
                    processed: Number(data.total) || 0,
                    updated: Number(data.updated) || 0,
                    unchanged: Number(data.unchanged) || 0,
                    skipped: Number(data.skipped) || 0,
                    failed: Number(data.failed) || 0,
                });
                completionRef.current.get(integrationId)?.();
            }
        } catch {
            return;
        }
    }, [patchRow]);

    useSSE({
        businessId: businessId ?? 0,
        eventTypes: GLOBAL_INVENTORY_EVENT_TYPES,
        onMessage: handleMessage,
        enabled: true,
    });

    const toggleSelected = (id: number) => {
        if (running) return;
        setSelected(prev => {
            const next = new Set(prev);
            if (next.has(id)) next.delete(id);
            else next.add(id);
            return next;
        });
    };

    const queue = useMemo(() => integrations.filter(i => selected.has(i.id)), [integrations, selected]);

    const run = async () => {
        if (running || queue.length === 0) return;
        setRunning(true);
        corrToIntegrationRef.current.clear();

        const initial: Record<number, RowState> = {};
        for (const integration of integrations) {
            initial[integration.id] = selected.has(integration.id)
                ? { ...EMPTY_ROW, status: 'queued' }
                : { ...EMPTY_ROW };
        }
        setRows(initial);

        for (const integration of queue) {
            const provider = getSyncProvider(integration.integration_type_id);
            if (!provider) continue;

            patchRow(integration.id, { status: 'starting' });
            let result: { success?: boolean; correlation_id?: string; message?: string } | null = null;
            try {
                result = await provider.syncInventory(integration.id, businessId ?? undefined) as typeof result;
            } catch {
                result = null;
            }

            if (!result?.success || !result?.correlation_id) {
                patchRow(integration.id, { status: 'error', message: result?.message || 'No se pudo iniciar la sincronizacion' });
                continue;
            }

            corrToIntegrationRef.current.set(result.correlation_id, integration.id);

            await new Promise<void>(resolve => {
                const timer = setTimeout(() => {
                    patchRow(integration.id, { status: 'background', message: 'Sin respuesta del canal, continua en segundo plano' });
                    completionRef.current.delete(integration.id);
                    resolve();
                }, SYNC_TIMEOUT_MS);
                completionRef.current.set(integration.id, () => {
                    clearTimeout(timer);
                    completionRef.current.delete(integration.id);
                    resolve();
                });
            });
        }

        setRunning(false);
        onCompleted?.();
    };

    if (integrations.length === 0) {
        return (
            <p className="py-4 text-center text-xs italic text-gray-400 dark:text-gray-500">
                No hay integraciones e-commerce activas para sincronizar
            </p>
        );
    }

    return (
        <div className="flex flex-col gap-2">
            <p className="text-xs text-gray-500 dark:text-gray-400">
                Envia el stock de Probability a los canales seleccionados, uno por uno.
            </p>

            <div className="flex flex-col gap-1.5">
                {integrations.map(integration => {
                    const row = rows[integration.id] || EMPTY_ROW;
                    const typeName = integration.integration_type?.name || integration.name;
                    const pct = row.total > 0 ? Math.min(100, Math.round((row.processed / row.total) * 100)) : 0;

                    return (
                        <label
                            key={integration.id}
                            className={`flex items-center gap-3 rounded-lg border border-gray-200 px-3 py-2 transition-colors dark:border-gray-700 ${running ? '' : 'cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700/50'}`}
                        >
                            <input
                                type="checkbox"
                                checked={selected.has(integration.id)}
                                onChange={() => toggleSelected(integration.id)}
                                disabled={running}
                                className="h-4 w-4 rounded border-gray-300 text-cyan-600 focus:ring-cyan-500"
                            />
                            {integration.integration_type?.image_url ? (
                                <img
                                    src={integration.integration_type.image_url}
                                    alt={typeName}
                                    className="h-7 w-7 flex-shrink-0 rounded-full object-contain ring-1 ring-gray-200 dark:ring-gray-600"
                                />
                            ) : (
                                <div className="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-full bg-gray-100 text-xs font-bold text-gray-500 ring-1 ring-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:ring-gray-600">
                                    {typeName.charAt(0).toUpperCase()}
                                </div>
                            )}
                            <div className="flex min-w-0 flex-1 flex-col">
                                <span className="truncate text-xs font-semibold leading-tight text-gray-800 dark:text-gray-100">{typeName}</span>
                                <span className="truncate text-[11px] leading-tight text-gray-500 dark:text-gray-400">{integration.name}</span>
                            </div>

                            <div className="flex flex-shrink-0 items-center gap-2">
                                {row.status === 'queued' && (
                                    <span className="text-[11px] text-gray-400 dark:text-gray-500">En cola</span>
                                )}
                                {row.status === 'starting' && (
                                    <span className="flex items-center gap-1 text-[11px] text-cyan-600 dark:text-cyan-400">
                                        <Loader2 size={12} className="animate-spin" /> Iniciando...
                                    </span>
                                )}
                                {row.status === 'running' && (
                                    <span className="flex items-center gap-2 text-[11px] text-cyan-700 dark:text-cyan-300">
                                        <span className="h-1.5 w-20 overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700">
                                            <span
                                                className="block h-full rounded-full bg-cyan-500 transition-all duration-300"
                                                style={{ width: `${pct}%` }}
                                            />
                                        </span>
                                        <span className="tabular-nums">{row.processed}/{row.total}</span>
                                    </span>
                                )}
                                {row.status === 'done' && (
                                    <span className="flex items-center gap-1.5 text-[11px]">
                                        <CheckCircle2 size={13} className="text-emerald-500" />
                                        <span className="text-emerald-600 dark:text-emerald-400">{row.updated} act.</span>
                                        <span className="text-gray-400">{row.unchanged} s/c</span>
                                        {row.failed > 0 && <span className="text-red-500">{row.failed} fall.</span>}
                                    </span>
                                )}
                                {row.status === 'error' && (
                                    <span className="flex items-center gap-1 text-[11px] text-red-500" title={row.message}>
                                        <AlertCircle size={13} /> {row.message || 'Error'}
                                    </span>
                                )}
                                {row.status === 'background' && (
                                    <span className="flex items-center gap-1 text-[11px] text-amber-600 dark:text-amber-400" title={row.message}>
                                        <AlertCircle size={13} /> Continua en segundo plano
                                    </span>
                                )}
                            </div>
                        </label>
                    );
                })}
            </div>

            <div className="mt-1 flex justify-end">
                <button
                    onClick={run}
                    disabled={running || queue.length === 0}
                    className="flex items-center gap-2 rounded-lg bg-cyan-600 px-4 py-2 text-xs font-semibold text-white transition-colors hover:bg-cyan-700 disabled:cursor-not-allowed disabled:opacity-50"
                >
                    <RefreshCw size={14} className={running ? 'animate-spin' : ''} />
                    {running ? 'Sincronizando...' : `Sincronizar ${queue.length === integrations.length ? 'todas' : 'seleccionadas'} (${queue.length})`}
                </button>
            </div>
        </div>
    );
}
