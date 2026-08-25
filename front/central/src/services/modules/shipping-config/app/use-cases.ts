import { IShippingConfigRepository } from '../domain/ports';
import { SaveShippingConfigRequest } from '../domain/types';

export class ShippingConfigUseCases {
    private repository: IShippingConfigRepository;

    constructor(repository: IShippingConfigRepository) {
        this.repository = repository;
    }

    async getOverview(businessId?: number) {
        return this.repository.getOverview(businessId);
    }

    async saveBusinessConfig(req: SaveShippingConfigRequest, businessId?: number) {
        return this.repository.saveBusinessConfig(req, businessId);
    }

    async saveWarehouseConfig(warehouseId: number, req: SaveShippingConfigRequest, businessId?: number) {
        return this.repository.saveWarehouseConfig(warehouseId, req, businessId);
    }

    async deleteWarehouseConfig(warehouseId: number, businessId?: number) {
        return this.repository.deleteWarehouseConfig(warehouseId, businessId);
    }

    async setDefaultWarehouse(warehouseId: number, businessId?: number) {
        return this.repository.setDefaultWarehouse(warehouseId, businessId);
    }
}
