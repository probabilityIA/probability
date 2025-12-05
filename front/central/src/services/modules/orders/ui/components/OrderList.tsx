'use client';

import { useState, useEffect } from 'react';
import { getOrdersAction, deleteOrderAction } from '../../infra/actions';
import { Order, GetOrdersParams } from '../../domain/types';
import { Button, Alert } from '@/shared/ui';

interface OrderListProps {
    onView?: (order: Order) => void;
    onEdit?: (order: Order) => void;
}

export default function OrderList({ onView, onEdit }: OrderListProps) {
    const [orders, setOrders] = useState<Order[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);

    // Filters
    const [filters, setFilters] = useState<GetOrdersParams>({
        page: 1,
        page_size: 20,
    });

    useEffect(() => {
        fetchOrders();
    }, [filters]);

    const fetchOrders = async () => {
        setLoading(true);
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
            setLoading(false);
        }
    };

    const handleDelete = async (id: string) => {
        if (!confirm('¿Estás seguro de que deseas eliminar esta orden?')) return;

        try {
            const response = await deleteOrderAction(id);
            if (response.success) {
                fetchOrders();
            } else {
                alert(response.message || 'Error al eliminar la orden');
            }
        } catch (err: any) {
            alert(err.message || 'Error al eliminar la orden');
        }
    };

    const formatCurrency = (amount: number, currency: string = 'USD') => {
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency,
        }).format(amount);
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('es-CO', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const getStatusBadge = (status: string) => {
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
    };

    if (loading) {
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
        <div className="space-y-4">
            {/* Filters */}
            <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <input
                        type="text"
                        placeholder="Buscar por email..."
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        onChange={(e) => setFilters({ ...filters, customer_email: e.target.value || undefined })}
                    />
                    <select
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        onChange={(e) => setFilters({ ...filters, status: e.target.value || undefined })}
                    >
                        <option value="">Todos los estados</option>
                        <option value="pending">Pendiente</option>
                        <option value="processing">Procesando</option>
                        <option value="shipped">Enviado</option>
                        <option value="delivered">Entregado</option>
                        <option value="cancelled">Cancelado</option>
                    </select>
                    <select
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        onChange={(e) => setFilters({ ...filters, platform: e.target.value || undefined })}
                    >
                        <option value="">Todas las plataformas</option>
                        <option value="shopify">Shopify</option>
                        <option value="woocommerce">WooCommerce</option>
                        <option value="manual">Manual</option>
                    </select>
                    <Button variant="outline" onClick={fetchOrders}>
                        🔄 Actualizar
                    </Button>
                </div>
            </div>

            {/* Table */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Orden
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Cliente
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Total
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Estado
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Plataforma
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Fecha
                                </th>
                                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {orders.length === 0 ? (
                                <tr>
                                    <td colSpan={7} className="px-6 py-8 text-center text-gray-500">
                                        No hay órdenes disponibles
                                    </td>
                                </tr>
                            ) : (
                                orders.map((order) => (
                                    <tr key={order.id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <div className="text-sm font-medium text-gray-900">
                                                {order.order_number}
                                            </div>
                                            <div className="text-xs text-gray-500">
                                                {order.internal_number}
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="text-sm text-gray-900">{order.customer_name}</div>
                                            <div className="text-xs text-gray-500">{order.customer_email}</div>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <div className="text-sm font-semibold text-gray-900">
                                                {formatCurrency(order.total_amount, order.currency)}
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            {getStatusBadge(order.status)}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className="text-sm text-gray-900 capitalize">
                                                {order.platform}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {formatDate(order.created_at)}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <div className="flex justify-end gap-2">
                                                {onView && (
                                                    <button
                                                        onClick={() => onView(order)}
                                                        className="text-blue-600 hover:text-blue-900"
                                                    >
                                                        Ver
                                                    </button>
                                                )}
                                                {onEdit && (
                                                    <button
                                                        onClick={() => onEdit(order)}
                                                        className="text-indigo-600 hover:text-indigo-900"
                                                    >
                                                        Editar
                                                    </button>
                                                )}
                                                <button
                                                    onClick={() => handleDelete(order.id)}
                                                    className="text-red-600 hover:text-red-900"
                                                >
                                                    Eliminar
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {/* Pagination */}
                {totalPages > 1 && (
                    <div className="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
                        <div className="flex-1 flex justify-between sm:hidden">
                            <Button
                                variant="outline"
                                onClick={() => setFilters({ ...filters, page: page - 1 })}
                                disabled={page === 1}
                            >
                                Anterior
                            </Button>
                            <Button
                                variant="outline"
                                onClick={() => setFilters({ ...filters, page: page + 1 })}
                                disabled={page === totalPages}
                            >
                                Siguiente
                            </Button>
                        </div>
                        <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                            <div>
                                <p className="text-sm text-gray-700">
                                    Mostrando <span className="font-medium">{(page - 1) * (filters.page_size || 20) + 1}</span> a{' '}
                                    <span className="font-medium">{Math.min(page * (filters.page_size || 20), total)}</span> de{' '}
                                    <span className="font-medium">{total}</span> resultados
                                </p>
                            </div>
                            <div>
                                <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page - 1 })}
                                        disabled={page === 1}
                                        className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Anterior
                                    </button>
                                    <span className="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700">
                                        Página {page} de {totalPages}
                                    </span>
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page + 1 })}
                                        disabled={page === totalPages}
                                        className="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Siguiente
                                    </button>
                                </nav>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
