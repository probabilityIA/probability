'use client';

import { useState, useEffect, useCallback } from 'react';

import { ArrowsRightLeftIcon, AdjustmentsHorizontalIcon, ArrowUpTrayIcon, BuildingStorefrontIcon, CubeIcon } from '@heroicons/react/24/outline';
import InventoryLevelList from './InventoryLevelList';
import ProductInventoryView from './ProductInventoryView';
import StockMovementList from './StockMovementList';
import AdjustStockModal from './AdjustStockModal';
import TransferStockModal from './TransferStockModal';
import BulkLoadInventoryModal from './BulkLoadInventoryModal';
import { Button, Spinner } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';

type StockView = 'warehouse' | 'product';
type Tab = 'stock' | 'movements';
type ModalType = 'adjust' | 'transfer' | 'bulk-load' | null;

export default function InventoryManager() {
    const { isSuperAdmin } = usePermissions();
    const { businesses, loading: loadingBusinesses } = useBusinessesSimple();

    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [selectedWarehouseId, setSelectedWarehouseId] = useState<number | null>(null);
    const [loadingWarehouses, setLoadingWarehouses] = useState(false);

    const [stockView, setStockView] = useState<StockView>('warehouse');
    const [activeTab, setActiveTab] = useState<Tab>('stock');
    const [modalType, setModalType] = useState<ModalType>(null);
    const [adjustProductId, setAdjustProductId] = useState<string | undefined>(undefined);
    const [adjustWarehouseId, setAdjustWarehouseId] = useState<number | null>(null);
    const [refreshStock, setRefreshStock] = useState<(() => void) | null>(null);
    const [refreshMovements, setRefreshMovements] = useState<(() => void) | null>(null);
    const [refreshProductView, setRefreshProductView] = useState<(() => void) | null>(null);

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
                const response = await getWarehousesAction({ page: 1, page_size: 100, is_active: true, business_id: effectiveBusinessId });
                const whs = response.data || [];
                setWarehouses(whs);
                if (whs.length > 0) setSelectedWarehouseId(whs[0].id);
                else setSelectedWarehouseId(null);
            } catch {
                setWarehouses([]);
                setSelectedWarehouseId(null);
            } finally {
                setLoadingWarehouses(false);
            }
        };
        loadWarehouses();
    }, [isSuperAdmin, selectedBusinessId, effectiveBusinessId]);

    const handleStockRefreshRef = useCallback((ref: () => void) => { setRefreshStock(() => ref); }, []);
    const handleMovementRefreshRef = useCallback((ref: () => void) => { setRefreshMovements(() => ref); }, []);
    const handleProductViewRefreshRef = useCallback((ref: () => void) => { setRefreshProductView(() => ref); }, []);

    const handleModalSuccess = () => {
        setModalType(null);
        setAdjustProductId(undefined);
        setAdjustWarehouseId(null);
        refreshStock?.();
        refreshMovements?.();
        refreshProductView?.();
    };

    const handleAdjustFromProduct = (productId: string, warehouseId: number) => {
        setAdjustProductId(productId);
        setAdjustWarehouseId(warehouseId);
        setModalType('adjust');
    };

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    const canShowActions = stockView === 'warehouse' ? (selectedWarehouseId && !requiresBusinessSelection) : !requiresBusinessSelection;

    const tabCls = (active: boolean) =>
        `py-3 px-1 text-sm font-medium border-b-2 transition-colors ${active
            ? 'border-blue-500 text-blue-600'
            : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:text-gray-200 hover:border-gray-300'}`;

    const viewBtnCls = (active: boolean) =>
        `flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg transition-colors ${active
            ? 'bg-teal-700 text-white'
            : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700'}`;

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Inventario</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        Gestiona el stock y movimientos de inventario
                    </p>
                </div>
                {canShowActions && (
                    <div className="flex gap-2">
                        <Button variant="outline" onClick={() => setModalType('bulk-load')}>
                            <ArrowUpTrayIcon className="w-4 h-4 mr-2" />
                            Cargar inventario
                        </Button>
                        <Button variant="outline" onClick={() => { setAdjustProductId(undefined); setAdjustWarehouseId(null); setModalType('adjust'); }}>
                            <AdjustmentsHorizontalIcon className="w-4 h-4 mr-2" />
                            Ajustar stock
                        </Button>
                        <Button variant="primary" onClick={() => setModalType('transfer')}>
                            <ArrowsRightLeftIcon className="w-4 h-4 mr-2" />
                            Transferir
                        </Button>
                    </div>
                )}
            </div>

            {isSuperAdmin && (
                <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                        Negocio <span className="text-red-500">*</span>
                        <span className="ml-1 text-xs text-gray-500 dark:text-gray-400 font-normal">(requerido para gestionar inventario)</span>
                    </label>
                    {loadingBusinesses ? (
                        <p className="text-sm text-gray-500 dark:text-gray-400">Cargando negocios...</p>
                    ) : (
                        <select
                            value={selectedBusinessId?.toString() ?? ''}
                            onChange={(e) => { const val = e.target.value; setSelectedBusinessId(val ? Number(val) : null); }}
                            className="w-full max-w-sm px-3 py-2 bg-white dark:bg-gray-800 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-transparent"
                        >
                            <option value="">— Selecciona un negocio —</option>
                            {businesses.map((b) => (
                                <option key={b.id} value={b.id}>{b.name} (ID: {b.id})</option>
                            ))}
                        </select>
                    )}
                </div>
            )}

            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                    <p className="text-gray-500 dark:text-gray-400 text-sm">Selecciona un negocio para ver y gestionar su inventario</p>
                </div>
            ) : (
                <>
                    <div className="border-b border-gray-200 dark:border-gray-700">
                        <nav className="flex items-center gap-6" aria-label="Tabs">
                            <button onClick={() => setActiveTab('stock')} className={tabCls(activeTab === 'stock')}>Stock</button>
                            <button onClick={() => setActiveTab('movements')} className={tabCls(activeTab === 'movements')}>Movimientos</button>

                            {activeTab === 'stock' && (
                                <div className="ml-auto flex gap-2 pb-2">
                                    <button onClick={() => setStockView('warehouse')} className={viewBtnCls(stockView === 'warehouse')}>
                                        <BuildingStorefrontIcon className="w-4 h-4" />
                                        Por bodega
                                    </button>
                                    <button onClick={() => setStockView('product')} className={viewBtnCls(stockView === 'product')}>
                                        <CubeIcon className="w-4 h-4" />
                                        Por producto
                                    </button>
                                </div>
                            )}
                        </nav>
                    </div>

                    {activeTab === 'stock' && stockView === 'warehouse' && (
                        <>
                            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Bodega</label>
                                {loadingWarehouses ? (
                                    <div className="flex items-center gap-2">
                                        <Spinner size="sm" />
                                        <span className="text-sm text-gray-500 dark:text-gray-400">Cargando bodegas...</span>
                                    </div>
                                ) : warehouses.length === 0 ? (
                                    <p className="text-sm text-gray-500 dark:text-gray-400">No hay bodegas activas. Crea una en el modulo de Bodegas.</p>
                                ) : (
                                    <select
                                        value={selectedWarehouseId?.toString() ?? ''}
                                        onChange={(e) => setSelectedWarehouseId(e.target.value ? Number(e.target.value) : null)}
                                        className="w-full max-w-sm px-3 py-2 bg-white dark:bg-gray-800 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-transparent"
                                    >
                                        {warehouses.map((w) => (
                                            <option key={w.id} value={w.id}>{w.name} ({w.code})</option>
                                        ))}
                                    </select>
                                )}
                            </div>

                            {selectedWarehouseId && (
                                <InventoryLevelList
                                    warehouseId={selectedWarehouseId}
                                    selectedBusinessId={effectiveBusinessId}
                                    onRefreshRef={handleStockRefreshRef}
                                />
                            )}
                        </>
                    )}

                    {activeTab === 'stock' && stockView === 'product' && (
                        <ProductInventoryView
                            businessId={effectiveBusinessId}
                            onAdjust={handleAdjustFromProduct}
                            onRefreshRef={handleProductViewRefreshRef}
                        />
                    )}

                    {activeTab === 'movements' && (
                        <StockMovementList
                            warehouseId={selectedWarehouseId ?? undefined}
                            selectedBusinessId={effectiveBusinessId}
                            onRefreshRef={handleMovementRefreshRef}
                        />
                    )}
                </>
            )}

            {modalType === 'bulk-load' && selectedWarehouseId && (
                <BulkLoadInventoryModal
                    warehouseId={selectedWarehouseId}
                    businessId={effectiveBusinessId}
                    onSuccess={handleModalSuccess}
                    onClose={() => setModalType(null)}
                />
            )}

            {modalType === 'adjust' && (
                <AdjustStockModal
                    warehouseId={adjustWarehouseId ?? selectedWarehouseId ?? warehouses[0]?.id}
                    businessId={effectiveBusinessId}
                    productId={adjustProductId}
                    onSuccess={handleModalSuccess}
                    onClose={() => { setModalType(null); setAdjustProductId(undefined); setAdjustWarehouseId(null); }}
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
    );
}
