'use client';

import { useRef, useEffect } from 'react';
import { CustomerInfo } from '../../../customers/domain/types';

interface ClientAutocompleteProps {
    results: CustomerInfo[];
    loading: boolean;
    /** Whether a search has completed (used to show "not found") */
    searched: boolean;
    visible: boolean;
    searchTerm?: string;
    onSelect: (client: CustomerInfo) => void;
    onClose: () => void;
}

function highlightMatch(text: string, term: string) {
    if (!term || !text) return text;
    const idx = text.toLowerCase().indexOf(term.toLowerCase());
    if (idx === -1) return text;
    return (
        <>
            {text.slice(0, idx)}
            <span className="bg-purple-100 text-purple-800 font-semibold">{text.slice(idx, idx + term.length)}</span>
            {text.slice(idx + term.length)}
        </>
    );
}

export default function ClientAutocomplete({
    results,
    loading,
    searched,
    visible,
    searchTerm = '',
    onSelect,
    onClose,
}: ClientAutocompleteProps) {
    const ref = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (ref.current && !ref.current.contains(event.target as Node)) {
                onClose();
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, [onClose]);

    // Show when: visible AND (loading OR has results OR searched-with-no-results)
    const showNotFound = searched && !loading && results.length === 0;
    if (!visible || (!loading && results.length === 0 && !showNotFound)) {
        return null;
    }

    return (
        <div
            ref={ref}
            className="absolute z-20 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg max-h-56 overflow-y-auto"
        >
            {/* Loading state */}
            {loading && (
                <div className="px-3 py-3 flex items-center gap-2.5">
                    <div className="w-4 h-4 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" />
                    <span className="text-sm text-gray-500">Buscando cliente...</span>
                </div>
            )}

            {/* No results */}
            {showNotFound && (
                <div className="px-3 py-3 flex items-center gap-2">
                    <svg className="w-4 h-4 text-gray-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span className="text-sm text-gray-500">
                        Cliente no encontrado — se creará uno nuevo al guardar la orden
                    </span>
                </div>
            )}

            {/* Results */}
            {results.map((client) => (
                <button
                    key={client.id}
                    type="button"
                    onClick={() => onSelect(client)}
                    className="w-full text-left px-3 py-2.5 hover:bg-purple-50 cursor-pointer border-b border-gray-100 last:border-b-0 transition-colors"
                >
                    <div className="flex items-center justify-between">
                        <span className="font-medium text-sm text-gray-800">{client.name}</span>
                        {client.dni && (
                            <span className="text-xs font-mono bg-gray-100 text-gray-600 px-1.5 py-0.5 rounded">
                                {highlightMatch(client.dni, searchTerm)}
                            </span>
                        )}
                    </div>
                    <div className="flex items-center gap-3 mt-0.5">
                        {client.email && (
                            <span className="text-xs text-gray-500">{client.email}</span>
                        )}
                        {client.phone && (
                            <span className="text-xs text-gray-400">{client.phone}</span>
                        )}
                    </div>
                </button>
            ))}
        </div>
    );
}
