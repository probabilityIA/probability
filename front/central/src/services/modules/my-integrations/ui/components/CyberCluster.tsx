'use client';

import type { Integration, IntegrationCategory } from '@/services/integrations/core/domain/types';
import { CyberIntegrationNode } from './CyberIntegrationNode';

interface CyberClusterProps {
    category: IntegrationCategory;
    color: string;
    integrations: Integration[];
    onToggle: (integration: Integration) => void;
    onEdit: (integration: Integration) => void;
    togglingId: number | null;
    editingId: number | null;
    anchorRef: (el: HTMLDivElement | null) => void;
}

export function CyberCluster({
    category,
    color,
    integrations,
    onToggle,
    onEdit,
    togglingId,
    editingId,
    anchorRef,
}: CyberClusterProps) {
    return (
        <div ref={anchorRef} className="relative min-w-[15rem] flex-1">
            <div className="rounded-2xl border border-gray-200 bg-white px-3 pb-3 pt-5 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                <div className="flex min-h-[3rem] flex-col justify-center gap-2">
                    {integrations.length === 0 ? (
                        <p className="py-2 text-center text-xs italic text-gray-400 dark:text-gray-500">
                            Sin configurar
                        </p>
                    ) : (
                        integrations.map(integration => (
                            <CyberIntegrationNode
                                key={integration.id}
                                integration={integration}
                                color={color}
                                onToggle={onToggle}
                                onEdit={onEdit}
                                togglingId={togglingId}
                                editingId={editingId}
                            />
                        ))
                    )}
                </div>
            </div>

            <div
                className="absolute -top-3 left-1/2 z-10 flex -translate-x-1/2 items-center gap-1.5 whitespace-nowrap rounded-full px-3 py-1 text-xs font-semibold text-white shadow-sm"
                style={{ backgroundColor: color }}
            >
                {category.name}
            </div>
        </div>
    );
}
