'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { AlertTriangle, ArrowRight, Filter, Loader2, RefreshCw, Search, Send } from 'lucide-react';
import type { Integration } from '@/services/integrations/core/domain/types';
import { fetchMatchMatrix, type MatrixRow } from '../../infra/repository/sync-findings';
import { compararInventario, type CompareRow } from '../../infra/repository/inventory-compare';
import { channelBrand } from '../../domain/types';
import { getSyncProvider } from '../providers';
import { useSyncActivity } from '../sync-activity-context';
import { ACCENT, ACCENT_BORDER, ACCENT_SOFT, CARD_BORDER } from '../panel-theme';
import { PanelPager } from './PanelPager';

interface InventoryMatrixTableProps {
    businessId: number | null;
    integrations: Integration[];
}

const GRUPO_INICIAL = 100;
const NO_VENTA: Record<string, boolean> = { siigo: true };

type PorCanal = Record<number, Record<string, CompareRow>>;

function LogoCanal({ integracion }: { integracion: Integration }) {
    const url = integracion.integration_type?.image_url;
    const clave = getSyncProvider(integracion.integration_type_id)?.key ?? '';
    if (url) {
        return (
            <img
                src={url}
                alt=""
                className="h-5 w-5 flex-shrink-0 rounded border bg-white object-contain p-0.5"
                style={{ borderColor: CARD_BORDER }}
            />
        );
    }
    return <span className="h-2 w-2 flex-shrink-0 rounded-full" style={{ backgroundColor: channelBrand(clave).dot }} />;
}

function Celda({ fila }: { fila?: CompareRow }) {
    if (!fila) return <span className="text-[11px] text-gray-300 dark:text-gray-600">—</span>;

    if (fila.action === 'update') {
        return (
            <span className="inline-flex items-center gap-1 text-[12px] font-semibold text-gray-700 dark:text-gray-200">
                {fila.channel_qty}
                <ArrowRight size={10} className="text-gray-400" />
                <span className="text-emerald-600 dark:text-emerald-400">{fila.probability_qty}</span>
            </span>
        );
    }
    if (fila.action === 'unchanged') {
        return <span className="text-[12px] tabular-nums text-gray-400">{fila.channel_qty} igual</span>;
    }
    return (
        <span className="text-[10.5px] italic text-amber-600 dark:text-amber-400" title={fila.reason}>
            {fila.channel_qty === null ? 'sin publicacion' : 'no aplica'}
        </span>
    );
}

export function InventoryMatrixTable({ businessId, integrations }: InventoryMatrixTableProps) {
    const { runInventoryOne, running, results } = useSyncActivity();

    const canales = useMemo(
        () => integrations.filter(i => i.is_active && getSyncProvider(i.integration_type_id)?.compareInventory),
        [integrations],
    );
    const canalesKey = canales.map(c => c.id).join(',');
    const esOrigen = (canal: Integration) => Boolean(NO_VENTA[getSyncProvider(canal.integration_type_id)?.key ?? '']);
    const origenes = canales.filter(esOrigen);
    const ventas = canales.filter(canal => !esOrigen(canal));

    const [pagina, setPagina] = useState(1);
    const [tamano, setTamano] = useState(GRUPO_INICIAL);
    const [productos, setProductos] = useState<MatrixRow[]>([]);
    const [totalPaginas, setTotalPaginas] = useState(1);
    const [total, setTotal] = useState(0);
    const [porCanal, setPorCanal] = useState<PorCanal>({});
    const [cargando, setCargando] = useState(false);
    const [leidoA, setLeidoA] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [seleccion, setSeleccion] = useState<Set<string>>(new Set());
    const [busqueda, setBusqueda] = useState('');
    const [soloCambios, setSoloCambios] = useState(false);
    const [enviando, setEnviando] = useState(false);
    const [barriendo, setBarriendo] = useState(false);
    const [avanceBarrido, setAvanceBarrido] = useState({ grupo: 0, grupos: 0 });
    const cortar = useRef(false);

    const leerGrupo = useCallback(async (page: number) => {
        const matriz = await fetchMatchMatrix(businessId ?? undefined, page, {}, undefined, tamano);
        const skus = matriz.rows.map(r => r.sku).filter(Boolean);
        if (skus.length === 0) return { matriz, mapa: {} as PorCanal };

        const lecturas = await Promise.all(
            canales.map(async canal => {
                try {
                    const respuesta = await compararInventario(
                        canal.integration_type_id,
                        canal.id,
                        businessId ?? undefined,
                        1,
                        skus,
                    );
                    return [canal.id, respuesta.rows] as const;
                } catch {
                    return [canal.id, [] as CompareRow[]] as const;
                }
            }),
        );

        const mapa: PorCanal = {};
        for (const [id, filas] of lecturas) {
            mapa[id] = Object.fromEntries(filas.map(f => [f.sku, f]));
        }
        return { matriz, mapa };
    }, [businessId, canales, tamano]);

    const comparar = useCallback(async (page: number) => {
        setCargando(true);
        setError(null);
        setSeleccion(new Set());
        try {
            const { matriz, mapa } = await leerGrupo(page);
            setProductos(matriz.rows);
            setTotalPaginas(matriz.total_pages || 1);
            setTotal(matriz.total);
            setPorCanal(mapa);
            setLeidoA(new Date().toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' }));
        } catch (e) {
            setError(e instanceof Error ? e.message : 'No se pudo comparar el inventario');
        } finally {
            setCargando(false);
        }
    }, [leerGrupo]);

    const barrerTodo = useCallback(async () => {
        setBarriendo(true);
        setError(null);
        setSeleccion(new Set());
        cortar.current = false;

        const acumulados: MatrixRow[] = [];
        const acumulado: PorCanal = {};
        let grupos = totalPaginas;

        try {
            for (let page = 1; page <= grupos; page++) {
                if (cortar.current) break;
                setAvanceBarrido({ grupo: page, grupos });

                const { matriz, mapa } = await leerGrupo(page);
                grupos = matriz.total_pages || grupos;
                setTotal(matriz.total);
                setTotalPaginas(matriz.total_pages || grupos);

                for (const [id, porSKU] of Object.entries(mapa)) {
                    acumulado[Number(id)] = { ...(acumulado[Number(id)] ?? {}), ...porSKU };
                }
                for (const producto of matriz.rows) {
                    const cambia = Object.values(mapa).some(porSKU => porSKU[producto.sku]?.action === 'update');
                    if (cambia) acumulados.push(producto);
                }

                setPorCanal({ ...acumulado });
                setProductos([...acumulados]);
                setLeidoA(new Date().toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' }));
            }
        } catch (e) {
            setError(e instanceof Error ? e.message : 'No se pudo comparar el inventario');
        } finally {
            setBarriendo(false);
            setAvanceBarrido({ grupo: 0, grupos: 0 });
        }
    }, [leerGrupo, totalPaginas]);

    useEffect(() => {
        if (barriendo || soloCambios) return;
        void comparar(pagina);
    }, [pagina, tamano, canalesKey]);

    useEffect(() => {
        if (enviando && !running) {
            setEnviando(false);
            void comparar(pagina);
        }
    }, [running, enviando, pagina, comparar]);

    const cambiaEnAlgunCanal = (sku: string) => canales.some(c => porCanal[c.id]?.[sku]?.action === 'update');

    const termino = busqueda.trim().toLowerCase();
    const visibles = productos.filter(p => {
        if (soloCambios && !cambiaEnAlgunCanal(p.sku)) return false;
        if (!termino) return true;
        return p.sku?.toLowerCase().includes(termino) || (p.name ?? '').toLowerCase().includes(termino);
    });
    const conCambios = productos.filter(p => cambiaEnAlgunCanal(p.sku));
    const todasMarcadas = conCambios.length > 0 && conCambios.every(p => seleccion.has(p.sku));

    const alternar = (sku: string) => {
        setSeleccion(prev => {
            const siguiente = new Set(prev);
            if (siguiente.has(sku)) siguiente.delete(sku);
            else siguiente.add(sku);
            return siguiente;
        });
    };

    const enviar = (skus: string[]) => {
        if (skus.length === 0) return;
        const elegidos = new Set(skus);
        setEnviando(true);
        for (const canal of canales) {
            const suyos = skus.filter(sku => porCanal[canal.id]?.[sku]?.action === 'update' && elegidos.has(sku));
            if (suyos.length > 0) runInventoryOne(canal.id, suyos);
        }
    };

    const totalPorCanal = (canalID: number) =>
        productos.filter(p => porCanal[canalID]?.[p.sku]?.action === 'update').length;

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-2">
            <p className="text-[12px] text-gray-500 dark:text-gray-400">
                Una fila por producto. Se lee de izquierda a derecha: primero el ERP que te sirve de origen de inventario,
                despues lo que tiene Probability, y despues cada canal de venta. Le preguntamos a cada uno cuanto stock
                tiene <span className="font-semibold text-gray-600 dark:text-gray-300">ahora mismo</span>, de a {tamano} productos.
                La flecha muestra en cuanto quedaria el canal si lo envias.
            </p>

            <div className="flex flex-wrap items-center gap-2">
                <div className="relative">
                    <Search size={13} className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input
                        value={busqueda}
                        onChange={event => setBusqueda(event.target.value)}
                        placeholder="Buscar SKU o producto en este grupo"
                        className="w-72 rounded-lg border bg-white py-1.5 pl-7 pr-2 text-[12px] text-gray-700 outline-none placeholder:text-gray-400 dark:bg-gray-800 dark:text-gray-200"
                        style={{ borderColor: CARD_BORDER }}
                    />
                </div>

                <button
                    onClick={() => {
                        if (barriendo) { cortar.current = true; return; }
                        const activando = !soloCambios;
                        setSoloCambios(activando);
                        if (activando) void barrerTodo();
                        else void comparar(pagina);
                    }}
                    className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-[11.5px] font-semibold transition-colors ${
                        soloCambios
                            ? 'border-emerald-300 bg-emerald-50 text-emerald-700 dark:border-emerald-500/50 dark:bg-emerald-900/30 dark:text-emerald-200'
                            : 'border-gray-200 text-gray-500 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-400'
                    }`}
                >
                    {barriendo ? <Loader2 size={11} className="animate-spin" /> : <Filter size={11} />}
                    {barriendo
                        ? `Revisando grupo ${avanceBarrido.grupo} de ${avanceBarrido.grupos} · parar`
                        : 'Solo los que no estan iguales'}
                </button>

                {[...origenes, ...ventas].map(canal => (
                    <span
                        key={canal.id}
                        className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-[11.5px] font-semibold ${channelBrand(getSyncProvider(canal.integration_type_id)?.key ?? '').chip}`}
                    >
                        <LogoCanal integracion={canal} />
                        {totalPorCanal(canal.id)} por igualar
                    </span>
                ))}

                <button
                    onClick={() => comparar(pagina)}
                    disabled={cargando}
                    className="ml-auto inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1.5 text-[11.5px] font-semibold transition-colors hover:opacity-80 disabled:opacity-50"
                    style={{ borderColor: ACCENT_BORDER, backgroundColor: ACCENT_SOFT, color: ACCENT }}
                >
                    {cargando ? <Loader2 size={12} className="animate-spin" /> : <RefreshCw size={12} />}
                    Volver a comparar
                </button>
                {leidoA && <span className="text-[11px] text-gray-400">Stock leido a las {leidoA}</span>}
            </div>

            <div className="min-h-0 flex-1 overflow-auto rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                <table className="w-full border-collapse">
                    <thead className="sticky top-0 z-20 bg-gray-50 dark:bg-gray-800">
                        <tr className="text-left">
                            <th className="w-9 px-2 py-2">
                                <input
                                    type="checkbox"
                                    checked={todasMarcadas}
                                    onChange={() => setSeleccion(todasMarcadas ? new Set() : new Set(conCambios.map(p => p.sku)))}
                                    disabled={conCambios.length === 0}
                                    className="h-3.5 w-3.5 accent-blue-600"
                                />
                            </th>
                            <th className="min-w-[20rem] px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-gray-400">
                                Probability
                            </th>
                            {origenes.map(canal => (
                                <th key={canal.id} className="min-w-[11rem] px-3 py-2 text-center">
                                    <span className="inline-flex items-center gap-1.5 whitespace-nowrap text-[11.5px] font-semibold text-gray-700 dark:text-gray-200">
                                        <LogoCanal integracion={canal} />
                                        {canal.name}
                                        <span className="rounded-full bg-gray-100 px-1.5 py-0.5 text-[9.5px] font-bold uppercase tracking-wide text-gray-500 dark:bg-gray-700 dark:text-gray-300">
                                            ERP origen
                                        </span>
                                    </span>
                                </th>
                            ))}
                            <th className="w-28 px-3 py-2 text-center text-[10px] font-bold uppercase tracking-wider text-gray-400">
                                Stock en Probability
                            </th>
                            {ventas.map(canal => (
                                <th key={canal.id} className="min-w-[11rem] px-3 py-2 text-center">
                                    <span className="inline-flex items-center gap-1.5 whitespace-nowrap text-[11.5px] font-semibold text-gray-700 dark:text-gray-200">
                                        <LogoCanal integracion={canal} />
                                        {canal.name}
                                    </span>
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {visibles.map(producto => {
                            const cambia = cambiaEnAlgunCanal(producto.sku);
                            const propio = canales
                                .map(c => porCanal[c.id]?.[producto.sku]?.probability_qty)
                                .find(v => v !== undefined && v !== null);
                            return (
                                <tr
                                    key={producto.product_id}
                                    className="border-t bg-white transition-colors hover:bg-gray-50/80 dark:bg-gray-900 dark:hover:bg-gray-800/60"
                                    style={{ borderColor: CARD_BORDER }}
                                >
                                    <td className="px-2 py-2 align-top">
                                        <input
                                            type="checkbox"
                                            checked={seleccion.has(producto.sku)}
                                            onChange={() => alternar(producto.sku)}
                                            disabled={!cambia}
                                            className="h-3.5 w-3.5 accent-blue-600 disabled:opacity-30"
                                        />
                                    </td>
                                    <td className="px-3 py-2 align-top">
                                        <span className="mb-0.5 block max-w-[22rem] truncate text-[11.5px] font-semibold text-gray-700 dark:text-gray-200" title={producto.name}>
                                            {producto.name || 'Sin nombre'}
                                        </span>
                                        <span className="flex items-baseline gap-1.5 font-mono text-[11.5px] leading-tight">
                                            <span className="w-6 flex-shrink-0 text-gray-400">sku</span>
                                            <span className="font-semibold" style={{ color: ACCENT }}>{producto.sku}</span>
                                        </span>
                                    </td>
                                    {origenes.map(canal => (
                                        <td key={canal.id} className="px-3 py-2 text-center align-top">
                                            <Celda fila={porCanal[canal.id]?.[producto.sku]} />
                                        </td>
                                    ))}
                                    <td className="px-3 py-2 text-center align-top text-[12px] font-semibold tabular-nums text-gray-800 dark:text-gray-100">
                                        {propio ?? <span className="text-[11px] italic text-gray-300">sin dato</span>}
                                    </td>
                                    {ventas.map(canal => (
                                        <td key={canal.id} className="px-3 py-2 text-center align-top">
                                            <Celda fila={porCanal[canal.id]?.[producto.sku]} />
                                        </td>
                                    ))}
                                </tr>
                            );
                        })}
                    </tbody>
                </table>

                {(cargando || barriendo) && (
                    <p className="flex items-center justify-center gap-2 py-10 text-[12px] text-gray-500">
                        <Loader2 size={14} className="animate-spin" />
                        {barriendo
                            ? `Revisando todo el catalogo, grupo ${avanceBarrido.grupo} de ${avanceBarrido.grupos}. Lo encontrado ya se ve arriba.`
                            : 'Preguntandole a cada canal cuanto stock tiene ahora mismo'}
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
                        {soloCambios
                            ? 'En este grupo todos los productos ya estan iguales en sus canales.'
                            : 'No hay productos en este grupo.'}
                    </p>
                )}
            </div>

            <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="flex flex-wrap items-center gap-2">
                    {seleccion.size > 0 ? (
                        <span className="text-[11.5px] text-gray-500 dark:text-gray-400">
                            <span className="font-bold text-gray-700 dark:text-gray-200">{seleccion.size}</span> seleccionados
                        </span>
                    ) : conCambios.length > 0 && (
                        <span className="text-[11px] italic text-gray-400 dark:text-gray-500">
                            marca los que quieras igualar, o envia de una los {conCambios.length} que cambian
                        </span>
                    )}
                    {Object.values(results).some(r => r?.kind === 'inventory') && !enviando && (
                        <span className="rounded-full bg-emerald-50 px-2.5 py-0.5 text-[11.5px] font-semibold text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300">
                            {Object.values(results).reduce((suma, r) => suma + (r?.kind === 'inventory' ? r.updated : 0), 0)} actualizados
                        </span>
                    )}
                </div>

                <div className="flex items-center gap-1.5">
                    <button
                        onClick={() => enviar(seleccion.size > 0 ? [...seleccion] : conCambios.map(p => p.sku))}
                        disabled={(seleccion.size === 0 && conCambios.length === 0) || running || cargando}
                        className="ml-1 inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-[11.5px] font-bold text-white shadow-sm transition-opacity hover:opacity-90 disabled:opacity-40"
                        style={{ backgroundColor: ACCENT }}
                    >
                        {running ? <Loader2 size={12} className="animate-spin" /> : <Send size={12} />}
                        {seleccion.size > 0
                            ? `Igualar ${seleccion.size} en sus canales`
                            : `Igualar los ${conCambios.length} que cambian`}
                    </button>
                </div>
            </div>

            {soloCambios ? (
                <div className="flex flex-wrap items-center justify-between gap-2">
                    <span
                        className="inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1 text-[11.5px] font-semibold"
                        style={{ borderColor: ACCENT_BORDER, backgroundColor: ACCENT_SOFT, color: ACCENT }}
                    >
                        <span className="text-[13px] font-black">{productos.length.toLocaleString('es-CO')}</span>
                        diferencias encontradas
                        <span className="font-normal opacity-70">
                            {barriendo
                                ? `revisando ${avanceBarrido.grupo} de ${avanceBarrido.grupos} grupos`
                                : `en ${total.toLocaleString('es-CO')} productos`}
                        </span>
                    </span>
                    {!barriendo && (
                        <button
                            onClick={() => void barrerTodo()}
                            className="rounded-lg border px-2.5 py-1 text-[11.5px] font-semibold text-gray-600 transition-colors hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
                            style={{ borderColor: CARD_BORDER }}
                        >
                            Volver a revisar todo
                        </button>
                    )}
                </div>
            ) : (
                <PanelPager
                    page={pagina}
                    totalPages={totalPaginas}
                    total={total}
                    shown={productos.length}
                    noun="productos"
                    pageSize={tamano}
                    onPage={destino => { if (!cargando) setPagina(destino); }}
                    onPageSize={nuevo => { setTamano(nuevo); setPagina(1); }}
                />
            )}

            <p className="text-[11px] text-gray-400">
                Cada producto se envia solo a los canales donde el numero difiere. Al enviar, Probability vuelve a leer su
                propio stock en ese instante. El stock nunca viaja del canal hacia Probability.
            </p>
        </div>
    );
}
