'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { AlertTriangle, Ban, Check, Filter, Loader2, PackageX, RefreshCw, Search } from 'lucide-react';
import type { Integration } from '@/services/integrations/core/domain/types';
import { compareChannelOrdersAction, applyChannelOrdersAction } from '../../infra/actions/orders-compare';
import { ORDERS_COMPARE_TYPE_IDS, type OrdersComparePage, type OrderCompareRow, type OrdersApplyResult } from '../../domain/orders-compare-types';
import { channelBrand } from '../../domain/types';
import { getSyncProvider } from '../providers';
import { ACCENT, ACCENT_BORDER, ACCENT_SOFT, CARD_BORDER, ghostButtonCls, inputCls, primaryButtonCls } from '../panel-theme';
import { PanelPager } from './PanelPager';

interface OrdersCompareTableProps {
    businessId: number | null;
    integrations: Integration[];
}

const PAGE_SIZE = 50;
const dinero = new Intl.NumberFormat('es-CO', { maximumFractionDigits: 0 });

function fecha(iso: string): string {
    if (!iso) return '';
    const parsed = new Date(iso);
    if (Number.isNaN(parsed.getTime()) || parsed.getFullYear() < 2000) return '';
    return parsed.toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: '2-digit' });
}

function hoyMenos(dias: number): string {
    const desde = new Date();
    desde.setDate(desde.getDate() - dias);
    return desde.toISOString().slice(0, 10);
}

function LogoCanal({ integracion }: { integracion: Integration }) {
    const url = integracion.integration_type?.image_url;
    const clave = getSyncProvider(integracion.integration_type_id)?.key ?? '';
    if (url) {
        return (
            <img
                src={url}
                alt=""
                className="h-4 w-4 flex-shrink-0 rounded border bg-white object-contain p-0.5"
                style={{ borderColor: CARD_BORDER }}
            />
        );
    }
    return <span className="h-2 w-2 flex-shrink-0 rounded-full" style={{ backgroundColor: channelBrand(clave).dot }} />;
}

function Kpi({ label, value, tone }: { label: string; value: number; tone?: string }) {
    return (
        <div
            className="rounded-xl border px-3 py-2 dark:bg-gray-800/60"
            style={{ borderColor: CARD_BORDER, backgroundColor: '#fafafd' }}
        >
            <span className="text-[9.5px] font-semibold uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">
                {label}
            </span>
            <p className={`mt-0.5 text-[20px] font-bold leading-none tabular-nums ${tone ?? 'text-gray-900 dark:text-white'}`}>
                {value.toLocaleString('es-CO')}
            </p>
        </div>
    );
}

function EstadoFila({ row }: { row: OrderCompareRow }) {
    if (row.action === 'create') {
        return (
            <span className="inline-flex items-center gap-1 rounded-full border border-sky-300 bg-sky-50 px-2 py-0.5 text-[10px] font-bold text-sky-700 dark:border-sky-500/40 dark:bg-sky-900/25 dark:text-sky-200">
                falta en Probability
            </span>
        );
    }
    if (row.action === 'only_in_probability') {
        return (
            <span className="inline-flex items-center gap-1 rounded-full border border-gray-200 bg-gray-50 px-2 py-0.5 text-[10px] font-bold text-gray-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400">
                solo en Probability
            </span>
        );
    }
    if (row.status_mismatch || row.total_mismatch) {
        return (
            <span className="inline-flex items-center gap-1 rounded-full border border-amber-300 bg-amber-50 px-2 py-0.5 text-[10px] font-bold text-amber-800 dark:border-amber-500/40 dark:bg-amber-900/25 dark:text-amber-200">
                <AlertTriangle size={9} />
                {row.status_mismatch ? 'estado distinto' : 'total distinto'}
            </span>
        );
    }
    return (
        <span className="inline-flex items-center gap-1 rounded-full border border-emerald-300 bg-emerald-50 px-2 py-0.5 text-[10px] font-bold text-emerald-700 dark:border-emerald-500/40 dark:bg-emerald-900/25 dark:text-emerald-200">
            <Check size={9} />
            en las dos
        </span>
    );
}

export function OrdersCompareTable({ businessId, integrations }: OrdersCompareTableProps) {
    const canales = useMemo(
        () => integrations.filter(i => ORDERS_COMPARE_TYPE_IDS.includes(i.integration_type_id)),
        [integrations],
    );

    const [canalID, setCanalID] = useState<number | null>(canales[0]?.id ?? null);
    const [desde, setDesde] = useState(hoyMenos(30));
    const [hasta, setHasta] = useState(() => new Date().toISOString().slice(0, 10));
    const [soloDiferencias, setSoloDiferencias] = useState(true);
    const [busqueda, setBusqueda] = useState('');
    const [pagina, setPagina] = useState(1);

    const [datos, setDatos] = useState<OrdersComparePage | null>(null);
    const [cargando, setCargando] = useState(false);
    const [error, setError] = useState('');
    const [seleccion, setSeleccion] = useState<Set<string>>(new Set());
    const [creando, setCreando] = useState(false);
    const [resultado, setResultado] = useState<OrdersApplyResult | null>(null);

    useEffect(() => {
        if (canalID === null && canales.length > 0) setCanalID(canales[0].id);
    }, [canales, canalID]);

    const comparar = useCallback(async (page: number) => {
        if (!canalID) return;
        setCargando(true);
        setError('');
        try {
            const respuesta = await compareChannelOrdersAction({
                integrationId: canalID,
                businessId: businessId ?? undefined,
                from: desde || undefined,
                to: hasta || undefined,
                page,
                pageSize: PAGE_SIZE,
                onlyDiff: soloDiferencias,
                search: busqueda.trim() || undefined,
            });
            setDatos(respuesta);
            setPagina(respuesta.page);
        } catch (err) {
            setDatos(null);
            setError(err instanceof Error ? err.message : 'No se pudo comparar');
        } finally {
            setCargando(false);
        }
    }, [canalID, businessId, desde, hasta, soloDiferencias, busqueda]);

    useEffect(() => {
        setSeleccion(new Set());
        setResultado(null);
        setDatos(null);
    }, [canalID]);

    const filas = datos?.rows ?? [];
    const creables = filas.filter(row => row.action === 'create');
    const sinInventario = creables.filter(row => !row.moves_inventory).length;

    const alternar = (externalID: string) => {
        setSeleccion(previa => {
            const siguiente = new Set(previa);
            if (siguiente.has(externalID)) siguiente.delete(externalID);
            else siguiente.add(externalID);
            return siguiente;
        });
    };

    const alternarTodas = () => {
        setSeleccion(previa =>
            previa.size === creables.length
                ? new Set()
                : new Set(creables.map(row => row.external_id)),
        );
    };

    const crear = async () => {
        if (!canalID || seleccion.size === 0) return;
        setCreando(true);
        setError('');
        try {
            const respuesta = await applyChannelOrdersAction(canalID, [...seleccion], businessId ?? undefined);
            setResultado(respuesta);
            setSeleccion(new Set());
            await new Promise(resolver => setTimeout(resolver, 2500));
            await comparar(pagina);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'No se pudieron crear las órdenes');
        } finally {
            setCreando(false);
        }
    };

    if (canales.length === 0) {
        return (
            <p className="py-20 text-center text-[12px] italic text-gray-400 dark:text-gray-500">
                Ninguno de tus canales conectados permite comparar ordenes todavia
            </p>
        );
    }

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto pb-2">
            <div className="flex flex-wrap items-end gap-2">
                <label className="flex flex-col gap-1">
                    <span className="text-[10px] font-semibold uppercase tracking-[0.14em] text-gray-400">Canal</span>
                    <select
                        value={canalID ?? ''}
                        onChange={event => setCanalID(Number(event.target.value))}
                        className={`${inputCls} min-w-[190px]`}
                        style={{ borderColor: CARD_BORDER }}
                    >
                        {canales.map(canal => (
                            <option key={canal.id} value={canal.id}>{canal.name}</option>
                        ))}
                    </select>
                </label>

                <label className="flex flex-col gap-1">
                    <span className="text-[10px] font-semibold uppercase tracking-[0.14em] text-gray-400">Desde</span>
                    <input
                        type="date"
                        value={desde}
                        onChange={event => setDesde(event.target.value)}
                        className={inputCls}
                        style={{ borderColor: CARD_BORDER }}
                    />
                </label>

                <label className="flex flex-col gap-1">
                    <span className="text-[10px] font-semibold uppercase tracking-[0.14em] text-gray-400">Hasta</span>
                    <input
                        type="date"
                        value={hasta}
                        onChange={event => setHasta(event.target.value)}
                        className={inputCls}
                        style={{ borderColor: CARD_BORDER }}
                    />
                </label>

                <label className="flex flex-1 flex-col gap-1">
                    <span className="text-[10px] font-semibold uppercase tracking-[0.14em] text-gray-400">Buscar</span>
                    <span className="relative">
                        <Search size={13} className="absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" />
                        <input
                            value={busqueda}
                            onChange={event => setBusqueda(event.target.value)}
                            onKeyDown={event => { if (event.key === 'Enter') comparar(1); }}
                            placeholder="número de orden o cliente"
                            className={`${inputCls} pl-8`}
                            style={{ borderColor: CARD_BORDER }}
                        />
                    </span>
                </label>

                <button
                    onClick={() => setSoloDiferencias(valor => !valor)}
                    className={ghostButtonCls}
                    style={{
                        borderColor: soloDiferencias ? ACCENT_BORDER : CARD_BORDER,
                        backgroundColor: soloDiferencias ? ACCENT_SOFT : '#ffffff',
                        color: soloDiferencias ? ACCENT : '#4b5563',
                    }}
                    title="Ocultar las órdenes que ya estan iguales en los dos lados"
                >
                    <Filter size={13} />
                    Solo diferencias
                </button>

                <button
                    onClick={() => comparar(1)}
                    disabled={cargando || !canalID}
                    className={primaryButtonCls}
                    style={{ backgroundColor: ACCENT }}
                >
                    {cargando ? <Loader2 size={13} className="animate-spin" /> : <RefreshCw size={13} />}
                    Comparar
                </button>
            </div>

            {error && (
                <p className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-[11.5px] font-semibold text-red-700 dark:border-red-500/40 dark:bg-red-900/25 dark:text-red-200">
                    {error}
                </p>
            )}

            {resultado && (
                <div className="rounded-lg border border-emerald-200 bg-emerald-50 px-3 py-2 text-[11.5px] text-emerald-800 dark:border-emerald-500/40 dark:bg-emerald-900/25 dark:text-emerald-100">
                    <p className="font-semibold">
                        {resultado.queued.length} orden{resultado.queued.length === 1 ? '' : 'es'} enviada{resultado.queued.length === 1 ? '' : 's'} a crear.
                        {resultado.skipped.length > 0 && ` ${resultado.skipped.length} ya existian.`}
                    </p>
                    {resultado.note && <p className="mt-0.5">{resultado.note}</p>}
                    {resultado.failed && Object.keys(resultado.failed).length > 0 && (
                        <p className="mt-0.5 text-red-700 dark:text-red-300">
                            Fallaron: {Object.entries(resultado.failed).map(([id, motivo]) => `${id} (${motivo})`).join(', ')}
                        </p>
                    )}
                    <p className="mt-0.5 text-[10.5px] italic">
                        La creacion es asincrona: vuelve a comparar en unos segundos para verlas como &quot;en las dos&quot;.
                    </p>
                </div>
            )}

            {datos && (
                <div className="grid gap-2 sm:grid-cols-3 lg:grid-cols-5">
                    <Kpi label="En el canal" value={datos.totals.in_sync + datos.totals.to_create} />
                    <Kpi label="Faltan aquí" value={datos.totals.to_create} tone="text-sky-600 dark:text-sky-300" />
                    <Kpi label="En las dos" value={datos.totals.in_sync} tone="text-emerald-600 dark:text-emerald-300" />
                    <Kpi label="Estado distinto" value={datos.totals.with_status_mismatch} tone="text-amber-600 dark:text-amber-300" />
                    <Kpi label="Sin mover stock" value={datos.totals.without_inventory} tone="text-gray-500 dark:text-gray-400" />
                </div>
            )}

            {sinInventario > 0 && (
                <p className="flex items-start gap-1.5 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-[11.5px] text-amber-900 dark:border-amber-500/40 dark:bg-amber-900/25 dark:text-amber-100">
                    <PackageX size={14} className="mt-0.5 flex-shrink-0" />
                    <span>
                        {sinInventario} de las ordenes por crear entran como historicas: el canal ya las dio por
                        entregadas, despachadas, canceladas o devueltas, asi que se crean completas pero
                        <strong> no mueven inventario</strong>. El resto sigue el flujo normal y reserva stock.
                    </span>
                </p>
            )}

            {cargando && (
                <div className="flex items-center justify-center gap-2 py-16 text-[12px] text-gray-500">
                    <Loader2 size={14} className="animate-spin" />
                    Preguntandole las ordenes al canal
                </div>
            )}

            {!cargando && datos && filas.length === 0 && (
                <p className="py-16 text-center text-[12px] italic text-gray-400 dark:text-gray-500">
                    No hay diferencias en ese rango de fechas
                </p>
            )}

            {!cargando && filas.length > 0 && (
                <>
                    <div className="overflow-x-auto rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                        <table className="w-full min-w-[860px] border-collapse text-[11.5px]">
                            <thead>
                                <tr className="bg-gray-50 text-left text-[10px] uppercase tracking-[0.12em] text-gray-500 dark:bg-gray-800 dark:text-gray-400">
                                    <th className="w-9 px-2 py-2">
                                        <input
                                            type="checkbox"
                                            checked={creables.length > 0 && seleccion.size === creables.length}
                                            onChange={alternarTodas}
                                            disabled={creables.length === 0}
                                            title="Seleccionar todas las que faltan"
                                        />
                                    </th>
                                    <th className="px-2 py-2">Orden</th>
                                    <th className="px-2 py-2">Cliente</th>
                                    <th className="px-2 py-2">Fecha</th>
                                    <th className="px-2 py-2">Estado canal</th>
                                    <th className="px-2 py-2">Estado aquí</th>
                                    <th className="px-2 py-2 text-right">Total</th>
                                    <th className="px-2 py-2">Situación</th>
                                    <th className="px-2 py-2">Inventario</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filas.map(row => {
                                    const creable = row.action === 'create';
                                    return (
                                        <tr
                                            key={`${row.action}-${row.external_id}`}
                                            className="border-t hover:bg-gray-50 dark:hover:bg-gray-800/60"
                                            style={{ borderColor: CARD_BORDER }}
                                        >
                                            <td className="px-2 py-2">
                                                <input
                                                    type="checkbox"
                                                    checked={seleccion.has(row.external_id)}
                                                    onChange={() => alternar(row.external_id)}
                                                    disabled={!creable}
                                                />
                                            </td>
                                            <td className="px-2 py-2 font-semibold text-gray-900 dark:text-white">
                                                {row.number || row.order_number || row.external_id}
                                                <span className="block text-[10px] font-normal text-gray-400">{row.external_id}</span>
                                            </td>
                                            <td className="px-2 py-2 text-gray-700 dark:text-gray-200">{row.customer_name || '-'}</td>
                                            <td className="px-2 py-2 whitespace-nowrap text-gray-500">{fecha(row.created_at)}</td>
                                            <td className="px-2 py-2 text-gray-700 dark:text-gray-200">{row.channel_status || '-'}</td>
                                            <td className="px-2 py-2 text-gray-700 dark:text-gray-200">{row.local_status || '-'}</td>
                                            <td className="px-2 py-2 text-right tabular-nums text-gray-900 dark:text-white">
                                                {dinero.format(row.total || row.local_total || 0)}
                                            </td>
                                            <td className="px-2 py-2"><EstadoFila row={row} /></td>
                                            <td className="px-2 py-2">
                                                {creable && !row.moves_inventory ? (
                                                    <span
                                                        title={row.inventory_note}
                                                        className="inline-flex items-center gap-1 rounded-full border border-gray-300 bg-gray-50 px-2 py-0.5 text-[10px] font-bold text-gray-600 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300"
                                                    >
                                                        <Ban size={9} />
                                                        no mueve stock
                                                    </span>
                                                ) : creable ? (
                                                    <span className="text-[10.5px] text-gray-500">reserva stock</span>
                                                ) : (
                                                    <span className="text-[10.5px] text-gray-300 dark:text-gray-600">-</span>
                                                )}
                                            </td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>

                    <div className="flex flex-wrap items-center justify-between gap-2">
                        <span className="flex items-center gap-2 text-[11.5px] text-gray-500">
                            {canalID && canales.find(c => c.id === canalID) && (
                                <LogoCanal integracion={canales.find(c => c.id === canalID)!} />
                            )}
                            {seleccion.size > 0
                                ? `${selección.size} seleccionada${selección.size === 1 ? '' : 's'} para crear en Probability`
                                : 'Marca las órdenes que faltan para crearlas aquí'}
                        </span>
                        <button
                            onClick={crear}
                            disabled={creando || seleccion.size === 0}
                            className={primaryButtonCls}
                            style={{ backgroundColor: ACCENT }}
                        >
                            {creando ? <Loader2 size={13} className="animate-spin" /> : <Check size={13} />}
                            Crear en Probability
                        </button>
                    </div>

                    {datos && datos.total_pages > 1 && (
                        <PanelPager
                            page={datos.page}
                            totalPages={datos.total_pages}
                            total={datos.total}
                            shown={filas.length}
                            noun="ordenes"
                            pageSize={datos.page_size}
                            onPage={page => comparar(page)}
                        />
                    )}
                </>
            )}
        </div>
    );
}
