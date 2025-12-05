'use client';

import { useState } from 'react';
import { Order, CreateOrderDTO, UpdateOrderDTO } from '../../domain/types';
import { Button, Input, Alert } from '@/shared/ui';

interface OrderFormProps {
    order?: Order;
    onSuccess?: () => void;
    onCancel?: () => void;
}

export default function OrderForm({ order, onSuccess, onCancel }: OrderFormProps) {
    const isEdit = !!order;

    const [formData, setFormData] = useState({
        // Integration
        integration_id: order?.integration_id || 0,
        platform: order?.platform || '',

        // Customer
        customer_name: order?.customer_name || '',
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
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            // Validation
            if (!formData.customer_name || !formData.total_amount) {
                throw new Error('Por favor completa los campos requeridos');
            }

            // TODO: Call create or update action
            console.log('Form data:', formData);

            if (onSuccess) onSuccess();
        } catch (err: any) {
            setError(err.message || 'Error al guardar la orden');
        } finally {
            setLoading(false);
        }
    };

    // Auto-calculate total
    const calculateTotal = () => {
        const total = formData.subtotal + formData.tax - formData.discount + formData.shipping_cost;
        setFormData({ ...formData, total_amount: total });
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            {/* Integration Info */}
            <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-base font-semibold text-gray-800 mb-4">Información de Integración</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            ID de Integración *
                        </label>
                        <Input
                            type="number"
                            required
                            value={formData.integration_id}
                            onChange={(e) => setFormData({ ...formData, integration_id: parseInt(e.target.value) })}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Plataforma *
                        </label>
                        <select
                            required
                            value={formData.platform}
                            onChange={(e) => setFormData({ ...formData, platform: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                            <option value="">Seleccionar...</option>
                            <option value="shopify">Shopify</option>
                            <option value="woocommerce">WooCommerce</option>
                            <option value="manual">Manual</option>
                        </select>
                    </div>
                </div>
            </div>

            {/* Customer Info */}
            <div className="bg-blue-50 p-4 rounded-lg">
                <h3 className="text-base font-semibold text-gray-800 mb-4">Información del Cliente</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Nombre *
                        </label>
                        <Input
                            type="text"
                            required
                            value={formData.customer_name}
                            onChange={(e) => setFormData({ ...formData, customer_name: e.target.value })}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Email
                        </label>
                        <Input
                            type="email"
                            value={formData.customer_email}
                            onChange={(e) => setFormData({ ...formData, customer_email: e.target.value })}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Teléfono
                        </label>
                        <Input
                            type="tel"
                            value={formData.customer_phone}
                            onChange={(e) => setFormData({ ...formData, customer_phone: e.target.value })}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            DNI
                        </label>
                        <Input
                            type="text"
                            value={formData.customer_dni}
                            onChange={(e) => setFormData({ ...formData, customer_dni: e.target.value })}
                        />
                    </div>
                </div>
            </div>

            {/* Shipping Address */}
            <div className="bg-purple-50 p-4 rounded-lg">
                <h3 className="text-base font-semibold text-gray-800 mb-4">Dirección de Envío</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="md:col-span-2">
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Dirección
                        </label>
                        <Input
                            type="text"
                            value={formData.shipping_street}
                            onChange={(e) => setFormData({ ...formData, shipping_street: e.target.value })}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Ciudad
                        </label>
                        <Input
                            type="text"
                            value={formData.shipping_city}
                            onChange={(e) => setFormData({ ...formData, shipping_city: e.target.value })}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Departamento
                        </label>
                        <Input
                            type="text"
                            value={formData.shipping_state}
                            onChange={(e) => setFormData({ ...formData, shipping_state: e.target.value })}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            País
                        </label>
                        <Input
                            type="text"
                            value={formData.shipping_country}
                            onChange={(e) => setFormData({ ...formData, shipping_country: e.target.value })}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
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

            {/* Financial */}
            <div className="bg-green-50 p-4 rounded-lg">
                <h3 className="text-base font-semibold text-gray-800 mb-4">Información Financiera</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
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
                        <label className="block text-sm font-medium text-gray-700 mb-1">
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
                        <label className="block text-sm font-medium text-gray-700 mb-1">
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
                        <label className="block text-sm font-medium text-gray-700 mb-1">
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
                        <label className="block text-sm font-medium text-gray-700 mb-1">
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
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Moneda
                        </label>
                        <select
                            value={formData.currency}
                            onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                            <option value="COP">COP</option>
                            <option value="USD">USD</option>
                        </select>
                    </div>
                </div>
            </div>

            {/* Payment & Status */}
            <div className="bg-yellow-50 p-4 rounded-lg">
                <h3 className="text-base font-semibold text-gray-800 mb-4">Pago y Estado</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="flex items-center">
                            <input
                                type="checkbox"
                                checked={formData.is_paid}
                                onChange={(e) => setFormData({ ...formData, is_paid: e.target.checked })}
                                className="mr-2"
                            />
                            <span className="text-sm font-medium text-gray-700">Orden Pagada</span>
                        </label>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Estado
                        </label>
                        <select
                            value={formData.status}
                            onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                            <option value="pending">Pendiente</option>
                            <option value="processing">Procesando</option>
                            <option value="shipped">Enviado</option>
                            <option value="delivered">Entregado</option>
                            <option value="cancelled">Cancelado</option>
                        </select>
                    </div>
                </div>
            </div>

            {/* Actions */}
            <div className="flex justify-end space-x-3 pt-4 border-t">
                {onCancel && (
                    <Button type="button" onClick={onCancel} variant="outline">
                        Cancelar
                    </Button>
                )}
                <Button type="submit" disabled={loading} loading={loading} variant="primary">
                    {isEdit ? 'Actualizar Orden' : 'Crear Orden'}
                </Button>
            </div>
        </form>
    );
}
