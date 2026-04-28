'use server';

import { cookies } from 'next/headers';
import { IntegrationApiRepository } from '../repository/api-repository';
import { IntegrationUseCases } from '../../app/use-cases';
import {
    GetIntegrationsParams,
    CreateIntegrationDTO,
    UpdateIntegrationDTO,
    CreateIntegrationTypeDTO,
    UpdateIntegrationTypeDTO,
    SyncOrdersParams
} from '../../domain/types';
import { env } from '@/shared/config/env';

async function getUseCases(tokenOverride?: string | null) {
    const cookieStore = await cookies();
    const token = tokenOverride || cookieStore.get('session_token')?.value || null;

    if (!token) {
        throw new Error('No se encontró sesión activa (Token ausente)');
    }
    const repository = new IntegrationApiRepository(token);
    return new IntegrationUseCases(repository);
}

/**
 * Helper: ejecuta una acción y retorna { success, message, data } sin lanzar throw.
 * Next.js producción enmascara errores de Server Actions con un mensaje genérico,
 * así que retornamos el error como dato para que el cliente vea el mensaje real.
 */
async function safeAction<T>(fn: () => Promise<T>): Promise<T | { success: false; message: string }> {
    try {
        return await fn();
    } catch (error: any) {
        const message = error?.message || 'Error desconocido';
        console.error('[Server Action Error]', message);
        return { success: false, message } as any;
    }
}

export const getIntegrationsAction = async (params?: GetIntegrationsParams, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getIntegrations(params);
    } catch (error: any) {
        console.error('Get Integrations Action Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al obtener integraciones',
            data: [],
            total: 0,
            page: 1,
            page_size: 10,
            total_pages: 0
        };
    }
};

export const getIntegrationByIdAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.getIntegrationById(id)));
};

export const getIntegrationByTypeAction = async (type: string, businessId?: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.getIntegrationByType(type, businessId)));
};

export const createIntegrationAction = async (data: CreateIntegrationDTO, token?: string | null) => {
    try {
        return await (await getUseCases(token)).createIntegration(data);
    } catch (error: any) {
        console.error('Create Integration Action Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al crear la integración',
            data: null as any
        };
    }
};

export const updateIntegrationAction = async (id: number, data: UpdateIntegrationDTO, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.updateIntegration(id, data)));
};

export const deleteIntegrationAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.deleteIntegration(id)));
};

export const testConnectionAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.testConnection(id)));
};

export const activateIntegrationAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.activateIntegration(id)));
};

export const deactivateIntegrationAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.deactivateIntegration(id)));
};

export const setAsDefaultAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.setAsDefault(id)));
};

export const syncOrdersAction = async (id: number, params?: SyncOrdersParams, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.syncOrders(id, params)));
};

export const getSyncStatusAction = async (id: number, businessId?: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.getSyncStatus(id, businessId)));
};

export const testIntegrationAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.testIntegration(id)));
};

export const testConnectionRawAction = async (typeCode: string, config: any, credentials: any, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.testConnectionRaw(typeCode, config, credentials)));
};

// Integration Types
export const getIntegrationTypesAction = async (categoryId?: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.getIntegrationTypes(categoryId)));
};

export const getActiveIntegrationTypesAction = async (token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.getActiveIntegrationTypes()));
};

export const getIntegrationTypeByIdAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.getIntegrationTypeById(id)));
};

export const getIntegrationTypeByCodeAction = async (code: string, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.getIntegrationTypeByCode(code)));
};

export const createIntegrationTypeAction = async (data: CreateIntegrationTypeDTO, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.createIntegrationType(data)));
};

export const updateIntegrationTypeAction = async (id: number, data: UpdateIntegrationTypeDTO, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.updateIntegrationType(id, data)));
};

export const deleteIntegrationTypeAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.deleteIntegrationType(id)));
};

export const getIntegrationTypePlatformCredentialsAction = async (id: number, token?: string | null): Promise<{ success: boolean; message: string; data: Record<string, unknown>; webhook_urls?: Record<string, string> }> => {
    try {
        return await (await getUseCases(token)).getIntegrationTypePlatformCredentials(id);
    } catch (error: any) {
        console.error('Get Integration Type Platform Credentials Action Error:', error.message);
        return { success: false, message: error.message, data: {} };
    }
};

// Webhooks
export const getWebhookUrlAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.getWebhookUrl(id)));
};

export const listWebhooksAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.listWebhooks(id)));
};

export const deleteWebhookAction = async (id: number, webhookId: string, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.deleteWebhook(id, webhookId)));
};

export const verifyWebhooksAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.verifyWebhooks(id)));
};

export const createWebhookAction = async (id: number, token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.createWebhook(id)));
};

// ============================================
// Simple Actions - Para Dropdowns/Selectores
// ============================================

export const getIntegrationsSimpleAction = async (businessId?: number): Promise<import('../../domain/types').IntegrationsSimpleResponse> => {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || '';

        const queryParams = new URLSearchParams();
        if (businessId) {
            queryParams.append('business_id', businessId.toString());
        }
        queryParams.append('is_active', 'true');

        const url = `${env.API_BASE_URL}/integrations/simple${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;

        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error('Error al obtener integraciones');
        }

        return await response.json();
    } catch (error: any) {
        console.error('Get Integrations Simple Action Error:', error.message);
        return {
            success: false,
            message: error.message,
            data: [],
        };
    }
};

// ============================================
// Integration Categories
// ============================================

export const getIntegrationCategoriesAction = async (token?: string | null) => {
    return safeAction(() => getUseCases(token).then(uc => uc.getIntegrationCategories()));
};
