'use client';

import { useState } from 'react';
import { ChevronRightIcon, ChevronDownIcon, PlusIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import { TreeZone, TreeAisle, TreeRack, TreeLevel, TreePosition } from '../../domain/hierarchy-types';

export type TreeNodeType = 'zone' | 'aisle' | 'rack' | 'level' | 'position';

interface Props {
    zones: TreeZone[];
    onCreateChild?: (parentType: TreeNodeType | 'root', parentId: number | null) => void;
    onEdit?: (type: TreeNodeType, id: number) => void;
    onDelete?: (type: TreeNodeType, id: number) => void;
    canEdit?: boolean;
}

export default function WarehouseHierarchyTree({ zones, onCreateChild, onEdit, onDelete, canEdit = true }: Props) {
    const [expanded, setExpanded] = useState<Record<string, boolean>>({});

    const toggle = (key: string) => setExpanded((e) => ({ ...e, [key]: !e[key] }));

    const NodeActions = ({ type, id, childLabel, onAddChild }: { type: TreeNodeType; id: number; childLabel?: string; onAddChild?: () => void }) => {
        if (!canEdit) return null;
        return (
            <div className="opacity-0 group-hover:opacity-100 transition flex items-center gap-1">
                {onAddChild && childLabel && (
                    <button
                        onClick={(e) => { e.stopPropagation(); onAddChild(); }}
                        className="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-500 hover:text-blue-600"
                        title={`Agregar ${childLabel}`}
                    >
                        <PlusIcon className="w-3.5 h-3.5" />
                    </button>
                )}
                <button
                    onClick={(e) => { e.stopPropagation(); onEdit?.(type, id); }}
                    className="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-500 hover:text-amber-600"
                    title="Editar"
                >
                    <PencilIcon className="w-3.5 h-3.5" />
                </button>
                <button
                    onClick={(e) => { e.stopPropagation(); onDelete?.(type, id); }}
                    className="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-500 hover:text-red-600"
                    title="Eliminar"
                >
                    <TrashIcon className="w-3.5 h-3.5" />
                </button>
            </div>
        );
    };

    const renderPosition = (p: TreePosition, levelId: number) => (
        <div key={`p-${p.id}`} className="group flex items-center justify-between pl-16 pr-3 py-1 hover:bg-gray-50 dark:hover:bg-gray-700/40 rounded">
            <div className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                <span className="text-xs px-1.5 py-0.5 rounded bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300">POS</span>
                <span className="font-mono">{p.code}</span>
                <span className="text-gray-400">{p.name}</span>
                {!p.is_active && <span className="text-xs text-red-500">inactivo</span>}
            </div>
            <NodeActions type="position" id={p.id} />
        </div>
    );

    const renderLevel = (l: TreeLevel) => {
        const key = `l-${l.id}`;
        const open = expanded[key] ?? true;
        return (
            <div key={key}>
                <div className="group flex items-center justify-between pl-12 pr-3 py-1.5 hover:bg-gray-50 dark:hover:bg-gray-700/40 rounded cursor-pointer" onClick={() => toggle(key)}>
                    <div className="flex items-center gap-2 text-sm">
                        {open ? <ChevronDownIcon className="w-4 h-4 text-gray-400" /> : <ChevronRightIcon className="w-4 h-4 text-gray-400" />}
                        <span className="text-xs px-1.5 py-0.5 rounded bg-indigo-100 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300">LVL</span>
                        <span className="font-mono text-gray-800 dark:text-gray-100">{l.code}</span>
                        <span className="text-xs text-gray-400">#{l.ordinal}</span>
                        <span className="text-xs text-gray-400">({l.positions.length})</span>
                    </div>
                    <NodeActions
                        type="level"
                        id={l.id}
                        childLabel="posicion"
                        onAddChild={() => onCreateChild?.('level', l.id)}
                    />
                </div>
                {open && l.positions.map((p) => renderPosition(p, l.id))}
            </div>
        );
    };

    const renderRack = (r: TreeRack) => {
        const key = `r-${r.id}`;
        const open = expanded[key] ?? true;
        return (
            <div key={key}>
                <div className="group flex items-center justify-between pl-8 pr-3 py-1.5 hover:bg-gray-50 dark:hover:bg-gray-700/40 rounded cursor-pointer" onClick={() => toggle(key)}>
                    <div className="flex items-center gap-2 text-sm">
                        {open ? <ChevronDownIcon className="w-4 h-4 text-gray-400" /> : <ChevronRightIcon className="w-4 h-4 text-gray-400" />}
                        <span className="text-xs px-1.5 py-0.5 rounded bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300">RCK</span>
                        <span className="font-mono text-gray-800 dark:text-gray-100">{r.code}</span>
                        <span className="text-gray-500 dark:text-gray-400">{r.name}</span>
                        <span className="text-xs text-gray-400">({r.levels.length} niveles)</span>
                    </div>
                    <NodeActions
                        type="rack"
                        id={r.id}
                        childLabel="nivel"
                        onAddChild={() => onCreateChild?.('rack', r.id)}
                    />
                </div>
                {open && r.levels.map(renderLevel)}
            </div>
        );
    };

    const renderAisle = (a: TreeAisle) => {
        const key = `a-${a.id}`;
        const open = expanded[key] ?? true;
        return (
            <div key={key}>
                <div className="group flex items-center justify-between pl-4 pr-3 py-2 hover:bg-gray-50 dark:hover:bg-gray-700/40 rounded cursor-pointer" onClick={() => toggle(key)}>
                    <div className="flex items-center gap-2 text-sm">
                        {open ? <ChevronDownIcon className="w-4 h-4 text-gray-400" /> : <ChevronRightIcon className="w-4 h-4 text-gray-400" />}
                        <span className="text-xs px-1.5 py-0.5 rounded bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300">PAS</span>
                        <span className="font-mono text-gray-800 dark:text-gray-100">{a.code}</span>
                        <span className="text-gray-500 dark:text-gray-400">{a.name}</span>
                        <span className="text-xs text-gray-400">({a.racks.length} racks)</span>
                    </div>
                    <NodeActions
                        type="aisle"
                        id={a.id}
                        childLabel="rack"
                        onAddChild={() => onCreateChild?.('aisle', a.id)}
                    />
                </div>
                {open && a.racks.map(renderRack)}
            </div>
        );
    };

    const renderZone = (z: TreeZone) => {
        const key = `z-${z.id}`;
        const open = expanded[key] ?? true;
        const color = z.color_hex || '#6366f1';
        return (
            <div key={key} className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden bg-white dark:bg-gray-800">
                <div className="group flex items-center justify-between px-3 py-2.5 hover:bg-gray-50 dark:hover:bg-gray-700/60 cursor-pointer" onClick={() => toggle(key)}>
                    <div className="flex items-center gap-2">
                        {open ? <ChevronDownIcon className="w-4 h-4 text-gray-400" /> : <ChevronRightIcon className="w-4 h-4 text-gray-400" />}
                        <span className="w-3 h-3 rounded-full" style={{ backgroundColor: color }} />
                        <span className="text-xs px-1.5 py-0.5 rounded font-medium" style={{ backgroundColor: `${color}22`, color }}>ZONA</span>
                        <span className="font-semibold text-gray-900 dark:text-white">{z.code}</span>
                        <span className="text-gray-500 dark:text-gray-300">{z.name}</span>
                        {z.purpose && <span className="text-xs text-gray-400 italic">{z.purpose}</span>}
                        <span className="text-xs text-gray-400 ml-1">({z.aisles.length} pasillos)</span>
                    </div>
                    <NodeActions
                        type="zone"
                        id={z.id}
                        childLabel="pasillo"
                        onAddChild={() => onCreateChild?.('zone', z.id)}
                    />
                </div>
                {open && (
                    <div className="pl-2 pr-1 py-1 bg-gray-50/40 dark:bg-gray-900/40">
                        {z.aisles.length === 0 ? (
                            <p className="pl-6 py-2 text-xs text-gray-400 italic">Sin pasillos</p>
                        ) : z.aisles.map(renderAisle)}
                    </div>
                )}
            </div>
        );
    };

    if (zones.length === 0) {
        return (
            <div className="text-center py-12 bg-white dark:bg-gray-800 rounded-lg border border-dashed border-gray-300 dark:border-gray-700">
                <p className="text-gray-500 dark:text-gray-400 mb-4 text-sm">Esta bodega aun no tiene jerarquia configurada.</p>
                {canEdit && (
                    <button
                        onClick={() => onCreateChild?.('root', null)}
                        className="inline-flex items-center gap-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-md text-sm font-medium"
                    >
                        <PlusIcon className="w-4 h-4" /> Crear primera zona
                    </button>
                )}
            </div>
        );
    }

    return (
        <div className="space-y-2">
            {canEdit && (
                <div className="flex justify-end">
                    <button
                        onClick={() => onCreateChild?.('root', null)}
                        className="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm bg-indigo-600 hover:bg-indigo-700 text-white rounded-md font-medium"
                    >
                        <PlusIcon className="w-4 h-4" /> Nueva zona
                    </button>
                </div>
            )}
            {zones.map(renderZone)}
        </div>
    );
}
