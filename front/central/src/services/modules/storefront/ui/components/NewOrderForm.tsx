'use client';

import { useState, useCallback, useEffect } from 'react';
import { StorefrontProduct, CreateStorefrontOrderDTO } from '../../domain/types';
import { getCatalogAction, createOrderAction } from '../../infra/actions';
import { CatalogGrid } from './CatalogGrid';
import { TrashIcon, MinusIcon, PlusIcon } from '@heroicons/react/24/outline';

interface CartItem {
    product: StorefrontProduct;
    quantity: number;
}

interface NewOrderFormProps {
    businessId?: number;
}

export function NewOrderForm({ businessId }: NewOrderFormProps) {
    const [products, setProducts] = useState<StorefrontProduct[]>([]);
    const [cart, setCart] = useState<CartItem[]>([]);
    const [notes, setNotes] = useState('');
    const [address, setAddress] = useState({
        first_name: '',
        last_name: '',
        phone: '',
        street: '',
        city: '',
        state: '',
    });
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [success, setSuccess] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [search, setSearch] = useState('');

    const loadProducts = useCallback(async (searchQuery?: string) => {
        setLoading(true);
        try {
            const result = await getCatalogAction({
                page: 1,
                page_size: 50,
                search: searchQuery,
                business_id: businessId,
            });
            setProducts(result.data);
        } catch {
            setError('Error cargando productos');
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => {
        loadProducts();
    }, [loadProducts]);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        loadProducts(search);
    };

    const addToCart = useCallback((product: StorefrontProduct) => {
        setCart(prev => {
            const existing = prev.find(item => item.product.id === product.id);
            if (existing) {
                return prev.map(item =>
                    item.product.id === product.id
                        ? { ...item, quantity: item.quantity + 1 }
                        : item
                );
            }
            return [...prev, { product, quantity: 1 }];
        });
    }, []);

    const updateQuantity = (productId: string, delta: number) => {
        setCart(prev =>
            prev
                .map(item =>
                    item.product.id === productId
                        ? { ...item, quantity: Math.max(0, item.quantity + delta) }
                        : item
                )
                .filter(item => item.quantity > 0)
        );
    };

    const removeFromCart = (productId: string) => {
        setCart(prev => prev.filter(item => item.product.id !== productId));
    };

    const total = cart.reduce((sum, item) => sum + item.product.price * item.quantity, 0);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (cart.length === 0) return;

        setSubmitting(true);
        setError(null);

        const orderData: CreateStorefrontOrderDTO = {
            items: cart.map(item => ({
                product_id: item.product.id,
                quantity: item.quantity,
            })),
            notes: notes || undefined,
            address: address.street ? {
                first_name: address.first_name,
                last_name: address.last_name,
                phone: address.phone,
                street: address.street,
                city: address.city,
                state: address.state,
            } : undefined,
        };

        try {
            const result = await createOrderAction(orderData, businessId);
            if ('success' in result && !result.success) {
                setError(result.message || 'Error al crear pedido');
            } else {
                setSuccess(true);
                setCart([]);
            }
        } catch (err: any) {
            setError(err.message || 'Error al crear pedido');
        } finally {
            setSubmitting(false);
        }
    };

    const formatPrice = (price: number) => {
        return new Intl.NumberFormat('es-CO', {
            style: 'currency',
            currency: 'COP',
            minimumFractionDigits: 0,
        }).format(price);
    };

    if (success) {
        return (
            <div className="text-center py-12">
                <div className="text-green-500 text-6xl mb-4">&#10003;</div>
                <h2 className="text-2xl font-bold text-gray-900 dark:text-white dark:text-white mb-2">Pedido enviado</h2>
                <p className="text-gray-500 dark:text-gray-400 dark:text-gray-400 mb-6">Tu pedido esta siendo procesado. Recibiras una notificacion cuando este listo.</p>
                <button
                    onClick={() => { setSuccess(false); setCart([]); }}
                    className="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
                >
                    Crear otro pedido
                </button>
            </div>
        );
    }

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Products */}
            <div className="lg:col-span-2">
                <form onSubmit={handleSearch} className="mb-4 flex gap-2">
                    <input
                        type="text"
                        value={search}
                        onChange={e => setSearch(e.target.value)}
                        placeholder="Buscar productos..."
                        className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white dark:text-white"
                    />
                    <button type="submit" className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700">
                        Buscar
                    </button>
                </form>

                {loading ? (
                    <div className="text-center py-8 text-gray-500 dark:text-gray-400">Cargando productos...</div>
                ) : (
                    <CatalogGrid products={products} onAddToCart={addToCart} />
                )}
            </div>

            {/* Cart Summary */}
            <div className="lg:col-span-1">
                <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 sticky top-4">
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white dark:text-white mb-4">Tu Pedido</h3>

                    {cart.length === 0 ? (
                        <p className="text-gray-500 dark:text-gray-400 dark:text-gray-400 text-sm">Agrega productos a tu pedido</p>
                    ) : (
                        <form onSubmit={handleSubmit}>
                            <div className="space-y-3 mb-4 max-h-64 overflow-y-auto">
                                {cart.map(item => (
                                    <div key={item.product.id} className="flex items-center justify-between gap-2 text-sm">
                                        <div className="flex-1 min-w-0">
                                            <p className="font-medium text-gray-900 dark:text-white dark:text-white truncate">{item.product.name}</p>
                                            <p className="text-gray-500 dark:text-gray-400">{formatPrice(item.product.price)}</p>
                                        </div>
                                        <div className="flex items-center gap-1">
                                            <button type="button" onClick={() => updateQuantity(item.product.id, -1)} className="p-1 text-gray-400 hover:text-gray-600 dark:text-gray-300">
                                                <MinusIcon className="w-4 h-4" />
                                            </button>
                                            <span className="w-6 text-center text-gray-900 dark:text-white dark:text-white">{item.quantity}</span>
                                            <button type="button" onClick={() => updateQuantity(item.product.id, 1)} className="p-1 text-gray-400 hover:text-gray-600 dark:text-gray-300">
                                                <PlusIcon className="w-4 h-4" />
                                            </button>
                                            <button type="button" onClick={() => removeFromCart(item.product.id)} className="p-1 text-red-400 hover:text-red-600">
                                                <TrashIcon className="w-4 h-4" />
                                            </button>
                                        </div>
                                    </div>
                                ))}
                            </div>

                            <div className="border-t border-gray-200 dark:border-gray-600 pt-3 mb-4">
                                <div className="flex justify-between text-lg font-bold text-gray-900 dark:text-white dark:text-white">
                                    <span>Total</span>
                                    <span>{formatPrice(total)}</span>
                                </div>
                            </div>

                            {/* Address */}
                            <div className="space-y-2 mb-4">
                                <h4 className="text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-300">Direccion de envio (opcional)</h4>
                                <input
                                    type="text"
                                    placeholder="Nombre"
                                    value={address.first_name}
                                    onChange={e => setAddress(prev => ({ ...prev, first_name: e.target.value }))}
                                    className="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white dark:text-white"
                                />
                                <input
                                    type="text"
                                    placeholder="Direccion"
                                    value={address.street}
                                    onChange={e => setAddress(prev => ({ ...prev, street: e.target.value }))}
                                    className="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white dark:text-white"
                                />
                                <div className="flex gap-2">
                                    <input
                                        type="text"
                                        placeholder="Ciudad"
                                        value={address.city}
                                        onChange={e => setAddress(prev => ({ ...prev, city: e.target.value }))}
                                        className="flex-1 px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white dark:text-white"
                                    />
                                    <input
                                        type="text"
                                        placeholder="Telefono"
                                        value={address.phone}
                                        onChange={e => setAddress(prev => ({ ...prev, phone: e.target.value }))}
                                        className="flex-1 px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white dark:text-white"
                                    />
                                </div>
                            </div>

                            {/* Notes */}
                            <textarea
                                placeholder="Notas del pedido (opcional)"
                                value={notes}
                                onChange={e => setNotes(e.target.value)}
                                rows={2}
                                className="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded mb-4 bg-white dark:bg-gray-700 text-gray-900 dark:text-white dark:text-white"
                            />

                            {error && (
                                <p className="text-red-500 text-sm mb-3">{error}</p>
                            )}

                            <button
                                type="submit"
                                disabled={submitting || cart.length === 0}
                                className="w-full py-3 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
                            >
                                {submitting ? 'Enviando...' : 'Enviar Pedido'}
                            </button>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
}
