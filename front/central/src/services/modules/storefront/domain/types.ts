export interface StorefrontProduct {
    id: string;
    name: string;
    description: string;
    short_description: string;
    price: number;
    compare_at_price?: number;
    currency: string;
    image_url: string;
    images?: string[];
    sku: string;
    stock_quantity: number;
    track_inventory: boolean;
    category: string;
    brand: string;
    is_featured: boolean;
    created_at: string;
}

export interface StorefrontOrder {
    id: string;
    order_number: string;
    status: string;
    total_amount: number;
    currency: string;
    created_at: string;
    items: StorefrontOrderItem[];
}

export interface StorefrontOrderItem {
    product_name: string;
    quantity: number;
    unit_price: number;
    total_price: number;
    image_url?: string;
}

export interface CreateStorefrontOrderDTO {
    items: { product_id: string; quantity: number }[];
    notes?: string;
    address?: {
        first_name: string;
        last_name?: string;
        phone?: string;
        street: string;
        street2?: string;
        city: string;
        state?: string;
        country?: string;
        postal_code?: string;
        instructions?: string;
    };
}

export interface StorefrontClient {
    id: number;
    name: string;
    email: string | null;
    phone: string;
    dni: string | null;
}

export interface CreateClientDTO {
    name: string;
    email: string;
    password?: string;
    phone?: string;
    dni?: string;
}

export interface CreateClientResult {
    client: StorefrontClient;
    temp_password?: string;
}

export interface StorefrontFamily {
    id: number;
    name: string;
    parent_family_id: number | null;
}

export interface CatalogLayout {
    columns: number;
    rows: number;
}

export interface CatalogBanner {
    enabled: boolean;
    image_url: string;
}

export interface CatalogFiltersResult {
    categories: string[];
    families: StorefrontFamily[];
    layout: CatalogLayout;
    banner: CatalogBanner;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}
