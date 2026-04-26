/** @jsxImportSource react */
import type { TrackingSearchResult } from '../../types/tracking';

interface TrackingDetailsProps {
  shipment: TrackingSearchResult;
}

const CARRIER_ICONS: Record<string, string> = {
  envioclik: '🚚',
  enviame: '📦',
  mipaquete: '🎁',
  interapidisimo: '⚡',
  default: '🚚',
};

export default function TrackingDetails({ shipment }: TrackingDetailsProps) {
  const carrierName = shipment.carrier?.toLowerCase() || '';
  const carrierIcon = CARRIER_ICONS[carrierName] || CARRIER_ICONS.default;

  return (
    <div class="space-y-6 mt-8">
      {/* Recipient Info */}
      <div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-2xl p-6 border border-slate-200 shadow-sm">
        <div class="flex items-start gap-4">
          <div class="w-14 h-14 bg-gradient-to-br from-blue-400 to-blue-600 rounded-full flex items-center justify-center text-2xl shadow-lg flex-shrink-0">
            {shipment.client_name ? shipment.client_name[0].toUpperCase() : '👤'}
          </div>
          <div class="flex-1">
            <h3 class="text-lg font-bold text-gray-900">
              {shipment.client_name || 'Cliente'}
            </h3>
            {shipment.destination_address && (
              <p class="text-sm text-gray-600 mt-2 flex items-start gap-2">
                <span class="text-base mt-0.5">📍</span>
                <span class="flex-1">{shipment.destination_address}</span>
              </p>
            )}
          </div>
        </div>
      </div>

      {/* Tracking Info Cards */}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Tracking Number */}
        <div class="bg-white rounded-xl p-5 border border-gray-200 shadow-sm hover:shadow-md transition-shadow">
          <div class="flex items-center gap-2 mb-2">
            <span class="text-xl">🏷️</span>
            <p class="text-xs font-bold text-gray-500 uppercase tracking-widest">Tracking</p>
          </div>
          <p class="text-base font-mono font-bold text-gray-900 break-all">
            {shipment.tracking_number}
          </p>
          <p class="text-xs text-gray-400 mt-2">Número único de rastreo</p>
        </div>

        {/* Carrier */}
        <div class="bg-white rounded-xl p-5 border border-gray-200 shadow-sm hover:shadow-md transition-shadow">
          <div class="flex items-center gap-2 mb-2">
            <span class="text-xl">{carrierIcon}</span>
            <p class="text-xs font-bold text-gray-500 uppercase tracking-widest">Transportista</p>
          </div>
          <p class="text-base font-semibold text-gray-900">
            {shipment.carrier || 'Sin asignar'}
          </p>
          <p class="text-xs text-gray-400 mt-2">Empresa de logística</p>
        </div>
      </div>

      {/* Action Buttons */}
      {(shipment.guide_url || shipment.tracking_url) && (
        <div class="pt-2 flex flex-col sm:flex-row gap-3">
          {shipment.guide_url && (
            <a
              href={shipment.guide_url}
              target="_blank"
              rel="noopener noreferrer"
              class="flex-1 flex items-center justify-center gap-2 px-6 py-3 rounded-xl bg-gradient-to-r from-emerald-600 to-emerald-700 hover:from-emerald-700 hover:to-emerald-800 text-white font-semibold transition-all shadow-lg hover:shadow-xl active:scale-95"
            >
              <span>📄</span>
              Ver Guía (PDF)
              <span>→</span>
            </a>
          )}
          {shipment.tracking_url && (
            <a
              href={shipment.tracking_url}
              target="_blank"
              rel="noopener noreferrer"
              class="flex-1 flex items-center justify-center gap-2 px-6 py-3 rounded-xl bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-semibold transition-all shadow-lg hover:shadow-xl active:scale-95"
            >
              <span>🔗</span>
              Ver en Sitio del Transportista
              <span>→</span>
            </a>
          )}
        </div>
      )}
    </div>
  );
}
