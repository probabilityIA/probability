'use client';

import { useState, useEffect, useCallback } from 'react';
import { AdjustmentsHorizontalIcon, BuildingStorefrontIcon, CubeIcon, ChartBarIcon } from '@heroicons/react/24/outline';
import { AdjustStockModal, TransferStockModal, ProductInventoryView, WarehouseInventoryView, InventoryAnalyticsView } from '@/services/modules/inventory/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';

type ModalType = 'adjust' | 'transfer' | null;
type StockView = 'warehouse' | 'product' | 'analytics';

export default function InventoryStockPage() {
    const { isSuperAdmin } = usePermissions();
    const { setActionButtons } = useNavbarActions();
    const { selectedBusinessId } = useInventoryBusiness();

    const [stockView, setStockView] = useState<StockView>('warehouse');
    const [modalType, setModalType] = useState<ModalType>(null);
    const [selectedProductId, setSelectedProductId] = useState<string | undefined>(undefined);
    const [adjustWarehouseId, setAdjustWarehouseId] = useState<number | null>(null);
    const [refreshWarehouseView, setRefreshWarehouseView] = useState<(() => void) | null>(null);
    const [refreshProductView, setRefreshProductView] = useState<(() => void) | null>(null);

    const effectiveBusinessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;
    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    const handleWarehouseViewRefreshRef = useCallback((ref: () => void) => {
        setRefreshWarehouseView(() => ref);
    }, []);

    const handleProductViewRefreshRef = useCallback((ref: () => void) => {
        setRefreshProductView(() => ref);
    }, []);

    const handleAdjustFromWarehouse = useCallback((productId: string, warehouseId: number) => {
        setSelectedProductId(productId);
        setAdjustWarehouseId(warehouseId);
        setModalType('adjust');
    }, []);

    const handleAdjustFromProduct = useCallback((productId: string, warehouseId: number) => {
        setSelectedProductId(productId);
        setAdjustWarehouseId(warehouseId);
        setModalType('adjust');
    }, []);

    const handleModalSuccess = () => {
        setModalType(null);
        setSelectedProductId(undefined);
        setAdjustWarehouseId(null);
        refreshWarehouseView?.();
        refreshProductView?.();
    };

    const canShowActions = !requiresBusinessSelection;

    useEffect(() => {
        if (!canShowActions) {
            setActionButtons(null);
            return;
        }
        const actionButtons = (
            <button
                onClick={() => { setSelectedProductId(undefined); setAdjustWarehouseId(null); setModalType('adjust'); }}
                title="Ajustar stock"
                className="p-2 rounded-md border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
                <AdjustmentsHorizontalIcon className="w-5 h-5" />
            </button>
        );
        setActionButtons(actionButtons);
        return () => setActionButtons(null);
    }, [setActionButtons, canShowActions]);

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
                        <p className="text-gray-500 dark:text-gray-400 text-sm">Selecciona un negocio para ver el inventario</p>
                    </div>
                ) : (
                    <>
                        <div className="flex gap-2">
                            <button onClick={() => setStockView('warehouse')} className={viewBtnCls(stockView === 'warehouse')}>
                                <BuildingStorefrontIcon className="w-4 h-4" />
                                Por bodega
                            </button>
                            <button onClick={() => setStockView('product')} className={viewBtnCls(stockView === 'product')}>
                                <CubeIcon className="w-4 h-4" />
                                Por producto
                            </button>
                            <button onClick={() => setStockView('analytics')} className={viewBtnCls(stockView === 'analytics')}>
                                <ChartBarIcon className="w-4 h-4" />
                                Analitica
                            </button>
                        </div>

                        {stockView === 'warehouse' && (
                            <WarehouseInventoryView
                                businessId={effectiveBusinessId}
                                onAdjust={handleAdjustFromWarehouse}
                                onRefreshRef={handleWarehouseViewRefreshRef}
                            />
                        )}

                        {stockView === 'product' && (
                            <ProductInventoryView
                                businessId={effectiveBusinessId}
                                onAdjust={handleAdjustFromProduct}
                                onRefreshRef={handleProductViewRefreshRef}
                            />
                        )}

                        {stockView === 'analytics' && (
                            <InventoryAnalyticsView businessId={effectiveBusinessId} />
                        )}
                    </>
                )}

                {modalType === 'adjust' && (
                    <AdjustStockModal
                        warehouseId={adjustWarehouseId ?? undefined}
                        businessId={effectiveBusinessId}
                        productId={selectedProductId}
                        onSuccess={handleModalSuccess}
                        onClose={() => { setModalType(null); setSelectedProductId(undefined); setAdjustWarehouseId(null); }}
                    />
                )}

                {modalType === 'transfer' && (
                    <TransferStockModal
                        fromWarehouseId={adjustWarehouseId ?? 0}
                        businessId={effectiveBusinessId}
                        productId={selectedProductId}
                        onSuccess={handleModalSuccess}
                        onClose={() => { setModalType(null); setSelectedProductId(undefined); }}
                    />
                )}
            </div>
        </div>
    );
}
