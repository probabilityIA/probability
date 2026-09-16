import type {
    AssistantAlert,
    AssistantAlertsUnread,
    AssistantHistoryMessage,
    AssistantReply,
    AssistantState,
    FeedbackValue,
    PaginatedResponse,
    ReviewFilters,
    ReviewMessage,
    ReviewSummary,
} from './types';

export interface IAssistantRepository {
    chat(messages: AssistantHistoryMessage[], conversationId: string, pathname: string, businessId?: number | null): Promise<AssistantReply>;
    getState(): Promise<AssistantState>;
    markIntroSeen(): Promise<void>;
    sendFeedback(messageId: string, value: FeedbackValue): Promise<void>;
    markClick(messageId: string): Promise<void>;
    listReviewMessages(filters: ReviewFilters): Promise<PaginatedResponse<ReviewMessage>>;
    getReviewSummary(filters: ReviewFilters): Promise<ReviewSummary>;
    listAlerts(businessId: number | null | undefined, page: number, pageSize: number): Promise<PaginatedResponse<AssistantAlert>>;
    getAlertsUnread(businessId: number | null | undefined): Promise<AssistantAlertsUnread>;
    markAlertsSeen(businessId: number | null | undefined): Promise<void>;
}
