'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { MapContainer, TileLayer, GeoJSON as GeoJSONLayer, useMap, useMapEvents } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { getGeozonesForDisplayAction } from '@/services/modules/geozones/infra/actions';
import { DisplayFeatureCollection, DisplayFeature } from '@/services/modules/geozones/domain/types';
import { Spinner } from '@/shared/ui';

type MetricType = 'orders' | 'percentage';
type DrillLevel = 'state' | 'city' | 'admin_district' | 'neighborhood';

interface LocationData {
    name: string;
    fullName: string;
    value: number;
}

interface DashboardInteractiveMapProps {
    ordersByDepartment: { name: string; value: number }[];
    ordersByLocation: LocationData[];
    selectedBusinessId?: number;
    height?: number;
}

const DRILL_CONFIG: Record<DrillLevel, { type: string; zoomLevel: number; nextLevel: DrillLevel | null }> = {
    state: { type: 'state', zoomLevel: 7, nextLevel: 'city' },
    city: { type: 'city', zoomLevel: 9, nextLevel: 'admin_district' },
    admin_district: { type: 'admin_district', zoomLevel: 11, nextLevel: 'neighborhood' },
    neighborhood: { type: 'neighborhood', zoomLevel: 13, nextLevel: null },
};

function normalizeDepartment(dept: string): string {
    const normalized = dept.toUpperCase().trim();
    if (normalized.includes('D.C') || normalized.includes('D.D') || normalized.includes('S.C') ||
        normalized === 'BOGOTA' || normalized === 'BOGOTA' || normalized === 'DC' || normalized === 'DD' || normalized === 'SC') {
        return 'BOGOTA';
    }
    return normalized;
}

function FitBounds({ features, fitKey }: { features: DisplayFeature[]; fitKey: string }) {
    const map = useMap();
    const lastSig = useRef<string>('');

    useEffect(() => {
        const ids = features.map((f) => f.properties.id).sort((a, b) => a - b).join(',');
        const sig = `${fitKey}|${ids}`;
        if (lastSig.current === sig || features.length === 0) return;

        const layers: L.Layer[] = [];
        features.forEach((f) => {
            if (f.geometry) {
                try {
                    layers.push(L.geoJSON(f.geometry as any));
                } catch {}
            }
        });

        if (layers.length === 0) return;
        const group = L.featureGroup(layers);
        const bounds = group.getBounds();
        if (bounds.isValid()) {
            map.fitBounds(bounds, { padding: [20, 20], maxZoom: 13 });
            lastSig.current = sig;
        }
    }, [features, map, fitKey]);

    return null;
}

function ZoomReporter({ onZoomChange }: { onZoomChange: (z: number) => void }) {
    const map = useMapEvents({
        zoomend: () => onZoomChange(map.getZoom()),
    });

    useEffect(() => {
        onZoomChange(map.getZoom());
    }, [map, onZoomChange]);

    return null;
}

export default function DashboardInteractiveMap({
    ordersByDepartment,
    ordersByLocation,
    selectedBusinessId,
    height = 600,
}: DashboardInteractiveMapProps) {
    const [drillLevel, setDrillLevel] = useState<DrillLevel>('state');
    const [breadcrumb, setBreadcrumb] = useState<{ id: number; name: string; type: string }[]>([]);
    const [geojsonData, setGeojsonData] = useState<DisplayFeatureCollection | null>(null);
    const [loading, setLoading] = useState(true);
    const [metricType, setMetricType] = useState<MetricType>('orders');
    const [zoom, setZoom] = useState(7);

    const ordersMap = useMemo(() => {
        if (drillLevel === 'state') {
            const map = new Map<string, number>();
            ordersByDepartment.forEach((item) => {
                const normalized = normalizeDepartment(item.name);
                map.set(normalized, item.value);
            });
            return map;
        }

        if (drillLevel === 'city' && breadcrumb.length > 0) {
            const parentDept = breadcrumb[breadcrumb.length - 1].name;
            const map = new Map<string, number>();
            ordersByLocation.forEach((item) => {
                const parts = item.fullName.split(', ');
                const state = parts[parts.length - 1] || '';
                if (normalizeDepartment(state) === normalizeDepartment(parentDept)) {
                    const cityNorm = item.name.toUpperCase().trim();
                    map.set(cityNorm, item.value);
                }
            });
            return map;
        }

        return new Map<string, number>();
    }, [drillLevel, breadcrumb, ordersByDepartment, ordersByLocation]);

    const totalOrders = useMemo(() => {
        return Array.from(ordersMap.values()).reduce((sum, val) => sum + val, 0) || 1;
    }, [ordersMap]);

    const getDensityColor = (featureName: string, metric: MetricType): string => {
        const normalized = featureName.toUpperCase().trim();
        const count = ordersMap.get(normalized) || 0;

        if (count === 0) return '#d1d5db';

        if (metric === 'percentage') {
            const percentage = (count / totalOrders) * 100;
            if (percentage > 66) return '#16a34a';
            if (percentage > 33) return '#ca8a04';
            return '#dc2626';
        }

        const maxCount = Math.max(...Array.from(ordersMap.values()));
        const ratio = maxCount > 0 ? count / maxCount : 0;
        if (ratio > 0.66) return '#16a34a';
        if (ratio > 0.33) return '#ca8a04';
        return '#dc2626';
    };

    const fetchGeozones = async (level: DrillLevel, parentId?: number) => {
        try {
            setLoading(true);
            const config = DRILL_CONFIG[level];
            const response = await getGeozonesForDisplayAction(
                config.type as any,
                config.zoomLevel,
                undefined,
                parentId || undefined
            );
            setGeojsonData(response);
            setDrillLevel(level);
        } catch (error) {
            console.error('Error loading geozones:', error);
            setGeojsonData(null);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchGeozones('state');
    }, [selectedBusinessId]);

    const handlePolygonClick = (feature: DisplayFeature) => {
        const nextLevel = DRILL_CONFIG[drillLevel].nextLevel;
        if (!nextLevel) return;

        const newBreadcrumb = [...breadcrumb, { id: feature.properties.id, name: feature.properties.name, type: drillLevel }];
        setBreadcrumb(newBreadcrumb);
        fetchGeozones(nextLevel, feature.properties.id);
    };

    const handleBreadcrumbClick = (index: number) => {
        if (index === -1) {
            setBreadcrumb([]);
            fetchGeozones('state');
            return;
        }

        const newBreadcrumb = breadcrumb.slice(0, index + 1);
        setBreadcrumb(newBreadcrumb);

        const parentId = newBreadcrumb.length > 0 ? newBreadcrumb[newBreadcrumb.length - 1].id : undefined;
        const levels: DrillLevel[] = ['state', 'city', 'admin_district', 'neighborhood'];
        const nextLevel = levels[newBreadcrumb.length] as DrillLevel;

        fetchGeozones(nextLevel, parentId);
    };

    const features = useMemo(() => {
        if (!geojsonData) return [];
        return geojsonData.features.filter((f) => !!f.geometry);
    }, [geojsonData]);

    const departmentOptions = useMemo(() => {
        if (!geojsonData || drillLevel !== 'state') return [];
        return geojsonData.features
            .map((f) => ({ id: f.properties.id, name: f.properties.name }))
            .sort((a, b) => a.name.localeCompare(b.name));
    }, [geojsonData, drillLevel]);

    const cityOptions = useMemo(() => {
        if (!geojsonData || drillLevel !== 'city') return [];
        return geojsonData.features
            .map((f) => ({ id: f.properties.id, name: f.properties.name }))
            .sort((a, b) => a.name.localeCompare(b.name));
    }, [geojsonData, drillLevel]);

    return (
        <div className="space-y-4">
            {/* Header con titulo y selector de metrica */}
            <div className="flex items-center justify-between">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Mapa Interactivo de Órdenes</h2>
                <div className="flex gap-2">
                    <button
                        onClick={() => setMetricType('orders')}
                        className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
                            metricType === 'orders'
                                ? 'bg-purple-600 text-white'
                                : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-200'
                        }`}
                    >
                        Órdenes
                    </button>
                    <button
                        onClick={() => setMetricType('percentage')}
                        className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
                            metricType === 'percentage'
                                ? 'bg-purple-600 text-white'
                                : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-200'
                        }`}
                    >
                        Porcentaje
                    </button>
                </div>
            </div>

            {/* Filtros por nivel */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                {/* Colombia */}
                <div>
                    <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">País</label>
                    <button
                        onClick={() => handleBreadcrumbClick(-1)}
                        className={`w-full px-3 py-2 rounded-lg text-sm font-medium transition-colors text-left ${
                            breadcrumb.length === 0
                                ? 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 border border-purple-300'
                                : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-200'
                        }`}
                    >
                        Colombia
                    </button>
                </div>

                {/* Departamentos */}
                <div>
                    <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">Departamento</label>
                    <select
                        value={breadcrumb.length > 0 ? breadcrumb[0].id : ''}
                        onChange={(e) => {
                            if (!e.target.value) {
                                handleBreadcrumbClick(-1);
                            } else {
                                const dept = departmentOptions.find((d) => d.id.toString() === e.target.value);
                                if (dept) {
                                    setBreadcrumb([{ id: dept.id, name: dept.name, type: 'state' }]);
                                    fetchGeozones('city', dept.id);
                                }
                            }
                        }}
                        className="w-full px-3 py-2 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600"
                    >
                        <option value="">Todos</option>
                        {departmentOptions.map((dept) => (
                            <option key={dept.id} value={dept.id}>
                                {dept.name}
                            </option>
                        ))}
                    </select>
                </div>

                {/* Ciudades */}
                <div>
                    <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">Ciudad</label>
                    <select
                        value={breadcrumb.length > 1 ? breadcrumb[1].id : ''}
                        onChange={(e) => {
                            if (!e.target.value) {
                                const newBreadcrumb = breadcrumb.slice(0, 1);
                                setBreadcrumb(newBreadcrumb);
                                if (newBreadcrumb.length > 0) {
                                    fetchGeozones('city', newBreadcrumb[0].id);
                                }
                            } else {
                                const city = cityOptions.find((c) => c.id.toString() === e.target.value);
                                if (city && breadcrumb.length > 0) {
                                    const newBreadcrumb = [breadcrumb[0], { id: city.id, name: city.name, type: 'city' }];
                                    setBreadcrumb(newBreadcrumb);
                                    fetchGeozones('admin_district', city.id);
                                }
                            }
                        }}
                        disabled={breadcrumb.length === 0}
                        className="w-full px-3 py-2 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 disabled:opacity-50"
                    >
                        <option value="">Todas</option>
                        {cityOptions.map((city) => (
                            <option key={city.id} value={city.id}>
                                {city.name}
                            </option>
                        ))}
                    </select>
                </div>
            </div>

            {/* Breadcrumb */}
            <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300 flex-wrap">
                <button
                    onClick={() => handleBreadcrumbClick(-1)}
                    className="hover:text-purple-600 transition-colors"
                >
                    Colombia
                </button>
                {breadcrumb.map((item, idx) => (
                    <div key={item.id} className="flex items-center gap-2">
                        <span className="text-gray-400">/</span>
                        <button
                            onClick={() => handleBreadcrumbClick(idx)}
                            className="hover:text-purple-600 transition-colors"
                        >
                            {item.name}
                        </button>
                    </div>
                ))}
            </div>

            {/* Mapa */}
            <div
                style={{ height, isolation: 'isolate', position: 'relative', zIndex: 0 }}
                className="rounded-xl overflow-hidden border border-gray-200 dark:border-gray-700 shadow-lg"
            >
                {loading && (
                    <div className="absolute inset-0 z-50 flex items-center justify-center bg-white/80 dark:bg-gray-800/80">
                        <Spinner size="lg" color="primary" text="Cargando mapa..." />
                    </div>
                )}
                <MapContainer
                    center={[4.5709, -74.2973]}
                    zoom={6}
                    minZoom={4}
                    maxZoom={15}
                    style={{ height: '100%', width: '100%' }}
                    scrollWheelZoom
                >
                    <TileLayer
                        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        maxZoom={19}
                    />
                    {features.map((f) => {
                        const color = getDensityColor(f.properties.name, metricType);
                        const count = ordersMap.get(f.properties.name.toUpperCase().trim()) || 0;
                        const percentage = totalOrders > 0 ? ((count / totalOrders) * 100).toFixed(1) : '0';

                        return (
                            <GeoJSONLayer
                                key={`${f.properties.id}`}
                                data={f.geometry as any}
                                style={() => ({
                                    color,
                                    weight: 2,
                                    fillColor: color,
                                    fillOpacity: 0.5,
                                    opacity: 0.9,
                                })}
                                onEachFeature={(_, layer) => {
                                    layer.bindTooltip(
                                        `<div style="font-family:system-ui;font-size:12px"><b>${f.properties.name}</b><br/>Órdenes: ${count.toLocaleString()}<br/>Porcentaje: ${percentage}%<br/><span style="color:${color};font-size:11px;">⚪ Clic para explorar</span></div>`,
                                        { sticky: true }
                                    );
                                    if (DRILL_CONFIG[drillLevel].nextLevel) {
                                        layer.on('click', () => handlePolygonClick(f));
                                        if ((layer as any)._path) {
                                            (layer as any)._path.style.cursor = 'pointer';
                                        }
                                    }
                                }}
                            />
                        );
                    })}
                    <FitBounds features={features} fitKey={`${drillLevel}|${breadcrumb.length}`} />
                    <ZoomReporter onZoomChange={setZoom} />
                </MapContainer>
            </div>

            {/* Leyenda de colores */}
            <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow border border-gray-200 dark:border-gray-700">
                <p className="text-sm font-semibold mb-3 text-gray-700 dark:text-gray-200">Escala de Densidad</p>
                <div className="space-y-2 text-sm">
                    <div className="flex items-center gap-2">
                        <div className="w-4 h-4 rounded" style={{ backgroundColor: '#16a34a' }}></div>
                        <span className="text-gray-600 dark:text-gray-300">
                            {metricType === 'percentage' ? '&gt; 66% del total' : 'Alta densidad'}
                        </span>
                    </div>
                    <div className="flex items-center gap-2">
                        <div className="w-4 h-4 rounded" style={{ backgroundColor: '#ca8a04' }}></div>
                        <span className="text-gray-600 dark:text-gray-300">
                            {metricType === 'percentage' ? '33-66% del total' : 'Densidad media'}
                        </span>
                    </div>
                    <div className="flex items-center gap-2">
                        <div className="w-4 h-4 rounded" style={{ backgroundColor: '#dc2626' }}></div>
                        <span className="text-gray-600 dark:text-gray-300">
                            {metricType === 'percentage' ? '&lt; 33% del total' : 'Baja densidad'}
                        </span>
                    </div>
                    <div className="flex items-center gap-2">
                        <div className="w-4 h-4 rounded" style={{ backgroundColor: '#d1d5db' }}></div>
                        <span className="text-gray-600 dark:text-gray-300">Sin Órdenes</span>
                    </div>
                </div>
            </div>
        </div>
    );
}
