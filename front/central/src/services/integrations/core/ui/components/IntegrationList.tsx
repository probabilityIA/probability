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
import { useSSE } from '@/shared/hooks/use-sse';
import { getActiveIntegrationTypesAction } from '../../infra/actions';
import {
    Integration,
    SyncOrdersParams
} from '../../domain/types';
import { getOrderStatusMappingsAction } from '@/services/modules/orderstatus/infra/actions';
import { OrderStatusMapping } from '@/services/modules/orderstatus/domain/types';
import { Input, Button, Badge, Spinner, Table, Alert, ConfirmModal, Select, DateRangePicker } from '@/shared/ui';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui/dynamic-filters';
import { useToast } from '@/shared/providers/toast-provider';
import { TokenStorage } from '@/shared/utils/token-storage';
import { playNotificationSound } from '@/shared/utils';

interface IntegrationListProps {
    onEdit?: (integration: Integration) => void;
    filterCategory?: string;
}

// Estado inicial para los filtros de sincronización
const initialSyncFilters: SyncOrdersParams = {
    created_at_min: '',
    created_at_max: '',
    status: 'any',
    financial_status: 'any',
    fulfillment_status: 'any'
};

export default function IntegrationList({ onEdit, filterCategory: propFilterCategory }: IntegrationListProps) {
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

    // Sincronizar filtro de categoría desde las props
    useEffect(() => {
        if (propFilterCategory !== undefined && propFilterCategory !== filterCategory) {
            setFilterCategory(propFilterCategory);
        }
    }, [propFilterCategory, filterCategory, setFilterCategory]);

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
                const token = TokenStorage.getSessionToken();
                const response = await getActiveIntegrationTypesAction(token);
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

    // Estados para el progreso de sincronización en tiempo real
    const [syncProgress, setSyncProgress] = useState<{
        total: number;
        created: number;
        rejected: number;
        updated: number;
        orders: Array<{
            orderNumber: string;
            status: 'created' | 'rejected' | 'updated';
            reason?: string;
            createdAt?: string;
            orderStatus?: string;
            timestamp: Date;
        }>;
    } | null>(null);

    // Estados para las opciones de filtros dinámicos desde la BD
    const [orderStatusOptions, setOrderStatusOptions] = useState<Array<{ value: string; label: string; mappedStatus?: string }>>([]);
    const [financialStatusOptions, setFinancialStatusOptions] = useState<Array<{ value: string; label: string; mappedStatus?: string }>>([]);
    const [fulfillmentStatusOptions, setFulfillmentStatusOptions] = useState<Array<{ value: string; label: string; mappedStatus?: string }>>([]);
    const [loadingMappings, setLoadingMappings] = useState(false);

    // Toast para notificaciones
    const { showToast } = useToast();

    // Obtener businessId del usuario actual
    const [currentBusinessId, setCurrentBusinessId] = useState<number | undefined>(undefined);

    useEffect(() => {
        const businessesData = TokenStorage.getBusinessesData();
        if (businessesData && businessesData.length > 0) {
            // Usar el primer business o el activo si existe
            const activeBusinessId = localStorage.getItem('active_business_id');
            const businessId = activeBusinessId
                ? parseInt(activeBusinessId, 10)
                : businessesData[0].id;
            setCurrentBusinessId(businessId);
        }
    }, []);

    // Hook para escuchar eventos de sincronización en tiempo real via SSE centralizado (modules/events)
    // Solo escuchar cuando el modal está abierto y hay una integración seleccionada
    const integrationEventTypes = syncModal.show ? [
        'integration.sync.order.created',
        'integration.sync.order.updated',
        'integration.sync.order.rejected',
        'integration.sync.started',
        'integration.sync.completed',
        'integration.sync.failed'
    ] : [];

    const { isConnected } = useSSE({
        businessId: currentBusinessId,
        integrationId: syncModal.show && syncModal.id ? syncModal.id : undefined,
        eventTypes: integrationEventTypes,
        onMessage: (messageEvent: MessageEvent) => {
            try {
                const event = JSON.parse(messageEvent.data);
                const eventType = event.event_type || event.type || messageEvent.type;

                // Solo procesar si es de la integración que está sincronizando
                if (syncModal.id && event.integration_id !== syncModal.id) return;

                const eventData = event.data?.data || event.data || {};

                switch (eventType) {
                    case 'integration.sync.order.created': {
                        const orderNumber = eventData.order_number || event.metadata?.order_number || 'Desconocida';
                        const createdAt = eventData.created_at || eventData.synced_at || event.metadata?.created_at || null;
                        const orderStatus = eventData.status || event.metadata?.status || null;

                        setSyncProgress(prev => {
                            if (!prev) {
                                return { total: 1, created: 1, rejected: 0, updated: 0, orders: [{ orderNumber, status: 'created' as const, createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date() }] };
                            }
                            return { ...prev, created: prev.created + 1, total: prev.total + 1, orders: [{ orderNumber, status: 'created' as const, createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date() }, ...prev.orders] };
                        });
                        playNotificationSound();
                        showToast(`Orden creada: #${orderNumber}`, 'success');
                        break;
                    }
                    case 'integration.sync.order.updated': {
                        const orderNumber = eventData.order_number || event.metadata?.order_number || 'Desconocida';
                        const createdAt = eventData.created_at || eventData.updated_at || event.metadata?.created_at || null;
                        const orderStatus = eventData.status || event.metadata?.status || null;

                        setSyncProgress(prev => {
                            if (!prev) {
                                return { total: 1, created: 0, rejected: 0, updated: 1, orders: [{ orderNumber, status: 'updated' as const, createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date() }] };
                            }
                            return { ...prev, updated: prev.updated + 1, total: prev.total + 1, orders: [{ orderNumber, status: 'updated' as const, createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date() }, ...prev.orders] };
                        });
                        playNotificationSound();
                        showToast(`Orden actualizada: #${orderNumber}`, 'info');
                        break;
                    }
                    case 'integration.sync.order.rejected': {
                        const orderNumber = eventData.order_number || event.metadata?.order_number || 'Desconocida';
                        const reason = eventData.reason || event.metadata?.reason || 'Error desconocido';
                        const error = eventData.error || event.metadata?.error || '';
                        const createdAt = eventData.created_at || eventData.rejected_at || event.metadata?.created_at || null;
                        const orderStatus = eventData.status || event.metadata?.status || null;

                        setSyncProgress(prev => {
                            if (!prev) {
                                return { total: 1, created: 0, rejected: 1, updated: 0, orders: [{ orderNumber, status: 'rejected' as const, reason: reason + (error ? `: ${error}` : ''), createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date() }] };
                            }
                            return { ...prev, rejected: prev.rejected + 1, total: prev.total + 1, orders: [{ orderNumber, status: 'rejected' as const, reason: reason + (error ? `: ${error}` : ''), createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date() }, ...prev.orders] };
                        });
                        playNotificationSound();
                        showToast(`Orden rechazada: #${orderNumber} - ${reason}${error ? `: ${error}` : ''}`, 'error');
                        break;
                    }
                    case 'integration.sync.started': {
                        const integrationId = event.integration_id;
                        const integration = integrations.find(i => i.id === integrationId);
                        const integrationName = integration?.name || `Integración ${integrationId}`;

                        setSyncProgress({ total: 0, created: 0, rejected: 0, updated: 0, orders: [] });
                        showToast(`Sincronización iniciada: ${integrationName}`, 'info');
                        break;
                    }
                    case 'integration.sync.completed': {
                        const integrationId = event.integration_id;
                        const integration = integrations.find(i => i.id === integrationId);
                        const integrationName = integration?.name || `Integración ${integrationId}`;
                        const totalOrders = Number(eventData.total_orders) || 0;
                        const createdOrders = Number(eventData.created_orders) || 0;
                        const updatedOrders = Number(eventData.updated_orders) || 0;
                        const rejectedOrders = Number(eventData.rejected_orders) || 0;

                        setSyncProgress(prev => {
                            if (!prev) return { total: totalOrders, created: createdOrders, rejected: rejectedOrders, updated: updatedOrders, orders: [] };
                            return { ...prev, total: totalOrders, created: createdOrders, rejected: rejectedOrders, updated: updatedOrders, orders: prev.orders };
                        });
                        setSyncing(false);
                        playNotificationSound();
                        showToast(`Sincronización completada: ${integrationName} - Total: ${totalOrders}, Creadas: ${createdOrders}, Actualizadas: ${updatedOrders}, Rechazadas: ${rejectedOrders}`, 'success');
                        break;
                    }
                    case 'integration.sync.failed': {
                        const integrationId = event.integration_id;
                        const integration = integrations.find(i => i.id === integrationId);
                        const integrationName = integration?.name || `Integración ${integrationId}`;
                        const error = eventData.error || event.metadata?.error || 'Error desconocido';

                        setSyncProgress(null);
                        playNotificationSound();
                        showToast(`Sincronización fallida: ${integrationName} - ${error}`, 'error');
                        break;
                    }
                }
            } catch (err) {
                console.error('Error parsing integration SSE event:', err);
            }
        },
        onError: () => {
            console.error('Error en conexión SSE de eventos de integraciones');
        },
        onOpen: () => {
            console.log('Conexión SSE de eventos de integraciones establecida');
        },
    });

    // Filtros dinámicos - sincronizar con el hook
    const [filters, setFilters] = useState<{
        search?: string;
        type?: string;
    }>({});

    // Definir filtros disponibles
    const availableFilters: FilterOption[] = useMemo(() => {
        return [
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
        ];
    }, [integrationTypes]);

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
            }
            return newFilters;
        });
        setPage(1);
    }, [setSearch, setFilterType, setPage]);

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
            }
            return newFilters;
        });
        setPage(1);
    }, [setSearch, setFilterType, setPage]);

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

        // Consultar si hay una sincronización en curso
        try {
            const { getSyncStatusAction } = await import('../../infra/actions');
            const token = TokenStorage.getSessionToken();
            const syncStatus = await getSyncStatusAction(id, currentBusinessId, token);
            if (syncStatus.success && syncStatus.in_progress && syncStatus.sync_state) {
                // Hay una sincronización en curso, mostrar el estado actual
                setSyncing(true);
                setSyncProgress({
                    total: 0,
                    created: 0,
                    rejected: 0,
                    updated: 0,
                    orders: []
                });
                showToast('🔄 Hay una sincronización en curso. Mostrando progreso actual...', 'info');
            }
        } catch (error: any) {
            console.error('Error al consultar estado de sincronización:', error);
            // Continuar normalmente si hay error
        }

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
            // Inicializar progreso
            setSyncProgress({
                total: 0,
                created: 0,
                rejected: 0,
                updated: 0,
                orders: []
            });
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
                    // No mostrar alert, el evento SSE mostrará la notificación en tiempo real
                    showToast('🔄 Sincronización iniciada. Recibirás notificaciones en tiempo real.', 'info');
                } else {
                    showToast(`❌ Error al iniciar sincronización: ${result.message}`, 'error');
                    setSyncProgress(null);
                    setSyncing(false);
                }
            } catch (error: any) {
                showToast(`❌ Error al iniciar sincronización: ${error.message}`, 'error');
                setSyncProgress(null);
                setSyncing(false);
            }
        }
    };

    // Cerrar modal de sincronización
    const handleSyncCancel = () => {
        setSyncModal({ show: false, id: null, name: '' });
        setSyncFilters(initialSyncFilters);
        setSyncProgress(null);
        setSyncing(false);
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
        { key: 'business', label: 'Empresa' },
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
        type: (
            <div className="text-sm text-gray-700">
                {integration.integration_type?.name || (
                    <span className="text-gray-400 text-xs">—</span>
                )}
            </div>
        ),
        category: (
            <span
                className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white"
                style={{ backgroundColor: integration.category_color || '#6B7280' }}
            >
                {integration.category_name || integration.category || 'Sin categoría'}
            </span>
        ),
        business: (
            <div className="text-sm text-gray-700">
                {integration.business_name || (
                    <span className="text-gray-400 text-xs">Sin empresa</span>
                )}
            </div>
        ),
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
                    className={`p-2 rounded-md transition-colors duration-200 focus:ring-2 focus:ring-offset-2 ${integration.is_active
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
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
                    <div className="bg-white rounded-lg shadow-xl w-full max-w-md flex flex-col max-h-[90vh] overflow-hidden">
                        {/* Header fijo */}
                        <div className="px-6 py-4 border-b border-gray-200 flex-shrink-0">
                            <div className="flex items-start justify-between">
                                <div className="flex-1">
                                    <h3 className="text-lg font-semibold text-gray-900">
                                        ↻ Sincronizar Órdenes
                                    </h3>
                                    <p className="text-sm text-gray-500 mt-1">
                                        Integración: <strong>{syncModal.name}</strong>
                                    </p>
                                </div>
                                <button
                                    onClick={() => {
                                        setSyncModal({ show: false, id: null, name: '' });
                                        setSyncProgress(null);
                                        setSyncing(false);
                                    }}
                                    className="ml-4 text-gray-400 hover:text-gray-600 transition-colors"
                                    aria-label="Cerrar modal"
                                >
                                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                    </svg>
                                </button>
                            </div>
                        </div>

                        {/* Contenido scrolleable */}
                        <div className="px-6 py-4 space-y-4 overflow-y-auto flex-1 min-h-0">
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

                            {/* Barra de Progreso en Tiempo Real */}
                            {syncProgress && (
                                <div className="space-y-4 pt-4 border-t border-gray-200">
                                    <div>
                                        <div className="flex justify-between items-center mb-2">
                                            <h4 className="text-sm font-semibold text-gray-700">Progreso de Sincronización</h4>
                                            <span className="text-xs text-gray-500">
                                                {syncProgress.total > 0
                                                    ? `${syncProgress.created + syncProgress.rejected + syncProgress.updated} / ${syncProgress.total}`
                                                    : 'Iniciando...'}
                                            </span>
                                        </div>

                                        {/* Barra de progreso total con colores */}
                                        <div className="w-full bg-gray-200 rounded-full h-4 overflow-hidden flex">
                                            {syncProgress.total > 0 && (
                                                <>
                                                    {/* Porción verde (creadas) */}
                                                    {syncProgress.created > 0 && (
                                                        <div
                                                            className="h-full bg-green-500 transition-all duration-300"
                                                            style={{
                                                                width: `${(syncProgress.created / syncProgress.total) * 100}%`
                                                            }}
                                                            title={`${syncProgress.created} órdenes creadas`}
                                                        ></div>
                                                    )}
                                                    {/* Porción roja (rechazadas) */}
                                                    {syncProgress.rejected > 0 && (
                                                        <div
                                                            className="h-full bg-red-500 transition-all duration-300"
                                                            style={{
                                                                width: `${(syncProgress.rejected / syncProgress.total) * 100}%`
                                                            }}
                                                            title={`${syncProgress.rejected} órdenes rechazadas`}
                                                        ></div>
                                                    )}
                                                    {/* Porción amarilla (actualizadas) */}
                                                    {syncProgress.updated > 0 && (
                                                        <div
                                                            className="h-full bg-yellow-500 transition-all duration-300"
                                                            style={{
                                                                width: `${(syncProgress.updated / syncProgress.total) * 100}%`
                                                            }}
                                                            title={`${syncProgress.updated} órdenes actualizadas`}
                                                        ></div>
                                                    )}
                                                </>
                                            )}
                                            {syncProgress.total === 0 && (
                                                <div
                                                    className="h-full bg-blue-500 animate-pulse"
                                                    style={{ width: '10%' }}
                                                ></div>
                                            )}
                                        </div>

                                        {/* Estadísticas */}
                                        <div className="flex gap-4 mt-3 text-xs">
                                            <div className="flex items-center gap-1">
                                                <div className="w-3 h-3 rounded-full bg-green-500"></div>
                                                <span className="text-gray-600">Creadas: <strong>{syncProgress.created}</strong></span>
                                            </div>
                                            <div className="flex items-center gap-1">
                                                <div className="w-3 h-3 rounded-full bg-red-500"></div>
                                                <span className="text-gray-600">Rechazadas: <strong>{syncProgress.rejected}</strong></span>
                                            </div>
                                            {syncProgress.updated > 0 && (
                                                <div className="flex items-center gap-1">
                                                    <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                                                    <span className="text-gray-600">Actualizadas: <strong>{syncProgress.updated}</strong></span>
                                                </div>
                                            )}
                                        </div>
                                    </div>

                                    {/* Lista de órdenes procesadas */}
                                    {syncProgress.orders.length > 0 && (
                                        <div className="max-h-64 overflow-y-auto border border-gray-200 rounded-lg">
                                            <div className="p-2 bg-gray-50 border-b border-gray-200 sticky top-0">
                                                <p className="text-xs font-semibold text-gray-700">Órdenes Procesadas</p>
                                            </div>
                                            <div className="divide-y divide-gray-100">
                                                {syncProgress.orders.slice(0, 50).map((order, index) => (
                                                    <div
                                                        key={index}
                                                        className={`p-3 text-xs ${order.status === 'created'
                                                            ? 'bg-green-50 hover:bg-green-100'
                                                            : order.status === 'updated'
                                                                ? 'bg-blue-50 hover:bg-blue-100'
                                                                : 'bg-red-50 hover:bg-red-100'
                                                            } transition-colors`}
                                                    >
                                                        <div className="flex items-start justify-between gap-2">
                                                            <div className="flex-1 min-w-0">
                                                                <div className="flex items-center gap-2 mb-1">
                                                                    <div className={`w-2 h-2 rounded-full flex-shrink-0 ${order.status === 'created' ? 'bg-green-500'
                                                                        : order.status === 'updated' ? 'bg-blue-500'
                                                                            : 'bg-red-500'
                                                                        }`}></div>
                                                                    <span className="font-medium text-gray-800">
                                                                        #{order.orderNumber}
                                                                    </span>
                                                                    {order.orderStatus && (
                                                                        <span className="px-2 py-0.5 bg-gray-200 text-gray-700 rounded text-xs font-medium">
                                                                            {order.orderStatus}
                                                                        </span>
                                                                    )}
                                                                </div>
                                                                {order.createdAt && (
                                                                    <div className="text-gray-600 text-xs ml-4">
                                                                        Creada: {(() => {
                                                                            try {
                                                                                const date = new Date(order.createdAt);
                                                                                return isNaN(date.getTime())
                                                                                    ? order.createdAt
                                                                                    : date.toLocaleString('es-ES', {
                                                                                        year: 'numeric',
                                                                                        month: '2-digit',
                                                                                        day: '2-digit',
                                                                                        hour: '2-digit',
                                                                                        minute: '2-digit'
                                                                                    });
                                                                            } catch {
                                                                                return order.createdAt;
                                                                            }
                                                                        })()}
                                                                    </div>
                                                                )}
                                                                {order.status === 'rejected' && order.reason && (
                                                                    <div className="text-red-600 text-xs ml-4 mt-1">
                                                                        {order.reason}
                                                                    </div>
                                                                )}
                                                            </div>
                                                            <span className="text-gray-400 text-xs flex-shrink-0">
                                                                {order.timestamp.toLocaleTimeString()}
                                                            </span>
                                                        </div>
                                                    </div>
                                                ))}
                                            </div>
                                            {syncProgress.orders.length > 50 && (
                                                <div className="p-2 bg-gray-50 text-xs text-gray-500 text-center">
                                                    Y {syncProgress.orders.length - 50} órdenes más...
                                                </div>
                                            )}
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>

                        {/* Footer fijo */}
                        <div className="px-6 py-4 border-t border-gray-200 flex justify-end gap-3 flex-shrink-0 bg-white">
                            {!syncing && syncProgress && syncProgress.total > 0 && (
                                <Button
                                    variant="primary"
                                    onClick={handleSyncCancel}
                                >
                                    Cerrar
                                </Button>
                            )}
                            {!syncing && (!syncProgress || syncProgress.total === 0) && (
                                <>
                                    <Button
                                        variant="outline"
                                        onClick={handleSyncCancel}
                                    >
                                        Cancelar
                                    </Button>
                                    <Button
                                        variant="primary"
                                        onClick={handleSyncConfirm}
                                    >
                                        ↻ Iniciar Sincronización
                                    </Button>
                                </>
                            )}
                            {syncing && (
                                <Button
                                    variant="outline"
                                    onClick={handleSyncCancel}
                                    disabled={true}
                                >
                                    Sincronizando... (No cerrar)
                                </Button>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
