'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useCategories } from '../hooks/useCategories';

const CATEGORY_RESOURCE_MAP: Record<string, string> = {
    'ecommerce': 'Integraciones-E-commerce',
    'invoicing': 'Integraciones-Facturacion-Electronica',
    'messaging': 'Integraciones-Mensajeria',
    'payment': 'Integraciones-Pagos',
    'shipping': 'Integraciones-Logistica',
    'platform': 'Integraciones-Platform',
};

const CATEGORY_DOT_COLORS: Record<string, string> = {
    'all': '#64748b',
    'platform': '#0ea5e9',
    'ecommerce': '#8b5cf6',
    'invoicing': '#3b82f6',
    'messaging': '#22c55e',
    'payment': '#f59e0b',
    'shipping': '#ef4444',
    'types': '#94a3b8',
    'environment': '#14b8a6',
};

const SUPER_ADMIN_ONLY_CATEGORIES = new Set(['storefront', 'internal']);

export function IntegrationCategoryTabs() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const { hasPermission, isSuperAdmin, isLoading, permissions } = usePermissions();
    const { categories, loading: categoriesLoading } = useCategories();

    const scrollElRef = useRef<HTMLDivElement | null>(null);
    const roRef = useRef<ResizeObserver | null>(null);
    const [canLeft, setCanLeft] = useState(false);
    const [canRight, setCanRight] = useState(false);
    const [isScrollable, setIsScrollable] = useState(false);

    const updateArrows = useCallback(() => {
        const el = scrollElRef.current;
        if (!el) return;
        setCanLeft(el.scrollLeft > 4);
        setCanRight(el.scrollLeft + el.clientWidth < el.scrollWidth - 4);
        setIsScrollable(el.scrollWidth > el.clientWidth + 4);
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
        el.scrollBy({ left: dir === 'left' ? -220 : 220, behavior: 'smooth' });
    }, []);

    const currentTab = searchParams.get('tab');
    const currentCategory = searchParams.get('category');
    const activeKey = currentTab || currentCategory || '';

    useEffect(() => {
        const el = scrollElRef.current;
        if (!el) return;
        const target = el.querySelector<HTMLElement>('[data-tab-active="true"]');
        if (target) {
            target.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'nearest' });
        }
    }, [activeKey, categoriesLoading]);

    if (categoriesLoading) {
        return null;
    }

    const permissionsNotLoaded = isLoading || !permissions || !permissions.resources || permissions.resources.length === 0;
    const canViewTypes = permissionsNotLoaded || isSuperAdmin || hasPermission('Integraciones-Tipos-de-integracion', 'Read');

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

    const isTypesActive = currentTab === 'types';
    const isEnvironmentActive = currentTab === 'environment';
    const isAllActive = currentTab === 'all';
    const activeCategoryCode = (isTypesActive || isEnvironmentActive || isAllActive)
        ? null
        : (currentCategory || allowedCategories[0]?.code || null);

    const dotColor = (code: string, fallback?: string) =>
        CATEGORY_DOT_COLORS[code] || (fallback && fallback.startsWith('#') ? fallback : '#94a3b8');

    const allItems = [
        {
            key: 'all',
            label: 'Todos',
            color: CATEGORY_DOT_COLORS['all'],
            isActive: isAllActive,
            onClick: () => router.push('/integrations?tab=all'),
        },
        ...allowedCategories.map(c => ({
            key: c.code,
            label: c.name,
            color: dotColor(c.code, c.color),
            isActive: !isTypesActive && !isEnvironmentActive && !isAllActive && activeCategoryCode === c.code,
            onClick: () => router.push(`/integrations?category=${c.code}`),
        })),
        ...(canViewTypes ? [{
            key: 'types',
            label: 'Tipos de Integración',
            color: CATEGORY_DOT_COLORS['types'],
            isActive: isTypesActive,
            onClick: () => router.push('/integrations?tab=types'),
        }] : []),
        ...(isSuperAdmin ? [{
            key: 'environment',
            label: 'Ambiente',
            color: CATEGORY_DOT_COLORS['environment'],
            isActive: isEnvironmentActive,
            onClick: () => router.push('/integrations?tab=environment'),
        }] : []),
    ];

    if (allItems.length === 0) {
        return null;
    }

    return (
        <div className="flex items-center gap-2 px-6 border-b border-gray-100 dark:border-gray-700">
            {isScrollable && (
                <button
                    onClick={() => scrollByDir('left')}
                    disabled={!canLeft}
                    aria-label="Desplazar pestanas a la izquierda"
                    className="flex-shrink-0 w-[30px] h-[30px] flex items-center justify-center rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors disabled:opacity-40"
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                    </svg>
                </button>
            )}
            <div
                ref={attachScroll}
                className="overflow-x-auto min-w-0 flex-1 [&::-webkit-scrollbar]:hidden"
                style={{ scrollbarWidth: 'none' }}
            >
                <div className="flex items-stretch gap-0.5 flex-nowrap w-max">
                    {allItems.map((item) => (
                        <button
                            key={item.key}
                            onClick={item.onClick}
                            data-tab-active={item.isActive}
                            className={`shrink-0 h-11 px-3.5 text-sm whitespace-nowrap flex items-center gap-2 border-b-2 -mb-px transition-colors ${
                                item.isActive
                                    ? 'font-bold text-purple-600 dark:text-purple-300 border-purple-600 dark:border-purple-400'
                                    : 'font-medium text-gray-600 dark:text-gray-300 border-transparent hover:text-gray-900 dark:hover:text-white'
                            }`}
                        >
                            <span
                                className="w-[7px] h-[7px] rounded-full flex-shrink-0"
                                style={{ backgroundColor: item.color }}
                            />
                            {item.label}
                        </button>
                    ))}
                </div>
            </div>
            {isScrollable && (
                <button
                    onClick={() => scrollByDir('right')}
                    disabled={!canRight}
                    aria-label="Desplazar pestanas a la derecha"
                    className="flex-shrink-0 w-[30px] h-[30px] flex items-center justify-center rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors disabled:opacity-40"
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                </button>
            )}
        </div>
    );
}
