'use client';

import { useEffect, useState } from 'react';
import { getDeliveryProbabilityAction } from '../../infra/actions';
import type { ProbabilityResult } from '../../domain/types';

interface Props {
    businessId: number;
    orderId?: string;
    lat?: number;
    lng?: number;
    carrier: string;
}

function barColor(rate: number): string {
    if (rate >= 0.9) return 'bg-emerald-500';
    if (rate >= 0.75) return 'bg-amber-500';
    return 'bg-red-500';
}

function Bar({ label, rate, sample }: { label: string; rate?: number; sample?: number }) {
    if (rate === undefined) {
        return (
            <div>
                <p className="text-[10px] text-gray-600 dark:text-gray-300 mb-0.5">{label}</p>
                <p className="text-[10px] text-gray-400">sin datos</p>
            </div>
        );
    }
    const pct = Math.round(rate * 100);
    return (
        <div>
            <div className="flex items-baseline justify-between mb-0.5">
                <p className="text-[10px] text-gray-600 dark:text-gray-300">{label}</p>
                <p className="text-xs font-bold tabular-nums text-gray-900 dark:text-white">{pct}%</p>
            </div>
            <div className="w-full h-1 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                <div className={`h-full ${barColor(rate)}`} style={{ width: `${pct}%` }} />
            </div>
            {sample !== undefined && sample > 0 && (
                <p className="text-[9px] text-gray-500 mt-0.5 tabular-nums">
                    {(rate * 10).toFixed(1)} de cada 10 envios
                </p>
            )}
        </div>
    );
}

export function CarrierEffectivenessRates({ businessId, orderId, lat, lng, carrier }: Props) {
    const [result, setResult] = useState<ProbabilityResult | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        let cancelled = false;
        if (!businessId || (!orderId && (lat === undefined || lng === undefined))) {
            setLoading(false);
            return;
        }
        setLoading(true);
        getDeliveryProbabilityAction({ business_id: businessId, order_id: orderId, lat, lng, carrier })
            .then((res) => { if (!cancelled) setResult(res); })
            .catch(() => { if (!cancelled) setResult(null); })
            .finally(() => { if (!cancelled) setLoading(false); });
        return () => { cancelled = true; };
    }, [businessId, orderId, lat, lng, carrier]);

    if (loading) {
        return <div className="text-[10px] text-gray-400 animate-pulse">Cargando efectividad...</div>;
    }
    return (
        <div className="space-y-2 w-full">
            <Bar
                label="Efectividad de recoleccion"
                rate={result?.global_rate}
                sample={result?.global_total}
            />
            <Bar
                label="Efectividad de entrega en zona"
                rate={result?.found ? result?.delivery_rate : undefined}
                sample={result?.stats?.total}
            />
        </div>
    );
}
