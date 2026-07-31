'use client';

import Link from 'next/link';
import { useCallback, useEffect, useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { getTiendaSessionAction, tiendaLogoutAction } from '../../infra/actions';
import { TiendaSession } from '../../domain/types';

interface TiendaAccountButtonProps {
    slug: string;
}

export function TiendaAccountButton({ slug }: TiendaAccountButtonProps) {
    const router = useRouter();
    const pathname = usePathname();
    const [session, setSession] = useState<TiendaSession | null>(null);
    const [loaded, setLoaded] = useState(false);
    const [open, setOpen] = useState(false);

    const refresh = useCallback(() => {
        getTiendaSessionAction(slug).then(s => {
            setSession(s);
            setLoaded(true);
        });
    }, [slug]);

    useEffect(() => {
        refresh();
    }, [refresh, pathname]);

    useEffect(() => {
        const handler = () => refresh();
        window.addEventListener('tienda-session-changed', handler);
        return () => window.removeEventListener('tienda-session-changed', handler);
    }, [refresh]);

    const handleLogout = async () => {
        await tiendaLogoutAction();
        setSession(null);
        setOpen(false);
        window.dispatchEvent(new Event('tienda-session-changed'));
        router.refresh();
    };

    if (!loaded) return <span className="w-20" />;

    if (!session) {
        return (
            <Link
                href={`/tienda/${slug}/cuenta`}
                className="text-sm font-medium opacity-90 hover:opacity-100 flex items-center gap-1"
                style={{ color: 'var(--brand-nav-text, #4b5563)' }}
            >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
                Iniciar sesion
            </Link>
        );
    }

    return (
        <div className="relative">
            <button
                type="button"
                onClick={() => setOpen(!open)}
                className="text-sm font-medium opacity-90 hover:opacity-100 flex items-center gap-1"
                style={{ color: 'var(--brand-nav-text, #4b5563)' }}
            >
                {session.avatar_url ? (
                    <img src={session.avatar_url} alt={session.name} className="w-6 h-6 rounded-full object-cover" />
                ) : (
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                )}
                {session.name.split(' ')[0]}
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
            </button>
            {open && (
                <div className="absolute right-0 mt-2 w-44 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-100 dark:border-gray-700 py-1 z-50">
                    <Link
                        href={`/tienda/${slug}/cuenta/perfil`}
                        className="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700"
                        onClick={() => setOpen(false)}
                    >
                        Mi perfil
                    </Link>
                    <Link
                        href={`/tienda/${slug}/cuenta/pedidos`}
                        className="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700"
                        onClick={() => setOpen(false)}
                    >
                        Mis pedidos
                    </Link>
                    <button
                        type="button"
                        onClick={handleLogout}
                        className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-50 dark:hover:bg-gray-700"
                    >
                        Cerrar sesion
                    </button>
                </div>
            )}
        </div>
    );
}
