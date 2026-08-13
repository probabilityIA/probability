'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Check, Download, Loader2, Search, TriangleAlert, X } from 'lucide-react';
import {
    fetchMatchMatrix,
    matchMatrixCsvUrl,
    type MatrixCell,
    type MatrixColumn,
    type MatrixRow,
} from '../../infra/repository/sync-findings';
import { channelBrand } from '../../domain/types';
import { PanelPager } from './PanelPager';
import { ACCENT, ACCENT_BORDER, ACCENT_SOFT, CARD_BORDER, inputCls } from '../panel-theme';

interface MatchMatrixTableProps {
    businessId: number | null;
}

const SEARCH_DEBOUNCE_MS = 400;

function Campo({ label, value, tone }: { label: string; value?: string; tone?: string }) {
    const vacio = !value;
    return (
        <div className="flex items-baseline gap-1.5 font-mono text-[10.5px] leading-tight">
            <span className="w-6 flex-shrink-0 text-gray-400">{label}</span>
            <span
                className={`truncate ${vacio ? 'italic text-gray-300 dark:text-gray-600' : tone ?? 'text-gray-500 dark:text-gray-400'}`}
                title={value}
            >
                {vacio ? 'null' : value}
            </span>
        </div>
    );
}

function Celda({ cell }: { cell: MatrixCell }) {
    if (!cell.present) {
        return <td className="px-3 py-2 align-top text-[12px] text-gray-300 dark:text-gray-600">—</td>;
    }

    const distinto = !cell.sku_matches;
    const tono = distinto ? 'text-amber-700 dark:text-amber-300' : 'text-emerald-700 dark:text-emerald-400';
    const marca = distinto
        ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300'
        : 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400';
    const id = cell.variant_id ? `${cell.external_id} / ${cell.variant_id}` : cell.external_id;

    return (
        <td className="px-3 py-2 align-top">
            <div className="flex items-baseline gap-1.5 font-mono text-[11.5px] leading-tight">
                <span className="w-6 flex-shrink-0 text-gray-400">sku</span>
                <span className={`flex min-w-0 items-center gap-1 font-bold ${tono}`}>
                    <span className={`flex h-3.5 w-3.5 flex-shrink-0 items-center justify-center rounded-full ${marca}`}>
                        {distinto ? <TriangleAlert size={8} /> : <Check size={8} strokeWidth={3} />}
                    </span>
                    <span className="truncate">{cell.sku}</span>
                </span>
            </div>
            <div className="mt-0.5 space-y-0.5">
                <Campo label="ean" value={cell.barcode} />
                <Campo label="id" value={id} />
            </div>
            {!cell.sku_known && <div className="mt-0.5 text-[10px] italic text-gray-400">emparejado por SKU</div>}
        </td>
    );
}

function Fila({ row, columns }: { row: MatrixRow; columns: MatrixColumn[] }) {
    return (
        <tr className="border-t transition-colors hover:bg-gray-50/70 dark:hover:bg-gray-800/40" style={{ borderColor: CARD_BORDER }}>
            <td className="sticky left-0 z-10 bg-white px-3 py-2 align-top dark:bg-gray-900">
                <span className="mb-1 block max-w-[20rem] truncate text-[11.5px] font-semibold text-gray-700 dark:text-gray-200" title={row.name}>
                    {row.name || 'Sin nombre'}
                </span>
                <div className="flex items-baseline gap-1.5 font-mono text-[11.5px] leading-tight">
                    <span className="w-6 flex-shrink-0 text-gray-400">sku</span>
                    <span className="truncate font-bold" style={{ color: ACCENT }}>{row.sku || 'null'}</span>
                </div>
                <div className="mt-0.5 space-y-0.5">
                    <Campo label="ean" value={row.barcode} />
                    <Campo label="id" value={row.product_id} />
                </div>
            </td>
            {columns.map(col => {
                const cell = row.cells.find(c => c.integration_id === col.integration_id);
                return cell
                    ? <Celda key={col.integration_id} cell={cell} />
                    : <td key={col.integration_id} className="px-3" />;
            })}
        </tr>
    );
}

function IntegrationToggle({ column, active, onToggle }: { column: MatrixColumn; active: boolean; onToggle: () => void }) {
    const brand = channelBrand(column.code);
    return (
        <button
            onClick={onToggle}
            aria-pressed={active}
            className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-[11.5px] font-semibold transition-all ${
                active ? brand.chip : 'border-gray-200 bg-white text-gray-500 hover:bg-gray-50 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300'
            }`}
        >
            <span
                className="h-2 w-2 flex-shrink-0 rounded-full transition-opacity"
                style={{ backgroundColor: brand.dot, opacity: active ? 1 : 0.35 }}
            />
            {column.name}
            {active && <Check size={11} className="flex-shrink-0" strokeWidth={3} />}
        </button>
    );
}

export function MatchMatrixTable({ businessId }: MatchMatrixTableProps) {
    const [columns, setColumns] = useState<MatrixColumn[]>([]);
    const [rows, setRows] = useState<MatrixRow[]>([]);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState('');
    const [term, setTerm] = useState('');
    const [integrationIds, setIntegrationIds] = useState<number[]>([]);
    const inFlight = useRef<AbortController | null>(null);

    useEffect(() => {
        const id = setTimeout(() => setTerm(search.trim()), SEARCH_DEBOUNCE_MS);
        return () => clearTimeout(id);
    }, [search]);

    const filters = useMemo(() => ({ search: term, integrationIds }), [term, integrationIds]);

    const load = useCallback(async (nextPage: number) => {
        inFlight.current?.abort();
        const controller = new AbortController();
        inFlight.current = controller;
        setLoading(true);
        try {
            const result = await fetchMatchMatrix(businessId ?? undefined, nextPage, filters, controller.signal);
            if (controller.signal.aborted) return;
            setColumns(result.columns);
            setRows(result.rows);
            setPage(result.page);
            setTotalPages(result.total_pages);
            setTotal(result.total);
        } finally {
            if (!controller.signal.aborted) setLoading(false);
        }
    }, [businessId, filters]);

    useEffect(() => {
        void load(1);
    }, [load]);

    useEffect(() => () => inFlight.current?.abort(), []);

    const hayFiltros = term !== '' || integrationIds.length > 0;

    const alternar = (id: number) =>
        setIntegrationIds(prev => (prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id]));

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-2">
            <div className="flex items-center gap-2">
                <div className="relative min-w-0 flex-1">
                    <Search size={13} className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input
                        value={search}
                        onChange={event => setSearch(event.target.value)}
                        placeholder="Buscar SKU, producto o codigo de barras"
                        className={`${inputCls} py-1.5 pl-7 pr-7`}
                        style={{ borderColor: CARD_BORDER }}
                    />
                    {search !== '' && (
                        <button
                            onClick={() => setSearch('')}
                            className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
                        >
                            <X size={12} />
                        </button>
                    )}
                </div>
                <a
                    href={matchMatrixCsvUrl(businessId ?? undefined, filters)}
                    className="inline-flex flex-shrink-0 items-center gap-1.5 rounded-lg border px-3 py-1.5 text-[12px] font-semibold transition-colors"
                    style={{ borderColor: ACCENT_BORDER, backgroundColor: ACCENT_SOFT, color: ACCENT }}
                >
                    <Download size={13} />
                    Excel
                </a>
            </div>

            {columns.length > 0 && (
                <div className="flex flex-wrap items-center gap-1.5">
                    <span className="text-[10px] font-bold uppercase tracking-wider text-gray-400">Coincide en</span>
                    {columns.map(col => (
                        <IntegrationToggle
                            key={col.integration_id}
                            column={col}
                            active={integrationIds.includes(col.integration_id)}
                            onToggle={() => alternar(col.integration_id)}
                        />
                    ))}
                    {integrationIds.length > 0 && (
                        <button
                            onClick={() => setIntegrationIds([])}
                            className="text-[11px] font-semibold text-gray-400 underline-offset-2 transition-colors hover:text-gray-700 hover:underline"
                        >
                            limpiar
                        </button>
                    )}
                </div>
            )}

            <div className="min-h-0 flex-1 overflow-auto rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                <table className="w-full border-collapse">
                    <thead className="sticky top-0 z-20 bg-gray-50 dark:bg-gray-800">
                        <tr className="text-left">
                            <th className="sticky left-0 z-30 min-w-[20rem] bg-gray-50 px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-gray-400 dark:bg-gray-800">
                                Probability
                            </th>
                            {columns.map(col => (
                                <th key={col.integration_id} className="min-w-[13rem] px-3 py-2">
                                    <span className="inline-flex items-center gap-1.5 whitespace-nowrap text-[11.5px] font-semibold text-gray-700 dark:text-gray-200">
                                        <span className="h-2 w-2 flex-shrink-0 rounded-full" style={{ backgroundColor: channelBrand(col.code).dot }} />
                                        {col.name}
                                    </span>
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {rows.map(row => (
                            <Fila key={row.product_id} row={row} columns={columns} />
                        ))}
                    </tbody>
                </table>

                {loading && (
                    <div className="flex items-center justify-center gap-2 py-8 text-[12px] text-gray-500">
                        <Loader2 size={14} className="animate-spin" />
                        Cargando
                    </div>
                )}

                {!loading && rows.length === 0 && (
                    <p className="px-3 py-8 text-center text-[12px] italic text-gray-400">
                        {hayFiltros ? 'Ningun producto coincide en todas las integraciones elegidas' : 'Sin productos'}
                    </p>
                )}
            </div>

            <PanelPager
                page={page}
                totalPages={totalPages}
                total={total}
                shown={rows.length}
                noun="productos"
                onPage={p => void load(p)}
            />
        </div>
    );
}
