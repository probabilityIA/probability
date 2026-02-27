'use server';

import { cookies } from 'next/headers';
import { PayGatewayApiRepository } from '../repository/api-repository';
import { PayGatewayUseCases } from '../../app/use-cases';
import { PaymentGatewayType } from '../../domain/types';

export async function getPaymentGatewayTypesAction(): Promise<PaymentGatewayType[]> {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || null;
        const repo = new PayGatewayApiRepository(token);
        const useCases = new PayGatewayUseCases(repo);
        return await useCases.getPaymentGatewayTypes();
    } catch (error: any) {
        console.error('getPaymentGatewayTypesAction error:', error.message);
        return [];
    }
}
