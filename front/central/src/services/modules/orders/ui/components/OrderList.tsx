'use client';

import { useState, useEffect, useCallback, useMemo, useRef, memo } from 'react';
import { getOrdersAction, deleteOrderAction, getOrderByIdAction } from '../../infra/actions';
import { getIntegrationsAction } from '@/services/integrations/core/infra/actions';
import { getOrderStatusesAction } from '@/services/modules/orderstatus/infra/actions';
import { getPaymentStatusesAction } from '@/services/modules/paymentstatus/infra/actions';
import { getFulfillmentStatusesAction } from '@/services/modules/fulfillmentstatus/infra/actions';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import { Order, GetOrdersParams } from '../../domain/types';
import { Button, Alert, DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { useSSE } from '@/shared/hooks/use-sse';
import { useToast } from '@/shared/providers/toast-provider';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { playNotificationSound } from '@/shared/utils';
import RawOrderModal from './RawOrderModal';

// Componente memoizado para las filas de la tabla
const OrderRow = memo(({
    order,
    onView,
    onEdit,
    onViewRecommendation, // NEW prop
    onDelete,
    onShowRaw,
    formatCurrency,
    formatDate,
    getStatusBadge,
    getProbabilityColor,
    isNew,
    businessesMap,
    isSuperAdmin
}: {
    order: Order;
    onView?: (order: Order) => void;
    onEdit?: (order: Order) => void;
    onViewRecommendation?: (order: Order) => void; // NEW prop definition
    onDelete: (id: string) => void;
    onShowRaw: (id: string) => void;
    formatCurrency: (amount: number, currency?: string, amountPresentment?: number, currencyPresentment?: string) => string;
    formatDate: (dateString: string) => { date: string; time: string };
    getStatusBadge: (status: string, color?: string) => React.ReactNode;
    getProbabilityColor: (probability: number) => string;
    isNew?: boolean;
    businessesMap: Map<number, string>;
    isSuperAdmin: boolean;
}) => {
    return (
        <tr className={`hover:bg-gray-50 transition-all duration-300 ${isNew ? 'animate-slide-in bg-green-50/50' : ''}`}>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                <div className="flex items-center gap-2">
                    <div
                        className="h-10 w-10 rounded-full shadow-md border-2 border-gray-200 hover:shadow-lg transition-all cursor-pointer bg-white flex items-center justify-center overflow-hidden"
                        title={`${order.platform} - Click para ver JSON crudo`}
                        onClick={() => onShowRaw(order.id)}
                    >
                        {order.integration_logo_url ? (
                            <img
                                src={order.integration_logo_url}
                                alt={order.platform}
                                className="h-full w-full object-contain p-1.5"
                                loading="lazy"
                            />
                        ) : (
                            <span className="text-xs font-medium text-gray-600 uppercase">
                                {order.platform.charAt(0)}
                            </span>
                        )}
                    </div>
                    {order.order_status_url && (
                        <a
                            href={order.order_status_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center justify-center w-8 h-8 rounded-md bg-blue-500 hover:bg-blue-600 text-white transition-colors duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                            title="Ver orden en Shopify"
                            aria-label="Ver orden en Shopify"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                            </svg>
                        </a>
                    )}
                </div>
            </td>
            <td className="px-3 sm:px-6 py-4">
                <div className="text-sm font-medium text-gray-900">
                    {order.order_number || order.external_id || order.id}
                </div>
                <div className="text-xs text-gray-500 sm:hidden">
                    {order.customer_name}
                </div>
                {order.internal_number && (
                    <div className="text-xs text-gray-500">
                        Interno: {order.internal_number}
                    </div>
                )}
                {order.external_id && order.order_number && order.external_id !== order.order_number && (
                    <div className="text-xs text-gray-400">
                        Ext: {order.external_id}
                    </div>
                )}
            </td>
            <td className="px-3 sm:px-6 py-4 hidden sm:table-cell">
                <div className="text-sm text-gray-900">{order.customer_name}</div>
                <div className="text-xs text-gray-500">{order.customer_email}</div>
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                <div className="text-sm font-semibold text-gray-900">
                    {formatCurrency(order.total_amount, order.currency, order.total_amount_presentment, order.currency_presentment)}
                </div>
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                {getStatusBadge(order.order_status?.name || order.status, order.order_status?.color) || (
                    <span className="text-xs text-gray-400">-</span>
                )}
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden lg:table-cell">
                {order.payment_status?.name ? (
                    getStatusBadge(order.payment_status.name, order.payment_status.color)
                ) : (
                    <span className="text-xs text-gray-400">-</span>
                )}
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden md:table-cell">
                {order.delivery_probability !== undefined && order.delivery_probability !== null ? (
                    <div className="flex items-center gap-2 min-w-[120px]">
                        <div className="flex-1 bg-gray-200 rounded-full h-3 overflow-hidden shadow-inner">
                            <div
                                className={`h-full rounded-full transition-all duration-300 ${getProbabilityColor(order.delivery_probability)}`}
                                style={{ width: `${Math.min(order.delivery_probability, 100)}%` }}
                                title={`Probabilidad de entrega: ${order.delivery_probability.toFixed(1)}%`}
                            ></div>
                        </div>
                        <span className="text-xs font-semibold text-gray-700 min-w-[40px] text-right">
                            {order.delivery_probability.toFixed(0)}%
                        </span>
                    </div>
                ) : (
                    <span className="text-xs text-gray-400">N/A</span>
                )}
            </td>
            <td className="px-3 sm:px-6 py-4 hidden lg:table-cell">
                <div className="flex flex-wrap gap-1">
                    {order.negative_factors && order.negative_factors.length > 0 ? (
                        order.negative_factors.map((factor, idx) => (
                            <span key={idx} className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">
                                {factor}
                            </span>
                        ))
                    ) : (
                        <span className="text-xs text-gray-400">-</span>
                    )}
                </div>
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-center">
                <div className="flex flex-col items-center gap-1">
                    {order.is_confirmed === true ? (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                            Sí
                        </span>
                    ) : order.is_confirmed === false ? (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                            No
                        </span>
                    ) : (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                            Pendiente
                        </span>
                    )}
                    {order.is_confirmed === true && order.novelty && (
                        <span className="text-[10px] font-medium text-orange-600 leading-tight" title={order.novelty}>
                            Novedad encontrada
                        </span>
                    )}
                </div>
            </td>

            <td className="px-3 sm:px-6 py-4 text-xs text-gray-500 hidden md:table-cell">
                <div className="leading-tight">
                    <div className="text-gray-900">{formatDate(order.created_at).date}</div>
                    <div className="text-gray-500">{formatDate(order.created_at).time}</div>
                </div>
            </td>
            {isSuperAdmin && (
                <td className="px-3 sm:px-6 py-4 hidden lg:table-cell">
                    <div className="text-sm text-gray-900">
                        {order.business_id && businessesMap.get(order.business_id)
                            ? businessesMap.get(order.business_id)
                            : order.business_id
                                ? `ID: ${order.business_id}`
                                : '-'
                        }
                    </div>
                </td>
            )}
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div className="flex flex-row justify-end gap-2">
                    {/* Botón de Recomendación Inteligente (Robot) */}
                    {onViewRecommendation && (
                        <button
                            onClick={() => onViewRecommendation(order)}
                            className="p-2 bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white rounded-md transition-all duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 flex items-center justify-center shadow-sm"
                            title="Recomendación Inteligente IA"
                            aria-label="Ver recomendación IA"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="16"
                                height="16"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                className="lucide lucide-bot"
                            >
                                <path d="M12 8V4H8" />
                                <rect width="16" height="12" x="4" y="8" rx="2" />
                                <path d="M2 14h2" />
                                <path d="M20 14h2" />
                                <path d="M15 13v2" />
                                <path d="M9 13v2" />
                            </svg>
                        </button>
                    )}

                    {onView && (
                        <button
                            onClick={() => onView(order)}
                            className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                            title="Ver orden"
                            aria-label="Ver orden"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                            </svg>
                        </button>
                    )}
                    {onEdit && (
                        <button
                            onClick={() => onEdit(order)}
                            className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                            title="Editar orden"
                            aria-label="Editar orden"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                            </svg>
                        </button>
                    )}
                    <button
                        onClick={() => onDelete(order.id)}
                        className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                        title="Eliminar orden"
                        aria-label="Eliminar orden"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
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
    onViewRecommendation?: (order: Order) => void;
    onCreate?: () => void;
    onTestGuide?: () => void;
    refreshKey?: number;
}

export default function OrderList({ onView, onEdit, onViewRecommendation, refreshKey, onCreate, onTestGuide }: OrderListProps) {
    const { isSuperAdmin, permissions } = usePermissions();
    const [orders, setOrders] = useState<Order[]>([]);
    const [initialLoading, setInitialLoading] = useState(true);
    const [tableLoading, setTableLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const isFirstLoad = useRef(true);
    const [newOrderIds, setNewOrderIds] = useState<Set<string>>(new Set());

    // Raw Data Modal
    const [selectedOrderId, setSelectedOrderId] = useState<string | null>(null);
    const [selectedOrderLogo, setSelectedOrderLogo] = useState<string | undefined>(undefined);
    const [selectedOrderPlatform, setSelectedOrderPlatform] = useState<string | undefined>(undefined);
    const [isRawModalOpen, setIsRawModalOpen] = useState(false);

    // Integrations for filter
    const [integrationsList, setIntegrationsList] = useState<{ value: string; label: string }[]>([]);
    // Order statuses for filter
    const [orderStatusesList, setOrderStatusesList] = useState<{ value: string; label: string }[]>([]);
    // Payment statuses for filter
    const [paymentStatusesList, setPaymentStatusesList] = useState<{ value: string; label: string }[]>([]);
    // Fulfillment statuses for filter
    const [fulfillmentStatusesList, setFulfillmentStatusesList] = useState<{ value: string; label: string }[]>([]);
    // Businesses for mapping (only load if super admin)
    const [businessesList, setBusinessesList] = useState<{ id: number; name: string }[]>([]);

    // Business map for quick lookup (memoized)
    const businessesMap = useMemo(() => {
        const map = new Map<number, string>();
        businessesList.forEach(business => {
            map.set(business.id, business.name);
        });
        return map;
    }, [businessesList]);

    useEffect(() => {
        const fetchIntegrations = async () => {
            try {
                const response = await getIntegrationsAction({ page: 1, page_size: 100 });
                if (response.success && response.data) {
                    const options = response.data.map((integration: any) => ({
                        value: String(integration.id),
                        label: integration.name,
                    }));
                    setIntegrationsList(options);
                }
            } catch (error) {
                console.error('Error fetching integrations for filter:', error);
            }
        };
        fetchIntegrations();
    }, []);

    useEffect(() => {
        // Solo cargar businesses si el usuario es super admin
        if (!isSuperAdmin) return;

        const fetchBusinesses = async () => {
            try {
                const response = await getBusinessesAction({ page: 1, per_page: 100 });
                if (response.success && response.data) {
                    const businesses = response.data.map((business: any) => ({
                        id: business.id,
                        name: business.name,
                    }));
                    setBusinessesList(businesses);
                }
            } catch (error) {
                console.error('Error fetching businesses:', error);
            }
        };
        fetchBusinesses();
    }, [isSuperAdmin]);

    useEffect(() => {
        const fetchOrderStatuses = async () => {
            try {
                const response = await getOrderStatusesAction(true); // Solo estados activos
                if (response.success && response.data) {
                    const options = response.data.map((status) => ({
                        value: String(status.id),
                        label: status.name,
                    }));
                    // Ordenar por nombre
                    options.sort((a, b) => a.label.localeCompare(b.label));
                    setOrderStatusesList(options);
                }
            } catch (error) {
                console.error('Error fetching order statuses for filter:', error);
            }
        };
        fetchOrderStatuses();
    }, []);

    useEffect(() => {
        const fetchPaymentStatuses = async () => {
            try {
                const response = await getPaymentStatusesAction(true); // Solo estados activos
                if (response.success && response.data) {
                    const options = response.data.map((status) => ({
                        value: String(status.id),
                        label: status.name,
                    }));
                    options.sort((a, b) => a.label.localeCompare(b.label));
                    setPaymentStatusesList(options);
                }
            } catch (error) {
                console.error('Error fetching payment statuses for filter:', error);
            }
        };
        fetchPaymentStatuses();
    }, []);

    useEffect(() => {
        const fetchFulfillmentStatuses = async () => {
            try {
                const response = await getFulfillmentStatusesAction(true); // Solo estados activos
                if (response.success && response.data) {
                    const options = response.data.map((status) => ({
                        value: String(status.id),
                        label: status.name,
                    }));
                    options.sort((a, b) => a.label.localeCompare(b.label));
                    setFulfillmentStatusesList(options);
                }
            } catch (error) {
                console.error('Error fetching fulfillment statuses for filter:', error);
            }
        };
        fetchFulfillmentStatuses();
    }, []);

    // Filters
    const [filters, setFilters] = useState<GetOrdersParams>({
        page: 1,
        page_size: 20,
    });

    const { showToast } = useToast();

    // Definir filtros disponibles (usar useMemo para actualizar cuando cambian las listas)
    const availableFilters: FilterOption[] = useMemo(() => [
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
            options: orderStatusesList,
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
            key: 'integration_id',
            label: 'Integración',
            type: 'select',
            options: integrationsList,
        },
        {
            key: 'is_paid',
            label: 'Estado de pago (boolean)',
            type: 'boolean',
        },
        {
            key: 'is_cod',
            label: 'Contra Entrega',
            type: 'boolean',
        },
        {
            key: 'payment_status_id',
            label: 'Estado de pago',
            type: 'select',
            options: paymentStatusesList,
        },
        {
            key: 'fulfillment_status_id',
            label: 'Estado de fulfillment',
            type: 'select',
            options: fulfillmentStatusesList,
        },
        {
            key: 'start_date',
            label: 'Rango de fechas',
            type: 'date-range',
        },
    ], [orderStatusesList, integrationsList, paymentStatusesList, fulfillmentStatusesList]);

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
            // Buscar el label del estado seleccionado
            const statusOption = orderStatusesList.find(opt => opt.value === String(filters.status));
            active.push({
                key: 'status',
                label: 'Estado',
                value: statusOption?.label || String(filters.status),
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

        if (filters.integration_id) {
            const integration = integrationsList.find(i => i.value === String(filters.integration_id));
            active.push({
                key: 'integration_id',
                label: 'Integración',
                value: integration ? integration.label : String(filters.integration_id),
                type: 'select',
            });
        }

        if (filters.is_paid !== undefined) {
            active.push({
                key: 'is_paid',
                label: 'Estado de pago (boolean)',
                value: filters.is_paid,
                type: 'boolean',
            });
        }

        if (filters.is_cod !== undefined) {
            active.push({
                key: 'is_cod',
                label: 'Contra Entrega',
                value: filters.is_cod,
                type: 'boolean',
            });
        }

        if (filters.payment_status_id) {
            const paymentStatus = paymentStatusesList.find(s => s.value === String(filters.payment_status_id));
            active.push({
                key: 'payment_status_id',
                label: 'Estado de pago',
                value: paymentStatus ? paymentStatus.label : String(filters.payment_status_id),
                type: 'select',
            });
        }

        if (filters.fulfillment_status_id) {
            const fulfillmentStatus = fulfillmentStatusesList.find(s => s.value === String(filters.fulfillment_status_id));
            active.push({
                key: 'fulfillment_status_id',
                label: 'Estado de fulfillment',
                value: fulfillmentStatus ? fulfillmentStatus.label : String(filters.fulfillment_status_id),
                type: 'select',
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
    }, [filters, orderStatusesList, integrationsList, paymentStatusesList, fulfillmentStatusesList]);

    // Manejar adición de filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev, page: 1 };

            if (filterKey === 'start_date' && typeof value === 'object') {
                newFilters.start_date = value.start;
                newFilters.end_date = value.end;
            } else if (filterKey === 'is_paid') {
                newFilters.is_paid = value === true;
            } else if (filterKey === 'is_cod') {
                newFilters.is_cod = value === true;
            } else if (filterKey === 'integration_id') {
                newFilters.integration_id = Number(value);
            } else if (filterKey === 'payment_status_id') {
                newFilters.payment_status_id = Number(value);
            } else if (filterKey === 'fulfillment_status_id') {
                newFilters.fulfillment_status_id = Number(value);
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
            sort_by: sortBy as 'created_at' | 'updated_at' | 'total_amount' | 'order_number',
            sort_order: sortOrder,
            page: 1,
        }));
    }, []);

    // SSE Integration - Agregar nueva orden sin recargar toda la tabla
    useSSE({
        eventTypes: ['order.created'],
        onMessage: async (event) => {
            try {
                const data = JSON.parse(event.data);
                // El evento SSE tiene order_id en data.data o data.metadata
                const orderId = data.data?.order_id || data.metadata?.order_id;
                const orderNumber = data.data?.order_number || 'Desconocida';

                if (data.type === 'order.created' && orderId) {
                    // Obtener la orden completa
                    try {
                        const response = await getOrderByIdAction(orderId);
                        if (response.success && response.data) {
                            const newOrder = response.data;

                            // Security Filter: Ignore orders from other businesses for non-super admins
                            if (!isSuperAdmin && permissions?.business_id && newOrder.business_id !== permissions.business_id) {
                                return;
                            }

                            // Optional: Filter for Super Admin if business filter is active (assuming filters.business_id exists)
                            // if (isSuperAdmin && filters.business_id && newOrder.business_id !== Number(filters.business_id)) {
                            //    return;
                            // }

                            // Verificar que la orden no esté ya en la lista
                            setOrders(prevOrders => {
                                if (prevOrders.some(o => o.id === newOrder.id)) {
                                    return prevOrders; // Ya existe, no hacer nada
                                }

                                // Marcar como nueva para la animación
                                setNewOrderIds(prev => new Set(prev).add(newOrder.id));

                                // Remover el flag de "nueva" después de la animación
                                setTimeout(() => {
                                    setNewOrderIds(prev => {
                                        const updated = new Set(prev);
                                        updated.delete(newOrder.id);
                                        return updated;
                                    });
                                }, 2000);

                                // Agregar al principio de la lista
                                return [newOrder, ...prevOrders];
                            });

                            // Actualizar el total solo si realmente agregamos una nueva orden
                            setTotal(prev => prev + 1);

                            // Reproducir sonido de notificación
                            playNotificationSound();

                            // Mostrar toast
                            showToast(`Nueva orden recibida: #${orderNumber}`, 'success');
                        }
                    } catch (err) {
                        console.error('Error al obtener orden completa:', err);
                        // Si falla, recargar la tabla como fallback
                        refreshTableOnly();
                    }
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

    const formatCurrency = useCallback((amount: number, currency: string = 'USD', amountPresentment?: number, currencyPresentment?: string) => {
        // Priorizar moneda local (presentment) si está disponible
        if (amountPresentment && amountPresentment > 0 && currencyPresentment) {
            return new Intl.NumberFormat('es-CO', {
                style: 'currency',
                currency: currencyPresentment,
            }).format(amountPresentment);
        }
        // Fallback a USD si no hay moneda local
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency,
        }).format(amount);
    }, []);

    const formatDate = useCallback((dateString: string): { date: string; time: string } => {
        const date = new Date(dateString);
        // Formato compacto: "20 dic 2025" (sin "de")
        const parts = date.toLocaleDateString('es-CO', {
            day: 'numeric',
            month: 'short',
            year: 'numeric',
        }).split(' ');
        const dateStr = `${parts[0]} ${parts[2]} ${parts[parts.length - 1]}`; // "20 dic 2025"
        // Hora compacta: "12:07 a.m."
        const timeStr = date.toLocaleTimeString('es-CO', {
            hour: '2-digit',
            minute: '2-digit',
            hour12: true,
        }).replace(/\s/g, '');
        return { date: dateStr, time: timeStr };
    }, []);

    const getStatusBadge = useCallback((status: string | undefined | null, color?: string) => {
        // Si no hay status o no es un string válido, retornar null
        if (!status || typeof status !== 'string' || status.trim() === '') {
            return null;
        }

        // Si hay color configurado en la BD, usarlo
        if (color) {
            // Convertir hex a RGB para calcular si es claro u oscuro
            const hex = color.replace('#', '');
            const r = parseInt(hex.substr(0, 2), 16);
            const g = parseInt(hex.substr(2, 2), 16);
            const b = parseInt(hex.substr(4, 2), 16);
            // Calcular luminosidad
            const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
            const textColor = luminance > 0.5 ? '#000000' : '#FFFFFF';

            return (
                <span
                    className="px-2 py-1 text-xs font-medium rounded-full"
                    style={{
                        backgroundColor: color,
                        color: textColor
                    }}
                >
                    {status}
                </span>
            );
        }

        // Fallback a colores por defecto si no hay color configurado
        const statusColors: Record<string, string> = {
            pending: 'bg-yellow-100 text-yellow-800',
            processing: 'bg-blue-100 text-blue-800',
            shipped: 'bg-purple-100 text-purple-800',
            delivered: 'bg-green-100 text-green-800',
            cancelled: 'bg-red-100 text-red-800',
        };

        const colorClass = statusColors[status.toLowerCase()] || 'bg-gray-100 text-gray-800';

        return (
            <span className={`px-2 py-1 text-xs font-medium rounded-full ${colorClass}`}>
                {status}
            </span>
        );
    }, []);

    const getProbabilityColor = useCallback((probability: number) => {
        if (probability >= 80) return 'bg-green-500';
        if (probability >= 70) return 'bg-green-400';
        if (probability >= 60) return 'bg-yellow-400';
        if (probability >= 50) return 'bg-yellow-500';
        if (probability >= 40) return 'bg-orange-500';
        if (probability >= 30) return 'bg-orange-600';
        return 'bg-red-500';
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
            <div className="flex items-start justify-between gap-4 mb-4">
                <div className="flex-1">
                    <DynamicFilters
                        availableFilters={availableFilters}
                        activeFilters={activeFilters}
                        onAddFilter={handleAddFilter}
                        onRemoveFilter={handleRemoveFilter}
                        sortBy={filters.sort_by || 'created_at'}
                        sortOrder={filters.sort_order || 'desc'}
                        onSortChange={handleSortChange}
                        onCreate={onCreate}
                        onTestGuide={onTestGuide}
                        sortOptions={[
                            { value: 'created_at', label: 'Ordenar por fecha' },
                            { value: 'updated_at', label: 'Ordenar por actualización' },
                            { value: 'total_amount', label: 'Ordenar por monto' },
                            { value: 'order_number', label: 'Ordenar por ID' },
                        ]}
                    />
                </div>
                {/* + Crear Orden button rendered inside DynamicFilters */}
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
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">
                                    {/* Columna del logo - sin título */}
                                </th>
                                <th
                                    className="px-3s sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors group"
                                    onClick={() => handleSortChange('order_number', filters.sort_order === 'asc' ? 'desc' : 'asc')}
                                >
                                    <div className="flex items-center gap-1">
                                        Orden
                                        {filters.sort_by === 'order_number' && (
                                            <span className="text-gray-400">
                                                {filters.sort_order === 'asc' ? '↑' : '↓'}
                                            </span>
                                        )}
                                    </div>
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
                                    Estado de Pago
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden md:table-cell">
                                    Probabilidad
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden lg:table-cell">
                                    Datos Faltantes
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Confirmado
                                </th>

                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden md:table-cell">
                                    Fecha
                                </th>
                                {isSuperAdmin && (
                                    <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden lg:table-cell">
                                        Business
                                    </th>
                                )}
                                <th className="px-3 sm:px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {orders.length === 0 ? (
                                <tr>
                                    <td colSpan={isSuperAdmin ? 11 : 10} className="px-4 sm:px-6 py-8 text-center text-gray-500">
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
                                        onViewRecommendation={onViewRecommendation}
                                        onDelete={handleDelete}
                                        onShowRaw={(id) => {
                                            setSelectedOrderId(id);
                                            setIsRawModalOpen(true);
                                        }}
                                        formatCurrency={formatCurrency}
                                        formatDate={formatDate}
                                        getStatusBadge={getStatusBadge}
                                        getProbabilityColor={getProbabilityColor}
                                        isNew={newOrderIds.has(order.id)}
                                        businessesMap={businessesMap}
                                        isSuperAdmin={isSuperAdmin}
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
                        setSelectedOrderLogo(undefined);
                        setSelectedOrderPlatform(undefined);
                    }}
                    integrationLogoUrl={selectedOrderLogo}
                    platform={selectedOrderPlatform}
                />
            )}
        </div>
    );
}
