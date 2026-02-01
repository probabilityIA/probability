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
    useIntegrations
} from '@/services/integrations/core/ui';
import { Button, Modal } from '@/shared/ui';
import { IntegrationType, Integration } from '@/services/integrations/core/domain/types';
import { getIntegrationByIdAction } from '@/services/integrations/core/infra/actions';

import { useSearchParams } from 'next/navigation';
import ShopifyOAuthCallback from '@/services/integrations/core/ui/components/shopify/ShopifyOAuthCallback';

export default function IntegrationsPage() {
    const searchParams = useSearchParams();
    const isOAuthCallback = searchParams.get('shopify_oauth');

    if (isOAuthCallback) {
        return <ShopifyOAuthCallback />;
    }

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
    const { setFilterCategory, refresh: refreshIntegrations } = useIntegrations();

    // Seleccionar primera categoría automáticamente cuando se cargan las categorías
    useEffect(() => {
        if (!categoriesLoading && categories.length > 0 && activeCategoryCode === null) {
            const firstCategory = categories
                .filter(c => c.is_visible && c.is_active)
                .sort((a, b) => (a.display_order || 0) - (b.display_order || 0))[0];
            if (firstCategory) {
                setActiveCategoryCode(firstCategory.code);
                setFilterCategory(firstCategory.code);
            }
        }
    }, [categories, categoriesLoading, activeCategoryCode]);

    // TODO: Obtener del contexto de usuario si es super admin
    const isSuperUser = true; // Por ahora en true, después conectar con el contexto de permisos

    const handleSuccess = () => {
        setShowCreateModal(false);
        setShowEditModal(false);
        setShowEditIntegrationModal(false);
        setSelectedType(undefined);
        setSelectedIntegration(undefined);
        setRefreshKey(prev => prev + 1);
        refreshIntegrations(); // Refresh integrations list
    };

    const handleCategoryChange = (categoryCode: string | null) => {
        setActiveCategoryCode(categoryCode);
        setFilterCategory(categoryCode || '');
        // No need to call refreshIntegrations() - the useEffect in useIntegrations will automatically
        // trigger when filterCategory changes
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
                    categories={categories}
                    activeCategory={activeCategoryCode}
                    activeTab={activeTab}
                    onSelectCategory={(code) => {
                        setActiveTab('integrations');
                        handleCategoryChange(code);
                    }}
                    onSelectTypes={() => setActiveTab('types')}
                    isSuperUser={isSuperUser}
                />
            )}

            {activeTab === 'integrations' ? (
                <IntegrationList
                    key={`list-${refreshKey}`}
                    filterCategory={activeCategoryCode || ''}
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
                    categories={categories}
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
