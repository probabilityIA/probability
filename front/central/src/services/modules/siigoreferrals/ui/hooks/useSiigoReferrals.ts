'use client';

import { useState, useEffect, useCallback } from 'react';
import { SiigoReferral } from '../../domain/types';
import { getSiigoReferralsAction } from '../../infra/actions';
import { getActionError } from '@/shared/utils/action-result';

export const useSiigoReferrals = () => {
    const [referrals, setReferrals] = useState<SiigoReferral[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);

    const loadReferrals = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getSiigoReferralsAction({ page, page_size: pageSize });
            if (response.success && response.data) {
                setReferrals(response.data.referrals || []);
                setTotal(response.data.total || 0);
                setTotalPages(response.data.total_pages || 1);
            } else {
                setError(response.message || 'Error al cargar los referidos');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar los referidos'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize]);

    useEffect(() => {
        loadReferrals();
    }, [loadReferrals]);

    return {
        referrals,
        loading,
        error,
        page,
        setPage,
        pageSize,
        setPageSize,
        totalPages,
        total,
        refresh: loadReferrals,
        setError,
    };
};
