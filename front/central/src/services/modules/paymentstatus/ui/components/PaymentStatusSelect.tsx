'use client';

import { usePaymentStatuses } from '../hooks/usePaymentStatuses';

interface PaymentStatusSelectProps {
    value?: string;
    onChange: (code: string) => void;
    placeholder?: string;
    className?: string;
    disabled?: boolean;
    /** Si true incluye estados inactivos también (default: solo activos) */
    includeInactive?: boolean;
}

/**
 * Selector de estado de pago con punto de color.
 * Carga el catálogo automáticamente desde el backend.
 *
 * @example
 * ```tsx
 * <PaymentStatusSelect
 *   value={filters.payment_status}
 *   onChange={(code) => setFilters({ ...filters, payment_status: code })}
 *   placeholder="Todos los estados"
 * />
 * ```
 */
export default function PaymentStatusSelect({
    value,
    onChange,
    placeholder = 'Seleccionar estado de pago',
    className = '',
    disabled = false,
    includeInactive = false,
}: PaymentStatusSelectProps) {
    const { paymentStatuses, loading } = usePaymentStatuses(!includeInactive);

    // Dot de color del estado seleccionado (para mostrar junto al select)
    const selectedStatus = paymentStatuses.find((s) => s.code === value);

    return (
        <div className="relative flex items-center">
            {/* Dot de color del estado seleccionado */}
            {selectedStatus?.color && (
                <span
                    className="absolute left-3 w-2.5 h-2.5 rounded-full flex-shrink-0 pointer-events-none z-10"
                    style={{ backgroundColor: selectedStatus.color }}
                />
            )}

            <select
                value={value ?? ''}
                onChange={(e) => onChange(e.target.value)}
                disabled={disabled || loading}
                className={`block w-full rounded-md border border-gray-300 bg-white py-2 pr-3 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:bg-gray-50 disabled:text-gray-400 ${selectedStatus?.color ? 'pl-8' : 'pl-3'} ${className}`}
            >
                <option value="">{loading ? 'Cargando...' : placeholder}</option>
                {paymentStatuses.map((status) => (
                    <option key={status.id} value={status.code}>
                        {status.name}
                    </option>
                ))}
            </select>
        </div>
    );
}
