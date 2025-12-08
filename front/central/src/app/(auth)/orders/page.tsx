'use client';

import { useState } from 'react';
import { OrderList, OrderDetails, OrderForm } from '@/services/modules/orders/ui';
import { Order } from '@/services/modules/orders/domain/types';
import { Button, Modal } from '@/shared/ui';

export default function OrdersPage() {
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showViewModal, setShowViewModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
    const [refreshKey, setRefreshKey] = useState(0);

    const handleView = (order: Order) => {
        setSelectedOrder(order);
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
                <Button
                    variant="primary"
                    onClick={() => setShowCreateModal(true)}
                    className="w-full sm:w-auto"
                >
                    + Crear Orden
                </Button>
            </div>

            <OrderList
                key={refreshKey}
                onView={handleView}
                onEdit={handleEdit}
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

            {/* View Modal */}
            <Modal
                isOpen={showViewModal}
                onClose={handleCancel}
                title="Detalles de la Orden"
                size="2xl"
            >
                {selectedOrder && <OrderDetails order={selectedOrder} />}
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
