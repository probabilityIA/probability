'use client';

import { useState, useEffect } from 'react';
import { CheckIcon } from '@heroicons/react/24/outline';
import {
    RouteInfo,
    CreateRouteDTO,
    UpdateRouteDTO,
    CreateRouteStopDTO,
    DriverOption,
    VehicleOption,
    AssignableOrder,
} from '../../domain/types';
import {
    createRouteAction,
    updateRouteAction,
    getAvailableDriversAction,
    getAvailableVehiclesAction,
    getAssignableOrdersAction,
} from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface RouteFormProps {
    route?: RouteInfo;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

export default function RouteForm({ route, onSuccess, onCancel, businessId }: RouteFormProps) {
    const isEditing = !!route;

    const [date, setDate] = useState(route?.date ? route.date.substring(0, 10) : '');
    const [driverId, setDriverId] = useState(route?.driver_id?.toString() || '');
    const [vehicleId, setVehicleId] = useState(route?.vehicle_id?.toString() || '');
    const [originAddress, setOriginAddress] = useState(route?.origin_address || '');
    const [notes, setNotes] = useState(route?.notes || '');

    // Form options
    const [drivers, setDrivers] = useState<DriverOption[]>([]);
    const [vehicles, setVehicles] = useState<VehicleOption[]>([]);
    const [assignableOrders, setAssignableOrders] = useState<AssignableOrder[]>([]);
    const [selectedOrderIds, setSelectedOrderIds] = useState<Set<string>>(new Set());
    const [loadingOptions, setLoadingOptions] = useState(false);

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    // Load form options on mount
    useEffect(() => {
        const loadOptions = async () => {
            setLoadingOptions(true);
            try {
                const [driversList, vehiclesList, ordersList] = await Promise.all([
                    getAvailableDriversAction(businessId),
                    getAvailableVehiclesAction(businessId),
                    isEditing ? Promise.resolve([]) : getAssignableOrdersAction(businessId),
                ]);
                setDrivers(driversList);
                setVehicles(vehiclesList);
                setAssignableOrders(ordersList);
            } catch {
                // Silently fail — selectors will just be empty
            } finally {
                setLoadingOptions(false);
            }
        };
        loadOptions();
    }, [businessId, isEditing]);

    const toggleOrder = (orderId: string) => {
        setSelectedOrderIds((prev) => {
            const next = new Set(prev);
            if (next.has(orderId)) {
                next.delete(orderId);
            } else {
                next.add(orderId);
            }
            return next;
        });
    };

    const selectAllOrders = () => {
        if (selectedOrderIds.size === assignableOrders.length) {
            setSelectedOrderIds(new Set());
        } else {
            setSelectedOrderIds(new Set(assignableOrders.map((o) => o.id)));
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            if (isEditing) {
                const updateData: UpdateRouteDTO = {
                    date: date || undefined,
                    driver_id: driverId ? Number(driverId) : undefined,
                    vehicle_id: vehicleId ? Number(vehicleId) : undefined,
                    origin_address: originAddress || undefined,
                    notes: notes || undefined,
                };
                await updateRouteAction(route.id, updateData, businessId);
            } else {
                // Build stops from selected orders
                const validStops: CreateRouteStopDTO[] = assignableOrders
                    .filter((o) => selectedOrderIds.has(o.id))
                    .map((o) => ({
                        order_id: o.id,
                        address: o.address || 'Sin direccion',
                        city: o.city || undefined,
                        lat: o.lat ?? undefined,
                        lng: o.lng ?? undefined,
                        customer_name: o.customer_name || 'Sin nombre',
                        customer_phone: o.customer_phone || undefined,
                    }));

                const createData: CreateRouteDTO = {
                    date,
                    driver_id: driverId ? Number(driverId) : undefined,
                    vehicle_id: vehicleId ? Number(vehicleId) : undefined,
                    origin_address: originAddress || undefined,
                    notes: notes || undefined,
                    stops: validStops.length > 0 ? validStops : undefined,
                };
                await createRouteAction(createData, businessId);
            }

            setSuccess(isEditing ? 'Ruta actualizada exitosamente' : 'Ruta creada exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar la ruta'));
        } finally {
            setLoading(false);
        }
    };

    const formatCurrency = (amount: number) =>
        new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(amount);

    return (
        <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}
            {success && (
                <Alert type="success" onClose={() => setSuccess(null)}>
                    {success}
                </Alert>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Fecha */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Fecha <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="date"
                        value={date}
                        onChange={(e) => setDate(e.target.value)}
                        required
                    />
                </div>

                {/* Conductor */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Conductor
                    </label>
                    <select
                        value={driverId}
                        onChange={(e) => setDriverId(e.target.value)}
                        disabled={loadingOptions}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                        <option value="">— Sin conductor —</option>
                        {drivers.map((d) => (
                            <option key={d.id} value={d.id}>
                                {d.first_name} {d.last_name} — {d.identification} {d.license_type ? `(${d.license_type})` : ''}
                            </option>
                        ))}
                    </select>
                </div>

                {/* Vehiculo */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Vehiculo
                    </label>
                    <select
                        value={vehicleId}
                        onChange={(e) => setVehicleId(e.target.value)}
                        disabled={loadingOptions}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                        <option value="">— Sin vehiculo —</option>
                        {vehicles.map((v) => (
                            <option key={v.id} value={v.id}>
                                {v.license_plate} — {v.brand} {v.vehicle_model} ({v.type})
                            </option>
                        ))}
                    </select>
                </div>

                {/* Direccion de origen */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Direccion de origen
                    </label>
                    <Input
                        type="text"
                        value={originAddress}
                        onChange={(e) => setOriginAddress(e.target.value)}
                        placeholder="Ej: Calle 100 #15-20"
                    />
                </div>

                {/* Notas */}
                <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Notas
                    </label>
                    <textarea
                        value={notes}
                        onChange={(e) => setNotes(e.target.value)}
                        placeholder="Notas sobre la ruta..."
                        rows={2}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                    />
                </div>
            </div>

            {/* Order selector - only for creation */}
            {!isEditing && (
                <div className="space-y-3">
                    <div className="flex items-center justify-between">
                        <h3 className="text-sm font-semibold text-gray-700">
                            Ordenes en procesamiento
                            {assignableOrders.length > 0 && (
                                <span className="ml-2 text-xs font-normal text-gray-400">
                                    ({selectedOrderIds.size} de {assignableOrders.length} seleccionadas)
                                </span>
                            )}
                        </h3>
                        {assignableOrders.length > 0 && (
                            <button
                                type="button"
                                onClick={selectAllOrders}
                                className="text-xs text-blue-600 hover:text-blue-700 font-medium"
                            >
                                {selectedOrderIds.size === assignableOrders.length ? 'Deseleccionar todas' : 'Seleccionar todas'}
                            </button>
                        )}
                    </div>

                    {loadingOptions ? (
                        <div className="text-sm text-gray-400 text-center py-6">
                            Cargando ordenes disponibles...
                        </div>
                    ) : assignableOrders.length === 0 ? (
                        <div className="text-sm text-gray-400 text-center py-6 border border-dashed border-gray-200 rounded-lg">
                            No hay ordenes en procesamiento disponibles para asignar
                        </div>
                    ) : (
                        <div className="max-h-64 overflow-y-auto border border-gray-200 rounded-lg divide-y divide-gray-100">
                            {assignableOrders.map((order) => {
                                const isSelected = selectedOrderIds.has(order.id);
                                return (
                                    <button
                                        key={order.id}
                                        type="button"
                                        onClick={() => toggleOrder(order.id)}
                                        className={`w-full flex items-center gap-3 px-3 py-2.5 text-left transition-colors ${
                                            isSelected
                                                ? 'bg-blue-50 hover:bg-blue-100'
                                                : 'bg-white hover:bg-gray-50'
                                        }`}
                                    >
                                        {/* Checkbox */}
                                        <div
                                            className={`flex-shrink-0 w-5 h-5 rounded border-2 flex items-center justify-center transition-colors ${
                                                isSelected
                                                    ? 'bg-blue-600 border-blue-600'
                                                    : 'border-gray-300'
                                            }`}
                                        >
                                            {isSelected && <CheckIcon className="w-3.5 h-3.5 text-white" />}
                                        </div>

                                        {/* Order info */}
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-center gap-2">
                                                <span className="text-sm font-medium text-gray-900">
                                                    #{order.order_number}
                                                </span>
                                                <span className="text-xs text-gray-400">
                                                    {order.item_count} {order.item_count === 1 ? 'item' : 'items'}
                                                </span>
                                                <span className="text-xs font-medium text-green-600">
                                                    {formatCurrency(order.total_amount)}
                                                </span>
                                            </div>
                                            <div className="text-xs text-gray-500 truncate">
                                                {order.customer_name}
                                                {order.address ? ` — ${order.address}` : ''}
                                                {order.city ? `, ${order.city}` : ''}
                                            </div>
                                        </div>
                                    </button>
                                );
                            })}
                        </div>
                    )}
                </div>
            )}

            <div className="flex justify-end gap-3 pt-4 border-t">
                <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
                    Cancelar
                </Button>
                <Button type="submit" variant="primary" disabled={loading}>
                    {loading ? 'Guardando...' : isEditing ? 'Actualizar' : 'Crear ruta'}
                </Button>
            </div>
        </form>
    );
}
