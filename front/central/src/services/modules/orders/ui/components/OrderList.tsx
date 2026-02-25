'use client';

import { useState, useEffect, useCallback, useMemo, useRef, memo } from 'react';
import { getOrdersAction, deleteOrderAction, getOrderByIdAction } from '../../infra/actions';
import { getIntegrationsAction } from '@/services/integrations/core/infra/actions';
import { getOrderStatusesAction } from '@/services/modules/orderstatus/infra/actions';
import { getPaymentStatusesAction } from '@/services/modules/paymentstatus/infra/actions';
import { getFulfillmentStatusesAction } from '@/services/modules/fulfillmentstatus/infra/actions';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
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
    onViewRecommendation,
    onDelete,
    onShowRaw,
    onShowGuide,
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
    onViewRecommendation?: (order: Order) => void;
    onDelete: (id: string) => void;
    onShowRaw: (id: string) => void;
    onShowGuide: (guideLink: string) => void;
    formatCurrency: (amount: number, currency?: string, amountPresentment?: number, currencyPresentment?: string) => string;
    formatDate: (dateString: string) => { date: string; time: string };
    getStatusBadge: (status: string, color?: string) => React.ReactNode;
    getProbabilityColor: (probability: number) => string;
    isNew?: boolean;
    businessesMap: Map<number, string>;
    isSuperAdmin: boolean;
}) => {
    return (
        <tr className={`transition-all duration-300 hover:bg-purple-50 cursor-pointer ${isNew ? 'animate-slide-in' : ''}`}>
            <td className="px-2 sm:px-3 py-2 whitespace-nowrap">
                <div className="flex items-center gap-1">
                    <div
                        className="h-8 w-8 rounded-full shadow-md border-2 border-gray-200 hover:shadow-lg transition-all cursor-pointer bg-white flex items-center justify-center overflow-hidden"
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
                            className="flex items-center justify-center w-8 h-8 rounded-md bg-purple-500 hover:bg-purple-600 text-white transition-colors duration-200 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
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
            <td className="px-2 sm:px-3 py-2">
                <div className="text-xs font-medium text-gray-900">
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
            <td className="px-2 sm:px-3 py-2 whitespace-nowrap">
                <div className="text-sm font-semibold text-gray-900">
                    {formatCurrency(order.total_amount, order.currency, order.total_amount_presentment, order.currency_presentment)}
                </div>
            </td>
            <td className="px-2 sm:px-3 py-2 whitespace-nowrap">
                {getStatusBadge(order.order_status?.name || order.status, order.order_status?.color) || (
                    <span className="text-xs text-gray-400">-</span>
                )}
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden lg:table-cell">
                {order.payment_status?.name ? (
                    getStatusBadge(order.payment_status.name, order.payment_status.color)
                ) : order.is_paid ? (
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        Pagado
                    </span>
                ) : (
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                        No pagado
                    </span>
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
                            className="p-2 bg-gradient-to-r from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700 text-white rounded-md transition-all duration-200 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 flex items-center justify-center shadow-sm"
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
                            className="p-2 bg-purple-500 hover:bg-purple-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
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
                            className="p-2 bg-purple-500 hover:bg-purple-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
                            title="Editar orden"
                            aria-label="Editar orden"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                            </svg>
                        </button>
                    )}
                    {order.guide_link && (
                        <button
                            onClick={() => onShowGuide(order.guide_link!)}
                            className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                            title="Ver guía de envío"
                            aria-label="Ver guía de envío"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
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
    const { businesses } = useBusinessesSimple();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [orders, setOrders] = useState<Order[]>([]);
    const [initialLoading, setInitialLoading] = useState(true);
    const [tableLoading, setTableLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const isFirstLoad = useRef(true);
    const [newOrderIds, setNewOrderIds] = useState<Set<string>>(new Set());

    // Guide Modal
    const [guideUrl, setGuideUrl] = useState<string | null>(null);

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
                    const options = response.data
                        .filter((status) => status?.name) // Filtrar registros sin nombre
                        .map((status) => ({
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
                    const options = response.data
                        .filter((status) => status?.name) // Filtrar registros sin nombre
                        .map((status) => ({
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
                    const options = response.data
                        .filter((status) => status?.name) // Filtrar registros sin nombre
                        .map((status) => ({
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
        // Super admin (business_id=0) recibe eventos de todos los businesses
        businessId: isSuperAdmin ? 0 : permissions?.business_id,
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
            const params: GetOrdersParams = { ...filters };
            if (isSuperAdmin && selectedBusinessId !== null) {
                params.business_id = selectedBusinessId;
            }
            const response = await getOrdersAction(params);
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
    }, [filters, isSuperAdmin, selectedBusinessId]);

    // Función unificada para cargar órdenes
    const loadOrders = useCallback(async (showInitialLoading = false) => {
        if (showInitialLoading) {
            setInitialLoading(true);
        } else {
            setTableLoading(true);
        }
        setError(null);
        try {
            const params: GetOrdersParams = { ...filters };
            if (isSuperAdmin && selectedBusinessId !== null) {
                params.business_id = selectedBusinessId;
            }
            const response = await getOrdersAction(params);
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
    }, [filters, isSuperAdmin, selectedBusinessId]);

    // Reset a página 1 cuando el super admin cambia el negocio seleccionado
    useEffect(() => {
        if (isSuperAdmin) {
            setFilters(prev => ({ ...prev, page: 1 }));
        }
    }, [selectedBusinessId, isSuperAdmin]);

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
            {/* Business Selector - Solo Super Admin */}
            {isSuperAdmin && businesses.length > 0 && (
                <div className="mb-4 bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Seleccionar Negocio (Super Admin)
                    </label>
                    <select
                        value={selectedBusinessId?.toString() ?? ''}
                        onChange={(e) => {
                            const val = e.target.value;
                            setSelectedBusinessId(val ? Number(val) : null);
                        }}
                        className="w-full max-w-xs px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                        <option value="">Todos los negocios</option>
                        {businesses.map((b) => (
                            <option key={b.id} value={b.id}>{b.name} (ID: {b.id})</option>
                        ))}
                    </select>
                </div>
            )}

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
                        sortOptions={[
                            { value: 'created_at', label: 'Ordenar por fecha' },
                            { value: 'updated_at', label: 'Ordenar por actualización' },
                            { value: 'total_amount', label: 'Ordenar por monto' },
                            { value: 'order_number', label: 'Ordenar por ID' },
                        ]}
                    />
                </div>
            </div>

            {/* Table */}
            <div className="ordersTable relative">
                {/* Overlay de carga solo para la tabla */}
                {tableLoading && (
                    <div className="absolute inset-0 bg-white/80 backdrop-blur-sm z-10 flex items-center justify-center transition-opacity duration-200">
                        <div className="flex flex-col items-center gap-2">
                            <div className="w-8 h-8 border-4 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
                            <p className="text-sm text-gray-600">Actualizando...</p>
                        </div>
                    </div>
                )}
                <div className="overflow-x-auto">
                    <table className={`min-w-full transition-opacity duration-200 ${tableLoading ? 'opacity-50' : 'opacity-100'}`}>
                        <thead style={{ background: 'linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%)' }}>
                            <tr>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest w-16" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)', borderTopLeftRadius: '14px', borderBottomLeftRadius: '14px' }}>
                                    {/* Columna del logo - sin título */}
                                </th>
                                <th
                                    className="px-3s sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest cursor-pointer transition-all group"
                                    onClick={() => handleSortChange('order_number', filters.sort_order === 'asc' ? 'desc' : 'asc')}
                                    style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}
                                >
                                    <div className="flex items-center gap-1">
                                        Orden
                                        {filters.sort_by === 'order_number' && (
                                            <span className="text-purple-100">
                                                {filters.sort_order === 'asc' ? '↑' : '↓'}
                                            </span>
                                        )}
                                    </div>
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden sm:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}>
                                    Cliente
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}>
                                    Total
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}>
                                    Estado
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden lg:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}>
                                    Estatus Pago
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden md:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}>
                                    Probabilidad
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden lg:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}>
                                    Datos Faltantes
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-center text-xs font-bold text-white uppercase tracking-widest" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}>
                                    Confirmado
                                </th>

                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden md:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}>
                                    Fecha
                                </th>
                                {isSuperAdmin && (
                                    <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden lg:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)' }}>
                                        Business
                                    </th>
                                )}
                                <th className="px-3 sm:px-6 py-3 text-right text-xs font-bold text-white uppercase tracking-widest" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em', boxShadow: '0 10px 25px rgba(124, 58, 237, 0.18)', borderTopRightRadius: '14px', borderBottomRightRadius: '14px' }}>
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody>
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
                                        onShowGuide={(link) => setGuideUrl(link)}
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
                    <div className="px-3 sm:px-4 lg:px-6 py-3 flex flex-col sm:flex-row items-center justify-between gap-3">
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
                                <div className="flex items-center gap-1">
                                    <label className="text-xs sm:text-sm text-gray-700 whitespace-nowrap">
                                        Mostrar:
                                    </label>
                                    <select
                                        value={filters.page_size || 20}
                                        onChange={(e) => {
                                            const newPageSize = parseInt(e.target.value);
                                            setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                        }}
                                        className="px-2 py-1.5 text-xs sm:text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 bg-white"
                                    >
                                        <option value="10">10</option>
                                        <option value="20">20</option>
                                        <option value="50">50</option>
                                        <option value="100">100</option>
                                    </select>
                                </div>
                            </div>
                            <div className="flex items-center gap-1">
                                <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page - 1 })}
                                        disabled={page === 1}
                                        className="relative inline-flex items-center px-3 sm:px-4 py-2 rounded-l-md border-2 border-purple-200 bg-white text-xs sm:text-sm font-medium text-purple-700 hover:bg-purple-50 hover:border-purple-300 disabled:opacity-50 transition-all"
                                    >
                                        Anterior
                                    </button>
                                    <span className="relative inline-flex items-center px-3 sm:px-4 py-2 border-2 border-purple-200 bg-purple-50 text-xs sm:text-sm font-semibold text-purple-900">
                                        Página {page} de {totalPages}
                                    </span>
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page + 1 })}
                                        disabled={page === totalPages}
                                        className="relative inline-flex items-center px-3 sm:px-4 py-2 rounded-r-md border-2 border-purple-200 bg-white text-xs sm:text-sm font-medium text-purple-700 hover:bg-purple-50 hover:border-purple-300 disabled:opacity-50 transition-all"
                                    >
                                        Siguiente
                                    </button>
                                </nav>
                            </div>
                        </div>

                        {/* Mobile: Page size selector */}
                        <div className="flex items-center justify-between w-full sm:hidden pt-2">
                            <div className="flex items-center gap-1">
                                <label className="text-xs text-gray-700 whitespace-nowrap">
                                    Mostrar:
                                </label>
                                <select
                                    value={filters.page_size || 20}
                                    onChange={(e) => {
                                        const newPageSize = parseInt(e.target.value);
                                        setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                    }}
                                    className="px-2 py-1.5 text-xs border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 bg-white"
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

                <style jsx>{`
                    /* Tabla más "card-like" fila por fila */
                    .ordersTable :global(.table) {
                        border-collapse: separate;
                        border-spacing: 0 10px; /* separación entre filas */
                        background: transparent;
                    }

                    /* Quitar el borde del contenedor global de Table SOLO aquí */
                    .ordersTable :global(div.overflow-hidden.w-full.rounded-lg.border.border-gray-200.bg-white) {
                        border: none !important;
                        background: transparent !important;
                    }

                    .ordersTable :global(.table th) {
                        background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%);
                        color: #fff;
                        position: sticky;
                        top: 0;
                        z-index: 1;
                    }

                    /* Header más llamativo + bordes redondeados */
                    .ordersTable :global(.table thead th) {
                        padding-top: 10px;
                        padding-bottom: 10px;
                        font-size: 0.75rem; /* más pequeño */
                        font-weight: 800;
                        letter-spacing: 0.06em;
                        text-transform: uppercase;
                        box-shadow: 0 10px 25px rgba(124, 58, 237, 0.18);
                    }

                    .ordersTable :global(.table thead th:first-child) {
                        border-top-left-radius: 14px;
                        border-bottom-left-radius: 14px;
                    }

                    .ordersTable :global(.table thead th:last-child) {
                        border-top-right-radius: 14px;
                        border-bottom-right-radius: 14px;
                    }

                    .ordersTable :global(.table tbody tr) {
                        background: rgba(255, 255, 255, 0.95);
                        box-shadow: 0 1px 0 rgba(17, 24, 39, 0.04);
                        transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
                    }

                    /* Zebra suave en morado */
                    .ordersTable :global(.table tbody tr:nth-child(even)) {
                        background: rgba(124, 58, 237, 0.03);
                    }

                    .ordersTable :global(.table tbody tr:hover) {
                        background: rgba(124, 58, 237, 0.06);
                        box-shadow: 0 10px 25px rgba(17, 24, 39, 0.08);
                        transform: translateY(-1px);
                    }

                    .ordersTable :global(.table td) {
                        border-top: none;
                    }

                    /* Redondeo de cada fila */
                    .ordersTable :global(.table tbody td:first-child) {
                        border-top-left-radius: 12px;
                        border-bottom-left-radius: 12px;
                    }
                    .ordersTable :global(.table tbody td:last-child) {
                        border-top-right-radius: 12px;
                        border-bottom-right-radius: 12px;
                    }

                    /* Acciones: focus consistente */
                    .ordersTable :global(a),
                    .ordersTable :global(button) {
                        outline-color: rgba(124, 58, 237, 0.35);
                    }
                `}</style>
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

            {/* Guide Modal */}
            {guideUrl && (
                <div
                    className="fixed inset-0 z-50 flex items-center justify-center"
                    onClick={() => setGuideUrl(null)}
                >
                    <div className="absolute inset-0 bg-black/60" />
                    <div
                        className="relative z-10 bg-white rounded-xl shadow-2xl flex flex-col overflow-hidden"
                        style={{ width: '480px', maxWidth: '90vw' }}
                        onClick={(e) => e.stopPropagation()}
                    >
                        {/* Header */}
                        <div className="flex items-center justify-between px-5 py-4 border-b bg-gray-50">
                            <div className="flex items-center gap-2 text-sm font-semibold text-gray-800">
                                <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                </svg>
                                Guía de envío
                            </div>
                            <button
                                onClick={() => setGuideUrl(null)}
                                className="p-1.5 text-gray-400 hover:text-red-500 transition-colors rounded-full hover:bg-red-50"
                                aria-label="Cerrar"
                            >
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                </svg>
                            </button>
                        </div>
                        {/* PDF download/open UI — iframes can't embed S3 PDFs reliably */}
                        <div className="p-8 flex flex-col items-center gap-6">
                            <div className="w-20 h-20 rounded-2xl bg-red-50 flex items-center justify-center">
                                <svg className="w-10 h-10 text-red-500" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zm-1 1.5L18.5 9H13V3.5zM6 4h5v7h7v9H6V4z"/>
                                    <path d="M8 12h8v1.5H8V12zm0 3h8v1.5H8V15zm0 3h5v1.5H8V18z"/>
                                </svg>
                            </div>
                            <div className="text-center">
                                <p className="font-semibold text-gray-800 text-lg">Guía de Envío lista</p>
                                <p className="text-sm text-gray-500 mt-1">El PDF está disponible para ver o descargar</p>
                            </div>
                            <div className="flex flex-col gap-3 w-full">
                                <a
                                    href={guideUrl}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="flex items-center justify-center gap-2 w-full py-3 px-6 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition-colors shadow-sm"
                                >
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                                    </svg>
                                    Abrir PDF en nueva pestaña
                                </a>
                                <a
                                    href={guideUrl}
                                    download
                                    className="flex items-center justify-center gap-2 w-full py-3 px-6 bg-gray-100 hover:bg-gray-200 text-gray-700 font-semibold rounded-xl transition-colors"
                                >
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                                    </svg>
                                    Descargar PDF
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
