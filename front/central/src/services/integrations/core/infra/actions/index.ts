'use server';

import { cookies } from 'next/headers';
import { IntegrationApiRepository } from '../repository/api-repository';
import { IntegrationUseCases } from '../../app/use-cases';
import {
    GetIntegrationsParams,
    CreateIntegrationDTO,
    UpdateIntegrationDTO
} from '../../domain/types';

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
