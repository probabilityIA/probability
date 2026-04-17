'use client';

import dynamic from 'next/dynamic';
import { MapPin } from 'lucide-react';

const MapComponent = dynamic(() => import('@/shared/ui/MapComponent'), {
    ssr: false,
    loading: () => (
        <div className="h-36 bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
            <MapPin size={16} className="text-gray-400 animate-pulse" />
        </div>
    ),
});

interface MiniAddressMapProps {
    address?: string;
    city?: string;
    color?: 'blue' | 'emerald';
}

function extractCity(address?: string): string {
    if (!address) return 'Colombia';
    const parts = address.split(',').map(s => s.trim()).filter(Boolean);
    if (parts.length >= 2) {
        const last = parts[parts.length - 1];
        if (last.toLowerCase() === 'colombia' && parts.length >= 3) {
            return parts[parts.length - 2];
        }
        return parts[parts.length - 1];
    }
    return 'Colombia';
}

export function MiniAddressMap({ address, city, color = 'blue' }: MiniAddressMapProps) {
    const bgClass = color === 'emerald' ? 'bg-emerald-100 dark:bg-emerald-900/30' : 'bg-blue-100 dark:bg-blue-900/30';
    const iconColor = color === 'emerald' ? 'text-emerald-500' : 'text-blue-500';

    if (!address) {
        return (
            <div className={`h-36 ${bgClass} flex items-center justify-center`}>
                <MapPin size={18} className={`${iconColor} opacity-40`} />
            </div>
        );
    }

    const resolvedCity = city || extractCity(address);

    return (
        <div className="h-36 pointer-events-none">
            <MapComponent address={address} city={resolvedCity} height="80px" />
        </div>
    );
}
