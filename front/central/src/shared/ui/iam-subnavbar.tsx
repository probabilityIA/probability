'use client';

import React, { memo, useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useResourceConfig } from '@/services/auth/business/ui/hooks/useResourceConfig';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';
import { TourLauncher } from '@/services/modules/tours/ui';

export const IAMSubNavbar = memo(function IAMSubNavbar() {
    const pathname = usePathname();
    const { isSuperAdmin, permissions } = usePermissions();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    const isInModule = pathname.startsWith('/resources') || pathname.startsWith('/roles') || pathname.startsWith('/businesses') || pathname.startsWith('/permissions') || pathname.startsWith('/users') || pathname.startsWith('/marketing-leads') || pathname.startsWith('/siigo-referrals');

    const businessIdForConfig = isSuperAdmin
        ? (selectedBusinessId || 0)
        : (permissions?.business_id || 0);
    const { config } = useResourceConfig(businessIdForConfig);

    const businessActiveResources = new Set<string>(
        (config?.resources || []).filter((r: any) => r.is_active).map((r: any) => r.resource_name)
    );

    const isResourceActive = (name: string) => {
        if (isSuperAdmin) return true;
        return businessActiveResources.has(name);
    };

    if (!isInModule) {
        return null;
    }

    const isActive = (path: string) => pathname === path || pathname.startsWith(path + '/');

    const allMenuItems = [
        { href: '/businesses', label: 'Empresas', icon: '🏢', resource: 'Empresas' },
        { href: '/users', label: 'Usuarios', icon: '👤', resource: 'Usuarios' },
        { href: '/resources', label: 'Recursos', icon: '📦', resource: 'Recursos' },
        { href: '/roles', label: 'Roles', icon: '🔐', resource: 'Roles' },
        { href: '/permissions', label: 'Permisos', icon: '📋', resource: 'Permisos' },
    ];
    const menuItems = allMenuItems.filter(i => isResourceActive(i.resource));
    if (isSuperAdmin) {
        menuItems.push({ href: '/marketing-leads', label: 'Leads de encuestas', icon: '📊', resource: 'MarketingLeads' });
        menuItems.push({ href: '/siigo-referrals', label: 'Referidos Siigo', icon: '🤝', resource: 'SiigoReferrals' });
    }

    return (
        <div className="subnav-surface border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center justify-between gap-4">
                    <div className="flex items-center gap-2 flex-wrap">
                        {menuItems.map((item) => (
                            <Link
                                key={item.href}
                                href={item.href}
                                className="px-4 py-3 text-base font-medium whitespace-nowrap transition-all rounded-lg flex items-center gap-2"
                                style={
                                    isActive(item.href)
                                        ? {
                                            backgroundColor: 'var(--color-secondary-500)',
                                            color: 'var(--color-on-secondary, white)',
                                            boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
                                          }
                                        : {
                                            color: 'var(--gray-700)',
                                          }
                                }
                            >
                                <span>{item.icon}</span>
                                {item.label}
                            </Link>
                        ))}
                    </div>
                    <div className="flex items-center gap-3 ml-auto">
                        <TourLauncher />
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
