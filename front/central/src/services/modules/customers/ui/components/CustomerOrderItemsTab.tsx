'use client';

import { useState, useEffect } from 'react';
import { CustomerOrderItem, CustomerOrderItemListResponse } from '../../domain/types';
import { getCustomerOrderItemsAction } from '../../infra/actions';
import { Spinner } from '@/shared/ui';

interface Props {
    customerId: number;
    businessId?: number;
}

const formatCurrency = (v: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', minimumFractionDigits: 0 }).format(v);

const statusColors: Record<string, string> = {
    delivered: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    completed: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    cancelled: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
    refunded: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
    pending: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
    processing: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    picking: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    shipped: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400',
};

function StatusBadge({ status }: { status: string }) {
    const colors = statusColors[status] || 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300';
    return (
        <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${colors}`}>
            {status}
        </span>
    );
}

export default function CustomerOrderItemsTab({ customerId, businessId }: Props) {
    const [data, setData] = useState<CustomerOrderItemListResponse | null>(null);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        setLoading(true);
        setError(null);
        getCustomerOrderItemsAction(customerId, { page, page_size: 10, business_id: businessId })
            .then(setData)
            .catch((e: any) => setError(e.message || 'Error al cargar items'))
            .finally(() => setLoading(false));
    }, [customerId, businessId, page]);

    if (loading) return <div className="flex justify-center p-8"><Spinner size="lg" /></div>;
    if (error) return <p className="text-sm text-red-500 p-4">{error}</p>;
    if (!data || data.data.length === 0) return <p className="text-sm text-gray-400 p-4">Sin items de ordenes</p>;

    return (
        <div className="space-y-3">
            {data.data.map((item: CustomerOrderItem) => (
                <div key={item.id} className="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                    <div className="flex items-start justify-between gap-2 mb-2">
                        <div className="min-w-0">
                            <p className="text-sm font-medium text-gray-900 dark:text-white truncate">
                                {item.product_name || 'Producto sin nombre'}
                            </p>
                            <p className="text-xs text-gray-400">
                                Orden {item.order_number || item.order_id.slice(0, 8)}
                            </p>
                        </div>
                        <StatusBadge status={item.order_status} />
                    </div>
                    <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
                        <span>Cant: {item.quantity}</span>
                        <span>Unitario: {formatCurrency(item.unit_price)}</span>
                        <span className="font-medium text-gray-700 dark:text-gray-300">Total: {formatCurrency(item.total_price)}</span>
                        <span>{new Date(item.ordered_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' })}</span>
                    </div>
                </div>
            ))}

            {data.total_pages > 1 && (
                <div className="flex items-center justify-between pt-2">
                    <button
                        onClick={() => setPage(p => Math.max(1, p - 1))}
                        disabled={page <= 1}
                        className="text-xs text-blue-600 hover:underline disabled:text-gray-300 disabled:no-underline"
                    >
                        Anterior
                    </button>
                    <span className="text-xs text-gray-400">{page} / {data.total_pages}</span>
                    <button
                        onClick={() => setPage(p => Math.min(data.total_pages, p + 1))}
                        disabled={page >= data.total_pages}
                        className="text-xs text-blue-600 hover:underline disabled:text-gray-300 disabled:no-underline"
                    >
                        Siguiente
                    </button>
                </div>
            )}
        </div>
    );
}
