import { ShippingConfigOverview, ShippingConfig, SaveShippingConfigRequest } from './types';

export interface IShippingConfigRepository {
    getOverview(businessId?: number): Promise<ShippingConfigOverview>;
    saveBusinessConfig(req: SaveShippingConfigRequest, businessId?: number): Promise<ShippingConfig>;
    saveWarehouseConfig(warehouseId: number, req: SaveShippingConfigRequest, businessId?: number): Promise<ShippingConfig>;
    deleteWarehouseConfig(warehouseId: number, businessId?: number): Promise<void>;
    setDefaultWarehouse(warehouseId: number, businessId?: number): Promise<void>;
}
