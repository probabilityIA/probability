'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { AlertTriangle, Lock, LockOpen, Search, X } from 'lucide-react';
import type { SyncRunKind } from '../../domain/types';
import { fetchSyncRunItems } from '../../infra/repository/sync-run-items';
import type { ProductApplyActions } from '../providers';
import type { DetailGroup, ProductActionKey, SyncDetailItem } from '../sync-activity-context';
import { ChannelDataStrip } from './ChannelDataStrip';

interface SyncDetailPanelProps {
    integrationId: number;
    kind: SyncRunKind;
    businessId?: number;
    counts: Partial<Record<DetailGroup, number>>;
    liveItems: SyncDetailItem[];
    providerLabel: string;
    apply: ProductApplyActions;
    isProducts: boolean;
    busyAction: ProductActionKey | null;
    onApply: (action: ProductActionKey, skus: string[]) => void;
}

interface GroupStyle {
    label: string;
    dot: string;
    text: string;
    chip: string;
}

const GROUP_STYLES: Record<DetailGroup, GroupStyle> = {
    both: {
        label: 'En ambos',
        dot: 'bg-emerald-500',
        text: 'text-emerald-700 dark:text-emerald-300',
        chip: 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-500/40 dark:bg-emerald-900/30 dark:text-emerald-300',
    },
    not_associated: {
        label: 'Sin asociar',
        dot: 'bg-blue-500',
        text: 'text-blue-700 dark:text-blue-300',
        chip: 'border-blue-200 bg-blue-50 text-blue-700 dark:border-blue-500/40 dark:bg-blue-900/30 dark:text-blue-300',
    },
    only_probability: {
        label: 'Solo en Probability',
        dot: 'bg-orange-500',
        text: 'text-orange-700 dark:text-orange-300',
        chip: 'border-orange-200 bg-orange-50 text-orange-700 dark:border-orange-500/40 dark:bg-orange-900/30 dark:text-orange-300',
    },
    only_channel: {
        label: 'Solo en el canal',
        dot: 'bg-fuchsia-500',
        text: 'text-fuchsia-700 dark:text-fuchsia-300',
        chip: 'border-fuchsia-200 bg-fuchsia-50 text-fuchsia-700 dark:border-fuchsia-500/40 dark:bg-fuchsia-900/30 dark:text-fuchsia-300',
    },
    updated: {
        label: 'Actualizados',
        dot: 'bg-emerald-500',
        text: 'text-emerald-700 dark:text-emerald-300',
        chip: 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-500/40 dark:bg-emerald-900/30 dark:text-emerald-300',
    },
    skipped: {
        label: 'Omitidos',
        dot: 'bg-amber-500',
        text: 'text-amber-700 dark:text-amber-300',
        chip: 'border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-500/40 dark:bg-amber-900/30 dark:text-amber-300',
    },
    sku_typo: {
        label: 'Posible error de digitacion',
        dot: 'bg-amber-500',
        text: 'text-amber-800 dark:text-amber-300',
        chip: 'border-amber-300 bg-amber-100 text-amber-800 dark:border-amber-500/50 dark:bg-amber-900/40 dark:text-amber-300',
    },
    sku_spacing: {
        label: 'SKU con espacios de sobra',
        dot: 'bg-sky-500',
        text: 'text-sky-700 dark:text-sky-300',
        chip: 'border-sky-300 bg-sky-50 text-sky-700 dark:border-sky-500/50 dark:bg-sky-900/30 dark:text-sky-300',
    },
    sku_changed: {
        label: 'SKU cambiado en el canal',
        dot: 'bg-red-600',
        text: 'text-red-700 dark:text-red-400',
        chip: 'border-red-400 bg-red-100 text-red-800 dark:border-red-500/60 dark:bg-red-900/40 dark:text-red-300',
    },
    channel_no_sku: {
        label: 'Sin SKU en el canal',
        dot: 'bg-red-500',
        text: 'text-red-600 dark:text-red-400',
        chip: 'border-red-300 bg-red-50 text-red-700 dark:border-red-500/50 dark:bg-red-900/30 dark:text-red-300',
    },
    failed: {
        label: 'Fallidos',
        dot: 'bg-red-500',
        text: 'text-red-600 dark:text-red-400',
        chip: 'border-red-200 bg-red-50 text-red-600 dark:border-red-500/40 dark:bg-red-900/30 dark:text-red-300',
    },
};

const PRODUCT_ORDER: DetailGroup[] = ['both', 'not_associated', 'only_probability', 'only_channel', 'channel_no_sku', 'sku_changed', 'sku_typo', 'sku_spacing'];
const INVENTORY_ORDER: DetailGroup[] = ['updated', 'skipped', 'failed'];

const PAGE_SIZE = 50;
const SEARCH_DEBOUNCE_MS = 400;
const MAX_CREACION_POR_LOTE = 50;

interface PanelAction {
    key: ProductActionKey;
    label: (channel: string) => string;
    creatable?: DetailGroup[];
}

const ALL_ACTION: PanelAction = {
    key: 'createBothSides',
    label: () => 'Crear en ambos lados',
    creatable: ['only_probability', 'only_channel'],
};

const GROUP_ACTION: Partial<Record<DetailGroup, PanelAction>> = {
    not_associated: { key: 'associate', label: () => 'Asociar al canal' },
    only_probability: { key: 'createInChannel', label: channel => `Crear en ${channel}` },
    only_channel: { key: 'createInProbability', label: () => 'Crear en Probability' },
    both: { key: 'updateInProbability', label: () => 'Actualizar en Probability' },
};

const MATCH_FIELD_LABELS: Record<string, string> = {
    sku: 'SKU',
    barcode: 'codigo de barras',
    external_id: 'ID externo',
    variant_id: 'ID de variante',
    name: 'nombre',
};

export const matchRuleLabel = (rule: string) => {
    const [probability, channel] = rule.split('->');
    const left = MATCH_FIELD_LABELS[probability] || probability;
    const right = MATCH_FIELD_LABELS[channel] || channel;
    return `${left} → ${right}`;
};

const groupLabel = (group: DetailGroup, providerLabel: string) => {
    if (group === 'only_channel') return `Solo en ${providerLabel}`;
    if (group === 'channel_no_sku') return `Sin SKU en ${providerLabel}`;
    if (group === 'sku_changed') return `SKU cambiado en ${providerLabel}`;
    if (group === 'sku_typo') return 'Posible error de digitacion';
    if (group === 'sku_spacing') return 'SKU con espacios de sobra';
    return GROUP_STYLES[group].label;
};

const PARENT_COLORS = [
    { bar: 'bg-sky-400', chip: 'bg-sky-50 text-sky-700 dark:bg-sky-900/40 dark:text-sky-300' },
    { bar: 'bg-violet-400', chip: 'bg-violet-50 text-violet-700 dark:bg-violet-900/40 dark:text-violet-300' },
    { bar: 'bg-amber-400', chip: 'bg-amber-50 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300' },
    { bar: 'bg-teal-400', chip: 'bg-teal-50 text-teal-700 dark:bg-teal-900/40 dark:text-teal-300' },
    { bar: 'bg-rose-400', chip: 'bg-rose-50 text-rose-700 dark:bg-rose-900/40 dark:text-rose-300' },
    { bar: 'bg-lime-400', chip: 'bg-lime-50 text-lime-700 dark:bg-lime-900/40 dark:text-lime-300' },
    { bar: 'bg-fuchsia-400', chip: 'bg-fuchsia-50 text-fuchsia-700 dark:bg-fuchsia-900/40 dark:text-fuchsia-300' },
    { bar: 'bg-cyan-400', chip: 'bg-cyan-50 text-cyan-700 dark:bg-cyan-900/40 dark:text-cyan-300' },
];

const parentColor = (ref: string) => {
    let hash = 0;
    for (let i = 0; i < ref.length; i += 1) hash = (hash * 31 + ref.charCodeAt(i)) % 100000;
    return PARENT_COLORS[hash % PARENT_COLORS.length];
};

export function SyncDetailPanel({
    integrationId,
    kind,
    businessId,
    counts,
    liveItems,
    providerLabel,
    apply,
    isProducts,
    busyAction,
    onApply,
}: SyncDetailPanelProps) {
    const [filter, setFilter] = useState<DetailGroup | 'all'>('all');
    const [selected, setSelected] = useState<Set<string>>(new Set());
    const [creacionDesbloqueada, setCreacionDesbloqueada] = useState(false);
    const [pidiendoDesbloqueo, setPidiendoDesbloqueo] = useState(false);
    const [search, setSearch] = useState('');
    const [term, setTerm] = useState('');
    const [remoteItems, setRemoteItems] = useState<SyncDetailItem[]>([]);
    const [remoteTotal, setRemoteTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(false);
    const [reloading, setReloading] = useState(false);
    const [hasMore, setHasMore] = useState(false);
    const sentinel = useRef<HTMLDivElement | null>(null);

    const isLive = liveItems.length > 0;
    const order = isProducts ? PRODUCT_ORDER : INVENTORY_ORDER;
    const totalItems = order.reduce((sum, group) => sum + (counts[group] ?? 0), 0);

    useEffect(() => {
        const timer = setTimeout(() => setTerm(search.trim()), SEARCH_DEBOUNCE_MS);
        return () => clearTimeout(timer);
    }, [search]);

    const requestSeq = useRef(0);
    const inFlight = useRef<AbortController | null>(null);

    const loadPage = useCallback(async (nextPage: number, replace: boolean) => {
        inFlight.current?.abort();
        const controller = new AbortController();
        inFlight.current = controller;
        const seq = ++requestSeq.current;

        setLoading(true);
        if (replace) setReloading(true);
        try {
            const result = await fetchSyncRunItems({
                integration_id: integrationId,
                kind,
                group: filter === 'all' ? undefined : filter,
                q: term || undefined,
                page: nextPage,
                page_size: PAGE_SIZE,
                business_id: businessId,
            }, controller.signal);

            if (seq !== requestSeq.current) return;

            const incoming = result.items.map(item => ({
                sku: item.sku,
                label: item.label,
                tone: item.tone,
                group: (item.group || 'both') as DetailGroup,
                matchedBy: item.matched_by,
                matchedValue: item.matched_value,
                parentRef: item.parent_ref,
                parentLabel: item.parent_label,
                variantLabel: item.variant_label,
            }));
            setRemoteItems(prev => (replace ? incoming : [...prev, ...incoming]));
            setRemoteTotal(result.total);
            setPage(result.page);
            setHasMore(result.page < result.total_pages);
        } catch {
            if (seq !== requestSeq.current) return;
            setHasMore(false);
        } finally {
            if (seq === requestSeq.current) {
                setLoading(false);
                setReloading(false);
            }
        }
    }, [integrationId, kind, filter, term, businessId]);

    useEffect(() => {
        if (isLive) return;
        setHasMore(false);
        loadPage(1, true);
    }, [isLive, loadPage]);

    useEffect(() => () => inFlight.current?.abort(), []);

    useEffect(() => {
        if (isLive || !hasMore || loading) return;
        const node = sentinel.current;
        if (!node) return;
        const observer = new IntersectionObserver(entries => {
            if (entries[0]?.isIntersecting) loadPage(page + 1, false);
        }, { rootMargin: '80px' });
        observer.observe(node);
        return () => observer.disconnect();
    }, [isLive, hasMore, loading, page, loadPage]);

    const liveCounts = useMemo(() => {
        const map = new Map<DetailGroup, number>();
        for (const item of liveItems) map.set(item.group, (map.get(item.group) ?? 0) + 1);
        return map;
    }, [liveItems]);

    const matches = (item: SyncDetailItem) =>
        term === '' || item.sku.toLowerCase().includes(term.toLowerCase()) || item.label.toLowerCase().includes(term.toLowerCase());

    const liveInGroup = filter === 'all' ? liveItems : liveItems.filter(item => item.group === filter);
    const liveVisible = liveInGroup.filter(matches);
    const visible = isLive ? liveVisible : remoteItems;

    const groupTotal = filter === 'all'
        ? totalItems
        : isLive ? (liveCounts.get(filter) ?? 0) : (counts[filter] ?? 0);
    const shownTotal = isLive ? liveVisible.length : remoteTotal;
    const staleDetail = !isLive && !loading && remoteTotal === 0 && groupTotal > 0;

    const elsewhere = !isLive || term === '' || liveVisible.length > 0
        ? []
        : Array.from(new Set(liveItems.filter(matches).map(item => item.group)));

    const action: PanelAction | null = !isProducts
        ? null
        : filter === 'all'
            ? ALL_ACTION
            : GROUP_ACTION[filter] ?? null;
    const canApply = action !== null && (
        action.key === 'associate'
            ? true
            : action.key === 'createBothSides'
                ? apply.createInChannel !== undefined || apply.createInProbability !== undefined
                : apply[action.key as keyof ProductApplyActions] !== undefined
    );
    const selectable = (item: SyncDetailItem) =>
        canApply && (action?.creatable === undefined || action.creatable.includes(item.group));
    const selectableItems = visible.filter(selectable);

    const applyAllTotal = action?.creatable === undefined
        ? groupTotal
        : action.creatable.reduce((sum, group) => sum + (counts[group] ?? 0), 0);

    const pick = (group: DetailGroup | 'all') => {
        setFilter(group);
        setSelected(new Set());
        setCreacionDesbloqueada(false);
    };

    const toggle = (sku: string) => {
        setSelected(prev => {
            const next = new Set(prev);
            if (next.has(sku)) next.delete(sku);
            else next.add(sku);
            return next;
        });
    };

    const typoCount = isLive ? (liveCounts.get('sku_typo') ?? 0) : (counts.sku_typo ?? 0);
    const skuChangedCount = isLive ? (liveCounts.get('sku_changed') ?? 0) : (counts.sku_changed ?? 0);
    const noSkuCount = isLive ? (liveCounts.get('channel_no_sku') ?? 0) : (counts.channel_no_sku ?? 0);

    const allVisibleSelected = selectableItems.length > 0 && selectableItems.every(item => selected.has(item.sku));

    const creaEnCanal = action?.key === 'createInChannel' || action?.key === 'createBothSides';
    const bloqueado = creaEnCanal && !creacionDesbloqueada;
    const excedeLote = creaEnCanal && selected.size > MAX_CREACION_POR_LOTE;
    const puedeAplicar = selected.size > 0 && !bloqueado && !excedeLote;
    const tope = creaEnCanal ? MAX_CREACION_POR_LOTE : selectableItems.length;
    const aSeleccionar = selectableItems.slice(0, tope).map(item => item.sku);

    return (
        <div className="space-y-1.5">
            <div className="flex flex-wrap items-center gap-1">
                <button
                    onClick={() => pick('all')}
                    className={`rounded-full border px-2 py-0.5 text-[10.5px] font-bold transition-colors ${
                        filter === 'all'
                            ? 'border-gray-400 bg-gray-700 text-white dark:border-gray-500 dark:bg-gray-600'
                            : 'border-gray-200 bg-white text-gray-500 hover:bg-gray-50 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300'
                    }`}
                >
                    Todos {totalItems}
                </button>
                {order.map(group => {
                    const count = isLive ? (liveCounts.get(group) ?? 0) : (counts[group] ?? 0);
                    if (count === 0) return null;
                    const style = GROUP_STYLES[group];
                    return (
                        <button
                            key={group}
                            onClick={() => pick(group)}
                            className={`rounded-full border px-2 py-0.5 text-[10.5px] font-bold transition-colors ${style.chip} ${
                                filter === group ? 'ring-2 ring-offset-1 ring-current dark:ring-offset-gray-800' : 'opacity-80 hover:opacity-100'
                            }`}
                        >
                            {groupLabel(group, providerLabel)} {count}
                        </button>
                    );
                })}
            </div>

            <div className="flex items-center gap-2">
                <div className="relative min-w-0 flex-1">
                    <Search size={12} className="pointer-events-none absolute left-2 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input
                        value={search}
                        onChange={event => setSearch(event.target.value)}
                        placeholder="Buscar SKU o nombre"
                        className="w-full rounded-lg border border-gray-200 bg-white py-1 pl-6 pr-6 text-[11px] text-gray-700 outline-none transition-colors focus:border-blue-400 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-200"
                    />
                    {search !== '' && (
                        <button
                            onClick={() => setSearch('')}
                            className="absolute right-1.5 top-1/2 -translate-y-1/2 text-gray-400 transition-colors hover:text-gray-700 dark:hover:text-gray-200"
                        >
                            <X size={11} />
                        </button>
                    )}
                </div>
                {term !== '' && (
                    <span className="whitespace-nowrap text-[10.5px] font-semibold text-gray-400 dark:text-gray-500">
                        {shownTotal} de {groupTotal}
                    </span>
                )}
            </div>

            {isProducts && (filter === 'all' || filter === 'both') && (
                <ChannelDataStrip businessId={businessId} integrationId={integrationId} />
            )}

            {typoCount > 0 && (filter === 'all' || filter === 'sku_typo') && (
                <div className="flex items-start gap-1.5 rounded-lg border border-amber-300 bg-amber-50 px-2 py-1.5 text-[10.5px] leading-snug text-amber-900 dark:border-amber-500/40 dark:bg-amber-900/25 dark:text-amber-300">
                    <AlertTriangle size={12} className="mt-px flex-shrink-0" />
                    <span>
                        <strong>{typoCount}</strong> {typoCount === 1 ? 'SKU se parece' : 'SKU se parecen'} mucho a un producto que quedo sin emparejar.
                        Suele ser el mismo producto escrito distinto en cada sistema. Revisa y corrige el SKU en {providerLabel} o aca.
                    </span>
                </div>
            )}

            {skuChangedCount > 0 && (filter === 'all' || filter === 'sku_changed') && (
                <div className="flex items-start gap-1.5 rounded-lg border border-red-300 bg-red-100 px-2 py-1.5 text-[10.5px] leading-snug text-red-800 dark:border-red-500/50 dark:bg-red-900/30 dark:text-red-300">
                    <AlertTriangle size={12} className="mt-px flex-shrink-0" />
                    <span>
                        A <strong>{skuChangedCount}</strong> {skuChangedCount === 1 ? 'variante le cambiaron' : 'variantes les cambiaron'} el SKU en {providerLabel}.
                        La asociacion sigue apuntando ahi, asi que les estamos reportando el inventario del producto anterior.
                        Hay que revisarlas y volver a asociarlas.
                    </span>
                </div>
            )}

            {noSkuCount > 0 && (filter === 'all' || filter === 'channel_no_sku') && (
                <div className="flex items-start gap-1.5 rounded-lg border border-red-200 bg-red-50 px-2 py-1.5 text-[10.5px] leading-snug text-red-700 dark:border-red-500/40 dark:bg-red-900/25 dark:text-red-300">
                    <AlertTriangle size={12} className="mt-px flex-shrink-0" />
                    <span>
                        <strong>{noSkuCount}</strong> {noSkuCount === 1 ? 'publicaci\u00f3n o variante' : 'publicaciones o variantes'} en {providerLabel} no {noSkuCount === 1 ? 'tiene' : 'tienen'} SKU.
                        El emparejamiento se hace por SKU, asi que no hay forma de asociarlas ni de sincronizarles stock.
                        Hay que asignarles el SKU en {providerLabel}.
                    </span>
                </div>
            )}

            <div className="max-h-40 overflow-y-auto rounded-lg bg-gray-50 p-1.5 dark:bg-gray-900/40">
                {visible.length === 0 && !loading && (
                    <p className="px-1 py-2 text-[11px] italic text-gray-400 dark:text-gray-500">
                        {term !== ''
                            ? elsewhere.length > 0
                                ? `Sin resultados aqui. "${term}" aparece en: ${elsewhere.map(group => groupLabel(group, providerLabel)).join(', ')}`
                                : `Sin resultados para "${term}"`
                            : staleDetail
                                ? `Los ${groupTotal} productos son de una comparaci\u00f3n anterior. Vuelve a comparar para ver el listado.`
                                : 'Sin productos en este grupo'}
                    </p>
                )}
                {!reloading && visible.map((item, index) => {
                    const style = GROUP_STYLES[item.group];
                    const parent = item.parentRef ? parentColor(item.parentRef) : null;
                    const firstOfParent =
                        parent != null && (index === 0 || visible[index - 1]?.parentRef !== item.parentRef);
                    return (
                        <label
                            key={`${item.group}-${item.sku}-${index}`}
                            className={`flex cursor-pointer items-start gap-2 border-b border-gray-100 py-1 text-[11px] last:border-0 hover:bg-white dark:border-gray-700/60 dark:hover:bg-gray-800/60 ${
                                parent ? 'pl-0 pr-1' : 'px-1'
                            }`}
                        >
                            {parent && <span className={`-my-1 w-0.5 flex-shrink-0 self-stretch rounded-full ${parent.bar}`} />}
                            {selectable(item) && (
                                <input
                                    type="checkbox"
                                    checked={selected.has(item.sku)}
                                    onChange={() => toggle(item.sku)}
                                    className="mt-0.5 h-3 w-3 flex-shrink-0 accent-blue-600"
                                />
                            )}
                            <span className={`mt-1 h-1.5 w-1.5 flex-shrink-0 rounded-full ${style.dot}`} />
                            <span
                                className={`w-28 flex-shrink-0 truncate ${
                                    item.group === 'channel_no_sku'
                                        ? 'italic font-semibold text-red-600 dark:text-red-400'
                                        : item.group === 'sku_typo'
                                            ? 'font-mono font-semibold text-amber-800 dark:text-amber-300'
                                        : item.group === 'sku_changed'
                                            ? 'font-mono font-semibold text-red-700 dark:text-red-400'
                                            : 'font-mono font-semibold text-gray-700 dark:text-gray-200'
                                }`}
                            >
                                {item.sku}
                            </span>
                            {parent ? (
                                <span className="flex min-w-0 flex-1 items-center gap-1">
                                    <span
                                        title={`Variante de la publicaci\u00f3n ${item.parentLabel || item.parentRef} (${item.parentRef})`}
                                        className={`flex-shrink-0 truncate rounded-full px-1.5 py-0.5 font-semibold ${parent.chip} ${
                                            firstOfParent ? 'max-w-[9rem]' : 'max-w-[9rem] opacity-60'
                                        }`}
                                    >
                                        {item.parentLabel || item.parentRef}
                                    </span>
                                    <span className={`min-w-0 flex-1 truncate ${style.text}`}>
                                        {item.variantLabel || item.label}
                                    </span>
                                </span>
                            ) : (
                                <span className={`min-w-0 flex-1 truncate ${style.text}`}>{item.label}</span>
                            )}
                            {item.matchedBy && (
                                <span
                                    title={`Coincidencia por ${matchRuleLabel(item.matchedBy)}${item.matchedValue ? `: ${item.matchedValue}` : ''}`}
                                    className="flex-shrink-0 rounded-full bg-indigo-50 px-1.5 py-0.5 font-semibold text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300"
                                >
                                    {matchRuleLabel(item.matchedBy)}
                                </span>
                            )}
                        </label>
                    );
                })}
                {!isLive && loading && (
                    <div className="space-y-1.5 px-1 py-1">
                        {(reloading ? [0, 1, 2, 3, 4] : [0, 1, 2]).map(row => (
                            <div key={row} className="h-3 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
                        ))}
                    </div>
                )}
                {!isLive && <div ref={sentinel} className="h-px" />}
            </div>

            {action && canApply && (
                <div className="flex flex-wrap items-center gap-2">
                    {selectableItems.length > 0 && (
                    <button
                        onClick={() => setSelected(allVisibleSelected ? new Set() : new Set(aSeleccionar))}
                        className="rounded-lg border border-gray-200 px-2 py-1 text-[11px] font-semibold text-gray-600 transition-colors hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                    >
                        {allVisibleSelected
                            ? 'Quitar selecci\u00f3n'
                            : creaEnCanal
                                ? `Seleccionar ${Math.min(tope, selectableItems.length)}`
                                : `Seleccionar los ${selectableItems.length} cargados`}
                    </button>
                    )}
                    {selected.size > 0 && (
                        <button
                            onClick={() => onApply(action.key, Array.from(selected))}
                            disabled={busyAction !== null || !puedeAplicar}
                            title={bloqueado ? 'Quita el seguro para poder crear en el canal' : undefined}
                            className="flex items-center gap-1.5 rounded-lg bg-[#0d5c80] px-3 py-1 text-[11px] font-bold text-white transition-colors hover:bg-[#0a4964] disabled:cursor-not-allowed disabled:opacity-40"
                        >
                            {busyAction === action.key && (
                                <span className="h-3 w-3 animate-spin rounded-full border border-transparent border-t-current" />
                            )}
                            {bloqueado && <Lock size={11} />}
                            {action.label(providerLabel)}
                            <span className="rounded-full bg-white/25 px-1.5 tabular-nums">{selected.size}</span>
                        </button>
                    )}
                    {selected.size === 0 && !creaEnCanal && term === '' && applyAllTotal > 0 && (
                        <button
                            onClick={() => onApply(action.key, [])}
                            disabled={busyAction !== null}
                            className="flex items-center gap-1.5 rounded-lg bg-[#0d5c80] px-3 py-1 text-[11px] font-bold text-white transition-colors hover:bg-[#0a4964] disabled:cursor-not-allowed disabled:opacity-40"
                        >
                            {busyAction === action.key && (
                                <span className="h-3 w-3 animate-spin rounded-full border border-transparent border-t-current" />
                            )}
                            {action.label(providerLabel)} los {applyAllTotal}
                        </button>
                    )}
                    {excedeLote && (
                        <span className="text-[11px] font-semibold text-amber-600 dark:text-amber-400">
                            Maximo {MAX_CREACION_POR_LOTE} por lote, tienes {selected.size}
                        </span>
                    )}
                    {selected.size === 0 && (
                        <span className="text-[11px] italic text-gray-400 dark:text-gray-500">
                            {creaEnCanal
                                ? `marca los productos que quieras crear, m\u00e1ximo ${MAX_CREACION_POR_LOTE} por lote`
                                : term === ''
                                    ? 'o marca productos para aplicar solo a esos'
                                    : 'marca los productos que quieras aplicar'}
                        </span>
                    )}

                    {creaEnCanal && (
                        <button
                            onClick={() => (creacionDesbloqueada ? setCreacionDesbloqueada(false) : setPidiendoDesbloqueo(true))}
                            aria-pressed={creacionDesbloqueada}
                            className={`ml-auto flex items-center gap-1.5 rounded-full border px-2 py-1 text-[11px] font-semibold transition-colors ${
                                creacionDesbloqueada
                                    ? 'border-amber-300 bg-amber-50 text-amber-700 dark:border-amber-500/50 dark:bg-amber-900/30 dark:text-amber-300'
                                    : 'border-gray-200 bg-white text-gray-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300'
                            }`}
                        >
                            {creacionDesbloqueada ? <LockOpen size={11} /> : <Lock size={11} />}
                            <span
                                className={`relative h-3.5 w-6 rounded-full transition-colors ${
                                    creacionDesbloqueada ? 'bg-amber-500' : 'bg-gray-300 dark:bg-gray-600'
                                }`}
                            >
                                <span
                                    className={`absolute top-0.5 h-2.5 w-2.5 rounded-full bg-white transition-all ${
                                        creacionDesbloqueada ? 'left-3' : 'left-0.5'
                                    }`}
                                />
                            </span>
                            {creacionDesbloqueada ? 'Creaci\u00f3n habilitada' : 'Creaci\u00f3n bloqueada'}
                        </button>
                    )}
                </div>
            )}

            {pidiendoDesbloqueo && createPortal(
                <div className="fixed inset-0 z-[210] flex items-center justify-center p-4">
                    <button
                        aria-label="Cancelar"
                        onClick={() => setPidiendoDesbloqueo(false)}
                        className="absolute inset-0 cursor-default bg-gray-900/50"
                    />
                    <div className="relative w-full max-w-md rounded-2xl bg-white p-5 shadow-2xl dark:bg-gray-900">
                        <div className="flex items-start gap-3">
                            <span className="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-amber-100 text-amber-600 dark:bg-amber-900/40 dark:text-amber-300">
                                <AlertTriangle size={18} />
                            </span>
                            <div className="min-w-0">
                                <h3 className="text-[14px] font-bold text-gray-900 dark:text-white">
                                    Vas a habilitar la creacion en {providerLabel}
                                </h3>
                                <p className="mt-1.5 text-[12px] leading-relaxed text-gray-600 dark:text-gray-300">
                                    Esto publica productos nuevos en {providerLabel}, no en Probability. Lo que se
                                    cree queda visible para tus compradores y hay que borrarlo a mano desde el canal.
                                </p>
                                <p className="mt-2 text-[12px] leading-relaxed text-gray-600 dark:text-gray-300">
                                    Se crean solo los productos que marques, maximo{' '}
                                    <span className="font-bold">{MAX_CREACION_POR_LOTE} por lote</span>.
                                </p>
                            </div>
                        </div>
                        <div className="mt-4 flex justify-end gap-2">
                            <button
                                onClick={() => setPidiendoDesbloqueo(false)}
                                className="rounded-lg border border-gray-200 px-3 py-1.5 text-[12px] font-semibold text-gray-600 transition-colors hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
                            >
                                Cancelar
                            </button>
                            <button
                                onClick={() => { setCreacionDesbloqueada(true); setPidiendoDesbloqueo(false); }}
                                className="rounded-lg bg-amber-500 px-3 py-1.5 text-[12px] font-bold text-white transition-colors hover:bg-amber-600"
                            >
                                Entiendo, habilitar
                            </button>
                        </div>
                    </div>
                </div>,
                document.body,
            )}
        </div>
    );
}
