'use client';

import { Order, OrderHistory } from '../../domain/types';
import MapComponent from '@/shared/ui/MapComponent';
import { getAIRecommendationAction, getOrderByIdAction, getOrderHistoryAction, updateOrderAction, requestWhatsAppConfirmationAction, checkWhatsAppIntegrationAction } from '../../infra/actions';
import { useState, useEffect } from 'react';
import ShipmentGuideModal from '@/shared/ui/modals/shipment-guide-modal';
import { ChangeStatusModal } from './ChangeStatusModal';
import { isTerminalStatus } from '../../domain/order-status-transitions';
import { useToast } from '@/shared/providers/toast-provider';

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

    // Fetch full order details and history on mount
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
    }, [initialOrder.id]);

    // Derived order object (prefer full, fallback to initial)
    const order = fullOrder || initialOrder;

    // AI Logic - Triggers when fullOrder (with address) is available
    // AI Logic - Triggers when fullOrder (with address) is available AND mode is 'recommendation'
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
            // Reset AI state if leaving recommendation mode (opt)
            setAIRecommendation(null);
            setLoadingAI(false);
        }
    }, [fullOrder, mode]);

    // Management State
    const [isConfirmed, setIsConfirmed] = useState<boolean | null>(false);
    const [novelty, setNovelty] = useState('');
    const [isSaving, setIsSaving] = useState(false);

    // WhatsApp Confirmation State
    const [hasWhatsApp, setHasWhatsApp] = useState(false);
    const [isSendingWhatsApp, setIsSendingWhatsApp] = useState(false);
    const [whatsAppSent, setWhatsAppSent] = useState(false);

    // Initialize management state
    useEffect(() => {
        if (order) {
            setIsConfirmed(order.is_confirmed ?? null);
            setNovelty(order.novelty || '');
        }
    }, [order]);

    // Check if business has WhatsApp integration
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
                // Update local state with the FULL updated order returned by backend
                // This ensures calculated fields like delivery_probability and negative_factors are updated
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

    // Helper para calcular color del texto basado en luminosidad
    const getTextColor = (bgColor: string): string => {
        const hex = bgColor.replace('#', '');
        const r = parseInt(hex.substr(0, 2), 16);
        const g = parseInt(hex.substr(2, 2), 16);
        const b = parseInt(hex.substr(4, 2), 16);
        const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
        return luminance > 0.5 ? '#000000' : '#FFFFFF';
    };

    // Parse items if they are JSON string or access directly
    const items = Array.isArray(order.items) ? order.items : [];

    // Address for Map
    const fullAddress = `${order.shipping_street || ''}`;
    const city = order.shipping_city || '';

    // If loading details, show a skeleton or loading state for critical sections
    const isReady = !loadingDetails && fullOrder;

    return (
        <div className="flex flex-col h-full bg-white">
            {/* AI Recommendation Section - Only shown in recommendation mode */}
            {mode === 'recommendation' && (
                <div className="bg-gradient-to-r from-blue-50 to-indigo-50 p-3 rounded-xl border border-blue-100 shadow-sm relative overflow-hidden transition-all hover:shadow-md">
                    <div className="absolute top-0 right-0 p-2 opacity-5 pointer-events-none">
                        <svg className="w-24 h-24" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm1 15h-2v-2h2zm0-4h-2V7h2z" /></svg>
                    </div>

                    {onClose && (
                        <button
                            onClick={onClose}
                            className="absolute top-2 right-2 z-20 p-1 bg-white dark:bg-gray-800/20 hover:bg-white dark:bg-gray-800/40 text-blue-900 rounded-full transition-colors backdrop-blur-sm"
                            title="Cerrar"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    )}

                    <div className="relative z-10">
                        <h3 className="text-lg font-bold text-blue-900 flex items-center gap-2 mb-2">
                            <span className="text-2xl">🤖</span> Recomendación Inteligente
                        </h3>

                        {isReady ? (
                            <>
                                {loadingAI ? (
                                    <div className="flex items-center gap-3 text-purple-600 bg-white dark:bg-gray-800/50 p-3 rounded-lg animate-pulse">
                                        <div className="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
                                        <span>Analizando mejores rutas y tarifas...</span>
                                    </div>
                                ) : aiRecommendation ? (
                                    <div className="flex flex-col gap-6">
                                        <div className="flex-1 space-y-4">
                                            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                                                <div>
                                                    <span className="text-xs font-bold text-purple-600 dark:text-purple-400 uppercase tracking-wider bg-purple-100 dark:bg-purple-900 px-2 py-1 rounded">
                                                        Mejor Opción
                                                    </span>
                                                    <p className="text-4xl font-extrabold text-purple-600 dark:text-purple-400 mt-2">
                                                        {aiRecommendation.recommended_carrier}
                                                    </p>
                                                </div>
                                                <button
                                                    onClick={() => setShowGuideModal(true)}
                                                    className="bg-purple-600 hover:bg-purple-700 text-white font-bold py-2 px-6 rounded-lg shadow-lg shadow-purple-200 transition-all flex items-center gap-2"
                                                >
                                                    <span>📦</span> Cotizar y Generar Guía
                                                </button>
                                            </div>

                                            <div className="bg-white dark:bg-gray-800/80 p-5 rounded-lg border border-blue-100 text-gray-700 dark:text-gray-200 text-sm leading-relaxed shadow-sm">
                                                <p className="font-semibold text-blue-900 mb-1">Análisis:</p>
                                                {aiRecommendation.reasoning}
                                            </div>
                                        </div>

                                        {aiRecommendation.quotations && aiRecommendation.quotations.length > 0 && (
                                            <div className="border-t border-blue-200 dark:border-blue-900 pt-6 mt-2">
                                                <h4 className="text-sm font-bold text-purple-600 dark:text-purple-400 uppercase tracking-wide mb-4 flex items-center gap-2">
                                                    <span>📊</span> Cotizaciones Estimadas
                                                </h4>
                                                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                                                    {aiRecommendation.quotations.map((quote, idx) => {
                                                        const isRecommended = quote.carrier === aiRecommendation.recommended_carrier;
                                                        return (
                                                            <div key={idx} className={`p-4 rounded-xl border flex flex-col justify-between transition-all hover:shadow-md ${isRecommended
                                                                ? 'bg-white dark:bg-gray-800 border-blue-300 dark:border-blue-600 shadow-sm ring-1 ring-blue-100 dark:ring-blue-900 relative overflow-hidden'
                                                                : 'bg-slate-50 dark:bg-gray-800 border-slate-200 dark:border-gray-700 hover:bg-white dark:hover:bg-gray-700'
                                                                }`}>
                                                                {isRecommended && <div className="absolute top-0 right-0 bg-purple-600 text-white text-[10px] px-2 py-0.5 rounded-bl-lg font-bold">RECOMENDADO</div>}
                                                                <div>
                                                                    <p className={`font-bold text-lg ${isRecommended ? 'text-blue-600 dark:text-blue-400' : 'text-gray-700 dark:text-gray-200'}`}>
                                                                        {quote.carrier}
                                                                    </p>
                                                                    <p className="text-xs text-gray-500 dark:text-gray-400 flex items-center gap-1 mt-1">
                                                                        <span>⏱️</span> {quote.estimated_delivery_days} días hábiles
                                                                    </p>
                                                                </div>
                                                                <div className="mt-4 pt-3 border-t border-gray-100 dark:border-gray-700">
                                                                    <p className="text-gray-500 dark:text-gray-400 text-xs uppercase mb-0.5">Costo Estimado</p>
                                                                    <div className="flex justify-between items-end">
                                                                        <p className={`font-bold text-xl ${isRecommended ? 'text-purple-600' : 'text-gray-600 dark:text-gray-300'}`}>
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
                            <div className="text-sm text-blue-400 mt-1 animate-pulse">Cargando datos de orden...</div>
                        )}

                        {showGuideModal && order && (
                            <ShipmentGuideModal
                                isOpen={showGuideModal}
                                onClose={() => setShowGuideModal(false)}
                                order={order}
                                recommendedCarrier={aiRecommendation?.recommended_carrier}
                                onGuideGenerated={() => {
                                    // Reload order details after guide is generated to get updated shipment data
                                    setShowGuideModal(false);
                                    fetchDetails();
                                }}
                            />
                        )}
                    </div>
                </div>
            )}

            {/* Header */}
            {mode === 'details' && (
                <>
                    <div className="bg-purple-700 dark:bg-purple-900 text-white px-6 py-4 flex items-center justify-between border-b border-purple-800">
                        <h2 className="text-xl font-black tracking-wide">Detalles de la Orden</h2>
                        <p className="text-lg opacity-70 font-semibold">{order.order_number || '-'}</p>
                        <div className="flex items-center gap-3">
                            {order.order_status?.color ? (
                                <span
                                    className="px-3 py-1 text-xs font-bold rounded-full"
                                    style={{
                                        backgroundColor: order.order_status.color,
                                        color: getTextColor(order.order_status.color)
                                    }}
                                >
                                    {order.order_status.name || order.status}
                                </span>
                            ) : (
                                <span className="px-3 py-1 text-xs font-bold rounded-full bg-purple-600 text-white">
                                    {order.order_status?.name || order.status || '-'}
                                </span>
                            )}
                            {onClose && (
                                <button
                                    onClick={onClose}
                                    className="p-1 hover:bg-white/20 rounded-full transition-colors"
                                    title="Cerrar"
                                >
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                    </svg>
                                </button>
                            )}
                        </div>
                    </div>

                    {/* Main Content - 3 Column Layout */}
                    <div className="flex flex-1 overflow-hidden">
                        {/* Sidebar Izquierdo - Información General */}
                        <div className="w-48 bg-purple-700 dark:bg-purple-900 text-white p-3 overflow-y-auto border-r border-purple-800">
                            <h3 className="text-xs font-black uppercase tracking-wider text-gray-300 mb-4">Información General</h3>
                            {loadingDetails ? (
                                <div className="py-4 text-center text-xs text-purple-200">Cargando...</div>
                            ) : (
                                <div className="space-y-4">
                                    <div>
                                        <p className="text-xs font-bold uppercase tracking-widest text-white mb-1">Nº Orden</p>
                                        <p className="text-lg font-black text-white">{order.order_number || '-'}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs font-bold uppercase tracking-widest text-white mb-1">Número Interno</p>
                                        <p className="text-xs font-semibold text-white break-all">{order.internal_number || '-'}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs font-bold uppercase tracking-widest text-white mb-1">Plataforma</p>
                                        {order.integration_logo_url ? (
                                            <img
                                                src={order.integration_logo_url}
                                                alt={order.platform}
                                                className="h-6 w-6 object-contain"
                                                title={order.platform}
                                            />
                                        ) : (
                                            <p className="text-xs font-semibold text-white capitalize">{order.platform || '-'}</p>
                                        )}
                                    </div>
                                    <div>
                                        <p className="text-xs font-bold uppercase tracking-widest text-white mb-1">Estado</p>
                                        <div className="flex flex-col items-center gap-2">
                                            {order.order_status?.color ? (
                                                <span
                                                    className="inline-block px-2 py-0.5 text-[10px] font-bold rounded-full w-fit"
                                                    style={{
                                                        backgroundColor: order.order_status.color,
                                                        color: getTextColor(order.order_status.color)
                                                    }}
                                                >
                                                    {order.order_status.name || order.status}
                                                </span>
                                            ) : (
                                                <span className="inline-block px-2 py-0.5 text-[10px] font-bold rounded-full bg-purple-600 text-white w-fit">
                                                    {order.status || '-'}
                                                </span>
                                            )}
                                            {!isTerminalStatus(order.order_status?.code || order.status || '') && (
                                                <button
                                                    onClick={() => setShowChangeStatus(true)}
                                                    className="px-2 py-0.5 text-[10px] font-bold text-green-300 bg-green-500/20 border border-green-500/40 rounded hover:bg-green-500/30 transition-all w-fit"
                                                >
                                                    Cambiar
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                    <div>
                                        <p className="text-xs font-bold uppercase tracking-widest text-white mb-1">Fecha</p>
                                        <div className="flex flex-row gap-2 items-center">
                                            <p className="text-xs font-bold text-white">{formatDateSeparated(order.occurred_at || order.created_at).date}</p>
                                            <p className="text-xs font-semibold text-white">{formatDateSeparated(order.occurred_at || order.created_at).time}</p>
                                        </div>
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* Centro - Contenido Principal */}
                        <div className="flex-1 overflow-y-auto">
                            <div className="p-3 space-y-3">
                                {/* Resumen Financiero - 3 tarjetas */}
                                <div className="flex gap-3">
                                    <div className="flex-1 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
                                        <p className="text-[9px] font-bold text-gray-400 uppercase tracking-wider mb-1">Subtotal</p>
                                        <p className="text-base font-black text-gray-900 dark:text-white">
                                            {formatCurrency(order.subtotal, order.currency, order.subtotal_presentment, order.currency_presentment)}
                                        </p>
                                    </div>
                                    <div className="flex-1 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
                                        <p className="text-[9px] font-bold text-gray-400 uppercase tracking-wider mb-1">Envío</p>
                                        <p className="text-base font-black text-gray-900 dark:text-white">
                                            {formatCurrency(
                                                order.shipment?.total_cost ?? order.shipping_cost,
                                                order.currency,
                                                order.shipment?.total_cost ?? order.shipping_cost_presentment,
                                                order.currency_presentment
                                            )}
                                        </p>
                                    </div>
                                    <div className="flex-1 bg-purple-100 dark:bg-purple-900/30 border border-purple-200 dark:border-purple-700 rounded-lg p-3">
                                        <p className="text-[9px] font-bold text-purple-600 dark:text-purple-300 uppercase tracking-wider mb-1">Total</p>
                                        <p className="text-xl font-black text-purple-700 dark:text-purple-400">
                                            {formatCurrency(
                                                (order.subtotal || 0) + (order.shipment?.total_cost ?? order.shipping_cost ?? 0),
                                                order.currency,
                                                (order.subtotal_presentment || 0) + (order.shipment?.total_cost ?? order.shipping_cost_presentment ?? 0),
                                                order.currency_presentment
                                            )}
                                        </p>
                                    </div>
                                </div>

                                {/* Productos del Pedido */}
                                <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
                                    <div className="px-4 py-3 bg-purple-100 dark:bg-purple-900/30 border-b border-gray-200 dark:border-gray-700">
                                        <h3 className="text-xs font-black text-gray-500 dark:text-gray-400 uppercase">Productos del Pedido</h3>
                                    </div>
                                    {loadingDetails ? (
                                        <div className="py-4 text-center text-xs text-gray-500 dark:text-gray-400">Cargando productos...</div>
                                    ) : (order.order_items || items).length > 0 ? (
                                        <div className="overflow-x-auto">
                                            <table className="w-full divide-y divide-gray-200 dark:divide-gray-700 text-xs">
                                                <thead className="bg-purple-100 dark:bg-purple-900/20">
                                                    <tr>
                                                        <th className="px-3 py-2 text-left font-bold text-purple-700 dark:text-purple-300 uppercase">Producto</th>
                                                        <th className="px-3 py-2 text-left font-bold text-purple-700 dark:text-purple-300 uppercase">SKU</th>
                                                        <th className="px-3 py-2 text-right font-bold text-purple-700 dark:text-purple-300 uppercase">Cant</th>
                                                        <th className="px-3 py-2 text-right font-bold text-purple-700 dark:text-purple-300 uppercase">Precio</th>
                                                        <th className="px-3 py-2 text-right font-bold text-purple-700 dark:text-purple-300 uppercase">Desc.</th>
                                                        <th className="px-3 py-2 text-right font-bold text-purple-700 dark:text-purple-300 uppercase">Total</th>
                                                    </tr>
                                                </thead>
                                                <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                                                    {(order.order_items || items).map((item: any, idx: number) => (
                                                        <tr key={idx} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                                                            <td className="px-3 py-2 text-gray-900 dark:text-white">{item.product_name || item.name || item.title || '-'}</td>
                                                            <td className="px-3 py-2 text-gray-600 dark:text-gray-300">{item.product_sku || item.sku || '-'}</td>
                                                            <td className="px-3 py-2 text-right text-gray-900 dark:text-white">{item.quantity || 0}</td>
                                                            <td className="px-3 py-2 text-right text-gray-900 dark:text-white">{formatCurrency(item.unit_price || item.price, order.currency, item.unit_price_presentment, order.currency_presentment)}</td>
                                                            <td className="px-3 py-2 text-right">
                                                                {(item.discount > 0 || (item.discount_presentment && item.discount_presentment > 0)) ? (
                                                                    <span className="text-green-600 font-bold">-{formatCurrency(item.discount, order.currency, item.discount_presentment, order.currency_presentment)}</span>
                                                                ) : (
                                                                    <span className="text-gray-400">-</span>
                                                                )}
                                                            </td>
                                                            <td className="px-3 py-2 text-right text-purple-700 dark:text-purple-300 font-bold">{formatCurrency(item.total_price || (parseFloat(item.unit_price || item.price || 0) * (item.quantity || 0)), order.currency, item.total_price_presentment, order.currency_presentment)}</td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        </div>
                                    ) : (
                                        <p className="text-xs text-gray-500 dark:text-gray-400 text-center py-3">No hay información de productos.</p>
                                    )}
                                </div>

                                {/* Cliente y Dirección */}
                                <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
                                    <h3 className="text-xs font-black text-gray-500 dark:text-gray-400 uppercase mb-3 pb-2 border-b border-gray-200 dark:border-gray-700">Cliente y Dirección</h3>
                                    {loadingDetails ? (
                                        <div className="py-2 text-center text-xs text-gray-500 dark:text-gray-400">Cargando...</div>
                                    ) : (
                                        <div className="grid grid-cols-4 gap-3">
                                            <div>
                                                <p className="text-[9px] font-bold text-gray-400 uppercase tracking-wider mb-1">Nombre</p>
                                                <p className="text-xs font-semibold text-gray-900 dark:text-white">{order.customer_name || '-'}</p>
                                            </div>
                                            <div className="col-span-1">
                                                <p className="text-[9px] font-bold text-gray-400 uppercase tracking-wider mb-1">Email</p>
                                                <p className="text-xs font-semibold text-gray-900 dark:text-white break-all">{order.customer_email || '-'}</p>
                                            </div>
                                            <div>
                                                <p className="text-[9px] font-bold text-gray-400 uppercase tracking-wider mb-1">Teléfono</p>
                                                <p className="text-xs font-semibold text-gray-900 dark:text-white">{order.customer_phone || '-'}</p>
                                                {order.customer_dni && (
                                                    <>
                                                        <p className="text-[9px] font-bold text-gray-400 uppercase tracking-wider mt-2 mb-1">DNI</p>
                                                        <p className="text-xs font-semibold text-gray-900 dark:text-white">{order.customer_dni}</p>
                                                    </>
                                                )}
                                            </div>
                                            <div>
                                                <p className="text-[9px] font-bold text-gray-400 uppercase tracking-wider mb-1">Dirección</p>
                                                <div className="space-y-1">
                                                    <p className="text-xs font-semibold text-gray-900 dark:text-white">{order.shipping_street || '-'}</p>
                                                    <p className="text-xs text-gray-700 dark:text-gray-200">
                                                        {order.shipping_city || ''}{order.shipping_state && ', ' + order.shipping_state}{order.shipping_postal_code && ' ' + order.shipping_postal_code}
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                    )}
                                </div>

                                {/* Gestión y Novedades */}
                                <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
                                    <h3 className="text-xs font-black text-gray-500 dark:text-gray-400 uppercase mb-3 pb-2 border-b border-gray-200 dark:border-gray-700">Gestión y Novedades</h3>
                                    <div className="grid grid-cols-2 gap-3 mb-3">
                                        <div className="flex flex-col">
                                            <label className="text-xs font-bold text-gray-700 dark:text-gray-200 mb-1 uppercase">Confirmación de Pedido</label>
                                            <select
                                                className={`block w-full px-2 py-2 text-xs border-1.5 rounded focus:outline-none focus:ring-purple-500 ${isConfirmed === true
                                                    ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400 border-green-500 dark:border-green-600 font-bold'
                                                    : isConfirmed === false
                                                        ? 'bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 border-red-500 dark:border-red-600 font-bold'
                                                        : 'bg-yellow-50 dark:bg-yellow-900/20 text-yellow-700 dark:text-yellow-400 border-yellow-500 dark:border-yellow-600 font-bold'
                                                    }`}
                                                value={isConfirmed === null ? 'pending' : (isConfirmed ? 'yes' : 'no')}
                                                onChange={(e) => {
                                                    const val = e.target.value;
                                                    setIsConfirmed(val === 'pending' ? null : val === 'yes');
                                                }}
                                            >
                                                <option value="yes">Sí, Confirmado</option>
                                                <option value="no">No, Rechazado/Cancelado</option>
                                                <option value="pending">Pendiente confirmación</option>
                                            </select>
                                        </div>
                                        <div className="flex flex-col">
                                            <label className="text-xs font-bold text-gray-700 dark:text-gray-200 mb-1 uppercase">Novedades / Notas</label>
                                            <textarea
                                                rows={2}
                                                className="w-full text-xs border-1.5 border-gray-200 dark:border-gray-600 rounded focus:ring-purple-500 focus:border-purple-500 p-2 text-gray-900 dark:text-white resize-vertical"
                                                placeholder="Escribe aquí novedades (ej: cambio de dirección, cliente contactado, etc.)"
                                                value={novelty}
                                                onChange={(e) => setNovelty(e.target.value)}
                                            />
                                        </div>
                                    </div>
                                    <div className="flex gap-2 mt-3">
                                        <button
                                            onClick={handleSaveManagement}
                                            disabled={isSaving}
                                            className="flex-1 px-4 py-2 border border-transparent text-xs font-bold rounded text-white bg-purple-700 hover:bg-purple-800 disabled:opacity-50 transition-colors"
                                        >
                                            {isSaving ? 'Guardando...' : 'Guardar'}
                                        </button>
                                        {hasWhatsApp && order.is_confirmed !== true && (
                                            <button
                                                onClick={handleWhatsAppConfirmation}
                                                disabled={isSendingWhatsApp || whatsAppSent}
                                                className={`px-4 py-2 border text-xs font-bold rounded transition-colors ${
                                                    whatsAppSent
                                                        ? 'bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-200 border-green-300 dark:border-green-600'
                                                        : 'text-white bg-green-600 hover:bg-green-700 border-transparent'
                                                }`}
                                                title={!order.customer_phone ? 'La orden no tiene teléfono de cliente' : 'Enviar confirmación por WhatsApp'}
                                            >
                                                {isSendingWhatsApp ? 'Enviando...' : whatsAppSent ? '✓ Enviado' : '💬 WhatsApp'}
                                            </button>
                                        )}
                                    </div>
                                </div>

                                {/* Historial de Estados */}
                                <div className="bg-gray-50 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
                                    <h3 className="text-xs font-black text-gray-500 dark:text-gray-400 uppercase mb-3">Historial de Estados</h3>
                                    {loadingHistory ? (
                                        <div className="py-4 text-center text-xs text-gray-500 dark:text-gray-400">Cargando historial...</div>
                                    ) : statusHistory.length === 0 ? (
                                        <p className="text-xs text-gray-500 dark:text-gray-400 text-center py-2">No hay cambios de estado registrados.</p>
                                    ) : (
                                        <div className="relative">
                                            <div className="absolute top-5 left-0 right-0 h-0.5 bg-gray-200 dark:bg-gray-600"></div>
                                            <div className="flex overflow-x-auto pb-2 gap-0" style={{ scrollbarWidth: 'thin' }}>
                                                <div className="flex flex-col items-center flex-shrink-0 relative" style={{ minWidth: '100px' }}>
                                                    <div className="w-8 h-8 rounded-full bg-gray-300 dark:bg-gray-600 flex items-center justify-center z-10 border-2 border-white dark:border-gray-800 shadow-sm">
                                                        <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                                                        </svg>
                                                    </div>
                                                    <span className="text-[9px] font-bold rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 px-1.5 py-0.5 mt-1">
                                                        {statusHistory[0].previous_status || 'Creada'}
                                                    </span>
                                                </div>
                                                {statusHistory.map((entry, idx) => {
                                                    const isLast = idx === statusHistory.length - 1;
                                                    return (
                                                        <div key={entry.id} className="flex flex-col items-center flex-shrink-0 relative" style={{ minWidth: '110px' }}>
                                                            <div className="absolute top-3 -left-3 text-gray-300 dark:text-gray-600">
                                                                <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                                                    <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
                                                                </svg>
                                                            </div>
                                                            <div className={`w-8 h-8 rounded-full flex items-center justify-center z-10 border-2 border-white dark:border-gray-800 shadow-sm ${
                                                                isLast ? 'bg-purple-500' : 'bg-blue-400'
                                                            }`}>
                                                                {isLast ? (
                                                                    <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                                                    </svg>
                                                                ) : (
                                                                    <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                                                                    </svg>
                                                                )}
                                                            </div>
                                                            <span className={`text-[9px] font-bold rounded-full px-1.5 py-0.5 mt-1 ${
                                                                isLast
                                                                    ? 'bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200'
                                                                    : 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200'
                                                            }`}>
                                                                {entry.new_status}
                                                            </span>
                                                            <p className="text-[9px] text-gray-500 dark:text-gray-400 mt-1 text-center leading-tight">
                                                                {new Date(entry.created_at).toLocaleDateString('es-CO', { day: '2-digit', month: 'short' })}
                                                                {' '}
                                                                {new Date(entry.created_at).toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' })}
                                                            </p>
                                                        </div>
                                                    );
                                                })}
                                            </div>
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>

                        {/* Sidebar Derecho - Cronología y Pago */}
                        <div className="w-64 bg-gray-50 dark:bg-gray-800 border-l border-gray-200 dark:border-gray-700 p-4 overflow-y-auto">
                            <h3 className="text-sm font-black uppercase tracking-wider text-gray-500 dark:text-gray-400 mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">Cronología y Pago</h3>

                            {/* Timeline */}
                            <div className="space-y-3 mb-4 pb-4 border-b border-gray-200 dark:border-gray-700">
                                <div className="flex gap-2">
                                    <div className="w-2 h-2 rounded-full bg-purple-600 flex-shrink-0 mt-1"></div>
                                    <div>
                                        <p className="text-xs font-bold text-purple-600 dark:text-purple-400 uppercase">Creado (DB)</p>
                                        <p className="text-xs text-gray-500 dark:text-gray-400">{formatDate(order.created_at)}</p>
                                    </div>
                                </div>
                                <div className="flex gap-2">
                                    <div className="w-2 h-2 rounded-full bg-purple-600 flex-shrink-0 mt-1"></div>
                                    <div>
                                        <p className="text-xs font-bold text-purple-600 dark:text-purple-400 uppercase">Importado</p>
                                        <p className="text-xs text-gray-500 dark:text-gray-400">{formatDate(order.imported_at)}</p>
                                    </div>
                                </div>
                                {order.updated_at && (
                                    <div className="flex gap-2">
                                        <div className="w-2 h-2 rounded-full bg-purple-600 flex-shrink-0 mt-1"></div>
                                        <div>
                                            <p className="text-xs font-bold text-purple-600 dark:text-purple-400 uppercase">Actualizado</p>
                                            <p className="text-xs text-gray-500 dark:text-gray-400">{formatDate(order.updated_at)}</p>
                                        </div>
                                    </div>
                                )}
                            </div>

                            {/* Estado de Pago */}
                            <div>
                                <p className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">Estado de Pago</p>
                                <span className={`inline-block px-2 py-1 text-sm font-bold rounded-full text-white uppercase ${(order.payment_details?.financial_status === 'paid' || order.is_paid) ? 'bg-green-600' :
                                    (order.payment_details?.financial_status === 'refunded') ? 'bg-red-600' :
                                        'bg-orange-500'
                                    }`}>
                                    {order.payment_details?.financial_status || (order.is_paid ? 'PAID' : 'PENDING')}
                                </span>
                                {order.paid_at && (
                                    <div className="mt-2">
                                        <p className="text-xs text-gray-600 dark:text-gray-300 font-semibold">Fecha de Pago</p>
                                        <p className="text-xs text-gray-500 dark:text-gray-400">{formatDate(order.paid_at)}</p>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                </>
            )}


            {/* Change Status Modal */}
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
