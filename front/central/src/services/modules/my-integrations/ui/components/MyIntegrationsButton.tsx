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
            <button
                onClick={() => setIsModalOpen(true)}
                disabled={disabled}
                className={`flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-lg transition-colors whitespace-nowrap ${
                    disabled
                        ? 'text-gray-400 dark:text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 cursor-not-allowed'
                        : 'subnav-active hover:opacity-90'
                }`}
                title={disabled ? 'Selecciona un negocio primero' : 'Ver tus integraciones'}
            >
                <span>🔗</span>
                <span className="hidden sm:inline">Tus Integraciones</span>
            </button>

            <MyIntegrationsModal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                businessId={businessId}
            />
        </>
    );
}
