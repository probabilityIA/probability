'use client';

import { MapContainer, TileLayer, Popup } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

interface QuotationMapProps {
    originAddress: string;
    originCity: string;
    destAddress: string;
    destCity: string;
}

const colombiaCoords = [4.5709, -74.2973];

export default function QuotationMap({ originAddress, originCity, destAddress, destCity }: QuotationMapProps) {
    return (
        <div className="flex-1 rounded-lg overflow-hidden bg-gray-100 dark:bg-gray-700 min-h-0" style={{ height: '220px' }}>
            <MapContainer
                center={[colombiaCoords[0], colombiaCoords[1]]}
                zoom={6}
                style={{ height: '100%', width: '100%' }}
            >
                <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='&copy; OpenStreetMap'
                />
            </MapContainer>
        </div>
    );
}
