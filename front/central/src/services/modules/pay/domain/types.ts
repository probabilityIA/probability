export interface PaymentGatewayType {
    id: number;
    name: string;
    code: string;
    image_url?: string;
    is_active: boolean;
    in_development: boolean;
}

export interface PaymentGatewayTypesResponse {
    success: boolean;
    data: PaymentGatewayType[];
    message?: string;
}
