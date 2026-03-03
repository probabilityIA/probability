'use client';

import React, { memo } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useNotificationBusiness } from '@/shared/contexts/notification-business-context';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';

export const NotificationsSubNavbar = memo(function NotificationsSubNavbar() {
    const pathname = usePathname();
    const { isSuperAdmin, isLoading, permissions } = usePermissions();
    const { actionButtons } = useNavbarActions();
    const { selectedBusinessId, setSelectedBusinessId } = useNotificationBusiness();

    const isInModule = pathname.startsWith('/notification-config') ||
                       pathname.startsWith('/notification-channels') ||
                       pathname.startsWith('/notification-event-types');

    if (!isInModule) {
        return null;
    }

    const isActive = (path: string) => pathname === path || pathname.startsWith(path + '/');

    const menuItems = [
        { href: '/notification-config', label: 'Configuraciones', icon: '🔔', enabled: true },
        { href: '/notification-channels', label: 'Canales', icon: '📡', enabled: isSuperAdmin },
        { href: '/notification-event-types', label: 'Tipos de Eventos', icon: '📋', enabled: isSuperAdmin },
    ];

    return (
        <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3 flex-wrap">
                        {menuItems.map((item) => {
                            if (!item.enabled) {
                                return (
                                    <span
                                        key={item.href}
                                        className="px-4 py-3 text-base font-medium whitespace-nowrap rounded-lg flex items-center gap-3 text-gray-400 dark:text-gray-500 cursor-not-allowed select-none"
                                        title="Solo disponible para Super Admin"
                                    >
                                        <span>{item.icon}</span>
                                        {item.label}
                                    </span>
                                );
                            }
                            return (
                                <Link
                                    key={item.href}
                                    href={item.href}
                                    className={`px-4 py-3 text-base font-medium whitespace-nowrap transition-all rounded-lg flex items-center gap-3 ${
                                        isActive(item.href)
                                            ? 'bg-purple-200 dark:bg-purple-900/50 text-purple-900 dark:text-purple-200'
                                            : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-gray-100 hover:shadow-md hover:scale-105'
                                    }`}
                                >
                                    <span>{item.icon}</span>
                                    {item.label}
                                </Link>
                            );
                        })}
                    </div>
                    <div className="flex items-center gap-2 ml-4">
                        <SuperAdminBusinessSelector
                            value={selectedBusinessId}
                            onChange={setSelectedBusinessId}
                            variant="navbar"
                            placeholder="— Selecciona un negocio —"
                        />
                        {actionButtons}
                    </div>
                </div>
            </div>
        </div>
    );
});
