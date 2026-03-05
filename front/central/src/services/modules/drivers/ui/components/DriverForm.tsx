'use client';

import { useState } from 'react';
import { DriverInfo, CreateDriverDTO, UpdateDriverDTO } from '../../domain/types';
import { createDriverAction, updateDriverAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';

interface DriverFormProps {
    driver?: DriverInfo;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

const LICENSE_TYPES = ['A1', 'A2', 'B1', 'B2', 'C1'];

export default function DriverForm({ driver, onSuccess, onCancel, businessId }: DriverFormProps) {
    const [formData, setFormData] = useState({
        first_name: driver?.first_name || '',
        last_name: driver?.last_name || '',
        identification: driver?.identification || '',
        phone: driver?.phone || '',
        email: driver?.email || '',
        license_type: driver?.license_type || '',
        license_expiry: driver?.license_expiry ? driver.license_expiry.split('T')[0] : '',
        notes: driver?.notes || '',
        status: driver?.status || 'active',
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
            if (driver) {
                const updateData: UpdateDriverDTO = {
                    first_name: formData.first_name,
                    last_name: formData.last_name,
                    identification: formData.identification,
                    phone: formData.phone,
                    email: formData.email || undefined,
                    license_type: formData.license_type || undefined,
                    license_expiry: formData.license_expiry || undefined,
                    notes: formData.notes || null,
                    status: formData.status,
                };
                await updateDriverAction(driver.id, updateData, businessId);
            } else {
                const createData: CreateDriverDTO = {
                    first_name: formData.first_name,
                    last_name: formData.last_name,
                    identification: formData.identification,
                    phone: formData.phone,
                    email: formData.email || undefined,
                    license_type: formData.license_type || undefined,
                    license_expiry: formData.license_expiry || undefined,
                    notes: formData.notes || undefined,
                };
                await createDriverAction(createData, businessId);
            }

            setSuccess(driver ? 'Conductor actualizado exitosamente' : 'Conductor creado exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(err.message || 'Error al guardar el conductor');
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
                {/* Nombre */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Nombre <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.first_name}
                        onChange={(e) => handleChange('first_name', e.target.value)}
                        placeholder="Nombre"
                        required
                        minLength={2}
                        maxLength={255}
                    />
                </div>

                {/* Apellido */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Apellido <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.last_name}
                        onChange={(e) => handleChange('last_name', e.target.value)}
                        placeholder="Apellido"
                        required
                        minLength={2}
                        maxLength={255}
                    />
                </div>

                {/* Identificacion */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Identificacion <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.identification}
                        onChange={(e) => handleChange('identification', e.target.value)}
                        placeholder="Numero de documento"
                        required
                        maxLength={30}
                    />
                </div>

                {/* Telefono */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Telefono <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="tel"
                        value={formData.phone}
                        onChange={(e) => handleChange('phone', e.target.value)}
                        placeholder="3001234567"
                        required
                        maxLength={20}
                    />
                </div>

                {/* Email */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Email
                    </label>
                    <Input
                        type="email"
                        value={formData.email}
                        onChange={(e) => handleChange('email', e.target.value)}
                        placeholder="correo@ejemplo.com"
                        maxLength={255}
                    />
                </div>

                {/* Tipo de licencia */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Tipo de licencia
                    </label>
                    <select
                        value={formData.license_type}
                        onChange={(e) => handleChange('license_type', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                        <option value="">Sin especificar</option>
                        {LICENSE_TYPES.map((type) => (
                            <option key={type} value={type}>{type}</option>
                        ))}
                    </select>
                </div>

                {/* Fecha de vencimiento de licencia */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Vencimiento de licencia
                    </label>
                    <Input
                        type="date"
                        value={formData.license_expiry}
                        onChange={(e) => handleChange('license_expiry', e.target.value)}
                    />
                </div>

                {/* Estado (solo en edicion) */}
                {driver && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Estado
                        </label>
                        <select
                            value={formData.status}
                            onChange={(e) => handleChange('status', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                            <option value="active">Activo</option>
                            <option value="inactive">Inactivo</option>
                        </select>
                    </div>
                )}

                {/* Notas */}
                <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Notas
                    </label>
                    <textarea
                        value={formData.notes}
                        onChange={(e) => handleChange('notes', e.target.value)}
                        placeholder="Notas adicionales sobre el conductor..."
                        rows={3}
                        maxLength={500}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                    />
                </div>
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t">
                <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
                    Cancelar
                </Button>
                <Button type="submit" variant="primary" disabled={loading}>
                    {loading ? 'Guardando...' : driver ? 'Actualizar' : 'Crear conductor'}
                </Button>
            </div>
        </form>
    );
}
