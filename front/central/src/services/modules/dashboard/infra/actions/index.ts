'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { DashboardApiRepository } from '../repository/api-repository';
import { DashboardUseCases } from '../../app/use-cases';

async function getUseCases() {
    const token = await getAuthToken();
    const repository = new DashboardApiRepository(token);
    return new DashboardUseCases(repository);
}

export const getDashboardStatsAction = async (
    businessId?: number,
    integrationId?: number
) => {
    try {
        return await (await getUseCases()).getStats(businessId, integrationId);
    } catch (error: any) {
        console.error('Get Dashboard Stats Action Error:', error.message);
        throw new Error(error.message);
    }
};
