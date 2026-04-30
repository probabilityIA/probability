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
    DollarSign, Box, User, Building2, Hash, StickyNote,
    PackageCheck, MapPinned, PauseCircle, RotateCcw
} from 'lucide-react';
import { getOrdersAction } from '@/services/modules/orders/infra/actions';
import { ManualShipmentModal } from './ManualShipmentModal';
import { SyncProgressModal } from './SyncProgressModal';
import { MiniAddressMap } from './MiniAddressMap';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';
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
    pending: { label: 'Pendiente', color: 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300 border-amber-200 dark:border-amber-600', icon: <Clock size={12} />, border: 'border-amber-400' },
    picked_up: { label: 'Recolectado', color: 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300 border-indigo-200 dark:border-indigo-600', icon: <PackageCheck size={12} />, border: 'border-indigo-400' },
    in_transit: { label: 'En tránsito', color: 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 border-blue-200 dark:border-blue-600', icon: <Truck size={12} />, border: 'border-blue-400' },
    out_for_delivery: { label: 'En reparto', color: 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 border-purple-200 dark:border-purple-600', icon: <MapPinned size={12} />, border: 'border-purple-400' },
    delivered: { label: 'Entregado', color: 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-600', icon: <CheckCircle2 size={12} />, border: 'border-emerald-400' },
    on_hold: { label: 'Novedad', color: 'bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-300 border-orange-200 dark:border-orange-600', icon: <PauseCircle size={12} />, border: 'border-orange-400' },
    returned: { label: 'Devuelto', color: 'bg-rose-100 dark:bg-rose-900/30 text-rose-700 dark:text-rose-300 border-rose-200 dark:border-rose-600', icon: <RotateCcw size={12} />, border: 'border-rose-400' },
    failed: { label: 'Fallido', color: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 border-red-200 dark:border-red-600', icon: <XCircle size={12} />, border: 'border-red-400' },
    cancelled: { label: 'Cancelado', color: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 border-red-200 dark:border-red-600', icon: <X size={12} />, border: 'border-red-400' },
};

const CHIP_STATUS_OPTIONS = [
    { value: 'pending', label: 'Pendiente', icon: Clock, activeClass: 'bg-amber-500 text-white' },
    { value: 'picked_up', label: 'Recolectado', icon: PackageCheck, activeClass: 'bg-indigo-500 text-white' },
    { value: 'in_transit', label: 'En tránsito', icon: Truck, activeClass: 'bg-blue-500 text-white' },
    { value: 'out_for_delivery', label: 'En reparto', icon: MapPinned, activeClass: 'bg-purple-500 text-white' },
    { value: 'delivered', label: 'Entregado', icon: CheckCircle2, activeClass: 'bg-emerald-500 text-white' },
    { value: 'on_hold', label: 'Novedad', icon: PauseCircle, activeClass: 'bg-orange-500 text-white' },
    { value: 'returned', label: 'Devuelto', icon: RotateCcw, activeClass: 'bg-rose-500 text-white' },
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

function SubStatusBadge({ status, carrier, detail }: { status: string; carrier?: string; detail?: string }) {
    if (!detail) return null;
    const cfg = STATUS_CONFIG[status] || { color: 'bg-gray-50 text-gray-600 dark:text-gray-300 border-gray-200 dark:border-gray-600' };
    return (
        <span
            className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-medium border-dashed border ${cfg.color}`}
            title={carrier ? `${carrier}: ${detail}` : detail}
        >
            {carrier && <span className="font-bold uppercase tracking-wider opacity-75">{carrier}</span>}
            <span className="truncate max-w-[180px]">{detail}</span>
        </span>
    );
}

function formatDate(dateStr?: string) {
    if (!dateStr) return null;
    return new Date(dateStr).toLocaleString('es-CO', {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
        timeZone: 'America/Bogota',
    });
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

    const carrierLogo = getCarrierLogo(shipment.carrier);

    return (
        <div className="flex flex-col h-full">
            <div className="flex items-start justify-between p-5 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-purple-50 to-indigo-50 dark:from-purple-950/40 dark:to-indigo-950/40">
                <div className="flex items-start gap-3 flex-1 min-w-0">
                    {carrierLogo ? (
                        <div className="w-12 h-12 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-1.5 flex items-center justify-center flex-shrink-0 shadow-sm">
                            <img src={carrierLogo} alt={shipment.carrier} className="max-w-full max-h-full object-contain" />
                        </div>
                    ) : shipment.carrier ? (
                        <div className="w-12 h-12 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 flex items-center justify-center flex-shrink-0 shadow-sm">
                            <Truck size={20} className="text-gray-400" />
                        </div>
                    ) : null}
                    <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1 flex-wrap">
                            {shipment.order_number && (
                                <span className="text-xs font-bold text-purple-700 dark:text-purple-300 bg-purple-100 dark:bg-purple-900/40 px-2 py-0.5 rounded border border-purple-200 dark:border-purple-700">
                                    {shipment.order_number}
                                </span>
                            )}
                            {shipment.carrier && (
                                <span className="text-xs text-gray-600 dark:text-gray-300 font-medium">{shipment.carrier}</span>
                            )}
                        </div>
                        <h3 className="text-base font-bold text-gray-900 dark:text-white truncate">
                            {shipment.customer_name || shipment.client_name || 'Cliente desconocido'}
                        </h3>
                        {shipment.tracking_number && (
                            <p className="text-xs font-mono text-gray-500 dark:text-gray-400 mt-0.5">#{shipment.tracking_number}</p>
                        )}
                    </div>
                </div>
                <button
                    onClick={onClose}
                    className="ml-3 p-1.5 rounded-full hover:bg-white/60 dark:hover:bg-gray-700 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:text-gray-300 transition-colors flex-shrink-0"
                >
                    <X size={16} />
                </button>
            </div>

            {/* Scrollable content */}
            <div className="flex-1 overflow-y-auto">
                {/* Direcciones — 2 columns */}
                <div className="px-4 pt-4 pb-2 border-b border-gray-50 dark:border-gray-700">
                    <div className="flex items-center gap-1.5 mb-2">
                        <MapPin size={12} className="text-gray-400 dark:text-gray-500" />
                        <p className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Direcciones</p>
                    </div>
                    <div className="grid grid-cols-2 gap-2">
                        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-900/40 rounded-lg overflow-hidden">
                            <div className="px-3 pt-2 pb-1.5">
                                <div className="flex items-center gap-1.5 mb-1">
                                    <Building2 size={11} className="text-blue-600 dark:text-blue-400" />
                                    <p className="text-[10px] text-blue-700 dark:text-blue-300 uppercase font-bold tracking-wider">Origen</p>
                                </div>
                                <p className="text-xs font-semibold text-blue-900 dark:text-blue-100 truncate" title={shipment.warehouse_name || 'Bodega principal'}>
                                    {shipment.warehouse_name || 'Bodega principal'}
                                </p>
                                {shipment.origin_address && (
                                    <p className="text-[11px] text-blue-800/80 dark:text-blue-200/80 truncate" title={`${shipment.origin_address}, ${shipment.origin_city || ''}`}>
                                        {shipment.origin_address}{shipment.origin_city ? `, ${shipment.origin_city}` : ''}
                                    </p>
                                )}
                            </div>
                            {shipment.origin_address ? (
                                <MiniAddressMap key={`origin-${shipment.id}`} address={`${shipment.origin_address}, ${shipment.origin_city || ''}`} city={shipment.origin_city || 'Colombia'} color="blue" />
                            ) : (
                                <MiniAddressMap key={`origin-empty-${shipment.id}`} color="blue" />
                            )}
                        </div>
                        <div className="bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-100 dark:border-emerald-900/40 rounded-lg overflow-hidden">
                            <div className="px-3 pt-2 pb-1.5">
                                <div className="flex items-center gap-1.5 mb-1">
                                    <MapPin size={11} className="text-emerald-600 dark:text-emerald-400" />
                                    <p className="text-[10px] text-emerald-700 dark:text-emerald-300 uppercase font-bold tracking-wider">Destino</p>
                                </div>
                                <p className="text-xs font-semibold text-emerald-900 dark:text-emerald-100 truncate" title={shipment.destination_address}>
                                    {shipment.destination_address || 'Sin destino'}
                                </p>
                                {shipment.destination_suburb && (
                                    <p className="text-[11px] text-emerald-800/80 dark:text-emerald-200/80 truncate">
                                        Barrio: {shipment.destination_suburb}
                                    </p>
                                )}
                                {(shipment.destination_city || shipment.destination_state) && (
                                    <p className="text-[11px] text-emerald-800/80 dark:text-emerald-200/80 truncate">
                                        {[shipment.destination_city, shipment.destination_state].filter(Boolean).join(', ')}
                                    </p>
                                )}
                            </div>
                            {shipment.destination_address ? (
                                <MiniAddressMap key={`dest-${shipment.id}`} address={`${shipment.destination_address}, ${shipment.destination_city || ''}`} city={shipment.destination_city || undefined} color="emerald" />
                            ) : (
                                <MiniAddressMap key={`dest-empty-${shipment.id}`} color="emerald" />
                            )}
                        </div>
                    </div>
                </div>

                {/* Info compacta: Cliente + Tracking + Estado + Fechas + Acciones */}
                <div className="px-4 py-3 border-b border-gray-50 dark:border-gray-700 space-y-3">
                    {/* Cliente + Tracking en una fila */}
                    {(shipment.customer_name || shipment.client_name || shipment.tracking_number) && (
                        <div className="grid grid-cols-2 gap-2">
                            {(shipment.customer_name || shipment.client_name) && (
                                <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-2.5">
                                    <div className="flex items-center gap-1 mb-1">
                                        <User size={10} className="text-gray-400" />
                                        <p className="text-[9px] text-gray-400 uppercase font-bold tracking-wider">Cliente</p>
                                    </div>
                                    <p className="text-xs font-semibold text-gray-900 dark:text-white truncate">
                                        {shipment.customer_name || shipment.client_name}
                                    </p>
                                    <div className="mt-1 space-y-0.5 text-[10px] text-gray-600 dark:text-gray-300">
                                        {shipment.customer_email && <p className="truncate">✉ {shipment.customer_email}</p>}
                                        {shipment.customer_phone && <p>📞 {shipment.customer_phone}</p>}
                                        {shipment.customer_dni && <p>Doc: {shipment.customer_dni}</p>}
                                    </div>
                                </div>
                            )}
                            <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-2.5">
                                <div className="flex items-center gap-1 mb-1">
                                    <Hash size={10} className="text-gray-400" />
                                    <p className="text-[9px] text-gray-400 uppercase font-bold tracking-wider">Tracking</p>
                                </div>
                                <p className="text-xs font-mono font-semibold text-gray-900 dark:text-white break-all">
                                    {shipment.tracking_number || 'Sin tracking'}
                                </p>
                                <div className="mt-1.5 flex items-center gap-1 flex-wrap">
                                    <StatusBadge status={shipment.status} />
                                    <SubStatusBadge status={shipment.status} carrier={shipment.carrier} detail={shipment.carrier_status_detail || shipment.carrier_status} />
                                    {shipment.is_test && (
                                        <span className="inline-flex items-center px-1 py-0.5 rounded text-[9px] font-bold bg-orange-100 text-orange-700 border border-orange-300 uppercase">TEST</span>
                                    )}
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Fechas en fila compacta */}
                    <div className="grid grid-cols-3 gap-2">
                        {shipment.created_at && (
                            <div className="bg-gray-50 dark:bg-gray-700/50 rounded p-2">
                                <p className="text-[9px] text-gray-400 uppercase font-bold tracking-wider">Creado</p>
                                <p className="text-[11px] font-medium text-gray-900 dark:text-white leading-tight mt-0.5">{formatDate(shipment.created_at)}</p>
                            </div>
                        )}
                        {shipment.shipped_at && (
                            <div className="bg-gray-50 dark:bg-gray-700/50 rounded p-2">
                                <p className="text-[9px] text-gray-400 uppercase font-bold tracking-wider">Enviado</p>
                                <p className="text-[11px] font-medium text-gray-900 dark:text-white leading-tight mt-0.5">{formatDate(shipment.shipped_at)}</p>
                            </div>
                        )}
                        {shipment.delivered_at ? (
                            <div className="bg-emerald-50 dark:bg-emerald-900/20 rounded p-2">
                                <p className="text-[9px] text-emerald-600 uppercase font-bold tracking-wider">Entregado</p>
                                <p className="text-[11px] font-medium text-emerald-700 dark:text-emerald-300 leading-tight mt-0.5">{formatDate(shipment.delivered_at)}</p>
                            </div>
                        ) : shipment.estimated_delivery ? (
                            <div className="bg-blue-50 dark:bg-blue-900/20 rounded p-2">
                                <p className="text-[9px] text-blue-600 uppercase font-bold tracking-wider">Entrega Est.</p>
                                <p className="text-[11px] font-medium text-blue-700 dark:text-blue-300 leading-tight mt-0.5">{formatDate(shipment.estimated_delivery)}</p>
                            </div>
                        ) : null}
                    </div>

                    {/* Botones compactos */}
                    <div className="flex gap-1.5 pt-1">
                        {shipment.guide_url && (
                            <a
                                href={shipment.guide_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex-1 flex items-center justify-center gap-1 py-1.5 px-2 rounded-md bg-blue-600 hover:bg-blue-700 text-white text-[11px] font-semibold transition-colors"
                            >
                                <FileText size={11} />
                                Ver Guía
                            </a>
                        )}
                        {shipment.tracking_url && (
                            <a
                                href={shipment.tracking_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex-1 flex items-center justify-center gap-1 py-1.5 px-2 rounded-md bg-purple-600 hover:bg-purple-700 text-white text-[11px] font-semibold transition-colors"
                            >
                                <Navigation size={11} />
                                Rastrear
                            </a>
                        )}
                        {isCancelled || (shipment.status === 'failed' && !shipment.guide_id) ? null : (
                            <button
                                onClick={() => onCancel(canelId)}
                                disabled={cancelingId === canelId}
                                className="flex items-center justify-center gap-1 py-1.5 px-2 rounded-md bg-red-600 hover:bg-red-700 text-white text-[11px] font-semibold transition-colors disabled:opacity-50"
                                title="Cancelar envío"
                            >
                                {cancelingId === canelId
                                    ? <RefreshCw size={11} className="animate-spin" />
                                    : <><X size={11} /> Cancelar</>
                                }
                            </button>
                        )}
                    </div>
                </div>

                {/* ─── Costos ─────────────────────────────────────────── */}
                {/* Costos + Paquete combinados en una fila */}
                {(shipment.shipping_cost != null || shipment.insurance_cost != null || shipment.total_cost != null || shipment.weight != null || shipment.length != null) && (
                    <div className="px-4 py-3 border-t border-gray-50 dark:border-gray-700">
                        <div className="flex flex-wrap items-center gap-1.5 text-[11px]">
                            {shipment.total_cost != null && (
                                <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300 font-semibold border border-emerald-100 dark:border-emerald-800">
                                    <DollarSign size={11} /> Total {formatMoney(shipment.total_cost)}
                                </span>
                            )}
                            {shipment.shipping_cost != null && (
                                <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-200">
                                    Envío {formatMoney(shipment.shipping_cost)}
                                </span>
                            )}
                            {shipment.insurance_cost != null && (
                                <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-200">
                                    Seguro {formatMoney(shipment.insurance_cost)}
                                </span>
                            )}
                            {shipment.weight != null && (
                                <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-200">
                                    <Box size={11} className="text-gray-400" /> {shipment.weight} kg
                                </span>
                            )}
                            {(shipment.length != null || shipment.width != null || shipment.height != null) && (
                                <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-200">
                                    {shipment.length ?? '?'}×{shipment.width ?? '?'}×{shipment.height ?? '?'} cm
                                </span>
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
                                {tracking.data.history.map((event: EnvioClickTrackHistory & { raw_status?: string; raw_status_detail?: string; carrier?: string }, idx: number) => {
                                    const primary = event.raw_status || STATUS_CONFIG[event.status]?.label || event.status;
                                    const secondary = event.raw_status_detail;
                                    const desc = event.description;
                                    const showDesc = desc && desc !== secondary && desc !== primary;
                                    return (
                                        <div key={idx} className="relative">
                                            <div className={`absolute -left-5 top-0.5 w-3 h-3 rounded-full ring-2 ring-white ${idx === 0 ? 'bg-blue-500' : 'bg-gray-300'}`} />
                                            <div>
                                                <div className="flex items-baseline justify-between gap-2">
                                                    <p className={`text-sm font-semibold ${idx === 0 ? 'text-blue-700' : 'text-gray-800 dark:text-gray-100'}`}>{primary}</p>
                                                    <p className="text-[10px] text-gray-400 dark:text-gray-500 flex-shrink-0">{formatDate(event.date) || event.date}</p>
                                                </div>
                                                {secondary && (
                                                    <p className="text-xs text-gray-700 dark:text-gray-200 mt-0.5 font-medium">{secondary}</p>
                                                )}
                                                {showDesc && <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{desc}</p>}
                                                {event.location && (
                                                    <div className="flex items-center gap-1 mt-1 text-[10px] text-gray-400 dark:text-gray-500">
                                                        <MapPin size={9} />
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
                        <p className="text-xs text-gray-400 dark:text-gray-500 italic">Este envío no tiene número de tracking.</p>
                    ) : (
                        <p className="text-xs text-gray-400 dark:text-gray-500 italic">No hay historial disponible.</p>
                    )}
                </div>

            </div>
        </div>
    );
}

interface OrderNumberFilterProps {
    value: string;
    onChange: (orderNumber: string | undefined) => void;
    businessId?: number;
}

function OrderNumberFilter({ value, onChange, businessId }: OrderNumberFilterProps) {
    const [input, setInput] = useState(value);
    const [suggestions, setSuggestions] = useState<{ order_number: string; customer_name?: string }[]>([]);
    const [open, setOpen] = useState(false);
    const [loading, setLoading] = useState(false);

    useEffect(() => { setInput(value); }, [value]);

    useEffect(() => {
        if (!open) return;
        const term = input.trim();
        if (term.length < 2) { setSuggestions([]); return; }
        const t = setTimeout(async () => {
            setLoading(true);
            try {
                const res = await getOrdersAction({
                    page: 1,
                    page_size: 10,
                    order_number: term,
                    business_id: businessId,
                } as any);
                if (res?.success && Array.isArray(res.data)) {
                    setSuggestions(res.data.map((o: any) => ({ order_number: o.order_number, customer_name: o.customer_name })));
                } else {
                    setSuggestions([]);
                }
            } catch {
                setSuggestions([]);
            } finally {
                setLoading(false);
            }
        }, 250);
        return () => clearTimeout(t);
    }, [input, open, businessId]);

    return (
        <div className="relative min-w-[220px]">
            <Hash size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500" />
            <input
                type="text"
                placeholder="Filtrar por # de orden..."
                className="w-full pl-9 pr-8 py-2 border border-gray-200 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-400 text-sm bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder:text-gray-400 dark:text-gray-500 transition-colors"
                value={input}
                onChange={(e) => { setInput(e.target.value); setOpen(true); }}
                onFocus={() => setOpen(true)}
                onBlur={() => setTimeout(() => setOpen(false), 150)}
                onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                        onChange(input.trim() || undefined);
                        setOpen(false);
                    } else if (e.key === 'Escape') {
                        setOpen(false);
                    }
                }}
            />
            {input && (
                <button
                    type="button"
                    onClick={() => { setInput(''); onChange(undefined); }}
                    className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                    <X size={14} />
                </button>
            )}
            {open && (suggestions.length > 0 || loading || input.trim().length >= 2) && (
                <div className="absolute z-30 mt-1 w-full bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                    {loading && (
                        <div className="px-3 py-2 text-xs text-gray-400 flex items-center gap-2">
                            <RefreshCw size={12} className="animate-spin" /> Buscando...
                        </div>
                    )}
                    {!loading && suggestions.length === 0 && input.trim().length >= 2 && (
                        <div className="px-3 py-2 text-xs text-gray-400">Sin coincidencias</div>
                    )}
                    {suggestions.map((s) => (
                        <button
                            key={s.order_number}
                            type="button"
                            onMouseDown={(e) => e.preventDefault()}
                            onClick={() => {
                                setInput(s.order_number);
                                onChange(s.order_number);
                                setOpen(false);
                            }}
                            className="w-full text-left px-3 py-2 text-sm hover:bg-gray-50 dark:hover:bg-gray-700 flex items-center justify-between gap-2"
                        >
                            <span className="font-semibold text-gray-700 dark:text-gray-200">{s.order_number}</span>
                            {s.customer_name && (
                                <span className="text-xs text-gray-400 truncate">{s.customer_name}</span>
                            )}
                        </button>
                    ))}
                </div>
            )}
        </div>
    );
}

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
        order_number: searchParams.get('order_number') || undefined,
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
                    <OrderNumberFilter
                        value={filters.order_number || ''}
                        onChange={(v) => updateFilters({ order_number: v })}
                        businessId={isSuperAdmin ? (selectedBusinessId ?? undefined) : defaultBusinessId}
                    />
                    {/* Select estado */}
                    <select
                        className="px-3 py-2 border border-gray-200 dark:border-gray-600 dark:border-gray-700 rounded-lg text-sm text-gray-600 dark:text-gray-300 bg-gray-50 dark:bg-gray-700 focus:ring-2 focus:ring-blue-500/20 min-w-[140px] transition-colors"
                        value={filters.status || ''}
                        onChange={(e) => updateFilters({ status: e.target.value || undefined })}
                    >
                        <option value="">Todos los estados</option>
                        <option value="pending">Pendiente</option>
                        <option value="picked_up">Recolectado</option>
                        <option value="in_transit">En Tránsito</option>
                        <option value="out_for_delivery">En Reparto</option>
                        <option value="delivered">Entregado</option>
                        <option value="on_hold">Novedad</option>
                        <option value="returned">Devuelto</option>
                        <option value="cancelled">Cancelado</option>
                    </select>
                    {isSuperAdmin && (
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
                    )}
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
                <div className="w-1/3 flex flex-col min-h-0 bg-white dark:bg-gray-800 dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
                    
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
                                            <SubStatusBadge status={shipment.status} carrier={shipment.carrier} detail={shipment.carrier_status_detail} />
                                            {shipment.is_test && (
                                                <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-bold bg-orange-100 text-orange-700 border border-orange-300 uppercase tracking-widest">TEST</span>
                                            )}
                                        </div>

                                        {/* Row 3: Order number + Tracking + carrier + date */}
                                        <div className="flex items-center gap-2 flex-wrap">
                                            {shipment.order_number && (
                                                <span className="text-[10px] font-semibold text-purple-700 dark:text-purple-300 bg-purple-50 dark:bg-purple-900/30 px-1.5 py-0.5 rounded border border-purple-100 dark:border-purple-800">
                                                    {shipment.order_number}
                                                </span>
                                            )}
                                            {shipment.tracking_number && (
                                                <span className="text-[10px] font-mono text-gray-400 dark:text-gray-500 bg-gray-100 px-1.5 py-0.5 rounded">
                                                    #{shipment.tracking_number.slice(-10)}
                                                </span>
                                            )}
                                            {shipment.carrier && (
                                                <span className="text-[10px] text-gray-500 dark:text-gray-400 flex items-center gap-1 bg-white dark:bg-gray-700 px-1.5 py-0.5 rounded border border-gray-200 dark:border-gray-600">
                                                    {getCarrierLogo(shipment.carrier) ? (
                                                        <img src={getCarrierLogo(shipment.carrier)!} alt={shipment.carrier} className="h-3 w-auto object-contain" />
                                                    ) : (
                                                        <Truck size={9} />
                                                    )}
                                                    {shipment.carrier}
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
                <div className="w-2/3 min-h-0 bg-white dark:bg-gray-800 dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
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
