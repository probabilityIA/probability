'use client';

import { useCallback, useEffect, useState } from 'react';
import { getMovementsAction } from '../../infra/actions';
import { GetMovementsParams, StockMovement } from '../../domain/types';
import { getActionError } from '@/shared/utils/action-result';

interface UseMovementsOptions extends GetMovementsParams {
    autoFetch?: boolean;
}

export function useMovements({ autoFetch = true, ...initialParams }: UseMovementsOptions = {}) {
    const [movements, setMovements] = useState<StockMovement[]>([]);
    const [loading, setLoading] = useState(autoFetch);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initialParams.page || 1);
    const [pageSize, setPageSize] = useState(initialParams.page_size || 20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [filters, setFilters] = useState<Omit<GetMovementsParams, 'page' | 'page_size'>>({
        product_id: initialParams.product_id,
        warehouse_id: initialParams.warehouse_id,
        type: initialParams.type,
        business_id: initialParams.business_id,
    });

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetMovementsParams = { page, page_size: pageSize, ...filters };
            const response = await getMovementsAction(params);
            setMovements(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar movimientos'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, filters]);

    useEffect(() => {
        if (autoFetch) fetch();
    }, [fetch, autoFetch]);

    return {
        movements,
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
