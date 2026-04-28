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

export async function getBoldSignatureAction(amount: number, businessId?: number): Promise<any> {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || null;
        const repo = new PayGatewayApiRepository(token);
        return await repo.getBoldSignature(amount, businessId);
    } catch (error: any) {
        console.error('getBoldSignatureAction error:', error.message);
        return { success: false, message: error.message };
    }
}

export async function syncBoldRechargeAction(orderId: string, businessId?: number): Promise<any> {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || null;
        const repo = new PayGatewayApiRepository(token);
        return await repo.syncBoldRecharge(orderId, businessId);
    } catch (error: any) {
        console.error('syncBoldRechargeAction error:', error.message);
        return { success: false, message: error.message };
    }
}

export async function simulateBoldPaymentAction(orderId: string, amount: number, businessId?: number): Promise<any> {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || null;
        const repo = new PayGatewayApiRepository(token);
        return await repo.simulateBoldPayment(orderId, amount, businessId);
    } catch (error: any) {
        console.error('simulateBoldPaymentAction error:', error.message);
        return { success: false, message: error.message };
    }
}
