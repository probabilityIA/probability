'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Spinner, Alert, Button } from '@/shared/ui';
import { ArrowLeftIcon, PencilSquareIcon, ArrowPathIcon } from '@heroicons/react/24/outline';
import { useSSE } from '@/shared/hooks/use-sse';
import { TokenStorage } from '@/shared/utils/token-storage';
import { getLayoutAction, getOccupancyAction } from '../../infra/actions/hierarchy';
import { getWarehouseByIdAction } from '../../infra/actions';
import { useWarehouseTree } from '../hooks/useWarehouseTree';
import { LayoutNode, WarehouseLayout } from '../../domain/hierarchy-types';

interface Props {
    warehouseId: number;
    businessId?: number;
}

interface OccCell { quantity: number; capacity: number | null }

function occColor(qty: number, cap: number | null): { fill: string; stroke: string } {
    if (qty <= 0) return { fill: '#e5e7eb', stroke: '#9ca3af' };
    if (cap && cap > 0) {
        const ratio = qty / cap;
        if (ratio >= 0.85) return { fill: '#f87171', stroke: '#dc2626' };
        if (ratio >= 0.5) return { fill: '#fbbf24', stroke: '#d97706' };
        return { fill: '#34d399', stroke: '#059669' };
    }
    return { fill: '#34d399', stroke: '#059669' };
}

function findRack(tree: any, rackId: number): any {
    if (!tree) return null;
    for (const z of tree.zones || []) {
        for (const a of z.aisles || []) {
            for (const r of a.racks || []) {
                if (r.id === rackId) return r;
            }
        }
    }
    return null;
}

export default function WarehouseOperationView({ warehouseId, businessId }: Props) {
    const router = useRouter();
    const [layout, setLayout] = useState<WarehouseLayout | null>(null);
    const [occupancy, setOccupancy] = useState<Record<number, OccCell>>({});
    const [warehouseName, setWarehouseName] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedRackId, setSelectedRackId] = useState<number | null>(null);
    const [zoom, setZoom] = useState(1);
    const [liveAt, setLiveAt] = useState<string | null>(null);
    const [pulseLoc, setPulseLoc] = useState<number | null>(null);
    const pulseTimer = useRef<any>(null);

    const { tree } = useWarehouseTree({ warehouseId, businessId });

    const sseBusinessId = businessId ?? TokenStorage.getBusinessesData()?.[0]?.id ?? 0;

    const handleSSE = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data);
            const type = parsed.type || parsed.metadata?.event_type;
            if (type !== 'inventory.location_changed') return;
            const data = parsed.data || {};
            const whId = Number(data.warehouse_id ?? parsed.metadata?.warehouse_id);
            if (whId !== warehouseId) return;
            const locationId = Number(data.location_id);
            const newQty = Number(data.new_quantity);
            if (!locationId) return;
            setOccupancy((prev) => ({ ...prev, [locationId]: { quantity: newQty, capacity: prev[locationId]?.capacity ?? null } }));
            setLiveAt(new Date().toLocaleTimeString());
            setPulseLoc(locationId);
            if (pulseTimer.current) clearTimeout(pulseTimer.current);
            pulseTimer.current = setTimeout(() => setPulseLoc(null), 1500);
        } catch {
            return;
        }
    }, [warehouseId]);

    useSSE({
        businessId: sseBusinessId,
        eventTypes: ['inventory.location_changed'],
        onMessage: handleSSE,
        enabled: sseBusinessId > 0,
    });

    const loadAll = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const [lay, occ, wh] = await Promise.all([
                getLayoutAction(warehouseId, businessId),
                getOccupancyAction(warehouseId, businessId),
                getWarehouseByIdAction(warehouseId, businessId).catch(() => null),
            ]);
            setLayout(lay);
            const map: Record<number, OccCell> = {};
            for (const l of occ.locations || []) map[l.location_id] = { quantity: l.quantity, capacity: l.capacity };
            setOccupancy(map);
            if (wh) setWarehouseName((wh as any).name || '');
        } catch (e: any) {
            setError(e?.message || 'Error al cargar la vista operativa');
        } finally {
            setLoading(false);
        }
    }, [warehouseId, businessId]);

    useEffect(() => { loadAll(); }, [loadAll]);

    const rackAgg = useCallback((rackId: number) => {
        const rack = findRack(tree, rackId);
        if (!rack) return { qty: 0, cap: 0, hasCap: false };
        let qty = 0, cap = 0, hasCap = false;
        for (const lv of rack.levels || []) {
            for (const p of lv.positions || []) {
                const o = occupancy[p.id];
                if (!o) continue;
                qty += o.quantity;
                if (o.capacity && o.capacity > 0) { cap += o.capacity; hasCap = true; }
            }
        }
        return { qty, cap, hasCap };
    }, [tree, occupancy]);

    const totals = useMemo(() => {
        let units = 0, occupied = 0, total = 0;
        for (const [, o] of Object.entries(occupancy)) {
            total += 1;
            units += o.quantity;
            if (o.quantity > 0) occupied += 1;
        }
        return { units, occupied, total };
    }, [occupancy]);

    const selectedRack = selectedRackId != null ? findRack(tree, selectedRackId) : null;

    if (loading) {
        return <div className="min-h-[60vh] flex items-center justify-center"><Spinner size="lg" /></div>;
    }
    if (error || !layout) {
        return <div className="p-6"><Alert type="error">{error || 'No se pudo cargar el plano'}</Alert></div>;
    }

    const W = layout.canvas_width;
    const H = layout.canvas_height;
    const grid = layout.grid_size || 20;
    const pct = totals.total > 0 ? Math.round((totals.occupied / totals.total) * 100) : 0;

    const rackColor = (n: LayoutNode): { fill: string; stroke: string; op: number } => {
        const agg = rackAgg(n.ref_id);
        if (agg.qty <= 0 && !agg.hasCap) return { fill: '#cbd5e1', stroke: '#94a3b8', op: 0.55 };
        const c = occColor(agg.qty, agg.hasCap ? agg.cap : null);
        return { fill: c.fill, stroke: c.stroke, op: 0.8 };
    };

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 px-4 sm:px-6 lg:px-8 py-4 sm:py-6">
            <div className="space-y-4">
                <div className="flex items-center justify-between flex-wrap gap-2">
                    <div className="flex items-center gap-3">
                        <button onClick={() => router.push(`/warehouses/${warehouseId}`)} className="p-2 rounded-md hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-500 dark:text-gray-300">
                            <ArrowLeftIcon className="w-5 h-5" />
                        </button>
                        <div>
                            <h1 className="text-xl font-semibold text-gray-900 dark:text-white">{warehouseName || `Bodega #${warehouseId}`}</h1>
                            <p className="text-sm text-gray-500 dark:text-gray-400">Vista operativa - ocupacion en vivo</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-2">
                        <span className="flex items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400">
                            <span className={`w-2 h-2 rounded-full ${sseBusinessId > 0 ? 'bg-emerald-500 animate-pulse' : 'bg-gray-400'}`} />
                            {liveAt ? `En vivo - ${liveAt}` : 'En vivo'}
                        </span>
                        <button className="px-2 py-1 rounded border border-gray-300 dark:border-gray-600 text-sm flex items-center gap-1" onClick={loadAll}>
                            <ArrowPathIcon className="w-4 h-4" /> Actualizar
                        </button>
                        <Button variant="outline" onClick={() => router.push(`/warehouses/${warehouseId}`)}>
                            <PencilSquareIcon className="w-4 h-4 mr-2" /> Editar plano
                        </Button>
                    </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                    {[
                        { label: 'Unidades', value: totals.units, color: 'text-indigo-600' },
                        { label: 'Ubicaciones ocupadas', value: `${totals.occupied}/${totals.total}`, color: 'text-emerald-600' },
                        { label: 'Ocupacion', value: `${pct}%`, color: 'text-amber-600' },
                        { label: 'Escala', value: `1 m = ${Math.round((layout.scale && layout.scale > 0 ? layout.scale : 40))} px`, color: 'text-gray-500' },
                    ].map((s) => (
                        <div key={s.label} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
                            <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">{s.label}</p>
                            <p className={`text-2xl font-semibold ${s.color}`}>{s.value}</p>
                        </div>
                    ))}
                </div>

                <div className="flex items-center gap-3 flex-wrap text-xs text-gray-600 dark:text-gray-300">
                    <span className="font-medium">Leyenda:</span>
                    {[
                        ['#e5e7eb', 'Vacio'],
                        ['#34d399', 'Bajo (<50%)'],
                        ['#fbbf24', 'Medio (50-85%)'],
                        ['#f87171', 'Lleno (>85%)'],
                    ].map(([c, l]) => (
                        <span key={l} className="flex items-center gap-1">
                            <span className="w-3 h-3 rounded-sm inline-block" style={{ background: c }} /> {l}
                        </span>
                    ))}
                    <span className="ml-auto flex items-center gap-1">
                        <button className="px-2 py-0.5 rounded border border-gray-300 dark:border-gray-600" onClick={() => setZoom((z) => Math.max(0.25, +(z - 0.1).toFixed(2)))}>-</button>
                        <span className="w-10 text-center">{Math.round(zoom * 100)}%</span>
                        <button className="px-2 py-0.5 rounded border border-gray-300 dark:border-gray-600" onClick={() => setZoom((z) => Math.min(2, +(z + 0.1).toFixed(2)))}>+</button>
                    </span>
                </div>

                <div className="flex gap-3">
                    <div className="flex-1 min-w-0 overflow-auto bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg" style={{ maxHeight: 680 }}>
                        <svg viewBox={`0 0 ${W} ${H}`} preserveAspectRatio="xMinYMin meet" onClick={() => setSelectedRackId(null)} style={{ width: `${zoom * 100}%`, height: 'auto', background: 'white', display: 'block' }}>
                            <defs>
                                <pattern id="ogrid" width={grid} height={grid} patternUnits="userSpaceOnUse">
                                    <path d={`M ${grid} 0 L 0 0 0 ${grid}`} fill="none" stroke="#eef2f7" strokeWidth="1" />
                                </pattern>
                                <linearGradient id="oRackRelief" x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="0%" stopColor="#ffffff" stopOpacity="0.4" />
                                    <stop offset="45%" stopColor="#ffffff" stopOpacity="0" />
                                    <stop offset="100%" stopColor="#000000" stopOpacity="0.25" />
                                </linearGradient>
                                <linearGradient id="oAisleSunken" x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="0%" stopColor="#000000" stopOpacity="0.22" />
                                    <stop offset="20%" stopColor="#000000" stopOpacity="0" />
                                    <stop offset="80%" stopColor="#000000" stopOpacity="0" />
                                    <stop offset="100%" stopColor="#000000" stopOpacity="0.22" />
                                </linearGradient>
                                <filter id="oRackShadow" x="-20%" y="-20%" width="140%" height="160%">
                                    <feDropShadow dx="0" dy="2.5" stdDeviation="2" floodColor="#0f172a" floodOpacity="0.3" />
                                </filter>
                            </defs>
                            <rect x={0} y={0} width={W} height={H} fill="url(#ogrid)" />

                            {layout.nodes.map((n) => {
                                if (n.ref_type === 'label') {
                                    return (
                                        <text key={n.node_id} x={n.x} y={n.y + 14} fontSize={12} fill={n.color}
                                            transform={`rotate(${n.rotation} ${n.x} ${n.y})`} style={{ userSelect: 'none' }}>{n.label}</text>
                                    );
                                }
                                const isRack = n.ref_type === 'rack' && n.ref_id > 0;
                                let fill = n.color;
                                let stroke = n.color;
                                let op = n.ref_type === 'zone' ? 0.18 : 0.7;
                                if (isRack) {
                                    const rc = rackColor(n);
                                    fill = rc.fill; stroke = rc.stroke; op = rc.op;
                                }
                                const isSel = isRack && n.ref_id === selectedRackId;
                                const isAisle = n.ref_type === 'aisle';
                                return (
                                    <g key={n.node_id} transform={`translate(${n.x} ${n.y}) rotate(${n.rotation} ${n.width / 2} ${n.height / 2})`}
                                        filter={isRack ? 'url(#oRackShadow)' : undefined}>
                                        <rect width={n.width} height={n.height} rx={4} fill={fill} fillOpacity={isAisle ? 0.55 : op}
                                            stroke={stroke} strokeWidth={isSel ? 3 : 1.2}
                                            style={{ cursor: isRack ? 'pointer' : 'default' }}
                                            onClick={isRack ? (e) => { e.stopPropagation(); setSelectedRackId(n.ref_id); } : undefined} />
                                        {(isRack || n.ref_type === 'dock') && (
                                            <rect width={n.width} height={n.height} rx={4} fill="url(#oRackRelief)" style={{ pointerEvents: 'none' }} />
                                        )}
                                        {isAisle && (
                                            <>
                                                <rect width={n.width} height={n.height} rx={4} fill="url(#oAisleSunken)" style={{ pointerEvents: 'none' }} />
                                                <line x1={6} y1={n.height / 2} x2={n.width - 6} y2={n.height / 2} stroke="#ffffff" strokeOpacity={0.7} strokeWidth={2} strokeDasharray="10 8" style={{ pointerEvents: 'none' }} />
                                            </>
                                        )}
                                        <text x={n.width / 2} y={n.height / 2} textAnchor="middle" dominantBaseline="middle"
                                            fontSize={n.ref_type === 'location' ? 9 : 11}
                                            fill={n.ref_type === 'zone' ? '#111827' : '#1f2937'}
                                            style={{ pointerEvents: 'none', userSelect: 'none' }}>{n.label}</text>
                                    </g>
                                );
                            })}
                        </svg>
                    </div>

                    <div className="w-72 shrink-0 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3 max-h-[680px] overflow-auto">
                        <p className="text-xs font-semibold text-gray-700 dark:text-gray-200 mb-2">Detalle</p>
                        {!selectedRack ? (
                            <p className="text-xs text-gray-400">Selecciona un rack del plano para ver su contenido</p>
                        ) : (
                            <div className="space-y-2">
                                <p className="text-sm font-semibold text-gray-800 dark:text-gray-100">{selectedRack.name || selectedRack.code}</p>
                                {(() => { const a = rackAgg(selectedRack.id); return (
                                    <p className="text-[11px] text-gray-500">{a.qty} unidades{a.hasCap ? ` / ${a.cap} cap (${Math.round((a.qty / a.cap) * 100)}%)` : ''}</p>
                                ); })()}
                                {(selectedRack.levels || []).length === 0 ? (
                                    <p className="text-[11px] text-gray-400">Sin niveles</p>
                                ) : (
                                    <svg viewBox={`0 0 200 ${(selectedRack.levels || []).length * 44 + 14}`} className="w-full bg-gray-50 dark:bg-gray-900 rounded border border-gray-200 dark:border-gray-700">
                                        {[...(selectedRack.levels || [])].sort((a: any, b: any) => (b.ordinal || 0) - (a.ordinal || 0)).map((lv: any, i: number) => {
                                            const y = i * 44 + 4;
                                            const positions = lv.positions || [];
                                            const areaX = 42, areaW = 152;
                                            const shown = positions.slice(0, 12);
                                            const cellW = shown.length ? areaW / shown.length : areaW;
                                            return (
                                                <g key={lv.id}>
                                                    <rect x={6} y={y} width={188} height={40} rx={3} fill="#f97316" fillOpacity={0.1} stroke="#f97316" strokeWidth={1} />
                                                    <text x={12} y={y + 16} fontSize={11} fontWeight={600} fill="#9a3412">{lv.code}</text>
                                                    <text x={12} y={y + 30} fontSize={8} fill="#9a3412">{positions.length} ubic</text>
                                                    {shown.length === 0 ? (
                                                        <text x={118} y={y + 24} fontSize={9} textAnchor="middle" fill="#9ca3af">vacio</text>
                                                    ) : shown.map((p: any, j: number) => {
                                                        const o = occupancy[p.id];
                                                        const qty = o?.quantity ?? 0;
                                                        const cap = o?.capacity ?? null;
                                                        const c = occColor(qty, cap);
                                                        return (
                                                            <g key={p.id} style={{ cursor: 'pointer' }} onClick={() => router.push(`/inventory?warehouse=${warehouseId}`)}>
                                                                <title>{p.code}{o ? ` — ${qty}${cap ? '/' + cap : ''}` : ''}</title>
                                                                <rect x={areaX + j * cellW + 1} y={y + 6} width={Math.max(cellW - 2, 4)} height={28} rx={2} fill={c.fill} fillOpacity={0.8} stroke={pulseLoc === p.id ? '#2563eb' : c.stroke} strokeWidth={pulseLoc === p.id ? 2.5 : 0.8} />
                                                                {cellW > 14 && <text x={areaX + j * cellW + cellW / 2} y={y + 17} fontSize={8} fontWeight={600} textAnchor="middle" fill="#1f2937">{qty}</text>}
                                                                {cellW > 30 && cap ? <text x={areaX + j * cellW + cellW / 2} y={y + 27} fontSize={6} textAnchor="middle" fill="#4b5563">/{cap}</text> : null}
                                                            </g>
                                                        );
                                                    })}
                                                </g>
                                            );
                                        })}
                                        <rect x={2} y={(selectedRack.levels || []).length * 44 + 6} width={196} height={5} fill="#475569" />
                                    </svg>
                                )}
                                <p className="text-[10px] text-gray-400">Vista de frente (nivel mas alto arriba). Click en una ubicacion para ver su stock.</p>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
