'use server';

import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { AuditApiRepository } from '../repository/audit-api-repository';
import {
    ApproveDiscrepancyInput,
    CreateCycleCountPlanDTO,
    GenerateCountTaskInput,
    GetCountLinesParams,
    GetCountPlansParams,
    GetCountTasksParams,
    GetDiscrepanciesParams,
    KardexQueryInput,
    RejectDiscrepancyInput,
    StartCountTaskInput,
    SubmitCountLineInput,
    UpdateCycleCountPlanDTO,
} from '../../domain/audit-types';

async function getRepo() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    return new AuditApiRepository(token);
}

function revalidateAudit() {
    revalidatePath('/inventory/audit');
    revalidatePath('/inventory/audit/plans');
    revalidatePath('/inventory/audit/tasks');
    revalidatePath('/inventory/audit/discrepancies');
    revalidatePath('/inventory/kardex');
    revalidatePath('/inventory');
}

export const listCountPlansAction = async (params: GetCountPlansParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listCountPlans(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createCountPlanAction = async (data: CreateCycleCountPlanDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.createCountPlan(data, businessId);
        revalidateAudit();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al crear plan de conteo' };
    }
};

export const updateCountPlanAction = async (id: number, data: UpdateCycleCountPlanDTO, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.updateCountPlan(id, data, businessId);
        revalidateAudit();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al actualizar plan de conteo' };
    }
};

export const deleteCountPlanAction = async (id: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        await repo.deleteCountPlan(id, businessId);
        revalidateAudit();
        return { success: true as const };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al eliminar plan de conteo' };
    }
};

export const listCountTasksAction = async (params: GetCountTasksParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listCountTasks(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const generateCountTaskAction = async (data: GenerateCountTaskInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.generateCountTask(data, businessId);
        revalidateAudit();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al generar tarea de conteo' };
    }
};

export const startCountTaskAction = async (id: number, data: StartCountTaskInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.startCountTask(id, data, businessId);
        revalidateAudit();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al iniciar tarea' };
    }
};

export const finishCountTaskAction = async (id: number, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.finishCountTask(id, businessId);
        revalidateAudit();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al finalizar tarea' };
    }
};

export const listCountLinesAction = async (params: GetCountLinesParams) => {
    try {
        const repo = await getRepo();
        return await repo.listCountLines(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const submitCountLineAction = async (id: number, data: SubmitCountLineInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.submitCountLine(id, data, businessId);
        revalidateAudit();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al enviar conteo' };
    }
};

export const listDiscrepanciesAction = async (params: GetDiscrepanciesParams = {}) => {
    try {
        const repo = await getRepo();
        return await repo.listDiscrepancies(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const approveDiscrepancyAction = async (id: number, data: ApproveDiscrepancyInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.approveDiscrepancy(id, data, businessId);
        revalidateAudit();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al aprobar discrepancia' };
    }
};

export const rejectDiscrepancyAction = async (id: number, data: RejectDiscrepancyInput, businessId?: number) => {
    try {
        const repo = await getRepo();
        const result = await repo.rejectDiscrepancy(id, data, businessId);
        revalidateAudit();
        return { success: true as const, data: result };
    } catch (error: any) {
        return { success: false as const, error: error.message || 'Error al rechazar discrepancia' };
    }
};

export const exportKardexAction = async (data: KardexQueryInput) => {
    try {
        const repo = await getRepo();
        return await repo.exportKardex(data);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
