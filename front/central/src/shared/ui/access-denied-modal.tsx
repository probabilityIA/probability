'use client';

import { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { useRouter } from 'next/navigation';
import { TokenStorage } from '@/shared/config';
import type { AccessDeniedInfo } from '@/shared/utils/access-denied';

interface AccessDeniedModalProps {
    info: AccessDeniedInfo | null;
    onClose: () => void;
}

const TITLES: Record<AccessDeniedInfo['code'], string> = {
    forbidden: 'No tienes permiso',
    subscription_suspended: 'Suscripción vencida',
    unauthenticated: 'Tu sesión expiró',
};

const ICONS: Record<AccessDeniedInfo['code'], string> = {
    forbidden: '\u{1F512}',
    subscription_suspended: '\u{1F4B3}',
    unauthenticated: '\u{23F0}',
};

export function AccessDeniedModal({ info, onClose }: AccessDeniedModalProps) {
    const router = useRouter();
    const [mounted, setMounted] = useState(false);

    useEffect(() => setMounted(true), []);

    useEffect(() => {
        if (!info) return;
        const onKey = (e: KeyboardEvent) => {
            if (e.key === 'Escape' && info.code !== 'unauthenticated') onClose();
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, [info, onClose]);

    if (!mounted || !info) return null;

    const goToLogin = () => {
        TokenStorage.clearSession();
        onClose();
        router.push('/login');
    };

    const goToSubscription = () => {
        onClose();
        router.push('/subscription');
    };

    return createPortal(
        <div
            role="dialog"
            aria-modal="true"
            aria-labelledby="access-denied-title"
            className="fixed inset-0 flex items-center justify-center bg-black/50 px-4"
            style={{ zIndex: 1100 }}
            onClick={info.code === 'unauthenticated' ? undefined : onClose}
        >
            <div
                className="w-full max-w-md rounded-2xl bg-white p-6 text-center shadow-xl dark:bg-gray-800"
                onClick={(e) => e.stopPropagation()}
            >
                <div className="mb-3 text-5xl">{ICONS[info.code]}</div>
                <h2 id="access-denied-title" className="mb-2 text-xl font-bold text-gray-900 dark:text-white">
                    {TITLES[info.code]}
                </h2>
                <p className="mb-2 text-sm text-gray-700 dark:text-gray-300">{info.message}</p>
                {info.code === 'forbidden' && (
                    <p className="mb-5 text-xs text-gray-500 dark:text-gray-400">
                        {'Si necesitas hacer esto, pídele acceso al administrador de tu negocio.'}
                    </p>
                )}
                <div className="mt-5 flex justify-center gap-3">
                    {info.code === 'unauthenticated' && (
                        <button
                            type="button"
                            onClick={goToLogin}
                            className="rounded-lg bg-[#5b21b6] px-4 py-2 text-sm font-semibold text-white hover:bg-[#4c1d95]"
                        >
                            {'Iniciar sesión'}
                        </button>
                    )}
                    {info.code === 'subscription_suspended' && (
                        <button
                            type="button"
                            onClick={goToSubscription}
                            className="rounded-lg bg-[#5b21b6] px-4 py-2 text-sm font-semibold text-white hover:bg-[#4c1d95]"
                        >
                            {'Ir a Suscripción'}
                        </button>
                    )}
                    {info.code !== 'unauthenticated' && (
                        <button
                            type="button"
                            onClick={onClose}
                            className="rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-700"
                        >
                            Entendido
                        </button>
                    )}
                </div>
            </div>
        </div>,
        document.body,
    );
}
