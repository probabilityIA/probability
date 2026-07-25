'use client';

import { ReactNode } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner } from '@/shared/ui';

export function AccountingGate({ children }: { children: ReactNode }) {
    const { isSuperAdmin, isLoading } = usePermissions();

    if (isLoading) {
        return (
            <div className="min-h-[50vh] flex items-center justify-center">
                <Spinner size="lg" color="primary" text="Cargando..." />
            </div>
        );
    }

    if (!isSuperAdmin) {
        return (
            <div className="min-h-[50vh] flex items-center justify-center px-4">
                <div className="max-w-md w-full text-center bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm p-8">
                    <div className="mx-auto w-12 h-12 rounded-full bg-purple-100 dark:bg-purple-900/40 flex items-center justify-center mb-4">
                        <svg className="w-6 h-6 text-purple-600 dark:text-purple-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                        </svg>
                    </div>
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Acceso restringido</h2>
                    <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                        La contabilidad interna de la plataforma solo esta disponible para super administradores.
                    </p>
                </div>
            </div>
        );
    }

    return <>{children}</>;
}
