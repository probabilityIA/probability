'use client';

import { useEffect, useState } from 'react';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { MyIntegrationsModal } from './MyIntegrationsModal';
import { OPEN_INTEGRATIONS_HUB_EVENT, consumePendingIntegrationsHub } from '../open-hub';
import type { SyncEnvironment } from '../sync-activity-context';

interface MyIntegrationsButtonProps {
    businessId?: number | null;
}

export function MyIntegrationsButton({ businessId }: MyIntegrationsButtonProps) {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [initialEnvironment, setInitialEnvironment] = useState<SyncEnvironment | null>(null);
    const { isSuperAdmin, can } = usePermissions();
    const pathname = usePathname();

    const allowed = can('integrations.read');
    const disabled = isSuperAdmin && !businessId;

    useEffect(() => {
        if (!allowed || disabled) return;
        const openPending = () => {
            const intent = consumePendingIntegrationsHub();
            if (!intent) return;
            setInitialEnvironment(intent.environment);
            setIsModalOpen(true);
        };
        openPending();
        window.addEventListener(OPEN_INTEGRATIONS_HUB_EVENT, openPending);
        return () => window.removeEventListener(OPEN_INTEGRATIONS_HUB_EVENT, openPending);
    }, [allowed, disabled, pathname]);

    if (!allowed) return null;

    const open = () => {
        setInitialEnvironment(null);
        setIsModalOpen(true);
    };

    const close = () => {
        setIsModalOpen(false);
        setInitialEnvironment(null);
    };

    return (
        <>
            <span className="relative group inline-flex">
                <button
                    onClick={open}
                    disabled={disabled}
                    data-tour="integrations.hub-button"
                    aria-label="Tus Integraciones"
                    className={`inline-flex items-center justify-center rounded-lg text-sm transition-colors ${
                        disabled
                            ? 'text-gray-400 dark:text-gray-500 bg-gray-100 dark:bg-gray-700 cursor-not-allowed'
                            : 'subnav-active hover:opacity-90'
                    }`}
                    style={{ width: '2rem', height: '2rem', padding: 0 }}
                >
                    <span>{'\u{1F517}'}</span>
                </button>
                <span className="pointer-events-none absolute left-1/2 -translate-x-1/2 top-full mt-1.5 z-50 whitespace-nowrap rounded-md bg-gray-900 px-2 py-1 text-[11px] font-medium text-white opacity-0 shadow-lg transition-opacity group-hover:opacity-100 dark:bg-gray-700">
                    {disabled ? 'Selecciona un negocio primero' : 'Tus Integraciones'}
                </span>
            </span>

            <MyIntegrationsModal
                isOpen={isModalOpen}
                onClose={close}
                businessId={businessId}
                initialEnvironment={initialEnvironment}
            />
        </>
    );
}
