'use client';

import { useState, useEffect, useCallback } from 'react';
import {
    ChevronRightIcon,
    ChevronDownIcon,
    PlusIcon,
    PencilIcon,
    TrashIcon,
    MapIcon,
} from '@heroicons/react/24/outline';
import { Alert, Spinner } from '@/shared/ui';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { getWarehousesAction, deleteWarehouseAction } from '../../infra/actions';
import { Warehouse } from '../../domain/types';
import {
    getWarehouseTreeAction,
    deleteZoneAction,
    deleteAisleAction,
    deleteRackAction,
    deleteRackLevelAction,
} from '../../infra/actions/hierarchy';
import { WarehouseTree, TreeZone, TreeAisle, TreeRack, TreeLevel, TreePosition } from '../../domain/hierarchy-types';
import HierarchyNodeModal, { NodeType } from './HierarchyNodeModal';

type NodeKind = 'zone' | 'aisle' | 'rack' | 'level' | 'position';
type StructureMode = 'simple' | 'zones' | 'wms';

interface Props {
    businessId?: number;
    onEditWarehouse: (w: Warehouse) => void;
    onNewWarehouse: () => void;
    refreshKey?: number;
}

type ModalState =
    | { mode: 'create'; warehouseId: number; type: NodeType; parentId: number | null }
    | { mode: 'edit'; warehouseId: number; type: NodeType; parentId: null; initial: Record<string, any> }
    | null;

export default function WarehouseTreeTable({ businessId, onEditWarehouse, onNewWarehouse, refreshKey }: Props) {
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');

    const [expanded, setExpanded] = useState<Record<number, boolean>>({});
    const [trees, setTrees] = useState<Record<number, WarehouseTree>>({});
    const [loadingTree, setLoadingTree] = useState<Record<number, boolean>>({});
    const [expandedNodes, setExpandedNodes] = useState<Record<string, boolean>>({});

    const [modal, setModal] = useState<ModalState>(null);
    const [deletingWh, setDeletingWh] = useState<Warehouse | null>(null);
    const [deletingNode, setDeletingNode] = useState<{ warehouseId: number; type: NodeKind; id: number; label: string } | null>(null);
    const fetchList = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await getWarehousesAction({ page: 1, page_size: 100, business_id: businessId, search: search || undefined });
            setWarehouses(r.data || []);
        } catch (e: any) {
            setError(e.message || 'Error al cargar bodegas');
        } finally {
            setLoading(false);
        }
    }, [businessId, search]);

    useEffect(() => { fetchList(); }, [fetchList, refreshKey]);

    const loadTree = useCallback(async (warehouseId: number) => {
        setLoadingTree((p) => ({ ...p, [warehouseId]: true }));
        try {
            const t = await getWarehouseTreeAction(warehouseId, businessId);
            setTrees((p) => ({ ...p, [warehouseId]: t }));
        } catch {} finally {
            setLoadingTree((p) => ({ ...p, [warehouseId]: false }));
        }
    }, [businessId]);

    const toggleWh = (wh: Warehouse) => {
        const next = !expanded[wh.id];
        setExpanded((p) => ({ ...p, [wh.id]: next }));
        if (next && !trees[wh.id]) loadTree(wh.id);
    };

    const toggleNode = (key: string) => setExpandedNodes((p) => ({ ...p, [key]: !p[key] }));

    const refreshTree = (warehouseId: number) => loadTree(warehouseId);

    const handleDeleteWh = async () => {
        if (!deletingWh) return;
        await deleteWarehouseAction(deletingWh.id, businessId);
        setDeletingWh(null);
        fetchList();
    };

    const handleDeleteNode = async () => {
        if (!deletingNode) return;
        const { warehouseId, type, id } = deletingNode;
        if (type === 'zone') await deleteZoneAction(id, warehouseId, businessId);
        else if (type === 'aisle') await deleteAisleAction(id, warehouseId, businessId);
        else if (type === 'rack') await deleteRackAction(id, warehouseId, businessId);
        else if (type === 'level') await deleteRackLevelAction(id, warehouseId, businessId);
        setDeletingNode(null);
        refreshTree(warehouseId);
    };

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(searchInput);
    };

    const countStats = (tree: WarehouseTree | undefined) => {
        if (!tree) return { z: 0, a: 0, r: 0, l: 0, p: 0 };
        let z = 0, a = 0, r = 0, l = 0, p = 0;
        for (const zone of tree.zones) {
            z++;
            for (const ai of zone.aisles) {
                a++;
                for (const rk of ai.racks) {
                    r++;
                    for (const lv of rk.levels) {
                        l++;
                        p += lv.positions.length;
                    }
                }
            }
        }
        return { z, a, r, l, p };
    };

    return (
        <div className="space-y-4">
            <form onSubmit={handleSearch} className="flex gap-2">
                <input
                    type="text"
                    value={searchInput}
                    onChange={(e) => setSearchInput(e.target.value)}
                    placeholder="Buscar por nombre o código..."
                    className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                />
                <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700">Buscar</button>
                {search && <button type="button" onClick={() => { setSearchInput(''); setSearch(''); }} className="px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 rounded-lg text-sm">Limpiar</button>}
            </form>

            {error && <Alert type="error">{error}</Alert>}

            {loading && warehouses.length === 0 ? (
                <div className="flex justify-center p-8"><Spinner size="lg" /></div>
            ) : (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="table w-full">
                            <thead>
                                <tr>
                                    <th style={{ width: 40 }}></th>
                                    <th>Bodega</th>
                                    <th>Código</th>
                                    <th>Tipo</th>
                                    <th>Ubicación</th>
                                    <th style={{ minWidth: 220 }}>Jerarquía</th>
                                    <th style={{ textAlign: 'center' }}>Estado</th>
                                    <th style={{ textAlign: 'right' }}>Acciones</th>
                                </tr>
                            </thead>
                            <tbody>
                                {warehouses.length === 0 && !loading && (
                                    <tr><td colSpan={8} className="text-center py-10 text-gray-400">No hay bodegas registradas</td></tr>
                                )}
                                {warehouses.map((w) => {
                                    const mode: StructureMode = (w.structure_type === 'zones' || w.structure_type === 'wms') ? w.structure_type : 'simple';
                                    const isSimple = mode === 'simple';
                                    const isOpen = !isSimple && !!expanded[w.id];
                                    const tree = trees[w.id];
                                    const stats = countStats(tree);
                                    return (
                                        <WhRowGroup key={`wh-${w.id}`}>
                                            <tr className="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                                <td className="px-2">
                                                    {isSimple ? (
                                                        <span className="block w-7 h-7" />
                                                    ) : (
                                                        <button
                                                            onClick={() => toggleWh(w)}
                                                            className="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-500 dark:text-gray-300"
                                                            title={isOpen ? 'Contraer' : 'Expandir jerarquía'}
                                                        >
                                                            {isOpen ? <ChevronDownIcon className="w-5 h-5" /> : <ChevronRightIcon className="w-5 h-5" />}
                                                        </button>
                                                    )}
                                                </td>
                                                <td><span className="font-medium text-gray-900 dark:text-white">{w.name}</span></td>
                                                <td><span className="text-sm font-mono text-gray-600 dark:text-gray-300">{w.code}</span></td>
                                                <td><StructureBadge mode={mode} /></td>
                                                <td>
                                                    <div className="text-sm">
                                                        {w.address && <p className="text-gray-900 dark:text-white truncate max-w-[200px]" title={w.address}>{w.address}</p>}
                                                        {(w.city || w.state) && <p className="text-xs text-gray-500 dark:text-gray-400">{[w.city, w.state].filter(Boolean).join(', ')}</p>}
                                                        {!w.address && !w.city && !w.state && <span className="text-gray-400">—</span>}
                                                    </div>
                                                </td>
                                                <td>
                                                    {isSimple ? (
                                                        <span className="text-gray-400 text-sm">—</span>
                                                    ) : tree ? (
                                                        <div className="flex items-center gap-2 text-xs">
                                                            <StatPill color="text-indigo-600 dark:text-indigo-300 bg-indigo-100 dark:bg-indigo-900/40" label="Z" value={stats.z} />
                                                            <StatPill color="text-emerald-600 dark:text-emerald-300 bg-emerald-100 dark:bg-emerald-900/40" label={mode === 'zones' ? 'S' : 'P'} value={stats.a} />
                                                            {mode === 'wms' && <>
                                                                <StatPill color="text-purple-600 dark:text-purple-300 bg-purple-100 dark:bg-purple-900/40" label="R" value={stats.r} />
                                                                <StatPill color="text-amber-600 dark:text-amber-300 bg-amber-100 dark:bg-amber-900/40" label="N" value={stats.l} />
                                                                <StatPill color="text-rose-600 dark:text-rose-300 bg-rose-100 dark:bg-rose-900/40" label="POS" value={stats.p} />
                                                            </>}
                                                        </div>
                                                    ) : (
                                                        <button onClick={() => toggleWh(w)} className="text-xs text-gray-400 hover:text-gray-600 italic">click para cargar</button>
                                                    )}
                                                </td>
                                                <td style={{ textAlign: 'center' }}>
                                                    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${w.is_active ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>
                                                        {w.is_active ? 'Activa' : 'Inactiva'}
                                                    </span>
                                                </td>
                                                <td>
                                                    <div className="flex justify-end gap-2">
                                                        <button onClick={() => onEditWarehouse(w)} className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md" title="Editar bodega"><PencilIcon className="w-4 h-4" /></button>
                                                        <button onClick={() => setDeletingWh(w)} className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md" title="Eliminar"><TrashIcon className="w-4 h-4" /></button>
                                                    </div>
                                                </td>
                                            </tr>

                                            {isOpen && (
                                                <tr className="bg-gray-50/60 dark:bg-gray-900/40">
                                                    <td></td>
                                                    <td colSpan={7} className="p-4">
                                                        {loadingTree[w.id] ? (
                                                            <div className="flex items-center gap-2 text-sm text-gray-500"><Spinner size="sm" /> Cargando jerarquía...</div>
                                                        ) : tree ? (
                                                            <TreeInline
                                                                tree={tree}
                                                                warehouseId={w.id}
                                                                mode={mode}
                                                                expandedNodes={expandedNodes}
                                                                toggleNode={toggleNode}
                                                                onAdd={(type, parentId) => setModal({ mode: 'create', warehouseId: w.id, type, parentId })}
                                                                onEdit={(type, id, initial) => setModal({ mode: 'edit', warehouseId: w.id, type, parentId: null, initial })}
                                                                onDelete={(type, id, label) => setDeletingNode({ warehouseId: w.id, type, id, label })}
                                                            />
                                                        ) : null}
                                                    </td>
                                                </tr>
                                            )}
                                        </WhRowGroup>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {modal && (
                <HierarchyNodeModal
                    warehouseId={modal.warehouseId}
                    businessId={businessId}
                    mode={modal.mode}
                    type={modal.type}
                    parentId={(modal as any).parentId ?? null}
                    initial={(modal as any).initial}
                    onClose={() => setModal(null)}
                    onSuccess={() => { const wId = modal.warehouseId; setModal(null); refreshTree(wId); }}
                />
            )}

            {deletingWh && (
                <ConfirmModal
                    isOpen={true}
                    onClose={() => setDeletingWh(null)}
                    onConfirm={handleDeleteWh}
                    title="Eliminar bodega"
                    message={`Se eliminará la bodega "${deletingWh.name}" y toda su jerarquía. Acción irreversible.`}
                    confirmText="Eliminar"
                    type="danger"
                />
            )}

            {deletingNode && (
                <ConfirmModal
                    isOpen={true}
                    onClose={() => setDeletingNode(null)}
                    onConfirm={handleDeleteNode}
                    title={`Eliminar ${deletingNode.label}`}
                    message={`Esta acción también elimina todos sus hijos. Irreversible.`}
                    confirmText="Eliminar"
                    type="danger"
                />
            )}
        </div>
    );
}

function WhRowGroup({ children }: { children: React.ReactNode }) {
    return <>{children}</>;
}

function StructureBadge({ mode }: { mode: StructureMode }) {
    const cfg: Record<StructureMode, { label: string; cls: string }> = {
        simple: { label: 'Simple', cls: 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300' },
        zones: { label: 'Con Zonas', cls: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300' },
        wms: { label: 'WMS', cls: 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-300' },
    };
    const { label, cls } = cfg[mode];
    return <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${cls}`}>{label}</span>;
}

function StatPill({ color, label, value }: { color: string; label: string; value: number }) {
    return (
        <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full ${color} text-xs font-semibold`}>
            <span className="opacity-70">{label}</span>
            <span>{value}</span>
        </span>
    );
}

function TreeInline({ tree, warehouseId, mode, expandedNodes, toggleNode, onAdd, onEdit, onDelete }: {
    tree: WarehouseTree;
    warehouseId: number;
    mode: StructureMode;
    expandedNodes: Record<string, boolean>;
    toggleNode: (key: string) => void;
    onAdd: (type: NodeType, parentId: number | null) => void;
    onEdit: (type: NodeType, id: number, initial: Record<string, any>) => void;
    onDelete: (type: NodeKind, id: number, label: string) => void;
}) {
    const key = (prefix: string, id: number) => `${warehouseId}-${prefix}-${id}`;

    if (tree.zones.length === 0) {
        return (
            <div className="p-4 text-center bg-white dark:bg-gray-800 border border-dashed border-gray-300 dark:border-gray-700 rounded">
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-2">
                    {mode === 'wms' ? 'Sin jerarquía configurada.' : 'Sin zonas configuradas.'}
                </p>
                <button onClick={() => onAdd('zone', null)} className="inline-flex items-center gap-1 px-3 py-1.5 btn-business-primary text-white text-xs rounded-md">
                    <PlusIcon className="w-3.5 h-3.5" /> Crear primera zona
                </button>
            </div>
        );
    }

    return (
        <div className="space-y-1">
            <div className="flex justify-end mb-1">
                <button onClick={() => onAdd('zone', null)} className="inline-flex items-center gap-1 px-2.5 py-1 btn-business-primary text-white text-xs rounded-md">
                    <PlusIcon className="w-3.5 h-3.5" /> Nueva zona
                </button>
            </div>

            {tree.zones.map((z) => {
                const zKey = key('z', z.id);
                const isZones = mode === 'zones';
                const zOpen = expandedNodes[zKey] ?? false;
                return (
                    <div key={zKey} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md overflow-hidden">
                        <Row
                            depth={0}
                            color={z.color_hex || '#6366f1'}
                            badge="ZONA"
                            code={z.code}
                            name={z.name}
                            meta={z.purpose ? `· ${z.purpose}` : ''}
                            count={isZones ? `(${z.aisles.length} secc.)` : `(${z.aisles.length} pasillos)`}
                            isOpen={zOpen}
                            onToggle={() => toggleNode(zKey)}
                            onAddChild={() => onAdd('aisle', z.id)}
                            addChildLabel={isZones ? 'seccion' : 'pasillo'}
                            onEdit={() => onEdit('zone', z.id, z)}
                            onDelete={() => onDelete('zone', z.id, 'zona')}
                        />
                        {zOpen && (
                            <div className="bg-gray-50/40 dark:bg-gray-900/30">
                                {z.aisles.length === 0 && (
                                    <div className="pl-10 py-2 text-xs text-gray-400 italic">
                                        {isZones ? 'Sin secciones' : 'Sin pasillos'}
                                    </div>
                                )}
                                {z.aisles.map((a) => {
                                    const aKey = key('a', a.id);
                                    if (isZones) {
                                        return (
                                            <Row
                                                key={aKey}
                                                depth={1}
                                                badge="SEC"
                                                badgeColor="text-emerald-700 bg-emerald-100 dark:text-emerald-300 dark:bg-emerald-900/40"
                                                code={a.code}
                                                name={a.name}
                                                isOpen={null}
                                                onToggle={() => {}}
                                                onEdit={() => onEdit('aisle', a.id, a)}
                                                onDelete={() => onDelete('aisle', a.id, 'seccion')}
                                            />
                                        );
                                    }
                                    const aOpen = expandedNodes[aKey] ?? false;
                                    return (
                                        <div key={aKey}>
                                            <Row
                                                depth={1}
                                                badge="PAS"
                                                badgeColor="text-emerald-700 bg-emerald-100 dark:text-emerald-300 dark:bg-emerald-900/40"
                                                code={a.code}
                                                name={a.name}
                                                count={`(${a.racks.length} racks)`}
                                                isOpen={aOpen}
                                                onToggle={() => toggleNode(aKey)}
                                                onAddChild={() => onAdd('rack', a.id)}
                                                addChildLabel="rack"
                                                onEdit={() => onEdit('aisle', a.id, a)}
                                                onDelete={() => onDelete('aisle', a.id, 'pasillo')}
                                            />
                                            {aOpen && a.racks.map((r) => {
                                                const rKey = key('r', r.id);
                                                const rOpen = expandedNodes[rKey] ?? false;
                                                return (
                                                    <div key={rKey}>
                                                        <Row
                                                            depth={2}
                                                            badge="RCK"
                                                            badgeColor="text-purple-700 bg-purple-100 dark:text-purple-300 dark:bg-purple-900/40"
                                                            code={r.code}
                                                            name={r.name}
                                                            count={`(${r.levels.length} niveles)`}
                                                            isOpen={rOpen}
                                                            onToggle={() => toggleNode(rKey)}
                                                            onAddChild={() => onAdd('level', r.id)}
                                                            addChildLabel="nivel"
                                                            onEdit={() => onEdit('rack', r.id, r)}
                                                            onDelete={() => onDelete('rack', r.id, 'rack')}
                                                        />
                                                        {rOpen && r.levels.map((l) => {
                                                            const lKey = key('l', l.id);
                                                            const lOpen = expandedNodes[lKey] ?? false;
                                                            return (
                                                                <div key={lKey}>
                                                                    <Row
                                                                        depth={3}
                                                                        badge="LVL"
                                                                        badgeColor="text-amber-700 bg-amber-100 dark:text-amber-300 dark:bg-amber-900/40"
                                                                        code={l.code}
                                                                        name={`Ordinal ${l.ordinal}`}
                                                                        count={`(${l.positions.length} posiciones)`}
                                                                        isOpen={lOpen}
                                                                        onToggle={() => toggleNode(lKey)}
                                                                        onEdit={() => onEdit('level', l.id, l)}
                                                                        onDelete={() => onDelete('level', l.id, 'nivel')}
                                                                    />
                                                                    {lOpen && l.positions.map((p) => (
                                                                        <Row
                                                                            key={key('p', p.id)}
                                                                            depth={4}
                                                                            badge="POS"
                                                                            badgeColor="text-rose-700 bg-rose-100 dark:text-rose-300 dark:bg-rose-900/40"
                                                                            code={p.code}
                                                                            name={p.name}
                                                                            meta={p.type ? `· ${p.type}` : ''}
                                                                            isOpen={null}
                                                                            onToggle={() => {}}
                                                                            leaf
                                                                        />
                                                                    ))}
                                                                </div>
                                                            );
                                                        })}
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    );
                                })}
                            </div>
                        )}
                    </div>
                );
            })}
        </div>
    );
}

function Row({
    depth, color, badge, badgeColor, code, name, meta, count, isOpen, onToggle,
    onAddChild, addChildLabel, onEdit, onDelete, leaf,
}: {
    depth: number; color?: string; badge: string; badgeColor?: string;
    code: string; name?: string; meta?: string; count?: string;
    isOpen: boolean | null; onToggle: () => void;
    onAddChild?: () => void; addChildLabel?: string;
    onEdit?: () => void; onDelete?: () => void;
    leaf?: boolean;
}) {
    const pad = 8 + depth * 20;
    return (
        <div
            className="group flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/40 px-3 py-1.5 border-b border-gray-100 dark:border-gray-700/50 last:border-b-0"
            style={{ paddingLeft: `${pad}px` }}
        >
            <div className="flex items-center gap-2 text-sm min-w-0 flex-1">
                {isOpen !== null ? (
                    <button onClick={onToggle} className="p-0.5 rounded hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-400 flex-shrink-0">
                        {isOpen ? <ChevronDownIcon className="w-3.5 h-3.5" /> : <ChevronRightIcon className="w-3.5 h-3.5" />}
                    </button>
                ) : <span className="w-4 h-4 flex-shrink-0" />}

                {color && <span className="w-3 h-3 rounded-full flex-shrink-0" style={{ backgroundColor: color }} />}

                <span className={`text-[10px] font-semibold px-1.5 py-0.5 rounded uppercase tracking-wide flex-shrink-0 ${badgeColor || 'text-indigo-700 bg-indigo-100 dark:text-indigo-300 dark:bg-indigo-900/40'}`}>
                    {badge}
                </span>
                <span className="font-mono font-medium text-gray-900 dark:text-white">{code}</span>
                {name && <span className="text-gray-600 dark:text-gray-300 truncate">{name}</span>}
                {meta && <span className="text-xs text-gray-400 italic">{meta}</span>}
                {count && <span className="text-xs text-gray-400 ml-auto flex-shrink-0">{count}</span>}
            </div>

            <div className="flex items-center gap-1 ml-2 flex-shrink-0">
                {onAddChild && addChildLabel && (
                    <button onClick={onAddChild} className="p-1 rounded hover:bg-blue-100 dark:hover:bg-blue-900/40 text-blue-600 dark:text-blue-300" title={`Agregar ${addChildLabel}`}>
                        <PlusIcon className="w-3.5 h-3.5" />
                    </button>
                )}
                {onEdit && (
                    <button onClick={onEdit} className="p-1 rounded hover:bg-amber-100 dark:hover:bg-amber-900/40 text-amber-600 dark:text-amber-300" title="Editar">
                        <PencilIcon className="w-3.5 h-3.5" />
                    </button>
                )}
                {onDelete && (
                    <button onClick={onDelete} className="p-1 rounded hover:bg-red-100 dark:hover:bg-red-900/40 text-red-600 dark:text-red-300" title="Eliminar">
                        <TrashIcon className="w-3.5 h-3.5" />
                    </button>
                )}
            </div>
        </div>
    );
}
