/** @jsxImportSource react */
import type { TrackingSearchResult } from '../../types/tracking';

interface TrackingDetailsProps {
  shipment: TrackingSearchResult;
}

function IconMapPin() {
  return (
    <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
    </svg>
  );
}

function IconExternalLink() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 3L21 3m0 0h-5.25M21 3v5.25" />
    </svg>
  );
}

export default function TrackingDetails({ shipment }: TrackingDetailsProps) {
  if (!shipment.destination_address && !shipment.tracking_url) return null;

  return (
    <div class="space-y-5">
      {shipment.destination_address && (
        <div class="flex items-start gap-3 bg-[#F6F3FD] border border-[#E7E2F5] rounded-2xl p-5">
          <div class="w-9 h-9 rounded-xl bg-white flex items-center justify-center text-[#7C3AED] flex-shrink-0 shadow-sm">
            <IconMapPin />
          </div>
          <div class="flex-1">
            <p class="text-[11px] font-bold text-[#7C3AED] uppercase tracking-widest mb-1">Dirección de entrega</p>
            <p class="text-sm text-[#3A2E5C]">{shipment.destination_address}</p>
          </div>
        </div>
      )}

      {shipment.tracking_url && (
        <a
          href={shipment.tracking_url}
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl bg-white border-2 border-[#E7E2F5] hover:border-[#C4B5FD] text-[#7C3AED] font-semibold transition-all active:scale-[0.98]"
        >
          <IconExternalLink />
          Ver en sitio del transportista
        </a>
      )}
    </div>
  );
}
