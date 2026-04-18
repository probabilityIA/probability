'use server';

import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { CaptureApiRepository } from '../repository/capture-api-repository';
import {
    AddToLPNInput,
    CreateLPNDTO,
    GetLPNsParams,
    GetSyncLogsParams,
    InboundSyncInput,
    MergeLPNInput,
    MoveLPNInput,
    ScanInput,
    UpdateLPNDTO,
} from '../../domain/capture-types';

async function getRepo() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    return new CaptureApiRepository(token);
}

function revalidateCapture() {
    revalidatePath('/inventory/lpn');
    revalidatePath('/inventory/mobile');
    revalidatePath('/inventory/sync/logs');
    revalidatePath('/inventory');
}

export const listLPNsAction = async (params: GetLPNsParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listLPNs(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createLPNAction = async (data: CreateLPNDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createLPN(data, businessId);
        revalidateCapture();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear LPN' };
    }
};

export const updateLPNAction = async (id: number, data: UpdateLPNDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.updateLPN(id, data, businessId);
        revalidateCapture();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al actualizar LPN' };
    }
};

export const deleteLPNAction = async (id: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.deleteLPN(id, businessId);
        revalidateCapture();
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al eliminar LPN' };
    }
};

export const addToLPNAction = async (id: number, data: AddToLPNInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.addToLPN(id, data, businessId);
        revalidateCapture();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al agregar a LPN' };
    }
};

export const moveLPNAction = async (id: number, data: MoveLPNInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.moveLPN(id, data, businessId);
        revalidateCapture();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al mover LPN' };
    }
};

export const dissolveLPNAction = async (id: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.dissolveLPN(id, businessId);
        revalidateCapture();
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al disolver LPN' };
    }
};

export const mergeLPNAction = async (id: number, data: MergeLPNInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.mergeLPN(id, data, businessId);
        revalidateCapture();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al mezclar LPN' };
    }
};

export const scanAction = async (data: ScanInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.scan(data, businessId);
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al escanear' };
    }
};

export const inboundSyncAction = async (integrationId: number, data: InboundSyncInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.inboundSync(integrationId, data, businessId);
        revalidateCapture();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error en sincronizacion entrante' };
    }
};

export const listSyncLogsAction = async (params: GetSyncLogsParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listSyncLogs(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
