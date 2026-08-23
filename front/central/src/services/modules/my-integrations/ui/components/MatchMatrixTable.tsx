'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Check, Download, Eye, Loader2, Search, TriangleAlert, X } from 'lucide-react';
import { ProductDetailModal } from '@/services/modules/products/ui/components/ProductDetailModal';
import {
    fetchMatchMatrix,
    matchMatrixCsvUrl,
    type MatrixCell,
    type MatrixColumn,
    type MatrixRow,
    type MatrixSearchBy,
} from '../../infra/repository/sync-findings';
import { DynamicFilters, type ActiveFilter, type FilterOption } from '@/shared/ui/dynamic-filters';
import { channelBrand } from '../../domain/types';
import { PanelPager } from './PanelPager';
import { ACCENT, ACCENT_BORDER, ACCENT_SOFT, CARD_BORDER, inputCls } from '../panel-theme';

interface MatchMatrixTableProps {
    businessId: number | null;
}

const SEARCH_DEBOUNCE_MS = 400;

function nombreDeArchivo(prefijo: string): string {
    const hoy = new Date();
    const dia = `${hoy.getFullYear()}-${String(hoy.getMonth() + 1).padStart(2, '0')}-${String(hoy.getDate()).padStart(2, '0')}`;
    return `${prefijo}-${dia}.csv`;
}

const PREFIJO_ESTA = 'esta_';
const PREFIJO_FALTA = 'falta_';

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

function Fila({ row, columns, onDetail }: { row: MatrixRow; columns: MatrixColumn[]; onDetail: (productId: string) => void }) {
    return (
        <tr className="border-t transition-colors hover:bg-gray-50/70 dark:hover:bg-gray-800/40" style={{ borderColor: CARD_BORDER }}>
            <td className="sticky left-0 z-10 w-9 bg-white px-1.5 py-2 align-top dark:bg-gray-900">
                <button
                    onClick={() => onDetail(row.product_id)}
                    title="Ver detalle del producto"
                    aria-label="Ver detalle del producto"
                    className="flex h-6 w-6 items-center justify-center rounded-md border transition-colors hover:opacity-80"
                    style={{ borderColor: ACCENT_BORDER, backgroundColor: ACCENT_SOFT, color: ACCENT }}
                >
                    <Eye size={12} />
                </button>
            </td>
            <td className="sticky left-9 z-10 bg-white px-3 py-2 align-top dark:bg-gray-900">
                <div className="flex gap-2">
                    {row.image_url ? (
                        <img
                            src={row.image_url}
                            alt=""
                            loading="lazy"
                            className="h-11 w-11 flex-shrink-0 rounded-md border bg-white object-cover"
                            style={{ borderColor: CARD_BORDER }}
                        />
                    ) : (
                        <span
                            className="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-md border text-[8.5px] italic text-gray-300"
                            style={{ borderColor: CARD_BORDER }}
                        >
                            sin foto
                        </span>
                    )}
                    <div className="min-w-0 flex-1">
                        <span className="mb-1 block max-w-[18rem] truncate text-[11.5px] font-semibold text-gray-700 dark:text-gray-200" title={row.name}>
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
                    </div>
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

const CAMPOS_BUSQUEDA: { key: MatrixSearchBy; label: string; hint: string }[] = [
    { key: 'all', label: 'Todo', hint: 'Buscar SKU, producto o ean' },
    { key: 'sku', label: 'SKU', hint: 'Buscar por SKU' },
    { key: 'name', label: 'Producto', hint: 'Buscar por nombre' },
    { key: 'barcode', label: 'Ean', hint: 'Buscar por codigo de barras' },
];

export function MatchMatrixTable({ businessId }: MatchMatrixTableProps) {
    const [columns, setColumns] = useState<MatrixColumn[]>([]);
    const [rows, setRows] = useState<MatrixRow[]>([]);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState('');
    const [term, setTerm] = useState('');
    const [searchBy, setSearchBy] = useState<MatrixSearchBy>('all');
    const [seleccion, setSeleccion] = useState<Record<string, string | boolean>>({});
    const [pageSize, setPageSize] = useState(10);
    const [detalleID, setDetalleID] = useState<string | null>(null);
    const inFlight = useRef<AbortController | null>(null);

    useEffect(() => {
        const id = setTimeout(() => setTerm(search.trim()), SEARCH_DEBOUNCE_MS);
        return () => clearTimeout(id);
    }, [search]);

    const idsDe = useCallback((prefijo: string) => Object.keys(seleccion)
        .filter(key => key.startsWith(prefijo) && seleccion[key])
        .map(key => Number(key.slice(prefijo.length)))
        .filter(id => id > 0), [seleccion]);

    const presentIn = useMemo(() => idsDe(PREFIJO_ESTA), [idsDe]);
    const missingIn = useMemo(() => idsDe(PREFIJO_FALTA), [idsDe]);

    const filters = useMemo(
        () => ({ search: term, searchBy, presentIn, missingIn }),
        [term, searchBy, presentIn, missingIn],
    );

    const opcionesFiltro: FilterOption[] = useMemo(() => {
        const verbo = (col: MatrixColumn) => (col.is_sales ? 'Se vende en' : 'Esta en');
        const negado = (col: MatrixColumn) => (col.is_sales ? 'No se vende en' : 'No esta en');
        return [
            ...columns.map(col => ({
                key: `${PREFIJO_ESTA}${col.integration_id}`,
                label: `${verbo(col)} ${col.name}`,
                type: 'boolean' as const,
                dot: channelBrand(col.code).dot,
                imageUrl: col.image_url,
            })),
            ...columns.map(col => ({
                key: `${PREFIJO_FALTA}${col.integration_id}`,
                label: `${negado(col)} ${col.name}`,
                type: 'boolean' as const,
                dot: channelBrand(col.code).dot,
                imageUrl: col.image_url,
            })),
        ];
    }, [columns]);

    const filtrosActivos: ActiveFilter[] = useMemo(
        () => Object.entries(seleccion)
            .filter(([, value]) => value !== '' && value !== false && value !== undefined)
            .map(([key, value]) => {
                const definicion = opcionesFiltro.find(opcion => opcion.key === key);
                return {
                    key,
                    label: definicion?.label ?? key,
                    value: value as ActiveFilter['value'],
                    type: definicion?.type ?? 'text',
                };
            }),
        [seleccion, opcionesFiltro],
    );

    const agregarFiltro = useCallback((key: string, value: string | boolean) => {
        setSeleccion(prev => ({ ...prev, [key]: value }));
    }, []);

    const quitarFiltro = useCallback((key: string) => {
        setSeleccion(prev => {
            const siguiente = { ...prev };
            delete siguiente[key];
            return siguiente;
        });
    }, []);

    const load = useCallback(async (nextPage: number) => {
        inFlight.current?.abort();
        const controller = new AbortController();
        inFlight.current = controller;
        setLoading(true);
        try {
            const result = await fetchMatchMatrix(businessId ?? undefined, nextPage, filters, controller.signal, pageSize);
            if (controller.signal.aborted) return;
            setColumns(result.columns);
            setRows(result.rows);
            setPage(result.page);
            setTotalPages(result.total_pages);
            setTotal(result.total);
        } catch (error) {
            if (!(error instanceof DOMException && error.name === 'AbortError')) throw error;
        } finally {
            if (!controller.signal.aborted) setLoading(false);
        }
    }, [businessId, filters, pageSize]);

    useEffect(() => {
        void load(1);
    }, [load]);

    useEffect(() => () => inFlight.current?.abort(), []);

    const hayFiltros = term !== '' || filtrosActivos.length > 0;
    const nombreArchivo = nombreDeArchivo(hayFiltros ? 'resumen-canales-filtrado' : 'resumen-canales');

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-2">
            <p className="text-[12px] text-gray-500 dark:text-gray-400">
                Una fila por producto, una columna por canal.{' '}
                <span className="font-semibold text-emerald-600 dark:text-emerald-400">Verde</span> = el SKU coincide.
                Combina <span className="font-semibold">se vende en</span> y{' '}
                <span className="font-semibold">no se vende en</span> para encontrar lo que falta publicar.
            </p>

            <div className="flex flex-wrap items-center gap-2">
                <div className="flex flex-shrink-0 items-center">
                    <select
                        value={searchBy}
                        onChange={event => setSearchBy(event.target.value as MatrixSearchBy)}
                        title="Por que campo buscar"
                        className={`${inputCls} w-auto rounded-r-none border-r-0 py-1.5 pl-2.5 pr-6 text-[12px]`}
                        style={{ borderColor: CARD_BORDER }}
                    >
                        {CAMPOS_BUSQUEDA.map(campo => (
                            <option key={campo.key} value={campo.key}>{campo.label}</option>
                        ))}
                    </select>
                    <div className="relative w-52">
                        <Search size={13} className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" />
                        <input
                            value={search}
                            onChange={event => setSearch(event.target.value)}
                            placeholder={CAMPOS_BUSQUEDA.find(c => c.key === searchBy)?.hint ?? ''}
                            className={`${inputCls} rounded-l-none py-1.5 pl-7 pr-7`}
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
                </div>

                <div className="min-w-0 flex-1">
                    <DynamicFilters
                        variant="bar"
                        availableFilters={opcionesFiltro}
                        activeFilters={filtrosActivos}
                        onAddFilter={agregarFiltro}
                        onRemoveFilter={quitarFiltro}
                        triggerClassName="!h-8 !w-8 !rounded-lg"
                    />
                </div>

                <a
                    href={matchMatrixCsvUrl(businessId ?? undefined, filters)}
                    download={nombreArchivo}
                    className="inline-flex flex-shrink-0 items-center gap-1.5 rounded-lg border px-3 py-1.5 text-[12px] font-semibold transition-colors"
                    style={{ borderColor: ACCENT_BORDER, backgroundColor: ACCENT_SOFT, color: ACCENT }}
                >
                    <Download size={13} />
                    Excel
                </a>
            </div>

            <div className="min-h-0 flex-1 overflow-auto rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                <table className="w-full border-collapse">
                    <thead className="sticky top-0 z-20 bg-gray-50 dark:bg-gray-800">
                        <tr className="text-left">
                            <th className="sticky left-0 z-30 w-9 bg-gray-50 px-1.5 py-2 dark:bg-gray-800" />
                            <th className="sticky left-9 z-30 min-w-[20rem] bg-gray-50 px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-gray-400 dark:bg-gray-800">
                                Probability
                            </th>
                            {columns.map(col => (
                                <th key={col.integration_id} className="min-w-[13rem] px-3 py-2">
                                    <span className="inline-flex items-center gap-1.5 whitespace-nowrap text-[11.5px] font-semibold text-gray-700 dark:text-gray-200">
                                        {col.image_url ? (
                                            <img
                                                src={col.image_url}
                                                alt=""
                                                className="h-5 w-5 flex-shrink-0 rounded border bg-white object-contain p-0.5"
                                                style={{ borderColor: CARD_BORDER }}
                                            />
                                        ) : (
                                            <span className="h-2 w-2 flex-shrink-0 rounded-full" style={{ backgroundColor: channelBrand(col.code).dot }} />
                                        )}
                                        {col.name}
                                        {!col.is_sales && (
                                            <span className="rounded-full bg-gray-100 px-1.5 py-0.5 text-[9.5px] font-bold uppercase tracking-wide text-gray-500 dark:bg-gray-700 dark:text-gray-300">
                                                ERP
                                            </span>
                                        )}
                                    </span>
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {rows.map(row => (
                            <Fila key={row.product_id} row={row} columns={columns} onDetail={setDetalleID} />
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
                        {hayFiltros ? 'Ningun producto cumple los filtros elegidos' : 'Sin productos'}
                    </p>
                )}
            </div>

            <PanelPager
                page={page}
                totalPages={totalPages}
                total={total}
                shown={rows.length}
                noun="productos"
                pageSize={pageSize}
                onPage={p => void load(p)}
                onPageSize={setPageSize}
            />

            <ProductDetailModal
                productId={detalleID}
                businessId={businessId ?? undefined}
                onClose={() => setDetalleID(null)}
            />
        </div>
    );
}
