'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { CodReportApiRepository } from '../repository/api-repository';
import {
    ReportFilters,
    CodOrdersParams,
    SaveCarrierConfigInput,
    UpdateCarrierFeeInput,
} from '../../domain/types';

const getRepo = async () => {
    const token = await getAuthToken();
    return new CodReportApiRepository(token);
};

export const getCodSummaryAction = async (filters: ReportFilters, bucket?: string) => {
    try {
        return await (await getRepo()).getSummary(filters, bucket);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al obtener el resumen', data: null as any };
    }
};

export const getCodOrdersAction = async (params: CodOrdersParams) => {
    try {
        return await (await getRepo()).getOrders(params);
    } catch (error: any) {
        return {
            success: false,
            message: error.message || 'Error al obtener las ordenes',
            data: [],
            total: 0,
            page: params.page || 1,
            page_size: params.page_size || 10,
            total_pages: 0,
        };
    }
};

export const getCodCutsAction = async (businessId?: number) => {
    try {
        return await (await getRepo()).getCuts(businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al obtener los cortes', data: [], can_confirm: false };
    }
};

export const getSelectableOrdersAction = async (periodStart: string, periodEnd: string, businessId?: number) => {
    try {
        return await (await getRepo()).getSelectableOrders(periodStart, periodEnd, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al obtener las ordenes de la semana', data: [] as any };
    }
};

export const getCutOrdersAction = async (cutId: number, businessId?: number) => {
    try {
        return await (await getRepo()).getCutOrders(cutId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al obtener las ordenes del corte', data: [] as any };
    }
};

export const deleteCodCutAction = async (cutId: number, businessId?: number) => {
    try {
        return await (await getRepo()).deleteCut(cutId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al eliminar el corte' };
    }
};

export const createDraftCutAction = async (periodStart: string, periodEnd: string, orderIds: string[], businessId?: number) => {
    try {
        return await (await getRepo()).createDraft(periodStart, periodEnd, orderIds, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al crear el borrador del corte', data: null as any };
    }
};

export const confirmCodCutAction = async (cutId: number, businessId?: number) => {
    try {
        return await (await getRepo()).confirmCut(cutId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al confirmar el corte' };
    }
};

export const getCodCutEmailHistoryAction = async (cutId: number, businessId?: number) => {
    try {
        return await (await getRepo()).getCutEmailHistory(cutId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al obtener el historico de correos', data: [] as any };
    }
};

export const getCodCutEmailPreviewAction = async (cutId: number, businessId?: number) => {
    try {
        return await (await getRepo()).getCutEmailPreview(cutId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al generar la vista previa', data: null as { subject: string; html: string } | null };
    }
};

export const sendCodCutEmailAction = async (cutId: number, emails: string[], businessId?: number) => {
    try {
        return await (await getRepo()).sendCutEmail(cutId, emails, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al enviar el corte por correo', data: undefined as { sent_to: string[] } | undefined };
    }
};

export const getCarrierConfigsAction = async (businessId?: number) => {
    try {
        return await (await getRepo()).getCarrierConfigs(businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al obtener la configuracion', data: [] };
    }
};

export const saveCarrierConfigAction = async (input: SaveCarrierConfigInput, businessId?: number) => {
    try {
        return await (await getRepo()).saveCarrierConfig(input, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al guardar la configuracion', data: null as any };
    }
};

export const updateCarrierFeeAction = async (input: UpdateCarrierFeeInput, businessId?: number) => {
    try {
        return await (await getRepo()).updateCarrierFee(input, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al actualizar la comision', data: null as any };
    }
};
