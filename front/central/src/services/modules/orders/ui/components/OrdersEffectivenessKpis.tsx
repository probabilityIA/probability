'use client';

import { useEffect, useState } from 'react';
import { getIntegrationStatsAction } from '@/services/integrations/core/infra/actions/stats';

interface Props {
    businessId: number | null;
}

interface Totals {
    total: number;
    inProgress: number;
    delivered: number;
    returned: number;
}

const numberFormat = new Intl.NumberFormat('es-CO');

function percent(value: number, base: number): string {
    if (base <= 0) return '0%';
    return `${((value / base) * 100).toFixed(1)}%`;
}

export function OrdersEffectivenessKpis({ businessId }: Props) {
    const [totals, setTotals] = useState<Totals | null>(null);

    useEffect(() => {
        let active = true;

        if (!businessId) {
            setTotals(null);
            return;
        }

        const load = async () => {
            const result = await getIntegrationStatsAction(businessId);
            if (!active || !result.success) return;

            const sum = (result.data || []).reduce<Totals>(
                (acc, item) => ({
                    total: acc.total + (item.orders_count || 0),
                    inProgress: acc.inProgress + (item.orders_in_progress || 0),
                    delivered: acc.delivered + (item.orders_delivered || 0),
                    returned: acc.returned + (item.orders_returned || 0),
                }),
                { total: 0, inProgress: 0, delivered: 0, returned: 0 },
            );
            setTotals(sum);
        };

        load();
        return () => {
            active = false;
        };
    }, [businessId]);

    if (!totals || totals.total === 0) return null;

    const closed = totals.delivered + totals.returned;
    if (closed === 0) return null;

    const cards = [
        {
            label: 'Efectividad',
            value: percent(totals.delivered, closed),
            detail: `${numberFormat.format(totals.delivered)} entregadas`,
            dot: '#22c55e',
            text: 'text-emerald-600 dark:text-emerald-400',
            title: `${numberFormat.format(totals.delivered)} entregadas de ${numberFormat.format(closed)} entre entregadas y devueltas`,
        },
        {
            label: 'Devoluciones',
            value: percent(totals.returned, closed),
            detail: `${numberFormat.format(totals.returned)} devueltas`,
            dot: '#f59e0b',
            text: 'text-amber-600 dark:text-amber-400',
            title: `${numberFormat.format(totals.returned)} devueltas de ${numberFormat.format(closed)} entre entregadas y devueltas`,
        },
    ];

    return (
        <div className="hidden xl:flex items-center gap-2">
            {cards.map(card => (
                <div
                    key={card.label}
                    title={card.title}
                    className="flex items-center gap-2 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-2.5 py-1 shadow-sm"
                >
                    <span className="h-1.5 w-1.5 flex-shrink-0 rounded-full" style={{ backgroundColor: card.dot }} />
                    <div className="leading-tight">
                        <p className="text-[8px] font-semibold uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                            {card.label}
                        </p>
                        <p className={`text-sm font-bold tabular-nums ${card.text}`}>{card.value}</p>
                    </div>
                    <span className="whitespace-nowrap text-[9px] text-gray-400 dark:text-gray-500">{card.detail}</span>
                </div>
            ))}
        </div>
    );
}
