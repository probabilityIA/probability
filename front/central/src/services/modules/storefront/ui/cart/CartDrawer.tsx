'use client';

import { useState } from 'react';
import Image from 'next/image';
import { XMarkIcon, MinusIcon, PlusIcon, TrashIcon, ShoppingBagIcon, MapPinIcon } from '@heroicons/react/24/outline';
import { useCart } from './cart-context';
import { isOptimizableProductImage } from '../utils/product-image';
import { createOrderAction } from '../../infra/actions';
import { useStorefrontBusiness } from '@/shared/contexts/storefront-business-context';
import { useToast } from '@/shared/providers/toast-provider';

const formatPrice = (price: number, currency: string) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: currency || 'COP', minimumFractionDigits: 0 }).format(price);

export function CartDrawer() {
    const { items, isOpen, close, totalPrice, currency, incrementItem, decrementItem, removeItem, clear, maxQuantity } = useCart();
    const { selectedBusinessId } = useStorefrontBusiness();
    const { showToast } = useToast();

    const [showAddress, setShowAddress] = useState(false);
    const [notes, setNotes] = useState('');
    const [address, setAddress] = useState({ first_name: '', last_name: '', phone: '', street: '', city: '', state: '' });
    const [submitting, setSubmitting] = useState(false);

    const handleCheckout = async () => {
        if (items.length === 0) return;
        setSubmitting(true);

        const result = await createOrderAction(
            {
                items: items.map(item => ({ product_id: item.product.id, quantity: item.quantity })),
                notes: notes || undefined,
                address: address.street ? address : undefined,
            },
            selectedBusinessId ?? undefined,
        );

        setSubmitting(false);

        if ('success' in result && result.success === false) {
            showToast(result.message || 'No se pudo enviar el pedido', 'error');
            return;
        }

        showToast('Pedido enviado correctamente', 'success');
        clear();
        setNotes('');
        setAddress({ first_name: '', last_name: '', phone: '', street: '', city: '', state: '' });
        close();
    };

    return (
        <>
            <div
                onClick={close}
                className={`fixed inset-0 bg-black/40 backdrop-blur-sm z-40 transition-opacity ${
                    isOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'
                }`}
            />

            <aside
                className={`fixed top-0 right-0 h-full w-full max-w-md bg-white dark:bg-gray-900 shadow-2xl z-50 flex flex-col transition-transform duration-300 ease-out ${
                    isOpen ? 'translate-x-0' : 'translate-x-full'
                }`}
            >
                <div className="flex items-center justify-between px-5 py-4 border-b border-gray-200 dark:border-gray-800">
                    <div className="flex items-center gap-2">
                        <ShoppingBagIcon className="w-5 h-5 text-indigo-600" />
                        <h2 className="text-lg font-bold text-gray-900 dark:text-white">Tu carrito</h2>
                        {items.length > 0 && (
                            <span className="text-xs font-medium bg-indigo-100 dark:bg-indigo-900/40 text-indigo-700 dark:text-indigo-300 px-2 py-0.5 rounded-full">
                                {items.length}
                            </span>
                        )}
                    </div>
                    <button onClick={close} className="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
                        <XMarkIcon className="w-5 h-5" />
                    </button>
                </div>

                <div className="flex-1 overflow-y-auto px-5 py-4">
                    {items.length === 0 ? (
                        <div className="h-full flex flex-col items-center justify-center text-center py-16">
                            <ShoppingBagIcon className="w-12 h-12 text-gray-300 dark:text-gray-700 mb-3" />
                            <p className="text-sm text-gray-500 dark:text-gray-400">Tu carrito esta vacio</p>
                            <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">Agrega productos desde el catalogo</p>
                        </div>
                    ) : (
                        <ul className="space-y-4">
                            {items.map(item => (
                                <li key={item.product.id} className="flex gap-3">
                                    <div className="relative w-16 h-16 rounded-lg overflow-hidden bg-gray-100 dark:bg-gray-800 flex-shrink-0">
                                        {item.product.image_url && isOptimizableProductImage(item.product.image_url) ? (
                                            <Image src={item.product.image_url} alt={item.product.name} fill className="object-cover" sizes="64px" />
                                        ) : item.product.image_url ? (
                                            <img src={item.product.image_url} alt={item.product.name} className="absolute inset-0 w-full h-full object-cover" />
                                        ) : (
                                            <div className="w-full h-full flex items-center justify-center text-gray-300 dark:text-gray-600">
                                                <ShoppingBagIcon className="w-6 h-6" />
                                            </div>
                                        )}
                                    </div>

                                    <div className="flex-1 min-w-0">
                                        <p className="text-sm font-medium text-gray-900 dark:text-white truncate">{item.product.name}</p>
                                        <p className="text-sm text-gray-500 dark:text-gray-400">{formatPrice(item.product.price, item.product.currency)}</p>

                                        <div className="flex items-center gap-2 mt-2">
                                            <button
                                                onClick={() => decrementItem(item.product.id)}
                                                className="w-6 h-6 flex items-center justify-center rounded border border-gray-300 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800"
                                            >
                                                <MinusIcon className="w-3.5 h-3.5" />
                                            </button>
                                            <span className="w-6 text-center text-sm font-medium text-gray-900 dark:text-white">{item.quantity}</span>
                                            <button
                                                onClick={() => incrementItem(item.product.id)}
                                                disabled={item.quantity >= maxQuantity(item.product)}
                                                className="w-6 h-6 flex items-center justify-center rounded border border-gray-300 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 disabled:opacity-30 disabled:cursor-not-allowed"
                                            >
                                                <PlusIcon className="w-3.5 h-3.5" />
                                            </button>

                                            <button
                                                onClick={() => removeItem(item.product.id)}
                                                className="ml-auto p-1 text-gray-300 hover:text-red-500 transition-colors"
                                            >
                                                <TrashIcon className="w-4 h-4" />
                                            </button>
                                        </div>
                                    </div>

                                    <p className="text-sm font-semibold text-gray-900 dark:text-white whitespace-nowrap">
                                        {formatPrice(item.product.price * item.quantity, item.product.currency)}
                                    </p>
                                </li>
                            ))}
                        </ul>
                    )}
                </div>

                {items.length > 0 && (
                    <div className="border-t border-gray-200 dark:border-gray-800 px-5 py-4 space-y-4">
                        <button
                            type="button"
                            onClick={() => setShowAddress(prev => !prev)}
                            className="flex items-center gap-2 text-sm font-medium text-indigo-600 dark:text-indigo-400 hover:text-indigo-700"
                        >
                            <MapPinIcon className="w-4 h-4" />
                            {showAddress ? 'Ocultar direccion de envio' : 'Agregar direccion de envio (opcional)'}
                        </button>

                        {showAddress && (
                            <div className="grid grid-cols-2 gap-2">
                                <input
                                    placeholder="Nombre"
                                    value={address.first_name}
                                    onChange={e => setAddress(prev => ({ ...prev, first_name: e.target.value }))}
                                    className="col-span-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                                />
                                <input
                                    placeholder="Telefono"
                                    value={address.phone}
                                    onChange={e => setAddress(prev => ({ ...prev, phone: e.target.value }))}
                                    className="col-span-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                                />
                                <input
                                    placeholder="Direccion"
                                    value={address.street}
                                    onChange={e => setAddress(prev => ({ ...prev, street: e.target.value }))}
                                    className="col-span-2 px-3 py-2 text-sm border border-gray-300 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                                />
                                <input
                                    placeholder="Ciudad"
                                    value={address.city}
                                    onChange={e => setAddress(prev => ({ ...prev, city: e.target.value }))}
                                    className="col-span-2 px-3 py-2 text-sm border border-gray-300 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                                />
                            </div>
                        )}

                        <textarea
                            value={notes}
                            onChange={e => setNotes(e.target.value)}
                            placeholder="Notas del pedido (opcional)"
                            rows={2}
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                        />

                        <div className="flex items-center justify-between text-base font-bold text-gray-900 dark:text-white">
                            <span>Total</span>
                            <span>{formatPrice(totalPrice, currency)}</span>
                        </div>

                        <button
                            onClick={handleCheckout}
                            disabled={submitting}
                            className="w-full py-3 bg-indigo-600 text-white font-semibold rounded-xl hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors shadow-sm"
                        >
                            {submitting ? 'Enviando pedido...' : 'Confirmar pedido'}
                        </button>
                    </div>
                )}
            </aside>
        </>
    );
}
