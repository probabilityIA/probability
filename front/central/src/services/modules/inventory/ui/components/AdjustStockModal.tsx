'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { format } from 'date-fns';
import { AdjustStockDTO } from '../../domain/types';
import { adjustStockAction, getProductInventoryAction } from '../../infra/actions';
import { getProductsAction } from '@/services/modules/products/infra/actions';
import { getLocationsAction, getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { WarehouseLocation, Warehouse } from '@/services/modules/warehouses/domain/types';
import { listLotsAction, listProductUoMsAction, listInventoryStatesAction } from '../../infra/actions/traceability';
import { InventoryLot, ProductUoM, InventoryState } from '../../domain/traceability-types';
import { Button, Alert, Input } from '@/shared/ui';
import { Package, Search, X, Plus, Minus, ChevronDown, Info } from 'lucide-react';

interface ProductOption {
    id: string;
    name: string;
    sku: string;
    variant_label?: string;
}

interface AdjustStockModalProps {
    warehouseId?: number;
    businessId?: number;
    productId?: string;
    onSuccess: () => void;
    onClose: () => void;
}

export default function AdjustStockModal({ warehouseId, businessId, productId, onSuccess, onClose }: AdjustStockModalProps) {
    const [selectedWarehouseId, setSelectedWarehouseId] = useState<number>(warehouseId ?? 0);
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [selectedProduct, setSelectedProduct] = useState<ProductOption | null>(null);
    const [searchName, setSearchName] = useState('');
    const [searchSku, setSearchSku] = useState('');
    const [searchResults, setSearchResults] = useState<ProductOption[]>([]);
    const [activeField, setActiveField] = useState<'name' | 'sku' | null>(null);
    const [searchLoading, setSearchLoading] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    const [quantity, setQuantity] = useState(0);
    const [isAdding, setIsAdding] = useState(true);
    const [reason, setReason] = useState('');
    const [notes, setNotes] = useState('');

    const [lots, setLots] = useState<InventoryLot[]>([]);
    const [selectedLotId, setSelectedLotId] = useState<number | null>(null);
    const [uoms, setUoms] = useState<ProductUoM[]>([]);
    const [selectedUomId, setSelectedUomId] = useState<number | null>(null);
    const [states, setStates] = useState<InventoryState[]>([]);
    const [selectedStateId, setSelectedStateId] = useState<number | null>(null);
    const [locations, setLocations] = useState<WarehouseLocation[]>([]);
    const [selectedLocationId, setSelectedLocationId] = useState<number | null>(null);
    const [loadingContext, setLoadingContext] = useState(false);

    const [productInWarehouse, setProductInWarehouse] = useState<boolean | null>(null);

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    useEffect(() => {
        if (productId) {
            setSelectedProduct({ id: productId, name: '', sku: '' });
        }
    }, [productId]);

    useEffect(() => {
        (async () => {
            try {
                const r = await listInventoryStatesAction();
                setStates(r.data || r || []);
                const avail = (r.data || r || []).find((s: InventoryState) => s.code === 'available');
                if (avail) setSelectedStateId(avail.id);
            } catch {}
        })();
    }, []);

    useEffect(() => {
        (async () => {
            try {
                const r = await getWarehousesAction({ page: 1, page_size: 100, is_active: true, business_id: businessId });
                const warehouseList = r.data || [];
                setWarehouses(warehouseList);
                if (warehouseList.length === 1 && !warehouseId) {
                    setSelectedWarehouseId(warehouseList[0].id);
                }
            } catch { setWarehouses([]); }
        })();
    }, [businessId]);

    useEffect(() => {
        setSelectedLocationId(null);
        (async () => {
            try {
                const r = await getLocationsAction(selectedWarehouseId, businessId);
                setLocations(r || []);
            } catch { setLocations([]); }
        })();
    }, [selectedWarehouseId, businessId]);

    useEffect(() => {
        if (!selectedProduct?.id) {
            setLots([]); setUoms([]); setSelectedLotId(null); setSelectedUomId(null);
            return;
        }
        (async () => {
            setLoadingContext(true);
            try {
                const [lotsRes, uomsRes] = await Promise.all([
                    listLotsAction({ product_id: selectedProduct.id, status: 'active', page: 1, page_size: 100, business_id: businessId }),
                    listProductUoMsAction(selectedProduct.id, businessId),
                ]);
                setLots(lotsRes.data || []);
                const pu: ProductUoM[] = uomsRes.data || uomsRes || [];
                setUoms(pu);
                const base = pu.find((u) => u.is_base);
                if (base) setSelectedUomId(base.id);
            } catch {
                setLots([]); setUoms([]);
            } finally {
                setLoadingContext(false);
            }
        })();
    }, [selectedProduct?.id, businessId]);

    useEffect(() => {
        if (!selectedProduct?.id) { setProductInWarehouse(null); return; }
        setProductInWarehouse(null);
        (async () => {
            try {
                const levels = await getProductInventoryAction(selectedProduct.id, businessId);
                const exists = Array.isArray(levels) && levels.some((l) => l.warehouse_id === selectedWarehouseId);
                setProductInWarehouse(exists);
            } catch { setProductInWarehouse(null); }
        })();
    }, [selectedProduct?.id, selectedWarehouseId, businessId]);

    const searchProducts = useCallback(async (params: { name?: string; sku?: string }) => {
        const term = params.name || params.sku || '';
        if (term.length < 2) { setSearchResults([]); return; }
        setSearchLoading(true);
        try {
            const response = await getProductsAction({ business_id: businessId, ...params, page: 1, page_size: 10 });
            if (response.success && response.data) {
                setSearchResults(response.data.map((p) => ({ id: p.id, name: p.name, sku: p.sku, variant_label: p.variant_label })));
            }
        } catch {} finally { setSearchLoading(false); }
    }, [businessId]);

    useEffect(() => {
        if (!searchName) { setSearchResults([]); return; }
        const timer = setTimeout(() => { if (searchName.length >= 2) searchProducts({ name: searchName }); }, 400);
        return () => clearTimeout(timer);
    }, [searchName, searchProducts]);

    useEffect(() => {
        if (!searchSku) { setSearchResults([]); return; }
        const timer = setTimeout(() => { if (searchSku.length >= 2) searchProducts({ sku: searchSku }); }, 400);
        return () => clearTimeout(timer);
    }, [searchSku, searchProducts]);

    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) setActiveField(null);
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleSelectProduct = (product: ProductOption) => {
        setSelectedProduct(product);
        setSearchName(''); setSearchSku(''); setActiveField(null); setSearchResults([]);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedProduct) return;
        setLoading(true); setError(null); setSuccess(null);

        const finalQuantity = isAdding ? Math.abs(quantity) : -Math.abs(quantity);
        const dto: AdjustStockDTO = {
            product_id: selectedProduct.id,
            warehouse_id: selectedWarehouseId,
            location_id: selectedLocationId,
            lot_id: selectedLotId,
            state_id: selectedStateId,
            uom_id: selectedUomId,
            quantity: finalQuantity,
            reason: reason.trim(),
            notes: notes.trim() || undefined,
        };
        const result = await adjustStockAction(dto, businessId);
        if (!result.success) setError(result.error);
        else {
            const loc = locations.find((l) => l.id === selectedLocationId);
            const lot = lots.find((l) => l.id === selectedLotId);
            const parts: string[] = [`${quantity > 0 ? '+' : ''}${quantity} uds`];
            if (loc) parts.push(`ubicación ${loc.code}`);
            else parts.push('stock general');
            if (lot) parts.push(`lote ${lot.lot_code}`);
            setSuccess(`Ajuste aplicado: ${parts.join(' · ')}`);
            setTimeout(() => onSuccess(), 1200);
        }
        setLoading(false);
    };

    const tracksLots = lots.length > 0;
    const hasMultipleUoms = uoms.length > 1;
    const reasonPresets = ['Conteo físico', 'Corrección', 'Merma', 'Devolución', 'Daño', 'Otro'];
    const isQuantityValid = quantity > 0;
    const isFormValid = selectedProduct && isQuantityValid && reason.trim().length > 0;

    const handleReasonChip = (preset: string) => {
        setReason(preset);
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/45 p-4" onClick={onClose}>
            <div className="bg-white dark:bg-slate-900 rounded-[18px] shadow-[0_24px_80px_-20px_rgba(15,23,42,0.45)] w-full max-w-[720px] max-h-[85vh] overflow-hidden flex flex-col" onClick={(e) => e.stopPropagation()}>
                <div className="flex items-start justify-between px-7 py-[22px]" style={{ borderBottomColor: 'var(--color-primary)' }}>
                    <div className="flex items-start gap-4">
                        <div className="w-[34px] h-[34px] rounded-[10px] flex items-center justify-center flex-shrink-0" style={{ backgroundColor: 'color-mix(in oklab, var(--color-primary) 10%, white)' }}>
                            <Package className="w-5 h-5" style={{ color: 'var(--color-primary)' }} />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold text-slate-900 dark:text-white" style={{ fontFamily: "'Plus Jakarta Sans', sans-serif" }}>Ajustar stock</h2>
                            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">Registra un movimiento de inventario con su justificación.</p>
                        </div>
                    </div>
                    <button
                        onClick={onClose}
                        className="p-1.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                    >
                        <X className="w-5 h-5" />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="overflow-y-auto flex-1 px-7 py-6 space-y-0">
                    {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}
                    {success && <Alert type="success" onClose={() => setSuccess(null)}>{success}</Alert>}

                    {warehouses.length > 1 && (
                        <div className="mb-6 pb-4">
                            <label htmlFor="warehouse" className="block text-sm font-semibold text-slate-700 dark:text-slate-200 mb-2">
                                Bodega <span className="text-rose-500">*</span>
                            </label>
                            <select
                                id="warehouse"
                                value={selectedWarehouseId}
                                onChange={(e) => setSelectedWarehouseId(Number(e.target.value))}
                                className="w-full px-3 py-3 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10"
                            >
                                {warehouses.map((w) => (
                                    <option key={w.id} value={w.id}>{w.name} ({w.code})</option>
                                ))}
                            </select>
                        </div>
                    )}

                    <div className="mb-6">
                        <div className="flex items-center gap-2 mb-5 pb-3 border-b border-slate-100 dark:border-slate-800">
                            <label className="text-xs font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500">
                                Producto <span className="text-rose-500">*</span>
                            </label>
                            <div className="flex-1 border-t border-slate-100 dark:border-slate-800"></div>
                        </div>

                        <div ref={dropdownRef} className="relative">
                            {selectedProduct && !productId ? (
                                <div className="flex items-center justify-between px-4 py-3 bg-slate-100 dark:bg-slate-800/20 border-[1.5px] border-slate-300 dark:border-slate-300 rounded-[10px]">
                                    <div className="flex-1">
                                        <span className="text-sm font-medium text-slate-900 dark:text-white">{selectedProduct.name}</span>
                                        <div className="flex items-center gap-2 mt-1">
                                            {selectedProduct.sku && <span className="text-xs text-slate-500 dark:text-slate-400">SKU: {selectedProduct.sku}</span>}
                                            {selectedProduct.variant_label && <span className="text-xs text-slate-500 dark:text-slate-400">• {selectedProduct.variant_label}</span>}
                                        </div>
                                    </div>
                                    <button type="button" onClick={() => setSelectedProduct(null)} className="text-slate-400 hover:text-slate-600 ml-2">
                                        <X className="w-4 h-4" />
                                    </button>
                                </div>
                            ) : productId ? (
                                <div className="px-4 py-3 bg-slate-100 dark:bg-slate-800 border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm text-slate-700 dark:text-slate-300 font-mono">{productId}</div>
                            ) : (
                                <>
                                    <div className="grid grid-cols-2 gap-3">
                                        <div>
                                            <label htmlFor="search-name" className="block text-xs font-semibold text-slate-600 dark:text-slate-400 mb-2">Nombre</label>
                                            <div className="relative">
                                                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                                <input
                                                    id="search-name"
                                                    type="text"
                                                    value={searchName}
                                                    onChange={(e) => { setSearchName(e.target.value); setSearchSku(''); setActiveField('name'); }}
                                                    onFocus={(e) => { if (searchName.length >= 2) setActiveField('name'); e.currentTarget.style.borderColor = 'var(--color-primary)'; e.currentTarget.style.boxShadow = `0 0 0 4px color-mix(in oklab, var(--color-primary) 10%, transparent)`; }}
                                                    onBlur={(e) => { e.currentTarget.style.borderColor = ''; e.currentTarget.style.boxShadow = ''; }}
                                                    placeholder="Buscar por nombre…"
                                                    className="w-full pl-10 pr-4 py-3 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none"
                                                />
                                            </div>
                                        </div>
                                        <div>
                                            <label htmlFor="search-sku" className="block text-xs font-semibold text-slate-600 dark:text-slate-400 mb-2">SKU</label>
                                            <div className="relative">
                                                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                                <input
                                                    id="search-sku"
                                                    type="text"
                                                    value={searchSku}
                                                    onChange={(e) => { setSearchSku(e.target.value); setSearchName(''); setActiveField('sku'); }}
                                                    onFocus={(e) => { if (searchSku.length >= 2) setActiveField('sku'); e.currentTarget.style.borderColor = 'var(--color-primary)'; e.currentTarget.style.boxShadow = `0 0 0 4px color-mix(in oklab, var(--color-primary) 10%, transparent)`; }}
                                                    onBlur={(e) => { e.currentTarget.style.borderColor = ''; e.currentTarget.style.boxShadow = ''; }}
                                                    placeholder="Buscar por SKU…"
                                                    className="w-full pl-10 pr-4 py-3 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none"
                                                />
                                            </div>
                                        </div>
                                    </div>
                                    {activeField && (searchName.length >= 2 || searchSku.length >= 2) && (
                                        <div className="absolute z-20 w-full mt-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-[10px] shadow-lg max-h-48 overflow-auto">
                                            {searchLoading ? (
                                                <div className="px-4 py-3 text-sm text-slate-500 text-center">Buscando...</div>
                                            ) : searchResults.length > 0 ? (
                                                <ul>
                                                    {searchResults.map((p) => (
                                                        <li key={p.id}>
                                                            <button type="button" onClick={() => handleSelectProduct(p)} className="w-full px-4 py-3 text-left hover:bg-slate-100 dark:hover:bg-teal-900/10 flex items-center justify-between gap-3">
                                                                <div className="flex-1 min-w-0">
                                                                    <span className="text-sm font-medium text-slate-900 dark:text-white truncate block">{p.name}</span>
                                                                    {p.variant_label && <span className="text-xs text-slate-500 dark:text-slate-400 truncate block">{p.variant_label}</span>}
                                                                </div>
                                                                <span className="text-xs text-slate-500 dark:text-slate-400 shrink-0">{p.sku}</span>
                                                            </button>
                                                        </li>
                                                    ))}
                                                </ul>
                                            ) : (
                                                <div className="px-4 py-3 text-sm text-slate-500 text-center">No se encontraron productos</div>
                                            )}
                                        </div>
                                    )}
                                </>
                            )}
                        </div>
                    </div>

                    {selectedProduct && productInWarehouse === false && (
                        <div className="mb-6 flex gap-3 p-4 bg-slate-100 dark:bg-slate-800/20 border border-slate-200 dark:border-slate-800 rounded-[10px]">
                            <Info className="w-5 h-5 text-slate-600 dark:text-slate-400 flex-shrink-0 mt-0.5" />
                            <div>
                                <p className="text-sm font-medium text-teal-800 dark:text-slate-200">Producto nuevo en esta bodega</p>
                                <p className="text-xs text-slate-700 dark:text-slate-300 mt-1">
                                    Este producto aún no tiene stock registrado aquí. El ajuste creará el registro e asignará el producto a esta bodega.
                                </p>
                            </div>
                        </div>
                    )}

                    {loadingContext && <p className="text-xs text-slate-400 animate-pulse mb-4">Cargando contexto del producto...</p>}

                    <div className="mb-6">
                        <div className="flex items-center gap-2 mb-5 pb-3 border-b border-slate-100 dark:border-slate-800">
                            <label className="text-xs font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500">
                                Movimiento <span className="text-rose-500">*</span>
                            </label>
                            <div className="flex-1 border-t border-slate-100 dark:border-slate-800"></div>
                        </div>

                        <div className="grid gap-3" style={{ gridTemplateColumns: '1.1fr 1fr' }}>
                            <div>
                                <label className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-3">Tipo</label>
                                <div className="flex gap-1 p-1 bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-[12px]">
                                    <button
                                        type="button"
                                        onClick={() => setIsAdding(true)}
                                        aria-pressed={isAdding}
                                        role="radio"
                                        className={`flex-1 flex items-center justify-center gap-1.5 py-2.5 px-3 rounded-[10px] text-sm font-semibold transition-all ${
                                            isAdding
                                                ? 'bg-white dark:bg-slate-700 shadow-sm'
                                                : 'bg-transparent text-slate-500 dark:text-slate-400'
                                        }`}
                                    >
                                        <Plus className="w-4 h-4" />
                                        Agregar
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => setIsAdding(false)}
                                        aria-pressed={!isAdding}
                                        role="radio"
                                        className={`flex-1 flex items-center justify-center gap-1.5 py-2.5 px-3 rounded-[10px] text-sm font-semibold transition-all ${
                                            !isAdding
                                                ? 'bg-white dark:bg-slate-700 shadow-sm'
                                                : 'bg-transparent text-slate-500 dark:text-slate-400'
                                        }`}
                                    >
                                        <Minus className="w-4 h-4" />
                                        Restar
                                    </button>
                                </div>
                            </div>

                            {states.length > 0 && (
                                <div>
                                    <label htmlFor="inventory-state" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-3">Estado del inventario</label>
                                    <div className="relative">
                                        <select
                                            id="inventory-state"
                                            value={selectedStateId ?? ''}
                                            onChange={(e) => setSelectedStateId(e.target.value ? Number(e.target.value) : null)}
                                            className="w-full px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10 appearance-none"
                                        >
                                            {states.map((s) => (
                                                <option key={s.id} value={s.id}>{s.name}</option>
                                            ))}
                                        </select>
                                        <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                    </div>
                                </div>
                            )}
                        </div>

                        <div className="mt-4">
                            <label className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-3">Cantidad <span className="text-rose-500">*</span></label>
                            <div className="flex gap-2 items-center">
                                <button
                                    type="button"
                                    onClick={() => setQuantity(Math.max(0, quantity - 1))}
                                    aria-label="Disminuir cantidad"
                                    className="w-12 h-12 flex items-center justify-center bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 rounded-[10px] text-slate-600 dark:text-slate-400 transition-colors"
                                >
                                    <Minus className="w-4 h-4" />
                                </button>
                                <div className="flex-1">
                                    <input
                                        type="number"
                                        value={quantity}
                                        onChange={(e) => setQuantity(Math.max(0, parseInt(e.target.value) || 0))}
                                        className="w-full px-4 py-3 text-center bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-2xl font-bold focus:outline-none"
                                        style={{
                                            fontFamily: "'JetBrains Mono', monospace",
                                            color: isAdding ? 'var(--color-primary)' : 'var(--color-secondary, #e11d48)',
                                            borderColor: 'var(--color-primary)'
                                        }}
                                        min="0"
                                        onFocus={(e) => {
                                            e.currentTarget.style.boxShadow = `0 0 0 4px color-mix(in oklab, var(--color-primary) 10%, transparent)`;
                                        }}
                                        onBlur={(e) => {
                                            e.currentTarget.style.boxShadow = '';
                                        }}
                                    />
                                </div>
                                <button
                                    type="button"
                                    onClick={() => setQuantity(quantity + 1)}
                                    aria-label="Aumentar cantidad"
                                    className="w-12 h-12 flex items-center justify-center bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 rounded-[10px] text-slate-600 dark:text-slate-400 transition-colors"
                                >
                                    <Plus className="w-4 h-4" />
                                </button>
                            </div>
                        </div>

                        {hasMultipleUoms && (
                            <div className="mt-4">
                                <label htmlFor="uom" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-2">Unidad</label>
                                <div className="relative">
                                    <select
                                        id="uom"
                                        value={selectedUomId ?? ''}
                                        onChange={(e) => setSelectedUomId(e.target.value ? Number(e.target.value) : null)}
                                        className="w-full px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10 appearance-none"
                                    >
                                        {uoms.map((u) => (
                                            <option key={u.id} value={u.id}>{u.uom_code} {u.is_base ? '(base)' : `x${u.conversion_factor}`}</option>
                                        ))}
                                    </select>
                                    <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                </div>
                            </div>
                        )}
                    </div>

                    {(locations.length > 0 || tracksLots) && (
                        <div className="mb-6">
                            <div className="flex items-center gap-2 mb-4 pb-3 border-b border-slate-100 dark:border-slate-800">
                                <label className="text-xs font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500">
                                    Contexto
                                </label>
                                <div className="flex-1 border-t border-slate-100 dark:border-slate-800"></div>
                            </div>
                            <div className="grid gap-3" style={{ gridTemplateColumns: '1fr 1fr' }}>
                                {locations.length > 0 && (
                                    <div>
                                        <label htmlFor="location" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-2">
                                            Ubicación <span className="text-slate-400">({locations.length})</span>
                                        </label>
                                        <div className="relative">
                                            <select
                                                id="location"
                                                value={selectedLocationId ?? ''}
                                                onChange={(e) => setSelectedLocationId(e.target.value ? Number(e.target.value) : null)}
                                                className="w-full px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10 appearance-none"
                                            >
                                                <option value="">Stock general de bodega</option>
                                                {locations.map((l) => (
                                                    <option key={l.id} value={l.id}>{l.code} - {l.name}</option>
                                                ))}
                                            </select>
                                            <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                        </div>
                                    </div>
                                )}
                                {tracksLots && (
                                    <div>
                                        <label htmlFor="lot" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-2">
                                            Lote {!isAdding && <span className="text-rose-500">*</span>}
                                        </label>
                                        <div className="relative">
                                            <select
                                                id="lot"
                                                value={selectedLotId ?? ''}
                                                onChange={(e) => setSelectedLotId(e.target.value ? Number(e.target.value) : null)}
                                                className="w-full px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10 appearance-none"
                                            >
                                                <option value="">Sin lote</option>
                                                {lots.map((l) => (
                                                    <option key={l.id} value={l.id}>
                                                        {l.lot_code}{l.expiration_date ? ` - vence ${format(new Date(l.expiration_date), 'dd/MM/yyyy')}` : ''}
                                                    </option>
                                                ))}
                                            </select>
                                            <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                        </div>
                                        {!isAdding && !selectedLotId && (
                                            <p className="text-xs text-amber-600 dark:text-amber-400 mt-1">Para retirar, elige el lote (FEFO).</p>
                                        )}
                                    </div>
                                )}
                            </div>
                        </div>
                    )}

                    <div className="mb-6">
                        <div className="flex items-center gap-2 mb-5 pb-3 border-b border-slate-100 dark:border-slate-800">
                            <label className="text-xs font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500">
                                Detalles
                            </label>
                            <div className="flex-1 border-t border-slate-100 dark:border-slate-800"></div>
                        </div>

                        <div className="grid gap-3" style={{ gridTemplateColumns: '1fr 1fr' }}>
                            <div>
                                <label className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-3">Razón <span className="text-rose-500">*</span></label>
                                <div className="flex flex-wrap gap-2 mb-3">
                                    {reasonPresets.map((preset) => (
                                        <button
                                            key={preset}
                                            type="button"
                                            onClick={() => handleReasonChip(preset)}
                                            className={`px-3 py-2 text-xs font-medium rounded-full transition-all border-[1.5px] ${
                                                reason === preset
                                                    ? 'border-slate-300 dark:border-slate-300 bg-slate-100 dark:bg-slate-800/30 text-slate-700 dark:text-slate-300'
                                                    : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-400 hover:border-slate-300 dark:hover:border-slate-600'
                                            }`}
                                        >
                                            {preset}
                                        </button>
                                    ))}
                                </div>
                                <input
                                    type="text"
                                    value={reason}
                                    onChange={(e) => setReason(e.target.value)}
                                    placeholder="Detalle de la razón…"
                                    className="w-full px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10"
                                />
                            </div>
                            <div>
                                <label htmlFor="notes" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-3">Notas</label>
                                <textarea
                                    id="notes"
                                    value={notes}
                                    onChange={(e) => setNotes(e.target.value)}
                                    placeholder="Notas adicionales…"
                                    className="w-full h-[96px] px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10 resize-none"
                                />
                            </div>
                        </div>
                    </div>
                </form>

                <div className="px-7 py-[18px] bg-slate-50 dark:bg-slate-800/50 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between">
                    <div className="flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                        <Info className="w-4 h-4 flex-shrink-0" />
                        <span>Este movimiento queda registrado en el historial.</span>
                    </div>
                    <div className="flex gap-2">
                        <button
                            type="button"
                            onClick={onClose}
                            disabled={loading}
                            className="px-5 py-3 border-[1.5px] border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-200 rounded-[10px] text-sm font-semibold hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors disabled:opacity-50"
                        >
                            Cancelar
                        </button>
                        <button
                            type="submit"
                            onClick={handleSubmit}
                            disabled={loading || !isFormValid}
                            className="px-5 py-3 text-white rounded-[10px] text-sm font-semibold shadow-sm hover:shadow-md transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                            style={{ backgroundColor: 'var(--color-primary)' }}
                            onMouseEnter={(e) => e.currentTarget.style.opacity = '0.9'}
                            onMouseLeave={(e) => e.currentTarget.style.opacity = '1'}
                        >
                            {loading ? 'Ajustando...' : 'Ajustar stock'}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
