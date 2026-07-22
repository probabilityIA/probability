'use client';

import type { CSSProperties } from 'react';
import type { Integration } from '@/services/integrations/core/domain/types';
import type { IntegrationStatsItem } from '@/services/integrations/core/infra/actions/stats';
import { Clock, SlidersHorizontal } from 'lucide-react';

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

    return (
        <div
            className="group relative overflow-hidden rounded-2xl p-px shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:shadow-[0_6px_18px_-8px_var(--neon)]"
            style={{ '--neon': color } as CSSProperties}
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
            <div className="relative z-10 flex h-full items-center gap-4 rounded-[15px] bg-white p-3 transition-colors group-hover:bg-gray-50/80 dark:bg-gray-800 dark:group-hover:bg-gray-700/60">
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
                    <p className="truncate text-sm font-bold leading-tight text-gray-900 dark:text-white">{typeName}</p>
                    <p className="truncate text-[11px] leading-tight text-gray-400 dark:text-gray-500">{integration.name}</p>
                    {lastOrder && (
                        <p className="flex items-center gap-1 text-[10px] leading-tight text-gray-400 dark:text-gray-500">
                            <Clock size={9} className="flex-shrink-0" />
                            <span className="truncate">Ult. orden {lastOrder}</span>
                        </p>
                    )}
                </div>
            </div>

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

            <div className="flex min-w-0 flex-1 flex-col justify-center gap-1.5">
                {hasBreakdown ? (
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
                                <span key={bucket.key} className="flex items-center gap-1 text-[11px] font-medium text-gray-600 dark:text-gray-300">
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
        </div>
    );
}
