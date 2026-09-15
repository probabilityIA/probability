import type {
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
}
