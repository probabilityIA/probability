'use client';

import React, { memo } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';

export const OrdersSubNavbar = memo(function OrdersSubNavbar() {
    const pathname = usePathname();
    const { hasPermission, isSuperAdmin, isLoading, permissions } = usePermissions();

    // Si está cargando, no hay permisos definidos, o resources es null/vacío, mostrar todo por defecto
    const permissionsNotLoaded = isLoading || !permissions || !permissions.resources || permissions.resources.length === 0;

    // Verificar permisos para cada recurso
    const canViewProducts = permissionsNotLoaded || isSuperAdmin || hasPermission('Productos', 'Read');
    const canViewOrders = permissionsNotLoaded || isSuperAdmin || hasPermission('Ordenes', 'Read');
    const canViewShipments = permissionsNotLoaded || isSuperAdmin || hasPermission('Envios', 'Read');
    const canViewOrderStatus = permissionsNotLoaded || isSuperAdmin || hasPermission('Estado de Ordenes', 'Read');
    const canViewNotifications = permissionsNotLoaded || isSuperAdmin || hasPermission('Configuración de Notificaciones', 'Read');
    const canViewOriginAddresses = permissionsNotLoaded || isSuperAdmin || hasPermission('Envios', 'Read');

    // Solo mostrar si estamos en alguna de estas secciones
    const isInOrdersModule = pathname.startsWith('/products') ||
                            pathname.startsWith('/orders') ||
                            pathname.startsWith('/shipments') ||
                            pathname.startsWith('/order-status') ||
                            pathname.startsWith('/notification-config') ||
                            pathname.startsWith('/shipments/origin-addresses');

    if (!isInOrdersModule) {
        return null;
    }

    const isActive = (path: string) => {
        // Exact match
        if (pathname === path) return true;

        // For /shipments, only match if NOT /shipments/origin-addresses
        if (path === '/shipments') {
            return pathname.startsWith('/shipments') && !pathname.startsWith('/shipments/origin-addresses');
        }

        // For other paths, use startsWith
        return pathname.startsWith(path);
    };

    // Todos los items del sidebar de órdenes en orden
    const menuItems = [
        { section: 'CATÁLOGO', items: [
            canViewProducts && { href: '/products', label: 'Productos', icon: '🛍️' },
        ].filter(Boolean) },
        { section: 'OPERACIONES', items: [
            canViewOrders && { href: '/orders', label: 'Órdenes', icon: '📦' },
            canViewShipments && { href: '/shipments', label: 'Envíos', icon: '🚚' },
        ].filter(Boolean) },
        { section: 'CONFIGURACIÓN', items: [
            canViewOrderStatus && { href: '/order-status', label: 'Estados de Orden', icon: '✅' },
            canViewNotifications && { href: '/notification-config', label: 'Notificaciones', icon: '🔔' },
            canViewOriginAddresses && { href: '/shipments/origin-addresses', label: 'Direcciones de Origen', icon: '📍' },
        ].filter(Boolean) },
    ];

    // Aplanar todos los items
    const allItems = menuItems.flatMap(section => section.items);

    if (allItems.length === 0) {
        return null;
    }

    return (
        <div className="bg-white border-b border-gray-200 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center overflow-x-auto gap-3">
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
            </div>
        </div>
    );
});
