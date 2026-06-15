'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button, Spinner, Alert } from '@/shared/ui';
import { getLayoutAction, saveLayoutAction, getOccupancyAction } from '../../infra/actions/hierarchy';
import { LayoutNode, LayoutRefType, WarehouseLayout, WarehouseTree } from '../../domain/hierarchy-types';

interface Props {
    warehouseId: number;
    businessId?: number;
    tree: WarehouseTree | null;
}

interface PaletteItem {
    refType: LayoutRefType;
    refId: number;
    label: string;
    parent: string;
}

const TYPE_COLOR: Record<LayoutRefType, string> = {
    zone: '#6366f1',
    aisle: '#10b981',
    rack: '#8b5cf6',
    level: '#f97316',
    location: '#f43f5e',
    wall: '#475569',
    dock: '#f59e0b',
    label: '#0ea5e9',
};

const TYPE_LABEL: Record<LayoutRefType, string> = {
    zone: 'Zona',
    aisle: 'Pasillo',
    rack: 'Rack',
    level: 'Nivel',
    location: 'Ubicacion',
    wall: 'Muro',
    dock: 'Muelle',
    label: 'Texto',
};

const DEFAULT_SIZE: Record<LayoutRefType, { w: number; h: number }> = {
    zone: { w: 240, h: 160 },
    aisle: { w: 200, h: 60 },
    rack: { w: 80, h: 60 },
    level: { w: 140, h: 16 },
    location: { w: 40, h: 40 },
    wall: { w: 200, h: 16 },
    dock: { w: 100, h: 40 },
    label: { w: 120, h: 30 },
};

function flattenTree(tree: WarehouseTree | null): PaletteItem[] {
    if (!tree) return [];
    const items: PaletteItem[] = [];
    for (const z of tree.zones || []) {
        items.push({ refType: 'zone', refId: z.id, label: z.name || z.code, parent: 'Bodega' });
        for (const a of z.aisles || []) {
            items.push({ refType: 'aisle', refId: a.id, label: a.name || a.code, parent: z.name || z.code });
            for (const r of a.racks || []) {
                items.push({ refType: 'rack', refId: r.id, label: r.name || r.code, parent: a.name || a.code });
                for (const l of r.levels || []) {
                    items.push({ refType: 'level', refId: l.id, label: l.code || `N${l.ordinal ?? ''}`, parent: r.name || r.code });
                    for (const p of l.positions || []) {
                        items.push({ refType: 'location', refId: p.id, label: p.code, parent: r.name || r.code });
                    }
                }
            }
        }
    }
    return items;
}

function findRackInTree(tree: WarehouseTree | null, rackId: number): { name: string; levels: any[]; width_cm: number; depth_cm: number; height_cm: number } | null {
    if (!tree) return null;
    for (const z of tree.zones || []) {
        for (const a of z.aisles || []) {
            for (const r of a.racks || []) {
                if (r.id === rackId) return { name: r.name || r.code, levels: r.levels || [], width_cm: r.width_cm || 0, depth_cm: r.depth_cm || 0, height_cm: r.height_cm || 0 };
            }
        }
    }
    return null;
}

function snap(value: number, grid: number): number {
    if (grid <= 0) return Math.round(value);
    return Math.round(value / grid) * grid;
}

interface OccCell { quantity: number; capacity: number | null }

function occColor(occ: OccCell | undefined): { fill: string; stroke: string } {
    if (!occ || occ.quantity <= 0) return { fill: '#e5e7eb', stroke: '#9ca3af' };
    if (occ.capacity && occ.capacity > 0) {
        const ratio = occ.quantity / occ.capacity;
        if (ratio >= 0.85) return { fill: '#f87171', stroke: '#dc2626' };
        if (ratio >= 0.5) return { fill: '#fbbf24', stroke: '#d97706' };
        return { fill: '#34d399', stroke: '#059669' };
    }
    return { fill: '#34d399', stroke: '#059669' };
}

export default function WarehouseLayout2D({ warehouseId, businessId, tree }: Props) {
    const router = useRouter();
    const [layout, setLayout] = useState<WarehouseLayout | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [saving, setSaving] = useState(false);
    const [dirty, setDirty] = useState(false);
    const [saved, setSaved] = useState(false);
    const [selectedId, setSelectedId] = useState<string | null>(null);
    const [occupancy, setOccupancy] = useState<Record<number, OccCell>>({});
    const [zoom, setZoom] = useState(0.7);
    const [paletteFilter, setPaletteFilter] = useState<LayoutRefType | 'all'>('all');

    const svgRef = useRef<SVGSVGElement | null>(null);
    const dragRef = useRef<{ id: string; mode: 'move' | 'resize'; startX: number; startY: number; origX: number; origY: number; origW: number; origH: number } | null>(null);

    useEffect(() => {
        let active = true;
        (async () => {
            setLoading(true);
            setError(null);
            try {
                const res = await getLayoutAction(warehouseId, businessId);
                if (active) setLayout(res);
            } catch (e: any) {
                if (active) setError(e?.message || 'Error al cargar el plano');
            } finally {
                if (active) setLoading(false);
            }
        })();
        return () => { active = false; };
    }, [warehouseId, businessId]);

    useEffect(() => {
        let active = true;
        (async () => {
            try {
                const res = await getOccupancyAction(warehouseId, businessId);
                if (!active) return;
                const map: Record<number, OccCell> = {};
                for (const l of res.locations || []) {
                    map[l.location_id] = { quantity: l.quantity, capacity: l.capacity };
                }
                setOccupancy(map);
            } catch {
                setOccupancy({});
            }
        })();
        return () => { active = false; };
    }, [warehouseId, businessId]);

    const placedRefs = useMemo(() => {
        const set = new Set<string>();
        for (const n of layout?.nodes || []) {
            if (n.ref_id > 0) set.add(`${n.ref_type}:${n.ref_id}`);
        }
        return set;
    }, [layout]);

    const palette = useMemo(() => {
        const all = flattenTree(tree).filter((it) => !placedRefs.has(`${it.refType}:${it.refId}`));
        if (paletteFilter === 'all') return all;
        return all.filter((it) => it.refType === paletteFilter);
    }, [tree, placedRefs, paletteFilter]);

    const updateNode = useCallback((id: string, patch: Partial<LayoutNode>) => {
        setLayout((prev) => {
            if (!prev) return prev;
            return { ...prev, nodes: prev.nodes.map((n) => (n.node_id === id ? { ...n, ...patch } : n)) };
        });
        setDirty(true);
        setSaved(false);
    }, []);

    const addNode = useCallback((refType: LayoutRefType, refId: number, label: string) => {
        setLayout((prev) => {
            if (!prev) return prev;
            const size = DEFAULT_SIZE[refType];
            const count = prev.nodes.length;
            const node: LayoutNode = {
                node_id: `${refType}-${refId || 'free'}-${count}-${prev.nodes.length}-${refId}${count}`,
                ref_type: refType,
                ref_id: refId,
                x: snap(40 + (count % 8) * 30, prev.grid_size),
                y: snap(40 + Math.floor(count / 8) * 30, prev.grid_size),
                width: size.w,
                height: size.h,
                rotation: 0,
                color: TYPE_COLOR[refType],
                label,
            };
            return { ...prev, nodes: [...prev.nodes, node] };
        });
        setDirty(true);
        setSaved(false);
    }, []);

    const removeNode = useCallback((id: string) => {
        setLayout((prev) => prev ? { ...prev, nodes: prev.nodes.filter((n) => n.node_id !== id) } : prev);
        setSelectedId(null);
        setDirty(true);
        setSaved(false);
    }, []);

    const autoArrange = useCallback(() => {
        const items = flattenTree(tree).filter((it) => !placedRefs.has(`${it.refType}:${it.refId}`) && (paletteFilter === 'all' ? it.refType === 'rack' : it.refType === paletteFilter));
        if (!items.length) return;
        setLayout((prev) => {
            if (!prev) return prev;
            const grid = prev.grid_size || 20;
            const perRow = Math.max(1, Math.floor((prev.canvas_width - 40) / (DEFAULT_SIZE[items[0].refType].w + grid)));
            const start = prev.nodes.length;
            const newNodes: LayoutNode[] = items.map((it, i) => {
                const size = DEFAULT_SIZE[it.refType];
                const col = i % perRow;
                const row = Math.floor(i / perRow);
                return {
                    node_id: `${it.refType}-${it.refId}-auto-${start + i}`,
                    ref_type: it.refType,
                    ref_id: it.refId,
                    x: snap(20 + col * (size.w + grid), grid),
                    y: snap(20 + row * (size.h + grid), grid),
                    width: size.w,
                    height: size.h,
                    rotation: 0,
                    color: TYPE_COLOR[it.refType],
                    label: it.label,
                };
            });
            return { ...prev, nodes: [...prev.nodes, ...newNodes] };
        });
        setDirty(true);
        setSaved(false);
    }, [tree, placedRefs, paletteFilter]);

    const buildAutoPlan = useCallback(() => {
        if (!tree) return;
        const W = 1200;
        const M = 30;
        const zoneGap = 24;
        const zonePadTop = 36;
        const zonePadBottom = 16;
        const aisleGap = 20;
        const rackGap = 10;
        const contentX = M + 20;
        const pxPerM = layout?.scale && layout.scale > 0 ? layout.scale : 40;
        const cmPx = (cm: number, def: number) => ((cm && cm > 0 ? cm : def) / 100) * pxPerM;
        const zoneColors = ['#6366f1', '#14b8a6', '#f59e0b', '#ec4899', '#0ea5e9', '#84cc16'];

        const mk = (refType: LayoutRefType, refId: number, x: number, y: number, w: number, h: number, color: string, label: string): LayoutNode => ({
            node_id: `${refType}-${refId || `${Math.round(x)}x${Math.round(y)}`}`,
            ref_type: refType, ref_id: refId, x, y, width: w, height: h, rotation: 0, color, label,
        });

        const nodes: LayoutNode[] = [];
        let y = M + 10;
        let canvasWidth = W;

        (tree.zones || []).forEach((z, zi) => {
            const zoneTop = y;
            let inner = zoneTop + zonePadTop;
            let zoneRight = contentX;
            (z.aisles || []).forEach((a) => {
                const racks = a.racks || [];
                const sideA = racks.filter((r) => r.side === 'A');
                const sideB = racks.filter((r) => r.side === 'B');
                racks.filter((r) => r.side !== 'A' && r.side !== 'B').forEach((r, i) => (i % 2 === 0 ? sideA : sideB).push(r));
                const rowDepthA = sideA.length ? Math.max(...sideA.map((r) => cmPx(r.depth_cm, 100))) : cmPx(0, 100);
                const rowDepthB = sideB.length ? Math.max(...sideB.map((r) => cmPx(r.depth_cm, 100))) : cmPx(0, 100);
                const corridorH = cmPx(a.width_cm, 300);

                let xa = contentX;
                sideA.forEach((r) => {
                    const w = cmPx(r.width_cm, 120);
                    const h = cmPx(r.depth_cm, 100);
                    nodes.push(mk('rack', r.id, xa, inner + (rowDepthA - h), w, h, TYPE_COLOR.rack, `${r.name || r.code} (A)`));
                    xa += w + rackGap;
                });
                const corrY = inner + rowDepthA + 4;
                const corridorW = Math.max(xa - contentX - rackGap, 200);
                const widthM = a.width_cm && a.width_cm > 0 ? a.width_cm / 100 : 3;
                nodes.push(mk('aisle', a.id, contentX, corrY, corridorW, corridorH, '#cbd5e1', `${a.name || a.code} (transito ${widthM} m)`));
                const botY = corrY + corridorH + 4;
                let xb = contentX;
                sideB.forEach((r) => {
                    const w = cmPx(r.width_cm, 120);
                    const h = cmPx(r.depth_cm, 100);
                    nodes.push(mk('rack', r.id, xb, botY, w, h, TYPE_COLOR.rack, `${r.name || r.code} (B)`));
                    xb += w + rackGap;
                });
                zoneRight = Math.max(zoneRight, xa, xb, contentX + corridorW);
                inner = botY + rowDepthB + aisleGap;
            });
            const zoneW = Math.max(W - 2 * M, zoneRight - M + 20);
            const zoneBottom = inner + zonePadBottom;
            nodes.unshift(mk('zone', z.id, M, zoneTop, zoneW, zoneBottom - zoneTop, zoneColors[zi % zoneColors.length], `ZONA ${z.name || z.code}`));
            canvasWidth = Math.max(canvasWidth, M + zoneW + M);
            y = zoneBottom + zoneGap;
        });

        const canvasHeight = Math.max(600, Math.round(y + M));
        setLayout((prev) => (prev ? { ...prev, canvas_width: Math.round(canvasWidth), canvas_height: canvasHeight, nodes } : prev));
        setSelectedId(null);
        setDirty(true);
        setSaved(false);
    }, [tree, layout?.scale]);

    const onPointerDownNode = (e: React.PointerEvent, id: string, mode: 'move' | 'resize') => {
        e.stopPropagation();
        const node = layout?.nodes.find((n) => n.node_id === id);
        if (!node) return;
        setSelectedId(id);
        (e.target as Element).setPointerCapture?.(e.pointerId);
        dragRef.current = {
            id, mode,
            startX: e.clientX, startY: e.clientY,
            origX: node.x, origY: node.y, origW: node.width, origH: node.height,
        };
    };

    const onPointerMove = (e: React.PointerEvent) => {
        const d = dragRef.current;
        if (!d || !layout) return;
        const dx = (e.clientX - d.startX) / zoom;
        const dy = (e.clientY - d.startY) / zoom;
        const grid = layout.grid_size || 20;
        if (d.mode === 'move') {
            updateNode(d.id, { x: snap(d.origX + dx, grid), y: snap(d.origY + dy, grid) });
        } else {
            updateNode(d.id, {
                width: Math.max(grid, snap(d.origW + dx, grid)),
                height: Math.max(grid, snap(d.origH + dy, grid)),
            });
        }
    };

    const onPointerUp = () => { dragRef.current = null; };

    const handleSave = async () => {
        if (!layout) return;
        setSaving(true);
        setError(null);
        const res = await saveLayoutAction(warehouseId, {
            canvas_width: layout.canvas_width,
            canvas_height: layout.canvas_height,
            grid_size: layout.grid_size,
            scale: layout.scale,
            nodes: layout.nodes,
        }, businessId);
        setSaving(false);
        if (res.success) {
            setDirty(false);
            setSaved(true);
        } else {
            setError(res.error);
        }
    };

    const selected = layout?.nodes.find((n) => n.node_id === selectedId) || null;
    const rackDetail = useMemo(
        () => (selected && selected.ref_type === 'rack' && selected.ref_id > 0 ? findRackInTree(tree, selected.ref_id) : null),
        [selected, tree],
    );

    if (loading) {
        return <div className="flex items-center justify-center py-16"><Spinner size="lg" /></div>;
    }
    if (!layout) {
        return <Alert type="error">{error || 'No se pudo cargar el plano'}</Alert>;
    }

    const W = layout.canvas_width;
    const H = layout.canvas_height;
    const grid = layout.grid_size || 20;
    const pxPerM = layout.scale && layout.scale > 0 ? layout.scale : 40;

    return (
        <div className="space-y-3">
            <div className="flex items-center justify-between flex-wrap gap-2">
                <div className="flex items-center gap-2">
                    <span className="text-sm text-gray-500 dark:text-gray-400">Zoom</span>
                    <button className="px-2 py-1 rounded border border-gray-300 dark:border-gray-600 text-sm" onClick={() => setZoom((z) => Math.max(0.25, +(z - 0.1).toFixed(2)))}>-</button>
                    <span className="text-sm w-12 text-center">{Math.round(zoom * 100)}%</span>
                    <button className="px-2 py-1 rounded border border-gray-300 dark:border-gray-600 text-sm" onClick={() => setZoom((z) => Math.min(2, +(z + 0.1).toFixed(2)))}>+</button>
                    <button className="px-2 py-1 rounded border border-gray-300 dark:border-gray-600 text-sm" onClick={autoArrange}>Auto acomodar</button>
                    <button className="px-2 py-1 rounded bg-indigo-600 text-white text-sm" onClick={buildAutoPlan}>Auto-plano</button>
                    <span className="ml-2 text-sm text-gray-500 dark:text-gray-400">Escala</span>
                    <input
                        type="number"
                        min={5}
                        max={200}
                        value={Math.round(pxPerM)}
                        onChange={(e) => { const v = Number(e.target.value) || 40; setLayout((p) => p ? { ...p, scale: v } : p); setDirty(true); setSaved(false); }}
                        className="w-16 px-2 py-1 rounded border border-gray-300 dark:border-gray-600 text-sm bg-transparent"
                    />
                    <span className="text-xs text-gray-400">px/m</span>
                </div>
                <div className="flex items-center gap-2">
                    {dirty && <span className="text-xs text-amber-600">Cambios sin guardar</span>}
                    {saved && !dirty && <span className="text-xs text-emerald-600">Guardado</span>}
                    <Button onClick={handleSave} disabled={saving || !dirty}>{saving ? 'Guardando...' : 'Guardar plano'}</Button>
                </div>
            </div>

            {error && <Alert type="error">{error}</Alert>}

            <div className="flex gap-3">
                <div className="w-56 shrink-0 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3 max-h-[640px] overflow-auto">
                    <p className="text-xs font-semibold text-gray-700 dark:text-gray-200 mb-2">Elementos</p>
                    <div className="flex flex-wrap gap-1 mb-2">
                        {(['all', 'zone', 'aisle', 'rack', 'level', 'location'] as const).map((t) => (
                            <button key={t}
                                className={`px-2 py-0.5 rounded text-xs border ${paletteFilter === t ? 'bg-indigo-600 text-white border-indigo-600' : 'border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300'}`}
                                onClick={() => setPaletteFilter(t)}>
                                {t === 'all' ? 'Todos' : TYPE_LABEL[t]}
                            </button>
                        ))}
                    </div>
                    <div className="flex flex-wrap gap-1 mb-2">
                        {(['wall', 'dock', 'label'] as LayoutRefType[]).map((t) => (
                            <button key={t}
                                className="px-2 py-0.5 rounded text-xs border border-dashed border-gray-400 dark:border-gray-500 text-gray-600 dark:text-gray-300"
                                onClick={() => addNode(t, 0, TYPE_LABEL[t])}>
                                + {TYPE_LABEL[t]}
                            </button>
                        ))}
                    </div>
                    {palette.length === 0 ? (
                        <p className="text-xs text-gray-400">Todo ubicado</p>
                    ) : (
                        <ul className="space-y-1">
                            {palette.slice(0, 300).map((it) => (
                                <li key={`${it.refType}:${it.refId}`}>
                                    <button
                                        className="w-full text-left px-2 py-1 rounded text-xs hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center gap-2"
                                        onClick={() => addNode(it.refType, it.refId, it.label)}>
                                        <span className="w-2.5 h-2.5 rounded-sm shrink-0" style={{ background: TYPE_COLOR[it.refType] }} />
                                        <span className="truncate">{it.label}</span>
                                        <span className="ml-auto text-[10px] text-gray-400">{TYPE_LABEL[it.refType]}</span>
                                    </button>
                                </li>
                            ))}
                        </ul>
                    )}
                </div>

                <div className="flex-1 overflow-auto bg-gray-100 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg" style={{ maxHeight: 640 }}>
                    <svg
                        ref={svgRef}
                        viewBox={`0 0 ${W} ${H}`}
                        width={W * zoom}
                        height={H * zoom}
                        onPointerMove={onPointerMove}
                        onPointerUp={onPointerUp}
                        onPointerLeave={onPointerUp}
                        style={{ background: 'white', display: 'block' }}
                    >
                        <defs>
                            <pattern id="grid" width={grid} height={grid} patternUnits="userSpaceOnUse">
                                <path d={`M ${grid} 0 L 0 0 0 ${grid}`} fill="none" stroke="#e5e7eb" strokeWidth="1" />
                            </pattern>
                        </defs>
                        <rect x={0} y={0} width={W} height={H} fill="url(#grid)" onPointerDown={() => setSelectedId(null)} />

                        <g transform={`translate(${grid} ${H - grid})`} style={{ pointerEvents: 'none' }}>
                            <line x1={0} y1={0} x2={pxPerM} y2={0} stroke="#111827" strokeWidth={2} />
                            <line x1={0} y1={-5} x2={0} y2={5} stroke="#111827" strokeWidth={2} />
                            <line x1={pxPerM} y1={-5} x2={pxPerM} y2={5} stroke="#111827" strokeWidth={2} />
                            <text x={pxPerM / 2} y={-8} textAnchor="middle" fontSize={12} fill="#111827">1 m</text>
                        </g>

                        {layout.nodes.map((n) => {
                            const isSel = n.node_id === selectedId;
                            const isLabel = n.ref_type === 'label';
                            return (
                                <g key={n.node_id} transform={`translate(${n.x} ${n.y}) rotate(${n.rotation} ${n.width / 2} ${n.height / 2})`}>
                                    {!isLabel && (
                                        <rect
                                            width={n.width} height={n.height} rx={4}
                                            fill={n.color} fillOpacity={n.ref_type === 'zone' ? 0.18 : 0.7}
                                            stroke={n.color} strokeWidth={isSel ? 3 : 1.5}
                                            style={{ cursor: 'move' }}
                                            onPointerDown={(e) => onPointerDownNode(e, n.node_id, 'move')}
                                        />
                                    )}
                                    <text
                                        x={isLabel ? 0 : n.width / 2} y={isLabel ? 14 : n.height / 2}
                                        textAnchor={isLabel ? 'start' : 'middle'} dominantBaseline="middle"
                                        fontSize={n.ref_type === 'location' ? 9 : 12}
                                        fill={isLabel ? n.color : (n.ref_type === 'zone' ? '#111827' : '#ffffff')}
                                        style={{ pointerEvents: isLabel ? 'auto' : 'none', cursor: isLabel ? 'move' : 'default', userSelect: 'none' }}
                                        onPointerDown={isLabel ? (e) => onPointerDownNode(e, n.node_id, 'move') : undefined}
                                    >
                                        {n.label}
                                    </text>
                                    {isSel && !isLabel && (
                                        <rect
                                            x={n.width - 10} y={n.height - 10} width={12} height={12}
                                            fill="#ffffff" stroke={n.color} strokeWidth={2}
                                            style={{ cursor: 'nwse-resize' }}
                                            onPointerDown={(e) => onPointerDownNode(e, n.node_id, 'resize')}
                                        />
                                    )}
                                </g>
                            );
                        })}
                    </svg>
                </div>

                <div className="w-64 shrink-0 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3 max-h-[640px] overflow-auto">
                    <p className="text-xs font-semibold text-gray-700 dark:text-gray-200 mb-2">Propiedades</p>
                    {!selected ? (
                        <p className="text-xs text-gray-400">Selecciona un elemento del plano</p>
                    ) : (
                        <div className="space-y-3 text-xs">
                            <div>
                                <label className="block text-gray-500 mb-1">Etiqueta</label>
                                <input
                                    className="w-full px-2 py-1 rounded border border-gray-300 dark:border-gray-600 bg-transparent"
                                    value={selected.label}
                                    onChange={(e) => updateNode(selected.node_id, { label: e.target.value })}
                                />
                            </div>
                            <div className="grid grid-cols-2 gap-2">
                                <div>
                                    <label className="block text-gray-500 mb-1">Ancho</label>
                                    <input type="number" className="w-full px-2 py-1 rounded border border-gray-300 dark:border-gray-600 bg-transparent"
                                        value={Math.round(selected.width)}
                                        onChange={(e) => updateNode(selected.node_id, { width: Math.max(grid, Number(e.target.value) || grid) })} />
                                </div>
                                <div>
                                    <label className="block text-gray-500 mb-1">Alto</label>
                                    <input type="number" className="w-full px-2 py-1 rounded border border-gray-300 dark:border-gray-600 bg-transparent"
                                        value={Math.round(selected.height)}
                                        onChange={(e) => updateNode(selected.node_id, { height: Math.max(grid, Number(e.target.value) || grid) })} />
                                </div>
                            </div>
                            <div>
                                <label className="block text-gray-500 mb-1">Rotacion {Math.round(selected.rotation)} grados</label>
                                <input type="range" min={0} max={350} step={10} className="w-full"
                                    value={selected.rotation}
                                    onChange={(e) => updateNode(selected.node_id, { rotation: Number(e.target.value) })} />
                            </div>
                            <div>
                                <label className="block text-gray-500 mb-1">Color</label>
                                <input type="color" className="w-full h-8 rounded border border-gray-300 dark:border-gray-600 bg-transparent"
                                    value={selected.color}
                                    onChange={(e) => updateNode(selected.node_id, { color: e.target.value })} />
                            </div>
                            {rackDetail && (
                                <div className="pt-2 border-t border-gray-200 dark:border-gray-700">
                                    <p className="text-[11px] font-semibold text-gray-700 dark:text-gray-200 mb-1">
                                        Elevacion frontal - {rackDetail.name}
                                    </p>
                                    <p className="text-[10px] text-gray-400 mb-1">
                                        {((rackDetail.width_cm || 0) / 100) || '?'} m ancho · {((rackDetail.depth_cm || 0) / 100) || '?'} m fondo · {((rackDetail.height_cm || 0) / 100) || '?'} m alto
                                    </p>
                                    {rackDetail.levels.length === 0 ? (
                                        <p className="text-[11px] text-gray-400">Este rack no tiene niveles</p>
                                    ) : (
                                        <svg viewBox={`0 0 200 ${rackDetail.levels.length * 44 + 14}`} className="w-full bg-gray-50 dark:bg-gray-900 rounded border border-gray-200 dark:border-gray-700">
                                            {[...rackDetail.levels]
                                                .sort((a, b) => (b.ordinal || 0) - (a.ordinal || 0))
                                                .map((lv, i) => {
                                                    const y = i * 44 + 4;
                                                    const positions = lv.positions || [];
                                                    const areaX = 42;
                                                    const areaW = 152;
                                                    const shown = positions.slice(0, 12);
                                                    const cellW = shown.length ? areaW / shown.length : areaW;
                                                    return (
                                                        <g key={lv.id}>
                                                            <rect x={6} y={y} width={188} height={40} rx={3} fill="#fb923c" fillOpacity={0.12} stroke="#f97316" strokeWidth={1.2} />
                                                            <text x={12} y={y + 16} fontSize={11} fontWeight={600} fill="#9a3412">{lv.code}</text>
                                                            <text x={12} y={y + 30} fontSize={8} fill="#9a3412">{positions.length} ubic</text>
                                                            {shown.length === 0 ? (
                                                                <text x={118} y={y + 24} fontSize={9} textAnchor="middle" fill="#9ca3af">vacio</text>
                                                                            ) : shown.map((p: any, j: number) => {
                                                                const occ = occupancy[p.id];
                                                                const col = occColor(occ);
                                                                const qty = occ?.quantity ?? 0;
                                                                const cap = occ?.capacity;
                                                                return (
                                                                    <g key={p.id} style={{ cursor: 'pointer' }} onClick={() => router.push(`/inventory?warehouse=${warehouseId}`)}>
                                                                        <title>{p.code}{occ ? ` — ${qty}${cap ? '/' + cap : ''}` : ''}</title>
                                                                        <rect x={areaX + j * cellW + 1} y={y + 6} width={Math.max(cellW - 2, 4)} height={28} rx={2}
                                                                            fill={col.fill} fillOpacity={0.75} stroke={col.stroke} strokeWidth={0.8} />
                                                                        {cellW > 14 && <text x={areaX + j * cellW + cellW / 2} y={y + 17} fontSize={8} fontWeight={600} textAnchor="middle" fill="#1f2937">{qty}</text>}
                                                                        {cellW > 30 && cap ? <text x={areaX + j * cellW + cellW / 2} y={y + 27} fontSize={6} textAnchor="middle" fill="#4b5563">/{cap}</text> : null}
                                                                    </g>
                                                                );
                                                            })}
                                                            {positions.length > shown.length && (
                                                                <text x={192} y={y + 12} fontSize={7} textAnchor="end" fill="#9a3412">+{positions.length - shown.length}</text>
                                                            )}
                                                        </g>
                                                    );
                                                })}
                                            <rect x={2} y={rackDetail.levels.length * 44 + 6} width={196} height={5} fill="#475569" />
                                        </svg>
                                    )}
                                    <p className="text-[10px] text-gray-400 mt-1">Vista de frente (nivel mas alto arriba)</p>
                                </div>
                            )}
                            <button className="w-full px-2 py-1 rounded bg-rose-600 text-white" onClick={() => removeNode(selected.node_id)}>
                                Quitar del plano
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
