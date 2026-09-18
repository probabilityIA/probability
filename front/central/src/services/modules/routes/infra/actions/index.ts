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

export type ActionResult<T> = { success: true; data: T } | { success: false; error: string };

async function run<T>(fn: () => Promise<T>): Promise<ActionResult<T>> {
    try {
        return { success: true, data: await fn() };
    } catch (error: unknown) {
        const message = error instanceof Error ? error.message : String(error || '');
        return { success: false, error: message || 'Ocurri\u00f3 un error. Por favor intenta de nuevo.' };
    }
}

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new RouteApiRepository(token);
    return new RouteUseCases(repository);
}

export const getRoutesAction = async (params?: GetRoutesParams) =>
    run(async () => (await getUseCases()).getRoutes(params));

export const getRouteByIdAction = async (id: number, businessId?: number) =>
    run(async () => (await getUseCases()).getRouteById(id, businessId));

export const createRouteAction = async (data: CreateRouteDTO, businessId?: number) =>
    run(async () => (await getUseCases()).createRoute(data, businessId));

export const updateRouteAction = async (id: number, data: UpdateRouteDTO, businessId?: number) =>
    run(async () => (await getUseCases()).updateRoute(id, data, businessId));

export const deleteRouteAction = async (id: number, businessId?: number) =>
    run(async () => (await getUseCases()).deleteRoute(id, businessId));

export const startRouteAction = async (id: number, businessId?: number) =>
    run(async () => (await getUseCases()).startRoute(id, businessId));

export const optimizeRouteAction = async (id: number, businessId?: number) =>
    run(async () => (await getUseCases()).optimizeRoute(id, businessId));

export const completeRouteAction = async (id: number, businessId?: number) =>
    run(async () => (await getUseCases()).completeRoute(id, businessId));

export const addStopAction = async (routeId: number, data: AddStopDTO, businessId?: number) =>
    run(async () => (await getUseCases()).addStop(routeId, data, businessId));

export const updateStopAction = async (routeId: number, stopId: number, data: UpdateStopDTO, businessId?: number) =>
    run(async () => (await getUseCases()).updateStop(routeId, stopId, data, businessId));

export const deleteStopAction = async (routeId: number, stopId: number, businessId?: number) =>
    run(async () => (await getUseCases()).deleteStop(routeId, stopId, businessId));

export const updateStopStatusAction = async (routeId: number, stopId: number, data: UpdateStopStatusDTO, businessId?: number) =>
    run(async () => (await getUseCases()).updateStopStatus(routeId, stopId, data, businessId));

export const reorderStopsAction = async (routeId: number, data: ReorderStopsDTO, businessId?: number) =>
    run(async () => (await getUseCases()).reorderStops(routeId, data, businessId));

export const getAvailableDriversAction = async (businessId?: number) =>
    run(async () => (await getUseCases()).getAvailableDrivers(businessId));

export const getAvailableVehiclesAction = async (businessId?: number) =>
    run(async () => (await getUseCases()).getAvailableVehicles(businessId));

export const getAssignableOrdersAction = async (businessId?: number) =>
    run(async () => (await getUseCases()).getAssignableOrders(businessId));
