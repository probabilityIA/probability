'use client';

import { useState } from 'react';
import { UserCircleIcon, EnvelopeIcon, PhoneIcon, IdentificationIcon, MapPinIcon, DocumentTextIcon } from '@heroicons/react/24/outline';
import { CustomerInfo, CreateCustomerDTO, UpdateCustomerDTO } from '../../domain/types';
import { createCustomerAction, updateCustomerAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface CustomerFormProps {
    customer?: CustomerInfo;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

interface FormState {
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    dni: string;
    address: string;
    city: string;
    notes: string;
}

function splitName(fullName: string): { firstName: string; lastName: string } {
    const parts = fullName.trim().split(/\s+/);
    if (parts.length <= 1) {
        return { firstName: parts[0] || '', lastName: '' };
    }
    return { firstName: parts[0], lastName: parts.slice(1).join(' ') };
}

function FieldLabel({ icon: Icon, children, required }: { icon: React.ComponentType<{ className?: string }>; children: React.ReactNode; required?: boolean }) {
    return (
        <label className="flex items-center gap-1.5 text-sm font-medium text-gray-700 dark:text-gray-200 mb-1.5">
            <Icon className="w-4 h-4 text-gray-400 dark:text-gray-500" />
            {children}
            {required && <span className="text-red-500">*</span>}
        </label>
    );
}

export default function CustomerForm({ customer, onSuccess, onCancel, businessId }: CustomerFormProps) {
    const initialSplit = splitName(customer?.name || '');
    const [formData, setFormData] = useState<FormState>({
        firstName: initialSplit.firstName,
        lastName: initialSplit.lastName,
        email: customer?.email || '',
        phone: customer?.phone || '',
        dni: customer?.dni || '',
        address: customer?.address || '',
        city: customer?.city || '',
        notes: customer?.notes || '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const handleChange = (field: keyof FormState, value: string) => {
        setFormData((prev) => ({ ...prev, [field]: value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        const fullName = [formData.firstName.trim(), formData.lastName.trim()].filter(Boolean).join(' ');

        try {
            if (customer) {
                const updateData: UpdateCustomerDTO = {
                    name: fullName,
                    email: formData.email || undefined,
                    phone: formData.phone || undefined,
                    dni: formData.dni || null,
                    address: formData.address || null,
                    city: formData.city || null,
                    notes: formData.notes || null,
                };
                await updateCustomerAction(customer.id, updateData, businessId);
            } else {
                const createData: CreateCustomerDTO = {
                    name: fullName,
                    email: formData.email || undefined,
                    phone: formData.phone || undefined,
                    dni: formData.dni || null,
                    address: formData.address || null,
                    city: formData.city || null,
                    notes: formData.notes || null,
                };
                await createCustomerAction(createData, businessId);
            }

            setSuccess(customer ? 'Cliente actualizado exitosamente' : 'Cliente creado exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar el cliente'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-7">
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

            <div>
                <h3 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3">
                    Datos personales
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <FieldLabel icon={UserCircleIcon} required>Nombre</FieldLabel>
                        <Input
                            type="text"
                            value={formData.firstName}
                            onChange={(e) => handleChange('firstName', e.target.value)}
                            placeholder="Juan"
                            required
                            minLength={2}
                            maxLength={150}
                        />
                    </div>
                    <div>
                        <FieldLabel icon={UserCircleIcon} required>Apellido</FieldLabel>
                        <Input
                            type="text"
                            value={formData.lastName}
                            onChange={(e) => handleChange('lastName', e.target.value)}
                            placeholder="Perez"
                            required
                            minLength={1}
                            maxLength={150}
                        />
                    </div>
                    <div>
                        <FieldLabel icon={IdentificationIcon} required>Cedula</FieldLabel>
                        <Input
                            type="text"
                            value={formData.dni}
                            onChange={(e) => handleChange('dni', e.target.value)}
                            placeholder="1234567890"
                            required
                            maxLength={30}
                        />
                    </div>
                </div>
            </div>

            <div>
                <h3 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3">
                    Contacto
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <FieldLabel icon={EnvelopeIcon} required>Correo</FieldLabel>
                        <Input
                            type="email"
                            value={formData.email}
                            onChange={(e) => handleChange('email', e.target.value)}
                            placeholder="correo@ejemplo.com"
                            required
                            maxLength={255}
                        />
                    </div>
                    <div>
                        <FieldLabel icon={PhoneIcon} required>Telefono</FieldLabel>
                        <Input
                            type="tel"
                            value={formData.phone}
                            onChange={(e) => handleChange('phone', e.target.value)}
                            placeholder="3001234567"
                            required
                            maxLength={20}
                        />
                    </div>
                </div>
            </div>

            <div>
                <h3 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3">
                    Direccion <span className="normal-case font-normal text-gray-400">(opcional)</span>
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="md:col-span-2">
                        <FieldLabel icon={MapPinIcon}>Direccion</FieldLabel>
                        <Input
                            type="text"
                            value={formData.address}
                            onChange={(e) => handleChange('address', e.target.value)}
                            placeholder="Calle 1 # 2-3"
                            maxLength={255}
                        />
                    </div>
                    <div>
                        <FieldLabel icon={MapPinIcon}>Ciudad</FieldLabel>
                        <Input
                            type="text"
                            value={formData.city}
                            onChange={(e) => handleChange('city', e.target.value)}
                            placeholder="Bogota"
                            maxLength={120}
                        />
                    </div>
                </div>
                <div className="mt-4">
                    <FieldLabel icon={DocumentTextIcon}>Notas</FieldLabel>
                    <textarea
                        value={formData.notes}
                        onChange={(e) => handleChange('notes', e.target.value)}
                        placeholder="Notas internas sobre este cliente..."
                        maxLength={1000}
                        rows={3}
                        className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-sm text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent resize-none"
                    />
                </div>
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t border-gray-100 dark:border-gray-700">
                <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
                    Cancelar
                </Button>
                <Button type="submit" variant="primary" disabled={loading}>
                    {loading ? 'Guardando...' : customer ? 'Actualizar' : 'Crear cliente'}
                </Button>
            </div>
        </form>
    );
}
