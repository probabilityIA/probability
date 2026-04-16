'use client';

import { useState, useEffect } from 'react';
import { CustomerDetail } from '../../domain/types';
import { getCustomerByIdAction } from '../../infra/actions';
import { Spinner, Alert } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';
import CustomerSummaryTab from './CustomerSummaryTab';
import CustomerAddressesTab from './CustomerAddressesTab';
import CustomerProductsTab from './CustomerProductsTab';
import CustomerOrderItemsTab from './CustomerOrderItemsTab';

interface CustomerDetailProps {
    customerId: number;
    businessId?: number;
}

type Tab = 'info' | 'summary' | 'addresses' | 'products' | 'orders';

const tabs: { key: Tab; label: string }[] = [
    { key: 'info', label: 'Info' },
    { key: 'summary', label: 'Resumen' },
    { key: 'addresses', label: 'Direcciones' },
    { key: 'products', label: 'Productos' },
    { key: 'orders', label: 'Ordenes' },
];

function Field({ label, value }: { label: string; value?: string | null }) {
    return (
        <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">{label}</p>
            <p className="text-sm text-gray-900 dark:text-white">{value || <span className="text-gray-300">--</span>}</p>
        </div>
    );
}

export default function CustomerDetailView({ customerId, businessId }: CustomerDetailProps) {
    const [customer, setCustomer] = useState<CustomerDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [activeTab, setActiveTab] = useState<Tab>('info');

    useEffect(() => {
        setLoading(true);
        setError(null);
        getCustomerByIdAction(customerId, businessId)
            .then(setCustomer)
            .catch((err: any) => setError(getActionError(err, 'Error al cargar el cliente')))
            .finally(() => setLoading(false));
    }, [customerId, businessId]);

    if (loading) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    if (error) {
        return <Alert type="error" onClose={() => setError(null)}>{error}</Alert>;
    }

    if (!customer) return null;

    return (
        <div className="space-y-4">
            <div className="flex gap-1 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
                {tabs.map(tab => (
                    <button
                        key={tab.key}
                        onClick={() => setActiveTab(tab.key)}
                        className={`px-3 py-2 text-sm font-medium whitespace-nowrap border-b-2 transition-colors ${
                            activeTab === tab.key
                                ? 'border-purple-600 text-purple-600 dark:text-purple-400'
                                : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
                        }`}
                    >
                        {tab.label}
                    </button>
                ))}
            </div>

            {activeTab === 'info' && (
                <div className="space-y-4">
                    <div className="grid grid-cols-3 gap-3">
                        <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4 text-center">
                            <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Ordenes</p>
                            <p className="text-lg font-semibold text-gray-900 dark:text-white">{customer.order_count}</p>
                        </div>
                        <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4 text-center">
                            <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Total gastado</p>
                            <p className="text-lg font-semibold text-gray-900 dark:text-white">
                                {new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', minimumFractionDigits: 0 }).format(customer.total_spent)}
                            </p>
                        </div>
                        <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4 text-center">
                            <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Ultima orden</p>
                            <p className="text-lg font-semibold text-gray-900 dark:text-white">
                                {customer.last_order_at
                                    ? new Date(customer.last_order_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' })
                                    : '--'}
                            </p>
                        </div>
                    </div>

                    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-5 space-y-4">
                        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-200 uppercase tracking-wide border-b pb-2">
                            Informacion de contacto
                        </h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <Field label="Nombre" value={customer.name} />
                            <Field label="Email" value={customer.email} />
                            <Field label="Telefono" value={customer.phone} />
                            <Field label="Documento" value={customer.dni} />
                        </div>
                    </div>

                    <div className="text-xs text-gray-400 flex gap-4">
                        <span>
                            Creado: {new Date(customer.created_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' })}
                        </span>
                        <span>
                            Actualizado: {new Date(customer.updated_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' })}
                        </span>
                    </div>
                </div>
            )}

            {activeTab === 'summary' && (
                <CustomerSummaryTab customerId={customerId} businessId={businessId} />
            )}

            {activeTab === 'addresses' && (
                <CustomerAddressesTab customerId={customerId} businessId={businessId} />
            )}

            {activeTab === 'products' && (
                <CustomerProductsTab customerId={customerId} businessId={businessId} />
            )}

            {activeTab === 'orders' && (
                <CustomerOrderItemsTab customerId={customerId} businessId={businessId} />
            )}
        </div>
    );
}
