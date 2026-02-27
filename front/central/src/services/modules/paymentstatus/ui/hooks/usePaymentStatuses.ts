'use client';

import { useState, useEffect } from 'react';
import { getPaymentStatusesAction } from '../../infra/actions';
import { PaymentStatusInfo } from '../../domain/types';

/**
 * Hook para obtener el catálogo de estados de pago.
 *
 * Se usa para:
 * - Selectores/dropdowns de estado de pago en formularios
 * - Filtros por estado de pago
 * - Mostrar el nombre/color de un estado a partir de su code
 *
 * @param isActive - Si true, solo retorna estados activos (default: true)
 *
 * @example
 * ```tsx
 * const { paymentStatuses, loading } = usePaymentStatuses();
 *
 * const status = paymentStatuses.find(s => s.code === order.payment_status);
 * return <PaymentStatusBadge status={status} />;
 * ```
 */
export function usePaymentStatuses(isActive: boolean = true) {
    const [paymentStatuses, setPaymentStatuses] = useState<PaymentStatusInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchPaymentStatuses = async () => {
        setLoading(true);
        setError(null);

        try {
            const result = await getPaymentStatusesAction(isActive !== undefined ? { is_active: isActive } : undefined);

            if (result.success) {
                setPaymentStatuses(result.data);
            } else {
                setError(result.message || 'Error al cargar estados de pago');
                setPaymentStatuses([]);
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar estados de pago');
            setPaymentStatuses([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchPaymentStatuses();
    }, [isActive]);

    return {
        paymentStatuses,
        loading,
        error,
        refresh: fetchPaymentStatuses,
    };
}
