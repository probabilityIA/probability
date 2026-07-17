'use client';

import React, { memo } from 'react';
import { usePathname } from 'next/navigation';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useIntegrationsBusiness } from '@/shared/contexts/integrations-business-context';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';

export const IntegrationsSubNavbar = memo(function IntegrationsSubNavbar() {
    const pathname = usePathname();
    const { actionButtons, secondaryContent } = useNavbarActions();
    const { selectedBusinessId, setSelectedBusinessId } = useIntegrationsBusiness();

    if (!pathname.startsWith('/integrations')) {
        return null;
    }

    return (
        <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center justify-end gap-2">
                    <SuperAdminBusinessSelector
                        value={selectedBusinessId}
                        onChange={setSelectedBusinessId}
                        variant="navbar"
                        placeholder="— Selecciona un negocio —"
                    />
                    {actionButtons}
                </div>
            </div>
            {secondaryContent && (
                <div className="border-t border-gray-200 dark:border-gray-700">
                    {secondaryContent}
                </div>
            )}
        </div>
    );
});
