'use client';

import { useState, useEffect } from 'react';
import { CustomerProduct, CustomerProductListResponse } from '../../domain/types';
import { getCustomerProductsAction } from '../../infra/actions';
import { Spinner } from '@/shared/ui';

interface Props {
    customerId: number;
    businessId?: number;
}

const formatCurrency = (v: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', minimumFractionDigits: 0 }).format(v);

export default function CustomerProductsTab({ customerId, businessId }: Props) {
    const [data, setData] = useState<CustomerProductListResponse | null>(null);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        setLoading(true);
        setError(null);
        getCustomerProductsAction(customerId, { page, page_size: 10, business_id: businessId })
            .then(setData)
            .catch((e: any) => setError(e.message || 'Error al cargar productos'))
            .finally(() => setLoading(false));
    }, [customerId, businessId, page]);

    if (loading) return <div className="flex justify-center p-8"><Spinner size="lg" /></div>;
    if (error) return <p className="text-sm text-red-500 p-4">{error}</p>;
    if (!data || data.data.length === 0) return <p className="text-sm text-gray-400 p-4">Sin productos comprados</p>;

    return (
        <div className="space-y-3">
            {data.data.map((prod: CustomerProduct) => (
                <div key={prod.id} className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 flex gap-3">
                    {prod.product_image ? (
                        <img
                            src={prod.product_image}
                            alt={prod.product_name}
                            className="w-12 h-12 rounded-lg object-cover flex-shrink-0"
                        />
                    ) : (
                        <div className="w-12 h-12 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center flex-shrink-0">
                            <span className="text-gray-400 text-xs">IMG</span>
                        </div>
                    )}
                    <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-gray-900 dark:text-white truncate">
                            {prod.product_name || prod.product_sku || prod.product_id}
                        </p>
                        {prod.product_sku && (
                            <p className="text-xs text-gray-400">SKU: {prod.product_sku}</p>
                        )}
                        <div className="flex flex-wrap gap-x-4 gap-y-1 mt-1 text-xs text-gray-500 dark:text-gray-400">
                            <span>{prod.times_ordered} {prod.times_ordered === 1 ? 'orden' : 'ordenes'}</span>
                            <span>{prod.total_quantity} uds</span>
                            <span>{formatCurrency(prod.total_spent)}</span>
                        </div>
                    </div>
                    <div className="text-right text-xs text-gray-400 flex-shrink-0">
                        <p>Primera: {new Date(prod.first_ordered_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short' })}</p>
                        <p>Ultima: {new Date(prod.last_ordered_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short' })}</p>
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
