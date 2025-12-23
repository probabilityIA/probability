'use server';

import { cookies } from 'next/headers';
import { DashboardApiRepository } from '../repository/api-repository';
import { DashboardUseCases } from '../../app/use-cases';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new DashboardApiRepository(token);
    return new DashboardUseCases(repository);
}

export const getDashboardStatsAction = async (businessId?: number, integrationId?: number) => {
    try {
        return await (await getUseCases()).getStats(businessId, integrationId);
    } catch (error: any) {
        console.error('Get Dashboard Stats Action Error:', error.message);
        throw new Error(error.message);
    }
};
