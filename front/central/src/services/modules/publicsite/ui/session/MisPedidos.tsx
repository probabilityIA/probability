'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { getTiendaOrdersAction, getTiendaSessionAction } from '../../infra/actions';
import { TiendaOrderSummary } from '../../domain/types';

const formatCOP = (n: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(n);

export function MisPedidos({ slug }: { slug: string }) {
    const [orders, setOrders] = useState<TiendaOrderSummary[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [state, setState] = useState<'loading' | 'anonymous' | 'ready'>('loading');

    useEffect(() => {
        let cancelled = false;
        (async () => {
            const session = await getTiendaSessionAction(slug);
            if (cancelled) return;
            if (!session) {
                setState('anonymous');
                return;
            }
            const result = await getTiendaOrdersAction(slug, page);
            if (cancelled) return;
            setOrders(result?.data || []);
            setTotal(result?.total || 0);
            setState('ready');
        })();
        return () => { cancelled = true; };
    }, [slug, page]);

    if (state === 'loading') {
        return <p className="text-center text-gray-400 py-16">Cargando tus pedidos...</p>;
    }

    if (state === 'anonymous') {
        return (
            <div className="text-center py-16 space-y-4">
                <p className="text-gray-500">Inicia sesión para ver tus pedidos.</p>
                <Link
                    href={`/tienda/${slug}/cuenta`}
                    className="inline-block px-6 py-3 rounded-lg text-white font-medium"
                    style={{ backgroundColor: 'var(--brand-secondary)' }}
                >
                    Iniciar sesion
                </Link>
            </div>
        );
    }

    if (orders.length === 0) {
        return (
            <div className="text-center py-16 space-y-4">
                <p className="text-gray-500">Aún no tienes pedidos en esta tienda.</p>
                <Link
                    href={`/tienda/${slug}/productos`}
                    className="inline-block px-6 py-3 rounded-lg text-white font-medium"
                    style={{ backgroundColor: 'var(--brand-secondary)' }}
                >
                    Explorar productos
                </Link>
            </div>
        );
    }

    return (
        <div className="max-w-3xl mx-auto space-y-3">
            {orders.map(o => (
                <div key={o.id} className="bg-white rounded-xl border border-gray-100 shadow-sm p-4 flex items-center justify-between gap-4">
                    <div className="min-w-0">
                        <p className="font-semibold text-gray-900 truncate">
                            {o.order_number || o.id}
                        </p>
                        <p className="text-xs text-gray-400">
                            {o.created_at?.slice(0, 10)} · {o.items_count} producto{o.items_count === 1 ? '' : 's'}
                        </p>
                    </div>
                    <div className="text-right shrink-0">
                        <p className="font-bold" style={{ color: 'var(--brand-secondary)' }}>{formatCOP(o.total)}</p>
                        {o.status && (
                            <span className="inline-block text-[11px] font-medium px-2 py-0.5 rounded-full bg-gray-100 text-gray-600 mt-1">
                                {o.status}
                            </span>
                        )}
                    </div>
                </div>
            ))}

            {total > 10 && (
                <div className="flex justify-center gap-2 pt-4">
                    <button
                        type="button"
                        disabled={page === 1}
                        onClick={() => setPage(p => p - 1)}
                        className="px-4 py-2 text-sm border border-gray-200 rounded-lg disabled:opacity-40"
                    >
                        Anterior
                    </button>
                    <button
                        type="button"
                        disabled={page * 10 >= total}
                        onClick={() => setPage(p => p + 1)}
                        className="px-4 py-2 text-sm border border-gray-200 rounded-lg disabled:opacity-40"
                    >
                        Siguiente
                    </button>
                </div>
            )}
        </div>
    );
}
