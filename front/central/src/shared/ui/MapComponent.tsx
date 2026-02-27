import React, { useEffect, useState } from 'react';
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
}

const RecenterAutomatically = ({ lat, lng }: { lat: number; lng: number }) => {
    const map = useMap();
    useEffect(() => {
        map.setView([lat, lng]);
    }, [lat, lng, map]);
    return null;
};

const MapComponent: React.FC<MapComponentProps> = ({ address, city, height = '400px' }) => {
    const [position, setPosition] = useState<[number, number] | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const geocodeAddress = async () => {
            if (!address || !city) return;

            setLoading(true);
            setError(null);

            try {
                // Llamamos a nuestro propio backend como proxy para evitar restricciones CORS/User-Agent
                // El backend hace la petición a Nominatim server-side sin restricciones
                const backendBase = process.env.NEXT_PUBLIC_API_BASE_URL?.replace('/api/v1', '') || 'http://localhost:3050';
                const url = `${backendBase}/geocode?address=${encodeURIComponent(address)}&city=${encodeURIComponent(city)}`;

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
    }, [address, city]);

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

    return (
        <div style={{ height, width: '100%', borderRadius: '0.475rem', overflow: 'hidden' }}>
            <MapContainer center={position} zoom={15} scrollWheelZoom={false} style={{ height: '100%', width: '100%' }}>
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
            {error && <div style={{ color: '#ecc94b', fontSize: '0.8rem', marginTop: 6 }}>{error}</div>}
        </div>
    );
};

export default MapComponent;
