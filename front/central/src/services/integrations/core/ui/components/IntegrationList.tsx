'use client';

import { useState, useEffect, useMemo, useCallback } from 'react';
import { 
    PlayIcon, 
    ArrowPathIcon, 
    PencilIcon, 
    PowerIcon, 
    TrashIcon 
} from '@heroicons/react/24/outline';
import { useIntegrations } from '../hooks/useIntegrations';
import { getActiveIntegrationTypesAction } from '../../infra/actions';
import {
    Integration,
    SyncOrdersParams
} from '../../domain/types';
import { getOrderStatusMappingsAction } from '@/services/modules/orderstatus/infra/actions';
import { OrderStatusMapping } from '@/services/modules/orderstatus/domain/types';
import { Input, Button, Badge, Spinner, Table, Alert, ConfirmModal, Select, DateRangePicker } from '@/shared/ui';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui/dynamic-filters';

interface IntegrationListProps {
    onEdit?: (integration: Integration) => void;
}

// Estado inicial para los filtros de sincronización
const initialSyncFilters: SyncOrdersParams = {
    created_at_min: '',
    created_at_max: '',
    status: 'any',
    financial_status: 'any',
    fulfillment_status: 'any'
};

export default function IntegrationList({ onEdit }: IntegrationListProps) {
    const {
        integrations,
        loading,
        error,
        page,
        setPage,
        totalPages,
        search,
        setSearch,
        filterType,
        setFilterType,
        filterCategory,
        setFilterCategory,
        deleteIntegration,
        toggleActive,
        testConnection,
        syncOrders,
        setError
    } = useIntegrations();

    const [deleteModal, setDeleteModal] = useState<{ show: boolean; id: number | null }>({
        show: false,
        id: null
    });

    // Estado para tipos de integración para filtros
    const [integrationTypes, setIntegrationTypes] = useState<Array<{ value: string; label: string }>>([]);

    // Obtener tipos de integración para el filtro
    useEffect(() => {
        const fetchIntegrationTypes = async () => {
            try {
                const response = await getActiveIntegrationTypesAction();
                if (response.success && response.data) {
                    const options = response.data.map((type) => ({
                        value: type.code || type.name.toLowerCase().replace(/\s+/g, '_'),
                        label: type.name,
                    }));
                    setIntegrationTypes(options);
                }
            } catch (err) {
                console.error('Error fetching integration types for filter:', err);
            }
        };
        fetchIntegrationTypes();
    }, []);

    // Estado para el modal de sincronización
    const [syncModal, setSyncModal] = useState<{ show: boolean; id: number | null; name: string; integrationTypeId?: number }>({
        show: false,
        id: null,
        name: ''
    });
    const [syncFilters, setSyncFilters] = useState<SyncOrdersParams>(initialSyncFilters);
    const [syncing, setSyncing] = useState(false);
    
    // Estados para las opciones de filtros dinámicos desde la BD
    const [orderStatusOptions, setOrderStatusOptions] = useState<Array<{ value: string; label: string; mappedStatus?: string }>>([]);
    const [financialStatusOptions, setFinancialStatusOptions] = useState<Array<{ value: string; label: string; mappedStatus?: string }>>([]);
    const [fulfillmentStatusOptions, setFulfillmentStatusOptions] = useState<Array<{ value: string; label: string; mappedStatus?: string }>>([]);
    const [loadingMappings, setLoadingMappings] = useState(false);

    // Filtros dinámicos - sincronizar con el hook
    const [filters, setFilters] = useState<{
        search?: string;
        type?: string;
        category?: string;
    }>({});

    // Definir filtros disponibles
    const availableFilters: FilterOption[] = useMemo(() => [
        {
            key: 'search',
            label: 'Nombre',
            type: 'text',
            placeholder: 'Buscar por nombre...',
        },
        {
            key: 'type',
            label: 'Tipo',
            type: 'select',
            options: integrationTypes,
        },
        {
            key: 'category',
            label: 'Categoría',
            type: 'select',
            options: [
                { value: 'internal', label: 'Interna' },
                { value: 'external', label: 'Externa' },
            ],
        },
    ], [integrationTypes]);

    // Convertir filtros a ActiveFilter[]
    const activeFilters: ActiveFilter[] = useMemo(() => {
        const active: ActiveFilter[] = [];

        if (filters.search) {
            active.push({
                key: 'search',
                label: 'Nombre',
                value: filters.search,
                type: 'text',
            });
        }

        if (filters.type) {
            active.push({
                key: 'type',
                label: 'Tipo',
                value: filters.type,
                type: 'select',
            });
        }

        if (filters.category) {
            active.push({
                key: 'category',
                label: 'Categoría',
                value: filters.category,
                type: 'select',
            });
        }

        return active;
    }, [filters]);

    // Manejar agregar filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev, [filterKey]: value };
            // Sincronizar con el estado del hook
            if (filterKey === 'search') {
                setSearch(value);
            } else if (filterKey === 'type') {
                setFilterType(value);
            } else if (filterKey === 'category') {
                setFilterCategory(value);
            }
            return newFilters;
        });
        setPage(1);
    }, [setSearch, setFilterType, setFilterCategory, setPage]);

    // Manejar eliminar filtro
    const handleRemoveFilter = useCallback((filterKey: string) => {
        setFilters((prev) => {
            const newFilters = { ...prev };
            delete (newFilters as any)[filterKey];
            // Sincronizar con el estado del hook
            if (filterKey === 'search') {
                setSearch('');
            } else if (filterKey === 'type') {
                setFilterType('');
            } else if (filterKey === 'category') {
                setFilterCategory('');
            }
            return newFilters;
        });
        setPage(1);
    }, [setSearch, setFilterType, setFilterCategory, setPage]);

    // Manejar cambio de ordenamiento
    const handleSortChange = useCallback((sortBy: string, sortOrder: 'asc' | 'desc') => {
        // Por ahora no tenemos ordenamiento en integraciones, pero podemos agregarlo después
        console.log('Sort changed:', sortBy, sortOrder);
    }, []);

    const handleDeleteClick = (id: number) => {
        setDeleteModal({ show: true, id });
    };

    const handleDeleteConfirm = async () => {
        if (deleteModal.id) {
            const success = await deleteIntegration(deleteModal.id);
            if (success) {
                setDeleteModal({ show: false, id: null });
            }
        }
    };

    const handleTest = async (id: number) => {
        const result = await testConnection(id);
        if (result.success) {
            alert('✅ Conexión exitosa');
        } else {
            alert(`❌ Error: ${result.message}`);
        }
    };

    // Estados de Shopify por categoría (para agrupar los mapeos)
    const SHOPIFY_ORDER_STATUSES = ['any', 'open', 'closed', 'cancelled'];
    const SHOPIFY_FINANCIAL_STATUSES = ['any', 'authorized', 'pending', 'paid', 'partially_paid', 'refunded', 'voided', 'partially_refunded', 'unpaid'];
    const SHOPIFY_FULFILLMENT_STATUSES = ['any', 'shipped', 'partial', 'unfulfilled', 'unshipped'];

    // Cargar mapeos de estados para la integración
    const loadStatusMappings = async (integrationTypeId: number) => {
        setLoadingMappings(true);
        try {
            const response = await getOrderStatusMappingsAction({ 
                integration_type_id: integrationTypeId,
                is_active: true 
            });
            
            // El backend devuelve { data: [...], total: number }
            // Puede venir con o sin el campo success dependiendo del endpoint
            const mappings: OrderStatusMapping[] = (response as any).data || response.data || [];
            
            console.log('Status mappings response:', response);
            console.log('Status mappings loaded:', mappings);
            
            if (mappings && mappings.length > 0) {
                // Agrupar por tipo de estado
                const orderStatusMap = mappings
                    .filter(m => SHOPIFY_ORDER_STATUSES.includes(m.original_status))
                    .map(m => ({
                        value: m.original_status,
                        label: `${m.original_status}${m.order_status ? ` → ${m.order_status.name}` : ''}`,
                        mappedStatus: m.order_status?.code
                    }));
                
                const financialStatusMap = mappings
                    .filter(m => SHOPIFY_FINANCIAL_STATUSES.includes(m.original_status))
                    .map(m => ({
                        value: m.original_status,
                        label: `${m.original_status}${m.order_status ? ` → ${m.order_status.name}` : ''}`,
                        mappedStatus: m.order_status?.code
                    }));
                
                const fulfillmentStatusMap = mappings
                    .filter(m => SHOPIFY_FULFILLMENT_STATUSES.includes(m.original_status))
                    .map(m => ({
                        value: m.original_status,
                        label: `${m.original_status}${m.order_status ? ` → ${m.order_status.name}` : ''}`,
                        mappedStatus: m.order_status?.code
                    }));
                
                console.log('Order status options:', orderStatusMap);
                console.log('Financial status options:', financialStatusMap);
                console.log('Fulfillment status options:', fulfillmentStatusMap);
                
                // Agregar opción "Todos" solo si no existe ya
                const hasAnyInOrder = orderStatusMap.some(opt => opt.value === 'any');
                const hasAnyInFinancial = financialStatusMap.some(opt => opt.value === 'any');
                const hasAnyInFulfillment = fulfillmentStatusMap.some(opt => opt.value === 'any');
                
                setOrderStatusOptions(
                    orderStatusMap.length > 0 
                        ? (hasAnyInOrder ? orderStatusMap : [{ value: 'any', label: 'Todos' }, ...orderStatusMap])
                        : [{ value: 'any', label: 'Todos' }]
                );
                setFinancialStatusOptions(
                    financialStatusMap.length > 0
                        ? (hasAnyInFinancial ? financialStatusMap : [{ value: 'any', label: 'Todos' }, ...financialStatusMap])
                        : [{ value: 'any', label: 'Todos' }]
                );
                setFulfillmentStatusOptions(
                    fulfillmentStatusMap.length > 0
                        ? (hasAnyInFulfillment ? fulfillmentStatusMap : [{ value: 'any', label: 'Todos' }, ...fulfillmentStatusMap])
                        : [{ value: 'any', label: 'Todos' }]
                );
            } else {
                // Si no hay mapeos, usar opciones por defecto
                console.warn('No mappings found, using default options');
                setOrderStatusOptions([{ value: 'any', label: 'Todos' }]);
                setFinancialStatusOptions([{ value: 'any', label: 'Todos' }]);
                setFulfillmentStatusOptions([{ value: 'any', label: 'Todos' }]);
            }
        } catch (err) {
            console.error('Error loading status mappings:', err);
            // En caso de error, usar opciones por defecto
            setOrderStatusOptions([{ value: 'any', label: 'Todos' }]);
            setFinancialStatusOptions([{ value: 'any', label: 'Todos' }]);
            setFulfillmentStatusOptions([{ value: 'any', label: 'Todos' }]);
        } finally {
            setLoadingMappings(false);
        }
    };

    // Abrir modal de sincronización
    const handleSyncClick = async (id: number, name: string) => {
        // Buscar la integración para obtener el integration_type_id
        const integration = integrations.find(i => i.id === id);
        if (!integration) {
            alert('Error: No se encontró la integración');
            return;
        }
        
        const integrationTypeId = integration.integration_type_id;
        
        setSyncModal({ show: true, id, name, integrationTypeId });
        
        // Cargar mapeos de estados
        await loadStatusMappings(integrationTypeId);
        
        // Establecer fecha mínima por defecto a 30 días atrás
        const thirtyDaysAgo = new Date();
        thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);
        setSyncFilters({
            ...initialSyncFilters,
            created_at_min: thirtyDaysAgo.toISOString().split('T')[0]
        });
    };

    // Ejecutar sincronización con filtros
    const handleSyncConfirm = async () => {
        if (syncModal.id) {
            setSyncing(true);
            try {
                // Preparar parámetros (solo enviar los que tienen valor)
                const params: SyncOrdersParams = {};
                if (syncFilters.created_at_min) params.created_at_min = syncFilters.created_at_min;
                if (syncFilters.created_at_max) params.created_at_max = syncFilters.created_at_max;
                if (syncFilters.status && syncFilters.status !== 'any') params.status = syncFilters.status;
                if (syncFilters.financial_status && syncFilters.financial_status !== 'any') params.financial_status = syncFilters.financial_status;
                if (syncFilters.fulfillment_status && syncFilters.fulfillment_status !== 'any') params.fulfillment_status = syncFilters.fulfillment_status;

                const result = await syncOrders(syncModal.id, Object.keys(params).length > 0 ? params : undefined);
                if (result.success) {
                    alert('✅ Sincronización iniciada correctamente');
                } else {
                    alert(`❌ Error al iniciar sincronización: ${result.message}`);
                }
            } finally {
                setSyncing(false);
                setSyncModal({ show: false, id: null, name: '' });
                setSyncFilters(initialSyncFilters);
            }
        }
    };

    // Cerrar modal de sincronización
    const handleSyncCancel = () => {
        setSyncModal({ show: false, id: null, name: '' });
        setSyncFilters(initialSyncFilters);
        // Limpiar opciones
        setOrderStatusOptions([]);
        setFinancialStatusOptions([]);
        setFulfillmentStatusOptions([]);
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    const columns = [
        { key: 'id', label: 'ID' },
        { key: 'logo', label: 'Logo' },
        { key: 'name', label: 'Nombre' },
        { key: 'type', label: 'Tipo' },
        { key: 'category', label: 'Categoría' },
        { key: 'status', label: 'Estado' },
        { key: 'actions', label: 'Acciones' }
    ];

    const renderRow = (integration: Integration) => ({
        id: integration.id,
        logo: (
            <div className="flex items-center justify-center">
                {integration.integration_type?.image_url ? (
                    <img
                        src={integration.integration_type.image_url}
                        alt={integration.integration_type.name || integration.name}
                        className="w-12 h-12 object-contain border border-gray-200 rounded-lg p-1 bg-white"
                        onError={(e) => {
                            // Si la imagen falla al cargar, mostrar un placeholder
                            (e.target as HTMLImageElement).style.display = 'none';
                            const parent = (e.target as HTMLImageElement).parentElement;
                            if (parent) {
                                parent.innerHTML = '<div class="w-12 h-12 flex items-center justify-center bg-gray-100 rounded-lg text-gray-400 text-xs">Sin logo</div>';
                            }
                        }}
                    />
                ) : (
                    <div className="w-12 h-12 flex items-center justify-center bg-gray-100 rounded-lg text-gray-400 text-xs">
                        {integration.name.charAt(0).toUpperCase()}
                    </div>
                )}
            </div>
        ),
        name: (
            <div>
                <div className="text-sm font-medium text-gray-900">{integration.name}</div>
                {integration.description && (
                    <div className="text-sm text-gray-500">{integration.description}</div>
                )}
            </div>
        ),
        type: integration.integration_type?.name || integration.type,
        category: integration.category,
        status: (
            <div className="flex items-center gap-2">
                <Badge type={integration.is_active ? 'success' : 'error'}>
                    {integration.is_active ? 'Activo' : 'Inactivo'}
                </Badge>
                {integration.is_default && (
                    <Badge type="primary">Por defecto</Badge>
                )}
            </div>
        ),
        actions: (
            <div className="flex gap-2 items-center">
                <button
                    onClick={() => handleTest(integration.id)}
                    className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                    title="Probar conexión"
                    aria-label="Probar conexión"
                >
                    <PlayIcon className="w-4 h-4" />
                </button>
                <button
                    onClick={() => handleSyncClick(integration.id, integration.name)}
                    className="p-2 bg-green-500 hover:bg-green-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-green-500 focus:ring-offset-2"
                    title="Sincronizar órdenes"
                    aria-label="Sincronizar órdenes"
                >
                    <ArrowPathIcon className="w-4 h-4" />
                </button>
                {onEdit && (
                    <button
                        onClick={() => onEdit(integration)}
                        className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                        title="Editar integración"
                        aria-label="Editar integración"
                    >
                        <PencilIcon className="w-4 h-4" />
                    </button>
                )}
                <button
                    onClick={() => toggleActive(integration.id, integration.is_active)}
                    className={`p-2 rounded-md transition-colors duration-200 focus:ring-2 focus:ring-offset-2 ${
                        integration.is_active
                            ? 'bg-orange-500 hover:bg-orange-600 text-white focus:ring-orange-500'
                            : 'bg-gray-500 hover:bg-gray-600 text-white focus:ring-gray-500'
                    }`}
                    title={integration.is_active ? 'Desactivar integración' : 'Activar integración'}
                    aria-label={integration.is_active ? 'Desactivar integración' : 'Activar integración'}
                >
                    <PowerIcon className="w-4 h-4" />
                </button>
                <button
                    onClick={() => handleDeleteClick(integration.id)}
                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                    title="Eliminar integración"
                    aria-label="Eliminar integración"
                >
                    <TrashIcon className="w-4 h-4" />
                </button>
            </div>
        )
    });

    return (
        <div className="space-y-4">
            {/* Dynamic Filters */}
            <DynamicFilters
                availableFilters={availableFilters}
                activeFilters={activeFilters}
                onAddFilter={handleAddFilter}
                onRemoveFilter={handleRemoveFilter}
                sortBy="created_at"
                sortOrder="desc"
                onSortChange={handleSortChange}
                sortOptions={[
                    { value: 'created_at', label: 'Ordenar por fecha de creación' },
                    { value: 'updated_at', label: 'Ordenar por fecha de actualización' },
                    { value: 'name', label: 'Ordenar por nombre' },
                ]}
            />

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="bg-white rounded-b-lg rounded-t-none shadow-sm border border-gray-200 border-t-0 overflow-hidden">
                <Table
                    columns={columns}
                    data={integrations.map(renderRow)}
                    emptyMessage="No hay integraciones disponibles"
                />
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
                <div className="flex justify-center items-center gap-2">
                    <Button
                        onClick={() => setPage(page - 1)}
                        disabled={page === 1}
                        variant="primary"
                    >
                        Anterior
                    </Button>
                    <span className="text-sm text-gray-700">
                        Página {page} de {totalPages}
                    </span>
                    <Button
                        onClick={() => setPage(page + 1)}
                        disabled={page === totalPages}
                        variant="primary"
                    >
                        Siguiente
                    </Button>
                </div>
            )}

            <ConfirmModal
                isOpen={deleteModal.show}
                onClose={() => setDeleteModal({ show: false, id: null })}
                onConfirm={handleDeleteConfirm}
                title="Eliminar Integración"
                message="¿Estás seguro de que deseas eliminar esta integración? Esta acción no se puede deshacer."
            />

            {/* Modal de Sincronización con Filtros */}
            {syncModal.show && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
                    <div className="bg-white rounded-lg shadow-xl w-full max-w-md mx-4 overflow-hidden">
                        <div className="px-6 py-4 border-b border-gray-200">
                            <h3 className="text-lg font-semibold text-gray-900">
                                ↻ Sincronizar Órdenes
                            </h3>
                            <p className="text-sm text-gray-500 mt-1">
                                Integración: <strong>{syncModal.name}</strong>
                            </p>
                        </div>

                        <div className="px-6 py-4 space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Rango de Fechas
                                </label>
                                <DateRangePicker
                                    startDate={syncFilters.created_at_min}
                                    endDate={syncFilters.created_at_max}
                                    onChange={(startDate, endDate) => {
                                        setSyncFilters(prev => ({
                                            ...prev,
                                            created_at_min: startDate || '',
                                            created_at_max: endDate || ''
                                        }));
                                    }}
                                    placeholder="Seleccionar rango de fechas (opcional)"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Estado de Orden
                                </label>
                                {loadingMappings ? (
                                    <div className="flex items-center gap-2 text-sm text-gray-500">
                                        <Spinner size="sm" /> Cargando estados...
                                    </div>
                                ) : (
                                    <Select
                                        value={syncFilters.status || 'any'}
                                        onChange={(e) => setSyncFilters(prev => ({
                                            ...prev,
                                            status: e.target.value
                                        }))}
                                        options={orderStatusOptions.length > 0 ? orderStatusOptions : [{ value: 'any', label: 'Todos' }]}
                                        disabled={orderStatusOptions.length === 0}
                                    />
                                )}
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Estado Financiero
                                </label>
                                {loadingMappings ? (
                                    <div className="flex items-center gap-2 text-sm text-gray-500">
                                        <Spinner size="sm" /> Cargando estados...
                                    </div>
                                ) : (
                                    <Select
                                        value={syncFilters.financial_status || 'any'}
                                        onChange={(e) => setSyncFilters(prev => ({
                                            ...prev,
                                            financial_status: e.target.value
                                        }))}
                                        options={financialStatusOptions.length > 0 ? financialStatusOptions : [{ value: 'any', label: 'Todos' }]}
                                        disabled={financialStatusOptions.length === 0}
                                    />
                                )}
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Estado de Envío
                                </label>
                                {loadingMappings ? (
                                    <div className="flex items-center gap-2 text-sm text-gray-500">
                                        <Spinner size="sm" /> Cargando estados...
                                    </div>
                                ) : (
                                    <Select
                                        value={syncFilters.fulfillment_status || 'any'}
                                        onChange={(e) => setSyncFilters(prev => ({
                                            ...prev,
                                            fulfillment_status: e.target.value
                                        }))}
                                        options={fulfillmentStatusOptions.length > 0 ? fulfillmentStatusOptions : [{ value: 'any', label: 'Todos' }]}
                                        disabled={fulfillmentStatusOptions.length === 0}
                                    />
                                )}
                            </div>

                            <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                                <p className="text-sm text-blue-700">
                                    💡 Los estados mostrados están mapeados a estados de Probability. Si no especificas filtros, se sincronizarán las órdenes de los últimos 30 días.
                                </p>
                            </div>
                        </div>

                        <div className="px-6 py-4 border-t border-gray-200 flex justify-end gap-3">
                            <Button
                                variant="outline"
                                onClick={handleSyncCancel}
                                disabled={syncing}
                            >
                                Cancelar
                            </Button>
                            <Button
                                variant="primary"
                                onClick={handleSyncConfirm}
                                disabled={syncing}
                            >
                                {syncing ? (
                                    <>
                                        <Spinner size="sm" /> Sincronizando...
                                    </>
                                ) : (
                                    '↻ Iniciar Sincronización'
                                )}
                            </Button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
