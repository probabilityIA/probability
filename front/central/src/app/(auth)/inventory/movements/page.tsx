'use client';

import { useState, useEffect, useCallback } from 'react';
import { ArrowsRightLeftIcon, AdjustmentsHorizontalIcon } from '@heroicons/react/24/outline';
import { StockMovementList, AdjustStockModal, TransferStockModal } from '@/services/modules/inventory/ui';
import { Button, Spinner } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';

type ModalType = 'adjust' | 'transfer' | null;

export default function InventoryMovementsPage() {
    const { isSuperAdmin } = usePermissions();
    const { setActionButtons } = useNavbarActions();
    const { selectedBusinessId } = useInventoryBusiness();
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [selectedWarehouseId, setSelectedWarehouseId] = useState<number | null>(null);
    const [loadingWarehouses, setLoadingWarehouses] = useState(false);

    const [modalType, setModalType] = useState<ModalType>(null);
    const [refreshMovements, setRefreshMovements] = useState<(() => void) | null>(null);

    const effectiveBusinessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    useEffect(() => {
        if (isSuperAdmin && selectedBusinessId === null) {
            setWarehouses([]);
            setSelectedWarehouseId(null);
            return;
        }

        const loadWarehouses = async () => {
            setLoadingWarehouses(true);
            try {
                const response = await getWarehousesAction({
                    page: 1,
                    page_size: 100,
                    is_active: true,
                    business_id: effectiveBusinessId,
                });
                const whs = response.data || [];
                setWarehouses(whs);
                if (whs.length > 0) {
                    setSelectedWarehouseId(whs[0].id);
                } else {
                    setSelectedWarehouseId(null);
                }
            } catch {
                setWarehouses([]);
                setSelectedWarehouseId(null);
            } finally {
                setLoadingWarehouses(false);
            }
        };
        loadWarehouses();
    }, [isSuperAdmin, selectedBusinessId, effectiveBusinessId]);

    const handleMovementRefreshRef = useCallback((ref: () => void) => {
        setRefreshMovements(() => ref);
    }, []);

    const handleModalSuccess = () => {
        setModalType(null);
        refreshMovements?.();
    };

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;
    const showActionButtons = selectedWarehouseId && !requiresBusinessSelection;

    // Set action buttons in navbar (business selector is in InventorySubNavbar)
    useEffect(() => {
        if (!showActionButtons) {
            setActionButtons(null);
            return;
        }
        const actionButtons = (
            <>
                <Button variant="outline" onClick={() => setModalType('adjust')}>
                    <AdjustmentsHorizontalIcon className="w-4 h-4 mr-2" />
                    Ajustar stock
                </Button>
                <Button variant="primary" onClick={() => setModalType('transfer')}>
                    <ArrowsRightLeftIcon className="w-4 h-4 mr-2" />
                    Transferir
                </Button>
            </>
        );
        setActionButtons(actionButtons);
        return () => setActionButtons(null);
    }, [setActionButtons, showActionButtons]);

    return (
        <div className="min-h-screen bg-gray-50 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="space-y-4">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900">Movimientos</h1>
                    <p className="text-sm text-gray-500 mt-0.5">
                        Historial de movimientos de inventario
                    </p>
                </div>

                {requiresBusinessSelection ? (
                    <div className="flex flex-col items-center justify-center py-16 text-center">
                        <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                        </svg>
                        <p className="text-gray-500 text-sm">Selecciona un negocio para ver los movimientos</p>
                    </div>
                ) : (
                    <>
                        <div className="bg-white border border-gray-200 rounded-lg p-4">
                            <label className="block text-sm font-medium text-gray-700 mb-2">Bodega</label>
                            {loadingWarehouses ? (
                                <div className="flex items-center gap-2">
                                    <Spinner size="sm" />
                                    <span className="text-sm text-gray-500">Cargando bodegas...</span>
                                </div>
                            ) : warehouses.length === 0 ? (
                                <p className="text-sm text-gray-500">No hay bodegas activas. Crea una en el módulo de Bodegas.</p>
                            ) : (
                                <select
                                    value={selectedWarehouseId?.toString() ?? ''}
                                    onChange={(e) => setSelectedWarehouseId(e.target.value ? Number(e.target.value) : null)}
                                    className="w-full max-w-sm px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                >
                                    {warehouses.map((w) => (
                                        <option key={w.id} value={w.id}>
                                            {w.name} ({w.code})
                                        </option>
                                    ))}
                                </select>
                            )}
                        </div>

                        {selectedWarehouseId && (
                            <StockMovementList
                                warehouseId={selectedWarehouseId}
                                selectedBusinessId={effectiveBusinessId}
                                onRefreshRef={handleMovementRefreshRef}
                            />
                        )}
                    </>
                )}

                {modalType === 'adjust' && selectedWarehouseId && (
                    <AdjustStockModal
                        warehouseId={selectedWarehouseId}
                        businessId={effectiveBusinessId}
                        onSuccess={handleModalSuccess}
                        onClose={() => setModalType(null)}
                    />
                )}

                {modalType === 'transfer' && selectedWarehouseId && (
                    <TransferStockModal
                        fromWarehouseId={selectedWarehouseId}
                        businessId={effectiveBusinessId}
                        onSuccess={handleModalSuccess}
                        onClose={() => setModalType(null)}
                    />
                )}
            </div>
        </div>
    );
}
