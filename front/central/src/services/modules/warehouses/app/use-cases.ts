import { IWarehouseRepository } from '../domain/ports';
import {
    GetWarehousesParams,
    CreateWarehouseDTO,
    UpdateWarehouseDTO,
    CreateLocationDTO,
    UpdateLocationDTO,
} from '../domain/types';

export class WarehouseUseCases {
    constructor(private repository: IWarehouseRepository) {}

    async getWarehouses(params?: GetWarehousesParams) {
        return this.repository.getWarehouses(params);
    }

    async getWarehouseById(id: number, businessId?: number) {
        return this.repository.getWarehouseById(id, businessId);
    }

    async createWarehouse(data: CreateWarehouseDTO, businessId?: number) {
        return this.repository.createWarehouse(data, businessId);
    }

    async updateWarehouse(id: number, data: UpdateWarehouseDTO, businessId?: number) {
        return this.repository.updateWarehouse(id, data, businessId);
    }

    async deleteWarehouse(id: number, businessId?: number) {
        return this.repository.deleteWarehouse(id, businessId);
    }

    async getLocations(warehouseId: number, businessId?: number) {
        return this.repository.getLocations(warehouseId, businessId);
    }

    async createLocation(warehouseId: number, data: CreateLocationDTO, businessId?: number) {
        return this.repository.createLocation(warehouseId, data, businessId);
    }

    async updateLocation(warehouseId: number, locationId: number, data: UpdateLocationDTO, businessId?: number) {
        return this.repository.updateLocation(warehouseId, locationId, data, businessId);
    }

    async deleteLocation(warehouseId: number, locationId: number, businessId?: number) {
        return this.repository.deleteLocation(warehouseId, locationId, businessId);
    }
}
