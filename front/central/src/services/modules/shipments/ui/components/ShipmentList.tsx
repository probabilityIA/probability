'use client';

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Badge, Button } from '@/shared/ui';
import { useHasPermission } from '@/shared/contexts/permissions-context';
import { getShipmentsAction, trackShipmentAction, cancelShipmentAction } from '../../infra/actions';
import { GetShipmentsParams, Shipment, EnvioClickTrackHistory } from '../../domain/types';
import { Search, Package, Truck, Calendar, MapPin, X, Eye, CreditCard, ExternalLink, RefreshCw, AlertTriangle, Plus } from 'lucide-react';
import { ManualShipmentModal } from './ManualShipmentModal';

export default function ShipmentList() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const canCreate = useHasPermission('Envios', 'Create');
    const canDelete = useHasPermission('Envios', 'Delete');
    const [loading, setLoading] = useState(true);
    const [shipments, setShipments] = useState<Shipment[]>([]);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);

    const [filters, setFilters] = useState<GetShipmentsParams>({
        page: Number(searchParams.get('page')) || 1,
        page_size: Number(searchParams.get('page_size')) || 20,
        tracking_number: searchParams.get('tracking_number') || undefined,
        order_id: searchParams.get('order_id') || undefined,
        carrier: searchParams.get('carrier') || undefined,
        status: searchParams.get('status') || undefined,
    });

    // Tracking State
    const [trackingModal, setTrackingModal] = useState<{ open: boolean; loading: boolean; data?: any; error?: string }>({
        open: false,
        loading: false
    });

    const [cancelingId, setCancelingId] = useState<string | null>(null);
    const [isManualModalOpen, setIsManualModalOpen] = useState(false);

    const fetchShipments = async () => {
        setLoading(true);
        try {
            const response = await getShipmentsAction(filters);
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
    }, [filters]);

    const updateFilters = (newFilters: Partial<GetShipmentsParams>) => {
        const updated = { ...filters, ...newFilters };
        if (!newFilters.page && newFilters.page !== 0) {
            updated.page = 1;
        }
        setFilters(updated);

        const params = new URLSearchParams();
        Object.entries(updated).forEach(([key, value]) => {
            if (value) params.set(key, String(value));
        });
        router.push(`?${params.toString()}`);
    };

    const handleTrack = async (trackingNumber: string) => {
        setTrackingModal({ open: true, loading: true });
        try {
            const response = await trackShipmentAction(trackingNumber);
            if ('data' in response && response.success) {
                setTrackingModal({ open: true, loading: false, data: response.data });
            } else {
                setTrackingModal({ open: true, loading: false, error: response.message });
            }
        } catch (error: any) {
            setTrackingModal({ open: true, loading: false, error: error.message });
        }
    };

    const handleCancel = async (id: string) => {
        if (!confirm('¿Estás seguro de que deseas cancelar este envío en EnvioClick?')) return;

        setCancelingId(id);
        try {
            const response = await cancelShipmentAction(id);
            if (response.success) {
                alert('Envío cancelado exitosamente');
                fetchShipments(); // Refresh list
            } else {
                alert(`Error: ${response.message}`);
            }
        } catch (error: any) {
            alert(`Error: ${error.message}`);
        } finally {
            setCancelingId(null);
        }
    };

    return (
        <div className="space-y-4">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-bold text-gray-800">Envíos</h2>
                {canCreate && (
                    <Button onClick={() => setIsManualModalOpen(true)}>
                        <Plus size={16} className="mr-2" />
                        Agregar Envío
                    </Button>
                )}
            </div>

            {/* Filters */}
            <div className="bg-white p-4 sm:p-6 rounded-lg shadow-sm border border-gray-200">
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
                    <div className="relative flex-1">
                        <input
                            type="text"
                            placeholder="Buscar por tracking..."
                            className="w-full px-3 py-2 pr-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder:text-gray-500 bg-white"
                            value={filters.tracking_number || ''}
                            onChange={(e) => updateFilters({ tracking_number: e.target.value || undefined })}
                        />
                        <button
                            onClick={() => filters.tracking_number && handleTrack(filters.tracking_number)}
                            className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 text-blue-600 hover:bg-blue-50 rounded-md transition-colors disabled:opacity-50 disabled:text-gray-400"
                            disabled={!filters.tracking_number}
                            title="Consultar en EnvioClick"
                        >
                            <Search size={18} />
                        </button>
                    </div>
                    <input
                        type="text"
                        placeholder="Buscar por ID de orden..."
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder:text-gray-500 bg-white"
                        value={filters.order_id || ''}
                        onChange={(e) => updateFilters({ order_id: e.target.value || undefined })}
                    />
                    <input
                        type="text"
                        placeholder="Buscar por transportista..."
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder:text-gray-500 bg-white"
                        value={filters.carrier || ''}
                        onChange={(e) => updateFilters({ carrier: e.target.value || undefined })}
                    />
                    <select
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                        value={filters.status || ''}
                        onChange={(e) => updateFilters({ status: e.target.value || undefined })}
                    >
                        <option value="">Todos los estados</option>
                        <option value="pending">Pendiente</option>
                        <option value="in_transit">En Tránsito</option>
                        <option value="delivered">Entregado</option>
                        <option value="failed">Fallido</option>
                    </select>
                </div>
            </div>

            {/* Table */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Tracking
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Orden
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden sm:table-cell">
                                    Transportista
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Estado
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden md:table-cell">
                                    Enviado
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden lg:table-cell">
                                    Entrega Est.
                                </th>
                                <th className="px-3 sm:px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Última Milla
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {loading ? (
                                <tr>
                                    <td colSpan={7} className="px-4 sm:px-6 py-8 text-center text-gray-500">
                                        Cargando envíos...
                                    </td>
                                </tr>
                            ) : shipments.length === 0 ? (
                                <tr>
                                    <td colSpan={7} className="px-4 sm:px-6 py-8 text-center text-gray-500">
                                        No hay envíos disponibles
                                    </td>
                                </tr>
                            ) : (
                                shipments.map((shipment) => (
                                    <tr key={shipment.id} className="hover:bg-gray-50">
                                        <td className="px-3 sm:px-6 py-4">
                                            <div className="font-medium text-gray-900 text-sm">
                                                {shipment.tracking_number || 'Sin tracking'}
                                            </div>
                                            {shipment.tracking_url && (
                                                <a
                                                    href={shipment.tracking_url}
                                                    target="_blank"
                                                    rel="noopener noreferrer"
                                                    className="text-xs text-blue-600 hover:underline flex items-center gap-1 mt-1"
                                                >
                                                    <ExternalLink size={12} /> External
                                                </a>
                                            )}
                                        </td>
                                        <td className="px-3 sm:px-6 py-4">
                                            {shipment.order_id ? (
                                                <span className="font-mono text-sm text-gray-900 bg-gray-100 px-2 py-1 rounded">
                                                    {shipment.order_id.substring(0, 8)}...
                                                </span>
                                            ) : (
                                                <div className="flex flex-col">
                                                    <span className="text-sm font-medium text-gray-900">{shipment.client_name || 'Desconocido'}</span>
                                                    <span className="text-xs text-gray-500">{shipment.destination_address || '-'}</span>
                                                </div>
                                            )}
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 hidden sm:table-cell">
                                            <div className="flex flex-col">
                                                <div className="flex items-center gap-2">
                                                    <Truck size={14} className="text-gray-400" />
                                                    <span className="text-sm font-medium text-gray-700">{shipment.carrier || 'Manual'}</span>
                                                </div>
                                                {shipment.carrier_code && (
                                                    <span className="text-[10px] text-gray-500 ml-5 uppercase font-bold tracking-wider">{shipment.carrier_code}</span>
                                                )}
                                            </div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap">
                                            <Badge type={
                                                shipment.status === 'delivered' ? 'success' :
                                                    shipment.status === 'in_transit' ? 'primary' :
                                                        shipment.status === 'failed' ? 'error' : 'warning'
                                            }>
                                                {shipment.status}
                                            </Badge>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500 hidden md:table-cell">
                                            <div className="flex items-center gap-2">
                                                <Calendar size={14} className="text-gray-400" />
                                                {shipment.shipped_at ? new Date(shipment.shipped_at).toLocaleDateString() : '-'}
                                            </div>
                                        </td>
                                        <td className="px-3 sm:px-6 py-4 whitespace-nowrap text-right">
                                            <div className="flex items-center justify-end gap-2">
                                                {canDelete && (
                                                    <Button
                                                        size="sm"
                                                        variant="outline"
                                                        className="h-8 w-8 p-0 text-red-600 hover:text-red-700 hover:bg-red-50 border-red-200"
                                                        onClick={() => handleCancel(shipment.tracking_number || shipment.id.toString())}
                                                        disabled={cancelingId === (shipment.tracking_number || shipment.id.toString())}
                                                        title="Cancelar envío"
                                                    >
                                                        {cancelingId === (shipment.tracking_number || shipment.id.toString()) ? (
                                                            <RefreshCw size={14} className="animate-spin" />
                                                        ) : (
                                                            <X size={14} />
                                                        )}
                                                    </Button>
                                                )}
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {/* Tracking Modal */}
                {trackingModal.open && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
                        <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg overflow-hidden flex flex-col max-h-[90vh]">
                            <div className="p-6 border-b border-gray-100 flex items-center justify-between bg-white sticky top-0">
                                <div className="flex items-center gap-3">
                                    <div className="p-2 bg-blue-50 rounded-lg text-blue-600">
                                        <Truck size={24} />
                                    </div>
                                    <div>
                                        <h3 className="text-lg font-bold text-gray-900">Estado del Envío</h3>
                                        <p className="text-sm text-gray-500">Rastreo oficial de EnvioClick</p>
                                    </div>
                                </div>
                                <button
                                    onClick={() => setTrackingModal({ open: false, loading: false })}
                                    className="p-2 hover:bg-gray-100 rounded-full transition-colors"
                                >
                                    <X size={20} className="text-gray-400" />
                                </button>
                            </div>

                            <div className="flex-1 overflow-y-auto p-6">
                                {trackingModal.loading ? (
                                    <div className="flex flex-col items-center justify-center py-12 gap-4">
                                        <RefreshCw size={40} className="text-blue-500 animate-spin" />
                                        <p className="text-gray-500 font-medium">Consultando API de EnvioClick...</p>
                                    </div>
                                ) : trackingModal.error ? (
                                    <div className="flex flex-col items-center justify-center py-8 gap-4 text-center">
                                        <div className="p-4 bg-red-50 rounded-full text-red-500">
                                            <AlertTriangle size={48} />
                                        </div>
                                        <div>
                                            <p className="text-gray-900 font-bold">No pudimos obtener el rastreo</p>
                                            <p className="text-red-600 text-sm mt-1">{trackingModal.error}</p>
                                        </div>
                                        <Button variant="outline" onClick={() => setTrackingModal({ open: false, loading: false })}>
                                            Cerrar
                                        </Button>
                                    </div>
                                ) : (
                                    <div className="space-y-8">
                                        <div className="grid grid-cols-2 gap-4">
                                            <div className="bg-gray-50 p-4 rounded-xl">
                                                <p className="text-[10px] text-gray-500 uppercase font-bold tracking-wider mb-1">Carrier</p>
                                                <p className="text-gray-900 font-semibold">{trackingModal.data?.carrier}</p>
                                            </div>
                                            <div className="bg-gray-50 p-4 rounded-xl">
                                                <p className="text-[10px] text-gray-500 uppercase font-bold tracking-wider mb-1">Estado Actual</p>
                                                <p className="text-blue-600 font-bold">{trackingModal.data?.status}</p>
                                            </div>
                                        </div>

                                        <div className="relative pl-8 space-y-8 before:absolute before:left-3 before:top-2 before:bottom-2 before:w-0.5 before:bg-gray-100">
                                            {trackingModal.data?.history?.map((event: EnvioClickTrackHistory, idx: number) => (
                                                <div key={idx} className="relative">
                                                    <div className={`absolute -left-8 p-1.5 rounded-full ring-4 ring-white z-10 ${idx === 0 ? 'bg-blue-500' : 'bg-gray-300'}`}>
                                                        <div className="w-2 h-2" />
                                                    </div>
                                                    <div>
                                                        <div className="flex items-center justify-between mb-1">
                                                            <p className={`text-sm font-bold ${idx === 0 ? 'text-blue-600' : 'text-gray-900'}`}>{event.status}</p>
                                                            <p className="text-xs text-gray-400">{event.date}</p>
                                                        </div>
                                                        <p className="text-sm text-gray-600 leading-relaxed">{event.description}</p>
                                                        {event.location && (
                                                            <div className="flex items-center gap-1 mt-2 text-xs text-gray-400">
                                                                <MapPin size={12} />
                                                                <span>{event.location}</span>
                                                            </div>
                                                        )}
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                )}

                {/* Pagination */}
                {(totalPages > 1 || total > 0) && (
                    <div className="bg-white px-3 sm:px-4 lg:px-6 py-3 flex flex-col sm:flex-row items-center justify-between gap-3 border-t border-gray-200">
                        {/* Mobile: Simple pagination */}
                        <div className="flex-1 flex justify-between sm:hidden w-full">
                            <Button
                                variant="outline"
                                onClick={() => updateFilters({ page: page - 1 })}
                                disabled={page === 1}
                                size="sm"
                            >
                                Anterior
                            </Button>
                            <Button
                                variant="outline"
                                onClick={() => updateFilters({ page: page + 1 })}
                                disabled={page === totalPages}
                                size="sm"
                            >
                                Siguiente
                            </Button>
                        </div>

                        {/* Desktop: Full pagination */}
                        <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between w-full">
                            <div className="flex items-center gap-3">
                                <p className="text-xs sm:text-sm text-gray-700">
                                    Mostrando <span className="font-medium">{(page - 1) * (filters.page_size || 20) + 1}</span> a{' '}
                                    <span className="font-medium">{Math.min(page * (filters.page_size || 20), total)}</span> de{' '}
                                    <span className="font-medium">{total}</span> resultados
                                </p>
                                <div className="flex items-center gap-2">
                                    <label className="text-xs sm:text-sm text-gray-700 whitespace-nowrap">
                                        Mostrar:
                                    </label>
                                    <select
                                        value={filters.page_size || 20}
                                        onChange={(e) => {
                                            const newPageSize = parseInt(e.target.value);
                                            updateFilters({ page_size: newPageSize, page: 1 });
                                        }}
                                        className="px-2 py-1.5 text-xs sm:text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                                    >
                                        <option value="10">10</option>
                                        <option value="20">20</option>
                                        <option value="50">50</option>
                                        <option value="100">100</option>
                                    </select>
                                </div>
                            </div>
                            <div className="flex items-center gap-2">
                                <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                                    <button
                                        onClick={() => updateFilters({ page: page - 1 })}
                                        disabled={page === 1}
                                        className="relative inline-flex items-center px-2 sm:px-3 py-2 rounded-l-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Anterior
                                    </button>
                                    <span className="relative inline-flex items-center px-3 sm:px-4 py-2 border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-700">
                                        Página {page} de {totalPages}
                                    </span>
                                    <button
                                        onClick={() => updateFilters({ page: page + 1 })}
                                        disabled={page === totalPages}
                                        className="relative inline-flex items-center px-2 sm:px-3 py-2 rounded-r-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Siguiente
                                    </button>
                                </nav>
                            </div>
                        </div>

                        {/* Mobile: Page size selector */}
                        <div className="flex items-center justify-between w-full sm:hidden pt-2 border-t border-gray-200">
                            <div className="flex items-center gap-2">
                                <label className="text-xs text-gray-700 whitespace-nowrap">
                                    Mostrar:
                                </label>
                                <select
                                    value={filters.page_size || 20}
                                    onChange={(e) => {
                                        const newPageSize = parseInt(e.target.value);
                                        updateFilters({ page_size: newPageSize, page: 1 });
                                    }}
                                    className="px-2 py-1.5 text-xs border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                                >
                                    <option value="10">10</option>
                                    <option value="20">20</option>
                                    <option value="50">50</option>
                                    <option value="100">100</option>
                                </select>
                            </div>
                            <p className="text-xs text-gray-500">
                                Página {page} de {totalPages}
                            </p>
                        </div>
                    </div>
                )}
            </div>

            <ManualShipmentModal
                isOpen={isManualModalOpen}
                onClose={() => setIsManualModalOpen(false)}
                onSuccess={fetchShipments}
            />
        </div>
    );
}
