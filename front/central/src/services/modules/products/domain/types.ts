export interface ProductFamilySummary {
    id: number;
    business_id: number;
    name: string;
    title?: string;
    description?: string;
    slug?: string;
    category?: string;
    brand?: string;
    image_url?: string;
    status: string;
    is_active: boolean;
    variant_axes?: any;
    created_at: string;
    updated_at: string;
}

export interface ProductFamily {
    id: number;
    business_id: number;
    name: string;
    title?: string;
    description?: string;
    slug?: string;
    category?: string;
    brand?: string;
    image_url?: string;
    status: string;
    is_active: boolean;
    variant_axes?: any;
    metadata?: any;
    variant_count: number;
    variants?: Product[];
    created_at: string;
    updated_at: string;
}

export interface Product {
    id: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;

    business_id: number;
    integration_id?: number;
    integration_type?: string;
    external_id?: string;

    sku: string;
    name: string;
    title?: string;
    description?: string;
    short_description?: string;
    slug?: string;

    barcode?: string;
    family_id?: number;
    family?: ProductFamilySummary;
    variant_label?: string;
    variant_attributes?: any;

    price: number;
    compare_at_price?: number;
    cost_price?: number;
    currency: string;

    stock: number;
    stock_quantity?: number;
    stock_status?: string;
    manage_stock: boolean;
    track_inventory: boolean;
    allow_backorder?: boolean;
    quantity?: number;
    low_stock_threshold?: number;

    weight?: number;
    weight_unit?: string;
    height?: number;
    width?: number;
    length?: number;
    dimension_unit?: string;

    image_url?: string;
    images?: string[];
    thumbnail?: string;
    video_url?: string;

    category?: string;
    tags?: any;
    brand?: string;

    status: string;
    is_active: boolean;
    is_featured?: boolean;

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

export interface UploadImageResponse {
    success: boolean;
    message: string;
    image_url: string;
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
    barcode?: string;
    family_id?: number;
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
    title?: string;
    description?: string;
    short_description?: string;
    slug?: string;
    barcode?: string;
    family_id?: number;
    family?: CreateProductFamilyDTO;
    variant_label?: string;
    variant_attributes?: any;
    price: number;
    compare_at_price?: number;
    cost_price?: number;
    currency?: string;
    stock?: number;
    stock_quantity?: number;
    stock_status?: string;
    manage_stock?: boolean;
    track_inventory?: boolean;
    allow_backorder?: boolean;
    low_stock_threshold?: number;
    weight?: number;
    weight_unit?: string;
    height?: number;
    width?: number;
    length?: number;
    dimension_unit?: string;
    images?: string[];
    thumbnail?: string;
    category?: string;
    brand?: string;
    status?: string;
    is_active?: boolean;
    is_featured?: boolean;
    metadata?: any;
}

export interface UpdateProductDTO {
    sku?: string;
    name?: string;
    title?: string;
    description?: string;
    short_description?: string;
    slug?: string;
    barcode?: string;
    family_id?: number;
    variant_label?: string;
    variant_attributes?: any;
    price?: number;
    compare_at_price?: number;
    cost_price?: number;
    currency?: string;
    stock?: number;
    stock_quantity?: number;
    stock_status?: string;
    manage_stock?: boolean;
    track_inventory?: boolean;
    allow_backorder?: boolean;
    low_stock_threshold?: number;
    weight?: number;
    weight_unit?: string;
    height?: number;
    width?: number;
    length?: number;
    dimension_unit?: string;
    image_url?: string;
    images?: string[];
    thumbnail?: string;
    category?: string;
    brand?: string;
    status?: string;
    is_active?: boolean;
    is_featured?: boolean;
    metadata?: any;
}

export interface GetFamiliesParams {
    page?: number;
    page_size?: number;
    business_id?: number;
    name?: string;
    category?: string;
    brand?: string;
    status?: string;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
}

export interface CreateProductFamilyDTO {
    name: string;
    title?: string;
    description?: string;
    slug?: string;
    category?: string;
    brand?: string;
    image_url?: string;
    status?: string;
    is_active?: boolean;
    variant_axes?: any;
    metadata?: any;
}

export interface UpdateProductFamilyDTO {
    name?: string;
    title?: string;
    description?: string;
    slug?: string;
    category?: string;
    brand?: string;
    image_url?: string;
    status?: string;
    is_active?: boolean;
    variant_axes?: any;
    metadata?: any;
}

export interface ProductIntegration {
    id: number;
    product_id: string;
    integration_id: number;
    integration_type?: string;
    integration_name?: string;
    external_product_id: string;
    external_variant_id?: string;
    external_sku?: string;
    external_barcode?: string;
    created_at: string;
    updated_at: string;
}

export interface AddProductIntegrationDTO {
    integration_id: number;
    external_product_id: string;
    external_variant_id?: string;
    external_sku?: string;
    external_barcode?: string;
}

export interface UpdateProductIntegrationDTO {
    external_product_id?: string;
    external_variant_id?: string;
    external_sku?: string;
    external_barcode?: string;
}

export interface ProductIntegrationsResponse {
    success: boolean;
    message: string;
    data: ProductIntegration[];
    total: number;
}
