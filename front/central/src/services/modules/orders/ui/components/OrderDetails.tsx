'use client';

import { Order } from '../../domain/types';
import MapComponent from '@/shared/ui/MapComponent';
import { getAIRecommendationAction, getOrderByIdAction, updateOrderAction } from '../../infra/actions';
import { useState, useEffect } from 'react';
import ShipmentGuideModal from '@/shared/ui/modals/shipment-guide-modal';

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

    // Fetch full order details on mount
    useEffect(() => {
        let isMounted = true;

        async function fetchDetails() {
            if (!initialOrder.id) return;

            setLoadingDetails(true);
            try {
                const response = await getOrderByIdAction(initialOrder.id);
                if (isMounted) {
                    if (response.success && response.data) {
                        setFullOrder(response.data);
                    } else if (!response.success) {
                        console.error("Failed to load order details:", response.message);
                    }
                }
            } catch (error) {
                console.error("Error loading order details:", error);
            } finally {
                if (isMounted) setLoadingDetails(false);
            }
        }

        fetchDetails();

        return () => { isMounted = false; };
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

    // Initialize management state
    useEffect(() => {
        if (order) {
            setIsConfirmed(order.is_confirmed ?? null);
            setNovelty(order.novelty || '');
        }
    }, [order]);

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
        return new Date(dateString).toLocaleString('es-CO');
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
        <div className="space-y-2 p-3">
            {/* AI Recommendation Section - Only shown in recommendation mode */}
            {mode === 'recommendation' && (
                <div className="bg-gradient-to-r from-blue-50 to-indigo-50 p-3 rounded-xl border border-blue-100 shadow-sm relative overflow-hidden transition-all hover:shadow-md">
                    <div className="absolute top-0 right-0 p-2 opacity-5 pointer-events-none">
                        <svg className="w-24 h-24" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm1 15h-2v-2h2zm0-4h-2V7h2z" /></svg>
                    </div>

                    {onClose && (
                        <button
                            onClick={onClose}
                            className="absolute top-2 right-2 z-20 p-1 bg-white/20 hover:bg-white/40 text-blue-900 rounded-full transition-colors backdrop-blur-sm"
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
                                    <div className="flex items-center gap-3 text-purple-600 bg-white/50 p-3 rounded-lg animate-pulse">
                                        <div className="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
                                        <span>Analizando mejores rutas y tarifas...</span>
                                    </div>
                                ) : aiRecommendation ? (
                                    <div className="flex flex-col gap-6">
                                        <div className="flex-1 space-y-4">
                                            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                                                <div>
                                                    <span className="text-xs font-bold text-purple-600 uppercase tracking-wider bg-purple-100 px-2 py-1 rounded">
                                                        Mejor Opción
                                                    </span>
                                                    <p className="text-4xl font-extrabold text-purple-800 mt-2">
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

                                            <div className="bg-white/80 p-5 rounded-lg border border-blue-100 text-gray-700 text-sm leading-relaxed shadow-sm">
                                                <p className="font-semibold text-blue-900 mb-1">Análisis:</p>
                                                {aiRecommendation.reasoning}
                                            </div>
                                        </div>

                                        {aiRecommendation.quotations && aiRecommendation.quotations.length > 0 && (
                                            <div className="border-t border-blue-200 pt-6 mt-2">
                                                <h4 className="text-sm font-bold text-purple-800 uppercase tracking-wide mb-4 flex items-center gap-2">
                                                    <span>📊</span> Cotizaciones Estimadas
                                                </h4>
                                                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                                                    {aiRecommendation.quotations.map((quote, idx) => {
                                                        const isRecommended = quote.carrier === aiRecommendation.recommended_carrier;
                                                        return (
                                                            <div key={idx} className={`p-4 rounded-xl border flex flex-col justify-between transition-all hover:shadow-md ${isRecommended
                                                                ? 'bg-white border-blue-300 shadow-sm ring-1 ring-blue-100 relative overflow-hidden'
                                                                : 'bg-slate-50 border-slate-200 hover:bg-white'
                                                                }`}>
                                                                {isRecommended && <div className="absolute top-0 right-0 bg-purple-600 text-white text-[10px] px-2 py-0.5 rounded-bl-lg font-bold">RECOMENDADO</div>}
                                                                <div>
                                                                    <p className={`font-bold text-lg ${isRecommended ? 'text-blue-900' : 'text-gray-700'}`}>
                                                                        {quote.carrier}
                                                                    </p>
                                                                    <p className="text-xs text-gray-500 flex items-center gap-1 mt-1">
                                                                        <span>⏱️</span> {quote.estimated_delivery_days} días hábiles
                                                                    </p>
                                                                </div>
                                                                <div className="mt-4 pt-3 border-t border-gray-100">
                                                                    <p className="text-gray-500 text-xs uppercase mb-0.5">Costo Estimado</p>
                                                                    <div className="flex justify-between items-end">
                                                                        <p className={`font-bold text-xl ${isRecommended ? 'text-purple-600' : 'text-gray-600'}`}>
                                                                            {formatCurrency(quote.estimated_cost, 'COP')}
                                                                        </p>
                                                                        <button
                                                                            onClick={() => setShowGuideModal(true)}
                                                                            className="text-xs bg-gray-200 hover:bg-gray-300 text-gray-700 px-2 py-1 rounded"
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
                                    <div className="text-sm text-gray-500 italic bg-gray-50 p-3 rounded border border-gray-100">
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
                            />
                        )}
                    </div>
                </div>
            )}

            {/* Order Details Sections - Only shown in details mode */}
            {mode === 'details' && (
                <>
                    {/* 3 Column Layout - Información Grande a la Izquierda */}
                    <div className="grid gap-3" style={{ gridTemplateColumns: '1fr 1fr 1fr', gridTemplateRows: 'auto auto auto' }}>
                        {/* ROW 1, COL 1: Información General - LARGER */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <h3 className="text-lg font-bold text-gray-900 mb-3">Información General</h3>
                            {loadingDetails ? (
                                <div className="py-4 text-center text-sm text-gray-500">Cargando información...</div>
                            ) : (
                                <div className="space-y-3">
                                    <div>
                                        <p className="text-xs text-gray-500 uppercase tracking-wide font-semibold">Nº Orden</p>
                                        <p className="text-sm font-bold text-gray-900 mt-0.5">{order.order_number || '-'}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-gray-500 uppercase tracking-wide font-semibold">Número Interno</p>
                                        <p className="text-sm font-medium text-gray-900 break-all mt-0.5">{order.internal_number || '-'}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-gray-500 uppercase tracking-wide font-semibold">Plataforma</p>
                                        {order.integration_logo_url ? (
                                            <img
                                                src={order.integration_logo_url}
                                                alt={order.platform}
                                                className="h-10 w-10 object-contain mt-1"
                                                title={order.platform}
                                            />
                                        ) : (
                                            <p className="text-sm font-medium text-gray-900 capitalize mt-1">{order.platform || '-'}</p>
                                        )}
                                    </div>
                                    <div>
                                        <p className="text-xs text-gray-500 uppercase tracking-wide font-semibold">Estado (Probability)</p>
                                        {order.order_status?.color ? (
                                            <span
                                                className="inline-block px-3 py-1 text-sm font-semibold rounded-full mt-1"
                                                style={{
                                                    backgroundColor: order.order_status.color,
                                                    color: getTextColor(order.order_status.color)
                                                }}
                                            >
                                                {order.order_status.name || order.status}
                                            </span>
                                        ) : (
                                            <span className="inline-block px-3 py-1 text-sm font-semibold rounded-full bg-purple-100 text-purple-800 mt-1">
                                                {order.order_status?.name || order.status || '-'}
                                            </span>
                                        )}
                                    </div>
                                    {order.original_status && (
                                        <div>
                                            <p className="text-xs text-gray-500 uppercase tracking-wide font-semibold">Estado Original (Shopify)</p>
                                            <span className="inline-block px-3 py-1 text-sm font-semibold rounded-full bg-gray-100 text-gray-800 mt-1">
                                                {order.original_status}
                                            </span>
                                        </div>
                                    )}
                                    <div>
                                        <p className="text-xs text-gray-500 uppercase tracking-wide font-semibold">Fecha</p>
                                        <p className="text-sm font-bold text-gray-900 mt-0.5">{formatDate(order.occurred_at || order.created_at)}</p>
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* ROW 1, COL 2: Resumen Financiero + Productos (STACKED) */}
                        <div className="space-y-3">
                            {/* Resumen Financiero */}
                            <div className="bg-gray-50 rounded-lg p-4">
                                <h3 className="text-lg font-bold text-gray-900 mb-3">Resumen Financiero</h3>
                                <div className="space-y-2">
                                    <div className="flex justify-between">
                                        <span className="text-sm text-gray-600 font-medium">Subtotal</span>
                                        <span className="text-sm font-semibold text-gray-900">
                                            {formatCurrency(order.subtotal, order.currency, order.subtotal_presentment, order.currency_presentment)}
                                        </span>
                                    </div>
                                    <div className="flex justify-between">
                                        <span className="text-sm text-gray-600 font-medium">Impuestos</span>
                                        <span className="text-sm font-semibold text-gray-900">
                                            {formatCurrency(order.tax, order.currency, order.tax_presentment, order.currency_presentment)}
                                        </span>
                                    </div>
                                    {(order.discount > 0 || (order.discount_presentment && order.discount_presentment > 0)) && (
                                        <div className="flex justify-between">
                                            <span className="text-sm text-gray-600 font-medium">Descuento</span>
                                            <span className="text-sm font-semibold text-gray-900 text-green-600">
                                                -{formatCurrency(order.discount, order.currency, order.discount_presentment, order.currency_presentment)}
                                            </span>
                                        </div>
                                    )}
                                    <div className="flex justify-between">
                                        <span className="text-sm text-gray-600 font-medium">Envío</span>
                                        <span className="text-sm font-semibold text-gray-900">
                                            {formatCurrency(order.shipping_cost, order.currency, order.shipping_cost_presentment, order.currency_presentment)}
                                        </span>
                                    </div>
                                    <div className="flex justify-between pt-2 border-t border-gray-300 mt-2">
                                        <span className="text-base font-bold text-gray-900">Total</span>
                                        <span className="text-base font-bold text-purple-600">
                                            {formatCurrency(order.total_amount, order.currency, order.total_amount_presentment, order.currency_presentment)}
                                        </span>
                                    </div>
                                </div>
                            </div>

                            {/* Productos del Pedido */}
                            <div className="bg-gray-50 rounded-lg p-4">
                                <h3 className="text-lg font-bold text-gray-900 mb-3">Productos del Pedido</h3>
                                {loadingDetails ? (
                                    <div className="py-4 text-center text-sm text-gray-500">Cargando productos...</div>
                                ) : (order.order_items || items).length > 0 ? (
                                    <div className="overflow-x-auto">
                                        <table className="w-full divide-y divide-gray-200 text-sm">
                                            <thead className="bg-gray-100">
                                                <tr>
                                                    <th className="px-2 py-2 text-left text-xs font-semibold text-gray-700 uppercase">Producto</th>
                                                    <th className="px-2 py-2 text-left text-xs font-semibold text-gray-700 uppercase">SKU</th>
                                                    <th className="px-2 py-2 text-right text-xs font-semibold text-gray-700 uppercase">Cant</th>
                                                    <th className="px-2 py-2 text-right text-xs font-semibold text-gray-700 uppercase">Precio</th>
                                                    <th className="px-2 py-2 text-right text-xs font-semibold text-gray-700 uppercase">Total</th>
                                                </tr>
                                            </thead>
                                            <tbody className="bg-white divide-y divide-gray-200">
                                                {(order.order_items || items).map((item: any, idx: number) => (
                                                    <tr key={idx} className="hover:bg-gray-50">
                                                        <td className="px-2 py-2 text-xs text-gray-900">{item.product_name || item.name || item.title || '-'}</td>
                                                        <td className="px-2 py-2 text-xs text-gray-600">{item.product_sku || item.sku || '-'}</td>
                                                        <td className="px-2 py-2 text-xs text-gray-900 text-right">{item.quantity || 0}</td>
                                                        <td className="px-2 py-2 text-xs text-gray-900 text-right">{formatCurrency(item.unit_price || item.price, order.currency, item.unit_price_presentment, order.currency_presentment)}</td>
                                                        <td className="px-2 py-2 text-xs text-gray-900 text-right font-semibold">{formatCurrency(item.total_price || (parseFloat(item.unit_price || item.price || 0) * (item.quantity || 0)), order.currency, item.total_price_presentment, order.currency_presentment)}</td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    </div>
                                ) : (
                                    <p className="text-sm text-gray-500 text-center py-2">No hay información de productos.</p>
                                )}
                            </div>
                        </div>

                        {/* ROW 1, COL 3: Cronología y Pago */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <h3 className="text-lg font-bold text-gray-900 mb-3">Cronología y Pago</h3>
                            <div className="space-y-4">
                                {/* Cronología */}
                                <div>
                                    <p className="text-xs text-gray-500 uppercase font-bold tracking-wide mb-2">Cronología</p>
                                    <div className="space-y-2 border-b border-gray-200 pb-3">
                                        <div>
                                            <p className="text-xs text-gray-600 font-medium">Creado (DB)</p>
                                            <p className="text-sm font-semibold text-gray-900">{formatDate(order.created_at)}</p>
                                        </div>
                                        <div>
                                            <p className="text-xs text-gray-600 font-medium">Importado</p>
                                            <p className="text-sm font-semibold text-gray-900">{formatDate(order.imported_at)}</p>
                                        </div>
                                        {order.updated_at && (
                                            <div>
                                                <p className="text-xs text-gray-600 font-medium">Actualizado</p>
                                                <p className="text-sm font-semibold text-gray-900">{formatDate(order.updated_at)}</p>
                                            </div>
                                        )}
                                    </div>
                                </div>

                                {/* Detalles de Pago */}
                                <div>
                                    <p className="text-xs text-gray-500 uppercase font-bold tracking-wide mb-2">Estado de Pago</p>
                                    <div className="space-y-2">
                                        <div>
                                            <p className="text-xs text-gray-600 font-medium">Estado</p>
                                            <span className={`inline-block px-3 py-1 text-sm font-semibold rounded-full ${(order.payment_details?.financial_status === 'paid' || order.is_paid) ? 'bg-green-100 text-green-800' :
                                                (order.payment_details?.financial_status === 'refunded') ? 'bg-red-100 text-red-800' :
                                                    'bg-yellow-100 text-yellow-800'
                                                }`}>
                                                {order.payment_details?.financial_status?.toUpperCase() || (order.is_paid ? 'PAID' : 'PENDING')}
                                            </span>
                                        </div>
                                        {order.paid_at && (
                                            <div>
                                                <p className="text-xs text-gray-600 font-medium">Fecha</p>
                                                <p className="text-sm font-semibold text-gray-900">{formatDate(order.paid_at)}</p>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* ROW 2, COL 1-3: Cliente y Dirección */}
                        <div className="bg-gray-50 rounded-lg p-4" style={{ gridColumn: '1 / 4' }}>
                            <h3 className="text-lg font-bold text-gray-900 mb-3">Cliente y Dirección</h3>
                            {loadingDetails ? (
                                <div className="py-2 text-center text-xs text-gray-500">Cargando...</div>
                            ) : (
                                <div className="grid grid-cols-4 gap-4">
                                    {/* Col 1: Nombre */}
                                    <div>
                                        <p className="text-gray-500 uppercase font-semibold text-xs mb-1">Nombre</p>
                                        <p className="font-bold text-gray-900 text-sm">{order.customer_name || '-'}</p>
                                    </div>

                                    {/* Col 2: Email */}
                                    <div>
                                        <p className="text-gray-500 uppercase font-semibold text-xs mb-1">Email</p>
                                        <p className="font-medium text-gray-900 break-all text-sm">{order.customer_email || '-'}</p>
                                    </div>

                                    {/* Col 3: Teléfono + DNI */}
                                    <div>
                                        <div>
                                            <p className="text-gray-500 uppercase font-semibold text-xs mb-1">Teléfono</p>
                                            <p className="font-medium text-gray-900 text-sm">{order.customer_phone || '-'}</p>
                                        </div>
                                        {order.customer_dni && (
                                            <div className="mt-2">
                                                <p className="text-gray-500 uppercase font-semibold text-xs mb-1">DNI</p>
                                                <p className="font-medium text-gray-900 text-sm">{order.customer_dni}</p>
                                            </div>
                                        )}
                                    </div>

                                    {/* Col 4: Dirección */}
                                    <div>
                                        <p className="text-gray-500 uppercase font-semibold text-xs mb-1">Dirección</p>
                                        <div className="space-y-1">
                                            <p className="font-medium text-gray-900 text-sm">{order.shipping_street || '-'}</p>
                                            <p className="text-gray-700 text-sm">
                                                {order.shipping_city || ''}{order.shipping_state && ', ' + order.shipping_state}{order.shipping_postal_code && ' ' + order.shipping_postal_code}
                                            </p>
                                            <p className="uppercase text-gray-700 text-sm">{order.shipping_country || '-'}</p>
                                        </div>
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* ROW 3, COL 1-3: Gestión y Novedades */}
                        <div className="bg-gray-50 rounded-lg p-4" style={{ gridColumn: '1 / 4' }}>
                            <h3 className="text-lg font-bold text-gray-900 mb-3">Gestión y Novedades</h3>
                            <div className="space-y-3">
                                <div className="flex flex-col">
                                    <label className="text-sm font-semibold text-gray-700 mb-2">Confirmación de Pedido</label>
                                    <select
                                        className={`block w-full pl-3 pr-3 py-2 text-sm border rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500 ${isConfirmed === true
                                            ? 'bg-green-50 text-green-700 border-green-500 font-semibold'
                                            : isConfirmed === false
                                                ? 'bg-red-50 text-red-700 border-red-500 font-semibold'
                                                : 'bg-yellow-50 text-yellow-700 border-yellow-500 font-semibold'
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
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 mb-2">Novedades / Notas</label>
                                    <textarea
                                        rows={3}
                                        className="shadow-sm focus:ring-purple-500 focus:border-purple-500 block w-full text-sm border-gray-300 rounded-md p-2 border text-gray-900"
                                        placeholder="Escribe aquí novedades (ej: cambio de dirección, cliente contactado, etc.)"
                                        value={novelty}
                                        onChange={(e) => setNovelty(e.target.value)}
                                    />
                                </div>
                                <button
                                    onClick={handleSaveManagement}
                                    disabled={isSaving}
                                    className="inline-flex items-center justify-center px-4 py-2 w-full border border-transparent text-sm font-semibold rounded-md shadow-sm text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 disabled:opacity-50"
                                >
                                    {isSaving ? 'Guardando...' : 'Guardar'}
                                </button>
                            </div>
                        </div>
                    </div>
                </>
            )}
        </div>
    );
}
