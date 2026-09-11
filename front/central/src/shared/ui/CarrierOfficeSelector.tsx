import React, { useState, useEffect } from 'react';
import { env } from '@/shared/config/env';
import { MapPin, Search, Loader2 } from 'lucide-react';
import { Input, Button } from '@/shared/ui';

interface OfficeResult {
    display_name: string;
    place_id: string;
    lat: number;
    lon: number;
    distance_km?: number;
}

interface CarrierOfficeSelectorProps {
    city: string;
    state?: string;
    initialCarrier?: string;
    selectedAddress?: string;
    onSelectAddress: (address: string, carrierId: string, coords?: { lat: number; lon: number }) => void;
    onClose: () => void;
}

const CARRIERS = [
    { id: 'coordinadora', name: 'Coordinadora' },
    { id: 'interrapidisimo', name: 'Inter Rapidísimo' },
    { id: 'servientrega', name: 'Servientrega' },
    { id: 'envia', name: 'Envía' },
    { id: 'tcc', name: 'TCC' }
];

export function CarrierOfficeSelector({ city, state = '', initialCarrier, selectedAddress = '', onSelectAddress, onClose }: CarrierOfficeSelectorProps) {
    const [selectedCarrier, setSelectedCarrier] = useState<string>(
        initialCarrier && CARRIERS.some((c) => c.id === initialCarrier) ? initialCarrier : CARRIERS[0].id,
    );
    const [results, setResults] = useState<OfficeResult[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [subSearch, setSubSearch] = useState<string>('');
    const cleanCityLabel = city.split('(')[0].trim();
    const normalizeAddr = (t: string) => t.normalize('NFD').replace(/[\u0300-\u036f]/g, '').toLowerCase().replace(/[^a-z0-9]/g, '');
    const selectedKey = normalizeAddr(selectedAddress);
    const isSelected = (displayName: string) => {
        if (!selectedKey) return false;
        const inner = displayName.match(/\((.*?)\)/)?.[1] || displayName;
        const key = normalizeAddr(inner);
        return key.includes(selectedKey) || selectedKey.includes(key);
    };

    const fetchOffices = async () => {
        if (!city) {
            setError('Por favor selecciona una ciudad primero');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const carrierName = CARRIERS.find(c => c.id === selectedCarrier)?.name || selectedCarrier;
            // Quitamos lo que esté en paréntesis de la ciudad if it looks like "Bogota (Bogota)"
            const cleanCity = city.split('(')[0].trim();
            const query = encodeURIComponent(`Oficina ${carrierName} ${subSearch} ${cleanCity} Colombia`);
            const params = new URLSearchParams({ query: decodeURIComponent(query), city: cleanCity, carrier: carrierName });
            if (state) params.set('state', state);

            const response = await fetch(`${env.API_BASE_URL}/places-search?${params.toString()}`);
            if (!response.ok) throw new Error('Error al buscar oficinas');
            
            const data = await response.json();
            if (data && data.length > 0) {
                setResults(data);
            } else {
                setResults([]);
                setError(`No se encontraron oficinas de ${carrierName} en ${cleanCity}`);
            }
        } catch (err: any) {
            setError(err.message || 'Ocurrió un error al buscar oficinas');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const timeoutId = setTimeout(() => {
            fetchOffices();
        }, 500);
        return () => clearTimeout(timeoutId);
    }, [selectedCarrier, city, state, subSearch]);

    return (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 p-4 mt-2">
            <div className="flex justify-between items-center mb-3">
                <h4 className="text-sm font-bold text-gray-700 dark:text-gray-200 flex items-center gap-2">
                    <MapPin size={16} className="text-purple-500" />
                    Buscar Oficina de Transportadora
                </h4>
                <button onClick={onClose} className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200">
                    &times;
                </button>
            </div>

            <p className="text-xs text-gray-500 dark:text-gray-400 mb-3">
                Selecciona una transportadora para ver sus sedes principales en <strong>{city}</strong>
            </p>
            {selectedAddress && !results.some((r) => isSelected(r.display_name)) && (
                <div className="mb-3 rounded-md border border-purple-300 bg-purple-50 px-2.5 py-2 dark:border-purple-500/40 dark:bg-purple-900/20">
                    <p className="text-[10px] font-bold uppercase tracking-wide text-purple-700 dark:text-purple-300">
                        Oficina elegida en la orden
                    </p>
                    <p className="text-xs text-gray-700 dark:text-gray-200">{selectedAddress}</p>
                </div>
            )}

            <div className="flex gap-2 mb-4 overflow-x-auto pb-1">
                {CARRIERS.map(carrier => (
                    <button
                        key={carrier.id}
                        type="button"
                        onClick={() => setSelectedCarrier(carrier.id)}
                        className={`px-3 py-1.5 text-xs font-semibold rounded-full whitespace-nowrap transition-colors ${
                            selectedCarrier === carrier.id 
                            ? 'bg-purple-100 text-purple-700 border border-purple-300 dark:bg-purple-900/30 dark:border-purple-600/50' 
                            : 'bg-gray-50 text-gray-600 border border-gray-200 hover:bg-gray-100 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600'
                        }`}
                    >
                        {carrier.name}
                    </button>
                ))}
            </div>

            <div className="relative mb-4">
                <Input 
                    placeholder="Busca por barrio o dirección específica..." 
                    value={subSearch}
                    onChange={(e) => setSubSearch(e.target.value)}
                    className="pl-9 text-xs"
                />
                <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            </div>

            <div className="min-h-[150px] max-h-[250px] overflow-y-auto pr-1">
                {loading ? (
                    <div className="h-32 flex flex-col items-center justify-center text-gray-400 gap-2">
                        <Loader2 className="animate-spin" size={24} />
                        <span className="text-xs">Buscando oficinas...</span>
                    </div>
                ) : error ? (
                    <div className="text-sm text-amber-600 bg-amber-50 p-3 rounded text-center">
                        {error}
                    </div>
                ) : results.length > 0 ? (
                    <div className="space-y-2">
                        {[...results].sort((a, b) => Number(isSelected(b.display_name)) - Number(isSelected(a.display_name))).map((r, i) => (
                            <div 
                                key={r.place_id || i}
                                onClick={() => {
                                    // Limpiamos los textos como "Nombre Oficina (Dirección)" a solo "Dirección" si es posible,
                                    // o usamos todo el display name como address
                                    let address = r.display_name;
                                    const match = address.match(/\((.*?)\)/);
                                    if (match && match[1]) {
                                        address = match[1];
                                    }
                                    const coords = r.lat != null && r.lon != null ? { lat: r.lat, lon: r.lon } : undefined;
                                    onSelectAddress(address, selectedCarrier, coords);
                                    onClose();
                                }}
                                className={`group p-3 border rounded-md cursor-pointer transition-colors ${
                                    isSelected(r.display_name)
                                        ? 'border-purple-500 bg-purple-50 dark:bg-purple-900/20 ring-1 ring-purple-400'
                                        : 'border-gray-100 dark:border-gray-700 hover:border-purple-300 dark:hover:border-purple-500 hover:bg-purple-50 dark:hover:bg-purple-900/10'
                                }`}
                            >
                                <div className="flex gap-3 items-start">
                                    <div className="bg-gray-100 dark:bg-gray-700 p-1.5 rounded-full mt-0.5 group-hover:bg-purple-100 dark:group-hover:bg-purple-800">
                                        <MapPin size={14} className="text-gray-500 dark:text-gray-400 group-hover:text-purple-600 dark:group-hover:text-purple-300" />
                                    </div>
                                    <div>
                                        <p className="text-xs font-semibold text-gray-800 dark:text-gray-200 leading-snug">
                                            {isSelected(r.display_name) && (
                                                <span className="mr-1.5 rounded bg-purple-600 px-1.5 py-0.5 text-[9px] font-bold text-white align-middle">
                                                    ELEGIDA
                                                </span>
                                            )}
                                            {r.display_name.split('(')[0].trim()}
                                            {typeof r.distance_km === 'number' && r.distance_km >= 1 && (
                                                <span className="ml-1.5 font-normal text-[10px] text-amber-600 dark:text-amber-400">
                                                    a {r.distance_km} km del centro de {cleanCityLabel}
                                                </span>
                                            )}
                                        </p>
                                        {r.display_name.includes('(') && (
                                            <p className="text-[10px] text-gray-500 dark:text-gray-400 mt-0.5">
                                                {r.display_name.match(/\((.*?)\)/)?.[1]}
                                            </p>
                                        )}
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                ) : null}
            </div>
        </div>
    );
}
