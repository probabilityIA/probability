'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { AlertTriangle, ArrowRight, Ban, Filter, History, Loader2, RefreshCw, Search, Send } from 'lucide-react';
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

type FotoCanal = Record<number, { guardada: boolean; cuando: string }>;

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

function estaEnCanal(producto: MatrixRow, integrationID: number): boolean {
    return producto.cells.some(celda => celda.integration_id === integrationID && celda.present);
}

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

function edad(iso: string) {
    if (!iso) return '';
    const fecha = new Date(iso);
    if (Number.isNaN(fecha.getTime()) || fecha.getFullYear() < 2000) return '';
    const minutos = Math.floor((Date.now() - fecha.getTime()) / 60000);
    if (minutos < 1) return 'hace un momento';
    if (minutos < 60) return `hace ${minutos} min`;
    const horas = Math.floor(minutos / 60);
    if (horas < 24) return `hace ${horas} h`;
    return `hace ${Math.floor(horas / 24)} d`;
}

function EstadoFoto({
    foto,
    onComparar,
    ocupado,
}: {
    foto?: { guardada: boolean; cuando: string };
    onComparar: () => void;
    ocupado: boolean;
}) {
    const guardada = foto?.guardada !== false;
    const cuando = edad(foto?.cuando ?? '');
    return (
        <button
            onClick={onComparar}
            disabled={ocupado}
            title={!guardada
                ? 'Reci\u00e9n consultado. Clic para volver a preguntar'
                : cuando === ''
                    ? 'Este canal no tiene comparaci\u00f3n guardada. Clic para preguntarle su stock ahora'
                    : 'Estos n\u00fameros son de la \u00faltima comparaci\u00f3n guardada. Clic para preguntarle el stock a este canal ahora'}
            className={`mx-auto mt-1 flex items-center gap-1 rounded-full border px-1.5 py-0.5 text-[9.5px] font-bold transition-colors disabled:opacity-40 ${
                !guardada
                    ? 'border-emerald-300 bg-emerald-50 text-emerald-700 hover:bg-emerald-100 dark:border-emerald-500/40 dark:bg-emerald-900/25 dark:text-emerald-200'
                    : cuando === ''
                        ? 'border-gray-200 bg-gray-50 text-gray-500 hover:bg-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400'
                        : 'border-amber-300 bg-amber-50 text-amber-800 hover:bg-amber-100 dark:border-amber-500/40 dark:bg-amber-900/25 dark:text-amber-200'
            }`}
        >
            {guardada ? <History size={9} /> : <RefreshCw size={9} />}
            {!guardada ? 'reci\u00e9n le\u00eddo' : cuando === '' ? 'sin comparar' : `guardado ${cuando}`}
        </button>
    );
}

function Celda({ fila, publicado, comparado }: { fila?: CompareRow; publicado: boolean; comparado: boolean }) {
    if (!fila) {
        if (!publicado) {
            return (
                <span
                    className="inline-flex items-center gap-1 text-[10.5px] italic text-gray-400 dark:text-gray-500"
                    title="Seg\u00fan el comparador de productos, este producto no est\u00e1 publicado en este canal"
                >
                    <Ban size={9} />
                    no esta aqui
                </span>
            );
        }
        if (!comparado) {
            return (
                <span
                    className="text-[10.5px] italic text-gray-300 dark:text-gray-600"
                    title="Esta publicado en este canal, pero todav\u00eda no se ha comparado su stock"
                >
                    sin comparar
                </span>
            );
        }
        return (
            <span
                className="text-[10.5px] italic text-amber-600 dark:text-amber-400"
                title="Esta publicado en este canal pero la comparaci\u00f3n no lo devolvi\u00f3. Vuelve a comparar este canal"
            >
                sin dato
            </span>
        );
    }

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
            {fila.channel_qty === null ? 'sin publicaci\u00f3n' : 'no aplica'}
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
    const [fotos, setFotos] = useState<FotoCanal>({});
    const cortar = useRef(false);

    const leerGrupo = useCallback(async (page: number, vivos?: Set<number>) => {
        const matriz = await fetchMatchMatrix(businessId ?? undefined, page, {}, undefined, tamano);
        const skus = matriz.rows.map(r => r.sku).filter(Boolean);
        if (skus.length === 0) return { matriz, mapa: {} as PorCanal, fotos: {} as FotoCanal };

        const lecturas = await Promise.all(
            canales.map(async canal => {
                const enVivo = vivos?.has(canal.id) === true;
                try {
                    const respuesta = await compararInventario(
                        canal.integration_type_id,
                        canal.id,
                        businessId ?? undefined,
                        1,
                        skus,
                        undefined,
                        enVivo ? undefined : { source: 'snapshot' },
                    );
                    return [canal.id, respuesta.rows, respuesta.from_cache === true, respuesta.checked_at] as const;
                } catch {
                    return [canal.id, [] as CompareRow[], !enVivo, ''] as const;
                }
            }),
        );

        const mapa: PorCanal = {};
        const fotos: FotoCanal = {};
        for (const [id, filas, deCache, cuando] of lecturas) {
            mapa[id] = Object.fromEntries(filas.map(f => [f.sku, f]));
            fotos[id] = { guardada: deCache, cuando };
        }
        return { matriz, mapa, fotos };
    }, [businessId, canales, tamano]);

    const compararCanal = useCallback(async (canalID: number) => {
        setCargando(true);
        setError(null);
        try {
            const { mapa, fotos } = await leerGrupo(pagina, new Set([canalID]));
            setPorCanal(previo => ({ ...previo, [canalID]: mapa[canalID] ?? {} }));
            setFotos(previo => ({ ...previo, [canalID]: fotos[canalID] ?? { guardada: false, cuando: '' } }));
            setLeidoA(new Date().toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' }));
        } catch (e) {
            setError(e instanceof Error ? e.message : 'No se pudo comparar el inventario');
        } finally {
            setCargando(false);
        }
    }, [leerGrupo, pagina]);

    const comparar = useCallback(async (page: number, vivos?: Set<number>) => {
        setCargando(true);
        setError(null);
        setSeleccion(new Set());
        try {
            const { matriz, mapa, fotos: foto } = await leerGrupo(page, vivos);
            setProductos(matriz.rows);
            setTotalPaginas(matriz.total_pages || 1);
            setTotal(matriz.total);
            setPorCanal(mapa);
            setFotos(foto);
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

                const { matriz, mapa, fotos: foto } = await leerGrupo(page);
                setFotos(previo => ({ ...previo, ...foto }));
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
                despues lo que tiene Probability, y despues cada canal de venta. Arranca con la
                <span className="font-semibold text-amber-700 dark:text-amber-300"> {'\u00faltima'} {'comparaci\u00f3n'} guardada</span>, para no
                pegarle a la API de cada canal cada vez que abres. Usa el boton bajo cada canal para preguntarle su stock
                ahora mismo. La flecha muestra en cuanto quedaria el canal si lo envias.
            </p>

            <p className="flex flex-wrap items-center gap-x-4 gap-y-1 text-[11px] text-gray-400 dark:text-gray-500">
                <span className="inline-flex items-center gap-1">
                    <Ban size={9} />
                    <span className="italic">{'no est\u00e1'} {'aqu\u00ed'}</span>: el producto no esta publicado en ese canal (sale del comparador de productos)
                </span>
                <span className="inline-flex items-center gap-1">
                    <span className="italic">sin comparar</span>: si esta publicado, pero todavia no le hemos preguntado su stock
                </span>
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
                        : 'Solo los que no est\u00e1n iguales'}
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
                    onClick={() => comparar(pagina, new Set(canales.map(canal => canal.id)))}
                    disabled={cargando}
                    title="Le pregunta el stock a todos los canales de este grupo"
                    className="ml-auto inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1.5 text-[11.5px] font-semibold transition-colors hover:opacity-80 disabled:opacity-50"
                    style={{ borderColor: ACCENT_BORDER, backgroundColor: ACCENT_SOFT, color: ACCENT }}
                >
                    {cargando ? <Loader2 size={12} className="animate-spin" /> : <RefreshCw size={12} />}
                    Comparar todos ahora
                </button>
                {leidoA && <span className="text-[11px] text-gray-400">Stock {'le\u00eddo'} a las {leidoA}</span>}
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
                                    <EstadoFoto foto={fotos[canal.id]} onComparar={() => void compararCanal(canal.id)} ocupado={cargando || barriendo} />
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
                                    <EstadoFoto foto={fotos[canal.id]} onComparar={() => void compararCanal(canal.id)} ocupado={cargando || barriendo} />
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
                                        <div className="flex gap-2">
                                            <FotoProducto url={producto.image_url} />
                                            <div className="min-w-0 flex-1">
                                                <span className="mb-0.5 block max-w-[20rem] truncate text-[11.5px] font-semibold text-gray-700 dark:text-gray-200" title={producto.name}>
                                                    {producto.name || 'Sin nombre'}
                                                </span>
                                                <span className="flex items-baseline gap-1.5 font-mono text-[11.5px] leading-tight">
                                                    <span className="w-6 flex-shrink-0 text-gray-400">sku</span>
                                                    <span className="font-semibold" style={{ color: ACCENT }}>{producto.sku}</span>
                                                </span>
                                            </div>
                                        </div>
                                    </td>
                                    {origenes.map(canal => (
                                        <td key={canal.id} className="px-3 py-2 text-center align-top">
                                            <Celda
                                                fila={porCanal[canal.id]?.[producto.sku]}
                                                publicado={estaEnCanal(producto, canal.id)}
                                                comparado={Boolean(porCanal[canal.id])}
                                            />
                                        </td>
                                    ))}
                                    <td className="px-3 py-2 text-center align-top text-[12px] font-semibold tabular-nums text-gray-800 dark:text-gray-100">
                                        {propio ?? <span className="text-[11px] italic text-gray-300">sin dato</span>}
                                    </td>
                                    {ventas.map(canal => (
                                        <td key={canal.id} className="px-3 py-2 text-center align-top">
                                            <Celda
                                                fila={porCanal[canal.id]?.[producto.sku]}
                                                publicado={estaEnCanal(producto, canal.id)}
                                                comparado={Boolean(porCanal[canal.id])}
                                            />
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
                            : 'Pregunt\u00e1ndole a cada canal cuanto stock tiene ahora mismo'}
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
                            ? 'En este grupo todos los productos ya est\u00e1n iguales en sus canales.'
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
                            ? `Igualar ${selecci\u00f3n.size} en sus canales`
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
