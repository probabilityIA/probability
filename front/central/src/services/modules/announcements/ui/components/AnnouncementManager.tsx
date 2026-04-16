'use client';

import { useState, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { AnnouncementInfo } from '../../domain/types';
import AnnouncementList from './AnnouncementList';
import AnnouncementForm from './AnnouncementForm';
import { SuperAdminBusinessSelector } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

type ModalMode = 'create' | 'edit' | null;

interface AnnouncementManagerProps {
    selectedBusinessId?: number | null;
    onBusinessChange?: (businessId: number | null) => void;
}

export default function AnnouncementManager({ selectedBusinessId = null, onBusinessChange }: AnnouncementManagerProps) {
    const { isSuperAdmin } = usePermissions();
    const [modalMode, setModalMode] = useState<ModalMode>(null);
    const [selectedAnnouncement, setSelectedAnnouncement] = useState<AnnouncementInfo | null>(null);
    const [refreshList, setRefreshList] = useState<(() => void) | null>(null);

    const openCreate = () => {
        setSelectedAnnouncement(null);
        setModalMode('create');
    };

    const openEdit = (announcement: AnnouncementInfo) => {
        setSelectedAnnouncement(announcement);
        setModalMode('edit');
    };

    const closeModal = () => {
        setModalMode(null);
        setSelectedAnnouncement(null);
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
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Anuncios</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        Gestiona avisos informativos para los negocios
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    {isSuperAdmin && (
                        <SuperAdminBusinessSelector
                            value={selectedBusinessId ?? null}
                            onChange={onBusinessChange || (() => {})}
                            variant="default"
                            placeholder="-- Todos los negocios --"
                        />
                    )}
                    {isSuperAdmin && (
                        <button
                            onClick={openCreate}
                            className="inline-flex items-center justify-center px-6 py-3 font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-all duration-300 hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-offset-2"
                        >
                            <PlusIcon className="w-4 h-4 mr-2" />
                            Nuevo anuncio
                        </button>
                    )}
                </div>
            </div>

            <AnnouncementList
                onEdit={openEdit}
                onRefreshRef={handleRefreshRef}
                selectedBusinessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
            />

            {(modalMode === 'create' || modalMode === 'edit') && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b">
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                                {modalMode === 'create' ? 'Nuevo anuncio' : 'Editar anuncio'}
                            </h2>
                            <button
                                onClick={closeModal}
                                className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl leading-none"
                            >
                                x
                            </button>
                        </div>
                        <div className="p-6">
                            <AnnouncementForm
                                announcement={selectedAnnouncement ?? undefined}
                                onSuccess={handleFormSuccess}
                                onCancel={closeModal}
                            />
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
