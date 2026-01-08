import { GetShipmentsParams, PaginatedResponse, Shipment, EnvioClickQuoteRequest, EnvioClickGenerateResponse, EnvioClickQuoteResponse } from './types';

export interface IShipmentRepository {
    getShipments(params?: GetShipmentsParams): Promise<PaginatedResponse<Shipment>>;
    quoteShipment(req: EnvioClickQuoteRequest): Promise<EnvioClickQuoteResponse>;
    generateGuide(req: EnvioClickQuoteRequest): Promise<EnvioClickGenerateResponse>;
}
