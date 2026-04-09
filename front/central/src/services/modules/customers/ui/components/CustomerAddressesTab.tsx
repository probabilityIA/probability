'use client';

import { useState, useEffect } from 'react';
import { CustomerAddress, CustomerAddressListResponse } from '../../domain/types';
import { getCustomerAddressesAction } from '../../infra/actions';
import { Spinner } from '@/shared/ui';
import { MapPinIcon } from '@heroicons/react/24/outline';

interface Props {
    customerId: number;
    businessId?: number;
}

export default function CustomerAddressesTab({ customerId, businessId }: Props) {
    const [data, setData] = useState<CustomerAddressListResponse | null>(null);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        setLoading(true);
        setError(null);
        getCustomerAddressesAction(customerId, { page, page_size: 10, business_id: businessId })
            .then(setData)
            .catch((e: any) => setError(e.message || 'Error al cargar direcciones'))
            .finally(() => setLoading(false));
    }, [customerId, businessId, page]);

    if (loading) return <div className="flex justify-center p-8"><Spinner size="lg" /></div>;
    if (error) return <p className="text-sm text-red-500 p-4">{error}</p>;
    if (!data || data.data.length === 0) return <p className="text-sm text-gray-400 p-4">Sin direcciones registradas</p>;

    return (
        <div className="space-y-3">
            {data.data.map((addr: CustomerAddress) => (
                <div key={addr.id} className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 flex gap-3">
                    <MapPinIcon className="w-5 h-5 text-gray-400 flex-shrink-0 mt-0.5" />
                    <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-gray-900 dark:text-white truncate">
                            {addr.street || '--'}
                        </p>
                        <p className="text-xs text-gray-500 dark:text-gray-400">
                            {[addr.city, addr.state, addr.country].filter(Boolean).join(', ')}
                            {addr.postal_code ? ` - ${addr.postal_code}` : ''}
                        </p>
                        <div className="flex gap-4 mt-1 text-xs text-gray-400">
                            <span>Usada {addr.times_used} {addr.times_used === 1 ? 'vez' : 'veces'}</span>
                            <span>Ultima: {new Date(addr.last_used_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' })}</span>
                        </div>
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
