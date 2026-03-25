'use client';

import { useState, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { RouteInfo } from '../../domain/types';
import RouteList from './RouteList';
import RouteForm from './RouteForm';
import RouteDetail from './RouteDetail';
import { Button } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

type ViewMode = 'list' | 'create' | 'detail';

interface RouteManagerProps {
    selectedBusinessId?: number | null;
}

export default function RouteManager({ selectedBusinessId = null }: RouteManagerProps) {
    const { isSuperAdmin } = usePermissions();
    const [viewMode, setViewMode] = useState<ViewMode>('list');
    const [selectedRouteId, setSelectedRouteId] = useState<number | null>(null);
    const [editingRoute, setEditingRoute] = useState<RouteInfo | null>(null);
    const [refreshList, setRefreshList] = useState<(() => void) | null>(null);

    const effectiveBusinessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    const openCreate = () => {
        setEditingRoute(null);
        setViewMode('create');
    };

    const openDetail = (route: RouteInfo) => {
        setSelectedRouteId(route.id);
        setViewMode('detail');
    };

    const openEdit = (route: RouteInfo) => {
        setEditingRoute(route);
        setViewMode('create');
    };

    const backToList = () => {
        setViewMode('list');
        setSelectedRouteId(null);
        setEditingRoute(null);
    };

    const handleFormSuccess = () => {
        backToList();
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
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Rutas de Entrega</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        Gestiona las rutas de entrega de ultima milla
                    </p>
                </div>
                {!requiresBusinessSelection && viewMode === 'list' && (
                    <Button variant="purple" onClick={openCreate}>
                        <PlusIcon className="w-4 h-4 mr-2" />
                        Nueva ruta
                    </Button>
                )}
                {viewMode !== 'list' && (
                    <Button variant="outline-purple" onClick={backToList}>
                        Volver a la lista
                    </Button>
                )}
            </div>

            {/* Gate: super admin debe seleccionar negocio */}
            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 6.75V15m6-6v8.25m.503 3.498l4.875-2.437c.381-.19.622-.58.622-1.006V4.82c0-.836-.88-1.38-1.628-1.006l-3.869 1.934c-.317.159-.69.159-1.006 0L9.503 3.252a1.125 1.125 0 00-1.006 0L3.622 5.689C3.24 5.88 3 6.27 3 6.695V19.18c0 .836.88 1.38 1.628 1.006l3.869-1.934c.317-.159.69-.159 1.006 0l4.994 2.497c.317.158.69.158 1.006 0z" />
                    </svg>
                    <p className="text-gray-500 dark:text-gray-400 text-sm">Selecciona un negocio para ver y gestionar sus rutas</p>
                </div>
            ) : (
                <>
                    {/* List view */}
                    {viewMode === 'list' && (
                        <RouteList
                            onView={openDetail}
                            onEdit={openEdit}
                            onRefreshRef={handleRefreshRef}
                            selectedBusinessId={effectiveBusinessId}
                        />
                    )}

                    {/* Create / Edit form modal */}
                    {viewMode === 'create' && (
                        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-3xl max-h-[90vh] overflow-y-auto">
                                <div className="flex items-center justify-between px-6 py-4 border-b">
                                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                                        {editingRoute ? 'Editar ruta' : 'Nueva ruta'}
                                    </h2>
                                    <button
                                        onClick={backToList}
                                        className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl leading-none"
                                    >
                                        &times;
                                    </button>
                                </div>
                                <div className="p-6">
                                    <RouteForm
                                        route={editingRoute ?? undefined}
                                        onSuccess={handleFormSuccess}
                                        onCancel={backToList}
                                        businessId={effectiveBusinessId}
                                    />
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Detail view */}
                    {viewMode === 'detail' && selectedRouteId && (
                        <RouteDetail
                            routeId={selectedRouteId}
                            businessId={effectiveBusinessId}
                            onBack={backToList}
                            onRefreshList={() => refreshList?.()}
                        />
                    )}
                </>
            )}
        </div>
    );
}
