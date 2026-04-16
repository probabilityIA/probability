'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { CustomerInfo } from '../../../customers/domain/types';
import { getCustomersAction } from '../../../customers/infra/actions';

interface UseClientSearchOptions {
    businessId: number;
    debounceMs?: number;
    minChars?: number;
    pageSize?: number;
}

interface UseClientSearchResult {
    results: CustomerInfo[];
    loading: boolean;
    /** true after a search completes (regardless of results) */
    searched: boolean;
    search: (term: string) => void;
    clear: () => void;
}

export function useClientSearch({
    businessId,
    debounceMs = 300,
    minChars = 2,
    pageSize = 5,
}: UseClientSearchOptions): UseClientSearchResult {
    const [results, setResults] = useState<CustomerInfo[]>([]);
    const [loading, setLoading] = useState(false);
    const [searched, setSearched] = useState(false);
    const timerRef = useRef<NodeJS.Timeout | null>(null);
    const abortRef = useRef(0);

    const search = useCallback((term: string) => {
        if (timerRef.current) {
            clearTimeout(timerRef.current);
        }

        if (!term || term.length < minChars || !businessId) {
            setResults([]);
            setLoading(false);
            setSearched(false);
            return;
        }

        setLoading(true);
        setSearched(false);
        const requestId = ++abortRef.current;

        timerRef.current = setTimeout(async () => {
            try {
                const response = await getCustomersAction({
                    search: term,
                    business_id: businessId,
                    page: 1,
                    page_size: pageSize,
                });
                if (requestId === abortRef.current) {
                    setResults(response.data || []);
                    setSearched(true);
                }
            } catch {
                if (requestId === abortRef.current) {
                    setResults([]);
                    setSearched(true);
                }
            } finally {
                if (requestId === abortRef.current) {
                    setLoading(false);
                }
            }
        }, debounceMs);
    }, [businessId, debounceMs, minChars, pageSize]);

    const clear = useCallback(() => {
        if (timerRef.current) {
            clearTimeout(timerRef.current);
        }
        setResults([]);
        setLoading(false);
        setSearched(false);
    }, []);

    useEffect(() => {
        return () => {
            if (timerRef.current) {
                clearTimeout(timerRef.current);
            }
        };
    }, []);

    return { results, loading, searched, search, clear };
}
