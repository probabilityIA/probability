'use client';

import { useState, useEffect, useCallback } from 'react';
import { Modal } from '@/shared/ui/modal';
import {
    getIntegrationCategoriesAction,
    getIntegrationsAction,
    activateIntegrationAction,
    deactivateIntegrationAction,
} from '@/services/integrations/core/infra/actions';
import type { IntegrationCategory, Integration } from '@/services/integrations/core/domain/types';
import { CHANNEL_CODES, SERVICE_CODES } from '../../domain/types';
import { CategoryCard } from './CategoryCard';
import { FlowConverge, FlowDiverge } from './FlowArrow';
import { OrdersHub } from './OrdersHub';

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
        if (isOpen) fetchData();
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
                        i.id === integration.id ? { ...i, is_active: !i.is_active } : i
                    )
                );
            }
        } catch (err) {
            console.error('Error toggling integration:', err);
        } finally {
            setTogglingId(null);
        }
    };

    const handleEdit = (integration: Integration) => {
        // TODO: abrir modal de configuración por tipo de integración
        console.log('Edit integration:', integration.id, integration.name);
    };

    const integrationsByCategory = categories.reduce<Record<string, Integration[]>>((acc, cat) => {
        acc[cat.code] = integrations.filter(i => i.category === cat.code);
        return acc;
    }, {});

    // Resolver categorías visibles por nivel
    const resolve = (codes: readonly string[]) =>
        codes
            .map(code => categories.find(c => c.code === code))
            .filter((c): c is IntegrationCategory => c !== undefined);

    const channels = resolve(CHANNEL_CODES);
    const services = resolve(SERVICE_CODES);

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Tus Integraciones" size="5xl">
            {loading ? (
                <div className="flex items-center justify-center py-16">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600" />
                </div>
            ) : categories.length === 0 ? (
                <p className="text-center text-gray-500 dark:text-gray-400 py-12">No hay categorías disponibles</p>
            ) : (
                <div className="flex flex-col items-center">
                    {/* Canales de venta (paralelos) */}
                    <div className="flex flex-wrap lg:flex-nowrap gap-3 w-full">
                        {channels.map(cat => (
                            <CategoryCard
                                key={cat.code}
                                category={cat}
                                integrations={integrationsByCategory[cat.code] || []}
                                onToggle={handleToggle}
                                onEdit={handleEdit}
                                togglingId={togglingId}
                            />
                        ))}
                    </div>

                    {/* Convergencia → Hub */}
                    <FlowConverge count={channels.length} />
                    <OrdersHub />
                    <FlowDiverge count={services.length} />

                    {/* Servicios (independientes) */}
                    <div className="flex flex-wrap lg:flex-nowrap gap-3 w-full">
                        {services.map(cat => (
                            <CategoryCard
                                key={cat.code}
                                category={cat}
                                integrations={integrationsByCategory[cat.code] || []}
                                onToggle={handleToggle}
                                onEdit={handleEdit}
                                togglingId={togglingId}
                            />
                        ))}
                    </div>
                </div>
            )}
        </Modal>
    );
}
