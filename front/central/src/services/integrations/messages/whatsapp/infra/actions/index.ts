'use server';

import { cookies } from 'next/headers';
import { WhatsAppApiRepository } from '../repository/api-repository';
import { WhatsAppUseCases } from '../../app/use-cases';
import { WhatsAppProvisionResponse, WhatsAppTemplatesResponse } from '../../domain/types';

async function getUseCases(tokenOverride?: string | null) {
    const cookieStore = await cookies();
    const token = tokenOverride || cookieStore.get('session_token')?.value || null;

    if (!token) {
        throw new Error('No se encontró sesión activa (Token ausente)');
    }

    return new WhatsAppUseCases(new WhatsAppApiRepository(token));
}

export const getWhatsAppTemplatesStatusAction = async (
    businessId?: number,
    refresh = false,
    token?: string | null
): Promise<WhatsAppTemplatesResponse> => {
    try {
        const useCases = await getUseCases(token);
        return await useCases.getTemplatesStatus(businessId, refresh);
    } catch (error: any) {
        return { success: false, message: error?.message || 'Error consultando las plantillas' };
    }
};

export const provisionWhatsAppTemplatesAction = async (
    businessId?: number,
    token?: string | null
): Promise<WhatsAppProvisionResponse> => {
    try {
        const useCases = await getUseCases(token);
        return await useCases.provisionTemplates(businessId);
    } catch (error: any) {
        return { success: false, message: error?.message || 'Error aprovisionando las plantillas' };
    }
};
