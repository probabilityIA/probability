'use client';

import { useState, useEffect, useCallback, useMemo, useRef, memo } from 'react';
import { getOrdersAction, deleteOrderAction } from '../../infra/actions';
import { Order, GetOrdersParams } from '../../domain/types';
import { Button, Alert, DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { useSSE } from '@/shared/hooks/use-sse';
import { useToast } from '@/shared/providers/toast-provider';
import RawOrderModal from './RawOrderModal';

// Componente memoizado para las filas de la tabla
const OrderRow = memo(({ 
    order, 
    onView, 
    onEdit, 
    onDelete, 
    onShowRaw,
    formatCurrency,
    formatDate,
    getStatusBadge
}: {
    order: Order;
    onView?: (order: Order) => void;
    onEdit?: (order: Order) => void;
    onDelete: (id: string) => void;
    onShowRaw: (id: string) => void;
    formatCurrency: (amount: number, currency?: string) => string;
    formatDate: (dateString: string) => string;
    getStatusBadge: (status: string) => React.ReactNode;
}) => {
    return (
        <tr className="hover:bg-gray-50 transition-colors">
            <td className="px-3 sm:px-6 py-4">
                <div className="text-sm font-medium text-gray-900">
                    {order.order_number}
                </div>
                <div className="text-xs text-gray-500 sm:hidden">
                    {order.customer_name}
                </div>
                <div className="text-xs text-gray-500">
                    {order.internal_number}
                </div>
            </td>
            <td className="px-3 sm:px-6 py-4 hidden sm:table-cell">
                <div className="text-sm text-gray-900">{order.customer_name}</div>
                <div className="text-xs text-gray-500">{order.customer_email}</div>
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                <div className="text-sm font-semibold text-gray-900">
                    {formatCurrency(order.total_amount, order.currency)}
                </div>
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                {getStatusBadge(order.status)}
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden lg:table-cell">
                <span className="text-sm text-gray-900 capitalize">
                    {order.platform}
                </span>
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500 hidden md:table-cell">
                {formatDate(order.created_at)}
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div className="flex flex-col sm:flex-row justify-end gap-1 sm:gap-2">
                    {onView && (
                        <button
                            onClick={() => onView(order)}
                            className="px-2 sm:px-3 py-1 sm:py-1.5 bg-blue-500 hover:bg-blue-600 text-white text-xs sm:text-sm font-medium rounded-md transition-colors duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                        >
                            Ver
                        </button>
                    )}
                    {onEdit && (
                        <button
                            onClick={() => onEdit(order)}
                            className="px-2 sm:px-3 py-1 sm:py-1.5 bg-yellow-500 hover:bg-yellow-600 text-white text-xs sm:text-sm font-medium rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                        >
                            Editar
                        </button>
                    )}
                    <button
                        onClick={() => onShowRaw(order.id)}
                        className="px-2 sm:px-3 py-1 sm:py-1.5 bg-gray-500 hover:bg-gray-600 text-white text-xs sm:text-sm font-medium rounded-md transition-colors duration-200 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
                    >
                        Original
                    </button>
                    <button
                        onClick={() => onDelete(order.id)}
                        className="px-2 sm:px-3 py-1 sm:py-1.5 bg-red-500 hover:bg-red-600 text-white text-xs sm:text-sm font-medium rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                    >
                        Eliminar
                    </button>
                </div>
            </td>
        </tr>
    );
});

OrderRow.displayName = 'OrderRow';

interface OrderListProps {
    onView?: (order: Order) => void;
    onEdit?: (order: Order) => void;
    refreshKey?: number;
}

export default function OrderList({ onView, onEdit, refreshKey }: OrderListProps) {
    const [orders, setOrders] = useState<Order[]>([]);
    const [initialLoading, setInitialLoading] = useState(true);
    const [tableLoading, setTableLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const isFirstLoad = useRef(true);

    // Raw Data Modal
    const [selectedOrderId, setSelectedOrderId] = useState<string | null>(null);
    const [isRawModalOpen, setIsRawModalOpen] = useState(false);

    // Filters
    const [filters, setFilters] = useState<GetOrdersParams>({
        page: 1,
        page_size: 20,
    });

    const { showToast } = useToast();

    // Definir filtros disponibles
    const availableFilters: FilterOption[] = [
        {
            key: 'order_number',
            label: 'ID de orden',
            type: 'text',
            placeholder: 'Buscar por ID de orden...',
        },
        {
            key: 'internal_number',
            label: 'Número interno',
            type: 'text',
            placeholder: 'Buscar por número interno...',
        },
        {
            key: 'status',
            label: 'Estado',
            type: 'select',
            options: [
                { value: 'pending', label: 'Pendiente' },
                { value: 'processing', label: 'Procesando' },
                { value: 'shipped', label: 'Enviado' },
                { value: 'delivered', label: 'Entregado' },
                { value: 'cancelled', label: 'Cancelado' },
            ],
        },
        {
            key: 'platform',
            label: 'Plataforma',
            type: 'select',
            options: [
                { value: 'shopify', label: 'Shopify' },
                { value: 'woocommerce', label: 'WooCommerce' },
                { value: 'manual', label: 'Manual' },
            ],
        },
        {
            key: 'is_paid',
            label: 'Estado de pago',
            type: 'boolean',
        },
        {
            key: 'start_date',
            label: 'Rango de fechas',
            type: 'date-range',
        },
    ];

    // Convertir filtros a ActiveFilter[]
    const activeFilters: ActiveFilter[] = useMemo(() => {
        const active: ActiveFilter[] = [];

        if (filters.order_number) {
            active.push({
                key: 'order_number',
                label: 'ID de orden',
                value: filters.order_number,
                type: 'text',
            });
        }

        if (filters.internal_number) {
            active.push({
                key: 'internal_number',
                label: 'Número interno',
                value: filters.internal_number,
                type: 'text',
            });
        }


        if (filters.status) {
            active.push({
                key: 'status',
                label: 'Estado',
                value: filters.status,
                type: 'select',
            });
        }

        if (filters.platform) {
            active.push({
                key: 'platform',
                label: 'Plataforma',
                value: filters.platform,
                type: 'select',
            });
        }

        if (filters.is_paid !== undefined) {
            active.push({
                key: 'is_paid',
                label: 'Estado de pago',
                value: filters.is_paid,
                type: 'boolean',
            });
        }

        if (filters.start_date || filters.end_date) {
            active.push({
                key: 'start_date',
                label: 'Rango de fechas',
                value: {
                    start: filters.start_date,
                    end: filters.end_date,
                },
                type: 'date-range',
            });
        }

        return active;
    }, [filters]);

    // Manejar adición de filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev, page: 1 };

            if (filterKey === 'start_date' && typeof value === 'object') {
                newFilters.start_date = value.start;
                newFilters.end_date = value.end;
            } else if (filterKey === 'is_paid') {
                newFilters.is_paid = value === true;
            } else {
                (newFilters as any)[filterKey] = value;
            }

            return newFilters;
        });
    }, []);

    // Manejar eliminación de filtro
    const handleRemoveFilter = useCallback((filterKey: string) => {
        setFilters((prev) => {
            const newFilters = { ...prev, page: 1 };

            if (filterKey === 'start_date') {
                delete newFilters.start_date;
                delete newFilters.end_date;
            } else {
                delete (newFilters as any)[filterKey];
            }

            return newFilters;
        });
    }, []);

    // Manejar cambio de ordenamiento
    const handleSortChange = useCallback((sortBy: string, sortOrder: 'asc' | 'desc') => {
        setFilters((prev) => ({
            ...prev,
            sort_by: sortBy as 'created_at' | 'updated_at' | 'total_amount',
            sort_order: sortOrder,
            page: 1,
        }));
    }, []);

    // SSE Integration - Actualizar sin recargar toda la página
    useSSE({
        eventTypes: ['order.created'],
        onMessage: (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.type === 'order.created') {
                    const orderNumber = data.data?.order_number || 'Desconocida';
                    showToast(`Nueva orden recibida: #${orderNumber}`, 'success');
                    // Actualizar solo la tabla sin recargar toda la página
                    refreshTableOnly();
                }
            } catch (e) {
                console.error('Error processing SSE message:', e);
            }
        },
    });

    // Función para actualizar solo la tabla (sin mostrar loading inicial)
    const refreshTableOnly = useCallback(async () => {
        setTableLoading(true);
        try {
            const response = await getOrdersAction(filters);
            if (response.success && response.data) {
                setOrders(response.data);
                setTotal(response.total || 0);
                setTotalPages(response.total_pages || 1);
                setPage(response.page || 1);
            }
        } catch (err: any) {
            console.error('Error al actualizar órdenes:', err);
        } finally {
            setTableLoading(false);
        }
    }, [filters]);

    // Función unificada para cargar órdenes
    const loadOrders = useCallback(async (showInitialLoading = false) => {
        if (showInitialLoading) {
            setInitialLoading(true);
        } else {
            setTableLoading(true);
        }
        setError(null);
        try {
            const response = await getOrdersAction(filters);
            if (response.success && response.data) {
                setOrders(response.data);
                setTotal(response.total || 0);
                setTotalPages(response.total_pages || 1);
                setPage(response.page || 1);
            } else {
                setError(response.message || 'Error al cargar las órdenes');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar las órdenes');
        } finally {
            setInitialLoading(false);
            setTableLoading(false);
        }
    }, [filters]);

    // Carga inicial - solo una vez
    useEffect(() => {
        if (isFirstLoad.current) {
            isFirstLoad.current = false;
            loadOrders(true);
        }
    }, [loadOrders]);

    // Actualizar cuando cambian los filtros (sin loading inicial, solo tabla)
    useEffect(() => {
        if (!isFirstLoad.current) {
            loadOrders(false);
        }
    }, [filters, loadOrders]);

    // Refresh cuando cambia el refreshKey (desde el padre, después de crear/editar)
    useEffect(() => {
        if (refreshKey !== undefined && refreshKey > 0) {
            refreshTableOnly();
        }
    }, [refreshKey, refreshTableOnly]);

    const handleDelete = async (id: string) => {
        if (!confirm('¿Estás seguro de que deseas eliminar esta orden?')) return;

        try {
            const response = await deleteOrderAction(id);
            if (response.success) {
                refreshTableOnly();
            } else {
                alert(response.message || 'Error al eliminar la orden');
            }
        } catch (err: any) {
            alert(err.message || 'Error al eliminar la orden');
        }
    };

    const formatCurrency = useCallback((amount: number, currency: string = 'USD') => {
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency,
        }).format(amount);
    }, []);

    const formatDate = useCallback((dateString: string) => {
        return new Date(dateString).toLocaleDateString('es-CO', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    }, []);

    const getStatusBadge = useCallback((status: string) => {
        const statusColors: Record<string, string> = {
            pending: 'bg-yellow-100 text-yellow-800',
            processing: 'bg-blue-100 text-blue-800',
            shipped: 'bg-purple-100 text-purple-800',
            delivered: 'bg-green-100 text-green-800',
            cancelled: 'bg-red-100 text-red-800',
        };

        const color = statusColors[status.toLowerCase()] || 'bg-gray-100 text-gray-800';

        return (
            <span className={`px-2 py-1 text-xs font-medium rounded-full ${color}`}>
                {status}
            </span>
        );
    }, []);

    if (initialLoading) {
        return <div className="text-center py-8">Cargando órdenes...</div>;
    }

    if (error) {
        return (
            <Alert type="error" onClose={() => setError(null)}>
                {error}
            </Alert>
        );
    }

    return (
        <div>
            {/* Dynamic Filters */}
            <div>
                <DynamicFilters
                    availableFilters={availableFilters}
                    activeFilters={activeFilters}
                    onAddFilter={handleAddFilter}
                    onRemoveFilter={handleRemoveFilter}
                    sortBy={filters.sort_by || 'created_at'}
                    sortOrder={filters.sort_order || 'desc'}
                    onSortChange={handleSortChange}
                    sortOptions={[
                        { value: 'created_at', label: 'Ordenar por fecha' },
                        { value: 'updated_at', label: 'Ordenar por actualización' },
                        { value: 'total_amount', label: 'Ordenar por monto' },
                    ]}
                />
            </div>

            {/* Table */}
            <div className="bg-white rounded-b-lg rounded-t-none shadow-sm border border-gray-200 border-t-0 overflow-hidden relative">
                {/* Overlay de carga solo para la tabla */}
                {tableLoading && (
                    <div className="absolute inset-0 bg-white/80 backdrop-blur-sm z-10 flex items-center justify-center transition-opacity duration-200">
                        <div className="flex flex-col items-center gap-2">
                            <div className="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
                            <p className="text-sm text-gray-600">Actualizando...</p>
                        </div>
                    </div>
                )}
                <div className="overflow-x-auto">
                    <table className={`min-w-full divide-y divide-gray-200 transition-opacity duration-200 ${tableLoading ? 'opacity-50' : 'opacity-100'}`}>
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Orden
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden sm:table-cell">
                                    Cliente
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Total
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Estado
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden lg:table-cell">
                                    Plataforma
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden md:table-cell">
                                    Fecha
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {orders.length === 0 ? (
                                <tr>
                                    <td colSpan={7} className="px-4 sm:px-6 py-8 text-center text-gray-500">
                                        No hay órdenes disponibles
                                    </td>
                                </tr>
                            ) : (
                                orders.map((order) => (
                                    <OrderRow
                                        key={order.id}
                                        order={order}
                                        onView={onView}
                                        onEdit={onEdit}
                                        onDelete={handleDelete}
                                        onShowRaw={(id) => {
                                            setSelectedOrderId(id);
                                            setIsRawModalOpen(true);
                                        }}
                                        formatCurrency={formatCurrency}
                                        formatDate={formatDate}
                                        getStatusBadge={getStatusBadge}
                                    />
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {/* Pagination */}
                {(totalPages > 1 || total > 0) && (
                    <div className="bg-white px-3 sm:px-4 lg:px-6 py-3 flex flex-col sm:flex-row items-center justify-between gap-3 border-t border-gray-200">
                        {/* Mobile: Simple pagination */}
                        <div className="flex-1 flex justify-between sm:hidden w-full">
                            <Button
                                variant="outline"
                                onClick={() => setFilters({ ...filters, page: page - 1 })}
                                disabled={page === 1}
                                size="sm"
                            >
                                Anterior
                            </Button>
                            <Button
                                variant="outline"
                                onClick={() => setFilters({ ...filters, page: page + 1 })}
                                disabled={page === totalPages}
                                size="sm"
                            >
                                Siguiente
                            </Button>
                        </div>

                        {/* Desktop: Full pagination */}
                        <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between w-full">
                            <div className="flex items-center gap-3">
                                <p className="text-xs sm:text-sm text-gray-700">
                                    Mostrando <span className="font-medium">{(page - 1) * (filters.page_size || 20) + 1}</span> a{' '}
                                    <span className="font-medium">{Math.min(page * (filters.page_size || 20), total)}</span> de{' '}
                                    <span className="font-medium">{total}</span> resultados
                                </p>
                                <div className="flex items-center gap-2">
                                    <label className="text-xs sm:text-sm text-gray-700 whitespace-nowrap">
                                        Mostrar:
                                    </label>
                                    <select
                                        value={filters.page_size || 20}
                                        onChange={(e) => {
                                            const newPageSize = parseInt(e.target.value);
                                            setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                        }}
                                        className="px-2 py-1.5 text-xs sm:text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                                    >
                                        <option value="10">10</option>
                                        <option value="20">20</option>
                                        <option value="50">50</option>
                                        <option value="100">100</option>
                                    </select>
                                </div>
                            </div>
                            <div className="flex items-center gap-2">
                                <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page - 1 })}
                                        disabled={page === 1}
                                        className="relative inline-flex items-center px-2 sm:px-3 py-2 rounded-l-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Anterior
                                    </button>
                                    <span className="relative inline-flex items-center px-3 sm:px-4 py-2 border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-700">
                                        Página {page} de {totalPages}
                                    </span>
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page + 1 })}
                                        disabled={page === totalPages}
                                        className="relative inline-flex items-center px-2 sm:px-3 py-2 rounded-r-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Siguiente
                                    </button>
                                </nav>
                            </div>
                        </div>

                        {/* Mobile: Page size selector */}
                        <div className="flex items-center justify-between w-full sm:hidden pt-2 border-t border-gray-200">
                            <div className="flex items-center gap-2">
                                <label className="text-xs text-gray-700 whitespace-nowrap">
                                    Mostrar:
                                </label>
                                <select
                                    value={filters.page_size || 20}
                                    onChange={(e) => {
                                        const newPageSize = parseInt(e.target.value);
                                        setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                    }}
                                    className="px-2 py-1.5 text-xs border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                                >
                                    <option value="10">10</option>
                                    <option value="20">20</option>
                                    <option value="50">50</option>
                                    <option value="100">100</option>
                                </select>
                            </div>
                            <p className="text-xs text-gray-500">
                                Página {page} de {totalPages}
                            </p>
                        </div>
                    </div>
                )}
            </div>

            {selectedOrderId && (
                <RawOrderModal
                    orderId={selectedOrderId}
                    isOpen={isRawModalOpen}
                    onClose={() => {
                        setIsRawModalOpen(false);
                        setSelectedOrderId(null);
                    }}
                />
            )}
        </div>
    );
}
