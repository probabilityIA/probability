'use client';

import { useState, useEffect, useMemo } from 'react';
import {
    IntegrationList,
    IntegrationForm,
    IntegrationTypeList,
    IntegrationTypeForm,
    CreateIntegrationModal,
    IntegrationCategoryTabs,
    useCategories,
} from '@/services/integrations/core/ui';
import { Button, Modal } from '@/shared/ui';
import { IntegrationType, Integration } from '@/services/integrations/core/domain/types';
import { getIntegrationByIdAction } from '@/services/integrations/core/infra/actions';

import { useSearchParams } from 'next/navigation';
import { ShopifyOAuthCallback } from '@/services/integrations/ecommerce/shopify/ui';
import { MercadoLibreOAuthCallback } from '@/services/integrations/ecommerce/mercadolibre/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useIntegrationsBusiness } from '@/shared/contexts/integrations-business-context';
import { WooStorePowerWidget } from '@/services/woostore/ui/components/WooStorePowerWidget';

const ALL_TAB_CATEGORIES = 'platform,ecommerce,invoicing,messaging';

const CATEGORY_RESOURCE_MAP: Record<string, string> = {
    'ecommerce': 'Integraciones-E-commerce',
    'invoicing': 'Integraciones-Facturacion-Electronica',
    'messaging': 'Integraciones-Mensajeria',
    'payment': 'Integraciones-Pagos',
    'shipping': 'Integraciones-Logistica',
    'platform': 'Integraciones-Platform',
};

export default function IntegrationsPage() {
    const searchParams = useSearchParams();
    const isOAuthCallback = searchParams.get('shopify_oauth');
    const isMeliOAuthCallback = searchParams.get('meli_oauth');

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showEditIntegrationModal, setShowEditIntegrationModal] = useState(false);
    const [selectedType, setSelectedType] = useState<IntegrationType | undefined>(undefined);
    const [selectedIntegration, setSelectedIntegration] = useState<Integration | undefined>(undefined);
    const [refreshKey, setRefreshKey] = useState(0);

    const { categories, loading: categoriesLoading } = useCategories();
    const { hasPermission, isSuperAdmin } = usePermissions();
    const { setActionButtons } = useNavbarActions();
    const { selectedBusinessId } = useIntegrationsBusiness();

    const needsBusiness = isSuperAdmin && !selectedBusinessId;
    const businessIdForList = isSuperAdmin ? selectedBusinessId : null;

    const currentTab = searchParams.get('tab');
    const currentCategory = searchParams.get('category');
    const isTypesTab = currentTab === 'types';
    const isEnvironmentTab = currentTab === 'environment' && isSuperAdmin;
    const isAllTab = currentTab === 'all';

    const allowedCategories = useMemo(() => {
        return categories.filter(c => {
            if (isSuperAdmin) return true;
            const resource = CATEGORY_RESOURCE_MAP[c.code];
            if (!resource) return true;
            return hasPermission(resource, 'Read');
        });
    }, [categories, isSuperAdmin, hasPermission]);

    const activeCategoryCode = useMemo(() => {
        if (isTypesTab || isEnvironmentTab || isAllTab) return null;
        if (currentCategory) return currentCategory;
        const first = allowedCategories
            .filter(c => c.is_visible && c.is_active)
            .sort((a, b) => (a.display_order || 0) - (b.display_order || 0))[0];
        return first?.code || null;
    }, [isTypesTab, isEnvironmentTab, isAllTab, currentCategory, allowedCategories]);

    useEffect(() => {
        if (isTypesTab) {
            setActionButtons(
                <Button
                    variant="primary"
                    onClick={() => setShowCreateModal(true)}
                >
                    Crear Tipo
                </Button>
            );
        } else {
            setActionButtons(null);
        }
        return () => setActionButtons(null);
    }, [isTypesTab, setActionButtons]);

    if (isOAuthCallback) {
        return <ShopifyOAuthCallback />;
    }

    if (isMeliOAuthCallback) {
        return <MercadoLibreOAuthCallback />;
    }

    const handleSuccess = () => {
        setShowCreateModal(false);
        setShowEditModal(false);
        setShowEditIntegrationModal(false);
        setSelectedType(undefined);
        setSelectedIntegration(undefined);
        setRefreshKey(prev => prev + 1);
    };

    const handleModalClose = () => {
        setShowCreateModal(false);
        setShowEditModal(false);
        setShowEditIntegrationModal(false);
        setSelectedType(undefined);
        setSelectedIntegration(undefined);
    };

    const handleEditType = (type: IntegrationType) => {
        setSelectedType(type);
        setShowEditModal(true);
    };

    const handleEditIntegration = async (integration: Integration) => {
        try {
            const response = await getIntegrationByIdAction(integration.id);
            if (response.success && response.data) {
                setSelectedIntegration(response.data);
                setShowEditIntegrationModal(true);
            } else {
                console.error('Error al obtener integracion:', response.message);
                alert('Error al cargar la integracion para editar');
            }
        } catch (error) {
            console.error('Error al obtener integracion:', error);
            alert('Error al cargar la integracion para editar');
        }
    };

    const withTabs = (content: React.ReactNode) => (
        <div className="mx-auto w-full max-w-7xl">
            <div className="rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-sm overflow-hidden">
                <IntegrationCategoryTabs />
                <div className="p-6">{content}</div>
            </div>
        </div>
    );

    return (
        <div className="w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            {isEnvironmentTab ? withTabs(
                <div className="max-w-3xl space-y-4">
                    <div>
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Ambiente de pruebas</h2>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            Enciende o apaga el ambiente de WooCommerce que usamos para probar las integraciones.
                        </p>
                    </div>
                    <WooStorePowerWidget />
                </div>
            ) : isTypesTab ? withTabs(
                <IntegrationTypeList
                    key={`types-${refreshKey}`}
                    onEdit={handleEditType}
                />
            ) : needsBusiness ? withTabs(
                <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-gray-300 dark:border-gray-700 py-16 text-center">
                    <p className="text-base font-medium text-gray-700 dark:text-gray-200">Selecciona un negocio</p>
                    <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                        Como super admin, elige un negocio en el selector de arriba para ver sus integraciones.
                    </p>
                </div>
            ) : isAllTab ? (
                <IntegrationList
                    key={`all-${businessIdForList ?? 'none'}-${refreshKey}`}
                    filterCategory={ALL_TAB_CATEGORIES}
                    businessId={businessIdForList}
                    onEdit={handleEditIntegration}
                    onCreate={() => setShowCreateModal(true)}
                />
            ) : activeCategoryCode !== null ? (
                <IntegrationList
                    key={`cat-${activeCategoryCode}-${businessIdForList ?? 'none'}-${refreshKey}`}
                    filterCategory={activeCategoryCode}
                    businessId={businessIdForList}
                    onEdit={handleEditIntegration}
                    onCreate={() => setShowCreateModal(true)}
                />
            ) : withTabs(
                <IntegrationTypeList
                    key={`types-${refreshKey}`}
                    onEdit={handleEditType}
                />
            )}

            {!isTypesTab ? (
                <CreateIntegrationModal
                    isOpen={showCreateModal}
                    onClose={handleModalClose}
                    categories={allowedCategories}
                    onSuccess={handleSuccess}
                />
            ) : (
                <Modal
                    isOpen={showCreateModal}
                    onClose={handleModalClose}
                    title="Nuevo Tipo de Integración"
                    size="full"
                >
                    <IntegrationTypeForm
                        onSuccess={handleSuccess}
                        onCancel={handleModalClose}
                    />
                </Modal>
            )}

            <Modal
                isOpen={showEditModal}
                onClose={handleModalClose}
                title="Editar Tipo de Integración"
                size="full"
            >
                <IntegrationTypeForm
                    integrationType={selectedType}
                    onSuccess={handleSuccess}
                    onCancel={handleModalClose}
                />
            </Modal>

            <Modal
                isOpen={showEditIntegrationModal}
                onClose={handleModalClose}
                title={(
                    <span className="inline-flex items-center justify-center gap-2">
                        <span className="h-2.5 w-2.5 rounded-full bg-green-400 animate-pulse shadow-[0_0_8px_rgba(74,222,128,0.9)]" />
                        Editar Integración
                    </span>
                )}
                size={selectedIntegration && [1, 3, 4, 8, 16, 33].includes(Number(selectedIntegration.integration_type_id)) ? '4xl' : '5xl'}
            >
                <div
                    style={
                        selectedIntegration && [1, 3, 4, 8, 16, 33].includes(Number(selectedIntegration.integration_type_id))
                            ? { width: 'min(768px, 92vw)' }
                            : undefined
                    }
                >
                    <IntegrationForm
                        integration={selectedIntegration}
                        onSuccess={handleSuccess}
                        onCancel={handleModalClose}
                    />
                </div>
            </Modal>
        </div>
    );
}
