'use client';

import React, { memo } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useWalletBusiness } from '@/shared/contexts/wallet-business-context';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';

export const WalletSubNavbar = memo(function WalletSubNavbar() {
    const pathname = usePathname();
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId, setSelectedBusinessId } = useWalletBusiness();

    if (!pathname.startsWith('/wallet')) {
        return null;
    }

    const isActive = (path: string) => {
        if (path === '/wallet') {
            return pathname === '/wallet';
        }
        return pathname.startsWith(path);
    };

    const menuItems = [
        { href: '/wallet/saldos', label: 'Saldos', icon: '💳' },
        { href: '/wallet/finanzas', label: 'Finanzas', icon: '📊' },
    ];

    return (
        <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center gap-4">
                    <div className="flex items-center gap-1 flex-nowrap flex-1 min-w-0 overflow-x-auto subnavbar-scroll">
                        {menuItems.map((item) => {
                            const active = isActive(item.href);
                            return (
                                <Link
                                    key={item.href}
                                    href={item.href}
                                    className={`flex-shrink-0 px-3 py-2 text-sm font-medium whitespace-nowrap transition-all rounded-md flex items-center gap-1.5 ${
                                        active
                                            ? 'subnav-active'
                                            : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-gray-100'
                                    }`}
                                >
                                    <span className="text-base">{item.icon}</span>
                                    {item.label}
                                </Link>
                            );
                        })}
                    </div>
                    {isSuperAdmin && (
                        <SuperAdminBusinessSelector
                            value={selectedBusinessId}
                            onChange={setSelectedBusinessId}
                            variant="navbar"
                            placeholder="— Selecciona un negocio —"
                        />
                    )}
                </div>
            </div>
        </div>
    );
});
