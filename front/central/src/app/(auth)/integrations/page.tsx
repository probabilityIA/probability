'use client';

import { useState } from 'react';
import {
    IntegrationList,
    IntegrationForm,
    IntegrationTypeList,
    IntegrationTypeForm
} from '@/services/integrations/core/ui';
import { Button, Modal } from '@/shared/ui';

export default function IntegrationsPage() {
    const [activeTab, setActiveTab] = useState<'integrations' | 'types'>('integrations');
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [refreshKey, setRefreshKey] = useState(0);
    const [modalSize, setModalSize] = useState<'md' | 'full'>('md');

    const handleSuccess = () => {
        setShowCreateModal(false);
        setRefreshKey(prev => prev + 1);
        setModalSize('md'); // Reset to small when closing
    };

    const handleTypeSelected = (hasTypeSelected: boolean) => {
        setModalSize(hasTypeSelected ? 'full' : 'md');
    };

    const handleModalClose = () => {
        setShowCreateModal(false);
        setModalSize('md'); // Reset to small when closing
    };

    return (
        <div className="w-full px-6 py-8">
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
                <IntegrationList key={`list-${refreshKey}`} />
            ) : (
                <IntegrationTypeList key={`types-${refreshKey}`} />
            )}

            <Modal
                isOpen={showCreateModal}
                onClose={handleModalClose}
                title={activeTab === 'integrations' ? "Nueva Integración" : "Nuevo Tipo de Integración"}
                size={modalSize}
            >
                {activeTab === 'integrations' ? (
                    <IntegrationForm
                        onSuccess={handleSuccess}
                        onCancel={handleModalClose}
                        onTypeSelected={handleTypeSelected}
                    />
                ) : (
                    <IntegrationTypeForm
                        onSuccess={handleSuccess}
                        onCancel={handleModalClose}
                    />
                )}
            </Modal>
        </div>
    );
}
