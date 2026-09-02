'use client';

import React, { memo } from 'react';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';

export const TICKETS_TABS_SLOT_ID = 'tickets-tabs-slot';
export const TICKETS_ACTIONS_SLOT_ID = 'tickets-actions-slot';
export const TICKETS_FILTERS_SLOT_ID = 'tickets-filters-slot';

export const TicketsSubNavbar = memo(function TicketsSubNavbar() {
    const pathname = usePathname();
    const { isSuperAdmin } = usePermissions();

    if (!pathname.startsWith('/tickets') || !isSuperAdmin) {
        return null;
    }

    return (
        <div className="subnav-surface border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 min-h-[56px] flex items-center gap-3">
                <span id={TICKETS_TABS_SLOT_ID} className="inline-flex items-center gap-2 shrink-0 empty:hidden" />
                <span id={TICKETS_ACTIONS_SLOT_ID} className="inline-flex items-center gap-2 shrink-0 empty:hidden" />
                <span id={TICKETS_FILTERS_SLOT_ID} className="flex-1 min-w-0 empty:hidden" />
            </div>
        </div>
    );
});
