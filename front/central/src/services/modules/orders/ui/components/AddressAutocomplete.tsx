'use client';

import { useState, useRef, useEffect, useCallback } from 'react';

export interface AddressSuggestion {
    display_name: string;
    street: string;
    house_number: string;
    neighbourhood: string;
    city: string;
    state: string;
    postcode: string;
    lat: number;
    lon: number;
}

interface AddressAutocompleteProps {
    value: string;
    onChange: (value: string) => void;
    onSelect: (suggestion: AddressSuggestion) => void;
    placeholder?: string;
    country?: string;
}

export default function AddressAutocomplete({
    value,
    onChange,
    onSelect,
    placeholder = 'Calle/Carrera número',
    country = 'co',
}: AddressAutocompleteProps) {
    const [suggestions, setSuggestions] = useState<AddressSuggestion[]>([]);
    const [showDropdown, setShowDropdown] = useState(false);
    const [loading, setLoading] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);
    const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const searchAddress = useCallback(async (query: string) => {
        if (query.length < 4) {
            setSuggestions([]);
            return;
        }

        setLoading(true);
        try {
            const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
            const params = new URLSearchParams({ q: query, country });
            const response = await fetch(`${apiBase}/address-search?${params}`);

            if (!response.ok) return;

            const data: AddressSuggestion[] = await response.json();
            setSuggestions(data);
            setShowDropdown(data.length > 0);
        } catch {
            setSuggestions([]);
        } finally {
            setLoading(false);
        }
    }, [country]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const val = e.target.value;
        onChange(val);

        if (debounceRef.current) clearTimeout(debounceRef.current);
        debounceRef.current = setTimeout(() => searchAddress(val), 400);
    };

    const handleSelect = (suggestion: AddressSuggestion) => {
        const street = suggestion.house_number
            ? `${suggestion.street} ${suggestion.house_number}`
            : suggestion.street;
        onChange(street || suggestion.display_name.split(',')[0]);
        setShowDropdown(false);
        setSuggestions([]);
        onSelect(suggestion);
    };

    // Close on click outside
    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
                setShowDropdown(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    return (
        <div ref={containerRef} className="relative">
            <div className="relative">
                <input
                    type="text"
                    value={value}
                    onChange={handleChange}
                    onFocus={() => suggestions.length > 0 && setShowDropdown(true)}
                    placeholder={placeholder}
                    className="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent text-black"
                />
                {loading && (
                    <div className="absolute right-3 top-1/2 -translate-y-1/2">
                        <div className="w-4 h-4 border-2 border-purple-400 border-t-transparent rounded-full animate-spin" />
                    </div>
                )}
            </div>

            {showDropdown && suggestions.length > 0 && (
                <div className="absolute z-20 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                    {suggestions.map((s, i) => (
                        <button
                            key={i}
                            type="button"
                            onClick={() => handleSelect(s)}
                            className="w-full text-left px-3 py-2.5 hover:bg-purple-50 cursor-pointer border-b border-gray-100 last:border-b-0 transition-colors"
                        >
                            <div className="flex items-start gap-2">
                                <svg className="w-4 h-4 text-purple-400 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                                </svg>
                                <div className="min-w-0">
                                    <p className="text-sm text-gray-800 font-medium truncate">
                                        {s.street ? `${s.street} ${s.house_number}`.trim() : s.display_name.split(',')[0]}
                                    </p>
                                    <p className="text-xs text-gray-500 truncate">
                                        {s.neighbourhood ? `${s.neighbourhood}, ` : ''}{s.city}{s.state ? `, ${s.state}` : ''}
                                    </p>
                                </div>
                            </div>
                        </button>
                    ))}
                </div>
            )}
        </div>
    );
}
