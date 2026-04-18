'use server';

import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { TraceabilityApiRepository } from '../repository/traceability-api-repository';
import {
    ChangeInventoryStateDTO,
    ConvertUoMInput,
    CreateLotDTO,
    CreateProductUoMDTO,
    CreateSerialDTO,
    GetLotsParams,
    GetSerialsParams,
    UpdateLotDTO,
    UpdateSerialDTO,
} from '../../domain/traceability-types';

async function getRepo() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    return new TraceabilityApiRepository(token);
}

function revalidateInventory() {
    revalidatePath('/inventory');
    revalidatePath('/inventory/lots');
    revalidatePath('/inventory/serials');
}

export const listLotsAction = async (params: GetLotsParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listLots(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createLotAction = async (data: CreateLotDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createLot(data, businessId);
        revalidateInventory();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear lote' };
    }
};

export const updateLotAction = async (id: number, data: UpdateLotDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.updateLot(id, data, businessId);
        revalidateInventory();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al actualizar lote' };
    }
};

export const deleteLotAction = async (id: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.deleteLot(id, businessId);
        revalidateInventory();
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al eliminar lote' };
    }
};

export const listSerialsAction = async (params: GetSerialsParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listSerials(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createSerialAction = async (data: CreateSerialDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createSerial(data, businessId);
        revalidateInventory();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear serie' };
    }
};

export const updateSerialAction = async (id: number, data: UpdateSerialDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.updateSerial(id, data, businessId);
        revalidateInventory();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al actualizar serie' };
    }
};

export const listInventoryStatesAction = async () => {
    try {
        const repo = await getRepo();
        return await repo.listStates();
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const changeInventoryStateAction = async (data: ChangeInventoryStateDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.changeState(data, businessId);
        revalidateInventory();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al cambiar estado' };
    }
};

export const listUoMsAction = async () => {
    try {
        const repo = await getRepo();
        return await repo.listUoMs();
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const listProductUoMsAction = async (productId: string, businessId?: number) => {
    try {
        const repo = await getRepo();
        return await repo.listProductUoMs(productId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createProductUoMAction = async (productId: string, data: CreateProductUoMDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createProductUoM(productId, data, businessId);
        revalidateInventory();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear UoM del producto' };
    }
};

export const deleteProductUoMAction = async (id: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.deleteProductUoM(id, businessId);
        revalidateInventory();
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al eliminar UoM' };
    }
};

export const convertUoMAction = async (data: ConvertUoMInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.convertUoM(data, businessId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al convertir UoM' };
    }
};
