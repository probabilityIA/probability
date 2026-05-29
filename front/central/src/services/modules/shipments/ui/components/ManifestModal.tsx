'use client';

import { useEffect, useMemo, useState } from 'react';
import { X, FileText, RefreshCw, Truck, Download, Search } from 'lucide-react';
import { getManifestPendingAction, generateManifestPdfAction, ManifestGroup, ManifestPendingShipment } from '../../infra/actions';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';

interface ManifestModalProps {
    isOpen: boolean;
    onClose: () => void;
    businessId: number | null;
}

export function ManifestModal({ isOpen, onClose, businessId }: ManifestModalProps) {
    const [loading, setLoading] = useState(false);
    const [groups, setGroups] = useState<ManifestGroup[]>([]);
    const [selected, setSelected] = useState<Set<number>>(new Set());
    const [carrierFilter, setCarrierFilter] = useState<string>('');
    const [search, setSearch] = useState('');
    const [generating, setGenerating] = useState(false);

    useEffect(() => {
        if (!isOpen || !businessId) return;
        setLoading(true);
        setSelected(new Set());
        getManifestPendingAction(businessId)
            .then((res) => {
                if (res.success) setGroups(res.data || []);
                else setGroups([]);
            })
            .finally(() => setLoading(false));
    }, [isOpen, businessId]);

    const filteredGroups = useMemo(() => {
        return groups
            .filter((g) => !carrierFilter || g.carrier === carrierFilter)
            .map((g) => ({
                ...g,
                shipments: g.shipments.filter((s) => {
                    if (!search.trim()) return true;
                    const q = search.trim().toLowerCase();
                    return (
                        s.order_number?.toLowerCase().includes(q) ||
                        s.tracking_number?.toLowerCase().includes(q) ||
                        s.customer_name?.toLowerCase().includes(q) ||
                        s.destination_city?.toLowerCase().includes(q)
                    );
                }),
            }))
            .filter((g) => g.shipments.length > 0);
    }, [groups, carrierFilter, search]);

    const allVisible = useMemo(() => filteredGroups.flatMap((g) => g.shipments), [filteredGroups]);
    const allSelected = allVisible.length > 0 && allVisible.every((s) => selected.has(s.shipment_id));

    const toggle = (id: number) => {
        const next = new Set(selected);
        if (next.has(id)) next.delete(id);
        else next.add(id);
        setSelected(next);
    };

    const toggleAll = () => {
        if (allSelected) {
            const next = new Set(selected);
            for (const s of allVisible) next.delete(s.shipment_id);
            setSelected(next);
        } else {
            const next = new Set(selected);
            for (const s of allVisible) next.add(s.shipment_id);
            setSelected(next);
        }
    };

    const toggleCarrierGroup = (g: ManifestGroup) => {
        const ids = g.shipments.map((s) => s.shipment_id);
        const allOn = ids.every((id) => selected.has(id));
        const next = new Set(selected);
        if (allOn) ids.forEach((id) => next.delete(id));
        else ids.forEach((id) => next.add(id));
        setSelected(next);
    };

    const handleGenerate = async () => {
        if (!businessId || selected.size === 0) return;
        setGenerating(true);
        try {
            const res = await generateManifestPdfAction(businessId, Array.from(selected), carrierFilter || undefined);
            if (res.success && res.blob) {
                const a = document.createElement('a');
                a.href = res.blob;
                a.download = res.filename || 'manifiesto.pdf';
                document.body.appendChild(a);
                a.click();
                a.remove();
            } else {
                alert(res.message || 'Error al generar PDF');
            }
        } finally {
            setGenerating(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-5xl max-h-[90vh] flex flex-col">
                <div className="flex items-center justify-between px-5 py-4 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-indigo-50 to-purple-50 dark:from-indigo-950/40 dark:to-purple-950/40 rounded-t-2xl">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-indigo-100 dark:bg-indigo-900/40 flex items-center justify-center">
                            <FileText size={20} className="text-indigo-600 dark:text-indigo-300" />
                        </div>
                        <div>
                            <h3 className="text-lg font-bold text-gray-900 dark:text-white">Manifiesto de Recoleccion</h3>
                            <p className="text-xs text-gray-500 dark:text-gray-400">Selecciona los envios pendientes a entregar al transportador</p>
                        </div>
                    </div>
                    <button onClick={onClose} className="p-2 rounded-lg hover:bg-white/60 dark:hover:bg-gray-700 text-gray-500">
                        <X size={18} />
                    </button>
                </div>

                <div className="px-5 py-3 border-b border-gray-100 dark:border-gray-700 flex items-center gap-3 flex-wrap">
                    <div className="relative flex-1 min-w-[200px]">
                        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                        <input
                            type="text"
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            placeholder="Buscar por orden, guia, cliente o ciudad..."
                            className="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 dark:border-gray-600 rounded-lg bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white"
                        />
                    </div>
                    <select
                        value={carrierFilter}
                        onChange={(e) => setCarrierFilter(e.target.value)}
                        className="px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-200 min-w-[180px]"
                    >
                        <option value="">Todas las transportadoras</option>
                        {groups.map((g) => (
                            <option key={g.carrier} value={g.carrier}>{g.carrier} ({g.count})</option>
                        ))}
                    </select>
                    <div className="text-xs font-semibold text-gray-600 dark:text-gray-300 ml-auto">
                        Seleccionados: <span className="text-indigo-600 dark:text-indigo-300">{selected.size}</span>
                    </div>
                </div>

                <div className="flex-1 overflow-y-auto px-5 py-3">
                    {loading ? (
                        <div className="flex items-center justify-center py-16 text-gray-400 gap-2">
                            <RefreshCw size={18} className="animate-spin" /> Cargando...
                        </div>
                    ) : !businessId ? (
                        <div className="text-center py-16 text-gray-400 text-sm">Selecciona un negocio para ver los envios pendientes.</div>
                    ) : filteredGroups.length === 0 ? (
                        <div className="text-center py-16 text-gray-400 text-sm">No hay envios pendientes de recoleccion.</div>
                    ) : (
                        <div className="space-y-4">
                            <label className="flex items-center gap-2 px-3 py-2 bg-gray-50 dark:bg-gray-700/50 rounded-lg cursor-pointer">
                                <input type="checkbox" checked={allSelected} onChange={toggleAll} className="rounded" />
                                <span className="text-sm font-semibold text-gray-700 dark:text-gray-200">Seleccionar todos los visibles ({allVisible.length})</span>
                            </label>

                            {filteredGroups.map((g) => {
                                const logo = getCarrierLogo(g.carrier);
                                const groupIds = g.shipments.map((s) => s.shipment_id);
                                const allOn = groupIds.every((id) => selected.has(id));
                                return (
                                    <div key={g.carrier} className="border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
                                        <div className="flex items-center justify-between px-4 py-2 bg-gray-50 dark:bg-gray-700/40">
                                            <label className="flex items-center gap-3 cursor-pointer">
                                                <input type="checkbox" checked={allOn} onChange={() => toggleCarrierGroup(g)} className="rounded" />
                                                {logo ? (
                                                    <img src={logo} alt={g.carrier} className="h-6 w-auto" />
                                                ) : (
                                                    <Truck size={16} className="text-gray-400" />
                                                )}
                                                <span className="font-bold text-sm text-gray-800 dark:text-gray-100">{g.carrier}</span>
                                                <span className="text-xs text-gray-500">({g.shipments.length})</span>
                                            </label>
                                        </div>
                                        <div className="divide-y divide-gray-100 dark:divide-gray-700">
                                            {g.shipments.map((s) => (
                                                <label key={s.shipment_id} className="flex items-center gap-3 px-4 py-2 hover:bg-gray-50 dark:hover:bg-gray-700/30 cursor-pointer text-sm">
                                                    <input type="checkbox" checked={selected.has(s.shipment_id)} onChange={() => toggle(s.shipment_id)} className="rounded" />
                                                    <span className="font-mono text-xs text-gray-500 dark:text-gray-400 w-28 truncate">{s.tracking_number || '—'}</span>
                                                    <span className="text-purple-700 dark:text-purple-300 font-semibold text-xs w-24 truncate">{s.order_number}</span>
                                                    <span className="flex-1 truncate text-gray-700 dark:text-gray-200">{s.customer_name || '—'}</span>
                                                    <span className="text-gray-500 dark:text-gray-400 text-xs w-32 truncate">{s.destination_city}</span>
                                                </label>
                                            ))}
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    )}
                </div>

                <div className="px-5 py-3 border-t border-gray-100 dark:border-gray-700 flex items-center justify-end gap-2 bg-gray-50 dark:bg-gray-800/60 rounded-b-2xl">
                    <button onClick={onClose} className="px-4 py-2 text-sm font-semibold text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg">
                        Cancelar
                    </button>
                    <button
                        onClick={handleGenerate}
                        disabled={selected.size === 0 || generating}
                        className="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 rounded-lg flex items-center gap-2"
                    >
                        {generating ? <RefreshCw size={14} className="animate-spin" /> : <Download size={14} />}
                        Generar PDF
                    </button>
                </div>
            </div>
        </div>
    );
}
