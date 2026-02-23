'use client';

import { useState, useEffect, useCallback } from 'react';
import { Product, GetProductsParams } from '../../domain/types';
import { getProductsAction } from '../../infra/actions';
import { Button, Input, Alert } from '@/shared/ui';
import { Search, Plus, X, Package } from 'lucide-react';

interface ProductSelectorProps {
    businessId: number;
    onSelect: (products: Product[]) => void;
    selectedProducts: Product[];
    onCreateNew?: () => void;
}

export default function ProductSelector({
    businessId,
    onSelect,
    selectedProducts,
    onCreateNew
}: ProductSelectorProps) {
    const [searchTerm, setSearchTerm] = useState('');
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(false);
    const [showResults, setShowResults] = useState(false);

    const searchProducts = useCallback(async (term: string) => {
        if (!term) {
            setProducts([]);
            return;
        }

        setLoading(true);
        try {
            const params: GetProductsParams = {
                business_id: businessId,
                name: term,
                page: 1,
                page_size: 10
            };
            const response = await getProductsAction(params);
            if (response.success) {
                setProducts(response.data);
            }
        } catch (error) {
            console.error('Error searching products:', error);
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => {
        const timer = setTimeout(() => {
            if (searchTerm) searchProducts(searchTerm);
        }, 500);

        return () => clearTimeout(timer);
    }, [searchTerm, searchProducts]);

    const handleAdd = (product: Product) => {
        const alreadySelected = selectedProducts.find(p => p.id === product.id);
        if (!alreadySelected) {
            onSelect([...selectedProducts, product]);
        }
        setShowResults(false);
        setSearchTerm('');
    };

    const handleRemove = (productId: string) => {
        onSelect(selectedProducts.filter(p => p.id !== productId));
    };

    return (
        <div className="space-y-4">
            <div className="relative">
                <div className="flex gap-2">
                    <div className="relative flex-1">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <Search className="h-4 w-4 text-gray-400" />
                        </div>
                        <Input
                            placeholder="Buscar producto por nombre o SKU..."
                            value={searchTerm}
                            onChange={(e) => {
                                setSearchTerm(e.target.value);
                                setShowResults(true);
                            }}
                            className="pl-10"
                            onFocus={() => setShowResults(true)}
                        />
                    </div>
                    {onCreateNew && (
                        <Button
                            type="button"
                            onClick={onCreateNew}
                            variant="secondary"
                            className="flex items-center gap-1"
                            style={{ background: '#ede9fe', color: '#7c3aed', border: '1px solid #e9d5ff' }}
                        >
                            <Plus className="h-4 w-4" />
                            Nuevo
                        </Button>
                    )}
                </div>

                {/* Results dropdown */}
                {showResults && (searchTerm || loading) && (
                    <div className="absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-xl max-h-64 overflow-auto">
                        {loading ? (
                            <div className="p-4 text-center text-gray-500">Buscando productos...</div>
                        ) : products.length > 0 ? (
                            <ul className="py-1">
                                {products.map((product) => (
                                    <li key={product.id}>
                                        <button
                                            type="button"
                                            onClick={() => handleAdd(product)}
                                            className="w-full px-4 py-2 text-left hover:bg-blue-50 flex items-center justify-between"
                                        >
                                            <div className="flex items-center gap-3">
                                                {product.thumbnail ? (
                                                    <img src={product.thumbnail} alt="" className="w-10 h-10 rounded object-cover" />
                                                ) : (
                                                    <div className="w-10 h-10 bg-gray-100 rounded flex items-center justify-center">
                                                        <Package className="h-5 w-5 text-gray-400" />
                                                    </div>
                                                )}
                                                <div>
                                                    <div className="text-sm font-medium text-gray-900">{product.name}</div>
                                                    <div className="text-xs text-gray-500">SKU: {product.sku}</div>
                                                </div>
                                            </div>
                                            <div className="text-sm font-semibold text-blue-600">
                                                {product.currency} {product.price.toLocaleString()}
                                            </div>
                                        </button>
                                    </li>
                                ))}
                            </ul>
                        ) : (
                            <div className="p-4 text-center text-gray-500">
                                No se encontraron productos.
                                <button
                                    type="button"
                                    onClick={onCreateNew}
                                    className="ml-1 text-blue-600 font-medium hover:underline"
                                >
                                    Crear uno nuevo
                                </button>
                            </div>
                        )}
                    </div>
                )}
            </div>

            {/* Selected products list */}
            {selectedProducts.length > 0 && (
                <div className="border border-gray-200 rounded-lg overflow-hidden">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Producto</th>
                                <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Precio</th>
                                <th className="px-4 py-2 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Acciones</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {selectedProducts.map((product) => (
                                <tr key={product.id}>
                                    <td className="px-4 py-2 whitespace-nowrap">
                                        <div className="flex items-center">
                                            <div className="text-sm font-medium text-gray-900">{product.name}</div>
                                            <div className="ml-2 text-xs text-gray-500">({product.sku})</div>
                                        </div>
                                    </td>
                                    <td className="px-4 py-2 whitespace-nowrap text-right text-sm text-gray-900">
                                        {product.currency} {product.price.toLocaleString()}
                                    </td>
                                    <td className="px-4 py-2 whitespace-nowrap text-center">
                                        <button
                                            type="button"
                                            onClick={() => handleRemove(product.id)}
                                            className="text-red-600 hover:text-red-900"
                                        >
                                            <X className="h-4 w-4" />
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
}
