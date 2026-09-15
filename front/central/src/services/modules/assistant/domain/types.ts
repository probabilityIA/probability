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

export type ChatEntry =
    | { id: string; kind: 'user'; text: string }
    | { id: string; kind: 'assistant'; text: string; destination: AssistantDestination | null }
    | { id: string; kind: 'system'; text: string }
    | { id: string; kind: 'notice'; text: string; tone: 'error' | 'warning'; retryable: boolean };
