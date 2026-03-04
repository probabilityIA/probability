'use client';

import { useState, useEffect } from 'react';
import { getOrderByIdAction } from '@/services/modules/orders/infra/actions';
import { deleteAllOrdersAction } from '@/services/modules/orders/infra/actions/testing-actions';
import { OrderList, OrderDetails, OrderForm, ShopifySimulatorModal } from '@/services/modules/orders/ui';
import { Order } from '@/services/modules/orders/domain/types';
import { Modal } from '@/shared/ui';
import ShipmentGuideModal from '@/shared/ui/modals/shipment-guide-modal';
import MassOrderUploadModal from '@/shared/ui/modals/mass-order-upload-modal';
import MassGuideGenerationModal from '@/shared/ui/modals/mass-guide-generation-modal';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useOrdersBusiness } from '@/shared/contexts/orders-business-context';
import { useToast } from '@/shared/providers/toast-provider';


export default function OrdersPage() {
    const { setActionButtons } = useNavbarActions();
    const { selectedBusinessId } = useOrdersBusiness();
    const { showToast } = useToast();
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showViewModal, setShowViewModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
    const [viewMode, setViewMode] = useState<'details' | 'recommendation'>('details');
    const [showTestGuideModal, setShowTestGuideModal] = useState(false);
    const [showGuideModal, setShowGuideModal] = useState(false);
    const [showMassUploadModal, setShowMassUploadModal] = useState(false);
    const [showMassGuideModal, setShowMassGuideModal] = useState(false);
    const [showShopifyModal, setShowShopifyModal] = useState(false);
    const [refreshKey, setRefreshKey] = useState(0);
    const [showDeleteAllModal, setShowDeleteAllModal] = useState(false);
    const [deleteStep, setDeleteStep] = useState<1 | 2 | 3>(1);
    const [deleteConfirmText, setDeleteConfirmText] = useState('');
    const [isDeletingAll, setIsDeletingAll] = useState(false);

    const handleDeleteAllOrders = async () => {
        if (!selectedBusinessId) return;
        setIsDeletingAll(true);
        try {
            const result = await deleteAllOrdersAction(selectedBusinessId);
            showToast(`${result.deleted} órdenes eliminadas correctamente`, 'success');
            setRefreshKey(prev => prev + 1);
        } catch (error: any) {
            showToast(error.message || 'Error al eliminar órdenes', 'error');
        } finally {
            setIsDeletingAll(false);
            setShowDeleteAllModal(false);
            setDeleteStep(1);
            setDeleteConfirmText('');
        }
    };

    // Set action buttons in navbar
    useEffect(() => {
        const actionButtons = (
            <>
                <button
                    onClick={() => setShowCreateModal(true)}
                    style={{ background: '#7c3aed' }}
                    className="px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all"
                >
                    + Nueva Orden
                </button>
                <button
                    onClick={() => setShowMassUploadModal(true)}
                    style={{ background: '#7c3aed' }}
                    className="px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all"
                >
                    Carga Masiva
                </button>
                <button
                    onClick={() => setShowMassGuideModal(true)}
                    style={{ background: '#7c3aed' }}
                    className="px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all"
                >
                    Guías Masivas
                </button>
                <button
                    onClick={() => setShowShopifyModal(true)}
                    style={{ background: '#059669' }}
                    className="px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all"
                >
                    Simular Shopify
                </button>
                <button
                    onClick={() => { setDeleteStep(1); setDeleteConfirmText(''); setShowDeleteAllModal(true); }}
                    style={{ background: '#dc2626' }}
                    className="px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all"
                >
                    Borrar todas
                </button>
            </>
        );
        setActionButtons(actionButtons);

        return () => setActionButtons(null);
    }, [setActionButtons]);

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

    const handleEdit = async (order: Order) => {
        try {
            const response = await getOrderByIdAction(order.id);
            if (response.success && response.data) {
                setSelectedOrder(response.data);
            } else {
                setSelectedOrder(order);
            }
        } catch {
            setSelectedOrder(order);
        }
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

            <OrderList
                refreshKey={refreshKey}
                onView={handleView}
                onViewRecommendation={handleViewRecommendation}
                onEdit={handleEdit}
                onCreate={() => setShowCreateModal(true)}
                onTestGuide={() => setShowTestGuideModal(true)}
                selectedBusinessId={selectedBusinessId}
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
                    selectedBusinessId={selectedBusinessId}
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
                        selectedBusinessId={selectedBusinessId}
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
                selectedBusinessId={selectedBusinessId}
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

            {/* Shopify Simulator Modal */}
            <ShopifySimulatorModal
                isOpen={showShopifyModal}
                onClose={() => setShowShopifyModal(false)}
                onSuccess={() => setRefreshKey(prev => prev + 1)}
            />

            {/* Delete All Orders Modal */}
            {showDeleteAllModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
                    <div className="bg-white rounded-xl shadow-2xl w-full max-w-md mx-4 p-6">
                        {deleteStep === 1 && (
                            <>
                                <div className="flex items-center gap-3 mb-4">
                                    <span className="text-3xl">⚠️</span>
                                    <h2 className="text-xl font-bold text-gray-900">Eliminar todas las órdenes</h2>
                                </div>
                                <p className="text-gray-600 mb-6">
                                    Esta acción eliminará <strong>permanentemente</strong> todas las órdenes del negocio seleccionado, incluyendo facturas, pagos, envíos y demás datos relacionados. Esta operación <strong>no se puede deshacer</strong>.
                                </p>
                                <div className="flex gap-3 justify-end">
                                    <button
                                        onClick={() => setShowDeleteAllModal(false)}
                                        className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
                                    >
                                        Cancelar
                                    </button>
                                    <button
                                        onClick={() => setDeleteStep(2)}
                                        className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
                                    >
                                        Continuar
                                    </button>
                                </div>
                            </>
                        )}

                        {deleteStep === 2 && (
                            <>
                                <div className="flex items-center gap-3 mb-4">
                                    <span className="text-3xl">🚨</span>
                                    <h2 className="text-xl font-bold text-gray-900">¿Estás seguro?</h2>
                                </div>
                                <p className="text-gray-600 mb-6">
                                    Se eliminarán <strong>todas</strong> las órdenes del negocio. Esta acción es irreversible y afectará también las facturas, pagos y envíos asociados.
                                </p>
                                <div className="flex gap-3 justify-end">
                                    <button
                                        onClick={() => setShowDeleteAllModal(false)}
                                        className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
                                    >
                                        Cancelar
                                    </button>
                                    <button
                                        onClick={() => setDeleteStep(3)}
                                        className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
                                    >
                                        Sí, continuar
                                    </button>
                                </div>
                            </>
                        )}

                        {deleteStep === 3 && (
                            <>
                                <div className="flex items-center gap-3 mb-4">
                                    <span className="text-3xl">🔴</span>
                                    <h2 className="text-xl font-bold text-gray-900">Confirmación final</h2>
                                </div>
                                <p className="text-gray-600 mb-4">
                                    Escribe <strong className="text-red-600">ELIMINAR</strong> para confirmar:
                                </p>
                                <input
                                    type="text"
                                    value={deleteConfirmText}
                                    onChange={(e) => setDeleteConfirmText(e.target.value)}
                                    placeholder="ELIMINAR"
                                    className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm mb-6 focus:outline-none focus:ring-2 focus:ring-red-500"
                                    autoFocus
                                />
                                <div className="flex gap-3 justify-end">
                                    <button
                                        onClick={() => setShowDeleteAllModal(false)}
                                        className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
                                        disabled={isDeletingAll}
                                    >
                                        Cancelar
                                    </button>
                                    <button
                                        onClick={handleDeleteAllOrders}
                                        disabled={deleteConfirmText !== 'ELIMINAR' || isDeletingAll}
                                        className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        {isDeletingAll ? 'Eliminando...' : 'Eliminar todo'}
                                    </button>
                                </div>
                            </>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}
