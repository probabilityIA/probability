export type DisplayType = 'modal_image' | 'modal_text' | 'ticker';
export type FrequencyType = 'once' | 'daily' | 'always' | 'requires_acceptance';
export type AnnouncementStatus = 'draft' | 'scheduled' | 'active' | 'inactive';
export type ViewAction = 'viewed' | 'closed' | 'clicked_link' | 'accepted';

export interface AnnouncementCategory {
    id: number;
    code: string;
    name: string;
    icon: string;
    color: string;
}

export interface AnnouncementImage {
    id: number;
    image_url: string;
    sort_order: number;
}

export interface AnnouncementLink {
    id: number;
    label: string;
    url: string;
    sort_order: number;
}

export interface AnnouncementTarget {
    id: number;
    business_id: number;
}

export interface AnnouncementInfo {
    id: number;
    business_id: number | null;
    category_id: number;
    category?: AnnouncementCategory;
    title: string;
    message: string;
    display_type: DisplayType;
    frequency_type: FrequencyType;
    priority: number;
    is_global: boolean;
    status: AnnouncementStatus;
    starts_at: string | null;
    ends_at: string | null;
    force_redisplay: boolean;
    created_by_id: number;
    created_at: string;
    updated_at: string;
    images: AnnouncementImage[];
    links: AnnouncementLink[];
    targets: AnnouncementTarget[];
}

export interface AnnouncementStats {
    total_views: number;
    unique_users: number;
    total_clicks: number;
    total_acceptances: number;
    total_closed: number;
}

export interface AnnouncementView {
    id: number;
    announcement_id: number;
    user_id: number;
    business_id: number;
    action: ViewAction;
    link_id: number | null;
    viewed_at: string;
    created_at: string;
}

export interface CreateLinkDTO {
    label: string;
    url: string;
    sort_order: number;
}

export interface CreateAnnouncementDTO {
    business_id?: number | null;
    category_id: number;
    title: string;
    message: string;
    display_type: DisplayType;
    frequency_type: FrequencyType;
    priority: number;
    is_global: boolean;
    starts_at?: string;
    ends_at?: string;
    links: CreateLinkDTO[];
    target_ids: number[];
}

export interface UpdateAnnouncementDTO {
    business_id?: number | null;
    category_id: number;
    title: string;
    message: string;
    display_type: DisplayType;
    frequency_type: FrequencyType;
    priority: number;
    is_global: boolean;
    starts_at?: string;
    ends_at?: string;
    links: CreateLinkDTO[];
    target_ids: number[];
}

export interface RegisterViewDTO {
    action: ViewAction;
    link_id?: number | null;
}

export interface ChangeStatusDTO {
    status: AnnouncementStatus;
}

export interface GetAnnouncementsParams {
    page?: number;
    page_size?: number;
    business_id?: number;
    status?: AnnouncementStatus;
    category_id?: number;
    search?: string;
}

export interface AnnouncementsListResponse {
    data: AnnouncementInfo[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface DeleteAnnouncementResponse {
    success: boolean;
    message: string;
}
