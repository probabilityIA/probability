'use server';

import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { OperationsApiRepository } from '../repository/operations-api-repository';
import {
    AssignReplenishmentInput,
    CompleteReplenishmentInput,
    ConfirmPutawayInput,
    CreateCrossDockLinkDTO,
    CreatePutawayRuleDTO,
    CreateReplenishmentTaskDTO,
    GetCrossDockLinksParams,
    GetPutawayRulesParams,
    GetPutawaySuggestionsParams,
    GetReplenishmentTasksParams,
    GetVelocitiesParams,
    RunSlottingInput,
    SuggestPutawayInput,
    UpdatePutawayRuleDTO,
} from '../../domain/operations-types';

async function getRepo() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    return new OperationsApiRepository(token);
}

function revalidateOperations() {
    revalidatePath('/inventory');
    revalidatePath('/inventory/operations');
    revalidatePath('/inventory/operations/putaway');
    revalidatePath('/inventory/operations/replenishment');
    revalidatePath('/inventory/operations/cross-dock');
    revalidatePath('/inventory/analytics/slotting');
}

export const listPutawayRulesAction = async (params: GetPutawayRulesParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listPutawayRules(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createPutawayRuleAction = async (data: CreatePutawayRuleDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createPutawayRule(data, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear regla' };
    }
};

export const updatePutawayRuleAction = async (id: number, data: UpdatePutawayRuleDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.updatePutawayRule(id, data, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al actualizar regla' };
    }
};

export const deletePutawayRuleAction = async (id: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.deletePutawayRule(id, businessId);
        revalidateOperations();
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al eliminar regla' };
    }
};

export const suggestPutawayAction = async (data: SuggestPutawayInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.suggestPutaway(data, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al sugerir put-away' };
    }
};

export const confirmPutawayAction = async (id: number, data: ConfirmPutawayInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.confirmPutaway(id, data, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al confirmar put-away' };
    }
};

export const listPutawaySuggestionsAction = async (params: GetPutawaySuggestionsParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listPutawaySuggestions(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const listReplenishmentTasksAction = async (params: GetReplenishmentTasksParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listReplenishmentTasks(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createReplenishmentTaskAction = async (data: CreateReplenishmentTaskDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createReplenishmentTask(data, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear tarea de reabastecimiento' };
    }
};

export const assignReplenishmentAction = async (id: number, data: AssignReplenishmentInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.assignReplenishment(id, data, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al asignar tarea' };
    }
};

export const completeReplenishmentAction = async (id: number, data: CompleteReplenishmentInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.completeReplenishment(id, data, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al completar tarea' };
    }
};

export const cancelReplenishmentAction = async (id: number, reason: string, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.cancelReplenishment(id, reason, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al cancelar tarea' };
    }
};

export const detectReplenishmentAction = async (businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.detectReplenishment(businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al detectar reabastecimiento' };
    }
};

export const listCrossDockLinksAction = async (params: GetCrossDockLinksParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listCrossDockLinks(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createCrossDockLinkAction = async (data: CreateCrossDockLinkDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createCrossDockLink(data, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear cross-dock' };
    }
};

export const executeCrossDockAction = async (id: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.executeCrossDock(id, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al ejecutar cross-dock' };
    }
};

export const runSlottingAction = async (data: RunSlottingInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.runSlotting(data, businessId);
        revalidateOperations();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error en analisis de slotting' };
    }
};

export const listVelocitiesAction = async (params: GetVelocitiesParams) => {
    try {
        const repo = await getRepo();
        return await repo.listVelocities(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
