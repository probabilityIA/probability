import { FullWidthModal } from '@/shared/ui/full-width-modal';
import { useEffect, useState } from 'react';
import { getOrderRawAction } from '../../infra/actions';

interface RawOrderModalProps {
    orderId: string;
    isOpen: boolean;
    onClose: () => void;
    integrationLogoUrl?: string;
    platform?: string;
}

interface ShopifyOrder {
    id?: number;
    name?: string;
    email?: string;
    contact_email?: string;
    phone?: string;
    created_at?: string;
    updated_at?: string;
    processed_at?: string;
    financial_status?: string;
    fulfillment_status?: string;
    total_price?: string;
    subtotal_price?: string;
    total_tax?: string;
    currency?: string;
    customer?: {
        id?: number;
        email?: string;
        first_name?: string;
        last_name?: string;
        phone?: string;
    };
    line_items?: Array<{
        id?: number;
        title?: string;
        name?: string;
        sku?: string;
        quantity?: number;
        price?: string;
        vendor?: string;
    }>;
    shipping_address?: {
        first_name?: string;
        last_name?: string;
        address1?: string;
        address2?: string;
        city?: string;
        province?: string;
        country?: string;
        zip?: string;
        phone?: string;
    };
    billing_address?: {
        first_name?: string;
        last_name?: string;
        address1?: string;
        address2?: string;
        city?: string;
        province?: string;
        country?: string;
        zip?: string;
        phone?: string;
    };
    fulfillments?: Array<{
        id?: number;
        status?: string;
        tracking_number?: string;
        tracking_company?: string;
        created_at?: string;
    }>;
    [key: string]: any;
}

export default function RawOrderModal({ orderId, isOpen, onClose, integrationLogoUrl, platform }: RawOrderModalProps) {
    const [data, setData] = useState<any>(null);
    const [shopifyOrder, setShopifyOrder] = useState<ShopifyOrder | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showRaw, setShowRaw] = useState(false);

    useEffect(() => {
        if (isOpen && orderId) {
            fetchRawData();
        } else {
            setData(null);
            setShopifyOrder(null);
            setError(null);
            setShowRaw(false);
        }
    }, [isOpen, orderId]);

    const fetchRawData = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getOrderRawAction(orderId);
            if (response.success && response.data) {
                setData(response.data);
                // Intentar parsear raw_data si existe
                if (response.data.raw_data) {
                    try {
                        const rawData = typeof response.data.raw_data === 'string' 
                            ? JSON.parse(response.data.raw_data) 
                            : response.data.raw_data;
                        setShopifyOrder(rawData);
                    } catch (e) {
                        console.error('Error parsing raw_data:', e);
                        setError('Error al parsear los datos crudos de la orden');
                    }
                } else {
                    setError('Esta orden no tiene datos crudos guardados. Los datos crudos solo están disponibles para órdenes creadas después de la implementación de esta funcionalidad.');
                }
            } else {
                setError(response.message || 'Error al cargar los datos crudos');
            }
        } catch (err: any) {
            // Si el error es 404 o "not found", mostrar mensaje más amigable
            const errorMessage = err.message || '';
            if (errorMessage.includes('not found') || 
                errorMessage.includes('no encontrado') ||
                errorMessage.includes('raw data not found') ||
                errorMessage.includes('Datos crudos no encontrados')) {
                setError('Esta orden no tiene datos crudos guardados. Los datos crudos solo están disponibles para órdenes creadas después de la implementación de esta funcionalidad.');
            } else {
                setError(err.message || 'Error al cargar los datos crudos');
            }
        } finally {
            setLoading(false);
        }
    };

    const formatDate = (dateString?: string) => {
        if (!dateString) return 'N/A';
        try {
            return new Date(dateString).toLocaleString('es-CO', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
            });
        } catch {
            return dateString;
        }
    };

    const formatCurrency = (amount?: string, currency?: string) => {
        if (!amount) return 'N/A';
        const num = parseFloat(amount);
        if (isNaN(num)) return amount;
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency || 'USD',
        }).format(num);
    };

    return (
        <FullWidthModal
            isOpen={isOpen}
            onClose={onClose}
            title={
                <div className="flex items-center gap-3">
                    {integrationLogoUrl ? (
                        <div className="h-8 w-8 rounded-full shadow-md border-2 border-gray-200 bg-white flex items-center justify-center overflow-hidden">
                            <img 
                                src={integrationLogoUrl} 
                                alt={platform || 'Integración'}
                                className="h-full w-full object-contain p-1"
                            />
                        </div>
                    ) : platform ? (
                        <div className="h-8 w-8 rounded-full shadow-md border-2 border-gray-200 bg-white flex items-center justify-center">
                            <span className="text-xs font-medium text-gray-600 uppercase">
                                {platform.charAt(0)}
                            </span>
                        </div>
                    ) : null}
                    <span>Orden Original de {platform || 'Shopify'}</span>
                </div>
            }
            width="95vw"
            height="90vh"
        >
            <div>
                {loading && <div className="text-center py-4">Cargando datos...</div>}
                {error && (
                    <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative mb-4">
                        {error}
                    </div>
                )}
                
                {shopifyOrder && !showRaw && (
                    <div className="space-y-4">
                        {/* Botón para ver JSON crudo */}
                        <div className="flex justify-end mb-4">
                            <button
                                onClick={() => setShowRaw(true)}
                                className="px-4 py-2 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md transition-colors font-medium"
                            >
                                Ver JSON Crudo
                            </button>
                        </div>

                        {/* Información General */}
                        <div className="bg-gray-50 rounded-lg p-5">
                            <h3 className="text-lg font-semibold text-gray-900 mb-4">Información General</h3>
                            <div className="grid grid-cols-3 gap-4">
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">ID de Orden</p>
                                    <p className="text-sm font-medium text-gray-900">{shopifyOrder.id || 'N/A'}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">Número de Orden</p>
                                    <p className="text-sm font-medium text-gray-900">{shopifyOrder.name || 'N/A'}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">Estado Financiero</p>
                                    <p className="text-sm font-medium text-gray-900 capitalize">{shopifyOrder.financial_status || 'N/A'}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">Estado de Fulfillment</p>
                                    <p className="text-sm font-medium text-gray-900 capitalize">{shopifyOrder.fulfillment_status || 'N/A'}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">Total</p>
                                    <p className="text-sm font-medium text-gray-900">{formatCurrency(shopifyOrder.total_price, shopifyOrder.currency)}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">Subtotal</p>
                                    <p className="text-sm font-medium text-gray-900">{formatCurrency(shopifyOrder.subtotal_price, shopifyOrder.currency)}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">Impuestos</p>
                                    <p className="text-sm font-medium text-gray-900">{formatCurrency(shopifyOrder.total_tax, shopifyOrder.currency)}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">Moneda</p>
                                    <p className="text-sm font-medium text-gray-900">{shopifyOrder.currency || 'N/A'}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">Fecha de Creación</p>
                                    <p className="text-sm font-medium text-gray-900">{formatDate(shopifyOrder.created_at)}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-gray-500 uppercase">Fecha de Procesamiento</p>
                                    <p className="text-sm font-medium text-gray-900">{formatDate(shopifyOrder.processed_at)}</p>
                                </div>
                            </div>
                        </div>

                        {/* Cliente */}
                        {shopifyOrder.customer && (
                            <div className="bg-gray-50 rounded-lg p-5">
                                <h3 className="text-lg font-semibold text-gray-900 mb-4">Cliente</h3>
                                <div className="grid grid-cols-4 gap-4">
                                    <div>
                                        <p className="text-xs text-gray-500 uppercase">ID</p>
                                        <p className="text-sm font-medium text-gray-900">{shopifyOrder.customer.id || 'N/A'}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-gray-500 uppercase">Email</p>
                                        <p className="text-sm font-medium text-gray-900">{shopifyOrder.customer.email || 'N/A'}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-gray-500 uppercase">Nombre</p>
                                        <p className="text-sm font-medium text-gray-900">
                                            {shopifyOrder.customer.first_name || ''} {shopifyOrder.customer.last_name || ''}
                                        </p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-gray-500 uppercase">Teléfono</p>
                                        <p className="text-sm font-medium text-gray-900">{shopifyOrder.customer.phone || 'N/A'}</p>
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* Items de Línea */}
                        {shopifyOrder.line_items && shopifyOrder.line_items.length > 0 && (
                            <div className="bg-gray-50 rounded-lg p-5">
                                <h3 className="text-lg font-semibold text-gray-900 mb-4">Items de la Orden</h3>
                                <div className="overflow-x-auto">
                                    <table className="w-full divide-y divide-gray-200">
                                        <thead className="bg-gray-100">
                                            <tr>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase">Producto</th>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase">SKU</th>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase">Cantidad</th>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase">Precio</th>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase">Vendor</th>
                                            </tr>
                                        </thead>
                                        <tbody className="bg-white divide-y divide-gray-200">
                                            {shopifyOrder.line_items.map((item, index) => (
                                                <tr key={item.id || index} className="hover:bg-gray-50">
                                                    <td className="px-4 py-3 text-sm text-gray-900">{item.title || item.name || 'N/A'}</td>
                                                    <td className="px-4 py-3 text-sm text-gray-600">{item.sku || 'N/A'}</td>
                                                    <td className="px-4 py-3 text-sm text-gray-900">{item.quantity || 0}</td>
                                                    <td className="px-4 py-3 text-sm text-gray-900">{formatCurrency(item.price, shopifyOrder.currency)}</td>
                                                    <td className="px-4 py-3 text-sm text-gray-600">{item.vendor || 'N/A'}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        )}

                        {/* Direcciones en grid */}
                        <div className="grid grid-cols-2 gap-4">
                            {/* Dirección de Envío */}
                            {shopifyOrder.shipping_address && (
                                <div className="bg-gray-50 rounded-lg p-5">
                                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Dirección de Envío</h3>
                                    <div className="text-sm text-gray-900">
                                        <p className="font-medium">
                                            {shopifyOrder.shipping_address.first_name || ''} {shopifyOrder.shipping_address.last_name || ''}
                                        </p>
                                        <p>{shopifyOrder.shipping_address.address1 || ''}</p>
                                        {shopifyOrder.shipping_address.address2 && (
                                            <p>{shopifyOrder.shipping_address.address2}</p>
                                        )}
                                        <p>
                                            {shopifyOrder.shipping_address.city || ''}, {shopifyOrder.shipping_address.province || ''}
                                        </p>
                                        <p>
                                            {shopifyOrder.shipping_address.country || ''} {shopifyOrder.shipping_address.zip || ''}
                                        </p>
                                        {shopifyOrder.shipping_address.phone && (
                                            <p className="mt-2 text-gray-600">Tel: {shopifyOrder.shipping_address.phone}</p>
                                        )}
                                    </div>
                                </div>
                            )}

                            {/* Dirección de Facturación */}
                            {shopifyOrder.billing_address && (
                                <div className="bg-gray-50 rounded-lg p-5">
                                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Dirección de Facturación</h3>
                                    <div className="text-sm text-gray-900">
                                        <p className="font-medium">
                                            {shopifyOrder.billing_address.first_name || ''} {shopifyOrder.billing_address.last_name || ''}
                                        </p>
                                        <p>{shopifyOrder.billing_address.address1 || ''}</p>
                                        {shopifyOrder.billing_address.address2 && (
                                            <p>{shopifyOrder.billing_address.address2}</p>
                                        )}
                                        <p>
                                            {shopifyOrder.billing_address.city || ''}, {shopifyOrder.billing_address.province || ''}
                                        </p>
                                        <p>
                                            {shopifyOrder.billing_address.country || ''} {shopifyOrder.billing_address.zip || ''}
                                        </p>
                                        {shopifyOrder.billing_address.phone && (
                                            <p className="mt-2 text-gray-600">Tel: {shopifyOrder.billing_address.phone}</p>
                                        )}
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* Fulfillments */}
                        {shopifyOrder.fulfillments && shopifyOrder.fulfillments.length > 0 && (
                            <div className="bg-gray-50 rounded-lg p-5">
                                <h3 className="text-lg font-semibold text-gray-900 mb-4">Fulfillments</h3>
                                <div className="grid grid-cols-2 gap-4">
                                    {shopifyOrder.fulfillments.map((fulfillment, index) => (
                                        <div key={fulfillment.id || index} className="bg-white rounded-lg p-4 border border-gray-200">
                                            <div className="grid grid-cols-2 gap-3">
                                                <div>
                                                    <p className="text-xs text-gray-500 uppercase">ID</p>
                                                    <p className="text-sm font-medium text-gray-900">{fulfillment.id || 'N/A'}</p>
                                                </div>
                                                <div>
                                                    <p className="text-xs text-gray-500 uppercase">Estado</p>
                                                    <p className="text-sm font-medium text-gray-900 capitalize">{fulfillment.status || 'N/A'}</p>
                                                </div>
                                                {fulfillment.tracking_number && (
                                                    <div>
                                                        <p className="text-xs text-gray-500 uppercase">Número de Rastreo</p>
                                                        <p className="text-sm font-medium text-gray-900">{fulfillment.tracking_number}</p>
                                                    </div>
                                                )}
                                                {fulfillment.tracking_company && (
                                                    <div>
                                                        <p className="text-xs text-gray-500 uppercase">Transportadora</p>
                                                        <p className="text-sm font-medium text-gray-900">{fulfillment.tracking_company}</p>
                                                    </div>
                                                )}
                                                {fulfillment.created_at && (
                                                    <div>
                                                        <p className="text-xs text-gray-500 uppercase">Fecha de Creación</p>
                                                        <p className="text-sm font-medium text-gray-900">{formatDate(fulfillment.created_at)}</p>
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>
                )}

                {/* Vista JSON Crudo */}
                {showRaw && data && (
                    <div>
                        <div className="flex justify-between items-center mb-4">
                            <h3 className="text-lg font-semibold text-gray-900">JSON Crudo</h3>
                            <button
                                onClick={() => setShowRaw(false)}
                                className="px-3 py-1.5 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md transition-colors"
                            >
                                Ver Vista Estructurada
                            </button>
                        </div>
                        <div className="bg-gray-900 rounded-lg p-6 overflow-auto" style={{ maxHeight: 'calc(90vh - 200px)' }}>
                            <pre className="text-green-400 font-mono text-xs sm:text-sm whitespace-pre-wrap">
                                {JSON.stringify(data, null, 2)}
                            </pre>
                        </div>
                    </div>
                )}

                {!loading && !error && !data && (
                    <div className="text-center py-4 text-gray-500">
                        No hay datos disponibles
                    </div>
                )}
            </div>
        </FullWidthModal>
    );
}
