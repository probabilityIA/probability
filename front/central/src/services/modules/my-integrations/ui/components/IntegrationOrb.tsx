'use client';

import type { Integration, IntegrationCategory } from '@/services/integrations/core/domain/types';
import { CATEGORY_ICONS } from '../../domain/types';
import { IntegrationToggle } from './IntegrationToggle';

interface IntegrationOrbProps {
    category: IntegrationCategory;
    integrations: Integration[];
    onToggle: (integration: Integration) => void;
    onEdit?: (integration: Integration) => void;
    togglingId: number | null;
}

export function IntegrationOrb({ category, integrations, onToggle, onEdit, togglingId }: IntegrationOrbProps) {
    const icon = CATEGORY_ICONS[category.code] || '🔗';
    const color = category.color || '#8B5CF6';
    const isEmpty = integrations.length === 0;

    return (
        <div className="relative flex-1 min-w-[12rem]">
            <div
                className="absolute -top-3 left-1/2 -translate-x-1/2 z-10 px-3 py-1 rounded-full text-white text-xs font-semibold shadow-md flex items-center gap-1.5 whitespace-nowrap"
                style={{ backgroundColor: color }}
            >
                <span>{icon}</span>
                <span>{category.name}</span>
            </div>

            <div
                className="rounded-[2rem] bg-white dark:bg-gray-800 pt-5 pb-3 px-3 shadow-md border-2 transition-all hover:shadow-lg"
                style={{ borderColor: `${color}33` }}
            >
                <div className="space-y-1 min-h-[3rem] flex flex-col justify-center">
                    {isEmpty ? (
                        <p className="text-xs text-gray-400 dark:text-gray-500 italic text-center py-2">Sin configurar</p>
                    ) : (
                        integrations.map(integration => (
                            <IntegrationToggle
                                key={integration.id}
                                integration={integration}
                                onToggle={onToggle}
                                onEdit={onEdit}
                                togglingId={togglingId}
                            />
                        ))
                    )}
                </div>
            </div>
        </div>
    );
}
