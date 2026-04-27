'use client';

import { useState, useEffect, useCallback } from 'react';
import { getMovementsAction } from '../../infra/actions';
import { StockMovement } from '../../domain/types';
import { ChevronRightIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { Spinner } from '@/shared/ui';
import MovementsInlineTable from './MovementsInlineTable';
import { getProductByIdAction } from '@/services/modules/products/infra/actions';

interface Props {
    businessId?: number;
}

interface ProductRow {
    id: string;
    name: string;
    sku: string;
    variant: string;
    variantLabel?: string;
    count: number;
}

const CLIENT_PAGE_SIZE = 15;

function formatVariantAttrs(attrs: any): string {
    if (!attrs || typeof attrs !== 'object') return '';
    return Object.values(attrs).filter(Boolean).join(' / ');
}

export default function MovementsByProductView({ businessId }: Props) {
    const [products, setProducts] = useState<ProductRow[]>([]);
    const [filtered, setFiltered] = useState<ProductRow[]>([]);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');
    const [page, setPage] = useState(1);
    const [selected, setSelected] = useState<ProductRow | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        try {
            const params: any = { page: 1, page_size: 500 };
            if (businessId) params.business_id = businessId;
            const res = await getMovementsAction(params);
            const movements: StockMovement[] = res.data ?? [];

            const map = new Map<string, ProductRow>();
            movements.forEach((m) => {
                const variantLabel = m.variant_label || 'sin-variante';
                const key = `${m.product_id}|${variantLabel}`;
                if (!map.has(key)) {
                    map.set(key, {
                        id: m.product_id,
                        name: m.product_name || m.product_id,
                        sku: m.product_sku || '',
                        variant: formatVariantAttrs((m as any)?.variant_attributes),
                        variantLabel: m.variant_label || undefined,
                        count: 0
                    });
                }
                map.get(key)!.count++;
            });

            const list = Array.from(map.values()).sort((a, b) => b.count - a.count);
            setProducts(list);
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => { load(); }, [load]);

    useEffect(() => {
        const q = search.toLowerCase();
        setFiltered(q ? products.filter((p) => p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q) || p.variant.toLowerCase().includes(q)) : products);
        setPage(1);
    }, [products, search]);

    const totalPages = Math.ceil(filtered.length / CLIENT_PAGE_SIZE);
    const paged = filtered.slice((page - 1) * CLIENT_PAGE_SIZE, page * CLIENT_PAGE_SIZE);

    const handleSearch = (e: React.FormEvent) => { e.preventDefault(); setSearch(searchInput); };

    return (
        <>
            <div className="space-y-3">
                <form onSubmit={handleSearch} className="flex gap-2 max-w-sm">
                    <input
                        type="text"
                        value={searchInput}
                        onChange={(e) => setSearchInput(e.target.value)}
                        placeholder="Buscar producto..."
                        className="flex-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                    <button type="submit" className="px-3 py-2 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors">
                        Buscar
                    </button>
                    {search && (
                        <button type="button" onClick={() => { setSearch(''); setSearchInput(''); }} className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                            Limpiar
                        </button>
                    )}
                </form>

                <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                    <table className="table w-full">
                        <thead>
                            <tr>
                                <th className="text-left">Producto</th>
                                <th className="text-left">SKU</th>
                                <th className="text-left">Variante</th>
                                <th className="text-center">Movimientos</th>
                                <th className="text-center w-12"></th>
                            </tr>
                        </thead>
                        <tbody>
                            {loading ? (
                                <tr><td colSpan={5} className="py-12 text-center"><div className="flex justify-center items-center gap-3"><div className="spinner"></div><span className="text-sm text-gray-500">Cargando...</span></div></td></tr>
                            ) : paged.length === 0 ? (
                                <tr><td colSpan={5} className="py-12 text-center text-sm text-gray-500">{search ? 'Sin resultados.' : 'Sin movimientos registrados.'}</td></tr>
                            ) : (
                                paged.map((p) => (
                                    <tr key={p.id} className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                                        <td className="font-medium text-gray-900 dark:text-white">{p.name}</td>
                                        <td className="text-sm text-gray-500 font-mono">{p.sku}</td>
                                        <td className="text-sm text-gray-500">
                                            {p.variant
                                                ? <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300">{p.variant}</span>
                                                : <span className="text-gray-300">&mdash;</span>
                                            }
                                        </td>
                                        <td className="text-center">
                                            <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-700">{p.count}</span>
                                        </td>
                                        <td className="text-center">
                                            <button onClick={() => setSelected(p)} className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors" title="Ver movimientos">
                                                <ChevronRightIcon className="w-4 h-4" />
                                            </button>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {totalPages > 1 && (
                    <div className="flex items-center justify-between px-1">
                        <span className="text-xs text-gray-500">{filtered.length} productos con movimientos</span>
                        <div className="flex items-center gap-1">
                            <button disabled={page <= 1} onClick={() => setPage((p) => p - 1)} className="px-3 py-1.5 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700">Anterior</button>
                            <span className="text-xs text-gray-600 px-2">{page} / {totalPages}</span>
                            <button disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)} className="px-3 py-1.5 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700">Siguiente</button>
                        </div>
                    </div>
                )}
            </div>

            {selected && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                    <div className="absolute inset-0 bg-black/50" onClick={() => setSelected(null)} />
                    <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-5xl max-h-[90vh] flex flex-col">
                        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                            <div>
                                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{selected.name}</h2>
                                <p className="text-sm text-gray-500 font-mono mt-0.5">
                                    {selected.sku}
                                    {selected.variant && <span className="ml-2 text-gray-400">&bull; {selected.variant}</span>}
                                </p>
                            </div>
                            <button onClick={() => setSelected(null)} className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"><XMarkIcon className="w-5 h-5" /></button>
                        </div>
                        <div className="overflow-auto flex-1 p-5">
                            <MovementsInlineTable productId={selected.id} businessId={businessId} />
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}
