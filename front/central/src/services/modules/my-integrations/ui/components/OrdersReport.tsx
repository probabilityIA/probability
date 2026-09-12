'use client';

import type { Integration } from '@/services/integrations/core/domain/types';
import type { IntegrationStatsItem } from '@/services/integrations/core/infra/actions/stats';
import { ChannelLogo } from './ChannelLogo';
import { CARD_BORDER } from '../panel-theme';
import { PanelSummary } from './PanelToolbar';

interface OrdersReportProps {
    integrations: Integration[];
    stats: Record<number, IntegrationStatsItem>;
    statsLoaded: boolean;
}

const numberFormat = new Intl.NumberFormat('es-CO');

const BUCKETS = [
    { key: 'orders_in_progress', label: 'En curso', color: '#3b82f6' },
    { key: 'orders_delivered', label: 'Entregadas', color: '#22c55e' },
    { key: 'orders_cancelled', label: 'Canceladas', color: '#ef4444' },
    { key: 'orders_returned', label: 'Devueltas', color: '#f59e0b' },
] as const;

const VACIO: IntegrationStatsItem = {
    integration_id: 0,
    orders_count: 0,
    orders_in_progress: 0,
    orders_delivered: 0,
    orders_cancelled: 0,
    orders_returned: 0,
    products_count: 0,
};

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

function Barra({ item, total }: { item: IntegrationStatsItem; total: number }) {
    if (total <= 0) {
        return <div className="h-2 w-full rounded-full bg-gray-100 dark:bg-gray-700" />;
    }
    return (
        <div className="flex h-2 w-full overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700">
            {BUCKETS.map(bucket => {
                const value = item[bucket.key];
                if (value <= 0) return null;
                return (
                    <span
                        key={bucket.key}
                        title={`${bucket.label}: ${numberFormat.format(value)}`}
                        style={{ width: `${(value / total) * 100}%`, backgroundColor: bucket.color }}
                    />
                );
            })}
        </div>
    );
}

export function OrdersReport({ integrations, stats, statsLoaded }: OrdersReportProps) {
    const filas = integrations
        .map(integration => ({
            integration,
            item: stats[integration.id] ?? { ...VACIO, integration_id: integration.id },
        }))
        .sort((a, b) => b.item.orders_count - a.item.orders_count);

    const totales = filas.reduce<IntegrationStatsItem>((acc, fila) => ({
        ...acc,
        orders_count: acc.orders_count + fila.item.orders_count,
        orders_in_progress: acc.orders_in_progress + fila.item.orders_in_progress,
        orders_delivered: acc.orders_delivered + fila.item.orders_delivered,
        orders_cancelled: acc.orders_cancelled + fila.item.orders_cancelled,
        orders_returned: acc.orders_returned + fila.item.orders_returned,
        products_count: acc.products_count + fila.item.products_count,
    }), { ...VACIO });

    const mayor = filas.reduce((max, fila) => Math.max(max, fila.item.orders_count), 0);

    if (!statsLoaded) {
        return (
            <p className="py-20 text-center text-[12px] text-gray-500 dark:text-gray-400">
                Cargando el resumen de ordenes
            </p>
        );
    }

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto pb-2">
            <PanelSummary>
                <div className="order-3 ml-auto flex flex-wrap items-center justify-end gap-x-6 gap-y-2">
                    <div className="leading-tight">
                        <p className="text-[9px] font-semibold uppercase tracking-[0.14em] text-gray-400 dark:text-gray-500">
                            Ordenes
                        </p>
                        <p className="flex items-baseline gap-1.5">
                            <span className="text-[17px] font-bold leading-none tabular-nums text-gray-900 dark:text-white">
                                {numberFormat.format(totales.orders_count)}
                            </span>
                            <span className="text-[10px] text-gray-400 dark:text-gray-500">
                                en {filas.length} {filas.length === 1 ? 'origen' : 'origenes'}
                            </span>
                        </p>
                    </div>

                    <span className="h-6 w-px bg-gray-200 dark:bg-gray-700" />

                    {BUCKETS.map(bucket => {
                        const value = totales[bucket.key];
                        const pct = totales.orders_count > 0 ? Math.round((value / totales.orders_count) * 100) : 0;
                        return (
                            <div key={bucket.key} className="flex items-center gap-2 leading-tight">
                                <span className="h-1.5 w-1.5 flex-shrink-0 rounded-full" style={{ backgroundColor: bucket.color }} />
                                <div>
                                    <p className="text-[9px] font-semibold uppercase tracking-[0.14em] text-gray-400 dark:text-gray-500">
                                        {bucket.label}
                                    </p>
                                    <p className="flex items-baseline gap-1.5">
                                        <span className="text-[17px] font-bold leading-none tabular-nums text-gray-900 dark:text-white">
                                            {numberFormat.format(value)}
                                        </span>
                                        <span className="text-[10px] text-gray-400 dark:text-gray-500">{pct}%</span>
                                    </p>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </PanelSummary>

            <div
                className="overflow-hidden rounded-2xl border dark:bg-gray-800/60"
                style={{ borderColor: CARD_BORDER, backgroundColor: '#ffffff' }}
            >
                <div
                    className="grid grid-cols-[minmax(10rem,1.4fr)_5rem_minmax(10rem,2fr)_6rem] items-center gap-3 px-4 py-2 text-[10px] font-bold uppercase tracking-wider text-gray-400 dark:text-gray-500"
                    style={{ backgroundColor: '#fafafd' }}
                >
                    <span>Origen</span>
                    <span className="text-right">Ordenes</span>
                    <span>Como van</span>
                    <span className="text-right">Productos</span>
                </div>
                {filas.map(({ integration, item }) => {
                    const ultima = relativeTime(item.last_order_at);
                    return (
                        <div
                            key={integration.id}
                            className="grid grid-cols-[minmax(10rem,1.4fr)_5rem_minmax(10rem,2fr)_6rem] items-center gap-3 border-t px-4 py-2.5"
                            style={{ borderColor: CARD_BORDER }}
                        >
                            <div className="min-w-0">
                                <span className="flex items-center gap-2 text-[12px] font-bold text-gray-800 dark:text-gray-100">
                                    <ChannelLogo
                                        url={integration.integration_type?.image_url}
                                        code={integration.integration_type?.code || integration.category}
                                    />
                                    <span className="truncate">
                                        {integration.integration_type?.name || integration.name}
                                    </span>
                                </span>
                                {ultima && (
                                    <span className="block truncate text-[10px] text-gray-400 dark:text-gray-500">
                                        Ult. orden {ultima}
                                    </span>
                                )}
                            </div>
                            <span className="text-right text-[17px] font-bold leading-none tabular-nums text-gray-900 dark:text-white">
                                {numberFormat.format(item.orders_count)}
                            </span>
                            <div>
                                <Barra item={item} total={mayor} />
                                {item.orders_count === 0 && (
                                    <span className="mt-1 block text-[10px] italic text-gray-300 dark:text-gray-600">
                                        Sin ordenes registradas
                                    </span>
                                )}
                            </div>
                            <span className="text-right text-[13px] font-semibold tabular-nums text-gray-600 dark:text-gray-300">
                                {numberFormat.format(item.products_count)}
                            </span>
                        </div>
                    );
                })}
            </div>

            <p className="text-[10.5px] text-gray-400 dark:text-gray-500">
                Cada barra usa la misma escala: la mas larga es el origen con mas ordenes. Los colores son los mismos
                del diagrama.
            </p>
        </div>
    );
}
