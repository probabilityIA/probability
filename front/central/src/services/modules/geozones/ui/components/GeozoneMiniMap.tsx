'use client';

import { useEffect, useState } from 'react';
import { MapContainer, TileLayer, GeoJSON as GeoJSONLayer, Marker, useMap } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { getOrderZoneAction } from '../../infra/actions';
import type { Geozone } from '../../domain/types';

interface Props {
    businessId: number;
    orderId?: string;
    geozone?: Geozone | null;
    lat?: number | null;
    lng?: number | null;
    height?: string;
    showHeader?: boolean;
}

const typeLabel: Record<string, string> = {
    barrio: 'Barrio',
    neighborhood: 'UPZ',
    admin_district: 'Localidad',
    locality: 'Corregimiento',
    city: 'Municipio',
    state: 'Departamento',
    country: 'Pais',
};

const typeColor: Record<string, string> = {
    barrio: '#dc2626',
    neighborhood: '#ef4444',
    admin_district: '#6366f1',
    locality: '#f59e0b',
    city: '#10b981',
    state: '#8b5cf6',
    country: '#0ea5e9',
};

const deliveryIcon = L.divIcon({
    className: 'probability-delivery-marker',
    html: `<div style="
        width: 28px;
        height: 28px;
        background: linear-gradient(135deg, #7c3aed, #4f46e5);
        border: 3px solid #ffffff;
        border-radius: 50%;
        box-shadow: 0 2px 8px rgba(0,0,0,0.35);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        color: white;
        font-weight: bold;
    ">P</div>`,
    iconSize: [28, 28],
    iconAnchor: [14, 14],
});

function FitBoundsToGeometry({ geometry, lat, lng }: { geometry: any; lat?: number | null; lng?: number | null }) {
    const map = useMap();
    useEffect(() => {
        if (!geometry) return;
        try {
            const layer = L.geoJSON(geometry);
            const bounds = layer.getBounds();
            if (lat != null && lng != null) {
                bounds.extend([lat, lng]);
            }
            if (bounds.isValid()) map.fitBounds(bounds, { padding: [12, 12], maxZoom: 14 });
        } catch { }
    }, [geometry, lat, lng, map]);
    return null;
}

export function GeozoneMiniMap({ businessId, orderId, geozone: geozoneProp, lat, lng, height = '220px', showHeader = true }: Props) {
    const [zone, setZone] = useState<Geozone | null>(geozoneProp ?? null);
    const [level, setLevel] = useState<string>(geozoneProp?.type || '');
    const [loading, setLoading] = useState(!geozoneProp && !!orderId);
    const [empty, setEmpty] = useState(false);

    useEffect(() => {
        if (geozoneProp) {
            setZone(geozoneProp);
            setLevel(geozoneProp.type || '');
            setLoading(false);
            return;
        }
        if (!businessId || !orderId) return;
        let cancelled = false;
        setLoading(true);
        setEmpty(false);
        (async () => {
            try {
                const g = await getOrderZoneAction(orderId, businessId);
                if (cancelled) return;
                if (!g) {
                    setEmpty(true);
                    setLoading(false);
                    return;
                }
                setZone(g);
                setLevel(g.type || '');
            } catch {
                if (!cancelled) setEmpty(true);
            } finally {
                if (!cancelled) setLoading(false);
            }
        })();
        return () => { cancelled = true; };
    }, [businessId, orderId, geozoneProp]);

    const hasPoint = lat != null && lng != null && Number.isFinite(lat) && Number.isFinite(lng);

    if (loading) {
        return (
            <div className="w-full rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 flex items-center justify-center text-xs text-gray-500" style={{ height }}>
                Cargando geozona...
            </div>
        );
    }
    if (empty || !zone || !zone.geometry) {
        if (hasPoint) {
            return (
                <div className="rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700" style={{ height, isolation: 'isolate' }}>
                    <MapContainer center={[lat as number, lng as number]} zoom={14} style={{ height: '100%', width: '100%' }} scrollWheelZoom={false} dragging={false} zoomControl={false} doubleClickZoom={false} attributionControl={false}>
                        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                        <Marker position={[lat as number, lng as number]} icon={deliveryIcon} />
                    </MapContainer>
                </div>
            );
        }
        return (
            <div className="w-full rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 flex items-center justify-center text-xs text-gray-500" style={{ height }}>
                Sin geozona disponible para esta direccion
            </div>
        );
    }
    const color = typeColor[level] || '#6366f1';
    return (
        <div className="rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700" style={{ isolation: 'isolate' }}>
            {showHeader && (
                <div className="px-3 py-2 bg-gray-50 dark:bg-gray-800 text-xs flex items-center justify-between">
                    <span className="font-semibold text-gray-700 dark:text-gray-200">
                        {typeLabel[level] || level}: {zone.name}
                    </span>
                    <span className="px-2 py-0.5 rounded-full text-[10px] font-bold text-white" style={{ backgroundColor: color }}>
                        {typeLabel[level] || level}
                    </span>
                </div>
            )}
            <div style={{ height }}>
                <MapContainer center={[4.6, -74.08]} zoom={5} style={{ height: '100%', width: '100%' }} scrollWheelZoom={false} dragging={false} zoomControl={false} doubleClickZoom={false} attributionControl={false}>
                    <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                    <GeoJSONLayer
                        key={`${zone.id}-${level}`}
                        data={zone.geometry as any}
                        style={() => ({ color, weight: 2, fillColor: color, fillOpacity: 0.35 })}
                    />
                    {hasPoint && (
                        <Marker position={[lat as number, lng as number]} icon={deliveryIcon} />
                    )}
                    <FitBoundsToGeometry geometry={zone.geometry} lat={lat} lng={lng} />
                </MapContainer>
            </div>
        </div>
    );
}
