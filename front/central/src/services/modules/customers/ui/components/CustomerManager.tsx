'use client';

import { useState, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { CustomerInfo } from '../../domain/types';
import CustomerList from './CustomerList';
import CustomerForm from './CustomerForm';
import CustomerDetailView from './CustomerDetail';
import { Button } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

type ModalMode = 'create' | 'edit' | 'view' | null;

interface CustomerManagerProps {
    selectedBusinessId?: number | null;
}

export default function CustomerManager({ selectedBusinessId = null }: CustomerManagerProps) {
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

    // Gate para super admin: debe seleccionar negocio antes de operar
    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    return (
        <div className="space-y-4">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900">Clientes</h1>
                    <p className="text-sm text-gray-500 mt-0.5">
                        Gestiona los clientes de tu negocio
                    </p>
                </div>
                {!requiresBusinessSelection && (
                    <Button variant="primary" onClick={openCreate}>
                        <PlusIcon className="w-4 h-4 mr-2" />
                        Nuevo cliente
                    </Button>
                )}
            </div>

            {/* Gate: super admin debe seleccionar negocio */}
            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    <p className="text-gray-500 text-sm">Selecciona un negocio para ver y gestionar sus clientes</p>
                </div>
            ) : (
                <CustomerList
                    onEdit={openEdit}
                    onView={openView}
                    onRefreshRef={handleRefreshRef}
                    selectedBusinessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
                />
            )}

            {/* Modal crear / editar */}
            {(modalMode === 'create' || modalMode === 'edit') && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white rounded-xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b">
                            <h2 className="text-lg font-semibold text-gray-900">
                                {modalMode === 'create' ? 'Nuevo cliente' : 'Editar cliente'}
                            </h2>
                            <button
                                onClick={closeModal}
                                className="text-gray-400 hover:text-gray-600 text-xl leading-none"
                            >
                                ×
                            </button>
                        </div>
                        <div className="p-6">
                            <CustomerForm
                                customer={selectedCustomer ?? undefined}
                                onSuccess={handleFormSuccess}
                                onCancel={closeModal}
                                businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
                            />
                        </div>
                    </div>
                </div>
            )}

            {/* Modal detalle */}
            {modalMode === 'view' && selectedCustomer && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white rounded-xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b">
                            <div>
                                <h2 className="text-lg font-semibold text-gray-900">
                                    {selectedCustomer.name}
                                </h2>
                                <p className="text-sm text-gray-500">ID #{selectedCustomer.id}</p>
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
                                    className="text-gray-400 hover:text-gray-600 text-xl leading-none ml-3"
                                >
                                    ×
                                </button>
                            </div>
                        </div>
                        <div className="p-6">
                            <CustomerDetailView
                                customerId={selectedCustomer.id}
                                businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
                            />
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
