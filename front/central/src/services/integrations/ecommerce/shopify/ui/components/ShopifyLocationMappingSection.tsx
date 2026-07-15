'use client';

import { useState, useEffect } from 'react';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { BuildingStorefrontIcon, InformationCircleIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';

export interface ShopifyLocationMapping {
    internal_warehouse_id: number;
    shopify_location_id: string;
}

interface NamedWarehouse {
    id: number;
    name: string;
    code?: string;
}

interface ShopifyLocationMappingSectionProps {
    mappings: ShopifyLocationMapping[];
    onChangeMappings: (v: ShopifyLocationMapping[]) => void;
    defaultLocationId: string;
    onChangeDefaultLocation: (v: string) => void;
    businessId: number | null;
}

const ACCENT = 'var(--color-primary)';
const CARD_BORDER = '#eceaf3';
const INPUT_BORDER = '#e9e9f0';

const fieldLabel = 'block text-[12px] font-semibold text-gray-700 dark:text-gray-200 mb-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';

export function ShopifyLocationMappingSection({ mappings, onChangeMappings, defaultLocationId, onChangeDefaultLocation, businessId }: ShopifyLocationMappingSectionProps) {
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

    const updateRow = (index: number, patch: Partial<ShopifyLocationMapping>) => {
        onChangeMappings(mappings.map((row, i) => (i === index ? { ...row, ...patch } : row)));
    };
    const addRow = () => onChangeMappings([...mappings, { internal_warehouse_id: 0, shopify_location_id: '' }]);
    const removeRow = (index: number) => onChangeMappings(mappings.filter((_, i) => i !== index));

    return (
        <div
            className="rounded-xl p-4 dark:bg-gray-800/60"
            style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
        >
            <div className="flex items-start gap-2 mb-3">
                <BuildingStorefrontIcon className="w-4 h-4 mt-0.5" style={{ color: ACCENT }} />
                <div>
                    <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100">Emparejamiento de locations (multi-ubicacion)</h4>
                    <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                        Acopla una bodega de Probability a una location de Shopify. Si defines emparejamientos,
                        el stock se envia por location; si no, se usa la location por defecto.
                    </p>
                </div>
            </div>

            <div className="mb-3">
                <label className={fieldLabel}>Location por defecto (opcional)</label>
                <input
                    type="text"
                    value={defaultLocationId}
                    onChange={(e) => onChangeDefaultLocation(e.target.value)}
                    placeholder="ID de la location por defecto de Shopify"
                    className={`${inputCls} font-mono`}
                    style={{ borderColor: INPUT_BORDER }}
                />
                <p className="text-[11px] text-gray-400 mt-1">Si lo dejas vacio, se usa la primera location de la tienda.</p>
            </div>

            {loading ? (
                <p className="text-[12px] text-gray-400">Cargando bodegas...</p>
            ) : (
                <div className="space-y-3">
                    {mappings.length === 0 && (
                        <p className="text-[12px] text-gray-400">
                            Sin emparejamientos. Se usara la location por defecto con el modo una bodega / suma.
                        </p>
                    )}

                    {mappings.map((row, index) => (
                        <div
                            key={index}
                            className="grid grid-cols-1 gap-2 sm:grid-cols-[1fr_1fr_auto] sm:items-end rounded-lg p-3"
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
                                <label className={fieldLabel}>Shopify Location ID</label>
                                <input
                                    type="text"
                                    value={row.shopify_location_id}
                                    onChange={(e) => updateRow(index, { shopify_location_id: e.target.value })}
                                    placeholder="location_id de Shopify"
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
                        style={{ border: `1px dashed ${INPUT_BORDER}`, color: ACCENT }}
                    >
                        <PlusIcon className="w-4 h-4" />
                        Agregar emparejamiento
                    </button>

                    <p className="text-[11px] text-gray-500 dark:text-gray-400 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>El Location ID lo obtienes en Shopify Admin &gt; Settings &gt; Locations.</span>
                    </p>
                </div>
            )}
        </div>
    );
}
