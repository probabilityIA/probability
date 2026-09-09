'use client';

import Image from 'next/image';
import { ShoppingBagIcon, MinusIcon, PlusIcon } from '@heroicons/react/24/outline';
import { StorefrontProduct } from '../../domain/types';
import { useCart } from '../cart/cart-context';
import { isOptimizableProductImage } from '../utils/product-image';

interface ProductCardProps {
    product: StorefrontProduct;
}

const formatPrice = (price: number, currency: string) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: currency || 'COP', minimumFractionDigits: 0 }).format(price);

export function ProductCard({ product }: ProductCardProps) {
    const { items, addItem, incrementItem, decrementItem, maxQuantity } = useCart();
    const inCart = items.find(item => item.product.id === product.id);
    const outOfStock = product.track_inventory && product.stock_quantity <= 0;
    const atLimit = !!inCart && inCart.quantity >= maxQuantity(product);

    return (
        <div className="group bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 overflow-hidden hover:shadow-xl hover:-translate-y-0.5 transition-all duration-200">
            <div className="relative aspect-square bg-gray-50 dark:bg-gray-900 overflow-hidden">
                {product.image_url && isOptimizableProductImage(product.image_url) ? (
                    <Image
                        src={product.image_url}
                        alt={product.name}
                        fill
                        className="object-cover group-hover:scale-105 transition-transform duration-300"
                        sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 25vw"
                    />
                ) : product.image_url ? (
                    <img
                        src={product.image_url}
                        alt={product.name}
                        className="absolute inset-0 w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                ) : (
                    <div className="flex items-center justify-center h-full text-gray-300 dark:text-gray-600">
                        <ShoppingBagIcon className="w-12 h-12" />
                    </div>
                )}

                {product.is_featured && (
                    <span className="absolute top-3 left-3 bg-indigo-600 text-white text-xs font-semibold px-2.5 py-1 rounded-full shadow-sm">
                        Destacado
                    </span>
                )}

                {outOfStock && (
                    <div className="absolute inset-0 bg-white/70 dark:bg-gray-900/70 flex items-center justify-center">
                        <span className="bg-gray-900 dark:bg-white text-white dark:text-gray-900 text-xs font-semibold px-3 py-1.5 rounded-full">
                            Sin stock
                        </span>
                    </div>
                )}
            </div>

            <div className="p-4">
                {product.category && (
                    <p className="text-[11px] font-medium uppercase tracking-wide text-indigo-500 dark:text-indigo-400 mb-1">
                        {product.category}
                    </p>
                )}
                <h3 className="text-sm font-semibold text-gray-900 dark:text-white line-clamp-2 mb-2 min-h-[2.5rem]">
                    {product.name}
                </h3>

                <div className="flex items-baseline gap-2 mb-1">
                    <span className="text-lg font-bold text-gray-900 dark:text-white">
                        {formatPrice(product.price, product.currency)}
                    </span>
                    {product.compare_at_price && product.compare_at_price > product.price && (
                        <span className="text-xs text-gray-400 line-through">
                            {formatPrice(product.compare_at_price, product.currency)}
                        </span>
                    )}
                </div>

                {product.track_inventory && !outOfStock && (
                    <p className="text-xs text-gray-400 dark:text-gray-500 mb-3">{product.stock_quantity} disponibles</p>
                )}
                {!product.track_inventory && <div className="mb-3" />}

                {outOfStock ? (
                    <button
                        disabled
                        className="w-full py-2.5 px-4 bg-gray-100 dark:bg-gray-700 text-gray-400 dark:text-gray-500 text-sm font-medium rounded-xl cursor-not-allowed"
                    >
                        Sin stock
                    </button>
                ) : inCart ? (
                    <div className="flex items-center justify-between bg-indigo-50 dark:bg-indigo-900/30 rounded-xl px-2 py-1.5">
                        <button
                            onClick={() => decrementItem(product.id)}
                            className="w-8 h-8 flex items-center justify-center rounded-lg bg-white dark:bg-gray-800 text-indigo-600 dark:text-indigo-300 shadow-sm hover:bg-indigo-100 dark:hover:bg-gray-700 transition-colors"
                        >
                            <MinusIcon className="w-4 h-4" />
                        </button>
                        <span className="text-sm font-bold text-indigo-700 dark:text-indigo-300">{inCart.quantity}</span>
                        <button
                            onClick={() => incrementItem(product.id)}
                            disabled={atLimit}
                            className="w-8 h-8 flex items-center justify-center rounded-lg bg-white dark:bg-gray-800 text-indigo-600 dark:text-indigo-300 shadow-sm hover:bg-indigo-100 dark:hover:bg-gray-700 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
                        >
                            <PlusIcon className="w-4 h-4" />
                        </button>
                    </div>
                ) : (
                    <button
                        onClick={() => addItem(product)}
                        className="w-full py-2.5 px-4 bg-indigo-600 text-white text-sm font-semibold rounded-xl hover:bg-indigo-700 active:scale-[0.98] transition-all"
                    >
                        Agregar al carrito
                    </button>
                )}
            </div>
        </div>
    );
}
