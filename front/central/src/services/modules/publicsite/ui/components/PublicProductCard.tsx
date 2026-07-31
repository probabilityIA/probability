'use client';

import Link from 'next/link';
import { PublicProduct } from '../../domain/types';
import { useCart } from '../cart/cart-context';
import { useCustomerPrice } from '../session/prices-context';

interface PublicProductCardProps {
    product: PublicProduct;
    slug: string;
}

export function PublicProductCard({ product, slug }: PublicProductCardProps) {
    const { addItem } = useCart();
    const { price, isCustom } = useCustomerPrice(product.id, product.price);

    const handleAddToCart = (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        addItem({
            product_id: product.id,
            name: product.name,
            sku: product.sku,
            price,
            image_url: product.image_url,
        });
    };

    return (
        <Link
            href={`/tienda/${slug}/producto/${product.id}`}
            className="group bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-md transition-shadow"
        >
            <div className="aspect-square bg-gray-100 overflow-hidden">
                {product.image_url ? (
                    <img
                        src={product.image_url}
                        alt={product.name}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-300">
                        <svg className="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                    </div>
                )}
            </div>
            <div className="p-4">
                {product.brand && (
                    <p className="text-xs text-gray-400 uppercase tracking-wide mb-1">{product.brand}</p>
                )}
                <h3 className="font-medium text-gray-900 dark:text-white text-sm mb-2 line-clamp-2">{product.name}</h3>
                <div className="flex items-baseline gap-2">
                    <span className="font-bold text-lg" style={{ color: 'var(--brand-secondary)' }}>
                        ${price.toLocaleString('es-CO')}
                    </span>
                    {isCustom && (
                        <span className="text-sm text-gray-400 line-through">
                            ${product.price.toLocaleString('es-CO')}
                        </span>
                    )}
                    {!isCustom && product.compare_at_price && product.compare_at_price > product.price && (
                        <span className="text-sm text-gray-400 line-through">
                            ${product.compare_at_price.toLocaleString('es-CO')}
                        </span>
                    )}
                </div>
                {isCustom && (
                    <p className="text-[11px] font-medium mt-0.5" style={{ color: 'var(--brand-tertiary)' }}>Tu precio</p>
                )}
                {product.category && (
                    <p className="text-xs text-gray-400 mt-2">{product.category}</p>
                )}
                <button
                    onClick={handleAddToCart}
                    className="mt-3 w-full py-2 rounded-lg text-sm font-semibold text-white transition-opacity hover:opacity-90"
                    style={{ backgroundColor: 'var(--brand-secondary)' }}
                >
                    Agregar al carrito
                </button>
            </div>
        </Link>
    );
}
