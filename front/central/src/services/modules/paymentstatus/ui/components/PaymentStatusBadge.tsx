'use client';

import { getStatusBadgeStyle } from '@/shared/utils/color-utils';
import { PaymentStatusInfo } from '../../domain/types';

interface PaymentStatusBadgeProps {
    status: PaymentStatusInfo | undefined | null;
    fallback?: string;
    className?: string;
}

/**
 * Badge de estado de pago con el color definido en Probability.
 *
 * @example
 * ```tsx
 * const { paymentStatuses } = usePaymentStatuses();
 * const status = paymentStatuses.find(s => s.code === order.payment_status);
 * <PaymentStatusBadge status={status} />
 * ```
 */
export default function PaymentStatusBadge({ status, fallback = '—', className = '' }: PaymentStatusBadgeProps) {
    if (!status) {
        return (
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-500 ${className}`}>
                {fallback}
            </span>
        );
    }

    return (
        <span
            className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium ${className}`}
            style={getStatusBadgeStyle(status.color)}
            title={status.description}
        >
            <span
                className="w-1.5 h-1.5 rounded-full flex-shrink-0"
                style={{ backgroundColor: status.color ?? '#9CA3AF' }}
            />
            {status.name}
        </span>
    );
}
