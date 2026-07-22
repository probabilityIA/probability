'use client';

import type { CSSProperties } from 'react';
import type { Integration } from '@/services/integrations/core/domain/types';
import { SlidersHorizontal } from 'lucide-react';

interface CyberIntegrationNodeProps {
    integration: Integration;
    color: string;
    onToggle: (integration: Integration) => void;
    onEdit: (integration: Integration) => void;
    togglingId: number | null;
    editingId: number | null;
}

export function CyberIntegrationNode({ integration, color, onToggle, onEdit, togglingId, editingId }: CyberIntegrationNodeProps) {
    const isToggling = togglingId === integration.id;
    const isEditing = editingId === integration.id;
    const active = integration.is_active;
    const typeName = integration.integration_type?.name || integration.name;

    return (
        <div
            className="group relative overflow-hidden rounded-xl p-px shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:shadow-[0_6px_18px_-8px_var(--neon)]"
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

            <div className="relative z-10 flex items-center gap-2 rounded-[11px] bg-white p-2 transition-colors group-hover:bg-gray-50/80 dark:bg-gray-800 dark:group-hover:bg-gray-700/60">
                {integration.integration_type?.image_url ? (
                    <img
                        src={integration.integration_type.image_url}
                        alt={typeName}
                        className="h-9 w-9 flex-shrink-0 rounded-full object-contain ring-1 ring-gray-200 transition-transform duration-200 group-hover:scale-105 dark:ring-gray-600"
                    />
                ) : (
                    <div
                        className="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-full text-sm font-bold ring-1 ring-gray-200 transition-transform duration-200 group-hover:scale-105 dark:ring-gray-600"
                        style={{ color, backgroundColor: `${color}14` }}
                    >
                        {typeName.charAt(0).toUpperCase()}
                    </div>
                )}

                <div className="flex min-w-0 flex-1 flex-col">
                    <span className="truncate text-xs font-semibold leading-tight text-gray-800 dark:text-gray-100">
                        {typeName}
                    </span>
                    <span className="truncate text-[11px] leading-tight text-gray-500 dark:text-gray-400">
                        {integration.name}
                    </span>
                </div>

                <button
                    onClick={() => onEdit(integration)}
                    disabled={isEditing}
                    title="Configurar integracion"
                    className={`flex-shrink-0 rounded-lg border border-indigo-200 bg-indigo-50 p-1.5 text-indigo-600 shadow-sm transition-all hover:scale-105 hover:bg-indigo-100 hover:shadow dark:border-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300 dark:hover:bg-indigo-900/60 ${isEditing ? 'cursor-wait opacity-60' : ''}`}
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
                    className={`relative inline-flex h-5 w-9 flex-shrink-0 items-center rounded-full transition-colors ${
                        active ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'
                    } ${isToggling ? 'cursor-wait opacity-50' : 'cursor-pointer'}`}
                >
                    <span
                        className={`inline-block h-3.5 w-3.5 rounded-full bg-white shadow-sm transition-transform dark:bg-gray-800 ${
                            active ? 'translate-x-4' : 'translate-x-0.5'
                        }`}
                    />
                </button>
            </div>
        </div>
    );
}
