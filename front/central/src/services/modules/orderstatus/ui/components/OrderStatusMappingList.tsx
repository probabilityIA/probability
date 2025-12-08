'use client';

import { useState, useEffect } from 'react';
import {
    getOrderStatusMappingsAction,
    deleteOrderStatusMappingAction,
    toggleOrderStatusMappingActiveAction
} from '../../infra/actions';
import { OrderStatusMapping, GetOrderStatusMappingsParams } from '../../domain/types';
import { Button, Alert, Badge } from '@/shared/ui';

interface OrderStatusMappingListProps {
    onView?: (mapping: OrderStatusMapping) => void;
    onEdit?: (mapping: OrderStatusMapping) => void;
}

export default function OrderStatusMappingList({ onView, onEdit }: OrderStatusMappingListProps) {
    const [mappings, setMappings] = useState<OrderStatusMapping[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [total, setTotal] = useState(0);

    // Filters
    const [filters, setFilters] = useState<GetOrderStatusMappingsParams>({});

    useEffect(() => {
        fetchMappings();
    }, [filters]);

    const fetchMappings = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getOrderStatusMappingsAction(filters);
            if (response.success && response.data) {
                setMappings(response.data);
                setTotal(response.total || 0);
            } else {
                setError(response.message || 'Error al cargar los mappings');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar los mappings');
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id: number) => {
        if (!confirm('¿Estás seguro de que deseas eliminar este mapping?')) return;

        try {
            const response = await deleteOrderStatusMappingAction(id);
            if (response.success) {
                fetchMappings();
            } else {
                alert(response.message || 'Error al eliminar el mapping');
            }
        } catch (err: any) {
            alert(err.message || 'Error al eliminar el mapping');
        }
    };

    const handleToggleActive = async (id: number) => {
        try {
            const response = await toggleOrderStatusMappingActiveAction(id);
            if (response.success) {
                fetchMappings();
            } else {
                alert(response.message || 'Error al cambiar el estado');
            }
        } catch (err: any) {
            alert(err.message || 'Error al cambiar el estado');
        }
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

    if (loading) {
        return <div className="text-center py-8">Cargando mappings...</div>;
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
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <select
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        onChange={(e) => setFilters({ ...filters, integration_type: e.target.value || undefined })}
                    >
                        <option value="">Todos los tipos de integración</option>
                        <option value="shopify">Shopify</option>
                        <option value="whatsapp">WhatsApp</option>
                        <option value="woocommerce">WooCommerce</option>
                    </select>
                    <select
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        onChange={(e) => setFilters({ ...filters, is_active: e.target.value === '' ? undefined : e.target.value === 'true' })}
                    >
                        <option value="">Todos los estados</option>
                        <option value="true">Activos</option>
                        <option value="false">Inactivos</option>
                    </select>
                    <Button variant="outline" onClick={fetchMappings}>
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
                                    Tipo de Integración
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Estado Original
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Estado Mapeado
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Prioridad
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Estado
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Fecha Creación
                                </th>
                                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {mappings.length === 0 ? (
                                <tr>
                                    <td colSpan={7} className="px-6 py-8 text-center text-gray-500">
                                        No hay mappings disponibles
                                    </td>
                                </tr>
                            ) : (
                                mappings.map((mapping) => (
                                    <tr key={mapping.id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className="text-sm font-medium text-gray-900 capitalize">
                                                {mapping.integration_type}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className="text-sm text-gray-900">
                                                {mapping.original_status}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className="text-sm font-medium text-gray-900">
                                                {mapping.mapped_status}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className="text-sm text-gray-900">
                                                {mapping.priority}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <Badge type={mapping.is_active ? 'success' : 'secondary'}>
                                                {mapping.is_active ? 'Activo' : 'Inactivo'}
                                            </Badge>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {formatDate(mapping.created_at)}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <div className="flex justify-end gap-2">
                                                {onView && (
                                                    <button
                                                        onClick={() => onView(mapping)}
                                                        className="text-blue-600 hover:text-blue-900"
                                                    >
                                                        Ver
                                                    </button>
                                                )}
                                                {onEdit && (
                                                    <button
                                                        onClick={() => onEdit(mapping)}
                                                        className="text-indigo-600 hover:text-indigo-900"
                                                    >
                                                        Editar
                                                    </button>
                                                )}
                                                <button
                                                    onClick={() => handleToggleActive(mapping.id)}
                                                    className="text-yellow-600 hover:text-yellow-900"
                                                >
                                                    {mapping.is_active ? 'Desactivar' : 'Activar'}
                                                </button>
                                                <button
                                                    onClick={() => handleDelete(mapping.id)}
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

                {/* Summary */}
                {total > 0 && (
                    <div className="bg-white px-4 py-3 border-t border-gray-200">
                        <p className="text-sm text-gray-700">
                            Total de mappings: <span className="font-medium">{total}</span>
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
}
