'use client';

import { useState, useEffect, useCallback } from 'react';
import { PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import { getDriversAction, deleteDriverAction } from '../../infra/actions';
import { DriverInfo, GetDriversParams } from '../../domain/types';
import { Alert, Table, Spinner } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface DriverListProps {
    onEdit?: (driver: DriverInfo) => void;
    onRefreshRef?: (ref: () => void) => void;
    selectedBusinessId?: number;
}

const statusConfig: Record<string, { label: string; className: string }> = {
    active: { label: 'Activo', className: 'bg-green-100 text-green-800' },
    inactive: { label: 'Inactivo', className: 'bg-gray-100 text-gray-800 dark:text-gray-100' },
    on_route: { label: 'En ruta', className: 'bg-blue-100 text-blue-800' },
};

export default function DriverList({ onEdit, onRefreshRef, selectedBusinessId }: DriverListProps) {
    const [drivers, setDrivers] = useState<DriverInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');

    const fetchDrivers = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetDriversParams = { page, page_size: pageSize };
            if (search) params.search = search;
            if (selectedBusinessId) params.business_id = selectedBusinessId;

            const response = await getDriversAction(params);
            setDrivers(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
            setPage(response.page || page);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar conductores'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, search, selectedBusinessId]);

    useEffect(() => {
        fetchDrivers();
    }, [fetchDrivers]);

    useEffect(() => {
        onRefreshRef?.(fetchDrivers);
    }, [fetchDrivers, onRefreshRef]);

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

    const handleDelete = async (driver: DriverInfo) => {
        if (!confirm(`¿Eliminar al conductor "${driver.first_name} ${driver.last_name}"? Esta accion no se puede deshacer.`)) return;
        try {
            await deleteDriverAction(driver.id, selectedBusinessId);
            fetchDrivers();
        } catch (err: any) {
            setError(getActionError(err, 'Error al eliminar el conductor'));
        }
    };

    const columns = [
        { key: 'name', label: 'Nombre' },
        { key: 'identification', label: 'Identificacion' },
        { key: 'phone', label: 'Telefono' },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'license_type', label: 'Licencia' },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (driver: DriverInfo) => {
        const status = statusConfig[driver.status] || { label: driver.status, className: 'bg-gray-100 text-gray-800 dark:text-gray-100' };

        return {
            name: (
                <span className="font-medium text-gray-900 dark:text-white">
                    {driver.first_name} {driver.last_name}
                </span>
            ),
            identification: (
                <span className="text-sm text-gray-600 dark:text-gray-300">
                    {driver.identification || <span className="text-gray-300">&mdash;</span>}
                </span>
            ),
            phone: (
                <span className="text-sm text-gray-600 dark:text-gray-300">
                    {driver.phone || <span className="text-gray-300">&mdash;</span>}
                </span>
            ),
            status: (
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${status.className}`}>
                    {status.label}
                </span>
            ),
            license_type: (
                <span className="text-sm text-gray-600 dark:text-gray-300">
                    {driver.license_type || <span className="text-gray-300">&mdash;</span>}
                </span>
            ),
            actions: (
                <div className="flex justify-end gap-2">
                    {onEdit && (
                        <button
                            onClick={() => onEdit(driver)}
                            className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors"
                            title="Editar"
                        >
                            <PencilIcon className="w-4 h-4" />
                        </button>
                    )}
                    <button
                        onClick={() => handleDelete(driver)}
                        className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors"
                        title="Eliminar"
                    >
                        <TrashIcon className="w-4 h-4" />
                    </button>
                </div>
            ),
        };
    };

    if (loading && drivers.length === 0) {
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
                    placeholder="Buscar por nombre, identificacion o telefono..."
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
                        className="px-4 py-2 bg-gray-100 text-gray-700 dark:text-gray-200 rounded-lg text-sm hover:bg-gray-200 transition-colors"
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
                    data={drivers.map(renderRow)}
                    keyExtractor={(_, index) => String(drivers[index]?.id || index)}
                    emptyMessage="No hay conductores registrados"
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
