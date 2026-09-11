import React, { useCallback, useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';

// Fix for default marker icon in React Leaflet
import iconImg from 'leaflet/dist/images/marker-icon.png';
import iconShadowImg from 'leaflet/dist/images/marker-shadow.png';

const DefaultIcon = L.icon({
    iconUrl: typeof iconImg === 'string' ? iconImg : iconImg.src,
    shadowUrl: typeof iconShadowImg === 'string' ? iconShadowImg : iconShadowImg.src,
    iconSize: [25, 41],
    iconAnchor: [12, 41]
});

L.Marker.prototype.options.icon = DefaultIcon;

interface MapComponentProps {
    address: string;
    city: string;
    height?: string;
    latitude?: number | null;
    longitude?: number | null;
    expandable?: boolean;
}

const RecenterAutomatically = ({ lat, lng }: { lat: number; lng: number }) => {
    const map = useMap();
    useEffect(() => {
        map.setView([lat, lng]);
    }, [lat, lng, map]);
    return null;
};

const MapComponent: React.FC<MapComponentProps> = ({ address, city, height = '400px', latitude, longitude, expandable = false }) => {
    const [position, setPosition] = useState<[number, number] | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [expanded, setExpanded] = useState(false);
    const [mounted, setMounted] = useState(false);

    useEffect(() => setMounted(true), []);

    const closeExpanded = useCallback(() => setExpanded(false), []);

    useEffect(() => {
        if (!expanded) return;
        const onKey = (e: KeyboardEvent) => {
            if (e.key === 'Escape') {
                e.stopPropagation();
                closeExpanded();
            }
        };
        document.addEventListener('keydown', onKey, true);
        const previous = document.body.style.overflow;
        document.body.style.overflow = 'hidden';
        return () => {
            document.removeEventListener('keydown', onKey, true);
            document.body.style.overflow = previous;
        };
    }, [expanded, closeExpanded]);

    useEffect(() => {
        // If direct coordinates are provided, use them without geocoding
        if (latitude != null && longitude != null) {
            setPosition([latitude, longitude]);
            setLoading(false);
            setError(null);
            return;
        }

        const geocodeAddress = async () => {
            if (!address || !city) return;

            setLoading(true);
            setError(null);

            try {
                // Llamamos a nuestro propio backend como proxy para evitar restricciones CORS/User-Agent
                // Usamos /api/v1/geocode para que funcione a través del proxy Nginx en producción
                const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:3050/api/v1';
                const url = `${apiBase}/geocode?address=${encodeURIComponent(address)}&city=${encodeURIComponent(city)}`;

                const response = await fetch(url);
                if (!response.ok) throw new Error('geocode request failed');

                const data: { lat: number; lon: number; found: boolean; fallback: boolean } = await response.json();

                if (data.found) {
                    setPosition([data.lat, data.lon]);
                    if (data.fallback) {
                        setError('Dirección exacta no encontrada, mostrando ubicación de la ciudad.');
                    }
                } else {
                    setError('No se pudo localizar la dirección.');
                }
            } catch (err) {
                console.error('Geocoding error:', err);
                setError('Error al cargar el mapa.');
            } finally {
                setLoading(false);
            }
        };

        geocodeAddress();
    }, [address, city, latitude, longitude]);

    if (loading) {
        return (
            <div
                style={{
                    height,
                    width: '100%',
                    borderRadius: '0.475rem',
                    background: 'linear-gradient(135deg, #1e1e2e 0%, #2a2a3e 100%)',
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '12px',
                    color: '#a0aec0',
                }}
            >
                <div
                    style={{
                        width: 36,
                        height: 36,
                        border: '3px solid #3b82f6',
                        borderTopColor: 'transparent',
                        borderRadius: '50%',
                        animation: 'spin 0.8s linear infinite',
                    }}
                />
                <span style={{ fontSize: '0.85rem' }}>Cargando mapa...</span>
                <style>{`@keyframes spin { to { transform: rotate(360deg); } }`}</style>
            </div>
        );
    }

    if (!position) {
        return (
            <div
                style={{
                    height,
                    width: '100%',
                    borderRadius: '0.475rem',
                    background: 'linear-gradient(135deg, #1e1e2e 0%, #2a2a3e 100%)',
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '8px',
                    color: '#718096',
                    fontSize: '0.875rem',
                }}
            >
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                    <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z" />
                    <circle cx="12" cy="9" r="2.5" />
                </svg>
                <span>{error || 'Ubicación no disponible'}</span>
            </div>
        );
    }

    const renderMap = (scrollZoom: boolean, zoom: number) => (
        <MapContainer center={position} zoom={zoom} scrollWheelZoom={scrollZoom} style={{ height: '100%', width: '100%' }}>
            <TileLayer
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            <Marker position={position}>
                <Popup>
                    {address}<br />{city}
                </Popup>
            </Marker>
            <RecenterAutomatically lat={position[0]} lng={position[1]} />
        </MapContainer>
    );

    return (
        <>
            <div style={{ position: 'relative', height, width: '100%', borderRadius: '0.475rem', overflow: 'hidden' }}>
                {renderMap(false, 15)}
                {expandable && (
                    <button
                        type="button"
                        onClick={() => setExpanded(true)}
                        title="Ampliar el mapa"
                        aria-label="Ampliar el mapa"
                        style={{
                            position: 'absolute',
                            top: 8,
                            right: 8,
                            zIndex: 500,
                            display: 'flex',
                            alignItems: 'center',
                            gap: 6,
                            padding: '6px 10px',
                            borderRadius: 8,
                            border: '1px solid rgba(148,163,184,0.35)',
                            background: 'rgba(15,23,42,0.82)',
                            color: '#e2e8f0',
                            fontSize: '0.72rem',
                            fontWeight: 600,
                            cursor: 'pointer',
                            backdropFilter: 'blur(4px)',
                        }}
                    >
                        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M15 3h6v6M9 21H3v-6M21 3l-7 7M3 21l7-7" />
                        </svg>
                        Ampliar
                    </button>
                )}
            </div>
            {error && <div style={{ color: '#ecc94b', fontSize: '0.8rem', marginTop: 6 }}>{error}</div>}

            {expandable && expanded && mounted && createPortal(
                <div
                    onClick={closeExpanded}
                    style={{
                        position: 'fixed',
                        inset: 0,
                        zIndex: 9999,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        padding: '4vh 4vw',
                        background: 'rgba(2,6,23,0.72)',
                        backdropFilter: 'blur(3px)',
                        animation: 'mapOverlayIn 180ms ease-out',
                    }}
                >
                    <div
                        onClick={(e) => e.stopPropagation()}
                        style={{
                            position: 'relative',
                            width: '100%',
                            maxWidth: 1100,
                            height: '100%',
                            maxHeight: 760,
                            borderRadius: 14,
                            overflow: 'hidden',
                            boxShadow: '0 24px 60px rgba(0,0,0,0.55)',
                            border: '1px solid rgba(148,163,184,0.25)',
                            transformOrigin: 'center',
                            animation: 'mapPanelIn 220ms cubic-bezier(0.16, 1, 0.3, 1)',
                        }}
                    >
                        {renderMap(true, 17)}
                        <div
                            style={{
                                position: 'absolute',
                                top: 12,
                                left: 12,
                                zIndex: 500,
                                maxWidth: 'calc(100% - 120px)',
                                padding: '8px 12px',
                                borderRadius: 10,
                                background: 'rgba(15,23,42,0.86)',
                                color: '#e2e8f0',
                                fontSize: '0.8rem',
                                lineHeight: 1.35,
                                backdropFilter: 'blur(4px)',
                            }}
                        >
                            <div style={{ fontWeight: 700 }}>{address}</div>
                            <div style={{ opacity: 0.75 }}>{city}</div>
                        </div>
                        <button
                            type="button"
                            onClick={closeExpanded}
                            title="Cerrar el mapa"
                            aria-label="Cerrar el mapa"
                            style={{
                                position: 'absolute',
                                top: 12,
                                right: 12,
                                zIndex: 500,
                                display: 'flex',
                                alignItems: 'center',
                                gap: 6,
                                padding: '8px 12px',
                                borderRadius: 10,
                                border: '1px solid rgba(148,163,184,0.35)',
                                background: 'rgba(15,23,42,0.86)',
                                color: '#e2e8f0',
                                fontSize: '0.78rem',
                                fontWeight: 600,
                                cursor: 'pointer',
                                backdropFilter: 'blur(4px)',
                            }}
                        >
                            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.4" strokeLinecap="round">
                                <path d="M18 6L6 18M6 6l12 12" />
                            </svg>
                            Cerrar
                        </button>
                    </div>
                    <style>{`
                        @keyframes mapOverlayIn { from { opacity: 0 } to { opacity: 1 } }
                        @keyframes mapPanelIn { from { opacity: 0; transform: scale(0.88) } to { opacity: 1; transform: scale(1) } }
                    `}</style>
                </div>,
                document.body,
            )}
        </>
    );
};

export default MapComponent;
