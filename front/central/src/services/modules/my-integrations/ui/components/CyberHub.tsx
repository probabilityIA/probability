'use client';

import { forwardRef, type CSSProperties } from 'react';
import { Package, Truck, Bell, Users, Store, ReceiptText, Cog, type LucideIcon } from 'lucide-react';
import type { Integration } from '@/services/integrations/core/domain/types';
import type { IntegrationStatsItem } from '@/services/integrations/core/infra/actions/stats';
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
    onCompareClick?: () => void;
    ownStats?: IntegrationStatsItem | null;
    lastOrderAt?: string;
}

const HUB_SIZE = 560;
const CORE_INSET = 173;
const ORBIT_RADIUS = 238;
const PILL_RADIUS = 188;
const CARD_OFFSET = 132;
const CHARGE_INSET = 163;
const SPARK_RADIUS = 128;

const numberFormat = new Intl.NumberFormat('es-CO');

const OWN_LABELS = [
    { key: 'orders_in_progress', label: 'en curso', color: '#3b82f6', angle: -131 },
    { key: 'orders_delivered', label: 'entregadas', color: '#22c55e', angle: -49 },
    { key: 'orders_cancelled', label: 'canceladas', color: '#ef4444', angle: 131 },
    { key: 'orders_returned', label: 'devueltas', color: '#f59e0b', angle: 49 },
] as const;

const CARD_CLASS =
    'absolute flex w-[90px] flex-col rounded-2xl bg-white px-3 py-2 shadow-[0_8px_24px_-12px_rgba(15,23,42,0.28)] ring-1 ring-gray-100 dark:bg-gray-800 dark:ring-gray-700';

function pillPosition(angle: number): CSSProperties {
    const rad = (angle * Math.PI) / 180;
    return {
        left: `calc(50% + ${Math.round(PILL_RADIUS * Math.cos(rad))}px)`,
        top: `calc(50% + ${Math.round(PILL_RADIUS * Math.sin(rad))}px)`,
        transform: 'translate(-50%, -50%)',
    };
}

function relativeTime(iso?: string): string | null {
    if (!iso) return null;
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) return null;
    const minutes = Math.floor((Date.now() - date.getTime()) / 60000);
    if (minutes < 1) return 'hace un momento';
    if (minutes < 60) return `hace ${minutes} min`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `hace ${hours} h`;
    const days = Math.floor(hours / 24);
    if (days < 30) return `hace ${days} d`;
    const months = Math.floor(days / 30);
    if (months < 12) return `hace ${months} mes${months > 1 ? 'es' : ''}`;
    return `hace ${Math.floor(months / 12)} a`;
}

export const CyberHub = forwardRef<HTMLDivElement, CyberHubProps>(function CyberHub(
    { integrations, resourceActive, onSyncClick, findingsCount = 0, onFindingsClick, onCompareClick, ownStats, lastOrderAt },
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
    const comparaInventario = environment === 'inventory' && Boolean(onCompareClick);
    const chargeColor = (mode === 'products' || (mode === 'idle' && environment === 'products')) ? '#8b5cf6' : '#22d3ee';
    const statusColor = mode === 'inventory'
        ? 'text-blue-600 dark:text-blue-400'
        : mode === 'products'
            ? 'text-violet-600 dark:text-violet-400'
            : 'text-emerald-600 dark:text-emerald-400';

    const ultimaOrden = relativeTime(lastOrderAt);

    const visibleModules = integrations.filter(integration => {
        const typeCode = integration.integration_type?.code || '';
        const resourceName = INTERNAL_MODULE_RESOURCE_NAME[typeCode];
        return resourceName ? resourceActive[resourceName] === true : false;
    });

    return (
        <div className="relative z-10 flex justify-center">
            <div ref={ref} className="relative" style={{ width: HUB_SIZE, height: HUB_SIZE }}>
                <div
                    className="pointer-events-none absolute inset-0 rounded-full"
                    style={{
                        background:
                            'radial-gradient(circle, rgba(148,163,184,0.13) 0%, rgba(148,163,184,0.05) 52%, rgba(148,163,184,0) 72%)',
                    }}
                />
                <div className="pointer-events-none absolute inset-0 rounded-full border border-dashed border-gray-200 dark:border-gray-700/70" />

                {visibleModules.length > 0 && (
                    <div className="orbit-ring absolute inset-0" style={{ animation: 'cyber-spin 60s linear infinite' }}>
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
                                        className="orbit-chip group relative flex h-9 w-9 cursor-default items-center justify-center rounded-full bg-white shadow-[0_6px_18px_-10px_rgba(15,23,42,0.5)] ring-1 ring-gray-100 transition-shadow hover:shadow-md dark:bg-gray-800 dark:ring-gray-700"
                                        style={{ animation: 'cyber-spin 60s linear infinite reverse' }}
                                    >
                                        <Icon size={15} className={isFunctional ? 'text-indigo-400 dark:text-indigo-300' : 'text-gray-300 dark:text-gray-500'} />
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
                    <div className="pointer-events-none absolute" style={{ inset: CHARGE_INSET }}>
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
                                    transform: `rotate(${deg}deg) translateY(-52%) translateX(${SPARK_RADIUS}px)`,
                                    transformOrigin: '0 0',
                                    animation: `cyber-spark ${(0.7 + i * 0.11).toFixed(2)}s ease-out infinite`,
                                    animationDelay: `-${(i * 0.17).toFixed(2)}s`,
                                }}
                            />
                        ))}
                    </div>
                )}

                <div className="absolute" style={{ inset: CORE_INSET }}>
                    <div
                        className="absolute inset-0 rounded-full"
                        style={{
                            background:
                                'conic-gradient(from 0deg, transparent 0%, #22d3ee 12%, transparent 28%, transparent 50%, #a855f7 62%, transparent 78%)',
                            WebkitMask: 'radial-gradient(farthest-side, transparent calc(100% - 3px), #000 calc(100% - 2px))',
                            mask: 'radial-gradient(farthest-side, transparent calc(100% - 3px), #000 calc(100% - 2px))',
                            animation: `cyber-spin ${busy ? '1.1s' : '8s'} linear infinite`,
                        }}
                    />
                    <div
                        className="absolute inset-[7px] flex flex-col items-center justify-center gap-2 rounded-full bg-white shadow-[0_24px_60px_-24px_rgba(15,23,42,0.45)] ring-1 ring-gray-50 dark:bg-gray-900 dark:ring-gray-700"
                        style={busy
                            ? ({ '--charge': chargeColor, animation: 'cyber-charge 1s ease-in-out infinite' } as CSSProperties)
                            : undefined}
                    >
                        {onFindingsClick && (
                            <button
                                onClick={onFindingsClick}
                                title={
                                    findingsCount > 0
                                        ? `Resumen general y ${findingsCount} ${findingsCount === 1 ? 'hallazgo' : 'hallazgos'} en tus integraciones`
                                        : 'Resumen general de tus integraciones'
                                }
                                aria-label="Ver el resumen general de tus integraciones"
                                className={`absolute -top-3.5 left-1/2 flex -translate-x-1/2 items-center gap-1.5 whitespace-nowrap rounded-full border-2 border-white px-3 py-1 text-[11px] font-bold shadow-lg transition-transform hover:scale-105 dark:border-gray-900 ${
                                    findingsCount > 0 ? 'bg-amber-500 text-white' : 'bg-gray-700 text-white dark:bg-gray-600'
                                }`}
                                style={findingsCount > 0 ? { animation: 'cyber-alert-pulse 2s ease-in-out infinite' } : undefined}
                            >
                                Resumen general
                                {findingsCount > 0 && (
                                    <span className="rounded-full bg-white/25 px-1.5 text-[10px] font-black">{findingsCount}</span>
                                )}
                            </button>
                        )}
                        <span className="text-[10px] uppercase tracking-[0.4em] text-gray-400 dark:text-gray-500">nucleo</span>
                        <span className="text-[23px] font-bold leading-none tracking-tight text-gray-800 dark:text-white">
                            Probability
                        </span>
                        <button
                            onClick={comparaInventario ? onCompareClick : (onSyncClick ?? runCurrent)}
                            disabled={running || (!comparaInventario && !onSyncClick && !canRun)}
                            title={
                                environment === 'products'
                                    ? 'Iniciar comparacion de productos'
                                    : comparaInventario
                                        ? 'Comparar el stock de cada canal contra Probability'
                                        : 'Elige arriba que quieres sincronizar'
                            }
                            className="mt-1 flex h-9 items-center justify-center rounded-full border border-cyan-200 bg-cyan-50/70 px-6 text-[12.5px] font-bold uppercase tracking-[0.12em] text-cyan-600 transition-all hover:scale-105 hover:bg-cyan-100 hover:shadow-[0_0_16px_rgba(34,211,238,0.55)] disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:scale-100 disabled:hover:shadow-none dark:border-cyan-500/40 dark:bg-cyan-900/30 dark:text-cyan-300 dark:hover:bg-cyan-900/50"
                        >
                            <span className={busy ? 'animate-pulse' : ''}>{comparaInventario ? 'Comparar' : 'Iniciar'}</span>
                        </button>
                    </div>
                </div>

                {ownStats && !busy && (
                    <div className="pointer-events-none absolute inset-0">
                        <div
                            className={`${CARD_CLASS} items-start text-left`}
                            style={{ right: `calc(50% + ${CARD_OFFSET}px)`, top: '50%', transform: 'translateY(-50%)' }}
                        >
                            <span className="text-[9px] font-semibold uppercase tracking-[0.18em] text-gray-400 dark:text-gray-500">Ordenes</span>
                            <span className="mt-1 text-[26px] font-bold leading-none tabular-nums text-gray-900 dark:text-white">
                                {numberFormat.format(ownStats.orders_count)}
                            </span>
                            {ultimaOrden && (
                                <span className="mt-1 whitespace-nowrap text-[9px] leading-none text-gray-400 dark:text-gray-500">
                                    Ult. {ultimaOrden}
                                </span>
                            )}
                        </div>

                        <div
                            className={`${CARD_CLASS} items-start text-left`}
                            style={{ left: `calc(50% + ${CARD_OFFSET}px)`, top: '50%', transform: 'translateY(-50%)' }}
                        >
                            <span className="text-[9px] font-semibold uppercase tracking-[0.18em] text-gray-400 dark:text-gray-500">Productos</span>
                            <span className="mt-1 text-[26px] font-bold leading-none tabular-nums text-gray-900 dark:text-white">
                                {numberFormat.format(ownStats.products_count)}
                            </span>
                        </div>

                        {OWN_LABELS.map(item => (
                            <span
                                key={item.key}
                                className="absolute flex items-center gap-2 whitespace-nowrap rounded-full bg-white px-3 py-1.5 text-[11.5px] shadow-[0_8px_24px_-14px_rgba(15,23,42,0.4)] ring-1 ring-gray-100 dark:bg-gray-800 dark:ring-gray-700"
                                style={pillPosition(item.angle)}
                            >
                                <span className="h-1.5 w-1.5 flex-shrink-0 rounded-full" style={{ backgroundColor: item.color }} />
                                <span className="font-bold tabular-nums text-gray-800 dark:text-gray-100">
                                    {numberFormat.format(ownStats[item.key])}
                                </span>
                                <span className="text-gray-400 dark:text-gray-500">{item.label}</span>
                            </span>
                        ))}
                    </div>
                )}
            </div>
            <div className="absolute -bottom-3 left-1/2 flex -translate-x-1/2 flex-col items-center gap-2">
                {statusText && (
                    <span className={`whitespace-nowrap text-[11.5px] font-semibold ${statusColor}`}>{statusText}</span>
                )}
            </div>
        </div>
    );
});
