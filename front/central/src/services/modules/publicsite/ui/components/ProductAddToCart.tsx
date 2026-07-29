'use client';

import { useState } from 'react';
import { PublicProduct } from '../../domain/types';
import { useCart } from '../cart/cart-context';

export function ProductAddToCart({ product }: { product: PublicProduct }) {
    const { addItem, open } = useCart();
    const [quantity, setQuantity] = useState(1);
    const [added, setAdded] = useState(false);

    const handleAdd = () => {
        addItem(
            {
                product_id: product.id,
                name: product.name,
                sku: product.sku,
                price: product.price,
                image_url: product.image_url,
            },
            quantity
        );
        setAdded(true);
        setTimeout(() => setAdded(false), 1500);
    };

    return (
        <div className="space-y-3">
            <div className="flex items-center gap-3">
                <span className="text-sm text-gray-500 dark:text-gray-400">Cantidad</span>
                <div className="flex items-center gap-2">
                    <button
                        onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                        className="w-8 h-8 rounded border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300"
                    >
                        -
                    </button>
                    <span className="w-6 text-center">{quantity}</span>
                    <button
                        onClick={() => setQuantity((q) => q + 1)}
                        className="w-8 h-8 rounded border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300"
                    >
                        +
                    </button>
                </div>
            </div>
            <button
                onClick={handleAdd}
                disabled={product.stock_quantity <= 0}
                className="inline-block w-full text-center px-8 py-4 rounded-lg text-white font-bold text-lg transition-transform hover:scale-[1.02] disabled:opacity-50 disabled:hover:scale-100"
                style={{ backgroundColor: 'var(--brand-secondary)' }}
            >
                {added ? 'Agregado ✓' : 'Agregar al carrito'}
            </button>
            <button
                onClick={() => { handleAdd(); open(); }}
                disabled={product.stock_quantity <= 0}
                className="w-full text-center px-8 py-3 rounded-lg font-semibold border-2 disabled:opacity-50"
                style={{ borderColor: 'var(--brand-secondary)', color: 'var(--brand-secondary)' }}
            >
                Comprar ahora
            </button>
        </div>
    );
}
