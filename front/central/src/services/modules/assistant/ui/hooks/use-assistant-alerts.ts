'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { useSelectedBusiness } from '@/shared/contexts/selected-business-context';
import { useSSE } from '@/shared/hooks/use-sse';
import { getAssistantAlertsAction, getAssistantAlertsUnreadAction, markAssistantAlertsSeenAction } from '../../infra/actions';
import { ALERT_EVENT_TYPES, ALERT_REFRESH_DELAY_MS, ALERT_TOAST_MS } from '../../app/use-cases';
import type { AssistantAlert } from '../../domain/types';

const PAGE_SIZE = 20;

export function useAssistantAlerts(enabled: boolean) {
    const { selectedBusinessId } = useSelectedBusiness();
    const [unread, setUnread] = useState(0);
    const [toast, setToast] = useState<AssistantAlert | null>(null);
    const [alerts, setAlerts] = useState<AssistantAlert[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(false);
    const [loaded, setLoaded] = useState(false);
    const announcedRef = useRef<string | null>(null);
    const toastTimer = useRef<number | null>(null);
    const refreshTimer = useRef<number | null>(null);

    const dismissToast = useCallback(() => {
        setToast(null);
        if (toastTimer.current) window.clearTimeout(toastTimer.current);
        toastTimer.current = null;
    }, []);

    const refreshUnread = useCallback(
        async (announce: boolean) => {
            const result = await getAssistantAlertsUnreadAction(selectedBusinessId).catch(() => null);
            if (!result?.success) return;
            setUnread(result.data.count);
            const latest = result.data.latest;
            if (announce && latest && announcedRef.current !== latest.id) {
                announcedRef.current = latest.id;
                setToast(latest);
                if (toastTimer.current) window.clearTimeout(toastTimer.current);
                toastTimer.current = window.setTimeout(() => setToast(null), ALERT_TOAST_MS);
            }
        },
        [selectedBusinessId],
    );

    const load = useCallback(
        async (nextPage = 1) => {
            setLoading(true);
            const result = await getAssistantAlertsAction(selectedBusinessId, nextPage, PAGE_SIZE).catch(() => null);
            setLoading(false);
            if (!result?.success) return;
            setAlerts((prev) => (nextPage === 1 ? result.data.data : [...prev, ...result.data.data]));
            setTotal(result.data.total);
            setPage(nextPage);
            setLoaded(true);
        },
        [selectedBusinessId],
    );

    const markSeen = useCallback(async () => {
        setUnread(0);
        dismissToast();
        await markAssistantAlertsSeenAction(selectedBusinessId).catch(() => null);
    }, [selectedBusinessId, dismissToast]);

    useEffect(() => {
        setAlerts([]);
        setTotal(0);
        setPage(1);
        setLoaded(false);
        setUnread(0);
        announcedRef.current = null;
        dismissToast();
        if (!enabled) return;
        void refreshUnread(false);
    }, [enabled, selectedBusinessId, refreshUnread, dismissToast]);

    useSSE({
        enabled: enabled && !!selectedBusinessId,
        businessId: selectedBusinessId ?? undefined,
        eventTypes: ALERT_EVENT_TYPES,
        onMessage: (event) => {
            let type = '';
            try {
                const parsed = JSON.parse(event.data);
                type = parsed?.type || parsed?.metadata?.event_type || '';
            } catch {
                return;
            }
            if (!ALERT_EVENT_TYPES.includes(type)) return;
            if (refreshTimer.current) window.clearTimeout(refreshTimer.current);
            refreshTimer.current = window.setTimeout(() => {
                void refreshUnread(true);
                if (loaded) void load(1);
            }, ALERT_REFRESH_DELAY_MS);
        },
    });

    useEffect(
        () => () => {
            if (toastTimer.current) window.clearTimeout(toastTimer.current);
            if (refreshTimer.current) window.clearTimeout(refreshTimer.current);
        },
        [],
    );

    const hasMore = alerts.length < total;
    const loadMore = useCallback(() => {
        if (!loading && hasMore) void load(page + 1);
    }, [loading, hasMore, load, page]);

    return { unread, toast, dismissToast, alerts, total, loading, loaded, hasMore, load, loadMore, markSeen };
}

export type AssistantAlerts = ReturnType<typeof useAssistantAlerts>;
