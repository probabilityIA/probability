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
    async getOriginAddresses(businessId?: number) {
        return this.repository.getOriginAddresses(businessId);
    }

    async createOriginAddress(req: any, businessId?: number) {
        return this.repository.createOriginAddress(req, businessId);
    }

    async updateOriginAddress(id: number, req: any, businessId?: number) {
        return this.repository.updateOriginAddress(id, req, businessId);
    }

    async deleteOriginAddress(id: number, businessId?: number) {
        return this.repository.deleteOriginAddress(id, businessId);
    }
}
