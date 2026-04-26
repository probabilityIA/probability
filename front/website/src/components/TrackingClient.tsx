/** @jsxImportSource react */
import { useEffect, useRef, useState } from 'react';
import TrackingSearchInput from './tracking/TrackingSearchInput';
import TrackingProgressBar from './tracking/TrackingProgressBar';
import TrackingDetails from './tracking/TrackingDetails';
import TrackingTimeline from './tracking/TrackingTimeline';
import type { TrackingSearchResult, TrackingHistory } from '../types/tracking';
import { getApiUrl } from '../config/api';

interface SearchResult {
  success: boolean;
  message: string;
  data?: {
    shipment?: TrackingSearchResult;
    history?: TrackingHistory[];
  };
}

export default function TrackingClient() {
  const [shipment, setShipment] = useState<TrackingSearchResult | null>(null);
  const [history, setHistory] = useState<TrackingHistory[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [initialQuery, setInitialQuery] = useState<string>('');
  const autoTriggered = useRef(false);

  useEffect(() => {
    if (autoTriggered.current) return;
    if (typeof window === 'undefined') return;
    const params = new URLSearchParams(window.location.search);
    const tracking = params.get('tracking') || params.get('q');
    if (tracking && tracking.trim()) {
      autoTriggered.current = true;
      setInitialQuery(tracking.trim());
      handleSearch(tracking.trim());
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleSearch = async (query: string) => {
    if (!query.trim()) return;

    setIsLoading(true);
    setError(null);
    setShipment(null);
    setHistory([]);

    try {
      // Llamar al endpoint del backend
      const apiUrl = getApiUrl();
      const response = await fetch(
        `${apiUrl}/tracking/search?tracking_number=${encodeURIComponent(query)}`
      );

      if (!response.ok) {
        throw new Error(`Error: ${response.status}`);
      }

      const result: SearchResult = await response.json();

      if (!result.success) {
        setError(result.message || 'No se encontró información del envío');
        return;
      }

      const foundShipment = result.data?.shipment;
      if (!foundShipment) {
        setError('No se encontró información del envío');
        return;
      }

      setShipment(foundShipment);

      // Obtener historial si hay tracking_number
      if (foundShipment.tracking_number) {
        try {
          const apiUrl = getApiUrl();
          const historyResponse = await fetch(
            `${apiUrl}/tracking/${encodeURIComponent(foundShipment.tracking_number)}/history`
          );

          if (historyResponse.ok) {
            const historyResult: SearchResult = await historyResponse.json();
            if (historyResult.data?.history) {
              setHistory(historyResult.data.history);
            }
          }
        } catch (err) {
          console.error('Error loading history:', err);
        }
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Error al buscar el envío';
      setError(message);
      setShipment(null);
      setHistory([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleReset = () => {
    setShipment(null);
    setHistory([]);
    setError(null);
  };

  return (
    <div class="space-y-8">
      {/* Search Section */}
      <div class="bg-white rounded-2xl shadow-lg p-8">
        <TrackingSearchInput onSearch={handleSearch} isLoading={isLoading} initialValue={initialQuery} />
      </div>

      {/* Error State */}
      {error && !shipment && (
        <div class="bg-red-50 border-2 border-red-200 rounded-xl p-6 flex gap-4">
          <svg class="w-6 h-6 text-red-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4v.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <div>
            <h3 class="font-bold text-red-900 mb-1">Envío no encontrado</h3>
            <p class="text-red-700 text-sm">{error}</p>
            <button
              onClick={handleReset}
              class="mt-3 text-sm font-semibold text-red-600 hover:text-red-700 underline"
            >
              Intentar de nuevo
            </button>
          </div>
        </div>
      )}

      {/* Results Section */}
      {shipment && (
        <div class="animate-fade-in space-y-8">
          {/* Progress Bar */}
          <TrackingProgressBar
            status={shipment.status}
            clientName={shipment.client_name}
            trackingNumber={shipment.tracking_number}
            carrier={shipment.carrier}
            hasGuide={!!(shipment.guide_url || shipment.tracking_number)}
          />

          {/* Details */}
          <div class="bg-white rounded-2xl shadow-lg p-8">
            <TrackingDetails shipment={shipment} />
          </div>

          {/* Timeline */}
          <div class="bg-white rounded-2xl shadow-lg p-8">
            <TrackingTimeline history={history} isLoading={false} />
          </div>

          {/* New Search Button */}
          <div class="flex justify-center">
            <button
              onClick={handleReset}
              class="px-6 py-3 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-semibold transition-colors"
            >
              Rastrear Otro Envío
            </button>
          </div>
        </div>
      )}

      {/* Initial State */}
      {!shipment && !error && !isLoading && (
        <div class="bg-white rounded-2xl shadow-lg p-12 text-center">
          <div class="flex justify-center mb-4">
            <div class="w-16 h-16 rounded-full bg-blue-100 flex items-center justify-center">
              <svg class="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9-4v4m0 0v4m0-4h4m-4 0H9"></path>
              </svg>
            </div>
          </div>
          <h3 class="text-xl font-bold text-gray-900 mb-2">Comienza a rastrear</h3>
          <p class="text-gray-600">
            Busca el número de tracking o de orden en la barra anterior para ver el estado de tu envío
          </p>
        </div>
      )}

      <style>{`
        @keyframes fade-in {
          from {
            opacity: 0;
            transform: translateY(10px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }

        .animate-fade-in {
          animation: fade-in 0.5s ease-out;
        }
      `}</style>
    </div>
  );
}
