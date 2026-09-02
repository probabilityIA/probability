'use client';

import Link from 'next/link';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner } from '@/shared/ui';
import { TicketsManager } from '@/services/modules/tickets/ui';

export default function TicketsPage() {
    const { isSuperAdmin, isLoading } = usePermissions();

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <Spinner />
            </div>
        );
    }

    if (!isSuperAdmin) {
        return (
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full flex items-center justify-center px-4">
                <div className="max-w-md w-full rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-8 flex flex-col items-center gap-3 text-center">
                    <div className="w-12 h-12 rounded-full bg-red-100 dark:bg-red-900/40 flex items-center justify-center">
                        <svg className="w-6 h-6 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" strokeWidth={1.8} viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z" />
                        </svg>
                    </div>
                    <h1 className="text-lg font-semibold text-gray-900 dark:text-gray-100">No autorizado</h1>
                    <p className="text-sm text-gray-600 dark:text-gray-300">
                        {'El m\u00f3dulo de tickets est\u00e1 disponible solo para super administradores.'}
                    </p>
                    <Link
                        href="/home"
                        className="mt-2 inline-flex items-center h-9 px-4 rounded-lg text-sm font-semibold text-white"
                        style={{ backgroundColor: 'var(--color-primary)' }}
                    >
                        Volver al inicio
                    </Link>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <TicketsManager />
        </div>
    );
}
