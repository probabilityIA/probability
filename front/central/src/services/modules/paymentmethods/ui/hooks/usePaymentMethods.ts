'use client';

import { useCallback, useEffect, useState } from 'react';
import { getPaymentMethodsAction } from '../../infra/actions';
import { PaymentMethodInfo } from '../../domain/types';
import { getActionError } from '@/shared/utils/action-result';

export function usePaymentMethods() {
    const [paymentMethods, setPaymentMethods] = useState<PaymentMethodInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchPaymentMethods = useCallback(async () => {
        setLoading(true);
        setError(null);

        try {
            const result = await getPaymentMethodsAction();

            if (result.success) {
                setPaymentMethods(result.data || []);
            } else {
                setError(result.message || 'Error al cargar los medios de pago');
                setPaymentMethods([]);
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar los medios de pago'));
            setPaymentMethods([]);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchPaymentMethods();
    }, [fetchPaymentMethods]);

    return {
        paymentMethods,
        loading,
        error,
        refresh: fetchPaymentMethods,
    };
}
