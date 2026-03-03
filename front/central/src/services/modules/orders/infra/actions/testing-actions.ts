'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { env } from '@/shared/config/env';
import { SimulateShopifyResult } from '../../domain/types';

interface SimulateShopifyResponse {
    success: boolean;
    data?: SimulateShopifyResult;
    error?: string;
}

export async function simulateShopifyAction(
    topic: string,
    count: number
): Promise<SimulateShopifyResponse> {
    try {
        const token = await getAuthToken();
        if (!token) {
            return { success: false, error: 'No autorizado' };
        }

        const response = await fetch(`${env.TESTING_API_URL}/orders/simulate-shopify`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({ topic, count }),
        });

        if (!response.ok) {
            const errorBody = await response.text();
            return {
                success: false,
                error: `Error ${response.status}: ${errorBody}`,
            };
        }

        const result = await response.json();
        return {
            success: true,
            data: result.data,
        };
    } catch (error: any) {
        console.error('Simulate Shopify Action Error:', error.message);
        return {
            success: false,
            error: error.message || 'Error desconocido',
        };
    }
}
