'use client';

import dynamic from 'next/dynamic';
import { useMemo, useState } from 'react';

const MapGLComponent = dynamic(() => import('./MapGLComponent'), {
    ssr: false,
    loading: () => <div className="w-full h-full flex items-center justify-center bg-gray-100">Cargando mapa...</div>
});

interface LocationData {
    name: string;
    fullName: string;
    value: number;
}

interface ColombiaMapProps {
    data: LocationData[];
    height?: number;
}

export function ColombiaMap({ data, height = 500 }: ColombiaMapProps) {
    const [departmentMap, setDepartmentMap] = useState<Map<string, { count: number; percentage: number }> | null>(null);

    return (
        <div className="space-y-4">
            {/* Contenedor del mapa */}
            <div
                style={{
                    height,
                    position: 'relative',
                    zIndex: 0,
                    isolation: 'isolate'
                }}
                className="rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700 bg-gray-100"
            >
                <MapGLComponent data={data} height={height} onDepartmentMapChange={setDepartmentMap} />
            </div>

            {/* Métricas debajo del mapa */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Escala de Órdenes */}
                <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow border border-gray-200 dark:border-gray-700">
                    <p className="text-sm font-semibold mb-3 text-gray-700 dark:text-gray-200">Escala de Órdenes</p>
                    <div className="space-y-2 text-sm">
                        <div className="flex items-center gap-2">
                            <div className="w-3 h-3 rounded-full" style={{ backgroundColor: '#6B2FA1' }}></div>
                            <span className="text-gray-600 dark:text-gray-300">80-100%</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <div className="w-3 h-3 rounded-full" style={{ backgroundColor: '#8A5CB6' }}></div>
                            <span className="text-gray-600 dark:text-gray-300">60-80%</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <div className="w-3 h-3 rounded-full" style={{ backgroundColor: '#A987CB' }}></div>
                            <span className="text-gray-600 dark:text-gray-300">40-60%</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <div className="w-3 h-3 rounded-full" style={{ backgroundColor: '#C9B2E0' }}></div>
                            <span className="text-gray-600 dark:text-gray-300">20-40%</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <div className="w-3 h-3 rounded-full" style={{ backgroundColor: '#E8DDF5' }}></div>
                            <span className="text-gray-600 dark:text-gray-300">0-20%</span>
                        </div>
                    </div>
                </div>

                {/* Top Departamentos */}
                <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow border border-gray-200 dark:border-gray-700">
                    <p className="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">Top Departamentos</p>
                    <div className="text-sm text-gray-600 dark:text-gray-300 space-y-2">
                        {departmentMap && Array.from(departmentMap.entries())
                            .sort((a, b) => b[1].count - a[1].count)
                            .slice(0, 5)
                            .map(([dept, stats]) => (
                                <div key={dept} className="flex justify-between items-center">
                                    <span className="text-gray-700 dark:text-gray-200">{dept}</span>
                                    <span className="font-semibold text-purple-600">{stats.percentage.toFixed(1)}%</span>
                                </div>
                            ))}
                        {!departmentMap && <p className="text-gray-500 dark:text-gray-400">Cargando datos...</p>}
                    </div>
                </div>
            </div>
        </div>
    );
}
