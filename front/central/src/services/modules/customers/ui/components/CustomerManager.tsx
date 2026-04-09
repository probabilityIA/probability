'use client';

import { useState, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { CustomerInfo } from '../../domain/types';
import CustomerList from './CustomerList';
import CustomerForm from './CustomerForm';
import CustomerDetailView from './CustomerDetail';
import CustomerSummaryTab from './CustomerSummaryTab';
import CustomerAddressesTab from './CustomerAddressesTab';
import CustomerProductsTab from './CustomerProductsTab';
import CustomerOrderItemsTab from './CustomerOrderItemsTab';
import { Button, SuperAdminBusinessSelector } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

type ModalMode = 'create' | 'edit' | 'view' | 'summary' | 'addresses' | 'products' | 'orders' | null;

interface CustomerManagerProps {
    selectedBusinessId?: number | null;
    onBusinessChange?: (businessId: number | null) => void;
}

const modalTitles: Record<string, string> = {
    summary: 'Resumen',
    addresses: 'Direcciones',
    products: 'Productos',
    orders: 'Ordenes',
};

export default function CustomerManager({ selectedBusinessId = null, onBusinessChange }: CustomerManagerProps) {
    const { isSuperAdmin } = usePermissions();
    const [modalMode, setModalMode] = useState<ModalMode>(null);
    const [selectedCustomer, setSelectedCustomer] = useState<CustomerInfo | null>(null);
    const [refreshList, setRefreshList] = useState<(() => void) | null>(null);

    const openCreate = () => {
        setSelectedCustomer(null);
        setModalMode('create');
    };

    const openEdit = (customer: CustomerInfo) => {
        setSelectedCustomer(customer);
        setModalMode('edit');
    };

    const openView = (customer: CustomerInfo) => {
        setSelectedCustomer(customer);
        setModalMode('view');
    };

    const openModal = (customer: CustomerInfo, mode: ModalMode) => {
        setSelectedCustomer(customer);
        setModalMode(mode);
    };

    const closeModal = () => {
        setModalMode(null);
        setSelectedCustomer(null);
    };

    const handleFormSuccess = () => {
        closeModal();
        refreshList?.();
    };

    const handleRefreshRef = useCallback((ref: () => void) => {
        setRefreshList(() => ref);
    }, []);

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    const renderHistoryModal = (mode: 'summary' | 'addresses' | 'products' | 'orders') => {
        if (!selectedCustomer) return null;
        const components = {
            summary: <CustomerSummaryTab customerId={selectedCustomer.id} businessId={businessId} />,
            addresses: <CustomerAddressesTab customerId={selectedCustomer.id} businessId={businessId} />,
            products: <CustomerProductsTab customerId={selectedCustomer.id} businessId={businessId} />,
            orders: <CustomerOrderItemsTab customerId={selectedCustomer.id} businessId={businessId} />,
        };
        return (
            <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
                    <div className="flex items-center justify-between px-6 py-4 border-b">
                        <div>
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                                {modalTitles[mode]} - {selectedCustomer.name}
                            </h2>
                            <p className="text-sm text-gray-500 dark:text-gray-400">ID #{selectedCustomer.id}</p>
                        </div>
                        <button
                            onClick={closeModal}
                            className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl leading-none"
                        >
                            x
                        </button>
                    </div>
                    <div className="p-6">
                        {components[mode]}
                    </div>
                </div>
            </div>
        );
    };

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Clientes</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        Gestiona los clientes de tu negocio
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    {isSuperAdmin && (
                        <SuperAdminBusinessSelector
                            value={selectedBusinessId ?? null}
                            onChange={onBusinessChange || (() => {})}
                            variant="default"
                            placeholder="-- Selecciona un negocio --"
                        />
                    )}
                    {isSuperAdmin && !requiresBusinessSelection && (
                        <button
                            onClick={openCreate}
                            className="inline-flex items-center justify-center px-6 py-3 font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-all duration-300 hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-offset-2"
                        >
                            <PlusIcon className="w-4 h-4 mr-2" />
                            Nuevo cliente
                        </button>
                    )}
                </div>
            </div>

            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    <p className="text-gray-500 dark:text-gray-400 text-sm">Selecciona un negocio para ver y gestionar sus clientes</p>
                </div>
            ) : (
                <CustomerList
                    onEdit={openEdit}
                    onView={openView}
                    onViewSummary={(c) => openModal(c, 'summary')}
                    onViewAddresses={(c) => openModal(c, 'addresses')}
                    onViewProducts={(c) => openModal(c, 'products')}
                    onViewOrders={(c) => openModal(c, 'orders')}
                    onRefreshRef={handleRefreshRef}
                    selectedBusinessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
                />
            )}

            {(modalMode === 'create' || modalMode === 'edit') && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b">
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                                {modalMode === 'create' ? 'Nuevo cliente' : 'Editar cliente'}
                            </h2>
                            <button
                                onClick={closeModal}
                                className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl leading-none"
                            >
                                x
                            </button>
                        </div>
                        <div className="p-6">
                            <CustomerForm
                                customer={selectedCustomer ?? undefined}
                                onSuccess={handleFormSuccess}
                                onCancel={closeModal}
                                businessId={businessId}
                            />
                        </div>
                    </div>
                </div>
            )}

            {modalMode === 'view' && selectedCustomer && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b">
                            <div>
                                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                                    {selectedCustomer.name}
                                </h2>
                                <p className="text-sm text-gray-500 dark:text-gray-400">ID #{selectedCustomer.id}</p>
                            </div>
                            <div className="flex gap-2">
                                <button
                                    onClick={() => openEdit(selectedCustomer)}
                                    className="text-sm text-blue-600 hover:underline"
                                >
                                    Editar
                                </button>
                                <button
                                    onClick={closeModal}
                                    className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl leading-none ml-3"
                                >
                                    x
                                </button>
                            </div>
                        </div>
                        <div className="p-6">
                            <CustomerDetailView
                                customerId={selectedCustomer.id}
                                businessId={businessId}
                            />
                        </div>
                    </div>
                </div>
            )}

            {modalMode === 'summary' && renderHistoryModal('summary')}
            {modalMode === 'addresses' && renderHistoryModal('addresses')}
            {modalMode === 'products' && renderHistoryModal('products')}
            {modalMode === 'orders' && renderHistoryModal('orders')}
        </div>
    );
}
