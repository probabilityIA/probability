'use client';

import { useEffect } from 'react';

export function GoogleMapsLoader() {
    useEffect(() => {
        const apiKey = process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY;

        console.log('🗺️ GoogleMapsLoader initialized');
        console.log('API Key available:', !!apiKey);

        if (!apiKey) {
            console.error('❌ Google Maps API key not found in environment variables');
            return;
        }

        if (window.google && window.google.maps) {
            console.log('✅ Google Maps already loaded');
            return;
        }

        console.log('📝 Creating script tag for Google Maps...');
        const script = document.createElement('script');
        script.src = `https://maps.googleapis.com/maps/api/js?key=${apiKey}&libraries=places&language=es`;
        script.async = true;
        script.defer = false;

        script.onload = () => {
            console.log('✅ Google Maps API loaded successfully');
            window.dispatchEvent(new Event('google-maps-loaded'));
        };

        script.onerror = (error) => {
            console.error('❌ Failed to load Google Maps API:', error);
        };

        console.log('📍 Appending script to head');
        document.head.appendChild(script);

        return () => {
            if (script && script.parentNode) {
                script.parentNode.removeChild(script);
            }
        };
    }, []);

    return null;
}
