'use client';

import { useCallback, useState } from 'react';
import { validateCubingAction } from '../../infra/actions/hierarchy';
import { CubingCheckResult, ValidateCubingInput } from '../../domain/hierarchy-types';

interface UseCubingState {
    result: CubingCheckResult | null;
    loading: boolean;
    error: string | null;
}

export function useCubing(businessId?: number) {
    const [state, setState] = useState<UseCubingState>({ result: null, loading: false, error: null });

    const validate = useCallback(async (input: ValidateCubingInput) => {
        setState({ result: null, loading: true, error: null });
        const response = await validateCubingAction(input, businessId);
        if (response.success) {
            setState({ result: response.data, loading: false, error: null });
            return response.data;
        }
        setState({ result: null, loading: false, error: response.error });
        return null;
    }, [businessId]);

    const reset = useCallback(() => setState({ result: null, loading: false, error: null }), []);

    return { ...state, validate, reset };
}
