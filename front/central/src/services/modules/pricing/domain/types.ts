export interface ClientGroup {
    id: number;
    business_id: number;
    name: string;
    description: string;
    is_active: boolean;
    member_count: number;
    created_at: string;
    updated_at: string;
}

export interface ClientSummary {
    id: number;
    name: string;
    email: string;
    phone: string;
    dni: string;
    group_id: number | null;
    group_name: string;
}

export interface CatalogPriceRow {
    product_id: string;
    product_name: string;
    product_sku: string;
    image_url: string;
    currency: string;
    base_price: number;
    custom_price: number | null;
    difference: number;
}

export interface Paginated<T> {
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface SaveClientGroupInput {
    id?: number;
    name: string;
    description: string;
    is_active: boolean;
}

export interface CatalogPriceTarget {
    client_group_id?: number;
    client_id?: number;
}

export interface CatalogPriceItem {
    product_id: string;
    price: number | null;
}

export interface ActionResult<T = unknown> {
    success: boolean;
    message?: string;
    data?: T;
}

export interface EffectivePrice {
    product_id: string;
    base_price: number;
    final_price: number;
    source: string;
    group_id: number | null;
    group_name: string;
}
