import React from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Select } from '@/shared/ui/select';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { Role } from '../../domain/types';
import { useRoleForm } from '../hooks/useRoleForm';

interface RoleFormProps {
    initialData?: Role;
    onSuccess: () => void;
    onCancel: () => void;
}

export const RoleForm: React.FC<RoleFormProps> = ({ initialData, onSuccess, onCancel }) => {
    const {
        formData,
        loading,
        error,
        handleChange,
        submit,
        setError
    } = useRoleForm(initialData, onSuccess);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        await submit();
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <Input
                label="Nombre"
                value={formData.name || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('name', e.target.value)}
                required
            />
            <Input
                label="Descripción"
                value={formData.description || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('description', e.target.value)}
                required
            />

            <div className="grid grid-cols-2 gap-4">
                <Input
                    label="Nivel (1-10)"
                    type="number"
                    min={1}
                    max={10}
                    value={formData.level || ''}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('level', Number(e.target.value))}
                    required
                />
                <label className="flex items-center gap-2 mt-8">
                    <input
                        type="checkbox"
                        checked={formData.is_system}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('is_system', e.target.checked)}
                    />
                    Rol del Sistema
                </label>
            </div>

            <Select
                label="Scope"
                value={formData.scope_id || ''}
                onChange={(e: React.ChangeEvent<HTMLSelectElement>) => handleChange('scope_id', Number(e.target.value))}
                options={[
                    { value: 1, label: 'Platform' },
                    { value: 2, label: 'Business' },
                ]}
                required
            />

            <div className="flex justify-end gap-2 mt-6">
                <Button type="button" variant="secondary" onClick={onCancel}>Cancelar</Button>
                <Button type="submit" disabled={loading}>
                    {loading ? <Spinner size="sm" /> : (initialData ? 'Actualizar' : 'Crear')}
                </Button>
            </div>
        </form>
    );
};
