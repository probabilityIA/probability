import { GetShipmentsParams, PaginatedResponse, Shipment, EnvioClickQuoteRequest, EnvioClickGenerateResponse, EnvioClickQuoteResponse, EnvioClickTrackingResponse, EnvioClickCancelResponse } from './types';

export interface IShipmentRepository {
    getShipments(params?: GetShipmentsParams): Promise<PaginatedResponse<Shipment>>;
    quoteShipment(req: EnvioClickQuoteRequest): Promise<EnvioClickQuoteResponse>;
    generateGuide(req: EnvioClickQuoteRequest): Promise<EnvioClickGenerateResponse>;
    trackShipment(trackingNumber: string): Promise<EnvioClickTrackingResponse>;
    cancelShipment(id: string): Promise<EnvioClickCancelResponse>;
}
