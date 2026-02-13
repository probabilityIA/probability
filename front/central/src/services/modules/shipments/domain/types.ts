export interface Shipment {
    id: number;
    created_at: string;
    updated_at: string;
    order_id: string;
    tracking_number?: string;
    tracking_url?: string;
    carrier?: string;
    carrier_code?: string;
    guide_id?: string;
    guide_url?: string;
    status: 'pending' | 'in_transit' | 'delivered' | 'failed';
    shipped_at?: string;
    delivered_at?: string;
    shipping_cost?: number;
    insurance_cost?: number;
    total_cost?: number;
    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    warehouse_name?: string;
    driver_name?: string;
    is_last_mile: boolean;
    estimated_delivery?: string;
    delivery_notes?: string;
}

export interface GetShipmentsParams {
    page?: number;
    page_size?: number;
    order_id?: string;
    tracking_number?: string;
    carrier?: string;
    status?: string;
    start_date?: string;
    end_date?: string;
    shipped_after?: string;
    shipped_before?: string;
    delivered_after?: string;
    delivered_before?: string;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
}

export interface PaginatedResponse<T> {
    success: boolean;
    message: string;
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface EnvioClickAddress {
    company: string;
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    address: string;
    suburb: string;
    crossStreet: string;
    reference: string;
    daneCode: string;
}

export interface EnvioClickPackage {
    weight: number;
    height: number;
    width: number;
    length: number;
}

export interface EnvioClickQuoteRequest {
    idRate: number;
    myShipmentReference: string;
    external_order_id: string;
    requestPickup: boolean;
    pickupDate: string;
    insurance: boolean;
    description: string;
    contentValue: number;
    codValue: number;
    includeGuideCost: boolean;
    codPaymentMethod: string;
    packages: EnvioClickPackage[];
    origin: EnvioClickAddress;
    destination: EnvioClickAddress;
}

export interface EnvioClickGenerateResponse {
    status: string;
    data: {
        tracker: string; // The API returns 'tracker'
        url: string; // PDF URL
        myShipmentReference?: string;
        // other fields if needed
    };
}

export interface EnvioClickRate {
    idRate: number;
    idProduct: number;
    product: string;
    idCarrier: number;
    carrier: string;
    flete: number;
    deliveryDays: number;
    quotationType: string;
}

export interface EnvioClickQuoteResponse {
    status: string;
    data: {
        rates: EnvioClickRate[];
    };
}

export interface EnvioClickTrackHistory {
    date: string;
    status: string;
    description: string;
    location: string;
}

export interface EnvioClickTrackingResponse {
    success: boolean;
    message: string;
    data: {
        trackingNumber: string;
        carrier: string;
        status: string;
        history: EnvioClickTrackHistory[];
    };
}

export interface EnvioClickCancelResponse {
    success: boolean;
    message: string;
    data: {
        status: string;
        message: string;
    };
}
