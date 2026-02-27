'use client';

import { useState } from 'react';
import { CustomerInfo, CreateCustomerDTO, UpdateCustomerDTO } from '../../domain/types';
import { createCustomerAction, updateCustomerAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';

interface CustomerFormProps {
    customer?: CustomerInfo;
    onSuccess: () => void;
    onCancel: () => void;
}

export default function CustomerForm({ customer, onSuccess, onCancel }: CustomerFormProps) {
    const [formData, setFormData] = useState<CreateCustomerDTO>({
        name: customer?.name || '',
        email: customer?.email || '',
        phone: customer?.phone || '',
        dni: customer?.dni || '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const handleChange = (field: keyof CreateCustomerDTO, value: string) => {
        setFormData((prev) => ({ ...prev, [field]: value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            if (customer) {
                const updateData: UpdateCustomerDTO = {
                    name: formData.name,
                    email: formData.email || undefined,
                    phone: formData.phone || undefined,
                    dni: formData.dni || null,
                };
                await updateCustomerAction(customer.id, updateData);
            } else {
                const createData: CreateCustomerDTO = {
                    name: formData.name,
                    email: formData.email || undefined,
                    phone: formData.phone || undefined,
                    dni: formData.dni || null,
                };
                await createCustomerAction(createData);
            }

            setSuccess(customer ? 'Cliente actualizado exitosamente' : 'Cliente creado exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(err.message || 'Error al guardar el cliente');
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
                <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Nombre <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.name}
                        onChange={(e) => handleChange('name', e.target.value)}
                        placeholder="Nombre completo"
                        required
                        minLength={2}
                        maxLength={255}
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

                {/* Teléfono */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Teléfono
                    </label>
                    <Input
                        type="tel"
                        value={formData.phone}
                        onChange={(e) => handleChange('phone', e.target.value)}
                        placeholder="3001234567"
                        maxLength={20}
                    />
                </div>

                {/* Documento (DNI) */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Documento de identidad
                    </label>
                    <Input
                        type="text"
                        value={formData.dni as string}
                        onChange={(e) => handleChange('dni', e.target.value)}
                        placeholder="CC, NIT, pasaporte..."
                        maxLength={30}
                    />
                </div>
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t">
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
