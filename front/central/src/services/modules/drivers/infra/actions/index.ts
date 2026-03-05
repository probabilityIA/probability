'use server';

import { cookies } from 'next/headers';
import { DriverApiRepository } from '../repository/api-repository';
import { DriverUseCases } from '../../app/use-cases';
import { GetDriversParams, CreateDriverDTO, UpdateDriverDTO } from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new DriverApiRepository(token);
    return new DriverUseCases(repository);
}

export const getDriversAction = async (params?: GetDriversParams) => {
    try {
        return await (await getUseCases()).getDrivers(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getDriverByIdAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getDriverById(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createDriverAction = async (data: CreateDriverDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createDriver(data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateDriverAction = async (id: number, data: UpdateDriverDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateDriver(id, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteDriverAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteDriver(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
