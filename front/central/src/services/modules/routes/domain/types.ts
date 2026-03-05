// ============================================
// ENTIDADES
// ============================================

export interface RouteInfo {
    id: number;
    business_id: number;
    driver_id: number | null;
    driver_name: string | null;
    vehicle_id: number | null;
    vehicle_plate: string | null;
    status: string;
    date: string;
    start_time: string | null;
    end_time: string | null;
    origin_address: string | null;
    total_stops: number;
    completed_stops: number;
    failed_stops: number;
    notes: string | null;
    created_at: string;
    updated_at: string;
}

export interface RouteStopInfo {
    id: number;
    route_id: number;
    order_id: string | null;
    sequence: number;
    status: string;
    address: string;
    city: string | null;
    lat: number | null;
    lng: number | null;
    customer_name: string;
    customer_phone: string | null;
    estimated_arrival: string | null;
    actual_arrival: string | null;
    actual_departure: string | null;
    delivery_notes: string | null;
    failure_reason: string | null;
    created_at: string;
    updated_at: string;
}

export interface RouteDetail extends RouteInfo {
    actual_start_time: string | null;
    actual_end_time: string | null;
    origin_warehouse_id: number | null;
    origin_lat: number | null;
    origin_lng: number | null;
    total_distance_km: number | null;
    total_duration_min: number | null;
    stops: RouteStopInfo[];
}

// ============================================
// DTOs
// ============================================

export interface CreateRouteStopDTO {
    order_id?: string;
    address: string;
    city?: string;
    lat?: number;
    lng?: number;
    customer_name: string;
    customer_phone?: string;
    delivery_notes?: string;
}

export interface CreateRouteDTO {
    date: string;
    driver_id?: number;
    vehicle_id?: number;
    origin_address?: string;
    origin_lat?: number;
    origin_lng?: number;
    notes?: string;
    stops?: CreateRouteStopDTO[];
}

export interface UpdateRouteDTO {
    driver_id?: number;
    vehicle_id?: number;
    date?: string;
    origin_address?: string;
    origin_lat?: number;
    origin_lng?: number;
    notes?: string;
}

export interface AddStopDTO {
    order_id?: string;
    address: string;
    city?: string;
    lat?: number;
    lng?: number;
    customer_name: string;
    customer_phone?: string;
    delivery_notes?: string;
}

export interface UpdateStopDTO {
    address?: string;
    city?: string;
    lat?: number;
    lng?: number;
    customer_name?: string;
    customer_phone?: string;
    delivery_notes?: string;
}

export interface UpdateStopStatusDTO {
    status: string;
    failure_reason?: string;
}

export interface ReorderStopsDTO {
    stop_ids: number[];
}

export interface GetRoutesParams {
    page?: number;
    page_size?: number;
    status?: string;
    driver_id?: number;
    date_from?: string;
    date_to?: string;
    search?: string;
    business_id?: number;
}

// ============================================
// RESPONSES
// ============================================

export interface RoutesListResponse {
    data: RouteInfo[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface DeleteRouteResponse {
    message?: string;
    error?: string;
}
