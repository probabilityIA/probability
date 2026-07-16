'use client';

import React, { memo, useCallback, useRef, useState } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useIntegrationsBusiness } from '@/shared/contexts/integrations-business-context';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';
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
    const { actionButtons, secondaryContent } = useNavbarActions();
    const { selectedBusinessId, setSelectedBusinessId } = useIntegrationsBusiness();
    const { categories, loading: categoriesLoading } = useCategories();

    const scrollElRef = useRef<HTMLDivElement | null>(null);
    const roRef = useRef<ResizeObserver | null>(null);
    const [canLeft, setCanLeft] = useState(false);
    const [canRight, setCanRight] = useState(false);

    const updateArrows = useCallback(() => {
        const el = scrollElRef.current;
        if (!el) return;
        setCanLeft(el.scrollLeft > 4);
        setCanRight(el.scrollLeft + el.clientWidth < el.scrollWidth - 4);
    }, []);

    const attachScroll = useCallback((node: HTMLDivElement | null) => {
        if (roRef.current) {
            roRef.current.disconnect();
            roRef.current = null;
        }
        if (scrollElRef.current) {
            scrollElRef.current.removeEventListener('scroll', updateArrows);
        }
        scrollElRef.current = node;
        if (node) {
            node.addEventListener('scroll', updateArrows, { passive: true });
            const ro = new ResizeObserver(() => updateArrows());
            ro.observe(node);
            if (node.firstElementChild) ro.observe(node.firstElementChild);
            roRef.current = ro;
            updateArrows();
        }
    }, [updateArrows]);

    const scrollByDir = useCallback((dir: 'left' | 'right') => {
        const el = scrollElRef.current;
        if (!el) return;
        el.scrollBy({ left: dir === 'left' ? -260 : 260, behavior: 'smooth' });
    }, []);

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

    const SUPER_ADMIN_ONLY_CATEGORIES = new Set(['storefront', 'internal']);

    const allowedCategories = categories
        .filter(c => c.is_visible && c.is_active)
        .filter(c => {
            if (isSuperAdmin) return true;
            if (SUPER_ADMIN_ONLY_CATEGORIES.has(c.code)) return false;
            if (permissionsNotLoaded) return true;
            const resource = CATEGORY_RESOURCE_MAP[c.code];
            if (!resource) return true;
            return hasPermission(resource, 'Read');
        })
        .sort((a, b) => (a.display_order || 0) - (b.display_order || 0));

    const currentTab = searchParams.get('tab');
    const currentCategory = searchParams.get('category');
    const isTypesActive = currentTab === 'types';
    const isEnvironmentActive = currentTab === 'environment';
    const isAllActive = currentTab === 'all';

    const activeCategoryCode = (isTypesActive || isEnvironmentActive || isAllActive) ? null : (currentCategory || allowedCategories[0]?.code || null);

    const allItems = [
        {
            key: 'all',
            label: 'Todos',
            icon: '📋',
            isActive: isAllActive,
            onClick: () => router.push('/integrations?tab=all'),
        },
        ...allowedCategories.map(c => ({
            key: c.code,
            label: c.name,
            icon: CATEGORY_ICONS[c.code] || '🔗',
            isActive: !isTypesActive && !isEnvironmentActive && !isAllActive && activeCategoryCode === c.code,
            onClick: () => handleTabClick(c.code),
        })),
        ...(canViewTypes ? [{
            key: 'types',
            label: 'Tipos de Integración',
            icon: '⚙️',
            isActive: isTypesActive,
            onClick: () => handleTabClick(null),
        }] : []),
        ...(isSuperAdmin ? [{
            key: 'environment',
            label: 'Ambiente',
            icon: '🧪',
            isActive: isEnvironmentActive,
            onClick: () => router.push('/integrations?tab=environment'),
        }] : []),
    ];

    if (allItems.length === 0) {
        return null;
    }

    return (
        <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center justify-between gap-4">
                    <div className="flex items-center gap-1 min-w-0 flex-1">
                        {canLeft && (
                            <button
                                onClick={() => scrollByDir('left')}
                                aria-label="Desplazar pestanas a la izquierda"
                                className="flex-shrink-0 p-1.5 rounded-md text-gray-500 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                            >
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                                </svg>
                            </button>
                        )}
                        <div
                            ref={attachScroll}
                            className="overflow-x-auto min-w-0 flex-1 [&::-webkit-scrollbar]:hidden"
                            style={{ scrollbarWidth: 'none' }}
                        >
                            <div className="flex items-center gap-3 flex-nowrap w-max">
                                {allItems.map((item) => (
                                    <button
                                        key={item.key}
                                        onClick={item.onClick}
                                        className={`shrink-0 px-4 py-3 text-base font-medium whitespace-nowrap transition-all rounded-lg flex items-center gap-3 ${
                                            item.isActive
                                                ? 'bg-purple-200 dark:bg-purple-900/50 text-purple-900 dark:text-purple-200'
                                                : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100 hover:shadow-md hover:scale-105'
                                        }`}
                                    >
                                        <span>{item.icon}</span>
                                        {item.label}
                                    </button>
                                ))}
                            </div>
                        </div>
                        {canRight && (
                            <button
                                onClick={() => scrollByDir('right')}
                                aria-label="Desplazar pestanas a la derecha"
                                className="flex-shrink-0 p-1.5 rounded-md text-gray-500 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                            >
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </button>
                        )}
                    </div>
                    <div className="flex items-center gap-2 ml-4 flex-shrink-0">
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
            {secondaryContent && (
                <div className="border-t border-gray-200 dark:border-gray-700">
                    {secondaryContent}
                </div>
            )}
        </div>
    );
});
