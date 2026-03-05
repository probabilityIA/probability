'use client';

import { useState, useEffect, useCallback } from 'react';
import { EyeIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import { getRoutesAction, deleteRouteAction } from '../../infra/actions';
import { RouteInfo, GetRoutesParams } from '../../domain/types';
import { Alert, Table, Spinner } from '@/shared/ui';

interface RouteListProps {
    onView?: (route: RouteInfo) => void;
    onEdit?: (route: RouteInfo) => void;
    onRefreshRef?: (ref: () => void) => void;
    selectedBusinessId?: number;
}

const STATUS_LABELS: Record<string, string> = {
    planned: 'Planificada',
    in_progress: 'En progreso',
    completed: 'Completada',
    cancelled: 'Cancelada',
};

const STATUS_COLORS: Record<string, string> = {
    planned: 'bg-gray-100 text-gray-700',
    in_progress: 'bg-blue-100 text-blue-700',
    completed: 'bg-green-100 text-green-700',
    cancelled: 'bg-red-100 text-red-700',
};

export default function RouteList({ onView, onEdit, onRefreshRef, selectedBusinessId }: RouteListProps) {
    const [routes, setRoutes] = useState<RouteInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');
    const [statusFilter, setStatusFilter] = useState('');

    const fetchRoutes = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetRoutesParams = { page, page_size: pageSize };
            if (search) params.search = search;
            if (statusFilter) params.status = statusFilter;
            if (selectedBusinessId) params.business_id = selectedBusinessId;

            const response = await getRoutesAction(params);
            setRoutes(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
            setPage(response.page || page);
        } catch (err: any) {
            setError(err.message || 'Error al cargar rutas');
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, search, statusFilter, selectedBusinessId]);

    useEffect(() => {
        fetchRoutes();
    }, [fetchRoutes]);

    useEffect(() => {
        onRefreshRef?.(fetchRoutes);
    }, [fetchRoutes, onRefreshRef]);

    // Resetear a pagina 1 cuando cambia el negocio seleccionado
    useEffect(() => {
        setPage(1);
        setSearch('');
        setSearchInput('');
        setStatusFilter('');
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

    const handleStatusFilterChange = (value: string) => {
        setStatusFilter(value);
        setPage(1);
    };

    const handleDelete = async (route: RouteInfo) => {
        if (!confirm(`Eliminar la ruta del ${formatDate(route.date)}? Esta accion no se puede deshacer.`)) return;
        try {
            await deleteRouteAction(route.id, selectedBusinessId);
            fetchRoutes();
        } catch (err: any) {
            setError(err.message || 'Error al eliminar la ruta');
        }
    };

    const formatDate = (dateStr: string) => {
        const date = new Date(dateStr);
        return date.toLocaleDateString('es-CO', { day: '2-digit', month: '2-digit', year: 'numeric' });
    };

    const columns = [
        { key: 'date', label: 'Fecha' },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'driver_name', label: 'Conductor' },
        { key: 'vehicle_plate', label: 'Vehiculo' },
        { key: 'progress', label: 'Progreso', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (route: RouteInfo) => ({
        date: (
            <span className="font-medium text-gray-900">{formatDate(route.date)}</span>
        ),
        status: (
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${STATUS_COLORS[route.status] || 'bg-gray-100 text-gray-700'}`}>
                {STATUS_LABELS[route.status] || route.status}
            </span>
        ),
        driver_name: (
            <span className="text-sm text-gray-600">{route.driver_name || <span className="text-gray-300">&mdash;</span>}</span>
        ),
        vehicle_plate: (
            <span className="text-sm text-gray-600">{route.vehicle_plate || <span className="text-gray-300">&mdash;</span>}</span>
        ),
        progress: (
            <div className="flex items-center gap-2">
                <div className="w-20 bg-gray-200 rounded-full h-2">
                    <div
                        className="bg-green-500 h-2 rounded-full transition-all"
                        style={{ width: route.total_stops > 0 ? `${(route.completed_stops / route.total_stops) * 100}%` : '0%' }}
                    />
                </div>
                <span className="text-xs text-gray-500 whitespace-nowrap">
                    {route.completed_stops}/{route.total_stops}
                </span>
            </div>
        ),
        actions: (
            <div className="flex justify-end gap-2">
                {onView && (
                    <button
                        onClick={() => onView(route)}
                        className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors"
                        title="Ver detalle"
                    >
                        <EyeIcon className="w-4 h-4" />
                    </button>
                )}
                {onEdit && (
                    <button
                        onClick={() => onEdit(route)}
                        className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors"
                        title="Editar"
                    >
                        <PencilIcon className="w-4 h-4" />
                    </button>
                )}
                <button
                    onClick={() => handleDelete(route)}
                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors"
                    title="Eliminar"
                >
                    <TrashIcon className="w-4 h-4" />
                </button>
            </div>
        ),
    });

    if (loading && routes.length === 0) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {/* Filters row */}
            <div className="flex flex-col sm:flex-row gap-2">
                <form onSubmit={handleSearch} className="flex gap-2 flex-1">
                    <input
                        type="text"
                        value={searchInput}
                        onChange={(e) => setSearchInput(e.target.value)}
                        placeholder="Buscar por conductor, vehiculo..."
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
                <select
                    value={statusFilter}
                    onChange={(e) => handleStatusFilterChange(e.target.value)}
                    className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                >
                    <option value="">Todos los estados</option>
                    <option value="planned">Planificada</option>
                    <option value="in_progress">En progreso</option>
                    <option value="completed">Completada</option>
                    <option value="cancelled">Cancelada</option>
                </select>
            </div>

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                <Table
                    columns={columns}
                    data={routes.map(renderRow)}
                    keyExtractor={(_, index) => String(routes[index]?.id || index)}
                    emptyMessage="No hay rutas registradas"
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
