'use client';

import { useState, useEffect, useCallback, useMemo, useRef, memo } from 'react';
import { createPortal } from 'react-dom';
import * as XLSX from 'xlsx';
import { getOrdersAction, deleteOrderAction, getOrderByIdAction } from '../../infra/actions';
import { getIntegrationsAction } from '@/services/integrations/core/infra/actions';
import { getOrderStatusesAction } from '@/services/modules/orderstatus/infra/actions';
import { getPaymentStatusesAction } from '@/services/modules/paymentstatus/infra/actions';
import { getFulfillmentStatusesAction } from '@/services/modules/fulfillmentstatus/infra/actions';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';

import { Order, GetOrdersParams } from '../../domain/types';
import { ProbabilityTooltip } from './ProbabilityTooltip';
import { NotificationBadge } from './NotificationBadge';
import { OrderNotificationsModal } from './OrderNotificationsModal';
import { Button, Alert, DynamicFilters, FilterOption, ActiveFilter, IconActionButton, ORDERS_FILTERS_SLOT_ID, ORDERS_ACTIONS_SLOT_ID } from '@/shared/ui';
import { ChartBarIcon, ArrowPathRoundedSquareIcon, ArrowDownTrayIcon } from '@heroicons/react/24/outline';
import { useSSE } from '@/shared/hooks/use-sse';
import { useToast } from '@/shared/providers/toast-provider';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { playNotificationSound } from '@/shared/utils';
import RawOrderModal from './RawOrderModal';
import { resolveCityState } from '@/shared/utils/dane-lookup';
import DownloadOrdersModal from '@/shared/ui/modals/download-orders-modal';
import { OrderStatusFlowModal } from './OrderStatusFlowModal';
import { ChangeStatusModal } from './ChangeStatusModal';
import { QuotationExpresModal } from './QuotationExpresModal';
import { isTerminalStatus } from '../../domain/order-status-transitions';
import { getActionError } from '@/shared/utils/action-result';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';

const COMBINING_MARK_MIN = 0x0300;
const COMBINING_MARK_MAX = 0x036f;

function normalizeCarrierLabel(value: string): string {
    const decomposed = value.normalize('NFD');
    let out = '';
    for (let i = 0; i < decomposed.length; i++) {
        const code = decomposed.charCodeAt(i);
        if (code < COMBINING_MARK_MIN || code > COMBINING_MARK_MAX) out += decomposed[i];
    }
    return out.toUpperCase().replace(/[\s\-_]/g, '');
}

const SHIPMENT_STATUS_CONFIG: Record<string, { label: string; className: string }> = {
    pending: { label: 'Pendiente', className: 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300 border-amber-200 dark:border-amber-700' },
    picked_up: { label: 'Recolectado', className: 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300 border-indigo-200 dark:border-indigo-700' },
    in_transit: { label: 'En transito', className: 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 border-blue-200 dark:border-blue-700' },
    out_for_delivery: { label: 'En reparto', className: 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 border-purple-200 dark:border-purple-700' },
    delivered: { label: 'Entregado', className: 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-700' },
    on_hold: { label: 'Novedad', className: 'bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-300 border-orange-200 dark:border-orange-700' },
    returned: { label: 'Devuelto', className: 'bg-rose-100 dark:bg-rose-900/30 text-rose-700 dark:text-rose-300 border-rose-200 dark:border-rose-700' },
    failed: { label: 'Fallido', className: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 border-red-200 dark:border-red-700' },
    cancelled: { label: 'Cancelado', className: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 border-red-200 dark:border-red-700' },
};

const SHIPMENT_STATUS_FALLBACK = {
    className: 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 border-gray-200 dark:border-gray-600',
};

const STATUS_SOURCE_STYLES: Record<string, string> = {
    user: 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300 ring-emerald-200 dark:ring-emerald-800',
    sales_channel: 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300 ring-indigo-200 dark:ring-indigo-800',
    carrier: 'bg-amber-50 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300 ring-amber-200 dark:ring-amber-800',
    inventory: 'bg-sky-50 dark:bg-sky-900/30 text-sky-700 dark:text-sky-300 ring-sky-200 dark:ring-sky-800',
    system: 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 ring-gray-200 dark:ring-gray-600',
};

const STATUS_SOURCE_LABELS: Record<string, string> = {
    user: 'Cambiado por',
    sales_channel: 'Cambiado desde',
    carrier: 'Cambiado desde',
    inventory: 'Cambiado desde',
    system: 'Cambiado desde',
};

function formatStatusChangedAt(value: string): string {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '';
    return date.toLocaleString('es-CO', {
        day: '2-digit',
        month: 'short',
        hour: '2-digit',
        minute: '2-digit',
    });
}

const OrderRow = memo(({
    order,
    onView,
    onEdit,
    onChangeStatus,
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
    isSuperAdmin,
    showNotifications,
    onShowNotifications
}: {
    order: Order;
    onView?: (order: Order) => void;
    onEdit?: (order: Order) => void;
    onChangeStatus?: (order: Order) => void;
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
    showNotifications: boolean;
    onShowNotifications: (order: Order) => void;
}) => {
    return (
        <tr className={`bg-white dark:bg-gray-800 transition-all duration-300 hover:bg-purple-50 dark:hover:bg-gray-700 cursor-pointer ${isNew ? 'animate-slide-in' : ''}`}>
            <td className="px-2 sm:px-3 py-2 whitespace-nowrap">
                <div className="flex items-center gap-1">
                    <div
                        className="h-8 w-8 rounded-full shadow-md border-2 border-gray-200 dark:border-gray-700 hover:shadow-lg transition-all cursor-pointer bg-white dark:bg-gray-800 flex items-center justify-center overflow-hidden"
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
                            <span className="text-xs font-medium text-gray-600 dark:text-gray-300 uppercase">
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
                <div className="flex items-center gap-1.5">
                    <span className="text-xs font-medium text-gray-900 dark:text-white">
                        {order.order_number || order.external_id || order.id}
                    </span>
                    {order.is_test && (
                        <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-bold bg-orange-100 text-orange-700 border border-orange-300 uppercase tracking-widest">
                            TEST
                        </span>
                    )}
                </div>
                <div className="text-xs text-gray-500 dark:text-gray-400 sm:hidden">
                    {order.customer_name}
                </div>
                {order.internal_number && (
                    <div className="text-xs text-gray-500 dark:text-gray-400">
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
                <div className="text-sm text-gray-900 dark:text-white">{order.customer_name}</div>
                <div className="text-xs text-gray-500 dark:text-gray-400">{order.customer_email}</div>
            </td>
            <td className="px-2 sm:px-3 py-2 whitespace-nowrap">
                {order.free_shipping && (
                    <span
                        className="mb-1 inline-flex items-center gap-1 rounded-full bg-emerald-100 px-2 py-px text-[9px] font-bold uppercase leading-tight tracking-wide text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300"
                        title="El canal no le cobro el envio al cliente"
                    >
                        Envio gratis
                    </span>
                )}
                <div className="flex items-center gap-1.5">
                    <span className={`inline-flex items-center px-1.5 py-0.5 rounded text-xs font-bold ${(order.currency_presentment || order.currency) === 'COP'
                            ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200'
                            : (order.currency_presentment || order.currency) === 'EUR'
                                ? 'bg-purple-100 text-purple-800'
                                : 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-200'
                        }`}>
                        {order.currency_presentment || order.currency || 'USD'}
                    </span>
                    <span className="text-sm font-semibold text-gray-900 dark:text-white">
                        {formatCurrency(order.total_amount, order.currency, order.total_amount_presentment, order.currency_presentment)}
                    </span>
                </div>
            </td>
            <td className="px-2 sm:px-3 py-2 whitespace-nowrap text-center">
                {getStatusBadge(order.order_status?.name || order.status, order.order_status?.color) || (
                    <span className="text-xs text-gray-400">-</span>
                )}
                {order.status_changed_by && (
                    <span
                        className={`mt-1 mx-auto flex w-fit items-center gap-1 rounded-full px-1.5 py-0.5 text-[9px] font-semibold leading-tight ring-1 ${STATUS_SOURCE_STYLES[order.status_source || 'system'] || STATUS_SOURCE_STYLES.system}`}
                        title={`${STATUS_SOURCE_LABELS[order.status_source || 'system']} ${order.status_changed_by}${order.status_changed_at ? ` el ${formatStatusChangedAt(order.status_changed_at)}` : ''}`}
                    >
                        {order.status_changed_by}
                        {order.status_changed_at && (
                            <span className="font-normal opacity-75">{formatStatusChangedAt(order.status_changed_at)}</span>
                        )}
                    </span>
                )}
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden lg:table-cell">
                <div className="flex flex-col items-center gap-1">
                    {(order.cod_total || 0) > 0 && (
                        order.cod_cut_confirmed ? (
                            <span
                                className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-emerald-100 dark:bg-emerald-900/30 text-emerald-800 dark:text-emerald-200 border border-emerald-300 dark:border-emerald-700"
                                title="Contra entrega liquidada en un corte de pago confirmado"
                            >
                                COD Pagado
                            </span>
                        ) : (
                            <span
                                className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200 border border-yellow-300 dark:border-yellow-700"
                                title="Contra entrega pendiente de corte de pago"
                            >
                                Contra Entrega
                            </span>
                        )
                    )}
                    {order.payment_status?.name ? (
                        getStatusBadge(order.payment_status.name, order.payment_status.color)
                    ) : order.is_paid ? (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200">
                            Pagado
                        </span>
                    ) : (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-200">
                            No pagado
                        </span>
                    )}
                    {order.invoice_status === 'issued' ? (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-green-100 text-green-700 dark:text-green-200">
                            Facturado
                        </span>
                    ) : order.invoice_status === 'pending' ? (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-yellow-100 text-yellow-700 dark:text-yellow-200">
                            Factura pendiente
                        </span>
                    ) : order.invoice_status === 'failed' ? (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-red-100 text-red-700">
                            Factura fallida
                        </span>
                    ) : order.invoice_status === 'cancelled' ? (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-gray-100 text-gray-600 dark:text-gray-300">
                            Factura cancelada
                        </span>
                    ) : order.invoice_status === 'none' ? (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-gray-100 text-gray-500 dark:text-gray-400">
                            Sin factura
                        </span>
                    ) : null}
                </div>
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden md:table-cell">
                {order.delivery_probability !== undefined && order.delivery_probability !== null ? (
                    <div className="flex flex-col gap-1 min-w-[120px]">
                        {order.shipping_geo_confidence && (
                            <span
                                className="self-start inline-flex items-center text-[9px] font-bold px-1.5 py-0.5 rounded-full"
                                style={
                                    order.shipping_geo_confidence === 'high'
                                        ? { backgroundColor: '#dcfce7', color: '#166534' }
                                        : order.shipping_geo_confidence === 'medium'
                                            ? { backgroundColor: '#fef9c3', color: '#854d0e' }
                                            : { backgroundColor: '#fee2e2', color: '#991b1b' }
                                }
                                title="Confianza del geocode de la direccion (afecta la probabilidad)"
                            >
                                {order.shipping_geo_confidence === 'high'
                                    ? 'Dir. confiable'
                                    : order.shipping_geo_confidence === 'medium'
                                        ? 'Dir. verificar'
                                        : 'Dir. dudosa'}
                            </span>
                        )}
                        <div className="flex items-center gap-2">
                        <div className="flex-1 bg-gray-200 rounded-full h-3 overflow-hidden shadow-inner">
                            <div
                                className={`h-full rounded-full transition-all duration-300 ${getProbabilityColor(order.delivery_probability)}`}
                                style={{ width: `${Math.min(order.delivery_probability, 100)}%` }}
                                title={`Probabilidad de entrega: ${order.delivery_probability.toFixed(1)}%`}
                            ></div>
                        </div>
                        <span className="text-xs font-semibold text-gray-700 dark:text-gray-200 min-w-[40px] text-right">
                            {order.delivery_probability.toFixed(0)}%
                        </span>
                        {(order.score_breakdown?.categories || (order.negative_factors && order.negative_factors.length > 0)) && (
                            <ProbabilityTooltip>
                                <div className="font-semibold mb-2 text-center">Análisis de probabilidad</div>

                                {order.score_breakdown?.categories ? (
                                    <div className="space-y-1.5 mb-2">
                                        {order.score_breakdown.categories.map((cat, idx) => (
                                            <div key={idx} className="flex items-center gap-2">
                                                <span className="w-[100px] truncate text-[10px] text-gray-300">{cat.name}</span>
                                                <div className="flex-1 bg-gray-700 rounded-full h-1.5 overflow-hidden">
                                                    <div
                                                        className={`h-full rounded-full ${cat.raw_score >= 80 ? 'bg-green-400' : cat.raw_score >= 60 ? 'bg-yellow-400' : 'bg-red-400'}`}
                                                        style={{ width: `${cat.raw_score}%` }}
                                                    />
                                                </div>
                                                <span className={`text-[10px] font-medium min-w-[28px] text-right ${cat.raw_score >= 80 ? 'text-green-400' : cat.raw_score >= 60 ? 'text-yellow-400' : 'text-red-400'}`}>
                                                    {cat.raw_score.toFixed(0)}
                                                </span>
                                            </div>
                                        ))}
                                    </div>
                                ) : null}

                                {order.negative_factors && order.negative_factors.length > 0 && (
                                    <div className="border-t border-gray-700 pt-2 mt-2">
                                        <div className="font-semibold mb-1 text-orange-400 text-[10px]">Datos faltantes:</div>
                                        {order.negative_factors.map((factor, idx) => (
                                            <div key={idx} className="flex items-center gap-1">
                                                <span className="text-orange-400">•</span>
                                                <span className="text-[10px]">{factor}</span>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </ProbabilityTooltip>
                        )}
                        </div>
                    </div>
                ) : (
                    <span className="text-xs text-gray-400">N/A</span>
                )}
            </td>
            <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-center">
                <div className="flex flex-col items-center gap-1">
                    {order.is_confirmed === true ? (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200">
                            Sí
                        </span>
                    ) : order.is_confirmed === false ? (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-200">
                            No
                        </span>
                    ) : (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200">
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

            {showNotifications && (
                <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-center">
                    {order.notifications && order.notifications.length > 0 ? (
                        <div className="flex items-center justify-center gap-1">
                            {order.notifications.map((counter) => (
                                <NotificationBadge
                                    key={counter.channel}
                                    counter={counter}
                                    onClick={() => onShowNotifications(order)}
                                />
                            ))}
                        </div>
                    ) : (
                        <span className="text-xs text-gray-400">-</span>
                    )}
                </td>
            )}

            <td className="px-3 sm:px-6 py-4 text-xs text-gray-500 dark:text-gray-400 hidden md:table-cell">
                <div className="leading-tight">
                    <div className="text-gray-900 dark:text-white">{formatDate(order.created_at).date}</div>
                    <div className="text-gray-500 dark:text-gray-400">{formatDate(order.created_at).time}</div>
                </div>
            </td>
            {isSuperAdmin && (
                <td className="px-3 sm:px-6 py-4 hidden lg:table-cell">
                    <div className="text-sm text-gray-900 dark:text-white">
                        {order.business_id && businessesMap.get(order.business_id)
                            ? businessesMap.get(order.business_id)
                            : order.business_id
                                ? `ID: ${order.business_id}`
                                : '-'
                        }
                    </div>
                </td>
            )}
            <td className="px-3 sm:px-6 py-4 text-center hidden md:table-cell">
                {order.shipment?.carrier ? (() => {
                    const logo = getCarrierLogo(order.shipment.carrier);
                    const hasGuide = !!order.shipment.tracking_number;
                    const isFreeShipping = normalizeCarrierLabel(order.shipment.carrier) === 'ENVIOGRATUITO';
                    const isQuoted = !!order.quoted_shipping && !hasGuide;
                    const title = hasGuide
                        ? order.shipment.carrier
                        : isQuoted
                            ? `Cotizado en el checkout: ${order.quoted_shipping!.title} ($${order.quoted_shipping!.price.toLocaleString('es-CO')}). Falta generar la guia.`
                            : `${order.shipment.carrier} (preseleccionada, sin guia generada)`;
                    const shipmentStatus = order.shipment.status || '';
                    const statusConfig = SHIPMENT_STATUS_CONFIG[shipmentStatus];
                    const carrierDetail = order.shipment.carrier_status_detail || order.shipment.carrier_status;
                    const carrierShortName = order.shipment.carrier.split(' - ')[0];
                    return (
                        <div
                            className="flex flex-col items-center justify-center gap-0.5 rounded-md py-1"
                            title={title}
                        >
                            {isQuoted && (
                                <span className="mb-0.5 rounded-full bg-indigo-100 px-1.5 py-px text-[9px] font-bold uppercase leading-tight tracking-wide text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300">
                                    Cotizado
                                </span>
                            )}
                            {logo ? (
                                <img
                                    src={logo}
                                    alt={order.shipment.carrier}
                                    className={`h-8 w-12 object-contain ${hasGuide ? '' : 'grayscale opacity-60'}`}
                                    loading="lazy"
                                />
                            ) : (
                                <span className={`text-xs font-medium px-2 py-1 rounded ${
                                    isFreeShipping
                                        ? 'text-emerald-700 dark:text-emerald-300 bg-emerald-100 dark:bg-emerald-900/40'
                                        : hasGuide
                                            ? 'text-gray-600 dark:text-gray-300 bg-gray-100'
                                            : 'text-amber-700 dark:text-amber-300 bg-amber-100 dark:bg-amber-900/40'
                                }`}>
                                    {carrierShortName.substring(0, 12)}
                                </span>
                            )}
                            {hasGuide && shipmentStatus && (
                                <span
                                    className={`mt-0.5 inline-flex max-w-[150px] items-center rounded-full border px-1.5 py-px text-[9px] font-semibold leading-tight ${statusConfig?.className || SHIPMENT_STATUS_FALLBACK.className}`}
                                    title={`Estado del envio en Probability: ${statusConfig?.label || shipmentStatus}`}
                                >
                                    {statusConfig?.label || shipmentStatus}
                                </span>
                            )}
                            {hasGuide && carrierDetail && (
                                <span
                                    className={`inline-flex max-w-[150px] items-center gap-1 rounded-full border border-dashed px-1.5 py-px text-[9px] font-medium leading-tight ${statusConfig?.className || SHIPMENT_STATUS_FALLBACK.className}`}
                                    title={`${carrierShortName}: ${carrierDetail}`}
                                >
                                    <span className="truncate">{carrierDetail}</span>
                                </span>
                            )}
                        </div>
                    );
                })() : (
                    <span className="text-xs text-gray-400 text-center block">Sin envío</span>
                )}
            </td>
            <td className="px-2 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div className="flex flex-row justify-end gap-1">

                    {onViewRecommendation && (
                        <button
                            onClick={() => onViewRecommendation(order)}
                            disabled={!!order.guide_link}
                            className={`p-2 rounded-md transition-all duration-200 flex items-center justify-center shadow-sm ${order.guide_link
                                    ? 'bg-gray-400 text-gray-600 dark:text-gray-300 cursor-not-allowed opacity-60'
                                    : 'bg-gradient-to-r from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700 text-white focus:ring-2 focus:ring-purple-500 focus:ring-offset-2'
                                }`}
                            title={order.guide_link ? 'Guía ya validada' : 'Recomendación Inteligente IA'}
                            aria-label={order.guide_link ? 'Guía ya validada' : 'Ver recomendación IA'}
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
                    {onChangeStatus && (
                        <button
                            onClick={() => onChangeStatus(order)}
                            disabled={isTerminalStatus(order.order_status?.code || order.status || '')}
                            className={`p-2 rounded-md transition-colors duration-200 focus:ring-2 focus:ring-offset-2 ${
                                isTerminalStatus(order.order_status?.code || order.status || '')
                                    ? 'bg-gray-300 dark:bg-gray-600 text-gray-500 dark:text-gray-400 cursor-not-allowed'
                                    : 'bg-amber-500 hover:bg-amber-600 text-white focus:ring-amber-500'
                            }`}
                            title={isTerminalStatus(order.order_status?.code || order.status || '')
                                ? 'Estado terminal — no se puede cambiar'
                                : 'Cambiar estado'}
                            aria-label="Cambiar estado"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" />
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
                    {isSuperAdmin && (
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
                    )}
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
    selectedBusinessId?: number | null;
}

export default function OrderList({ onView, onEdit, onViewRecommendation, refreshKey, onCreate, onTestGuide, selectedBusinessId = null }: OrderListProps) {
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

    const [guideUrl, setGuideUrl] = useState<string | null>(null);
    const [notificationsOrder, setNotificationsOrder] = useState<Order | null>(null);

    const [selectedOrderId, setSelectedOrderId] = useState<string | null>(null);
    const [selectedOrderLogo, setSelectedOrderLogo] = useState<string | undefined>(undefined);
    const [selectedOrderPlatform, setSelectedOrderPlatform] = useState<string | undefined>(undefined);
    const [isRawModalOpen, setIsRawModalOpen] = useState(false);

    const [isDownloadModalOpen, setIsDownloadModalOpen] = useState(false);
    const [isStatusFlowOpen, setIsStatusFlowOpen] = useState(false);
    const [isQuotationExpresOpen, setIsQuotationExpresOpen] = useState(false);
    const [changeStatusOrder, setChangeStatusOrder] = useState<Order | null>(null);
    const [filtersSlot, setFiltersSlot] = useState<HTMLElement | null>(null);
    const [actionsSlot, setActionsSlot] = useState<HTMLElement | null>(null);

    useEffect(() => {
        setFiltersSlot(document.getElementById(ORDERS_FILTERS_SLOT_ID));
        setActionsSlot(document.getElementById(ORDERS_ACTIONS_SLOT_ID));
    }, []);

    const [integrationsList, setIntegrationsList] = useState<{ value: string; label: string }[]>([]);

    const [orderStatusesList, setOrderStatusesList] = useState<{ value: string; label: string }[]>([]);

    const [paymentStatusesList, setPaymentStatusesList] = useState<{ value: string; label: string }[]>([]);

    const [fulfillmentStatusesList, setFulfillmentStatusesList] = useState<{ value: string; label: string }[]>([]);

    const [businessesList, setBusinessesList] = useState<{ id: number; name: string }[]>([]);

    const showNotificationsColumn = useMemo(
        () => orders.some((o) => (o.notifications?.length ?? 0) > 0),
        [orders]
    );

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
                const response = await getOrderStatusesAction(true);
                if (response.success && response.data) {
                    const options = response.data
                        .filter((status) => status?.name)
                        .map((status) => ({
                            value: String(status.id),
                            label: status.name,
                        }));

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
                const response = await getPaymentStatusesAction({ is_active: true });
                if (response.success && response.data) {
                    const options = response.data
                        .filter((status) => status?.name)
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
                const response = await getFulfillmentStatusesAction(true);
                if (response.success && response.data) {
                    const options = response.data
                        .filter((status) => status?.name)
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

    const [filters, setFilters] = useState<GetOrdersParams>({
        page: 1,
        page_size: 10,
    });

    const { showToast } = useToast();

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
            key: 'customer_name',
            label: 'Nombre del cliente',
            type: 'text',
            placeholder: 'Buscar por nombre...',
        },
        {
            key: 'customer_email',
            label: 'Correo del cliente',
            type: 'text',
            placeholder: 'Buscar por correo...',
        },
        {
            key: 'customer_phone',
            label: 'Teléfono del cliente',
            type: 'text',
            placeholder: 'Buscar por teléfono...',
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
            key: 'currency_presentment',
            label: 'Moneda',
            type: 'select',
            options: [
                { value: 'COP', label: 'COP' },
                { value: 'USD', label: 'USD' },
                { value: 'EUR', label: 'EUR' },
            ],
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
            key: 'invoice_status',
            label: 'Estado de factura',
            type: 'select',
            options: [
                { value: 'none', label: 'Sin factura' },
                { value: 'pending', label: 'Factura pendiente' },
                { value: 'issued', label: 'Facturado' },
                { value: 'failed', label: 'Factura fallida' },
                { value: 'cancelled', label: 'Factura cancelada' },
            ],
        },
        {
            key: 'start_date',
            label: 'Rango de fechas',
            type: 'date-range',
        },
    ], [orderStatusesList, integrationsList, paymentStatusesList, fulfillmentStatusesList]);

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

        if (filters.customer_name) {
            active.push({
                key: 'customer_name',
                label: 'Nombre del cliente',
                value: filters.customer_name,
                type: 'text',
            });
        }

        if (filters.customer_email) {
            active.push({
                key: 'customer_email',
                label: 'Correo del cliente',
                value: filters.customer_email,
                type: 'text',
            });
        }

        if (filters.customer_phone) {
            active.push({
                key: 'customer_phone',
                label: 'Teléfono del cliente',
                value: filters.customer_phone,
                type: 'text',
            });
        }

        if (filters.status) {

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

        if (filters.currency_presentment) {
            active.push({
                key: 'currency_presentment',
                label: 'Moneda',
                value: filters.currency_presentment,
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

        if (filters.invoice_status) {
            const invoiceLabels: Record<string, string> = {
                none: 'Sin factura',
                pending: 'Factura pendiente',
                issued: 'Facturado',
                failed: 'Factura fallida',
                cancelled: 'Factura cancelada',
            };
            active.push({
                key: 'invoice_status',
                label: 'Estado de factura',
                value: invoiceLabels[filters.invoice_status] || filters.invoice_status,
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

    const handleSortChange = useCallback((sortBy: string, sortOrder: 'asc' | 'desc') => {
        setFilters((prev) => ({
            ...prev,
            sort_by: sortBy as 'created_at' | 'updated_at' | 'total_amount' | 'order_number',
            sort_order: sortOrder,
            page: 1,
        }));
    }, []);

    useSSE({
        eventTypes: ['order.created'],

        businessId: isSuperAdmin ? 0 : permissions?.business_id,
        onMessage: async (event) => {
            try {
                const data = JSON.parse(event.data);

                const orderId = data.data?.order_id || data.metadata?.order_id;
                const orderNumber = data.data?.order_number || 'Desconocida';

                if (data.type === 'order.created' && orderId) {

                    try {
                        const response = await getOrderByIdAction(orderId);
                        if (response.success && response.data) {
                            const newOrder = response.data;

                            if (!isSuperAdmin && permissions?.business_id && newOrder.business_id !== permissions.business_id) {
                                return;
                            }

                            setOrders(prevOrders => {
                                if (prevOrders.some(o => o.id === newOrder.id)) {
                                    return prevOrders;
                                }

                                setNewOrderIds(prev => new Set(prev).add(newOrder.id));

                                setTimeout(() => {
                                    setNewOrderIds(prev => {
                                        const updated = new Set(prev);
                                        updated.delete(newOrder.id);
                                        return updated;
                                    });
                                }, 2000);

                                return [newOrder, ...prevOrders];
                            });

                            setTotal(prev => prev + 1);

                            playNotificationSound();

                            showToast(`Nueva orden recibida: #${orderNumber}`, 'success');
                        }
                    } catch (err) {
                        console.error('Error al obtener orden completa:', err);

                        refreshTableOnly();
                    }
                }
            } catch (e) {
                console.error('Error processing SSE message:', e);
            }
        },
    });

    useSSE({
        eventTypes: ['order.score_calculated'],
        businessId: isSuperAdmin ? 0 : permissions?.business_id,
        onMessage: async (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.type !== 'order.score_calculated') return;

                const orderId = data.data?.order_id || data.metadata?.order_id;
                if (!orderId) return;

                const response = await getOrderByIdAction(orderId);
                if (response.success && response.data) {
                    setOrders(prev => prev.map(o =>
                        o.id === orderId ? { ...o, delivery_probability: response.data!.delivery_probability, negative_factors: response.data!.negative_factors } : o
                    ));
                }
            } catch {
            }
        },
    });

    useSSE({
        eventTypes: ['order.status_changed'],
        businessId: isSuperAdmin ? 0 : permissions?.business_id,
        onMessage: async (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.type !== 'order.status_changed') return;

                const orderId = data.data?.order_id || data.metadata?.order_id;
                if (!orderId) return;

                const response = await getOrderByIdAction(orderId);
                if (response.success && response.data) {
                    const updatedOrder = response.data;

                    setOrders(prev => prev.map(o =>
                        o.id === orderId ? updatedOrder : o
                    ));

                    if (updatedOrder.is_confirmed === true) {
                        showToast(`Pedido #${updatedOrder.order_number} confirmado por WhatsApp`, 'success');
                    }
                }
            } catch {
            }
        },
    });

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
            setError(getActionError(err, 'Error al cargar las órdenes'));
        } finally {
            setInitialLoading(false);
            setTableLoading(false);
        }
    }, [filters, isSuperAdmin, selectedBusinessId]);

    const handleDownloadOrders = useCallback(async (startDate: string, endDate: string) => {
        try {
            let ordersToDownload: Order[] = [];
            let page = 1;
            let hasMore = true;

            while (hasMore) {
                const params: GetOrdersParams = {
                    page,
                    page_size: 50000,
                    start_date: startDate,
                    end_date: endDate,
                };

                if (isSuperAdmin && selectedBusinessId !== null) {
                    params.business_id = selectedBusinessId;
                }

                const response = await getOrdersAction(params);

                if (!response.success || !response.data) {
                    throw new Error(response.message || 'Error al obtener órdenes');
                }

                ordersToDownload = [...ordersToDownload, ...response.data];

                const currentPage = response.page || 1;
                const totalPages = response.total_pages || 1;
                hasMore = currentPage < totalPages;
                page++;
            }

            const excelData = ordersToDownload.map((order: Order) => {
                const resolved = resolveCityState(order.shipping_city || '', order.shipping_state || '', order.destination_dane_code);
                return {
                    'ID Orden': order.order_number || order.external_id || order.id,
                    'Cliente': order.customer_name || '',
                    'Email': order.customer_email || '',
                    'Teléfono': order.customer_phone || '',
                    'Plataforma': order.platform || '',
                    'Estado': order.status || '',
                    'Total': order.total_amount || 0,
                    'Moneda': order.currency || '',
                    'Pagado': order.is_paid ? 'Sí' : 'No',
                    'Contra Entrega COD': order.cod_total || 0,
                    'Dirección': order.shipping_street || '',
                    'Ciudad': resolved?.ciudad || order.shipping_city || '',
                    'Departamento': resolved?.departamento || order.shipping_state || '',
                    'Valor Envío': order.shipment?.total_cost || order.quoted_shipping?.price || order.shipping_cost || 0,
                    'Guía': order.tracking_number || order.shipment?.tracking_number || '',
                    'Fecha Creación': order.created_at ? new Date(order.created_at).toLocaleString('es-CO') : '',
                };
            });

            const ws = XLSX.utils.json_to_sheet(excelData);
            const wb = XLSX.utils.book_new();
            XLSX.utils.book_append_sheet(wb, ws, 'Órdenes');

            const colWidths = [
                { wch: 15 },
                { wch: 20 },
                { wch: 25 },
                { wch: 15 },
                { wch: 12 },
                { wch: 15 },
                { wch: 12 },
                { wch: 10 },
                { wch: 10 },
                { wch: 15 },
                { wch: 30 },
                { wch: 15 },
                { wch: 15 },
                { wch: 12 },
                { wch: 30 },
                { wch: 20 },
            ];
            ws['!cols'] = colWidths;

            const fileName = `ordenes_${startDate}_a_${endDate}.xlsx`;
            XLSX.writeFile(wb, fileName);

            showToast(`Se descargaron ${excelData.length} órdenes correctamente`, 'success');
        } catch (error) {
            console.error('Error al descargar órdenes:', error);
            throw error;
        }
    }, [filters, isSuperAdmin, selectedBusinessId, showToast]);

    useEffect(() => {
        if (isSuperAdmin) {
            setFilters(prev => ({ ...prev, page: 1 }));
        }
    }, [selectedBusinessId, isSuperAdmin]);

    useEffect(() => {
        if (isFirstLoad.current) {
            isFirstLoad.current = false;
            loadOrders(true);
        }
    }, [loadOrders]);

    useEffect(() => {
        if (!isFirstLoad.current) {
            loadOrders(false);
        }
    }, [filters, loadOrders]);

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

        const finalAmount = (amountPresentment && amountPresentment > 0 && currencyPresentment) ? amountPresentment : amount;
        return new Intl.NumberFormat('es-CO', {
            style: 'decimal',
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        }).format(finalAmount);
    }, []);

    const formatDate = useCallback((dateString: string): { date: string; time: string } => {
        const date = new Date(dateString);

        const parts = date.toLocaleDateString('es-CO', {
            day: 'numeric',
            month: 'short',
            year: 'numeric',
        }).split(' ');
        const dateStr = `${parts[0]} ${parts[2]} ${parts[parts.length - 1]}`;

        const timeStr = date.toLocaleTimeString('es-CO', {
            hour: '2-digit',
            minute: '2-digit',
            hour12: true,
        }).replace(/\s/g, '');
        return { date: dateStr, time: timeStr };
    }, []);

    const getStatusBadge = useCallback((status: string | undefined | null, color?: string) => {

        const cleanStatus = status && typeof status === 'string' ? status.trim() : '';
        if (!cleanStatus) {
            return null;
        }

        if (color) {

            const hex = color.replace('#', '');
            const r = parseInt(hex.substr(0, 2), 16);
            const g = parseInt(hex.substr(2, 2), 16);
            const b = parseInt(hex.substr(4, 2), 16);

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
                    {cleanStatus}
                </span>
            );
        }

        const statusColors: Record<string, string> = {
            pending: 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200',
            processing: 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-200',
            shipped: 'bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-200',
            delivered: 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200',
            cancelled: 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-200',
        };

        const colorClass = statusColors[cleanStatus.toLowerCase()] || 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200';

        return (
            <span className={`px-2 py-1 text-xs font-medium rounded-full ${colorClass}`}>
                {cleanStatus}
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

    const filtersBar = (
        <DynamicFilters
            variant="bar"
            availableFilters={availableFilters}
            activeFilters={activeFilters}
            onAddFilter={handleAddFilter}
            onRemoveFilter={handleRemoveFilter}
            sortBy={filters.sort_by || 'created_at'}
            sortOrder={filters.sort_order || 'desc'}
            onSortChange={handleSortChange}
            sortOptions={[
                { value: 'created_at', label: 'Por fecha' },
                { value: 'updated_at', label: 'Por actualizacion' },
                { value: 'total_amount', label: 'Por monto' },
                { value: 'order_number', label: 'Por ID' },
            ]}
        />
    );

    const barActions = (
        <>
            <IconActionButton
                label="Cotizador Expres"
                variant="secondary"
                onClick={() => setIsQuotationExpresOpen(true)}
                icon={<ChartBarIcon className="w-4 h-4" />}
            />
            <IconActionButton
                label="Flujo de Estados"
                variant="secondary"
                onClick={() => setIsStatusFlowOpen(true)}
                icon={<ArrowPathRoundedSquareIcon className="w-4 h-4" />}
            />
            <IconActionButton
                label="Descargar ordenes en Excel"
                variant="tertiary"
                onClick={() => setIsDownloadModalOpen(true)}
                icon={<ArrowDownTrayIcon className="w-4 h-4" />}
            />
        </>
    );

    return (
        <div>
            {actionsSlot
                ? createPortal(barActions, actionsSlot)
                : <div className="mb-2 flex items-center gap-2">{barActions}</div>}

            {filtersSlot
                ? createPortal(filtersBar, filtersSlot)
                : <div className="mb-4">{filtersBar}</div>}

            <div className="ordersTable relative rounded-xl overflow-hidden shadow-sm bg-white dark:bg-gray-800">

                {tableLoading && (
                    <div className="absolute inset-0 bg-white dark:bg-gray-800/80 backdrop-blur-sm z-10 flex items-center justify-center transition-opacity duration-200">
                        <div className="flex flex-col items-center gap-2">
                            <div className="w-8 h-8 border-4 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
                            <p className="text-sm text-gray-600 dark:text-gray-300">Actualizando...</p>
                        </div>
                    </div>
                )}
                <div className="overflow-x-auto">
                    <table className={`min-w-full transition-opacity duration-200 ${tableLoading ? 'opacity-50' : 'opacity-100'}`}>
                        <thead style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}>
                            <tr>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest w-16" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>

                                </th>
                                <th
                                    className="px-3s sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest cursor-pointer transition-all group"
                                    onClick={() => handleSortChange('order_number', filters.sort_order === 'asc' ? 'desc' : 'asc')}
                                    style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}
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
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden sm:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                    Cliente
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                    Total
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-center text-xs font-bold text-white uppercase tracking-widest" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                    Estado
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-center text-xs font-bold text-white uppercase tracking-widest hidden lg:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                    Estatus Pago
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden md:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                    Probabilidad
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-center text-xs font-bold text-white uppercase tracking-widest" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                    Confirmado
                                </th>

                                {showNotificationsColumn && (
                                    <th className="px-3 sm:px-6 py-3 text-center text-xs font-bold text-white uppercase tracking-widest" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                        Notif.
                                    </th>
                                )}
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden md:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                    Fecha
                                </th>
                                {isSuperAdmin && (
                                    <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden lg:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                        Business
                                    </th>
                                )}
                                <th className="px-3 sm:px-6 py-3 text-center text-xs font-bold text-white uppercase tracking-widest hidden md:table-cell" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                    Transportadora
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-right text-xs font-bold text-white uppercase tracking-widest" style={{ paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' }}>
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {orders.length === 0 ? (
                                <tr>
                                    <td colSpan={isSuperAdmin ? 10 : 9} className="px-4 sm:px-6 py-8 text-center text-gray-500 dark:text-gray-400">
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
                                        onChangeStatus={(o) => setChangeStatusOrder(o)}
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
                                        showNotifications={showNotificationsColumn}
                                        onShowNotifications={setNotificationsOrder}
                                    />
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {(totalPages > 1 || total > 0) && (
                    <div className="px-3 py-2 flex flex-col sm:flex-row items-center justify-between gap-2" style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}>

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

                        <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between w-full">
                            <div className="flex items-center gap-3">
                                <p className="text-[11px] text-white">
                                    Mostrando <span className="font-medium">{(page - 1) * (filters.page_size || 10) + 1}</span> a{' '}
                                    <span className="font-medium">{Math.min(page * (filters.page_size || 10), total)}</span> de{' '}
                                    <span className="font-medium">{total}</span> resultados
                                </p>
                                <div className="flex items-center gap-1">
                                    <label className="text-[11px] text-white whitespace-nowrap">
                                        Mostrar:
                                    </label>
                                    <select
                                        value={filters.page_size || 10}
                                        onChange={(e) => {
                                            const newPageSize = parseInt(e.target.value);
                                            setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                        }}
                                        className="px-1.5 py-1 text-[11px] border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 dark:text-white bg-white dark:bg-gray-800"
                                    >
                                        <option value="10">10</option>
                                        <option value="20">20</option>
                                        <option value="50">50</option>
                                        <option value="100">100</option>
                                    </select>
                                </div>
                            </div>
                            <div className="flex items-center gap-1">
                                <nav className="relative z-0 inline-flex items-center gap-1 flex-wrap justify-end">

                                    <button
                                        onClick={() => setFilters({ ...filters, page: 1 })}
                                        disabled={page === 1}
                                        className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                        title="Primera página"
                                    >
                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 19l-7-7 7-7m8 14l-7-7 7-7" /></svg>
                                    </button>

                                    <button
                                        onClick={() => setFilters({ ...filters, page: page - 1 })}
                                        disabled={page === 1}
                                        className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                        title="Anterior"
                                    >
                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" /></svg>
                                    </button>

                                    {(() => {
                                        const pages: (number | string)[] = [];
                                        const maxVisible = 8;

                                        if (totalPages <= maxVisible) {
                                            for (let i = 1; i <= totalPages; i++) pages.push(i);
                                        } else {
                                            const windowSize = maxVisible - 2;
                                            let start = Math.max(2, page - Math.floor(windowSize / 2));
                                            let end = start + windowSize - 1;
                                            if (end >= totalPages) {
                                                end = totalPages - 1;
                                                start = Math.max(2, end - windowSize + 1);
                                            }

                                            pages.push(1);
                                            if (start > 2) pages.push('...');
                                            for (let i = start; i <= end; i++) pages.push(i);
                                            if (end < totalPages - 1) pages.push('...');
                                            pages.push(totalPages);
                                        }

                                        return pages.map((p, idx) =>
                                            typeof p === 'string' ? (
                                                <span key={`ellipsis-${idx}`} className="relative inline-flex items-center px-1 py-1 text-[11px] font-bold text-white/70">
                                                    ...
                                                </span>
                                            ) : (
                                                <button
                                                    key={p}
                                                    onClick={() => setFilters({ ...filters, page: p })}
                                                    className={`relative inline-flex items-center justify-center min-w-7 px-2 py-1 rounded-md border text-[11px] font-semibold transition-all ${
                                                        p === page
                                                            ? 'page-btn-active z-10 shadow-md scale-105'
                                                            : 'page-btn'
                                                    }`}
                                                    style={p === page
                                                        ? { backgroundColor: 'var(--color-secondary)', borderColor: 'var(--color-secondary)', color: 'white' }
                                                        : undefined}
                                                >
                                                    {p}
                                                </button>
                                            )
                                        );
                                    })()}

                                    <button
                                        onClick={() => setFilters({ ...filters, page: page + 1 })}
                                        disabled={page === totalPages}
                                        className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                        title="Siguiente"
                                    >
                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
                                    </button>

                                    <button
                                        onClick={() => setFilters({ ...filters, page: totalPages })}
                                        disabled={page === totalPages}
                                        className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                        title="Última página"
                                    >
                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 5l7 7-7 7M5 5l7 7-7 7" /></svg>
                                    </button>
                                </nav>
                            </div>
                        </div>

                        <div className="flex items-center justify-between w-full sm:hidden pt-2">
                            <div className="flex items-center gap-1">
                                <label className="text-xs text-white whitespace-nowrap">
                                    Mostrar:
                                </label>
                                <select
                                    value={filters.page_size || 10}
                                    onChange={(e) => {
                                        const newPageSize = parseInt(e.target.value);
                                        setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                    }}
                                    className="px-2 py-1.5 text-xs border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 dark:text-white bg-white dark:bg-gray-800"
                                >
                                    <option value="10">10</option>
                                    <option value="20">20</option>
                                    <option value="50">50</option>
                                    <option value="100">100</option>
                                </select>
                            </div>
                            <p className="text-xs text-white/80">
                                Página {page} de {totalPages}
                            </p>
                        </div>
                    </div>
                )}

                <style jsx>{`
                    .ordersTable :global(table) {
                        zoom: 0.95;
                    }

                    /* Tabla más "card-like" fila por fila */
                    .ordersTable :global(.table) {
                        border-collapse: separate;
                        border-spacing: 0 10px; /* separación entre filas */
                        background: transparent;
                    }

                    /* Quitar el borde del contenedor global de Table SOLO aquí */
                    .ordersTable :global(div.overflow-hidden.w-full.rounded-lg.border.border-gray-200 dark:border-gray-700.bg-white dark:bg-gray-800) {
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
                        transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
                    }

                    .ordersTable :global(.table tbody tr:hover) {
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

            {notificationsOrder && (
                <OrderNotificationsModal
                    orderId={notificationsOrder.id}
                    orderNumber={notificationsOrder.order_number}
                    businessId={isSuperAdmin ? notificationsOrder.business_id : undefined}
                    onClose={() => setNotificationsOrder(null)}
                    onChanged={() => loadOrders()}
                />
            )}

            {guideUrl && (
                <div
                    className="fixed inset-0 z-50 flex items-center justify-center"
                    onClick={() => setGuideUrl(null)}
                >
                    <div className="absolute inset-0 bg-black/60" />
                    <div
                        className="relative z-10 bg-white dark:bg-gray-800 rounded-xl shadow-2xl flex flex-col overflow-hidden"
                        style={{ width: '480px', maxWidth: '90vw' }}
                        onClick={(e) => e.stopPropagation()}
                    >

                        <div className="flex items-center justify-between px-5 py-4 border-b bg-gray-50">
                            <div className="flex items-center gap-2 text-sm font-semibold text-gray-800 dark:text-gray-100">
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

                        <div className="p-8 flex flex-col items-center gap-6">
                            <div className="w-20 h-20 rounded-2xl bg-red-50 flex items-center justify-center">
                                <svg className="w-10 h-10 text-red-500" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zm-1 1.5L18.5 9H13V3.5zM6 4h5v7h7v9H6V4z" />
                                    <path d="M8 12h8v1.5H8V12zm0 3h8v1.5H8V15zm0 3h5v1.5H8V18z" />
                                </svg>
                            </div>
                            <div className="text-center">
                                <p className="font-semibold text-gray-800 dark:text-gray-100 text-lg">Guía de Envío lista</p>
                                <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">El PDF está disponible para ver o descargar</p>
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
                                    className="flex items-center justify-center gap-2 w-full py-3 px-6 bg-gray-100 hover:bg-gray-200 text-gray-700 dark:text-gray-200 font-semibold rounded-xl transition-colors"
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

            <DownloadOrdersModal
                isOpen={isDownloadModalOpen}
                onClose={() => setIsDownloadModalOpen(false)}
                onDownload={handleDownloadOrders}
            />

            <OrderStatusFlowModal
                isOpen={isStatusFlowOpen}
                onClose={() => setIsStatusFlowOpen(false)}
            />

            <QuotationExpresModal
                isOpen={isQuotationExpresOpen}
                onClose={() => setIsQuotationExpresOpen(false)}
                business_id={selectedBusinessId || undefined}
            />

            {changeStatusOrder && (
                <ChangeStatusModal
                    isOpen={!!changeStatusOrder}
                    onClose={() => setChangeStatusOrder(null)}
                    order={changeStatusOrder}
                    onSuccess={() => {
                        showToast(`Estado de #${changeStatusOrder.order_number} actualizado`, 'success');
                        setChangeStatusOrder(null);
                    }}
                />
            )}
        </div>
    );
}
