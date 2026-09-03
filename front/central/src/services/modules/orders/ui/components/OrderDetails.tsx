'use client';

import { Order, OrderHistory } from '../../domain/types';
import MapComponent from '@/shared/ui/MapComponent';
import { getAIRecommendationAction, getOrderByIdAction, getOrderHistoryAction, updateOrderAction, requestWhatsAppConfirmationAction, checkWhatsAppIntegrationAction } from '../../infra/actions';
import { useState, useEffect } from 'react';
import Link from 'next/link';
import ShipmentGuideModal from '@/shared/ui/modals/shipment-guide-modal';
import { ChangeStatusModal } from './ChangeStatusModal';
import { isTerminalStatus } from '../../domain/order-status-transitions';
import { getOrderStatusLabel, isPaymentStatus } from '../../domain/order-status-labels';
import { salesChannelLabel } from '../../domain/sales-channel-labels';
import { useToast } from '@/shared/providers/toast-provider';
import { IVAIncludedBadge } from './IVAIncludedBadge';
import { useDynamicBusinessColors } from '../hooks/useDynamicBusinessColors';
import { resolveCityState } from '@/shared/utils/dane-lookup';
import dynamic from 'next/dynamic';
import { Package, Copy, Check, X, Calendar, CreditCard, Truck, Receipt, Percent, Wallet, ShoppingBag, User, Phone, Mail, MapPin, ClipboardList, Save, MessageCircle, History, AlertTriangle, RefreshCw, Clock, FileText, ChevronRight, ArrowRight, Plus, CircleCheck, Scissors } from 'lucide-react';
const GeozoneMiniMap = dynamic(() => import('@/services/modules/geozones/ui/components/GeozoneMiniMap').then(m => m.GeozoneMiniMap), { ssr: false });

interface Quotation {
    carrier: string;
    estimated_cost: number;
    estimated_delivery_days: number;
}

interface AIRecommendation {
    recommended_carrier: string;
    reasoning: string;
    alternatives: string[];
    quotations?: Quotation[];
}

interface OrderDetailsProps {
    initialOrder: Order;
    onClose?: () => void;
    mode?: 'details' | 'recommendation';
}

const SOURCE_LABELS: Record<string, string> = {
    user: 'Usuario',
    sales_channel: 'Canal de venta',
    carrier: 'Transportadora',
    inventory: 'Inventario',
    system: 'Sistema',
};

function statusSourceStyle(source?: string): string {
    switch (source) {
        case 'sales_channel':
            return 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300';
        case 'carrier':
            return 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300';
        case 'inventory':
            return 'bg-sky-100 text-sky-700 dark:bg-sky-900/40 dark:text-sky-300';
        case 'user':
            return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300';
        default:
            return 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300';
    }
}

export default function OrderDetails({ initialOrder, onClose, mode = 'details' }: OrderDetailsProps) {
    const [fullOrder, setFullOrder] = useState<Order | null>(null);
    const [aiRecommendation, setAIRecommendation] = useState<AIRecommendation | null>(null);
    const [loadingAI, setLoadingAI] = useState(false);
    const [loadingDetails, setLoadingDetails] = useState(false);
    const [showGuideModal, setShowGuideModal] = useState(false);
    const [showChangeStatus, setShowChangeStatus] = useState(false);
    const [statusHistory, setStatusHistory] = useState<OrderHistory[]>([]);
    const [loadingHistory, setLoadingHistory] = useState(false);
    const { showToast } = useToast();
    const { colors: businessColors } = useDynamicBusinessColors(initialOrder.business_id);

    const primaryColor = businessColors?.primary_color || '#5b21b6';
    const secondaryColor = businessColors?.secondary_color || '#7c3aed';
    const tertiaryColor = businessColors?.tertiary_color || '#c4b5fd';
    const quaternaryColor = businessColors?.quaternary_color || '#ede9fe';

    const fetchDetails = async () => {
        if (!initialOrder.id) return;

        setLoadingDetails(true);
        setLoadingHistory(true);
        try {
            const [orderResponse, historyResponse] = await Promise.all([
                getOrderByIdAction(initialOrder.id),
                getOrderHistoryAction(initialOrder.id),
            ]);
            if (orderResponse.success && orderResponse.data) {
                setFullOrder(orderResponse.data);
            } else if (!orderResponse.success) {
                console.error("Failed to load order details:", orderResponse.message);
            }
            if (historyResponse.success && historyResponse.data) {
                setStatusHistory(historyResponse.data);
            }
        } catch (error) {
            console.error("Error loading order details:", error);
        } finally {
            setLoadingDetails(false);
            setLoadingHistory(false);
        }
    };

    useEffect(() => {
        fetchDetails();
    }, [initialOrder.id, initialOrder.updated_at]);

    const order = fullOrder || initialOrder;

    useEffect(() => {
        if (mode === 'recommendation' && fullOrder && fullOrder.shipping_city && fullOrder.shipping_state) {
            setLoadingAI(true);
            getAIRecommendationAction(fullOrder.shipping_city, fullOrder.shipping_state)
                .then(data => {
                    if (data && data.recommended_carrier) {
                        setAIRecommendation(data);
                    } else {
                        setAIRecommendation(null);
                    }
                })
                .catch(err => {
                    console.warn("Recomendación AI no disponible:", err);
                    setAIRecommendation(null);
                })
                .finally(() => setLoadingAI(false));
        } else if (mode !== 'recommendation') {

            setAIRecommendation(null);
            setLoadingAI(false);
        }
    }, [fullOrder, mode]);

    const [isConfirmed, setIsConfirmed] = useState<boolean | null>(false);
    const [novelty, setNovelty] = useState('');
    const [isSaving, setIsSaving] = useState(false);
    const [copiedNumber, setCopiedNumber] = useState(false);

    const [hasWhatsApp, setHasWhatsApp] = useState(false);
    const [isSendingWhatsApp, setIsSendingWhatsApp] = useState(false);
    const [whatsAppSent, setWhatsAppSent] = useState(false);

    useEffect(() => {
        if (order) {
            setIsConfirmed(order.is_confirmed ?? null);
            setNovelty(order.novelty || '');
        }
    }, [order]);

    useEffect(() => {
        if (order?.business_id != null && order.business_id > 0) {
            checkWhatsAppIntegrationAction(order.business_id).then((result) => {
                setHasWhatsApp(result);
            });
        }
    }, [order?.business_id]);

    const handleWhatsAppConfirmation = async () => {
        if (!order.id) return;
        setIsSendingWhatsApp(true);
        try {
            const result = await requestWhatsAppConfirmationAction(order.id);
            if (result.success) {
                setWhatsAppSent(true);
                alert('Mensaje de confirmación enviado por WhatsApp');
            } else {
                alert(result.message || 'Error al enviar confirmación por WhatsApp');
            }
        } catch (error: any) {
            console.error('Error sending WhatsApp confirmation:', error);
            alert('Error al enviar confirmación por WhatsApp');
        } finally {
            setIsSendingWhatsApp(false);
        }
    };

    const handleSaveManagement = async () => {
        if (!order.id) return;
        setIsSaving(true);
        try {
            const status = isConfirmed === true ? 'yes' : isConfirmed === false ? 'no' : 'pending';
            const result = await updateOrderAction(order.id, {
                confirmation_status: status,
                novelty: novelty
            });

            if (result.success && result.data) {

                setFullOrder(result.data);

                alert('Cambios guardados correctamente');
            } else {
                alert('Error al guardar cambios');
            }
        } catch (error) {
            console.error('Error saving management:', error);
            alert('Error al guardar cambios');
        } finally {
            setIsSaving(false);
        }
    };

    const formatCurrency = (amount: number | string, currency: string = 'USD', amountPresentment?: number, currencyPresentment?: string) => {
        const num = typeof amount === 'string' ? parseFloat(amount) : amount;
        if (isNaN(num) || num === undefined) return '-';

        if (amountPresentment && amountPresentment > 0 && currencyPresentment) {
            return new Intl.NumberFormat('es-CO', {
                style: 'currency',
                currency: currencyPresentment,
            }).format(amountPresentment);
        }

        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency || 'USD',
        }).format(num);
    };

    const formatDate = (dateString: string) => {
        if (!dateString) return '-';
        const date = new Date(dateString);
        const dateOnly = date.toLocaleDateString('es-CO');
        const timeOnly = date.toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
        return `${dateOnly} - ${timeOnly}`;
    };

    const formatDateSeparated = (dateString: string) => {
        if (!dateString) return { date: '-', time: '-' };
        const date = new Date(dateString);
        const dateOnly = date.toLocaleDateString('es-CO');
        const timeOnly = date.toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
        return { date: dateOnly, time: timeOnly };
    };

    const getTextColor = (bgColor: string): string => {
        const hex = bgColor.replace('#', '');
        const r = parseInt(hex.substr(0, 2), 16);
        const g = parseInt(hex.substr(2, 2), 16);
        const b = parseInt(hex.substr(4, 2), 16);
        const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
        return luminance > 0.5 ? '#000000' : '#FFFFFF';
    };

    const items = Array.isArray(order.items) ? order.items : [];

    const isReady = !loadingDetails && fullOrder;

    const orderItemsList: any[] = (order.order_items || items) as any[];

    const itemsTotals = orderItemsList.reduce(
        (acc: any, item: any) => {
            const unitPrice = item.unit_price || item.price || 0;
            const unitPricePresentment = item.unit_price_presentment || 0;
            const quantity = item.quantity || 0;
            const discount = item.discount || 0;
            const discountPresentment = item.discount_presentment || 0;
            acc.qty += quantity;
            acc.price += unitPrice * quantity;
            acc.pricePresentment += unitPricePresentment * quantity;
            acc.priceNoTax += (unitPrice / 1.19) * quantity;
            acc.priceNoTaxPresentment += (unitPricePresentment / 1.19) * quantity;
            acc.discounts += discount;
            acc.discountsPresentment += discountPresentment;
            acc.total += unitPrice * quantity - discount;
            acc.totalPresentment += unitPricePresentment * quantity - discountPresentment;
            return acc;
        },
        { qty: 0, price: 0, pricePresentment: 0, priceNoTax: 0, priceNoTaxPresentment: 0, discounts: 0, discountsPresentment: 0, total: 0, totalPresentment: 0 }
    );

    const totalItemsQty = itemsTotals.qty;

    const isCodOrder = order.is_cod === true || (order.cod_total || 0) > 0;
    const codCarrierFee = order.shipment?.cod_carrier_fee || 0;
    const codNet = order.cod_total || 0;
    const codToCollect = codNet + codCarrierFee;

    const paymentTone = (order.payment_details?.financial_status === 'paid' || order.is_paid)
        ? 'bg-emerald-600'
        : order.payment_details?.financial_status === 'refunded'
            ? 'bg-rose-600'
            : 'bg-amber-500';

    const customerInitials = (order.customer_name || '')
        .split(' ')
        .filter(Boolean)
        .slice(0, 2)
        .map((part: string) => part.charAt(0).toUpperCase())
        .join('') || '?';

    const handleCopyOrderNumber = () => {
        navigator.clipboard?.writeText(order.order_number || '').catch(() => {});
        setCopiedNumber(true);
        setTimeout(() => setCopiedNumber(false), 1500);
    };

    return (
        <div className={`flex h-full flex-col ${mode === 'details' ? '' : 'bg-white'}`}>

            {mode === 'recommendation' && businessColors && (
                <div
                    className="p-3 rounded-xl border shadow-sm relative overflow-hidden transition-all hover:shadow-md"
                    style={{
                        background: `linear-gradient(to right, ${quaternaryColor}, ${tertiaryColor})`,
                        borderColor: tertiaryColor
                    }}
                >
                    <div className="absolute top-0 right-0 p-2 opacity-5 pointer-events-none">
                        <svg className="w-24 h-24" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm1 15h-2v-2h2zm0-4h-2V7h2z" /></svg>
                    </div>

                    {onClose && (
                        <button
                            onClick={onClose}
                            className="absolute top-2 right-2 z-20 p-1 bg-white dark:bg-gray-800/20 hover:bg-white dark:bg-gray-800/40 rounded-full transition-colors backdrop-blur-sm"
                            style={{ color: primaryColor }}
                            title="Cerrar"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    )}

                    <div className="relative z-10">
                        <h3 className="text-lg font-bold flex items-center gap-2 mb-2" style={{ color: primaryColor }}>
                            <span className="text-2xl">🤖</span> Recomendación Inteligente
                        </h3>

                        {isReady ? (
                            <>
                                {loadingAI ? (
                                    <div className="flex items-center gap-3 bg-white dark:bg-gray-800/50 p-3 rounded-lg animate-pulse" style={{ color: secondaryColor }}>
                                        <div className="w-5 h-5 border-2 border-t-transparent rounded-full animate-spin" style={{ borderColor: secondaryColor, borderTopColor: 'transparent' }}></div>
                                        <span>Analizando mejores rutas y tarifas...</span>
                                    </div>
                                ) : aiRecommendation ? (
                                    <div className="flex flex-col gap-6">
                                        <div className="flex-1 space-y-4">
                                            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                                                <div>
                                                    <span
                                                        className="text-xs font-bold uppercase tracking-wider px-2 py-1 rounded"
                                                        style={{
                                                            color: secondaryColor,
                                                            backgroundColor: quaternaryColor + '80'
                                                        }}
                                                    >
                                                        Mejor Opción
                                                    </span>
                                                    <p className="text-4xl font-extrabold mt-2" style={{ color: secondaryColor }}>
                                                        {aiRecommendation.recommended_carrier}
                                                    </p>
                                                </div>
                                                <button
                                                    onClick={() => setShowGuideModal(true)}
                                                    className="text-white font-bold py-2 px-6 rounded-lg shadow-lg transition-all flex items-center gap-2"
                                                    style={{
                                                        backgroundColor: secondaryColor,
                                                        boxShadow: `0 0 0 20px ${secondaryColor}20`
                                                    }}
                                                >
                                                    <span>📦</span> Cotizar y Generar Guía
                                                </button>
                                            </div>

                                            <div className="bg-white dark:bg-gray-800/80 p-5 rounded-lg border text-gray-700 dark:text-gray-200 text-sm leading-relaxed shadow-sm" style={{ borderColor: tertiaryColor }}>
                                                <p className="font-semibold mb-1" style={{ color: primaryColor }}>Análisis:</p>
                                                {aiRecommendation.reasoning}
                                            </div>
                                        </div>

                                        {aiRecommendation.quotations && aiRecommendation.quotations.length > 0 && (
                                            <div className="pt-6 mt-2" style={{ borderTop: `1px solid ${tertiaryColor}` }}>
                                                <h4 className="text-sm font-bold uppercase tracking-wide mb-4 flex items-center gap-2" style={{ color: secondaryColor }}>
                                                    <span>📊</span> Cotizaciones Estimadas
                                                </h4>
                                                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                                                    {aiRecommendation.quotations.map((quote, idx) => {
                                                        const isRecommended = quote.carrier === aiRecommendation.recommended_carrier;
                                                        return (
                                                            <div
                                                                key={idx}
                                                                className={`p-4 rounded-xl border flex flex-col justify-between transition-all hover:shadow-md ${isRecommended ? 'bg-white dark:bg-gray-800 shadow-sm relative overflow-hidden' : 'bg-slate-50 dark:bg-gray-800 hover:bg-white dark:hover:bg-gray-700'}`}
                                                                style={{
                                                                    borderColor: isRecommended ? secondaryColor : tertiaryColor,
                                                                    ...(isRecommended && {
                                                                        boxShadow: `0 0 0 1px ${secondaryColor}20`
                                                                    })
                                                                }}
                                                            >
                                                                {isRecommended && <div className="absolute top-0 right-0 text-white text-[10px] px-2 py-0.5 rounded-bl-lg font-bold" style={{ backgroundColor: secondaryColor }}>RECOMENDADO</div>}
                                                                <div>
                                                                    <p className="font-bold text-lg" style={{ color: isRecommended ? secondaryColor : '#374151' }}>
                                                                        {quote.carrier}
                                                                    </p>
                                                                    <p className="text-xs text-gray-500 dark:text-gray-400 flex items-center gap-1 mt-1">
                                                                        <span>⏱️</span> {quote.estimated_delivery_days} días hábiles
                                                                    </p>
                                                                </div>
                                                                <div className="mt-4 pt-3 border-t border-gray-100 dark:border-gray-700">
                                                                    <p className="text-gray-500 dark:text-gray-400 text-xs uppercase mb-0.5">Costo Estimado</p>
                                                                    <div className="flex justify-between items-end">
                                                                        <p className="font-bold text-xl" style={{ color: isRecommended ? secondaryColor : '#4b5563' }}>
                                                                            {formatCurrency(quote.estimated_cost, 'COP')}
                                                                        </p>
                                                                        <button
                                                                            onClick={() => setShowGuideModal(true)}
                                                                            className="text-xs bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200 px-2 py-1 rounded"
                                                                        >
                                                                            Elegir
                                                                        </button>
                                                                    </div>
                                                                </div>
                                                            </div>
                                                        );
                                                    })}
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                ) : (
                                    <div className="text-sm text-gray-500 dark:text-gray-400 italic bg-gray-50 dark:bg-gray-800 p-3 rounded border border-gray-100 dark:border-gray-700">
                                        No hay recomendación disponible. Verifique que la orden tenga dirección completa (Ciudad y Departamento).
                                    </div>
                                )}
                            </>
                        ) : (
                            <div className="text-sm mt-1 animate-pulse" style={{ color: tertiaryColor }}>Cargando datos de orden...</div>
                        )}

                        {showGuideModal && order && (
                            <ShipmentGuideModal
                                isOpen={showGuideModal}
                                onClose={() => setShowGuideModal(false)}
                                order={order}
                                recommendedCarrier={aiRecommendation?.recommended_carrier}
                                onGuideGenerated={() => {

                                    setShowGuideModal(false);
                                    fetchDetails();
                                }}
                            />
                        )}
                    </div>
                </div>
            )}

            {mode === 'details' && (
                <div className="flex h-full flex-col bg-slate-100 dark:bg-gray-900">
                    <header className="shrink-0 border-b border-slate-200 bg-white dark:border-gray-700 dark:bg-gray-800">
                        <div className="flex flex-wrap items-center gap-2.5 px-4 py-2">
                            <span
                                className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-white"
                                style={{ backgroundColor: primaryColor }}
                            >
                                <Package className="h-4 w-4" />
                            </span>
                            <div className="min-w-0">
                                <div className="flex items-center gap-1.5">
                                    <h2 className="truncate text-base font-bold tracking-tight text-slate-900 dark:text-white">
                                        Orden {order.order_number || '-'}
                                    </h2>
                                    <button
                                        onClick={handleCopyOrderNumber}
                                        title={'Copiar n\u00famero de orden'}
                                        className="rounded-md p-1 text-slate-400 transition hover:bg-slate-100 hover:text-slate-600 dark:hover:bg-gray-700"
                                    >
                                        {copiedNumber ? <Check className="h-4 w-4 text-emerald-500" /> : <Copy className="h-4 w-4" />}
                                    </button>
                                </div>
                                <p className="truncate font-mono text-[11px] text-slate-400">{order.internal_number || '-'}</p>
                            </div>

                            {order.channel_pack_id && (
                                <span
                                    className="inline-flex items-center gap-1 rounded-full border border-amber-300 bg-amber-100 px-2.5 py-0.5 text-[11px] font-bold uppercase text-amber-700"
                                    title={`Carrito de MercadoLibre (pack ${order.channel_pack_id}): agrupa varias \u00f3rdenes del canal en un solo env\u00edo`}
                                >
                                    <ShoppingBag className="h-3 w-3" />
                                    Pack
                                </span>
                            )}

                            {order.order_status?.color ? (
                                <span
                                    className="rounded-full px-2.5 py-0.5 text-[11px] font-bold"
                                    style={{
                                        backgroundColor: order.order_status.color,
                                        color: getTextColor(order.order_status.color),
                                    }}
                                >
                                    {order.order_status.name || order.status}
                                </span>
                            ) : (
                                <span className="rounded-full px-2.5 py-0.5 text-[11px] font-bold text-white" style={{ backgroundColor: secondaryColor }}>
                                    {order.order_status?.name || order.status || '-'}
                                </span>
                            )}

                            <span
                                className={`inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-[11px] font-bold uppercase ${isCodOrder
                                    ? 'border-amber-300 bg-amber-100 text-amber-700'
                                    : 'border-slate-200 bg-slate-100 text-slate-500'
                                    }`}
                            >
                                <Wallet className="h-3 w-3" />
                                {isCodOrder ? 'Contra entrega' : 'Pago anticipado'}
                            </span>

                            <span
                                className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-[11px] font-bold uppercase text-white ${paymentTone}`}
                            >
                                <CreditCard className="h-3 w-3" />
                                {order.payment_details?.financial_status || (order.is_paid ? 'PAID' : 'PENDING')}
                            </span>

                            <div className="ml-auto flex items-center gap-3">
                                <span className="hidden items-center gap-1.5 text-xs text-slate-400 md:flex">
                                    <Calendar className="h-3.5 w-3.5" />
                                    {formatDate(order.occurred_at || order.created_at)}
                                </span>
                                <span className="flex items-center gap-1.5">
                                    {order.integration_logo_url && (
                                        <img src={order.integration_logo_url} alt={order.platform} className="h-5 w-5 object-contain" />
                                    )}
                                    <span className="text-xs font-semibold text-slate-600 dark:text-gray-300">{salesChannelLabel(order.platform)}</span>
                                </span>
                                {onClose && (
                                    <button
                                        onClick={onClose}
                                        title="Cerrar"
                                        className="rounded-full p-1.5 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700 dark:hover:bg-gray-700"
                                    >
                                        <X className="h-5 w-5" />
                                    </button>
                                )}
                            </div>
                        </div>
                    </header>

                    <div className="flex-1 overflow-y-auto">
                        <main className="grid grid-cols-1 gap-3 p-3 xl:grid-cols-[minmax(0,1fr)_300px]">
                            <div className="min-w-0 space-y-3">
                                <div className={`grid grid-cols-2 gap-px overflow-hidden rounded-xl border border-slate-200 bg-slate-200 dark:border-gray-700 dark:bg-gray-700 ${isCodOrder ? 'lg:grid-cols-3 xl:grid-cols-6' : 'lg:grid-cols-4'}`}>
                                    <div className="bg-white p-2.5 dark:bg-gray-800">
                                        <p className="flex items-center gap-1.5 text-[10px] font-semibold uppercase tracking-wide text-slate-400">
                                            <Receipt className="h-3.5 w-3.5" /> Subtotal
                                        </p>
                                        <p className="mt-1 text-base font-bold tabular-nums text-slate-800 dark:text-white">
                                            {formatCurrency(order.subtotal, order.currency, order.subtotal_presentment, order.currency_presentment)}
                                        </p>
                                        {order.tax > 0 && ['shopify', 'amazon', 'mercadolibre'].includes(order.platform?.toLowerCase() || '') && (
                                            <IVAIncludedBadge
                                                ivaAmount={order.tax}
                                                currency={order.currency}
                                                currencyPresentment={order.currency_presentment}
                                                amountPresentment={order.tax_presentment}
                                            />
                                        )}
                                    </div>
                                    <div className="bg-white p-2.5 dark:bg-gray-800">
                                        <p className="flex items-center gap-1.5 text-[10px] font-semibold uppercase tracking-wide text-slate-400">
                                            <Truck className="h-3.5 w-3.5" /> Env&iacute;o
                                        </p>
                                        <p className="mt-1 text-base font-bold tabular-nums text-slate-800 dark:text-white">
                                            {formatCurrency(
                                                (order.shipment?.total_cost ?? order.shipping_cost ?? 0) - (order.shipping_discount ?? 0),
                                                order.currency,
                                                (order.shipment?.total_cost ?? order.shipping_cost_presentment ?? 0) - (order.shipping_discount_presentment ?? 0),
                                                order.currency_presentment
                                            )}
                                        </p>
                                        {((order.shipment?.total_cost ?? order.shipping_cost ?? 0) - (order.shipping_discount ?? 0)) > 0 && ['shopify', 'amazon', 'mercadolibre'].includes(order.platform?.toLowerCase() || '') && (
                                            <IVAIncludedBadge
                                                ivaAmount={((order.shipment?.total_cost ?? order.shipping_cost ?? 0) - (order.shipping_discount ?? 0)) / 1.19 * 0.19}
                                                currency={order.currency}
                                                currencyPresentment={order.currency_presentment}
                                                amountPresentment={((order.shipment?.total_cost ?? order.shipping_cost_presentment ?? 0) - (order.shipping_discount_presentment ?? 0)) / 1.19 * 0.19}
                                            />
                                        )}
                                    </div>
                                    <div className="bg-white p-2.5 dark:bg-gray-800">
                                        <p className="flex items-center gap-1.5 text-[10px] font-semibold uppercase tracking-wide text-slate-400">
                                            <Percent className="h-3.5 w-3.5" /> Rete fuente
                                        </p>
                                        <p className="mt-1 text-base font-bold tabular-nums text-slate-800 dark:text-white">
                                            {order.invoice?.retention_amount ? formatCurrency(order.invoice.retention_amount, order.currency) : '-'}
                                        </p>
                                    </div>
                                    <div className="p-2.5" style={{ backgroundColor: primaryColor }}>
                                        <p className="flex items-center gap-1.5 text-[10px] font-semibold uppercase tracking-wide text-white/60">
                                            <Wallet className="h-3.5 w-3.5" /> {isCodOrder ? 'Total para el negocio' : 'Total'}
                                        </p>
                                        <p className="mt-1 text-lg font-bold tabular-nums text-white">
                                            {isCodOrder && codNet > 0
                                                ? formatCurrency(codNet, order.currency)
                                                : formatCurrency(
                                                    (order.subtotal || 0) + ((order.shipment?.total_cost ?? order.shipping_cost ?? 0) - (order.shipping_discount ?? 0)) - (order.invoice?.retention_amount || 0),
                                                    order.currency,
                                                    (order.subtotal_presentment || 0) + ((order.shipment?.total_cost ?? order.shipping_cost_presentment ?? 0) - (order.shipping_discount_presentment ?? 0)) - (order.invoice?.retention_amount || 0),
                                                    order.currency_presentment
                                                )}
                                        </p>
                                        {isCodOrder && (
                                            <p className="mt-1 text-[10px] leading-snug text-white/70">
                                                Neto a pagar al negocio, sin comisi&oacute;n
                                            </p>
                                        )}
                                    </div>

                                    {isCodOrder && (
                                        <div className="bg-white p-2.5 dark:bg-gray-800">
                                            <p className="flex items-center gap-1.5 text-[10px] font-semibold uppercase tracking-wide text-slate-400">
                                                <Truck className="h-3.5 w-3.5" /> Comisi&oacute;n carrier
                                            </p>
                                            <p className="mt-1 text-base font-bold tabular-nums text-slate-800 dark:text-white">
                                                {codCarrierFee > 0 ? formatCurrency(codCarrierFee, order.currency) : 'Pendiente'}
                                            </p>
                                            <p className="text-[10px] leading-snug text-slate-400">
                                                {codCarrierFee > 0
                                                    ? `${order.shipment?.carrier || 'Transportadora'} - contra entrega`
                                                    : 'Se liquida al confirmar el recaudo'}
                                            </p>
                                        </div>
                                    )}
                                    {isCodOrder && (
                                        <div className="bg-amber-50 p-2.5 dark:bg-amber-900/20">
                                            <p className="flex items-center gap-1.5 text-[10px] font-semibold uppercase tracking-wide text-amber-600">
                                                <Truck className="h-3.5 w-3.5" /> Cobro al cliente final
                                            </p>
                                            <p className="mt-1 text-lg font-bold tabular-nums text-amber-700 dark:text-amber-300">
                                                {formatCurrency(codToCollect, order.currency)}
                                            </p>
                                            <p className="mt-1 text-[10px] leading-snug text-amber-600/80">
                                                {codCarrierFee > 0
                                                    ? `Lo cobra ${order.shipment?.carrier || 'la transportadora'} al entregar`
                                                    : 'Falta sumar la comisi\u00f3n por liquidar'}
                                            </p>
                                        </div>
                                    )}
                                </div>

                                <section className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm dark:border-gray-700 dark:bg-gray-800">
                                    <div className="flex items-center justify-between border-b border-slate-100 px-3 py-2 dark:border-gray-700">
                                        <h3 className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">
                                            <ShoppingBag className="h-4 w-4 text-slate-400" /> Productos del pedido
                                        </h3>
                                        <span className="rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-semibold tabular-nums text-slate-600 dark:bg-gray-700 dark:text-gray-200">
                                            {totalItemsQty}{totalItemsQty === 1 ? ' \u00edtem' : ' \u00edtems'}
                                        </span>
                                    </div>
                                    {loadingDetails ? (
                                        <div className="py-6 text-center text-xs text-slate-400">Cargando productos...</div>
                                    ) : (order.order_items || items).length > 0 ? (
                                        <div className="overflow-x-auto">
                                            <table className="w-full text-xs">
                                                <thead className="bg-slate-50 dark:bg-gray-900/50">
                                                    <tr className="text-[11px] uppercase tracking-wide text-slate-400">
                                                        <th className="px-3 py-2 text-left font-semibold">Producto</th>
                                                        <th className="px-3 py-2 text-left font-semibold">SKU</th>
                                                        <th className="px-3 py-2 text-right font-semibold">Cant</th>
                                                        <th className="px-3 py-2 text-right font-semibold">Precio</th>
                                                        <th className="px-3 py-2 text-right font-semibold">Precio s/IVA</th>
                                                        <th className="px-3 py-2 text-right font-semibold">Desc.</th>
                                                        <th className="px-3 py-2 text-right font-semibold">Total</th>
                                                    </tr>
                                                </thead>
                                                <tbody className="divide-y divide-slate-100 dark:divide-gray-700">
                                                    {(order.order_items || items).map((item: any, idx: number) => {
                                                        const unitPrice = item.unit_price || item.price || 0;
                                                        const unitPricePresentment = item.unit_price_presentment || 0;
                                                        const quantity = item.quantity || 0;
                                                        const discount = item.discount || 0;
                                                        const discountPresentment = item.discount_presentment || 0;
                                                        const priceWithoutTax = unitPrice / 1.19;
                                                        const priceWithoutTaxPresentment = unitPricePresentment / 1.19;
                                                        const itemTotal = (unitPrice * quantity) - discount;
                                                        const itemTotalPresentment = (unitPricePresentment * quantity) - discountPresentment;

                                                        return (
                                                            <tr key={idx} className="transition hover:bg-slate-50 dark:hover:bg-gray-700/50">
                                                                <td className="px-3 py-1.5">
                                                                    <div className="flex items-center gap-3">
                                                                        <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-slate-100 text-slate-400 dark:bg-gray-700">
                                                                            <Package className="h-4 w-4" />
                                                                        </span>
                                                                        <div className="min-w-0">
                                                                            <p className="font-medium text-slate-800 dark:text-white">{item.product_name || item.name || item.title || '-'}</p>
                                                                            {item.variant_label ? <p className="text-[11px] text-slate-400">{item.variant_label}</p> : null}
                                                                        </div>
                                                                    </div>
                                                                </td>
                                                                <td className="px-3 py-1.5 font-mono text-[11px] text-slate-500 dark:text-gray-300">{item.product_sku || item.sku || '-'}</td>
                                                                <td className="px-3 py-1.5 text-right tabular-nums text-slate-700 dark:text-gray-200">{quantity}</td>
                                                                <td className="px-3 py-1.5 text-right tabular-nums text-slate-700 dark:text-gray-200">{formatCurrency(unitPrice, order.currency, unitPricePresentment, order.currency_presentment)}</td>
                                                                <td className="px-3 py-1.5 text-right tabular-nums text-slate-400">{formatCurrency(priceWithoutTax, order.currency, priceWithoutTaxPresentment, order.currency_presentment)}</td>
                                                                <td className="px-3 py-1.5 text-right tabular-nums">
                                                                    {(discount > 0 || discountPresentment > 0) ? (
                                                                        <span className="font-semibold text-emerald-600">-{formatCurrency(discount, order.currency, discountPresentment, order.currency_presentment)}</span>
                                                                    ) : (
                                                                        <span className="text-slate-300">-</span>
                                                                    )}
                                                                </td>
                                                                <td className="px-3 py-1.5 text-right font-semibold tabular-nums" style={{ color: primaryColor }}>
                                                                    {formatCurrency(itemTotal, order.currency, itemTotalPresentment, order.currency_presentment)}
                                                                </td>
                                                            </tr>
                                                        );
                                                    })}
                                                    <tr className="bg-slate-50 font-semibold dark:bg-gray-900/50">
                                                        <td className="px-3 py-2"></td>
                                                        <td className="px-3 py-2 text-right text-[11px] uppercase tracking-wide text-slate-400">Total</td>
                                                        <td className="px-3 py-2 text-right tabular-nums" style={{ color: primaryColor }}>{itemsTotals.qty}</td>
                                                        <td className="px-3 py-2 text-right tabular-nums" style={{ color: primaryColor }}>{formatCurrency(itemsTotals.price, order.currency, itemsTotals.pricePresentment, order.currency_presentment)}</td>
                                                        <td className="px-3 py-2 text-right tabular-nums" style={{ color: primaryColor }}>{formatCurrency(itemsTotals.priceNoTax, order.currency, itemsTotals.priceNoTaxPresentment, order.currency_presentment)}</td>
                                                        <td className="px-3 py-2 text-right tabular-nums" style={{ color: primaryColor }}>{formatCurrency(itemsTotals.discounts, order.currency, itemsTotals.discountsPresentment, order.currency_presentment)}</td>
                                                        <td className="px-3 py-2 text-right tabular-nums" style={{ color: primaryColor }}>{formatCurrency(itemsTotals.total, order.currency, itemsTotals.totalPresentment, order.currency_presentment)}</td>
                                                    </tr>
                                                </tbody>
                                            </table>
                                        </div>
                                    ) : (
                                        <div className="flex flex-col items-center gap-2 py-8 text-slate-400">
                                            <ShoppingBag className="h-7 w-7" />
                                            <p className="text-sm">No hay informaci&oacute;n de productos.</p>
                                        </div>
                                    )}
                                </section>

                                <div className="grid grid-cols-1 gap-3 lg:grid-cols-2">
                                    <section className="rounded-xl border border-slate-200 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                                        <h3 className="mb-2.5 flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">
                                            <User className="h-4 w-4 text-slate-400" /> Cliente
                                        </h3>
                                        {loadingDetails ? (
                                            <div className="py-4 text-center text-xs text-slate-400">Cargando...</div>
                                        ) : (
                                            <div className="space-y-4">
                                                <div className="flex items-center gap-3">
                                                    <span
                                                        className="flex h-9 w-9 items-center justify-center rounded-full text-xs font-bold text-white"
                                                        style={{ backgroundColor: primaryColor }}
                                                    >
                                                        {customerInitials}
                                                    </span>
                                                    <div className="min-w-0">
                                                        <p className="truncate text-sm font-semibold text-slate-800 dark:text-white">{order.customer_name || '-'}</p>
                                                        <p className="text-xs text-slate-400">DNI {order.customer_dni || '-'}</p>
                                                    </div>
                                                </div>
                                                <div className="space-y-2 border-t border-slate-100 pt-3 dark:border-gray-700">
                                                    <div className="flex items-center gap-2.5 text-xs text-slate-600 dark:text-gray-200">
                                                        <Phone className="h-4 w-4 shrink-0 text-slate-400" />
                                                        <span className="tabular-nums">{order.customer_phone || 'Sin tel\u00e9fono registrado'}</span>
                                                    </div>
                                                    <div className="flex items-center gap-2.5 text-xs text-slate-600 dark:text-gray-200">
                                                        <Mail className="h-4 w-4 shrink-0 text-slate-400" />
                                                        <span className="break-all">{order.customer_email || 'Sin correo registrado'}</span>
                                                    </div>
                                                </div>
                                            </div>
                                        )}
                                    </section>

                                    <section className="rounded-xl border border-slate-200 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                                        <div className="mb-3 flex items-center justify-between gap-2">
                                            <h3 className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">
                                                <MapPin className="h-4 w-4 text-slate-400" /> Direcci&oacute;n de entrega
                                            </h3>
                                            {order.shipping_geo_confidence && (
                                                <span
                                                    className="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
                                                    style={
                                                        order.shipping_geo_confidence === 'high'
                                                            ? { backgroundColor: '#dcfce7', color: '#166534' }
                                                            : order.shipping_geo_confidence === 'medium'
                                                                ? { backgroundColor: '#fef9c3', color: '#854d0e' }
                                                                : { backgroundColor: '#fee2e2', color: '#991b1b' }
                                                    }
                                                    title={'Confianza del geocode de la direcci\u00f3n'}
                                                >
                                                    {order.shipping_geo_confidence === 'high'
                                                        ? 'Direcci\u00f3n confiable'
                                                        : order.shipping_geo_confidence === 'medium'
                                                            ? 'Verificar direcci\u00f3n'
                                                            : 'Direcci\u00f3n dudosa'}
                                                </span>
                                            )}
                                        </div>
                                        {loadingDetails ? (
                                            <div className="py-4 text-center text-xs text-slate-400">Cargando...</div>
                                        ) : (
                                            <div className="space-y-2.5">
                                                <div>
                                                    <p className="text-sm font-medium text-slate-800 dark:text-white">{order.shipping_street || '-'}</p>
                                                    <p className="text-xs text-slate-400">
                                                        {(() => {
                                                            const resolved = resolveCityState(order.shipping_city || '', order.shipping_state || '', order.destination_dane_code);
                                                            const cityName = resolved?.ciudad || order.shipping_city || '';
                                                            const stateName = resolved?.departamento || order.shipping_state || '';
                                                            return `${cityName}${stateName ? ', ' + stateName : ''}${order.shipping_postal_code ? ' ' + order.shipping_postal_code : ''}`;
                                                        })()}
                                                    </p>
                                                </div>
                                                {order.business_id && order.business_id > 0 && order.id && (
                                                    <div className="overflow-hidden rounded-xl border border-slate-200 dark:border-gray-700">
                                                        <GeozoneMiniMap businessId={order.business_id} orderId={order.id} lat={order.shipping_lat ?? null} lng={order.shipping_lng ?? null} height="190px" />
                                                    </div>
                                                )}
                                            </div>
                                        )}
                                    </section>
                                </div>

                                <section className="rounded-xl border border-slate-200 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                                    <h3 className="mb-2.5 flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">
                                        <ClipboardList className="h-4 w-4 text-slate-400" /> Gesti&oacute;n y novedades
                                    </h3>
                                    <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                                        <div className="space-y-1.5">
                                            <label className="text-xs font-medium text-slate-500 dark:text-gray-300">Confirmaci&oacute;n de pedido</label>
                                            <select
                                                className={`w-full rounded-lg border px-3 py-2 text-sm font-semibold shadow-sm outline-none transition ${isConfirmed === true
                                                    ? 'border-emerald-400 bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300'
                                                    : isConfirmed === false
                                                        ? 'border-rose-400 bg-rose-50 text-rose-700 dark:bg-rose-900/20 dark:text-rose-300'
                                                        : 'border-amber-400 bg-amber-50 text-amber-700 dark:bg-amber-900/20 dark:text-amber-300'
                                                    }`}
                                                value={isConfirmed === null ? 'pending' : (isConfirmed ? 'yes' : 'no')}
                                                onChange={(e) => {
                                                    const val = e.target.value;
                                                    setIsConfirmed(val === 'pending' ? null : val === 'yes');
                                                }}
                                            >
                                                <option value="yes">S&iacute;, Confirmado</option>
                                                <option value="no">No, Rechazado/Cancelado</option>
                                                <option value="pending">Pendiente confirmaci&oacute;n</option>
                                            </select>
                                            {isConfirmed === false && (
                                                <p className="flex items-center gap-1.5 text-xs text-rose-500">
                                                    <AlertTriangle className="h-3.5 w-3.5" /> La orden quedar&aacute; marcada como cancelada.
                                                </p>
                                            )}
                                        </div>
                                        <div className="space-y-1.5">
                                            <label className="text-xs font-medium text-slate-500 dark:text-gray-300">Novedades / Notas</label>
                                            <textarea
                                                rows={2}
                                                className="w-full resize-none rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 shadow-sm outline-none transition placeholder:text-slate-400 focus:border-slate-400 dark:border-gray-600 dark:bg-gray-900 dark:text-white"
                                                placeholder={'Escribe aqu\u00ed novedades (ej: cambio de direcci\u00f3n, cliente contactado, etc.)'}
                                                value={novelty}
                                                onChange={(e) => setNovelty(e.target.value)}
                                            />
                                        </div>
                                    </div>
                                    <div className="mt-3 flex flex-col gap-2 sm:flex-row">
                                        <button
                                            onClick={handleSaveManagement}
                                            disabled={isSaving}
                                            className="flex flex-1 items-center justify-center gap-1.5 rounded-lg px-4 py-2 text-sm font-semibold text-white transition disabled:opacity-50"
                                            style={{ backgroundColor: primaryColor }}
                                        >
                                            <Save className="h-4 w-4" />
                                            {isSaving ? 'Guardando...' : 'Guardar gesti\u00f3n'}
                                        </button>
                                        {hasWhatsApp && order.is_confirmed !== true && (
                                            <button
                                                onClick={handleWhatsAppConfirmation}
                                                disabled={isSendingWhatsApp || whatsAppSent}
                                                title={!order.customer_phone ? 'La orden no tiene tel\u00e9fono de cliente' : 'Enviar confirmaci\u00f3n por WhatsApp'}
                                                className={`flex items-center justify-center gap-1.5 rounded-lg border px-4 py-2 text-sm font-semibold transition ${whatsAppSent
                                                    ? 'border-emerald-300 bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300'
                                                    : 'border-emerald-300 text-emerald-700 hover:bg-emerald-50 dark:text-emerald-300'
                                                    }`}
                                            >
                                                <MessageCircle className="h-4 w-4" />
                                                {isSendingWhatsApp ? 'Enviando...' : whatsAppSent ? 'Enviado' : 'Contactar por WhatsApp'}
                                            </button>
                                        )}
                                    </div>
                                </section>

                                <section className="rounded-xl border border-slate-200 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                                    <h3 className="mb-2.5 flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">
                                        <History className="h-4 w-4 text-slate-400" /> Historial de estados
                                    </h3>
                                    {loadingHistory ? (
                                        <div className="py-4 text-center text-xs text-slate-400">Cargando historial...</div>
                                    ) : statusHistory.length === 0 ? (
                                        <div className="flex flex-col items-center gap-2 rounded-xl border border-dashed border-slate-300 py-8 text-slate-400 dark:border-gray-600">
                                            <History className="h-7 w-7" />
                                            <p className="text-sm">No hay cambios de estado registrados.</p>
                                        </div>
                                    ) : (
                                        <div className="relative">
                                            <div className="absolute left-0 right-0 top-5 h-0.5 bg-slate-200 dark:bg-gray-600"></div>
                                            <div className="flex gap-0 overflow-x-auto pb-2" style={{ scrollbarWidth: 'thin' }}>
                                                <div className="relative flex flex-shrink-0 flex-col items-center" style={{ minWidth: '100px' }}>
                                                    <div className="z-10 flex h-8 w-8 items-center justify-center rounded-full border-2 border-white bg-slate-300 shadow-sm dark:border-gray-800 dark:bg-gray-600">
                                                        <Plus className="h-4 w-4 text-white" />
                                                    </div>
                                                    <span
                                                        className="mt-1.5 rounded-full bg-slate-100 px-2 py-0.5 text-[10px] font-semibold text-slate-600 dark:bg-gray-700 dark:text-gray-300"
                                                        title={statusHistory[0].previous_status || ''}
                                                    >
                                                        {getOrderStatusLabel(statusHistory[0].previous_status) || 'Creada'}
                                                    </span>
                                                </div>
                                                {statusHistory.map((entry, idx) => {
                                                    const isLast = idx === statusHistory.length - 1;
                                                    return (
                                                        <div key={entry.id} className="relative flex flex-shrink-0 flex-col items-center" style={{ minWidth: '104px' }}>
                                                            <div className="absolute -left-3 top-3 text-slate-300 dark:text-gray-600">
                                                                <ChevronRight className="h-4 w-4" />
                                                            </div>
                                                            <div
                                                                className="z-10 flex h-8 w-8 items-center justify-center rounded-full border-2 border-white shadow-sm dark:border-gray-800"
                                                                style={{ backgroundColor: isLast ? secondaryColor : tertiaryColor }}
                                                            >
                                                                {isLast ? <Check className="h-4 w-4 text-white" /> : <ArrowRight className="h-4 w-4 text-white" />}
                                                            </div>
                                                            <span
                                                                className="mt-1.5 rounded-full px-2 py-0.5 text-center text-[10px] font-semibold leading-tight"
                                                                style={{
                                                                    backgroundColor: isLast ? quaternaryColor : tertiaryColor + '40',
                                                                    color: isLast ? secondaryColor : primaryColor,
                                                                }}
                                                                title={entry.new_status}
                                                            >
                                                                {getOrderStatusLabel(entry.new_status) || entry.new_status}
                                                            </span>
                                                            {isPaymentStatus(entry.new_status) && (
                                                                <span className="mt-1 rounded-full bg-slate-100 px-1.5 py-px text-[9px] font-semibold text-slate-500 dark:bg-gray-700 dark:text-gray-300">
                                                                    pago
                                                                </span>
                                                            )}
                                                            <p className="mt-1 text-center text-[10px] leading-tight text-slate-400">
                                                                {new Date(entry.created_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short' })}
                                                                {' '}
                                                                {new Date(entry.created_at).toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' })}
                                                            </p>
                                                            {entry.changed_by_name && (
                                                                <span
                                                                    className={`mt-1 max-w-[112px] truncate rounded-full px-1.5 py-px text-[9px] font-semibold leading-tight ${statusSourceStyle(entry.source)}`}
                                                                    title={[
                                                                        `${entry.source === 'user' ? 'Registrado por' : 'Registrado desde'}: ${entry.changed_by_name}`,
                                                                        entry.source ? `Origen: ${SOURCE_LABELS[entry.source] || entry.source}` : '',
                                                                        entry.reason ? `Motivo: ${entry.reason}` : '',
                                                                    ].filter(Boolean).join('\n')}
                                                                >
                                                                    {entry.changed_by_name}
                                                                </span>
                                                            )}
                                                        </div>
                                                    );
                                                })}
                                            </div>
                                        </div>
                                    )}
                                </section>
                            </div>

                            <aside className="space-y-3 xl:sticky xl:top-0 xl:self-start">
                                <section className="rounded-xl border border-slate-200 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                                    <h3 className="mb-2.5 flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">
                                        <CircleCheck className="h-4 w-4 text-slate-400" /> Estado de la orden
                                    </h3>
                                    <div className="space-y-2.5">
                                        <div className="flex items-center justify-between rounded-lg bg-slate-50 px-2.5 py-2 dark:bg-gray-900/50">
                                            <span className="flex items-center gap-2 text-xs font-medium text-slate-600 dark:text-gray-200">
                                                <Clock className="h-4 w-4 text-slate-400" /> Estado
                                            </span>
                                            {order.order_status?.color ? (
                                                <span
                                                    className="rounded-full px-2.5 py-0.5 text-xs font-bold"
                                                    style={{ backgroundColor: order.order_status.color, color: getTextColor(order.order_status.color) }}
                                                >
                                                    {order.order_status.name || order.status}
                                                </span>
                                            ) : (
                                                <span className="rounded-full px-2.5 py-0.5 text-xs font-bold text-white" style={{ backgroundColor: secondaryColor }}>
                                                    {order.status || '-'}
                                                </span>
                                            )}
                                        </div>
                                        <div className="flex items-center justify-between rounded-lg bg-slate-50 px-2.5 py-2 dark:bg-gray-900/50">
                                            <span className="flex items-center gap-2 text-xs font-medium text-slate-600 dark:text-gray-200">
                                                <CreditCard className="h-4 w-4 text-slate-400" /> Pago
                                            </span>
                                            <span className={`rounded-full px-2.5 py-0.5 text-xs font-bold uppercase text-white ${paymentTone}`}>
                                                {order.payment_details?.financial_status || (order.is_paid ? 'PAID' : 'PENDING')}
                                            </span>
                                        </div>
                                        {isCodOrder && (
                                            <div className="flex items-center justify-between rounded-lg bg-slate-50 px-2.5 py-2 dark:bg-gray-900/50">
                                                <span className="flex items-center gap-2 text-xs font-medium text-slate-600 dark:text-gray-200">
                                                    <Scissors className="h-4 w-4 text-slate-400" /> Corte de pago
                                                </span>
                                                <span
                                                    className={`rounded-full px-2.5 py-0.5 text-xs font-bold ${order.cod_cut_confirmed
                                                        ? 'bg-emerald-100 text-emerald-700'
                                                        : 'bg-amber-100 text-amber-700'
                                                        }`}
                                                >
                                                    {order.cod_cut_confirmed ? 'Confirmado' : 'Sin confirmar'}
                                                </span>
                                            </div>
                                        )}
                                        {order.paid_at && (
                                            <div className="flex items-center justify-between px-1 text-xs">
                                                <span className="text-slate-400">Fecha de pago</span>
                                                <span className="tabular-nums text-slate-600 dark:text-gray-300">{formatDate(order.paid_at)}</span>
                                            </div>
                                        )}
                                        {!isTerminalStatus(order.order_status?.code || order.status || '') && (
                                            <button
                                                onClick={() => setShowChangeStatus(true)}
                                                className="flex w-full items-center justify-center gap-1.5 rounded-lg border border-slate-200 px-3 py-2 text-sm font-semibold text-slate-600 transition hover:bg-slate-50 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-700"
                                            >
                                                <RefreshCw className="h-4 w-4" /> Cambiar estado
                                            </button>
                                        )}
                                    </div>
                                </section>

                                <section className="rounded-xl border border-slate-200 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                                    <h3 className="mb-2.5 flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">
                                        <Clock className="h-4 w-4 text-slate-400" /> Cronolog&iacute;a
                                    </h3>
                                    <ol className="relative space-y-3.5 border-l border-slate-200 pl-4 dark:border-gray-700">
                                        {[
                                            { label: 'Creado (DB)', time: formatDate(order.created_at), tone: 'bg-emerald-100 text-emerald-600' },
                                            { label: 'Importado', time: formatDate(order.imported_at), tone: 'bg-violet-100 text-violet-600' },
                                            ...(order.updated_at ? [{ label: 'Actualizado', time: formatDate(order.updated_at), tone: 'bg-sky-100 text-sky-600' }] : []),
                                        ].map((event) => (
                                            <li key={event.label} className="relative">
                                                <span className={`absolute -left-[27px] flex h-5 w-5 items-center justify-center rounded-full ring-4 ring-white dark:ring-gray-800 ${event.tone}`}>
                                                    <Clock className="h-3 w-3" />
                                                </span>
                                                <p className="text-[13px] font-semibold text-slate-800 dark:text-white">{event.label}</p>
                                                <p className="text-xs tabular-nums text-slate-400">{event.time}</p>
                                            </li>
                                        ))}
                                    </ol>
                                </section>

                                <section className="rounded-xl border border-slate-200 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                                    <h3 className="mb-2.5 flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">
                                        <FileText className="h-4 w-4 text-slate-400" /> Informaci&oacute;n general
                                    </h3>
                                    <div className="space-y-2.5">
                                        <div className="flex items-center justify-between gap-2">
                                            <div>
                                                <p className="text-[11px] font-medium uppercase tracking-wide text-slate-400">Factura</p>
                                                <p className="mt-0.5 text-[13px] font-medium text-slate-800 dark:text-white">
                                                    {order.invoice && order.invoice.status === 'issued' ? 'Emitida' : 'Sin emitir'}
                                                </p>
                                                {order.invoice?.id ? (
                                                    <Link
                                                        href={`/invoicing/invoices?invoice_id=${order.invoice.id}`}
                                                        title="Ver esta factura en el m&oacute;dulo de facturaci&oacute;n"
                                                        className="mt-0.5 inline-block text-[13px] font-semibold text-[var(--color-primary)] hover:underline"
                                                    >
                                                        {order.invoice.invoice_number || `ID ${order.invoice.id}`}
                                                    </Link>
                                                ) : null}
                                            </div>
                                            <span
                                                className={`rounded-full px-2.5 py-1 text-[10px] font-bold ${order.invoice && order.invoice.status === 'issued'
                                                    ? 'bg-emerald-50 text-emerald-600'
                                                    : 'bg-rose-50 text-rose-600'
                                                    }`}
                                            >
                                                {order.invoice && order.invoice.status === 'issued' ? 'Facturada' : 'Requiere factura'}
                                            </span>
                                        </div>
                                        <div className="border-t border-slate-100 pt-3 dark:border-gray-700">
                                            <p className="text-[11px] font-medium uppercase tracking-wide text-slate-400">N&uacute;mero interno</p>
                                            <p className="mt-0.5 break-all font-mono text-[13px] text-slate-800 dark:text-white">{order.internal_number || '-'}</p>
                                        </div>
                                        <div className="flex items-center justify-between border-t border-slate-100 pt-3 dark:border-gray-700">
                                            <div>
                                                <p className="text-[11px] font-medium uppercase tracking-wide text-slate-400">Plataforma</p>
                                                <p className="mt-0.5 text-[13px] font-medium text-slate-800 dark:text-white">{salesChannelLabel(order.platform)}</p>
                                            </div>
                                            {order.integration_logo_url && (
                                                <img src={order.integration_logo_url} alt={order.platform} className="h-7 w-7 object-contain" />
                                            )}
                                        </div>
                                        <div className="border-t border-slate-100 pt-3 dark:border-gray-700">
                                            <p className="text-[11px] font-medium uppercase tracking-wide text-slate-400">Fecha</p>
                                            <p className="mt-0.5 text-[13px] tabular-nums text-slate-800 dark:text-white">
                                                {formatDateSeparated(order.occurred_at || order.created_at).date} {formatDateSeparated(order.occurred_at || order.created_at).time}
                                            </p>
                                        </div>
                                    </div>
                                </section>
                            </aside>
                        </main>
                    </div>
                </div>
            )}

            {showChangeStatus && (
                <ChangeStatusModal
                    isOpen={showChangeStatus}
                    onClose={() => setShowChangeStatus(false)}
                    order={order}
                    onSuccess={() => {
                        showToast(`Estado de #${order.order_number} actualizado`, 'success');
                        setShowChangeStatus(false);
                        fetchDetails();
                    }}
                />
            )}
        </div>
    );
}
