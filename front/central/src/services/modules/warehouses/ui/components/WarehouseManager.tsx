'use client';

import { useState, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { Warehouse } from '../../domain/types';
import WarehouseList from './WarehouseList';
import WarehouseForm from './WarehouseForm';
import { Button } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';

type ModalMode = 'create' | 'edit' | null;

export default function WarehouseManager() {
    const { isSuperAdmin } = usePermissions();
    const { businesses, loading: loadingBusinesses } = useBusinessesSimple();

    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [modalMode, setModalMode] = useState<ModalMode>(null);
    const [selectedWarehouse, setSelectedWarehouse] = useState<Warehouse | null>(null);
    const [refreshList, setRefreshList] = useState<(() => void) | null>(null);

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
        refreshList?.();
    };

    const handleRefreshRef = useCallback((ref: () => void) => {
        setRefreshList(() => ref);
    }, []);

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    return (
        <div className="space-y-4">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900">Bodegas</h1>
                    <p className="text-sm text-gray-500 mt-0.5">
                        Gestiona las bodegas y ubicaciones de tu negocio
                    </p>
                </div>
                {!requiresBusinessSelection && (
                    <Button variant="primary" onClick={openCreate}>
                        <PlusIcon className="w-4 h-4 mr-2" />
                        Nueva bodega
                    </Button>
                )}
            </div>

            {/* Selector de negocio para super admin */}
            {isSuperAdmin && (
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Negocio <span className="text-red-500">*</span>
                        <span className="ml-1 text-xs text-gray-500 font-normal">(requerido para gestionar bodegas)</span>
                    </label>
                    {loadingBusinesses ? (
                        <p className="text-sm text-gray-500">Cargando negocios...</p>
                    ) : (
                        <select
                            value={selectedBusinessId?.toString() ?? ''}
                            onChange={(e) => {
                                const val = e.target.value;
                                setSelectedBusinessId(val ? Number(val) : null);
                            }}
                            className="w-full max-w-sm px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                            <option value="">— Selecciona un negocio —</option>
                            {businesses.map((b) => (
                                <option key={b.id} value={b.id}>
                                    {b.name} (ID: {b.id})
                                </option>
                            ))}
                        </select>
                    )}
                </div>
            )}

            {/* Gate: super admin debe seleccionar negocio */}
            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M2.25 21h19.5m-18-18v18m10.5-18v18m6-13.5V21M6.75 6.75h.75m-.75 3h.75m-.75 3h.75m3-6h.75m-.75 3h.75m-.75 3h.75M6.75 21v-3.375c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21M3 3h12m-.75 4.5H21m-3.75 7.5h.008v.008h-.008v-.008zm0 3h.008v.008h-.008v-.008z" />
                    </svg>
                    <p className="text-gray-500 text-sm">Selecciona un negocio para ver y gestionar sus bodegas</p>
                </div>
            ) : (
                <WarehouseList
                    onEdit={openEdit}
                    onRefreshRef={handleRefreshRef}
                    selectedBusinessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
                />
            )}

            {/* Modal crear / editar */}
            {(modalMode === 'create' || modalMode === 'edit') && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b">
                            <h2 className="text-lg font-semibold text-gray-900">
                                {modalMode === 'create' ? 'Nueva bodega' : 'Editar bodega'}
                            </h2>
                            <button
                                onClick={closeModal}
                                className="text-gray-400 hover:text-gray-600 text-xl leading-none"
                            >
                                &times;
                            </button>
                        </div>
                        <div className="p-6">
                            <WarehouseForm
                                warehouse={selectedWarehouse ?? undefined}
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
