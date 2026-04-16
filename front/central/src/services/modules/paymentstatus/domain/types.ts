export interface PaymentStatusInfo {
    id: number;
    code: string;
    name: string;
    description?: string;
    category?: string;
    color?: string;
    icon?: string;
    is_active: boolean;
}

export interface GetPaymentStatusesParams {
    is_active?: boolean;
}

export interface PaymentStatusesResponse {
    success: boolean;
    message?: string;
    data: PaymentStatusInfo[];
}
