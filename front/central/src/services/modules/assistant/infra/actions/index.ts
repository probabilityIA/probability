'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { AssistantApiError, AssistantApiRepository } from '../repository/api-repository';
import type {
    AssistantHistoryMessage,
    AssistantReply,
    AssistantResult,
    AssistantState,
} from '../../domain/types';

async function withRepository<T>(run: (repo: AssistantApiRepository) => Promise<T>): Promise<AssistantResult<T>> {
    const token = await getAuthToken();
    if (!token) {
        return { success: false, code: 'unauthorized', message: 'Tu sesi\u00f3n termin\u00f3. Vuelve a iniciar sesi\u00f3n.' };
    }
    try {
        return { success: true, data: await run(new AssistantApiRepository(token)) };
    } catch (error) {
        if (error instanceof AssistantApiError) {
            return { success: false, code: error.code, message: error.message, resetAt: error.resetAt };
        }
        return { success: false, code: 'assistant_unavailable', message: 'No pude conectarme para responder.' };
    }
}

export async function sendAssistantMessageAction(
    messages: AssistantHistoryMessage[],
): Promise<AssistantResult<AssistantReply>> {
    return withRepository((repo) => repo.chat(messages));
}

export async function getAssistantStateAction(): Promise<AssistantResult<AssistantState>> {
    return withRepository((repo) => repo.getState());
}

export async function markAssistantIntroSeenAction(): Promise<AssistantResult<null>> {
    return withRepository(async (repo) => {
        await repo.markIntroSeen();
        return null;
    });
}
