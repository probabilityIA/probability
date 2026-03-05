import { IRouteRepository } from '../domain/ports';
import {
    GetRoutesParams,
    CreateRouteDTO,
    UpdateRouteDTO,
    AddStopDTO,
    UpdateStopDTO,
    UpdateStopStatusDTO,
    ReorderStopsDTO,
} from '../domain/types';

export class RouteUseCases {
    constructor(private repository: IRouteRepository) {}

    // Route CRUD
    async getRoutes(params?: GetRoutesParams) {
        return this.repository.getRoutes(params);
    }

    async getRouteById(id: number, businessId?: number) {
        return this.repository.getRouteById(id, businessId);
    }

    async createRoute(data: CreateRouteDTO, businessId?: number) {
        return this.repository.createRoute(data, businessId);
    }

    async updateRoute(id: number, data: UpdateRouteDTO, businessId?: number) {
        return this.repository.updateRoute(id, data, businessId);
    }

    async deleteRoute(id: number, businessId?: number) {
        return this.repository.deleteRoute(id, businessId);
    }

    // Route lifecycle
    async startRoute(id: number, businessId?: number) {
        return this.repository.startRoute(id, businessId);
    }

    async completeRoute(id: number, businessId?: number) {
        return this.repository.completeRoute(id, businessId);
    }

    // Stop management
    async addStop(routeId: number, data: AddStopDTO, businessId?: number) {
        return this.repository.addStop(routeId, data, businessId);
    }

    async updateStop(routeId: number, stopId: number, data: UpdateStopDTO, businessId?: number) {
        return this.repository.updateStop(routeId, stopId, data, businessId);
    }

    async deleteStop(routeId: number, stopId: number, businessId?: number) {
        return this.repository.deleteStop(routeId, stopId, businessId);
    }

    async updateStopStatus(routeId: number, stopId: number, data: UpdateStopStatusDTO, businessId?: number) {
        return this.repository.updateStopStatus(routeId, stopId, data, businessId);
    }

    async reorderStops(routeId: number, data: ReorderStopsDTO, businessId?: number) {
        return this.repository.reorderStops(routeId, data, businessId);
    }
}
