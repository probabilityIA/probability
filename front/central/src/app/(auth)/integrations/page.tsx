'use client';

import { useState, useEffect } from 'react';
import {
    IntegrationList,
    IntegrationForm,
    IntegrationTypeList,
    IntegrationTypeForm,
    CategoryTabs,
    CreateIntegrationModal,
    useCategories,
} from '@/services/integrations/core/ui';
import { Button, Modal } from '@/shared/ui';
import { IntegrationType, Integration } from '@/services/integrations/core/domain/types';
import { getIntegrationByIdAction } from '@/services/integrations/core/infra/actions';

import { useSearchParams } from 'next/navigation';
import { ShopifyOAuthCallback } from '@/services/integrations/ecommerce/shopify/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

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

    // Hooks deben estar antes de cualquier return condicional
    const [activeTab, setActiveTab] = useState<'integrations' | 'types'>('integrations');
    const [activeCategoryCode, setActiveCategoryCode] = useState<string | null>(null);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showEditIntegrationModal, setShowEditIntegrationModal] = useState(false);
    const [selectedType, setSelectedType] = useState<IntegrationType | undefined>(undefined);
    const [selectedIntegration, setSelectedIntegration] = useState<Integration | undefined>(undefined);
    const [refreshKey, setRefreshKey] = useState(0);

    // Hooks for categories and integrations
    const { categories, loading: categoriesLoading } = useCategories();
    const { hasPermission, isSuperAdmin } = usePermissions();

    // Permisos por tab/categoría
    const canViewTypes = isSuperAdmin || hasPermission('Integraciones-Tipos-de-integracion', 'Read');

    // Filtrar categorías según permisos del usuario
    const allowedCategories = categories.filter(c => {
        if (isSuperAdmin) return true;
        const resource = CATEGORY_RESOURCE_MAP[c.code];
        if (!resource) return true; // categoría sin recurso mapeado: visible por defecto
        return hasPermission(resource, 'Read');
    });

    // Seleccionar primera categoría permitida automáticamente
    useEffect(() => {
        if (!categoriesLoading && allowedCategories.length > 0 && activeCategoryCode === null) {
            const firstCategory = allowedCategories
                .filter(c => c.is_visible && c.is_active)
                .sort((a, b) => (a.display_order || 0) - (b.display_order || 0))[0];
            if (firstCategory) {
                setActiveCategoryCode(firstCategory.code);
            }
        }
    }, [categoriesLoading, allowedCategories, activeCategoryCode]);

    // Return condicional después de todos los hooks
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

    const handleCategoryChange = (categoryCode: string | null) => {
        setActiveCategoryCode(categoryCode);
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
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold text-gray-900">Integraciones</h1>
                <Button
                    variant="primary"
                    onClick={() => setShowCreateModal(true)}
                >
                    {activeTab === 'integrations' ? 'Crear Integración' : 'Crear Tipo'}
                </Button>
            </div>

            {/* Category Tabs - navegación única */}
            {!categoriesLoading && (
                <CategoryTabs
                    categories={allowedCategories}
                    activeCategory={activeCategoryCode}
                    activeTab={activeTab}
                    onSelectCategory={(code) => {
                        setActiveTab('integrations');
                        handleCategoryChange(code);
                    }}
                    onSelectTypes={() => setActiveTab('types')}
                    canViewTypes={canViewTypes}
                />
            )}

            {activeTab === 'integrations' && activeCategoryCode !== null ? (
                <IntegrationList
                    key={`cat-${activeCategoryCode}-${refreshKey}`}
                    filterCategory={activeCategoryCode}
                    onEdit={async (integration) => {
                        try {
                            // Obtener la integración completa por ID
                            // Si el usuario es super admin, el endpoint devolverá credenciales desencriptadas automáticamente
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
            {activeTab === 'integrations' ? (
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
