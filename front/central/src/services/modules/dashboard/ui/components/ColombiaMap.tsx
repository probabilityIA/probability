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
    return (
        <div style={{ height }} className="rounded-lg overflow-hidden border border-gray-200 relative bg-gray-100">
            <MapGLComponent data={data} height={height} />
        </div>
    );
}
