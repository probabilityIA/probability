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
    const zoneEstimated = !result?.found || !!result?.is_estimated;
    const sample = result?.stats?.total;

    const levelLabel: Record<string, string> = {
        barrio: 'barrio',
        neighborhood: 'UPZ',
        admin_district: 'localidad',
        locality: 'corregimiento',
        city: 'municipio',
        state: 'departamento',
        country: 'pais',
    };
    const lvl = result?.level ? (levelLabel[result.level] || result.level) : '';
    const zoneName = result?.stats?.geozone_name || '';

    let source = '';
    if (result?.found && sample) {
        source = zoneName ? `${lvl} ${zoneName} - ${sample} envios` : `${lvl} - ${sample} envios`;
    } else if (result?.estimate_source === 'global_carrier' && result?.global_total) {
        source = `tasa nacional - ${result.global_total} envios`;
    } else if (result?.estimate_source === 'carrier_baseline') {
        source = 'transportadora nueva, sin historial';
    }

    const tooltip = result?.estimate_source === 'global_carrier'
        ? 'Aun no tenemos envios en tu zona con este carrier; mostramos la tasa nacional del carrier.'
        : result?.estimate_source === 'zone_low_sample'
            ? 'Muestra pequena en esta zona; se ajusta a medida que llegan mas envios.'
            : result?.estimate_source === 'carrier_baseline'
                ? 'Transportadora sin historial en el sistema. Valor preliminar.'
                : undefined;

    return (
        <div className="space-y-1 w-full">
            <Bar
                label="Efectividad de entrega en zona"
                rate={zoneRate}
                sample={sample}
                estimated={zoneEstimated}
                title={tooltip}
            />
            {source && (
                <p className="text-[9px] text-gray-400 truncate" title={source}>{source}</p>
            )}
        </div>
    );
}
