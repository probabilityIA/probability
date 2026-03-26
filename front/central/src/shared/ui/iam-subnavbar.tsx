'use client';

import React, { memo, useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';

export const IAMSubNavbar = memo(function IAMSubNavbar() {
    const pathname = usePathname();
    const { isSuperAdmin } = usePermissions();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    const isInModule = pathname.startsWith('/resources') || pathname.startsWith('/roles') || pathname.startsWith('/businesses') || pathname.startsWith('/permissions') || pathname.startsWith('/users');

    if (!isInModule) {
        return null;
    }

    const isActive = (path: string) => pathname === path || pathname.startsWith(path + '/');

    const menuItems = [
        { href: '/businesses', label: 'Empresas', icon: '🏢' },
        { href: '/users', label: 'Usuarios', icon: '👤' },
        { href: '/resources', label: 'Recursos', icon: '📦' },
        { href: '/roles', label: 'Roles', icon: '🔐' },
        { href: '/permissions', label: 'Permisos', icon: '📋' },
    ];

    return (
        <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center justify-between gap-4">
                    <div className="flex items-center gap-2 flex-wrap">
                        {menuItems.map((item) => (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={`px-4 py-3 text-base font-medium whitespace-nowrap transition-all rounded-lg flex items-center gap-2 ${
                                    isActive(item.href)
                                        ? 'bg-purple-600 dark:bg-purple-600 text-white shadow-md'
                                        : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-white hover:shadow-md'
                                }`}
                            >
                                <span>{item.icon}</span>
                                {item.label}
                            </Link>
                        ))}
                    </div>
                    <div className="flex items-center gap-3 ml-auto">
                        {isSuperAdmin && (
                            <SuperAdminBusinessSelector
                                value={selectedBusinessId}
                                onChange={setSelectedBusinessId}
                                variant="navbar"
                            />
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
});
