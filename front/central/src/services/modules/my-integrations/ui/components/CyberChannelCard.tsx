'use client';

import { useState, type CSSProperties } from 'react';
import type { Integration } from '@/services/integrations/core/domain/types';
import type { IntegrationStatsItem } from '@/services/integrations/core/infra/actions/stats';
import { Clock, SlidersHorizontal, ChevronDown } from 'lucide-react';
import { useSyncActivity } from '../sync-activity-context';

interface CyberChannelCardProps {
    integration: Integration;
    color: string;
    stats?: IntegrationStatsItem;
    onToggle: (integration: Integration) => void;
    onEdit: (integration: Integration) => void;
    togglingId: number | null;
    editingId: number | null;
}

const numberFormat = new Intl.NumberFormat('es-CO');

const BUCKETS = [
    { key: 'orders_in_progress', label: 'en curso', dot: '#3b82f6' },
    { key: 'orders_delivered', label: 'entregadas', dot: '#22c55e' },
    { key: 'orders_cancelled', label: 'canceladas', dot: '#ef4444' },
    { key: 'orders_returned', label: 'devueltas', dot: '#f59e0b' },
] as const;

function relativeTime(iso?: string): string | null {
    if (!iso) return null;
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) return null;
    const diffMs = Date.now() - date.getTime();
    const minutes = Math.floor(diffMs / 60000);
    if (minutes < 1) return 'hace un momento';
    if (minutes < 60) return `hace ${minutes} min`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `hace ${hours} h`;
    const days = Math.floor(hours / 24);
    if (days < 30) return `hace ${days} d`;
    const months = Math.floor(days / 30);
    if (months < 12) return `hace ${months} mes${months > 1 ? 'es' : ''}`;
    return `hace ${Math.floor(months / 12)} a`;
}

export function CyberChannelCard({ integration, color, stats, onToggle, onEdit, togglingId, editingId }: CyberChannelCardProps) {
    const isToggling = togglingId === integration.id;
    const isEditing = editingId === integration.id;
    const active = integration.is_active;
    const typeName = integration.integration_type?.name || integration.name;
    const lastOrder = relativeTime(stats?.last_order_at);
    const total = stats?.orders_count ?? 0;
    const hasBreakdown = stats !== undefined && total > 0;

    const { nodes, progress, results, details } = useSyncActivity();
    const [openDetail, setOpenDetail] = useState(false);
    const syncState = nodes[integration.id] || 'idle';
    const syncProgress = progress[integration.id] || { processed: 0, total: 0 };
    const syncResult = results[integration.id];
    const syncBusy = syncState !== 'idle' || syncResult !== undefined;
    const syncDetail = details[integration.id] || [];
    const syncPct = syncProgress.total > 0
        ? Math.min(100, Math.round((syncProgress.processed / syncProgress.total) * 100))
        : 0;
    const stateRing =
        syncState === 'active' ? '0 0 0 4px rgba(59,130,246,.16), 0 12px 32px rgba(37,99,235,.22)'
        : syncState === 'scan' ? '0 0 0 4px rgba(139,92,246,.16), 0 12px 32px rgba(124,58,237,.20)'
        : syncState === 'done' ? '0 0 0 3px rgba(34,197,94,.16)'
        : undefined;

    return (
        <div
            data-node-id={integration.id}
            className="group relative overflow-hidden rounded-2xl p-px shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:shadow-[0_6px_18px_-8px_var(--neon)]"
            style={{
                '--neon': color,
                boxShadow: stateRing,
                opacity: syncState === 'queued' ? 0.6 : 1,
            } as CSSProperties}
        >
            <div className="absolute inset-0" style={{ backgroundColor: `${color}2c` }} />
            <div
                className="absolute inset-0"
                style={{
                    background: `linear-gradient(110deg, transparent 25%, ${color} 50%, transparent 75%)`,
                    backgroundSize: '250% 100%',
                    animation: 'cyber-sweep 3.2s linear infinite',
                    animationDelay: `${(integration.id % 5) * -0.65}s`,
                }}
            />
            <div className="relative z-10 flex h-full flex-col rounded-[15px] bg-white transition-colors group-hover:bg-gray-50/80 dark:bg-gray-800 dark:group-hover:bg-gray-700/60">
            <div className="flex items-center gap-4 p-3">
            <div className="flex w-48 flex-shrink-0 items-center gap-2.5" title={lastOrder ? `Ultima orden ${lastOrder}` : undefined}>
                <div className="relative flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-xl bg-gray-50 ring-1 ring-gray-200 dark:bg-gray-700/60 dark:ring-gray-600">
                    {integration.integration_type?.image_url ? (
                        <img
                            src={integration.integration_type.image_url}
                            alt={typeName}
                            className="h-7 w-7 rounded-md object-contain"
                        />
                    ) : (
                        <span className="text-sm font-bold text-gray-500 dark:text-gray-300">{typeName.charAt(0).toUpperCase()}</span>
                    )}
                    <span
                        className={`absolute -right-1 -top-1 h-2.5 w-2.5 rounded-full border-2 border-white dark:border-gray-800 ${active ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'}`}
                    />
                </div>
                <div className="min-w-0">
                    <p className="flex items-center gap-1.5 truncate text-sm font-bold leading-tight text-gray-900 dark:text-white">
                        <span className="truncate">{typeName}</span>
                        {integration.is_testing && (
                            <span
                                title="Integracion en modo pruebas: apunta al simulador, no a la tienda real"
                                className="flex-shrink-0 rounded-md border border-amber-300 bg-amber-100 px-1.5 py-px text-[9px] font-extrabold uppercase tracking-wider text-amber-700 dark:border-amber-500/40 dark:bg-amber-500/20 dark:text-amber-300"
                            >
                                Test
                            </span>
                        )}
                    </p>
                    <p className="truncate text-[11px] leading-tight text-gray-400 dark:text-gray-500">{integration.name}</p>
                    {lastOrder && (
                        <p className="flex items-center gap-1 text-[10px] leading-tight text-gray-400 dark:text-gray-500">
                            <Clock size={9} className="flex-shrink-0" />
                            <span className="truncate">Ult. orden {lastOrder}</span>
                        </p>
                    )}
                </div>
            </div>

            {!syncBusy && (
            <div className="flex flex-shrink-0 gap-5">
                <div>
                    <p className="text-[9px] font-semibold uppercase tracking-[0.15em] text-gray-400 dark:text-gray-500">Ordenes</p>
                    <p className="text-[22px] font-bold leading-none tabular-nums text-gray-900 dark:text-white">
                        {stats ? numberFormat.format(stats.orders_count) : '-'}
                    </p>
                </div>
                <div>
                    <p className="text-[9px] font-semibold uppercase tracking-[0.15em] text-gray-400 dark:text-gray-500">Productos</p>
                    <p className="text-[22px] font-bold leading-none tabular-nums text-gray-900 dark:text-white">
                        {stats ? numberFormat.format(stats.products_count) : '-'}
                    </p>
                </div>
            </div>
            )}

            <div className="flex min-w-0 flex-1 flex-col justify-center gap-1.5">
                {syncState === 'queued' ? (
                    <div className="flex items-center justify-end gap-2">
                        <span className="h-2 w-2 animate-pulse rounded-full bg-gray-300 dark:bg-gray-600" />
                        <span className="text-xs font-semibold text-gray-400 dark:text-gray-500">En cola</span>
                    </div>
                ) : syncState === 'active' ? (
                    <>
                        <div className="flex items-center justify-between text-xs">
                            <span className="flex items-center gap-1.5 font-bold text-blue-600 dark:text-blue-400">
                                <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-blue-500" />
                                Sincronizando inventario
                            </span>
                            <span className="font-bold tabular-nums text-blue-600 dark:text-blue-400">
                                {syncProgress.processed}/{syncProgress.total || '?'}
                            </span>
                        </div>
                        <div className="h-2 w-full overflow-hidden rounded-full bg-blue-100 dark:bg-blue-900/40">
                            <span
                                className="block h-full rounded-full bg-gradient-to-r from-blue-600 to-cyan-400 transition-[width] duration-150"
                                style={{ width: `${syncPct}%` }}
                            />
                        </div>
                    </>
                ) : syncState === 'scan' ? (
                    <>
                        <span className="flex items-center gap-1.5 text-[11px] font-bold text-violet-600 dark:text-violet-400">
                            <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-violet-500" />
                            Comparando catalogo...
                        </span>
                        <div
                            className="h-2 w-full rounded-full"
                            style={{
                                background: 'linear-gradient(90deg,#ede9fe 25%,#c4b5fd 50%,#ede9fe 75%)',
                                backgroundSize: '200% 100%',
                                animation: 'cyber-shimmer 1.1s linear infinite',
                            }}
                        />
                    </>
                ) : syncResult?.kind === 'inventory' ? (
                    <div className="flex flex-wrap items-center justify-end gap-1.5">
                        <span className="rounded-full bg-emerald-50 px-2 py-0.5 text-[10.5px] font-bold text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300">
                            {syncResult.updated} actualizados
                        </span>
                        {syncResult.unchanged > 0 && (
                            <span className="rounded-full bg-gray-100 px-2 py-0.5 text-[10.5px] font-bold text-gray-500 dark:bg-gray-700 dark:text-gray-300">
                                {syncResult.unchanged} sin cambio
                            </span>
                        )}
                        {syncResult.failed > 0 && (
                            <span className="rounded-full bg-red-50 px-2 py-0.5 text-[10.5px] font-bold text-red-600 dark:bg-red-900/30 dark:text-red-300">
                                {syncResult.failed} fallidos
                            </span>
                        )}
                    </div>
                ) : syncResult?.kind === 'products' ? (
                    <div className="flex flex-wrap items-center justify-end gap-1.5">
                        <span className="rounded-full bg-emerald-50 px-2 py-0.5 text-[10.5px] font-bold text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300">
                            {syncResult.matched} ok
                        </span>
                        {syncResult.notAssociated > 0 && (
                            <span className="rounded-full bg-blue-50 px-2 py-0.5 text-[10.5px] font-bold text-blue-700 dark:bg-blue-900/30 dark:text-blue-300">
                                {syncResult.notAssociated} sin asociar
                            </span>
                        )}
                        <span className="rounded-full bg-orange-50 px-2 py-0.5 text-[10.5px] font-bold text-orange-700 dark:bg-orange-900/30 dark:text-orange-300">
                            {syncResult.onlyInProbability} solo Prob.
                        </span>
                        <span className="rounded-full bg-fuchsia-50 px-2 py-0.5 text-[10.5px] font-bold text-fuchsia-700 dark:bg-fuchsia-900/30 dark:text-fuchsia-300">
                            {syncResult.onlyInChannel} solo canal
                        </span>
                    </div>
                ) : syncResult?.kind === 'error' ? (
                    <span className="text-right text-[11px] font-semibold text-red-500">{syncResult.message}</span>
                ) : hasBreakdown ? (
                    <>
                        <div className="flex h-1.5 w-full overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700">
                            {BUCKETS.map(bucket => {
                                const value = stats[bucket.key];
                                if (value <= 0) return null;
                                return (
                                    <span
                                        key={bucket.key}
                                        style={{ width: `${(value / total) * 100}%`, backgroundColor: bucket.dot }}
                                    />
                                );
                            })}
                        </div>
                        <div className="flex flex-wrap items-center gap-x-3 gap-y-0.5">
                            {BUCKETS.map(bucket => (
                                <span key={bucket.key} className="flex items-center gap-1 whitespace-nowrap text-[11px] font-medium text-gray-600 dark:text-gray-300">
                                    <span className="h-1.5 w-1.5 flex-shrink-0 rounded-full" style={{ backgroundColor: bucket.dot }} />
                                    {numberFormat.format(stats[bucket.key])} {bucket.label}
                                </span>
                            ))}
                        </div>
                    </>
                ) : (
                    <p className="text-[11px] italic text-gray-300 dark:text-gray-600">Sin ordenes registradas</p>
                )}
            </div>

            <div className="flex flex-shrink-0 flex-col items-center gap-2 border-l border-gray-100 pl-3 dark:border-gray-700">
                <button
                    onClick={() => onEdit(integration)}
                    disabled={isEditing}
                    title="Configurar integracion"
                    className={`flex h-8 w-8 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 shadow-sm transition-all hover:border-gray-300 hover:bg-gray-50 hover:text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 ${isEditing ? 'cursor-wait opacity-60' : ''}`}
                >
                    {isEditing ? (
                        <span className="block h-4 w-4 animate-spin rounded-full border border-transparent border-t-current" />
                    ) : (
                        <SlidersHorizontal className="h-4 w-4" />
                    )}
                </button>
                <button
                    onClick={() => onToggle(integration)}
                    disabled={isToggling}
                    title={active ? 'Desactivar' : 'Activar'}
                    className={`relative inline-flex h-5 w-10 items-center rounded-full transition-colors ${
                        active ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'
                    } ${isToggling ? 'cursor-wait opacity-50' : 'cursor-pointer'}`}
                >
                    <span
                        className={`inline-block h-3.5 w-3.5 rounded-full bg-white shadow-sm transition-transform duration-200 ${
                            active ? 'translate-x-6' : 'translate-x-1'
                        }`}
                    />
                </button>
            </div>
            </div>
            {syncDetail.length > 0 && (
                <div className="border-t border-gray-100 px-3 pb-2 dark:border-gray-700">
                    <button
                        onClick={() => setOpenDetail(v => !v)}
                        className="flex w-full items-center justify-center gap-1 py-1.5 text-[11px] font-semibold text-gray-500 transition-colors hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-100"
                    >
                        <ChevronDown size={12} className={`transition-transform ${openDetail ? 'rotate-180' : ''}`} />
                        {openDetail ? 'Ocultar detalle' : `Ver detalle (${syncDetail.length})`}
                    </button>
                    {openDetail && (
                        <div className="max-h-40 overflow-y-auto rounded-lg bg-gray-50 p-1.5 dark:bg-gray-900/40">
                            {syncDetail.map((d, i) => (
                                <div
                                    key={`${d.sku}-${i}`}
                                    className="flex items-start gap-2 border-b border-gray-100 px-1 py-1 text-[11px] last:border-0 dark:border-gray-700/60"
                                >
                                    <span
                                        className={`mt-1 h-1.5 w-1.5 flex-shrink-0 rounded-full ${
                                            d.tone === 'error' ? 'bg-red-500' : d.tone === 'warn' ? 'bg-amber-500' : 'bg-emerald-500'
                                        }`}
                                    />
                                    <span className="w-28 flex-shrink-0 truncate font-mono font-semibold text-gray-700 dark:text-gray-200">{d.sku}</span>
                                    <span className={`min-w-0 flex-1 truncate ${d.tone === 'error' ? 'text-red-600 dark:text-red-400' : 'text-gray-500 dark:text-gray-400'}`}>
                                        {d.label}
                                    </span>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            )}
            </div>
        </div>
    );
}