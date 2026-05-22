'use client';

import { FinancialStatsView } from '../financial-stats';

export default function WalletFinanzasPage() {
    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm dark:shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                <div className="p-6">
                    <FinancialStatsView />
                </div>
            </div>
        </div>
    );
}
