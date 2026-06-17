'use client';

import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { useWalletBusiness } from '@/shared/contexts/wallet-business-context';
import { AdminWalletView, BusinessWalletView } from '../wallet-views';

export default function WalletSaldosPage() {
    const { isSuperAdmin } = usePermissions();
    const { businesses } = useBusinessesSimple();
    const { selectedBusinessId } = useWalletBusiness();

    const selectedBusiness = businesses.find(b => b.id === selectedBusinessId);

    console.log('WalletSaldosPage render:', { isSuperAdmin, selectedBusinessId, selectedBusiness: selectedBusiness?.name });

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm dark:shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                <div className="p-6">
                    {isSuperAdmin ? (
                        selectedBusinessId ? (
                            <BusinessWalletView
                                businessId={selectedBusinessId}
                                businessName={selectedBusiness?.name}
                            />
                        ) : (
                            <AdminWalletView />
                        )
                    ) : (
                        <BusinessWalletView />
                    )}
                </div>
            </div>
        </div>
    );
}
