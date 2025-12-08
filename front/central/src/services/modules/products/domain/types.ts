export interface Product {
    id: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;

    // Identificadores
    business_id: number;
    integration_id?: number;
    integration_type?: string;
    external_id?: string;

    // Información básica
    sku: string;
    name: string;
    description?: string;

    // Precios
    price: number;
    compare_at_price?: number;
    cost_price?: number;
    currency: string;

    // Inventario
    stock: number;
    stock_status?: string;
    manage_stock: boolean;

    // Dimensiones y peso
    weight?: number;
    height?: number;
    width?: number;
    length?: number;

    // Multimedia
    images?: string[];
    thumbnail?: string;

    // Estado
    status: string;
    is_active: boolean;

    // Metadatos
    metadata?: any;
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

export interface SingleResponse<T> {
    success: boolean;
    message: string;
    data: T;
}

export interface ActionResponse {
    success: boolean;
    message: string;
    error?: string;
}

export interface GetProductsParams {
    page?: number;
    page_size?: number;
    business_id?: number;
    integration_id?: number;
    integration_type?: string;
    sku?: string;
    skus?: string;
    name?: string;
    external_id?: string;
    external_ids?: string;
    start_date?: string;
    end_date?: string;
    created_after?: string;
    created_before?: string;
    updated_after?: string;
    updated_before?: string;
    sort_by?: 'id' | 'sku' | 'name' | 'created_at' | 'updated_at' | 'business_id';
    sort_order?: 'asc' | 'desc';
}

export interface CreateProductDTO {
    business_id: number;
    integration_id?: number;
    integration_type?: string;
    external_id?: string;
    sku: string;
    name: string;
    description?: string;
    price: number;
    compare_at_price?: number;
    cost_price?: number;
    currency?: string;
    stock: number;
    stock_status?: string;
    manage_stock?: boolean;
    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    images?: string[];
    thumbnail?: string;
    status?: string;
    is_active?: boolean;
    metadata?: any;
}

export interface UpdateProductDTO {
    sku?: string;
    name?: string;
    description?: string;
    price?: number;
    compare_at_price?: number;
    cost_price?: number;
    currency?: string;
    stock?: number;
    stock_status?: string;
    manage_stock?: boolean;
    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    images?: string[];
    thumbnail?: string;
    status?: string;
    is_active?: boolean;
    metadata?: any;
}
