'use client';

import { useState } from 'react';
import { BuildingStorefrontIcon, CubeIcon, RectangleGroupIcon } from '@heroicons/react/24/outline';
import { MovementsByProductView, MovementsByWarehouseView, MovementsByFamilyView } from '@/services/modules/inventory/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';

type MovView = 'product' | 'warehouse' | 'family';

export default function InventoryMovementsPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const [view, setView] = useState<MovView>('warehouse');

    const effectiveBusinessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;
    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    const viewBtnCls = (active: boolean) =>
        `flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-lg transition-colors ${active
            ? 'btn-business-primary'
            : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700'}`;

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="space-y-4">
                {requiresBusinessSelection ? (
                    <div className="flex flex-col items-center justify-center py-16 text-center">
                        <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                        </svg>
                        <p className="text-gray-500 dark:text-gray-400 text-sm">Selecciona un negocio para ver los movimientos</p>
                    </div>
                ) : (
                    <>
                        <div className="flex gap-2">
                            <button onClick={() => setView('warehouse')} className={viewBtnCls(view === 'warehouse')}>
                                <BuildingStorefrontIcon className="w-4 h-4" />
                                Por bodega
                            </button>
                            <button onClick={() => setView('product')} className={viewBtnCls(view === 'product')}>
                                <CubeIcon className="w-4 h-4" />
                                Por producto
                            </button>
                            <button onClick={() => setView('family')} className={viewBtnCls(view === 'family')}>
                                <RectangleGroupIcon className="w-4 h-4" />
                                Por familia
                            </button>
                        </div>

                        {view === 'warehouse' && (
                            <MovementsByWarehouseView businessId={effectiveBusinessId} />
                        )}
                        {view === 'product' && (
                            <MovementsByProductView businessId={effectiveBusinessId} />
                        )}
                        {view === 'family' && (
                            <MovementsByFamilyView businessId={effectiveBusinessId} />
                        )}
                    </>
                )}
            </div>
        </div>
    );
}
