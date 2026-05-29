'use client';

import type { Integration, IntegrationCategory } from '@/services/integrations/core/domain/types';
import { CATEGORY_ICONS, INTERNAL_MODULE_RESOURCE_NAME } from '../../domain/types';

interface InternalModulesOrbProps {
    category: IntegrationCategory;
    integrations: Integration[];
    resourceActive: Record<string, boolean>;
}

const MODULE_ICONS: Record<string, string> = {
    inventory: '📦',
    delivery: '🚚',
    notifications: '🔔',
    customers: '👥',
    storefront_module: '🛍️',
    invoicing_module: '🧾',
};

export function InternalModulesOrb({ category, integrations, resourceActive }: InternalModulesOrbProps) {
    const icon = CATEGORY_ICONS[category.code] || '⚙️';
    const color = category.color || '#6366F1';

    return (
        <div className="relative">
            <div
                className="absolute -top-3 left-1/2 -translate-x-1/2 z-10 px-4 py-1 rounded-full text-white text-xs font-semibold shadow-md flex items-center gap-1.5 whitespace-nowrap"
                style={{ backgroundColor: color }}
            >
                <span>{icon}</span>
                <span>{category.name}</span>
            </div>

            <div
                className="rounded-full bg-gradient-to-br from-purple-50 to-indigo-50 dark:from-purple-950/30 dark:to-indigo-950/30 pt-6 pb-4 px-6 shadow-lg border-2 transition-all"
                style={{ borderColor: `${color}55` }}
            >
                <div className="flex flex-wrap items-center justify-center gap-2 max-w-[42rem]">
                    {integrations.map(integration => {
                        const typeCode = integration.integration_type?.code || '';
                        const resourceName = INTERNAL_MODULE_RESOURCE_NAME[typeCode];
                        const isActive = resourceName ? resourceActive[resourceName] === true : false;
                        const moduleIcon = MODULE_ICONS[typeCode] || '⚙️';
                        const displayName = (integration.integration_type?.name || integration.name).replace(/\s*\(Modulo\)\s*$/i, '');

                        return (
                            <div
                                key={integration.id}
                                className="flex items-center gap-2 px-2 py-1.5 rounded-full bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 shadow-sm"
                            >
                                <div className="w-7 h-7 rounded-full bg-gradient-to-br from-purple-100 to-indigo-100 dark:from-purple-900/40 dark:to-indigo-900/40 ring-1 ring-purple-200 dark:ring-purple-800 flex items-center justify-center flex-shrink-0">
                                    <span className="text-sm">{moduleIcon}</span>
                                </div>
                                <span className="text-xs font-semibold text-gray-800 dark:text-gray-100 whitespace-nowrap">
                                    {displayName}
                                </span>
                                <span
                                    className={`w-2 h-2 rounded-full flex-shrink-0 ${
                                        isActive ? 'bg-green-500 shadow-[0_0_6px_rgba(34,197,94,0.6)]' : 'bg-gray-300 dark:bg-gray-600'
                                    }`}
                                    title={isActive ? 'Modulo activo' : 'Modulo inactivo'}
                                />
                            </div>
                        );
                    })}
                </div>
            </div>
        </div>
    );
}
