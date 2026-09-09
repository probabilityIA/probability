export interface WhatsAppTemplateStatus {
    name: string;
    language: string;
    status: string;
    category: string;
    meta_id: string;
    reason: string;
    provisioned: boolean;
    updated_at: string;
}

export interface WhatsAppTemplatesSnapshot {
    integration_id: number;
    business_id: number;
    waba_id: string;
    hosted_by_platform: boolean;
    refreshed_at: string;
    templates: WhatsAppTemplateStatus[];
}

export interface WhatsAppProvisionResult {
    integration_id: number;
    business_id: number;
    waba_id: string;
    created: string[];
    already_exists: string[];
    skipped: string[];
    failed: Record<string, string>;
    templates: WhatsAppTemplateStatus[];
}

export interface WhatsAppTemplatesResponse {
    success: boolean;
    message?: string;
    data?: WhatsAppTemplatesSnapshot;
}

export interface WhatsAppProvisionResponse {
    success: boolean;
    message?: string;
    data?: WhatsAppProvisionResult;
}

export interface WhatsAppConnectionResult {
    integration_id: number;
    business_id: number;
    own_number: boolean;
    waba_id: string;
    phone_number_id: string;
    display_phone_number: string;
    verified_name: string;
    quality_rating: string;
    platform_token: boolean;
    hosted_by_platform: boolean;
}

export interface WhatsAppConnectionResponse {
    success: boolean;
    message?: string;
    data?: WhatsAppConnectionResult;
}

export interface WhatsAppConnectionValues {
    hosted: boolean;
    waba_id: string;
    phone_number_id: string;
    access_token: string;
    use_platform_token: boolean;
}

export interface WhatsAppNumberState {
    integration_id: number;
    business_id: number;
    phone_number_id: string;
    status: string;
    display_phone_number: string;
    verified_name: string;
    name_status: string;
    code_verification_status: string;
    quality_rating: string;
    hosted_by_platform: boolean;
    active: boolean;
    pin: string;
}

export interface WhatsAppNumberResponse {
    success: boolean;
    message?: string;
    data?: WhatsAppNumberState;
}

export interface WhatsAppAddNumberValues {
    country_code: string;
    phone_number: string;
    verified_name: string;
}

export interface WhatsAppEmbeddedSignupConfig {
    enabled: boolean;
    app_id?: string;
    config_id?: string;
    graph_version?: string;
}

export interface WhatsAppEmbeddedSignupConfigResponse {
    success: boolean;
    message?: string;
    data?: WhatsAppEmbeddedSignupConfig;
}

export interface WhatsAppEmbeddedSignupResult {
    integration_id: number;
    business_id: number;
    waba_id: string;
    phone_number_id: string;
    display_phone_number: string;
    verified_name: string;
    quality_rating: string;
    registered: boolean;
    pin: string;
    warning: string;
}

export interface WhatsAppEmbeddedSignupResponse {
    success: boolean;
    message?: string;
    data?: WhatsAppEmbeddedSignupResult;
}

export interface WhatsAppEmbeddedSignupPayload {
    code: string;
    waba_id: string;
    phone_number_id: string;
}

export interface WhatsAppBusinessProfile {
    about: string;
    address: string;
    description: string;
    email: string;
    profile_picture_url: string;
    websites: string[] | null;
    vertical: string;
}

export interface WhatsAppBusinessProfileResponse {
    success: boolean;
    message?: string;
    data?: WhatsAppBusinessProfile;
}

export interface WhatsAppBusinessProfileValues {
    about?: string;
    address?: string;
    description?: string;
    email?: string;
    vertical?: string;
    websites?: string[];
}
