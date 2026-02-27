'use client';

import { useState } from 'react';
import { Warehouse, CreateWarehouseDTO, UpdateWarehouseDTO } from '../../domain/types';
import { createWarehouseAction, updateWarehouseAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';

interface WarehouseFormProps {
    warehouse?: Warehouse;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

export default function WarehouseForm({ warehouse, onSuccess, onCancel, businessId }: WarehouseFormProps) {
    const [formData, setFormData] = useState({
        name: warehouse?.name || '',
        code: warehouse?.code || '',
        address: warehouse?.address || '',
        city: warehouse?.city || '',
        state: warehouse?.state || '',
        country: warehouse?.country || '',
        zip_code: warehouse?.zip_code || '',
        phone: warehouse?.phone || '',
        contact_name: warehouse?.contact_name || '',
        contact_email: warehouse?.contact_email || '',
        is_default: warehouse?.is_default || false,
        is_fulfillment: warehouse?.is_fulfillment || false,
        is_active: warehouse?.is_active ?? true,
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const handleChange = (field: string, value: string | boolean) => {
        setFormData((prev) => ({ ...prev, [field]: value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            if (warehouse) {
                const updateData: UpdateWarehouseDTO = {
                    name: formData.name,
                    code: formData.code,
                    address: formData.address || undefined,
                    city: formData.city || undefined,
                    state: formData.state || undefined,
                    country: formData.country || undefined,
                    zip_code: formData.zip_code || undefined,
                    phone: formData.phone || undefined,
                    contact_name: formData.contact_name || undefined,
                    contact_email: formData.contact_email || undefined,
                    is_active: formData.is_active,
                    is_default: formData.is_default,
                    is_fulfillment: formData.is_fulfillment,
                };
                await updateWarehouseAction(warehouse.id, updateData, businessId);
            } else {
                const createData: CreateWarehouseDTO = {
                    name: formData.name,
                    code: formData.code,
                    address: formData.address || undefined,
                    city: formData.city || undefined,
                    state: formData.state || undefined,
                    country: formData.country || undefined,
                    zip_code: formData.zip_code || undefined,
                    phone: formData.phone || undefined,
                    contact_name: formData.contact_name || undefined,
                    contact_email: formData.contact_email || undefined,
                    is_default: formData.is_default,
                    is_fulfillment: formData.is_fulfillment,
                };
                await createWarehouseAction(createData, businessId);
            }

            setSuccess(warehouse ? 'Bodega actualizada exitosamente' : 'Bodega creada exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(err.message || 'Error al guardar la bodega');
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
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Nombre <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.name}
                        onChange={(e) => handleChange('name', e.target.value)}
                        placeholder="Bodega principal"
                        required
                        minLength={2}
                        maxLength={255}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Código <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.code}
                        onChange={(e) => handleChange('code', e.target.value.toUpperCase())}
                        placeholder="BOD-001"
                        required
                        minLength={1}
                        maxLength={50}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Dirección</label>
                    <Input
                        type="text"
                        value={formData.address}
                        onChange={(e) => handleChange('address', e.target.value)}
                        placeholder="Calle 123 #45-67"
                        maxLength={255}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Ciudad</label>
                    <Input
                        type="text"
                        value={formData.city}
                        onChange={(e) => handleChange('city', e.target.value)}
                        placeholder="Bogotá"
                        maxLength={100}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Departamento / Estado</label>
                    <Input
                        type="text"
                        value={formData.state}
                        onChange={(e) => handleChange('state', e.target.value)}
                        placeholder="Cundinamarca"
                        maxLength={100}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">País</label>
                    <Input
                        type="text"
                        value={formData.country}
                        onChange={(e) => handleChange('country', e.target.value)}
                        placeholder="Colombia"
                        maxLength={100}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Código postal</label>
                    <Input
                        type="text"
                        value={formData.zip_code}
                        onChange={(e) => handleChange('zip_code', e.target.value)}
                        placeholder="110111"
                        maxLength={20}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Teléfono</label>
                    <Input
                        type="tel"
                        value={formData.phone}
                        onChange={(e) => handleChange('phone', e.target.value)}
                        placeholder="3001234567"
                        maxLength={20}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Nombre contacto</label>
                    <Input
                        type="text"
                        value={formData.contact_name}
                        onChange={(e) => handleChange('contact_name', e.target.value)}
                        placeholder="Juan Pérez"
                        maxLength={255}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Email contacto</label>
                    <Input
                        type="email"
                        value={formData.contact_email}
                        onChange={(e) => handleChange('contact_email', e.target.value)}
                        placeholder="contacto@empresa.com"
                        maxLength={255}
                    />
                </div>
            </div>

            {/* Toggles */}
            <div className="space-y-3 border-t pt-4">
                <label className="flex items-center gap-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={formData.is_default}
                        onChange={(e) => handleChange('is_default', e.target.checked)}
                        className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <span className="text-sm text-gray-700">Bodega principal (por defecto para nuevas órdenes)</span>
                </label>

                <label className="flex items-center gap-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={formData.is_fulfillment}
                        onChange={(e) => handleChange('is_fulfillment', e.target.checked)}
                        className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                    />
                    <span className="text-sm text-gray-700">Bodega de fulfillment (despacho de pedidos)</span>
                </label>

                {warehouse && (
                    <label className="flex items-center gap-3 cursor-pointer">
                        <input
                            type="checkbox"
                            checked={formData.is_active}
                            onChange={(e) => handleChange('is_active', e.target.checked)}
                            className="w-4 h-4 rounded border-gray-300 text-green-600 focus:ring-green-500"
                        />
                        <span className="text-sm text-gray-700">Bodega activa</span>
                    </label>
                )}
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t">
                <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
                    Cancelar
                </Button>
                <Button type="submit" variant="primary" disabled={loading}>
                    {loading ? 'Guardando...' : warehouse ? 'Actualizar' : 'Crear bodega'}
                </Button>
            </div>
        </form>
    );
}
