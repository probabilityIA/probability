'use server';

import { cookies } from 'next/headers';
import { DashboardApiRepository } from '../repository/api-repository';
import { DashboardUseCases } from '../../app/use-cases';

/**
 * Get use cases con token desde cookie o parámetro explícito
 * @param explicitToken Token opcional para iframes donde cookies están bloqueadas
 */
async function getUseCases(explicitToken?: string | null) {
    let token = explicitToken;

    // Si no hay token explícito, intentar leer de cookies
    if (!token) {
        const cookieStore = await cookies();
        token = cookieStore.get('session_token')?.value || null;
    }

    const repository = new DashboardApiRepository(token);
    return new DashboardUseCases(repository);
}

/**
 * Get dashboard stats action
 * @param businessId Business ID opcional
 * @param integrationId Integration ID opcional
 * @param token Token opcional para iframes (donde cookies están bloqueadas)
 */
export const getDashboardStatsAction = async (
    businessId?: number,
    integrationId?: number,
    token?: string | null
) => {
    try {
        return await (await getUseCases(token)).getStats(businessId, integrationId);
    } catch (error: any) {
        console.error('Get Dashboard Stats Action Error:', error.message);
        throw new Error(error.message);
    }
};
