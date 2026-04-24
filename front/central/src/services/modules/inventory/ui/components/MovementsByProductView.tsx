'use client';

import { useState, useEffect, useCallback } from 'react';
import { getProductsAction } from '@/services/modules/products/infra/actions';
import { Product } from '@/services/modules/products/domain/types';
import { ChevronRightIcon, XMarkIcon } from '@heroicons/react/24/outline';
import MovementsInlineTable from './MovementsInlineTable';

interface Props {
    businessId?: number;
}

const PAGE_SIZE = 15;

export default function MovementsByProductView({ businessId }: Props) {
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');
    const [selected, setSelected] = useState<Product | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        try {
            const params: any = { page, page_size: PAGE_SIZE };
            if (businessId) params.business_id = businessId;
            if (search) params.search = search;
            const res = await getProductsAction(params);
            setProducts((res as any).data ?? []);
            setTotal((res as any).total ?? 0);
            setTotalPages((res as any).total_pages ?? 1);
        } finally {
            setLoading(false);
        }
    }, [page, search, businessId]);

    useEffect(() => { load(); }, [load]);
    useEffect(() => { setPage(1); }, [search, businessId]);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(searchInput);
    };

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
                        <button
                            type="button"
                            onClick={() => { setSearch(''); setSearchInput(''); }}
                            className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                        >
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
                                <th className="text-left">Familia</th>
                                <th className="text-center w-12"></th>
                            </tr>
                        </thead>
                        <tbody>
                            {loading ? (
                                <tr>
                                    <td colSpan={4} className="py-12 text-center">
                                        <div className="flex justify-center items-center gap-3">
                                            <div className="spinner"></div>
                                            <span className="text-sm text-gray-500">Cargando...</span>
                                        </div>
                                    </td>
                                </tr>
                            ) : products.length === 0 ? (
                                <tr>
                                    <td colSpan={4} className="py-12 text-center text-sm text-gray-500">
                                        Sin productos.
                                    </td>
                                </tr>
                            ) : (
                                products.map((p) => (
                                    <tr key={p.id} className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                                        <td className="font-medium text-gray-900 dark:text-white">{p.name}</td>
                                        <td className="text-sm text-gray-500 font-mono">{p.sku}</td>
                                        <td className="text-sm text-gray-500">{(p as any).family?.name || <span className="text-gray-300">&mdash;</span>}</td>
                                        <td className="text-center">
                                            <button
                                                onClick={() => setSelected(p)}
                                                className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                                                title="Ver movimientos"
                                            >
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
                        <span className="text-xs text-gray-500">{total} productos</span>
                        <div className="flex items-center gap-1">
                            <button
                                disabled={page <= 1}
                                onClick={() => setPage((p) => p - 1)}
                                className="px-3 py-1.5 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700"
                            >
                                Anterior
                            </button>
                            <span className="text-xs text-gray-600 px-2">{page} / {totalPages}</span>
                            <button
                                disabled={page >= totalPages}
                                onClick={() => setPage((p) => p + 1)}
                                className="px-3 py-1.5 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700"
                            >
                                Siguiente
                            </button>
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
                                <p className="text-sm text-gray-500 font-mono mt-0.5">{selected.sku}</p>
                            </div>
                            <button
                                onClick={() => setSelected(null)}
                                className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"
                            >
                                <XMarkIcon className="w-5 h-5" />
                            </button>
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
