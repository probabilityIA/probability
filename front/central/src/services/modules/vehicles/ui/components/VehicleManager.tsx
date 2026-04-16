'use client';

import { useState, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { VehicleInfo } from '../../domain/types';
import VehicleList from './VehicleList';
import VehicleForm from './VehicleForm';
import { Button } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

type ModalMode = 'create' | 'edit' | null;

interface VehicleManagerProps {
    selectedBusinessId?: number | null;
}

export default function VehicleManager({ selectedBusinessId = null }: VehicleManagerProps) {
    const { isSuperAdmin } = usePermissions();
    const [modalMode, setModalMode] = useState<ModalMode>(null);
    const [selectedVehicle, setSelectedVehicle] = useState<VehicleInfo | null>(null);
    const [refreshList, setRefreshList] = useState<(() => void) | null>(null);

    const openCreate = () => {
        setSelectedVehicle(null);
        setModalMode('create');
    };

    const openEdit = (vehicle: VehicleInfo) => {
        setSelectedVehicle(vehicle);
        setModalMode('edit');
    };

    const closeModal = () => {
        setModalMode(null);
        setSelectedVehicle(null);
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
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Vehiculos</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        Gestiona los vehiculos de tu negocio
                    </p>
                </div>
                {!requiresBusinessSelection && (
                    <Button variant="purple" onClick={openCreate}>
                        <PlusIcon className="w-4 h-4 mr-2" />
                        Nuevo vehiculo
                    </Button>
                )}
            </div>

            {/* Gate: super admin debe seleccionar negocio */}
            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8.25 18.75a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m3 0h6m-9 0H3.375a1.125 1.125 0 01-1.125-1.125V14.25m17.25 4.5a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m3 0h1.125c.621 0 1.125-.504 1.125-1.125v-3.026a2.999 2.999 0 00-.879-2.121l-3.121-3.121A3 3 0 0014.379 8H14V5.625c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v12.25" />
                    </svg>
                    <p className="text-gray-500 dark:text-gray-400 text-sm">Selecciona un negocio para ver y gestionar sus vehiculos</p>
                </div>
            ) : (
                <VehicleList
                    onEdit={openEdit}
                    onRefreshRef={handleRefreshRef}
                    selectedBusinessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
                />
            )}

            {/* Modal crear / editar */}
            {(modalMode === 'create' || modalMode === 'edit') && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b">
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                                {modalMode === 'create' ? 'Nuevo vehiculo' : 'Editar vehiculo'}
                            </h2>
                            <button
                                onClick={closeModal}
                                className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl leading-none"
                            >
                                &times;
                            </button>
                        </div>
                        <div className="p-6">
                            <VehicleForm
                                vehicle={selectedVehicle ?? undefined}
                                onSuccess={handleFormSuccess}
                                onCancel={closeModal}
                                businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
                            />
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
