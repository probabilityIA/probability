'use client';

import React, { useState, useEffect, useTransition } from 'react';
import {
    getOriginAddressesAction,
    createOriginAddressAction,
    updateOriginAddressAction,
    deleteOriginAddressAction
} from '@/services/modules/shipments/infra/actions';
import { OriginAddress, CreateOriginAddressRequest } from '@/services/modules/shipments/domain/types';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { useToast } from '@/shared/providers/toast-provider';
import { useHasPermission } from '@/shared/contexts/permissions-context';
import { Plus, Edit2, Trash2, CheckCircle, MapPin, Phone, Mail, Building, X } from 'lucide-react';
import danes from '@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json';

export function OriginAddressManager() {
    const { showToast } = useToast();
    const canCreate = useHasPermission('Envios', 'Create');
    const canUpdate = useHasPermission('Envios', 'Update');
    const canDelete = useHasPermission('Envios', 'Delete');
    const [addresses, setAddresses] = useState<OriginAddress[]>([]);
    const [isFormOpen, setIsFormOpen] = useState(false);
    const [editingAddress, setEditingAddress] = useState<OriginAddress | null>(null);
    const [isPending, startTransition] = useTransition();

    // Form state
    const [formData, setFormData] = useState<CreateOriginAddressRequest>({
        alias: '',
        company: '',
        first_name: '',
        last_name: '',
        email: '',
        phone: '',
        street: '',
        suburb: '',
        city_dane_code: '',
        city: '',
        state: '',
        postal_code: '',
    });

    // Search city state
    const [citySearch, setCitySearch] = useState('');
    const [citySuggestions, setCitySuggestions] = useState<any[]>([]);

    const loadAddresses = async () => {
        const result = await getOriginAddressesAction();
        if (result.success && result.data) {
            setAddresses(result.data);
        } else {
            showToast(result.message || 'Error al cargar direcciones', 'error');
        }
    };

    useEffect(() => {
        loadAddresses();
    }, []);

    const handleCitySearch = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setCitySearch(value);
        if (value.length > 2) {
            const results = Object.entries(danes)
                .filter(([_, data]: [string, any]) =>
                    data.ciudad.toUpperCase().includes(value.toUpperCase())
                )
                .slice(0, 5)
                .map(([code, data]: [string, any]) => ({
                    code,
                    ...data
                }));
            setCitySuggestions(results);
        } else {
            setCitySuggestions([]);
        }
    };

    const selectCity = (suggestion: any) => {
        setFormData({
            ...formData,
            city_dane_code: suggestion.code,
            city: suggestion.ciudad,
            state: suggestion.departamento,
        });
        setCitySearch(suggestion.ciudad);
        setCitySuggestions([]);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        startTransition(async () => {
            let result;
            if (editingAddress) {
                result = await updateOriginAddressAction(editingAddress.id, formData);
            } else {
                result = await createOriginAddressAction(formData);
            }

            if (result.success) {
                showToast(editingAddress ? 'Dirección actualizada' : 'Dirección creada', 'success');
                setIsFormOpen(false);
                setEditingAddress(null);
                setFormData({
                    alias: '', company: '', first_name: '', last_name: '',
                    email: '', phone: '', street: '', suburb: '',
                    city_dane_code: '', city: '', state: '', postal_code: ''
                });
                setCitySearch('');
                loadAddresses();
            } else {
                showToast(result.message, 'error');
            }
        });
    };

    const handleEdit = (address: OriginAddress) => {
        setEditingAddress(address);
        setFormData({
            alias: address.alias,
            company: address.company,
            first_name: address.first_name,
            last_name: address.last_name,
            email: address.email,
            phone: address.phone,
            street: address.street,
            suburb: address.suburb || '',
            city_dane_code: address.city_dane_code,
            city: address.city,
            state: address.state,
            postal_code: address.postal_code || '',
        });
        setCitySearch(address.city);
        setIsFormOpen(true);
    };

    const handleDelete = async (id: number) => {
        if (window.confirm('¿Estás seguro de eliminar esta dirección?')) {
            const result = await deleteOriginAddressAction(id);
            if (result.success) {
                showToast(result.message, 'success');
                loadAddresses();
            } else {
                showToast(result.message, 'error');
            }
        }
    };

    const handleSetDefault = async (id: number) => {
        const result = await updateOriginAddressAction(id, { is_default: true });
        if (result.success) {
            showToast('Dirección predeterminada actualizada', 'success');
            loadAddresses();
        } else {
            showToast(result.message, 'error');
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h2 className="text-2xl font-bold text-gray-800">Direcciones de Origen</h2>
                    <p className="text-gray-500 text-sm">Gestiona los lugares desde donde envías tus productos.</p>
                </div>
                {!isFormOpen && canCreate && (
                    <Button onClick={() => setIsFormOpen(true)} leftIcon={<Plus className="w-4 h-4" />}>
                        Nueva Dirección
                    </Button>
                )}
            </div>

            {isFormOpen && (
                <div className="bg-white border-2 border-orange-100 rounded-xl overflow-hidden shadow-sm transition-all duration-300">
                    <div className="bg-orange-50 px-6 py-4 flex justify-between items-center border-b border-orange-100">
                        <div>
                            <h3 className="font-bold text-gray-800">{editingAddress ? 'Editar Dirección' : 'Nueva Dirección de Origen'}</h3>
                            <p className="text-xs text-orange-700">Esta dirección podrá ser seleccionada al generar guías.</p>
                        </div>
                        <button onClick={() => setIsFormOpen(false)} className="text-gray-400 hover:text-gray-600 transition-colors">
                            <X className="w-5 h-5" />
                        </button>
                    </div>
                    <div className="px-6 py-6">
                        <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                            <div className="space-y-2">
                                <Label>Alias de la dirección</Label>
                                <Input
                                    value={formData.alias}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, alias: e.target.value })}
                                    placeholder="Ej: Bodega Central, Oficina Norte"
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label>Empresa</Label>
                                <Input
                                    value={formData.company}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, company: e.target.value })}
                                    placeholder="Nombre de la empresa remitente"
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label>Nombre de contacto</Label>
                                <Input
                                    value={formData.first_name}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, first_name: e.target.value })}
                                    placeholder="Nombre"
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label>Apellido de contacto</Label>
                                <Input
                                    value={formData.last_name}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, last_name: e.target.value })}
                                    placeholder="Apellido"
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label>Correo electrónico</Label>
                                <Input
                                    type="email"
                                    value={formData.email}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, email: e.target.value })}
                                    placeholder="email@ejemplo.com"
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label>Teléfono</Label>
                                <Input
                                    value={formData.phone}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, phone: e.target.value })}
                                    placeholder="3001234567"
                                    required
                                />
                            </div>
                            <div className="space-y-2 col-span-full">
                                <Label>Dirección (Calle, Carrera, No.)</Label>
                                <Input
                                    value={formData.street}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, street: e.target.value })}
                                    placeholder="Ej: Calle 100 # 20 - 30, Bodega 5"
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label>Colonia / Barrio</Label>
                                <Input
                                    value={formData.suburb}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, suburb: e.target.value })}
                                    placeholder="Nombre del barrio"
                                />
                            </div>
                            <div className="space-y-2 relative">
                                <Label>Ciudad / Municipio</Label>
                                <Input
                                    value={citySearch}
                                    onChange={handleCitySearch}
                                    placeholder="Escribe para buscar..."
                                    required
                                />
                                {citySuggestions.length > 0 && (
                                    <div className="absolute z-10 w-full bg-white border border-gray-200 rounded-lg shadow-xl mt-1 max-h-48 overflow-auto animate-in fade-in slide-in-from-top-2">
                                        {citySuggestions.map(s => (
                                            <div
                                                key={s.code}
                                                className="p-3 hover:bg-orange-50 cursor-pointer text-sm border-b last:border-0"
                                                onClick={() => selectCity(s)}
                                            >
                                                <div className="font-semibold text-gray-800">{s.ciudad}</div>
                                                <div className="text-xs text-gray-500">{s.departamento} - {s.code}</div>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>

                            <div className="flex justify-end col-span-full gap-3 mt-6 pt-4 border-t border-gray-100">
                                <Button type="button" variant="outline" onClick={() => setIsFormOpen(false)}>Cancelar</Button>
                                <Button type="submit" loading={isPending}>
                                    {editingAddress ? 'Actualizar Dirección' : 'Guardar Dirección'}
                                </Button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {addresses.map(address => (
                    <div
                        key={address.id}
                        className={`group relative bg-white border-2 rounded-xl p-5 transition-all duration-300 ${address.is_default ? 'border-orange-400 shadow-md ring-1 ring-orange-100' : 'border-gray-100 hover:border-orange-200 hover:shadow-sm'
                            }`}
                    >
                        <div className="flex justify-between items-start mb-4">
                            <div className="space-y-1">
                                <div className="flex items-center gap-2">
                                    <h3 className="font-bold text-lg text-gray-800">{address.alias}</h3>
                                    {address.is_default && (
                                        <span className="bg-orange-600 text-white text-[9px] px-2 py-0.5 rounded-full font-bold tracking-wider uppercase">
                                            PRINCIPAL
                                        </span>
                                    )}
                                </div>
                                <div className="flex items-center gap-1.5 text-xs text-gray-500 font-medium">
                                    <Building className="w-3.5 h-3.5" /> {address.company}
                                </div>
                            </div>
                            <div className="flex gap-1">
                                {canUpdate && (
                                    <button
                                        onClick={() => handleEdit(address)}
                                        className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
                                        title="Editar"
                                    >
                                        <Edit2 className="w-4 h-4" />
                                    </button>
                                )}
                                {canDelete && !address.is_default && (
                                    <button
                                        onClick={() => handleDelete(address.id)}
                                        className="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                                        title="Eliminar"
                                    >
                                        <Trash2 className="w-4 h-4" />
                                    </button>
                                )}
                            </div>
                        </div>

                        <div className="space-y-3">
                            <div className="flex items-start gap-2.5">
                                <div className="mt-1 bg-gray-50 p-1 rounded">
                                    <MapPin className="w-4 h-4 text-gray-400" />
                                </div>
                                <div className="text-sm">
                                    <p className="font-semibold text-gray-800 leading-tight">{address.street}</p>
                                    <p className="text-gray-500">{address.city}, {address.state}</p>
                                    {address.suburb && <p className="text-xs text-gray-400 mt-0.5">{address.suburb}</p>}
                                </div>
                            </div>

                            <div className="grid grid-cols-1 gap-2 border-t border-gray-50 pt-3">
                                <div className="flex items-center gap-2.5 text-sm text-gray-600">
                                    <Phone className="w-3.5 h-3.5 text-gray-400" />
                                    {address.phone}
                                </div>
                                <div className="flex items-center gap-2.5 text-sm text-gray-600">
                                    <Mail className="w-3.5 h-3.5 text-gray-400" />
                                    <span className="truncate" title={address.email}>{address.email}</span>
                                </div>
                                <div className="flex items-center gap-2.5 text-[13px] text-gray-700 font-medium">
                                    <CheckCircle className="w-3.5 h-3.5 text-gray-400" />
                                    {address.first_name} {address.last_name}
                                </div>
                            </div>
                        </div>

                        {canUpdate && !address.is_default && (
                            <button
                                onClick={() => handleSetDefault(address.id)}
                                className="w-full mt-5 py-2 text-[11px] font-bold border border-gray-200 rounded-lg text-gray-500 hover:border-orange-300 hover:text-orange-600 hover:bg-orange-50 transition-all uppercase tracking-wide"
                            >
                                Establecer como principal
                            </button>
                        )}
                    </div>
                ))}

                {addresses.length === 0 && !isFormOpen && (
                    <div className="col-span-full py-20 bg-gray-50/50 border-2 border-dashed border-gray-200 rounded-2xl flex flex-col items-center justify-center text-center">
                        <div className="bg-white p-4 rounded-full shadow-sm mb-4">
                            <MapPin className="w-10 h-10 text-gray-300" />
                        </div>
                        <h4 className="text-gray-900 font-bold mb-1">Sin direcciones guardadas</h4>
                        <p className="text-gray-500 text-sm max-w-xs mb-6">Configura tus bodegas o domicilios para generar guías con un solo clic.</p>
                        {canCreate && (
                            <Button onClick={() => setIsFormOpen(true)} leftIcon={<Plus className="w-4 h-4" />}>
                                Agregar mi primera dirección
                            </Button>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}
