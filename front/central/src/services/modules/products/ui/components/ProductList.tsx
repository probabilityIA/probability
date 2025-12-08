'use client';

import { useState, useEffect } from 'react';
import { getProductsAction, deleteProductAction } from '../../infra/actions';
import { Product, GetProductsParams } from '../../domain/types';
import { Button, Alert, Badge } from '@/shared/ui';
import ProductIntegrationsModal from './ProductIntegrationsModal';

interface ProductListProps {
    onView?: (product: Product) => void;
    onEdit?: (product: Product) => void;
}

export default function ProductList({ onView, onEdit }: ProductListProps) {
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);

    // Modal de integraciones
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
    const [isIntegrationsModalOpen, setIsIntegrationsModalOpen] = useState(false);

    // Filters
    const [filters, setFilters] = useState<GetProductsParams>({
        page: 1,
        page_size: 20,
    });

    useEffect(() => {
        fetchProducts();
    }, [filters]);

    const fetchProducts = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getProductsAction(filters);
            if (response.success && response.data) {
                setProducts(response.data);
                setTotal(response.total || 0);
                setTotalPages(response.total_pages || 1);
                setPage(response.page || 1);
            } else {
                setError(response.message || 'Error al cargar los productos');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar los productos');
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id: string) => {
        if (!confirm('¿Estás seguro de que deseas eliminar este producto?')) return;

        try {
            const response = await deleteProductAction(id);
            if (response.success) {
                fetchProducts();
            } else {
                alert(response.message || 'Error al eliminar el producto');
            }
        } catch (err: any) {
            alert(err.message || 'Error al eliminar el producto');
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
            {/* Filters */}
            <div className="bg-white p-4 sm:p-6 rounded-lg shadow-sm border border-gray-200">
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
                    <input
                        type="text"
                        placeholder="Buscar por nombre..."
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder:text-gray-500 bg-white"
                        onChange={(e) => setFilters({ ...filters, name: e.target.value || undefined })}
                    />
                    <input
                        type="text"
                        placeholder="Buscar por SKU..."
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder:text-gray-500 bg-white"
                        onChange={(e) => setFilters({ ...filters, sku: e.target.value || undefined })}
                    />
                    <select
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                        onChange={(e) => setFilters({ ...filters, integration_type: e.target.value || undefined })}
                    >
                        <option value="">Todas las integraciones</option>
                        <option value="shopify">Shopify</option>
                        <option value="woocommerce">WooCommerce</option>
                        <option value="whatsapp">WhatsApp</option>
                    </select>
                    <button
                        onClick={fetchProducts}
                        className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg transition-colors duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 w-full sm:w-auto"
                    >
                        🔄 Actualizar
                    </button>
                </div>
            </div>

            {/* Table */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Producto
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden sm:table-cell">
                                    SKU
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Precio
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Stock
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden lg:table-cell">
                                    Estado
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden md:table-cell">
                                    Fecha
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {products.length === 0 ? (
                                <tr>
                                    <td colSpan={7} className="px-4 sm:px-6 py-8 text-center text-gray-500">
                                        No hay productos disponibles
                                    </td>
                                </tr>
                            ) : (
                                products.map((product) => (
                                    <tr key={product.id} className="hover:bg-gray-50">
                                        <td className="px-3 sm:px-6 py-4">
                                            <div className="flex items-center">
                                                {product.thumbnail && (
                                                    <img src={product.thumbnail} alt={product.name} className="h-10 w-10 rounded-full mr-3 object-cover" />
                                                )}
                                                <div>
                                                    <div className="text-sm font-medium text-gray-900">
                                                        {product.name}
                                                    </div>
                                                    <div className="text-xs text-gray-500 sm:hidden">
                                                        {product.sku}
                                                    </div>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 hidden sm:table-cell">
                                            <div className="text-sm text-gray-900">{product.sku}</div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                                            <div className="text-sm font-semibold text-gray-900">
                                                {formatCurrency(product.price, product.currency)}
                                            </div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                                            <div className="text-sm text-gray-900">
                                                {product.manage_stock ? product.stock : '∞'}
                                            </div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap hidden lg:table-cell">
                                            <Badge type={product.is_active ? 'success' : 'secondary'}>
                                                {product.is_active ? 'Activo' : 'Inactivo'}
                                            </Badge>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500 hidden md:table-cell">
                                            {formatDate(product.created_at)}
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <div className="flex flex-col sm:flex-row justify-end gap-1 sm:gap-2">
                                                {onView && (
                                                    <button
                                                        onClick={() => onView(product)}
                                                        className="px-2 sm:px-3 py-1 sm:py-1.5 bg-blue-500 hover:bg-blue-600 text-white text-xs sm:text-sm font-medium rounded-md transition-colors duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                                                    >
                                                        Ver
                                                    </button>
                                                )}
                                                {onEdit && (
                                                    <button
                                                        onClick={() => onEdit(product)}
                                                        className="px-2 sm:px-3 py-1 sm:py-1.5 bg-yellow-500 hover:bg-yellow-600 text-white text-xs sm:text-sm font-medium rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                                                    >
                                                        Editar
                                                    </button>
                                                )}
                                                <button
                                                    onClick={() => {
                                                        setSelectedProduct(product);
                                                        setIsIntegrationsModalOpen(true);
                                                    }}
                                                    className="px-2 sm:px-3 py-1 sm:py-1.5 bg-purple-500 hover:bg-purple-600 text-white text-xs sm:text-sm font-medium rounded-md transition-colors duration-200 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
                                                    >
                                                        Integraciones
                                                    </button>
                                                <button
                                                    onClick={() => handleDelete(product.id)}
                                                    className="px-2 sm:px-3 py-1 sm:py-1.5 bg-red-500 hover:bg-red-600 text-white text-xs sm:text-sm font-medium rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                                                >
                                                    Eliminar
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {/* Pagination */}
                {(totalPages > 1 || total > 0) && (
                    <div className="bg-white px-3 sm:px-4 lg:px-6 py-3 flex flex-col sm:flex-row items-center justify-between gap-3 border-t border-gray-200">
                        {/* Mobile: Simple pagination */}
                        <div className="flex-1 flex justify-between sm:hidden w-full">
                            <Button
                                variant="outline"
                                onClick={() => setFilters({ ...filters, page: page - 1 })}
                                disabled={page === 1}
                                size="sm"
                            >
                                Anterior
                            </Button>
                            <Button
                                variant="outline"
                                onClick={() => setFilters({ ...filters, page: page + 1 })}
                                disabled={page === totalPages}
                                size="sm"
                            >
                                Siguiente
                            </Button>
                        </div>

                        {/* Desktop: Full pagination */}
                        <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between w-full">
                            <div className="flex items-center gap-3">
                                <p className="text-xs sm:text-sm text-gray-700">
                                    Mostrando <span className="font-medium">{(page - 1) * (filters.page_size || 20) + 1}</span> a{' '}
                                    <span className="font-medium">{Math.min(page * (filters.page_size || 20), total)}</span> de{' '}
                                    <span className="font-medium">{total}</span> resultados
                                </p>
                                <div className="flex items-center gap-2">
                                    <label className="text-xs sm:text-sm text-gray-700 whitespace-nowrap">
                                        Mostrar:
                                    </label>
                                    <select
                                        value={filters.page_size || 20}
                                        onChange={(e) => {
                                            const newPageSize = parseInt(e.target.value);
                                            setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                        }}
                                        className="px-2 py-1.5 text-xs sm:text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                                    >
                                        <option value="10">10</option>
                                        <option value="20">20</option>
                                        <option value="50">50</option>
                                        <option value="100">100</option>
                                    </select>
                                </div>
                            </div>
                            <div className="flex items-center gap-2">
                                <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page - 1 })}
                                        disabled={page === 1}
                                        className="relative inline-flex items-center px-2 sm:px-3 py-2 rounded-l-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Anterior
                                    </button>
                                    <span className="relative inline-flex items-center px-3 sm:px-4 py-2 border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-700">
                                        Página {page} de {totalPages}
                                    </span>
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page + 1 })}
                                        disabled={page === totalPages}
                                        className="relative inline-flex items-center px-2 sm:px-3 py-2 rounded-r-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Siguiente
                                    </button>
                                </nav>
                            </div>
                        </div>

                        {/* Mobile: Page size selector */}
                        <div className="flex items-center justify-between w-full sm:hidden pt-2 border-t border-gray-200">
                            <div className="flex items-center gap-2">
                                <label className="text-xs text-gray-700 whitespace-nowrap">
                                    Mostrar:
                                </label>
                                <select
                                    value={filters.page_size || 20}
                                    onChange={(e) => {
                                        const newPageSize = parseInt(e.target.value);
                                        setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                    }}
                                    className="px-2 py-1.5 text-xs border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                                >
                                    <option value="10">10</option>
                                    <option value="20">20</option>
                                    <option value="50">50</option>
                                    <option value="100">100</option>
                                </select>
                            </div>
                            <p className="text-xs text-gray-500">
                                Página {page} de {totalPages}
                            </p>
                        </div>
                    </div>
                )}
            </div>

            {/* Modal de Integraciones */}
            {selectedProduct && (
                <ProductIntegrationsModal
                    product={selectedProduct}
                    isOpen={isIntegrationsModalOpen}
                    onClose={() => {
                        setIsIntegrationsModalOpen(false);
                        setSelectedProduct(null);
                    }}
                    onSuccess={() => {
                        // Recargar productos si es necesario
                        fetchProducts();
                    }}
                />
            )}
        </div>
    );
}
