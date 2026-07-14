'use client';

import { useEffect } from 'react';
import { usePaymentMethods } from '../hooks/usePaymentMethods';

interface PaymentMethodSelectProps {
    value?: number | null;
    onChange: (id: number) => void;
    defaultCode?: string;
    placeholder?: string;
    className?: string;
    disabled?: boolean;
    hasError?: boolean;
}

const BASE_CLASS =
    'w-full px-3 py-2.5 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-600 focus:border-transparent bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100';

export default function PaymentMethodSelect({
    value,
    onChange,
    defaultCode,
    placeholder = 'Selecciona un medio de pago...',
    className = '',
    disabled = false,
    hasError = false,
}: PaymentMethodSelectProps) {
    const { paymentMethods, loading } = usePaymentMethods();

    useEffect(() => {
        if (loading || !defaultCode || (value && value > 0)) return;
        const match = paymentMethods.find((m) => m.code === defaultCode);
        if (match) onChange(match.id);
    }, [loading, defaultCode, value, paymentMethods, onChange]);

    return (
        <select
            value={value && value > 0 ? String(value) : ''}
            onChange={(e) => onChange(parseInt(e.target.value, 10) || 0)}
            disabled={disabled || loading}
            className={`${BASE_CLASS} ${hasError ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'} ${className}`}
        >
            <option value="">{loading ? 'Cargando medios de pago...' : placeholder}</option>
            {paymentMethods.map((method) => (
                <option key={method.id} value={method.id}>
                    {method.name}
                </option>
            ))}
        </select>
    );
}
