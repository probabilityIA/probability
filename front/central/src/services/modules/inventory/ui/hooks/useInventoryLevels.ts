'use client';

import { useCallback, useEffect, useState } from 'react';
import { getWarehouseInventoryAction } from '../../infra/actions';
import { GetInventoryParams, InventoryLevel } from '../../domain/types';
import { getActionError } from '@/shared/utils/action-result';

interface UseInventoryLevelsOptions {
    warehouseId: number;
    businessId?: number;
    pageSize?: number;
    autoFetch?: boolean;
}

interface UseInventoryLevelsState {
    levels: InventoryLevel[];
    loading: boolean;
    error: string | null;
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
    search: string;
    lowStockOnly: boolean;
}

export function useInventoryLevels({
    warehouseId,
    businessId,
    pageSize: initialPageSize = 20,
    autoFetch = true,
}: UseInventoryLevelsOptions) {
    const [state, setState] = useState<UseInventoryLevelsState>({
        levels: [],
        loading: autoFetch,
        error: null,
        page: 1,
        pageSize: initialPageSize,
        total: 0,
        totalPages: 1,
        search: '',
        lowStockOnly: false,
    });

    const fetchLevels = useCallback(async () => {
        setState((s) => ({ ...s, loading: true, error: null }));
        try {
            const params: GetInventoryParams = {
                page: state.page,
                page_size: state.pageSize,
            };
            if (state.search) params.search = state.search;
            if (state.lowStockOnly) params.low_stock = true;
            if (businessId) params.business_id = businessId;

            const response = await getWarehouseInventoryAction(warehouseId, params);
            setState((s) => ({
                ...s,
                levels: response.data || [],
                total: response.total || 0,
                totalPages: response.total_pages || 1,
                page: response.page || s.page,
                loading: false,
            }));
        } catch (err: any) {
            setState((s) => ({
                ...s,
                loading: false,
                error: getActionError(err, 'Error al cargar inventario'),
            }));
        }
    }, [warehouseId, businessId, state.page, state.pageSize, state.search, state.lowStockOnly]);

    useEffect(() => {
        if (autoFetch) fetchLevels();
    }, [fetchLevels, autoFetch]);

    useEffect(() => {
        setState((s) => ({ ...s, page: 1, search: '' }));
    }, [warehouseId, businessId]);

    return {
        ...state,
        setPage: (page: number) => setState((s) => ({ ...s, page })),
        setPageSize: (pageSize: number) => setState((s) => ({ ...s, pageSize, page: 1 })),
        setSearch: (search: string) => setState((s) => ({ ...s, search, page: 1 })),
        setLowStockOnly: (lowStockOnly: boolean) => setState((s) => ({ ...s, lowStockOnly, page: 1 })),
        refresh: fetchLevels,
    };
}
