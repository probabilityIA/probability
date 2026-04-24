'use client';

import { useCallback, useEffect, useState } from 'react';
import { getMovementTypesAction } from '../../infra/actions';
import { MovementType } from '../../domain/types';
import { getActionError } from '@/shared/utils/action-result';

interface UseMovementTypesOptions {
    activeOnly?: boolean;
    businessId?: number;
    pageSize?: number;
    autoFetch?: boolean;
}

export function useMovementTypes({
    activeOnly = true,
    businessId,
    pageSize = 50,
    autoFetch = true,
}: UseMovementTypesOptions = {}) {
    const [types, setTypes] = useState<MovementType[]>([]);
    const [loading, setLoading] = useState(autoFetch);
    const [error, setError] = useState<string | null>(null);

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getMovementTypesAction({
                page: 1,
                page_size: pageSize,
                active_only: activeOnly,
                business_id: businessId,
            });
            setTypes(response.data || []);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar tipos de movimiento'));
        } finally {
            setLoading(false);
        }
    }, [activeOnly, businessId, pageSize]);

    useEffect(() => {
        if (autoFetch) fetch();
    }, [fetch, autoFetch]);

    return { types, loading, error, refresh: fetch };
}
