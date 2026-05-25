'use client';

import { useEffect, useState } from 'react';
import { MapContainer, TileLayer, GeoJSON as GeoJSONLayer, Marker, useMap } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { getOrderZoneAction } from '../../infra/actions';
import type { Geozone } from '../../domain/types';

interface OriginInfo {
    address?: string;
    lat?: number | null;
    lng?: number | null;
}

interface DestinationInfo {
    address?: string;
}

interface Props {
    businessId: number;
    orderId?: string;
    geozone?: Geozone | null;
    lat?: number | null;
    lng?: number | null;
    height?: string;
    showHeader?: boolean;
    origin?: OriginInfo | null;
    destination?: DestinationInfo | null;
    carrierRate?: number | null;
    carrierName?: string | null;
    carrierEstimated?: boolean;
    viewMode?: 'origin-destination' | 'destination-only';
}

function rateColor(rate: number): string {
    if (rate >= 0.9) return '#16a34a';
    if (rate >= 0.8) return '#65a30d';
    if (rate >= 0.7) return '#ca8a04';
    if (rate >= 0.6) return '#ea580c';
    return '#dc2626';
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

const originIcon = L.divIcon({
    className: 'probability-origin-marker',
    html: `<div style="
        width: 24px;
        height: 24px;
        background: linear-gradient(135deg, #059669, #047857);
        border: 3px solid #ffffff;
        border-radius: 50%;
        box-shadow: 0 2px 8px rgba(0,0,0,0.35);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 12px;
        color: white;
        font-weight: bold;
    ">O</div>`,
    iconSize: [24, 24],
    iconAnchor: [12, 12],
});

function FitBoundsToGeometry({ geometry, lat, lng, originLat, originLng, viewMode = 'origin-destination' }: { geometry: any; lat?: number | null; lng?: number | null; originLat?: number | null; originLng?: number | null; viewMode?: 'origin-destination' | 'destination-only' }) {
    const map = useMap();
    useEffect(() => {
        if (!geometry) return;
        try {
            const layer = L.geoJSON(geometry);
            const bounds = layer.getBounds();
            if (lat != null && lng != null) bounds.extend([lat, lng]);
            if (viewMode === 'origin-destination' && originLat != null && originLng != null) bounds.extend([originLat, originLng]);
            if (bounds.isValid()) map.fitBounds(bounds, { padding: [16, 16], maxZoom: 14 });
        } catch { }
    }, [geometry, lat, lng, originLat, originLng, viewMode, map]);
    return null;
}

export function GeozoneMiniMap({ businessId, orderId, geozone: geozoneProp, lat, lng, height = '220px', showHeader = true, origin, destination, carrierRate, carrierName, carrierEstimated, viewMode = 'origin-destination' }: Props) {
    const hasCarrierRate = carrierRate != null && Number.isFinite(carrierRate);
    const carrierPct = hasCarrierRate ? Math.round((carrierRate as number) * 100) : null;
    const carrierColor = hasCarrierRate ? rateColor(carrierRate as number) : null;
    const originLat = origin?.lat;
    const originLng = origin?.lng;
    const hasOriginPoint = originLat != null && originLng != null && Number.isFinite(originLat) && Number.isFinite(originLng);
    const destinationBanner = destination?.address ? (
        <div className="px-3 py-2 bg-violet-50 border-b border-violet-200 text-xs flex items-center gap-2">
            <span className="inline-flex items-center justify-center w-5 h-5 rounded-full bg-violet-600 text-white font-bold text-[10px] shrink-0">D</span>
            <span className="font-semibold text-violet-900 shrink-0">Destino:</span>
            <span className="text-violet-900 truncate">{destination.address}</span>
        </div>
    ) : null;

    const originBanner = origin?.address && viewMode === 'origin-destination' ? (
        <div className="px-3 py-1.5 bg-emerald-50/60 border-b border-emerald-100 text-[11px] flex items-center gap-2">
            <span className="inline-flex items-center justify-center w-4 h-4 rounded-full bg-emerald-600 text-white font-bold text-[9px] shrink-0">O</span>
            <span className="text-emerald-700 shrink-0">Origen:</span>
            <span className="text-emerald-800 truncate">{origin.address}</span>
        </div>
    ) : null;

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
                <div className="rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700" style={{ isolation: 'isolate' }}>
                    {destinationBanner}
                    {originBanner}
                    <div style={{ height }}>
                        <MapContainer center={[lat as number, lng as number]} zoom={14} style={{ height: '100%', width: '100%' }} scrollWheelZoom={false} dragging={false} zoomControl={false} doubleClickZoom={false} attributionControl={false}>
                            <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                            <Marker position={[lat as number, lng as number]} icon={deliveryIcon} />
                            {viewMode === 'origin-destination' && hasOriginPoint && <Marker position={[originLat as number, originLng as number]} icon={originIcon} />}
                        </MapContainer>
                    </div>
                </div>
            );
        }
        return (
            <div className="w-full rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 flex items-center justify-center text-xs text-gray-500" style={{ height }}>
                Sin geozona disponible para esta direccion
            </div>
        );
    }
    const defaultColor = typeColor[level] || '#6366f1';
    const polygonColor = carrierColor ?? defaultColor;
    const polygonOpacity = carrierColor ? (carrierEstimated ? 0.35 : 0.55) : 0.35;
    return (
        <div className="rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700" style={{ isolation: 'isolate' }}>
            {showHeader && (
                <div className="px-3 py-2 bg-gray-50 dark:bg-gray-800 text-xs flex items-center justify-between">
                    <span className="font-semibold text-gray-700 dark:text-gray-200">
                        {typeLabel[level] || level}: {zone.name}
                    </span>
                    <span className="px-2 py-0.5 rounded-full text-[10px] font-bold text-white" style={{ backgroundColor: defaultColor }}>
                        {typeLabel[level] || level}
                    </span>
                </div>
            )}
            {destinationBanner}
            {originBanner}
            <div className="relative" style={{ height }}>
                <MapContainer center={[4.6, -74.08]} zoom={5} style={{ height: '100%', width: '100%' }} scrollWheelZoom={false} dragging={false} zoomControl={false} doubleClickZoom={false} attributionControl={false}>
                    <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                    <GeoJSONLayer
                        key={`${zone.id}-${level}-${carrierPct ?? 'none'}-${carrierEstimated ? 'est' : 'real'}`}
                        data={zone.geometry as any}
                        style={() => ({ color: polygonColor, weight: 2, fillColor: polygonColor, fillOpacity: polygonOpacity })}
                    />
                    {hasPoint && (
                        <Marker position={[lat as number, lng as number]} icon={deliveryIcon} />
                    )}
                    {viewMode === 'origin-destination' && hasOriginPoint && (
                        <Marker position={[originLat as number, originLng as number]} icon={originIcon} />
                    )}
                    <FitBoundsToGeometry geometry={zone.geometry} lat={lat} lng={lng} originLat={originLat} originLng={originLng} viewMode={viewMode} />
                </MapContainer>
                {hasCarrierRate && carrierColor && (
                    <div
                        className="absolute top-3 right-3 z-[400] flex items-center gap-2 rounded-full px-3 py-1.5 shadow-lg pointer-events-none"
                        style={{ backgroundColor: carrierColor, color: '#fff' }}
                    >
                        <span className="text-xs font-semibold tracking-wide opacity-90">
                            {carrierName || 'Efectividad'}{carrierEstimated ? ' (est.)' : ''}
                        </span>
                        <span className="text-lg font-extrabold tabular-nums">{carrierPct}%</span>
                    </div>
                )}
            </div>
        </div>
    );
}
