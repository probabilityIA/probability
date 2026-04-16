'use client';

import { useState, useEffect } from 'react';
import { CustomerSummary } from '../../domain/types';
import { getCustomerSummaryAction } from '../../infra/actions';
import { Spinner } from '@/shared/ui';

interface Props {
    customerId: number;
    businessId?: number;
}

function StatCard({ label, value, sub }: { label: string; value: string; sub?: string }) {
    return (
        <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4">
            <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">{label}</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">{value}</p>
            {sub && <p className="text-xs text-gray-400 mt-0.5">{sub}</p>}
        </div>
    );
}

const formatCurrency = (v: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', minimumFractionDigits: 0 }).format(v);

const formatDate = (d: string | null) =>
    d ? new Date(d).toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' }) : '--';

export default function CustomerSummaryTab({ customerId, businessId }: Props) {
    const [summary, setSummary] = useState<CustomerSummary | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        setLoading(true);
        setError(null);
        getCustomerSummaryAction(customerId, businessId)
            .then(setSummary)
            .catch((e: any) => setError(e.message || 'Error al cargar resumen'))
            .finally(() => setLoading(false));
    }, [customerId, businessId]);

    if (loading) return <div className="flex justify-center p-8"><Spinner size="lg" /></div>;
    if (error) return <p className="text-sm text-red-500 p-4">{error}</p>;
    if (!summary) return <p className="text-sm text-gray-400 p-4">Sin datos de resumen</p>;

    const deliveryPct = summary.total_orders > 0
        ? Math.round((summary.delivered_orders / summary.total_orders) * 100)
        : 0;

    return (
        <div className="space-y-4">
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                <StatCard label="Total ordenes" value={String(summary.total_orders)} />
                <StatCard label="Entregadas" value={String(summary.delivered_orders)} sub={`${deliveryPct}%`} />
                <StatCard label="En progreso" value={String(summary.in_progress_orders)} />
                <StatCard label="Canceladas" value={String(summary.cancelled_orders)} />
            </div>

            <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
                <StatCard label="Total gastado" value={formatCurrency(summary.total_spent)} />
                <StatCard label="Ticket promedio" value={formatCurrency(summary.avg_ticket)} />
                <StatCard label="Ordenes pagadas" value={String(summary.total_paid_orders)} />
            </div>

            <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
                <StatCard label="Score entrega" value={`${summary.avg_delivery_score.toFixed(1)}%`} />
                <StatCard label="Plataforma preferida" value={summary.preferred_platform || '--'} />
                <StatCard label="Primera orden" value={formatDate(summary.first_order_at)} />
            </div>

            <div className="text-xs text-gray-400 text-right">
                Ultima orden: {formatDate(summary.last_order_at)}
            </div>
        </div>
    );
}
