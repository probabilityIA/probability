'use server';

import { cookies } from 'next/headers';
import { WhatsAppApiRepository } from '../repository/api-repository';
import { WhatsAppUseCases } from '../../app/use-cases';
import {
    WhatsAppAddNumberValues,
    WhatsAppConnectionResponse,
    WhatsAppConnectionValues,
    WhatsAppEmbeddedSignupConfigResponse,
    WhatsAppEmbeddedSignupPayload,
    WhatsAppEmbeddedSignupResponse,
    WhatsAppNumberResponse,
    WhatsAppBusinessProfileResponse,
    WhatsAppBusinessProfileValues,
    WhatsAppProvisionResponse,
    WhatsAppTemplatesResponse,
} from '../../domain/types';

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

export const saveWhatsAppConnectionAction = async (
    values: WhatsAppConnectionValues,
    businessId?: number,
    token?: string | null
): Promise<WhatsAppConnectionResponse> => {
    try {
        const useCases = await getUseCases(token);
        return await useCases.saveConnection(values, businessId);
    } catch (error: any) {
        return { success: false, message: error?.message || 'Error guardando la conexión de WhatsApp' };
    }
};

const numberAction = async (
    run: (useCases: WhatsAppUseCases) => Promise<WhatsAppNumberResponse>,
    token?: string | null
): Promise<WhatsAppNumberResponse> => {
    try {
        const useCases = await getUseCases(token);
        return await run(useCases);
    } catch (error: any) {
        return { success: false, message: error?.message || 'Error con el n\u00famero de WhatsApp' };
    }
};

export const getWhatsAppNumberStateAction = async (businessId?: number, token?: string | null) =>
    numberAction((useCases) => useCases.getNumberState(businessId), token);

const profileAction = async (
    run: (useCases: WhatsAppUseCases) => Promise<WhatsAppBusinessProfileResponse>,
    token?: string | null
): Promise<WhatsAppBusinessProfileResponse> => {
    try {
        const useCases = await getUseCases(token);
        return await run(useCases);
    } catch (error: any) {
        return { success: false, message: error?.message || 'Error con el perfil de WhatsApp' };
    }
};

export const getWhatsAppBusinessProfileAction = async (businessId?: number, token?: string | null) =>
    profileAction((useCases) => useCases.getBusinessProfile(businessId), token);

export const updateWhatsAppBusinessProfileAction = async (
    values: WhatsAppBusinessProfileValues,
    businessId?: number,
    token?: string | null
) => profileAction((useCases) => useCases.updateBusinessProfile(values, businessId), token);

export const updateWhatsAppBusinessProfilePhotoAction = async (
    file: File,
    businessId?: number,
    token?: string | null
) => profileAction((useCases) => useCases.updateBusinessProfilePhoto(file, businessId), token);

export const addWhatsAppNumberAction = async (
    values: WhatsAppAddNumberValues,
    businessId?: number,
    token?: string | null
) => numberAction((useCases) => useCases.addNumber(values, businessId), token);

export const requestWhatsAppNumberCodeAction = async (
    method: string,
    businessId?: number,
    token?: string | null
) => numberAction((useCases) => useCases.requestNumberCode(method, businessId), token);

export const verifyWhatsAppNumberCodeAction = async (
    code: string,
    businessId?: number,
    token?: string | null
) => numberAction((useCases) => useCases.verifyNumberCode(code, businessId), token);

export const registerWhatsAppNumberAction = async (businessId?: number, token?: string | null) =>
    numberAction((useCases) => useCases.registerNumber(businessId), token);

export const removeWhatsAppNumberAction = async (businessId?: number, token?: string | null) =>
    numberAction((useCases) => useCases.removeNumber(businessId), token);

export const getWhatsAppEmbeddedSignupConfigAction = async (
    token?: string | null
): Promise<WhatsAppEmbeddedSignupConfigResponse> => {
    try {
        const useCases = await getUseCases(token);
        return await useCases.getEmbeddedSignupConfig();
    } catch (error: any) {
        return { success: true, data: { enabled: false }, message: error?.message };
    }
};

export const completeWhatsAppEmbeddedSignupAction = async (
    payload: WhatsAppEmbeddedSignupPayload,
    businessId?: number,
    token?: string | null
): Promise<WhatsAppEmbeddedSignupResponse> => {
    try {
        const useCases = await getUseCases(token);
        return await useCases.completeEmbeddedSignup(payload, businessId);
    } catch (error: any) {
        return { success: false, message: error?.message || 'No se pudo completar el registro' };
    }
};
