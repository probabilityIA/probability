'use client';

import { useState, useRef, useEffect } from 'react';
import { Modal } from '@/shared/ui';
import ProductList from '@/services/modules/products/ui/components/ProductList';
import ProductForm from '@/services/modules/products/ui/components/ProductForm';
import ProductFamilyList, { ProductFamilyListHandle } from '@/services/modules/products/ui/components/ProductFamilyList';
import ProductFamilyForm from '@/services/modules/products/ui/components/ProductFamilyForm';
import ProductTour from '@/services/modules/products/ui/components/ProductTour';
import { CatalogPricingModal } from '@/services/modules/pricing/ui/components/CatalogPricingModal';
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
    const [showIntegrationMenu, setShowIntegrationMenu] = useState(false);
    const integrationMenuRef = useRef<HTMLDivElement>(null);
    const productListRef = useRef<any>(null);

    const integrationOptions = [
        { value: '', label: 'Todas las integraciones' },
        { value: 'shopify', label: 'Shopify' },
        { value: 'woocommerce', label: 'WooCommerce' },
        { value: 'whatsapp', label: 'WhatsApp' },
    ];

    const [isFamilyModalOpen, setIsFamilyModalOpen] = useState(false);
    const [selectedFamily, setSelectedFamily] = useState<ProductFamily | undefined>(undefined);
    const familyListRef = useRef<ProductFamilyListHandle>(null);

    const [isTourOpen, setIsTourOpen] = useState(false);
    const [pulseTour, setPulseTour] = useState(false);
    const [isPricingModalOpen, setIsPricingModalOpen] = useState(false);

    useEffect(() => {
        try {
            const seen = localStorage.getItem('products_tour_seen_v1');
            if (!seen) setPulseTour(true);
        } catch {}
    }, []);

    useEffect(() => {
        if (!showIntegrationMenu) return;
        const handleClickOutside = (e: MouseEvent) => {
            if (integrationMenuRef.current && !integrationMenuRef.current.contains(e.target as Node)) {
                setShowIntegrationMenu(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, [showIntegrationMenu]);

    const effectiveBusinessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;
    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

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
            ? 'btn-business-primary'
            : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
        }`;

    return (
        <div className="space-y-6 p-8">
            {!requiresBusinessSelection && (
                <div className="flex justify-between items-start gap-4">
                            <div className="flex gap-2 w-fit flex-shrink-0">
                                <button onClick={() => setActiveTab('products')} className={tabClass('products')}>
                                    SKUs / Productos
                                </button>
                                <button onClick={() => setActiveTab('families')} className={tabClass('families')}>
                                    Familias de variantes
                                </button>
                            </div>
                            {activeTab === 'products' ? (
                                <div className="flex-1 grid grid-cols-1 sm:grid-cols-2 gap-3">
                                    <div>
                                        <label className="block text-xs font-bold text-gray-700 dark:text-gray-200 mb-2.5">Nombre</label>
                                        <input
                                            type="text"
                                            placeholder="Ej: Camiseta..."
                                            value={searchName}
                                            onChange={e => setSearchName(e.target.value)}
                                            className="w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-business-primary focus:border-business-primary text-gray-900 dark:text-white placeholder:text-gray-400 bg-white dark:bg-gray-700 transition-all text-sm"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-xs font-bold text-gray-700 dark:text-gray-200 mb-2.5">SKU</label>
                                        <input
                                            type="text"
                                            placeholder="Ej: PROD-001..."
                                            value={searchSku}
                                            onChange={e => setSearchSku(e.target.value)}
                                            className="w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-business-primary focus:border-business-primary text-gray-900 dark:text-white placeholder:text-gray-400 bg-white dark:bg-gray-700 transition-all text-sm"
                                        />
                                    </div>
                                </div>
                            ) : (
                                <div className="flex-1 flex items-center">
                                    <p className="text-sm text-gray-500 dark:text-gray-400">
                                        Agrupa tus SKUs en familias para organizar variantes por color, talla u otro eje.
                                    </p>
                                </div>
                            )}
                            <div className="flex gap-3 flex-shrink-0">
                                <div className="relative" ref={integrationMenuRef}>
                                    <button
                                        onClick={() => setShowIntegrationMenu(v => !v)}
                                        title="Filtrar por integracion"
                                        className="w-14 h-14 flex items-center justify-center rounded-2xl bg-[#4fc3c9] hover:bg-[#3fb0b6] text-white shadow-sm transition-all duration-200"
                                    >
                                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z" />
                                        </svg>
                                    </button>
                                    {showIntegrationMenu && (
                                        <div className="absolute right-0 mt-2 w-56 bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 py-2 z-20">
                                            {integrationOptions.map(opt => (
                                                <button
                                                    key={opt.value}
                                                    onClick={() => { setSearchIntegration(opt.value); setShowIntegrationMenu(false); }}
                                                    className={`w-full text-left px-4 py-2 text-sm transition-colors ${searchIntegration === opt.value
                                                        ? 'bg-[#4fc3c9]/10 text-[#3fb0b6] font-semibold'
                                                        : 'text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700'
                                                        }`}
                                                >
                                                    {opt.label}
                                                </button>
                                            ))}
                                        </div>
                                    )}
                                </div>
                                <button
                                    onClick={() => setIsPricingModalOpen(true)}
                                    title="Grupos de clientes y precios por catalogo"
                                    className="p-3 rounded-lg border border-gray-300 dark:border-gray-600 text-business-primary bg-white dark:bg-gray-700 hover:border-business-primary shadow-sm transition-all duration-200"
                                >
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a4 4 0 00-3-3.87M9 20H4v-2a4 4 0 013-3.87m6-3.13a4 4 0 10-4-4 4 4 0 004 4zm6 0a3 3 0 10-3-3" />
                                    </svg>
                                </button>
                                <button
                                    onClick={() => { setIsTourOpen(true); setPulseTour(false); }}
                                    title={pulseTour ? 'Nuevo! Tutorial guiado de productos' : 'Tutorial guiado'}
                                    className={`p-3 rounded-lg border border-gray-300 dark:border-gray-600 text-business-primary bg-white dark:bg-gray-700 hover:border-business-primary shadow-sm transition-all duration-200 ${pulseTour ? 'animate-pulse' : ''}`}
                                >
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l9-5-9-5-9 5 9 5z" />
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" />
                                    </svg>
                                </button>
                                <button
                                    onClick={handleCreate}
                                    className="px-5 py-2.5 btn-business-primary text-white font-bold rounded-lg shadow-md transition-all duration-200 text-2xl"
                                >
                                    +
                                </button>
                            </div>
                        </div>
            )}
            {requiresBusinessSelection ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <svg className="w-12 h-12 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                    <p className="text-gray-500 dark:text-gray-400 text-sm">Selecciona un negocio para ver y gestionar sus productos</p>
                </div>
            ) : (
                <>
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

            <CatalogPricingModal
                isOpen={isPricingModalOpen}
                onClose={() => setIsPricingModalOpen(false)}
                businessId={effectiveBusinessId}
            />
        </div>
    );
}
