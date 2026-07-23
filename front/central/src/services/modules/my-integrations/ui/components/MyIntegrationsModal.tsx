'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { Modal } from '@/shared/ui/modal';
import {
    getIntegrationCategoriesAction,
    getIntegrationsAction,
    getIntegrationByIdAction,
    activateIntegrationAction,
    deactivateIntegrationAction,
} from '@/services/integrations/core/infra/actions';
import { IntegrationForm, CreateIntegrationModal } from '@/services/integrations/core/ui';
import { getIntegrationStatsAction, type IntegrationStatsItem } from '@/services/integrations/core/infra/actions/stats';
import type { IntegrationCategory, Integration } from '@/services/integrations/core/domain/types';
import { getBusinessConfiguredResourcesAction } from '@/services/auth/business/infra/actions';
import { CHANNEL_CODES, SERVICE_CODES, INTERNAL_CODES, CATEGORY_COLORS, CHANNELS_COLOR } from '../../domain/types';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { CyberCluster } from './CyberCluster';
import { CyberChannelsCluster } from './CyberChannelsCluster';
import { CyberHub } from './CyberHub';
import { GlobalSyncModal } from './GlobalSyncModal';
import { NetworkLinks, type NetworkTarget } from './NetworkLinks';

interface MyIntegrationsModalProps {
    isOpen: boolean;
    onClose: () => void;
    businessId?: number | null;
}

const WIDE_FORM_TYPE_IDS = [1, 3, 4, 8, 16, 33];

const HUB_KEYFRAMES = `
@keyframes cyber-dash { to { stroke-dashoffset: -24; } }
@keyframes cyber-spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
@keyframes cyber-sweep { from { background-position: 200% 0; } to { background-position: -100% 0; } }
.orbit-ring:has(.orbit-chip:hover) { animation-play-state: paused !important; }
.orbit-ring:has(.orbit-chip:hover) .orbit-chip { animation-play-state: paused !important; }
`;

export function MyIntegrationsModal({ isOpen, onClose, businessId }: MyIntegrationsModalProps) {
    const { permissions, isSuperAdmin } = usePermissions();
    const effectiveBusinessId = businessId ?? (isSuperAdmin ? null : permissions?.business_id ?? null);

    const [categories, setCategories] = useState<IntegrationCategory[]>([]);
    const [integrations, setIntegrations] = useState<Integration[]>([]);
    const [stats, setStats] = useState<Record<number, IntegrationStatsItem>>({});
    const [statsLoaded, setStatsLoaded] = useState(false);
    const [resourceActive, setResourceActive] = useState<Record<string, boolean>>({});
    const [loading, setLoading] = useState(true);
    const [togglingId, setTogglingId] = useState<number | null>(null);
    const [editLoadingId, setEditLoadingId] = useState<number | null>(null);
    const [editingIntegration, setEditingIntegration] = useState<Integration | null>(null);
    const [syncModalOpen, setSyncModalOpen] = useState(false);
    const [createModalOpen, setCreateModalOpen] = useState(false);

    const containerRef = useRef<HTMLDivElement | null>(null);
    const hubRef = useRef<HTMLDivElement | null>(null);
    const clusterRefs = useRef<Map<string, HTMLDivElement>>(new Map());

    const fetchData = useCallback(async () => {
        setLoading(true);
        try {
            const intParams: Record<string, unknown> = { page_size: 100 };
            if (effectiveBusinessId) intParams.business_id = effectiveBusinessId;

            const [catRes, intRes, resourcesRes, statsRes] = await Promise.all([
                getIntegrationCategoriesAction(),
                getIntegrationsAction(intParams),
                effectiveBusinessId ? getBusinessConfiguredResourcesAction(effectiveBusinessId) : Promise.resolve(null),
                getIntegrationStatsAction(effectiveBusinessId ?? undefined),
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

            if (statsRes.success && statsRes.data) {
                const map: Record<number, IntegrationStatsItem> = {};
                for (const item of statsRes.data) {
                    map[item.integration_id] = item;
                }
                setStats(map);
                setStatsLoaded(true);
            } else {
                setStats({});
                setStatsLoaded(false);
            }

            if (resourcesRes?.success && resourcesRes.data) {
                const map: Record<string, boolean> = {};
                for (const r of resourcesRes.data.resources || []) {
                    map[r.resource_name] = r.is_active;
                }
                setResourceActive(map);
            } else {
                setResourceActive({});
            }
        } catch (err) {
            console.error('Error fetching integrations data:', err);
        } finally {
            setLoading(false);
        }
    }, [effectiveBusinessId]);

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
            if (res && (res as { success?: boolean }).success !== false) {
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

    const handleEdit = async (integration: Integration) => {
        setEditLoadingId(integration.id);
        try {
            const res = await getIntegrationByIdAction(integration.id);
            if (res.success && res.data) {
                setEditingIntegration(res.data as Integration);
            } else {
                console.error('Error al obtener integracion:', res.message);
            }
        } catch (err) {
            console.error('Error al obtener integracion:', err);
        } finally {
            setEditLoadingId(null);
        }
    };

    const handleEditClose = () => setEditingIntegration(null);

    const handleEditSuccess = () => {
        setEditingIntegration(null);
        fetchData();
    };

    const setClusterRef = useCallback((code: string) => (el: HTMLDivElement | null) => {
        if (el) clusterRefs.current.set(code, el);
        else clusterRefs.current.delete(code);
    }, []);

    const getTargets = useCallback((): NetworkTarget[] => {
        const targets: NetworkTarget[] = [];
        clusterRefs.current.forEach((el, code) => {
            const isChannel = code === 'channels';
            targets.push({
                key: code,
                el,
                dir: isChannel ? 'in' : 'out',
                color: isChannel ? CHANNELS_COLOR : CATEGORY_COLORS[code] || '#6366f1',
            });
        });
        return targets;
    }, []);

    const integrationsByCategory = categories.reduce<Record<string, Integration[]>>((acc, cat) => {
        acc[cat.code] = integrations.filter(i => i.category === cat.code);
        return acc;
    }, {});

    const resolve = (codes: readonly string[]) =>
        codes
            .map(code => categories.find(c => c.code === code))
            .filter((c): c is IntegrationCategory => c !== undefined);

    const channels = resolve(CHANNEL_CODES);
    const services = resolve(SERVICE_CODES);
    const internal = resolve(INTERNAL_CODES);
    const channelIntegrations = channels.flatMap(cat => integrationsByCategory[cat.code] || []);
    const createCategories = categories.filter(c => c.code === 'ecommerce' || c.code === 'invoicing');

    const internalIntegrations = internal.flatMap(cat => integrationsByCategory[cat.code] || []);
    const revision = loading ? 0 : categories.length * 1000 + integrations.length + 1;

    const editIsWide = editingIntegration
        ? WIDE_FORM_TYPE_IDS.includes(Number(editingIntegration.integration_type_id))
        : false;

    const renderCluster = (cat: IntegrationCategory) => (
        <CyberCluster
            key={cat.code}
            category={cat}
            color={CATEGORY_COLORS[cat.code] || '#6366f1'}
            integrations={integrationsByCategory[cat.code] || []}
            onToggle={handleToggle}
            onEdit={handleEdit}
            togglingId={togglingId}
            editingId={editLoadingId}
            anchorRef={setClusterRef(cat.code)}
        />
    );

    return (
        <>
            <Modal
                isOpen={isOpen}
                onClose={onClose}
                title={(
                    <span className="relative block w-full">
                        Tus Integraciones
                        <button
                            onClick={() => setCreateModalOpen(true)}
                            className="absolute right-8 top-1/2 flex -translate-y-1/2 items-center gap-1.5 rounded-lg bg-white/15 px-3 py-1.5 text-sm font-semibold text-white transition-colors hover:bg-white/25"
                        >
                            <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                            Crear Integracion
                        </button>
                    </span>
                )}
                size="4xl"
            >
                <style>{HUB_KEYFRAMES}</style>
                <div style={{ width: 'min(80rem, 92vw)' }}>
                {loading ? (
                    <div className="flex items-center justify-center py-16">
                        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-purple-600" />
                    </div>
                ) : categories.length === 0 ? (
                    <p className="py-12 text-center text-gray-500 dark:text-gray-400">
                        No hay categorias disponibles
                    </p>
                ) : (
                    <div ref={containerRef} className="relative">
                        <NetworkLinks
                            container={containerRef}
                            hub={hubRef}
                            getTargets={getTargets}
                            revision={revision}
                        />
                        <div className="relative z-10 flex flex-col gap-12 pt-3">
                            <CyberChannelsCluster
                                integrations={channelIntegrations}
                                stats={stats}
                                statsLoaded={statsLoaded}
                                color={CHANNELS_COLOR}
                                onToggle={handleToggle}
                                onEdit={handleEdit}
                                togglingId={togglingId}
                                editingId={editLoadingId}
                                anchorRef={setClusterRef('channels')}
                            />
                            <CyberHub
                                ref={hubRef}
                                integrations={internalIntegrations}
                                resourceActive={resourceActive}
                                onSyncClick={() => setSyncModalOpen(true)}
                            />
                            <div className="flex flex-wrap gap-8 lg:flex-nowrap">
                                {services.map(renderCluster)}
                            </div>
                        </div>
                    </div>
                )}
                </div>
            </Modal>

            <Modal
                isOpen={!!editingIntegration}
                onClose={handleEditClose}
                title={(
                    <span className="inline-flex items-center justify-center gap-2">
                        <span className="h-2.5 w-2.5 animate-pulse rounded-full bg-green-400 shadow-[0_0_8px_rgba(74,222,128,0.9)]" />
                        Editar Integracion
                    </span>
                )}
                size={editIsWide ? '4xl' : '5xl'}
                zIndex={60}
            >
                <div style={editIsWide ? { width: 'min(768px, 92vw)' } : undefined}>
                    {editingIntegration && (
                        <IntegrationForm
                            integration={editingIntegration}
                            onSuccess={handleEditSuccess}
                            onCancel={handleEditClose}
                        />
                    )}
                </div>
            </Modal>

            <GlobalSyncModal
                isOpen={syncModalOpen}
                onClose={() => setSyncModalOpen(false)}
                integrations={integrations}
                businessId={effectiveBusinessId}
            />

            <CreateIntegrationModal
                isOpen={createModalOpen}
                onClose={() => setCreateModalOpen(false)}
                zIndex={60}
                categories={createCategories}
                onSuccess={() => {
                    setCreateModalOpen(false);
                    fetchData();
                }}
            />
        </>
    );
}
