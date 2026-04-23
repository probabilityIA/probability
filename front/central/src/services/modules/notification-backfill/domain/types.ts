export interface BackfillEvent {
    event_code: string;
    event_name: string;
    channel: string;
}

export interface OrderCandidate {
    order_id: string;
    order_number: string;
    customer_phone: string;
    tracking_number?: string;
    status: string;
    carrier?: string;
    carrier_logo_url?: string;
}

export interface BusinessGroup {
    business_id: number;
    business_name: string;
    count: number;
    orders: OrderCandidate[];
}

export interface PreviewRequest {
    event_code: string;
    business_id?: number;
    days?: number;
    limit?: number;
}

export interface PreviewResponse {
    event_code: string;
    total_eligible: number;
    businesses: BusinessGroup[];
}

export interface RunRequest extends PreviewRequest { }

export interface RunResponse {
    job_id: string;
}

export interface JobState {
    id: string;
    event_code: string;
    business_id?: number;
    status: 'running' | 'completed' | 'failed';
    total_eligible: number;
    sent: number;
    skipped: number;
    failed: number;
    started_at: string;
    finished_at?: string;
    error_message?: string;
}

export interface BackfillProgressEvent {
    job_id: string;
    event_code: string;
    status: string;
    total_eligible: number;
    sent: number;
    skipped: number;
    failed: number;
    dry_run: boolean;
    error_message: string;
}
