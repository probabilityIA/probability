'use client';

import { useState, useRef, useEffect } from 'react';
import { Modal } from '@/shared/ui';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui/dynamic-filters';
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

    const [productFilters, setProductFilters] = useState<Record<string, string | boolean>>({});

    const productFilterOptions: FilterOption[] = [
        {
            key: 'status', label: 'Estado', type: 'select', options: [
                { value: 'active', label: 'Activo' },
                { value: 'draft', label: 'Borrador' },
                { value: 'archived', label: 'Archivado' },
            ]
        },
        { key: 'category', label: 'Categoria', type: 'text', placeholder: 'Ej: Calzado' },
        { key: 'brand', label: 'Marca', type: 'text', placeholder: 'Ej: Nike' },
        { key: 'has_family', label: 'Con familia de variantes', type: 'boolean' },
        { key: 'price_min', label: 'Precio minimo', type: 'text', placeholder: 'Ej: 10000' },
        { key: 'price_max', label: 'Precio maximo', type: 'text', placeholder: 'Ej: 200000' },
        { key: 'stock_min', label: 'Stock minimo', type: 'text', placeholder: 'Ej: 1' },
        { key: 'stock_max', label: 'Stock maximo', type: 'text', placeholder: 'Ej: 100' },
        { key: 'weight_min', label: 'Peso minimo (kg)', type: 'text', placeholder: 'Ej: 0.5' },
        { key: 'weight_max', label: 'Peso maximo (kg)', type: 'text', placeholder: 'Ej: 5' },
    ];

    const productActiveFilters: ActiveFilter[] = Object.entries(productFilters)
        .filter(([, value]) => value !== undefined && value !== '')
        .map(([key, value]) => {
            const def = productFilterOptions.find(f => f.key === key)!;
            return { key, label: def.label, value: value as any, type: def.type };
        });

    const handleAddProductFilter = (key: string, value: any) => {
        setProductFilters(prev => ({ ...prev, [key]: value }));
    };

    const handleRemoveProductFilter = (key: string) => {
        setProductFilters(prev => {
            const next = { ...prev };
            delete next[key];
            return next;
        });
    };

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
        `h-11 px-5 flex items-center text-sm font-bold rounded-lg transition-all duration-200 ${activeTab === tab
            ? 'btn-business-primary'
            : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
        }`;

    return (
        <div className="space-y-6 p-8">
            {!requiresBusinessSelection && (
                <div className="flex justify-between items-end gap-4">
                            <div className="flex gap-2 w-fit flex-shrink-0 self-end">
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
                                            className="w-full h-11 px-4 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-business-primary focus:border-business-primary text-gray-900 dark:text-white placeholder:text-gray-400 bg-white dark:bg-gray-700 transition-all text-sm"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-xs font-bold text-gray-700 dark:text-gray-200 mb-2.5">SKU</label>
                                        <input
                                            type="text"
                                            placeholder="Ej: PROD-001..."
                                            value={searchSku}
                                            onChange={e => setSearchSku(e.target.value)}
                                            className="w-full h-11 px-4 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-business-primary focus:border-business-primary text-gray-900 dark:text-white placeholder:text-gray-400 bg-white dark:bg-gray-700 transition-all text-sm"
                                        />
                                    </div>
                                </div>
                            ) : (
                                <div className="flex-1 flex items-end pb-2.5">
                                    <p className="text-sm text-gray-500 dark:text-gray-400">
                                        Agrupa tus SKUs en familias para organizar variantes por color, talla u otro eje.
                                    </p>
                                </div>
                            )}
                            <div className="flex gap-3 flex-shrink-0 self-end">
                                {activeTab === 'products' && (
                                    <div className="h-11 flex items-center">
                                        <DynamicFilters
                                            variant="bar"
                                            availableFilters={productFilterOptions}
                                            activeFilters={productActiveFilters}
                                            onAddFilter={handleAddProductFilter}
                                            onRemoveFilter={handleRemoveProductFilter}
                                            triggerClassName="!w-11 !h-11 !rounded-lg"
                                        />
                                    </div>
                                )}
                                <div className="relative flex items-center" ref={integrationMenuRef}>
                                    <button
                                        onClick={() => setShowIntegrationMenu(v => !v)}
                                        title="Filtrar por integracion"
                                        className="btn-business-primary h-11 w-11 flex items-center justify-center rounded-lg text-white shadow-sm transition-all duration-200"
                                    >
                                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
                                                        ? 'btn-business-primary-soft font-semibold'
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
                                    className="h-11 w-11 flex items-center justify-center rounded-lg border border-gray-300 dark:border-gray-600 text-business-primary bg-white dark:bg-gray-700 hover:border-business-primary shadow-sm transition-all duration-200"
                                >
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a4 4 0 00-3-3.87M9 20H4v-2a4 4 0 013-3.87m6-3.13a4 4 0 10-4-4 4 4 0 004 4zm6 0a3 3 0 10-3-3" />
                                    </svg>
                                </button>
                                <button
                                    onClick={() => { setIsTourOpen(true); setPulseTour(false); }}
                                    title={pulseTour ? 'Nuevo! Tutorial guiado de productos' : 'Tutorial guiado'}
                                    className={`h-11 w-11 flex items-center justify-center rounded-lg border border-gray-300 dark:border-gray-600 text-business-primary bg-white dark:bg-gray-700 hover:border-business-primary shadow-sm transition-all duration-200 ${pulseTour ? 'animate-pulse' : ''}`}
                                >
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l9-5-9-5-9 5 9 5z" />
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" />
                                    </svg>
                                </button>
                                <button
                                    onClick={handleCreate}
                                    className="h-11 w-11 flex items-center justify-center btn-business-primary text-white font-bold rounded-lg shadow-md transition-all duration-200 text-2xl leading-none"
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
                            advancedFilters={productFilters}
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
                size="4xl"
            >
                <div className="p-4">
                    {viewMode === 'view' && selectedProduct ? (
                        <div className="space-y-5">
                            <div className="flex items-center gap-4 rounded-xl border border-gray-200 dark:border-gray-700 p-4" style={{ background: 'linear-gradient(135deg, color-mix(in srgb, var(--color-primary) 10%, transparent), transparent)' }}>
                                {selectedProduct.image_url ? (
                                    <img src={selectedProduct.image_url} alt={selectedProduct.name} className="w-16 h-16 rounded-xl object-cover border border-gray-200 dark:border-gray-700 flex-shrink-0" />
                                ) : (
                                    <div className="w-16 h-16 rounded-xl bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 flex items-center justify-center flex-shrink-0">
                                        <svg className="w-7 h-7 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>
                                    </div>
                                )}
                                <div className="min-w-0 flex-1">
                                    <h3 className="text-lg font-bold text-gray-900 dark:text-white truncate">{selectedProduct.name}</h3>
                                    <p className="text-sm font-mono text-gray-500 dark:text-gray-400">{selectedProduct.sku}</p>
                                </div>
                                <span className={`flex-shrink-0 px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wide ${selectedProduct.status === 'active'
                                    ? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300'
                                    : selectedProduct.status === 'archived'
                                        ? 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300'
                                        : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'
                                    }`}>
                                    {selectedProduct.status === 'active' ? 'Activo' : selectedProduct.status === 'archived' ? 'Archivado' : selectedProduct.status === 'draft' ? 'Borrador' : selectedProduct.status}
                                </span>
                            </div>

                            <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
                                <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
                                    <p className="text-[11px] font-bold text-gray-400 uppercase tracking-wide">Precio</p>
                                    <p className="text-base font-semibold text-gray-900 dark:text-white mt-0.5">
                                        {new Intl.NumberFormat('es-CO', { style: 'currency', currency: selectedProduct.currency }).format(selectedProduct.price)}
                                    </p>
                                </div>
                                <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
                                    <p className="text-[11px] font-bold text-gray-400 uppercase tracking-wide">Stock</p>
                                    <p className="text-base font-semibold text-gray-900 dark:text-white mt-0.5">{selectedProduct.manage_stock ? selectedProduct.stock : '∞'}</p>
                                </div>
                                {selectedProduct.barcode && (
                                    <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
                                        <p className="text-[11px] font-bold text-gray-400 uppercase tracking-wide">Barcode</p>
                                        <p className="text-base font-semibold text-gray-900 dark:text-white mt-0.5 font-mono">{selectedProduct.barcode}</p>
                                    </div>
                                )}
                                {selectedProduct.family && (
                                    <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
                                        <p className="text-[11px] font-bold text-gray-400 uppercase tracking-wide">Familia</p>
                                        <p className="text-base font-semibold text-gray-900 dark:text-white mt-0.5 truncate">{selectedProduct.family.name}</p>
                                    </div>
                                )}
                                {selectedProduct.variant_label && (
                                    <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
                                        <p className="text-[11px] font-bold text-gray-400 uppercase tracking-wide">Variante</p>
                                        <p className="text-base font-semibold text-gray-900 dark:text-white mt-0.5">{selectedProduct.variant_label}</p>
                                    </div>
                                )}
                            </div>

                            {selectedProduct.description && (
                                <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4">
                                    <p className="text-[11px] font-bold text-gray-400 uppercase tracking-wide mb-1">Descripcion</p>
                                    <p className="text-sm text-gray-700 dark:text-gray-200 whitespace-pre-wrap">{selectedProduct.description}</p>
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
                size="4xl"
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
