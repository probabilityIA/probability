'use client';

import { useState, useEffect, useRef } from 'react';
import { TrashIcon, PlusIcon, BuildingStorefrontIcon, ArrowUpTrayIcon } from '@heroicons/react/24/outline';
import AddressAutocomplete, { AddressSuggestion } from '@/services/modules/orders/ui/components/AddressAutocomplete';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { uploadWebsiteImageAction } from '../../infra/actions';

type ContentRecord = Record<string, any>;

interface FieldProps {
    label: string;
    value: string;
    onChange: (v: string) => void;
    placeholder?: string;
    textarea?: boolean;
}

function Field({ label, value, onChange, placeholder, textarea }: FieldProps) {
    const cls = 'w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500';
    return (
        <div>
            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{label}</label>
            {textarea ? (
                <textarea value={value} onChange={(e) => onChange(e.target.value)} placeholder={placeholder} rows={3} className={cls} />
            ) : (
                <input type="text" value={value} onChange={(e) => onChange(e.target.value)} placeholder={placeholder} className={cls} />
            )}
        </div>
    );
}

interface EditorProps {
    content: ContentRecord | null;
    onChange: (content: ContentRecord) => void;
    businessId?: number;
    onImageDeleted?: () => void;
}

interface ImageFieldProps {
    label: string;
    value: string;
    onChange: (v: string) => void;
    businessId?: number;
    onDeleted?: () => void;
}

export function ImageField({ label, value, onChange, businessId, onDeleted }: ImageFieldProps) {
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState('');
    const fileRef = useRef<HTMLInputElement>(null);

    const handleFile = async (file: File | undefined) => {
        if (!file) return;
        setUploading(true);
        setError('');
        const formData = new FormData();
        formData.append('image', file);
        const result = await uploadWebsiteImageAction(formData, businessId);
        setUploading(false);
        if (result.success && result.image_url) {
            onChange(result.image_url);
        } else {
            setError('message' in result && result.message ? result.message : 'Error al subir imagen');
        }
        if (fileRef.current) fileRef.current.value = '';
    };

    return (
        <div>
            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{label}</label>
            <div className="flex gap-2">
                <input
                    type="text"
                    value={value}
                    onChange={(e) => onChange(e.target.value)}
                    placeholder="https://... o sube un archivo"
                    className="flex-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
                <button
                    type="button"
                    onClick={() => fileRef.current?.click()}
                    disabled={uploading}
                    className="shrink-0 inline-flex items-center gap-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                    title="Subir imagen"
                >
                    <ArrowUpTrayIcon className="w-4 h-4" />
                    {uploading ? 'Subiendo...' : 'Subir'}
                </button>
                <input
                    ref={fileRef}
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={(e) => handleFile(e.target.files?.[0])}
                />
            </div>
            {error && <p className="text-xs text-red-500 mt-1">{error}</p>}
            {value && (
                <div className="mt-2 relative inline-block">
                    <img src={value} alt="" className="h-20 rounded-md object-cover border border-gray-200 dark:border-gray-700" />
                    <button
                        type="button"
                        onClick={() => { onChange(''); onDeleted?.(); }}
                        className="absolute -top-2 -right-2 p-1 bg-red-500 text-white rounded-full shadow hover:bg-red-600"
                        title="Eliminar imagen definitivamente"
                    >
                        <TrashIcon className="w-3.5 h-3.5" />
                    </button>
                </div>
            )}
        </div>
    );
}

const HERO_BLOCK_TYPES: Array<{ value: string; label: string }> = [
    { value: 'empty', label: 'Vacio' },
    { value: 'text', label: 'Texto' },
    { value: 'image', label: 'Imagen' },
    { value: 'button', label: 'Bot\u00f3n' },
];

function resizeBlocks(blocks: ContentRecord[], rows: number, columns: number): ContentRecord[] {
    const total = rows * columns;
    const next = Array.from({ length: total }, (_, i) => blocks[i] || { type: 'empty' });
    return next;
}

function resizeColumnWidths(widths: number[], columns: number): number[] {
    return Array.from({ length: columns }, (_, i) => widths[i] || 1);
}

export function HeroEditor({ content, onChange, businessId, onImageDeleted }: EditorProps) {
    const c = content || {};
    const set = (k: string, v: string) => onChange({ ...c, [k]: v });
    const rows = c.layout_rows || 1;
    const columns = c.layout_columns || 1;
    const isGrid = rows > 1 || columns > 1;
    const blocks: ContentRecord[] = c.blocks || [];
    const columnWidths: number[] = resizeColumnWidths(c.column_widths || [], columns);

    const setLayout = (nextRows: number, nextColumns: number) => {
        const nextBlocks = (nextRows > 1 || nextColumns > 1) ? resizeBlocks(blocks, nextRows, nextColumns) : blocks;
        const nextWidths = (nextRows > 1 || nextColumns > 1) ? resizeColumnWidths(c.column_widths || [], nextColumns) : c.column_widths;
        onChange({ ...c, layout_rows: nextRows, layout_columns: nextColumns, blocks: nextBlocks, column_widths: nextWidths });
    };

    const setBlock = (index: number, patch: ContentRecord) => {
        const next = blocks.map((b, i) => i === index ? { ...b, ...patch } : b);
        onChange({ ...c, blocks: next });
    };

    const setColumnWidth = (index: number, width: number) => {
        const next = columnWidths.map((w, i) => i === index ? width : w);
        onChange({ ...c, column_widths: next });
    };

    return (
        <div className="space-y-4">
            <div>
                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{'Dise\u00f1o'}</label>
                <div className="flex gap-2 mb-2">
                    <button
                        type="button"
                        onClick={() => setLayout(1, 1)}
                        className={`px-3 py-1.5 text-xs rounded-lg border ${!isGrid ? 'border-blue-500 ring-1 ring-blue-500 text-blue-600' : 'border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300'}`}
                    >
                        Clasico
                    </button>
                    <button
                        type="button"
                        onClick={() => setLayout(rows > 1 || columns > 1 ? rows : 1, columns > 1 ? columns : 2)}
                        className={`px-3 py-1.5 text-xs rounded-lg border ${isGrid ? 'border-blue-500 ring-1 ring-blue-500 text-blue-600' : 'border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300'}`}
                    >
                        Filas y columnas
                    </button>
                </div>
                {isGrid && (
                    <div className="flex gap-3">
                        <label className="text-xs text-gray-600 dark:text-gray-300 flex items-center gap-2">
                            Filas
                            <select
                                value={rows}
                                onChange={(e) => setLayout(Number(e.target.value), columns)}
                                className="text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white px-2 py-1"
                            >
                                {[1, 2, 3].map(n => <option key={n} value={n}>{n}</option>)}
                            </select>
                        </label>
                        <label className="text-xs text-gray-600 dark:text-gray-300 flex items-center gap-2">
                            Columnas
                            <select
                                value={columns}
                                onChange={(e) => setLayout(rows, Number(e.target.value))}
                                className="text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white px-2 py-1"
                            >
                                {[1, 2, 3].map(n => <option key={n} value={n}>{n}</option>)}
                            </select>
                        </label>
                    </div>
                )}
                {isGrid && columns > 1 && (
                    <div className="mt-3">
                        <p className="text-xs text-gray-600 dark:text-gray-300 mb-1">Ancho de columnas</p>
                        <div className="flex gap-3">
                            {columnWidths.map((width, i) => (
                                <label key={i} className="flex-1 text-xs text-gray-500 dark:text-gray-400">
                                    Col. {i + 1}
                                    <input
                                        type="range"
                                        min={1}
                                        max={5}
                                        step={0.5}
                                        value={width}
                                        onChange={(e) => setColumnWidth(i, Number(e.target.value))}
                                        className="w-full"
                                    />
                                </label>
                            ))}
                        </div>
                    </div>
                )}
            </div>

            {!isGrid && (
                <>
                    <Field label="T\u00edtulo" value={c.title || ''} onChange={(v) => set('title', v)} placeholder="Bienvenido a tu tienda" />
                    <label className="block text-xs text-gray-500 dark:text-gray-400">
                        Tamano del titulo: {c.title_size || 40}px
                        <input
                            type="range" min={20} max={64} step={2}
                            value={c.title_size || 40}
                            onChange={(e) => onChange({ ...c, title_size: Number(e.target.value) })}
                            className="w-full"
                        />
                    </label>
                    <Field label="Subtitulo" value={c.subtitle || ''} onChange={(v) => set('subtitle', v)} placeholder="Descubre nuestros productos" />
                    <label className="block text-xs text-gray-500 dark:text-gray-400">
                        Tamano del subtitulo: {c.subtitle_size || 20}px
                        <input
                            type="range" min={12} max={36} step={2}
                            value={c.subtitle_size || 20}
                            onChange={(e) => onChange({ ...c, subtitle_size: Number(e.target.value) })}
                            className="w-full"
                        />
                    </label>
                    <Field label="Texto del bot\u00f3n" value={c.cta_text || ''} onChange={(v) => set('cta_text', v)} placeholder="Ver Productos" />
                    <button
                        type="button"
                        onClick={() => onChange({ ...c, position: { x: 50, y: 50 } })}
                        className="text-xs text-blue-600 hover:underline"
                    >
                        Centrar contenido
                    </button>
                </>
            )}

            <ImageField label="Imagen de fondo" value={c.background_image || ''} onChange={(v) => set('background_image', v)} businessId={businessId} onDeleted={onImageDeleted} />

            {isGrid && (
                <div className="space-y-3 border-t border-gray-200 dark:border-gray-700 pt-3">
                    <p className="text-xs font-medium text-gray-600 dark:text-gray-300">Celdas ({rows} x {columns})</p>
                    <p className="text-xs text-gray-400">Arrastra el contenido dentro de la vista previa para ubicarlo; se ajusta solo al centro de cada celda.</p>
                    <div className="grid gap-3" style={{ gridTemplateColumns: `repeat(${columns}, 1fr)` }}>
                        {blocks.map((block, i) => (
                            <div key={i} className="border border-gray-200 dark:border-gray-600 rounded-lg p-3 space-y-2">
                                <label className="block text-xs text-gray-600 dark:text-gray-300">
                                    Tipo de contenido
                                    <select
                                        value={block.type || 'empty'}
                                        onChange={(e) => setBlock(i, { type: e.target.value })}
                                        className="mt-1 w-full text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white px-2 py-1.5"
                                    >
                                        {HERO_BLOCK_TYPES.map(t => <option key={t.value} value={t.value}>{t.label}</option>)}
                                    </select>
                                </label>

                                {block.type === 'text' && (
                                    <>
                                        <Field label="T\u00edtulo" value={block.title || ''} onChange={(v) => setBlock(i, { title: v })} />
                                        <label className="block text-xs text-gray-500 dark:text-gray-400">
                                            Tamano del titulo: {block.title_size || 24}px
                                            <input
                                                type="range" min={14} max={48} step={2}
                                                value={block.title_size || 24}
                                                onChange={(e) => setBlock(i, { title_size: Number(e.target.value) })}
                                                className="w-full"
                                            />
                                        </label>
                                        <Field label="Texto" value={block.subtitle || ''} onChange={(v) => setBlock(i, { subtitle: v })} textarea />
                                        <label className="block text-xs text-gray-500 dark:text-gray-400">
                                            Tamano del texto: {block.subtitle_size || 16}px
                                            <input
                                                type="range" min={10} max={28} step={2}
                                                value={block.subtitle_size || 16}
                                                onChange={(e) => setBlock(i, { subtitle_size: Number(e.target.value) })}
                                                className="w-full"
                                            />
                                        </label>
                                    </>
                                )}

                                {block.type === 'image' && (
                                    <ImageField label="Imagen" value={block.image || ''} onChange={(v) => setBlock(i, { image: v })} businessId={businessId} onDeleted={onImageDeleted} />
                                )}

                                {block.type === 'button' && (
                                    <>
                                        <Field label="Texto del bot\u00f3n" value={block.button_text || ''} onChange={(v) => setBlock(i, { button_text: v })} placeholder="Ver Productos" />
                                        <Field label="Enlace (ruta relativa o URL)" value={block.button_link || ''} onChange={(v) => setBlock(i, { button_link: v })} placeholder="/productos" />
                                    </>
                                )}

                                {block.type !== 'empty' && (
                                    <button
                                        type="button"
                                        onClick={() => setBlock(i, { position: { x: 50, y: 50 } })}
                                        className="text-xs text-blue-600 hover:underline"
                                    >
                                        Centrar en la celda
                                    </button>
                                )}
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}

export function AboutEditor({ content, onChange, businessId, onImageDeleted }: EditorProps) {
    const c = content || {};
    const set = (k: string, v: string) => onChange({ ...c, [k]: v });
    return (
        <div className="space-y-3">
            <Field label="Texto principal" value={c.text || ''} onChange={(v) => set('text', v)} textarea placeholder="Quienes somos..." />
            <Field label="Mision" value={c.mission || ''} onChange={(v) => set('mission', v)} textarea />
            <Field label="Vision" value={c.vision || ''} onChange={(v) => set('vision', v)} textarea />
            <ImageField label="Imagen" value={c.image || ''} onChange={(v) => set('image', v)} businessId={businessId} onDeleted={onImageDeleted} />
        </div>
    );
}

interface TestimonialsEditorProps {
    content: ContentRecord[] | null;
    onChange: (content: ContentRecord[]) => void;
}

export function TestimonialsEditor({ content, onChange }: TestimonialsEditorProps) {
    const items = content || [];
    const setItem = (i: number, k: string, v: string | number) => {
        const next = items.map((item, idx) => idx === i ? { ...item, [k]: v } : item);
        onChange(next);
    };
    return (
        <div className="space-y-4">
            {items.map((t, i) => (
                <div key={i} className="border border-gray-200 dark:border-gray-600 rounded-lg p-3 space-y-2 relative">
                    <button
                        type="button"
                        onClick={() => onChange(items.filter((_, idx) => idx !== i))}
                        className="absolute top-2 right-2 text-red-400 hover:text-red-600"
                        title="Eliminar testimonio"
                    >
                        <TrashIcon className="w-4 h-4" />
                    </button>
                    <Field label="Nombre" value={t.name || ''} onChange={(v) => setItem(i, 'name', v)} />
                    <Field label="Testimonio" value={t.text || ''} onChange={(v) => setItem(i, 'text', v)} textarea />
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Calificacion (1-5)</label>
                        <input
                            type="number" min={1} max={5}
                            value={t.rating || 5}
                            onChange={(e) => setItem(i, 'rating', Number(e.target.value))}
                            className="w-20 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                        />
                    </div>
                </div>
            ))}
            <button
                type="button"
                onClick={() => onChange([...items, { name: '', text: '', rating: 5 }])}
                className="w-full flex items-center justify-center gap-1 py-2 text-sm text-blue-600 border border-dashed border-blue-300 rounded-lg hover:bg-blue-50 dark:hover:bg-blue-900/20"
            >
                <PlusIcon className="w-4 h-4" /> Agregar testimonio
            </button>
        </div>
    );
}

interface WarehouseOption {
    id: number;
    name: string;
    address: string;
    city: string;
    latitude: number | null;
    longitude: number | null;
}

interface LocationEditorProps extends EditorProps {
    businessId?: number;
}

export function LocationEditor({ content, onChange, businessId }: LocationEditorProps) {
    const c = content || {};
    const latest = useRef<ContentRecord>(c);
    latest.current = c;
    const set = (patch: ContentRecord) => {
        latest.current = { ...latest.current, ...patch };
        onChange(latest.current);
    };
    const [warehouses, setWarehouses] = useState<WarehouseOption[]>([]);

    useEffect(() => {
        let cancelled = false;
        getWarehousesAction({ page: 1, page_size: 10, is_active: true, business_id: businessId })
            .then((res: any) => {
                if (cancelled) return;
                const list = (res?.data || []) as WarehouseOption[];
                setWarehouses(list.filter(w => w.address));
            })
            .catch(() => setWarehouses([]));
        return () => { cancelled = true; };
    }, [businessId]);

    const applyWarehouse = (w: WarehouseOption) => {
        set({
            address: w.city ? `${w.address}, ${w.city}` : w.address,
            lat: w.latitude ?? undefined,
            lng: w.longitude ?? undefined,
        });
    };

    return (
        <div className="space-y-3">
            {warehouses.length > 0 && (
                <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3 space-y-2">
                    <p className="text-xs font-medium text-blue-800 dark:text-blue-200">Sugerencia: usa la {'direcci\u00f3n'} de tu bodega</p>
                    {warehouses.map(w => (
                        <button
                            key={w.id}
                            type="button"
                            onClick={() => applyWarehouse(w)}
                            className="w-full flex items-start gap-2 text-left px-2 py-1.5 rounded-md hover:bg-blue-100 dark:hover:bg-blue-900/40"
                        >
                            <BuildingStorefrontIcon className="w-4 h-4 text-blue-500 mt-0.5 shrink-0" />
                            <span className="text-xs text-blue-900 dark:text-blue-100">
                                <span className="font-medium">{w.name}</span>: {w.address}{w.city ? `, ${w.city}` : ''}
                            </span>
                        </button>
                    ))}
                </div>
            )}

            <div>
                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{'Direcci\u00f3n'}</label>
                <AddressAutocomplete
                    value={c.address || ''}
                    onChange={(v) => set({ address: v })}
                    onSelect={(s: AddressSuggestion) => set({ lat: s.lat, lng: s.lon })}
                    placeholder="Busca la direcci\u00f3n de tu negocio"
                />
                <p className="text-xs text-gray-400 mt-1">Al elegir una sugerencia se ubica el punto en el mapa automaticamente.</p>
            </div>

            <Field label="Horarios" value={c.hours || ''} onChange={(v) => set({ hours: v })} textarea placeholder={'Lun-Vie: 9am - 6pm\nSab: 9am - 1pm'} />

            <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                <input
                    type="checkbox"
                    checked={c.show_map !== false}
                    onChange={(e) => set({ show_map: e.target.checked })}
                    className="rounded border-gray-300"
                    disabled={!c.lat || !c.lng}
                />
                Mostrar mapa (OpenStreetMap)
            </label>
            <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                <input
                    type="checkbox"
                    checked={!!c.show_directions}
                    onChange={(e) => set({ show_directions: e.target.checked })}
                    className="rounded border-gray-300"
                    disabled={!c.lat || !c.lng}
                />
                Mostrar boton &quot;Como llegar&quot;
            </label>
            {(!c.lat || !c.lng) && (
                <p className="text-xs text-amber-600">El mapa y el {'bot\u00f3n'} de como llegar requieren ubicar la {'direcci\u00f3n'} (usa el buscador o una bodega).</p>
            )}
        </div>
    );
}

export function ContactEditor({ content, onChange }: EditorProps) {
    const c = content || {};
    const set = (k: string, v: string) => onChange({ ...c, [k]: v });
    return (
        <div className="space-y-3">
            <Field label="Email de contacto" value={c.email || ''} onChange={(v) => set('email', v)} placeholder="ventas@negocio.com" />
            <Field label="Tel\u00e9fono" value={c.phone || ''} onChange={(v) => set('phone', v)} placeholder="+57 300 000 0000" />
            <p className="text-xs text-gray-400">El formulario de contacto se muestra siempre que la {'secci\u00f3n'} este visible.</p>
        </div>
    );
}

export function SocialMediaEditor({ content, onChange }: EditorProps) {
    const c = content || {};
    const set = (k: string, v: string) => onChange({ ...c, [k]: v });
    return (
        <div className="space-y-3">
            <Field label="Facebook (URL)" value={c.facebook || ''} onChange={(v) => set('facebook', v)} placeholder="https://facebook.com/..." />
            <Field label="Instagram (URL)" value={c.instagram || ''} onChange={(v) => set('instagram', v)} placeholder="https://instagram.com/..." />
            <Field label="Twitter / X (URL)" value={c.twitter || ''} onChange={(v) => set('twitter', v)} placeholder="https://x.com/..." />
            <Field label="TikTok (URL)" value={c.tiktok || ''} onChange={(v) => set('tiktok', v)} placeholder="https://tiktok.com/@..." />
        </div>
    );
}

export function WhatsAppEditor({ content, onChange }: EditorProps) {
    const c = content || {};
    const set = (k: string, v: string | boolean) => onChange({ ...c, [k]: v });
    return (
        <div className="space-y-3">
            <Field label="N\u00famero (con indicativo)" value={c.number || ''} onChange={(v) => set('number', v)} placeholder="573000000000" />
            <Field label="Mensaje inicial" value={c.message || ''} onChange={(v) => set('message', v)} textarea placeholder="Hola, me gustaria m\u00e1s informaci\u00f3n" />
            <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                <input
                    type="checkbox"
                    checked={!!c.show_floating_button}
                    onChange={(e) => set('show_floating_button', e.target.checked)}
                    className="rounded border-gray-300"
                />
                Mostrar boton flotante
            </label>
        </div>
    );
}

const NAVBAR_POSITIONS: Array<{ value: string; label: string }> = [
    { value: 'top', label: 'Arriba' },
    { value: 'left', label: 'Lateral izquierdo' },
    { value: 'right', label: 'Lateral derecho' },
];

const NAVBAR_ALIGNMENTS: Array<{ value: string; label: string }> = [
    { value: 'left', label: 'Izquierda' },
    { value: 'center', label: 'Centro' },
    { value: 'right', label: 'Derecha' },
];

const NAV_ICON_PRESETS = ['⌂', '▤', '✉', '☎', '★', '●', '■', '▲', '✓', '→', '♥', '☰', '⚙', '⚑'];

export function NavbarEditor({ content, onChange, businessId, onImageDeleted }: EditorProps) {
    const c = content || {};
    const position = c.position || 'top';
    const links: ContentRecord[] = c.links && c.links.length > 0 ? c.links : [
        { label: 'Inicio', href: '' },
        { label: 'Productos', href: '/productos' },
        { label: 'Contacto', href: '/contacto' },
    ];

    const setLink = (i: number, patch: ContentRecord) => {
        const next = links.map((l, idx) => idx === i ? { ...l, ...patch } : l);
        onChange({ ...c, links: next });
    };

    const moveLink = (i: number, dir: -1 | 1) => {
        const target = i + dir;
        if (target < 0 || target >= links.length) return;
        const next = [...links];
        [next[i], next[target]] = [next[target], next[i]];
        onChange({ ...c, links: next });
    };

    const removeLink = (i: number) => {
        onChange({ ...c, links: links.filter((_, idx) => idx !== i) });
    };

    return (
        <div className="space-y-4">
            <div>
                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{'Posici\u00f3n'}</label>
                <select
                    value={position}
                    onChange={(e) => onChange({ ...c, position: e.target.value })}
                    className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                >
                    {NAVBAR_POSITIONS.map(p => <option key={p.value} value={p.value}>{p.label}</option>)}
                </select>
            </div>

            {position === 'top' && (
                <div>
                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Alineacion de los links</label>
                    <select
                        value={c.alignment || 'right'}
                        onChange={(e) => onChange({ ...c, alignment: e.target.value })}
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                    >
                        {NAVBAR_ALIGNMENTS.map(a => <option key={a.value} value={a.value}>{a.label}</option>)}
                    </select>
                </div>
            )}

            {(position === 'left' || position === 'right') && (
                <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                    <input
                        type="checkbox"
                        checked={!!c.collapsible}
                        onChange={(e) => onChange({ ...c, collapsible: e.target.checked })}
                        className="rounded border-gray-300"
                    />
                    Modo desplegable (se expande al pasar el mouse, muestra emojis cuando esta contraido)
                </label>
            )}

            <div className="space-y-2">
                <p className="text-xs font-medium text-gray-600 dark:text-gray-300">Enlaces del menu</p>
                {links.map((link, i) => (
                    <div key={i} className="border border-gray-200 dark:border-gray-600 rounded-lg p-3 space-y-2">
                        <div className="flex gap-2">
                            <input
                                type="text"
                                value={link.label || ''}
                                onChange={(e) => setLink(i, { label: e.target.value })}
                                placeholder="Texto"
                                className="flex-1 px-2 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                            />
                            <input
                                type="text"
                                value={link.href || ''}
                                onChange={(e) => setLink(i, { href: e.target.value })}
                                placeholder="/ruta o https://..."
                                className="flex-1 px-2 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                            />
                        </div>

                        {(position === 'left' || position === 'right') && c.collapsible && (
                            <div className="space-y-2 border-t border-gray-100 dark:border-gray-700 pt-2">
                                <div className="flex gap-2 text-xs">
                                    <button
                                        type="button"
                                        onClick={() => setLink(i, { icon_type: 'symbol' })}
                                        className={`px-2 py-1 rounded border ${link.icon_type !== 'image' ? 'border-blue-500 text-blue-600' : 'border-gray-200 dark:border-gray-600 text-gray-500'}`}
                                    >
                                        Simbolo
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => setLink(i, { icon_type: 'image' })}
                                        className={`px-2 py-1 rounded border ${link.icon_type === 'image' ? 'border-blue-500 text-blue-600' : 'border-gray-200 dark:border-gray-600 text-gray-500'}`}
                                    >
                                        Imagen propia
                                    </button>
                                </div>

                                {link.icon_type === 'image' ? (
                                    <ImageField
                                        label="Icono del negocio"
                                        value={link.icon_image || ''}
                                        onChange={(v) => setLink(i, { icon_image: v })}
                                        businessId={businessId}
                                        onDeleted={onImageDeleted}
                                    />
                                ) : (
                                    <div className="flex flex-wrap gap-1.5">
                                        {NAV_ICON_PRESETS.map(sym => (
                                            <button
                                                key={sym}
                                                type="button"
                                                onClick={() => setLink(i, { icon: sym })}
                                                className={`w-8 h-8 flex items-center justify-center rounded border text-base ${link.icon === sym ? 'border-blue-500 ring-1 ring-blue-500' : 'border-gray-200 dark:border-gray-600'}`}
                                            >
                                                {sym}
                                            </button>
                                        ))}
                                    </div>
                                )}
                            </div>
                        )}

                        <div className="flex items-center justify-between">
                            <div className="flex gap-1">
                                <button type="button" onClick={() => moveLink(i, -1)} disabled={i === 0} className="px-2 py-1 text-xs border border-gray-200 dark:border-gray-600 rounded disabled:opacity-30">
                                    Subir
                                </button>
                                <button type="button" onClick={() => moveLink(i, 1)} disabled={i === links.length - 1} className="px-2 py-1 text-xs border border-gray-200 dark:border-gray-600 rounded disabled:opacity-30">
                                    Bajar
                                </button>
                            </div>
                            <button type="button" onClick={() => removeLink(i)} className="text-red-400 hover:text-red-600">
                                <TrashIcon className="w-4 h-4" />
                            </button>
                        </div>
                    </div>
                ))}
                <button
                    type="button"
                    onClick={() => onChange({ ...c, links: [...links, { label: 'Nuevo', href: '/' }] })}
                    className="w-full flex items-center justify-center gap-1 py-2 text-sm text-blue-600 border border-dashed border-blue-300 rounded-lg hover:bg-blue-50 dark:hover:bg-blue-900/20"
                >
                    <PlusIcon className="w-4 h-4" /> Agregar enlace
                </button>
            </div>

            <div className="space-y-2 border-t border-gray-200 dark:border-gray-700 pt-3">
                <Field label="Barra de anuncio (opcional)" value={c.announcement_text || ''} onChange={(v) => onChange({ ...c, announcement_text: v })} placeholder="Env\u00edo gratis en compras +$100.000" />
                {c.announcement_text && (
                    <>
                        <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                            Color de fondo
                            <input
                                type="color"
                                value={c.announcement_color || '#111827'}
                                onChange={(e) => onChange({ ...c, announcement_color: e.target.value })}
                                className="w-8 h-8 rounded cursor-pointer border border-gray-300 dark:border-gray-600 bg-transparent p-0"
                            />
                        </label>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Animacion</label>
                            <select
                                value={c.announcement_animation || 'none'}
                                onChange={(e) => onChange({ ...c, announcement_animation: e.target.value })}
                                className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                            >
                                <option value="none">Sin animacion</option>
                                <option value="fade">Desvanecer al cargar</option>
                                <option value="slide">Deslizar desde arriba</option>
                                <option value="marquee">Texto en movimiento (marquesina)</option>
                            </select>
                        </div>
                    </>
                )}
            </div>
        </div>
    );
}
