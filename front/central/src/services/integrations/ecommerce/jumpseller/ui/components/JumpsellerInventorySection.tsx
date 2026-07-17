'use client';

import { useState, useEffect } from 'react';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { ExclamationTriangleIcon, InformationCircleIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';

export type JumpsellerInventoryMode = 'single' | 'mapped';

export interface JumpsellerLocationMapping {
    internal_warehouse_id: number;
    jumpseller_location_id: number;
}

export interface JumpsellerInventoryConfig {
    enabled: boolean;
    mode: JumpsellerInventoryMode;
    single_warehouse_id: number;
    default_location_id: string;
    mappings: JumpsellerLocationMapping[];
}

interface NamedWarehouse {
    id: number;
    name: string;
    code?: string;
}

export interface JumpsellerLocation {
    id: number;
    name: string;
    city?: string;
    country?: string;
    main?: boolean;
    is_stock_origin?: boolean;
}

export interface JumpsellerLocationsInfo {
    locations: JumpsellerLocation[];
    multi_location: boolean;
    subscription_plan: string;
    stock_origin_name: string;
}

interface JumpsellerInventorySectionProps {
    value: JumpsellerInventoryConfig;
    onChange: (v: JumpsellerInventoryConfig) => void;
    businessId: number | null;
    locationsInfo: JumpsellerLocationsInfo | null;
}

const ACCENT = 'var(--color-primary)';
const ACCENT_SOFT = 'color-mix(in srgb, var(--color-primary) 10%, white)';
const INPUT_BORDER = '#e9e9f0';

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const rowLabel = 'block text-[12px] font-semibold text-gray-700 dark:text-gray-200 mb-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';

const MODES: Array<{ id: JumpsellerInventoryMode; title: string; hint: string }> = [
    { id: 'single', title: 'Una bodega', hint: 'El stock de una bodega' },
    { id: 'mapped', title: 'Emparejar bodegas', hint: 'Cada bodega a una de Jumpseller' },
];

const locationLabel = (loc: JumpsellerLocation) => {
    const place = loc.city ? `${loc.name} - ${loc.city}` : loc.name;
    const tags: string[] = [];
    if (loc.is_stock_origin) tags.push('origen de stock');
    if (loc.main) tags.push('principal');
    return tags.length > 0 ? `${place} (${tags.join(', ')})` : place;
};

const warehouseLabel = (w: NamedWarehouse) => `${w.name}${w.code ? ` (${w.code})` : ''}`;

const countDistinctLocations = (mappings: JumpsellerLocationMapping[]) => {
    const ids = new Set<number>();
    mappings.forEach((m) => {
        if (m.jumpseller_location_id > 0) ids.add(m.jumpseller_location_id);
    });
    return ids.size;
};

export function JumpsellerInventorySection({ value, onChange, businessId, locationsInfo }: JumpsellerInventorySectionProps) {
    const [warehouses, setWarehouses] = useState<NamedWarehouse[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (!value.enabled) return;
        let cancelled = false;
        setLoading(true);
        getWarehousesAction(businessId ? ({ business_id: businessId, page_size: 100 } as any) : ({ page_size: 100 } as any))
            .then((res: any) => {
                if (cancelled) return;
                const list = res?.data || res?.items || (Array.isArray(res) ? res : []);
                setWarehouses(list.map((w: any) => ({ id: w.id, name: w.name, code: w.code })));
            })
            .catch(() => { })
            .finally(() => { if (!cancelled) setLoading(false); });
        return () => { cancelled = true; };
    }, [value.enabled, businessId]);

    const set = (patch: Partial<JumpsellerInventoryConfig>) => onChange({ ...value, ...patch });

    const updateRow = (index: number, patch: Partial<JumpsellerLocationMapping>) => {
        set({ mappings: value.mappings.map((row, i) => (i === index ? { ...row, ...patch } : row)) });
    };
    const addRow = () => set({ mappings: [...value.mappings, { internal_warehouse_id: 0, jumpseller_location_id: 0 }] });
    const removeRow = (index: number) => set({ mappings: value.mappings.filter((_, i) => i !== index) });

    if (!value.enabled) return null;

    const locations = locationsInfo?.locations ?? [];
    const tooManyLocations = value.mode === 'mapped' && countDistinctLocations(value.mappings) > 1;

    return (
        <div className="mt-4 space-y-4">
            <div>
                <label className={fieldLabel}>Que stock enviar a Jumpseller</label>
                <div className="grid grid-cols-2 gap-2">
                    {MODES.map((m) => {
                        const disabled = m.id === 'mapped';
                        const active = !disabled && value.mode === m.id;
                        return (
                            <button
                                key={m.id}
                                type="button"
                                onClick={() => { if (!disabled) set({ mode: m.id }); }}
                                disabled={disabled}
                                aria-pressed={active}
                                title={disabled ? 'Jumpseller todavia no permite enviar stock por bodega' : undefined}
                                className="rounded-lg border px-3 py-2 text-left text-[12px] transition-colors disabled:cursor-not-allowed disabled:opacity-60"
                                style={active
                                    ? { borderColor: ACCENT, backgroundColor: ACCENT_SOFT }
                                    : { borderColor: INPUT_BORDER }}
                            >
                                <span className="font-semibold text-gray-900 dark:text-white block">
                                    {m.title}
                                    {disabled && (
                                        <span className="ml-1.5 rounded px-1.5 py-0.5 text-[9px] font-semibold uppercase tracking-wide text-amber-700 bg-amber-100 dark:text-amber-300 dark:bg-amber-900/40">
                                            En desarrollo
                                        </span>
                                    )}
                                </span>
                                <span className="text-[11px] text-gray-500 dark:text-gray-400">{m.hint}</span>
                            </button>
                        );
                    })}
                </div>
                <p className="mt-1.5 text-[11px] text-gray-400 dark:text-gray-500 leading-snug">
                    Emparejar bodegas esta en desarrollo: la API de Jumpseller todavia no permite enviar el stock a una
                    bodega especifica.
                </p>
            </div>

            {locationsInfo && !locationsInfo.multi_location && (
                <p className="text-[11px] text-gray-500 dark:text-gray-400 leading-snug">
                    Tu tienda Jumpseller{locationsInfo.subscription_plan ? ` (plan ${locationsInfo.subscription_plan})` : ''} maneja
                    una sola bodega: Jumpseller habilita multi-bodega desde el plan Premium. Por eso en la lista aparece
                    una sola bodega de Jumpseller.
                </p>
            )}

            {loading ? (
                <p className="text-[12px] text-gray-400">Cargando bodegas...</p>
            ) : value.mode === 'single' ? (
                <>
                    <div>
                        <label className={fieldLabel}>Bodega</label>
                        <select
                            value={String(value.single_warehouse_id || 0)}
                            onChange={(e) => set({ single_warehouse_id: Number(e.target.value) })}
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        >
                            <option value="0">-- Selecciona una bodega --</option>
                            {warehouses.map((w) => (
                                <option key={w.id} value={w.id}>{warehouseLabel(w)}</option>
                            ))}
                        </select>
                    </div>

                    {locations.length > 0 && (
                        <div>
                            <label className={fieldLabel}>Bodega de Jumpseller de destino (opcional)</label>
                            <select
                                value={value.default_location_id}
                                onChange={(e) => set({ default_location_id: e.target.value })}
                                className={inputCls}
                                style={{ borderColor: INPUT_BORDER }}
                            >
                                <option value="">-- Sin definir --</option>
                                {locations.map((loc) => (
                                    <option key={loc.id} value={String(loc.id)}>{locationLabel(loc)}</option>
                                ))}
                            </select>
                            <p className="text-[11px] text-gray-400 mt-1">
                                Si la dejas vacia, se usa la bodega de origen de stock de tu tienda.
                            </p>
                        </div>
                    )}
                </>
            ) : (
                <div className="space-y-3">
                    <p className="text-[11px] text-gray-500 dark:text-gray-400 leading-snug">
                        Cada pareja envia el stock de esa bodega de Probability a esa bodega de Jumpseller. Si emparejas
                        dos bodegas de Probability a la misma bodega de Jumpseller, se suma el stock de ambas.
                    </p>

                    {value.mappings.length === 0 && (
                        <p className="text-[12px] text-gray-400">
                            Sin parejas todavia. Agrega una para elegir que bodega envia stock y hacia donde.
                        </p>
                    )}

                    {value.mappings.map((row, index) => (
                        <div key={index} className="grid grid-cols-1 gap-2 sm:grid-cols-[1fr_1fr_auto] sm:items-end">
                            <div>
                                <label className={rowLabel}>Bodega Probability</label>
                                <select
                                    value={String(row.internal_warehouse_id || 0)}
                                    onChange={(e) => updateRow(index, { internal_warehouse_id: Number(e.target.value) })}
                                    className={inputCls}
                                    style={{ borderColor: INPUT_BORDER }}
                                >
                                    <option value="0">-- Selecciona --</option>
                                    {warehouses.map((w) => (
                                        <option key={w.id} value={w.id}>{warehouseLabel(w)}</option>
                                    ))}
                                </select>
                            </div>
                            <div>
                                <label className={rowLabel}>Bodega Jumpseller</label>
                                <select
                                    value={String(row.jumpseller_location_id || 0)}
                                    onChange={(e) => updateRow(index, { jumpseller_location_id: Number(e.target.value) })}
                                    className={inputCls}
                                    style={{ borderColor: INPUT_BORDER }}
                                >
                                    <option value="0">-- Selecciona --</option>
                                    {locations.map((loc) => (
                                        <option key={loc.id} value={String(loc.id)}>{locationLabel(loc)}</option>
                                    ))}
                                </select>
                            </div>
                            <button
                                type="button"
                                onClick={() => removeRow(index)}
                                className="inline-flex items-center justify-center rounded-lg p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20"
                                aria-label="Eliminar pareja"
                            >
                                <TrashIcon className="w-4 h-4" />
                            </button>
                        </div>
                    ))}

                    <button
                        type="button"
                        onClick={addRow}
                        className="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-[12px] font-semibold transition-colors"
                        style={{ border: `1px dashed ${INPUT_BORDER}`, color: ACCENT }}
                    >
                        <PlusIcon className="w-4 h-4" />
                        Agregar pareja
                    </button>

                    {tooManyLocations && (
                        <p className="text-[11px] text-amber-700 dark:text-amber-500 flex items-start gap-1">
                            <ExclamationTriangleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>
                                Jumpseller todavia no permite enviar stock por bodega. Deja un solo destino o usa
                                &quot;Una bodega&quot;.
                            </span>
                        </p>
                    )}
                </div>
            )}

            <p className="text-[11px] text-gray-500 dark:text-gray-400 flex items-start gap-1">
                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                <span>Solo se actualizan productos ya vinculados (creados desde &quot;Sincronizar productos&quot;). El match es por SKU.</span>
            </p>
        </div>
    );
}
