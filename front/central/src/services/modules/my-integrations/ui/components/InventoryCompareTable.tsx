'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { AlertTriangle, ArrowRight, Filter, History, Loader2, RefreshCw, Search, Send } from 'lucide-react';
import type { Integration } from '@/services/integrations/core/domain/types';
import {
    compararInventario,
    GRUPO_INVENTARIO,
    type ComparePage,
    type CompareRow,
} from '../../infra/repository/inventory-compare';
import { channelBrand } from '../../domain/types';
import { getSyncProvider } from '../providers';
import { useSyncActivity } from '../sync-activity-context';
import { ACCENT, ACCENT_BORDER, ACCENT_SOFT, CARD_BORDER } from '../panel-theme';
import { PanelPager } from './PanelPager';

interface InventoryCompareTableProps {
    businessId: number | null;
    integrations: Integration[];
    fixedIntegrationId?: number;
}

const ETIQUETA_ACCION: Record<string, { texto: string; clase: string }> = {
    update: {
        texto: 'Se actualiza',
        clase: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300',
    },
    unchanged: {
        texto: 'Sin cambio',
        clase: 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400',
    },
    skip: {
        texto: 'No aplica',
        clase: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300',
    },
};

const VENTAS: Record<string, boolean> = { siigo: false };

function hora(iso: string) {
    if (!iso) return '';
    return new Date(iso).toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
}

function fechaGuardada(iso: string) {
    if (!iso) return 'de una comparacion anterior';
    const fecha = new Date(iso);
    if (Number.isNaN(fecha.getTime()) || fecha.getFullYear() < 2000) return 'de una comparacion anterior';
    const minutos = Math.floor((Date.now() - fecha.getTime()) / 60000);
    if (minutos < 1) return 'de hace un momento';
    if (minutos < 60) return `de hace ${minutos} min`;
    const horas = Math.floor(minutos / 60);
    if (horas < 24) return `de hace ${horas} h`;
    const dias = Math.floor(horas / 24);
    return `de hace ${dias} d`;
}

function LogoCanal({ integracion, size = 20 }: { integracion: Integration; size?: number }) {
    const url = integracion.integration_type?.image_url;
    const clave = getSyncProvider(integracion.integration_type_id)?.key ?? '';
    if (url) {
        return (
            <img
                src={url}
                alt=""
                className="flex-shrink-0 rounded border bg-white object-contain p-0.5"
                style={{ borderColor: CARD_BORDER, height: size, width: size }}
            />
        );
    }
    return <span className="h-2 w-2 flex-shrink-0 rounded-full" style={{ backgroundColor: channelBrand(clave).dot }} />;
}

function FotoProducto({ url }: { url?: string }) {
    if (!url) {
        return (
            <span
                className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-md border text-[8.5px] italic text-gray-300 dark:text-gray-600"
                style={{ borderColor: CARD_BORDER }}
            >
                sin foto
            </span>
        );
    }
    return (
        <img
            src={url}
            alt=""
            loading="lazy"
            className="h-10 w-10 flex-shrink-0 rounded-md border bg-white object-cover"
            style={{ borderColor: CARD_BORDER }}
        />
    );
}

function Cantidad({ valor }: { valor: number | null }) {
    if (valor === null) return <span className="text-[11px] italic text-gray-300 dark:text-gray-600">sin dato</span>;
    return <span className="font-semibold tabular-nums">{valor}</span>;
}

export function InventoryCompareTable({ businessId, integrations, fixedIntegrationId }: InventoryCompareTableProps) {
    const { runInventoryOne, running, progress, results } = useSyncActivity();

    const comparables = useMemo(
        () => integrations.filter(i => i.is_active && getSyncProvider(i.integration_type_id)?.compareInventory),
        [integrations],
    );

    const [canalID, setCanalID] = useState<number | null>(fixedIntegrationId ?? comparables[0]?.id ?? null);
    const [pagina, setPagina] = useState(1);
    const [tamano, setTamano] = useState(GRUPO_INVENTARIO);
    const [datos, setDatos] = useState<ComparePage | null>(null);
    const [cargando, setCargando] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [seleccion, setSeleccion] = useState<Set<string>>(new Set());
    const [busqueda, setBusqueda] = useState('');
    const [soloCambios, setSoloCambios] = useState(false);
    const [enviado, setEnviado] = useState(false);

    const canal = comparables.find(i => i.id === canalID) ?? null;
    const tipoCanal = canal?.integration_type_id;
    const claveCanal = canal ? getSyncProvider(canal.integration_type_id)?.key ?? '' : '';
    const esVenta = VENTAS[claveCanal] ?? true;

    const [termino, setTermino] = useState('');

    useEffect(() => {
        const reloj = window.setTimeout(() => setTermino(busqueda.trim()), 350);
        return () => window.clearTimeout(reloj);
    }, [busqueda]);

    const leerGuardado = useCallback(async (page: number) => {
        if (!canalID || !tipoCanal) return;
        setCargando(true);
        setError(null);
        try {
            const respuesta = await compararInventario(
                tipoCanal, canalID, businessId ?? undefined, page, undefined, tamano,
                { source: 'snapshot', only_diff: soloCambios, q: termino },
            );
            setDatos(respuesta);
            setSeleccion(new Set());
        } catch (e) {
            setDatos(null);
            setError(e instanceof Error ? e.message : 'No se pudo leer la ultima comparacion guardada');
        } finally {
            setCargando(false);
        }
    }, [businessId, canalID, tipoCanal, tamano, soloCambios, termino]);

    const comparar = useCallback(async (page: number) => {
        if (!canalID || !tipoCanal) return;
        setCargando(true);
        setError(null);
        try {
            const respuesta = await compararInventario(tipoCanal, canalID, businessId ?? undefined, page, undefined, tamano);
            setDatos(respuesta);
            setSeleccion(new Set());
        } catch (e) {
            setDatos(null);
            setError(e instanceof Error ? e.message : 'No se pudo leer el stock del canal');
        } finally {
            setCargando(false);
        }
    }, [businessId, canalID, tipoCanal, tamano]);

    useEffect(() => {
        void leerGuardado(pagina);
    }, [leerGuardado, pagina]);

    const resultado = canalID ? results[canalID] : undefined;
    const avance = canalID ? progress[canalID] : undefined;

    useEffect(() => {
        if (enviado && resultado?.kind === 'inventory') {
            setEnviado(false);
            void comparar(pagina);
        }
    }, [resultado, enviado, pagina, comparar]);

    const guardado = datos?.from_cache === true;
    const filas = datos?.rows ?? [];
    const buscado = termino.toLowerCase();
    const visibles = guardado ? filas : filas.filter(f => {
        if (soloCambios && f.action !== 'update') return false;
        if (!buscado) return true;
        return f.sku.toLowerCase().includes(buscado) || f.name.toLowerCase().includes(buscado);
    });
    const seleccionables = filas.filter(f => f.action !== 'skip');
    const cambian = filas.filter(f => f.action === 'update');
    const todasMarcadas = seleccionables.length > 0 && seleccionables.every(f => seleccion.has(f.sku));

    const alternar = (sku: string) => {
        setSeleccion(prev => {
            const siguiente = new Set(prev);
            if (siguiente.has(sku)) siguiente.delete(sku);
            else siguiente.add(sku);
            return siguiente;
        });
    };

    const alternarTodas = () => {
        setSeleccion(todasMarcadas ? new Set() : new Set(seleccionables.map(f => f.sku)));
    };

    const sincronizar = (skus: string[]) => {
        if (!canalID || skus.length === 0) return;
        setEnviado(true);
        runInventoryOne(canalID, skus);
    };

    const cambiarPagina = (destino: number) => {
        if (!datos || destino < 1 || destino > datos.total_pages || cargando) return;
        setPagina(destino);
    };

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-2">
            {!fixedIntegrationId && (
                <p className="text-[12px] text-gray-500 dark:text-gray-400">
                    Le preguntamos al canal cuanto stock tiene <span className="font-semibold text-gray-600 dark:text-gray-300">ahora mismo</span>,
                    de a {tamano} productos, y tu eliges cuales se igualan con el stock de Probability.
                </p>
            )}

            <div className="flex flex-wrap items-center gap-2">
                {!fixedIntegrationId && comparables.map(integracion => {
                    const activo = integracion.id === canalID;
                    const suClave = getSyncProvider(integracion.integration_type_id)?.key ?? '';
                    return (
                        <button
                            key={integracion.id}
                            onClick={() => { setCanalID(integracion.id); setPagina(1); }}
                            className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-[11.5px] font-semibold transition-all ${
                                activo ? channelBrand(suClave).chip : 'border-gray-200 text-gray-500 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-400'
                            }`}
                        >
                            <LogoCanal integracion={integracion} size={16} />
                            {integracion.name}
                            {(VENTAS[suClave] ?? true) === false && (
                                <span className="rounded-full bg-gray-100 px-1.5 py-0.5 text-[9px] font-bold uppercase tracking-wide text-gray-500 dark:bg-gray-700 dark:text-gray-300">
                                    ERP
                                </span>
                            )}
                        </button>
                    );
                })}

                <div className="relative ml-auto">
                    <Search size={13} className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input
                        value={busqueda}
                        onChange={event => setBusqueda(event.target.value)}
                        placeholder="Buscar SKU o producto en este grupo"
                        className="w-64 rounded-lg border bg-white py-1.5 pl-7 pr-2 text-[12px] text-gray-700 outline-none placeholder:text-gray-400 dark:bg-gray-800 dark:text-gray-200"
                        style={{ borderColor: CARD_BORDER }}
                    />
                </div>

                <button
                    onClick={() => setSoloCambios(v => !v)}
                    className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-[11.5px] font-semibold transition-colors ${
                        soloCambios
                            ? 'border-emerald-300 bg-emerald-50 text-emerald-700 dark:border-emerald-500/50 dark:bg-emerald-900/30 dark:text-emerald-200'
                            : 'border-gray-200 text-gray-500 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-400'
                    }`}
                >
                    <Filter size={11} />
                    Solo los que no estan iguales
                </button>

                <button
                    onClick={() => comparar(pagina)}
                    disabled={cargando || !canalID}
                    title={`Le pregunta a ${canal?.name ?? 'el canal'} cuanto stock tiene ahora mismo, de a ${tamano} productos`}
                    className="inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1.5 text-[11.5px] font-semibold transition-colors hover:opacity-80 disabled:opacity-50"
                    style={{ borderColor: ACCENT_BORDER, backgroundColor: ACCENT_SOFT, color: ACCENT }}
                >
                    {cargando ? <Loader2 size={12} className="animate-spin" /> : <RefreshCw size={12} />}
                    Comparar ahora
                </button>
            </div>

            {datos && (datos.totals.total > 0 || !guardado) && (
                <div className="flex flex-wrap items-center gap-2 text-[11.5px]">
                    <span className="rounded-full bg-emerald-100 px-2.5 py-0.5 font-bold text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
                        {datos.totals.to_update} se actualizarian
                    </span>
                    <span className="rounded-full bg-gray-100 px-2.5 py-0.5 font-semibold text-gray-500 dark:bg-gray-700 dark:text-gray-400">
                        {datos.totals.unchanged} ya estan iguales
                    </span>
                    <span className="rounded-full bg-amber-100 px-2.5 py-0.5 font-semibold text-amber-700 dark:bg-amber-900/30 dark:text-amber-300">
                        {datos.totals.skipped} no aplican
                    </span>
                    {guardado ? (
                        <span
                            className="inline-flex items-center gap-1.5 rounded-full border border-amber-300 bg-amber-50 px-2.5 py-0.5 font-bold text-amber-800 dark:border-amber-500/40 dark:bg-amber-900/25 dark:text-amber-200"
                            title="Estos numeros salen de la ultima comparacion guardada, no se le pregunto al canal ahora"
                        >
                            <History size={11} />
                            Foto guardada {fechaGuardada(datos.checked_at)} · vuelve a comparar para confirmar
                        </span>
                    ) : (
                        <span className="text-gray-400">Stock del canal leido a las {hora(datos.checked_at)}</span>
                    )}
                </div>
            )}

            <div className="min-h-0 flex-1 overflow-auto rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                <table className="w-full border-collapse">
                    <thead className="sticky top-0 z-20 bg-gray-50 dark:bg-gray-800">
                        <tr className="text-left">
                            <th className="w-9 px-2 py-2">
                                <input
                                    type="checkbox"
                                    checked={todasMarcadas}
                                    onChange={alternarTodas}
                                    disabled={seleccionables.length === 0}
                                    className="h-3.5 w-3.5 accent-blue-600"
                                />
                            </th>
                            <th className="min-w-[22rem] px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-gray-400">
                                Probability
                            </th>
                            <th className="w-28 px-3 py-2 text-center text-[10px] font-bold uppercase tracking-wider text-gray-400">
                                Stock aqui
                            </th>
                            <th className="min-w-[11rem] px-3 py-2 text-center">
                                <span className="inline-flex items-center gap-1.5 whitespace-nowrap text-[11.5px] font-semibold text-gray-700 dark:text-gray-200">
                                    {canal && <LogoCanal integracion={canal} />}
                                    {canal?.name ?? 'Canal'}
                                    {!esVenta && (
                                        <span className="rounded-full bg-gray-100 px-1.5 py-0.5 text-[9.5px] font-bold uppercase tracking-wide text-gray-500 dark:bg-gray-700 dark:text-gray-300">
                                            ERP
                                        </span>
                                    )}
                                </span>
                            </th>
                            <th className="w-32 px-3 py-2 text-center text-[10px] font-bold uppercase tracking-wider text-gray-400">
                                Quedaria
                            </th>
                            <th className="min-w-[13rem] px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-gray-400">
                                Que pasaria
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {visibles.map((fila: CompareRow) => {
                            const etiqueta = ETIQUETA_ACCION[fila.action];
                            const bloqueada = fila.action === 'skip';
                            return (
                                <tr
                                    key={`${fila.product_id}-${fila.external_item_id}-${fila.external_variant_id ?? ''}`}
                                    className="border-t bg-white transition-colors hover:bg-gray-50/80 dark:bg-gray-900 dark:hover:bg-gray-800/60"
                                    style={{ borderColor: CARD_BORDER }}
                                >
                                    <td className="px-2 py-2 align-top">
                                        <input
                                            type="checkbox"
                                            checked={seleccion.has(fila.sku)}
                                            onChange={() => alternar(fila.sku)}
                                            disabled={bloqueada}
                                            className="h-3.5 w-3.5 accent-blue-600 disabled:opacity-30"
                                        />
                                    </td>
                                    <td className="px-3 py-2 align-top">
                                        <div className="flex gap-2">
                                            <FotoProducto url={fila.image_url} />
                                            <div className="min-w-0 flex-1">
                                                <span className="mb-0.5 block max-w-[22rem] truncate text-[11.5px] font-semibold text-gray-700 dark:text-gray-200" title={fila.name}>
                                                    {fila.name || 'Sin nombre'}
                                                </span>
                                                <span className="flex items-baseline gap-1.5 font-mono text-[11.5px] leading-tight">
                                                    <span className="w-6 flex-shrink-0 text-gray-400">sku</span>
                                                    <span className="font-semibold" style={{ color: ACCENT }}>{fila.sku}</span>
                                                </span>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-3 py-2 text-center align-top text-[12px] text-gray-800 dark:text-gray-100">
                                        <Cantidad valor={fila.probability_qty} />
                                    </td>
                                    <td className="px-3 py-2 text-center align-top text-[12px] text-gray-800 dark:text-gray-100">
                                        <Cantidad valor={fila.channel_qty} />
                                    </td>
                                    <td className="px-3 py-2 text-center align-top text-[12px]">
                                        {fila.action === 'update' ? (
                                            <span className="inline-flex items-center gap-1 font-semibold text-gray-700 dark:text-gray-200">
                                                {fila.channel_qty}
                                                <ArrowRight size={11} className="text-gray-400" />
                                                <span className="text-emerald-600 dark:text-emerald-400">{fila.probability_qty}</span>
                                            </span>
                                        ) : (
                                            <span className="text-gray-300 dark:text-gray-600">—</span>
                                        )}
                                    </td>
                                    <td className="px-3 py-2 align-top">
                                        <span className={`inline-flex rounded-full px-2 py-0.5 text-[10.5px] font-bold ${etiqueta.clase}`}>
                                            {etiqueta.texto}
                                        </span>
                                        {fila.reason && (
                                            <span className="mt-0.5 block text-[10.5px] italic text-gray-400">{fila.reason}</span>
                                        )}
                                    </td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>

                {cargando && (
                    <p className="flex items-center justify-center gap-2 py-10 text-[12px] text-gray-500">
                        <Loader2 size={14} className="animate-spin" />
                        {guardado
                            ? 'Leyendo la ultima comparacion guardada'
                            : `Preguntandole a ${canal?.name ?? 'el canal'} cuanto stock tiene ahora mismo`}
                    </p>
                )}

                {!cargando && error && (
                    <p className="flex items-center justify-center gap-2 py-10 text-[12px] text-red-500">
                        <AlertTriangle size={14} />
                        {error}
                    </p>
                )}

                {!cargando && !error && visibles.length === 0 && (
                    <p className="px-3 py-10 text-center text-[12px] italic text-gray-400">
                        {filas.length === 0 && guardado && (datos?.totals.total ?? 0) === 0
                            ? 'Todavia no hay una comparacion guardada de este canal. Dale a "Comparar ahora" para preguntarle su stock.'
                            : filas.length === 0
                                ? 'Este canal no tiene productos asociados para comparar.'
                                : soloCambios
                                    ? 'En este grupo todo esta igual, no hay nada que enviar.'
                                    : 'Ningun producto de este grupo coincide con la busqueda.'}
                    </p>
                )}
            </div>

            <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="flex flex-wrap items-center gap-2">
                    {seleccionables.length > 0 && (
                        <button
                            onClick={alternarTodas}
                            className="rounded-lg border px-2 py-1 text-[11px] font-semibold text-gray-600 transition-colors hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
                            style={{ borderColor: CARD_BORDER }}
                        >
                            {todasMarcadas ? 'Quitar seleccion' : `Seleccionar los ${seleccionables.length} de este grupo`}
                        </button>
                    )}
                    {seleccion.size > 0 && (
                        <span className="text-[11.5px] text-gray-500 dark:text-gray-400">
                            <span className="font-bold text-gray-700 dark:text-gray-200">{seleccion.size}</span> seleccionados
                        </span>
                    )}
                    {seleccion.size === 0 && cambian.length > 0 && (
                        <span className="text-[11px] italic text-gray-400 dark:text-gray-500">
                            marca los que quieras enviar, o envia de una los {cambian.length} que cambian
                        </span>
                    )}
                    {avance && avance.total > 0 && enviado && (
                        <span className="text-[11.5px] tabular-nums text-blue-600 dark:text-blue-400">
                            enviando {avance.processed}/{avance.total}
                        </span>
                    )}
                    {resultado?.kind === 'inventory' && (
                        <span className="rounded-full bg-emerald-50 px-2.5 py-0.5 text-[11.5px] font-semibold text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300">
                            {resultado.updated} actualizados{resultado.failed > 0 ? `, ${resultado.failed} fallidos` : ''}
                        </span>
                    )}
                </div>

                <div className="flex items-center gap-1.5">
                    {seleccion.size > 0 ? (
                        <button
                            onClick={() => sincronizar([...seleccion])}
                            disabled={running || cargando}
                            className="ml-1 inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-[11.5px] font-bold text-white shadow-sm transition-opacity hover:opacity-90 disabled:opacity-40"
                            style={{ backgroundColor: ACCENT }}
                        >
                            {running && enviado ? <Loader2 size={12} className="animate-spin" /> : <Send size={12} />}
                            Enviar a {canal?.name ?? 'el canal'}
                            <span className="rounded-full bg-white/25 px-1.5 tabular-nums">{seleccion.size}</span>
                        </button>
                    ) : cambian.length > 0 && (
                        <button
                            onClick={() => sincronizar(cambian.map(f => f.sku))}
                            disabled={running || cargando}
                            className="ml-1 inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-[11.5px] font-bold text-white shadow-sm transition-opacity hover:opacity-90 disabled:opacity-40"
                            style={{ backgroundColor: ACCENT }}
                        >
                            {running && enviado ? <Loader2 size={12} className="animate-spin" /> : <Send size={12} />}
                            Enviar los {cambian.length} que cambian
                        </button>
                    )}
                </div>
            </div>

            <PanelPager
                page={datos?.page ?? 1}
                totalPages={datos?.total_pages ?? 1}
                total={datos?.total ?? 0}
                shown={filas.length}
                noun="asociados"
                pageSize={tamano}
                onPage={cambiarPagina}
                onPageSize={nuevo => { setTamano(nuevo); setPagina(1); }}
            />

            <p className="text-[11px] text-gray-400">
                Al enviar, Probability vuelve a leer su propio stock en ese instante, por si entro una orden mientras revisabas.
                El stock nunca viaja del canal hacia Probability.
            </p>
        </div>
    );
}
