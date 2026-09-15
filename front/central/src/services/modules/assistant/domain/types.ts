export type AssistantRole = 'user' | 'assistant';

export interface AssistantHistoryMessage {
    role: AssistantRole;
    text: string;
}

export interface AssistantDestination {
    key: string;
    label: string;
    route: string;
    description: string;
}

export interface AssistantReply {
    message_id: string;
    conversation_id: string;
    message: string;
    destination: AssistantDestination | null;
}

export interface AssistantState {
    intro_seen: boolean;
    limit: number;
    remaining: number;
    reset_at: string | null;
}

export type AssistantErrorCode =
    | 'rate_limited'
    | 'assistant_unavailable'
    | 'message_too_long'
    | 'invalid_message'
    | 'unauthorized'
    | 'internal_error';

export type AssistantResult<T> =
    | { success: true; data: T }
    | { success: false; code: AssistantErrorCode; message: string; resetAt?: string | null };

export type AvatarMood = 'idle' | 'thinking' | 'pointing';

export type FeedbackValue = -1 | 0 | 1;

export type ChatEntry =
    | { id: string; kind: 'user'; text: string }
    | { id: string; kind: 'assistant'; text: string; destination: AssistantDestination | null; messageId?: string; feedback?: FeedbackValue }
    | { id: string; kind: 'system'; text: string }
    | { id: string; kind: 'notice'; text: string; tone: 'error' | 'warning'; retryable: boolean };

export type ReviewKind = 'all' | 'no_destination' | 'negative' | 'positive' | 'errors' | 'not_clicked';

export interface ReviewFilters {
    business_id?: number;
    kind?: ReviewKind;
    search?: string;
    from?: string;
    to?: string;
    page?: number;
    page_size?: number;
}

export interface ReviewMessage {
    id: string;
    conversation_id: string;
    business_id: number | null;
    business_name: string;
    user_id: number;
    user_name: string;
    user_email: string;
    pathname: string;
    question: string;
    answer: string;
    destination_key: string;
    destination_route: string;
    error_code: string;
    model: string;
    input_tokens: number;
    output_tokens: number;
    latency_ms: number;
    feedback: FeedbackValue;
    feedback_at: string | null;
    clicked_at: string | null;
    created_at: string;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface ReviewSummary {
    messages: number;
    conversations: number;
    users: number;
    input_tokens: number;
    output_tokens: number;
    no_destination: number;
    errors: number;
    positive: number;
    negative: number;
    with_destination: number;
    clicked: number;
    estimated_cost_usd: number;
    top_destinations: { key: string; count: number }[];
}
