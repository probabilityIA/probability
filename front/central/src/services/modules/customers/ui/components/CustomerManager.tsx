'use client';

import { useState, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { CustomerInfo } from '../../domain/types';
import CustomerList from './CustomerList';
import CustomerForm from './CustomerForm';
import CustomerDetailView from './CustomerDetail';
import { Button } from '@/shared/ui';

type ModalMode = 'create' | 'edit' | 'view' | null;

export default function CustomerManager() {
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

    return (
        <div className="space-y-4">
            {/* Header con botón crear */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900">Clientes</h1>
                    <p className="text-sm text-gray-500 mt-0.5">
                        Gestiona los clientes de tu negocio
                    </p>
                </div>
                <Button variant="primary" onClick={openCreate}>
                    <PlusIcon className="w-4 h-4 mr-2" />
                    Nuevo cliente
                </Button>
            </div>

            {/* Lista */}
            <CustomerList
                onEdit={openEdit}
                onView={openView}
                onRefreshRef={handleRefreshRef}
            />

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
                            <CustomerDetailView customerId={selectedCustomer.id} />
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
