'use client';

import { getStatusBadgeStyle } from '@/shared/utils/color-utils';
import { OrderStatusInfo } from '../../domain/types';

interface OrderStatusBadgeProps {
    status: OrderStatusInfo | undefined | null;
    fallback?: string;
    className?: string;
}

/**
 * Badge de estado de orden con el color definido en Probability.
 *
 * @example
 * ```tsx
 * <OrderStatusBadge status={mapping.order_status} />
 * ```
 */
export default function OrderStatusBadge({ status, fallback = '—', className = '' }: OrderStatusBadgeProps) {
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
