import type { AssistantHistoryMessage, AssistantReply, AssistantState } from './types';

export interface IAssistantRepository {
    chat(messages: AssistantHistoryMessage[]): Promise<AssistantReply>;
    getState(): Promise<AssistantState>;
    markIntroSeen(): Promise<void>;
}
