'use client';

import type { Integration } from '@/services/integrations/core/domain/types';
import type { IntegrationStatsItem } from '@/services/integrations/core/infra/actions/stats';
import { CATEGORY_COLORS } from '../../domain/types';
import { CyberChannelCard } from './CyberChannelCard';

interface CyberChannelsClusterProps {
    integrations: Integration[];
    stats: Record<number, IntegrationStatsItem>;
    statsLoaded: boolean;
    color: string;
    onToggle: (integration: Integration) => void;
    onEdit: (integration: Integration) => void;
    togglingId: number | null;
    editingId: number | null;
    anchorRef: (el: HTMLDivElement | null) => void;
}

export function CyberChannelsCluster({
    integrations,
    stats,
    statsLoaded,
    color,
    onToggle,
    onEdit,
    togglingId,
    editingId,
    anchorRef,
}: CyberChannelsClusterProps) {
    return (
        <div ref={anchorRef} className="relative w-full">
            <div className="rounded-2xl border border-gray-200 bg-white px-4 pb-4 pt-6 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                {integrations.length === 0 ? (
                    <p className="py-3 text-center text-xs italic text-gray-400 dark:text-gray-500">
                        Sin configurar
                    </p>
                ) : (
                    <div className="grid gap-2.5 xl:grid-cols-2">
                        {integrations.map(integration => (
                            <CyberChannelCard
                                key={integration.id}
                                integration={integration}
                                color={CATEGORY_COLORS[integration.category] || color}
                                stats={
                                    stats[integration.id] ??
                                    (statsLoaded
                                        ? {
                                            integration_id: integration.id,
                                            orders_count: 0,
                                            orders_in_progress: 0,
                                            orders_delivered: 0,
                                            orders_cancelled: 0,
                                            orders_returned: 0,
                                            products_count: 0,
                                        }
                                        : undefined)
                                }
                                onToggle={onToggle}
                                onEdit={onEdit}
                                togglingId={togglingId}
                                editingId={editingId}
                            />
                        ))}
                    </div>
                )}
            </div>

            <div
                className="absolute -top-3 left-1/2 z-10 flex -translate-x-1/2 items-center gap-1.5 whitespace-nowrap rounded-full px-3 py-1 text-xs font-semibold text-white shadow-sm"
                style={{ backgroundColor: color }}
            >
                Canales de venta
            </div>
        </div>
    );
}
