'use client';

import React, { memo, useCallback } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useCategories } from '@/services/integrations/core/ui/hooks/useCategories';

// Mapeo de código de categoría → nombre del recurso en BD
const CATEGORY_RESOURCE_MAP: Record<string, string> = {
    'ecommerce': 'Integraciones-E-commerce',
    'invoicing': 'Integraciones-Facturacion-Electronica',
    'messaging': 'Integraciones-Mensajeria',
    'payment': 'Integraciones-Pagos',
    'shipping': 'Integraciones-Logistica',
    'platform': 'Integraciones-Platform',
};

// Emojis por categoría
const CATEGORY_ICONS: Record<string, string> = {
    'platform': '🧩',
    'ecommerce': '🛒',
    'invoicing': '🧾',
    'messaging': '💬',
    'payment': '💳',
    'shipping': '🚚',
};

export const IntegrationsSubNavbar = memo(function IntegrationsSubNavbar() {
    const pathname = usePathname();
    const router = useRouter();
    const searchParams = useSearchParams();
    const { hasPermission, isSuperAdmin, isLoading, permissions } = usePermissions();
    const { actionButtons } = useNavbarActions();
    const { categories, loading: categoriesLoading } = useCategories();

    const isInIntegrationsModule = pathname.startsWith('/integrations');

    const permissionsNotLoaded = isLoading || !permissions || !permissions.resources || permissions.resources.length === 0;
    const canViewTypes = permissionsNotLoaded || isSuperAdmin || hasPermission('Integraciones-Tipos-de-integracion', 'Read');

    const handleTabClick = useCallback((categoryCode: string | null) => {
        if (categoryCode === null) {
            // Tipos de Integración
            router.push('/integrations?tab=types');
        } else {
            router.push(`/integrations?category=${categoryCode}`);
        }
    }, [router]);

    if (!isInIntegrationsModule || categoriesLoading) {
        return null;
    }

    // Filter and sort categories by permissions
    const allowedCategories = categories
        .filter(c => c.is_visible && c.is_active)
        .filter(c => {
            if (isSuperAdmin || permissionsNotLoaded) return true;
            const resource = CATEGORY_RESOURCE_MAP[c.code];
            if (!resource) return true;
            return hasPermission(resource, 'Read');
        })
        .sort((a, b) => (a.display_order || 0) - (b.display_order || 0));

    const currentTab = searchParams.get('tab');
    const currentCategory = searchParams.get('category');
    const isTypesActive = currentTab === 'types';

    // If no category selected and not types, first category is active
    const activeCategoryCode = isTypesActive ? null : (currentCategory || allowedCategories[0]?.code || null);

    const allItems = [
        ...allowedCategories.map(c => ({
            key: c.code,
            label: c.name,
            icon: CATEGORY_ICONS[c.code] || '🔗',
            isActive: !isTypesActive && activeCategoryCode === c.code,
            onClick: () => handleTabClick(c.code),
        })),
        ...(canViewTypes ? [{
            key: 'types',
            label: 'Tipos de Integración',
            icon: '⚙️',
            isActive: isTypesActive,
            onClick: () => handleTabClick(null),
        }] : []),
    ];

    if (allItems.length === 0) {
        return null;
    }

    return (
        <div className="bg-white border-b border-gray-200 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3 flex-wrap">
                        {allItems.map((item) => (
                            <button
                                key={item.key}
                                onClick={item.onClick}
                                className={`px-4 py-3 text-base font-medium whitespace-nowrap transition-all rounded-lg flex items-center gap-3 ${
                                    item.isActive
                                        ? 'bg-purple-200 text-purple-900'
                                        : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900 hover:shadow-md hover:scale-105'
                                }`}
                            >
                                <span>{item.icon}</span>
                                {item.label}
                            </button>
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
