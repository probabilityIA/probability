'use client';

import { useEffect, useMemo } from 'react';
import { MapContainer, TileLayer, Marker, Polyline, Tooltip, useMap } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { RouteDetail, RouteStopInfo } from '../../domain/types';

interface RouteMapProps {
    route: RouteDetail;
    height?: string;
    selectedStopId?: number | null;
    onSelectStop?: (stopId: number) => void;
}

const STOP_COLORS: Record<string, string> = {
    pending: '#6b7280',
    arrived: '#ca8a04',
    delivered: '#16a34a',
    failed: '#dc2626',
    skipped: '#9ca3af',
};

function stopIcon(sequence: number, status: string, selected: boolean): L.DivIcon {
    const color = STOP_COLORS[status] || STOP_COLORS.pending;
    const size = selected ? 34 : 28;
    return L.divIcon({
        className: 'probability-route-stop',
        html: `<div style="width:${size}px;height:${size}px;background:${color};border:3px solid #fff;border-radius:50%;box-shadow:0 2px 8px rgba(0,0,0,.35);display:flex;align-items:center;justify-content:center;color:#fff;font-weight:700;font-size:13px;">${sequence}</div>`,
        iconSize: [size, size],
        iconAnchor: [size / 2, size / 2],
    });
}

const originIcon = L.divIcon({
    className: 'probability-route-origin',
    html: `<div style="width:32px;height:32px;background:linear-gradient(135deg,#7c3aed,#4f46e5);border:3px solid #fff;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,.35);display:flex;align-items:center;justify-content:center;color:#fff;font-weight:700;font-size:13px;">B</div>`,
    iconSize: [32, 32],
    iconAnchor: [16, 16],
});

function FitBounds({ points }: { points: [number, number][] }) {
    const map = useMap();
    useEffect(() => {
        if (points.length === 0) return;
        if (points.length === 1) {
            map.setView(points[0], 14);
            return;
        }
        map.fitBounds(L.latLngBounds(points), { padding: [40, 40] });
    }, [map, points]);
    return null;
}

function FocusStop({ stop }: { stop: RouteStopInfo | null }) {
    const map = useMap();
    useEffect(() => {
        if (stop?.lat == null || stop?.lng == null) return;
        map.panTo([stop.lat, stop.lng]);
    }, [map, stop]);
    return null;
}

export default function RouteMap({ route, height = '420px', selectedStopId, onSelectStop }: RouteMapProps) {
    const stops = useMemo(
        () => [...(route.stops || [])].sort((a, b) => a.sequence - b.sequence),
        [route.stops]
    );
    const located = useMemo(() => stops.filter((s) => s.lat != null && s.lng != null), [stops]);
    const origin: [number, number] | null =
        route.origin_lat != null && route.origin_lng != null ? [route.origin_lat, route.origin_lng] : null;

    const path = useMemo(() => {
        const pts: [number, number][] = located.map((s) => [s.lat as number, s.lng as number]);
        return origin ? [origin, ...pts] : pts;
    }, [located, origin?.[0], origin?.[1]]);

    const selectedStop = stops.find((s) => s.id === selectedStopId) || null;
    const missing = stops.length - located.length;

    if (path.length === 0) {
        return (
            <div className="flex items-center justify-center rounded-lg border border-dashed border-gray-300 dark:border-gray-600 text-sm text-gray-500 dark:text-gray-400" style={{ height }}>
                {'Ninguna parada tiene ubicación para mostrar en el mapa'}
            </div>
        );
    }

    return (
        <div className="relative">
            <div className="rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700" style={{ height }}>
                <MapContainer center={path[0]} zoom={12} style={{ height: '100%', width: '100%' }} scrollWheelZoom>
                    <TileLayer
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        attribution="&copy; OpenStreetMap"
                    />
                    <FitBounds points={path} />
                    <FocusStop stop={selectedStop} />
                    {path.length > 1 && (
                        <Polyline positions={path} pathOptions={{ color: '#6d28d9', weight: 4, opacity: 0.75, dashArray: '8 8' }} />
                    )}
                    {origin && (
                        <Marker position={origin} icon={originIcon}>
                            <Tooltip direction="top" offset={[0, -14]}>
                                {`Origen: ${route.origin_address || 'Bodega'}`}
                            </Tooltip>
                        </Marker>
                    )}
                    {located.map((s) => (
                        <Marker
                            key={s.id}
                            position={[s.lat as number, s.lng as number]}
                            icon={stopIcon(s.sequence, s.status, s.id === selectedStopId)}
                            eventHandlers={{ click: () => onSelectStop?.(s.id) }}
                        >
                            <Tooltip direction="top" offset={[0, -14]}>
                                {`${s.sequence}. ${s.customer_name} - ${s.address}`}
                            </Tooltip>
                        </Marker>
                    ))}
                </MapContainer>
            </div>
            <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
                {origin && <span>{'B = bodega de origen'}</span>}
                {route.total_distance_km != null && (
                    <span>{`Distancia: ${route.total_distance_km.toFixed(1)} km`}</span>
                )}
                {route.total_duration_min != null && <span>{`Duración: ${route.total_duration_min} min`}</span>}
                {!origin && <span className="text-amber-600">{'Sin bodega de origen con ubicación'}</span>}
                {missing > 0 && (
                    <span className="text-amber-600">{`${missing} parada(s) sin ubicación no aparecen en el mapa`}</span>
                )}
                <span>{'La línea muestra el orden de visita, no el trazado por las calles'}</span>
            </div>
        </div>
    );
}
