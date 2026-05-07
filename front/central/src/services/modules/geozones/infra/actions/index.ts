'use server';

import { cookies } from 'next/headers';
import { GeozoneApiRepository } from '../repository/api-repository';
import { GeozoneUseCases } from '../../app/use-cases';
import {
    GetGeozonesParams,
    CreateGeozoneDTO,
    LookupParams,
    BulkImportRequest,
    GeozoneType,
    ProbabilityRequest,
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    return new GeozoneUseCases(new GeozoneApiRepository(token));
}

export const listGeozonesAction = async (params?: GetGeozonesParams) => {
    try { return await (await getUseCases()).list(params); }
    catch (error: any) { throw new Error(error.message); }
};

export const getGeozoneAction = async (id: number, includeGeom = true, businessId?: number) => {
    try { return await (await getUseCases()).getById(id, includeGeom, businessId); }
    catch (error: any) { throw new Error(error.message); }
};

export const createGeozoneAction = async (data: CreateGeozoneDTO, businessId?: number) => {
    try { return await (await getUseCases()).create(data, businessId); }
    catch (error: any) { throw new Error(error.message); }
};

export const bulkImportGeozonesAction = async (data: BulkImportRequest, businessId?: number) => {
    try { return await (await getUseCases()).bulkImport(data, businessId); }
    catch (error: any) { throw new Error(error.message); }
};

export const lookupGeozoneAction = async (params: LookupParams) => {
    try { return await (await getUseCases()).lookup(params); }
    catch (error: any) { throw new Error(error.message); }
};

export const deleteGeozoneAction = async (id: number, businessId?: number) => {
    try { await (await getUseCases()).remove(id, businessId); return { success: true }; }
    catch (error: any) { throw new Error(error.message); }
};

export const getGeozonesForDisplayAction = async (type: GeozoneType | '', zoom: number, bbox?: string) => {
    try { return await (await getUseCases()).getForDisplay(type, zoom, bbox); }
    catch (error: any) { throw new Error(error.message); }
};

export const getDeliveryProbabilityAction = async (req: ProbabilityRequest) => {
    try { return await (await getUseCases()).probability(req); }
    catch (error: any) { throw new Error(error.message); }
};
