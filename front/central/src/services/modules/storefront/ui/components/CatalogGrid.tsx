'use client';

import { StorefrontProduct } from '../../domain/types';
import { ProductCard } from './ProductCard';

interface CatalogGridProps {
    products: StorefrontProduct[];
    onAddToCart?: (product: StorefrontProduct) => void;
}

export function CatalogGrid({ products, onAddToCart }: CatalogGridProps) {
    if (products.length === 0) {
        return (
            <div className="text-center py-12">
                <p className="text-gray-500 dark:text-gray-400 dark:text-gray-400 text-lg">No se encontraron productos</p>
            </div>
        );
    }

    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {products.map((product) => (
                <ProductCard key={product.id} product={product} onAddToCart={onAddToCart} />
            ))}
        </div>
    );
}
