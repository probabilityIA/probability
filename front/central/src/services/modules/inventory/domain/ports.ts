import {
    InventoryLevel,
    InventoryListResponse,
    MovementListResponse,
    MovementTypeListResponse,
    StockMovement,
    GetInventoryParams,
    GetMovementsParams,
    AdjustStockDTO,
    TransferStockDTO,
} from './types';

export interface IInventoryRepository {
    getProductInventory(productId: string, businessId?: number): Promise<InventoryLevel[]>;
    getWarehouseInventory(warehouseId: number, params?: GetInventoryParams): Promise<InventoryListResponse>;
    adjustStock(data: AdjustStockDTO, businessId?: number): Promise<StockMovement>;
    transferStock(data: TransferStockDTO, businessId?: number): Promise<{ message: string }>;
    getMovements(params?: GetMovementsParams): Promise<MovementListResponse>;
    getMovementTypes(params?: { page?: number; page_size?: number; active_only?: boolean; business_id?: number }): Promise<MovementTypeListResponse>;
}
