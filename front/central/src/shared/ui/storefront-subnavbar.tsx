'use client';

import React, { memo, useCallback } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useStorefrontBusiness } from '@/shared/contexts/storefront-business-context';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';
import { MyIntegrationsButton } from '@/services/modules/my-integrations/ui';

export const StorefrontSubNavbar = memo(function StorefrontSubNavbar() {
    const pathname = usePathname();
    const router = useRouter();
    const { actionButtons } = useNavbarActions();
    const { selectedBusinessId, setSelectedBusinessId } = useStorefrontBusiness();
    const { isSuperAdmin, permissions } = usePermissions();

    const handleBusinessChange = useCallback((id: number | null) => {
        setSelectedBusinessId(id);
        router.refresh();
    }, [setSelectedBusinessId, router]);

    const isInModule = pathname.startsWith('/storefront') || pathname.startsWith('/website-config');

    if (!isInModule) {
        return null;
    }

    const isActive = (path: string) => pathname.startsWith(path);

    const canViewWebsiteConfig = isSuperAdmin || permissions?.role_name === 'Administrador';

    const menuItems = [
        { href: '/storefront/catalogo', label: 'Catalogo', icon: '🛍️' },
        { href: '/storefront/nuevo-pedido', label: 'Nuevo Pedido', icon: '➕' },
        { href: '/storefront/pedidos', label: 'Pedidos', icon: '📋' },
        ...(canViewWebsiteConfig ? [{ href: '/website-config', label: 'Mi Sitio Web', icon: '🌐' }] : []),
    ];

    return (
        <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3 flex-wrap">
                        {menuItems.map((item) => (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={`px-4 py-3 text-base font-medium whitespace-nowrap transition-all rounded-lg flex items-center gap-3 ${
                                    isActive(item.href)
                                        ? 'bg-purple-200 dark:bg-purple-900/50 text-purple-900 dark:text-purple-200'
                                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:text-white dark:hover:text-gray-100 hover:shadow-md hover:scale-105'
                                }`}
                            >
                                <span>{item.icon}</span>
                                {item.label}
                            </Link>
                        ))}
                    </div>
                    <div className="flex items-center gap-2 ml-4">
                        <MyIntegrationsButton businessId={selectedBusinessId} />
                        <SuperAdminBusinessSelector
                            value={selectedBusinessId}
                            onChange={handleBusinessChange}
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
