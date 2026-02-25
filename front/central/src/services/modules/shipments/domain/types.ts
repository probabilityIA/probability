export interface Shipment {
    id: number;
    created_at: string;
    updated_at: string;
    order_id?: string;
    client_name?: string;
    destination_address?: string;
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
    is_test: boolean;
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
    business_id?: number;
    start_date?: string;
    end_date?: string;
    shipped_after?: string;
    shipped_before?: string;
    delivered_after?: string;
    delivered_before?: string;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
    is_test?: boolean;
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
    company?: string;
    firstName?: string;
    lastName?: string;
    email?: string;
    phone?: string;
    address: string;
    suburb?: string;
    crossStreet?: string;
    reference?: string;
    daneCode: string;
}

export interface EnvioClickPackage {
    weight: number;
    height: number;
    width: number;
    length: number;
}

export interface EnvioClickQuoteRequest {
    business_id?: number;
    idRate?: number;
    myShipmentReference?: string;
    external_order_id?: string;
    order_uuid?: string;
    requestPickup?: boolean;
    pickupDate?: string;
    insurance?: boolean;
    description: string;
    contentValue: number;
    codValue?: number;
    includeGuideCost: boolean;
    codPaymentMethod: string;
    totalCost?: number;
    packages: EnvioClickPackage[];
    origin: EnvioClickAddress;
    destination: EnvioClickAddress;
}

export interface EnvioClickGenerateResponse {
    status?: string;
    // Async 202 fields
    success?: boolean;
    message?: string;
    correlation_id?: string;
    shipment_id?: number;
    // Sync response fields (legacy)
    data?: {
        tracker: string;
        url: string;
        myShipmentReference?: string;
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
    minimumInsurance?: number;
    extraInsurance?: number;
    cod?: boolean;
}

// Response from POST /shipments/quote (202 Accepted - async, result via SSE)
export interface EnvioClickQuoteResponse {
    success: boolean;
    message: string;
    correlation_id: string;
    data?: {
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

export interface CreateShipmentRequest {
    order_id?: string;
    client_name?: string;
    destination_address?: string;
    tracking_number?: string;
    carrier?: string;
    status?: string;
}

export interface OriginAddress {
    id: number;
    business_id: number;
    alias: string;
    company: string;
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
    street: string;
    suburb?: string;
    city_dane_code: string;
    city: string;
    state: string;
    postal_code?: string;
    is_default: boolean;
    created_at: string;
    updated_at: string;
}

export type CreateOriginAddressRequest = Omit<OriginAddress, 'id' | 'business_id' | 'created_at' | 'updated_at' | 'is_default'> & {
    is_default?: boolean;
};

export type UpdateOriginAddressRequest = Partial<CreateOriginAddressRequest>;

// ===================================
// SSE EVENTS (Tiempo Real)
// ===================================

export type ShipmentSSEEventType =
  | 'shipment.quote_received'
  | 'shipment.quote_failed'
  | 'shipment.guide_generated'
  | 'shipment.guide_failed'
  | 'shipment.tracking_updated'
  | 'shipment.tracking_failed'
  | 'shipment.cancelled'
  | 'shipment.cancel_failed';

export interface ShipmentSSEEvent {
  id: string;
  type: string;
  business_id: string;
  timestamp: string;
  data: ShipmentSSEEventData;
  metadata: Record<string, any>;
}

export interface ShipmentSSEEventData {
  shipment_id?: number;
  correlation_id?: string;
  tracking_number?: string;
  label_url?: string;
  error_message?: string;
  quotes?: Record<string, any>;
  tracking?: Record<string, any>;
}

