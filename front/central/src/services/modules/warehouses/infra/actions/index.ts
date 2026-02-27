'use server';

import { cookies } from 'next/headers';
import { WarehouseApiRepository } from '../repository/api-repository';
import { WarehouseUseCases } from '../../app/use-cases';
import {
    GetWarehousesParams,
    CreateWarehouseDTO,
    UpdateWarehouseDTO,
    CreateLocationDTO,
    UpdateLocationDTO,
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new WarehouseApiRepository(token);
    return new WarehouseUseCases(repository);
}

export const getWarehousesAction = async (params?: GetWarehousesParams) => {
    try {
        return await (await getUseCases()).getWarehouses(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getWarehouseByIdAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getWarehouseById(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createWarehouseAction = async (data: CreateWarehouseDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createWarehouse(data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateWarehouseAction = async (id: number, data: UpdateWarehouseDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateWarehouse(id, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteWarehouseAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteWarehouse(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getLocationsAction = async (warehouseId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getLocations(warehouseId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createLocationAction = async (warehouseId: number, data: CreateLocationDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createLocation(warehouseId, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateLocationAction = async (warehouseId: number, locationId: number, data: UpdateLocationDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateLocation(warehouseId, locationId, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteLocationAction = async (warehouseId: number, locationId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteLocation(warehouseId, locationId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
