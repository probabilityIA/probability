'use client';

import { useState, useEffect, useMemo, useCallback } from 'react';
import {
    PencilIcon,
    PowerIcon,
    TrashIcon
} from '@heroicons/react/24/outline';
import {
    getOrderStatusMappingsAction,
    deleteOrderStatusMappingAction,
    toggleOrderStatusMappingActiveAction
} from '../../infra/actions';
import { getActiveIntegrationTypesAction } from '@/services/integrations/core/infra/actions';
import { OrderStatusMapping, GetOrderStatusMappingsParams } from '../../domain/types';
import { Alert, Badge, Table, Spinner } from '@/shared/ui';
import OrderStatusBadge from './OrderStatusBadge';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui/dynamic-filters';

interface OrderStatusMappingListProps {
    onView?: (mapping: OrderStatusMapping) => void;
    onEdit?: (mapping: OrderStatusMapping) => void;
}

export default function OrderStatusMappingList({ onView, onEdit }: OrderStatusMappingListProps) {
    const [mappings, setMappings] = useState<OrderStatusMapping[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [totalPages, setTotalPages] = useState(1);

    // Filters
    const [filters, setFilters] = useState<GetOrderStatusMappingsParams>({});
    const [integrationTypes, setIntegrationTypes] = useState<Array<{ value: string; label: string }>>([]);

    // Cargar tipos de integración para el filtro
    useEffect(() => {
        const fetchIntegrationTypes = async () => {
            try {
                const response = await getActiveIntegrationTypesAction();
                if (response.success && response.data) {
                    const options = response.data.map((type) => ({
                        value: String(type.id),
                        label: type.name,
                    }));
                    setIntegrationTypes(options);
                }
            } catch (err) {
                console.error('Error fetching integration types:', err);
            }
        };
        fetchIntegrationTypes();
    }, []);

    useEffect(() => {
        fetchMappings();
    }, [page, pageSize, filters]);

    const fetchMappings = async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetOrderStatusMappingsParams = {
                page,
                page_size: pageSize,
                ...filters
            };
            const response = await getOrderStatusMappingsAction(params);
            // El backend devuelve { data: [...], total: number, page, page_size, total_pages }
            const mappingsData = (response as any).data || response.data || [];
            const totalCount = (response as any).total || response.total || 0;
            const currentPage = (response as any).page || page;
            const currentPageSize = (response as any).page_size || pageSize;
            const pagesTotal = (response as any).total_pages || Math.ceil(totalCount / currentPageSize) || 1;
            
            if (mappingsData && mappingsData.length >= 0) {
                setMappings(mappingsData);
                setTotal(totalCount);
                setPage(currentPage);
                setPageSize(currentPageSize);
                setTotalPages(pagesTotal);
            } else {
                setError('Error al cargar los mappings');
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

    // Definir filtros disponibles
    const availableFilters: FilterOption[] = useMemo(() => [
        {
            key: 'integration_type_id',
            label: 'Tipo de Integración',
            type: 'select',
            options: integrationTypes,
        },
        {
            key: 'is_active',
            label: 'Estado',
            type: 'select',
            options: [
                { value: 'true', label: 'Activo' },
                { value: 'false', label: 'Inactivo' },
            ],
        },
    ], [integrationTypes]);

    // Convertir filtros a ActiveFilter[]
    const activeFilters: ActiveFilter[] = useMemo(() => {
        const active: ActiveFilter[] = [];

        if (filters.integration_type_id) {
            const type = integrationTypes.find(t => t.value === String(filters.integration_type_id));
            active.push({
                key: 'integration_type_id',
                label: 'Tipo de Integración',
                value: type?.label || String(filters.integration_type_id),
                type: 'select',
            });
        }

        if (filters.is_active !== undefined) {
            active.push({
                key: 'is_active',
                label: 'Estado',
                value: filters.is_active ? 'Activo' : 'Inactivo',
                type: 'select',
            });
        }

        return active;
    }, [filters, integrationTypes]);

    // Manejar agregar filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev };
            if (filterKey === 'integration_type_id') {
                newFilters.integration_type_id = value ? parseInt(value) : undefined;
            } else if (filterKey === 'is_active') {
                newFilters.is_active = value === 'true' ? true : value === 'false' ? false : undefined;
            }
            return newFilters;
        });
        setPage(1);
    }, []);

    // Manejar eliminar filtro
    const handleRemoveFilter = useCallback((filterKey: string) => {
        setFilters((prev) => {
            const newFilters = { ...prev };
            if (filterKey === 'integration_type_id') {
                delete (newFilters as any).integration_type_id;
            } else if (filterKey === 'is_active') {
                delete (newFilters as any).is_active;
            }
            return newFilters;
        });
        setPage(1);
    }, []);

    // Definir columnas de la tabla
    const columns = [
        { key: 'integration_type', label: 'Tipo de Integración' },
        { key: 'original_status', label: 'Estado Original' },
        { key: 'order_status', label: 'Estado de Probability' },
        { key: 'description', label: 'Descripción' },
        { key: 'priority', label: 'Prioridad', align: 'center' as const },
        { key: 'is_active', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    // Renderizar filas
    const renderRow = (mapping: OrderStatusMapping) => ({
        integration_type: (
            <div className="flex items-center justify-center">
                {mapping.integration_type?.image_url ? (
                    <img
                        src={mapping.integration_type.image_url}
                        alt={mapping.integration_type.name || `ID: ${mapping.integration_type_id}`}
                        className="h-10 w-10 object-contain border border-gray-200 rounded-lg p-1 bg-white"
                        onError={(e) => {
                            // Si la imagen falla al cargar, mostrar un placeholder
                            (e.target as HTMLImageElement).style.display = 'none';
                            const parent = (e.target as HTMLImageElement).parentElement;
                            if (parent) {
                                parent.innerHTML = `<span class="text-xs font-medium text-gray-600 uppercase">${(mapping.integration_type?.name || `ID: ${mapping.integration_type_id}`).charAt(0)}</span>`;
                            }
                        }}
                    />
                ) : (
                    <div className="h-10 w-10 flex items-center justify-center bg-gray-100 rounded-lg text-gray-400 text-xs">
                        {(mapping.integration_type?.name || `ID: ${mapping.integration_type_id}`).charAt(0).toUpperCase()}
                    </div>
                )}
            </div>
        ),
        original_status: (
            <span className="text-sm text-gray-900 font-mono">
                {mapping.original_status}
            </span>
        ),
        order_status: (
            <OrderStatusBadge
                status={mapping.order_status}
                fallback={`ID: ${mapping.order_status_id}`}
            />
        ),
        description: (
            <span className="text-sm text-gray-500">
                {mapping.description || <span className="text-gray-300">—</span>}
            </span>
        ),
        priority: (
            <span className="text-sm text-gray-900">
                {mapping.order_status?.priority ?? '—'}
            </span>
        ),
        is_active: (
            <Badge type={mapping.is_active ? 'success' : 'secondary'}>
                {mapping.is_active ? 'Activo' : 'Inactivo'}
            </Badge>
        ),
        actions: (
            <div className="flex justify-end gap-2">
                {onEdit && (
                    <button
                        onClick={() => onEdit(mapping)}
                        className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                        title="Editar mapping"
                        aria-label="Editar mapping"
                    >
                        <PencilIcon className="w-4 h-4" />
                    </button>
                )}
                <button
                    onClick={() => handleToggleActive(mapping.id)}
                    className={`p-2 rounded-md transition-colors duration-200 focus:ring-2 focus:ring-offset-2 ${
                        mapping.is_active
                            ? 'bg-orange-500 hover:bg-orange-600 text-white focus:ring-orange-500'
                            : 'bg-gray-500 hover:bg-gray-600 text-white focus:ring-gray-500'
                    }`}
                    title={mapping.is_active ? 'Desactivar mapping' : 'Activar mapping'}
                    aria-label={mapping.is_active ? 'Desactivar mapping' : 'Activar mapping'}
                >
                    <PowerIcon className="w-4 h-4" />
                </button>
                <button
                    onClick={() => handleDelete(mapping.id)}
                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                    title="Eliminar mapping"
                    aria-label="Eliminar mapping"
                >
                    <TrashIcon className="w-4 h-4" />
                </button>
            </div>
        ),
    });

    if (loading && mappings.length === 0) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
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
            <DynamicFilters
                availableFilters={availableFilters}
                activeFilters={activeFilters}
                onAddFilter={handleAddFilter}
                onRemoveFilter={handleRemoveFilter}
            />

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="bg-white rounded-b-lg rounded-t-none shadow-sm border border-gray-200 border-t-0 overflow-hidden">
                <Table
                    columns={columns}
                    data={mappings.map(renderRow)}
                    keyExtractor={(_, index) => String(mappings[index]?.id || index)}
                    emptyMessage="No hay mappings disponibles"
                    loading={loading}
                    pagination={{
                        currentPage: page,
                        totalPages: totalPages,
                        totalItems: total,
                        itemsPerPage: pageSize,
                        onPageChange: (newPage) => setPage(newPage),
                        onItemsPerPageChange: (newPageSize) => {
                            setPageSize(newPageSize);
                            setPage(1);
                        },
                    }}
                />
            </div>
        </div>
    );
}
