'use client';

import { useEffect, useState } from 'react';
import { getProbabilityByCarrierAction } from '../../infra/actions';
import type { ProbabilityResult } from '../../domain/types';
import { getCarrierLogo } from '@/shared/utils/carrier-logo';

interface Props {
    businessId: number;
    orderId: string;
}

function colorFor(rate: number): string {
    if (rate >= 0.9) return 'bg-emerald-500';
    if (rate >= 0.75) return 'bg-amber-500';
    return 'bg-red-500';
}

function textColorFor(rate: number): string {
    if (rate >= 0.9) return 'text-emerald-700 dark:text-emerald-300';
    if (rate >= 0.75) return 'text-amber-700 dark:text-amber-300';
    return 'text-red-700 dark:text-red-400';
}

export function DeliveryProbabilityByCarrier({ businessId, orderId }: Props) {
    const [results, setResults] = useState<ProbabilityResult[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        let cancelled = false;
        if (!businessId || !orderId) return;
        setLoading(true);
        getProbabilityByCarrierAction(orderId, businessId)
            .then((res) => { if (!cancelled) setResults(res); })
            .catch(() => { if (!cancelled) setResults([]); })
            .finally(() => { if (!cancelled) setLoading(false); });
        return () => { cancelled = true; };
    }, [businessId, orderId]);

    if (loading) {
        return <div className="text-xs text-gray-500 dark:text-gray-400">Cargando probabilidades por transportadora...</div>;
    }
    if (!results.length) return null;

    const sorted = [...results].sort((a, b) => {
        if (a.found && !b.found) return -1;
        if (!a.found && b.found) return 1;
        return (b.delivery_rate || 0) - (a.delivery_rate || 0);
    });

    return (
        <div className="space-y-2">
            <p className="text-[10px] font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Probabilidad por transportadora
            </p>
            <ul className="space-y-1">
                {sorted.map((r) => {
                    const logo = getCarrierLogo(r.carrier || '');
                    const initials = (r.carrier || '?').slice(0, 3).toUpperCase();
                    const rate = r.delivery_rate;
                    const pct = rate !== undefined ? Math.round(rate * 100) : null;
                    return (
                        <li
                            key={r.carrier}
                            className="flex items-center gap-2 px-2 py-1.5 bg-white dark:bg-gray-800 rounded border border-gray-200 dark:border-gray-700"
                            title={r.stats?.geozone_name ? `Zona: ${r.stats.geozone_name} (${r.stats.total} envios)` : 'Sin muestra suficiente'}
                        >
                            <div className="w-7 h-7 rounded bg-gray-50 dark:bg-gray-700 flex items-center justify-center overflow-hidden flex-shrink-0">
                                {logo ? (
                                    <img
                                        src={logo}
                                        alt={r.carrier}
                                        className="w-full h-full object-contain"
                                        onError={(e) => {
                                            e.currentTarget.style.display = 'none';
                                            e.currentTarget.parentElement!.innerHTML = `<span class="text-[9px] font-bold text-gray-500">${initials}</span>`;
                                        }}
                                    />
                                ) : (
                                    <span className="text-[9px] font-bold text-gray-500">{initials}</span>
                                )}
                            </div>
                            <div className="flex-1 text-xs font-medium text-gray-900 dark:text-white truncate">{r.carrier}</div>
                            {r.found && rate !== undefined && pct !== null ? (
                                <div className="flex items-center gap-1.5 flex-shrink-0">
                                    <div className="w-12 h-1.5 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                                        <div className={`h-full ${colorFor(rate)}`} style={{ width: `${pct}%` }} />
                                    </div>
                                    <span className={`text-xs font-bold tabular-nums ${textColorFor(rate)}`}>{pct}%</span>
                                </div>
                            ) : (
                                <span className="text-[10px] text-gray-400">sin datos</span>
                            )}
                        </li>
                    );
                })}
            </ul>
        </div>
    );
}
