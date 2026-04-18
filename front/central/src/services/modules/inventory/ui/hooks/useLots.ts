'use client';

import { useCallback, useEffect, useState } from 'react';
import { listLotsAction } from '../../infra/actions/traceability';
import { GetLotsParams, InventoryLot } from '../../domain/traceability-types';
import { getActionError } from '@/shared/utils/action-result';

interface UseLotsOptions extends GetLotsParams {
    autoFetch?: boolean;
}

export function useLots({ autoFetch = true, ...initialParams }: UseLotsOptions = {}) {
    const [lots, setLots] = useState<InventoryLot[]>([]);
    const [loading, setLoading] = useState(autoFetch);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initialParams.page || 1);
    const [pageSize, setPageSize] = useState(initialParams.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [filters, setFilters] = useState<Omit<GetLotsParams, 'page' | 'page_size'>>({
        product_id: initialParams.product_id,
        status: initialParams.status,
        expiring_in_days: initialParams.expiring_in_days,
        business_id: initialParams.business_id,
    });

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await listLotsAction({ page, page_size: pageSize, ...filters });
            setLots(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar lotes'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, filters]);

    useEffect(() => {
        if (autoFetch) fetch();
    }, [fetch, autoFetch]);

    return {
        lots,
        loading,
        error,
        page,
        pageSize,
        total,
        totalPages,
        filters,
        setPage,
        setPageSize: (size: number) => { setPageSize(size); setPage(1); },
        setFilters: (next: Partial<typeof filters>) => { setFilters((f) => ({ ...f, ...next })); setPage(1); },
        refresh: fetch,
    };
}
