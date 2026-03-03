'use client';

import { useState } from 'react';
import { ChevronDownIcon, ChevronRightIcon } from '@heroicons/react/24/outline';
import { Warehouse, CreateWarehouseDTO, UpdateWarehouseDTO } from '../../domain/types';
import { createWarehouseAction, updateWarehouseAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';

interface WarehouseFormProps {
    warehouse?: Warehouse;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

function CollapsibleSection({ title, description, defaultOpen = false, children }: {
    title: string;
    description?: string;
    defaultOpen?: boolean;
    children: React.ReactNode;
}) {
    const [open, setOpen] = useState(defaultOpen);
    return (
        <div className="border border-gray-200 rounded-lg">
            <button
                type="button"
                onClick={() => setOpen(!open)}
                className="w-full flex items-center justify-between px-4 py-3 text-left hover:bg-gray-50 transition-colors rounded-lg"
            >
                <div>
                    <span className="text-sm font-medium text-gray-900">{title}</span>
                    {description && <p className="text-xs text-gray-500 mt-0.5">{description}</p>}
                </div>
                {open
                    ? <ChevronDownIcon className="w-4 h-4 text-gray-400" />
                    : <ChevronRightIcon className="w-4 h-4 text-gray-400" />
                }
            </button>
            {open && <div className="px-4 pb-4 pt-1">{children}</div>}
        </div>
    );
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
        // Carrier fields
        company: warehouse?.company || '',
        first_name: warehouse?.first_name || '',
        last_name: warehouse?.last_name || '',
        email: warehouse?.email || '',
        street: warehouse?.street || '',
        suburb: warehouse?.suburb || '',
        city_dane_code: warehouse?.city_dane_code || '',
        postal_code: warehouse?.postal_code || '',
        // GPS
        latitude: warehouse?.latitude != null ? String(warehouse.latitude) : '',
        longitude: warehouse?.longitude != null ? String(warehouse.longitude) : '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const handleChange = (field: string, value: string | boolean) => {
        setFormData((prev) => ({ ...prev, [field]: value }));
    };

    const hasCarrierData = !!(warehouse?.company || warehouse?.first_name || warehouse?.last_name || warehouse?.email || warehouse?.street || warehouse?.suburb || warehouse?.city_dane_code || warehouse?.postal_code);
    const hasGpsData = warehouse?.latitude != null || warehouse?.longitude != null;

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        const carrierFields = {
            company: formData.company || undefined,
            first_name: formData.first_name || undefined,
            last_name: formData.last_name || undefined,
            email: formData.email || undefined,
            street: formData.street || undefined,
            suburb: formData.suburb || undefined,
            city_dane_code: formData.city_dane_code || undefined,
            postal_code: formData.postal_code || undefined,
            latitude: formData.latitude ? parseFloat(formData.latitude) : null,
            longitude: formData.longitude ? parseFloat(formData.longitude) : null,
        };

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
                    ...carrierFields,
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
                    ...carrierFields,
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
        <form onSubmit={handleSubmit} className="space-y-5">
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

            {/* Información básica */}
            <div>
                <h3 className="text-sm font-medium text-gray-700 mb-3">Información básica</h3>
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
                </div>
            </div>

            {/* Dirección general */}
            <div>
                <h3 className="text-sm font-medium text-gray-700 mb-3">Dirección</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="md:col-span-2">
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
                </div>
            </div>

            {/* Contacto */}
            <div>
                <h3 className="text-sm font-medium text-gray-700 mb-3">Contacto</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
            </div>

            {/* Datos de transportadora (colapsable) */}
            <CollapsibleSection
                title="Datos de transportadora"
                description="Estos datos se usan al generar guías de envío"
                defaultOpen={hasCarrierData}
            >
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Empresa</label>
                        <Input
                            type="text"
                            value={formData.company}
                            onChange={(e) => handleChange('company', e.target.value)}
                            placeholder="Mi empresa S.A.S"
                            maxLength={255}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
                        <Input
                            type="email"
                            value={formData.email}
                            onChange={(e) => handleChange('email', e.target.value)}
                            placeholder="bodega@empresa.com"
                            maxLength={255}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Nombre</label>
                        <Input
                            type="text"
                            value={formData.first_name}
                            onChange={(e) => handleChange('first_name', e.target.value)}
                            placeholder="Juan"
                            maxLength={255}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Apellido</label>
                        <Input
                            type="text"
                            value={formData.last_name}
                            onChange={(e) => handleChange('last_name', e.target.value)}
                            placeholder="Pérez"
                            maxLength={255}
                        />
                    </div>
                    <div className="md:col-span-2">
                        <label className="block text-sm font-medium text-gray-700 mb-1">Calle</label>
                        <Input
                            type="text"
                            value={formData.street}
                            onChange={(e) => handleChange('street', e.target.value)}
                            placeholder="Calle 100 #15-20"
                            maxLength={255}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Barrio / Suburb</label>
                        <Input
                            type="text"
                            value={formData.suburb}
                            onChange={(e) => handleChange('suburb', e.target.value)}
                            placeholder="Chapinero"
                            maxLength={255}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Código DANE ciudad</label>
                        <Input
                            type="text"
                            value={formData.city_dane_code}
                            onChange={(e) => handleChange('city_dane_code', e.target.value)}
                            placeholder="11001"
                            maxLength={20}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Código postal (carrier)</label>
                        <Input
                            type="text"
                            value={formData.postal_code}
                            onChange={(e) => handleChange('postal_code', e.target.value)}
                            placeholder="110111"
                            maxLength={20}
                        />
                    </div>
                </div>
            </CollapsibleSection>

            {/* Ubicación GPS (colapsable) */}
            <CollapsibleSection
                title="Ubicación GPS"
                description="Coordenadas para mostrar en el mapa"
                defaultOpen={hasGpsData}
            >
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Latitud</label>
                        <Input
                            type="number"
                            step="any"
                            value={formData.latitude}
                            onChange={(e) => handleChange('latitude', e.target.value)}
                            placeholder="4.7110"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Longitud</label>
                        <Input
                            type="number"
                            step="any"
                            value={formData.longitude}
                            onChange={(e) => handleChange('longitude', e.target.value)}
                            placeholder="-74.0721"
                        />
                    </div>
                </div>
            </CollapsibleSection>

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
