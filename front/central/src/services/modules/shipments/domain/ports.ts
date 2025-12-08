import { GetShipmentsParams, PaginatedResponse, Shipment } from './types';

export interface IShipmentRepository {
    getShipments(params?: GetShipmentsParams): Promise<PaginatedResponse<Shipment>>;
}
