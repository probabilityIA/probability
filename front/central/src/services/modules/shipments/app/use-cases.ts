import { IShipmentRepository } from '../domain/ports';
import { GetShipmentsParams, PaginatedResponse, Shipment, CreateShipmentRequest } from '../domain/types';

export class ShipmentUseCases {
    private repository: IShipmentRepository;

    constructor(repository: IShipmentRepository) {
        this.repository = repository;
    }

    async getShipments(params?: GetShipmentsParams): Promise<PaginatedResponse<Shipment>> {
        return this.repository.getShipments(params);
    }

    async quoteShipment(req: any) {
        return this.repository.quoteShipment(req);
    }

    async generateGuide(req: any) {
        return this.repository.generateGuide(req);
    }

    async trackShipment(trackingNumber: string) {
        return this.repository.trackShipment(trackingNumber);
    }

    async cancelShipment(id: string) {
        return this.repository.cancelShipment(id);
    }

    async createShipment(req: CreateShipmentRequest) {
        return this.repository.createShipment(req);
    }

    // Origin Addresses
    async getOriginAddresses() {
        return this.repository.getOriginAddresses();
    }

    async createOriginAddress(req: any) {
        return this.repository.createOriginAddress(req);
    }

    async updateOriginAddress(id: number, req: any) {
        return this.repository.updateOriginAddress(id, req);
    }

    async deleteOriginAddress(id: number) {
        return this.repository.deleteOriginAddress(id);
    }
}
