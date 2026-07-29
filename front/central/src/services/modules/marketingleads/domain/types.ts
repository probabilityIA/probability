export interface MarketingLead {
    id: number;
    name: string;
    email: string;
    phone: string;
    survey_slug: string;
    score_total: number;
    score_max: number;
    level: string;
    whats_app_message_id?: string;
    created_at: string;
}

export interface GetMarketingLeadsParams {
    page?: number;
    page_size?: number;
}

export interface PaginatedLeadsResponse {
    success: boolean;
    message?: string;
    data: {
        leads: MarketingLead[];
        total: number;
        page: number;
        page_size: number;
        total_pages: number;
    };
}
