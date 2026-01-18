'use client';

import { useState } from 'react';
import { OrderList, OrderDetails, OrderForm } from '@/services/modules/orders/ui';
import { Order } from '@/services/modules/orders/domain/types';
import { Button, Modal } from '@/shared/ui';
import ShipmentGuideModal from '@/shared/ui/modals/shipment-guide-modal';


export default function OrdersPage() {
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showViewModal, setShowViewModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
    const [viewMode, setViewMode] = useState<'details' | 'recommendation'>('details'); // NEW state
    const [showTestGuideModal, setShowTestGuideModal] = useState(false);
    const [refreshKey, setRefreshKey] = useState(0);

    const handleView = (order: Order) => {
        setSelectedOrder(order);
        setViewMode('details'); // Set mode to details
        setShowViewModal(true);
    };

    const handleViewRecommendation = (order: Order) => { // NEW handler
        setSelectedOrder(order);
        setViewMode('recommendation'); // Set mode to recommendation
        setShowViewModal(true);
    };

    const handleEdit = (order: Order) => {
        setSelectedOrder(order);
        setShowEditModal(true);
    };

    const handleSuccess = () => {
        setShowCreateModal(false);
        setShowEditModal(false);
        setShowViewModal(false);
        setSelectedOrder(null);
        setRefreshKey(prev => prev + 1);
    };

    const handleCancel = () => {
        setShowCreateModal(false);
        setShowEditModal(false);
        setShowViewModal(false);
        setSelectedOrder(null);
    };

    return (
        <div className="min-h-screen bg-gray-50 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
                <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Órdenes</h1>
            </div>

            <OrderList
                refreshKey={refreshKey}
                onView={handleView}
                onViewRecommendation={handleViewRecommendation}
                onEdit={handleEdit}
                onCreate={() => setShowCreateModal(true)}
                onTestGuide={() => setShowTestGuideModal(true)}
            />

            {/* Test Guide Modal */}
            <ShipmentGuideModal
                isOpen={showTestGuideModal}
                onClose={() => setShowTestGuideModal(false)}
            />

            {/* Create Modal */}
            <Modal
                isOpen={showCreateModal}
                onClose={handleCancel}
                title="Nueva Orden"
                size="full"
            >
                <OrderForm
                    onSuccess={handleSuccess}
                    onCancel={handleCancel}
                />
            </Modal>

            {/* View Modal - Dynamic Title based on mode */}
            <Modal
                isOpen={showViewModal}
                onClose={handleCancel}
                title={viewMode === 'recommendation' ? undefined : 'Detalles de la Orden'} // Remove title for recommendation
                transparent={viewMode === 'recommendation'} // Enable transparent mode
                size="full"
            >
                {selectedOrder && (
                    <OrderDetails
                        initialOrder={selectedOrder}
                        onClose={handleCancel}
                        mode={viewMode} // Pass mode prop
                    />
                )}
            </Modal>

            {/* Edit Modal */}
            <Modal
                isOpen={showEditModal}
                onClose={handleCancel}
                title="Editar Orden"
                size="full"
            >
                {selectedOrder && (
                    <OrderForm
                        order={selectedOrder}
                        onSuccess={handleSuccess}
                        onCancel={handleCancel}
                    />
                )}
            </Modal>
        </div>
    );
}
