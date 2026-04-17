'use client';

import { useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import { useRouter, useSearchParams } from 'next/navigation';
import { Badge, Button, Modal } from '@/shared/ui';
import { useHasPermission } from '@/shared/contexts/permissions-context';
import { getShipmentsAction, trackShipmentAction, cancelShipmentAction, cancelBatchShipmentAction, syncShipmentStatusAction } from '../../infra/actions';
import { GetShipmentsParams, Shipment, EnvioClickTrackHistory } from '../../domain/types';
import { useShipmentSSE } from '../hooks/useShipmentSSE';
import {
    Search, Package, Truck, Calendar, MapPin, X, RefreshCw,
    AlertTriangle, Plus, ChevronLeft, ChevronRight, FileText,
    Download, CheckCircle2, Clock, XCircle, Navigation,
    DollarSign, Box, User, Building2, Hash, StickyNote
} from 'lucide-react';
import { ManualShipmentModal } from './ManualShipmentModal';
import { SyncProgressModal } from './SyncProgressModal';
import { usePermissions } from '@/shared/contexts/permissions-context';

// Carga dinámica del mapa para evitar SSR issues
const MapComponent = dynamic(() => import('@/shared/ui/MapComponent'), {
    ssr: false,
    loading: () => (
        <div className="flex items-center justify-center h-full bg-gray-50 dark:bg-gray-700 rounded-lg">
            <RefreshCw size={20} className="animate-spin text-gray-400 dark:text-gray-500" />
        </div>
    ),
});

// ─── Helpers ───────────────────────────────────────────────────────────────

const STATUS_CONFIG: Record<string, { label: string; color: string; icon: React.ReactNode; border: string }> = {
    delivered: { label: 'Entregado', color: 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-600', icon: <CheckCircle2 size={12} />, border: 'border-emerald-400' },
    in_transit: { label: 'En tránsito', color: 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 border-blue-200 dark:border-blue-600', icon: <Truck size={12} />, border: 'border-blue-400' },
    pending: { label: 'Pendiente', color: 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300 border-amber-200 dark:border-amber-600', icon: <Clock size={12} />, border: 'border-amber-400' },
    failed: { label: 'Fallido', color: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 border-red-200 dark:border-red-600', icon: <XCircle size={12} />, border: 'border-red-400' },
    cancelled: { label: 'Cancelado', color: 'bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 border-gray-200 dark:border-gray-600', icon: <X size={12} />, border: 'border-gray-400' },
};

const CHIP_STATUS_OPTIONS = [
    { value: 'pending', label: 'Pendiente', icon: Clock, activeClass: 'bg-amber-500 text-white' },
    { value: 'in_transit', label: 'En tránsito', icon: Truck, activeClass: 'bg-blue-500 text-white' },
    { value: 'delivered', label: 'Entregado', icon: CheckCircle2, activeClass: 'bg-emerald-500 text-white' },
    { value: 'failed', label: 'Fallido', icon: XCircle, activeClass: 'bg-red-500 text-white' },
    { value: 'cancelled', label: 'Cancelado', icon: X, activeClass: 'bg-gray-500 text-white' },
];

function StatusBadge({ status }: { status: string }) {
    const cfg = STATUS_CONFIG[status] || { label: status, color: 'bg-gray-100 text-gray-600 dark:text-gray-300 border-gray-200 dark:border-gray-600 dark:border-gray-700', icon: null };
    return (
        <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold border ${cfg.color}`}>
            {cfg.icon}
            {cfg.label}
        </span>
    );
}

function formatDate(dateStr?: string) {
    if (!dateStr) return null;
    return new Date(dateStr).toLocaleDateString('es-CO', { day: 'numeric', month: 'short', year: 'numeric' });
}

function formatMoney(amount?: number) {
    if (amount == null) return '—';
    return new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(amount);
}

// Extract city from destination_address.
// Tries multiple strategies:
// 1. Last segment after a comma (e.g., "Calle 28 # 85-36, Montería")
// 2. Second-to-last segment if last is a country (e.g., "..., Montería, Colombia")
// 3. Falls back to the full address as city hint (Nominatim handles fuzzy matching)
function extractCity(destination?: string): string | null {
    if (!destination) return null;
    const parts = destination.split(',').map(s => s.trim()).filter(Boolean);
    if (parts.length >= 2) {
        // If last part is "Colombia" (country), use the one before it
        const last = parts[parts.length - 1];
        if (last.toLowerCase() === 'colombia' && parts.length >= 3) {
            return parts[parts.length - 2];
        }
        return last;
    }
    // No commas: return the full address string and let Nominatim figure it out
    return destination.trim() || null;
}

// ─── Tracking Detail Panel ──────────────────────────────────────────────────

interface TrackingDetailProps {
    shipment: Shipment;
    onClose: () => void;
    onCancel: (id: string) => void;
    cancelingId: string | null;
    isCancelled: boolean;
}

function TrackingDetail({ shipment, onClose, onCancel, cancelingId, isCancelled }: TrackingDetailProps) {
    const [tracking, setTracking] = useState<{ loading: boolean; data?: any; error?: string }>({ loading: false });

    // Parse destination address for map
    const destination = shipment.destination_address || '';
    const city = extractCity(destination) || '';

    useEffect(() => {
        if (shipment.tracking_number) {
            setTracking({ loading: true });
            trackShipmentAction(shipment.tracking_number).then(res => {
                if ('data' in res && res.success) {
                    setTracking({ loading: false, data: res.data });
                } else {
                    setTracking({ loading: false, error: (res as any).message || 'No disponible' });
                }
            }).catch(err => {
                setTracking({ loading: false, error: err.message });
            });
        }
    }, [shipment.id, shipment.tracking_number]);

    const canelId = shipment.tracking_number || shipment.id.toString();

    return (
        <div className="flex flex-col h-full">
            {/* Header */}
            <div className="flex items-start justify-between p-5 border-b border-gray-100 dark:border-gray-700">
                <div className="flex-1 min-w-0">
                    <p className="text-xs text-gray-400 dark:text-gray-500 font-medium uppercase tracking-wider mb-1">Detalle de Envío</p>
                    <h3 className="text-base font-bold text-gray-900 dark:text-white truncate">
                        {shipment.customer_name || shipment.client_name || 'Cliente desconocido'}
                    </h3>
                    {shipment.destination_address && (
                        <div className="flex items-center gap-1 mt-0.5 text-xs text-gray-500 dark:text-gray-400 dark:text-gray-500">
                            <MapPin size={11} className="flex-shrink-0" />
                            <span className="truncate">{shipment.destination_address}</span>
                        </div>
                    )}
                </div>
                <button
                    onClick={onClose}
                    className="ml-3 p-1.5 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:text-gray-300 transition-colors flex-shrink-0"
                >
                    <X size={16} />
                </button>
            </div>

            {/* Scrollable content */}
            <div className="flex-1 overflow-y-auto">
                {/* Info strip */}
                <div className="grid grid-cols-2 gap-3 p-4 border-b border-gray-50">
                    <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
                        <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold tracking-wider mb-1">Tracking</p>
                        <p className="text-sm font-mono font-semibold text-gray-900 dark:text-white break-all">
                            {shipment.tracking_number || 'Sin tracking'}
                        </p>
                    </div>
                    <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
                        <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold tracking-wider mb-1">Estado</p>
                        <div className="flex items-center gap-1.5 flex-wrap">
                            <StatusBadge status={shipment.status} />
                            {shipment.is_test && (
                                <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-bold bg-orange-100 text-orange-700 border border-orange-300 uppercase tracking-widest">TEST</span>
                            )}
                        </div>
                    </div>
                    {shipment.carrier && (
                        <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
                            <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold tracking-wider mb-1">Transportista</p>
                            <div className="flex items-center gap-1.5">
                                <Truck size={13} className="text-gray-400 dark:text-gray-500" />
                                <p className="text-sm font-semibold text-gray-900 dark:text-white">
                                    {shipment.carrier.split('(')[0].trim()}
                                </p>
                            </div>
                        </div>
                    )}
                    {shipment.shipped_at && (
                        <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
                            <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold tracking-wider mb-1">Enviado</p>
                            <div className="flex items-center gap-1.5">
                                <Calendar size={13} className="text-gray-400 dark:text-gray-500" />
                                <p className="text-sm font-semibold text-gray-900 dark:text-white">{formatDate(shipment.shipped_at)}</p>
                            </div>
                        </div>
                    )}
                    {shipment.delivered_at && (
                        <div className="bg-emerald-50 rounded-lg p-3">
                            <p className="text-[10px] text-emerald-600 uppercase font-bold tracking-wider mb-1">Entregado</p>
                            <div className="flex items-center gap-1.5">
                                <CheckCircle2 size={13} className="text-emerald-500" />
                                <p className="text-sm font-semibold text-emerald-700">{formatDate(shipment.delivered_at)}</p>
                            </div>
                        </div>
                    )}
                    {shipment.estimated_delivery && !shipment.delivered_at && (
                        <div className="bg-blue-50 rounded-lg p-3">
                            <p className="text-[10px] text-blue-600 uppercase font-bold tracking-wider mb-1">Entrega Est.</p>
                            <div className="flex items-center gap-1.5">
                                <Clock size={13} className="text-blue-500" />
                                <p className="text-sm font-semibold text-blue-700">{formatDate(shipment.estimated_delivery)}</p>
                            </div>
                        </div>
                    )}
                    {shipment.created_at && (
                        <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
                            <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold tracking-wider mb-1">Creado</p>
                            <div className="flex items-center gap-1.5">
                                <Calendar size={13} className="text-gray-400 dark:text-gray-500" />
                                <p className="text-sm font-semibold text-gray-900 dark:text-white">{formatDate(shipment.created_at)}</p>
                            </div>
                        </div>
                    )}
                    {shipment.order_id && (
                        <div className="col-span-2 bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
                            <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold tracking-wider mb-1">ID Orden</p>
                            <div className="flex items-center gap-1.5">
                                <Hash size={13} className="text-gray-400 dark:text-gray-500 flex-shrink-0" />
                                <p className="text-xs font-mono text-gray-700 dark:text-gray-200 break-all">{shipment.order_id}</p>
                            </div>
                        </div>
                    )}
                    {/* Sección de contacto del cliente */}
                    {(shipment.customer_email || shipment.customer_phone) && (
                        <>
                            {shipment.customer_email && (
                                <div className="bg-blue-50 rounded-lg p-3">
                                    <p className="text-[10px] text-blue-600 uppercase font-bold tracking-wider mb-1">Email</p>
                                    <p className="text-xs text-blue-900 break-all">{shipment.customer_email}</p>
                                </div>
                            )}
                            {shipment.customer_phone && (
                                <div className="bg-green-50 rounded-lg p-3">
                                    <p className="text-[10px] text-green-600 uppercase font-bold tracking-wider mb-1">Teléfono</p>
                                    <p className="text-xs text-green-900 break-all">{shipment.customer_phone}</p>
                                </div>
                            )}
                        </>
                    )}
                </div>

                {/* Action buttons */}
                <div className="flex gap-2 px-4 pt-3 pb-2">
                    {shipment.guide_url && (
                        <a
                            href={shipment.guide_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex-1 flex items-center justify-center gap-2 py-2 px-3 rounded-lg bg-blue-600 hover:bg-blue-700 text-white text-xs font-semibold transition-colors shadow-sm"
                        >
                            <FileText size={13} />
                            Ver Guía
                        </a>
                    )}
                    {shipment.tracking_url && (
                        <a
                            href={shipment.tracking_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex-1 flex items-center justify-center gap-2 py-2 px-3 rounded-lg bg-purple-600 hover:bg-purple-700 text-white border border-purple-600 text-xs font-semibold transition-colors dark:bg-purple-700 dark:hover:bg-purple-800"
                        >
                            <Navigation size={13} />
                            Rastrear
                        </a>
                    )}
                    {isCancelled ? (
                        <div className="flex items-center justify-center gap-1.5 py-2 px-3 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-400 dark:text-gray-500 text-xs font-semibold border border-gray-200 dark:border-gray-600 cursor-default">
                            <X size={13} /> Cancelado
                        </div>
                    ) : (
                        <button
                            onClick={() => onCancel(canelId)}
                            disabled={cancelingId === canelId}
                            className="flex items-center justify-center gap-1 py-2 px-3 rounded-lg bg-red-600 hover:bg-red-700 text-white border border-red-600 text-xs font-semibold transition-colors disabled:opacity-50 dark:bg-red-700 dark:hover:bg-red-800"
                            title="Cancelar envío"
                        >
                            {cancelingId === canelId
                                ? <RefreshCw size={13} className="animate-spin" />
                                : <><X size={13} /> Cancelar</>
                            }
                        </button>
                    )}
                </div>

                {/* ─── Costos ─────────────────────────────────────────── */}
                {(shipment.shipping_cost != null || shipment.insurance_cost != null || shipment.total_cost != null) && (
                    <div className="px-4 py-3 border-t border-gray-50">
                        <div className="flex items-center gap-1.5 mb-2">
                            <DollarSign size={12} className="text-gray-400 dark:text-gray-500" />
                            <p className="text-xs font-bold text-gray-500 dark:text-gray-400 dark:text-gray-500 uppercase tracking-wider">Costos</p>
                        </div>
                        <div className="grid grid-cols-3 gap-2">
                            {shipment.shipping_cost != null && (
                                <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-2.5">
                                    <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold mb-1">Envío</p>
                                    <p className="text-sm font-semibold text-gray-900 dark:text-white">{formatMoney(shipment.shipping_cost)}</p>
                                </div>
                            )}
                            {shipment.insurance_cost != null && (
                                <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-2.5">
                                    <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold mb-1">Seguro</p>
                                    <p className="text-sm font-semibold text-gray-900 dark:text-white">{formatMoney(shipment.insurance_cost)}</p>
                                </div>
                            )}
                            {shipment.total_cost != null && (
                                <div className="bg-emerald-50 rounded-lg p-2.5">
                                    <p className="text-[10px] text-emerald-600 uppercase font-bold mb-1">Total</p>
                                    <p className="text-sm font-bold text-emerald-700">{formatMoney(shipment.total_cost)}</p>
                                </div>
                            )}
                        </div>
                    </div>
                )}

                {/* ─── Paquete ─────────────────────────────────────────── */}
                {(shipment.weight != null || shipment.length != null || shipment.width != null || shipment.height != null) && (
                    <div className="px-4 py-3 border-t border-gray-50">
                        <div className="flex items-center gap-1.5 mb-2">
                            <Box size={12} className="text-gray-400 dark:text-gray-500" />
                            <p className="text-xs font-bold text-gray-500 dark:text-gray-400 dark:text-gray-500 uppercase tracking-wider">Paquete</p>
                        </div>
                        <div className="flex flex-wrap gap-2">
                            {shipment.weight != null && (
                                <div className="flex items-center gap-2 bg-gray-50 dark:bg-gray-700 rounded-lg px-3 py-2">
                                    <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold">Peso</p>
                                    <p className="text-sm font-semibold text-gray-900 dark:text-white">{shipment.weight} kg</p>
                                </div>
                            )}
                            {(shipment.length != null || shipment.width != null || shipment.height != null) && (
                                <div className="flex items-center gap-2 bg-gray-50 dark:bg-gray-700 rounded-lg px-3 py-2">
                                    <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold">Dim.</p>
                                    <p className="text-sm font-semibold text-gray-900 dark:text-white">
                                        {shipment.length ?? '?'} × {shipment.width ?? '?'} × {shipment.height ?? '?'} cm
                                    </p>
                                </div>
                            )}
                        </div>
                    </div>
                )}

                {/* ─── Logística ───────────────────────────────────────── */}
                {(shipment.warehouse_name || shipment.driver_name || shipment.is_last_mile) && (
                    <div className="px-4 py-3 border-t border-gray-50">
                        <div className="flex items-center gap-1.5 mb-2">
                            <Building2 size={12} className="text-gray-400 dark:text-gray-500" />
                            <p className="text-xs font-bold text-gray-500 dark:text-gray-400 dark:text-gray-500 uppercase tracking-wider">Logística</p>
                        </div>
                        <div className="grid grid-cols-2 gap-2">
                            {shipment.warehouse_name && (
                                <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-2.5">
                                    <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold mb-1">Almacén</p>
                                    <div className="flex items-center gap-1.5">
                                        <Building2 size={12} className="text-gray-400 dark:text-gray-500 flex-shrink-0" />
                                        <p className="text-sm font-semibold text-gray-900 dark:text-white">{shipment.warehouse_name}</p>
                                    </div>
                                </div>
                            )}
                            {shipment.driver_name && (
                                <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-2.5">
                                    <p className="text-[10px] text-gray-400 dark:text-gray-500 uppercase font-bold mb-1">Conductor</p>
                                    <div className="flex items-center gap-1.5">
                                        <User size={12} className="text-gray-400 dark:text-gray-500 flex-shrink-0" />
                                        <p className="text-sm font-semibold text-gray-900 dark:text-white">{shipment.driver_name}</p>
                                    </div>
                                </div>
                            )}
                            {shipment.is_last_mile && (
                                <div className="col-span-2 flex items-center gap-2 bg-purple-50 border border-purple-100 rounded-lg px-3 py-2">
                                    <div className="w-2 h-2 rounded-full bg-purple-500 flex-shrink-0" />
                                    <p className="text-xs font-semibold text-purple-700">Envío de Última Milla</p>
                                </div>
                            )}
                        </div>
                    </div>
                )}

                {/* ─── Notas de entrega ────────────────────────────────── */}
                {shipment.delivery_notes && (
                    <div className="px-4 py-3 border-t border-gray-50">
                        <div className="flex items-center gap-1.5 mb-2">
                            <StickyNote size={12} className="text-gray-400 dark:text-gray-500" />
                            <p className="text-xs font-bold text-gray-500 dark:text-gray-400 dark:text-gray-500 uppercase tracking-wider">Notas de Entrega</p>
                        </div>
                        <p className="text-sm text-gray-700 dark:text-gray-200 bg-amber-50 border border-amber-100 rounded-lg p-3 leading-relaxed">
                            {shipment.delivery_notes}
                        </p>
                    </div>
                )}

                {/* Tracking Timeline */}
                <div className="px-4 py-3">
                    <p className="text-xs font-bold text-gray-500 dark:text-gray-400 dark:text-gray-500 uppercase tracking-wider mb-3">Historial de rastreo</p>
                    {tracking.loading ? (
                        <div className="flex items-center gap-2 text-sm text-gray-400 dark:text-gray-500 py-4">
                            <RefreshCw size={16} className="animate-spin" />
                            <span>Consultando rastreo...</span>
                        </div>
                    ) : tracking.error ? (
                        <div className="flex items-start gap-2 bg-amber-50 border border-amber-200 rounded-lg p-3 text-xs text-amber-700">
                            <AlertTriangle size={14} className="flex-shrink-0 mt-0.5" />
                            <span>{tracking.error}</span>
                        </div>
                    ) : tracking.data?.history?.length > 0 ? (
                        <div className="relative pl-5">
                            <div className="absolute left-1.5 top-1 bottom-2 w-px bg-gray-200" />
                            <div className="space-y-4">
                                {tracking.data.history.map((event: EnvioClickTrackHistory, idx: number) => (
                                    <div key={idx} className="relative">
                                        <div className={`absolute -left-5 top-0.5 w-3 h-3 rounded-full ring-2 ring-white ${idx === 0 ? 'bg-blue-500' : 'bg-gray-300'}`} />
                                        <div>
                                            <div className="flex items-baseline justify-between gap-2">
                                                <p className={`text-sm font-semibold ${idx === 0 ? 'text-blue-700' : 'text-gray-800 dark:text-gray-100'}`}>{event.status}</p>
                                                <p className="text-[10px] text-gray-400 dark:text-gray-500 flex-shrink-0">{event.date}</p>
                                            </div>
                                            {event.description && <p className="text-xs text-gray-500 dark:text-gray-400 dark:text-gray-500 mt-0.5">{event.description}</p>}
                                            {event.location && (
                                                <div className="flex items-center gap-1 mt-1 text-[10px] text-gray-400 dark:text-gray-500">
                                                    <MapPin size={9} />
                                                    <span>{event.location}</span>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    ) : !shipment.tracking_number ? (
                        <p className="text-xs text-gray-400 dark:text-gray-500 italic">Este envío no tiene número de tracking.</p>
                    ) : (
                        <p className="text-xs text-gray-400 dark:text-gray-500 italic">No hay historial disponible.</p>
                    )}
                </div>

                {/* Reference Map */}
                {destination && (
                    <div className="px-4 pb-5">
                        <p className="text-xs font-bold text-gray-500 dark:text-gray-400 dark:text-gray-500 uppercase tracking-wider mb-2">Mapa de referencia</p>
                        <div style={{ height: '200px' }} className="rounded-xl overflow-hidden border border-gray-200 dark:border-gray-600 dark:border-gray-700">
                            <MapComponent
                                address={destination}
                                city={city}
                                height="200px"
                            />
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}

// ─── Main Component ─────────────────────────────────────────────────────────

interface ShipmentListProps {
    selectedBusinessId?: number | null;
}

export default function ShipmentList({ selectedBusinessId = null }: ShipmentListProps) {
    const router = useRouter();
    const searchParams = useSearchParams();
    const { permissions, isSuperAdmin } = usePermissions();
    const canCreate = useHasPermission('Envios', 'Create');
    const canDelete = useHasPermission('Envios', 'Delete');
    const [loading, setLoading] = useState(true);
    const [shipments, setShipments] = useState<Shipment[]>([]);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const [selectedShipment, setSelectedShipment] = useState<Shipment | null>(null);
    const [cancelingId, setCancelingId] = useState<string | null>(null);
    const [isManualModalOpen, setIsManualModalOpen] = useState(false);
    const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
    const [isSyncing, setIsSyncing] = useState(false);
    const [isSyncModalOpen, setIsSyncModalOpen] = useState(false);
    const [isCancelingBatch, setIsCancelingBatch] = useState(false);
    const [cancelModalData, setCancelModalData] = useState<{ isOpen: boolean; type: 'single' | 'batch'; shipmentId?: string } | null>(null);

    // Auto-inject business_id for non-super-admins
    const defaultBusinessId = (!isSuperAdmin && permissions?.business_id) ? permissions.business_id : undefined;
    const sseBusinessId = permissions?.business_id || 0;

    // Optimistic local cancellation via SSE: mark shipment as cancelled without full reload
    useShipmentSSE({
        businessId: sseBusinessId,
        onShipmentCancelled: (data) => {
            const shipmentId: number | undefined = (data as any).shipment_id;
            if (shipmentId) {
                setShipments(prev => prev.map(s =>
                    s.id === shipmentId ? { ...s, status: 'cancelled' } : s
                ));
                setSelectedShipment(prev =>
                    prev?.id === shipmentId ? { ...prev, status: 'cancelled' } : prev
                );
            }
        },
        onCancelFailed: (data) => {
            const errorMsg = data.error_message || 'Error al cancelar el envio';
            alert(`Cancelacion fallida: ${errorMsg}`);
        },
    });

    const [filters, setFilters] = useState<GetShipmentsParams>({
        page: Number(searchParams.get('page')) || 1,
        page_size: Number(searchParams.get('page_size')) || 20,
        tracking_number: searchParams.get('tracking_number') || undefined,
        order_id: searchParams.get('order_id') || undefined,
        carrier: searchParams.get('carrier') || undefined,
        status: searchParams.get('status') || undefined,
        customer_name: searchParams.get('customer_name') || undefined,
        is_test: searchParams.get('is_test') !== null ? searchParams.get('is_test') === 'true' : undefined,
        business_id: defaultBusinessId,
    });

    const fetchShipments = async () => {
        setLoading(true);
        try {
            const params: GetShipmentsParams = {
                ...filters,
                business_id: isSuperAdmin
                    ? (selectedBusinessId !== null ? selectedBusinessId : undefined)
                    : defaultBusinessId,
            };
            const response = await getShipmentsAction(params);
            if (response.success) {
                setShipments(response.data);
                setPage(response.page);
                setTotalPages(response.total_pages);
                setTotal(response.total);
            }
        } catch (error) {
            console.error('Error fetching shipments:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchShipments();
    }, [filters, selectedBusinessId]);

    const updateFilters = (newFilters: Partial<GetShipmentsParams>) => {
        const updated = { ...filters, ...newFilters };
        if (!newFilters.page && newFilters.page !== 0) updated.page = 1;
        setFilters(updated);

        const params = new URLSearchParams();
        Object.entries(updated).forEach(([key, value]) => {
            if (value !== undefined && value !== null && value !== '' && key !== 'business_id') {
                params.set(key, String(value));
            }
        });
        router.push(`?${params.toString()}`);
    };

    // Resolve business_id for transport operations (super admin needs explicit business context)

    const handleCancel = async (id: string) => {
        setCancelModalData({ isOpen: true, type: 'single', shipmentId: id });
    };

    const handleBatchCancel = async () => {
        if (selectedIds.size === 0) return;
        setCancelModalData({ isOpen: true, type: 'batch' });
    };

    const confirmCancel = async () => {
        if (!cancelModalData) return;

        if (cancelModalData.type === 'single' && cancelModalData.shipmentId) {
            setCancelingId(cancelModalData.shipmentId);
            setCancelModalData(null);
            try {
                const response = await cancelShipmentAction(cancelModalData.shipmentId);
                if (!response.success) {
                    alert(`Error: ${response.message}`);
                }
            } catch (error: any) {
                alert(`Error: ${error.message}`);
            } finally {
                setCancelingId(null);
            }
        } else if (cancelModalData.type === 'batch') {
            setIsCancelingBatch(true);
            setCancelModalData(null);
            try {
                const orders = Array.from(selectedIds).map(id => {
                    const shipment = shipments.find(s => s.id === id);
                    return {
                        trackingCode: shipment?.tracking_number || id.toString(),
                        motivo: 'Cancelado por el usuario'
                    };
                });

                const response = await cancelBatchShipmentAction({ orders });
                if (!response.success) {
                    alert(`Error: ${response.message}`);
                } else {
                    setSelectedIds(new Set());
                }
            } catch (error: any) {
                alert(`Error: ${error.message}`);
            } finally {
                setIsCancelingBatch(false);
            }
        }
    };

    const toggleSelection = (id: number) => {
        const newSet = new Set(selectedIds);
        if (newSet.has(id)) newSet.delete(id);
        else newSet.add(id);
        setSelectedIds(newSet);
    };

    const toggleAll = () => {
        if (selectedIds.size === shipments.length) {
            setSelectedIds(new Set());
        } else {
            setSelectedIds(new Set(shipments.map(s => s.id)));
        }
    };

    return (
        <div className="flex flex-col h-full gap-4" style={{ height: 'calc(100vh - 120px)' }}>

            {/* ─── Top bar ─── */}
            <div className="flex items-start flex-shrink-0">
                {/* Left: icon + title + subtitle */}
                <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0">
                        <Package size={20} className="text-blue-600" />
                    </div>
                    <div>
                        <div className="flex items-center gap-2">
                            <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Envíos</h2>
                            {total > 0 && (
                                <span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-blue-50 text-blue-600 border border-blue-100">
                                    {total}
                                </span>
                            )}
                        </div>
                        <p className="text-sm text-gray-400 dark:text-gray-500 mt-0.5">Gestiona y rastrea todos tus envíos</p>
                    </div>
                </div>
            </div>

            {/* ─── Status chips ─── */}
            <div className="flex items-center gap-2 flex-shrink-0 overflow-x-auto pb-0.5">
                {/* Chip "Todos" */}
                <button
                    onClick={() => updateFilters({ status: undefined })}
                    className={`flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition-all ${!filters.status ? 'bg-purple-600 text-white dark:bg-purple-700' : 'bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-600 dark:hover:bg-gray-600'
                        }`}
                >
                    Todos <span className="opacity-70">{total}</span>
                </button>
                {/* Chips por estado */}
                {CHIP_STATUS_OPTIONS.map(({ value, label, icon: Icon, activeClass }) => {
                    const count = shipments.filter((s) => s.status === value).length;
                    const isActive = filters.status === value;
                    return (
                        <button
                            key={value}
                            onClick={() => updateFilters({ status: isActive ? undefined : value })}
                            className={`flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition-all ${isActive ? activeClass : 'bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-600 dark:hover:bg-gray-600'
                                }`}
                        >
                            <Icon size={11} />
                            {label}
                            {count > 0 && <span className="opacity-70">{count}</span>}
                        </button>
                    );
                })}
            </div>

            {/* ─── Filters ─── */}
            <div className="flex-shrink-0 bg-white dark:bg-gray-800 dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 px-4 py-3">
                <div className="flex gap-3">
                    {/* Búsqueda por nombre del cliente */}
                    <div className="relative flex-1">
                        <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500" />
                        <input
                            type="text"
                            placeholder="Buscar por nombre del cliente..."
                            className="w-full pl-9 pr-3 py-2 border border-gray-200 dark:border-gray-600 dark:border-gray-700 rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-400 text-sm bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder:text-gray-400 dark:text-gray-500 transition-colors"
                            value={filters.customer_name || ''}
                            onChange={(e) => updateFilters({ customer_name: e.target.value || undefined })}
                        />
                    </div>
                    {/* Select estado */}
                    <select
                        className="px-3 py-2 border border-gray-200 dark:border-gray-600 dark:border-gray-700 rounded-lg text-sm text-gray-600 dark:text-gray-300 bg-gray-50 dark:bg-gray-700 focus:ring-2 focus:ring-blue-500/20 min-w-[140px] transition-colors"
                        value={filters.status || ''}
                        onChange={(e) => updateFilters({ status: e.target.value || undefined })}
                    >
                        <option value="">Todos los estados</option>
                        <option value="pending">Pendiente</option>
                        <option value="in_transit">En Tránsito</option>
                        <option value="delivered">Entregado</option>
                        <option value="failed">Fallido</option>
                    </select>
                    {/* Select entorno */}
                    <select
                        className="px-3 py-2 border border-gray-200 dark:border-gray-600 dark:border-gray-700 rounded-lg text-sm text-gray-600 dark:text-gray-300 bg-gray-50 dark:bg-gray-700 focus:ring-2 focus:ring-orange-500/20 min-w-[140px] transition-colors"
                        value={filters.is_test === undefined ? '' : filters.is_test ? 'test' : 'production'}
                        onChange={(e) => {
                            const val = e.target.value;
                            updateFilters({ is_test: val === '' ? undefined : val === 'test' });
                        }}
                    >
                        <option value="">Prod + TEST</option>
                        <option value="production">Solo producción</option>
                        <option value="test">Solo TEST</option>
                    </select>
                    <button
                        onClick={() => setIsSyncModalOpen(true)}
                        className="px-3 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-sm font-semibold flex items-center gap-2 transition-colors whitespace-nowrap"
                        title="Consulta el carrier y actualiza los estados de las guías activas"
                    >
                        <RefreshCw size={14} />
                        Sincronizar Estados
                    </button>
                </div>
            </div>

            {/* ─── Split Panel ─── */}
            <div className="flex gap-4 flex-1 min-h-0">

                {/* LEFT — Shipment cards list */}
                <div className="w-1/2 flex flex-col min-h-0 bg-white dark:bg-gray-800 dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
                    
                    {/* Batch Actions Header */}
                    {shipments.length > 0 && (
                        <div className="flex items-center justify-between px-4 py-2 border-b border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50 flex-shrink-0">
                            <div className="flex items-center gap-2">
                                <input 
                                    type="checkbox" 
                                    checked={selectedIds.size === shipments.length && shipments.length > 0} 
                                    onChange={toggleAll}
                                    className="rounded border-gray-300 text-purple-600 focus:ring-purple-500 w-4 h-4"
                                />
                                <span className="text-xs font-semibold text-gray-600 dark:text-gray-300">
                                    {selectedIds.size > 0 ? `${selectedIds.size} seleccionados` : 'Seleccionar todos'}
                                </span>
                            </div>
                            {selectedIds.size > 0 && (
                                <button
                                    onClick={handleBatchCancel}
                                    disabled={isCancelingBatch}
                                    className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-red-50 hover:bg-red-100 text-red-600 border border-red-200 text-xs font-semibold transition-colors disabled:opacity-50"
                                >
                                    {isCancelingBatch ? <RefreshCw size={12} className="animate-spin" /> : <XCircle size={12} />}
                                    Cancelar Seleccionados
                                </button>
                            )}
                        </div>
                    )}

                    {/* List */}
                    <div className="flex-1 overflow-y-auto divide-y divide-gray-50">
                        {loading ? (
                            <div className="flex flex-col items-center justify-center h-40 gap-3 text-gray-400 dark:text-gray-500">
                                <RefreshCw size={24} className="animate-spin" />
                                <p className="text-sm">Cargando envíos...</p>
                            </div>
                        ) : shipments.length === 0 ? (
                            <div className="flex flex-col items-center justify-center h-40 gap-2 text-gray-400 dark:text-gray-500">
                                <Package size={32} strokeWidth={1.5} />
                                <p className="text-sm">No hay envíos disponibles</p>
                            </div>
                        ) : (
                            shipments.map((shipment) => {
                                const isSelected = selectedShipment?.id === shipment.id;
                                const city = extractCity(shipment.destination_address);
                                const clientName = (shipment.customer_name || shipment.client_name)?.trim() || null;
                                const statusCfg = STATUS_CONFIG[shipment.status] || { label: shipment.status, color: 'bg-gray-100 text-gray-600 dark:text-gray-300 border-gray-200 dark:border-gray-600 dark:border-gray-700', icon: null, border: 'border-gray-300' };

                                return (
                                    <button
                                        key={shipment.id}
                                        onClick={() => setSelectedShipment(isSelected ? null : shipment)}
                                        className={`w-full text-left transition-all duration-150 hover:bg-gray-50 dark:hover:bg-gray-700 ${isSelected
                                                ? 'bg-purple-100/50 dark:bg-purple-900/30 border-l-[3px] border-purple-500'
                                                : `border-l-[3px] ${statusCfg.border}`
                                            }`}
                                    >
                                        <div className="flex items-stretch w-full">
                                            {/* Checkbox zone */}
                                            <div 
                                                className="px-3 py-4 flex items-start justify-center cursor-pointer border-r border-transparent hover:border-gray-100 dark:hover:border-gray-600"
                                                onClick={(e) => { e.stopPropagation(); toggleSelection(shipment.id); }}
                                            >
                                                <input 
                                                    type="checkbox" 
                                                    checked={selectedIds.has(shipment.id)} 
                                                    onChange={() => {}} // Controlled by the div click
                                                    className="mt-1 rounded border-gray-300 text-purple-600 focus:ring-purple-500 w-4 h-4 cursor-pointer"
                                                />
                                            </div>
                                            
                                            {/* Main content click zone */}
                                            <div 
                                                className="flex-1 px-3 py-3.5 min-w-0"
                                                onClick={() => setSelectedShipment(isSelected ? null : shipment)}
                                            >
                                                {/* Row 1: Client name + destination city */}
                                        <div className="flex items-center justify-between gap-2 mb-1.5">
                                            <div className="flex items-center gap-1.5 min-w-0">
                                                {!clientName && <Package size={11} className="text-gray-300 flex-shrink-0" />}
                                                <p className={`text-sm font-semibold truncate ${clientName ? 'text-gray-900 dark:text-white' : 'text-gray-400 dark:text-gray-500 italic'}`}>
                                                    {clientName || 'Sin destinatario'}
                                                </p>
                                            </div>
                                            <div className="flex items-center gap-1 flex-shrink-0">
                                                {city && (
                                                    <span className="flex items-center gap-0.5 text-[10px] text-gray-500 dark:text-gray-400 dark:text-gray-500 bg-gray-100 px-1.5 py-0.5 rounded-full">
                                                        <MapPin size={8} />{city}
                                                    </span>
                                                )}
                                                {shipment.total_cost != null && (
                                                    <span className="text-[10px] font-semibold text-emerald-700 bg-emerald-50 px-1.5 py-0.5 rounded-full border border-emerald-100">
                                                        {formatMoney(shipment.total_cost)}
                                                    </span>
                                                )}
                                            </div>
                                        </div>

                                        {/* Row 2: Status + TEST badge */}
                                        <div className="flex items-center gap-1.5 mb-1.5 flex-wrap">
                                            <StatusBadge status={shipment.status} />
                                            {shipment.is_test && (
                                                <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-bold bg-orange-100 text-orange-700 border border-orange-300 uppercase tracking-widest">TEST</span>
                                            )}
                                        </div>

                                        {/* Row 3: Tracking + carrier + date */}
                                        <div className="flex items-center gap-2 flex-wrap">
                                            {shipment.tracking_number && (
                                                <span className="text-[10px] font-mono text-gray-400 dark:text-gray-500 bg-gray-100 px-1.5 py-0.5 rounded">
                                                    #{shipment.tracking_number.slice(-10)}
                                                </span>
                                            )}
                                            {shipment.carrier && (
                                                <span className="text-[10px] text-gray-400 dark:text-gray-500 flex items-center gap-0.5">
                                                    <Truck size={9} />{shipment.carrier}
                                                </span>
                                            )}
                                            {(shipment.shipped_at || shipment.created_at) && (
                                                <span className="text-[10px] text-gray-400 dark:text-gray-500 flex items-center gap-0.5">
                                                    <Calendar size={9} />{formatDate(shipment.shipped_at || shipment.created_at)}
                                                </span>
                                            )}
                                        </div>
                                    </div>
                                    </div>
                                    </button>
                                );
                            })
                        )}
                    </div>

                    {/* Pagination */}
                    {totalPages > 1 && (
                        <div className="flex items-center justify-between px-4 py-3 border-t border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50 flex-shrink-0">
                            <p className="text-xs text-gray-500 dark:text-gray-400 dark:text-gray-500">
                                Pág. <span className="font-semibold text-gray-700 dark:text-gray-200">{page}</span> de {totalPages}
                            </p>
                            <div className="flex items-center gap-1">
                                <button
                                    onClick={() => updateFilters({ page: page - 1 })}
                                    disabled={page === 1}
                                    className="p-1.5 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-600 dark:text-gray-300 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                                >
                                    <ChevronLeft size={15} />
                                </button>
                                <button
                                    onClick={() => updateFilters({ page: page + 1 })}
                                    disabled={page === totalPages}
                                    className="p-1.5 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-600 dark:text-gray-300 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                                >
                                    <ChevronRight size={15} />
                                </button>
                            </div>
                        </div>
                    )}
                </div>

                {/* RIGHT — Detail panel */}
                <div className="w-1/2 min-h-0 bg-white dark:bg-gray-800 dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
                    {selectedShipment ? (
                        <TrackingDetail
                            shipment={selectedShipment}
                            onClose={() => setSelectedShipment(null)}
                            onCancel={handleCancel}
                            cancelingId={cancelingId}
                            isCancelled={selectedShipment.status === 'cancelled'}
                        />
                    ) : (
                        <div className="flex flex-col items-center justify-center h-full gap-4 text-gray-400 dark:text-gray-500 select-none">
                            <div className="p-5 bg-gray-50 dark:bg-gray-700 rounded-full">
                                <Package size={36} strokeWidth={1.2} className="text-gray-300" />
                            </div>
                            <div className="text-center">
                                <p className="text-sm font-medium text-gray-500 dark:text-gray-400 dark:text-gray-500">Selecciona un envío</p>
                                <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">para ver su información de rastreo</p>
                            </div>
                            <div className="flex items-center gap-1 text-xs text-gray-300">
                                <ChevronLeft size={12} />
                                <span>Elige de la lista</span>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Modal para envío manual */}
            <ManualShipmentModal
                isOpen={isManualModalOpen}
                onClose={() => setIsManualModalOpen(false)}
                onSuccess={fetchShipments}
            />

            <SyncProgressModal
                isOpen={isSyncModalOpen}
                onClose={() => setIsSyncModalOpen(false)}
                businessId={selectedBusinessId}
                onCompleted={fetchShipments}
            />

            {/* Modal de confirmación de cancelación */}
            {cancelModalData?.isOpen && (
                <Modal 
                    isOpen={cancelModalData.isOpen} 
                    onClose={() => setCancelModalData(null)} 
                    title={cancelModalData.type === 'single' ? "Cancelar Envío" : "Cancelar Envíos Masivos"}
                    size="md"
                >
                    <div className="p-5">
                        {/* Icon + message */}
                        <div className="flex items-start gap-4 mb-5">
                            <div className="w-12 h-12 rounded-2xl bg-red-100 dark:bg-red-900/30 flex items-center justify-center flex-shrink-0 shadow-sm">
                                <AlertTriangle size={22} className="text-red-600 dark:text-red-400" />
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-1">
                                    {cancelModalData.type === 'single' 
                                        ? '¿Cancelar este envío?' 
                                        : `¿Cancelar ${selectedIds.size} envíos?`}
                                </p>
                                <p className="text-xs text-gray-500 dark:text-gray-400 leading-relaxed">
                                    {cancelModalData.type === 'single'
                                        ? 'Esta acción notificará a la transportadora y no podrá revertirse.'
                                        : 'Esta acción notificará a la transportadora por cada envío seleccionado.'}
                                </p>
                            </div>
                        </div>

                        {/* Shipment details card */}
                        {cancelModalData.type === 'single' && cancelModalData.shipmentId && (() => {
                            const s = shipments.find(sh => (sh.tracking_number || sh.id.toString()) === cancelModalData.shipmentId);
                            if (!s) return null;
                            return (
                                <div className="rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden mb-4">
                                    <div className="bg-gray-50 dark:bg-gray-800 px-4 py-2 border-b border-gray-200 dark:border-gray-700">
                                        <span className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Detalles del envío</span>
                                    </div>
                                    <div className="divide-y divide-gray-100 dark:divide-gray-700">
                                        <div className="flex justify-between items-center px-4 py-2.5 text-sm">
                                            <span className="text-gray-500 dark:text-gray-400 flex items-center gap-1.5"><User size={13}/>Destinatario</span>
                                            <span className="font-semibold text-gray-900 dark:text-gray-100 truncate max-w-[180px] text-right">{s.customer_name || s.client_name || 'Desconocido'}</span>
                                        </div>
                                        <div className="flex justify-between items-center px-4 py-2.5 text-sm">
                                            <span className="text-gray-500 dark:text-gray-400 flex items-center gap-1.5"><MapPin size={13}/>Destino</span>
                                            <span className="font-semibold text-gray-900 dark:text-gray-100 truncate max-w-[180px] text-right" title={s.destination_address}>{s.destination_address || '—'}</span>
                                        </div>
                                        <div className="flex justify-between items-center px-4 py-2.5 text-sm">
                                            <span className="text-gray-500 dark:text-gray-400 flex items-center gap-1.5"><Truck size={13}/>Transportadora</span>
                                            <span className="font-semibold text-gray-900 dark:text-gray-100">{s.carrier?.split('(')[0].trim() || '—'}</span>
                                        </div>
                                        {s.total_cost != null && (
                                            <div className="flex justify-between items-center px-4 py-2.5 text-sm bg-emerald-50 dark:bg-emerald-900/10">
                                                <span className="text-gray-500 dark:text-gray-400 flex items-center gap-1.5"><DollarSign size={13}/>Costo Total</span>
                                                <span className="font-bold text-emerald-600 dark:text-emerald-400">{formatMoney(s.total_cost)}</span>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            );
                        })()}

                        {/* 72h refund notice */}
                        <div className="flex gap-3 p-3 rounded-xl bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-700/50 mb-5">
                            <span className="text-lg leading-none flex-shrink-0">⏰</span>
                            <p className="text-xs text-amber-800 dark:text-amber-300 leading-relaxed">
                                <span className="font-semibold">Reembolso en hasta 72 horas hábiles.</span>{' '}
                                El valor del envío será devuelto a tu monedero dentro de 3 días hábiles una vez confirmada la cancelación.
                            </p>
                        </div>

                        {/* Action buttons */}
                        <div className="flex justify-end gap-3">
                            <Button variant="outline" onClick={() => setCancelModalData(null)}>
                                Mantener
                            </Button>
                            <button
                                onClick={confirmCancel}
                                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-red-600 hover:bg-red-700 active:bg-red-800 text-white text-sm font-semibold shadow-sm transition-colors"
                            >
                                <X size={15} />
                                Sí, Cancelar
                            </button>
                        </div>
                    </div>
                </Modal>
            )}
        </div>
    );
}
