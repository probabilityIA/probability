'use client';

import { useState } from 'react';
import { OrderStatusMappingList, OrderStatusMappingDetails, OrderStatusMappingForm, OrderStatusCatalogModal, ChannelStatusManager } from '@/services/modules/orderstatus/ui';
import { OrderStatusMapping } from '@/services/modules/orderstatus/domain/types';
import { Button, Modal } from '@/shared/ui';

export default function OrderStatusPage() {
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showViewModal, setShowViewModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showCatalogModal, setShowCatalogModal] = useState(false);
    const [showChannelStatusModal, setShowChannelStatusModal] = useState(false);
    const [selectedMapping, setSelectedMapping] = useState<OrderStatusMapping | null>(null);
    const [refreshKey, setRefreshKey] = useState(0);

    const handleView = (mapping: OrderStatusMapping) => {
        setSelectedMapping(mapping);
        setShowViewModal(true);
    };

    const handleEdit = (mapping: OrderStatusMapping) => {
        setSelectedMapping(mapping);
        setShowEditModal(true);
    };

    const handleSuccess = () => {
        setShowCreateModal(false);
        setShowEditModal(false);
        setShowViewModal(false);
        setSelectedMapping(null);
        setRefreshKey(prev => prev + 1);
    };

    const handleCancel = () => {
        setShowCreateModal(false);
        setShowEditModal(false);
        setShowViewModal(false);
        setSelectedMapping(null);
    };

    return (
        <div className="w-full px-6 py-8">
            <div className="flex justify-between items-center mb-6">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900">Order Status Mappings</h1>
                    <p className="text-gray-600 mt-1">
                        Gestiona los mapeos de estados de órdenes desde diferentes integraciones
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button
                        variant="outline"
                        onClick={() => setShowChannelStatusModal(true)}
                    >
                        Estados por Canal
                    </Button>
                    <Button
                        variant="outline"
                        onClick={() => setShowCatalogModal(true)}
                    >
                        Estados de Probability
                    </Button>
                    <Button
                        variant="primary"
                        onClick={() => setShowCreateModal(true)}
                    >
                        + Crear Mapping
                    </Button>
                </div>
            </div>

            <OrderStatusMappingList
                key={refreshKey}
                onView={handleView}
                onEdit={handleEdit}
            />

            {/* Create Modal */}
            <Modal
                isOpen={showCreateModal}
                onClose={handleCancel}
                title="Nuevo Mapping de Estado"
                size="lg"
            >
                <OrderStatusMappingForm
                    onSuccess={handleSuccess}
                    onCancel={handleCancel}
                />
            </Modal>

            {/* View Modal */}
            <Modal
                isOpen={showViewModal}
                onClose={handleCancel}
                title="Detalles del Mapping"
                size="lg"
            >
                {selectedMapping && <OrderStatusMappingDetails mapping={selectedMapping} />}
            </Modal>

            {/* Edit Modal */}
            <Modal
                isOpen={showEditModal}
                onClose={handleCancel}
                title="Editar Mapping de Estado"
                size="lg"
            >
                {selectedMapping && (
                    <OrderStatusMappingForm
                        mapping={selectedMapping}
                        onSuccess={handleSuccess}
                        onCancel={handleCancel}
                    />
                )}
            </Modal>

            {/* Catalog Modal */}
            <OrderStatusCatalogModal
                isOpen={showCatalogModal}
                onClose={() => setShowCatalogModal(false)}
            />

            {/* Channel Status Manager */}
            <ChannelStatusManager
                isOpen={showChannelStatusModal}
                onClose={() => setShowChannelStatusModal(false)}
            />
        </div>
    );
}
