'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { Spinner } from './spinner';

export function AccessRouteGuard({ children }: { children: React.ReactNode }) {
    const pathname = usePathname() ?? '/';
    const { isLoading, access, loadError, canAccessRoute, reloadAccess } = usePermissions();

    if (isLoading && !access) {
        return (
            <div className="flex min-h-[60vh] items-center justify-center">
                <Spinner size="lg" color="primary" text="Cargando permisos..." />
            </div>
        );
    }

    if (!access) {
        return (
            <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4 px-4 text-center">
                <h1 className="text-xl font-bold text-gray-900 dark:text-white">{'No pudimos cargar tus permisos'}</h1>
                <p className="text-sm text-gray-600 dark:text-gray-300">{loadError}</p>
                <button
                    type="button"
                    onClick={() => reloadAccess()}
                    className="rounded-lg bg-[#5b21b6] px-4 py-2 text-sm font-semibold text-white hover:bg-[#4c1d95]"
                >
                    Reintentar
                </button>
            </div>
        );
    }

    if (!canAccessRoute(pathname)) {
        return (
            <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4 px-4 text-center">
                <div className="text-5xl">{'\u{1F512}'}</div>
                <h1 className="text-xl font-bold text-gray-900 dark:text-white">{'No tienes acceso a este módulo'}</h1>
                <p className="max-w-md text-sm text-gray-600 dark:text-gray-300">
                    {'Tu rol o tu plan no incluyen esta sección. Si la necesitas, pídele acceso al administrador de tu negocio.'}
                </p>
                <Link
                    href="/home"
                    className="rounded-lg bg-[#5b21b6] px-4 py-2 text-sm font-semibold text-white hover:bg-[#4c1d95]"
                >
                    Ir al inicio
                </Link>
            </div>
        );
    }

    return <>{children}</>;
}
