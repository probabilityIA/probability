'use server';

import { updateTag } from 'next/cache';
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
        const result = await (await getUseCases()).updateConfig(data, businessId);
        updateTag('public-tienda');
        return result;
    } catch (error: any) {
        console.error('Update Website Config Action Error:', error.message);
        return { success: false, message: error.message || 'Error al actualizar configuración' };
    }
};

export const getWebsiteCategoriesAction = async (businessId?: number) => {
    try {
        return await (await getUseCases()).getCategories(businessId);
    } catch (error: any) {
        console.error('Get Website Categories Action Error:', error.message);
        return [];
    }
};

export const deleteWebsiteImageAction = async (imageUrl: string, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteImage(imageUrl, businessId);
    } catch (error: any) {
        console.error('Delete Website Image Action Error:', error.message);
        return { success: false as const, message: error.message || 'Error al eliminar imagen' };
    }
};

export const uploadWebsiteImageAction = async (formData: FormData, businessId?: number) => {
    try {
        return await (await getUseCases()).uploadImage(formData, businessId);
    } catch (error: any) {
        console.error('Upload Website Image Action Error:', error.message);
        return { success: false as const, image_url: '', message: error.message || 'Error al subir imagen' };
    }
};
