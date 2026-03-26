'use client';

import { useState, useEffect, useCallback } from 'react';
import { PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import { getVehiclesAction, deleteVehicleAction } from '../../infra/actions';
import { VehicleInfo, GetVehiclesParams } from '../../domain/types';
import { Alert, Table, Spinner } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface VehicleListProps {
    onEdit?: (vehicle: VehicleInfo) => void;
    onRefreshRef?: (ref: () => void) => void;
    selectedBusinessId?: number;
}

const VEHICLE_TYPE_ICONS: Record<string, string> = {
    motorcycle: '\uD83C\uDFCD\uFE0F',
    car: '\uD83D\uDE97',
    van: '\uD83D\uDE90',
    truck: '\uD83D\uDE9B',
};

const STATUS_STYLES: Record<string, string> = {
    active: 'bg-green-100 text-green-800',
    inactive: 'bg-gray-100 text-gray-800 dark:text-gray-100',
    in_maintenance: 'bg-yellow-100 text-yellow-800',
};

const STATUS_LABELS: Record<string, string> = {
    active: 'Activo',
    inactive: 'Inactivo',
    in_maintenance: 'En mantenimiento',
};

export default function VehicleList({ onEdit, onRefreshRef, selectedBusinessId }: VehicleListProps) {
    const [vehicles, setVehicles] = useState<VehicleInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');

    const fetchVehicles = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetVehiclesParams = { page, page_size: pageSize };
            if (search) params.search = search;
            if (selectedBusinessId) params.business_id = selectedBusinessId;

            const response = await getVehiclesAction(params);
            setVehicles(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
            setPage(response.page || page);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar vehiculos'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, search, selectedBusinessId]);

    useEffect(() => {
        fetchVehicles();
    }, [fetchVehicles]);

    useEffect(() => {
        onRefreshRef?.(fetchVehicles);
    }, [fetchVehicles, onRefreshRef]);

    // Resetear a pagina 1 cuando cambia el negocio seleccionado
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

    const handleDelete = async (vehicle: VehicleInfo) => {
        if (!confirm(`Eliminar el vehiculo "${vehicle.license_plate}"? Esta accion no se puede deshacer.`)) return;
        try {
            await deleteVehicleAction(vehicle.id, selectedBusinessId);
            fetchVehicles();
        } catch (err: any) {
            setError(getActionError(err, 'Error al eliminar el vehiculo'));
        }
    };

    const columns = [
        { key: 'type', label: 'Tipo' },
        { key: 'license_plate', label: 'Placa' },
        { key: 'brand_model', label: 'Marca / Modelo' },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (vehicle: VehicleInfo) => ({
        type: (
            <span className="text-sm text-gray-900 dark:text-white">
                <span className="mr-1">{VEHICLE_TYPE_ICONS[vehicle.type] || vehicle.type}</span>
                {vehicle.type}
            </span>
        ),
        license_plate: (
            <span className="font-medium text-gray-900 dark:text-white">{vehicle.license_plate}</span>
        ),
        brand_model: (
            <span className="text-sm text-gray-600 dark:text-gray-300">
                {vehicle.brand || vehicle.model
                    ? `${vehicle.brand || ''}${vehicle.brand && vehicle.model ? ' ' : ''}${vehicle.model || ''}`
                    : <span className="text-gray-300">&mdash;</span>
                }
            </span>
        ),
        status: (
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${STATUS_STYLES[vehicle.status] || 'bg-gray-100 text-gray-800 dark:text-gray-100'}`}>
                {STATUS_LABELS[vehicle.status] || vehicle.status}
            </span>
        ),
        actions: (
            <div className="flex justify-end gap-2">
                {onEdit && (
                    <button
                        onClick={() => onEdit(vehicle)}
                        className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors"
                        title="Editar"
                    >
                        <PencilIcon className="w-4 h-4" />
                    </button>
                )}
                <button
                    onClick={() => handleDelete(vehicle)}
                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors"
                    title="Eliminar"
                >
                    <TrashIcon className="w-4 h-4" />
                </button>
            </div>
        ),
    });

    if (loading && vehicles.length === 0) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {/* Buscador */}
            <form onSubmit={handleSearch} className="flex gap-2">
                <input
                    type="text"
                    value={searchInput}
                    onChange={(e) => setSearchInput(e.target.value)}
                    placeholder="Buscar por placa, marca o modelo..."
                    className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
                />
                <button
                    type="submit"
                    className="px-4 py-2 bg-purple-600 dark:bg-purple-600 text-white rounded-lg text-sm hover:bg-purple-700 dark:hover:bg-purple-700 transition-colors"
                >
                    Buscar
                </button>
                {search && (
                    <button
                        type="button"
                        onClick={handleClearSearch}
                        className="px-4 py-2 bg-purple-600 dark:bg-purple-600 text-white rounded-lg text-sm hover:bg-purple-700 dark:hover:bg-purple-700 transition-colors"
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

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                <Table
                    columns={columns}
                    data={vehicles.map(renderRow)}
                    keyExtractor={(_, index) => String(vehicles[index]?.id || index)}
                    emptyMessage="No hay vehiculos registrados"
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
