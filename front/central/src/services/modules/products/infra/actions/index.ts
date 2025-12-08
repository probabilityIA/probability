'use server';

import { cookies } from 'next/headers';
import { ProductApiRepository } from '../repository/api-repository';
import { ProductUseCases } from '../../app/use-cases';
import {
    GetProductsParams,
    CreateProductDTO,
    UpdateProductDTO
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
        console.error('Get Products Action Error:', error.message);
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

export const getProductByIdAction = async (id: string) => {
    try {
        return await (await getUseCases()).getProductById(id);
    } catch (error: any) {
        console.error('Get Product By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createProductAction = async (data: CreateProductDTO) => {
    try {
        return await (await getUseCases()).createProduct(data);
    } catch (error: any) {
        console.error('Create Product Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateProductAction = async (id: string, data: UpdateProductDTO) => {
    try {
        return await (await getUseCases()).updateProduct(id, data);
    } catch (error: any) {
        console.error('Update Product Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteProductAction = async (id: string) => {
    try {
        return await (await getUseCases()).deleteProduct(id);
    } catch (error: any) {
        console.error('Delete Product Action Error:', error.message);
        throw new Error(error.message);
    }
};
