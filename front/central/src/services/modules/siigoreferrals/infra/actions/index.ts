'use server';

import { cookies } from 'next/headers';
import { SiigoReferralsApiRepository } from '../repository/api-repository';
import { SiigoReferralsUseCases } from '../../app/use-cases';
import { GetSiigoReferralsParams } from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new SiigoReferralsApiRepository(token);
    return new SiigoReferralsUseCases(repository);
}

export const getSiigoReferralsAction = async (params?: GetSiigoReferralsParams) => {
    try {
        return await (await getUseCases()).getReferrals(params);
    } catch (error: any) {
        console.error('Get Siigo Referrals Action Error:', error.message);
        throw new Error(error.message);
    }
};
