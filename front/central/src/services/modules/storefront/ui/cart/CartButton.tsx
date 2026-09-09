'use client';

import { ShoppingCartIcon } from '@heroicons/react/24/outline';
import { useCart } from './cart-context';

export function CartButton() {
    const { totalItems, toggle } = useCart();

    return (
        <button
            type="button"
            onClick={toggle}
            aria-label="Ver carrito"
            className="fixed bottom-6 right-6 z-40 flex items-center justify-center w-14 h-14 rounded-full bg-indigo-600 text-white shadow-lg shadow-indigo-600/30 hover:bg-indigo-700 hover:scale-105 active:scale-95 transition-all"
        >
            <ShoppingCartIcon className="w-6 h-6" />
            {totalItems > 0 && (
                <span className="absolute -top-1 -right-1 flex items-center justify-center min-w-[22px] h-[22px] px-1 rounded-full bg-red-500 text-white text-xs font-bold border-2 border-white dark:border-gray-900">
                    {totalItems > 99 ? '99+' : totalItems}
                </span>
            )}
        </button>
    );
}
