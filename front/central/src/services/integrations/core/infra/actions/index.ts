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
    try {
        return await (await getUseCases(token)).getIntegrationById(id);
    } catch (error: any) {
        console.error('Get Integration By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getIntegrationByTypeAction = async (type: string, businessId?: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getIntegrationByType(type, businessId);
    } catch (error: any) {
        console.error('Get Integration By Type Action Error:', error.message);
        throw new Error(error.message);
    }
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
    try {
        return await (await getUseCases(token)).updateIntegration(id, data);
    } catch (error: any) {
        console.error('Update Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteIntegrationAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).deleteIntegration(id);
    } catch (error: any) {
        console.error('Delete Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const testConnectionAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).testConnection(id);
    } catch (error: any) {
        console.error('Test Connection Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const activateIntegrationAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).activateIntegration(id);
    } catch (error: any) {
        console.error('Activate Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deactivateIntegrationAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).deactivateIntegration(id);
    } catch (error: any) {
        console.error('Deactivate Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const setAsDefaultAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).setAsDefault(id);
    } catch (error: any) {
        console.error('Set As Default Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const syncOrdersAction = async (id: number, params?: SyncOrdersParams, token?: string | null) => {
    try {
        return await (await getUseCases(token)).syncOrders(id, params);
    } catch (error: any) {
        console.error('Sync Orders Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getSyncStatusAction = async (id: number, businessId?: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getSyncStatus(id, businessId);
    } catch (error: any) {
        console.error('Get Sync Status Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const testIntegrationAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).testIntegration(id);
    } catch (error: any) {
        console.error('Test Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const testConnectionRawAction = async (typeCode: string, config: any, credentials: any, token?: string | null) => {
    try {
        return await (await getUseCases(token)).testConnectionRaw(typeCode, config, credentials);
    } catch (error: any) {
        console.error('Test Connection Raw Action Error:', error.message);
        throw new Error(error.message);
    }
};

// Integration Types
export const getIntegrationTypesAction = async (token?: string | null) => {
    try {
        return await (await getUseCases(token)).getIntegrationTypes();
    } catch (error: any) {
        console.error('Get Integration Types Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getActiveIntegrationTypesAction = async (token?: string | null) => {
    try {
        return await (await getUseCases(token)).getActiveIntegrationTypes();
    } catch (error: any) {
        console.error('Get Active Integration Types Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getIntegrationTypeByIdAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getIntegrationTypeById(id);
    } catch (error: any) {
        console.error('Get Integration Type By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getIntegrationTypeByCodeAction = async (code: string, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getIntegrationTypeByCode(code);
    } catch (error: any) {
        console.error('Get Integration Type By Code Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createIntegrationTypeAction = async (data: CreateIntegrationTypeDTO, token?: string | null) => {
    try {
        return await (await getUseCases(token)).createIntegrationType(data);
    } catch (error: any) {
        console.error('Create Integration Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateIntegrationTypeAction = async (id: number, data: UpdateIntegrationTypeDTO, token?: string | null) => {
    try {
        return await (await getUseCases(token)).updateIntegrationType(id, data);
    } catch (error: any) {
        console.error('Update Integration Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteIntegrationTypeAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).deleteIntegrationType(id);
    } catch (error: any) {
        console.error('Delete Integration Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getIntegrationTypePlatformCredentialsAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getIntegrationTypePlatformCredentials(id);
    } catch (error: any) {
        console.error('Get Integration Type Platform Credentials Action Error:', error.message);
        return { success: false, message: error.message, data: {} as Record<string, string> };
    }
};

export const getWebhookUrlAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getWebhookUrl(id);
    } catch (error: any) {
        console.error('Get Webhook URL Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const listWebhooksAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).listWebhooks(id);
    } catch (error: any) {
        console.error('List Webhooks Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteWebhookAction = async (id: number, webhookId: string, token?: string | null) => {
    try {
        return await (await getUseCases(token)).deleteWebhook(id, webhookId);
    } catch (error: any) {
        console.error('Delete Webhook Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const verifyWebhooksAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).verifyWebhooks(id);
    } catch (error: any) {
        console.error('Verify Webhooks Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createWebhookAction = async (id: number, token?: string | null) => {
    try {
        return await (await getUseCases(token)).createWebhook(id);
    } catch (error: any) {
        console.error('Create Webhook Action Error:', error.message);
        throw new Error(error.message);
    }
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
    try {
        return await (await getUseCases(token)).getIntegrationCategories();
    } catch (error: any) {
        console.error('Get Integration Categories Action Error:', error.message);
        throw new Error(error.message);
    }
};
