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
    refreshed_at: string;
    templates: WhatsAppTemplateStatus[];
}

export interface WhatsAppProvisionResult {
    integration_id: number;
    business_id: number;
    waba_id: string;
    created: string[];
    already_exists: string[];
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

export interface WhatsAppConnectionValues {
    waba_id: string;
    phone_number_id: string;
    access_token: string;
    use_platform_token: boolean;
}
