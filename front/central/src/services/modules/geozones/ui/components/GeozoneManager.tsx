'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import dynamic from 'next/dynamic';
import { PlusIcon, MagnifyingGlassIcon, ArrowPathIcon, GlobeAltIcon, SparklesIcon, ChevronRightIcon, HomeIcon } from '@heroicons/react/24/outline';
import { Geozone, DrillState, DisplayFeature } from '../../domain/types';
import { listGeozonesAction, deleteGeozoneAction } from '../../infra/actions';
import { Alert, Button, Spinner } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import GeozoneList from './GeozoneList';
import GeozoneForm from './GeozoneForm';

const GeozoneMap = dynamic(() => import('./GeozoneMap'), { ssr: false, loading: () => (
    <div className="h-[600px] flex items-center justify-center bg-gray-50 dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700">
        <Spinner size="lg" />
    </div>
)});

interface GeozoneManagerProps { selectedBusinessId?: number | null }

export default function GeozoneManager({ selectedBusinessId = null }: GeozoneManagerProps) {
    const { isSuperAdmin } = usePermissions();

    const [drill, setDrill] = useState<DrillState>({ level: 'country' });
    const [items, setItems] = useState<Geozone[]>([]);
    const [customItems, setCustomItems] = useState<Geozone[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');
    const [showCustom, setShowCustom] = useState(false);
    const [showForm, setShowForm] = useState(false);
    const [selectedId, setSelectedId] = useState<number | null>(null);
    const [loadMs, setLoadMs] = useState<number | null>(null);
    const [bytesKB, setBytesKB] = useState<number | null>(null);

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    const drillConfig = useMemo(() => {
        switch (drill.level) {
            case 'country':
                return { type: 'state', parentId: undefined as number | undefined, fetchZoom: 7, label: 'Departamentos' };
            case 'state':
                return { type: 'city', parentId: drill.state?.id, fetchZoom: 9, label: `Municipios de ${drill.state?.name}` };
            case 'city':
                return { type: '', parentId: drill.city?.id, fetchZoom: 11, label: `Zonas de ${drill.city?.name}` };
            case 'admin_district':
                return { type: 'neighborhood', parentId: drill.adminDistrict?.id, fetchZoom: 13, label: `UPZ de ${drill.adminDistrict?.name}` };
            case 'neighborhood':
                return { type: 'barrio', parentId: drill.neighborhood?.id, fetchZoom: 14, label: `Barrios de ${drill.neighborhood?.name}` };
        }
    }, [drill]);

    const fetchDisplay = useCallback(async () => {
        setLoading(true);
        setError(null);
        const t0 = performance.now();
        try {
            const sp = new URLSearchParams({ zoom: String(drillConfig.fetchZoom), type: drillConfig.type });
            if (drillConfig.parentId) sp.append('parent_id', String(drillConfig.parentId));
            const r = await fetch(`/api/geozones-display?${sp.toString()}`, { cache: 'no-cache' });
            if (!r.ok) throw new Error(`HTTP ${r.status}`);
            const fc = await r.json();
            const features = fc.features || [];
            const mapped: Geozone[] = features.map((f: DisplayFeature) => ({
                id: f.properties.id,
                business_id: 0,
                parent_id: drillConfig.parentId ?? null,
                type: f.properties.type,
                code: f.properties.code ?? null,
                name: f.properties.name,
                geometry: f.geometry,
                centroid: null,
                properties: {},
                is_active: true,
            }));
            setItems(mapped);
            setLoadMs(Math.round(performance.now() - t0));
            try { setBytesKB(Math.round(JSON.stringify(fc).length / 1024)); } catch {}
        } catch (err: any) {
            setError(err.message || 'Error al cargar geozonas');
        } finally {
            setLoading(false);
        }
    }, [drillConfig.type, drillConfig.parentId, drillConfig.fetchZoom]);

    const fetchCustom = useCallback(async () => {
        const businessId = isSuperAdmin && selectedBusinessId ? selectedBusinessId : undefined;
        if (!businessId || !showCustom) { setCustomItems([]); return; }
        try {
            const resp = await listGeozonesAction({
                page: 1, page_size: 100, include_geometry: true, business_id: businessId,
            });
            const onlyCustom = (resp.data || []).filter((g) => g.business_id !== 0);
            setCustomItems(onlyCustom);
        } catch {}
    }, [isSuperAdmin, selectedBusinessId, showCustom]);

    useEffect(() => { if (!requiresBusinessSelection) fetchDisplay(); }, [fetchDisplay, requiresBusinessSelection]);
    useEffect(() => { if (!requiresBusinessSelection) fetchCustom(); }, [fetchCustom, requiresBusinessSelection]);

    const filtered = useMemo(() => {
        let combined: Geozone[] = items.slice();
        if (showCustom) combined = combined.concat(customItems);
        if (search) {
            const s = search.toLowerCase();
            combined = combined.filter((g) => g.name.toLowerCase().includes(s));
        }
        return combined;
    }, [items, customItems, search, showCustom]);

    const handlePolygonClick = useCallback((g: Geozone) => {
        if (g.business_id !== 0) { setSelectedId(g.id); return; }
        if (drill.level === 'country' && g.type === 'state') {
            setDrill({ level: 'state', state: { id: g.id, name: g.name } });
            setSelectedId(null);
        } else if (drill.level === 'state' && g.type === 'city') {
            setDrill({ level: 'city', state: drill.state, city: { id: g.id, name: g.name } });
            setSelectedId(null);
        } else if (drill.level === 'city' && g.type === 'admin_district') {
            setDrill({ level: 'admin_district', state: drill.state, city: drill.city, adminDistrict: { id: g.id, name: g.name } });
            setSelectedId(null);
        } else if (drill.level === 'admin_district' && g.type === 'neighborhood') {
            setDrill({ level: 'neighborhood', state: drill.state, city: drill.city, adminDistrict: drill.adminDistrict, neighborhood: { id: g.id, name: g.name } });
            setSelectedId(null);
        } else {
            setSelectedId(g.id);
        }
    }, [drill]);

    const goToCountry = () => { setDrill({ level: 'country' }); setSelectedId(null); };
    const goToState = () => { if (drill.state) setDrill({ level: 'state', state: drill.state }); setSelectedId(null); };
    const goToCity = () => { if (drill.state && drill.city) setDrill({ level: 'city', state: drill.state, city: drill.city }); setSelectedId(null); };
    const goToAdminDistrict = () => { if (drill.adminDistrict) setDrill({ level: 'admin_district', state: drill.state, city: drill.city, adminDistrict: drill.adminDistrict }); setSelectedId(null); };

    const handleDelete = async (g: Geozone) => {
        if (g.business_id === 0) { alert('No puedes eliminar geozonas oficiales DANE'); return; }
        if (!confirm(`Eliminar "${g.name}"?`)) return;
        try {
            await deleteGeozoneAction(g.id, isSuperAdmin && selectedBusinessId ? selectedBusinessId : undefined);
            fetchCustom();
        } catch (err: any) {
            setError(err.message || 'Error al eliminar');
        }
    };

    const handleSearch = (e: React.FormEvent) => { e.preventDefault(); setSearch(searchInput); };

    if (requiresBusinessSelection) {
        return (
            <div className="flex flex-col items-center justify-center py-20 text-center">
                <GlobeAltIcon className="w-16 h-16 text-gray-300 mb-4" />
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-1">Selecciona un negocio</h2>
                <p className="text-sm text-gray-500 dark:text-gray-400">Elige un negocio en el navbar para ver y gestionar sus geozonas</p>
            </div>
        );
    }

    return (
        <div className="space-y-5">
            {/* Header */}
            <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-purple-600 via-pink-500 to-orange-400 p-6 shadow-xl">
                <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_right,rgba(255,255,255,0.2),transparent)]" />
                <div className="relative flex items-start justify-between gap-4 flex-wrap">
                    <div>
                        <h1 className="text-3xl font-bold text-white flex items-center gap-2">
                            <span>Geozonas</span>
                            <SparklesIcon className="w-7 h-7" />
                        </h1>
                        <p className="text-white/90 text-sm mt-1 max-w-xl">
                            Navegacion jerarquica: click en un departamento para ver sus municipios, click en un municipio para ver sus corregimientos.
                        </p>
                    </div>
                    <button
                        onClick={() => setShowForm(true)}
                        className="inline-flex items-center gap-1 px-4 py-2 bg-white text-purple-700 hover:bg-purple-50 rounded-lg shadow-md hover:shadow-lg transition-all text-sm font-semibold"
                    >
                        <PlusIcon className="w-4 h-4" />
                        Nueva geozona
                    </button>
                </div>
            </div>

            {/* Breadcrumb */}
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 px-4 py-3 flex items-center gap-2 text-sm">
                <button
                    onClick={goToCountry}
                    className={`flex items-center gap-1 px-2 py-1 rounded transition-colors ${drill.level === 'country' ? 'text-purple-700 dark:text-purple-300 font-semibold' : 'text-gray-600 dark:text-gray-300 hover:text-purple-700 dark:hover:text-purple-300'}`}
                >
                    <HomeIcon className="w-4 h-4" />
                    Colombia
                </button>
                {drill.level !== 'country' && drill.state && (
                    <>
                        <ChevronRightIcon className="w-4 h-4 text-gray-400" />
                        <button
                            onClick={goToState}
                            className={`px-2 py-1 rounded transition-colors ${drill.level === 'state' ? 'text-purple-700 dark:text-purple-300 font-semibold' : 'text-gray-600 dark:text-gray-300 hover:text-purple-700 dark:hover:text-purple-300'}`}
                        >
                            {drill.state.name}
                        </button>
                    </>
                )}
                {(drill.level === 'city' || drill.level === 'admin_district') && drill.city && (
                    <>
                        <ChevronRightIcon className="w-4 h-4 text-gray-400" />
                        <button
                            onClick={goToCity}
                            className={`px-2 py-1 rounded transition-colors ${drill.level === 'city' ? 'text-purple-700 dark:text-purple-300 font-semibold' : 'text-gray-600 dark:text-gray-300 hover:text-purple-700 dark:hover:text-purple-300'}`}
                        >
                            {drill.city.name}
                        </button>
                    </>
                )}
                {(drill.level === 'admin_district' || drill.level === 'neighborhood') && drill.adminDistrict && (
                    <>
                        <ChevronRightIcon className="w-4 h-4 text-gray-400" />
                        <button
                            onClick={goToAdminDistrict}
                            className={`px-2 py-1 rounded transition-colors ${drill.level === 'admin_district' ? 'text-purple-700 dark:text-purple-300 font-semibold' : 'text-gray-600 dark:text-gray-300 hover:text-purple-700 dark:hover:text-purple-300'}`}
                        >
                            {drill.adminDistrict.name}
                        </button>
                    </>
                )}
                {drill.level === 'neighborhood' && drill.neighborhood && (
                    <>
                        <ChevronRightIcon className="w-4 h-4 text-gray-400" />
                        <span className="px-2 py-1 text-purple-700 dark:text-purple-300 font-semibold">{drill.neighborhood.name}</span>
                    </>
                )}
                <span className="ml-auto text-xs text-gray-500 dark:text-gray-400">
                    Mostrando: <b className="text-gray-900 dark:text-white">{drillConfig.label}</b>
                </span>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                <StatCard label="Visibles" value={filtered.length} accent="from-blue-500 to-cyan-500" />
                <StatCard label="Personalizadas" value={showCustom ? customItems.length : '—'} accent="from-pink-500 to-rose-500" />
                <StatCard
                    label={bytesKB !== null ? `Payload` : 'Cargando'}
                    value={bytesKB !== null ? `${bytesKB} KB` : '...'}
                    sub={loadMs !== null ? `${loadMs} ms` : undefined}
                    accent="from-emerald-500 to-teal-500"
                />
                <StatCard label="Nivel" value={drill.level} accent="from-violet-500 to-purple-500" />
            </div>

            {/* Filtros */}
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-4 space-y-3">
                <div className="flex items-center gap-3 flex-wrap">
                    <label className="inline-flex items-center gap-2 cursor-pointer">
                        <input
                            type="checkbox"
                            checked={showCustom}
                            onChange={(e) => setShowCustom(e.target.checked)}
                            className="w-4 h-4 text-purple-600 rounded"
                        />
                        <span className="text-sm text-gray-700 dark:text-gray-200">Mostrar mis zonas personalizadas</span>
                    </label>
                </div>

                <form onSubmit={handleSearch} className="flex gap-2">
                    <div className="relative flex-1">
                        <MagnifyingGlassIcon className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                        <input
                            value={searchInput}
                            onChange={(e) => setSearchInput(e.target.value)}
                            placeholder={`Buscar en ${drillConfig.label.toLowerCase()}...`}
                            className="w-full pl-9 pr-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                        />
                    </div>
                    <Button variant="purple" type="submit">Buscar</Button>
                    <Button variant="secondary" type="button" onClick={() => { fetchDisplay(); fetchCustom(); }} title="Recargar">
                        <ArrowPathIcon className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                    </Button>
                </form>
            </div>

            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            {/* Mapa + Lista */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
                <div className="lg:col-span-2">
                    {loading && filtered.length === 0 ? (
                        <div className="h-[600px] flex items-center justify-center bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700">
                            <Spinner size="lg" />
                        </div>
                    ) : (
                        <GeozoneMap
                            geozones={filtered}
                            selectedId={selectedId}
                            onSelect={handlePolygonClick}
                            fitKey={`${drill.level}-${drill.state?.id || 0}-${drill.city?.id || 0}-${drill.adminDistrict?.id || 0}-${drill.neighborhood?.id || 0}`}
                        />
                    )}
                </div>
                <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm overflow-hidden">
                    <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50 flex items-center justify-between">
                        <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                            {filtered.length} {drill.level === 'country' ? 'departamento' : drill.level === 'state' ? 'municipio' : drill.level === 'admin_district' ? 'UPZ' : drill.level === 'neighborhood' ? 'barrio' : 'zona'}{filtered.length !== 1 ? 's' : ''}
                        </h3>
                        {drill.level !== 'country' && (
                            <button
                                onClick={drill.level === 'state' ? goToCountry : drill.level === 'city' ? goToState : drill.level === 'admin_district' ? goToCity : goToAdminDistrict}
                                className="text-xs text-purple-600 dark:text-purple-300 hover:underline"
                            >
                                ← Atras
                            </button>
                        )}
                    </div>
                    <div className="max-h-[540px] overflow-y-auto">
                        <GeozoneList
                            items={filtered}
                            selectedId={selectedId}
                            onSelect={handlePolygonClick}
                            onDelete={handleDelete}
                            canDelete={(g) => g.business_id !== 0}
                        />
                    </div>
                </div>
            </div>

            {/* Modal Form */}
            {showForm && (
                <div className="fixed inset-0 z-[2000] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 overflow-y-auto">
                    <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-4xl max-h-[95vh] overflow-y-auto">
                        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 sticky top-0 bg-white dark:bg-gray-800 z-10">
                            <div>
                                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Nueva geozona</h2>
                                <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Dibuja sobre el mapa o pega un GeoJSON</p>
                            </div>
                            <button onClick={() => setShowForm(false)} className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 text-2xl leading-none">&times;</button>
                        </div>
                        <div className="p-6">
                            <GeozoneForm
                                onSuccess={() => { setShowForm(false); fetchCustom(); }}
                                onCancel={() => setShowForm(false)}
                                businessId={isSuperAdmin && selectedBusinessId ? selectedBusinessId : undefined}
                                contextLayers={items.slice(0, 50)}
                            />
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

function StatCard({ label, value, sub, accent }: { label: string; value: number | string; sub?: string; accent: string }) {
    return (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-4 shadow-sm relative overflow-hidden">
            <div className={`absolute -right-4 -top-4 w-20 h-20 rounded-full bg-gradient-to-br ${accent} opacity-15`} />
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">{label}</p>
            <p className="text-2xl font-bold text-gray-900 dark:text-white">{value}</p>
            {sub && <p className="text-[10px] text-gray-500 dark:text-gray-400 mt-0.5">{sub}</p>}
        </div>
    );
}
