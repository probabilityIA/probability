'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { Download, Loader2, Search, X } from 'lucide-react';
import {
    fetchFindingItems,
    findingItemsCsvUrl,
    type FindingItem,
} from '../../infra/repository/sync-findings';

interface FindingItemsTableProps {
    code: string;
    businessId: number | null;
    total: number;
}

const SEARCH_DEBOUNCE_MS = 400;

export function FindingItemsTable({ code, businessId, total }: FindingItemsTableProps) {
    const [items, setItems] = useState<FindingItem[]>([]);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState('');
    const [term, setTerm] = useState('');
    const [shown, setShown] = useState(0);
    const inFlight = useRef<AbortController | null>(null);

    useEffect(() => {
        const id = setTimeout(() => setTerm(search.trim()), SEARCH_DEBOUNCE_MS);
        return () => clearTimeout(id);
    }, [search]);

    const load = useCallback(async (nextPage: number, replace: boolean) => {
        inFlight.current?.abort();
        const controller = new AbortController();
        inFlight.current = controller;
        setLoading(true);
        try {
            const result = await fetchFindingItems(code, businessId ?? undefined, nextPage, term, controller.signal);
            if (controller.signal.aborted) return;
            setItems(prev => (replace ? result.items : [...prev, ...result.items]));
            setShown(prev => (replace ? result.items.length : prev + result.items.length));
            setPage(result.page);
            setTotalPages(result.total_pages);
        } finally {
            if (!controller.signal.aborted) setLoading(false);
        }
    }, [code, businessId, term]);

    useEffect(() => {
        void load(1, true);
    }, [load]);

    useEffect(() => () => inFlight.current?.abort(), []);

    const hayMas = page < totalPages;

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-2">
            <div className="flex items-center gap-2">
                <div className="relative min-w-0 flex-1">
                    <Search size={12} className="pointer-events-none absolute left-2 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input
                        value={search}
                        onChange={event => setSearch(event.target.value)}
                        placeholder="Buscar SKU o producto"
                        className="w-full rounded-lg border border-gray-200 bg-white py-1.5 pl-6 pr-6 text-[11.5px] text-gray-700 outline-none transition-colors focus:border-blue-400 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-200"
                    />
                    {search !== '' && (
                        <button
                            onClick={() => setSearch('')}
                            className="absolute right-1.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
                        >
                            <X size={11} />
                        </button>
                    )}
                </div>
                <a
                    href={findingItemsCsvUrl(code, businessId ?? undefined, term)}
                    className="flex flex-shrink-0 items-center gap-1.5 rounded-lg border border-emerald-300 bg-emerald-50 px-2.5 py-1.5 text-[11.5px] font-semibold text-emerald-700 transition-colors hover:bg-emerald-100 dark:border-emerald-500/40 dark:bg-emerald-900/30 dark:text-emerald-300 dark:hover:bg-emerald-900/50"
                >
                    <Download size={12} />
                    Excel
                </a>
            </div>

            <div className="min-h-0 flex-1 overflow-y-auto rounded-lg border border-gray-100 dark:border-gray-700">
                <table className="w-full border-collapse text-[11.5px]">
                    <thead className="sticky top-0 z-10 bg-gray-50 dark:bg-gray-800">
                        <tr className="text-left text-[10.5px] uppercase tracking-wider text-gray-400">
                            <th className="px-2 py-1.5 font-bold">SKU</th>
                            <th className="px-2 py-1.5 font-bold">Producto</th>
                            <th className="px-2 py-1.5 font-bold">Detalle</th>
                            <th className="px-2 py-1.5 font-bold">Canales</th>
                        </tr>
                    </thead>
                    <tbody>
                        {items.map((item, index) => (
                            <tr
                                key={`${item.sku}-${index}`}
                                className="border-t border-gray-100 align-top hover:bg-gray-50 dark:border-gray-700/60 dark:hover:bg-gray-800/50"
                            >
                                <td className="px-2 py-1.5 font-mono font-semibold text-gray-800 dark:text-gray-100">{item.sku}</td>
                                <td className="max-w-[14rem] px-2 py-1.5 text-gray-600 dark:text-gray-300">
                                    <span className="line-clamp-2">{item.name || '—'}</span>
                                </td>
                                <td className="max-w-[18rem] px-2 py-1.5 text-gray-600 dark:text-gray-300">
                                    <span className="line-clamp-2">{item.detail || '—'}</span>
                                </td>
                                <td className="px-2 py-1.5">
                                    <div className="flex flex-wrap gap-1">
                                        {(item.channels ?? []).map(channel => (
                                            <span
                                                key={channel}
                                                className="whitespace-nowrap rounded-full border border-gray-200 bg-white px-1.5 py-0.5 text-[10px] font-semibold text-gray-600 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300"
                                            >
                                                {channel}
                                            </span>
                                        ))}
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>

                {loading && (
                    <div className="flex items-center justify-center gap-2 py-6 text-[11.5px] text-gray-500">
                        <Loader2 size={13} className="animate-spin" />
                        Cargando
                    </div>
                )}

                {!loading && items.length === 0 && (
                    <p className="px-2 py-6 text-center text-[11.5px] italic text-gray-400">
                        {term ? `Sin resultados para "${term}"` : 'Sin productos en esta categoria'}
                    </p>
                )}

                {!loading && hayMas && (
                    <button
                        onClick={() => void load(page + 1, false)}
                        className="w-full border-t border-gray-100 py-2 text-[11.5px] font-semibold text-blue-600 transition-colors hover:bg-blue-50 dark:border-gray-700 dark:text-blue-400 dark:hover:bg-blue-900/20"
                    >
                        Cargar mas
                    </button>
                )}
            </div>

            <p className="text-[10.5px] text-gray-400">
                Mostrando {shown} de {term ? 'los filtrados' : total}. El Excel baja la lista completa.
            </p>
        </div>
    );
}
