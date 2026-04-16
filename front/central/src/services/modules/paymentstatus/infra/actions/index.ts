'use server';

import { cookies } from 'next/headers';
import { PaymentStatusApiRepository } from '../repository/api-repository';
import { PaymentStatusUseCases } from '../../app/use-cases';
import { GetPaymentStatusesParams, PaymentStatusesResponse } from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new PaymentStatusApiRepository(token);
    return new PaymentStatusUseCases(repository);
}

export const getPaymentStatusesAction = async (
    params?: GetPaymentStatusesParams
): Promise<PaymentStatusesResponse> => {
    try {
        return await (await getUseCases()).getPaymentStatuses(params);
    } catch (error: any) {
        console.error('Get Payment Statuses Action Error:', error.message);
        throw new Error(error.message);
    }
};
