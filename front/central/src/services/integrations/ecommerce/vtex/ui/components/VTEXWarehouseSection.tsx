'use client';

import { useState, useEffect } from 'react';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { InformationCircleIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';

export interface VTEXWarehouseMapping {
    internal_warehouse_id: number;
    vtex_warehouse_id: string;
}

export interface VTEXWarehouse {
    id: string;
    name: string;
    is_active?: boolean;
}

export interface VTEXWarehousesInfo {
    warehouses: VTEXWarehouse[];
}

interface NamedWarehouse {
    id: number;
    name: string;
    code?: string;
}

interface VTEXWarehouseSectionProps {
    enabled: boolean;
    value: VTEXWarehouseMapping[];
    onChange: (v: VTEXWarehouseMapping[]) => void;
    businessId: number | null;
    warehousesInfo: VTEXWarehousesInfo | null;
}

const ACCENT = 'var(--color-primary)';
const INPUT_BORDER = '#e9e9f0';

const rowLabel = 'block text-[12px] font-semibold text-gray-700 dark:text-gray-200 mb-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';

const warehouseLabel = (w: NamedWarehouse) => `${w.name}${w.code ? ` (${w.code})` : ''}`;

const vtexWarehouseLabel = (w: VTEXWarehouse) => (w.is_active === false ? `${w.name} (inactiva)` : w.name);

export function VTEXWarehouseSection({ enabled, value, onChange, businessId, warehousesInfo }: VTEXWarehouseSectionProps) {
    const [warehouses, setWarehouses] = useState<NamedWarehouse[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (!enabled) return;
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
    }, [enabled, businessId]);

    if (!enabled) return null;

    const vtexWarehouses = warehousesInfo?.warehouses ?? [];

    const updateRow = (index: number, patch: Partial<VTEXWarehouseMapping>) => {
        onChange(value.map((row, i) => (i === index ? { ...row, ...patch } : row)));
    };
    const addRow = () => onChange([...value, { internal_warehouse_id: 0, vtex_warehouse_id: '' }]);
    const removeRow = (index: number) => onChange(value.filter((_, i) => i !== index));

    return (
        <div className="mt-4 space-y-3">
            <p className="text-[11px] text-gray-500 dark:text-gray-400 leading-snug">
                Cada pareja envia el stock de esa bodega de Probability a esa bodega de VTEX. Si emparejas dos bodegas de
                Probability a la misma bodega de VTEX, se suma el stock de ambas.
            </p>

            {loading ? (
                <p className="text-[12px] text-gray-400">Cargando bodegas...</p>
            ) : (
                <>
                    {vtexWarehouses.length === 0 && (
                        <p className="text-[11px] text-amber-700 dark:text-amber-500 leading-snug">
                            Guarda la integracion con credenciales validas para poder listar las bodegas de VTEX.
                        </p>
                    )}

                    {value.length === 0 && (
                        <p className="text-[12px] text-gray-400">
                            Sin parejas todavia. Agrega una para elegir que bodega envia stock y hacia donde.
                        </p>
                    )}

                    {value.map((row, index) => (
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
                                <label className={rowLabel}>Bodega VTEX</label>
                                <select
                                    value={row.vtex_warehouse_id}
                                    onChange={(e) => updateRow(index, { vtex_warehouse_id: e.target.value })}
                                    className={inputCls}
                                    style={{ borderColor: INPUT_BORDER }}
                                >
                                    <option value="">-- Selecciona --</option>
                                    {vtexWarehouses.map((w) => (
                                        <option key={w.id} value={w.id}>{vtexWarehouseLabel(w)}</option>
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
                </>
            )}

            <p className="text-[11px] text-gray-500 dark:text-gray-400 flex items-start gap-1">
                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                <span>Solo se actualizan productos ya vinculados (asociados desde &quot;Sincronizar productos&quot;). El match es por SKU contra el RefId de VTEX.</span>
            </p>
        </div>
    );
}
