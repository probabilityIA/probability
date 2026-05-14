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

function barColor(rate: number, estimated?: boolean): string {
    if (estimated) return 'bg-gray-400';
    if (rate >= 0.9) return 'bg-emerald-500';
    if (rate >= 0.75) return 'bg-amber-500';
    return 'bg-red-500';
}

interface BarProps {
    label: string;
    rate?: number;
    sample?: number;
    estimated?: boolean;
    title?: string;
}

function Bar({ label, rate, sample, estimated, title }: BarProps) {
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
        <div title={title}>
            <div className="flex items-baseline justify-between mb-0.5">
                <p className="text-[10px] text-gray-600 dark:text-gray-300">
                    {label}
                    {estimated && (
                        <span className="ml-1 px-1 py-px rounded text-[9px] font-semibold bg-gray-100 text-gray-500 border border-gray-200">
                            estimado
                        </span>
                    )}
                </p>
                <p className={`text-xs font-bold tabular-nums ${estimated ? 'text-gray-500 italic' : 'text-gray-900 dark:text-white'}`}>
                    {pct}%
                </p>
            </div>
            <div className={`w-full h-1 rounded-full overflow-hidden ${estimated ? 'bg-gray-100' : 'bg-gray-200 dark:bg-gray-700'}`}>
                <div
                    className={`h-full ${barColor(rate, estimated)} ${estimated ? 'opacity-60' : ''}`}
                    style={{ width: `${pct}%` }}
                />
            </div>
            {!estimated && sample !== undefined && sample > 0 && (
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

    const zoneRate = result?.delivery_rate;
    const zoneEstimated = !result?.found || result?.is_estimated;
    const globalRate = result?.global_rate;
    const globalEstimated = (result?.global_total ?? 0) < 20;

    const estimateTooltip = result?.estimate_source === 'global_carrier'
        ? 'Estimado: tasa global de este transportador (aun no hay datos de tu zona).'
        : result?.estimate_source === 'carrier_baseline'
            ? 'Estimado preliminar de la transportadora. Se ajusta con tus envios reales.'
            : undefined;

    return (
        <div className="space-y-2 w-full">
            <Bar
                label="Efectividad de recoleccion"
                rate={globalRate}
                sample={result?.global_total}
                estimated={globalRate !== undefined && globalEstimated}
                title={globalEstimated ? 'Muestra pequena, valor preliminar.' : undefined}
            />
            <Bar
                label="Efectividad de entrega en zona"
                rate={zoneRate}
                sample={result?.stats?.total}
                estimated={zoneEstimated}
                title={estimateTooltip}
            />
        </div>
    );
}
