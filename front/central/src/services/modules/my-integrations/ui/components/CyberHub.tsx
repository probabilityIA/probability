'use client';

import { forwardRef, type CSSProperties } from 'react';
import { Package, Truck, Bell, Users, Store, ReceiptText, Cog, type LucideIcon } from 'lucide-react';
import type { Integration } from '@/services/integrations/core/domain/types';
import { INTERNAL_MODULE_RESOURCE_NAME } from '../../domain/types';
import { useSyncActivity } from '../sync-activity-context';

const MODULE_ICONS: Record<string, LucideIcon> = {
    inventory: Package,
    delivery: Truck,
    notifications: Bell,
    customers: Users,
    storefront_module: Store,
    invoicing_module: ReceiptText,
};

interface CyberHubProps {
    integrations: Integration[];
    resourceActive: Record<string, boolean>;
    onSyncClick?: () => void;
    findingsCount?: number;
    onFindingsClick?: () => void;
}

const ORBIT_RADIUS = 120;

export const CyberHub = forwardRef<HTMLDivElement, CyberHubProps>(function CyberHub(
    { integrations, resourceActive, onSyncClick, findingsCount = 0, onFindingsClick },
    ref,
) {
    const { mode, nodes, running, environment, canRun, runCurrent } = useSyncActivity();
    const states = Object.values(nodes);
    const busy = mode !== 'idle' || states.some(s => s === 'active' || s === 'scan');
    const finished = !busy && states.length > 0 && states.every(s => s === 'done' || s === 'error');
    const statusText = mode === 'inventory'
        ? 'Enviando stock a los canales...'
        : mode === 'products'
            ? 'Recibiendo catalogos de los canales...'
            : finished
                ? 'Sincronizacion completada'
                : '';
    const chargeColor = (mode === 'products' || (mode === 'idle' && environment === 'products')) ? '#8b5cf6' : '#22d3ee';
    const statusColor = mode === 'inventory'
        ? 'text-blue-600 dark:text-blue-400'
        : mode === 'products'
            ? 'text-violet-600 dark:text-violet-400'
            : 'text-emerald-600 dark:text-emerald-400';

    const visibleModules = integrations.filter(integration => {
        const typeCode = integration.integration_type?.code || '';
        const resourceName = INTERNAL_MODULE_RESOURCE_NAME[typeCode];
        return resourceName ? resourceActive[resourceName] === true : false;
    });

    return (
        <div className="relative z-10 flex justify-center">
            <div ref={ref} className="relative h-72 w-72">
                {visibleModules.length > 0 && (
                    <div className="orbit-ring absolute inset-0" style={{ animation: 'cyber-spin 45s linear infinite' }}>
                        <div className="absolute inset-6 rounded-full border border-dashed border-indigo-300/50 dark:border-indigo-500/30" />
                        {visibleModules.map((integration, i) => {
                            const theta = (360 / visibleModules.length) * i - 90;
                            const typeCode = integration.integration_type?.code || '';
                            const displayName = (integration.integration_type?.name || integration.name)
                                .replace(/\s*\(Modulo\)\s*$/i, '');
                            const isFunctional = integration.is_active === true;
                            const Icon = MODULE_ICONS[typeCode] || Cog;
                            return (
                                <div
                                    key={integration.id}
                                    className="absolute left-1/2 top-1/2"
                                    style={{
                                        transform: `rotate(${theta}deg) translateX(${ORBIT_RADIUS}px) rotate(${-theta}deg) translate(-50%, -50%)`,
                                        transformOrigin: '0 0',
                                    }}
                                >
                                    <div
                                        className="orbit-chip group relative flex h-9 w-9 cursor-default items-center justify-center rounded-full border border-gray-200 bg-white shadow-sm transition-shadow hover:shadow-md dark:border-gray-700 dark:bg-gray-800"
                                        style={{ animation: 'cyber-spin 45s linear infinite reverse' }}
                                    >
                                        <Icon size={15} className={isFunctional ? 'text-indigo-500 dark:text-indigo-400' : 'text-gray-400 dark:text-gray-500'} />
                                        <span
                                            className={`absolute -right-0.5 -top-0.5 h-2 w-2 rounded-full border border-white dark:border-gray-800 ${isFunctional ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'}`}
                                        />
                                        <span className="pointer-events-none absolute left-1/2 top-full z-20 mt-1.5 -translate-x-1/2 whitespace-nowrap rounded-md bg-gray-900 px-2 py-0.5 text-[10px] font-semibold text-white opacity-0 shadow-lg transition-opacity duration-150 group-hover:opacity-100 dark:bg-gray-700">
                                            {displayName}
                                            <span className={`ml-1 ${isFunctional ? 'text-green-400' : 'text-gray-400'}`}>
                                                {isFunctional ? 'activo' : 'inactivo'}
                                            </span>
                                        </span>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}

                {busy && (
                    <div className="pointer-events-none absolute inset-10">
                        <div
                            className="absolute inset-0 rounded-full"
                            style={{
                                background: `conic-gradient(from 0deg, transparent 0deg, ${chargeColor} 18deg, transparent 40deg, transparent 160deg, ${chargeColor} 182deg, transparent 205deg)`,
                                WebkitMask: 'radial-gradient(farthest-side, transparent calc(100% - 5px), #000 calc(100% - 4px))',
                                mask: 'radial-gradient(farthest-side, transparent calc(100% - 5px), #000 calc(100% - 4px))',
                                filter: 'blur(1px)',
                                animation: 'cyber-arc .7s linear infinite, cyber-flicker .45s steps(3) infinite',
                            }}
                        />
                        <div
                            className="absolute inset-2 rounded-full"
                            style={{
                                background: `conic-gradient(from 180deg, transparent 0deg, ${chargeColor} 12deg, transparent 30deg)`,
                                WebkitMask: 'radial-gradient(farthest-side, transparent calc(100% - 3px), #000 calc(100% - 2px))',
                                mask: 'radial-gradient(farthest-side, transparent calc(100% - 3px), #000 calc(100% - 2px))',
                                animation: 'cyber-arc .45s linear infinite reverse, cyber-flicker .3s steps(2) infinite',
                            }}
                        />
                        {[0, 60, 120, 180, 240, 300].map((deg, i) => (
                            <span
                                key={deg}
                                className="absolute left-1/2 top-1/2 h-1.5 w-1.5 rounded-full"
                                style={{
                                    background: chargeColor,
                                    boxShadow: `0 0 10px 3px ${chargeColor}`,
                                    transform: `rotate(${deg}deg) translateY(-52%) translateX(74px)`,
                                    transformOrigin: '0 0',
                                    animation: `cyber-spark ${(0.7 + i * 0.11).toFixed(2)}s ease-out infinite`,
                                    animationDelay: `-${(i * 0.17).toFixed(2)}s`,
                                }}
                            />
                        ))}
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
                            animation: `cyber-spin ${busy ? '1.1s' : '5s'} linear infinite`,
                        }}
                    />
                    <div
                        className="absolute inset-3 rounded-full border border-dashed border-gray-300 dark:border-gray-600"
                        style={{ animation: 'cyber-spin 20s linear infinite reverse' }}
                    />
                    <div
                        className="absolute inset-6 flex flex-col items-center justify-center gap-1 rounded-full border border-gray-100 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-900"
                        style={busy
                            ? ({ '--charge': chargeColor, animation: 'cyber-charge 1s ease-in-out infinite' } as CSSProperties)
                            : undefined}
                    >
                        {findingsCount > 0 && onFindingsClick && (
                            <button
                                onClick={onFindingsClick}
                                title={`${findingsCount} ${findingsCount === 1 ? 'hallazgo' : 'hallazgos'} en tus integraciones`}
                                aria-label={`Ver ${findingsCount} hallazgos de tus integraciones`}
                                className="absolute -right-1 -top-1 flex h-6 min-w-6 items-center justify-center rounded-full border-2 border-white bg-amber-500 px-1 text-[11px] font-black text-white shadow-lg transition-transform hover:scale-110 dark:border-gray-900"
                                style={{ animation: 'cyber-alert-pulse 2s ease-in-out infinite' }}
                            >
                                !
                            </button>
                        )}
                        <span className="text-[9px] uppercase tracking-[0.3em] text-gray-400">nucleo</span>
                        <span className="text-sm font-bold tracking-wide text-gray-800 dark:text-white">
                            Probability
                        </span>
                        <button
                            onClick={onSyncClick ?? runCurrent}
                            disabled={running || (!onSyncClick && !canRun)}
                            title={
                                environment === 'products'
                                    ? 'Iniciar comparacion de productos'
                                    : environment === 'inventory'
                                        ? 'Iniciar sincronizacion de inventario'
                                        : 'Elige arriba que quieres sincronizar'
                            }
                            className="mt-1 flex h-7 items-center justify-center rounded-full border border-cyan-400/60 bg-cyan-50 px-3.5 text-[11px] font-bold uppercase tracking-wider text-cyan-600 transition-all hover:scale-110 hover:bg-cyan-100 hover:shadow-[0_0_12px_rgba(34,211,238,0.6)] disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:scale-100 disabled:hover:shadow-none dark:bg-cyan-900/30 dark:text-cyan-300 dark:hover:bg-cyan-900/50"
                        >
                            <span className={busy ? 'animate-pulse' : ''}>Iniciar</span>
                        </button>
                    </div>
                </div>
            </div>
            <div className="absolute -bottom-3 left-1/2 flex -translate-x-1/2 flex-col items-center gap-2">
                {statusText && (
                    <span className={`whitespace-nowrap text-[11.5px] font-semibold ${statusColor}`}>{statusText}</span>
                )}
            </div>
        </div>
    );
});
