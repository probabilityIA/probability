'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { getShipmentsAction } from '@/services/modules/shipments/infra/actions';
import { Shipment } from '@/services/modules/shipments/domain/types';
import { Spinner } from '@/shared/ui';
import { Clock, PackageCheck, Truck, MapPinned, CheckCircle2, PauseCircle, RotateCcw, X } from 'lucide-react';

interface ShipmentStatusWidgetProps {
  selectedBusinessId?: number;
}

interface StatusConfig {
  label: string;
  color: string;
  bgColor: string;
  iconType: string;
}

const STATUS_CONFIG: Record<string, StatusConfig> = {
  pending: { label: 'Pendiente', color: 'text-amber-700', bgColor: 'bg-amber-50 border-amber-200', iconType: 'clock' },
  picked_up: { label: 'Recolectado', color: 'text-indigo-700', bgColor: 'bg-indigo-50 border-indigo-200', iconType: 'package' },
  in_transit: { label: 'En tránsito', color: 'text-blue-700', bgColor: 'bg-blue-50 border-blue-200', iconType: 'truck' },
  out_for_delivery: { label: 'En reparto', color: 'text-purple-700', bgColor: 'bg-purple-50 border-purple-200', iconType: 'mappin' },
  delivered: { label: 'Entregado', color: 'text-emerald-700', bgColor: 'bg-emerald-50 border-emerald-200', iconType: 'check' },
  on_hold: { label: 'Novedad', color: 'text-orange-700', bgColor: 'bg-orange-50 border-orange-200', iconType: 'pause' },
  returned: { label: 'Devuelto', color: 'text-rose-700', bgColor: 'bg-rose-50 border-rose-200', iconType: 'rotate' },
  cancelled: { label: 'Cancelado', color: 'text-red-700', bgColor: 'bg-red-50 border-red-200', iconType: 'x' },
};

function getIcon(iconType: string) {
  switch (iconType) {
    case 'clock': return <Clock size={20} />;
    case 'package': return <PackageCheck size={20} />;
    case 'truck': return <Truck size={20} />;
    case 'mappin': return <MapPinned size={20} />;
    case 'check': return <CheckCircle2 size={20} />;
    case 'pause': return <PauseCircle size={20} />;
    case 'rotate': return <RotateCcw size={20} />;
    case 'x': return <X size={20} />;
    default: return null;
  }
}

export function ShipmentStatusWidget({ selectedBusinessId }: ShipmentStatusWidgetProps) {
  console.log('ShipmentStatusWidget: Mounted');
  const [counts, setCounts] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [totalShipments, setTotalShipments] = useState(0);

  useEffect(() => {
    const loadShipments = async () => {
      try {
        console.log('ShipmentStatusWidget: Loading shipments, selectedBusinessId:', selectedBusinessId);
        setLoading(true);
        const params: any = {
          page: 1,
          page_size: 10000,
        };
        if (selectedBusinessId) {
          params.business_id = selectedBusinessId;
        }

        console.log('ShipmentStatusWidget: Calling getShipmentsAction with params:', params);
        const result = await getShipmentsAction(params);
        console.log('ShipmentStatusWidget: Got result:', result);

        if (!result || typeof result !== 'object' || !('data' in result)) {
          throw new Error('Invalid response format');
        }

        const shipments = (result as any).data || [];
        const statusCounts: Record<string, number> = {};

        shipments.forEach((shipment: Shipment) => {
          if (shipment.status) {
            statusCounts[shipment.status] = (statusCounts[shipment.status] || 0) + 1;
          }
        });

        setCounts(statusCounts);
        setTotalShipments(shipments.length);
        setError(null);
      } catch (err: any) {
        setError(err.message);
        setCounts({});
      } finally {
        setLoading(false);
      }
    };

    loadShipments();
  }, [selectedBusinessId]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Spinner />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8 text-red-600">
        <p>Error cargando envíos: {error}</p>
      </div>
    );
  }

  const statusEntries = Object.entries(counts)
    .filter(([_, count]) => count > 0)
    .map(([status, count]) => ({ status, count, config: STATUS_CONFIG[status] }))
    .filter(item => item.config)
    .sort((a, b) => b.count - a.count);

  return (
    <div>
      {/* Total Summary */}
      <div className="mb-6 p-4 bg-gradient-to-r from-slate-50 to-slate-100 rounded-lg border border-slate-200">
        <p className="text-sm text-slate-600 mb-1">Total de Envíos</p>
        <p className="text-3xl font-bold text-slate-900">{totalShipments.toLocaleString()}</p>
      </div>

      {/* Status Cards Grid */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {statusEntries.map(({ status, count, config }) => (
          <Link
            key={status}
            href={`/shipments?status=${status}`}
            className={`p-4 rounded-lg border cursor-pointer transition-all hover:shadow-md hover:scale-105 ${config.bgColor}`}
          >
            <div className={`flex items-start justify-between mb-3 ${config.color}`}>
              {getIcon(config.iconType)}
            </div>
            <p className="text-2xl font-bold text-gray-900 mb-1">{count.toLocaleString()}</p>
            <p className="text-sm text-gray-600">{config.label}</p>
          </Link>
        ))}
      </div>

      {statusEntries.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          <p>No hay envíos registrados</p>
        </div>
      )}
    </div>
  );
}
