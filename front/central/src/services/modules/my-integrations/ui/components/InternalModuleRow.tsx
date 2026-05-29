'use client';

import type { Integration } from '@/services/integrations/core/domain/types';

interface InternalModuleRowProps {
    integration: Integration;
    isActive: boolean;
}

const MODULE_ICONS: Record<string, string> = {
    inventory: '📦',
    delivery: '🚚',
    notifications: '🔔',
    customers: '👥',
    storefront_module: '🛍️',
    invoicing_module: '🧾',
};

export function InternalModuleRow({ integration, isActive }: InternalModuleRowProps) {
    const rawTypeName = integration.integration_type?.name || integration.name;
    const displayName = rawTypeName.replace(/\s*\(Modulo\)\s*$/i, '');
    const typeCode = integration.integration_type?.code || '';
    const icon = MODULE_ICONS[typeCode] || '⚙️';

    return (
        <div className="flex items-center gap-3 p-2 pr-3 rounded-full bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 shadow-sm">
            <div className="w-10 h-10 rounded-full bg-gradient-to-br from-purple-100 to-indigo-100 dark:from-purple-900/40 dark:to-indigo-900/40 ring-1 ring-purple-200 dark:ring-purple-800 flex items-center justify-center flex-shrink-0">
                <span className="text-lg">{icon}</span>
            </div>

            <span className="text-sm font-semibold text-gray-800 dark:text-gray-100 truncate flex-1">
                {displayName}
            </span>

            <div
                className={`relative inline-flex h-5 w-9 items-center rounded-full flex-shrink-0 ${
                    isActive ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'
                }`}
                title={isActive ? 'Modulo activo' : 'Modulo inactivo'}
                aria-label={isActive ? 'Modulo activo' : 'Modulo inactivo'}
            >
                <span
                    className={`inline-block h-3.5 w-3.5 rounded-full bg-white dark:bg-gray-800 shadow-sm ${
                        isActive ? 'translate-x-4' : 'translate-x-0.5'
                    }`}
                />
            </div>
        </div>
    );
}
