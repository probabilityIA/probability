'use client';

import Image from 'next/image';
import { StorefrontProduct } from '../../domain/types';

interface ProductCardProps {
    product: StorefrontProduct;
    onAddToCart?: (product: StorefrontProduct) => void;
}

export function ProductCard({ product, onAddToCart }: ProductCardProps) {
    const formatPrice = (price: number, currency: string) => {
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: currency || 'COP',
            minimumFractionDigits: 0,
        }).format(price);
    };

    return (
        <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden hover:shadow-lg transition-shadow">
            <div className="relative aspect-square bg-gray-100 dark:bg-gray-700">
                {product.image_url ? (
                    <Image
                        src={product.image_url}
                        alt={product.name}
                        fill
                        className="object-cover"
                        sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 25vw"
                    />
                ) : (
                    <div className="flex items-center justify-center h-full text-gray-400">
                        <ShoppingBagIcon className="w-12 h-12" />
                    </div>
                )}
                {product.is_featured && (
                    <span className="absolute top-2 left-2 bg-indigo-600 text-white text-xs font-medium px-2 py-1 rounded">
                        Destacado
                    </span>
                )}
            </div>

            <div className="p-4">
                <h3 className="text-sm font-medium text-gray-900 dark:text-white dark:text-white line-clamp-2 mb-1">
                    {product.name}
                </h3>
                {product.category && (
                    <p className="text-xs text-gray-500 dark:text-gray-400 dark:text-gray-400 mb-2">{product.category}</p>
                )}
                <div className="flex items-center gap-2 mb-3">
                    <span className="text-lg font-bold text-gray-900 dark:text-white dark:text-white">
                        {formatPrice(product.price, product.currency)}
                    </span>
                    {product.compare_at_price && product.compare_at_price > product.price && (
                        <span className="text-sm text-gray-400 line-through">
                            {formatPrice(product.compare_at_price, product.currency)}
                        </span>
                    )}
                </div>

                {onAddToCart && (
                    <button
                        onClick={() => onAddToCart(product)}
                        disabled={product.stock_quantity <= 0}
                        className="w-full py-2 px-4 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
                    >
                        {product.stock_quantity > 0 ? 'Agregar al pedido' : 'Sin stock'}
                    </button>
                )}
            </div>
        </div>
    );
}

function ShoppingBagIcon({ className }: { className?: string }) {
    return (
        <svg className={className} fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 10.5V6a3.75 3.75 0 10-7.5 0v4.5m11.356-1.993l1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 01-1.12-1.243l1.264-12A1.125 1.125 0 015.513 7.5h12.974c.576 0 1.059.435 1.119 1.007zM8.625 10.5a.375.375 0 11-.75 0 .375.375 0 01.75 0zm7.5 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
        </svg>
    );
}
