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

export function HeroEditor({ content, onChange, businessId, onImageDeleted }: EditorProps) {
    const c = content || {};
    const set = (k: string, v: string) => onChange({ ...c, [k]: v });
    return (
        <div className="space-y-3">
            <Field label="Titulo" value={c.title || ''} onChange={(v) => set('title', v)} placeholder="Bienvenido a tu tienda" />
            <Field label="Subtitulo" value={c.subtitle || ''} onChange={(v) => set('subtitle', v)} placeholder="Descubre nuestros productos" />
            <Field label="Texto del boton" value={c.cta_text || ''} onChange={(v) => set('cta_text', v)} placeholder="Ver Productos" />
            <ImageField label="Imagen de fondo" value={c.background_image || ''} onChange={(v) => set('background_image', v)} businessId={businessId} onDeleted={onImageDeleted} />
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
                    <p className="text-xs font-medium text-blue-800 dark:text-blue-200">Sugerencia: usa la direccion de tu bodega</p>
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
                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Direccion</label>
                <AddressAutocomplete
                    value={c.address || ''}
                    onChange={(v) => set({ address: v })}
                    onSelect={(s: AddressSuggestion) => set({ lat: s.lat, lng: s.lon })}
                    placeholder="Busca la direccion de tu negocio"
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
                <p className="text-xs text-amber-600">El mapa y el boton de como llegar requieren ubicar la direccion (usa el buscador o una bodega).</p>
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
            <Field label="Telefono" value={c.phone || ''} onChange={(v) => set('phone', v)} placeholder="+57 300 000 0000" />
            <p className="text-xs text-gray-400">El formulario de contacto se muestra siempre que la seccion este visible.</p>
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
            <Field label="Numero (con indicativo)" value={c.number || ''} onChange={(v) => set('number', v)} placeholder="573000000000" />
            <Field label="Mensaje inicial" value={c.message || ''} onChange={(v) => set('message', v)} textarea placeholder="Hola, me gustaria mas informacion" />
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
