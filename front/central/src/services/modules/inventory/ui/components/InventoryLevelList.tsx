'use client';

import { useState, useEffect, useCallback } from 'react';
import { getWarehouseInventoryAction } from '../../infra/actions';
import { InventoryLevel, GetInventoryParams } from '../../domain/types';
import { Alert, Table, Spinner } from '@/shared/ui';

interface InventoryLevelListProps {
    warehouseId: number;
    selectedBusinessId?: number;
    onRefreshRef?: (ref: () => void) => void;
}

export default function InventoryLevelList({ warehouseId, selectedBusinessId, onRefreshRef }: InventoryLevelListProps) {
    const [levels, setLevels] = useState<InventoryLevel[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');
    const [lowStockOnly, setLowStockOnly] = useState(false);

    const fetchLevels = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetInventoryParams = { page, page_size: pageSize };
            if (search) params.search = search;
            if (lowStockOnly) params.low_stock = true;
            if (selectedBusinessId) params.business_id = selectedBusinessId;

            const response = await getWarehouseInventoryAction(warehouseId, params);
            setLevels(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
            setPage(response.page || page);
        } catch (err: any) {
            setError(err.message || 'Error al cargar inventario');
        } finally {
            setLoading(false);
        }
    }, [warehouseId, page, pageSize, search, lowStockOnly, selectedBusinessId]);

    useEffect(() => {
        fetchLevels();
    }, [fetchLevels]);

    useEffect(() => {
        onRefreshRef?.(fetchLevels);
    }, [fetchLevels, onRefreshRef]);

    useEffect(() => {
        setPage(1);
        setSearch('');
        setSearchInput('');
    }, [warehouseId]);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(searchInput);
        setPage(1);
    };

    const handleClearSearch = () => {
        setSearchInput('');
        setSearch('');
        setPage(1);
    };

    const isLowStock = (level: InventoryLevel) => {
        if (level.reorder_point != null) return level.available_qty <= level.reorder_point;
        if (level.min_stock != null) return level.available_qty <= level.min_stock;
        return false;
    };

    const columns = [
        { key: 'product', label: 'Producto' },
        { key: 'quantity', label: 'Cantidad', align: 'center' as const },
        { key: 'reserved', label: 'Reservado', align: 'center' as const },
        { key: 'available', label: 'Disponible', align: 'center' as const },
        { key: 'limits', label: 'Min / Max', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
    ];

    const renderRow = (level: InventoryLevel) => ({
        product: (
            <div>
                <span className="font-medium text-gray-900">{level.product_name || level.product_id}</span>
                {level.product_sku && (
                    <span className="block text-xs text-gray-500 font-mono">{level.product_sku}</span>
                )}
            </div>
        ),
        quantity: (
            <span className="text-sm font-medium text-gray-900">{level.quantity}</span>
        ),
        reserved: (
            <span className={`text-sm ${level.reserved_qty > 0 ? 'text-orange-600 font-medium' : 'text-gray-400'}`}>
                {level.reserved_qty}
            </span>
        ),
        available: (
            <span className="text-sm font-semibold text-gray-900">{level.available_qty}</span>
        ),
        limits: (
            <span className="text-xs text-gray-500">
                {level.min_stock != null ? level.min_stock : '—'} / {level.max_stock != null ? level.max_stock : '—'}
            </span>
        ),
        status: (
            <div className="flex justify-center">
                {isLowStock(level) ? (
                    <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                        Stock bajo
                    </span>
                ) : (
                    <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        OK
                    </span>
                )}
            </div>
        ),
    });

    if (loading && levels.length === 0) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div className="space-y-4">
            <div className="flex gap-2 items-center flex-wrap">
                <form onSubmit={handleSearch} className="flex gap-2 flex-1 min-w-[200px]">
                    <input
                        type="text"
                        value={searchInput}
                        onChange={(e) => setSearchInput(e.target.value)}
                        placeholder="Buscar por producto o SKU..."
                        className="flex-1 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                    <button
                        type="submit"
                        className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 transition-colors"
                    >
                        Buscar
                    </button>
                    {search && (
                        <button
                            type="button"
                            onClick={handleClearSearch}
                            className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg text-sm hover:bg-gray-200 transition-colors"
                        >
                            Limpiar
                        </button>
                    )}
                </form>
                <label className="flex items-center gap-2 text-sm text-gray-600 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={lowStockOnly}
                        onChange={(e) => { setLowStockOnly(e.target.checked); setPage(1); }}
                        className="w-4 h-4 rounded border-gray-300 text-red-600 focus:ring-red-500"
                    />
                    Solo stock bajo
                </label>
            </div>

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                <Table
                    columns={columns}
                    data={levels.map(renderRow)}
                    keyExtractor={(_, index) => String(levels[index]?.id || index)}
                    emptyMessage="No hay inventario en esta bodega"
                    loading={loading}
                    pagination={{
                        currentPage: page,
                        totalPages: totalPages,
                        totalItems: total,
                        itemsPerPage: pageSize,
                        onPageChange: (newPage) => setPage(newPage),
                        onItemsPerPageChange: (newSize) => {
                            setPageSize(newSize);
                            setPage(1);
                        },
                    }}
                />
            </div>
        </div>
    );
}
