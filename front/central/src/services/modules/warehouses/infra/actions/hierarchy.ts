'use server';

import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { HierarchyApiRepository } from '../repository/hierarchy-api-repository';
import {
    CreateAisleDTO,
    CreateRackDTO,
    CreateRackLevelDTO,
    CreateZoneDTO,
    UpdateAisleDTO,
    UpdateRackDTO,
    UpdateRackLevelDTO,
    UpdateZoneDTO,
    ValidateCubingInput,
} from '../../domain/hierarchy-types';

async function getRepo() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    return new HierarchyApiRepository(token);
}

function revalidateTree(warehouseId: number) {
    revalidatePath(`/inventory/warehouses`);
    revalidatePath(`/inventory/warehouses/${warehouseId}`);
}

export const getWarehouseTreeAction = async (warehouseId: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        return await repo.getTree(warehouseId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const listZonesAction = async (warehouseId: number, params: { page?: number; page_size?: number; active_only?: boolean } = {}, businessId?: number) => {
    try {
        const repo = await getRepo();
        return await repo.listZones(warehouseId, params, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createZoneAction = async (data: CreateZoneDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createZone(data, businessId);
        revalidateTree(data.warehouse_id);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear zona' };
    }
};

export const updateZoneAction = async (zoneId: number, warehouseId: number, data: UpdateZoneDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.updateZone(zoneId, data, businessId);
        revalidateTree(warehouseId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al actualizar zona' };
    }
};

export const deleteZoneAction = async (zoneId: number, warehouseId: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.deleteZone(zoneId, businessId);
        revalidateTree(warehouseId);
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al eliminar zona' };
    }
};

export const createAisleAction = async (data: CreateAisleDTO, warehouseId: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createAisle(data, businessId);
        revalidateTree(warehouseId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear pasillo' };
    }
};

export const updateAisleAction = async (aisleId: number, warehouseId: number, data: UpdateAisleDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.updateAisle(aisleId, data, businessId);
        revalidateTree(warehouseId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al actualizar pasillo' };
    }
};

export const deleteAisleAction = async (aisleId: number, warehouseId: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.deleteAisle(aisleId, businessId);
        revalidateTree(warehouseId);
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al eliminar pasillo' };
    }
};

export const createRackAction = async (data: CreateRackDTO, warehouseId: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createRack(data, businessId);
        revalidateTree(warehouseId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear rack' };
    }
};

export const updateRackAction = async (rackId: number, warehouseId: number, data: UpdateRackDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.updateRack(rackId, data, businessId);
        revalidateTree(warehouseId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al actualizar rack' };
    }
};

export const deleteRackAction = async (rackId: number, warehouseId: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.deleteRack(rackId, businessId);
        revalidateTree(warehouseId);
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al eliminar rack' };
    }
};

export const createRackLevelAction = async (data: CreateRackLevelDTO, warehouseId: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createRackLevel(data, businessId);
        revalidateTree(warehouseId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear nivel' };
    }
};

export const updateRackLevelAction = async (levelId: number, warehouseId: number, data: UpdateRackLevelDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.updateRackLevel(levelId, data, businessId);
        revalidateTree(warehouseId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al actualizar nivel' };
    }
};

export const deleteRackLevelAction = async (levelId: number, warehouseId: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.deleteRackLevel(levelId, businessId);
        revalidateTree(warehouseId);
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al eliminar nivel' };
    }
};

export const validateCubingAction = async (input: ValidateCubingInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.validateCubing(input, businessId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al validar cubicaje' };
    }
};
