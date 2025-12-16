'use client';

import React from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { Resource } from '../../domain/types';
import { useResourceForm } from '../hooks/useResourceForm';

interface ResourceFormProps {
    initialData?: Resource;
    onSuccess: () => void;
    onCancel: () => void;
}

export const ResourceForm: React.FC<ResourceFormProps> = ({ initialData, onSuccess, onCancel }) => {
    const {
        formData,
        loading,
        error,
        handleChange,
        submit,
        setError
    } = useResourceForm(initialData, onSuccess);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        await submit();
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <Input
                label="Nombre *"
                value={formData.name || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('name', e.target.value)}
                placeholder="Nombre del recurso"
                required
            />

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Descripción
                </label>
                <textarea
                    value={formData.description || ''}
                    onChange={(e) => handleChange('description', e.target.value)}
                    placeholder="Descripción del recurso"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    rows={3}
                />
            </div>

            <p className="text-xs text-gray-500">
                Los recursos representan entidades del sistema sobre las cuales se pueden definir permisos (ej: orders, users, products).
            </p>

            <div className="flex justify-end gap-2 pt-4">
                <Button type="button" variant="secondary" onClick={onCancel}>
                    Cancelar
                </Button>
                <Button type="submit" disabled={loading}>
                    {loading ? <Spinner size="sm" /> : (initialData ? 'Actualizar' : 'Crear')}
                </Button>
            </div>
        </form>
    );
};
