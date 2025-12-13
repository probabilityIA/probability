import React, { useEffect, useState } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';

// Fix for default marker icon in React Leaflet
import icon from 'leaflet/dist/images/marker-icon.png';
import iconShadow from 'leaflet/dist/images/marker-shadow.png';

let DefaultIcon = L.icon({
    iconUrl: icon,
    shadowUrl: iconShadow,
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

            const query = `${address}, ${city}, Colombia`;
            const url = `https://nominatim.openstreetmap.org/search?format=json&addressdetails=1&q=${encodeURIComponent(query)}`;

            try {
                const response = await fetch(url, {
                    headers: { "User-Agent": "ProbabilityApp_v1.0" }
                });
                const data = await response.json();

                if (data && data.length > 0) {
                    const lat = parseFloat(data[0].lat);
                    const lon = parseFloat(data[0].lon);
                    setPosition([lat, lon]);
                } else {
                    // Fallback to just city
                    const cityUrl = `https://nominatim.openstreetmap.org/search?format=json&addressdetails=1&q=${encodeURIComponent(city + ", Colombia")}`;
                    const cityRes = await fetch(cityUrl, { headers: { "User-Agent": "ProbabilityApp_v1.0" } });
                    const cityData = await cityRes.json();

                    if (cityData && cityData.length > 0) {
                        const lat = parseFloat(cityData[0].lat);
                        const lon = parseFloat(cityData[0].lon);
                        setPosition([lat, lon]);
                        setError("Dirección exacta no encontrada, mostrando ubicación de la ciudad.");
                    } else {
                        setError("No se pudo localizar la dirección.");
                    }
                }
            } catch (err) {
                console.error("Geocoding error:", err);
                setError("Error al cargar el mapa.");
            } finally {
                setLoading(false);
            }
        };

        geocodeAddress();
    }, [address, city]);

    if (loading) {
        return (
            <div className="d-flex align-items-center justify-content-center bg-light rounded" style={{ height }}>
                <div className="spinner-border text-primary" role="status">
                    <span className="visually-hidden">Cargando mapa...</span>
                </div>
            </div>
        );
    }

    if (!position) {
        return (
            <div className="d-flex align-items-center justify-content-center bg-light rounded text-muted" style={{ height }}>
                {error || "Ubicación no disponible"}
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
            {error && <div className="text-warning small mt-1">{error}</div>}
        </div>
    );
};

export default MapComponent;
