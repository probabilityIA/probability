/** @jsxImportSource react */
import { useEffect, useRef, useState } from 'react';
import TrackingSearchInput from './tracking/TrackingSearchInput';
import TrackingProgressBar from './tracking/TrackingProgressBar';
import TrackingDetails from './tracking/TrackingDetails';
import TrackingTimeline from './tracking/TrackingTimeline';
import type { TrackingSearchResult, TrackingHistory, OrderPublicTracking } from '../types/tracking';
import OrderOnlyView from './tracking/OrderOnlyView';
import { getApiUrl } from '../config/api';

interface SearchResult {
  success: boolean;
  message: string;
  data?: {
    shipment?: TrackingSearchResult;
    history?: TrackingHistory[];
    order?: OrderPublicTracking;
  };
}

export default function TrackingClient() {
  const [shipment, setShipment] = useState<TrackingSearchResult | null>(null);
  const [orderOnly, setOrderOnly] = useState<OrderPublicTracking | null>(null);
  const [history, setHistory] = useState<TrackingHistory[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [initialQuery, setInitialQuery] = useState<string>('');
  const autoTriggered = useRef(false);

  const businessIdFromURL = useRef<string>('');

  useEffect(() => {
    if (autoTriggered.current) return;
    if (typeof window === 'undefined') return;
    const params = new URLSearchParams(window.location.search);
    const tracking = params.get('tracking') || params.get('q') || params.get('order');
    const bid = params.get('b') || params.get('business_id');
    if (bid) businessIdFromURL.current = bid;
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
    setOrderOnly(null);
    setHistory([]);

    try {
      const apiUrl = getApiUrl();
      const isOrderNumber = /^prob-\d+$/i.test(query.trim());
      const param = isOrderNumber ? 'order_number' : 'tracking_number';
      const bidParam = businessIdFromURL.current ? `&business_id=${encodeURIComponent(businessIdFromURL.current)}` : '';
      const response = await fetch(
        `${apiUrl}/tracking/search?${param}=${encodeURIComponent(query)}${bidParam}`
      );

      if (!response.ok && response.status !== 404) {
        throw new Error(`Error: ${response.status}`);
      }

      const result: SearchResult = await response.json();

      if (!result.success) {
        setError(result.message || 'No se encontró información del envío');
        return;
      }

      const foundShipment = result.data?.shipment;
      const foundOrder = result.data?.order;

      if (!foundShipment && foundOrder) {
        setOrderOnly(foundOrder);
        return;
      }

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
    setOrderOnly(null);
    setHistory([]);
    setError(null);
  };

  return (
    <div class="space-y-8">
      {/* Search Section */}
      <div class="bg-white rounded-3xl border border-[#EFEAFB] shadow-[0_20px_50px_-30px_rgba(124,58,237,0.3)] p-6 sm:p-8">
        <TrackingSearchInput onSearch={handleSearch} isLoading={isLoading} initialValue={initialQuery} />
      </div>

      {/* Order-only View (no shipment yet) */}
      {orderOnly && !shipment && (
        <div class="animate-fade-in">
          <OrderOnlyView order={orderOnly} />
          <div class="flex justify-center mt-6">
            <button
              onClick={handleReset}
              class="px-6 py-3 rounded-xl bg-gradient-to-r from-[#A855F7] to-[#7C3AED] hover:shadow-[0_14px_28px_-10px_rgba(124,58,237,0.55)] text-white font-semibold transition-all shadow-[0_10px_20px_-10px_rgba(124,58,237,0.5)]"
            >
              Rastrear otro pedido
            </button>
          </div>
        </div>
      )}

      {/* Error State */}
      {error && !shipment && !orderOnly && (
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
            guideUrl={shipment.guide_url}
            hasGuide={!!(shipment.guide_url || shipment.tracking_number)}
          />

          {/* Details */}
          {(shipment.destination_address || shipment.tracking_url) && (
            <div class="bg-white rounded-3xl border border-[#EFEAFB] shadow-[0_20px_50px_-30px_rgba(124,58,237,0.3)] p-6 sm:p-8">
              <TrackingDetails shipment={shipment} />
            </div>
          )}

          {/* Timeline */}
          <div class="bg-white rounded-3xl border border-[#EFEAFB] shadow-[0_20px_50px_-30px_rgba(124,58,237,0.3)] pt-6 sm:pt-8 px-6 sm:px-8">
            <TrackingTimeline history={history} isLoading={false} />
          </div>

          {/* New Search Button */}
          <div class="flex justify-center">
            <button
              onClick={handleReset}
              class="px-6 py-3 rounded-xl bg-gradient-to-r from-[#A855F7] to-[#7C3AED] hover:shadow-[0_14px_28px_-10px_rgba(124,58,237,0.55)] text-white font-semibold transition-all shadow-[0_10px_20px_-10px_rgba(124,58,237,0.5)]"
            >
              Rastrear otro envío
            </button>
          </div>
        </div>
      )}

      {/* Initial State */}
      {!shipment && !orderOnly && !error && !isLoading && (
        <div class="bg-white rounded-3xl border border-[#EFEAFB] shadow-[0_20px_50px_-30px_rgba(124,58,237,0.3)] p-12 text-center">
          <div class="flex justify-center mb-4">
            <div class="w-16 h-16 rounded-full bg-gradient-to-br from-[#F0EAFD] to-[#F6F3FD] flex items-center justify-center">
              <svg class="w-8 h-8 text-[#7C3AED]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9-4v4m0 0v4m0-4h4m-4 0H9"></path>
              </svg>
            </div>
          </div>
          <h3 class="font-space-grotesk text-xl font-bold text-[#181225] mb-2">Comienza a rastrear</h3>
          <p class="text-[#6B6480]">
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
