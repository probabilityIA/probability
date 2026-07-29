'use client';

import { useState, useEffect, useCallback } from 'react';
import { MarketingLead } from '../../domain/types';
import { getMarketingLeadsAction } from '../../infra/actions';
import { getActionError } from '@/shared/utils/action-result';

export const useMarketingLeads = () => {
    const [leads, setLeads] = useState<MarketingLead[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);

    const loadLeads = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getMarketingLeadsAction({ page, page_size: pageSize });
            if (response.success && response.data) {
                setLeads(response.data.leads || []);
                setTotal(response.data.total || 0);
                setTotalPages(response.data.total_pages || 1);
            } else {
                setError(response.message || 'Error al cargar los leads');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar los leads'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize]);

    useEffect(() => {
        loadLeads();
    }, [loadLeads]);

    return {
        leads,
        loading,
        error,
        page,
        setPage,
        pageSize,
        setPageSize,
        totalPages,
        total,
        refresh: loadLeads,
        setError,
    };
};
