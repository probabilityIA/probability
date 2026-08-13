export interface SiigoReferral {
    id: number;
    name: string;
    email: string;
    phone: string;
    order_range: string;
    created_at: string;
}

export interface GetSiigoReferralsParams {
    page?: number;
    page_size?: number;
}

export interface PaginatedSiigoReferralsResponse {
    success: boolean;
    message?: string;
    data: {
        referrals: SiigoReferral[];
        total: number;
        page: number;
        page_size: number;
        total_pages: number;
    };
}
