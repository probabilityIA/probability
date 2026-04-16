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

export interface RegisterDTO {
    name: string;
    email: string;
    password: string;
    phone?: string;
    dni?: string;
    business_code: string;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}
