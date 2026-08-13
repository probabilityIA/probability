'use server';

import { cookies } from 'next/headers';
import { updateTag } from 'next/cache';
import { getAuthToken } from '@/shared/utils/server-auth';
import { env } from '@/shared/config/env';
import { ProductApiRepository } from '../repository/api-repository';
import { ProductUseCases } from '../../app/use-cases';
import {
    GetProductsParams,
    GetFamiliesParams,
    CreateProductDTO,
    UpdateProductDTO,
    CreateProductFamilyDTO,
    UpdateProductFamilyDTO,
    AddProductIntegrationDTO,
    UpdateProductIntegrationDTO
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new ProductApiRepository(token);
    return new ProductUseCases(repository);
}

export const getProductsAction = async (params?: GetProductsParams) => {
    try {
        return await (await getUseCases()).getProducts(params);
    } catch (error: any) {
        return {
            success: false,
            message: error.message || 'Error al obtener productos',
            data: [],
            total: 0,
            page: params?.page || 1,
            page_size: params?.page_size || 20,
            total_pages: 0
        };
    }
};

export const getProductByIdAction = async (id: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getProductById(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getProductDetailAction = async (id: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getProductDetail(id, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al obtener el detalle', data: null };
    }
};

export const createProductAction = async (data: CreateProductDTO, businessId?: number) => {
    try {
        const result = await (await getUseCases()).createProduct(data, businessId);
        updateTag('public-tienda');
        return result;
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al crear producto', data: null };
    }
};

export const updateProductAction = async (id: string, data: UpdateProductDTO, businessId?: number) => {
    try {
        const result = await (await getUseCases()).updateProduct(id, data, businessId);
        updateTag('public-tienda');
        return result;
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteProductAction = async (id: string, businessId?: number) => {
    try {
        const result = await (await getUseCases()).deleteProduct(id, businessId);
        updateTag('public-tienda');
        return result;
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const uploadProductImageAction = async (productId: string, formData: FormData, businessId?: number) => {
    try {
        const result = await (await getUseCases()).uploadProductImage(productId, formData, businessId);
        updateTag('public-tienda');
        return result;
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al subir imagen', image_url: '' };
    }
};

export const getProductFamiliesAction = async (params?: GetFamiliesParams) => {
    try {
        return await (await getUseCases()).getProductFamilies(params);
    } catch (error: any) {
        return {
            success: false,
            message: error.message || 'Error al obtener familias',
            data: [],
            total: 0,
            page: params?.page || 1,
            page_size: params?.page_size || 20,
            total_pages: 0
        };
    }
};

export const getProductFamilyByIdAction = async (familyId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getProductFamilyById(familyId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getFamilyVariantsAction = async (familyId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getFamilyVariants(familyId, businessId);
    } catch (error: any) {
        return { success: false, data: [] };
    }
};

export const createProductFamilyAction = async (data: CreateProductFamilyDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createProductFamily(data, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al crear familia', data: null };
    }
};

export const updateProductFamilyAction = async (familyId: number, data: UpdateProductFamilyDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateProductFamily(familyId, data, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al actualizar familia', data: null };
    }
};

export const deleteProductFamilyAction = async (familyId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteProductFamily(familyId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al eliminar familia', error: error.message };
    }
};

export const addProductIntegrationAction = async (
    productId: string,
    data: AddProductIntegrationDTO,
    businessId?: number
) => {
    try {
        return await (await getUseCases()).addProductIntegration(productId, data, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al asociar integración', data: null };
    }
};

export const updateProductIntegrationAction = async (
    productId: string,
    integrationId: number,
    data: UpdateProductIntegrationDTO,
    businessId?: number
) => {
    try {
        return await (await getUseCases()).updateProductIntegration(productId, integrationId, data, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al actualizar mapping', data: null };
    }
};

export const removeProductIntegrationAction = async (
    productId: string,
    integrationId: number,
    businessId?: number
) => {
    try {
        return await (await getUseCases()).removeProductIntegration(productId, integrationId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al remover integración', error: error.message };
    }
};

export const getProductIntegrationsAction = async (productId: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getProductIntegrations(productId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al obtener integraciones', data: [], total: 0 };
    }
};

export const getSKUsAction = async (prefix?: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getSKUs(prefix, businessId);
    } catch (error: any) {
        return { success: false, data: [] };
    }
};

export const getNextSKUAction = async (prefix: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getNextSKU(prefix, businessId);
    } catch (error: any) {
        return { success: false, data: '' };
    }
};

export const getNextSKUBatchAction = async (prefix: string, count: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getNextSKUBatch(prefix, count, businessId);
    } catch (error: any) {
        return { success: false, data: [] };
    }
};

export const exportProductDimensionsAction = async (businessId?: number) => {
    try {
        const token = await getAuthToken();
        const url = businessId
            ? `${env.API_BASE_URL}/products/dimensions/export?business_id=${businessId}`
            : `${env.API_BASE_URL}/products/dimensions/export`;

        const response = await fetch(url, {
            headers: { 'Authorization': `Bearer ${token}` },
            cache: 'no-store',
        });

        if (!response.ok) {
            const text = await response.text();
            return { success: false, message: text || 'Error al generar el archivo' };
        }

        const buf = Buffer.from(await response.arrayBuffer());
        const disp = response.headers.get('content-disposition') || '';
        const match = disp.match(/filename="([^"]+)"/);
        const filename = match ? match[1] : 'dimensiones-productos.xlsx';
        const contentType = response.headers.get('content-type') || 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet';
        const dataUrl = `data:${contentType};base64,${buf.toString('base64')}`;

        return { success: true, blob: dataUrl, filename };
    } catch (error: any) {
        console.error('Export Product Dimensions Error:', error.message);
        return { success: false, message: error.message || 'Error al generar el archivo' };
    }
};

export const importProductDimensionsAction = async (file: File, businessId?: number) => {
    try {
        const token = await getAuthToken();

        const formData = new FormData();
        formData.append('file', file);

        const url = businessId
            ? `${env.API_BASE_URL}/products/dimensions/import?business_id=${businessId}`
            : `${env.API_BASE_URL}/products/dimensions/import`;

        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${token}` },
            body: formData,
            cache: 'no-store',
        });

        const result = await response.json();

        if (response.ok && result.success) {
            return {
                success: true,
                message: result.message,
                data: result.data || { total_rows: 0, success_count: 0, failed_count: 0, errors: [] },
            };
        }
        return { success: false, message: result.message || 'Error al procesar el archivo', data: result.data };
    } catch (error: any) {
        console.error('Import Product Dimensions Error:', error.message);
        return { success: false, message: error.message || 'Error al cargar el archivo' };
    }
};

export const importFamilyVariantsPreviewAction = async (familyId: number, file: File, businessId?: number) => {
    try {
        const token = await getAuthToken();

        const formData = new FormData();
        formData.append('file', file);

        const url = businessId
            ? `${env.API_BASE_URL}/products/families/${familyId}/import-variants?business_id=${businessId}`
            : `${env.API_BASE_URL}/products/families/${familyId}/import-variants`;

        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${token}` },
            body: formData,
            cache: 'no-store',
        });

        const result = await response.json();

        if (response.ok && result.success) {
            return {
                success: true,
                message: result.message,
                data: result.data || { matched: [], not_found: [] },
            };
        }
        return { success: false, message: result.message || 'Error al procesar el archivo' };
    } catch (error: any) {
        console.error('Import Family Variants Error:', error.message);
        return { success: false, message: error.message || 'Error al cargar el archivo' };
    }
};
