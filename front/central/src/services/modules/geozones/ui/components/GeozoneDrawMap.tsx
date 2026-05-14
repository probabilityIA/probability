'use client';

import { useEffect, useState } from 'react';
import { MapContainer, TileLayer, Polygon, CircleMarker, useMapEvents, GeoJSON as GeoJSONLayer } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import { Geozone, GeoJSONPolygon } from '../../domain/types';

interface DrawHandlerProps {
    onClick: (lat: number, lng: number) => void;
}
function DrawHandler({ onClick }: DrawHandlerProps) {
    useMapEvents({ click: (e) => onClick(e.latlng.lat, e.latlng.lng) });
    return null;
}

interface GeozoneDrawMapProps {
    points: Array<[number, number]>;
    onChange: (points: Array<[number, number]>) => void;
    contextLayers?: Geozone[];
    height?: string;
    initialCenter?: [number, number];
    initialZoom?: number;
}

export default function GeozoneDrawMap({
    points,
    onChange,
    contextLayers = [],
    height = '420px',
    initialCenter = [4.7110, -74.0721],
    initialZoom = 11,
}: GeozoneDrawMapProps) {
    const [hoverIdx, setHoverIdx] = useState<number | null>(null);

    useEffect(() => {
        const onKey = (e: KeyboardEvent) => {
            if (e.key === 'Backspace' && points.length > 0) {
                e.preventDefault();
                onChange(points.slice(0, -1));
            }
            if (e.key === 'Escape' && points.length > 0) {
                onChange([]);
            }
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, [points, onChange]);

    const polygonPositions: Array<[number, number]> = points.length >= 3 ? points : [];

    return (
        <div className="relative">
            <div style={{ height }} className="rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700">
                <MapContainer center={initialCenter} zoom={initialZoom} style={{ height: '100%', width: '100%' }} scrollWheelZoom>
                    <TileLayer
                        attribution='&copy; OpenStreetMap'
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    />
                    {contextLayers.map((g) => g.geometry && (
                        <GeoJSONLayer
                            key={`ctx-${g.id}`}
                            data={g.geometry as any}
                            style={() => ({ color: '#94a3b8', weight: 1, fillOpacity: 0.05, dashArray: '4 4' })}
                        />
                    ))}
                    <DrawHandler onClick={(lat, lng) => onChange([...points, [lat, lng]])} />
                    {polygonPositions.length > 0 && (
                        <Polygon
                            positions={polygonPositions}
                            pathOptions={{ color: '#ec4899', weight: 2, fillColor: '#ec4899', fillOpacity: 0.25 }}
                        />
                    )}
                    {points.length > 0 && points.length < 3 && (
                        <Polygon
                            positions={points}
                            pathOptions={{ color: '#ec4899', weight: 2, dashArray: '4 4', fill: false }}
                        />
                    )}
                    {points.map((p, i) => (
                        <CircleMarker
                            key={i}
                            center={p}
                            radius={hoverIdx === i ? 8 : 6}
                            pathOptions={{
                                color: i === 0 ? '#10b981' : '#ec4899',
                                fillColor: i === 0 ? '#10b981' : '#ec4899',
                                fillOpacity: 1,
                                weight: 2,
                            }}
                            eventHandlers={{
                                mouseover: () => setHoverIdx(i),
                                mouseout: () => setHoverIdx(null),
                                click: (e) => {
                                    (e as any).originalEvent?.stopPropagation?.();
                                    onChange(points.filter((_, idx) => idx !== i));
                                },
                            }}
                        />
                    ))}
                </MapContainer>
            </div>

            <div className="absolute top-3 right-3 bg-white/95 dark:bg-gray-800/95 backdrop-blur rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 p-3 text-xs z-[1000] max-w-xs">
                <div className="font-semibold text-gray-900 dark:text-white mb-1.5">Como dibujar</div>
                <ul className="space-y-1 text-gray-600 dark:text-gray-300">
                    <li>&bull; Click en el mapa para agregar vertices</li>
                    <li>&bull; Click sobre un vertice para borrarlo</li>
                    <li>&bull; <kbd className="px-1 py-0.5 bg-gray-100 dark:bg-gray-700 rounded">Backspace</kbd> deshacer</li>
                    <li>&bull; <kbd className="px-1 py-0.5 bg-gray-100 dark:bg-gray-700 rounded">Esc</kbd> limpiar</li>
                </ul>
                <div className="mt-2 pt-2 border-t border-gray-200 dark:border-gray-700">
                    <span className="text-gray-500 dark:text-gray-400">Vertices: </span>
                    <span className={`font-bold ${points.length >= 3 ? 'text-green-600' : 'text-orange-500'}`}>{points.length}</span>
                    <span className="text-gray-500 dark:text-gray-400"> (min 3)</span>
                </div>
            </div>
        </div>
    );
}

export function pointsToPolygon(points: Array<[number, number]>): GeoJSONPolygon | null {
    if (points.length < 3) return null;
    const ring = points.map(([lat, lng]) => [lng, lat]);
    const first = ring[0];
    const last = ring[ring.length - 1];
    if (first[0] !== last[0] || first[1] !== last[1]) ring.push([first[0], first[1]]);
    return { type: 'Polygon', coordinates: [ring] };
}
