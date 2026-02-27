// ============================================
// ENTIDADES
// ============================================

export interface CustomerInfo {
    id: number;
    business_id: number;
    name: string;
    email: string;
    phone: string;
    dni: string | null;
    created_at: string;
    updated_at: string;
}

export interface CustomerDetail extends CustomerInfo {
    order_count: number;
    total_spent: number;
    last_order_at: string | null;
}

// ============================================
// DTOs
// ============================================

export interface CreateCustomerDTO {
    name: string;
    email?: string;
    phone?: string;
    dni?: string | null;
}

export interface UpdateCustomerDTO {
    name: string;
    email?: string;
    phone?: string;
    dni?: string | null;
}

export interface GetCustomersParams {
    page?: number;
    page_size?: number;
    search?: string;
}

// ============================================
// RESPONSES (coinciden con el backend exactamente)
// ============================================

export interface CustomersListResponse {
    data: CustomerInfo[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface DeleteCustomerResponse {
    message?: string;
    error?: string;
}
