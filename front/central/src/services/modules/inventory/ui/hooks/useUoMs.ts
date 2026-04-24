'use client';

import { useCallback, useEffect, useState } from 'react';
import { listUoMsAction, listProductUoMsAction, convertUoMAction } from '../../infra/actions/traceability';
import { ConvertUoMInput, ConvertUoMResult, ProductUoM, UnitOfMeasure } from '../../domain/traceability-types';
import { getActionError } from '@/shared/utils/action-result';

export function useUoMs() {
    const [uoms, setUoMs] = useState<UnitOfMeasure[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await listUoMsAction();
            setUoMs(response.data || []);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar unidades de medida'));
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => { fetch(); }, [fetch]);

    return { uoms, loading, error, refresh: fetch };
}

export function useProductUoMs(productId: string | null, businessId?: number) {
    const [productUoMs, setProductUoMs] = useState<ProductUoM[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetch = useCallback(async () => {
        if (!productId) return;
        setLoading(true);
        setError(null);
        try {
            const response = await listProductUoMsAction(productId, businessId);
            setProductUoMs(response.data || []);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar UoM del producto'));
        } finally {
            setLoading(false);
        }
    }, [productId, businessId]);

    useEffect(() => { fetch(); }, [fetch]);

    return { productUoMs, loading, error, refresh: fetch };
}

interface UseUoMConverterState {
    result: ConvertUoMResult | null;
    loading: boolean;
    error: string | null;
}

export function useUoMConverter(businessId?: number) {
    const [state, setState] = useState<UseUoMConverterState>({ result: null, loading: false, error: null });

    const convert = useCallback(async (input: ConvertUoMInput) => {
        setState({ result: null, loading: true, error: null });
        const response = await convertUoMAction(input, businessId);
        if (response.success) {
            setState({ result: response.data, loading: false, error: null });
            return response.data;
        }
        setState({ result: null, loading: false, error: response.error });
        return null;
    }, [businessId]);

    return { ...state, convert };
}
