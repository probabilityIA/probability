'use client';

import { useState, useMemo } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { useWalletBusiness } from '@/shared/contexts/wallet-business-context';
import ShippingProfitMonthly from '@/services/modules/shipping-margins/ui/components/ShippingProfitMonthly';
import { BusinessWalletFinancialSummary } from '@/services/modules/wallet/ui/components';
import { GlobalFinancialSummaryCard } from '@/services/modules/accounting/ui';

export default function WalletFinanzasPage() {
    const { isSuperAdmin } = usePermissions();
    const { businesses } = useBusinessesSimple();
    const { selectedBusinessId: contextBusinessId } = useWalletBusiness();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    const userBusinessId = useMemo(() => {
        if (businesses && businesses.length > 0) {
            return businesses[0].id;
        }
        return null;
    }, [businesses]);

    const contextBusiness = businesses.find((b) => b.id === contextBusinessId);

    if (!isSuperAdmin) {
        return (
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">Mi Resumen Financiero</h2>
                <BusinessWalletFinancialSummary />
            </div>
        );
    }

    if (contextBusinessId) {
        return (
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
                <div className="flex items-center gap-3 bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-700 rounded-lg px-4 py-3 mb-6">
                    <p className="text-sm text-blue-800 dark:text-blue-300">
                        Viendo lo mismo que ve <strong>{contextBusiness?.name || `negocio #${contextBusinessId}`}</strong> en su
                        Resumen Financiero — sin margen ni ganancia de Probability.
                    </p>
                </div>
                <BusinessWalletFinancialSummary businessId={contextBusinessId} />
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="mb-6 bg-white dark:bg-gray-800 rounded-lg shadow-sm dark:shadow-lg border border-gray-200 dark:border-gray-700 p-6">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Ver finanzas de un negocio especifico (opcional)
                </label>
                <select
                    value={selectedBusinessId || ''}
                    onChange={(e) => setSelectedBusinessId(e.target.value ? parseInt(e.target.value) : null)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                    <option value="">Todos los negocios (resumen de Probability)</option>
                    {businesses.map(business => (
                        <option key={business.id} value={business.id}>
                            {business.name}
                        </option>
                    ))}
                </select>
            </div>

            {selectedBusinessId ? (
                <ShippingProfitMonthly selectedBusinessId={selectedBusinessId} />
            ) : (
                <GlobalFinancialSummaryCard />
            )}
        </div>
    );
}
