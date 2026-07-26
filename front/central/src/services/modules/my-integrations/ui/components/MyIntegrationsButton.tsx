'use client';

import { useState } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { MyIntegrationsModal } from './MyIntegrationsModal';

interface MyIntegrationsButtonProps {
    businessId?: number | null;
}

export function MyIntegrationsButton({ businessId }: MyIntegrationsButtonProps) {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const { isSuperAdmin } = usePermissions();

    const disabled = isSuperAdmin && !businessId;

    return (
        <>
            <span className="relative group inline-flex">
                <button
                    onClick={() => setIsModalOpen(true)}
                    disabled={disabled}
                    aria-label="Tus Integraciones"
                    className={`inline-flex items-center justify-center rounded-lg text-sm transition-colors ${
                        disabled
                            ? 'text-gray-400 dark:text-gray-500 bg-gray-100 dark:bg-gray-700 cursor-not-allowed'
                            : 'subnav-active hover:opacity-90'
                    }`}
                    style={{ width: '2rem', height: '2rem', padding: 0 }}
                >
                    <span>🔗</span>
                </button>
                <span className="pointer-events-none absolute left-1/2 -translate-x-1/2 top-full mt-1.5 z-50 whitespace-nowrap rounded-md bg-gray-900 px-2 py-1 text-[11px] font-medium text-white opacity-0 shadow-lg transition-opacity group-hover:opacity-100 dark:bg-gray-700">
                    {disabled ? 'Selecciona un negocio primero' : 'Tus Integraciones'}
                </span>
            </span>

            <MyIntegrationsModal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                businessId={businessId}
            />
        </>
    );
}
