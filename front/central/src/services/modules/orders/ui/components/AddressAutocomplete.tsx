'use client';

import { useEffect, useRef, useState } from 'react';

export interface AddressSuggestion {
    address: string;
    city?: string;
    state?: string;
    postcode?: string;
    neighbourhood?: string;
    lat?: number;
    lon?: number;
}

interface AddressAutocompleteProps {
    value: string;
    onChange: (value: string) => void;
    onSelect?: (suggestion: AddressSuggestion) => void;
    city?: string;
    placeholder?: string;
    className?: string;
}

export default function AddressAutocomplete({
    value,
    onChange,
    onSelect,
    city = '',
    placeholder = 'Calle y número',
    className = '',
}: AddressAutocompleteProps) {
    const [predictions, setPredictions] = useState<any[]>([]);
    const [showPredictions, setShowPredictions] = useState(false);
    const [loading, setLoading] = useState(false);
    const autocompleteServiceRef = useRef<google.maps.places.AutocompleteService | null>(null);
    const placesServiceRef = useRef<google.maps.places.PlacesService | null>(null);
    const containerRef = useRef<HTMLDivElement>(null);
    const mapRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const loadGoogleMaps = async () => {
            if (window.google && window.google.maps) {
                console.log('✅ Google Maps already loaded');
                initializeServices();
                return;
            }

            const apiKey = process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY;
            console.log('🔍 Checking for API key:', !!apiKey);

            if (!apiKey) {
                console.error('❌ No API key found');
                return;
            }

            if (document.querySelector('script[src*="maps.googleapis.com"]')) {
                console.log('⏳ Google Maps script already loading...');
                let attempts = 0;
                const checkInterval = setInterval(() => {
                    if (window.google && window.google.maps) {
                        clearInterval(checkInterval);
                        initializeServices();
                    }
                    attempts++;
                    if (attempts > 100) {
                        clearInterval(checkInterval);
                        console.warn('Google Maps failed to load');
                    }
                }, 100);
                return;
            }

            console.log('📝 Loading Google Maps script...');
            const script = document.createElement('script');
            script.src = `https://maps.googleapis.com/maps/api/js?key=${apiKey}&libraries=places&language=es`;
            script.async = true;

            script.onload = () => {
                console.log('✅ Google Maps script loaded');
                initializeServices();
            };

            script.onerror = (error) => {
                console.error('❌ Failed to load Google Maps:', error);
            };

            document.head.appendChild(script);
        };

        const initializeServices = () => {
            if (!window.google || !window.google.maps) {
                console.error('❌ window.google.maps not available');
                return;
            }

            if (!autocompleteServiceRef.current) {
                try {
                    autocompleteServiceRef.current = new window.google.maps.places.AutocompleteService();
                    console.log('✅ AutocompleteService initialized');
                } catch (e) {
                    console.error('❌ Error initializing AutocompleteService:', e);
                }
            }

            if (!placesServiceRef.current && mapRef.current) {
                try {
                    const hiddenMap = new window.google.maps.Map(mapRef.current, {
                        center: { lat: 4.5709, lng: -74.2973 },
                        zoom: 8,
                    });
                    placesServiceRef.current = new window.google.maps.places.PlacesService(hiddenMap);
                    console.log('✅ PlacesService initialized');
                } catch (e) {
                    console.error('❌ Error initializing PlacesService:', e);
                }
            }
        };

        loadGoogleMaps();
    }, []);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
                setShowPredictions(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleInputChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const val = e.target.value;
        onChange(val);

        if (!val || val.length < 3) {
            setPredictions([]);
            setShowPredictions(false);
            return;
        }

        if (!autocompleteServiceRef.current) {
            if (!window.google || !window.google.maps || !window.google.maps.places) {
                console.warn('Google Maps API not ready yet');
                setPredictions([]);
                return;
            }
            try {
                autocompleteServiceRef.current = new window.google.maps.places.AutocompleteService();
            } catch (e) {
                console.error('Error initializing AutocompleteService:', e);
                return;
            }
        }

        setLoading(true);
        try {
            console.log('🔍 Fetching predictions for:', val);
            const results = await autocompleteServiceRef.current.getPlacePredictions({
                input: val,
                componentRestrictions: { country: 'co' },
            });

            console.log('📋 Got predictions:', results.predictions?.length ?? 0);
            setPredictions(results.predictions || []);
            setShowPredictions(true);
        } catch (error) {
            console.error('❌ Error fetching predictions:', error);
            setPredictions([]);
        } finally {
            setLoading(false);
        }
    };

    const handleSelectPrediction = async (prediction: any) => {
        console.log('🎯 Selected prediction:', prediction);
        onChange(prediction.description);
        setShowPredictions(false);
        setPredictions([]);

        if (!placesServiceRef.current || !mapRef.current) {
            console.log('⚠️ PlacesService not ready, calling onSelect with basic data');
            onSelect?.({
                address: prediction.description,
            });
            return;
        }

        try {
            placesServiceRef.current.getDetails(
                { placeId: prediction.place_id },
                (place: google.maps.places.PlaceResult | null) => {
                    if (!place || !place.formatted_address) return;

                    const suggestion: AddressSuggestion = {
                        address: place.formatted_address,
                        lat: place.geometry?.location?.lat(),
                        lon: place.geometry?.location?.lng(),
                    };

                    const addressComponents = place.address_components || [];
                    addressComponents.forEach(component => {
                        if (component.types.includes('locality')) {
                            suggestion.city = component.long_name;
                        }
                        if (component.types.includes('administrative_area_level_1')) {
                            suggestion.state = component.long_name;
                        }
                        if (component.types.includes('postal_code')) {
                            suggestion.postcode = component.long_name;
                        }
                        if (component.types.includes('neighborhood')) {
                            suggestion.neighbourhood = component.long_name;
                        }
                    });

                    console.log('✅ Calling onSelect with suggestion:', suggestion);
                    onSelect?.(suggestion);
                }
            );
        } catch (error) {
            console.error('Error getting place details:', error);
            onSelect?.({
                address: prediction.description,
            });
        }
    };

    return (
        <>
            <div ref={containerRef} className="relative">
                <input
                    type="text"
                    value={value}
                    onChange={handleInputChange}
                    onFocus={() => value && predictions.length > 0 && setShowPredictions(true)}
                    placeholder={placeholder}
                    className={`shipment-input ${className}`}
                    autoComplete="off"
                />
                {loading && (
                    <div className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
                        <div className="w-4 h-4 border-2 border-gray-300 border-t-gray-600 rounded-full animate-spin" />
                    </div>
                )}
                {showPredictions && predictions.length > 0 && (
                    <div className="absolute z-10 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md shadow-lg max-h-60 overflow-y-auto">
                        {predictions.map((pred) => (
                            <div
                                key={pred.place_id}
                                onClick={() => handleSelectPrediction(pred)}
                                className="px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer text-sm"
                            >
                                <div className="font-medium text-gray-900 dark:text-white">
                                    {pred.structured_formatting?.main_text || pred.description}
                                </div>
                                {(pred.structured_formatting?.secondary_text || pred.description) && (
                                    <div className="text-xs text-gray-500 dark:text-gray-400">
                                        {pred.structured_formatting?.secondary_text}
                                    </div>
                                )}
                            </div>
                        ))}
                    </div>
                )}
            </div>
            <div ref={mapRef} style={{ display: 'none' }} />
        </>
    );
}
