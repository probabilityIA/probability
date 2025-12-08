'use client';

import { useState } from 'react';
import { Button, Modal } from '@/shared/ui';
import ProductList from '@/services/modules/products/ui/components/ProductList';
import ProductForm from '@/services/modules/products/ui/components/ProductForm';
import { Product } from '@/services/modules/products/domain/types';

export default function ProductsPage() {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedProduct, setSelectedProduct] = useState<Product | undefined>(undefined);
    const [viewMode, setViewMode] = useState<'list' | 'create' | 'edit' | 'view'>('list');

    const handleCreate = () => {
        setSelectedProduct(undefined);
        setViewMode('create');
        setIsModalOpen(true);
    };

    const handleEdit = (product: Product) => {
        setSelectedProduct(product);
        setViewMode('edit');
        setIsModalOpen(true);
    };

    const handleView = (product: Product) => {
        setSelectedProduct(product);
        setViewMode('view');
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setSelectedProduct(undefined);
        setViewMode('list');
    };

    const handleSuccess = () => {
        handleCloseModal();
        // Trigger refresh in list (handled by useEffect dependency in ProductList)
        // For now, closing modal is enough as list auto-refreshes on mount/update
        // Ideally we'd pass a refresh trigger
        window.location.reload(); // Simple refresh for now
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Productos</h1>
                    <p className="text-sm text-gray-500">
                        Gestiona el catálogo de productos
                    </p>
                </div>
                <Button onClick={handleCreate}>
                    + Nuevo Producto
                </Button>
            </div>

            <ProductList
                onView={handleView}
                onEdit={handleEdit}
            />

            <Modal
                isOpen={isModalOpen}
                onClose={handleCloseModal}
                title={
                    viewMode === 'create' ? 'Crear Producto' :
                        viewMode === 'edit' ? 'Editar Producto' :
                            'Detalles del Producto'
                }
                size="xl"
            >
                <div className="p-4">
                    {viewMode === 'view' && selectedProduct ? (
                        <div className="space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="text-sm font-medium text-gray-500">Nombre</label>
                                    <p className="text-gray-900">{selectedProduct.name}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-gray-500">SKU</label>
                                    <p className="text-gray-900">{selectedProduct.sku}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-gray-500">Precio</label>
                                    <p className="text-gray-900">
                                        {new Intl.NumberFormat('es-CO', { style: 'currency', currency: selectedProduct.currency }).format(selectedProduct.price)}
                                    </p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-gray-500">Stock</label>
                                    <p className="text-gray-900">{selectedProduct.manage_stock ? selectedProduct.stock : '∞'}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-gray-500">Estado</label>
                                    <p className="text-gray-900">{selectedProduct.status}</p>
                                </div>
                            </div>
                            {selectedProduct.description && (
                                <div>
                                    <label className="text-sm font-medium text-gray-500">Descripción</label>
                                    <p className="text-gray-900 whitespace-pre-wrap">{selectedProduct.description}</p>
                                </div>
                            )}
                        </div>
                    ) : (
                        <ProductForm
                            product={selectedProduct}
                            onSuccess={handleSuccess}
                            onCancel={handleCloseModal}
                        />
                    )}
                </div>
            </Modal>
        </div>
    );
}
