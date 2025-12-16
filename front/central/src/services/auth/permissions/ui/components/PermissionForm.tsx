'use client';

import React from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Select } from '@/shared/ui/select';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { Permission } from '../../domain/types';
import { usePermissionForm, SCOPE_OPTIONS } from '../hooks/usePermissionForm';

interface PermissionFormProps {
    initialData?: Permission;
    onSuccess: () => void;
    onCancel: () => void;
}

export const PermissionForm: React.FC<PermissionFormProps> = ({ initialData, onSuccess, onCancel }) => {
    const {
        formData,
        loading,
        loadingData,
        error,
        resources,
        actions,
        handleChange,
        submit,
        setError
    } = usePermissionForm(initialData, onSuccess);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        await submit();
    };

    if (loadingData) {
        return (
            <div className="flex justify-center items-center py-8">
                <Spinner size="lg" />
                <span className="ml-2 text-gray-600">Cargando datos...</span>
            </div>
        );
    }

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <Input
                label="Nombre *"
                value={formData.name || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('name', e.target.value)}
                placeholder="Nombre del permiso"
                required
            />

            <Input
                label="Descripción"
                value={formData.description || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('description', e.target.value)}
                placeholder="Descripción del permiso"
            />

            <div className="grid grid-cols-2 gap-4">
                <Select
                    label="Recurso *"
                    value={formData.resource_id || ''}
                    onChange={(e) => handleChange('resource_id', e.target.value ? Number(e.target.value) : null)}
                    options={resources.map(r => ({ value: r.id, label: r.name }))}
                    placeholder="Seleccionar recurso..."
                    required
                />

                <Select
                    label="Acción *"
                    value={formData.action_id || ''}
                    onChange={(e) => handleChange('action_id', e.target.value ? Number(e.target.value) : null)}
                    options={actions.map(a => ({ value: a.id, label: a.name }))}
                    placeholder="Seleccionar acción..."
                    required
                />
            </div>

            <Select
                label="Scope *"
                value={formData.scope_id || ''}
                onChange={(e) => handleChange('scope_id', e.target.value ? Number(e.target.value) : null)}
                options={SCOPE_OPTIONS.map(s => ({ value: s.id, label: `${s.name} (${s.code})` }))}
                placeholder="Seleccionar scope..."
                required
                helperText="Platform: permisos para super administradores. Business: permisos para usuarios de negocio."
            />

            <div className="flex justify-end gap-2 mt-6">
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
