'use client';

import { useCallback, useEffect, useState } from 'react';
import { getWarehouseTreeAction } from '../../infra/actions/hierarchy';
import { WarehouseTree } from '../../domain/hierarchy-types';
import { getActionError } from '@/shared/utils/action-result';

interface UseWarehouseTreeOptions {
    warehouseId: number;
    businessId?: number;
    autoFetch?: boolean;
}

export function useWarehouseTree({ warehouseId, businessId, autoFetch = true }: UseWarehouseTreeOptions) {
    const [tree, setTree] = useState<WarehouseTree | null>(null);
    const [loading, setLoading] = useState(autoFetch);
    const [error, setError] = useState<string | null>(null);

    const fetch = useCallback(async () => {
        if (!warehouseId) return;
        setLoading(true);
        setError(null);
        try {
            const result = await getWarehouseTreeAction(warehouseId, businessId);
            setTree(result);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar arbol de bodega'));
        } finally {
            setLoading(false);
        }
    }, [warehouseId, businessId]);

    useEffect(() => {
        if (autoFetch) fetch();
    }, [fetch, autoFetch]);

    return { tree, loading, error, refresh: fetch };
}
