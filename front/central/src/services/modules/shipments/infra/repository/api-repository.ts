import { IShipmentRepository } from '../../domain/ports';
import { GetShipmentsParams, PaginatedResponse, Shipment, EnvioClickQuoteRequest, EnvioClickGenerateResponse, EnvioClickQuoteResponse } from '../../domain/types';
import { TokenStorage } from '@/shared/config';
import { envPublic } from '@/shared/config/env';

export class ShipmentApiRepository implements IShipmentRepository {
    private baseUrl: string;

    constructor() {
        this.baseUrl = envPublic.API_BASE_URL;
    }

    private async fetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
        const token = TokenStorage.getSessionToken();
        const headers = {
            'Content-Type': 'application/json',
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
            ...options.headers,
        };

        const response = await fetch(`${this.baseUrl}${endpoint}`, {
            ...options,
            headers,
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.message || data.error || 'An error occurred');
        }

        return data;
    }

    async getShipments(params?: GetShipmentsParams): Promise<PaginatedResponse<Shipment>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        return this.fetch<PaginatedResponse<Shipment>>(`/shipments?${searchParams.toString()}`);
    }

    async quoteShipment(req: EnvioClickQuoteRequest): Promise<EnvioClickQuoteResponse> {
        return this.fetch<EnvioClickQuoteResponse>('/shipments/quote', {
            method: 'POST',
            body: JSON.stringify(req),
        });
    }

    async generateGuide(req: EnvioClickQuoteRequest): Promise<EnvioClickGenerateResponse> {
        return this.fetch<EnvioClickGenerateResponse>('/shipments/generate', {
            method: 'POST',
            body: JSON.stringify(req),
        });
    }
}
