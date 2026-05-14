'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { MapContainer, TileLayer, GeoJSON as GeoJSONLayer, useMap, useMapEvents } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { Geozone, GeozoneType } from '../../domain/types';

const TYPE_COLORS: Record<GeozoneType, string> = {
    country: '#0ea5e9',
    state: '#8b5cf6',
    city: '#10b981',
    admin_district: '#6366f1',
    locality: '#f59e0b',
    neighborhood: '#ef4444',
    barrio: '#dc2626',
    custom: '#ec4899',
};

interface FitBoundsProps {
    geozones: Geozone[];
    fitKey: string;
}

function FocusOnSelected({ geozones, selectedId }: { geozones: Geozone[]; selectedId: number | null | undefined }) {
    const map = useMap();
    const lastFocused = useRef<number | null>(null);
    useEffect(() => {
        if (!selectedId || lastFocused.current === selectedId) return;
        const target = geozones.find((g) => g.id === selectedId);
        if (!target?.geometry) return;
        try {
            const layer = L.geoJSON(target.geometry as any);
            const bounds = layer.getBounds();
            if (bounds.isValid()) {
                map.fitBounds(bounds, { padding: [40, 40], maxZoom: 14 });
                lastFocused.current = selectedId;
            }
        } catch {}
    }, [selectedId, geozones, map]);
    return null;
}

function FitBounds({ geozones, fitKey }: FitBoundsProps) {
    const map = useMap();
    const lastSig = useRef<string>('');
    useEffect(() => {
        const ids = geozones.map((g) => g.id).sort((a, b) => a - b).join(',');
        const sig = `${fitKey}|${ids}`;
        if (lastSig.current === sig) return;
        const layers: L.Layer[] = [];
        geozones.forEach((g) => {
            if (g.geometry) {
                try { layers.push(L.geoJSON(g.geometry as any)); } catch {}
            }
        });
        if (layers.length === 0) return;
        const group = L.featureGroup(layers);
        const bounds = group.getBounds();
        if (bounds.isValid()) {
            map.fitBounds(bounds, { padding: [20, 20], maxZoom: 13 });
            lastSig.current = sig;
        }
    }, [geozones, map, fitKey]);
    return null;
}

function ZoomReporter({ onZoomChange }: { onZoomChange: (z: number) => void }) {
    const map = useMapEvents({
        zoomend: () => onZoomChange(map.getZoom()),
    });
    useEffect(() => { onZoomChange(map.getZoom()); }, [map, onZoomChange]);
    return null;
}

interface GeozoneMapProps {
    geozones: Geozone[];
    selectedId?: number | null;
    onSelect?: (g: Geozone) => void;
    onZoomChange?: (zoom: number) => void;
    height?: string;
    fitKey?: string;
}

export default function GeozoneMap({ geozones, selectedId, onSelect, onZoomChange, height = '600px', fitKey = 'default' }: GeozoneMapProps) {
    const items = useMemo(() => geozones.filter((g) => !!g.geometry), [geozones]);
    const layersRef = useRef<Map<number, L.Layer>>(new Map());

    return (
        <div style={{ height, isolation: 'isolate', position: 'relative', zIndex: 0 }} className="rounded-xl overflow-hidden border border-gray-200 dark:border-gray-700 shadow-lg">
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
                {items.map((g) => {
                    const color = TYPE_COLORS[g.type] || '#6b7280';
                    const isSelected = selectedId === g.id;
                    return (
                        <GeoJSONLayer
                            key={`${g.id}-${isSelected ? 'sel' : 'idle'}`}
                            data={g.geometry as any}
                            style={() => ({
                                color,
                                weight: isSelected ? 3 : 2,
                                fillColor: color,
                                fillOpacity: isSelected ? 0.5 : 0.3,
                                opacity: 0.9,
                            })}
                            onEachFeature={(_, layer) => {
                                layersRef.current.set(g.id, layer);
                                layer.bindTooltip(
                                    `<div style="font-family:system-ui;font-size:12px"><b>${g.name}</b><br/><span style="color:${color}">${g.type}</span>${g.code ? ` &middot; ${g.code}` : ''}</div>`,
                                    { sticky: true }
                                );
                                if (onSelect) layer.on('click', () => onSelect(g));
                            }}
                        />
                    );
                })}
                <FitBounds geozones={items} fitKey={fitKey} />
                <FocusOnSelected geozones={items} selectedId={selectedId} />
                {onZoomChange && <ZoomReporter onZoomChange={onZoomChange} />}
            </MapContainer>
        </div>
    );
}
