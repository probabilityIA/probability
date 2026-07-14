export interface PaymentMethodInfo {
    id: number;
    code: string;
    name: string;
}

export interface PaymentMethodsResponse {
    success: boolean;
    message?: string;
    data: PaymentMethodInfo[];
}

export interface PaginatedPaymentMethodsResponse {
    success?: boolean;
    message?: string;
    data: PaymentMethodInfo[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}
