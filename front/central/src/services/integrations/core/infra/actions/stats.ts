'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';

export interface IntegrationStatsItem {
    integration_id: number;
    orders_count: number;
    orders_in_progress: number;
    orders_delivered: number;
    orders_cancelled: number;
    orders_returned: number;
    products_count: number;
    last_order_at?: string;
}

export interface IntegrationStatsResponse {
    success: boolean;
    message?: string;
    data: IntegrationStatsItem[];
}

export const getIntegrationStatsAction = async (businessId?: number): Promise<IntegrationStatsResponse> => {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || '';

        const queryParams = new URLSearchParams();
        if (businessId) {
            queryParams.append('business_id', businessId.toString());
        }

        const url = `${env.API_BASE_URL}/integrations/stats${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;

        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            cache: 'no-store',
        });

        if (!response.ok) {
            throw new Error('Error al obtener estadisticas de integraciones');
        }

        return await response.json();
    } catch (error) {
        return {
            success: false,
            message: error instanceof Error ? error.message : 'Error desconocido',
            data: [],
        };
    }
};
