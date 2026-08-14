'use client';

import { useState, useRef, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { Modal, INVENTORY_FILTERS_SLOT_ID } from '@/shared/ui';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui/dynamic-filters';
import ProductList from '@/services/modules/products/ui/components/ProductList';
import ProductForm from '@/services/modules/products/ui/components/ProductForm';
import ProductFamilyList, { ProductFamilyListHandle } from '@/services/modules/products/ui/components/ProductFamilyList';
import ProductFamilyForm from '@/services/modules/products/ui/components/ProductFamilyForm';
import ProductTour from '@/services/modules/products/ui/components/ProductTour';
import { ProductDetailModal } from '@/services/modules/products/ui/components/ProductDetailModal';
import { CatalogPricingModal } from '@/services/modules/pricing/ui/components/CatalogPricingModal';
import ProductDimensionsUploadModal from '@/services/modules/products/ui/components/ProductDimensionsUploadModal';
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
    const [viewMode, setViewMode] = useState<'create' | 'edit'>('create');
    const [detalleID, setDetalleID] = useState<string | null>(null);
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
    const [searchFamily, setSearchFamily] = useState('');
    const [familyFilters, setFamilyFilters] = useState<Record<string, string | boolean>>({});

    const familyFilterOptions: FilterOption[] = [
        {
            key: 'status', label: 'Estado', type: 'select', options: [
                { value: 'active', label: 'Activo' },
                { value: 'draft', label: 'Borrador' },
                { value: 'archived', label: 'Archivado' },
            ]
        },
        { key: 'category', label: 'Categoria', type: 'text', placeholder: 'Ej: Buzos' },
        { key: 'brand', label: 'Marca', type: 'text', placeholder: 'Ej: Viga' },
        { key: 'has_variants', label: 'Con variantes asociadas', type: 'boolean' },
    ];

    const familyActiveFilters: ActiveFilter[] = Object.entries(familyFilters)
        .filter(([, value]) => value !== undefined && value !== '')
        .map(([key, value]) => {
            const def = familyFilterOptions.find(f => f.key === key)!;
            return { key, label: def.label, value: value as any, type: def.type };
        });

    const handleAddFamilyFilter = (key: string, value: any) => {
        setFamilyFilters(prev => ({ ...prev, [key]: value }));
    };

    const handleRemoveFamilyFilter = (key: string) => {
        setFamilyFilters(prev => {
            const next = { ...prev };
            delete next[key];
            return next;
        });
    };

    const [isTourOpen, setIsTourOpen] = useState(false);
    const [pulseTour, setPulseTour] = useState(false);
    const [isPricingModalOpen, setIsPricingModalOpen] = useState(false);
    const [isDimensionsModalOpen, setIsDimensionsModalOpen] = useState(false);
    const [filtersSlot, setFiltersSlot] = useState<HTMLElement | null>(null);

    useEffect(() => {
        setFiltersSlot(document.getElementById(INVENTORY_FILTERS_SLOT_ID));
    }, []);

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
        setDetalleID(product.id);
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
        `h-9 px-4 flex items-center text-[13px] font-bold rounded-lg whitespace-nowrap transition-all duration-200 ${activeTab === tab
            ? 'btn-business-primary'
            : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
        }`;

    const inputClass = 'h-9 px-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-business-primary focus:border-business-primary text-gray-900 dark:text-white placeholder:text-gray-400 bg-white dark:bg-gray-700 transition-all text-[13px]';
    const iconBtnClass = 'h-9 w-9 flex-shrink-0 flex items-center justify-center rounded-lg border border-gray-300 dark:border-gray-600 text-business-primary bg-white dark:bg-gray-700 hover:border-business-primary shadow-sm transition-all duration-200';

    const toolbar = (
        <div className="flex items-center gap-2">
            <div className="flex gap-1.5 flex-shrink-0">
                <button onClick={() => setActiveTab('products')} className={tabClass('products')}>
                    SKUs / Productos
                </button>
                <button onClick={() => setActiveTab('families')} className={tabClass('families')}>
                    Familias de variantes
                </button>
            </div>

            {activeTab === 'products' ? (
                <div className="flex flex-shrink-0 items-center gap-2">
                    <input
                        type="text"
                        placeholder="Buscar por nombre"
                        value={searchName}
                        onChange={e => setSearchName(e.target.value)}
                        className={`${inputClass} w-52`}
                    />
                    <input
                        type="text"
                        placeholder="Buscar por SKU"
                        value={searchSku}
                        onChange={e => setSearchSku(e.target.value)}
                        className={`${inputClass} w-40`}
                    />
                    <DynamicFilters
                        variant="bar"
                        availableFilters={productFilterOptions}
                        activeFilters={productActiveFilters}
                        onAddFilter={handleAddProductFilter}
                        onRemoveFilter={handleRemoveProductFilter}
                        triggerClassName="!w-9 !h-9 !rounded-lg"
                    />
                </div>
            ) : (
                <div className="flex min-w-0 flex-shrink items-center gap-2">
                    <input
                        type="text"
                        placeholder="Buscar familia"
                        value={searchFamily}
                        onChange={e => setSearchFamily(e.target.value)}
                        className={`${inputClass} w-52`}
                    />
                    <DynamicFilters
                        variant="bar"
                        availableFilters={familyFilterOptions}
                        activeFilters={familyActiveFilters}
                        onAddFilter={handleAddFamilyFilter}
                        onRemoveFilter={handleRemoveFamilyFilter}
                        triggerClassName="!w-9 !h-9 !rounded-lg"
                    />
                    <p className="hidden min-w-0 truncate text-[13px] text-gray-500 dark:text-gray-400 2xl:block">
                        Agrupa tus SKUs en familias para organizar variantes por color, talla u otro eje.
                    </p>
                </div>
            )}

            <div className="ml-auto flex flex-shrink-0 items-center gap-2">
                <div className="relative flex items-center" ref={integrationMenuRef}>
                    <button
                        onClick={() => setShowIntegrationMenu(v => !v)}
                        title="Filtrar por integracion"
                        className="btn-business-primary h-9 w-9 flex items-center justify-center rounded-lg text-white shadow-sm transition-all duration-200"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z" />
                        </svg>
                    </button>
                    {showIntegrationMenu && (
                        <div className="absolute right-0 top-full mt-2 w-56 bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 py-2 z-50">
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
                    onClick={() => setIsDimensionsModalOpen(true)}
                    title="Cargar/descargar peso y dimensiones en Excel"
                    className={iconBtnClass}
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3M4 4h16v4H4V4zm0 8h16v8H4v-8z" />
                    </svg>
                </button>
                <button
                    onClick={() => setIsPricingModalOpen(true)}
                    title="Grupos de clientes y precios por catalogo"
                    className={iconBtnClass}
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a4 4 0 00-3-3.87M9 20H4v-2a4 4 0 013-3.87m6-3.13a4 4 0 10-4-4 4 4 0 004 4zm6 0a3 3 0 10-3-3" />
                    </svg>
                </button>
                <button
                    onClick={() => { setIsTourOpen(true); setPulseTour(false); }}
                    title={pulseTour ? 'Nuevo! Tutorial guiado de productos' : 'Tutorial guiado'}
                    className={`${iconBtnClass} ${pulseTour ? 'animate-pulse' : ''}`}
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l9-5-9-5-9 5 9 5z" />
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" />
                    </svg>
                </button>
                <button
                    onClick={handleCreate}
                    title="Crear"
                    className="h-9 w-9 flex-shrink-0 flex items-center justify-center btn-business-primary text-white font-bold rounded-lg shadow-md transition-all duration-200 text-xl leading-none"
                >
                    +
                </button>
            </div>
        </div>
    );

    return (
        <div className="space-y-6 p-8 pt-4">
            {!requiresBusinessSelection && (
                filtersSlot ? createPortal(toolbar, filtersSlot) : toolbar
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
                            searchName={searchFamily}
                            advancedFilters={familyFilters}
                        />
                    )}
                </>
            )}

            <Modal
                isOpen={isModalOpen}
                onClose={handleCloseModal}
                title={viewMode === 'create' ? 'Crear Producto' : 'Editar Producto'}
                size="4xl"
            >
                <div className="p-4">
                    <ProductForm
                        product={selectedProduct}
                        onSuccess={handleSuccess}
                        onCancel={handleCloseModal}
                        businessId={effectiveBusinessId}
                    />
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

            <ProductDetailModal
                productId={detalleID}
                businessId={effectiveBusinessId}
                onClose={() => setDetalleID(null)}
            />

            <ProductTour isOpen={isTourOpen} onClose={() => setIsTourOpen(false)} />

            <CatalogPricingModal
                isOpen={isPricingModalOpen}
                onClose={() => setIsPricingModalOpen(false)}
                businessId={effectiveBusinessId}
            />

            <ProductDimensionsUploadModal
                isOpen={isDimensionsModalOpen}
                onClose={() => setIsDimensionsModalOpen(false)}
                onUploadComplete={() => productListRef.current?.refreshProducts()}
                businessId={effectiveBusinessId}
            />
        </div>
    );
}
