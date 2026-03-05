import {
    RouteInfo,
    RouteDetail,
    RouteStopInfo,
    RoutesListResponse,
    GetRoutesParams,
    CreateRouteDTO,
    UpdateRouteDTO,
    AddStopDTO,
    UpdateStopDTO,
    UpdateStopStatusDTO,
    ReorderStopsDTO,
    DeleteRouteResponse,
} from './types';

export interface IRouteRepository {
    // Route CRUD
    getRoutes(params?: GetRoutesParams): Promise<RoutesListResponse>;
    getRouteById(id: number, businessId?: number): Promise<RouteDetail>;
    createRoute(data: CreateRouteDTO, businessId?: number): Promise<RouteInfo>;
    updateRoute(id: number, data: UpdateRouteDTO, businessId?: number): Promise<RouteInfo>;
    deleteRoute(id: number, businessId?: number): Promise<DeleteRouteResponse>;

    // Route lifecycle
    startRoute(id: number, businessId?: number): Promise<RouteDetail>;
    completeRoute(id: number, businessId?: number): Promise<RouteDetail>;

    // Stop management
    addStop(routeId: number, data: AddStopDTO, businessId?: number): Promise<RouteStopInfo>;
    updateStop(routeId: number, stopId: number, data: UpdateStopDTO, businessId?: number): Promise<RouteStopInfo>;
    deleteStop(routeId: number, stopId: number, businessId?: number): Promise<DeleteRouteResponse>;
    updateStopStatus(routeId: number, stopId: number, data: UpdateStopStatusDTO, businessId?: number): Promise<RouteStopInfo>;
    reorderStops(routeId: number, data: ReorderStopsDTO, businessId?: number): Promise<RouteDetail>;
}
