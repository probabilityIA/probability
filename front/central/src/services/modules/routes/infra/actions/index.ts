'use server';

import { cookies } from 'next/headers';
import { RouteApiRepository } from '../repository/api-repository';
import { RouteUseCases } from '../../app/use-cases';
import {
    GetRoutesParams,
    CreateRouteDTO,
    UpdateRouteDTO,
    AddStopDTO,
    UpdateStopDTO,
    UpdateStopStatusDTO,
    ReorderStopsDTO,
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new RouteApiRepository(token);
    return new RouteUseCases(repository);
}

// ============================================
// Route CRUD
// ============================================

export const getRoutesAction = async (params?: GetRoutesParams) => {
    try {
        return await (await getUseCases()).getRoutes(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getRouteByIdAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getRouteById(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createRouteAction = async (data: CreateRouteDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createRoute(data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateRouteAction = async (id: number, data: UpdateRouteDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateRoute(id, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteRouteAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteRoute(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

// ============================================
// Route lifecycle
// ============================================

export const startRouteAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).startRoute(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const completeRouteAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).completeRoute(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

// ============================================
// Stop management
// ============================================

export const addStopAction = async (routeId: number, data: AddStopDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).addStop(routeId, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateStopAction = async (routeId: number, stopId: number, data: UpdateStopDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateStop(routeId, stopId, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteStopAction = async (routeId: number, stopId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteStop(routeId, stopId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateStopStatusAction = async (routeId: number, stopId: number, data: UpdateStopStatusDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateStopStatus(routeId, stopId, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const reorderStopsAction = async (routeId: number, data: ReorderStopsDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).reorderStops(routeId, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
