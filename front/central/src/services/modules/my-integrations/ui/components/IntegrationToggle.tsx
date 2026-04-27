'use client';

import type { Integration } from '@/services/integrations/core/domain/types';
import { CATEGORY_ICONS } from '../../domain/types';
import { PencilSquareIcon } from '@heroicons/react/24/outline';

interface IntegrationToggleProps {
    integration: Integration;
    onToggle: (integration: Integration) => void;
    onEdit?: (integration: Integration) => void;
    togglingId: number | null;
}

export function IntegrationToggle({ integration, onToggle, onEdit, togglingId }: IntegrationToggleProps) {
    const icon = CATEGORY_ICONS[integration.category] || '🔗';
    const isToggling = togglingId === integration.id;
    const isEnvioClick = integration.integration_type?.code === 'envioclick';
    const typeName = isEnvioClick ? 'Transportadora' : integration.integration_type?.name;
    const hideLogo = isEnvioClick;
    const hideInstanceName = isEnvioClick;

    return (
        <div className="flex items-center gap-2 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            {hideLogo ? (
                <div className="w-7 h-7 rounded bg-gray-200 dark:bg-gray-600 flex items-center justify-center flex-shrink-0">
                    <span className="text-xs">{icon}</span>
                </div>
            ) : integration.integration_type?.image_url ? (
                <img
                    src={integration.integration_type.image_url}
                    alt={typeName || integration.name}
                    className="w-7 h-7 rounded object-contain flex-shrink-0"
                />
            ) : (
                <div className="w-7 h-7 rounded bg-gray-200 dark:bg-gray-600 flex items-center justify-center flex-shrink-0">
                    <span className="text-xs">{icon}</span>
                </div>
            )}

            <div className="flex flex-col min-w-0 flex-1">
                {typeName && (
                    <span className="text-xs font-semibold text-gray-800 dark:text-gray-100 dark:text-gray-200 truncate leading-tight">
                        {typeName}
                    </span>
                )}
                {!hideInstanceName && (
                    <span className="text-[11px] text-gray-500 dark:text-gray-400 truncate leading-tight">
                        {integration.name}
                    </span>
                )}
            </div>

            {/* Editar */}
            <button
                onClick={() => onEdit?.(integration)}
                className="p-1 text-gray-400 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors flex-shrink-0"
                title="Editar integración"
            >
                <PencilSquareIcon className="w-4 h-4" />
            </button>

            {/* Toggle */}
            <button
                onClick={() => onToggle(integration)}
                disabled={isToggling}
                className={`relative inline-flex h-5 w-9 items-center rounded-full flex-shrink-0 transition-colors ${
                    integration.is_active
                        ? 'bg-green-500'
                        : 'bg-gray-300 dark:bg-gray-600'
                } ${isToggling ? 'opacity-50 cursor-wait' : 'cursor-pointer'}`}
            >
                <span
                    className={`inline-block h-3.5 w-3.5 rounded-full bg-white dark:bg-gray-800 transition-transform shadow-sm ${
                        integration.is_active ? 'translate-x-4' : 'translate-x-0.5'
                    }`}
                />
            </button>
        </div>
    );
}
