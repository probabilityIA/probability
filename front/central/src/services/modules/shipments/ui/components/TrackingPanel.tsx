'use client';

import { useEffect, useState } from 'react';
import {
  X, AlertTriangle, RefreshCw, FileText, Navigation, MapPin,
  Package, Truck, Clock, Home, CheckCircle2, XCircle, Calendar
} from 'lucide-react';
import { Shipment, EnvioClickTrackHistory } from '../../domain/types';
import { trackShipmentAction } from '../../infra/actions';
import { ColombiaMap } from './ColombiaMap';

interface TrackingPanelProps {
  shipment: Shipment;
  onClose: () => void;
  onCancel: (id: string) => void;
  cancelingId: string | null;
}

// Status badge styles
const STATUS_CONFIG: Record<string, { label: string; color: string; bgColor: string }> = {
  pending: { label: 'Pendiente', color: 'bg-amber-100 text-amber-700 border-amber-200', bgColor: '#fbbf24' },
  in_transit: { label: 'En Tránsito', color: 'bg-blue-100 text-blue-700 border-blue-200', bgColor: '#3b82f6' },
  delivered: { label: 'Entregado', color: 'bg-emerald-100 text-emerald-700 border-emerald-200', bgColor: '#10b981' },
  failed: { label: 'Fallido', color: 'bg-red-100 text-red-700 border-red-200', bgColor: '#ef4444' },
};

// Map status to progress step (1-5)
const STATUS_TO_STEP: Record<string, number> = {
  pending: 1,
  picked_up: 2,
  in_transit: 3,
  out_for_delivery: 4,
  delivered: 5,
};

// Progress step icons
const STEP_ICONS = [
  { icon: Package, label: 'Creado' },
  { icon: Truck, label: 'Recogido' },
  { icon: MapPin, label: 'En Tránsito' },
  { icon: Home, label: 'En Reparto' },
  { icon: CheckCircle2, label: 'Entregado' },
];

function formatDate(dateStr?: string) {
  if (!dateStr) return null;
  return new Date(dateStr).toLocaleDateString('es-CO', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  });
}

function formatCurrency(value?: number): string {
  if (!value) return '-';
  return new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', minimumFractionDigits: 0 }).format(value);
}

function extractDepartment(address?: string): string | null {
  if (!address) return null;
  const parts = address.split(',').map(s => s.trim()).filter(Boolean);
  if (parts.length >= 2) {
    const last = parts[parts.length - 1].toLowerCase();
    if (last === 'colombia') return parts[parts.length - 2];
    return parts[parts.length - 1];
  }
  return null;
}

export function TrackingPanel({ shipment, onClose, onCancel, cancelingId }: TrackingPanelProps) {
  const [tracking, setTracking] = useState<{ loading: boolean; data?: any; error?: string }>({ loading: false });
  const [trackingUpdatedAt, setTrackingUpdatedAt] = useState<Date | null>(null);

  // Load tracking history on mount
  useEffect(() => {
    if (shipment.tracking_number) {
      setTracking({ loading: true });
      trackShipmentAction(shipment.tracking_number)
        .then(res => {
          if ('data' in res && res.success) {
            setTracking({ loading: false, data: res.data });
            setTrackingUpdatedAt(new Date());
          } else {
            setTracking({ loading: false, error: (res as any).message || 'No disponible' });
          }
        })
        .catch(err => {
          setTracking({ loading: false, error: err.message });
        });
    }
  }, [shipment.id, shipment.tracking_number]);

  const cancelId = shipment.tracking_number || shipment.id.toString();
  const statusConfig = STATUS_CONFIG[shipment.status] || STATUS_CONFIG.pending;
  const currentStep = STATUS_TO_STEP[shipment.status] || 1;
  const destination = shipment.destination_address || '';
  const departmentName = extractDepartment(destination);
  const progressPercentage = ((currentStep - 1) / 4) * 100;

  // Time since last tracking update
  const lastUpdateText = trackingUpdatedAt
    ? `Hace ${Math.floor((new Date().getTime() - trackingUpdatedAt.getTime()) / 60000)}m`
    : 'Desconocido';

  return (
    <div className="flex flex-col h-full bg-white">
      {/* ─── HEADER CON GRADIENTE ─── */}
      <div
        className="text-white p-5 flex justify-between items-start gap-3"
        style={{
          background: 'linear-gradient(135deg, #1e1b4b 0%, #312e81 50%, #1e40af 100%)',
        }}
      >
        <div className="flex-1 min-w-0">
          <h3 className="text-lg font-bold truncate">{shipment.client_name || 'Cliente desconocido'}</h3>
          {destination && (
            <div className="flex items-center gap-1.5 mt-1 text-sm text-blue-100">
              <MapPin size={14} className="flex-shrink-0" />
              <span className="truncate">{destination}</span>
            </div>
          )}
          <div className="mt-2 flex items-center gap-2 flex-wrap">
            <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-semibold border ${statusConfig.color}`}>
              {statusConfig.label}
            </span>
            {shipment.is_test && (
              <span className="text-[10px] font-bold bg-orange-400 text-white px-2 py-1 rounded uppercase tracking-widest">
                TEST
              </span>
            )}
          </div>
          {shipment.tracking_number && (
            <p className="mt-2 text-[11px] font-mono text-blue-100 break-all">{shipment.tracking_number}</p>
          )}
        </div>
        <button
          onClick={onClose}
          className="p-1.5 rounded-full hover:bg-blue-600 text-blue-100 hover:text-white transition-colors flex-shrink-0"
        >
          <X size={18} />
        </button>
      </div>

      {/* ─── SCROLLABLE CONTENT ─── */}
      <div className="flex-1 overflow-y-auto">
        {/* ─── PROGRESS BAR (5 STEPS) ─── */}
        <div className="px-5 py-6 border-b border-gray-100 bg-gray-50/30">
          <p className="text-xs font-bold text-gray-500 uppercase tracking-wider mb-4">Progreso del envío</p>

          <div className="flex items-center justify-between gap-2">
            {STEP_ICONS.map((step, idx) => {
              const Icon = step.icon;
              const isActive = idx + 1 === currentStep;
              const isCompleted = idx + 1 < currentStep;
              const isFailed = shipment.status === 'failed';

              return (
                <div key={idx} className="flex-1 flex flex-col items-center gap-1">
                  <div
                    className={`w-8 h-8 rounded-full flex items-center justify-center transition-all duration-300 ${
                      isFailed
                        ? 'bg-red-500 text-white'
                        : isActive
                          ? 'bg-blue-500 text-white scale-110 animate-pulse'
                          : isCompleted
                            ? 'bg-emerald-500 text-white'
                            : 'bg-gray-300 text-gray-600'
                    }`}
                  >
                    <Icon size={16} />
                  </div>
                  <span className="text-[10px] font-semibold text-gray-600 text-center">{step.label}</span>
                </div>
              );
            })}
          </div>

          {/* Progress line */}
          <div className="mt-4 h-1 bg-gray-300 rounded-full overflow-hidden">
            <div
              className={`h-full transition-all duration-700 ${
                shipment.status === 'failed'
                  ? 'bg-red-500'
                  : 'bg-gradient-to-r from-blue-400 to-emerald-400'
              }`}
              style={{ width: `${shipment.status === 'failed' ? 100 : progressPercentage}%` }}
            />
          </div>
        </div>

        {/* ─── STATS CARDS (2x2 GRID) ─── */}
        <div className="grid grid-cols-2 gap-3 p-4 bg-white">
          {/* Estado */}
          <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-3 border border-blue-200">
            <p className="text-[10px] font-bold text-blue-600 uppercase tracking-widest mb-1">Estado</p>
            <p className={`text-xs font-semibold ${statusConfig.color.split(' ')[1]}`}>{statusConfig.label}</p>
            <p className="text-[10px] text-blue-600 mt-0.5">{lastUpdateText}</p>
          </div>

          {/* Transportista */}
          <div className="bg-gradient-to-br from-purple-50 to-purple-100 rounded-lg p-3 border border-purple-200">
            <p className="text-[10px] font-bold text-purple-600 uppercase tracking-widest mb-1">Transportista</p>
            <div className="flex items-center gap-1.5">
              <Truck size={14} className="text-purple-600" />
              <p className="text-xs font-semibold text-purple-900">{shipment.carrier || '-'}</p>
            </div>
          </div>

          {/* Tiempo estimado */}
          <div className="bg-gradient-to-br from-amber-50 to-amber-100 rounded-lg p-3 border border-amber-200">
            <p className="text-[10px] font-bold text-amber-600 uppercase tracking-widest mb-1">Entrega</p>
            <div className="flex items-center gap-1.5">
              <Calendar size={14} className="text-amber-600" />
              <p className="text-xs font-semibold text-amber-900">
                {shipment.estimated_delivery ? formatDate(shipment.estimated_delivery) : '-'}
              </p>
            </div>
          </div>

          {/* Costo total */}
          <div className="bg-gradient-to-br from-emerald-50 to-emerald-100 rounded-lg p-3 border border-emerald-200">
            <p className="text-[10px] font-bold text-emerald-600 uppercase tracking-widest mb-1">Costo</p>
            <p className="text-xs font-mono font-bold text-emerald-900">{formatCurrency(shipment.total_cost || shipment.shipping_cost)}</p>
          </div>
        </div>

        {/* ─── MAPA COLOMBIA ─── */}
        <div className="px-4 py-5 border-t border-gray-100">
          <p className="text-xs font-bold text-gray-500 uppercase tracking-wider mb-3">Departamento de destino</p>
          <div className="rounded-xl overflow-hidden border border-gray-200 bg-gray-50 p-3">
            <ColombiaMap highlightedDepartment={departmentName} status={shipment.status} />
          </div>
        </div>

        {/* ─── TIMELINE ─── */}
        <div className="px-4 py-5 border-t border-gray-100">
          <p className="text-xs font-bold text-gray-500 uppercase tracking-wider mb-3">Historial de rastreo</p>

          {tracking.loading ? (
            <div className="flex items-center gap-2 text-sm text-gray-400 py-6">
              <RefreshCw size={16} className="animate-spin" />
              <span>Consultando rastreo...</span>
            </div>
          ) : tracking.error ? (
            <div className="flex items-start gap-2 bg-amber-50 border border-amber-200 rounded-lg p-3 text-xs text-amber-700">
              <AlertTriangle size={14} className="flex-shrink-0 mt-0.5" />
              <span>{tracking.error}</span>
            </div>
          ) : tracking.data?.history?.length > 0 ? (
            <div className="space-y-0">
              {/* Vertical line gradient */}
              <div className="relative pl-7">
                {/* Gradient line */}
                <div className="absolute left-2.5 top-0 bottom-0 w-1 bg-gradient-to-b from-blue-400 via-blue-300 to-gray-200 rounded-full" />

                {/* Events */}
                {tracking.data.history.map((event: EnvioClickTrackHistory, idx: number) => {
                  const isFirst = idx === 0;
                  return (
                    <div
                      key={idx}
                      className="relative pb-4 animate-fade-in"
                      style={{ animationDelay: `${idx * 50}ms`, animationFillMode: 'both' }}
                    >
                      {/* Circle dot */}
                      <div
                        className={`absolute -left-5 top-1 w-4 h-4 rounded-full ring-2 ring-white transition-all duration-300 ${
                          isFirst ? 'bg-blue-500 scale-110 shadow-md shadow-blue-500/50' : 'bg-gray-300'
                        }`}
                      />

                      {/* Content */}
                      <div className="pt-0.5">
                        <div className="flex items-baseline justify-between gap-2">
                          <p className={`text-sm font-bold ${isFirst ? 'text-blue-700' : 'text-gray-800'}`}>
                            {event.status}
                          </p>
                          <p className="text-[10px] text-gray-400 flex-shrink-0">{event.date}</p>
                        </div>
                        {event.description && <p className="text-xs text-gray-600 mt-0.5">{event.description}</p>}
                        {event.location && (
                          <div className="flex items-center gap-1 mt-1.5 text-[11px] text-gray-500">
                            <MapPin size={11} className="flex-shrink-0" />
                            <span>{event.location}</span>
                          </div>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          ) : !shipment.tracking_number ? (
            <p className="text-xs text-gray-400 italic py-4">Este envío no tiene número de tracking.</p>
          ) : (
            <p className="text-xs text-gray-400 italic py-4">No hay historial disponible.</p>
          )}
        </div>

        {/* Space for buttons */}
        <div className="h-24" />
      </div>

      {/* ─── ACTION BUTTONS (STICKY) ─── */}
      <div className="border-t border-gray-100 bg-white p-4 flex gap-2 sticky bottom-0">
        {shipment.guide_url && (
          <a
            href={shipment.guide_url}
            target="_blank"
            rel="noopener noreferrer"
            className="flex-1 flex items-center justify-center gap-2 py-2.5 px-3 rounded-lg bg-blue-600 hover:bg-blue-700 text-white text-xs font-semibold transition-colors shadow-sm hover:shadow-md"
          >
            <FileText size={14} />
            Ver Guía
          </a>
        )}
        {shipment.tracking_url && (
          <a
            href={shipment.tracking_url}
            target="_blank"
            rel="noopener noreferrer"
            className="flex-1 flex items-center justify-center gap-2 py-2.5 px-3 rounded-lg bg-indigo-50 hover:bg-indigo-100 text-indigo-700 border border-indigo-200 text-xs font-semibold transition-colors"
          >
            <Navigation size={14} />
            Rastrear
          </a>
        )}
        <button
          onClick={() => onCancel(cancelId)}
          disabled={cancelingId === cancelId}
          className="flex items-center justify-center gap-1.5 py-2.5 px-3 rounded-lg bg-red-50 hover:bg-red-100 text-red-600 border border-red-200 text-xs font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          title="Cancelar envío"
        >
          {cancelingId === cancelId ? <RefreshCw size={14} className="animate-spin" /> : <XCircle size={14} />}
          {!cancelingId && 'Cancelar'}
        </button>
      </div>
    </div>
  );
}
