'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { Modal } from '@/shared/ui/modal';
import {
    getIntegrationCategoriesAction,
    getIntegrationsAction,
    getIntegrationByIdAction,
    activateIntegrationAction,
    deactivateIntegrationAction,
    updateIntegrationAction,
} from '@/services/integrations/core/infra/actions';
import { IntegrationForm, CreateIntegrationModal } from '@/services/integrations/core/ui';
import { getIntegrationStatsAction, type IntegrationStatsItem } from '@/services/integrations/core/infra/actions/stats';
import type { IntegrationCategory, Integration } from '@/services/integrations/core/domain/types';
import { getBusinessConfiguredResourcesAction } from '@/services/auth/business/infra/actions';
import { CHANNEL_CODES, SERVICE_CODES, INTERNAL_CODES, PLATFORM_CODE, CATEGORY_COLORS, CHANNELS_COLOR } from '../../domain/types';
import { getSyncProvider } from '../providers';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { CyberCluster } from './CyberCluster';
import { CyberChannelsCluster } from './CyberChannelsCluster';
import { CyberHub } from './CyberHub';
import { NetworkLinks, type NetworkTarget } from './NetworkLinks';
import { SyncActivityProvider, type HubView } from '../sync-activity-context';
import { SyncActions } from './SyncActions';
import { ReportView } from './ReportView';
import { fetchSyncFindings } from '../../infra/repository/sync-findings';
import type { FindingsReport } from '../../domain/types';

interface MyIntegrationsModalProps {
    isOpen: boolean;
    onClose: () => void;
    businessId?: number | null;
}

const WIDE_FORM_TYPE_IDS = [1, 3, 4, 8, 16, 33];

const HUB_KEYFRAMES = `
@keyframes cyber-dash { to { stroke-dashoffset: -24; } }
@keyframes cyber-dashflow { to { stroke-dashoffset: -32; } }
@keyframes cyber-travel { from { offset-distance: 0%; } to { offset-distance: 100%; } }
@keyframes cyber-shimmer { to { background-position: -200% 0; } }
@keyframes cyber-arc { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
@keyframes cyber-flicker { 0%,100% { opacity: .95; } 12% { opacity: .35; } 24% { opacity: 1; } 40% { opacity: .5; } 55% { opacity: .9; } 70% { opacity: .3; } 85% { opacity: .85; } }
@keyframes cyber-charge { 0%,100% { box-shadow: 0 0 18px 2px var(--charge), inset 0 0 12px -4px var(--charge); } 50% { box-shadow: 0 0 42px 10px var(--charge), inset 0 0 22px -2px var(--charge); } }
@keyframes cyber-spark { 0% { transform: scale(.4); opacity: 0; } 15% { opacity: 1; } 45% { transform: scale(1.15); opacity: .9; } 100% { transform: scale(1.5); opacity: 0; } }
@keyframes cyber-core-pulse { 0%,100% { box-shadow: 0 0 0 0 rgba(37,99,235,.28), 0 18px 44px rgba(15,50,110,.16); } 50% { box-shadow: 0 0 0 16px rgba(37,99,235,0), 0 18px 44px rgba(15,50,110,.16); } }
@keyframes cyber-spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
@keyframes cyber-alert-pulse { 0%,100% { box-shadow: 0 0 0 0 rgba(245,158,11,.55); } 50% { box-shadow: 0 0 0 7px rgba(245,158,11,0); } }
@keyframes cyber-sweep { from { background-position: 200% 0; } to { background-position: -100% 0; } }
.orbit-ring:has(.orbit-chip:hover) { animation-play-state: paused !important; }
.orbit-ring:has(.orbit-chip:hover) .orbit-chip { animation-play-state: paused !important; }
`;

export function MyIntegrationsModal({ isOpen, onClose, businessId }: MyIntegrationsModalProps) {
    const { permissions, isSuperAdmin } = usePermissions();
    const effectiveBusinessId = businessId ?? (isSuperAdmin ? null : permissions?.business_id ?? null);
    const [findings, setFindings] = useState<FindingsReport | null>(null);
    const [view, setView] = useState<HubView>('diagrama');

    useEffect(() => {
        if (!isOpen) {
            setView('diagrama');
            return;
        }
        const controller = new AbortController();
        fetchSyncFindings(effectiveBusinessId ?? undefined, controller.signal)
            .then(setFindings)
            .catch(() => setFindings(null));
        return () => controller.abort();
    }, [isOpen, effectiveBusinessId]);

    const [categories, setCategories] = useState<IntegrationCategory[]>([]);
    const [integrations, setIntegrations] = useState<Integration[]>([]);
    const [stats, setStats] = useState<Record<number, IntegrationStatsItem>>({});
    const [statsLoaded, setStatsLoaded] = useState(false);
    const [resourceActive, setResourceActive] = useState<Record<string, boolean>>({});
    const [loading, setLoading] = useState(true);
    const [togglingId, setTogglingId] = useState<number | null>(null);
    const [inventoryTogglingId, setInventoryTogglingId] = useState<number | null>(null);
    const [editLoadingId, setEditLoadingId] = useState<number | null>(null);
    const [editingIntegration, setEditingIntegration] = useState<Integration | null>(null);
    const [createModalOpen, setCreateModalOpen] = useState(false);

    const containerRef = useRef<HTMLDivElement | null>(null);
    const channelNodesRef = useRef<HTMLDivElement | null>(null);
    const hubRef = useRef<HTMLDivElement | null>(null);
    const clusterRefs = useRef<Map<string, HTMLDivElement>>(new Map());

    const fetchData = useCallback(async () => {
        setLoading(true);
        try {
            const intParams: Record<string, unknown> = { page_size: 100 };
            if (effectiveBusinessId) intParams.business_id = effectiveBusinessId;

            const settled = await Promise.allSettled([
                getIntegrationCategoriesAction(),
                getIntegrationsAction(intParams),
                effectiveBusinessId ? getBusinessConfiguredResourcesAction(effectiveBusinessId) : Promise.resolve(null),
                getIntegrationStatsAction(effectiveBusinessId ?? undefined),
            ]);

            const settledValue = <T,>(index: number): T | null => {
                const item = settled[index];
                if (item.status === 'rejected') {
                    console.error('Error cargando datos de integraciones:', item.reason);
                    return null;
                }
                return item.value as T;
            };

            const catRes = settledValue<Awaited<ReturnType<typeof getIntegrationCategoriesAction>>>(0);
            const intRes = settledValue<Awaited<ReturnType<typeof getIntegrationsAction>>>(1);
            const resourcesRes = settledValue<Awaited<ReturnType<typeof getBusinessConfiguredResourcesAction>>>(2);
            const statsRes = settledValue<Awaited<ReturnType<typeof getIntegrationStatsAction>>>(3);

            if (catRes?.success && catRes.data) {
                const visible = (catRes.data as IntegrationCategory[])
                    .filter(c => c.is_visible && c.is_active)
                    .sort((a, b) => a.display_order - b.display_order);
                setCategories(visible);
            }

            if (intRes?.success && intRes.data) {
                setIntegrations(intRes.data as Integration[]);
            }

            if (statsRes?.success && statsRes.data) {
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

    const handleToggleInventory = async (integration: Integration) => {
        setInventoryTogglingId(integration.id);
        const enabled = integration.config?.inventory_sync_enabled === true;
        try {
            const res = await updateIntegrationAction(integration.id, {
                config: { ...(integration.config || {}), inventory_sync_enabled: !enabled },
            });
            if (res && (res as { success?: boolean }).success !== false) {
                setIntegrations(prev =>
                    prev.map(i =>
                        i.id === integration.id
                            ? { ...i, config: { ...(i.config || {}), inventory_sync_enabled: !enabled } }
                            : i
                    )
                );
            }
        } catch (err) {
            console.error('Error toggling inventory sync:', err);
        } finally {
            setInventoryTogglingId(null);
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
            if (code === 'channels') {
                const cards = Array.from(el.querySelectorAll<HTMLElement>('[data-node-id]'));
                cards.forEach((card, index) => {
                    const nodeId = Number(card.getAttribute('data-node-id'));
                    if (!Number.isFinite(nodeId)) return;
                    targets.push({
                        key: `channel-${nodeId}`,
                        el: card,
                        dir: 'in',
                        color: CHANNELS_COLOR,
                        nodeId,
                        lane: index,
                        laneCount: cards.length,
                    });
                });
                return;
            }
            targets.push({
                key: code,
                el,
                dir: 'out',
                color: CATEGORY_COLORS[code] || '#6366f1',
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

    const platformIntegrations = integrations.filter(i => i.category === PLATFORM_CODE);
    const ownStats = platformIntegrations.reduce<IntegrationStatsItem | null>((acc, integration) => {
        const item = stats[integration.id];
        if (!item) return acc;
        if (!acc) return { ...item };
        return {
            ...acc,
            orders_count: acc.orders_count + item.orders_count,
            orders_in_progress: acc.orders_in_progress + item.orders_in_progress,
            orders_delivered: acc.orders_delivered + item.orders_delivered,
            orders_cancelled: acc.orders_cancelled + item.orders_cancelled,
            orders_returned: acc.orders_returned + item.orders_returned,
            products_count: acc.products_count + item.products_count,
        };
    }, null);
    const lastOwnOrder = platformIntegrations
        .map(integration => stats[integration.id]?.last_order_at)
        .filter((value): value is string => Boolean(value))
        .sort()
        .pop();

    const orderSources = [...platformIntegrations, ...channelIntegrations];

    const internalIntegrations = internal.flatMap(cat => integrationsByCategory[cat.code] || []);
    const serviceIntegrations = services.flatMap(cat => integrationsByCategory[cat.code] || []);
    const syncableIntegrations = [...channelIntegrations, ...serviceIntegrations.filter(i => getSyncProvider(i.integration_type_id))];
    const revision = loading ? 0 : categories.length * 1000 + integrations.length * 10 + channelIntegrations.length + 1;

    const editIsWide = editingIntegration
        ? WIDE_FORM_TYPE_IDS.includes(Number(editingIntegration.integration_type_id))
        : false;

    const renderCluster = (cat: IntegrationCategory) => (
        <CyberCluster
            key={cat.code}
            category={cat}
            color={CATEGORY_COLORS[cat.code] || '#6366f1'}
            integrations={integrationsByCategory[cat.code] || []}
            stats={stats}
            onToggle={handleToggle}
            onToggleInventory={handleToggleInventory}
            onEdit={handleEdit}
            togglingId={togglingId}
            inventoryTogglingId={inventoryTogglingId}
            editingId={editLoadingId}
            anchorRef={setClusterRef(cat.code)}
        />
    );

    return (
        <SyncActivityProvider
            integrations={syncableIntegrations}
            businessId={effectiveBusinessId}
            view={view}
            onViewChange={setView}
        >
            <Modal
                isOpen={isOpen}
                onClose={onClose}
                title={(
                    <span className="flex w-full items-center gap-4 px-8">
                        <span className="flex-shrink-0">
                            <SyncActions />
                        </span>
                        <span className="min-w-0 flex-1 truncate text-center">Tus Integraciones</span>
                        <button
                            onClick={() => setCreateModalOpen(true)}
                            className="flex flex-shrink-0 items-center gap-1.5 rounded-lg bg-white/15 px-3 py-1.5 text-sm font-semibold text-white transition-colors hover:bg-white/25"
                        >
                            <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                            Crear Integracion
                        </button>
                    </span>
                )}
                size="6xl"
            >
                <style>{HUB_KEYFRAMES}</style>
                <div
                    className="mx-auto"
                    style={{ width: view === 'informe' ? 'min(96rem, 96vw)' : 'min(80rem, 92vw)', maxWidth: '100%' }}
                >
                {loading ? (
                    <div className="flex items-center justify-center py-16">
                        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-purple-600" />
                    </div>
                ) : categories.length === 0 ? (
                    <p className="py-12 text-center text-gray-500 dark:text-gray-400">
                        No hay categorias disponibles
                    </p>
                ) : view === 'informe' ? (
                    <ReportView
                        businessId={effectiveBusinessId}
                        integrations={syncableIntegrations}
                        orderSources={orderSources}
                        stats={stats}
                        statsLoaded={statsLoaded}
                    />
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
                                onToggleInventory={handleToggleInventory}
                                onEdit={handleEdit}
                                togglingId={togglingId}
                                inventoryTogglingId={inventoryTogglingId}
                                editingId={editLoadingId}
                                anchorRef={setClusterRef('channels')}
                            />
                            <div className="pb-12">
                            <CyberHub
                                ref={hubRef}
                                integrations={internalIntegrations}
                                resourceActive={resourceActive}
                                ownStats={statsLoaded ? ownStats : null}
                                lastOrderAt={lastOwnOrder}
                                findingsCount={findings?.findings.length ?? 0}
                            />
                            </div>
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
        </SyncActivityProvider>
    );
}
