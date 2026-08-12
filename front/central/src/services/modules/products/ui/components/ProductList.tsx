'use client';

import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { getProductsAction, deleteProductAction, updateProductAction } from '../../infra/actions';
import { Product, GetProductsParams } from '../../domain/types';
import { Alert } from '@/shared/ui';
import ProductIntegrationsModal from './ProductIntegrationsModal';
import { getActionError } from '@/shared/utils/action-result';

interface ProductListProps {
    onView?: (product: Product) => void;
    onEdit?: (product: Product) => void;
    searchName?: string;
    searchSku?: string;
    searchIntegration?: string;
    advancedFilters?: Record<string, string | boolean>;
    selectedBusinessId?: number;
}

const ProductList = forwardRef(function ProductList(
    { onView, onEdit, searchName = '', searchSku = '', searchIntegration = '', advancedFilters = {}, selectedBusinessId }: ProductListProps,
    ref: any
) {
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);

    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
    const [isIntegrationsModalOpen, setIsIntegrationsModalOpen] = useState(false);

    const [filters, setFilters] = useState<GetProductsParams>({
        page: 1,
        page_size: 20,
    });

    useEffect(() => {
        const numericFields = ['price_min', 'price_max', 'stock_min', 'stock_max', 'weight_min', 'weight_max'];
        const advanced: Partial<GetProductsParams> = {};
        for (const [key, value] of Object.entries(advancedFilters)) {
            if (numericFields.includes(key)) {
                const num = Number(value);
                if (!isNaN(num) && num >= 0) (advanced as any)[key] = num;
            } else {
                (advanced as any)[key] = value;
            }
        }

        setFilters(prev => ({
            ...prev,
            name: searchName || undefined,
            sku: searchSku || undefined,
            integration_type: searchIntegration || undefined,
            status: undefined,
            category: undefined,
            brand: undefined,
            has_family: undefined,
            price_min: undefined,
            price_max: undefined,
            stock_min: undefined,
            stock_max: undefined,
            weight_min: undefined,
            weight_max: undefined,
            ...advanced,
            page: 1,
        }));
    }, [searchName, searchSku, searchIntegration, advancedFilters]);

    useEffect(() => {
        setFilters(prev => ({ ...prev, page: 1 }));
    }, [selectedBusinessId]);

    useEffect(() => {
        fetchProducts();
    }, [filters]);

    useImperativeHandle(ref, () => ({
        refreshProducts: fetchProducts,
    }));

    const fetchProducts = async () => {
        setLoading(true);
        setError(null);
        try {
            const params = { ...filters };
            if (selectedBusinessId) params.business_id = selectedBusinessId;
            const response = await getProductsAction(params);
            if (response.success && response.data) {
                setProducts(response.data);
                setTotal(response.total || 0);
                setTotalPages(response.total_pages || 1);
                setPage(response.page || 1);
            } else {
                setError(response.message || 'Error al cargar los productos');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar los productos'));
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id: string) => {
        if (!confirm('¿Estás seguro de que deseas eliminar este producto?')) return;

        try {
            const response = await deleteProductAction(id, selectedBusinessId);
            if (response.success) {
                fetchProducts();
            } else {
                alert(response.message || 'Error al eliminar el producto');
            }
        } catch (err: any) {
            alert(err.message || 'Error al eliminar el producto');
        }
    };

    const handleToggleActive = async (product: Product) => {
        try {
            const response = await updateProductAction(product.id, { is_active: !product.is_active }, selectedBusinessId);
            if (response.success) {
                setProducts(prev => prev.map(p => p.id === product.id ? { ...p, is_active: !p.is_active } : p));
            } else {
                alert(response.message || 'Error al actualizar el producto');
            }
        } catch (err: any) {
            alert(err.message || 'Error al actualizar el producto');
        }
    };

    const formatCurrency = (amount: number, currency: string = 'USD') => {
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency,
        }).format(amount);
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('es-CO', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    if (loading) {
        return <div className="text-center py-8">Cargando productos...</div>;
    }

    if (error) {
        return (
            <Alert type="error" onClose={() => setError(null)}>
                {error}
            </Alert>
        );
    }

    return (
        <div className="space-y-4">
            <div className="relative rounded-xl overflow-hidden shadow-sm bg-white dark:bg-gray-800">
                <div className="overflow-x-auto">
                    <table className="min-w-full">
                        <thead style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}>
                            <tr>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest">
                                    Producto
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden sm:table-cell">
                                    SKU
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest">
                                    Precio
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest">
                                    Stock
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden lg:table-cell">
                                    Familia
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden lg:table-cell">
                                    Inventario
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden lg:table-cell">
                                    Estado
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest hidden md:table-cell">
                                    Fecha
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-right text-xs font-bold text-white uppercase tracking-widest">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {products.length === 0 ? (
                                <tr>
                                    <td colSpan={9} className="px-4 sm:px-6 py-8 text-center text-gray-500 dark:text-gray-400">
                                        No hay productos disponibles
                                    </td>
                                </tr>
                            ) : (
                                products.map((product) => (
                                    <tr key={product.id} className="bg-white dark:bg-gray-800 border-b border-gray-100 dark:border-gray-700 hover:bg-purple-50 dark:hover:bg-gray-700 transition-colors">
                                        <td className="px-3 sm:px-6 py-4">
                                            <div className="flex items-center">
                                                {product.image_url || product.family?.image_url ? (
                                                    <img src={product.image_url || product.family?.image_url || ''} alt={product.name} className="h-10 w-10 rounded-full mr-3 object-cover" />
                                                ) : (
                                                    <div className="h-10 w-10 rounded-full mr-3 bg-gray-100 flex items-center justify-center text-gray-400 text-xs">N/A</div>
                                                )}
                                                <div>
                                                    <div className="text-sm font-medium text-gray-900 dark:text-white">
                                                        {product.name}
                                                    </div>
                                                    <div className="text-xs text-gray-500 dark:text-gray-400 sm:hidden">
                                                        {product.sku}
                                                    </div>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 hidden sm:table-cell">
                                            <div className="text-sm text-gray-900 dark:text-white">{product.sku}</div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                                            <div className="flex flex-col gap-0.5">
                                                <div className="text-sm font-semibold text-gray-900 dark:text-white">
                                                    {formatCurrency(product.compare_at_price ?? product.price, product.currency)}
                                                </div>
                                                {product.compare_at_price && product.compare_at_price !== product.price && (
                                                    <div className="text-xs text-red-600 dark:text-red-400 line-through">
                                                        {formatCurrency(product.price, product.currency)}
                                                    </div>
                                                )}
                                            </div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                                            <div className="text-sm text-gray-900 dark:text-white">
                                                {product.track_inventory ? (product.stock_quantity ?? product.stock ?? 0) : '∞'}
                                            </div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden lg:table-cell">
                                            {product.family ? (
                                                <div className="flex flex-col gap-0.5">
                                                    <span className="inline-flex items-center gap-1 px-2 py-0.5 text-xs font-semibold rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 max-w-[140px] truncate">
                                                        <svg className="w-3 h-3 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20"><path d="M7 3a1 1 0 000 2h6a1 1 0 100-2H7zM4 7a1 1 0 011-1h10a1 1 0 110 2H5a1 1 0 01-1-1zM2 11a2 2 0 012-2h12a2 2 0 012 2v4a2 2 0 01-2 2H4a2 2 0 01-2-2v-4z" /></svg>
                                                        <span className="truncate">{product.family.name}</span>
                                                    </span>
                                                    {product.variant_label && (
                                                        <span className="text-xs text-gray-400 dark:text-gray-500 pl-1 truncate max-w-[140px]">{product.variant_label}</span>
                                                    )}
                                                </div>
                                            ) : (
                                                <span className="text-xs text-gray-400 dark:text-gray-500">Sin familia</span>
                                            )}
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden lg:table-cell">
                                            <span className={`inline-flex px-2 py-0.5 text-xs font-medium rounded-full ${
                                                product.track_inventory
                                                    ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300'
                                                    : 'bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400'
                                            }`}>
                                                {product.track_inventory ? 'Habilitado' : 'No'}
                                            </span>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden lg:table-cell">
                                            <button
                                                onClick={() => handleToggleActive(product)}
                                                className={`inline-flex px-3 py-1 text-xs font-semibold rounded-full transition-colors duration-200 cursor-pointer ${
                                                    product.is_active
                                                        ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 hover:bg-green-200 dark:hover:bg-green-800'
                                                        : 'bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-200 hover:bg-red-200 dark:hover:bg-red-800'
                                                }`}
                                                title={product.is_active ? 'Click para desactivar' : 'Click para activar'}
                                            >
                                                {product.is_active ? 'Activo' : 'Inactivo'}
                                            </button>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400 hidden md:table-cell">
                                            {formatDate(product.created_at)}
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <div className="flex flex-row justify-end gap-1.5">
                                                {onView && (
                                                    <button
                                                        onClick={() => onView(product)}
                                                        title="Ver detalle"
                                                        className="p-1.5 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors duration-200"
                                                    >
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                                                    </button>
                                                )}
                                                {onEdit && (
                                                    <button
                                                        onClick={() => onEdit(product)}
                                                        title="Editar"
                                                        className="p-1.5 bg-amber-500 hover:bg-amber-600 text-white rounded-md transition-colors duration-200"
                                                    >
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                                                    </button>
                                                )}
                                                <button
                                                    onClick={() => { setSelectedProduct(product); setIsIntegrationsModalOpen(true); }}
                                                    title="Integraciones"
                                                    className="p-1.5 btn-business-primary rounded-md transition-colors duration-200"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" /></svg>
                                                </button>
                                                <button
                                                    onClick={() => handleDelete(product.id)}
                                                    title="Eliminar"
                                                    className="p-1.5 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {(totalPages > 1 || total > 0) && (
                    <div className="px-3 py-2 flex flex-col sm:flex-row items-center justify-between gap-2" style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}>
                        <div className="flex items-center gap-3">
                            <p className="text-[11px] text-white">
                                Mostrando <span className="font-medium">{(page - 1) * (filters.page_size || 20) + 1}</span> a{' '}
                                <span className="font-medium">{Math.min(page * (filters.page_size || 20), total)}</span> de{' '}
                                <span className="font-medium">{total}</span> resultados
                            </p>
                            <div className="flex items-center gap-1">
                                <label className="text-[11px] text-white whitespace-nowrap">Mostrar:</label>
                                <select
                                    value={filters.page_size || 20}
                                    onChange={(e) => {
                                        const newPageSize = parseInt(e.target.value);
                                        setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                    }}
                                    className="px-1.5 py-1 text-[11px] border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 dark:text-white bg-white dark:bg-gray-800"
                                >
                                    <option value="10">10</option>
                                    <option value="20">20</option>
                                    <option value="50">50</option>
                                    <option value="100">100</option>
                                </select>
                            </div>
                        </div>
                        <div className="flex items-center gap-1">
                            <nav className="relative z-0 inline-flex items-center gap-1 flex-wrap justify-end">
                                <button
                                    onClick={() => setFilters({ ...filters, page: 1 })}
                                    disabled={page === 1}
                                    className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                    title="Primera pagina"
                                >
                                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 19l-7-7 7-7m8 14l-7-7 7-7" /></svg>
                                </button>

                                <button
                                    onClick={() => setFilters({ ...filters, page: page - 1 })}
                                    disabled={page === 1}
                                    className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                    title="Anterior"
                                >
                                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" /></svg>
                                </button>

                                {(() => {
                                    const pages: (number | string)[] = [];
                                    const maxVisible = 8;

                                    if (totalPages <= maxVisible) {
                                        for (let i = 1; i <= totalPages; i++) pages.push(i);
                                    } else {
                                        const windowSize = maxVisible - 2;
                                        let start = Math.max(2, page - Math.floor(windowSize / 2));
                                        let end = start + windowSize - 1;
                                        if (end >= totalPages) {
                                            end = totalPages - 1;
                                            start = Math.max(2, end - windowSize + 1);
                                        }

                                        pages.push(1);
                                        if (start > 2) pages.push('...');
                                        for (let i = start; i <= end; i++) pages.push(i);
                                        if (end < totalPages - 1) pages.push('...');
                                        pages.push(totalPages);
                                    }

                                    return pages.map((p, idx) =>
                                        typeof p === 'string' ? (
                                            <span key={`ellipsis-${idx}`} className="relative inline-flex items-center px-1 py-1 text-[11px] font-bold text-white/70">
                                                ...
                                            </span>
                                        ) : (
                                            <button
                                                key={p}
                                                onClick={() => setFilters({ ...filters, page: p })}
                                                className={`relative inline-flex items-center justify-center min-w-7 px-2 py-1 rounded-md border text-[11px] font-semibold transition-all ${p === page ? 'page-btn-active z-10 shadow-md scale-105' : 'page-btn'
                                                    }`}
                                                style={p === page
                                                    ? { backgroundColor: 'var(--color-secondary)', borderColor: 'var(--color-secondary)', color: 'white' }
                                                    : undefined}
                                            >
                                                {p}
                                            </button>
                                        )
                                    );
                                })()}

                                <button
                                    onClick={() => setFilters({ ...filters, page: page + 1 })}
                                    disabled={page === totalPages}
                                    className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                    title="Siguiente"
                                >
                                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
                                </button>

                                <button
                                    onClick={() => setFilters({ ...filters, page: totalPages })}
                                    disabled={page === totalPages}
                                    className="page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all"
                                    title="Ultima pagina"
                                >
                                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 5l7 7-7 7M5 5l7 7-7 7" /></svg>
                                </button>
                            </nav>
                        </div>
                    </div>
                )}
            </div>

            {selectedProduct && (
                <ProductIntegrationsModal
                    product={selectedProduct}
                    isOpen={isIntegrationsModalOpen}
                    onClose={() => {
                        setIsIntegrationsModalOpen(false);
                        setSelectedProduct(null);
                    }}
                    onSuccess={() => {
                        fetchProducts();
                    }}
                />
            )}

        </div>
    );
});

export default ProductList;
