'use server';

import { cookies } from 'next/headers';
import { PaymentMethodApiRepository } from '../repository/api-repository';
import { PaymentMethodUseCases } from '../../app/use-cases';
import { PaymentMethodsResponse } from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new PaymentMethodApiRepository(token);
    return new PaymentMethodUseCases(repository);
}

export const getPaymentMethodsAction = async (): Promise<PaymentMethodsResponse> => {
    try {
        return await (await getUseCases()).getPaymentMethods();
    } catch (error: any) {
        console.error('Get Payment Methods Action Error:', error.message);
        throw new Error(error.message);
    }
};
