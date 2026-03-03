'use client';

import { useState, useEffect, useCallback } from 'react';
import { PencilIcon, TrashIcon, EyeIcon } from '@heroicons/react/24/outline';
import { getWarehousesAction, deleteWarehouseAction } from '../../infra/actions';
import { Warehouse, GetWarehousesParams } from '../../domain/types';
import { Alert, Table, Spinner } from '@/shared/ui';

interface WarehouseListProps {
    onView?: (warehouse: Warehouse) => void;
    onEdit?: (warehouse: Warehouse) => void;
    onRefreshRef?: (ref: () => void) => void;
    selectedBusinessId?: number;
}

export default function WarehouseList({ onView, onEdit, onRefreshRef, selectedBusinessId }: WarehouseListProps) {
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');

    const fetchWarehouses = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetWarehousesParams = { page, page_size: pageSize };
            if (search) params.search = search;
            if (selectedBusinessId) params.business_id = selectedBusinessId;

            const response = await getWarehousesAction(params);
            setWarehouses(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
            setPage(response.page || page);
        } catch (err: any) {
            setError(err.message || 'Error al cargar bodegas');
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, search, selectedBusinessId]);

    useEffect(() => {
        fetchWarehouses();
    }, [fetchWarehouses]);

    useEffect(() => {
        onRefreshRef?.(fetchWarehouses);
    }, [fetchWarehouses, onRefreshRef]);

    useEffect(() => {
        setPage(1);
        setSearch('');
        setSearchInput('');
    }, [selectedBusinessId]);

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

    const handleDelete = async (warehouse: Warehouse) => {
        if (!confirm(`¿Eliminar la bodega "${warehouse.name}"? Esta acción no se puede deshacer.`)) return;
        try {
            await deleteWarehouseAction(warehouse.id, selectedBusinessId);
            fetchWarehouses();
        } catch (err: any) {
            setError(err.message || 'Error al eliminar la bodega');
        }
    };

    const columns = [
        { key: 'name', label: 'Nombre' },
        { key: 'code', label: 'Código' },
        { key: 'location', label: 'Ubicación' },
        { key: 'contact', label: 'Contacto' },
        { key: 'badges', label: 'Tipo', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (warehouse: Warehouse) => ({
        name: (
            <span className="font-medium text-gray-900">{warehouse.name}</span>
        ),
        code: (
            <span className="text-sm text-gray-600 font-mono">{warehouse.code}</span>
        ),
        location: (
            <div className="text-sm">
                {warehouse.address ? (
                    <p className="text-gray-900 truncate max-w-[200px]" title={warehouse.address}>{warehouse.address}</p>
                ) : (
                    <p className="text-gray-300">&mdash;</p>
                )}
                {(warehouse.city || warehouse.state) && (
                    <p className="text-xs text-gray-500">
                        {[warehouse.city, warehouse.state].filter(Boolean).join(', ')}
                    </p>
                )}
            </div>
        ),
        contact: (
            <div className="text-sm">
                {warehouse.contact_name ? (
                    <p className="text-gray-900 truncate max-w-[180px]">{warehouse.contact_name}</p>
                ) : null}
                {warehouse.phone ? (
                    <p className="text-xs text-gray-500">{warehouse.phone}</p>
                ) : null}
                {!warehouse.contact_name && !warehouse.phone && (
                    <span className="text-gray-300">&mdash;</span>
                )}
            </div>
        ),
        badges: (
            <div className="flex justify-center gap-1">
                {warehouse.is_default && (
                    <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                        Principal
                    </span>
                )}
                {warehouse.is_fulfillment && (
                    <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                        Fulfillment
                    </span>
                )}
            </div>
        ),
        status: (
            <div className="flex justify-center">
                <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${warehouse.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'}`}>
                    {warehouse.is_active ? 'Activa' : 'Inactiva'}
                </span>
            </div>
        ),
        actions: (
            <div className="flex justify-end gap-2">
                {onView && (
                    <button
                        onClick={() => onView(warehouse)}
                        className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors"
                        title="Ver detalle"
                    >
                        <EyeIcon className="w-4 h-4" />
                    </button>
                )}
                {onEdit && (
                    <button
                        onClick={() => onEdit(warehouse)}
                        className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors"
                        title="Editar"
                    >
                        <PencilIcon className="w-4 h-4" />
                    </button>
                )}
                <button
                    onClick={() => handleDelete(warehouse)}
                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors"
                    title="Eliminar"
                >
                    <TrashIcon className="w-4 h-4" />
                </button>
            </div>
        ),
    });

    if (loading && warehouses.length === 0) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div className="space-y-4">
            <form onSubmit={handleSearch} className="flex gap-2">
                <input
                    type="text"
                    value={searchInput}
                    onChange={(e) => setSearchInput(e.target.value)}
                    placeholder="Buscar por nombre o código..."
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

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                <Table
                    columns={columns}
                    data={warehouses.map(renderRow)}
                    keyExtractor={(_, index) => String(warehouses[index]?.id || index)}
                    emptyMessage="No hay bodegas registradas"
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
