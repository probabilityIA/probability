'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { OrderApiRepository } from '../repository/api-repository';
import { OrderUseCases } from '../../app/use-cases';
import {
    GetOrdersParams,
    CreateOrderDTO,
    UpdateOrderDTO,
    ChangeOrderStatusDTO
} from '../../domain/types';

async function getUseCases() {
    const token = await getAuthToken();
    const repository = new OrderApiRepository(token);
    return new OrderUseCases(repository);
}

export const getOrdersAction = async (params?: GetOrdersParams) => {
    try {
        return await (await getUseCases()).getOrders(params);
    } catch (error: any) {
        console.error('Get Orders Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getOrderByIdAction = async (id: string) => {
    try {
        return await (await getUseCases()).getOrderById(id);
    } catch (error: any) {
        console.error('Get Order By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getOrderHistoryAction = async (orderId: string) => {
    try {
        return await (await getUseCases()).getOrderHistory(orderId);
    } catch (error: any) {
        console.error('Get Order History Action Error:', error.message);
        return { success: false, data: [] };
    }
};

export const createOrderAction = async (data: CreateOrderDTO) => {
    try {
        return await (await getUseCases()).createOrder(data);
    } catch (error: any) {
        console.error('Create Order Action Error:', error.message);
        return { success: false, message: error.message || 'Error al crear orden', data: null };
    }
};

export const updateOrderAction = async (id: string, data: UpdateOrderDTO) => {
    try {
        return await (await getUseCases()).updateOrder(id, data);
    } catch (error: any) {
        console.error('Update Order Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const changeOrderStatusAction = async (id: string, data: ChangeOrderStatusDTO) => {
    try {
        return await (await getUseCases()).changeOrderStatus(id, data);
    } catch (error: any) {
        console.error('Change Order Status Action Error:', error.message);
        return { success: false, message: error.message || 'Error al cambiar estado', data: null };
    }
};

export const deleteOrderAction = async (id: string) => {
    try {
        return await (await getUseCases()).deleteOrder(id);
    } catch (error: any) {
        console.error('Delete Order Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getOrderRawAction = async (id: string) => {
    try {
        return await (await getUseCases()).getOrderRaw(id);
    } catch (error: any) {
        console.error('Get Order Raw Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getAIRecommendationAction = async (origin: string, destination: string) => {
    try {
        return await (await getUseCases()).getAIRecommendation(origin, destination);
    } catch (error: any) {
        console.warn('AI Recommendation no disponible:', error.message);
        return null;
    }
};

export const requestWhatsAppConfirmationAction = async (orderId: string) => {
    try {
        return await (await getUseCases()).requestConfirmation(orderId);
    } catch (error: any) {
        console.error('Request WhatsApp Confirmation Error:', error.message);
        return { success: false, message: error.message || 'Error al enviar confirmación WhatsApp' };
    }
};

export const uploadBulkOrdersAction = async (file: File, businessId?: number) => {
    try {
        const token = await getAuthToken();
        const { env } = await import('@/shared/config/env');

        const formData = new FormData();
        formData.append('file', file);

        const url = businessId
            ? `${env.API_BASE_URL}/orders/upload-bulk?business_id=${businessId}`
            : `${env.API_BASE_URL}/orders/upload-bulk`;

        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
            body: formData,
            cache: 'no-store',
        });

        const result = await response.json();

        if (response.ok && result.success) {
            return {
                success: true,
                data: result.data || {
                    total_rows: 0,
                    success_count: 0,
                    failed_count: 0,
                    errors: []
                }
            };
        } else {
            return {
                success: false,
                message: result.message || 'Error al procesar el archivo',
                data: result.data
            };
        }
    } catch (error: any) {
        console.error('Upload Bulk Orders Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al cargar el archivo',
        };
    }
};

export const checkWhatsAppIntegrationAction = async (businessId: number): Promise<boolean> => {
    try {
        const token = await getAuthToken();
        const { env } = await import('@/shared/config/env');
        const res = await fetch(
            `${env.API_BASE_URL}/integrations/check?integration_type_id=2&business_id=${businessId}`,
            {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'application/json',
                },
                cache: 'no-store',
            }
        );
        if (!res.ok) return false;
        const data = await res.json();
        return data.exists === true;
    } catch {
        return false;
    }
};

export const downloadOrderTemplateAction = async () => {
    const headers = [
        'order_number',
        'customer_name',
        'customer_first_name',
        'customer_last_name',
        'customer_email',
        'customer_phone',
        'customer_dni',
        'shipping_street',
        'shipping_city',
        'shipping_state',
        'shipping_country',
        'shipping_postal_code',
        'shipping_lat',
        'shipping_lng',
        'subtotal',
        'tax',
        'discount',
        'shipping_cost',
        'shipping_discount',
        'total_amount',
        'currency',
        'weight',
        'height',
        'width',
        'length',
        'platform',
        'status',
        'payment_method_id',
        'is_paid',
        'tracking_number',
        'guide_id',
        'warehouse_name',
        'driver_name',
        'notes',
        'order_type_name',
        'invoiceable'
    ];

    const exampleRows = [
        [
            'ORD-001',
            'Juan Perez',
            'Juan',
            'Perez',
            'juan@example.com',
            '3001234567',
            '1234567890',
            'Calle 1 # 2-3',
            'Bogota',
            'Cundinamarca',
            'Colombia',
            '110111',
            '4.7110',
            '-74.0721',
            '45000',
            '5000',
            '0',
            '2000',
            '0',
            '50000',
            'COP',
            '1',
            '10',
            '10',
            '10',
            'manual',
            'pending',
            '1',
            'false',
            'TRACK001',
            'GUIDE001',
            'Warehouse 1',
            'Driver 1',
            'Nota de ejemplo',
            'standard',
            'true'
        ],
        [
            'ORD-002',
            'Maria Lopez',
            'Maria',
            'Lopez',
            'maria@example.com',
            '3109876543',
            '9876543210',
            'Carrera 4 # 5-6',
            'Medellin',
            'Antioquia',
            'Colombia',
            '50001',
            '6.2442',
            '-75.5812',
            '100000',
            '15000',
            '5000',
            '3000',
            '500',
            '120000',
            'COP',
            '2',
            '15',
            '15',
            '15',
            'shopify',
            'shipped',
            '1',
            'true',
            'TRACK002',
            'GUIDE002',
            'Warehouse 2',
            'Driver 2',
            '',
            'express',
            'true'
        ]
    ];

    const csvContent = [
        headers.join(','),
        ...exampleRows.map(row => row.map(cell => `"${cell}"`).join(','))
    ].join('\n');

    return {
        success: true,
        data: csvContent,
    };
};
