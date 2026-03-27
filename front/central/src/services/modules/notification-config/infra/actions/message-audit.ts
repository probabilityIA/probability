"use server";

import { cookies } from "next/headers";
import { env } from "@/shared/config/env";
import { MessageAuditApiRepository } from "../repository/message-audit-repository";
import type { MessageAuditFilter, ConversationListFilter } from "../../domain/types";

const getRepository = async () => {
    const cookieStore = await cookies();
    const token = cookieStore.get("session_token")?.value || "";
    return new MessageAuditApiRepository(env.API_BASE_URL, token);
};

export async function getMessageAuditLogsAction(filter: MessageAuditFilter) {
    try {
        const repo = await getRepository();
        const result = await repo.list(filter);
        return { success: true, ...result };
    } catch (error: any) {
        return {
            success: false,
            error: error.message,
            data: [],
            total: 0,
            page: 1,
            page_size: 20,
            total_pages: 0,
        };
    }
}

export async function getMessageAuditStatsAction(
    businessId: number,
    dateFrom?: string,
    dateTo?: string
) {
    try {
        const repo = await getRepository();
        const stats = await repo.getStats(businessId, dateFrom, dateTo);
        return { success: true, data: stats };
    } catch (error: any) {
        return { success: false, error: error.message };
    }
}

export async function listConversationsAction(filter: ConversationListFilter) {
    try {
        const repo = await getRepository();
        const result = await repo.listConversations(filter);
        return { success: true, ...result };
    } catch (error: any) {
        return {
            success: false,
            error: error.message,
            data: [],
            total: 0,
            page: 1,
            page_size: 20,
            total_pages: 0,
        };
    }
}

export async function getConversationMessagesAction(
    conversationId: string,
    businessId: number
) {
    try {
        const repo = await getRepository();
        const result = await repo.getConversationMessages(conversationId, businessId);
        return { success: true, data: result };
    } catch (error: any) {
        return { success: false, error: error.message };
    }
}
