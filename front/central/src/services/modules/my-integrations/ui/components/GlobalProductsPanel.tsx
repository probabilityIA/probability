'use client';

import { useCallback, useEffect, useState } from 'react';
import { RefreshCw, AlertCircle, Loader2, ArrowRightLeft, Table2, List } from 'lucide-react';
import type { Integration } from '@/services/integrations/core/domain/types';
import { getSyncProvider } from '../providers';
import { GlobalProductsMatrix, type ChannelDiffData, type ProductBrief } from './GlobalProductsMatrix';

interface GlobalProductsPanelProps {
    integrations: Integration[];
    businessId: number | null;
}

interface DiffSummary {
    status: 'analyzing' | 'ready' | 'error';
    matched: number;
    notAssociated: ProductBrief[];
    onlyInProbability: ProductBrief[];
    onlyInChannel: ProductBrief[];
    message?: string;
}

const ANALYZING: DiffSummary = { status: 'analyzing', matched: 0, notAssociated: [], onlyInProbability: [], onlyInChannel: [] };

function toBriefs(value: unknown): ProductBrief[] {
    if (!Array.isArray(value)) return [];
    return value.map(item => ({
        sku: String((item as ProductBrief)?.sku || ''),
        name: String((item as ProductBrief)?.name || ''),
    })).filter(b => b.sku !== '');
}

export function GlobalProductsPanel({ integrations, businessId }: GlobalProductsPanelProps) {
    const [rows, setRows] = useState<Record<number, DiffSummary>>({});
    const [reviewing, setReviewing] = useState<Integration | null>(null);
    const [view, setView] = useState<'channels' | 'matrix'>('channels');

    const analyzeOne = useCallback(async (integration: Integration) => {
        const provider = getSyncProvider(integration.integration_type_id);
        if (!provider) return;
        setRows(prev => ({ ...prev, [integration.id]: ANALYZING }));
        let summary: DiffSummary;
        try {
            const res = await provider.reconcileProducts(integration.id, businessId ?? undefined) as Record<string, unknown>;
            if (res?.success === false) {
                summary = { ...ANALYZING, status: 'error', message: String(res?.message || 'No se pudo comparar') };
            } else {
                summary = {
                    status: 'ready',
                    matched: Number(res?.matched) || 0,
                    notAssociated: toBriefs(res?.matched_not_associated),
                    onlyInProbability: toBriefs(res?.only_in_probability),
                    onlyInChannel: toBriefs(res?.[provider.onlyInChannelField]),
                };
            }
        } catch {
            summary = { ...ANALYZING, status: 'error', message: 'No se pudo comparar' };
        }
        setRows(prev => ({ ...prev, [integration.id]: summary }));
    }, [businessId]);

    const analyzeAll = useCallback(() => {
        integrations.forEach(integration => { analyzeOne(integration); });
    }, [integrations, analyzeOne]);

    useEffect(() => {
        analyzeAll();
    }, [analyzeAll]);

    if (integrations.length === 0) {
        return (
            <p className="py-4 text-center text-xs italic text-gray-400 dark:text-gray-500">
                No hay integraciones e-commerce activas para comparar
            </p>
        );
    }

    const ReviewModal = reviewing ? getSyncProvider(reviewing.integration_type_id)?.ProductSyncModal : null;

    const readyChannels: ChannelDiffData[] = integrations
        .filter(i => rows[i.id]?.status === 'ready')
        .map(i => ({
            integration: i,
            notAssociated: rows[i.id].notAssociated,
            onlyInProbability: rows[i.id].onlyInProbability,
            onlyInChannel: rows[i.id].onlyInChannel,
        }));

    return (
        <div className="flex flex-col gap-2">
            <div className="flex items-center justify-between gap-2">
                <p className="text-xs text-gray-500 dark:text-gray-400">
                    Compara el catalogo de cada canal contra Probability.
                </p>
                <div className="flex items-center gap-1">
                    <button
                        onClick={() => setView(view === 'matrix' ? 'channels' : 'matrix')}
                        disabled={readyChannels.length === 0}
                        className={`flex items-center gap-1 rounded-md px-2 py-1 text-[11px] font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-40 ${
                            view === 'matrix'
                                ? 'bg-indigo-50 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300'
                                : 'text-gray-500 hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-400'
                        }`}
                    >
                        {view === 'matrix' ? <List size={12} /> : <Table2 size={12} />}
                        {view === 'matrix' ? 'Ver por canal' : 'Ver matriz'}
                    </button>
                    <button
                        onClick={analyzeAll}
                        className="flex items-center gap-1 rounded-md px-2 py-1 text-[11px] font-semibold text-gray-500 transition-colors hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-400"
                    >
                        <RefreshCw size={12} /> Reanalizar
                    </button>
                </div>
            </div>

            {view === 'matrix' ? (
                <GlobalProductsMatrix
                    channels={readyChannels}
                    businessId={businessId}
                    onAssociated={analyzeAll}
                />
            ) : (
                <div className="flex flex-col gap-1.5">
                    {integrations.map(integration => {
                        const row = rows[integration.id] || ANALYZING;
                        const typeName = integration.integration_type?.name || integration.name;
                        const pending = row.notAssociated.length + row.onlyInProbability.length + row.onlyInChannel.length;

                        return (
                            <div
                                key={integration.id}
                                className="flex items-center gap-3 rounded-lg border border-gray-200 px-3 py-2 dark:border-gray-700"
                            >
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
                                    {row.status === 'analyzing' && (
                                        <span className="flex items-center gap-1 text-[11px] text-gray-400 dark:text-gray-500">
                                            <Loader2 size={12} className="animate-spin" /> Comparando...
                                        </span>
                                    )}
                                    {row.status === 'error' && (
                                        <span className="flex items-center gap-1 text-[11px] text-red-500" title={row.message}>
                                            <AlertCircle size={13} /> {row.message || 'Error'}
                                        </span>
                                    )}
                                    {row.status === 'ready' && (
                                        <>
                                            <span className="hidden items-center gap-1.5 text-[11px] sm:flex">
                                                <span className="rounded-full bg-emerald-50 px-2 py-0.5 font-semibold text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
                                                    {row.matched} ok
                                                </span>
                                                {row.notAssociated.length > 0 && (
                                                    <span className="rounded-full bg-indigo-50 px-2 py-0.5 font-semibold text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300">
                                                        {row.notAssociated.length} sin asociar
                                                    </span>
                                                )}
                                                {row.onlyInProbability.length > 0 && (
                                                    <span className="rounded-full bg-amber-50 px-2 py-0.5 font-semibold text-amber-700 dark:bg-amber-900/40 dark:text-amber-300">
                                                        {row.onlyInProbability.length} solo Prob.
                                                    </span>
                                                )}
                                                {row.onlyInChannel.length > 0 && (
                                                    <span className="rounded-full bg-purple-50 px-2 py-0.5 font-semibold text-purple-700 dark:bg-purple-900/40 dark:text-purple-300">
                                                        {row.onlyInChannel.length} solo canal
                                                    </span>
                                                )}
                                                {pending === 0 && (
                                                    <span className="text-[11px] text-gray-400">al dia</span>
                                                )}
                                            </span>
                                            <button
                                                onClick={() => setReviewing(integration)}
                                                className="flex items-center gap-1 rounded-lg border border-indigo-200 px-2.5 py-1 text-[11px] font-semibold text-indigo-600 transition-colors hover:bg-indigo-50 dark:border-indigo-800 dark:text-indigo-300 dark:hover:bg-indigo-900/30"
                                            >
                                                <ArrowRightLeft size={12} /> Revisar
                                            </button>
                                        </>
                                    )}
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}

            {ReviewModal && reviewing && (
                <ReviewModal
                    isOpen
                    onClose={() => setReviewing(null)}
                    integrationId={reviewing.id}
                    businessId={businessId}
                    onCompleted={() => analyzeOne(reviewing)}
                />
            )}
        </div>
    );
}
