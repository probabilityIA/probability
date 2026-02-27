'use client';

import React, { memo } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';

export const InventorySubNavbar = memo(function InventorySubNavbar() {
    const pathname = usePathname();
    const { hasPermission, isSuperAdmin, isLoading, permissions } = usePermissions();
    const { actionButtons } = useNavbarActions();

    const permissionsNotLoaded = isLoading || !permissions || !permissions.resources || permissions.resources.length === 0;

    const canViewProducts = permissionsNotLoaded || isSuperAdmin || hasPermission('Productos', 'Read') || hasPermission('Products', 'Read');
    const canViewWarehouses = permissionsNotLoaded || isSuperAdmin || hasPermission('Bodegas', 'Read') || hasPermission('Warehouses', 'Read');
    const canViewInventory = permissionsNotLoaded || isSuperAdmin || hasPermission('Inventario', 'Read') || hasPermission('Inventory', 'Read');

    const isInInventoryModule = pathname.startsWith('/products') ||
                                pathname.startsWith('/warehouses') ||
                                pathname.startsWith('/inventory');

    if (!isInInventoryModule) {
        return null;
    }

    const isActive = (path: string) => {
        if (pathname === path) return true;

        // /inventory exact match to avoid overlap with /inventory/movements
        if (path === '/inventory') {
            return pathname === '/inventory';
        }

        return pathname.startsWith(path);
    };

    const menuItems = [
        { section: 'CATALOGO', items: [
            canViewProducts && { href: '/products', label: 'Productos', icon: '📦' },
            canViewWarehouses && { href: '/warehouses', label: 'Bodegas', icon: '🏭' },
        ].filter(Boolean) },
        { section: 'INVENTARIO', items: [
            canViewInventory && { href: '/inventory', label: 'Stock', icon: '📊' },
            canViewInventory && { href: '/inventory/movements', label: 'Movimientos', icon: '🔄' },
        ].filter(Boolean) },
    ];

    const allItems = menuItems.flatMap(section => section.items);

    if (allItems.length === 0) {
        return null;
    }

    return (
        <div className="bg-white border-b border-gray-200 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3 flex-wrap">
                        {allItems.map((item: any) => (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={`px-4 py-3 text-base font-medium whitespace-nowrap transition-all rounded-lg flex items-center gap-3 ${
                                    isActive(item.href)
                                        ? 'bg-purple-200 text-purple-900'
                                        : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900 hover:shadow-md hover:scale-105'
                                }`}
                            >
                                <span>{item.icon}</span>
                                {item.label}
                            </Link>
                        ))}
                    </div>
                    {actionButtons && (
                        <div className="flex gap-2 ml-4">
                            {actionButtons}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
});
