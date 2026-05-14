'use client';

import { useEffect, useState } from 'react';
import { MapContainer, TileLayer, GeoJSON as GeoJSONLayer, useMap } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { getOrderZoneAction } from '../../infra/actions';
import type { Geozone } from '../../domain/types';

interface Props {
    businessId: number;
    orderId: string;
    height?: string;
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

function FitBoundsToGeometry({ geometry }: { geometry: any }) {
    const map = useMap();
    useEffect(() => {
        if (!geometry) return;
        try {
            const layer = L.geoJSON(geometry);
            const bounds = layer.getBounds();
            if (bounds.isValid()) map.fitBounds(bounds, { padding: [10, 10], maxZoom: 14 });
        } catch { }
    }, [geometry, map]);
    return null;
}

export function GeozoneMiniMap({ businessId, orderId, height = '220px' }: Props) {
    const [zone, setZone] = useState<Geozone | null>(null);
    const [level, setLevel] = useState<string>('');
    const [loading, setLoading] = useState(true);
    const [empty, setEmpty] = useState(false);

    useEffect(() => {
        let cancelled = false;
        if (!businessId || !orderId) return;
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
    }, [businessId, orderId]);

    if (loading) {
        return (
            <div className="w-full rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 flex items-center justify-center text-xs text-gray-500" style={{ height }}>
                Cargando geozona...
            </div>
        );
    }
    if (empty || !zone || !zone.geometry) {
        return (
            <div className="w-full rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 flex items-center justify-center text-xs text-gray-500" style={{ height }}>
                Sin geozona disponible para esta direccion
            </div>
        );
    }
    const color = typeColor[level] || '#6366f1';
    return (
        <div className="rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700" style={{ isolation: 'isolate' }}>
            <div className="px-3 py-2 bg-gray-50 dark:bg-gray-800 text-xs flex items-center justify-between">
                <span className="font-semibold text-gray-700 dark:text-gray-200">
                    {typeLabel[level] || level}: {zone.name}
                </span>
                <span className="px-2 py-0.5 rounded-full text-[10px] font-bold text-white" style={{ backgroundColor: color }}>
                    {typeLabel[level] || level}
                </span>
            </div>
            <div style={{ height }}>
                <MapContainer center={[4.6, -74.08]} zoom={5} style={{ height: '100%', width: '100%' }} scrollWheelZoom={false} dragging={false} zoomControl={false} doubleClickZoom={false} attributionControl={false}>
                    <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                    <GeoJSONLayer
                        key={`${zone.id}-${level}`}
                        data={zone.geometry as any}
                        style={() => ({ color, weight: 2, fillColor: color, fillOpacity: 0.35 })}
                    />
                    <FitBoundsToGeometry geometry={zone.geometry} />
                </MapContainer>
            </div>
        </div>
    );
}
