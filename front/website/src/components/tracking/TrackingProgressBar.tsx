/** @jsxImportSource react */
import type { TrackingStatus } from '../../types/tracking';

interface TrackingProgressBarProps {
  status: TrackingStatus;
  clientName?: string;
  trackingNumber?: string;
  carrier?: string;
  hasGuide?: boolean;
}

const STATUS_CONFIG: Record<TrackingStatus, { label: string; color: string }> = {
  pending: { label: 'Pendiente', color: '#F59E0B' },
  picked_up: { label: 'Recogido', color: '#10B981' },
  in_transit: { label: 'En Tránsito', color: '#3B82F6' },
  out_for_delivery: { label: 'En Reparto', color: '#8B5CF6' },
  delivered: { label: 'Entregado', color: '#22C55E' },
  failed: { label: 'Fallido', color: '#EF4444' },
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
  { label: 'Guía Generada', icon: '📄' },
  { label: 'Recogido', icon: '✅' },
  { label: 'En Tránsito', icon: '🚚' },
  { label: 'En Reparto', icon: '🚙' },
  { label: 'Entregado', icon: '✅' },
];

const STEPS_NO_GUIDE = [
  { label: 'Creado', icon: '📦' },
  { label: 'Recogido', icon: '✅' },
  { label: 'En Tránsito', icon: '🚚' },
  { label: 'En Reparto', icon: '🚙' },
  { label: 'Entregado', icon: '✅' },
];

export default function TrackingProgressBar({
  status,
  clientName,
  trackingNumber,
  carrier,
  hasGuide,
}: TrackingProgressBarProps) {
  const STEPS = hasGuide ? STEPS_WITH_GUIDE : STEPS_NO_GUIDE;
  const config = STATUS_CONFIG[status];
  const currentStep = STATUS_TO_STEP[status];
  const isFailed = status === 'failed';

  const getStatusMessage = () => {
    switch (status) {
      case 'delivered':
        return {
          icon: '🎉',
          title: '¡Envío entregado!',
          description: 'Tu paquete ha llegado exitosamente',
        };
      case 'failed':
        return {
          icon: '⚠️',
          title: 'No se pudo entregar',
          description: 'Contacta al transportista para más información',
        };
      case 'pending':
        return hasGuide
          ? {
              icon: '📄',
              title: '¡Guía generada!',
              description: 'Tu envío fue registrado con la transportadora. Pronto será recogido',
            }
          : {
              icon: '📦',
              title: 'Pedido registrado',
              description: 'Tu pedido fue creado y está pendiente de despacho',
            };
      default:
        return {
          icon: '📍',
          title: 'En camino',
          description: 'Tu envío está siendo entregado. Pronto recibirás más actualizaciones',
        };
    }
  };

  const statusMsg = getStatusMessage();

  return (
    <div class="tracking-container">
      {/* Header */}
      <div class="tracking-header">
        <div>
          <h2 class="tracking-title">Estado del Envío</h2>
          <p class="tracking-subtitle">Rastreo en tiempo real</p>
        </div>
        <div class="tracking-badge" style={{ backgroundColor: config.color }}>
          {config.label}
        </div>
      </div>

      {/* Progress Section */}
      <div class="tracking-progress-section">
        {/* Progress Track */}
        <div class="tracking-progress-track">
          {/* Lines */}
          {[...Array(4)].map((_, i) => (
            <div
              key={`line-${i}`}
              class={`tracking-line ${i < currentStep ? 'completed' : ''}`}
              style={{
                left: `${40 + i * 190}px`,
                width: '190px',
              }}
            />
          ))}

          {/* Nodes */}
          {STEPS.map((step, idx) => {
            const isCompleted = idx < currentStep;
            const isActive = idx === currentStep;

            return (
              <div
                key={`node-${idx}`}
                class={`tracking-node ${isActive ? 'active' : isCompleted ? 'completed' : 'pending'}`}
                style={{ left: `${20 + idx * 190}px` }}
              >
                {step.icon}
              </div>
            );
          })}

          {/* Labels */}
          {STEPS.map((step, idx) => {
            const isCompleted = idx < currentStep;
            const isActive = idx === currentStep;

            return (
              <div
                key={`label-${idx}`}
                class={`tracking-label ${isActive ? 'active' : isCompleted ? 'completed' : ''}`}
                style={{ left: `${45 + idx * 190}px` }}
              >
                {step.label}
              </div>
            );
          })}
        </div>
      </div>

      {/* Status Card */}
      <div class="tracking-status-card">
        <div class="tracking-status-icon">{statusMsg.icon}</div>
        <div class="tracking-status-text">
          <h4>{statusMsg.title}</h4>
          <p>{statusMsg.description}</p>
        </div>
      </div>

      <div class="tracking-divider" />

      {/* Client Row */}
      {clientName && (
        <>
          <div class="tracking-client-row">
            <div class="tracking-avatar">
              {clientName[0].toUpperCase()}
            </div>
            <h3 class="tracking-client-name">{clientName}</h3>
          </div>
        </>
      )}

      {/* Bottom Grid */}
      {(trackingNumber || carrier) && (
        <div class="tracking-bottom-grid">
          {trackingNumber && (
            <div class="tracking-info-card">
              <div class="tracking-info-header">
                <span class="tracking-info-icon">🏷️</span>
                <span class="tracking-info-label">Tracking</span>
              </div>
              <div class="tracking-info-value">{trackingNumber}</div>
              <div class="tracking-info-subtext">Número único de rastreo</div>
            </div>
          )}

          {carrier && (
            <div class="tracking-info-card">
              <div class="tracking-info-header">
                <span class="tracking-info-icon">🚚</span>
                <span class="tracking-info-label">Transportista</span>
              </div>
              <div class="tracking-info-value">{carrier}</div>
              <div class="tracking-info-subtext">Empresa de logística</div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
