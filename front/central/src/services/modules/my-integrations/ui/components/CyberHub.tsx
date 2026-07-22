'use client';

import { forwardRef } from 'react';
import { RefreshCw } from 'lucide-react';
import type { Integration } from '@/services/integrations/core/domain/types';
import { INTERNAL_MODULE_RESOURCE_NAME } from '../../domain/types';

interface CyberHubProps {
    integrations: Integration[];
    resourceActive: Record<string, boolean>;
    onSyncClick?: () => void;
}

const ORBIT_RADIUS = 120;

export const CyberHub = forwardRef<HTMLDivElement, CyberHubProps>(function CyberHub(
    { integrations, resourceActive, onSyncClick },
    ref,
) {
    const visibleModules = integrations.filter(integration => {
        const typeCode = integration.integration_type?.code || '';
        const resourceName = INTERNAL_MODULE_RESOURCE_NAME[typeCode];
        return resourceName ? resourceActive[resourceName] === true : false;
    });

    return (
        <div className="relative z-10 flex justify-center">
            <div ref={ref} className="relative h-72 w-72">
                {visibleModules.length > 0 && (
                    <div className="absolute inset-0" style={{ animation: 'cyber-spin 45s linear infinite' }}>
                        <div className="absolute inset-6 rounded-full border border-dashed border-indigo-300/50 dark:border-indigo-500/30" />
                        {visibleModules.map((integration, i) => {
                            const theta = (360 / visibleModules.length) * i - 90;
                            const displayName = (integration.integration_type?.name || integration.name)
                                .replace(/\s*\(Modulo\)\s*$/i, '');
                            const isFunctional = integration.is_active === true;
                            return (
                                <span
                                    key={integration.id}
                                    className="absolute left-1/2 top-1/2 flex items-center gap-1.5 whitespace-nowrap text-[9px] font-bold uppercase tracking-[0.18em] text-gray-500 dark:text-gray-400"
                                    style={{
                                        transform: `rotate(${theta}deg) translateX(${ORBIT_RADIUS}px) rotate(-90deg) translate(-50%, -50%)`,
                                        transformOrigin: '0 0',
                                    }}
                                    title={isFunctional ? 'Modulo activo' : 'Modulo inactivo'}
                                >
                                    <span
                                        className={`h-1.5 w-1.5 rounded-full ${isFunctional ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'}`}
                                    />
                                    {displayName}
                                </span>
                            );
                        })}
                    </div>
                )}

                <div className="absolute inset-16">
                    <div
                        className="absolute inset-0 rounded-full"
                        style={{
                            background:
                                'conic-gradient(from 0deg, transparent 0%, #22d3ee 12%, transparent 28%, transparent 50%, #a855f7 62%, transparent 78%)',
                            WebkitMask: 'radial-gradient(farthest-side, transparent calc(100% - 3px), #000 calc(100% - 2px))',
                            mask: 'radial-gradient(farthest-side, transparent calc(100% - 3px), #000 calc(100% - 2px))',
                            animation: 'cyber-spin 5s linear infinite',
                        }}
                    />
                    <div
                        className="absolute inset-3 rounded-full border border-dashed border-gray-300 dark:border-gray-600"
                        style={{ animation: 'cyber-spin 20s linear infinite reverse' }}
                    />
                    <div className="absolute inset-6 flex flex-col items-center justify-center gap-1 rounded-full border border-gray-100 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-900">
                        <span className="text-[9px] uppercase tracking-[0.3em] text-gray-400">nucleo</span>
                        <span className="text-sm font-bold tracking-wide text-gray-800 dark:text-white">
                            Probability
                        </span>
                        <button
                            onClick={onSyncClick}
                            title="Sincronizacion global"
                            className="mt-1 flex h-8 w-8 items-center justify-center rounded-full border border-cyan-400/60 bg-cyan-50 text-cyan-600 transition-all hover:scale-110 hover:bg-cyan-100 hover:shadow-[0_0_12px_rgba(34,211,238,0.6)] dark:bg-cyan-900/30 dark:text-cyan-300 dark:hover:bg-cyan-900/50"
                        >
                            <RefreshCw size={14} />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
});
