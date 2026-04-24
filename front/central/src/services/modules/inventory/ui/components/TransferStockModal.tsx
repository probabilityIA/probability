'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { format } from 'date-fns';
import { TransferStockDTO } from '../../domain/types';
import { transferStockAction } from '../../infra/actions';
import { listLotsAction, listProductUoMsAction, listInventoryStatesAction } from '../../infra/actions/traceability';
import { InventoryLot, ProductUoM, InventoryState } from '../../domain/traceability-types';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { getProductsAction } from '@/services/modules/products/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';
import { Button, Alert, Input } from '@/shared/ui';

interface ProductOption {
    id: string;
    name: string;
    sku: string;
}

interface TransferStockModalProps {
    fromWarehouseId: number;
    businessId?: number;
    productId?: string;
    onSuccess: () => void;
    onClose: () => void;
}

export default function TransferStockModal({ fromWarehouseId, businessId, productId, onSuccess, onClose }: TransferStockModalProps) {
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [selectedProduct, setSelectedProduct] = useState<ProductOption | null>(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [searchResults, setSearchResults] = useState<ProductOption[]>([]);
    const [showDropdown, setShowDropdown] = useState(false);
    const [searchLoading, setSearchLoading] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    const [toWarehouseId, setToWarehouseId] = useState(0);
    const [quantity, setQuantity] = useState(1);
    const [reason, setReason] = useState('');
    const [notes, setNotes] = useState('');

    const [lots, setLots] = useState<InventoryLot[]>([]);
    const [selectedLotId, setSelectedLotId] = useState<number | null>(null);
    const [uoms, setUoms] = useState<ProductUoM[]>([]);
    const [selectedUomId, setSelectedUomId] = useState<number | null>(null);
    const [states, setStates] = useState<InventoryState[]>([]);
    const [selectedStateId, setSelectedStateId] = useState<number | null>(null);
    const [loadingContext, setLoadingContext] = useState(false);

    const [loading, setLoading] = useState(false);
    const [loadingWarehouses, setLoadingWarehouses] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    useEffect(() => {
        (async () => {
            try {
                const r = await getWarehousesAction({ page: 1, page_size: 100, is_active: true, business_id: businessId });
                setWarehouses((r.data || []).filter((w: Warehouse) => w.id !== fromWarehouseId));
            } catch { setError('Error al cargar bodegas'); }
            finally { setLoadingWarehouses(false); }
        })();
    }, [businessId, fromWarehouseId]);

    useEffect(() => {
        if (productId) setSelectedProduct({ id: productId, name: '', sku: '' });
    }, [productId]);

    useEffect(() => {
        (async () => {
            try {
                const r = await listInventoryStatesAction();
                const list: InventoryState[] = r.data || r || [];
                setStates(list);
                const avail = list.find((s) => s.code === 'available');
                if (avail) setSelectedStateId(avail.id);
            } catch {}
        })();
    }, []);

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
            } catch { setLots([]); setUoms([]); }
            finally { setLoadingContext(false); }
        })();
    }, [selectedProduct?.id, businessId]);

    const searchProducts = useCallback(async (term: string) => {
        if (term.length < 2) { setSearchResults([]); return; }
        setSearchLoading(true);
        try {
            const response = await getProductsAction({ business_id: businessId, name: term, page: 1, page_size: 10 });
            if (response.success && response.data) {
                setSearchResults(response.data.map((p) => ({ id: p.id, name: p.name, sku: p.sku })));
            }
        } catch {} finally { setSearchLoading(false); }
    }, [businessId]);

    useEffect(() => {
        const timer = setTimeout(() => { if (searchTerm.length >= 2) searchProducts(searchTerm); }, 400);
        return () => clearTimeout(timer);
    }, [searchTerm, searchProducts]);

    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) setShowDropdown(false);
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleSelectProduct = (product: ProductOption) => {
        setSelectedProduct(product);
        setSearchTerm(''); setShowDropdown(false); setSearchResults([]);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedProduct) return;
        setLoading(true); setError(null); setSuccess(null);

        const dto: TransferStockDTO = {
            product_id: selectedProduct.id,
            from_warehouse_id: fromWarehouseId,
            to_warehouse_id: toWarehouseId,
            lot_id: selectedLotId,
            state_id: selectedStateId,
            uom_id: selectedUomId,
            quantity,
            reason: reason.trim() || undefined,
            notes: notes.trim() || undefined,
        };
        const result = await transferStockAction(dto, businessId);
        if (!result.success) setError(result.error);
        else {
            setSuccess('Transferencia realizada exitosamente');
            setTimeout(() => onSuccess(), 800);
        }
        setLoading(false);
    };

    const tracksLots = lots.length > 0;
    const hasMultipleUoms = uoms.length > 1;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Transferir stock</h2>
                    <button onClick={onClose} className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 text-xl leading-none">&times;</button>
                </div>
                <form onSubmit={handleSubmit} className="p-6 space-y-4">
                    {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}
                    {success && <Alert type="success" onClose={() => setSuccess(null)}>{success}</Alert>}

                    <div ref={dropdownRef} className="relative">
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Producto <span className="text-red-500">*</span>
                        </label>
                        {selectedProduct && !productId ? (
                            <div className="flex items-center justify-between px-3 py-2 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
                                <div>
                                    <span className="text-sm font-medium text-gray-900 dark:text-white">{selectedProduct.name}</span>
                                    {selectedProduct.sku && <span className="ml-2 text-xs text-gray-500 dark:text-gray-400">SKU: {selectedProduct.sku}</span>}
                                </div>
                                <button type="button" onClick={() => setSelectedProduct(null)} className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 ml-2">&times;</button>
                            </div>
                        ) : productId ? (
                            <div className="px-3 py-2 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-700 dark:text-gray-300 font-mono">{productId}</div>
                        ) : (
                            <>
                                <Input
                                    type="text"
                                    value={searchTerm}
                                    onChange={(e) => { setSearchTerm(e.target.value); setShowDropdown(true); }}
                                    onFocus={() => { if (searchTerm.length >= 2) setShowDropdown(true); }}
                                    placeholder="Buscar por nombre o SKU..."
                                />
                                {showDropdown && searchTerm.length >= 2 && (
                                    <div className="absolute z-20 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-xl max-h-48 overflow-auto">
                                        {searchLoading ? (
                                            <div className="px-4 py-3 text-sm text-gray-500 text-center">Buscando...</div>
                                        ) : searchResults.length > 0 ? (
                                            <ul>
                                                {searchResults.map((p) => (
                                                    <li key={p.id}>
                                                        <button type="button" onClick={() => handleSelectProduct(p)} className="w-full px-4 py-2.5 text-left hover:bg-blue-50 dark:hover:bg-blue-900/20 flex items-center justify-between">
                                                            <span className="text-sm font-medium text-gray-900 dark:text-white truncate">{p.name}</span>
                                                            <span className="text-xs text-gray-500 dark:text-gray-400 ml-2 shrink-0">{p.sku}</span>
                                                        </button>
                                                    </li>
                                                ))}
                                            </ul>
                                        ) : (
                                            <div className="px-4 py-3 text-sm text-gray-500 text-center">No se encontraron productos</div>
                                        )}
                                    </div>
                                )}
                            </>
                        )}
                    </div>

                    {loadingContext && <p className="text-xs text-gray-400 animate-pulse">Cargando contexto del producto...</p>}

                    {tracksLots && (
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                                Lote a transferir <span className="text-red-500">*</span>
                                <span className="ml-2 text-xs text-gray-400 font-normal">Producto con lotes</span>
                            </label>
                            <select
                                value={selectedLotId ?? ''}
                                onChange={(e) => setSelectedLotId(e.target.value ? Number(e.target.value) : null)}
                                className="w-full px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-lg text-sm"
                                required
                            >
                                <option value="">— Elige un lote —</option>
                                {lots.map((l) => (
                                    <option key={l.id} value={l.id}>
                                        {l.lot_code}
                                        {l.expiration_date ? ` · vence ${format(new Date(l.expiration_date), 'dd/MM/yyyy')}` : ''}
                                    </option>
                                ))}
                            </select>
                            <p className="text-xs text-gray-500 mt-1">FEFO: prioriza el lote más próximo a vencer.</p>
                        </div>
                    )}

                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Bodega destino <span className="text-red-500">*</span>
                        </label>
                        {loadingWarehouses ? (
                            <p className="text-sm text-gray-500 dark:text-gray-400">Cargando bodegas...</p>
                        ) : (
                            <select
                                value={toWarehouseId || ''}
                                onChange={(e) => setToWarehouseId(Number(e.target.value))}
                                className="w-full px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                                required
                            >
                                <option value="">-- Selecciona bodega destino --</option>
                                {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name} ({w.code})</option>)}
                            </select>
                        )}
                    </div>

                    <div className="grid grid-cols-2 gap-3">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                                Cantidad <span className="text-red-500">*</span>
                            </label>
                            <Input type="number" value={quantity.toString()} onChange={(e) => setQuantity(parseInt(e.target.value) || 0)} min="1" required />
                        </div>
                        {hasMultipleUoms && (
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Unidad</label>
                                <select
                                    value={selectedUomId ?? ''}
                                    onChange={(e) => setSelectedUomId(e.target.value ? Number(e.target.value) : null)}
                                    className="w-full px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-lg text-sm"
                                >
                                    {uoms.map((u) => (
                                        <option key={u.id} value={u.id}>
                                            {u.uom_code} {u.is_base ? '(base)' : `×${u.conversion_factor}`}
                                        </option>
                                    ))}
                                </select>
                            </div>
                        )}
                    </div>

                    {states.length > 0 && (
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Estado del inventario</label>
                            <select
                                value={selectedStateId ?? ''}
                                onChange={(e) => setSelectedStateId(e.target.value ? Number(e.target.value) : null)}
                                className="w-full px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-lg text-sm"
                            >
                                {states.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
                            </select>
                        </div>
                    )}

                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Razón</label>
                        <Input type="text" value={reason} onChange={(e) => setReason(e.target.value)} placeholder="Reabastecimiento, redistribución..." />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Notas</label>
                        <textarea value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Notas adicionales..." rows={2} className="w-full px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 resize-none" />
                    </div>

                    <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                        <Button type="button" variant="outline" onClick={onClose} disabled={loading}>Cancelar</Button>
                        <Button type="submit" variant="primary" disabled={loading || !selectedProduct || (tracksLots && !selectedLotId)}>
                            {loading ? 'Transfiriendo...' : 'Transferir'}
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
}
