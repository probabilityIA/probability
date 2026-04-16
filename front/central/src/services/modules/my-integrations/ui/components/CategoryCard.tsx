'use client';

import type { Integration, IntegrationCategory } from '@/services/integrations/core/domain/types';
import { CATEGORY_ICONS } from '../../domain/types';
import { IntegrationToggle } from './IntegrationToggle';

interface CategoryCardProps {
    category: IntegrationCategory;
    integrations: Integration[];
    onToggle: (integration: Integration) => void;
    onEdit?: (integration: Integration) => void;
    togglingId: number | null;
}

export function CategoryCard({ category, integrations, onToggle, onEdit, togglingId }: CategoryCardProps) {
    const icon = CATEGORY_ICONS[category.code] || '🔗';
    const color = category.color || '#8B5CF6';

    return (
        <div className="flex-1 min-w-[13rem] rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden shadow-sm">
            <div
                className="px-3 py-2 text-white font-semibold flex items-center gap-2 text-sm"
                style={{ backgroundColor: color }}
            >
                <span>{icon}</span>
                <span className="truncate">{category.name}</span>
            </div>

            <div className="p-2 space-y-1 bg-white dark:bg-gray-800 min-h-[60px]">
                {integrations.length === 0 ? (
                    <p className="text-xs text-gray-400 dark:text-gray-500 dark:text-gray-400 italic text-center py-4">Sin configurar</p>
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
    );
}
