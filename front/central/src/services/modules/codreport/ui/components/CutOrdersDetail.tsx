'use client';

import { useEffect, useState } from 'react';
import { RefreshCw, Package } from 'lucide-react';
import { getCutOrdersAction, getSelectableOrdersAction } from '../../infra/actions';
import { CodOrder } from '../../domain/types';
import { formatMoney, formatDateTime, carrierLabel } from './helpers';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';

const STATUS_PILL: Record<string, { bg: string; c: string; label: string }> = {
    pending: { bg: '#f1f5f9', c: '#475569', label: 'Pendiente' },
    picked_up: { bg: '#eef2ff', c: '#4f46e5', label: 'Recogido' },
    in_transit: { bg: '#eef2ff', c: '#4f46e5', label: 'En transito' },
    out_for_delivery: { bg: '#eef2ff', c: '#4338ca', label: 'En reparto' },
    on_hold: { bg: '#fef3c7', c: '#b45309', label: 'En espera' },
    delivered: { bg: '#dcfce7', c: '#15803d', label: 'Entregado' },
    failed: { bg: '#fee2e2', c: '#b91c1c', label: 'Fallido' },
    returned: { bg: '#fee2e2', c: '#b91c1c', label: 'Devuelto' },
    cancelled: { bg: '#f1f5f9', c: '#64748b', label: 'Cancelado' },
};

const CHIP_COLORS = ['#e11d48', '#ea580c', '#0891b2', '#7c3aed', '#0ea5e9', '#16a34a', '#db2777', '#4f46e5', '#f59e0b', '#14b8a6'];
function chipColor(name: string): string {
    let h = 0;
    for (let i = 0; i < name.length; i++) h = (h * 31 + name.charCodeAt(i)) >>> 0;
    return CHIP_COLORS[h % CHIP_COLORS.length];
}

interface Props {
    cutId?: number;
    periodStart?: string;
    periodEnd?: string;
    businessId?: number | null;
}

export function CutOrdersDetail({ cutId, periodStart, periodEnd, businessId }: Props) {
    const [orders, setOrders] = useState<CodOrder[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        let cancelled = false;
        setLoading(true);
        setError(null);
        const req = cutId
            ? getCutOrdersAction(cutId, businessId || undefined)
            : getSelectableOrdersAction(periodStart || '', periodEnd || '', businessId || undefined);
        req.then(res => {
            if (cancelled) return;
            if (res.success) setOrders((res.data || []) as CodOrder[]);
            else setError((res as any).message || 'Error al cargar las ordenes del corte');
            setLoading(false);
        });
        return () => { cancelled = true; };
    }, [cutId, periodStart, periodEnd, businessId]);

    if (loading) {
        return (
            <div className="flex items-center justify-center py-8 text-gray-400 text-sm">
                <RefreshCw size={16} className="animate-spin mr-2" /> Cargando ordenes del corte...
            </div>
        );
    }

    if (error) {
        return <div className="text-xs text-red-600 py-3">{error}</div>;
    }

    if (orders.length === 0) {
        return (
            <div className="text-center py-8 text-gray-400 text-sm">
                <Package size={24} className="mx-auto mb-2 opacity-50" />
                Este corte no tiene ordenes registradas.
            </div>
        );
    }

    return (
        <div className="overflow-x-auto">
            <table className="w-full text-xs min-w-[760px]">
                <thead>
                    <tr className="text-[10px] uppercase text-gray-400 dark:text-gray-500">
                        <th className="text-left px-2 py-1.5 font-semibold">Orden</th>
                        <th className="text-left px-2 py-1.5 font-semibold">Cliente</th>
                        <th className="text-left px-2 py-1.5 font-semibold">Transportadora</th>
                        <th className="text-left px-2 py-1.5 font-semibold">Estado</th>
                        <th className="text-right px-2 py-1.5 font-semibold">COD orden</th>
                        <th className="text-right px-2 py-1.5 font-semibold">Cargo carrier</th>
                        <th className="text-right px-2 py-1.5 font-semibold">Total cliente</th>
                        <th className="text-left pl-6 pr-2 py-1.5 font-semibold">Entregado</th>
                    </tr>
                </thead>
                <tbody>
                    {orders.map(o => {
                        const cc = chipColor(o.carrier || '');
                        const logo = getCarrierLogo(o.carrier);
                        return (
                            <tr key={o.order_id} className="border-t border-gray-100 dark:border-gray-700/50">
                                <td className="px-2 py-2 font-mono text-[12px] font-bold text-[#6d28d9] whitespace-nowrap">#{o.order_number || o.order_id.slice(0, 8)}</td>
                                <td className="px-2 py-2">
                                    <div className="font-semibold text-gray-800 dark:text-gray-200 truncate max-w-[160px]">{o.customer_name || '-'}</div>
                                    <div className="text-[10.5px] text-gray-400 whitespace-nowrap">{formatDateTime(o.created_at)}</div>
                                </td>
                                <td className="px-2 py-2">
                                    <div className="flex items-center gap-1.5">
                                        {logo ? (
                                            <span className="inline-flex items-center justify-center h-6 w-9 rounded border border-gray-200 dark:border-gray-600 bg-white p-0.5 shrink-0" title={carrierLabel(o.carrier)}>
                                                <img src={logo} alt={carrierLabel(o.carrier)} className="max-h-full max-w-full object-contain" />
                                            </span>
                                        ) : (
                                            <span className="w-[18px] h-[18px] rounded text-white inline-flex items-center justify-center text-[9px] font-extrabold shrink-0" style={{ background: cc }}>
                                                {(o.carrier || '?').charAt(0).toUpperCase()}
                                            </span>
                                        )}
                                        <span className="text-gray-700 dark:text-gray-300 whitespace-nowrap">{carrierLabel(o.carrier)}</span>
                                    </div>
                                </td>
                                <td className="px-2 py-2">
                                    {(() => {
                                        const st = STATUS_PILL[o.status] || { bg: '#f1f5f9', c: '#475569', label: o.status };
                                        return <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-bold whitespace-nowrap" style={{ background: st.bg, color: st.c }}>{st.label}</span>;
                                    })()}
                                </td>
                                <td className="px-2 py-2 text-right text-gray-600 dark:text-gray-300 tabular-nums whitespace-nowrap">{formatMoney(o.cod_total, o.currency)}</td>
                                <td className="px-2 py-2 text-right text-[#c2410c] tabular-nums whitespace-nowrap">{o.cod_carrier_fee > 0 ? formatMoney(o.cod_carrier_fee, o.currency) : '-'}</td>
                                <td className="px-2 py-2 text-right font-bold text-gray-900 dark:text-white tabular-nums whitespace-nowrap">{formatMoney(o.cod_total + (o.cod_carrier_fee || 0), o.currency)}</td>
                                <td className="pl-6 pr-2 py-2 text-[11px] text-gray-500 whitespace-nowrap">{o.delivered_at ? formatDateTime(o.delivered_at) : '-'}</td>
                            </tr>
                        );
                    })}
                </tbody>
            </table>
        </div>
    );
}
