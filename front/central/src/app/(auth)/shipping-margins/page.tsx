'use client';

import { useOrdersBusiness } from '@/shared/contexts/orders-business-context';
import { ShippingMarginManager } from '@/services/modules/shipping-margins/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

export default function ShippingMarginsPage() {
    const { selectedBusinessId } = useOrdersBusiness();
    const { isSuperAdmin, isLoading } = usePermissions();

    if (isLoading) {
        return null;
    }

    if (!isSuperAdmin) {
        return (
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-12 flex flex-col items-center justify-center text-center">
                <svg className="w-16 h-16 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 15v2m0 0v2m0-2h2m-2 0h-2m6.364-1.636A9 9 0 105.636 5.636a9 9 0 0012.728 12.728z" />
                </svg>
                <h2 className="text-xl font-semibold text-gray-700 dark:text-gray-200">Acceso restringido</h2>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-2 max-w-md">
                    La configuracion de margenes de envio esta disponible solo para super administradores.
                </p>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <ShippingMarginManager selectedBusinessId={selectedBusinessId} />
        </div>
    );
}
