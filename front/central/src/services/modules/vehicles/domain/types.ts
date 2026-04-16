// ============================================
// ENTIDADES
// ============================================

export interface VehicleInfo {
    id: number;
    business_id: number;
    type: string;
    license_plate: string;
    brand: string;
    model: string;
    year: number | null;
    color: string;
    status: string;
    weight_capacity_kg: number | null;
    volume_capacity_m3: number | null;
    photo_url: string;
    insurance_expiry: string | null;
    registration_expiry: string | null;
    created_at: string;
    updated_at: string;
}

// ============================================
// DTOs
// ============================================

export interface CreateVehicleDTO {
    type: string;
    license_plate: string;
    brand?: string;
    model?: string;
    year?: number;
    color?: string;
    weight_capacity_kg?: number;
    volume_capacity_m3?: number;
    insurance_expiry?: string;
    registration_expiry?: string;
}

export interface UpdateVehicleDTO {
    type?: string;
    license_plate?: string;
    brand?: string;
    model?: string;
    year?: number | null;
    color?: string;
    status?: string;
    weight_capacity_kg?: number | null;
    volume_capacity_m3?: number | null;
    insurance_expiry?: string | null;
    registration_expiry?: string | null;
}

export interface GetVehiclesParams {
    page?: number;
    page_size?: number;
    search?: string;
    type?: string;
    status?: string;
    business_id?: number;
}

// ============================================
// RESPONSES (coinciden con el backend exactamente)
// ============================================

export interface VehiclesListResponse {
    data: VehicleInfo[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface DeleteVehicleResponse {
    message?: string;
    error?: string;
}
