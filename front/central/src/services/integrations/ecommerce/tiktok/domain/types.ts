export interface TikTokConfig {
    shop_name: string;
}

export interface TikTokCredentials {
    client_key: string;
    client_secret: string;
}

export interface TikTokIntegrationData {
    name: string;
    config: TikTokConfig;
    credentials: TikTokCredentials;
    is_active: boolean;
}

export interface TikTokShopSnapshot {
    id: number;
    integration_id: number;
    shop_id: number;
    shop_name: string;
    shop_rating: string;
    shop_review_count: number;
    item_sold_count: number;
    shop_performance_value: number;
    created_at: string;
}
