'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { CodReportApiRepository } from '../repository/api-repository';
import {
    ReportFilters,
    CodOrdersParams,
    SaveCarrierConfigInput,
} from '../../domain/types';

const getRepo = async () => {
    const token = await getAuthToken();
    return new CodReportApiRepository(token);
};

export const getCodSummaryAction = async (filters: ReportFilters) => {
    try {
        return await (await getRepo()).getSummary(filters);
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

export const confirmCodCutAction = async (periodStart: string, periodEnd: string, businessId?: number) => {
    try {
        return await (await getRepo()).confirmCut(periodStart, periodEnd, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al confirmar el corte', data: null as any };
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
