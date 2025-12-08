import { IShipmentRepository } from '../domain/ports';
import { GetShipmentsParams, PaginatedResponse, Shipment } from '../domain/types';

export class ShipmentUseCases {
    private repository: IShipmentRepository;

    constructor(repository: IShipmentRepository) {
        this.repository = repository;
    }

    async getShipments(params?: GetShipmentsParams): Promise<PaginatedResponse<Shipment>> {
        return this.repository.getShipments(params);
    }
}
