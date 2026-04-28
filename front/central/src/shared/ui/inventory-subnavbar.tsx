'use client';

import React, { memo, useState, useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { AcademicCapIcon } from '@heroicons/react/24/outline';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';
import { MyIntegrationsButton } from '@/services/modules/my-integrations/ui';
import InventoryTour from '@/services/modules/inventory/ui/components/InventoryTour';
import { useResourceConfig } from '@/services/auth/business/ui/hooks/useResourceConfig';

export const InventorySubNavbar = memo(function InventorySubNavbar() {
    const [tourOpen, setTourOpen] = useState(false);
    const [pulseTour, setPulseTour] = useState(false);

    useEffect(() => {
        try {
            const seen = localStorage.getItem('inventory_tour_seen_v1');
            if (!seen) setPulseTour(true);
        } catch {}
    }, []);
    const pathname = usePathname();
    const { hasPermission, isSuperAdmin, isLoading, permissions } = usePermissions();
    const { actionButtons } = useNavbarActions();
    const { selectedBusinessId, setSelectedBusinessId } = useInventoryBusiness();

    const permissionsNotLoaded = isLoading || !permissions || !permissions.resources || permissions.resources.length === 0;

    const businessIdForConfig = isSuperAdmin
        ? (selectedBusinessId || 0)
        : (permissions?.business_id || 0);
    const { config: businessConfig, loading: businessConfigLoading } = useResourceConfig(businessIdForConfig);
    const businessActiveResources = businessConfig?.resources
        ?.filter((r: any) => r.is_active)
        .map((r: any) => r.resource_name) ?? [];
    const businessActiveSet = new Set<string>(businessActiveResources);

    const allow = (resource: string) => {
        if (permissionsNotLoaded) return true;
        if (isSuperAdmin && !selectedBusinessId) return true;
        if (businessConfigLoading) return true;
        if (!businessActiveSet.has(resource)) return false;
        if (isSuperAdmin) return true;
        return hasPermission(resource, 'Read');
    };

    const canViewProducts     = allow('Productos') || allow('Products');
    const canViewWarehouses   = allow('Bodegas')   || allow('Warehouses');
    const canViewStock        = allow('Inventario-Stock');
    const canViewMovements    = allow('Inventario-Movimientos');
    const canViewTraceability = allow('Inventario-Trazabilidad');
    const canViewKardex       = allow('Inventario-Kardex');
    const canViewOperations   = allow('Inventario-Operaciones');
    const canViewSlotting     = allow('Inventario-Slotting');
    const canViewAudit        = allow('Inventario-Auditoria');
    const canViewLPN          = allow('Inventario-LPN');
    const canViewScan         = allow('Inventario-Scan');
    const canViewSyncLogs     = allow('Inventario-Sync-Logs');

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
            canViewStock && { href: '/inventory', label: 'Stock', icon: '📊' },
            canViewMovements && { href: '/inventory/movements', label: 'Movimientos', icon: '🔄' },
            canViewTraceability && { href: '/inventory/traceability', label: 'Trazabilidad', icon: '🏷️' },
            canViewKardex && { href: '/inventory/kardex', label: 'Kardex', icon: '📑' },
        ].filter(Boolean) },
        { section: 'OPERACIONES', items: [
            canViewOperations && { href: '/inventory/operations', label: 'Operaciones', icon: '📥' },
            canViewSlotting && { href: '/inventory/analytics/slotting', label: 'Slotting ABC', icon: '📈' },
            canViewAudit && { href: '/inventory/audit', label: 'Auditoría', icon: '✅' },
        ].filter(Boolean) },
        { section: 'CAPTURE', items: [
            canViewLPN && { href: '/inventory/lpn', label: 'LPN', icon: '📦' },
            canViewScan && { href: '/inventory/mobile', label: 'Scan', icon: '📱' },
            canViewSyncLogs && { href: '/inventory/sync/logs', label: 'Sync Logs', icon: '🔄' },
        ].filter(Boolean) },
    ];

    const populatedSections = menuItems.filter((s) => s.items.length > 0);

    if (populatedSections.length === 0) {
        return null;
    }

    return (
        <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center gap-4">
                    <div className="flex items-center gap-1 flex-nowrap flex-1 min-w-0 overflow-x-auto subnavbar-scroll">
                        {populatedSections.map((section, sIdx) => (
                            <div key={section.section} className="flex items-center gap-1 flex-nowrap flex-shrink-0">
                                {sIdx > 0 && <span className="h-6 w-px bg-gray-200 dark:bg-gray-700 mx-1 flex-shrink-0" />}
                                {section.items.map((item: any) => {
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
                        ))}
                    </div>
                    <div className="flex items-center gap-2 flex-shrink-0">
                        <button
                            onClick={() => { setTourOpen(true); setPulseTour(false); }}
                            className={`p-2 rounded-md transition-all text-white btn-business-primary ${pulseTour ? 'tour-pulse' : ''}`}
                            title={pulseTour ? '¡Nuevo! Tutorial guiado' : 'Tutorial guiado'}
                        >
                            <AcademicCapIcon className="w-5 h-5" />
                        </button>
                        <MyIntegrationsButton businessId={selectedBusinessId} />
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
            <InventoryTour isOpen={tourOpen} onClose={() => setTourOpen(false)} />
        </div>
    );
});
