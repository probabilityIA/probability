'use server';

import { cookies } from 'next/headers';
import { ShippingMarginApiRepository } from '../repository/api-repository';
import { ShippingMarginUseCases } from '../../app/use-cases';
import { GetShippingMarginsParams, CreateShippingMarginDTO, UpdateShippingMarginDTO } from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new ShippingMarginApiRepository(token);
    return new ShippingMarginUseCases(repository);
}

export const listShippingMarginsAction = async (params?: GetShippingMarginsParams) => {
    try {
        return await (await getUseCases()).list(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getShippingMarginAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getById(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createShippingMarginAction = async (data: CreateShippingMarginDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).create(data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateShippingMarginAction = async (id: number, data: UpdateShippingMarginDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).update(id, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

