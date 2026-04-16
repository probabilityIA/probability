'use client';

import { useState, useEffect, useMemo } from 'react';
import {
    IntegrationList,
    IntegrationForm,
    IntegrationTypeList,
    IntegrationTypeForm,
    CreateIntegrationModal,
    useCategories,
} from '@/services/integrations/core/ui';
import { Button, Modal } from '@/shared/ui';
import { IntegrationType, Integration } from '@/services/integrations/core/domain/types';
import { getIntegrationByIdAction } from '@/services/integrations/core/infra/actions';

import { useSearchParams } from 'next/navigation';
import { ShopifyOAuthCallback } from '@/services/integrations/ecommerce/shopify/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNavbarActions } from '@/shared/contexts/navbar-context';

// Mapeo de código de categoría → nombre del recurso en BD
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

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showEditIntegrationModal, setShowEditIntegrationModal] = useState(false);
    const [selectedType, setSelectedType] = useState<IntegrationType | undefined>(undefined);
    const [selectedIntegration, setSelectedIntegration] = useState<Integration | undefined>(undefined);
    const [refreshKey, setRefreshKey] = useState(0);

    const { categories, loading: categoriesLoading } = useCategories();
    const { hasPermission, isSuperAdmin } = usePermissions();
    const { setActionButtons } = useNavbarActions();

    // Read active tab and category from URL (driven by subnavbar)
    const currentTab = searchParams.get('tab');
    const currentCategory = searchParams.get('category');
    const isTypesTab = currentTab === 'types';

    // Filter categories by permissions
    const allowedCategories = useMemo(() => {
        return categories.filter(c => {
            if (isSuperAdmin) return true;
            const resource = CATEGORY_RESOURCE_MAP[c.code];
            if (!resource) return true;
            return hasPermission(resource, 'Read');
        });
    }, [categories, isSuperAdmin, hasPermission]);

    // Determine active category code
    const activeCategoryCode = useMemo(() => {
        if (isTypesTab) return null;
        if (currentCategory) return currentCategory;
        const first = allowedCategories
            .filter(c => c.is_visible && c.is_active)
            .sort((a, b) => (a.display_order || 0) - (b.display_order || 0))[0];
        return first?.code || null;
    }, [isTypesTab, currentCategory, allowedCategories]);

    // Set action buttons in navbar
    useEffect(() => {
        setActionButtons(
            <Button
                variant="primary"
                onClick={() => setShowCreateModal(true)}
            >
                {isTypesTab ? 'Crear Tipo' : 'Crear Integración'}
            </Button>
        );
        return () => setActionButtons(null);
    }, [isTypesTab, setActionButtons]);

    if (isOAuthCallback) {
        return <ShopifyOAuthCallback />;
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

    return (
        <div className="w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            {!isTypesTab && activeCategoryCode !== null ? (
                <IntegrationList
                    key={`cat-${activeCategoryCode}-${refreshKey}`}
                    filterCategory={activeCategoryCode}
                    onEdit={async (integration) => {
                        try {
                            const response = await getIntegrationByIdAction(integration.id);
                            if (response.success && response.data) {
                                setSelectedIntegration(response.data);
                                setShowEditIntegrationModal(true);
                            } else {
                                console.error('Error al obtener integración:', response.message);
                                alert('Error al cargar la integración para editar');
                            }
                        } catch (error) {
                            console.error('Error al obtener integración:', error);
                            alert('Error al cargar la integración para editar');
                        }
                    }}
                />
            ) : (
                <IntegrationTypeList
                    key={`types-${refreshKey}`}
                    onEdit={handleEditType}
                />
            )}

            {/* Create Modal */}
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

            {/* Edit Modal for Integration Types */}
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

            {/* Edit Modal for Integrations */}
            <Modal
                isOpen={showEditIntegrationModal}
                onClose={handleModalClose}
                title="Editar Integración"
                size="full"
            >
                <IntegrationForm
                    integration={selectedIntegration}
                    onSuccess={handleSuccess}
                    onCancel={handleModalClose}
                />
            </Modal>
        </div>
    );
}
