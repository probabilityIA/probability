'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import {
    CubeIcon,
    InformationCircleIcon,
    PlusIcon,
    TrashIcon,
    ArrowDownTrayIcon,
    ArrowPathIcon,
} from '@heroicons/react/24/outline';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { useSSE } from '@/shared/hooks/use-sse';
import { listSiigoWarehousesAction } from '../../infra/actions';
import {
    GREEN,
    GREEN_SOFT,
    GREEN_BORDER,
    INPUT_BORDER,
    fieldLabel,
    fieldHint,
    inputCls,
    SectionCard,
    ToggleRow,
} from './SiigoFormKit';

export interface WarehousePair {
    velocity_warehouse_id: number;
    siigo_warehouse_id: number;
}

export interface InventorySyncConfig {
    enabled: boolean;
    mode: 'single' | 'mapped';
    single_warehouse_id: number;
    mappings: WarehousePair[];
    product_sync_enabled: boolean;
}

interface NamedWarehouse {
    id: number;
    name: string;
    code?: string;
}

interface SiigoInventorySectionProps {
    value: InventorySyncConfig;
    onChange: (v: InventorySyncConfig) => void;
    businessId: number | null;
    integrationId?: number;
    onSyncNow?: () => void;
    canSyncNow?: boolean;
}

export function SiigoInventorySection({ value, onChange, businessId, integrationId, onSyncNow, canSyncNow }: SiigoInventorySectionProps) {
    const [warehouses, setWarehouses] = useState<NamedWarehouse[]>([]);
    const [loading, setLoading] = useState(false);

    const [siigoWarehouses, setSiigoWarehouses] = useState<NamedWarehouse[]>([]);
    const [loadingSiigo, setLoadingSiigo] = useState(false);
    const siigoCorrelationRef = useRef<string | null>(null);

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

    const handleSSE = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data);
            const eventType = parsed.type || parsed.metadata?.event_type;
            const data = parsed.data;
            if (eventType !== 'invoice.siigo_warehouses_ready' || !data) return;
            if (siigoCorrelationRef.current && data.correlation_id !== siigoCorrelationRef.current) return;
            const results = Array.isArray(data.results) ? data.results : [];
            setSiigoWarehouses(results.map((w: any) => ({ id: Number(w.id), name: w.name })));
            setLoadingSiigo(false);
        } catch {
            return;
        }
    }, []);

    const { isConnected } = useSSE({
        businessId: businessId ?? 0,
        eventTypes: ['invoice.siigo_warehouses_ready'],
        onMessage: handleSSE,
        enabled: value.enabled && value.mode === 'mapped' && !!integrationId,
    });

    const loadSiigoWarehouses = useCallback(async () => {
        if (!integrationId) return;
        setLoadingSiigo(true);
        const res: any = await listSiigoWarehousesAction(integrationId, businessId ?? undefined);
        if (!res?.success || !res?.correlation_id) {
            setLoadingSiigo(false);
            return;
        }
        siigoCorrelationRef.current = res.correlation_id;
    }, [integrationId, businessId]);

    useEffect(() => {
        if (value.enabled && value.mode === 'mapped' && integrationId && isConnected && siigoWarehouses.length === 0 && !loadingSiigo) {
            loadSiigoWarehouses();
        }
    }, [value.enabled, value.mode, integrationId, isConnected]);

    const set = (patch: Partial<InventorySyncConfig>) => onChange({ ...value, ...patch });

    const updateMapping = (idx: number, patch: Partial<WarehousePair>) => {
        const mappings = value.mappings.map((m, i) => (i === idx ? { ...m, ...patch } : m));
        set({ mappings });
    };

    const addMapping = () => set({ mappings: [...value.mappings, { velocity_warehouse_id: 0, siigo_warehouse_id: 0 }] });
    const removeMapping = (idx: number) => set({ mappings: value.mappings.filter((_, i) => i !== idx) });

    const velocityOptions = (
        <>
            <option value="0">-- Selecciona una bodega --</option>
            {warehouses.map((w) => (
                <option key={w.id} value={w.id}>{w.name}{w.code ? ` (${w.code})` : ''}</option>
            ))}
        </>
    );

    const siigoOptions = (
        <>
            <option value="0">-- Bodega Siigo --</option>
            {siigoWarehouses.map((w) => (
                <option key={w.id} value={w.id}>{w.name} (#{w.id})</option>
            ))}
        </>
    );

    return (
        <SectionCard icon={<CubeIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Sincronizacion de Inventario">
            <div className="rounded-lg bg-white dark:bg-gray-800" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                <ToggleRow
                    icon={<CubeIcon className="w-4 h-4" style={{ color: GREEN }} />}
                    title="Sincronizar inventario con Siigo"
                    subtitle="El stock de Probability se actualiza con el de Siigo. Solo lectura, en una sola direccion: Siigo -> Probability."
                    checked={value.enabled}
                    onToggle={() => set({ enabled: !value.enabled })}
                />
                <div className="border-t border-gray-100 dark:border-gray-700" />
                <ToggleRow
                    icon={<CubeIcon className="w-4 h-4" style={{ color: GREEN }} />}
                    title="Sincronizar productos de Siigo a Probability"
                    subtitle="Cuando se crea o actualiza un producto en Siigo, se crea o actualiza tambien en Probability (aplicando la logica de bodegas configurada)."
                    checked={value.product_sync_enabled}
                    onToggle={() => set({ product_sync_enabled: !value.product_sync_enabled })}
                />
            </div>

            {value.enabled && (
                <div className="mt-3 space-y-3">
                    <p className={`${fieldHint} !text-[12px]`}>
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Siigo es la fuente de verdad del stock. Probability nunca escribe inventario en Siigo.</span>
                    </p>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                        <button
                            type="button"
                            onClick={() => set({ mode: 'single' })}
                            className="text-left rounded-lg px-3 py-2.5"
                            style={{
                                border: `1px solid ${value.mode === 'single' ? GREEN_BORDER : INPUT_BORDER}`,
                                backgroundColor: value.mode === 'single' ? GREEN_SOFT : 'transparent',
                            }}
                        >
                            <p className="text-[13px] font-semibold text-gray-900 dark:text-white">Una sola bodega</p>
                            <p className="text-[11px] text-gray-500 dark:text-gray-400">Todo el stock de Siigo entra a una sola bodega.</p>
                        </button>
                        <button
                            type="button"
                            onClick={() => set({ mode: 'mapped' })}
                            className="text-left rounded-lg px-3 py-2.5"
                            style={{
                                border: `1px solid ${value.mode === 'mapped' ? GREEN_BORDER : INPUT_BORDER}`,
                                backgroundColor: value.mode === 'mapped' ? GREEN_SOFT : 'transparent',
                            }}
                        >
                            <p className="text-[13px] font-semibold text-gray-900 dark:text-white">Mapear bodegas</p>
                            <p className="text-[11px] text-gray-500 dark:text-gray-400">Relaciona cada bodega con su bodega de Siigo.</p>
                        </button>
                    </div>

                    {value.mode === 'single' ? (
                        <div>
                            <label className={fieldLabel}>Bodega destino</label>
                            <select
                                value={String(value.single_warehouse_id || 0)}
                                onChange={(e) => set({ single_warehouse_id: Number(e.target.value) })}
                                className={inputCls}
                                style={{ borderColor: INPUT_BORDER }}
                                disabled={loading}
                            >
                                {velocityOptions}
                            </select>
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>El inventario total de cada producto en Siigo se cargara en esta bodega.</span>
                            </p>
                        </div>
                    ) : (
                        <div className="space-y-2">
                            <div className="flex items-center justify-between">
                                <label className={`${fieldLabel} mb-0`}>Parejas de bodegas</label>
                                {integrationId && (
                                    <button
                                        type="button"
                                        onClick={loadSiigoWarehouses}
                                        disabled={loadingSiigo}
                                        className="inline-flex items-center gap-1 text-[11px] font-semibold disabled:opacity-50"
                                        style={{ color: GREEN }}
                                    >
                                        <ArrowPathIcon className={`w-3.5 h-3.5 ${loadingSiigo ? 'animate-spin' : ''}`} />
                                        {loadingSiigo ? 'Cargando bodegas Siigo...' : 'Recargar bodegas Siigo'}
                                    </button>
                                )}
                            </div>

                            {!integrationId && (
                                <p className="text-[11px] text-amber-600">Guarda la integracion primero para poder elegir bodegas de Siigo.</p>
                            )}

                            <div className="grid grid-cols-[1fr_auto] gap-x-2 items-center px-0.5">
                                <div className="grid grid-cols-2 gap-2">
                                    <span className="text-[11px] font-semibold text-gray-400">Bodega</span>
                                    <span className="text-[11px] font-semibold text-gray-400">Siigo</span>
                                </div>
                                <span />
                            </div>

                            {value.mappings.length === 0 && (
                                <p className="text-[11px] text-gray-400">No hay parejas. Agrega al menos una.</p>
                            )}

                            {value.mappings.map((m, idx) => (
                                <div key={idx} className="grid grid-cols-[1fr_auto] gap-x-2 items-center">
                                    <div className="grid grid-cols-2 gap-2 min-w-0">
                                        <select
                                            value={String(m.velocity_warehouse_id || 0)}
                                            onChange={(e) => updateMapping(idx, { velocity_warehouse_id: Number(e.target.value) })}
                                            className={`${inputCls} min-w-0`}
                                            style={{ borderColor: INPUT_BORDER }}
                                            disabled={loading}
                                        >
                                            {velocityOptions}
                                        </select>
                                        <select
                                            value={String(m.siigo_warehouse_id || 0)}
                                            onChange={(e) => updateMapping(idx, { siigo_warehouse_id: Number(e.target.value) })}
                                            className={`${inputCls} min-w-0`}
                                            style={{ borderColor: INPUT_BORDER }}
                                            disabled={!integrationId || loadingSiigo || siigoWarehouses.length === 0}
                                        >
                                            {siigoOptions}
                                        </select>
                                    </div>
                                    <button
                                        type="button"
                                        onClick={() => removeMapping(idx)}
                                        className="p-2 rounded-lg text-gray-400 hover:text-red-500 shrink-0"
                                    >
                                        <TrashIcon className="w-4 h-4" />
                                    </button>
                                </div>
                            ))}

                            <button
                                type="button"
                                onClick={addMapping}
                                className="inline-flex items-center gap-1.5 text-[13px] font-semibold rounded-lg px-3 py-1.5"
                                style={{ border: `1px solid ${GREEN_BORDER}`, color: GREEN }}
                            >
                                <PlusIcon className="w-4 h-4" /> Agregar pareja
                            </button>
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Cada pareja: una bodega propia y su bodega de Siigo. El stock de esa bodega Siigo se refleja en la bodega elegida.</span>
                            </p>
                        </div>
                    )}

                    {onSyncNow && (
                        <button
                            type="button"
                            onClick={onSyncNow}
                            disabled={!canSyncNow}
                            className="w-full flex items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-[13px] font-semibold text-white disabled:opacity-50"
                            style={{ backgroundColor: GREEN }}
                        >
                            <ArrowDownTrayIcon className="w-4 h-4" /> Sincronizar inventario ahora
                        </button>
                    )}
                    {onSyncNow && !canSyncNow && (
                        <p className="text-[11px] text-gray-400 text-center">Guarda la integracion para poder sincronizar.</p>
                    )}
                </div>
            )}
        </SectionCard>
    );
}
