'use client';

import { useState } from 'react';
import { PlusIcon, TrashIcon } from '@heroicons/react/24/outline';
import { RouteInfo, CreateRouteDTO, UpdateRouteDTO, CreateRouteStopDTO } from '../../domain/types';
import { createRouteAction, updateRouteAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';

interface RouteFormProps {
    route?: RouteInfo;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

interface StopEntry {
    order_id: string;
    address: string;
    customer_name: string;
    customer_phone: string;
}

const emptyStop = (): StopEntry => ({
    order_id: '',
    address: '',
    customer_name: '',
    customer_phone: '',
});

export default function RouteForm({ route, onSuccess, onCancel, businessId }: RouteFormProps) {
    const isEditing = !!route;

    const [date, setDate] = useState(route?.date ? route.date.substring(0, 10) : '');
    const [driverId, setDriverId] = useState(route?.driver_id?.toString() || '');
    const [vehicleId, setVehicleId] = useState(route?.vehicle_id?.toString() || '');
    const [originAddress, setOriginAddress] = useState(route?.origin_address || '');
    const [notes, setNotes] = useState(route?.notes || '');

    const [stops, setStops] = useState<StopEntry[]>([emptyStop()]);

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const addStop = () => {
        setStops((prev) => [...prev, emptyStop()]);
    };

    const removeStop = (index: number) => {
        setStops((prev) => prev.filter((_, i) => i !== index));
    };

    const updateStopField = (index: number, field: keyof StopEntry, value: string) => {
        setStops((prev) => prev.map((s, i) => (i === index ? { ...s, [field]: value } : s)));
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
                // Build stops array from entries that have at minimum an address and customer name
                const validStops: CreateRouteStopDTO[] = stops
                    .filter((s) => s.address.trim() && s.customer_name.trim())
                    .map((s) => ({
                        order_id: s.order_id.trim() || undefined,
                        address: s.address.trim(),
                        customer_name: s.customer_name.trim(),
                        customer_phone: s.customer_phone.trim() || undefined,
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
            setError(err.message || 'Error al guardar la ruta');
        } finally {
            setLoading(false);
        }
    };

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

                {/* Driver ID */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        ID Conductor
                    </label>
                    <Input
                        type="number"
                        value={driverId}
                        onChange={(e) => setDriverId(e.target.value)}
                        placeholder="ID del conductor"
                        min={1}
                    />
                </div>

                {/* Vehicle ID */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        ID Vehiculo
                    </label>
                    <Input
                        type="number"
                        value={vehicleId}
                        onChange={(e) => setVehicleId(e.target.value)}
                        placeholder="ID del vehiculo"
                        min={1}
                    />
                </div>

                {/* Origin Address */}
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

                {/* Notes */}
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

            {/* Stops section - only for creation */}
            {!isEditing && (
                <div className="space-y-3">
                    <div className="flex items-center justify-between">
                        <h3 className="text-sm font-semibold text-gray-700">Paradas</h3>
                        <button
                            type="button"
                            onClick={addStop}
                            className="inline-flex items-center gap-1 text-sm text-blue-600 hover:text-blue-700"
                        >
                            <PlusIcon className="w-4 h-4" />
                            Agregar parada
                        </button>
                    </div>

                    {stops.length === 0 && (
                        <p className="text-sm text-gray-400 text-center py-4">
                            No hay paradas. Puedes agregar paradas despues de crear la ruta.
                        </p>
                    )}

                    {stops.map((stop, index) => (
                        <div key={index} className="border border-gray-200 rounded-lg p-3 space-y-2 relative">
                            <div className="flex items-center justify-between mb-1">
                                <span className="text-xs font-medium text-gray-500">Parada {index + 1}</span>
                                {stops.length > 1 && (
                                    <button
                                        type="button"
                                        onClick={() => removeStop(index)}
                                        className="text-red-400 hover:text-red-600"
                                        title="Eliminar parada"
                                    >
                                        <TrashIcon className="w-4 h-4" />
                                    </button>
                                )}
                            </div>
                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                                <Input
                                    type="text"
                                    value={stop.order_id}
                                    onChange={(e) => updateStopField(index, 'order_id', e.target.value)}
                                    placeholder="ID de orden (opcional)"
                                />
                                <Input
                                    type="text"
                                    value={stop.address}
                                    onChange={(e) => updateStopField(index, 'address', e.target.value)}
                                    placeholder="Direccion *"
                                />
                                <Input
                                    type="text"
                                    value={stop.customer_name}
                                    onChange={(e) => updateStopField(index, 'customer_name', e.target.value)}
                                    placeholder="Nombre del cliente *"
                                />
                                <Input
                                    type="text"
                                    value={stop.customer_phone}
                                    onChange={(e) => updateStopField(index, 'customer_phone', e.target.value)}
                                    placeholder="Telefono del cliente"
                                />
                            </div>
                        </div>
                    ))}
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
