'use client';

import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';

interface SuperAdminBusinessSelectorProps {
    value: number | null;
    onChange: (businessId: number | null) => void;
    /** Variante visual: 'navbar' para uso compacto en barras, 'default' para uso en contenido */
    variant?: 'navbar' | 'default';
    /** Texto del option por defecto (sin selección) */
    placeholder?: string;
}

export function SuperAdminBusinessSelector({
    value,
    onChange,
    variant = 'default',
    placeholder = 'Todos los negocios',
}: SuperAdminBusinessSelectorProps) {
    const { isSuperAdmin } = usePermissions();
    const { businesses } = useBusinessesSimple();

    if (!isSuperAdmin || businesses.length === 0) {
        return null;
    }

    const options = (
        <>
            <option value="">{placeholder}</option>
            {businesses.map((b) => (
                <option key={b.id} value={b.id}>{b.name}</option>
            ))}
        </>
    );

    const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const val = e.target.value;
        onChange(val ? Number(val) : null);
    };

    if (variant === 'navbar') {
        return (
            <div className="flex items-center gap-2 bg-purple-100 dark:bg-purple-900/30 border border-purple-300 dark:border-purple-700 rounded-lg px-3 py-1.5">
                <span className="px-2 py-0.5 text-xs font-bold text-white bg-purple-700 rounded select-none whitespace-nowrap">
                    SUPER ADMIN
                </span>
                <select
                    value={value?.toString() ?? ''}
                    onChange={handleChange}
                    className="px-2 py-1.5 border border-purple-400 dark:border-purple-600 rounded-md text-sm font-medium focus:outline-none focus:ring-2 focus:ring-purple-600 bg-white dark:bg-gray-800 text-purple-900 dark:text-purple-200 cursor-pointer"
                >
                    {options}
                </select>
            </div>
        );
    }

    return (
        <div className="flex items-center gap-3 bg-purple-50 dark:bg-purple-900/20 border-2 border-purple-300 dark:border-purple-700 rounded-lg px-4 py-3">
            <span className="px-2.5 py-1 text-xs font-bold text-white bg-purple-700 rounded-md select-none whitespace-nowrap">
                SUPER ADMIN
            </span>
            <select
                value={value?.toString() ?? ''}
                onChange={handleChange}
                className="flex-1 max-w-xs px-3 py-2 border-2 border-purple-400 dark:border-purple-600 rounded-lg text-sm font-medium focus:outline-none focus:ring-2 focus:ring-purple-600 bg-white dark:bg-gray-800 text-purple-900 dark:text-purple-200 cursor-pointer"
            >
                {options}
            </select>
        </div>
    );
}
