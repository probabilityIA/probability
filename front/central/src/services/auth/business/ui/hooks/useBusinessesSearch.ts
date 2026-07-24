'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { getBusinessesSimpleAction } from '../../infra/actions';
import { BusinessSimple } from '../../domain/types';

const PAGE_SIZE = 20;
const MIN_SEARCH_LENGTH = 4;
const SEARCH_DEBOUNCE_MS = 300;

export const useBusinessesSearch = () => {
    const [businesses, setBusinesses] = useState<BusinessSimple[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [search, setSearch] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const requestId = useRef(0);

    const trimmed = search.trim();
    const effectiveSearch = trimmed.length >= MIN_SEARCH_LENGTH ? trimmed : '';

    const fetchPage = useCallback(async (targetPage: number, append: boolean) => {
        const id = ++requestId.current;
        setLoading(true);
        setError(null);
        try {
            const response = await getBusinessesSimpleAction({
                page: targetPage,
                page_size: PAGE_SIZE,
                search: effectiveSearch || undefined,
            });
            if (id !== requestId.current) return;
            if (response.success) {
                setBusinesses(prev => append ? [...prev, ...response.data] : response.data);
                setTotal(response.total ?? response.data.length);
                setPage(targetPage);
            } else {
                setError(response.message);
            }
        } catch (err: any) {
            if (id === requestId.current) {
                setError(err?.message ?? 'Error al obtener negocios');
            }
        } finally {
            if (id === requestId.current) {
                setLoading(false);
            }
        }
    }, [effectiveSearch]);

    useEffect(() => {
        const timer = setTimeout(() => fetchPage(1, false), SEARCH_DEBOUNCE_MS);
        return () => clearTimeout(timer);
    }, [fetchPage]);

    const hasMore = businesses.length < total;

    const loadMore = useCallback(() => {
        if (!loading && hasMore) {
            fetchPage(page + 1, true);
        }
    }, [loading, hasMore, page, fetchPage]);

    return {
        businesses,
        total,
        loading,
        error,
        search,
        setSearch,
        hasMore,
        loadMore,
        minSearchLength: MIN_SEARCH_LENGTH,
    };
};
