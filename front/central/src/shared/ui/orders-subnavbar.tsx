'use client';

import React, { memo } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useOrdersBusiness } from '@/shared/contexts/orders-business-context';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';
import { MyIntegrationsButton } from '@/services/modules/my-integrations/ui';
import { OrdersEffectivenessKpis } from '@/services/modules/orders/ui';

export const ORDERS_FILTERS_SLOT_ID = 'orders-filters-slot';
export const ORDERS_ACTIONS_SLOT_ID = 'orders-actions-slot';

export const OrdersSubNavbar = memo(function OrdersSubNavbar() {
    const pathname = usePathname();
    const { hasPermission, isSuperAdmin, isLoading, permissions } = usePermissions();
    const { actionButtons } = useNavbarActions();
    const { selectedBusinessId, setSelectedBusinessId } = useOrdersBusiness();

    const permissionsNotLoaded = isLoading || !permissions || !permissions.resources || permissions.resources.length === 0;

    const canViewOrders = permissionsNotLoaded || isSuperAdmin || hasPermission('Ordenes', 'Read');
    const canViewShipments = permissionsNotLoaded || isSuperAdmin || hasPermission('Envios', 'Read');

    const isInOrdersModule = pathname.startsWith('/orders') ||
                            pathname.startsWith('/shipments') ||
                            pathname.startsWith('/shipping-margins');

    if (!isInOrdersModule) {
        return null;
    }

    const isActive = (path: string) => {
        if (pathname === path) return true;

        if (path === '/shipments') {
            return pathname.startsWith('/shipments')
                && !pathname.startsWith('/shipments/origin-addresses')
                && !pathname.startsWith('/shipments/cod')
                && !pathname.startsWith('/shipments/quotes');
        }

        return pathname.startsWith(path);
    };

    const menuItems = [
        { section: 'OPERACIONES', items: [
            canViewOrders && { href: '/orders', label: 'Órdenes', icon: '\u{1F4E6}' },
            canViewShipments && { href: '/shipments', label: 'Envíos', icon: '\u{1F69A}' },
            canViewShipments && { href: '/shipments/cod', label: 'Recaudo contra entrega', icon: '\u{1F4B5}' },
            canViewShipments && { href: '/shipments/quotes', label: 'Cotizaciones', icon: '\u{1F9FE}' },
        ].filter(Boolean) },
        { section: 'CONFIGURACION', items: [
            isSuperAdmin && { href: '/shipping-margins', label: 'Margenes de envio', icon: '\u{1F4B0}' },
        ].filter(Boolean) },
    ];

    const allItems = menuItems.flatMap(section => section.items);

    if (allItems.length === 0) {
        return null;
    }

    return (
        <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 pt-2">
                <div className="flex items-center gap-3">
                    <div className="flex items-center gap-1.5 flex-wrap">
                        {allItems.map((item: any) => (
                            <Link
                                key={item.href}
                                href={item.href}
                                style={isActive(item.href) ? { backgroundColor: 'var(--color-secondary-500)', color: 'white' } : {}}
                                className={`px-3 py-2 text-sm font-medium whitespace-nowrap transition-all rounded-lg flex items-center gap-2 ${
                                    isActive(item.href)
                                        ? ''
                                        : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-gray-100'
                                }`}
                            >
                                <span>{item.icon}</span>
                                {item.label}
                            </Link>
                        ))}
                    </div>
                    <div className="flex items-center gap-2 ml-auto">
                        <OrdersEffectivenessKpis businessId={selectedBusinessId} />
                        <MyIntegrationsButton businessId={selectedBusinessId} />
                        <SuperAdminBusinessSelector
                            value={selectedBusinessId}
                            onChange={setSelectedBusinessId}
                            variant="navbar"
                            placeholder="— Selecciona un negocio —"
                        />
                        <span id={ORDERS_ACTIONS_SLOT_ID} className="inline-flex items-center gap-2 empty:hidden" />
                        {actionButtons}
                    </div>
                </div>
            </div>
            <div
                id={ORDERS_FILTERS_SLOT_ID}
                className="empty:hidden px-4 sm:px-6 lg:px-8 py-2 mt-1 border-t border-gray-100 dark:border-gray-700/60"
            />
        </div>
    );
});
