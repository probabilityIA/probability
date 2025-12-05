'use client';

import { Order } from '../../domain/types';

interface OrderDetailsProps {
    order: Order;
}

export default function OrderDetails({ order }: OrderDetailsProps) {
    const formatCurrency = (amount: number, currency: string = 'USD') => {
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency,
        }).format(amount);
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('es-CO');
    };

    return (
        <div className="space-y-6">
            {/* Order Information */}
            <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">Información de la Orden</h3>
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <p className="text-sm text-gray-500">Número de Orden</p>
                        <p className="text-sm font-medium text-gray-900">{order.order_number}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500">Número Interno</p>
                        <p className="text-sm font-medium text-gray-900">{order.internal_number}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500">Plataforma</p>
                        <p className="text-sm font-medium text-gray-900 capitalize">{order.platform}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500">Estado</p>
                        <span className="inline-block px-2 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-800">
                            {order.status}
                        </span>
                    </div>
                </div>
            </div>

            {/* Customer Information */}
            <div className="bg-blue-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">Información del Cliente</h3>
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <p className="text-sm text-gray-500">Nombre</p>
                        <p className="text-sm font-medium text-gray-900">{order.customer_name}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500">Email</p>
                        <p className="text-sm font-medium text-gray-900">{order.customer_email}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500">Teléfono</p>
                        <p className="text-sm font-medium text-gray-900">{order.customer_phone}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500">DNI</p>
                        <p className="text-sm font-medium text-gray-900">{order.customer_dni}</p>
                    </div>
                </div>
            </div>

            {/* Shipping Address */}
            <div className="bg-purple-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">Dirección de Envío</h3>
                <div className="space-y-2">
                    <p className="text-sm text-gray-900">{order.shipping_street}</p>
                    <p className="text-sm text-gray-900">
                        {order.shipping_city}, {order.shipping_state} {order.shipping_postal_code}
                    </p>
                    <p className="text-sm text-gray-900">{order.shipping_country}</p>
                </div>
            </div>

            {/* Financial Summary */}
            <div className="bg-green-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">Resumen Financiero</h3>
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
                    <div className="flex justify-between">
                        <span className="text-sm text-gray-600">Descuento</span>
                        <span className="text-sm font-medium text-gray-900">
                            -{formatCurrency(order.discount, order.currency)}
                        </span>
                    </div>
                    <div className="flex justify-between">
                        <span className="text-sm text-gray-600">Envío</span>
                        <span className="text-sm font-medium text-gray-900">
                            {formatCurrency(order.shipping_cost, order.currency)}
                        </span>
                    </div>
                    <div className="flex justify-between pt-2 border-t border-green-200">
                        <span className="text-base font-semibold text-gray-900">Total</span>
                        <span className="text-base font-bold text-gray-900">
                            {formatCurrency(order.total_amount, order.currency)}
                        </span>
                    </div>
                </div>
            </div>

            {/* Payment Details */}
            <div className="bg-yellow-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">Detalles de Pago</h3>
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <p className="text-sm text-gray-500">Estado de Pago</p>
                        <span className={`inline-block px-2 py-1 text-xs font-medium rounded-full ${order.is_paid ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                            }`}>
                            {order.is_paid ? 'Pagado' : 'Pendiente'}
                        </span>
                    </div>
                    {order.paid_at && (
                        <div>
                            <p className="text-sm text-gray-500">Fecha de Pago</p>
                            <p className="text-sm font-medium text-gray-900">{formatDate(order.paid_at)}</p>
                        </div>
                    )}
                </div>
            </div>

            {/* Timestamps */}
            <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">Fechas</h3>
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <p className="text-sm text-gray-500">Creado</p>
                        <p className="text-sm font-medium text-gray-900">{formatDate(order.created_at)}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500">Actualizado</p>
                        <p className="text-sm font-medium text-gray-900">{formatDate(order.updated_at)}</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
