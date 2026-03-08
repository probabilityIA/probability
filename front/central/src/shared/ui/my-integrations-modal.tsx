'use client';

import { useState, useEffect, useCallback } from 'react';
import { Modal } from '@/shared/ui/modal';
import { getIntegrationCategoriesAction, getIntegrationsAction, activateIntegrationAction, deactivateIntegrationAction } from '@/services/integrations/core/infra/actions';
import type { IntegrationCategory, Integration } from '@/services/integrations/core/domain/types';

const CATEGORY_ICONS: Record<string, string> = {
    'platform': '🧩',
    'ecommerce': '🛒',
    'invoicing': '🧾',
    'messaging': '💬',
    'payment': '💳',
    'shipping': '🚚',
};

interface MyIntegrationsModalProps {
    isOpen: boolean;
    onClose: () => void;
    businessId?: number | null;
}

export function MyIntegrationsModal({ isOpen, onClose, businessId }: MyIntegrationsModalProps) {
    const [categories, setCategories] = useState<IntegrationCategory[]>([]);
    const [integrations, setIntegrations] = useState<Integration[]>([]);
    const [loading, setLoading] = useState(true);
    const [togglingId, setTogglingId] = useState<number | null>(null);

    const fetchData = useCallback(async () => {
        setLoading(true);
        try {
            const intParams: Record<string, any> = { page_size: 100 };
            if (businessId) intParams.business_id = businessId;

            const [catRes, intRes] = await Promise.all([
                getIntegrationCategoriesAction(),
                getIntegrationsAction(intParams),
            ]);

            if (catRes.success && catRes.data) {
                const visible = (catRes.data as IntegrationCategory[])
                    .filter(c => c.is_visible && c.is_active)
                    .sort((a, b) => a.display_order - b.display_order);
                setCategories(visible);
            }

            if (intRes.success && intRes.data) {
                setIntegrations(intRes.data as Integration[]);
            }
        } catch (err) {
            console.error('Error fetching integrations data:', err);
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => {
        if (isOpen) {
            fetchData();
        }
    }, [isOpen, fetchData]);

    const handleToggle = async (integration: Integration) => {
        setTogglingId(integration.id);
        try {
            const action = integration.is_active
                ? deactivateIntegrationAction
                : activateIntegrationAction;
            const res = await action(integration.id);
            if (res && (res as any).success !== false) {
                setIntegrations(prev =>
                    prev.map(i =>
                        i.id === integration.id
                            ? { ...i, is_active: !i.is_active }
                            : i
                    )
                );
            }
        } catch (err) {
            console.error('Error toggling integration:', err);
        } finally {
            setTogglingId(null);
        }
    };

    const grouped = categories.reduce<Record<string, Integration[]>>((acc, cat) => {
        acc[cat.code] = integrations.filter(i => i.category === cat.code);
        return acc;
    }, {});

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Tus Integraciones" size="5xl">
            {loading ? (
                <div className="flex items-center justify-center py-16">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600" />
                </div>
            ) : categories.length === 0 ? (
                <p className="text-center text-gray-500 dark:text-gray-400 py-12">No hay categorías disponibles</p>
            ) : (
                <div className="flex flex-col lg:flex-row items-start lg:items-stretch gap-4 lg:gap-0 overflow-x-auto pb-4">
                    {categories.map((cat, i) => (
                        <div key={cat.code} className="flex flex-col lg:flex-row items-center lg:items-stretch w-full lg:w-auto">
                            <CategoryNode
                                category={cat}
                                integrations={grouped[cat.code] || []}
                                onToggle={handleToggle}
                                togglingId={togglingId}
                            />
                            {i < categories.length - 1 && <FlowArrow />}
                        </div>
                    ))}
                </div>
            )}
        </Modal>
    );
}

function CategoryNode({
    category,
    integrations,
    onToggle,
    togglingId,
}: {
    category: IntegrationCategory;
    integrations: Integration[];
    onToggle: (i: Integration) => void;
    togglingId: number | null;
}) {
    const icon = CATEGORY_ICONS[category.code] || '🔗';
    const color = category.color || '#8B5CF6';

    return (
        <div className="w-full lg:w-56 min-w-[14rem] flex-shrink-0 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden shadow-sm">
            {/* Header */}
            <div
                className="px-4 py-3 text-white font-semibold flex items-center gap-2"
                style={{ backgroundColor: color }}
            >
                <span className="text-lg">{icon}</span>
                <span className="truncate">{category.name}</span>
            </div>

            {/* Integrations list */}
            <div className="p-3 space-y-2 bg-white dark:bg-gray-800 min-h-[80px]">
                {integrations.length === 0 ? (
                    <p className="text-sm text-gray-400 dark:text-gray-500 italic text-center py-4">Sin configurar</p>
                ) : (
                    integrations.map(integration => (
                        <div
                            key={integration.id}
                            className="flex items-center gap-2 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
                        >
                            {/* Logo */}
                            {integration.integration_type?.image_url ? (
                                <img
                                    src={integration.integration_type.image_url}
                                    alt={integration.name}
                                    className="w-6 h-6 rounded object-contain flex-shrink-0"
                                />
                            ) : (
                                <div className="w-6 h-6 rounded bg-gray-200 dark:bg-gray-600 flex items-center justify-center flex-shrink-0">
                                    <span className="text-xs">{icon}</span>
                                </div>
                            )}

                            {/* Name */}
                            <span className="text-sm text-gray-700 dark:text-gray-300 truncate flex-1">
                                {integration.name}
                            </span>

                            {/* Toggle */}
                            <button
                                onClick={() => onToggle(integration)}
                                disabled={togglingId === integration.id}
                                className={`relative inline-flex h-5 w-9 items-center rounded-full flex-shrink-0 transition-colors ${
                                    integration.is_active
                                        ? 'bg-green-500'
                                        : 'bg-gray-300 dark:bg-gray-600'
                                } ${togglingId === integration.id ? 'opacity-50 cursor-wait' : 'cursor-pointer'}`}
                            >
                                <span
                                    className={`inline-block h-3.5 w-3.5 rounded-full bg-white transition-transform shadow-sm ${
                                        integration.is_active ? 'translate-x-4' : 'translate-x-0.5'
                                    }`}
                                />
                            </button>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}

function FlowArrow() {
    return (
        <>
            {/* Horizontal arrow - desktop */}
            <div className="hidden lg:flex items-center justify-center px-2 flex-shrink-0">
                <svg width="32" height="24" viewBox="0 0 32 24" fill="none" className="text-gray-400 dark:text-gray-500">
                    <path d="M0 12H28M28 12L20 4M28 12L20 20" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
            </div>
            {/* Vertical arrow - mobile */}
            <div className="flex lg:hidden items-center justify-center py-1 flex-shrink-0">
                <svg width="24" height="32" viewBox="0 0 24 32" fill="none" className="text-gray-400 dark:text-gray-500">
                    <path d="M12 0V28M12 28L4 20M12 28L20 20" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
            </div>
        </>
    );
}
