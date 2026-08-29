'use client';

import { useState, useEffect } from 'react';
import { useCart } from './cart-context';
import { createCheckoutSessionAction, getCheckoutStatusAction, getTiendaSessionAction } from '../../infra/actions';

const formatCOP = (n: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(n);

type PayState = 'idle' | 'starting' | 'waiting' | 'success' | 'failed' | 'timeout';

export function CartWidget({ slug }: { slug: string }) {
    const { items, removeItem, updateQuantity, clear, total, count, isOpen, open, close } = useCart();

    const [step, setStep] = useState<'cart' | 'form'>('cart');
    const [customerName, setCustomerName] = useState('');
    const [customerEmail, setCustomerEmail] = useState('');
    const [customerPhone, setCustomerPhone] = useState('');
    const [customerDni, setCustomerDni] = useState('');
    const [street, setStreet] = useState('');
    const [city, setCity] = useState('');
    const [state, setState] = useState('');
    const [fromAccount, setFromAccount] = useState(false);

    const [payState, setPayState] = useState<PayState>('idle');
    const [error, setError] = useState<string | null>(null);
    const [reference, setReference] = useState<string | null>(null);

    useEffect(() => {
        let cancelled = false;
        const load = () => {
            getTiendaSessionAction(slug).then(session => {
                if (cancelled || !session) return;
                setFromAccount(true);
                setCustomerName(prev => prev || session.name);
                setCustomerEmail(prev => prev || session.email);
                setCustomerPhone(prev => prev || session.phone);
                setCustomerDni(prev => prev || session.dni);
                if (session.address) {
                    setStreet(prev => prev || session.address!.street);
                    setCity(prev => prev || session.address!.city);
                    setState(prev => prev || session.address!.state);
                }
            });
        };
        load();
        window.addEventListener('tienda-session-changed', load);
        return () => {
            cancelled = true;
            window.removeEventListener('tienda-session-changed', load);
        };
    }, [slug]);

    if (count === 0 && !isOpen) {
        return (
            <button
                onClick={open}
                className="fixed bottom-6 right-6 z-40 w-14 h-14 rounded-full shadow-lg flex items-center justify-center text-white opacity-40 cursor-default"
                style={{ backgroundColor: 'var(--brand-secondary)' }}
                aria-hidden
            >
                <CartIcon />
            </button>
        );
    }

    const handlePay = async () => {
        setError(null);
        setPayState('starting');
        try {
            const res = await createCheckoutSessionAction(slug, {
                items: items.map((i) => ({ product_id: i.product_id, quantity: i.quantity })),
                customer_name: customerName,
                customer_email: customerEmail || undefined,
                customer_phone: customerPhone || undefined,
                customer_dni: customerDni || undefined,
                payment_method: 'agree',
                address: street ? { street, city, state, country: 'Colombia' } : undefined,
            });
            if (!res.success || !res.data) {
                throw new Error(res.error || 'No se pudo enviar el pedido');
            }

            const session = res.data;
            setReference(session.reference);

            if (session.payment_method === 'agree') {
                setPayState('success');
                clear();
                return;
            }

            if (!(window as any).BoldCheckout) {
                await new Promise<void>((resolve, reject) => {
                    const script = document.createElement('script');
                    script.src = 'https://checkout.bold.co/library/boldPaymentButton.js';
                    script.async = true;
                    script.onload = () => resolve();
                    script.onerror = () => reject(new Error('No se pudo cargar el script de pago'));
                    document.body.appendChild(script);
                    setTimeout(() => reject(new Error('Timeout cargando el script de pago')), 10000);
                });
            }

            const checkoutConfig: Record<string, unknown> = {
                orderId: session.reference,
                currency: session.currency,
                amount: session.amount,
                apiKey: session.public_key,
                integritySignature: session.hash,
                description: `Compra en tienda - Orden ${session.reference}`,
            };
            if (session.redirection_url) checkoutConfig.redirectionUrl = session.redirection_url;

            // @ts-expect-error BoldCheckout is loaded from an external script
            const checkout = new BoldCheckout(checkoutConfig);
            checkout.open();

            setPayState('waiting');
            pollStatus(session.reference);
        } catch (err: any) {
            setError(err.message || 'Error al iniciar el pago');
            setPayState('failed');
        }
    };

    const pollStatus = (ref: string) => {
        const startedAt = Date.now();
        const interval = setInterval(async () => {
            if (Date.now() - startedAt > 90_000) {
                clearInterval(interval);
                setPayState('timeout');
                return;
            }
            const res = await getCheckoutStatusAction(slug, ref);
            if (res.success && res.data.status === 'paid') {
                clearInterval(interval);
                setPayState('success');
                clear();
            } else if (res.success && res.data.status === 'failed') {
                clearInterval(interval);
                setPayState('failed');
            }
        }, 4000);
    };

    return (
        <>
            <button
                onClick={open}
                className="fixed bottom-6 right-6 z-40 w-14 h-14 rounded-full shadow-lg flex items-center justify-center text-white"
                style={{ backgroundColor: 'var(--brand-secondary)' }}
            >
                <CartIcon />
                {count > 0 && (
                    <span className="absolute -top-1 -right-1 bg-red-500 text-white text-[11px] font-bold rounded-full w-5 h-5 flex items-center justify-center">
                        {count}
                    </span>
                )}
            </button>

            {isOpen && (
                <div className="fixed inset-0 z-50 flex justify-end">
                    <div className="absolute inset-0 bg-black/40" onClick={close} />
                    <div className="relative w-full max-w-md h-full bg-white dark:bg-gray-900 shadow-xl overflow-y-auto">
                        <div className="p-5 border-b border-gray-100 dark:border-gray-800 flex items-center justify-between">
                            <h2 className="text-lg font-bold text-gray-900 dark:text-white">
                                {payState !== 'idle' ? 'Pago' : step === 'cart' ? 'Tu carrito' : 'Datos de entrega'}
                            </h2>
                            <button onClick={close} className="text-gray-400 hover:text-gray-600">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none"><path d="M5 5L19 19M19 5L5 19" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg>
                            </button>
                        </div>

                        <div className="p-5 space-y-4">
                            {payState === 'idle' && step === 'cart' && (
                                <>
                                    {items.length === 0 ? (
                                        <p className="text-sm text-gray-400 text-center py-10">Tu carrito está vacío</p>
                                    ) : (
                                        <div className="space-y-3">
                                            {items.map((it) => (
                                                <div key={it.product_id} className="flex gap-3 items-center border-b border-gray-100 dark:border-gray-800 pb-3">
                                                    <div className="w-14 h-14 rounded-lg bg-gray-100 dark:bg-gray-800 flex-shrink-0 overflow-hidden">
                                                        {it.image_url && <img src={it.image_url} alt={it.name} className="w-full h-full object-cover" />}
                                                    </div>
                                                    <div className="flex-1 min-w-0">
                                                        <p className="text-sm font-medium text-gray-900 dark:text-white truncate">{it.name}</p>
                                                        <p className="text-xs text-gray-400">{formatCOP(it.price)}</p>
                                                        <div className="flex items-center gap-2 mt-1">
                                                            <button onClick={() => updateQuantity(it.product_id, it.quantity - 1)} className="w-6 h-6 rounded border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300">-</button>
                                                            <span className="text-sm w-5 text-center">{it.quantity}</span>
                                                            <button onClick={() => updateQuantity(it.product_id, it.quantity + 1)} className="w-6 h-6 rounded border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300">+</button>
                                                        </div>
                                                    </div>
                                                    <button onClick={() => removeItem(it.product_id)} className="text-gray-300 hover:text-red-500 text-xs">Quitar</button>
                                                </div>
                                            ))}
                                        </div>
                                    )}

                                    {items.length > 0 && (
                                        <>
                                            <div className="flex items-center justify-between pt-2">
                                                <span className="font-semibold text-gray-900 dark:text-white">Total</span>
                                                <span className="font-bold text-lg" style={{ color: 'var(--brand-secondary)' }}>{formatCOP(total)}</span>
                                            </div>
                                            <button
                                                onClick={() => setStep('form')}
                                                className="w-full py-3 rounded-lg text-white font-semibold"
                                                style={{ backgroundColor: 'var(--brand-secondary)' }}
                                            >
                                                Continuar
                                            </button>
                                        </>
                                    )}
                                </>
                            )}

                            {payState === 'idle' && step === 'form' && (
                                <div className="space-y-3">
                                    {fromAccount && (
                                        <p className="text-xs rounded-lg bg-green-50 text-green-700 px-3 py-2">
                                            Usando los datos de tu cuenta. Puedes ajustarlos para este pedido.
                                        </p>
                                    )}
                                    <Field label="Nombre completo *" value={customerName} onChange={setCustomerName} />
                                    <Field label="Correo" value={customerEmail} onChange={setCustomerEmail} type="email" />
                                    <Field label="Teléfono" value={customerPhone} onChange={setCustomerPhone} />
                                    <Field label="Dirección" value={street} onChange={setStreet} />
                                    <div className="grid grid-cols-2 gap-2">
                                        <Field label="Ciudad" value={city} onChange={setCity} />
                                        <Field label="Departamento" value={state} onChange={setState} />
                                    </div>

                                    {error && <p className="text-xs text-red-500">{error}</p>}

                                    <p className="text-xs rounded-lg bg-blue-50 text-blue-700 px-3 py-2">
                                        El pago se coordina directamente con el vendedor: te contactara para acordar el metodo de pago y la entrega.
                                    </p>

                                    <div className="flex items-center justify-between pt-1">
                                        <span className="font-semibold text-gray-900 dark:text-white">Total a pagar</span>
                                        <span className="font-bold text-lg" style={{ color: 'var(--brand-secondary)' }}>{formatCOP(total)}</span>
                                    </div>
                                    <button
                                        onClick={handlePay}
                                        disabled={!customerName}
                                        className="w-full py-3 rounded-lg text-white font-semibold disabled:opacity-50"
                                        style={{ backgroundColor: 'var(--brand-secondary)' }}
                                    >
                                        Enviar pedido
                                    </button>
                                    <button onClick={() => setStep('cart')} className="w-full py-2 text-sm text-gray-500">Volver al carrito</button>
                                </div>
                            )}

                            {payState === 'starting' && <StatusBlock text="Enviando tu pedido..." spinner />}
                            {payState === 'waiting' && (
                                <StatusBlock text="Esperando confirmación del pago. No cierres esta ventana." spinner reference={reference} />
                            )}
                            {payState === 'success' && (
                                <StatusBlock text="Pedido enviado. El vendedor se pondra en contacto contigo para coordinar el pago y la entrega." success reference={reference} />
                            )}
                            {payState === 'failed' && (
                                <StatusBlock text={error || 'El pago no se pudo confirmar.'} failed reference={reference} />
                            )}
                            {payState === 'timeout' && (
                                <StatusBlock text="El pago sigue en proceso. Si ya pagaste, en unos minutos debería confirmarse." reference={reference} />
                            )}
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}

function Field({ label, value, onChange, type = 'text' }: { label: string; value: string; onChange: (v: string) => void; type?: string }) {
    return (
        <div>
            <label className="block text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">{label}</label>
            <input
                type={type}
                value={value}
                onChange={(e) => onChange(e.target.value)}
                className="w-full px-3 py-2 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-offset-0"
            />
        </div>
    );
}

function StatusBlock({ text, spinner, success, failed, reference }: { text: string; spinner?: boolean; success?: boolean; failed?: boolean; reference?: string | null }) {
    return (
        <div className="py-10 text-center space-y-3">
            {spinner && <div className="mx-auto w-8 h-8 border-2 border-gray-200 border-t-gray-600 rounded-full animate-spin" />}
            {success && <div className="mx-auto w-10 h-10 rounded-full bg-green-100 text-green-600 flex items-center justify-center">✓</div>}
            {failed && <div className="mx-auto w-10 h-10 rounded-full bg-red-100 text-red-600 flex items-center justify-center">✕</div>}
            <p className="text-sm text-gray-600 dark:text-gray-300">{text}</p>
            {reference && <p className="text-xs text-gray-400 font-mono">{reference}</p>}
        </div>
    );
}

function CartIcon() {
    return (
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
            <path d="M3 3H4.5L6.68 14.39A2 2 0 008.63 16H17.4A2 2 0 0019.35 14.39L21 6H6" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" />
            <circle cx="9" cy="20" r="1.3" fill="currentColor" />
            <circle cx="17" cy="20" r="1.3" fill="currentColor" />
        </svg>
    );
}
