'use client';

import { useEffect, useState } from 'react';
import { getDeliveryProbabilityAction } from '../../infra/actions';
import type { ProbabilityResult } from '../../domain/types';

interface Props {
    businessId: number;
    orderId?: string;
    lat?: number;
    lng?: number;
    carrier?: string;
    label?: string;
    compact?: boolean;
}

function colorFor(rate: number): string {
    if (rate >= 0.9) return 'bg-emerald-100 text-emerald-800 border-emerald-300 dark:bg-emerald-900/30 dark:text-emerald-300 dark:border-emerald-700';
    if (rate >= 0.75) return 'bg-amber-100 text-amber-800 border-amber-300 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-700';
    return 'bg-red-100 text-red-800 border-red-300 dark:bg-red-900/30 dark:text-red-300 dark:border-red-700';
}

export function DeliveryProbabilityBadge({ businessId, orderId, lat, lng, carrier, label = 'Probabilidad de entrega', compact = false }: Props) {
    const [result, setResult] = useState<ProbabilityResult | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        let cancelled = false;
        if (!businessId || (!orderId && (lat === undefined || lng === undefined))) {
            setLoading(false);
            return;
        }
        setLoading(true);
        setError(null);
        getDeliveryProbabilityAction({ business_id: businessId, order_id: orderId, lat, lng, carrier })
            .then((res) => { if (!cancelled) setResult(res); })
            .catch((err) => { if (!cancelled) setError(err.message); })
            .finally(() => { if (!cancelled) setLoading(false); });
        return () => { cancelled = true; };
    }, [businessId, orderId, lat, lng, carrier]);

    if (loading) {
        return (
            <span className={`inline-flex items-center gap-1 px-2 py-1 rounded-md border text-xs font-medium bg-gray-50 text-gray-500 border-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-700 ${compact ? '' : ''}`}>
                <span className="animate-pulse">{label}: ...</span>
            </span>
        );
    }
    if (error) return null;
    if (!result || !result.found || result.delivery_rate === undefined) {
        return (
            <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md border text-xs font-medium bg-gray-50 text-gray-500 border-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-700">
                {label}: sin datos
            </span>
        );
    }
    const pct = Math.round(result.delivery_rate * 100);
    return (
        <span
            className={`inline-flex items-center gap-1 px-2 py-1 rounded-md border text-xs font-semibold ${colorFor(result.delivery_rate)}`}
            title={result.stats?.geozone_name ? `Zona: ${result.stats.geozone_name}` : ''}
        >
            {compact ? `${pct}%` : `${label}: ${pct}%`}
        </span>
    );
}
