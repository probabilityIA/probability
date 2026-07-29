'use server';

import { cookies } from 'next/headers';
import { MarketingLeadsApiRepository } from '../repository/api-repository';
import { MarketingLeadsUseCases } from '../../app/use-cases';
import { GetMarketingLeadsParams } from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new MarketingLeadsApiRepository(token);
    return new MarketingLeadsUseCases(repository);
}

export const getMarketingLeadsAction = async (params?: GetMarketingLeadsParams) => {
    try {
        return await (await getUseCases()).getLeads(params);
    } catch (error: any) {
        console.error('Get Marketing Leads Action Error:', error.message);
        throw new Error(error.message);
    }
};
