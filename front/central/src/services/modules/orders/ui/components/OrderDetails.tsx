'use client';

import { Order } from '../../domain/types';
import { AccordionItem } from '@/shared/ui/accordion';
import MapComponent from '@/shared/ui/MapComponent';
import { getAIRecommendationAction, getOrderByIdAction } from '../../infra/actions';
import { useState, useEffect } from 'react';

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
    initialOrder: Order; // Renamed from order to match usage in page.tsx
    onClose?: () => void; // Added onClose prop
}

export default function OrderDetails({ initialOrder, onClose }: OrderDetailsProps) {
    const [fullOrder, setFullOrder] = useState<Order | null>(null);
    const [aiRecommendation, setAIRecommendation] = useState<AIRecommendation | null>(null);
    const [loadingAI, setLoadingAI] = useState(false);
    const [loadingDetails, setLoadingDetails] = useState(false);

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
                        // Only log specifics if it's an actual failure flag
                        console.error("Failed to load order details:", response.message);
                    }
                    // If success=true but data is missing (rare after backend fix), we silently fallback to initialOrder
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

    // DEBUG: Log the order data to the console as requested
    useEffect(() => {
        console.log("📦 [DEBUG] Order Data (Initial):", initialOrder);
        if (fullOrder) {
            console.log("📦 [DEBUG] Order Data (Full Fetched):", fullOrder);
            console.log("📍 [DEBUG] Address Info:", {
                street: fullOrder.shipping_street,
                city: fullOrder.shipping_city,
                state: fullOrder.shipping_state,
                items: fullOrder.items,
                itemsType: typeof fullOrder.items,
                isArray: Array.isArray(fullOrder.items)
            });
        }
    }, [initialOrder, fullOrder]);

    // AI Logic - Triggers when fullOrder (with address) is available
    useEffect(() => {
        if (fullOrder && fullOrder.shipping_city && fullOrder.shipping_state) {
            setLoadingAI(true);
            getAIRecommendationAction(fullOrder.shipping_city, fullOrder.shipping_state)
                .then(data => setAIRecommendation(data))
                .catch(err => console.error("Error fetching AI:", err))
                .finally(() => setLoadingAI(false));
        }
    }, [fullOrder]); // Depend on fullOrder to ensure we have address

    const formatCurrency = (amount: number | string, currency: string = 'USD') => {
        const num = typeof amount === 'string' ? parseFloat(amount) : amount;
        if (isNaN(num) || num === undefined) return '-';
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency || 'COP',
        }).format(num);
    };

    const formatDate = (dateString: string) => {
        if (!dateString) return '-';
        return new Date(dateString).toLocaleString('es-CO');
    };

    // Parse items if they are JSON string or access directly
    // Ensure we handle both scenarios (array of objects or potentially JSON string from some backends)
    const items = Array.isArray(order.items) ? order.items : [];

    // Address for Map
    const fullAddress = `${order.shipping_street || ''}`;
    const city = order.shipping_city || '';

    // If loading details, show a skeleton or loading state for critical sections
    const isReady = !loadingDetails && fullOrder;

    return (
        <div className="space-y-4 max-h-[80vh] overflow-y-auto p-1">

            {/* AI Recommendation Section */}
            <div className="bg-gradient-to-r from-blue-50 to-indigo-50 p-5 rounded-xl border border-blue-100 shadow-sm relative overflow-hidden transition-all hover:shadow-md">
                <div className="absolute top-0 right-0 p-2 opacity-5 pointer-events-none">
                    <svg className="w-32 h-32" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm1 15h-2v-2h2zm0-4h-2V7h2z" /></svg>
                </div>
                <div className="relative z-10">
                    <h3 className="text-xl font-bold text-blue-900 flex items-center gap-2 mb-4">
                        <span className="text-2xl">🤖</span> Recomendación Inteligente
                    </h3>

                    {isReady ? (
                        <>
                            {loadingAI ? (
                                <div className="flex items-center gap-3 text-blue-600 bg-white/50 p-3 rounded-lg animate-pulse">
                                    <div className="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
                                    <span>Analizando mejores rutas y tarifas...</span>
                                </div>
                            ) : aiRecommendation ? (
                                <div className="flex flex-col gap-6">
                                    {/* Main Recommendation */}
                                    <div className="flex-1 space-y-4">
                                        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                                            <div>
                                                <span className="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-100 px-2 py-1 rounded">
                                                    Mejor Opción
                                                </span>
                                                <p className="text-4xl font-extrabold text-blue-800 mt-2">
                                                    {aiRecommendation.recommended_carrier}
                                                </p>
                                            </div>
                                            {/* Allow Quotation for recommended to be shown here if needed, but keeping it simple for now */}
                                        </div>

                                        <div className="bg-white/80 p-5 rounded-lg border border-blue-100 text-gray-700 text-sm leading-relaxed shadow-sm">
                                            <p className="font-semibold text-blue-900 mb-1">Análisis:</p>
                                            {aiRecommendation.reasoning}
                                        </div>
                                    </div>

                                    {/* Quotations / Alternatives - Now below or wider grid */}
                                    {aiRecommendation.quotations && aiRecommendation.quotations.length > 0 && (
                                        <div className="border-t border-blue-200 pt-6 mt-2">
                                            <h4 className="text-sm font-bold text-blue-800 uppercase tracking-wide mb-4 flex items-center gap-2">
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
                                                            {isRecommended && <div className="absolute top-0 right-0 bg-blue-600 text-white text-[10px] px-2 py-0.5 rounded-bl-lg font-bold">RECOMENDADO</div>}
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
                                                                <p className={`font-bold text-xl ${isRecommended ? 'text-blue-600' : 'text-gray-600'}`}>
                                                                    {formatCurrency(quote.estimated_cost, 'COP')}
                                                                </p>
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
                </div>
            </div>

            {/* Order Information */}
            <AccordionItem title="Información de la Orden" defaultOpen={true}>
                <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
                    <div>
                        <p className="text-xs text-gray-500 uppercase">Número de Orden</p>
                        <p className="text-sm font-medium text-gray-900">{order.order_number}</p>
                    </div>
                    <div>
                        <p className="text-xs text-gray-500 uppercase">Plataforma</p>
                        <p className="text-sm font-medium text-gray-900 capitalize">{order.platform}</p>
                    </div>
                    <div>
                        <p className="text-xs text-gray-500 uppercase">Estado</p>
                        <span className="inline-block px-2 py-0.5 text-xs font-medium rounded-full bg-blue-100 text-blue-800 mt-1">
                            {order.status}
                        </span>
                    </div>
                    <div>
                        <p className="text-xs text-gray-500 uppercase">Fecha</p>
                        <p className="text-sm text-gray-900">{formatDate(order.occurred_at || order.created_at)}</p>
                    </div>
                </div>
            </AccordionItem>

            {/* Order Items */}
            <AccordionItem title="Productos del Pedido">
                {loadingDetails ? (
                    <div className="py-4 text-center text-sm text-gray-500">Cargando productos...</div>
                ) : items.length > 0 ? (
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase">Producto</th>
                                    <th className="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase">SKU</th>
                                    <th className="px-3 py-2 text-right text-xs font-medium text-gray-500 uppercase">Cant.</th>
                                    <th className="px-3 py-2 text-right text-xs font-medium text-gray-500 uppercase">Precio</th>
                                    <th className="px-3 py-2 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {items.map((item: any, idx: number) => (
                                    <tr key={idx}>
                                        <td className="px-3 py-2 text-sm text-gray-900">{item.name || item.title || item.product_name}</td>
                                        <td className="px-3 py-2 text-sm text-gray-500">{item.sku || item.product_sku || '-'}</td>
                                        <td className="px-3 py-2 text-sm text-gray-900 text-right">{item.quantity}</td>
                                        <td className="px-3 py-2 text-sm text-gray-900 text-right">{formatCurrency(item.price || item.unit_price, order.currency)}</td>
                                        <td className="px-3 py-2 text-sm text-gray-900 text-right">{formatCurrency((parseFloat(item.price || item.unit_price) * item.quantity), order.currency)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                ) : (
                    <p className="text-sm text-gray-500 text-center py-2">No hay información de productos.</p>
                )}
            </AccordionItem>

            {/* Customer Information */}
            <AccordionItem title="Información del Cliente">
                {loadingDetails ? (
                    <div className="py-4 text-center text-sm text-gray-500">Cargando cliente...</div>
                ) : (
                    <div className="grid grid-cols-2 lg:grid-cols-4 gap-x-4 gap-y-3">
                        <div className="col-span-2 sm:col-span-1">
                            <p className="text-xs text-gray-500 uppercase">Nombre</p>
                            <p className="text-sm font-medium text-gray-900">{order.customer_name || '-'}</p>
                        </div>
                        <div className="col-span-2 sm:col-span-1">
                            <p className="text-xs text-gray-500 uppercase">Email</p>
                            <p className="text-sm font-medium text-gray-900 break-all">{order.customer_email || '-'}</p>
                        </div>
                        <div>
                            <p className="text-xs text-gray-500 uppercase">Teléfono</p>
                            <p className="text-sm font-medium text-gray-900">{order.customer_phone || '-'}</p>
                        </div>

                    </div>
                )}
            </AccordionItem>
            {/* Shipping Address */}
            <div className="bg-purple-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">Dirección de Envío</h3>
                <div className="space-y-2">
                    <p className="text-sm text-gray-900">{order.shipping_street}</p>
                    <p className="text-sm text-gray-900">
                        {order.shipping_city}, {order.shipping_state} {order.shipping_postal_code}
                    </p>
                    <p className="text-sm text-gray-900">{order.shipping_country}</p>
                    {order.delivery_probability !== undefined && order.delivery_probability !== null && (
                        <div className="mt-2 pt-2 border-t border-purple-200">
                            <p className="text-sm text-gray-500 mb-1">Probabilidad de Entrega</p>
                            <div className="flex items-center gap-2">
                                <div className="flex-1 bg-white rounded-full h-2.5 border border-purple-200">
                                    <div
                                        className={`h-2.5 rounded-full ${order.delivery_probability < 30 ? 'bg-red-500' :
                                            order.delivery_probability < 70 ? 'bg-yellow-500' : 'bg-green-500'
                                            }`}
                                        style={{ width: `${order.delivery_probability}%` }}
                                    ></div>
                                </div>
                                <span className="text-sm font-medium text-gray-900">{order.delivery_probability}%</span>
                            </div>
                        </div>
                    )}
                </div>
            </div>
            {/* Shipping Address & Map */}
            <AccordionItem title="Dirección de Envío y Mapa">
                {loadingDetails ? (
                    <div className="py-4 text-center text-sm text-gray-500">Cargando dirección...</div>
                ) : (
                    <div className="space-y-4">
                        <div className="bg-gray-50 p-3 rounded border border-gray-100">
                            <p className="text-sm font-medium text-gray-900">{order.shipping_street}</p>
                            <p className="text-sm text-gray-600">
                                {order.shipping_city}, {order.shipping_state} {order.shipping_postal_code}
                            </p>
                            <p className="text-sm text-gray-600 uppercase mt-1">{order.shipping_country}</p>
                        </div>

                        {order.shipping_street || order.shipping_city ? (
                            <div className="w-full rounded-lg border border-gray-200 overflow-hidden">
                                <MapComponent
                                    address={fullAddress}
                                    city={city}
                                    height="300px"
                                />
                            </div>
                        ) : (
                            <div className="text-sm text-gray-500 italic">No hay dirección para mostrar mapa.</div>
                        )}
                    </div>
                )}
            </AccordionItem>

            {/* Financial Summary */}
            <AccordionItem title="Resumen Financiero">
                <div className="space-y-2">
                    <div className="flex justify-between">
                        <span className="text-sm text-gray-600">Subtotal</span>
                        <span className="text-sm font-medium text-gray-900">
                            {formatCurrency(order.subtotal, order.currency)}
                        </span>
                    </div>
                    <div className="flex justify-between">
                        <span className="text-sm text-gray-600">Impuestos</span>
                        <span className="text-sm font-medium text-gray-900">
                            {formatCurrency(order.tax, order.currency)}
                        </span>
                    </div>
                    {/* Only show discount if > 0 */}
                    {order.discount > 0 && (
                        <div className="flex justify-between">
                            <span className="text-sm text-gray-600">Descuento</span>
                            <span className="text-sm font-medium text-gray-900 text-green-600">
                                -{formatCurrency(order.discount, order.currency)}
                            </span>
                        </div>
                    )}
                    <div className="flex justify-between">
                        <span className="text-sm text-gray-600">Envío</span>
                        <span className="text-sm font-medium text-gray-900">
                            {formatCurrency(order.shipping_cost, order.currency)}
                        </span>
                    </div>
                    <div className="flex justify-between pt-3 border-t border-gray-100 mt-2">
                        <span className="text-base font-semibold text-gray-900">Total</span>
                        <span className="text-base font-bold text-blue-600">
                            {formatCurrency(order.total_amount, order.currency)}
                        </span>
                    </div>
                </div>
            </AccordionItem>

            {/* Payment & Dates Group */}
            <div className="grid grid-cols-1 gap-4">
                <AccordionItem title="Detalles de Pago">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-xs text-gray-500 uppercase">Estado Financiero</p>
                            <span className={`inline-block px-2 py-1 text-xs font-medium rounded-full mt-1 ${(order.payment_details?.financial_status === 'paid' || order.is_paid) ? 'bg-green-100 text-green-800' :
                                (order.payment_details?.financial_status === 'refunded') ? 'bg-red-100 text-red-800' :
                                    'bg-yellow-100 text-yellow-800'
                                }`}>
                                {order.payment_details?.financial_status?.toUpperCase() || (order.is_paid ? 'PAID' : 'PENDING')}
                            </span>
                        </div>
                        {order.paid_at && (
                            <div className="text-right">
                                <p className="text-xs text-gray-500 uppercase">Fecha de Pago</p>
                                <p className="text-sm font-medium text-gray-900">{formatDate(order.paid_at)}</p>
                            </div>
                        )}
                    </div>
                </AccordionItem>

                <AccordionItem title="Cronología">
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <p className="text-xs text-gray-500 uppercase">Creado (DB)</p>
                            <p className="text-sm text-gray-700">{formatDate(order.created_at)}</p>
                        </div>
                        <div>
                            <p className="text-xs text-gray-500 uppercase">Importado</p>
                            <p className="text-sm text-gray-700">{formatDate(order.imported_at)}</p>
                        </div>
                    </div>
                </AccordionItem>
            </div>
        </div >
    );
}
