import React from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { BusinessType } from '../../domain/types';
import { useBusinessTypeForm } from '../hooks/useBusinessTypes';

interface BusinessTypeFormProps {
    initialData?: BusinessType;
    onSuccess: () => void;
    onCancel: () => void;
}

export const BusinessTypeForm: React.FC<BusinessTypeFormProps> = ({ initialData, onSuccess, onCancel }) => {
    const {
        formData,
        setFormData,
        loading,
        error,
        submit,
        setError
    } = useBusinessTypeForm(initialData, onSuccess);

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
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, name: e.target.value })}
                required
            />
            <Input
                label="Code"
                value={formData.code || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, code: e.target.value })}
                required
            />
            <Input
                label="Description"
                value={formData.description || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, description: e.target.value })}
            />
            <Input
                label="Icon"
                value={formData.icon || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, icon: e.target.value })}
            />

            <label className="flex items-center gap-2">
                <input
                    type="checkbox"
                    checked={formData.is_active}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, is_active: e.target.checked })}
                />
                Active
            </label>

            <div className="flex justify-end gap-2 mt-6">
                <Button type="button" variant="secondary" onClick={onCancel}>Cancel</Button>
                <Button type="submit" disabled={loading}>
                    {loading ? <Spinner size="sm" /> : (initialData ? 'Update' : 'Create')}
                </Button>
            </div>
        </form>
    );
};
