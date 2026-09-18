/** @jsxImportSource react */
import type { TrackingStatus } from '../../types/tracking';

interface TrackingProgressBarProps {
  status: TrackingStatus;
  clientName?: string;
  trackingNumber?: string;
  carrier?: string;
  guideUrl?: string;
  hasGuide?: boolean;
}

function IconDocument() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M6 2a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8.828a2 2 0 00-.586-1.414l-4.828-4.828A2 2 0 0012.172 2H6z" />
    </svg>
  );
}

function IconBox() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
    </svg>
  );
}

function IconCheck() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
      <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
    </svg>
  );
}

function IconTruck() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M3 16V6a1 1 0 011-1h9v11M3 16h10m0 0h4m-4 0V9h4l3 3.5V16m0 0h-3m-9 0a2 2 0 11-4 0 2 2 0 014 0zm10 0a2 2 0 11-4 0 2 2 0 014 0z" />
    </svg>
  );
}

function IconMapPin() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
    </svg>
  );
}

function IconFlag() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M5 21V4m0 0h11l-1.5 4L16 12H5" />
    </svg>
  );
}

function IconAlert() {
  return (
    <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
    </svg>
  );
}

function IconParty() {
  return (
    <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  );
}

function IconTag() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.83.699 2.528 0l4.318-4.318a1.788 1.788 0 000-2.528l-9.581-9.581A2.25 2.25 0 009.568 3z" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M6 6h.008v.008H6V6z" />
    </svg>
  );
}

function IconArrow() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
      <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
    </svg>
  );
}

const CARRIER_LOGOS: Record<string, string> = {
  coordinadora: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_coordinadora.png',
  servientrega: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_servientrega.png',
  interrapidisimo: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_inerapidisimo.png',
  envia: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_envia.png',
  deprisa: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_deprisa.png',
  pibox: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_pibox.png',
  'mensajeros urbanos': 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_mensajerosUrbanos.png',
  tcc: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_TCC.png',
};

function getCarrierLogo(carrier?: string): string | null {
  if (!carrier) return null;
  return CARRIER_LOGOS[carrier.trim().toLowerCase()] ?? null;
}

const STATUS_LABEL: Record<TrackingStatus, string> = {
  pending: 'Pendiente',
  picked_up: 'Recogido',
  in_transit: 'En Tránsito',
  out_for_delivery: 'En Reparto',
  delivered: 'Entregado',
  failed: 'Fallido',
};

const STATUS_TO_STEP: Record<TrackingStatus, number> = {
  pending: 0,
  picked_up: 1,
  in_transit: 2,
  out_for_delivery: 3,
  delivered: 4,
  failed: -1,
};

const STEPS_WITH_GUIDE = [
  { label: 'Guía Generada', Icon: IconDocument },
  { label: 'Recogido', Icon: IconCheck },
  { label: 'En Tránsito', Icon: IconTruck },
  { label: 'En Reparto', Icon: IconMapPin },
  { label: 'Entregado', Icon: IconFlag },
];

const STEPS_NO_GUIDE = [
  { label: 'Creado', Icon: IconBox },
  { label: 'Recogido', Icon: IconCheck },
  { label: 'En Tránsito', Icon: IconTruck },
  { label: 'En Reparto', Icon: IconMapPin },
  { label: 'Entregado', Icon: IconFlag },
];

export default function TrackingProgressBar({
  status,
  clientName,
  trackingNumber,
  carrier,
  guideUrl,
  hasGuide,
}: TrackingProgressBarProps) {
  const STEPS = hasGuide ? STEPS_WITH_GUIDE : STEPS_NO_GUIDE;
  const currentStep = STATUS_TO_STEP[status];
  const isFailed = status === 'failed';

  const getStatusMessage = () => {
    switch (status) {
      case 'delivered':
        return { Icon: IconParty, title: '¡Envío entregado!', description: 'Tu paquete ha llegado exitosamente' };
      case 'failed':
        return { Icon: IconAlert, title: 'No se pudo entregar', description: 'Contacta al transportista para más información' };
      case 'pending':
        return hasGuide
          ? { Icon: IconDocument, title: '¡Guía generada!', description: 'Tu envío fue registrado con la transportadora. Pronto será recogido' }
          : { Icon: IconBox, title: 'Pedido registrado', description: 'Tu pedido fue creado y está pendiente de despacho' };
      default:
        return { Icon: IconMapPin, title: 'En camino', description: 'Tu envío está siendo entregado. Pronto recibirás más actualizaciones' };
    }
  };

  const statusMsg = getStatusMessage();
  const carrierLogo = getCarrierLogo(carrier);

  return (
    <div class="bg-white rounded-3xl border border-[#EFEAFB] shadow-[0_20px_50px_-30px_rgba(124,58,237,0.3)] overflow-hidden">
      <div class="flex items-start justify-between gap-4 p-6 sm:p-8 pb-0">
        <div>
          <h2 class="font-space-grotesk font-bold text-xl sm:text-2xl text-[#181225]">Estado del envío</h2>
          <p class="text-sm text-[#8B85A0] mt-1">Rastreo en tiempo real</p>
        </div>
        <span class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-bold text-white bg-gradient-to-r from-[#A855F7] to-[#7C3AED] shadow-[0_8px_18px_-8px_rgba(124,58,237,0.6)] whitespace-nowrap">
          {STATUS_LABEL[status]}
        </span>
      </div>

      <div class="px-4 sm:px-8 pt-8 pb-2">
        <div class="flex items-start">
          {STEPS.map((step, idx) => {
            const isCompleted = idx < currentStep;
            const isActive = idx === currentStep && !isFailed;
            const StepIcon = step.Icon;
            return (
              <div key={idx} class="contents">
                <div class="flex flex-col items-center gap-2 w-10 sm:w-14 shrink-0">
                  <div
                    class={`w-10 h-10 rounded-full flex items-center justify-center transition-all duration-300 ${
                      isCompleted
                        ? 'bg-gradient-to-br from-[#A855F7] to-[#7C3AED] text-white shadow-[0_8px_20px_-8px_rgba(124,58,237,0.6)]'
                        : isActive
                          ? 'bg-white border-2 border-[#7C3AED] text-[#7C3AED] ring-4 ring-[#F0EAFD]'
                          : 'bg-[#F6F3FD] border border-[#E7E2F5] text-[#C7BEDB]'
                    }`}
                  >
                    {isCompleted ? <IconCheck /> : <StepIcon />}
                  </div>
                  <span
                    class={`text-[10px] sm:text-xs font-semibold text-center leading-tight ${
                      isActive ? 'text-[#7C3AED]' : isCompleted ? 'text-[#3A2E5C]' : 'text-[#B9AED0]'
                    }`}
                  >
                    {step.label}
                  </span>
                </div>
                {idx < STEPS.length - 1 && (
                  <div
                    class={`flex-1 h-[3px] rounded-full mt-5 transition-colors duration-700 ${
                      isCompleted ? 'bg-gradient-to-r from-[#A855F7] to-[#7C3AED]' : 'bg-[#EFEAFB]'
                    }`}
                  />
                )}
              </div>
            );
          })}
        </div>
      </div>

      <div class="px-6 sm:px-8 pt-6">
        <div class="flex items-start gap-3 bg-[#F6F3FD] border border-[#E7E2F5] rounded-2xl p-5">
          <div class="w-9 h-9 rounded-xl bg-white flex items-center justify-center text-[#7C3AED] flex-shrink-0 shadow-sm">
            <statusMsg.Icon />
          </div>
          <div>
            <h4 class="font-bold text-[#181225] text-sm">{statusMsg.title}</h4>
            <p class="text-[#6B6480] text-sm mt-0.5">{statusMsg.description}</p>
          </div>
        </div>
      </div>

      {clientName && (
        <div class="px-6 sm:px-8 pt-6">
          <div class="flex items-center gap-3 pb-6 border-b border-[#F1EDFA]">
            <div class="w-10 h-10 rounded-full bg-gradient-to-br from-[#A855F7] to-[#7C3AED] text-white flex items-center justify-center font-bold text-sm flex-shrink-0">
              {clientName[0].toUpperCase()}
            </div>
            <p class="font-semibold text-[#181225] text-sm">{clientName}</p>
          </div>
        </div>
      )}

      {(trackingNumber || carrier || guideUrl) && (
        <div class={`grid grid-cols-1 sm:grid-cols-2 ${guideUrl ? 'lg:grid-cols-3' : ''} gap-4 p-6 sm:p-8 pt-6`}>
          {trackingNumber && (
            <div class="bg-white rounded-2xl p-5 border border-[#EFEAFB] shadow-[0_6px_16px_-12px_rgba(124,58,247,0.25)]">
              <div class="flex items-center gap-2 mb-2 text-[#7C3AED]">
                <IconTag />
                <p class="text-[11px] font-bold text-[#8B85A0] uppercase tracking-widest">Tracking</p>
              </div>
              <p class="text-base font-mono font-bold text-[#181225] break-all">{trackingNumber}</p>
              <p class="text-xs text-[#B0A9C2] mt-1.5">Número único de rastreo</p>
            </div>
          )}

          {carrier && (
            <div class="bg-white rounded-2xl p-5 border border-[#EFEAFB] shadow-[0_6px_16px_-12px_rgba(124,58,247,0.25)]">
              <div class="flex items-center gap-2 mb-2 text-[#7C3AED]">
                <IconTruck />
                <p class="text-[11px] font-bold text-[#8B85A0] uppercase tracking-widest">Transportista</p>
              </div>
              {carrierLogo ? (
                <img src={carrierLogo} alt={carrier} class="h-7 max-w-[120px] object-contain object-left" />
              ) : (
                <p class="text-base font-semibold text-[#181225]">{carrier}</p>
              )}
              <p class="text-xs text-[#B0A9C2] mt-1.5">Empresa de logística</p>
            </div>
          )}

          {guideUrl && (
            <a
              href={guideUrl}
              target="_blank"
              rel="noopener noreferrer"
              class="group flex flex-col justify-between bg-gradient-to-br from-[#A855F7] to-[#7C3AED] rounded-2xl p-5 text-white shadow-[0_14px_30px_-14px_rgba(124,58,237,0.6)] hover:shadow-[0_18px_36px_-12px_rgba(124,58,237,0.65)] transition-all"
            >
              <div class="flex items-center gap-2 mb-2 text-white/80">
                <IconDocument />
                <p class="text-[11px] font-bold uppercase tracking-widest">Guía de envío</p>
              </div>
              <div class="flex items-center justify-between">
                <p class="text-base font-bold">Ver guía (PDF)</p>
                <span class="transition-transform group-hover:translate-x-1">
                  <IconArrow />
                </span>
              </div>
            </a>
          )}
        </div>
      )}
    </div>
  );
}
