export interface ShippingMargin {
    id: number;
    business_id: number;
    carrier_code: string;
    carrier_name: string;
    margin_amount: number;
    insurance_margin: number;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface CreateShippingMarginDTO {
    carrier_code: string;
    carrier_name: string;
    margin_amount: number;
    insurance_margin: number;
    is_active?: boolean;
}

export interface UpdateShippingMarginDTO {
    carrier_name: string;
    margin_amount: number;
    insurance_margin: number;
    is_active?: boolean;
}

export interface GetShippingMarginsParams {
    page?: number;
    page_size?: number;
    carrier_code?: string;
    business_id?: number;
}

export interface ShippingMarginsListResponse {
    data: ShippingMargin[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface DeleteShippingMarginResponse {
    message?: string;
    error?: string;
}

export interface CarrierOption {
    code: string;
    name: string;
}

export const CARRIER_OPTIONS: CarrierOption[] = [
    { code: 'servientrega', name: 'Servientrega' },
    { code: 'interrapidisimo', name: 'Interrapidisimo' },
    { code: 'coordinadora', name: 'Coordinadora' },
    { code: 'mipaquete', name: 'MiPaquete' },
    { code: 'enviame', name: 'Enviame' },
    { code: 'tcc', name: 'TCC' },
    { code: 'envia', name: 'Envia' },
];
