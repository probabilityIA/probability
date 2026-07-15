'use client';

import { useState, useEffect } from 'react';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { BuildingStorefrontIcon, InformationCircleIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';

export interface MeliWarehouseMapping {
    internal_warehouse_id: number;
    ml_store_id: string;
    ml_network_node_id: string;
}

interface NamedWarehouse {
    id: number;
    name: string;
    code?: string;
}

interface MercadoLibreWarehouseMappingSectionProps {
    value: MeliWarehouseMapping[];
    onChange: (v: MeliWarehouseMapping[]) => void;
    businessId: number | null;
}

const GREEN = 'var(--color-primary)';
const CARD_BORDER = '#eceaf3';
const INPUT_BORDER = '#e9e9f0';

const fieldLabel = 'block text-[12px] font-semibold text-gray-700 dark:text-gray-200 mb-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';

export function MercadoLibreWarehouseMappingSection({ value, onChange, businessId }: MercadoLibreWarehouseMappingSectionProps) {
    const [warehouses, setWarehouses] = useState<NamedWarehouse[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
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
    }, [businessId]);

    const updateRow = (index: number, patch: Partial<MeliWarehouseMapping>) => {
        onChange(value.map((row, i) => (i === index ? { ...row, ...patch } : row)));
    };

    const addRow = () => {
        onChange([...value, { internal_warehouse_id: 0, ml_store_id: '', ml_network_node_id: '' }]);
    };

    const removeRow = (index: number) => {
        onChange(value.filter((_, i) => i !== index));
    };

    return (
        <div
            className="rounded-xl p-4 dark:bg-gray-800/60"
            style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
        >
            <div className="flex items-start gap-2 mb-3">
                <BuildingStorefrontIcon className="w-4 h-4 mt-0.5" style={{ color: GREEN }} />
                <div>
                    <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100">Emparejamiento de bodegas (multi-bodega / Full)</h4>
                    <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                        Acopla una bodega de Probability a una bodega/deposito de MercadoLibre. Si defines emparejamientos,
                        el stock se envia por bodega usando la API de stock de MercadoLibre en vez de un total unico.
                    </p>
                </div>
            </div>

            {loading ? (
                <p className="text-[12px] text-gray-400">Cargando bodegas...</p>
            ) : (
                <div className="space-y-3">
                    {value.length === 0 && (
                        <p className="text-[12px] text-gray-400">
                            Sin emparejamientos. Se usara el modo de una bodega / suma configurado arriba.
                        </p>
                    )}

                    {value.map((row, index) => (
                        <div
                            key={index}
                            className="grid grid-cols-1 gap-2 sm:grid-cols-[1fr_1fr_1fr_auto] sm:items-end rounded-lg p-3"
                            style={{ border: `1px solid ${INPUT_BORDER}` }}
                        >
                            <div>
                                <label className={fieldLabel}>Bodega Probability</label>
                                <select
                                    value={String(row.internal_warehouse_id || 0)}
                                    onChange={(e) => updateRow(index, { internal_warehouse_id: Number(e.target.value) })}
                                    className={inputCls}
                                    style={{ borderColor: INPUT_BORDER }}
                                >
                                    <option value="0">-- Selecciona --</option>
                                    {warehouses.map((w) => (
                                        <option key={w.id} value={w.id}>{w.name}{w.code ? ` (${w.code})` : ''}</option>
                                    ))}
                                </select>
                            </div>
                            <div>
                                <label className={fieldLabel}>ML Store ID</label>
                                <input
                                    type="text"
                                    value={row.ml_store_id}
                                    onChange={(e) => updateRow(index, { ml_store_id: e.target.value })}
                                    placeholder="store_id de MercadoLibre"
                                    className={`${inputCls} font-mono`}
                                    style={{ borderColor: INPUT_BORDER }}
                                />
                            </div>
                            <div>
                                <label className={fieldLabel}>ML Network Node ID</label>
                                <input
                                    type="text"
                                    value={row.ml_network_node_id}
                                    onChange={(e) => updateRow(index, { ml_network_node_id: e.target.value })}
                                    placeholder="opcional"
                                    className={`${inputCls} font-mono`}
                                    style={{ borderColor: INPUT_BORDER }}
                                />
                            </div>
                            <button
                                type="button"
                                onClick={() => removeRow(index)}
                                className="inline-flex items-center justify-center rounded-lg p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20"
                                aria-label="Eliminar emparejamiento"
                            >
                                <TrashIcon className="w-4 h-4" />
                            </button>
                        </div>
                    ))}

                    <button
                        type="button"
                        onClick={addRow}
                        className="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-[12px] font-semibold transition-colors"
                        style={{ border: `1px dashed ${INPUT_BORDER}`, color: GREEN }}
                    >
                        <PlusIcon className="w-4 h-4" />
                        Agregar emparejamiento
                    </button>

                    <p className="text-[11px] text-gray-500 dark:text-gray-400 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>El Store ID y Network Node ID los obtienes de tus bodegas en MercadoLibre (stock por deposito / Full).</span>
                    </p>
                </div>
            )}
        </div>
    );
}
