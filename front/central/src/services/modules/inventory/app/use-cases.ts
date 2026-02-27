import { IInventoryRepository } from '../domain/ports';
import {
    GetInventoryParams,
    GetMovementsParams,
    AdjustStockDTO,
    TransferStockDTO,
} from '../domain/types';

export class InventoryUseCases {
    constructor(private repository: IInventoryRepository) {}

    async getProductInventory(productId: string, businessId?: number) {
        return this.repository.getProductInventory(productId, businessId);
    }

    async getWarehouseInventory(warehouseId: number, params?: GetInventoryParams) {
        return this.repository.getWarehouseInventory(warehouseId, params);
    }

    async adjustStock(data: AdjustStockDTO, businessId?: number) {
        return this.repository.adjustStock(data, businessId);
    }

    async transferStock(data: TransferStockDTO, businessId?: number) {
        return this.repository.transferStock(data, businessId);
    }

    async getMovements(params?: GetMovementsParams) {
        return this.repository.getMovements(params);
    }

    async getMovementTypes(params?: { page?: number; page_size?: number; active_only?: boolean; business_id?: number }) {
        return this.repository.getMovementTypes(params);
    }
}
