'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { ArrowRight, Download, Loader2, Search, X } from 'lucide-react';
import {
    fetchFindingItems,
    findingItemsCsvUrl,
    type FindingItem,
} from '../../infra/repository/sync-findings';
import { channelBrand } from '../../domain/types';
import { PanelPager } from './PanelPager';
import { ACCENT, ACCENT_BORDER, ACCENT_SOFT, CARD_BORDER, inputCls } from '../panel-theme';

interface FindingItemsTableProps {
    code: string;
    detail?: string;
    businessId: number | null;
    total: number;
    channels?: string[];
    codeByChannel?: Record<string, string>;
}

const SEARCH_DEBOUNCE_MS = 400;

const COMPARATIVOS = new Set(['sku_spacing', 'sku_typo']);

const CRUZADOS: Record<string, { esta: string; estaNota?: string; falta: string; faltaNota?: string }> = {
    not_published: {
        esta: 'Si esta en tu inventario',
        estaNota: 'ERP, no es canal de venta',
        falta: 'Falta publicarlo en',
        faltaNota: 'canales de venta',
    },
    sold_not_owned: { esta: 'Se vende en', falta: 'No existe en' },
    channel_imbalance: { esta: 'Si esta publicado en', falta: 'Falta publicarlo en' },
};

const PATRON: Record<string, string> = {
    spacing: 'Sobra un espacio',
    transposed: 'Letras trocadas',
    one_char: 'Un carácter distinto',
};

function unidades(value?: number) {
    if (value === undefined || value === null) return null;
    return `${value} und`;
}

function LadoSKU({ sku, name, qty, malo }: { sku: string; name?: string; qty?: number; malo: boolean }) {
    const tono = malo
        ? 'border-red-200 bg-red-50 text-red-700 dark:border-red-500/40 dark:bg-red-900/20 dark:text-red-300'
        : 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-500/40 dark:bg-emerald-900/20 dark:text-emerald-300';
    const cantidad = unidades(qty);
    return (
        <div className={`rounded-lg border px-2 py-1.5 ${tono}`}>
            <div className="flex items-baseline justify-between gap-2">
                <span className="font-mono text-[12px] font-bold">{sku || '—'}</span>
                {cantidad && <span className="flex-shrink-0 text-[10px] font-semibold opacity-70">{cantidad}</span>}
            </div>
            {name && <span className="mt-0.5 line-clamp-1 text-[10.5px] opacity-80">{name}</span>}
        </div>
    );
}

function ChannelChip({ name, code }: { name: string; code?: string }) {
    const brand = channelBrand(code);
    return (
        <span className={`inline-flex items-center gap-1 whitespace-nowrap rounded-full border px-1.5 py-0.5 text-[10px] font-semibold ${brand.chip}`}>
            <span className="h-1.5 w-1.5 flex-shrink-0 rounded-full" style={{ backgroundColor: brand.dot }} />
            {name}
        </span>
    );
}

function FilaComparativa({ item, canal }: { item: FindingItem; canal: string }) {
    const corregirEnCanal = item.fix_side !== 'probability';
    const donde = corregirEnCanal ? canal : 'Probability';
    const correcto = corregirEnCanal ? item.counterpart_sku : item.sku;

    return (
        <tr className="border-t align-top transition-colors hover:bg-gray-50/70 dark:hover:bg-gray-800/40" style={{ borderColor: CARD_BORDER }}>
            <td className="w-[28%] px-3 py-2">
                <LadoSKU sku={item.sku} name={item.name} qty={item.channel_qty} malo={corregirEnCanal} />
            </td>
            <td className="w-6 px-0 py-2 text-center">
                <ArrowRight size={12} className="mx-auto text-gray-300 dark:text-gray-600" />
            </td>
            <td className="w-[28%] px-3 py-2">
                <LadoSKU
                    sku={item.counterpart_sku ?? ''}
                    name={item.counterpart_name}
                    qty={item.own_qty}
                    malo={!corregirEnCanal}
                />
            </td>
            <td className="px-3 py-2">
                <span className="text-[11.5px] font-semibold text-gray-800 dark:text-gray-100">
                    Corrige en {donde}
                </span>
                <div className="mt-0.5 flex flex-wrap items-center gap-1">
                    <span className="rounded bg-gray-100 px-1.5 py-0.5 font-mono text-[10.5px] font-bold text-gray-700 dark:bg-gray-700 dark:text-gray-200">
                        {correcto}
                    </span>
                    {item.pattern && (
                        <span className="text-[10.5px] text-gray-500 dark:text-gray-400">
                            {PATRON[item.pattern] ?? item.pattern}
                        </span>
                    )}
                </div>
            </td>
        </tr>
    );
}

function ListaCanales({ canales, codeByChannel }: { canales?: string[]; codeByChannel?: Record<string, string> }) {
    if (!canales || canales.length === 0) {
        return <span className="text-[11px] text-gray-300 dark:text-gray-600">—</span>;
    }
    return (
        <div className="flex flex-wrap gap-1">
            {canales.map(canal => (
                <ChannelChip key={canal} name={canal} code={codeByChannel?.[canal]} />
            ))}
        </div>
    );
}

function FilaCruzada({ item, codeByChannel }: { item: FindingItem; codeByChannel?: Record<string, string> }) {
    return (
        <tr className="border-t align-top transition-colors hover:bg-gray-50/70 dark:hover:bg-gray-800/40" style={{ borderColor: CARD_BORDER }}>
            <td className="px-3 py-2 font-mono text-[12px] font-bold" style={{ color: ACCENT }}>{item.sku}</td>
            <td className="max-w-[20rem] px-3 py-2 text-gray-600 dark:text-gray-300">
                <span className="line-clamp-2">{item.name || '—'}</span>
            </td>
            <td className="px-3 py-2">
                <ListaCanales canales={item.present_in} codeByChannel={codeByChannel} />
            </td>
            <td className="px-3 py-2">
                <ListaCanales canales={item.missing_in} codeByChannel={codeByChannel} />
            </td>
        </tr>
    );
}

function FilaSimple({ item, codeByChannel }: { item: FindingItem; codeByChannel?: Record<string, string> }) {
    return (
        <tr className="border-t align-top transition-colors hover:bg-gray-50/70 dark:hover:bg-gray-800/40" style={{ borderColor: CARD_BORDER }}>
            <td className="px-3 py-2 font-mono text-[12px] font-bold" style={{ color: ACCENT }}>{item.sku}</td>
            <td className="max-w-[18rem] px-3 py-2 text-gray-600 dark:text-gray-300">
                <span className="line-clamp-2">{item.name || '—'}</span>
            </td>
            <td className="max-w-[24rem] px-3 py-2 text-gray-600 dark:text-gray-300">
                <span className="line-clamp-2">{item.detail || '—'}</span>
            </td>
            <td className="px-3 py-2">
                <div className="flex flex-wrap gap-1">
                    {(item.channels ?? []).map(channel => (
                        <ChannelChip key={channel} name={channel} code={codeByChannel?.[channel]} />
                    ))}
                </div>
            </td>
        </tr>
    );
}

export function FindingItemsTable({ code, detail, businessId, total, channels, codeByChannel }: FindingItemsTableProps) {
    const [items, setItems] = useState<FindingItem[]>([]);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState('');
    const [term, setTerm] = useState('');
    const inFlight = useRef<AbortController | null>(null);

    useEffect(() => {
        const id = setTimeout(() => setTerm(search.trim()), SEARCH_DEBOUNCE_MS);
        return () => clearTimeout(id);
    }, [search]);

    const load = useCallback(async (nextPage: number) => {
        inFlight.current?.abort();
        const controller = new AbortController();
        inFlight.current = controller;
        setLoading(true);
        try {
            const result = await fetchFindingItems(code, businessId ?? undefined, nextPage, term, controller.signal);
            if (controller.signal.aborted) return;
            setItems(result.items);
            setPage(result.page);
            setTotalPages(result.total_pages);
        } finally {
            if (!controller.signal.aborted) setLoading(false);
        }
    }, [code, businessId, term]);

    useEffect(() => {
        void load(1);
    }, [load]);

    useEffect(() => () => inFlight.current?.abort(), []);

    const comparativo = COMPARATIVOS.has(code);
    const cruzado = CRUZADOS[code];
    const canal = channels && channels.length === 1 ? channels[0] : 'el canal';

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-2">
            {detail && <p className="text-[12px] text-gray-500 dark:text-gray-400">{detail}</p>}
            <div className="flex items-center gap-2">
                <div className="relative min-w-0 flex-1">
                    <Search size={13} className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input
                        value={search}
                        onChange={event => setSearch(event.target.value)}
                        placeholder="Buscar SKU o producto"
                        className={`${inputCls} py-1.5 pl-7 pr-7`}
                        style={{ borderColor: CARD_BORDER }}
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
                    download={`${code}-${new Date().toISOString().slice(0, 10)}.csv`}
                    className="inline-flex flex-shrink-0 items-center gap-1.5 rounded-lg border px-3 py-1.5 text-[12px] font-semibold transition-colors"
                    style={{ borderColor: ACCENT_BORDER, backgroundColor: ACCENT_SOFT, color: ACCENT }}
                >
                    <Download size={13} />
                    Excel
                </a>
            </div>

            <div className="min-h-0 flex-1 overflow-y-auto rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                <table className="w-full border-collapse text-[11.5px]">
                    <thead className="sticky top-0 z-10 bg-gray-50 dark:bg-gray-800">
                        <tr className="text-left text-[10px] uppercase tracking-wider text-gray-400">
                            {comparativo ? (
                                <>
                                    <th className="px-3 py-2 font-bold">En {canal}</th>
                                    <th className="px-0 py-1" />
                                    <th className="px-2 py-1 font-bold">En Probability</th>
                                    <th className="px-2 py-1 font-bold">Que hacer</th>
                                </>
                            ) : cruzado ? (
                                <>
                                    <th className="px-3 py-2 font-bold">SKU</th>
                                    <th className="px-3 py-2 font-bold">Producto</th>
                                    <th className="px-3 py-2 font-bold text-emerald-600 dark:text-emerald-400">
                                        {cruzado.esta}
                                        {cruzado.estaNota && (
                                            <span className="ml-1 font-normal normal-case italic text-gray-400">({cruzado.estaNota})</span>
                                        )}
                                    </th>
                                    <th className="px-3 py-2 font-bold text-red-500 dark:text-red-400">
                                        {cruzado.falta}
                                        {cruzado.faltaNota && (
                                            <span className="ml-1 font-normal normal-case italic text-gray-400">({cruzado.faltaNota})</span>
                                        )}
                                    </th>
                                </>
                            ) : (
                                <>
                                    <th className="px-3 py-2 font-bold">SKU</th>
                                    <th className="px-3 py-2 font-bold">Producto</th>
                                    <th className="px-3 py-2 font-bold">Detalle</th>
                                    <th className="px-3 py-2 font-bold">Canales</th>
                                </>
                            )}
                        </tr>
                    </thead>
                    <tbody>
                        {items.map((item, index) =>
                            comparativo ? (
                                <FilaComparativa key={`${item.sku}-${index}`} item={item} canal={canal} />
                            ) : cruzado ? (
                                <FilaCruzada key={`${item.sku}-${index}`} item={item} codeByChannel={codeByChannel} />
                            ) : (
                                <FilaSimple key={`${item.sku}-${index}`} item={item} codeByChannel={codeByChannel} />
                            ),
                        )}
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
                        {term ? `Sin resultados para "${term}"` : 'Sin productos en esta categoría'}
                    </p>
                )}

            </div>

            <PanelPager
                page={page}
                totalPages={totalPages}
                total={total}
                shown={items.length}
                noun="productos"
                pageSize={50}
                onPage={p => void load(p)}
            />
        </div>
    );
}
