'use client';

import { useState } from 'react';
import { Order, CreateOrderDTO, UpdateOrderDTO } from '../../domain/types';
import { Product } from '../../../products/domain/types';
import { Button, Input, Alert, Modal } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useToast } from '@/shared/providers/toast-provider';
import ProductSelector from '../../../products/ui/components/ProductSelector';
import ProductForm from '../../../products/ui/components/ProductForm';
import { createOrderAction, updateOrderAction } from '../../infra/actions';

interface OrderFormProps {
    order?: Order;
    onSuccess?: () => void;
    onCancel?: () => void;
}

export default function OrderForm({ order, onSuccess, onCancel }: OrderFormProps) {
    const isEdit = !!order;
    const { permissions } = usePermissions();
    const defaultBusinessId = permissions?.business_id || 0;
    const { showToast } = useToast();

    const [formData, setFormData] = useState({
        // Integration
        integration_id: order?.integration_id || 0,
        platform: order?.platform || 'manual',
        business_id: order?.business_id || defaultBusinessId,

        // Customer — fallback: split customer_name if first/last aren't set separately
        customer_name: order?.customer_name || '',
        customer_first_name: order?.customer_first_name || (order?.customer_name ? order.customer_name.split(' ')[0] : ''),
        customer_last_name: order?.customer_last_name || (order?.customer_name ? order.customer_name.split(' ').slice(1).join(' ') : ''),
        customer_email: order?.customer_email || '',
        customer_phone: order?.customer_phone || '',
        customer_dni: order?.customer_dni || '',

        // Shipping
        shipping_street: order?.shipping_street || '',
        shipping_city: order?.shipping_city || '',
        shipping_state: order?.shipping_state || '',
        shipping_country: order?.shipping_country || 'Colombia',
        shipping_postal_code: order?.shipping_postal_code || '',

        // Financial
        subtotal: order?.subtotal || 0,
        tax: order?.tax || 0,
        discount: order?.discount || 0,
        shipping_cost: order?.shipping_cost || 0,
        total_amount: order?.total_amount || 0,
        currency: order?.currency || 'COP',

        // Payment
        payment_method_id: order?.payment_method_id || 1,
        is_paid: order?.is_paid || false,

        // Status
        status: order?.status || 'pending',

        // Logistics (preserved on update)
        tracking_number: order?.tracking_number || '',
        tracking_link: order?.tracking_link || '',
        guide_id: order?.guide_id || '',
        warehouse_name: order?.warehouse_name || '',
        driver_name: order?.driver_name || '',
        is_last_mile: order?.is_last_mile || false,

        // Additional
        notes: order?.notes || '',
        invoiceable: order?.invoiceable ?? false,
        is_confirmed: order?.is_confirmed ?? null,
        novelty: order?.novelty || '',

        // Items
        items: order?.items || [],

        // Extra
        integration_type: order?.integration_type || '',
        external_id: order?.external_id || '',
    });

    const [selectedProducts, setSelectedProducts] = useState<Product[]>(() => {
        if (!order?.items || !Array.isArray(order.items)) return [];
        return (order.items as any[])
            .map((item: any) => ({
                ...item,
                // Normalize field names that differ between stored format and Product interface
                stock: item.stock ?? item.stock_quantity ?? 0,
                manage_stock: item.manage_stock ?? item.track_inventory ?? false,
                thumbnail: item.thumbnail || item.image_url || undefined,
            } as Product))
            .filter((p: any) => p.id);
    });
    const [showProductModal, setShowProductModal] = useState(false);

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            // Validation
            if ((!formData.customer_name && !formData.customer_first_name) || !formData.total_amount) {
                throw new Error('Por favor completa los campos requeridos');
            }

            const baseData = {
                ...formData,
                items: selectedProducts.length > 0 ? selectedProducts : formData.items,
                customer_name: formData.customer_name || `${formData.customer_first_name} ${formData.customer_last_name}`.trim()
            };

            let response;
            if (isEdit && order) {
                const updateData: UpdateOrderDTO = {
                    ...baseData,
                    // Map is_confirmed → confirmation_status so clearing to null works
                    confirmation_status: formData.is_confirmed === true ? 'yes' : formData.is_confirmed === false ? 'no' : 'pending',
                    is_confirmed: undefined,
                };
                response = await updateOrderAction(order.id, updateData);
            } else {
                response = await createOrderAction(baseData as CreateOrderDTO);
            }

            if (response.success) {
                showToast(isEdit ? 'Orden actualizada exitosamente' : 'Orden creada exitosamente', 'success');
                if (onSuccess) onSuccess();
            } else {
                setError(response.message || 'Error al guardar la orden');
                showToast(response.message || 'Error al guardar la orden', 'error');
            }
        } catch (err: any) {
            setError(err.message || 'Error al guardar la orden');
            showToast(err.message || 'Error al guardar la orden', 'error');
        } finally {
            setLoading(false);
        }
    };

    const handleProductsChange = (products: Product[]) => {
        setSelectedProducts(products);
        const subtotal = products.reduce((acc, p) => acc + p.price, 0);
        setFormData(prev => ({
            ...prev,
            subtotal,
            total_amount: subtotal + prev.tax - prev.discount + prev.shipping_cost,
        }));
    };

    // Auto-calculate total
    const calculateTotal = () => {
        setFormData(prev => ({
            ...prev,
            total_amount: prev.subtotal + prev.tax - prev.discount + prev.shipping_cost,
        }));
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-2 bg-purple-50 p-2 rounded-lg">
            {isEdit && order && (
                <div className="flex items-center gap-3 px-2 py-2 mb-1 bg-purple-100 rounded-lg border border-purple-200">
                    {order.integration_logo_url && (
                        <img src={order.integration_logo_url} alt="" className="h-8 w-8 object-contain rounded" />
                    )}
                    <div>
                        <p className="text-xs text-purple-600 font-medium">
                            {order.integration_name || order.integration_type || 'Integración'}
                        </p>
                        <p className="text-sm font-bold text-gray-800">
                            {order.order_number || order.internal_number || order.id}
                        </p>
                    </div>
                </div>
            )}
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            {/* 3-Column Layout with 2-Column below */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-2">
                {/* Column 1: Customer Info */}
                <div>
                    {/* Customer Info */}
                    <div className="bg-purple-300 p-5 rounded-lg">
                        <h3 className="text-base font-semibold text-black mb-4">Información del Cliente</h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Nombre *
                                </label>
                                <Input
                                    type="text"
                                    required
                                    value={formData.customer_first_name}
                                    onChange={(e) => setFormData({ ...formData, customer_first_name: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Apellido *
                                </label>
                                <Input
                                    type="text"
                                    required
                                    value={formData.customer_last_name}
                                    onChange={(e) => setFormData({ ...formData, customer_last_name: e.target.value })}
                                />
                            </div>
                            <div className="md:col-span-2">
                                <label className="block text-sm font-medium text-black mb-2">
                                    📧 Email *
                                </label>
                                <Input
                                    type="email"
                                    value={formData.customer_email}
                                    onChange={(e) => setFormData({ ...formData, customer_email: e.target.value })}
                                />
                            </div>
                            <div className="md:col-span-2">
                                <label className="block text-sm font-medium text-black mb-2">
                                    📱 Teléfono *
                                </label>
                                <div className="flex items-center w-full">
                                    <span className="px-3 py-2 bg-purple-200 text-black font-medium rounded-l-lg border border-r-0 border-gray-300">
                                        +57
                                    </span>
                                    <div className="flex-1">
                                        <Input
                                            type="tel"
                                            value={formData.customer_phone}
                                            onChange={(e) => setFormData({ ...formData, customer_phone: e.target.value })}
                                            placeholder="300 1234567"
                                            className="rounded-l-none w-full"
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                {/* Column 2: Shipping Address */}
                <div>
                    {/* Shipping Address */}
                    <div className="bg-purple-300 p-5 rounded-lg">
                        <h3 className="text-base font-semibold text-black mb-4">Dirección de Envío</h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="md:col-span-2">
                                <label className="block text-sm font-medium text-black mb-2">
                                    Dirección
                                </label>
                                <Input
                                    type="text"
                                    value={formData.shipping_street}
                                    onChange={(e) => setFormData({ ...formData, shipping_street: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Ciudad
                                </label>
                                <Input
                                    type="text"
                                    value={formData.shipping_city}
                                    onChange={(e) => setFormData({ ...formData, shipping_city: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Departamento
                                </label>
                                <Input
                                    type="text"
                                    value={formData.shipping_state}
                                    onChange={(e) => setFormData({ ...formData, shipping_state: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    País
                                </label>
                                <Input
                                    type="text"
                                    value={formData.shipping_country}
                                    onChange={(e) => setFormData({ ...formData, shipping_country: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Código Postal
                                </label>
                                <Input
                                    type="text"
                                    value={formData.shipping_postal_code}
                                    onChange={(e) => setFormData({ ...formData, shipping_postal_code: e.target.value })}
                                />
                            </div>
                        </div>
                    </div>
                </div>

                {/* Column 3: Financial Info only */}
                <div>
                    {/* Financial */}
                    <div className="bg-purple-300 p-5 rounded-lg">
                        <h3 className="text-base font-semibold text-black mb-4">Información Financiera</h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Subtotal *
                                </label>
                                <Input
                                    type="number"
                                    step="0.01"
                                    required
                                    value={formData.subtotal}
                                    onChange={(e) => setFormData({ ...formData, subtotal: parseFloat(e.target.value) || 0 })}
                                    onBlur={calculateTotal}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Impuestos
                                </label>
                                <Input
                                    type="number"
                                    step="0.01"
                                    value={formData.tax}
                                    onChange={(e) => setFormData({ ...formData, tax: parseFloat(e.target.value) || 0 })}
                                    onBlur={calculateTotal}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Descuento
                                </label>
                                <Input
                                    type="number"
                                    step="0.01"
                                    value={formData.discount}
                                    onChange={(e) => setFormData({ ...formData, discount: parseFloat(e.target.value) || 0 })}
                                    onBlur={calculateTotal}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Costo de Envío
                                </label>
                                <Input
                                    type="number"
                                    step="0.01"
                                    value={formData.shipping_cost}
                                    onChange={(e) => setFormData({ ...formData, shipping_cost: parseFloat(e.target.value) || 0 })}
                                    onBlur={calculateTotal}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Total *
                                </label>
                                <Input
                                    type="number"
                                    step="0.01"
                                    required
                                    value={formData.total_amount}
                                    onChange={(e) => setFormData({ ...formData, total_amount: parseFloat(e.target.value) || 0 })}
                                    className="font-bold"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Moneda
                                </label>
                                <select
                                    value={formData.currency}
                                    onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                                >
                                    <option value="COP">COP</option>
                                    <option value="USD">USD</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Product Selection - spans 2 columns */}
                <div className="lg:col-span-2">
                    <div className="bg-purple-300 p-5 rounded-lg">
                        <h3 className="text-base font-semibold text-black mb-4">Productos</h3>
                        <ProductSelector
                            businessId={formData.business_id || 0}
                            selectedProducts={selectedProducts}
                            onSelect={handleProductsChange}
                            onCreateNew={() => setShowProductModal(true)}
                        />
                    </div>
                </div>

                {/* Payment & Status - right side */}
                <div>
                    <div className="bg-purple-300 p-5 rounded-lg">
                        <h3 className="text-base font-semibold text-black mb-4">Pago y Estado</h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label className="flex items-center">
                                    <input
                                        type="checkbox"
                                        checked={formData.is_paid}
                                        onChange={(e) => setFormData({ ...formData, is_paid: e.target.checked })}
                                        className="mr-2"
                                    />
                                    <span className="text-sm font-medium text-black">Orden Pagada</span>
                                </label>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Estado
                                </label>
                                <select
                                    value={formData.status}
                                    onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                                >
                                    <option value="pending">Pendiente</option>
                                    <option value="processing">Procesando</option>
                                    <option value="shipped">Enviado</option>
                                    <option value="delivered">Entregado</option>
                                    <option value="cancelled">Cancelado</option>
                                </select>
                            </div>
                            <div>
                                <label className="flex items-center">
                                    <input
                                        type="checkbox"
                                        checked={formData.invoiceable}
                                        onChange={(e) => setFormData({ ...formData, invoiceable: e.target.checked })}
                                        className="mr-2"
                                    />
                                    <span className="text-sm font-medium text-black">Facturable</span>
                                </label>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Confirmación
                                </label>
                                <select
                                    value={formData.is_confirmed === true ? 'yes' : formData.is_confirmed === false ? 'no' : 'pending'}
                                    onChange={(e) => {
                                        const v = e.target.value;
                                        setFormData({ ...formData, is_confirmed: v === 'yes' ? true : v === 'no' ? false : null });
                                    }}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                                >
                                    <option value="pending">Pendiente</option>
                                    <option value="yes">Confirmado</option>
                                    <option value="no">No confirmado</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Logistics */}
                <div className="lg:col-span-2">
                    <div className="bg-purple-300 p-5 rounded-lg">
                        <h3 className="text-base font-semibold text-black mb-4">Logística</h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Número de Guía
                                </label>
                                <Input
                                    type="text"
                                    value={formData.tracking_number}
                                    onChange={(e) => setFormData({ ...formData, tracking_number: e.target.value })}
                                    placeholder="ej: 123456789"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    ID Guía
                                </label>
                                <Input
                                    type="text"
                                    value={formData.guide_id}
                                    onChange={(e) => setFormData({ ...formData, guide_id: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Bodega
                                </label>
                                <Input
                                    type="text"
                                    value={formData.warehouse_name}
                                    onChange={(e) => setFormData({ ...formData, warehouse_name: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-black mb-2">
                                    Conductor
                                </label>
                                <Input
                                    type="text"
                                    value={formData.driver_name}
                                    onChange={(e) => setFormData({ ...formData, driver_name: e.target.value })}
                                />
                            </div>
                        </div>
                    </div>
                </div>

                {/* Notes */}
                <div className="lg:col-span-1">
                    <div className="bg-purple-300 p-5 rounded-lg h-full">
                        <h3 className="text-base font-semibold text-black mb-4">Notas</h3>
                        <textarea
                            value={formData.notes}
                            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                            rows={5}
                            placeholder="Notas internas sobre la orden..."
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent resize-none text-sm"
                        />
                    </div>
                </div>
            </div>

            {/* Actions */}
            <div className="flex justify-end space-x-3 pt-2 border-t">
                {onCancel && (
                    <button
                        type="button"
                        onClick={onCancel}
                        style={{ background: '#6d28d9' }}
                        className="px-6 py-3 text-base text-white font-semibold rounded-lg hover:shadow-lg hover:scale-105 transition-all"
                    >
                        Cancelar
                    </button>
                )}
                <button
                    type="submit"
                    disabled={loading}
                    style={{ background: '#6d28d9' }}
                    className="px-6 py-3 text-base text-white font-semibold rounded-lg hover:shadow-lg hover:scale-105 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                    {loading && <div className="spinner w-4 h-4" />}
                    {isEdit ? 'Actualizar Orden' : 'Crear Orden'}
                </button>
            </div>
            <Modal
                isOpen={showProductModal}
                onClose={() => setShowProductModal(false)}
                title="Crear Nuevo Producto"
            >
                <ProductForm
                    onSuccess={() => {
                        setShowProductModal(false);
                        // Optional: Refresh products list if needed, but ProductSelector handles it
                    }}
                    onCancel={() => setShowProductModal(false)}
                />
            </Modal>
        </form>
    );
}
