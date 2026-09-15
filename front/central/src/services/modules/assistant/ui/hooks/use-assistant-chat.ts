'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { sendAssistantMessageAction } from '../../infra/actions';
import {
    GREETING,
    buildHistory,
    entryId,
    isCurrentRoute,
    rateLimitText,
} from '../../app/use-cases';
import type {
    AssistantDestination,
    AssistantHistoryMessage,
    AssistantState,
    AvatarMood,
    ChatEntry,
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
    const moodTimer = useRef<number | null>(null);

    useEffect(() => {
        entriesRef.current = entries;
    }, [entries]);

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
            const result = await sendAssistantMessageAction(history).catch(() => null);
            setPending(false);

            if (result?.success) {
                append({
                    id: entryId(),
                    kind: 'assistant',
                    text: result.data.message,
                    destination: result.data.destination,
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
        [append, limit, point],
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
        (destination: AssistantDestination) => {
            if (isCurrentRoute(pathname, destination.route)) return;
            router.push(destination.route);
            append({ id: entryId(), kind: 'system', text: `Te llev\u00e9 a ${destination.label} \u00b7 ${destination.route}` });
            point();
        },
        [append, pathname, point, router],
    );

    const reset = useCallback(() => {
        setEntries([greetingEntry()]);
        setMood('idle');
    }, []);

    const applyState = useCallback((state: AssistantState) => {
        setLimit(state.limit);
        if (state.remaining <= 0 && state.reset_at) setBlockedUntil(state.reset_at);
    }, []);

    return { entries, pending, mood, blocked, blockedUntil, pathname, send, retry, goTo, reset, applyState };
}

export type AssistantChat = ReturnType<typeof useAssistantChat>;
