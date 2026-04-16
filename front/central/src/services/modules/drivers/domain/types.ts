// ============================================
// ENTIDADES
// ============================================

export interface DriverInfo {
    id: number;
    business_id: number;
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
    identification: string;
    status: string;
    photo_url: string;
    license_type: string;
    license_expiry: string | null;
    warehouse_id: number | null;
    notes: string | null;
    created_at: string;
    updated_at: string;
}

// ============================================
// DTOs
// ============================================

export interface CreateDriverDTO {
    first_name: string;
    last_name: string;
    email?: string;
    phone: string;
    identification: string;
    license_type?: string;
    license_expiry?: string;
    warehouse_id?: number;
    notes?: string;
}

export interface UpdateDriverDTO {
    first_name?: string;
    last_name?: string;
    email?: string;
    phone?: string;
    identification?: string;
    status?: string;
    license_type?: string;
    license_expiry?: string;
    warehouse_id?: number | null;
    notes?: string | null;
}

export interface GetDriversParams {
    page?: number;
    page_size?: number;
    search?: string;
    status?: string;
    business_id?: number;
}

// ============================================
// RESPONSES (coinciden con el backend exactamente)
// ============================================

export interface DriversListResponse {
    data: DriverInfo[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface DeleteDriverResponse {
    message?: string;
    error?: string;
}
