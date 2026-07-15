'use client';

import { useState, useEffect } from 'react';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { CubeIcon, InformationCircleIcon, ArrowPathIcon, ArrowDownTrayIcon } from '@heroicons/react/24/outline';

export interface WooInventoryConfig {
    enabled: boolean;
    mode: 'single' | 'sum';
    single_warehouse_id: number;
    warehouse_ids: number[];
}

interface NamedWarehouse {
    id: number;
    name: string;
    code?: string;
}

interface WooCommerceInventorySectionProps {
    value: WooInventoryConfig;
    onChange: (v: WooInventoryConfig) => void;
    businessId: number | null;
    integrationId?: number;
    onSyncNow?: () => void;
    canSyncNow?: boolean;
}

const ACCENT = 'var(--color-primary)';
const ACCENT_SOFT = 'color-mix(in srgb, var(--color-primary) 10%, white)';
const CARD_BORDER = '#eceaf3';
const INPUT_BORDER = '#e9e9f0';

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';

export function WooCommerceInventorySection({ value, onChange, businessId, integrationId, onSyncNow, canSyncNow }: WooCommerceInventorySectionProps) {
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

    const set = (patch: Partial<WooInventoryConfig>) => onChange({ ...value, ...patch });

    const toggleWarehouse = (id: number) => {
        const exists = value.warehouse_ids.includes(id);
        set({ warehouse_ids: exists ? value.warehouse_ids.filter((w) => w !== id) : [...value.warehouse_ids, id] });
    };

    return (
        <div
            className="rounded-xl p-4 dark:bg-gray-800/60"
            style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
        >
            <div className="flex items-center justify-between gap-3">
                <div className="flex items-start gap-2">
                    <CubeIcon className="w-4 h-4 mt-0.5" style={{ color: ACCENT }} />
                    <div>
                        <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100">Sincronizacion de inventario</h4>
                        <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                            Empuja el stock de Probability a tus productos de WooCommerce (stock_quantity).
                        </p>
                    </div>
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
                <div className="mt-4 space-y-4">
                    <div>
                        <label className={fieldLabel}>Que stock enviar a WooCommerce</label>
                        <div className="grid grid-cols-2 gap-2">
                            <button
                                type="button"
                                onClick={() => set({ mode: 'single' })}
                                className="rounded-lg border px-3 py-2 text-left text-[12px] transition-colors"
                                style={value.mode === 'single'
                                    ? { borderColor: ACCENT, backgroundColor: ACCENT_SOFT }
                                    : { borderColor: INPUT_BORDER }}
                            >
                                <span className="font-semibold text-gray-900 dark:text-white block">Una bodega</span>
                                <span className="text-[11px] text-gray-500 dark:text-gray-400">El stock de una bodega</span>
                            </button>
                            <button
                                type="button"
                                onClick={() => set({ mode: 'sum' })}
                                className="rounded-lg border px-3 py-2 text-left text-[12px] transition-colors"
                                style={value.mode === 'sum'
                                    ? { borderColor: ACCENT, backgroundColor: ACCENT_SOFT }
                                    : { borderColor: INPUT_BORDER }}
                            >
                                <span className="font-semibold text-gray-900 dark:text-white block">Sumar bodegas</span>
                                <span className="text-[11px] text-gray-500 dark:text-gray-400">Suma el stock de varias</span>
                            </button>
                        </div>
                    </div>

                    {loading ? (
                        <p className="text-[12px] text-gray-400">Cargando bodegas...</p>
                    ) : value.mode === 'single' ? (
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
                                    <option key={w.id} value={w.id}>{w.name}{w.code ? ` (${w.code})` : ''}</option>
                                ))}
                            </select>
                        </div>
                    ) : (
                        <div>
                            <label className={fieldLabel}>Bodegas a sumar</label>
                            <p className="text-[11px] text-gray-400 mb-2">Si no seleccionas ninguna, se suma el stock de todas.</p>
                            <div className="space-y-1.5 max-h-40 overflow-y-auto rounded-lg p-2" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                                {warehouses.length === 0 && <p className="text-[12px] text-gray-400">No hay bodegas.</p>}
                                {warehouses.map((w) => (
                                    <label key={w.id} className="flex items-center gap-2 text-[12px] text-gray-700 dark:text-gray-200 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            checked={value.warehouse_ids.includes(w.id)}
                                            onChange={() => toggleWarehouse(w.id)}
                                            className="h-4 w-4 rounded border-gray-300 dark:border-gray-600 accent-[var(--color-primary)]"
                                        />
                                        <span>{w.name}{w.code ? ` (${w.code})` : ''}</span>
                                    </label>
                                ))}
                            </div>
                        </div>
                    )}

                    <p className="text-[11px] text-gray-500 dark:text-gray-400 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Solo se actualizan productos ya vinculados (creados desde "Sincronizar productos"). El match es por SKU.</span>
                    </p>

                    {onSyncNow && (
                        <button
                            type="button"
                            onClick={onSyncNow}
                            disabled={!canSyncNow}
                            className="w-full inline-flex items-center justify-center gap-2 rounded-lg py-2 text-[13px] font-semibold text-white transition-colors disabled:opacity-50"
                            style={{ backgroundColor: ACCENT }}
                        >
                            <ArrowDownTrayIcon className="w-4 h-4" />
                            Sincronizar inventario ahora
                        </button>
                    )}
                    {onSyncNow && !canSyncNow && (
                        <p className="text-[11px] text-gray-400 flex items-center gap-1">
                            <ArrowPathIcon className="w-3.5 h-3.5" />
                            Guarda la integracion para poder sincronizar.
                        </p>
                    )}
                </div>
            )}
        </div>
    );
}
