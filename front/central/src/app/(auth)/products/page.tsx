'use client';

import { useState, useRef } from 'react';
import { Modal } from '@/shared/ui';
import ProductList from '@/services/modules/products/ui/components/ProductList';
import ProductForm from '@/services/modules/products/ui/components/ProductForm';
import { Product } from '@/services/modules/products/domain/types';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';

export default function ProductsPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedProduct, setSelectedProduct] = useState<Product | undefined>(undefined);
    const [viewMode, setViewMode] = useState<'create' | 'edit' | 'view'>('create');
    const [searchName, setSearchName] = useState('');
    const [searchSku, setSearchSku] = useState('');
    const [searchIntegration, setSearchIntegration] = useState('');
    const productListRef = useRef<any>(null);

    const handleRefresh = () => {
        productListRef.current?.refreshProducts();
    };

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
    };

    const handleSuccess = () => {
        handleCloseModal();
        handleRefresh();
    };

    // Gate para super admin: debe seleccionar negocio antes de operar
    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;
    const effectiveBusinessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    return (
        <div className="space-y-8 p-8">
            {/* Header: Título a la izquierda, Filtros a la derecha */}
            <div className="flex items-start gap-6">
                {/* Título y Descripción */}
                <div className="flex-shrink-0 min-w-fit">
                    <h1 className="text-4xl font-bold bg-gradient-to-r from-slate-900 to-slate-700 bg-clip-text text-transparent mb-2">
                        Productos
                    </h1>
                    <p className="text-slate-600 text-base">
                        Gestiona el catálogo de productos de tu negocio
                    </p>
                </div>

                {/* Filtros y Botones (solo cuando no requiere selección de negocio) */}
                {!requiresBusinessSelection && (
                    <div className="bg-gradient-to-br from-[#7c3aed]/8 to-[#6d28d9]/8 px-6 py-4 rounded-xl shadow-lg hover:shadow-xl border-2 border-[#7c3aed]/40 transition-all duration-300 backdrop-blur-sm flex-1">
                        <div className="flex justify-between items-end gap-4">
                            <div className="flex-1 grid grid-cols-1 sm:grid-cols-3 gap-3">
                                <div>
                                    <label className="block text-xs font-bold text-slate-700 mb-2.5">Nombre</label>
                                    <input
                                        type="text"
                                        placeholder="Ej: Camiseta..."
                                        value={searchName}
                                        onChange={(e) => setSearchName(e.target.value)}
                                        className="w-full px-5 py-3 border-2 border-slate-200 rounded-lg focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed] focus:shadow-lg focus:shadow-[#7c3aed]/20 text-slate-900 placeholder:text-slate-400 bg-white transition-all duration-300 hover:border-slate-300 text-sm"
                                    />
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-700 mb-2.5">SKU</label>
                                    <input
                                        type="text"
                                        placeholder="Ej: PROD-001..."
                                        value={searchSku}
                                        onChange={(e) => setSearchSku(e.target.value)}
                                        className="w-full px-5 py-3 border-2 border-slate-200 rounded-lg focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed] focus:shadow-lg focus:shadow-[#7c3aed]/20 text-slate-900 placeholder:text-slate-400 bg-white transition-all duration-300 hover:border-slate-300 text-sm"
                                    />
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-700 mb-2.5">Integraciones</label>
                                    <select
                                        value={searchIntegration}
                                        onChange={(e) => setSearchIntegration(e.target.value)}
                                        className="w-full px-5 py-3 border-2 border-slate-200 rounded-lg focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed] focus:shadow-lg focus:shadow-[#7c3aed]/20 text-slate-900 bg-white transition-all duration-300 hover:border-slate-300 text-sm"
                                    >
                                        <option value="">Todas las integraciones</option>
                                        <option value="shopify">Shopify</option>
                                        <option value="woocommerce">WooCommerce</option>
                                        <option value="whatsapp">WhatsApp</option>
                                    </select>
                                </div>
                            </div>
                            <div className="flex gap-3 flex-shrink-0">
                                <button
                                    onClick={handleRefresh}
                                    className="group px-6 py-3 bg-gradient-to-r from-[#a855f7]/10 to-[#9333ea]/10 border-2 border-[#7c3aed]/40 hover:border-[#7c3aed] hover:from-[#7c3aed]/20 hover:to-[#6d28d9]/20 text-[#7c3aed] font-bold rounded-lg transition-all duration-300 text-sm shadow-sm hover:shadow-md transform hover:scale-105 flex items-center gap-2 whitespace-nowrap"
                                >
                                    <span className="inline-block transition-transform duration-500 group-hover:rotate-180">↻</span>
                                    Actualizar
                                </button>
                                <button
                                    onClick={handleCreate}
                                    className="px-5 py-3 bg-gradient-to-r from-[#7c3aed] to-[#6d28d9] hover:from-[#6d28d9] hover:to-[#5b21b6] text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-110 hover:-translate-y-1 text-2xl font-black"
                                >
                                    +
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </div>

            {/* Gate: super admin debe seleccionar negocio */}
            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                    <p className="text-gray-500 text-sm">Selecciona un negocio para ver y gestionar sus productos</p>
                </div>
            ) : (
                <ProductList
                    ref={productListRef}
                    onView={handleView}
                    onEdit={handleEdit}
                    searchName={searchName}
                    searchSku={searchSku}
                    searchIntegration={searchIntegration}
                    selectedBusinessId={effectiveBusinessId}
                />
            )}

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
                            businessId={effectiveBusinessId}
                        />
                    )}
                </div>
            </Modal>
        </div>
    );
}
