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

interface StatusCount {
  status: string;
  label: string;
  count: number;
  color: string;
  bgColor: string;
  icon: React.ReactNode;
}

const countByStatus = (shipments: Shipment[]): StatusCount[] => {
  const statusMap: Record<string, StatusCount> = {
    pending: { status: 'pending', label: 'Pendiente', count: 0, color: 'text-amber-700', bgColor: 'bg-amber-50 border-amber-200', icon: <Clock size={20} /> },
    picked_up: { status: 'picked_up', label: 'Recolectado', count: 0, color: 'text-indigo-700', bgColor: 'bg-indigo-50 border-indigo-200', icon: <PackageCheck size={20} /> },
    in_transit: { status: 'in_transit', label: 'En tránsito', count: 0, color: 'text-blue-700', bgColor: 'bg-blue-50 border-blue-200', icon: <Truck size={20} /> },
    out_for_delivery: { status: 'out_for_delivery', label: 'En reparto', count: 0, color: 'text-purple-700', bgColor: 'bg-purple-50 border-purple-200', icon: <MapPinned size={20} /> },
    delivered: { status: 'delivered', label: 'Entregado', count: 0, color: 'text-emerald-700', bgColor: 'bg-emerald-50 border-emerald-200', icon: <CheckCircle2 size={20} /> },
    on_hold: { status: 'on_hold', label: 'Novedad', count: 0, color: 'text-orange-700', bgColor: 'bg-orange-50 border-orange-200', icon: <PauseCircle size={20} /> },
    returned: { status: 'returned', label: 'Devuelto', count: 0, color: 'text-rose-700', bgColor: 'bg-rose-50 border-rose-200', icon: <RotateCcw size={20} /> },
    cancelled: { status: 'cancelled', label: 'Cancelado', count: 0, color: 'text-red-700', bgColor: 'bg-red-50 border-red-200', icon: <X size={20} /> },
  };

  shipments.forEach((shipment) => {
    if (shipment.status && statusMap[shipment.status]) {
      statusMap[shipment.status].count++;
    }
  });

  return Object.values(statusMap).filter(s => s.count > 0).sort((a, b) => b.count - a.count);
};

export function ShipmentStatusWidget({ selectedBusinessId }: ShipmentStatusWidgetProps) {
  const [statusCounts, setStatusCounts] = useState<StatusCount[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [totalShipments, setTotalShipments] = useState(0);

  useEffect(() => {
    const loadShipments = async () => {
      try {
        setLoading(true);
        const result = await getShipmentsAction({
          business_id: selectedBusinessId,
          page: 1,
          page_size: 10000,
        });

        if (!result || typeof result !== 'object' || !('data' in result)) {
          throw new Error('Invalid response format');
        }

        const shipments = (result as any).data || [];
        const counts = countByStatus(shipments);
        setStatusCounts(counts);
        setTotalShipments(shipments.length);
        setError(null);
      } catch (err: any) {
        setError(err.message);
        setStatusCounts([]);
      } finally {
        setLoading(false);
      }
    };

    if (selectedBusinessId) {
      loadShipments();
    }
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

  return (
    <div>
      {/* Total Summary */}
      <div className="mb-6 p-4 bg-gradient-to-r from-slate-50 to-slate-100 rounded-lg border border-slate-200">
        <p className="text-sm text-slate-600 mb-1">Total de Envíos</p>
        <p className="text-3xl font-bold text-slate-900">{totalShipments.toLocaleString()}</p>
      </div>

      {/* Status Cards Grid */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {statusCounts.map((status) => (
          <Link
            key={status.status}
            href={`/shipments?status=${status.status}`}
            className={`p-4 rounded-lg border cursor-pointer transition-all hover:shadow-md hover:scale-105 ${status.bgColor}`}
          >
            <div className="flex items-start justify-between mb-3">
              <div className={status.color}>
                {status.icon}
              </div>
            </div>
            <p className="text-2xl font-bold text-gray-900 mb-1">{status.count.toLocaleString()}</p>
            <p className="text-sm text-gray-600">{status.label}</p>
          </Link>
        ))}
      </div>

      {statusCounts.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          <p>No hay envíos registrados</p>
        </div>
      )}
    </div>
  );
}
