'use client';

import { useState } from 'react';
import dynamic from 'next/dynamic';
import { Warehouse, CreateWarehouseDTO, UpdateWarehouseDTO } from '../../domain/types';
import { createWarehouseAction, updateWarehouseAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';
import AddressAutocomplete, { AddressSuggestion } from '@/services/modules/orders/ui/components/AddressAutocomplete';

const MapComponent = dynamic(() => import('@/shared/ui/MapComponent'), { ssr: false });

interface WarehouseFormProps {
    warehouse?: Warehouse;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

type StructureMode = 'simple' | 'zones' | 'wms';

const STRUCTURE_OPTIONS: {
    id: StructureMode;
    icon: React.ReactNode;
    title: string;
    description: string;
    levels: string;
}[] = [
    {
        id: 'simple',
        icon: (
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                    d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
        ),
        title: 'Simple',
        description: 'Sin zonas ni ubicaciones.',
        levels: 'Solo la bodega',
    },
    {
        id: 'zones',
        icon: (
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                    d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
            </svg>
        ),
        title: 'Con Zonas',
        description: 'Areas funcionales diferenciadas.',
        levels: 'Bodega → Zonas → Ubicaciones',
    },
    {
        id: 'wms',
        icon: (
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
        ),
        title: 'WMS Completo',
        description: 'Jerarquia completa avanzada.',
        levels: 'Zonas → Pasillos → Racks → Niveles → Posiciones',
    },
];

function generateCode(name: string): string {
    const slug = name.trim().toUpperCase().replace(/[^A-Z0-9]+/g, '-').replace(/^-|-$/g, '').slice(0, 20);
    return slug || 'WH-' + Math.random().toString(36).slice(2, 6).toUpperCase();
}


export default function WarehouseForm({ warehouse, onSuccess, onCancel, businessId }: WarehouseFormProps) {
    const [structureMode, setStructureMode] = useState<StructureMode>(() => {
        if (warehouse?.id) {
            try {
                const saved = localStorage.getItem(`wh_struct_${warehouse.id}`);
                if (saved === 'simple' || saved === 'zones' || saved === 'wms') return saved as StructureMode;
            } catch {}
        }
        return 'simple';
    });

    const [addressCoords, setAddressCoords] = useState<{ lat: number; lon: number } | null>(() => {
        if (warehouse?.latitude != null && warehouse?.longitude != null) {
            return { lat: warehouse.latitude, lon: warehouse.longitude };
        }
        return null;
    });

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
        is_active: warehouse?.is_active ?? true,
        latitude: warehouse?.latitude != null ? String(warehouse.latitude) : '',
        longitude: warehouse?.longitude != null ? String(warehouse.longitude) : '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const handleChange = (field: string, value: string | boolean) => {
        setFormData((prev) => ({ ...prev, [field]: value }));
    };

    const handleAddressSelect = (s: AddressSuggestion) => {
        if (s.lat && s.lon) {
            setAddressCoords({ lat: s.lat, lon: s.lon });
            setFormData(prev => ({
                ...prev,
                latitude: String(s.lat),
                longitude: String(s.lon),
                ...(s.city ? { city: s.city } : {}),
                ...(s.state ? { state: s.state } : {}),
                ...(s.postcode ? { zip_code: s.postcode } : {}),
            }));
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        const resolvedCode = formData.code || generateCode(formData.name);

        try {
            if (warehouse) {
                const updateData: UpdateWarehouseDTO = {
                    name: formData.name,
                    code: resolvedCode,
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
                    latitude: formData.latitude ? parseFloat(formData.latitude) : null,
                    longitude: formData.longitude ? parseFloat(formData.longitude) : null,
                };
                await updateWarehouseAction(warehouse.id, updateData, businessId);
                try { localStorage.setItem(`wh_struct_${warehouse.id}`, structureMode); } catch {}
            } else {
                const createData: CreateWarehouseDTO = {
                    name: formData.name,
                    code: resolvedCode,
                    address: formData.address || undefined,
                    city: formData.city || undefined,
                    state: formData.state || undefined,
                    country: formData.country || undefined,
                    zip_code: formData.zip_code || undefined,
                    phone: formData.phone || undefined,
                    contact_name: formData.contact_name || undefined,
                    contact_email: formData.contact_email || undefined,
                    is_default: formData.is_default,
                    latitude: formData.latitude ? parseFloat(formData.latitude) : null,
                    longitude: formData.longitude ? parseFloat(formData.longitude) : null,
                };
                const result = await createWarehouseAction(createData, businessId);
                try {
                    const newId = (result as any)?.data?.id ?? (result as any)?.id;
                    if (newId) localStorage.setItem(`wh_struct_${newId}`, structureMode);
                } catch {}
            }

            setSuccess(warehouse ? 'Bodega actualizada exitosamente' : 'Bodega creada exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar la bodega'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-5">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}
            {success && <Alert type="success" onClose={() => setSuccess(null)}>{success}</Alert>}

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">

                {/* LEFT COLUMN */}
                <div className="space-y-5">

                    <div>
                        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">Informacion basica</h3>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">
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
                    </div>

                    <div>
                        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-1">Estructura de la bodega</h3>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mb-3">Define como organizar el espacio fisico.</p>
                        <div className="flex flex-col gap-2">
                            {STRUCTURE_OPTIONS.map((opt) => {
                                const selected = structureMode === opt.id;
                                return (
                                    <button
                                        key={opt.id}
                                        type="button"
                                        onClick={() => setStructureMode(opt.id)}
                                        className={`text-left px-4 py-3 rounded-xl border-2 transition-all duration-200 flex items-start gap-3 ${
                                            selected
                                                ? 'border-[#7c3aed] bg-purple-50 dark:bg-purple-900/20 shadow-sm'
                                                : 'border-gray-200 dark:border-gray-700 hover:border-purple-300 dark:hover:border-purple-600 bg-white dark:bg-gray-800'
                                        }`}
                                    >
                                        <div className={`mt-0.5 flex-shrink-0 ${selected ? 'text-[#7c3aed]' : 'text-gray-400 dark:text-gray-500'}`}>
                                            {opt.icon}
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-center gap-2">
                                                <span className={`text-sm font-bold ${selected ? 'text-[#7c3aed]' : 'text-gray-800 dark:text-gray-200'}`}>
                                                    {opt.title}
                                                </span>
                                                {selected && (
                                                    <svg className="w-4 h-4 text-[#7c3aed] flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                                                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                                                    </svg>
                                                )}
                                            </div>
                                            <p className="text-xs text-gray-500 dark:text-gray-400">{opt.description}</p>
                                            <p className="text-xs text-gray-400 dark:text-gray-500 mt-0.5">
                                                <span className="font-medium">Niveles:</span> {opt.levels}
                                            </p>
                                        </div>
                                    </button>
                                );
                            })}
                        </div>
                    </div>

                    <div className="border-t dark:border-gray-700 pt-4 space-y-3">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm font-medium text-gray-700 dark:text-gray-200">Bodega principal</p>
                                <p className="text-xs text-gray-500 dark:text-gray-400">Por defecto para nuevas ordenes</p>
                            </div>
                            <button
                                type="button"
                                onClick={() => handleChange('is_default', !formData.is_default)}
                                className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none ${formData.is_default ? 'bg-[#7c3aed]' : 'bg-gray-200 dark:bg-gray-600'}`}
                            >
                                <span className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${formData.is_default ? 'translate-x-5' : 'translate-x-0'}`} />
                            </button>
                        </div>
                        {warehouse && (
                            <div className="flex items-center justify-between">
                                <div>
                                    <p className="text-sm font-medium text-gray-700 dark:text-gray-200">Bodega activa</p>
                                    <p className="text-xs text-gray-500 dark:text-gray-400">Disponible para operaciones</p>
                                </div>
                                <button
                                    type="button"
                                    onClick={() => handleChange('is_active', !formData.is_active)}
                                    className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none ${formData.is_active ? 'bg-green-500' : 'bg-gray-200 dark:bg-gray-600'}`}
                                >
                                    <span className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${formData.is_active ? 'translate-x-5' : 'translate-x-0'}`} />
                                </button>
                            </div>
                        )}
                    </div>
                </div>

                {/* RIGHT COLUMN */}
                <div className="space-y-4">

                    <div>
                        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">Direccion</h3>
                        <div className="space-y-3">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Direccion</label>
                                <AddressAutocomplete
                                    value={formData.address}
                                    onChange={(val) => handleChange('address', val)}
                                    onSelect={handleAddressSelect}
                                    placeholder="Calle 123 #45-67..."
                                    city={formData.city}
                                />
                            </div>
                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Ciudad</label>
                                    <Input
                                        type="text"
                                        value={formData.city}
                                        onChange={(e) => handleChange('city', e.target.value)}
                                        placeholder="Bogota"
                                        maxLength={100}
                                    />
                                </div>
                                <div>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Departamento</label>
                                    <Input
                                        type="text"
                                        value={formData.state}
                                        onChange={(e) => handleChange('state', e.target.value)}
                                        placeholder="Cundinamarca"
                                        maxLength={100}
                                    />
                                </div>
                                <div>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Pais</label>
                                    <Input
                                        type="text"
                                        value={formData.country}
                                        onChange={(e) => handleChange('country', e.target.value)}
                                        placeholder="Colombia"
                                        maxLength={100}
                                    />
                                </div>
                                <div>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Cod. postal</label>
                                    <Input
                                        type="text"
                                        value={formData.zip_code}
                                        onChange={(e) => handleChange('zip_code', e.target.value)}
                                        placeholder="110111"
                                        maxLength={20}
                                    />
                                </div>
                            </div>

                            {addressCoords && (
                                <MapComponent
                                    address={formData.address}
                                    city={formData.city}
                                    latitude={addressCoords.lat}
                                    longitude={addressCoords.lon}
                                    height="200px"
                                />
                            )}
                        </div>
                    </div>

                    <div>
                        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">Contacto</h3>
                        <div className="grid grid-cols-2 gap-3">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Telefono</label>
                                <Input
                                    type="tel"
                                    value={formData.phone}
                                    onChange={(e) => handleChange('phone', e.target.value)}
                                    placeholder="3001234567"
                                    maxLength={20}
                                />
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Nombre contacto</label>
                                <Input
                                    type="text"
                                    value={formData.contact_name}
                                    onChange={(e) => handleChange('contact_name', e.target.value)}
                                    placeholder="Juan Perez"
                                    maxLength={255}
                                />
                            </div>
                            <div className="col-span-2">
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Email contacto</label>
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

                </div>
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t dark:border-gray-700">
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
