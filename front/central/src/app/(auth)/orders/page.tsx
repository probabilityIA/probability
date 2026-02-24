'use client';

import { useState } from 'react';
import { getOrderByIdAction } from '@/services/modules/orders/infra/actions';
import { OrderList, OrderDetails, OrderForm } from '@/services/modules/orders/ui';
import { Order } from '@/services/modules/orders/domain/types';
import { Button, Modal } from '@/shared/ui';
import ShipmentGuideModal from '@/shared/ui/modals/shipment-guide-modal';
import MassOrderUploadModal from '@/shared/ui/modals/mass-order-upload-modal';
import MassGuideGenerationModal from '@/shared/ui/modals/mass-guide-generation-modal';


export default function OrdersPage() {
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showViewModal, setShowViewModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
    const [viewMode, setViewMode] = useState<'details' | 'recommendation'>('details');
    const [showTestGuideModal, setShowTestGuideModal] = useState(false);
    const [showGuideModal, setShowGuideModal] = useState(false);
    const [showMassUploadModal, setShowMassUploadModal] = useState(false);
    const [showMassGuideModal, setShowMassGuideModal] = useState(false);
    const [refreshKey, setRefreshKey] = useState(0);

    const handleView = (order: Order) => {
        setSelectedOrder(order);
        setViewMode('details'); // Set mode to details
        setShowViewModal(true);
    };

    const handleViewRecommendation = async (order: Order) => {
        try {
            const response = await getOrderByIdAction(order.id);
            if (response.success && response.data) {
                setSelectedOrder(response.data);
                setShowGuideModal(true);
            } else {
                alert('No se pudieron cargar los detalles de la orden');
            }
        } catch (error) {
            console.error('Error fetching order details:', error);
            alert('Error al cargar la orden completa');
        }
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
                <div className="flex gap-2">
                    <button
                        onClick={() => setShowCreateModal(true)}
                        style={{ background: '#7c3aed' }}
                        className="px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all"
                    >
                        ➕ Nueva Orden
                    </button>
                    <button
                        onClick={() => setShowMassUploadModal(true)}
                        style={{ background: '#7c3aed' }}
                        className="px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all"
                    >
                        📤 Carga Masiva
                    </button>
                    <button
                        onClick={() => setShowMassGuideModal(true)}
                        style={{ background: '#7c3aed' }}
                        className="px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all"
                    >
                        📦 Guías Masivas
                    </button>
                </div>
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
            {/* Test Guide Modal (No specific order) */}
            <ShipmentGuideModal
                isOpen={showTestGuideModal}
                onClose={() => setShowTestGuideModal(false)}
            />

            {/* Guide Modal (Selected order) */}
            <ShipmentGuideModal
                isOpen={showGuideModal}
                onClose={() => {
                    setShowGuideModal(false);
                    setSelectedOrder(null);
                }}
                order={selectedOrder || undefined}
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

            {/* Mass Upload Modal */}
            <MassOrderUploadModal
                isOpen={showMassUploadModal}
                onClose={() => setShowMassUploadModal(false)}
                onUploadComplete={(count) => {
                    setRefreshKey(prev => prev + 1);
                    setShowMassUploadModal(false);
                }}
            />

            {/* Mass Guide Generation Modal */}
            <MassGuideGenerationModal
                isOpen={showMassGuideModal}
                onClose={() => setShowMassGuideModal(false)}
                onComplete={(count) => {
                    setRefreshKey(prev => prev + 1);
                    setShowMassGuideModal(false);
                }}
            />
        </div>
    );
}
