import { GetShipmentsParams, PaginatedResponse, Shipment, EnvioClickQuoteRequest, EnvioClickGenerateResponse, EnvioClickQuoteResponse, EnvioClickTrackingResponse, EnvioClickCancelResponse, CreateShipmentRequest, OriginAddress, CreateOriginAddressRequest, UpdateOriginAddressRequest } from './types';

export interface IShipmentRepository {
    getShipments(params?: GetShipmentsParams): Promise<PaginatedResponse<Shipment>>;
    quoteShipment(req: EnvioClickQuoteRequest): Promise<EnvioClickQuoteResponse>;
    generateGuide(req: EnvioClickQuoteRequest): Promise<EnvioClickGenerateResponse>;
    trackShipment(trackingNumber: string): Promise<EnvioClickTrackingResponse>;
    cancelShipment(id: string): Promise<EnvioClickCancelResponse>;
    createShipment(req: CreateShipmentRequest): Promise<{ success: boolean; message: string; data?: Shipment }>;

    // Direcciones de Origen
    getOriginAddresses(): Promise<OriginAddress[]>;
    createOriginAddress(req: CreateOriginAddressRequest): Promise<OriginAddress>;
    updateOriginAddress(id: number, req: UpdateOriginAddressRequest): Promise<OriginAddress>;
    deleteOriginAddress(id: number): Promise<{ message: string }>;
}
