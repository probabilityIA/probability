'use client';

import { useMemo } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

interface LocationData {
    name: string;
    fullName: string;
    value: number;
}

interface MapGLComponentProps {
    data: LocationData[];
    height?: number;
}

// Colombia department coordinates (latitude, longitude)
const DEPARTMENT_COORDS: Record<string, [number, number]> = {
    'AMAZONAS': [3.1190, -71.6399],
    'ANTIOQUIA': [7.0000, -75.5000],
    'ARAUCA': [7.0783, -70.7596],
    'ATLÁNTICO': [10.9472, -74.7583],
    'BOGOTÁ': [4.7110, -74.0721],
    'BOLÍVAR': [10.3910, -75.5140],
    'BOYACÁ': [5.5300, -72.5300],
    'CALDAS': [5.2667, -75.5667],
    'CAQUETÁ': [2.6306, -72.8311],
    'CASANARE': [5.2981, -71.1781],
    'CAUCA': [2.5521, -76.4432],
    'CESAR': [10.2372, -73.1652],
    'CHOCÓ': [5.7321, -77.3144],
    'CÓRDOBA': [8.7500, -75.8830],
    'CUNDINAMARCA': [5.0000, -74.2500],
    'GUAINÍA': [3.0842, -67.8589],
    'GUAVIARE': [2.3045, -72.6407],
    'HUILA': [2.2667, -75.5000],
    'LA GUAJIRA': [11.5000, -72.6000],
    'MAGDALENA': [11.2381, -74.1921],
    'META': [3.8400, -72.3000],
    'NARIÑO': [1.2136, -77.2833],
    'NORTE DE SANTANDER': [7.8862, -72.6479],
    'PUTUMAYO': [1.0000, -76.5000],
    'QUINDÍO': [4.5306, -75.6794],
    'RISARALDA': [4.8128, -75.7300],
    'SANTANDER': [6.8000, -73.1500],
    'SUCRE': [9.3045, -74.7970],
    'TOLIMA': [4.7500, -75.3333],
    'VALLE DEL CAUCA': [4.5981, -76.0383],
    'VAUPÉS': [1.9425, -70.3742],
    'VICHADA': [5.6698, -68.1193],
};

interface MapGLComponentWithMetricsProps extends MapGLComponentProps {
    onDepartmentMapChange?: (map: Map<string, { count: number; percentage: number }>) => void;
}

export default function MapGLComponent({ data, height = 500, onDepartmentMapChange }: MapGLComponentWithMetricsProps) {
    const { departmentMap, totalOrders } = useMemo(() => {
        const total = data.reduce((sum, item) => sum + item.value, 0);
        const map = new Map<string, { count: number; percentage: number }>();

        // Normalizar variaciones de Bogotá
        const normalizeDepartment = (dept: string): string => {
            const normalized = dept.toUpperCase().trim();
            if (normalized.includes('D.C') || normalized.includes('D.D') || normalized.includes('S.C') ||
                normalized === 'BOGOTÁ' || normalized === 'BOGOTA' || normalized === 'DC' || normalized === 'DD' || normalized === 'SC') {
                return 'BOGOTÁ';
            }
            return normalized;
        };

        data.forEach(item => {
            const state = item.fullName.split(', ')[1] || item.fullName;
            const upperState = normalizeDepartment(state);

            if (map.has(upperState)) {
                const current = map.get(upperState)!;
                current.count += item.value;
                current.percentage = (current.count / total) * 100;
            } else {
                map.set(upperState, {
                    count: item.value,
                    percentage: total > 0 ? (item.value / total) * 100 : 0,
                });
            }
        });

        if (onDepartmentMapChange) {
            onDepartmentMapChange(map);
        }

        return { departmentMap: map, totalOrders: total };
    }, [data, onDepartmentMapChange]);

    const getColorForPercentage = (percentage: number): string => {
        if (percentage >= 80) return '#6B2FA1'; // Muy oscuro
        if (percentage >= 60) return '#8A5CB6'; // Oscuro
        if (percentage >= 40) return '#A987CB'; // Medio
        if (percentage >= 20) return '#C9B2E0'; // Claro
        return '#E8DDF5'; // Muy claro
    };

    return (
        <div style={{ height, position: 'relative', zIndex: 0 }} className="relative w-full overflow-hidden">
            <MapContainer
                center={[4.5709, -74.2973]}
                zoom={5}
                style={{ width: '100%', height: '100%', position: 'absolute', top: 0, left: 0 }}
                className="rounded-lg"
            >
                <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
                />

                {Array.from(departmentMap.entries()).map(([dept, stats]) => {
                    const coords = DEPARTMENT_COORDS[dept];
                    if (!coords) return null;

                    const icon = L.divIcon({
                        html: `<div style="
                            background-color: ${getColorForPercentage(stats.percentage)};
                            width: 30px;
                            height: 30px;
                            border-radius: 50%;
                            border: 2px solid white;
                            box-shadow: 0 2px 8px rgba(0,0,0,0.3);
                            display: flex;
                            align-items: center;
                            justify-content: center;
                            font-weight: bold;
                            font-size: 12px;
                            color: white;
                        "></div>`,
                        iconSize: [30, 30],
                        className: '',
                    });

                    return (
                        <Marker
                            key={dept}
                            position={[coords[0], coords[1]]}
                            icon={icon}
                        >
                            <Popup>
                                <div className="text-center">
                                    <p className="font-semibold">{dept}</p>
                                    <p className="text-sm">{stats.count} órdenes</p>
                                    <p className="text-sm text-gray-600 dark:text-gray-300">{stats.percentage.toFixed(1)}%</p>
                                </div>
                            </Popup>
                        </Marker>
                    );
                })}
            </MapContainer>

        </div>
    );
}
