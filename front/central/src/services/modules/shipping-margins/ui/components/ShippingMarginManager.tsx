'use client';

import { useState, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { ShippingMargin } from '../../domain/types';
import ShippingMarginList from './ShippingMarginList';
import ShippingMarginForm from './ShippingMarginForm';
import { Button } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

type ModalMode = 'create' | 'edit' | null;

interface Props {
    selectedBusinessId?: number | null;
}

export default function ShippingMarginManager({ selectedBusinessId = null }: Props) {
    const { isSuperAdmin } = usePermissions();
    const [modalMode, setModalMode] = useState<ModalMode>(null);
    const [selected, setSelected] = useState<ShippingMargin | null>(null);
    const [refreshList, setRefreshList] = useState<(() => void) | null>(null);

    const openCreate = () => {
        setSelected(null);
        setModalMode('create');
    };

    const openEdit = (m: ShippingMargin) => {
        setSelected(m);
        setModalMode('edit');
    };

    const closeModal = () => {
        setModalMode(null);
        setSelected(null);
    };

    const handleSuccess = () => {
        closeModal();
        refreshList?.();
    };

    const handleRefreshRef = useCallback((ref: () => void) => {
        setRefreshList(() => ref);
    }, []);

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Margenes de envio</h2>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        Configura el margen comercial por transportadora aplicado a cada guia
                    </p>
                </div>
                {!requiresBusinessSelection && (
                    <Button variant="purple" onClick={openCreate}>
                        <PlusIcon className="w-4 h-4 mr-2" />
                        Nuevo margen
                    </Button>
                )}
            </div>

            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                    </svg>
                    <p className="text-gray-500 dark:text-gray-400 text-sm">
                        Selecciona un negocio para ver y configurar sus margenes de envio
                    </p>
                </div>
            ) : (
                <ShippingMarginList
                    onEdit={openEdit}
                    onRefreshRef={handleRefreshRef}
                    selectedBusinessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
                />
            )}

            {(modalMode === 'create' || modalMode === 'edit') && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                                {modalMode === 'create' ? 'Nuevo margen' : `Editar margen - ${selected?.carrier_name}`}
                            </h3>
                            <button
                                onClick={closeModal}
                                className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl leading-none"
                            >
                                &times;
                            </button>
                        </div>
                        <div className="p-6">
                            <ShippingMarginForm
                                margin={selected ?? undefined}
                                onSuccess={handleSuccess}
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
