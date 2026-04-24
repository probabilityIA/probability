'use client';

import { useState, useRef, useEffect } from 'react';
import { Modal } from '@/shared/ui';
import ProductList from '@/services/modules/products/ui/components/ProductList';
import ProductForm from '@/services/modules/products/ui/components/ProductForm';
import ProductFamilyList, { ProductFamilyListHandle } from '@/services/modules/products/ui/components/ProductFamilyList';
import ProductFamilyForm from '@/services/modules/products/ui/components/ProductFamilyForm';
import ProductTour from '@/services/modules/products/ui/components/ProductTour';
import { Product, ProductFamily } from '@/services/modules/products/domain/types';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';

type Tab = 'products' | 'families';

export default function ProductsPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const [activeTab, setActiveTab] = useState<Tab>('products');

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedProduct, setSelectedProduct] = useState<Product | undefined>(undefined);
    const [viewMode, setViewMode] = useState<'create' | 'edit' | 'view'>('create');
    const [searchName, setSearchName] = useState('');
    const [searchSku, setSearchSku] = useState('');
    const [searchIntegration, setSearchIntegration] = useState('');
    const productListRef = useRef<any>(null);

    const [isFamilyModalOpen, setIsFamilyModalOpen] = useState(false);
    const [selectedFamily, setSelectedFamily] = useState<ProductFamily | undefined>(undefined);
    const familyListRef = useRef<ProductFamilyListHandle>(null);

    const [isTourOpen, setIsTourOpen] = useState(false);
    const [pulseTour, setPulseTour] = useState(false);

    useEffect(() => {
        try {
            const seen = localStorage.getItem('products_tour_seen_v1');
            if (!seen) setPulseTour(true);
        } catch {}
    }, []);

    const effectiveBusinessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;
    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    const handleRefresh = () => {
        if (activeTab === 'products') productListRef.current?.refreshProducts();
        else familyListRef.current?.refresh();
    };

    const handleCreate = () => {
        if (activeTab === 'products') {
            setSelectedProduct(undefined);
            setViewMode('create');
            setIsModalOpen(true);
        } else {
            setSelectedFamily(undefined);
            setIsFamilyModalOpen(true);
        }
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
        productListRef.current?.refreshProducts();
    };

    const handleEditFamily = (family: ProductFamily) => {
        setSelectedFamily(family);
        setIsFamilyModalOpen(true);
    };

    const handleCloseFamilyModal = () => {
        setIsFamilyModalOpen(false);
        setSelectedFamily(undefined);
    };

    const handleFamilySuccess = () => {
        handleCloseFamilyModal();
        familyListRef.current?.refresh();
    };

    const tabClass = (tab: Tab) =>
        `px-5 py-2.5 text-sm font-bold rounded-lg transition-all duration-200 ${activeTab === tab
            ? 'bg-[#7c3aed] text-white shadow-md'
            : 'text-slate-600 dark:text-slate-300 hover:bg-purple-50 dark:hover:bg-gray-700'
        }`;

    return (
        <div className="space-y-6 p-8">
            <div className="flex items-start gap-6">
                <div className="flex-shrink-0 min-w-fit">
                    <h1 className="text-4xl font-bold text-slate-900 dark:text-white mb-2">Productos</h1>
                    <p className="text-slate-600 dark:text-slate-300 text-base">
                        Gestiona el catalogo de productos de tu negocio
                    </p>
                </div>

                {!requiresBusinessSelection && (
                    <div className="bg-white dark:bg-gray-800 px-6 py-4 rounded-xl shadow-lg hover:shadow-xl border-2 border-[#7c3aed]/40 dark:border-[#7c3aed]/60 transition-all duration-300 backdrop-blur-sm flex-1">
                        <div className="flex justify-between items-end gap-4">
                            {activeTab === 'products' ? (
                                <div className="flex-1 grid grid-cols-1 sm:grid-cols-3 gap-3">
                                    <div>
                                        <label className="block text-xs font-bold text-slate-700 dark:text-slate-200 mb-2.5">Nombre</label>
                                        <input
                                            type="text"
                                            placeholder="Ej: Camiseta..."
                                            value={searchName}
                                            onChange={e => setSearchName(e.target.value)}
                                            className="w-full px-5 py-3 border-2 border-slate-200 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed] text-slate-900 dark:text-white placeholder:text-slate-400 bg-white dark:bg-gray-700 transition-all text-sm"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-xs font-bold text-slate-700 dark:text-slate-200 mb-2.5">SKU</label>
                                        <input
                                            type="text"
                                            placeholder="Ej: PROD-001..."
                                            value={searchSku}
                                            onChange={e => setSearchSku(e.target.value)}
                                            className="w-full px-5 py-3 border-2 border-slate-200 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed] text-slate-900 dark:text-white placeholder:text-slate-400 bg-white dark:bg-gray-700 transition-all text-sm"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-xs font-bold text-slate-700 dark:text-slate-200 mb-2.5">Integraciones</label>
                                        <select
                                            value={searchIntegration}
                                            onChange={e => setSearchIntegration(e.target.value)}
                                            className="w-full px-5 py-3 border-2 border-slate-200 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-[#7c3aed] text-slate-900 dark:text-white bg-white dark:bg-gray-700 transition-all text-sm"
                                        >
                                            <option value="">Todas las integraciones</option>
                                            <option value="shopify">Shopify</option>
                                            <option value="woocommerce">WooCommerce</option>
                                            <option value="whatsapp">WhatsApp</option>
                                        </select>
                                    </div>
                                </div>
                            ) : (
                                <div className="flex-1 flex items-center">
                                    <p className="text-sm text-slate-500 dark:text-slate-400">
                                        Agrupa tus SKUs en familias para organizar variantes por color, talla u otro eje.
                                    </p>
                                </div>
                            )}
                            <div className="flex gap-3 flex-shrink-0">
                                <button
                                    onClick={() => { setIsTourOpen(true); setPulseTour(false); }}
                                    title={pulseTour ? '¡Nuevo! Tutorial guiado de productos' : 'Tutorial guiado'}
                                    className={`p-3 rounded-lg border-2 transition-all duration-300 text-[#7c3aed] dark:text-[#a855f7] border-[#7c3aed]/40 hover:border-[#7c3aed] bg-white dark:bg-gray-700 shadow-sm hover:shadow-md ${pulseTour ? 'animate-pulse' : ''}`}
                                >
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l9-5-9-5-9 5 9 5z" />
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" />
                                    </svg>
                                </button>
                                <button
                                    onClick={handleRefresh}
                                    className="group px-6 py-3 bg-white dark:bg-gray-700 border-2 border-[#7c3aed]/40 hover:border-[#7c3aed] text-[#7c3aed] dark:text-[#a855f7] font-bold rounded-lg transition-all duration-300 text-sm shadow-sm hover:shadow-md flex items-center gap-2 whitespace-nowrap"
                                >
                                    <span className="inline-block transition-transform duration-500 group-hover:rotate-180">&#8635;</span>
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

            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                    <p className="text-gray-500 dark:text-gray-400 text-sm">Selecciona un negocio para ver y gestionar sus productos</p>
                </div>
            ) : (
                <>
                    <div className="flex gap-2 bg-slate-100 dark:bg-gray-800 p-1.5 rounded-xl w-fit">
                        <button onClick={() => setActiveTab('products')} className={tabClass('products')}>
                            SKUs / Productos
                        </button>
                        <button onClick={() => setActiveTab('families')} className={tabClass('families')}>
                            Familias de variantes
                        </button>
                    </div>

                    {activeTab === 'products' ? (
                        <ProductList
                            ref={productListRef}
                            onView={handleView}
                            onEdit={handleEdit}
                            searchName={searchName}
                            searchSku={searchSku}
                            searchIntegration={searchIntegration}
                            selectedBusinessId={effectiveBusinessId}
                        />
                    ) : (
                        <ProductFamilyList
                            ref={familyListRef}
                            onEdit={handleEditFamily}
                            selectedBusinessId={effectiveBusinessId}
                        />
                    )}
                </>
            )}

            <Modal
                isOpen={isModalOpen}
                onClose={handleCloseModal}
                title={viewMode === 'create' ? 'Crear Producto' : viewMode === 'edit' ? 'Editar Producto' : 'Detalles del Producto'}
                size="xl"
            >
                <div className="p-4">
                    {viewMode === 'view' && selectedProduct ? (
                        <div className="space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Nombre</label>
                                    <p className="text-gray-900 dark:text-white">{selectedProduct.name}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-gray-500 dark:text-gray-400">SKU</label>
                                    <p className="text-gray-900 dark:text-white">{selectedProduct.sku}</p>
                                </div>
                                {selectedProduct.barcode && (
                                    <div>
                                        <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Barcode</label>
                                        <p className="text-gray-900 dark:text-white font-mono">{selectedProduct.barcode}</p>
                                    </div>
                                )}
                                {selectedProduct.family && (
                                    <div>
                                        <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Familia</label>
                                        <p className="text-gray-900 dark:text-white">{selectedProduct.family.name}</p>
                                    </div>
                                )}
                                {selectedProduct.variant_label && (
                                    <div>
                                        <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Variante</label>
                                        <p className="text-gray-900 dark:text-white">{selectedProduct.variant_label}</p>
                                    </div>
                                )}
                                <div>
                                    <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Precio</label>
                                    <p className="text-gray-900 dark:text-white">
                                        {new Intl.NumberFormat('es-CO', { style: 'currency', currency: selectedProduct.currency }).format(selectedProduct.price)}
                                    </p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Stock</label>
                                    <p className="text-gray-900 dark:text-white">{selectedProduct.manage_stock ? selectedProduct.stock : '∞'}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Estado</label>
                                    <p className="text-gray-900 dark:text-white">{selectedProduct.status}</p>
                                </div>
                            </div>
                            {selectedProduct.description && (
                                <div>
                                    <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Descripcion</label>
                                    <p className="text-gray-900 dark:text-white whitespace-pre-wrap">{selectedProduct.description}</p>
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

            <Modal
                isOpen={isFamilyModalOpen}
                onClose={handleCloseFamilyModal}
                title={selectedFamily ? 'Editar Familia' : 'Nueva Familia de Variantes'}
                size="lg"
            >
                <div className="p-4">
                    <ProductFamilyForm
                        family={selectedFamily}
                        onSuccess={handleFamilySuccess}
                        onCancel={handleCloseFamilyModal}
                        businessId={effectiveBusinessId}
                    />
                </div>
            </Modal>

            <ProductTour isOpen={isTourOpen} onClose={() => setIsTourOpen(false)} />
        </div>
    );
}
