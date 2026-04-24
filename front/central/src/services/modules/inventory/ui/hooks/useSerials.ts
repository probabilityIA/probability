'use client';

import { useCallback, useEffect, useState } from 'react';
import { listSerialsAction } from '../../infra/actions/traceability';
import { GetSerialsParams, InventorySerial } from '../../domain/traceability-types';
import { getActionError } from '@/shared/utils/action-result';

interface UseSerialsOptions extends GetSerialsParams {
    autoFetch?: boolean;
}

export function useSerials({ autoFetch = true, ...initialParams }: UseSerialsOptions = {}) {
    const [serials, setSerials] = useState<InventorySerial[]>([]);
    const [loading, setLoading] = useState(autoFetch);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initialParams.page || 1);
    const [pageSize, setPageSize] = useState(initialParams.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [filters, setFilters] = useState<Omit<GetSerialsParams, 'page' | 'page_size'>>({
        product_id: initialParams.product_id,
        lot_id: initialParams.lot_id,
        state_id: initialParams.state_id,
        location_id: initialParams.location_id,
        business_id: initialParams.business_id,
    });

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await listSerialsAction({ page, page_size: pageSize, ...filters });
            setSerials(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar series'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, filters]);

    useEffect(() => {
        if (autoFetch) fetch();
    }, [fetch, autoFetch]);

    return {
        serials,
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
