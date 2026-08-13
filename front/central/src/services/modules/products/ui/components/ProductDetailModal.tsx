'use client';

import { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { getProductDetailAction } from '../../infra/actions';
import { channelBrand } from '@/shared/utils/channel-brand';
import type { ProductChannel, ProductDetail, ProductFieldOrigin } from '../../domain/types';

const ETIQUETA_CAMPO: Record<string, string> = {
    name: 'Nombre',
    description: 'Descripcion',
    barcode: 'Codigo de barras',
    brand: 'Marca',
    category: 'Categoria',
    price: 'Precio',
    cost_price: 'Costo',
    image_url: 'Imagen',
    weight: 'Peso',
    dimensions: 'Medidas',
    variant_label: 'Variante',
};

interface ProductDetailModalProps {
    productId: string | null;
    businessId?: number;
    onClose: () => void;
}

const CARD_BG = '#fafafd';
const CARD_BORDER = '#eceaf3';
const ACCENT = 'var(--color-primary)';

function fecha(value?: string): string {
    if (!value) return 'nunca';
    const d = new Date(value);
    if (Number.isNaN(d.getTime())) return 'nunca';
    return d.toLocaleString('es-CO', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    });
}

function moneda(value: number, currency?: string): string {
    return new Intl.NumberFormat('es-CO', { style: 'currency', currency: currency || 'COP', maximumFractionDigits: 0 }).format(value);
}

function Copiar({ value, label }: { value: string; label: string }) {
    const [copiado, setCopiado] = useState(false);

    const copiar = async () => {
        try {
            await navigator.clipboard.writeText(value);
            setCopiado(true);
            setTimeout(() => setCopiado(false), 1400);
        } catch {
            setCopiado(false);
        }
    };

    return (
        <button
            onClick={copiar}
            title={copiado ? 'Copiado' : `Copiar ${label}`}
            aria-label={`Copiar ${label}`}
            className={`flex h-4 w-4 flex-shrink-0 items-center justify-center rounded transition-colors ${copiado ? 'text-emerald-600' : 'text-gray-300 hover:text-gray-600 dark:hover:text-gray-200'}`}
        >
            {copiado ? (
                <svg className="h-3 w-3" fill="none" stroke="currentColor" strokeWidth={3} viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" /></svg>
            ) : (
                <svg className="h-3 w-3" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
            )}
        </button>
    );
}

function Campo({ label, value, mono, hint, copiable }: { label: string; value?: string | number | null; mono?: boolean; hint?: string; copiable?: boolean }) {
    const vacio = value === undefined || value === null || value === '';
    return (
        <div className="flex items-baseline gap-1.5 text-[11.5px] leading-tight">
            <span className="w-[4.5rem] flex-shrink-0 font-mono text-gray-400">{label}</span>
            <span
                className={`min-w-0 truncate ${mono ? 'font-mono' : ''} ${vacio ? 'italic text-gray-300 dark:text-gray-600' : 'text-gray-700 dark:text-gray-200'}`}
                title={vacio ? undefined : String(value)}
            >
                {vacio ? 'null' : String(value)}
            </span>
            {!vacio && hint && <span className="flex-shrink-0 text-[10px] italic text-gray-400">{hint}</span>}
            {!vacio && copiable && <Copiar value={String(value)} label={label} />}
        </div>
    );
}

function Kpi({ label, value, hint, tone }: { label: string; value: string; hint?: string; tone?: string }) {
    return (
        <div className="rounded-xl border px-3 py-2 dark:bg-gray-800/60" style={{ borderColor: CARD_BORDER, backgroundColor: CARD_BG }}>
            <span className="block text-[10.5px] font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">{label}</span>
            <span className={`mt-0.5 block text-[17px] font-black leading-none ${tone ?? 'text-gray-900 dark:text-white'}`}>{value}</span>
            {hint && <span className="mt-0.5 block text-[10.5px] text-gray-500 dark:text-gray-400">{hint}</span>}
        </div>
    );
}

function Seccion({ title, count, children }: { title: string; count?: number; children: React.ReactNode }) {
    return (
        <section className="space-y-2">
            <h4 className="flex items-center gap-2 text-[11px] font-bold uppercase tracking-wider text-gray-400">
                {title}
                {count !== undefined && (
                    <span
                        className="rounded-full px-1.5 py-0.5 text-[10px] font-bold"
                        style={{ backgroundColor: `color-mix(in srgb, ${ACCENT} 12%, transparent)`, color: ACCENT }}
                    >
                        {count}
                    </span>
                )}
            </h4>
            {children}
        </section>
    );
}

function CanalCard({ channel, ownSKU }: { channel: ProductChannel; ownSKU: string }) {
    const brand = channelBrand(channel.integration_code);
    const empujado = channel.last_pushed_qty !== undefined && channel.last_pushed_qty !== null;
    const skuGuardado = Boolean(channel.external_sku);
    const skuUsado = channel.external_sku || ownSKU;

    return (
        <div className="rounded-xl border p-3 dark:bg-gray-800/60" style={{ borderColor: CARD_BORDER, backgroundColor: CARD_BG }}>
            <div className="mb-2 flex flex-wrap items-center gap-2">
                {channel.integration_image_url ? (
                    <img
                        src={channel.integration_image_url}
                        alt={channel.integration_name}
                        className="h-6 w-6 flex-shrink-0 rounded-md border bg-white object-contain p-0.5"
                        style={{ borderColor: CARD_BORDER }}
                    />
                ) : (
                    <span className="h-2.5 w-2.5 flex-shrink-0 rounded-full" style={{ backgroundColor: brand.dot }} />
                )}
                <span className="text-[12.5px] font-bold text-gray-900 dark:text-white">{channel.integration_name}</span>
                {!channel.is_active && (
                    <span className="rounded-full bg-gray-200 px-1.5 py-0.5 text-[10px] font-bold text-gray-600 dark:bg-gray-700 dark:text-gray-300">
                        integracion inactiva
                    </span>
                )}
                {channel.sku_matches && (
                    <span className="rounded-full bg-emerald-100 px-1.5 py-0.5 text-[10px] font-bold text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
                        SKU coincide
                    </span>
                )}
                {!channel.sku_matches && channel.external_sku && (
                    <span className="rounded-full bg-amber-100 px-1.5 py-0.5 text-[10px] font-bold text-amber-700 dark:bg-amber-900/40 dark:text-amber-300">
                        SKU distinto
                    </span>
                )}
                {!channel.sku_matches && !channel.external_sku && (
                    <span className="rounded-full bg-gray-200 px-1.5 py-0.5 text-[10px] font-bold text-gray-600 dark:bg-gray-700 dark:text-gray-300">
                        emparejado por SKU
                    </span>
                )}
                {channel.logistic_type && (
                    <span className="rounded-full bg-sky-100 px-1.5 py-0.5 text-[10px] font-bold text-sky-700 dark:bg-sky-900/40 dark:text-sky-300">
                        {channel.logistic_type}
                    </span>
                )}
            </div>

            <div className="grid gap-1 sm:grid-cols-2">
                <div className="space-y-1">
                    <Campo label="sku" value={skuUsado} mono copiable hint={skuGuardado ? undefined : 'del catalogo'} />
                    <Campo label="ean" value={channel.external_barcode} mono copiable />
                    <Campo label="id" value={channel.external_product_id} mono copiable />
                    <Campo label="variante" value={channel.external_variant_id} mono copiable />
                </div>
                <div className="space-y-1">
                    <Campo label="asociado" value={fecha(channel.linked_at)} />
                    <Campo label="stock env." value={empujado ? String(channel.last_pushed_qty) : 'nunca empujado'} />
                    <Campo label="ult. inv." value={`${fecha(channel.inventory_sync_at)}${channel.inventory_sync_status ? ` (${channel.inventory_sync_status})` : ''}`} />
                    <Campo label="ult. comp." value={`${fecha(channel.compare_at)}${channel.compare_status ? ` (${channel.compare_status})` : ''}`} />
                </div>
            </div>
        </div>
    );
}

function OrigenFila({ origin }: { origin: ProductFieldOrigin }) {
    const quien = origin.source === 'manual'
        ? origin.user_name || 'edicion manual'
        : origin.integration_name || 'un canal';
    const tono = origin.source === 'manual'
        ? 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'
        : 'bg-sky-100 text-sky-700 dark:bg-sky-900/40 dark:text-sky-300';

    return (
        <div className="flex items-baseline gap-2 border-t py-1.5 text-[11.5px] leading-tight first:border-t-0" style={{ borderColor: CARD_BORDER }}>
            <span className="w-[7rem] flex-shrink-0 font-semibold text-gray-700 dark:text-gray-200">
                {ETIQUETA_CAMPO[origin.field] ?? origin.field}
            </span>
            <span className={`flex-shrink-0 rounded-full px-1.5 py-0.5 text-[10px] font-bold ${tono}`}>{quien}</span>
            <span className="flex-shrink-0 text-gray-400">{fecha(origin.changed_at)}</span>
            {origin.previous_value && (
                <span className="min-w-0 truncate text-gray-400" title={origin.previous_value}>
                    antes: {origin.previous_value}
                </span>
            )}
        </div>
    );
}

export function ProductDetailModal({ productId, businessId, onClose }: ProductDetailModalProps) {
    const [detail, setDetail] = useState<ProductDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [montado, setMontado] = useState(false);

    useEffect(() => setMontado(true), []);

    useEffect(() => {
        if (!productId) return;
        let vivo = true;
        setLoading(true);
        setError(null);
        getProductDetailAction(productId, businessId)
            .then(res => {
                if (!vivo) return;
                if (res.success && res.data) setDetail(res.data as ProductDetail);
                else setError(res.message || 'No se pudo cargar el producto');
            })
            .catch(() => vivo && setError('No se pudo cargar el producto'))
            .finally(() => vivo && setLoading(false));
        return () => { vivo = false; };
    }, [productId, businessId]);

    useEffect(() => {
        const escape = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose(); };
        window.addEventListener('keydown', escape);
        return () => window.removeEventListener('keydown', escape);
    }, [onClose]);

    if (!montado || !productId) return null;

    const imagen = detail?.image_url || detail?.family?.image_url;
    const canales = detail?.channels ?? [];

    return createPortal(
        <div className="fixed inset-0 z-[300] flex items-center justify-center p-3 sm:p-5" role="dialog" aria-modal="true">
            <button aria-label="Cerrar detalle" onClick={onClose} className="absolute inset-0 cursor-default bg-gray-900/50 backdrop-blur-[2px]" />

            <div className="relative flex max-h-[92vh] w-full max-w-5xl flex-col overflow-hidden rounded-2xl bg-white shadow-2xl dark:bg-gray-900">
                <div className="flex items-start gap-3 border-b px-5 py-4" style={{ borderColor: CARD_BORDER }}>
                    {imagen ? (
                        <img src={imagen} alt={detail?.name ?? ''} className="h-14 w-14 flex-shrink-0 rounded-xl border object-cover" style={{ borderColor: CARD_BORDER }} />
                    ) : (
                        <div className="flex h-14 w-14 flex-shrink-0 items-center justify-center rounded-xl border text-[10px] text-gray-300" style={{ borderColor: CARD_BORDER, backgroundColor: CARD_BG }}>
                            sin foto
                        </div>
                    )}
                    <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-1.5">
                            <h3 className="truncate text-[16px] font-bold text-gray-900 dark:text-white">
                                {loading ? 'Cargando producto' : detail?.name || 'Sin nombre'}
                            </h3>
                            {detail?.name && <Copiar value={detail.name} label="el nombre" />}
                        </div>
                        {detail?.sku && (
                            <div className="mt-0.5 flex items-center gap-1.5">
                                <span className="truncate font-mono text-[12px] font-bold" style={{ color: ACCENT }}>{detail.sku}</span>
                                <Copiar value={detail.sku} label="el SKU" />
                            </div>
                        )}
                        {detail?.id && (
                            <div className="mt-0.5 flex items-center gap-1.5">
                                <span className="truncate font-mono text-[10.5px] text-gray-400">{detail.id}</span>
                                <Copiar value={detail.id} label="el id" />
                            </div>
                        )}
                    </div>
                    {detail && (
                        <span className={`flex-shrink-0 rounded-full px-2.5 py-1 text-[10.5px] font-bold uppercase tracking-wide ${detail.is_active
                            ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300'
                            : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300'}`}>
                            {detail.is_active ? 'Activo' : 'Inactivo'}
                        </span>
                    )}
                    <button
                        onClick={onClose}
                        aria-label="Cerrar"
                        className="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-gray-100 text-gray-500 transition-colors hover:bg-gray-200 hover:text-gray-700 dark:bg-gray-800 dark:text-gray-400"
                    >
                        <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>
                </div>

                <div className="min-h-0 flex-1 space-y-5 overflow-y-auto px-5 py-4">
                    {loading && <p className="py-10 text-center text-[12px] text-gray-500">Cargando el detalle del producto</p>}
                    {!loading && error && <p className="py-10 text-center text-[12px] text-red-600">{error}</p>}

                    {!loading && detail && (
                        <>
                            <div className="grid grid-cols-2 gap-2 sm:grid-cols-5">
                                <Kpi label="Precio" value={moneda(detail.price, detail.currency)} />
                                <Kpi
                                    label="Stock"
                                    value={detail.track_inventory ? String(detail.stock_total) : 'sin control'}
                                    hint={detail.track_inventory ? `${detail.reserved_total} reservado` : undefined}
                                />
                                <Kpi label="Disponible" value={detail.track_inventory ? String(detail.available_total) : '-'} />
                                <Kpi label="Canales" value={String(canales.length)} hint={canales.length === 0 ? 'sin asociar' : undefined} />
                                <Kpi label="Bodegas" value={String(detail.stock_levels.length)} />
                            </div>

                            <Seccion title="Canales asociados" count={canales.length}>
                                {canales.length === 0 ? (
                                    <p className="rounded-xl border px-3 py-6 text-center text-[12px] italic text-gray-400" style={{ borderColor: CARD_BORDER }}>
                                        Este producto no esta asociado a ninguna integracion
                                    </p>
                                ) : (
                                    <div className="grid gap-2 lg:grid-cols-2">
                                        {canales.map(channel => (
                                            <CanalCard
                                                key={`${channel.integration_id}-${channel.external_product_id}-${channel.external_variant_id ?? ''}`}
                                                channel={channel}
                                                ownSKU={detail.sku}
                                            />
                                        ))}
                                    </div>
                                )}
                            </Seccion>

                            <Seccion title="Inventario por bodega" count={detail.stock_levels.length}>
                                {detail.stock_levels.length === 0 ? (
                                    <p className="rounded-xl border px-3 py-6 text-center text-[12px] italic text-gray-400" style={{ borderColor: CARD_BORDER }}>
                                        Sin existencias registradas en bodega
                                    </p>
                                ) : (
                                    <div className="overflow-hidden rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                                        <table className="w-full border-collapse text-[11.5px]">
                                            <thead style={{ backgroundColor: CARD_BG }}>
                                                <tr className="text-left text-[10px] font-bold uppercase tracking-wider text-gray-400">
                                                    <th className="px-3 py-2">Bodega</th>
                                                    <th className="px-3 py-2 text-right">Cantidad</th>
                                                    <th className="px-3 py-2 text-right">Reservado</th>
                                                    <th className="px-3 py-2 text-right">Disponible</th>
                                                    <th className="px-3 py-2">Actualizado</th>
                                                </tr>
                                            </thead>
                                            <tbody>
                                                {detail.stock_levels.map(level => (
                                                    <tr key={level.warehouse_id} className="border-t" style={{ borderColor: CARD_BORDER }}>
                                                        <td className="px-3 py-2 text-gray-700 dark:text-gray-200">
                                                            {level.warehouse_name || `Bodega ${level.warehouse_id}`}
                                                            {level.warehouse_code && <span className="ml-1 font-mono text-[10px] text-gray-400">{level.warehouse_code}</span>}
                                                        </td>
                                                        <td className="px-3 py-2 text-right font-bold text-gray-900 dark:text-white">{level.quantity}</td>
                                                        <td className="px-3 py-2 text-right text-gray-500">{level.reserved_qty}</td>
                                                        <td className="px-3 py-2 text-right text-gray-500">{level.available_qty}</td>
                                                        <td className="px-3 py-2 text-gray-500">{fecha(level.updated_at)}</td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    </div>
                                )}
                            </Seccion>

                            <Seccion title="Ultimos movimientos" count={detail.movements.length}>
                                {detail.movements.length === 0 ? (
                                    <p className="rounded-xl border px-3 py-6 text-center text-[12px] italic text-gray-400" style={{ borderColor: CARD_BORDER }}>
                                        Sin movimientos de inventario
                                    </p>
                                ) : (
                                    <div className="overflow-hidden rounded-xl border" style={{ borderColor: CARD_BORDER }}>
                                        <table className="w-full border-collapse text-[11.5px]">
                                            <thead style={{ backgroundColor: CARD_BG }}>
                                                <tr className="text-left text-[10px] font-bold uppercase tracking-wider text-gray-400">
                                                    <th className="px-3 py-2">Fecha</th>
                                                    <th className="px-3 py-2">Tipo</th>
                                                    <th className="px-3 py-2">Origen</th>
                                                    <th className="px-3 py-2 text-right">Cantidad</th>
                                                    <th className="px-3 py-2 text-right">Antes</th>
                                                    <th className="px-3 py-2 text-right">Despues</th>
                                                </tr>
                                            </thead>
                                            <tbody>
                                                {detail.movements.map((mov, idx) => (
                                                    <tr key={`${mov.created_at}-${idx}`} className="border-t" style={{ borderColor: CARD_BORDER }}>
                                                        <td className="whitespace-nowrap px-3 py-2 text-gray-500">{fecha(mov.created_at)}</td>
                                                        <td className="px-3 py-2 text-gray-700 dark:text-gray-200">{mov.type_name || '-'}</td>
                                                        <td className="px-3 py-2 text-gray-500">{mov.integration_name || mov.reason || mov.warehouse_name || '-'}</td>
                                                        <td className={`px-3 py-2 text-right font-bold ${mov.direction === 'out' ? 'text-red-600' : 'text-emerald-600'}`}>
                                                            {mov.direction === 'out' ? '-' : '+'}{Math.abs(mov.quantity)}
                                                        </td>
                                                        <td className="px-3 py-2 text-right text-gray-400">{mov.previous_qty}</td>
                                                        <td className="px-3 py-2 text-right text-gray-700 dark:text-gray-200">{mov.new_qty}</td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    </div>
                                )}
                            </Seccion>

                            {detail.origins.length > 0 && (
                                <Seccion title="De donde salio cada dato" count={detail.origins.length}>
                                    <div className="rounded-xl border px-3 py-1.5" style={{ borderColor: CARD_BORDER, backgroundColor: CARD_BG }}>
                                        {detail.origins.map(origin => (
                                            <OrigenFila key={origin.field} origin={origin} />
                                        ))}
                                    </div>
                                </Seccion>
                            )}

                            <Seccion title="Ficha del producto">
                                <div className="grid gap-1.5 rounded-xl border p-3 sm:grid-cols-2" style={{ borderColor: CARD_BORDER, backgroundColor: CARD_BG }}>
                                    <Campo label="ean" value={detail.barcode} mono copiable />
                                    <Campo label="estado" value={detail.status} />
                                    <Campo label="familia" value={detail.family?.name} />
                                    <Campo label="variante" value={detail.variant_label} />
                                    <Campo label="hermanas" value={detail.sibling_variants > 0 ? `${detail.sibling_variants} variantes` : null} />
                                    <Campo label="categoria" value={detail.category} />
                                    <Campo label="marca" value={detail.brand} />
                                    <Campo label="costo" value={detail.cost_price ? moneda(detail.cost_price, detail.currency) : null} />
                                    <Campo label="peso" value={detail.weight ? `${detail.weight} ${detail.weight_unit || ''}`.trim() : null} />
                                    <Campo
                                        label="medidas"
                                        value={detail.length || detail.width || detail.height
                                            ? `${detail.length ?? 0} x ${detail.width ?? 0} x ${detail.height ?? 0} ${detail.dimension_unit || ''}`.trim()
                                            : null}
                                    />
                                    <Campo label="creado" value={fecha(detail.created_at)} />
                                    <Campo label="editado" value={fecha(detail.updated_at)} />
                                </div>
                                {detail.description && (
                                    <p className="rounded-xl border p-3 text-[12px] text-gray-600 dark:text-gray-300" style={{ borderColor: CARD_BORDER }}>
                                        {detail.description}
                                    </p>
                                )}
                            </Seccion>
                        </>
                    )}
                </div>
            </div>
        </div>,
        document.body,
    );
}

export default ProductDetailModal;
