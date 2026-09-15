'use client';

import { useCallback, useEffect, useState } from 'react';
import { getAssistantReviewMessagesAction, getAssistantReviewSummaryAction } from '../../infra/actions';
import { reviewRange } from '../../app/use-cases';
import type { PaginatedResponse, ReviewFilters, ReviewKind, ReviewMessage, ReviewSummary } from '../../domain/types';

const PAGE_SIZE = 20;
const SEARCH_DEBOUNCE_MS = 400;

export function useAssistantReview() {
    const [days, setDays] = useState(30);
    const [businessId, setBusinessId] = useState<number | undefined>(undefined);
    const [kind, setKind] = useState<ReviewKind>('all');
    const [searchInput, setSearchInput] = useState('');
    const [search, setSearch] = useState('');
    const [page, setPage] = useState(1);

    const [messages, setMessages] = useState<PaginatedResponse<ReviewMessage> | null>(null);
    const [summary, setSummary] = useState<ReviewSummary | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const timer = window.setTimeout(() => {
            setSearch(searchInput.trim());
            setPage(1);
        }, SEARCH_DEBOUNCE_MS);
        return () => window.clearTimeout(timer);
    }, [searchInput]);

    const load = useCallback(async () => {
        setLoading(true);
        setError(null);
        const base: ReviewFilters = { ...reviewRange(days), business_id: businessId, search: search || undefined };
        const [list, totals] = await Promise.all([
            getAssistantReviewMessagesAction({ ...base, kind, page, page_size: PAGE_SIZE }),
            getAssistantReviewSummaryAction(base),
        ]);
        if (list.success) setMessages(list.data);
        if (totals.success) setSummary(totals.data);
        if (!list.success) setError(list.message);
        else if (!totals.success) setError(totals.message);
        setLoading(false);
    }, [businessId, days, kind, page, search]);

    useEffect(() => {
        void load();
    }, [load]);

    const changeDays = (value: number) => {
        setDays(value);
        setPage(1);
    };
    const changeBusiness = (value: number | undefined) => {
        setBusinessId(value);
        setPage(1);
    };
    const changeKind = (value: ReviewKind) => {
        setKind(value);
        setPage(1);
    };

    return {
        days,
        businessId,
        kind,
        searchInput,
        page,
        messages,
        summary,
        loading,
        error,
        setPage,
        setSearchInput,
        changeDays,
        changeBusiness,
        changeKind,
        reload: load,
    };
}
