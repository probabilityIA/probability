'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { ArrowLeftIcon, MapIcon } from '@heroicons/react/24/outline';
import { Spinner, Alert, Button, ConfirmModal } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { useWarehouseTree } from '@/services/modules/warehouses/ui/hooks/useWarehouseTree';
import { getWarehouseByIdAction } from '@/services/modules/warehouses/infra/actions';
import {
    deleteZoneAction,
    deleteAisleAction,
    deleteRackAction,
    deleteRackLevelAction,
} from '@/services/modules/warehouses/infra/actions/hierarchy';
import WarehouseHierarchyTree, { TreeNodeType } from '@/services/modules/warehouses/ui/components/WarehouseHierarchyTree';
import HierarchyNodeModal, { NodeType } from '@/services/modules/warehouses/ui/components/HierarchyNodeModal';
import { Warehouse } from '@/services/modules/warehouses/domain/types';

type ModalState =
    | { mode: 'create'; type: NodeType; parentId: number | null; initial?: undefined }
    | { mode: 'edit'; type: NodeType; parentId: null; initial: Record<string, any> }
    | null;

function findNode(tree: any, type: NodeType, id: number): any {
    if (!tree) return null;
    for (const z of tree.zones || []) {
        if (type === 'zone' && z.id === id) return z;
        for (const a of z.aisles || []) {
            if (type === 'aisle' && a.id === id) return a;
            for (const r of a.racks || []) {
                if (type === 'rack' && r.id === id) return r;
                for (const l of r.levels || []) {
                    if (type === 'level' && l.id === id) return l;
                }
            }
        }
    }
    return null;
}

export default function WarehouseDetailPage() {
    const params = useParams<{ id: string }>();
    const router = useRouter();
    const warehouseId = Number(params?.id);
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    const [warehouse, setWarehouse] = useState<Warehouse | null>(null);
    const [loadingWh, setLoadingWh] = useState(true);
    const [modal, setModal] = useState<ModalState>(null);
    const [confirm, setConfirm] = useState<{ type: NodeType; id: number; label: string } | null>(null);

    const { tree, loading, error, refresh } = useWarehouseTree({ warehouseId, businessId });

    useEffect(() => {
        if (!warehouseId) return;
        if (isSuperAdmin && selectedBusinessId === null) {
            setLoadingWh(false);
            return;
        }
        (async () => {
            setLoadingWh(true);
            try {
                const res = await getWarehouseByIdAction(warehouseId, businessId);
                setWarehouse(res as Warehouse);
            } catch {
                setWarehouse(null);
            } finally {
                setLoadingWh(false);
            }
        })();
    }, [warehouseId, businessId, isSuperAdmin, selectedBusinessId]);

    const handleCreateChild = (parentType: TreeNodeType | 'root', parentId: number | null) => {
        const typeMap: Record<string, NodeType> = { root: 'zone', zone: 'aisle', aisle: 'rack', rack: 'level' };
        const childType = typeMap[parentType];
        if (!childType) return;
        setModal({ mode: 'create', type: childType, parentId });
    };

    const handleEdit = (type: TreeNodeType, id: number) => {
        if (type === 'position') return;
        const node = findNode(tree, type, id);
        if (!node) return;
        setModal({ mode: 'edit', type, parentId: null, initial: node });
    };

    const handleDelete = (type: TreeNodeType, id: number) => {
        if (type === 'position') return;
        const labelMap: Record<NodeType, string> = { zone: 'zona', aisle: 'pasillo', rack: 'rack', level: 'nivel' };
        setConfirm({ type, id, label: labelMap[type] });
    };

    const confirmDelete = async () => {
        if (!confirm) return;
        const { type, id } = confirm;
        if (type === 'zone') await deleteZoneAction(id, warehouseId, businessId);
        else if (type === 'aisle') await deleteAisleAction(id, warehouseId, businessId);
        else if (type === 'rack') await deleteRackAction(id, warehouseId, businessId);
        else if (type === 'level') await deleteRackLevelAction(id, warehouseId, businessId);
        setConfirm(null);
        refresh();
    };

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    if (requiresBusinessSelection) {
        return (
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900 px-6 py-8">
                <Alert type="info">Selecciona un negocio para ver la jerarquia.</Alert>
            </div>
        );
    }

    if (loadingWh || loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
                <Spinner size="lg" />
            </div>
        );
    }

    const totalZones = tree?.zones.length || 0;
    const totalAisles = tree?.zones.reduce((a, z) => a + z.aisles.length, 0) || 0;
    const totalRacks = tree?.zones.reduce((a, z) => a + z.aisles.reduce((b, ai) => b + ai.racks.length, 0), 0) || 0;
    const totalLevels = tree?.zones.reduce((a, z) => a + z.aisles.reduce((b, ai) => b + ai.racks.reduce((c, r) => c + r.levels.length, 0), 0), 0) || 0;
    const totalPositions = tree?.zones.reduce((a, z) => a + z.aisles.reduce((b, ai) => b + ai.racks.reduce((c, r) => c + r.levels.reduce((d, l) => d + l.positions.length, 0), 0), 0), 0) || 0;

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 px-4 sm:px-6 lg:px-8 py-4 sm:py-6">
            <div className="space-y-4">
                <div className="flex items-center justify-between flex-wrap gap-2">
                    <div className="flex items-center gap-3">
                        <button
                            onClick={() => router.push('/warehouses')}
                            className="p-2 rounded-md hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-500 dark:text-gray-300"
                        >
                            <ArrowLeftIcon className="w-5 h-5" />
                        </button>
                        <div>
                            <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
                                {warehouse?.name || `Bodega #${warehouseId}`}
                            </h1>
                            {warehouse && (
                                <p className="text-sm text-gray-500 dark:text-gray-400">
                                    {warehouse.code} · {warehouse.city}
                                </p>
                            )}
                        </div>
                    </div>
                    <Button variant="outline" onClick={() => router.push(`/inventory?warehouse=${warehouseId}`)}>
                        <MapIcon className="w-4 h-4 mr-2" /> Ver stock
                    </Button>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
                    {[
                        { label: 'Zonas', value: totalZones, color: 'text-indigo-600' },
                        { label: 'Pasillos', value: totalAisles, color: 'text-emerald-600' },
                        { label: 'Racks', value: totalRacks, color: 'text-purple-600' },
                        { label: 'Niveles', value: totalLevels, color: 'text-amber-600' },
                        { label: 'Posiciones', value: totalPositions, color: 'text-rose-600' },
                    ].map((s) => (
                        <div key={s.label} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
                            <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">{s.label}</p>
                            <p className={`text-2xl font-semibold ${s.color}`}>{s.value}</p>
                        </div>
                    ))}
                </div>

                {error && <Alert type="error">{error}</Alert>}

                <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                    <h2 className="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">Jerarquia fisica</h2>
                    <WarehouseHierarchyTree
                        zones={tree?.zones || []}
                        onCreateChild={handleCreateChild}
                        onEdit={handleEdit}
                        onDelete={handleDelete}
                    />
                </div>
            </div>

            {modal && (
                <HierarchyNodeModal
                    warehouseId={warehouseId}
                    businessId={businessId}
                    mode={modal.mode}
                    type={modal.type}
                    parentId={modal.parentId}
                    initial={(modal as any).initial}
                    onClose={() => setModal(null)}
                    onSuccess={() => { setModal(null); refresh(); }}
                />
            )}

            {confirm && (
                <ConfirmModal
                    isOpen={true}
                    onClose={() => setConfirm(null)}
                    onConfirm={confirmDelete}
                    title={`Eliminar ${confirm.label}`}
                    message={`Esta accion es irreversible. Eliminar ${confirm.label} tambien eliminara todos sus hijos.`}
                    confirmText="Eliminar"
                    type="danger"
                />
            )}
        </div>
    );
}
