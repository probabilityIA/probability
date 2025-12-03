import React from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';

import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { Permission } from '../../domain/types';
import { usePermissionForm } from '../hooks/usePermissionForm';

interface PermissionFormProps {
    initialData?: Permission;
    onSuccess: () => void;
    onCancel: () => void;
}

export const PermissionForm: React.FC<PermissionFormProps> = ({ initialData, onSuccess, onCancel }) => {
    const {
        formData,
        loading,
        error,
        handleChange,
        submit,
        setError
    } = usePermissionForm(initialData, onSuccess);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        await submit();
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <Input
                label="Name"
                value={formData.name || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('name', e.target.value)}
                required
            />
            <Input
                label="Code"
                value={formData.code || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('code', e.target.value)}
                placeholder="Auto-generated if empty"
            />
            <Input
                label="Description"
                value={formData.description || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('description', e.target.value)}
            />

            <div className="grid grid-cols-2 gap-4">
                <Input
                    label="Resource ID"
                    type="number"
                    value={formData.resource_id || ''}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('resource_id', Number(e.target.value))}
                    required
                />
                <Input
                    label="Action ID"
                    type="number"
                    value={formData.action_id || ''}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('action_id', Number(e.target.value))}
                    required
                />
            </div>

            <div className="grid grid-cols-2 gap-4">
                <Input
                    label="Scope ID"
                    type="number"
                    value={formData.scope_id || ''}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('scope_id', Number(e.target.value))}
                    required
                />
                <Input
                    label="Business Type ID"
                    type="number"
                    value={formData.business_type_id || ''}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('business_type_id', e.target.value ? Number(e.target.value) : null)}
                    placeholder="Leave empty for generic"
                />
            </div>

            <div className="flex justify-end gap-2 mt-6">
                <Button type="button" variant="secondary" onClick={onCancel}>Cancel</Button>
                <Button type="submit" disabled={loading}>
                    {loading ? <Spinner size="sm" /> : (initialData ? 'Update' : 'Create')}
                </Button>
            </div>
        </form>
    );
};
