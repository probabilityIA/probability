'use client';

import { useEffect, useState } from 'react';
import { usePathname } from 'next/navigation';
import Link from 'next/link';
import { usePermissions } from '@/shared/contexts/permissions-context';

const SUBSCRIPTION_GUARD_ENABLED = true;

export function SubscriptionGuard({ children }: { children: React.ReactNode }) {
    const { permissions, isSuperAdmin } = usePermissions();
    const pathname = usePathname();

    if (!SUBSCRIPTION_GUARD_ENABLED) return <>{children}</>;
    if (isSuperAdmin) return <>{children}</>;

    const isExpired = permissions?.subscription_status === 'expired' ||
        permissions?.subscription_status === 'cancelled';

    if (!isExpired) return <>{children}</>;
    if (pathname?.startsWith('/subscription')) return <>{children}</>;

    return <SuspendedWall />;
}

function SuspendedWall() {
    const [tick, setTick] = useState(0);

    useEffect(() => {
        const id = setInterval(() => setTick((t) => t + 1), 1000);
        return () => clearInterval(id);
    }, []);

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900 px-4">
            <div className="max-w-md w-full text-center">
                <div
                    className="text-7xl mb-6 transition-transform duration-300"
                    style={{ transform: tick % 2 === 0 ? 'scale(1)' : 'scale(1.07)' }}
                >
                    🔒
                </div>

                <h1 className="text-2xl font-bold text-gray-900 dark:text-white dark:text-white mb-3">
                    Cuenta suspendida
                </h1>

                <p className="text-gray-600 dark:text-gray-300 dark:text-gray-400 mb-6 leading-relaxed">
                    No puedes usar las funciones de la plataforma hasta que no te pongas al día con tu pago.
                    Ingresa al módulo de <strong>Suscripción</strong> para ver los detalles y la información de pago.
                </p>

                <div className="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-700 rounded-xl p-4 mb-6 text-left">
                    <p className="text-sm font-semibold text-amber-800 dark:text-amber-300 mb-2">
                        ¿Qué puedes hacer?
                    </p>
                    <ul className="text-sm text-amber-700 dark:text-amber-400 space-y-1">
                        <li>✅ Ver información de pago en el módulo de Suscripción</li>
                        <li>✅ Contactar a tu asesor</li>
                        <li>❌ Crear órdenes o cotizaciones</li>
                        <li>❌ Recargar billetera o crear guías</li>
                    </ul>
                </div>

                <Link
                    href="/subscription"
                    className="inline-flex items-center gap-2 px-6 py-3 bg-violet-600 hover:bg-violet-700 text-white font-semibold rounded-xl transition-all shadow-lg hover:shadow-violet-200 dark:hover:shadow-violet-900"
                >
                    <span>💳</span>
                    <span>Ir a mi Suscripción</span>
                </Link>
            </div>
        </div>
    );
}
