'use client';

import { useState } from 'react';
import { VehicleInfo, CreateVehicleDTO, UpdateVehicleDTO } from '../../domain/types';
import { createVehicleAction, updateVehicleAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface VehicleFormProps {
    vehicle?: VehicleInfo;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

const VEHICLE_TYPES = [
    { value: 'motorcycle', label: 'Motocicleta' },
    { value: 'car', label: 'Carro' },
    { value: 'van', label: 'Van' },
    { value: 'truck', label: 'Camion' },
];

const VEHICLE_STATUSES = [
    { value: 'active', label: 'Activo' },
    { value: 'inactive', label: 'Inactivo' },
    { value: 'in_maintenance', label: 'En mantenimiento' },
];

export default function VehicleForm({ vehicle, onSuccess, onCancel, businessId }: VehicleFormProps) {
    const [formData, setFormData] = useState({
        type: vehicle?.type || 'motorcycle',
        license_plate: vehicle?.license_plate || '',
        brand: vehicle?.brand || '',
        model: vehicle?.model || '',
        year: vehicle?.year != null ? String(vehicle.year) : '',
        color: vehicle?.color || '',
        status: vehicle?.status || 'active',
        weight_capacity_kg: vehicle?.weight_capacity_kg != null ? String(vehicle.weight_capacity_kg) : '',
        volume_capacity_m3: vehicle?.volume_capacity_m3 != null ? String(vehicle.volume_capacity_m3) : '',
        insurance_expiry: vehicle?.insurance_expiry ? vehicle.insurance_expiry.split('T')[0] : '',
        registration_expiry: vehicle?.registration_expiry ? vehicle.registration_expiry.split('T')[0] : '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const handleChange = (field: string, value: string) => {
        setFormData((prev) => ({ ...prev, [field]: value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            if (vehicle) {
                const updateData: UpdateVehicleDTO = {
                    type: formData.type,
                    license_plate: formData.license_plate,
                    brand: formData.brand || undefined,
                    model: formData.model || undefined,
                    year: formData.year ? Number(formData.year) : null,
                    color: formData.color || undefined,
                    status: formData.status,
                    weight_capacity_kg: formData.weight_capacity_kg ? Number(formData.weight_capacity_kg) : null,
                    volume_capacity_m3: formData.volume_capacity_m3 ? Number(formData.volume_capacity_m3) : null,
                    insurance_expiry: formData.insurance_expiry ? `${formData.insurance_expiry}T00:00:00Z` : null,
                    registration_expiry: formData.registration_expiry ? `${formData.registration_expiry}T00:00:00Z` : null,
                };
                await updateVehicleAction(vehicle.id, updateData, businessId);
            } else {
                const createData: CreateVehicleDTO = {
                    type: formData.type,
                    license_plate: formData.license_plate,
                    brand: formData.brand || undefined,
                    model: formData.model || undefined,
                    year: formData.year ? Number(formData.year) : undefined,
                    color: formData.color || undefined,
                    weight_capacity_kg: formData.weight_capacity_kg ? Number(formData.weight_capacity_kg) : undefined,
                    volume_capacity_m3: formData.volume_capacity_m3 ? Number(formData.volume_capacity_m3) : undefined,
                    insurance_expiry: formData.insurance_expiry ? `${formData.insurance_expiry}T00:00:00Z` : undefined,
                    registration_expiry: formData.registration_expiry ? `${formData.registration_expiry}T00:00:00Z` : undefined,
                };
                await createVehicleAction(createData, businessId);
            }

            setSuccess(vehicle ? 'Vehiculo actualizado exitosamente' : 'Vehiculo creado exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar el vehiculo'));
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

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Tipo */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Tipo <span className="text-red-500">*</span>
                    </label>
                    <select
                        value={formData.type}
                        onChange={(e) => handleChange('type', e.target.value)}
                        required
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                        {VEHICLE_TYPES.map((t) => (
                            <option key={t.value} value={t.value}>{t.label}</option>
                        ))}
                    </select>
                </div>

                {/* Placa */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Placa <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.license_plate}
                        onChange={(e) => handleChange('license_plate', e.target.value)}
                        placeholder="ABC123"
                        required
                        maxLength={20}
                    />
                </div>

                {/* Marca */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Marca
                    </label>
                    <Input
                        type="text"
                        value={formData.brand}
                        onChange={(e) => handleChange('brand', e.target.value)}
                        placeholder="Yamaha, Chevrolet..."
                        maxLength={100}
                    />
                </div>

                {/* Modelo */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Modelo
                    </label>
                    <Input
                        type="text"
                        value={formData.model}
                        onChange={(e) => handleChange('model', e.target.value)}
                        placeholder="FZ25, Spark..."
                        maxLength={100}
                    />
                </div>

                {/* Ano */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Ano
                    </label>
                    <Input
                        type="number"
                        value={formData.year}
                        onChange={(e) => handleChange('year', e.target.value)}
                        placeholder="2024"
                        min={1900}
                        max={2100}
                    />
                </div>

                {/* Color */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Color
                    </label>
                    <Input
                        type="text"
                        value={formData.color}
                        onChange={(e) => handleChange('color', e.target.value)}
                        placeholder="Rojo, Azul..."
                        maxLength={50}
                    />
                </div>

                {/* Estado (solo en edicion) */}
                {vehicle && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Estado
                        </label>
                        <select
                            value={formData.status}
                            onChange={(e) => handleChange('status', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                            {VEHICLE_STATUSES.map((s) => (
                                <option key={s.value} value={s.value}>{s.label}</option>
                            ))}
                        </select>
                    </div>
                )}

                {/* Capacidad de peso */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Capacidad de peso (kg)
                    </label>
                    <Input
                        type="number"
                        value={formData.weight_capacity_kg}
                        onChange={(e) => handleChange('weight_capacity_kg', e.target.value)}
                        placeholder="50"
                        min={0}
                        step="0.1"
                    />
                </div>

                {/* Capacidad de volumen */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Capacidad de volumen (m3)
                    </label>
                    <Input
                        type="number"
                        value={formData.volume_capacity_m3}
                        onChange={(e) => handleChange('volume_capacity_m3', e.target.value)}
                        placeholder="0.5"
                        min={0}
                        step="0.01"
                    />
                </div>

                {/* Vencimiento seguro */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Vencimiento del seguro
                    </label>
                    <Input
                        type="date"
                        value={formData.insurance_expiry}
                        onChange={(e) => handleChange('insurance_expiry', e.target.value)}
                    />
                </div>

                {/* Vencimiento registro */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Vencimiento del registro
                    </label>
                    <Input
                        type="date"
                        value={formData.registration_expiry}
                        onChange={(e) => handleChange('registration_expiry', e.target.value)}
                    />
                </div>
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t">
                <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
                    Cancelar
                </Button>
                <Button type="submit" variant="primary" disabled={loading}>
                    {loading ? 'Guardando...' : vehicle ? 'Actualizar' : 'Crear vehiculo'}
                </Button>
            </div>
        </form>
    );
}
