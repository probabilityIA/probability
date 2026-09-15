'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import {
    markAssistantClickAction,
    sendAssistantMessageAction,
    submitAssistantFeedbackAction,
} from '../../infra/actions';
import {
    GREETING,
    buildHistory,
    entryId,
    hubEnvironmentFor,
    isCurrentRoute,
    isHubDestination,
    newConversationId,
    nextFeedback,
    rateLimitText,
} from '../../app/use-cases';
import { queueIntegrationsHub, requestIntegrationsHub } from '@/services/modules/my-integrations/ui/open-hub';
import type {
    AssistantDestination,
    AssistantHistoryMessage,
    AssistantState,
    AvatarMood,
    ChatEntry,
    FeedbackValue,
} from '../../domain/types';

const POINTING_MS = 1800;

function greetingEntry(): ChatEntry {
    return { id: entryId(), kind: 'assistant', text: GREETING, destination: null };
}

function withoutNotices(entries: ChatEntry[]): ChatEntry[] {
    return entries.filter((e) => e.kind !== 'notice');
}

export function useAssistantChat() {
    const router = useRouter();
    const pathname = usePathname();
    const [entries, setEntries] = useState<ChatEntry[]>(() => [greetingEntry()]);
    const [pending, setPending] = useState(false);
    const [mood, setMood] = useState<AvatarMood>('idle');
    const [limit, setLimit] = useState(30);
    const [blockedUntil, setBlockedUntil] = useState<string | null>(null);
    const entriesRef = useRef(entries);
    const conversationRef = useRef<string>('');
    const moodTimer = useRef<number | null>(null);

    useEffect(() => {
        entriesRef.current = entries;
    }, [entries]);

    useEffect(() => {
        if (!conversationRef.current) conversationRef.current = newConversationId();
    }, []);

    useEffect(() => {
        if (!blockedUntil) return;
        const wait = new Date(blockedUntil).getTime() - Date.now();
        if (Number.isNaN(wait) || wait <= 0) {
            setBlockedUntil(null);
            return;
        }
        const timer = window.setTimeout(() => setBlockedUntil(null), wait);
        return () => window.clearTimeout(timer);
    }, [blockedUntil]);

    useEffect(() => () => {
        if (moodTimer.current) window.clearTimeout(moodTimer.current);
    }, []);

    const point = useCallback(() => {
        setMood('pointing');
        if (moodTimer.current) window.clearTimeout(moodTimer.current);
        moodTimer.current = window.setTimeout(() => setMood('idle'), POINTING_MS);
    }, []);

    const append = useCallback((entry: ChatEntry) => {
        setEntries((prev) => [...prev, entry]);
    }, []);

    const request = useCallback(
        async (history: AssistantHistoryMessage[]) => {
            setPending(true);
            setMood('thinking');
            if (!conversationRef.current) conversationRef.current = newConversationId();
            const result = await sendAssistantMessageAction(history, conversationRef.current, pathname).catch(() => null);
            setPending(false);

            if (result?.success) {
                if (result.data.conversation_id) conversationRef.current = result.data.conversation_id;
                append({
                    id: entryId(),
                    kind: 'assistant',
                    text: result.data.message,
                    destination: result.data.destination,
                    messageId: result.data.message_id || undefined,
                    feedback: 0,
                });
                if (result.data.destination) point();
                else setMood('idle');
                return;
            }

            setMood('idle');
            if (result && result.code === 'rate_limited') {
                setBlockedUntil(result.resetAt ?? null);
                append({ id: entryId(), kind: 'notice', tone: 'warning', retryable: false, text: rateLimitText(limit, result.resetAt) });
                return;
            }
            if (result && (result.code === 'message_too_long' || result.code === 'invalid_message' || result.code === 'unauthorized')) {
                append({ id: entryId(), kind: 'notice', tone: 'error', retryable: false, text: result.message });
                return;
            }
            append({
                id: entryId(),
                kind: 'notice',
                tone: 'error',
                retryable: true,
                text: 'No pude conectarme para responder. Tu mensaje qued\u00f3 guardado.',
            });
        },
        [append, limit, pathname, point],
    );

    const blocked = blockedUntil !== null;

    const send = useCallback(
        (text: string) => {
            const trimmed = text.trim();
            if (!trimmed || pending || blocked) return;
            const next = [...withoutNotices(entriesRef.current), { id: entryId(), kind: 'user' as const, text: trimmed }];
            setEntries(next);
            void request(buildHistory(next));
        },
        [blocked, pending, request],
    );

    const retry = useCallback(() => {
        if (pending || blocked) return;
        const clean = withoutNotices(entriesRef.current);
        setEntries(clean);
        void request(buildHistory(clean));
    }, [blocked, pending, request]);

    const goTo = useCallback(
        (destination: AssistantDestination, messageId?: string) => {
            if (messageId) void markAssistantClickAction(messageId).catch(() => null);

            const onRoute = isCurrentRoute(pathname, destination.route);
            if (isHubDestination(destination.key)) {
                const intent = { environment: hubEnvironmentFor(destination.key) };
                if (onRoute) {
                    requestIntegrationsHub(intent);
                } else {
                    queueIntegrationsHub(intent);
                    router.push(destination.route);
                }
                append({ id: entryId(), kind: 'system', text: `Te abr\u00ed ${destination.label}` });
                point();
                return;
            }
            if (onRoute) return;
            router.push(destination.route);
            append({ id: entryId(), kind: 'system', text: `Te llev\u00e9 a ${destination.label} \u00b7 ${destination.route}` });
            point();
        },
        [append, pathname, point, router],
    );

    const rate = useCallback((id: string, pressed: 1 | -1) => {
        const entry = entriesRef.current.find((e) => e.id === id);
        if (!entry || entry.kind !== 'assistant' || !entry.messageId) return;

        const previous: FeedbackValue = entry.feedback ?? 0;
        const value = nextFeedback(previous, pressed);
        const messageId = entry.messageId;
        const apply = (feedback: FeedbackValue) =>
            setEntries((prev) => prev.map((e) => (e.id === id && e.kind === 'assistant' ? { ...e, feedback } : e)));

        apply(value);
        void submitAssistantFeedbackAction(messageId, value)
            .then((result) => {
                if (!result.success) apply(previous);
            })
            .catch(() => apply(previous));
    }, []);

    const reset = useCallback(() => {
        conversationRef.current = newConversationId();
        setEntries([greetingEntry()]);
        setMood('idle');
    }, []);

    const applyState = useCallback((state: AssistantState) => {
        setLimit(state.limit);
        if (state.remaining <= 0 && state.reset_at) setBlockedUntil(state.reset_at);
    }, []);

    return { entries, pending, mood, blocked, blockedUntil, pathname, send, retry, goTo, rate, reset, applyState };
}

export type AssistantChat = ReturnType<typeof useAssistantChat>;
