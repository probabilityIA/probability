'use server';

import { cookies } from 'next/headers';
import { VehicleApiRepository } from '../repository/api-repository';
import { VehicleUseCases } from '../../app/use-cases';
import { GetVehiclesParams, CreateVehicleDTO, UpdateVehicleDTO } from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new VehicleApiRepository(token);
    return new VehicleUseCases(repository);
}

export const getVehiclesAction = async (params?: GetVehiclesParams) => {
    try {
        return await (await getUseCases()).getVehicles(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getVehicleByIdAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getVehicleById(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createVehicleAction = async (data: CreateVehicleDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createVehicle(data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateVehicleAction = async (id: number, data: UpdateVehicleDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateVehicle(id, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteVehicleAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteVehicle(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
