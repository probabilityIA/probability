'use server';

import { cookies } from 'next/headers';
import { InventoryApiRepository } from '../repository/api-repository';
import { InventoryUseCases } from '../../app/use-cases';
import {
    GetInventoryParams,
    GetMovementsParams,
    AdjustStockDTO,
    TransferStockDTO,
    BulkLoadDTO,
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new InventoryApiRepository(token);
    return new InventoryUseCases(repository);
}

export const getProductInventoryAction = async (productId: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getProductInventory(productId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getWarehouseInventoryAction = async (warehouseId: number, params?: GetInventoryParams) => {
    try {
        return await (await getUseCases()).getWarehouseInventory(warehouseId, params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const adjustStockAction = async (data: AdjustStockDTO, businessId?: number) => {
    try {
        const result = await (await getUseCases()).adjustStock(data, businessId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al ajustar stock' };
    }
};

export const transferStockAction = async (data: TransferStockDTO, businessId?: number) => {
    try {
        const result = await (await getUseCases()).transferStock(data, businessId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al transferir stock' };
    }
};

export const bulkLoadInventoryAction = async (data: BulkLoadDTO, businessId?: number) => {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || null;
        const repository = new InventoryApiRepository(token);
        const result = await repository.bulkLoadInventory(data, businessId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error en carga masiva' };
    }
};

export const getMovementsAction = async (params?: GetMovementsParams) => {
    try {
        return await (await getUseCases()).getMovements(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getMovementTypesAction = async (params?: { page?: number; page_size?: number; active_only?: boolean; business_id?: number }) => {
    try {
        return await (await getUseCases()).getMovementTypes(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
