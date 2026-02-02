'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { BusinessApiRepository } from '../repository/api-repository';
import { BusinessUseCases } from '../../app/use-cases';
import {
    GetBusinessesParams,
    CreateBusinessDTO,
    UpdateBusinessDTO,
    GetConfiguredResourcesParams,
    CreateBusinessTypeDTO,
    UpdateBusinessTypeDTO,
    BusinessesSimpleResponse
} from '../../domain/types';
import { env } from '@/shared/config/env';

async function getUseCases() {
    const token = await getAuthToken();
    const repository = new BusinessApiRepository(token);
    return new BusinessUseCases(repository);
}

// Business Actions
export const getBusinessesAction = async (params?: GetBusinessesParams) => {
    try {
        return (await getUseCases()).getBusinesses(params);
    } catch (error: any) {
        console.error('Get Businesses Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getBusinessByIdAction = async (id: number) => {
    try {
        return (await getUseCases()).getBusinessById(id);
    } catch (error: any) {
        console.error('Get Business By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createBusinessAction = async (data: CreateBusinessDTO) => {
    try {
        return (await getUseCases()).createBusiness(data);
    } catch (error: any) {
        console.error('Create Business Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateBusinessAction = async (id: number, data: UpdateBusinessDTO) => {
    try {
        return (await getUseCases()).updateBusiness(id, data);
    } catch (error: any) {
        console.error('Update Business Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteBusinessAction = async (id: number) => {
    try {
        return (await getUseCases()).deleteBusiness(id);
    } catch (error: any) {
        console.error('Delete Business Action Error:', error.message);
        throw new Error(error.message);
    }
};

// Configured Resources Actions
export const getConfiguredResourcesAction = async (params?: GetConfiguredResourcesParams) => {
    try {
        return (await getUseCases()).getConfiguredResources(params);
    } catch (error: any) {
        console.error('Get Configured Resources Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getBusinessConfiguredResourcesAction = async (id: number) => {
    try {
        return (await getUseCases()).getBusinessConfiguredResources(id);
    } catch (error: any) {
        console.error('Get Business Configured Resources Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const activateResourceAction = async (resourceId: number, businessId?: number) => {
    try {
        return (await getUseCases()).activateResource(resourceId, businessId);
    } catch (error: any) {
        console.error('Activate Resource Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deactivateResourceAction = async (resourceId: number, businessId?: number) => {
    try {
        return (await getUseCases()).deactivateResource(resourceId, businessId);
    } catch (error: any) {
        console.error('Deactivate Resource Action Error:', error.message);
        throw new Error(error.message);
    }
};

// Business Types Actions
export const getBusinessTypesAction = async () => {
    try {
        return (await getUseCases()).getBusinessTypes();
    } catch (error: any) {
        console.error('Get Business Types Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getBusinessTypeByIdAction = async (id: number) => {
    try {
        return (await getUseCases()).getBusinessTypeById(id);
    } catch (error: any) {
        console.error('Get Business Type By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createBusinessTypeAction = async (data: CreateBusinessTypeDTO) => {
    try {
        return (await getUseCases()).createBusinessType(data);
    } catch (error: any) {
        console.error('Create Business Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateBusinessTypeAction = async (id: number, data: UpdateBusinessTypeDTO) => {
    try {
        return (await getUseCases()).updateBusinessType(id, data);
    } catch (error: any) {
        console.error('Update Business Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteBusinessTypeAction = async (id: number) => {
    try {
        return (await getUseCases()).deleteBusinessType(id);
    } catch (error: any) {
        console.error('Delete Business Type Action Error:', error.message);
        throw new Error(error.message);
    }
};

// ============================================
// Simple Actions - Para Dropdowns/Selectores
// ============================================

export const getBusinessesSimpleAction = async (): Promise<BusinessesSimpleResponse> => {
    try {
        const token = await getAuthToken();

        const response = await fetch(`${env.API_BASE_URL}/businesses/simple`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error('Error al obtener negocios');
        }

        return await response.json();
    } catch (error: any) {
        console.error('Get Businesses Simple Action Error:', error.message);
        return {
            success: false,
            message: error.message,
            data: [],
        };
    }
};
