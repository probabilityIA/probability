'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { WebsiteConfigApiRepository } from '../repository/api-repository';
import { WebsiteConfigUseCases } from '../../app/use-cases';
import { UpdateWebsiteConfigDTO } from '../../domain/types';

async function getUseCases() {
    const token = await getAuthToken();
    const repository = new WebsiteConfigApiRepository(token);
    return new WebsiteConfigUseCases(repository);
}

export const getWebsiteConfigAction = async (businessId?: number) => {
    try {
        return await (await getUseCases()).getConfig(businessId);
    } catch (error: any) {
        console.error('Get Website Config Action Error:', error.message);
        return null;
    }
};

export const updateWebsiteConfigAction = async (data: UpdateWebsiteConfigDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateConfig(data, businessId);
    } catch (error: any) {
        console.error('Update Website Config Action Error:', error.message);
        return { success: false, message: error.message || 'Error al actualizar configuración' };
    }
};
