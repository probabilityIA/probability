'use client';

import { useState } from 'react';
import {
    IntegrationList,
    IntegrationForm,
    IntegrationTypeList,
    IntegrationTypeForm
} from '@/services/integrations/core/ui';
import { Button, Modal } from '@/shared/ui';
import { WideModal } from '@/shared/ui/wide-modal';
import { IntegrationType, Integration } from '@/services/integrations/core/domain/types';
import { getIntegrationByIdAction } from '@/services/integrations/core/infra/actions';

export default function IntegrationsPage() {
    const [activeTab, setActiveTab] = useState<'integrations' | 'types'>('integrations');
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showEditIntegrationModal, setShowEditIntegrationModal] = useState(false);
    const [selectedType, setSelectedType] = useState<IntegrationType | undefined>(undefined);
    const [selectedIntegration, setSelectedIntegration] = useState<Integration | undefined>(undefined);
    const [refreshKey, setRefreshKey] = useState(0);
    const [modalSize, setModalSize] = useState<'md' | '4xl' | '5xl' | '6xl' | 'full'>('5xl');

    const handleSuccess = () => {
        setShowCreateModal(false);
        setShowEditModal(false);
        setShowEditIntegrationModal(false);
        setSelectedType(undefined);
        setSelectedIntegration(undefined);
        setRefreshKey(prev => prev + 1);
        setModalSize('5xl'); // Reset to large when closing
    };

    const handleTypeSelected = (hasTypeSelected: boolean) => {
        setModalSize(hasTypeSelected ? 'full' : 'md');
    };

    const handleModalClose = () => {
        setShowCreateModal(false);
        setShowEditModal(false);
        setShowEditIntegrationModal(false);
        setSelectedType(undefined);
        setSelectedIntegration(undefined);
        setModalSize('5xl'); // Reset to large when closing
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

            {/* Tabs */}
            <div className="border-b border-gray-200 mb-6">
                <nav className="-mb-px flex space-x-8">
                    <button
                        onClick={() => setActiveTab('integrations')}
                        className={`
                            whitespace-nowrap pb-4 px-1 border-b-2 font-medium text-sm
                            ${activeTab === 'integrations'
                                ? 'border-blue-500 text-blue-600'
                                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                        `}
                    >
                        Mis Integraciones
                    </button>
                    <button
                        onClick={() => setActiveTab('types')}
                        className={`
                            whitespace-nowrap pb-4 px-1 border-b-2 font-medium text-sm
                            ${activeTab === 'types'
                                ? 'border-blue-500 text-blue-600'
                                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                        `}
                    >
                        Tipos de Integración
                    </button>
                </nav>
            </div>

            {activeTab === 'integrations' ? (
                <IntegrationList 
                    key={`list-${refreshKey}`}
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
                <WideModal
                    isOpen={showCreateModal}
                    onClose={handleModalClose}
                    title="Nueva Integración"
                    width="90vw"
                >
                    <IntegrationForm
                        onSuccess={handleSuccess}
                        onCancel={handleModalClose}
                        onTypeSelected={handleTypeSelected}
                    />
                </WideModal>
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
                    onTypeSelected={handleTypeSelected}
                />
            </Modal>
        </div>
    );
}
