'use server';

import { cookies } from 'next/headers';
import { getAuthToken } from '@/shared/utils/server-auth';
import { CustomerApiRepository } from '../repository/api-repository';
import { CustomerUseCases } from '../../app/use-cases';
import { GetCustomersParams, PaginationParams, CreateCustomerDTO, UpdateCustomerDTO, BulkCustomerResult } from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new CustomerApiRepository(token);
    return new CustomerUseCases(repository);
}

export const getCustomersAction = async (params?: GetCustomersParams) => {
    try {
        return await (await getUseCases()).getCustomers(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getCustomerByIdAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getCustomerById(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createCustomerAction = async (data: CreateCustomerDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createCustomer(data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateCustomerAction = async (id: number, data: UpdateCustomerDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateCustomer(id, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteCustomerAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteCustomer(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getCustomerSummaryAction = async (customerId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getCustomerSummary(customerId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getCustomerAddressesAction = async (customerId: number, params?: PaginationParams) => {
    try {
        return await (await getUseCases()).getCustomerAddresses(customerId, params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getCustomerProductsAction = async (customerId: number, params?: PaginationParams) => {
    try {
        return await (await getUseCases()).getCustomerProducts(customerId, params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getCustomerOrderItemsAction = async (customerId: number, params?: PaginationParams) => {
    try {
        return await (await getUseCases()).getCustomerOrderItems(customerId, params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const uploadBulkCustomersAction = async (file: File, businessId?: number) => {
    try {
        const token = await getAuthToken();
        const { env } = await import('@/shared/config/env');

        const formData = new FormData();
        formData.append('file', file);

        const url = businessId
            ? `${env.API_BASE_URL}/customers/upload-bulk?business_id=${businessId}`
            : `${env.API_BASE_URL}/customers/upload-bulk`;

        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
            body: formData,
            cache: 'no-store',
        });

        const result = await response.json();

        if (response.ok && result.success) {
            return {
                success: true,
                data: result.data as BulkCustomerResult,
            };
        }
        return {
            success: false,
            message: result.message || 'Error al procesar el archivo',
        };
    } catch (error: any) {
        return {
            success: false,
            message: error.message || 'Error al cargar el archivo',
        };
    }
};

export const downloadCustomerTemplateAction = async () => {
    const headers = ['nombre', 'apellido', 'cedula', 'correo', 'telefono', 'direccion', 'ciudad', 'notas'];
    const exampleRows = [
        ['Juan', 'Perez', '1234567890', 'juan@example.com', '3001234567', 'Calle 1 # 2-3', 'Bogota', 'Cliente frecuente'],
        ['Maria', 'Gomez', '0987654321', 'maria@example.com', '3007654321', '', '', ''],
    ];
    return { success: true, data: { headers, exampleRows } };
};
