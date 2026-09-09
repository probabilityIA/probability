import { ShoppingBagIcon } from '@heroicons/react/24/outline';
import { StorefrontProduct } from '../../domain/types';
import { ProductCard } from './ProductCard';

interface CatalogGridProps {
    products: StorefrontProduct[];
    columns?: number;
}

export function CatalogGrid({ products, columns = 4 }: CatalogGridProps) {
    if (products.length === 0) {
        return (
            <div className="flex flex-col items-center justify-center text-center py-20">
                <ShoppingBagIcon className="w-12 h-12 text-gray-300 dark:text-gray-700 mb-3" />
                <p className="text-gray-500 dark:text-gray-400 text-lg font-medium">No se encontraron productos</p>
                <p className="text-gray-400 dark:text-gray-500 text-sm mt-1">Prueba con otra busqueda o categoria</p>
            </div>
        );
    }

    const safeColumns = Math.min(Math.max(columns, 1), 6);

    return (
        <>
            <div className="catalog-grid" style={{ ['--catalog-cols' as string]: safeColumns }}>
                {products.map((product) => (
                    <ProductCard key={product.id} product={product} />
                ))}
            </div>
            <style>{`
                .catalog-grid { display: grid; gap: 1rem; grid-template-columns: repeat(2, minmax(0, 1fr)); }
                @media (min-width: 640px) { .catalog-grid { grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 1.25rem; } }
                @media (min-width: 1024px) { .catalog-grid { grid-template-columns: repeat(var(--catalog-cols), minmax(0, 1fr)); } }
            `}</style>
        </>
    );
}
