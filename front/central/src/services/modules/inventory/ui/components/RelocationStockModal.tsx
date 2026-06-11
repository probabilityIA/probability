'use client';

import { useState, useEffect } from 'react';
import { ChevronDown, X, Package } from 'lucide-react';
import { TransferStockDTO } from '../../domain/types';
import { transferStockAction, getProductInventoryAction } from '../../infra/actions';
import { listZonesAction, listAislesAction, listRacksAction, listRackLevelsAction } from '@/services/modules/warehouses/infra/actions/hierarchy';
import { getLocationsAction, createLocationAction } from '@/services/modules/warehouses/infra/actions';
import { Zone, Aisle, Rack, RackLevel } from '@/services/modules/warehouses/domain/hierarchy-types';
import { WarehouseLocation } from '@/services/modules/warehouses/domain/types';
import { Alert } from '@/shared/ui';

interface RelocationStockModalProps {
    productId: string;
    warehouseId: number;
    businessId?: number;
    onSuccess: () => void;
    onClose: () => void;
}

export default function RelocationStockModal({
    productId,
    warehouseId,
    businessId,
    onSuccess,
    onClose,
}: RelocationStockModalProps) {
    const [stockGeneral, setStockGeneral] = useState(0);
    const [quantity, setQuantity] = useState(0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const [zones, setZones] = useState<Zone[]>([]);
    const [selectedZoneId, setSelectedZoneId] = useState<number | null>(null);
    const [aisles, setAisles] = useState<Aisle[]>([]);
    const [selectedAisleId, setSelectedAisleId] = useState<number | null>(null);
    const [racks, setRacks] = useState<Rack[]>([]);
    const [selectedRackId, setSelectedRackId] = useState<number | null>(null);
    const [levels, setLevels] = useState<RackLevel[]>([]);
    const [selectedLevelId, setSelectedLevelId] = useState<number | null>(null);
    const [hierarchyLoading, setHierarchyLoading] = useState(false);
    const [locations, setLocations] = useState<WarehouseLocation[]>([]);

    useEffect(() => {
        (async () => {
            try {
                const inventoryLevels = await getProductInventoryAction(productId, businessId);
                const general = inventoryLevels.find(
                    (l: any) => l.warehouse_id === warehouseId && !l.location_id
                );
                setStockGeneral(general?.available_qty || 0);

                const locs = await getLocationsAction(warehouseId, businessId);
                setLocations(locs || []);

                const zonesRes = await listZonesAction(warehouseId, { page: 1, page_size: 100 }, businessId);
                setZones(zonesRes.data || []);
            } catch (err) {
                setError('Error al cargar datos');
            }
        })();
    }, [productId, warehouseId, businessId]);

    useEffect(() => {
        setSelectedAisleId(null);
        setSelectedRackId(null);
        setSelectedLevelId(null);
        setAisles([]);
        setRacks([]);
        setLevels([]);
        if (!selectedZoneId) return;
        (async () => {
            setHierarchyLoading(true);
            try {
                const aislesRes = await listAislesAction(selectedZoneId, businessId);
                setAisles(aislesRes.data || []);
            } catch {
                setAisles([]);
            } finally {
                setHierarchyLoading(false);
            }
        })();
    }, [selectedZoneId, businessId]);

    useEffect(() => {
        setSelectedRackId(null);
        setSelectedLevelId(null);
        setRacks([]);
        setLevels([]);
        if (!selectedAisleId) return;
        (async () => {
            setHierarchyLoading(true);
            try {
                const racksRes = await listRacksAction(selectedAisleId, businessId);
                setRacks(racksRes.data || []);
            } catch {
                setRacks([]);
            } finally {
                setHierarchyLoading(false);
            }
        })();
    }, [selectedAisleId, businessId]);

    useEffect(() => {
        setSelectedLevelId(null);
        setLevels([]);
        if (!selectedRackId) return;
        (async () => {
            setHierarchyLoading(true);
            try {
                const levelsRes = await listRackLevelsAction(selectedRackId, businessId);
                setLevels(levelsRes.data || []);
            } catch {
                setLevels([]);
            } finally {
                setHierarchyLoading(false);
            }
        })();
    }, [selectedRackId, businessId]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedLevelId || !quantity || quantity > stockGeneral) {
            setError('Selecciona un nivel y cantidad válida');
            return;
        }

        const selectedLevel = levels.find((l) => l.id === selectedLevelId);
        if (!selectedLevel) {
            setError('Nivel no encontrado');
            return;
        }

        const selectedRack = racks.find((r) => r.id === selectedRackId);
        const selectedAisle = aisles.find((a) => a.id === selectedAisleId);
        const selectedZone = zones.find((z) => z.id === selectedZoneId);

        if (!selectedRack || !selectedAisle || !selectedZone) {
            setError('Validación de jerarquía fallida');
            return;
        }

        let existingLocation = locations.find((loc) => loc.level_id === selectedLevelId);
        let toLocationId = existingLocation?.id;

        setLoading(true);
        setError(null);

        try {
            if (!toLocationId) {
                const locationCode = `LOC-${selectedZone.code}-${selectedAisle.code}-${selectedRack.code}-${String(selectedLevel.ordinal).padStart(2, '0')}`;
                const locationName = `${selectedZone.name} / ${selectedAisle.name} / ${selectedRack.name} / Nivel ${selectedLevel.ordinal}`;
                const newLoc = await createLocationAction(warehouseId, {
                    name: locationName,
                    code: locationCode,
                    type: 'storage',
                    level_id: selectedLevelId,
                }, businessId);
                if (!newLoc || !newLoc.id) {
                    setLoading(false);
                    setError('Error al crear ubicación para el nivel');
                    return;
                }
                toLocationId = newLoc.id;
            }

            const dto: TransferStockDTO = {
                product_id: productId,
                from_warehouse_id: warehouseId,
                to_warehouse_id: warehouseId,
                from_location_id: null,
                to_location_id: toLocationId,
                quantity: quantity,
                reason: 'Reubicación de stock general a nivel específico',
            };

            const result = await transferStockAction(dto, businessId);
            if (result.success) {
                setSuccess(`${quantity} unidades reubicadas a ${selectedZone.name} / ${selectedAisle.code} / ${selectedRack.code} / Nivel ${selectedLevel.ordinal}`);
                setTimeout(() => onSuccess(), 1500);
            } else {
                setError(result.error || 'Error al transferir stock');
            }
        } catch (err: any) {
            setError(err.message || 'Error desconocido');
        } finally {
            setLoading(false);
        }
    };

    const isFormValid = selectedLevelId && quantity > 0 && quantity <= stockGeneral;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/45 p-4" onClick={onClose}>
            <div className="bg-white dark:bg-slate-900 rounded-[18px] shadow-[0_24px_80px_-20px_rgba(15,23,42,0.45)] w-full max-w-2xl overflow-hidden flex flex-col" onClick={(e) => e.stopPropagation()}>
                <div className="flex items-start justify-between px-7 py-[22px] border-b border-slate-200 dark:border-slate-700">
                    <div className="flex items-start gap-4">
                        <div className="w-[34px] h-[34px] rounded-[10px] flex items-center justify-center flex-shrink-0 bg-amber-100 dark:bg-amber-900/20">
                            <Package className="w-5 h-5 text-amber-600 dark:text-amber-400" />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold text-slate-900 dark:text-white">Reubicar Stock</h2>
                            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">Mover stock de general a una ubicación específica en la jerarquía.</p>
                        </div>
                    </div>
                    <button
                        onClick={onClose}
                        className="p-1.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                    >
                        <X className="w-5 h-5" />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="overflow-y-auto flex-1 px-7 py-6 space-y-6">
                    {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}
                    {success && <Alert type="success" onClose={() => setSuccess(null)}>{success}</Alert>}

                    <div className="bg-amber-50 dark:bg-amber-900/10 border border-amber-200 dark:border-amber-800 rounded-lg p-4">
                        <p className="text-sm font-semibold text-amber-900 dark:text-amber-200">
                            Stock disponible en STOCK GENERAL: <span className="text-lg font-bold">{stockGeneral} unidades</span>
                        </p>
                    </div>

                    <div>
                        <label className="text-xs font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500 block mb-4 pb-3 border-b border-slate-100 dark:border-slate-800">
                            Selecciona Destino
                        </label>
                        <div className="grid gap-3" style={{ gridTemplateColumns: '1fr 1fr 1fr 1fr' }}>
                            <div>
                                <label htmlFor="zone-select" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-2">Zona</label>
                                <div className="relative">
                                    <select
                                        id="zone-select"
                                        value={selectedZoneId ?? ''}
                                        onChange={(e) => setSelectedZoneId(e.target.value ? Number(e.target.value) : null)}
                                        className="w-full px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10 appearance-none"
                                    >
                                        <option value="">Selecciona zona</option>
                                        {zones.map((z) => (
                                            <option key={z.id} value={z.id}>{z.code} - {z.name}</option>
                                        ))}
                                    </select>
                                    <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                </div>
                            </div>
                            <div>
                                <label htmlFor="aisle-select" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-2">Pasillo</label>
                                <div className="relative">
                                    <select
                                        id="aisle-select"
                                        value={selectedAisleId ?? ''}
                                        onChange={(e) => setSelectedAisleId(e.target.value ? Number(e.target.value) : null)}
                                        disabled={!selectedZoneId}
                                        className="w-full px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10 appearance-none disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        <option value="">Selecciona pasillo</option>
                                        {aisles.map((a) => (
                                            <option key={a.id} value={a.id}>{a.code} - {a.name}</option>
                                        ))}
                                    </select>
                                    <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                </div>
                            </div>
                            <div>
                                <label htmlFor="rack-select" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-2">Estantería</label>
                                <div className="relative">
                                    <select
                                        id="rack-select"
                                        value={selectedRackId ?? ''}
                                        onChange={(e) => setSelectedRackId(e.target.value ? Number(e.target.value) : null)}
                                        disabled={!selectedAisleId}
                                        className="w-full px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10 appearance-none disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        <option value="">Selecciona estantería</option>
                                        {racks.map((r) => (
                                            <option key={r.id} value={r.id}>{r.code} - {r.name}</option>
                                        ))}
                                    </select>
                                    <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                </div>
                            </div>
                            <div>
                                <label htmlFor="level-select" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-2">Nivel</label>
                                <div className="relative">
                                    <select
                                        id="level-select"
                                        value={selectedLevelId ?? ''}
                                        onChange={(e) => setSelectedLevelId(e.target.value ? Number(e.target.value) : null)}
                                        disabled={!selectedRackId}
                                        className="w-full px-3 py-2.5 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10 appearance-none disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        <option value="">Selecciona nivel</option>
                                        {levels.map((l) => (
                                            <option key={l.id} value={l.id}>Nivel {l.ordinal}</option>
                                        ))}
                                    </select>
                                    <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                                </div>
                            </div>
                        </div>
                        {hierarchyLoading && <p className="text-xs text-slate-400 animate-pulse mt-2">Cargando opciones...</p>}
                    </div>

                    <div>
                        <label htmlFor="quantity" className="text-xs font-semibold text-slate-600 dark:text-slate-400 block mb-2">Cantidad a Reubicar <span className="text-rose-500">*</span></label>
                        <input
                            id="quantity"
                            type="number"
                            min="1"
                            max={stockGeneral}
                            value={quantity}
                            onChange={(e) => setQuantity(Math.max(0, Math.min(stockGeneral, parseInt(e.target.value) || 0)))}
                            className="w-full px-4 py-3 bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-[1.5px] border-slate-200 dark:border-slate-700 rounded-[10px] text-sm focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-500/10"
                        />
                        <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">Máximo: {stockGeneral} unidades</p>
                    </div>
                </form>

                <div className="flex items-center justify-end gap-3 px-7 py-4 border-t border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/50">
                    <button
                        onClick={onClose}
                        disabled={loading}
                        className="px-4 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 bg-white dark:bg-slate-700 border border-slate-200 dark:border-slate-600 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-600 transition-colors disabled:opacity-50"
                    >
                        Cancelar
                    </button>
                    <button
                        onClick={handleSubmit}
                        disabled={!isFormValid || loading}
                        className="px-6 py-2 text-sm font-semibold text-white bg-amber-600 hover:bg-amber-700 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {loading ? 'Reubiando...' : 'Reubicar Stock'}
                    </button>
                </div>
            </div>
        </div>
    );
}
