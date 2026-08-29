'use client';

import { useState, useEffect } from 'react';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { ExclamationTriangleIcon, InformationCircleIcon } from '@heroicons/react/24/outline';

export interface TiendanubeInventoryConfig {
    enabled: boolean;
    single_warehouse_id: number;
}

interface NamedWarehouse {
    id: number;
    name: string;
    code?: string;
}

interface TiendanubeInventorySectionProps {
    value: TiendanubeInventoryConfig;
    onChange: (v: TiendanubeInventoryConfig) => void;
    businessId: number | null;
}

const ACCENT = 'var(--color-primary)';
const INPUT_BORDER = '#e9e9f0';

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';

export function TiendanubeInventorySection({ value, onChange, businessId }: TiendanubeInventorySectionProps) {
    const [warehouses, setWarehouses] = useState<NamedWarehouse[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        let cancelled = false;
        const load = async () => {
            setLoading(true);
            try {
                const res: any = await getWarehousesAction(
                    businessId ? ({ business_id: businessId, page_size: 100 } as any) : ({ page_size: 100 } as any)
                );
                if (cancelled) return;
                const list = res?.data || res?.warehouses || [];
                setWarehouses(Array.isArray(list) ? list : []);
            } catch {
                if (!cancelled) setWarehouses([]);
            } finally {
                if (!cancelled) setLoading(false);
            }
        };
        load();
        return () => { cancelled = true; };
    }, [businessId]);

    const set = (patch: Partial<TiendanubeInventoryConfig>) => onChange({ ...value, ...patch });

    return (
        <div className="space-y-3">
            <div
                className="flex items-start justify-between gap-3 rounded-lg bg-white dark:bg-gray-800 px-3 py-2.5"
                style={{ border: `1px solid ${INPUT_BORDER}` }}
            >
                <div>
                    <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100">Sincronizar inventario hacia Tiendanube</h4>
                    <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                        Cuando cambie el stock en Probability, se empuja a las variantes de tu tienda que compartan el SKU.
                    </p>
                </div>
                <button
                    type="button"
                    role="switch"
                    aria-checked={value.enabled}
                    onClick={() => set({ enabled: !value.enabled })}
                    className={`relative inline-flex h-6 w-11 flex-shrink-0 items-center rounded-full transition-colors ${value.enabled ? '' : 'bg-gray-300 dark:bg-gray-600'}`}
                    style={value.enabled ? { backgroundColor: ACCENT } : undefined}
                >
                    <span className={`inline-block h-5 w-5 transform rounded-full bg-white shadow transition-transform ${value.enabled ? 'translate-x-5' : 'translate-x-0.5'}`} />
                </button>
            </div>

            {value.enabled && (
                <div>
                    <label className={fieldLabel}>Bodega de origen</label>
                    {loading ? (
                        <div className="rounded-lg px-3 py-2 text-sm text-gray-500 dark:text-gray-300 bg-white dark:bg-gray-800" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                            Cargando bodegas...
                        </div>
                    ) : (
                        <select
                            value={value.single_warehouse_id ? String(value.single_warehouse_id) : ''}
                            onChange={(e) => set({ single_warehouse_id: Number(e.target.value) })}
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        >
                            <option value="">-- Selecciona una bodega --</option>
                            {warehouses.map((w) => (
                                <option key={w.id} value={String(w.id)}>{w.name}{w.code ? ` (${w.code})` : ''}</option>
                            ))}
                        </select>
                    )}
                    <p className="mt-1 flex items-start gap-1 text-[11px] text-gray-400 dark:text-gray-500">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Tiendanube maneja un solo stock por variante, así que se envía el de esta bodega.</span>
                    </p>

                    {!loading && warehouses.length === 0 && (
                        <p className="mt-2 flex items-start gap-1 text-[11px] text-amber-600 dark:text-amber-400">
                            <ExclamationTriangleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>No hay bodegas configuradas para este negocio: crea una antes de activar la sincronización.</span>
                        </p>
                    )}

                    {!loading && warehouses.length > 0 && !value.single_warehouse_id && (
                        <p className="mt-2 flex items-start gap-1 text-[11px] text-amber-600 dark:text-amber-400">
                            <ExclamationTriangleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>Sin bodega seleccionada no se envía stock a Tiendanube.</span>
                        </p>
                    )}
                </div>
            )}
        </div>
    );
}
