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

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new IntegrationApiRepository(token);
    return new IntegrationUseCases(repository);
}

export const getIntegrationsAction = async (params?: GetIntegrationsParams) => {
    try {
        return await (await getUseCases()).getIntegrations(params);
    } catch (error: any) {
        console.error('Get Integrations Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getIntegrationByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getIntegrationById(id);
    } catch (error: any) {
        console.error('Get Integration By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getIntegrationByTypeAction = async (type: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getIntegrationByType(type, businessId);
    } catch (error: any) {
        console.error('Get Integration By Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createIntegrationAction = async (data: CreateIntegrationDTO) => {
    try {
        return await (await getUseCases()).createIntegration(data);
    } catch (error: any) {
        console.error('Create Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateIntegrationAction = async (id: number, data: UpdateIntegrationDTO) => {
    try {
        return await (await getUseCases()).updateIntegration(id, data);
    } catch (error: any) {
        console.error('Update Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteIntegrationAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteIntegration(id);
    } catch (error: any) {
        console.error('Delete Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const testConnectionAction = async (id: number) => {
    try {
        return await (await getUseCases()).testConnection(id);
    } catch (error: any) {
        console.error('Test Connection Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const activateIntegrationAction = async (id: number) => {
    try {
        return await (await getUseCases()).activateIntegration(id);
    } catch (error: any) {
        console.error('Activate Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deactivateIntegrationAction = async (id: number) => {
    try {
        return await (await getUseCases()).deactivateIntegration(id);
    } catch (error: any) {
        console.error('Deactivate Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const setAsDefaultAction = async (id: number) => {
    try {
        return await (await getUseCases()).setAsDefault(id);
    } catch (error: any) {
        console.error('Set As Default Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const syncOrdersAction = async (id: number, params?: SyncOrdersParams) => {
    try {
        return await (await getUseCases()).syncOrders(id, params);
    } catch (error: any) {
        console.error('Sync Orders Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getSyncStatusAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getSyncStatus(id, businessId);
    } catch (error: any) {
        console.error('Get Sync Status Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const testIntegrationAction = async (id: number) => {
    try {
        return await (await getUseCases()).testIntegration(id);
    } catch (error: any) {
        console.error('Test Integration Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const testConnectionRawAction = async (typeCode: string, config: any, credentials: any) => {
    try {
        return await (await getUseCases()).testConnectionRaw(typeCode, config, credentials);
    } catch (error: any) {
        console.error('Test Connection Raw Action Error:', error.message);
        throw new Error(error.message);
    }
};

// Integration Types
export const getIntegrationTypesAction = async () => {
    try {
        return await (await getUseCases()).getIntegrationTypes();
    } catch (error: any) {
        console.error('Get Integration Types Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getActiveIntegrationTypesAction = async () => {
    try {
        return await (await getUseCases()).getActiveIntegrationTypes();
    } catch (error: any) {
        console.error('Get Active Integration Types Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getIntegrationTypeByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getIntegrationTypeById(id);
    } catch (error: any) {
        console.error('Get Integration Type By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getIntegrationTypeByCodeAction = async (code: string) => {
    try {
        return await (await getUseCases()).getIntegrationTypeByCode(code);
    } catch (error: any) {
        console.error('Get Integration Type By Code Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createIntegrationTypeAction = async (data: CreateIntegrationTypeDTO) => {
    try {
        return await (await getUseCases()).createIntegrationType(data);
    } catch (error: any) {
        console.error('Create Integration Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateIntegrationTypeAction = async (id: number, data: UpdateIntegrationTypeDTO) => {
    try {
        return await (await getUseCases()).updateIntegrationType(id, data);
    } catch (error: any) {
        console.error('Update Integration Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteIntegrationTypeAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteIntegrationType(id);
    } catch (error: any) {
        console.error('Delete Integration Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getWebhookUrlAction = async (id: number) => {
    try {
        return await (await getUseCases()).getWebhookUrl(id);
    } catch (error: any) {
        console.error('Get Webhook URL Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const listWebhooksAction = async (id: number) => {
    try {
        return await (await getUseCases()).listWebhooks(id);
    } catch (error: any) {
        console.error('List Webhooks Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteWebhookAction = async (id: number, webhookId: string) => {
    try {
        return await (await getUseCases()).deleteWebhook(id, webhookId);
    } catch (error: any) {
        console.error('Delete Webhook Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const verifyWebhooksAction = async (id: number) => {
    try {
        return await (await getUseCases()).verifyWebhooks(id);
    } catch (error: any) {
        console.error('Verify Webhooks Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createWebhookAction = async (id: number) => {
    try {
        return await (await getUseCases()).createWebhook(id);
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
