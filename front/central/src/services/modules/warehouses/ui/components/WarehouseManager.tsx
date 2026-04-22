'use client';

import { useState, useEffect } from 'react';
import { PlusIcon, AcademicCapIcon } from '@heroicons/react/24/outline';
import { Warehouse } from '../../domain/types';
import WarehouseForm from './WarehouseForm';
import WarehouseTreeTable from './WarehouseTreeTable';
import WarehouseTour from './WarehouseTour';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';

type ModalMode = 'create' | 'edit' | null;

export default function WarehouseManager() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const [modalMode, setModalMode] = useState<ModalMode>(null);
    const [selectedWarehouse, setSelectedWarehouse] = useState<Warehouse | null>(null);
    const [refreshKey, setRefreshKey] = useState(0);
    const [tourOpen, setTourOpen] = useState(false);
    const [pulseTour, setPulseTour] = useState(false);

    useEffect(() => {
        try {
            const seen = localStorage.getItem('warehouse_tour_seen_v1');
            if (!seen) setPulseTour(true);
        } catch {}
    }, []);

    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;
    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    const openCreate = () => {
        setSelectedWarehouse(null);
        setModalMode('create');
    };

    const openEdit = (warehouse: Warehouse) => {
        setSelectedWarehouse(warehouse);
        setModalMode('edit');
    };

    const closeModal = () => {
        setModalMode(null);
        setSelectedWarehouse(null);
    };

    const handleFormSuccess = () => {
        closeModal();
        setRefreshKey((k) => k + 1);
    };

    const isFormModal = modalMode === 'create' || modalMode === 'edit';

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-end">
                {!requiresBusinessSelection && (
                    <div className="flex items-center gap-2">
                        <button
                            onClick={() => { setTourOpen(true); setPulseTour(false); }}
                            className={`p-2 rounded-md transition-all text-white btn-business-primary ${pulseTour ? 'tour-pulse' : ''}`}
                            title={pulseTour ? '¡Nuevo! Tutorial de jerarquía' : 'Tutorial de jerarquía'}
                        >
                            <AcademicCapIcon className="w-5 h-5" />
                        </button>
                        <button
                            onClick={openCreate}
                            className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-110 hover:-translate-y-1 flex items-center gap-2"
                        >
                            <PlusIcon className="w-4 h-4" />
                            Nueva bodega
                        </button>
                    </div>
                )}
            </div>

            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M2.25 21h19.5m-18-18v18m10.5-18v18m6-13.5V21M6.75 6.75h.75m-.75 3h.75m-.75 3h.75m3-6h.75m-.75 3h.75m-.75 3h.75M6.75 21v-3.375c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21M3 3h12m-.75 4.5H21m-3.75 7.5h.008v.008h-.008v-.008zm0 3h.008v.008h-.008v-.008z" />
                    </svg>
                    <p className="text-gray-500 dark:text-gray-400 text-sm">Selecciona un negocio para ver y gestionar sus bodegas</p>
                </div>
            ) : (
                <WarehouseTreeTable
                    businessId={businessId}
                    onEditWarehouse={openEdit}
                    onNewWarehouse={openCreate}
                    refreshKey={refreshKey}
                />
            )}

            <WarehouseTour isOpen={tourOpen} onClose={() => setTourOpen(false)} />

            {isFormModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                                {modalMode === 'create' ? 'Nueva bodega' : 'Editar bodega'}
                            </h2>
                            <button onClick={closeModal} className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl leading-none">
                                &times;
                            </button>
                        </div>
                        <div className="p-6">
                            <WarehouseForm
                                warehouse={selectedWarehouse ?? undefined}
                                onSuccess={handleFormSuccess}
                                onCancel={closeModal}
                                businessId={businessId}
                            />
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
