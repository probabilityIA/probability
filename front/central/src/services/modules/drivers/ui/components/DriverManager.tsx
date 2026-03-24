'use client';

import { useState, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { DriverInfo } from '../../domain/types';
import DriverList from './DriverList';
import DriverForm from './DriverForm';
import { Button } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

type ModalMode = 'create' | 'edit' | null;

interface DriverManagerProps {
    selectedBusinessId?: number | null;
}

export default function DriverManager({ selectedBusinessId = null }: DriverManagerProps) {
    const { isSuperAdmin } = usePermissions();
    const [modalMode, setModalMode] = useState<ModalMode>(null);
    const [selectedDriver, setSelectedDriver] = useState<DriverInfo | null>(null);
    const [refreshList, setRefreshList] = useState<(() => void) | null>(null);

    const openCreate = () => {
        setSelectedDriver(null);
        setModalMode('create');
    };

    const openEdit = (driver: DriverInfo) => {
        setSelectedDriver(driver);
        setModalMode('edit');
    };

    const closeModal = () => {
        setModalMode(null);
        setSelectedDriver(null);
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
                    <h1 className="text-xl font-semibold text-gray-900">Conductores</h1>
                    <p className="text-sm text-gray-500 mt-0.5">
                        Gestiona los conductores de tu negocio
                    </p>
                </div>
                {!requiresBusinessSelection && (
                    <Button variant="primary" onClick={openCreate}>
                        <PlusIcon className="w-4 h-4 mr-2" />
                        Nuevo conductor
                    </Button>
                )}
            </div>

            {/* Gate: super admin debe seleccionar negocio */}
            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                    </svg>
                    <p className="text-gray-500 text-sm">Selecciona un negocio para ver y gestionar sus conductores</p>
                </div>
            ) : (
                <DriverList
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
                            <h2 className="text-lg font-semibold text-gray-900">
                                {modalMode === 'create' ? 'Nuevo conductor' : 'Editar conductor'}
                            </h2>
                            <button
                                onClick={closeModal}
                                className="text-gray-400 hover:text-gray-600 text-xl leading-none"
                            >
                                &times;
                            </button>
                        </div>
                        <div className="p-6">
                            <DriverForm
                                driver={selectedDriver ?? undefined}
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
