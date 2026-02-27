import {
    Warehouse,
    WarehouseDetail,
    WarehouseLocation,
    WarehousesListResponse,
    GetWarehousesParams,
    CreateWarehouseDTO,
    UpdateWarehouseDTO,
    CreateLocationDTO,
    UpdateLocationDTO,
} from './types';

export interface IWarehouseRepository {
    getWarehouses(params?: GetWarehousesParams): Promise<WarehousesListResponse>;
    getWarehouseById(id: number, businessId?: number): Promise<WarehouseDetail>;
    createWarehouse(data: CreateWarehouseDTO, businessId?: number): Promise<Warehouse>;
    updateWarehouse(id: number, data: UpdateWarehouseDTO, businessId?: number): Promise<Warehouse>;
    deleteWarehouse(id: number, businessId?: number): Promise<void>;
    // Locations
    getLocations(warehouseId: number, businessId?: number): Promise<WarehouseLocation[]>;
    createLocation(warehouseId: number, data: CreateLocationDTO, businessId?: number): Promise<WarehouseLocation>;
    updateLocation(warehouseId: number, locationId: number, data: UpdateLocationDTO, businessId?: number): Promise<WarehouseLocation>;
    deleteLocation(warehouseId: number, locationId: number, businessId?: number): Promise<void>;
}
