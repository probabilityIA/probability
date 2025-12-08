import React, { useEffect, useState } from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { AvatarUpload } from '@/shared/ui/avatar-upload';
import { User } from '../../domain/types';
import { useUserForm } from '../hooks/useUserForm';

interface UserFormProps {
    initialData?: User;
    onSuccess: () => void;
    onCancel: () => void;
}

export const UserForm: React.FC<UserFormProps> = ({ initialData, onSuccess, onCancel }) => {
    const {
        formData,
        loading,
        error,
        successMessage,
        handleChange,
        handleFileChange,
        submit,
        setError
    } = useUserForm(initialData, onSuccess);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        await submit();
        // If success and no password message to show (e.g. update), or if user acknowledges password, we close.
        // But here we rely on successMessage state.
        // If it's an update, successMessage is null and we called onSuccess in hook.
        // If it's create, successMessage is set.
    };

    if (successMessage) {
        return (
            <div className="space-y-4">
                <Alert type="success">{successMessage}</Alert>
                <div className="flex justify-end">
                    <Button onClick={onSuccess}>Close</Button>
                </div>
            </div>
        );
    }

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
                label="Email"
                type="email"
                value={formData.email || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('email', e.target.value)}
                required
            />
            <Input
                label="Phone"
                value={formData.phone || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('phone', e.target.value)}
            />

            <div className="flex flex-col items-center gap-4 py-4">
                <AvatarUpload
                    currentAvatarUrl={initialData?.avatar_url || null}
                    onFileSelect={handleFileChange}
                    size="lg"
                />
            </div>

            <div className="flex items-center gap-2 mt-4">
                <input
                    type="checkbox"
                    checked={formData.is_active}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('is_active', e.target.checked)}
                    id="is_active"
                />
                <label htmlFor="is_active">Active</label>
            </div>

            <div className="mt-4">
                <label className="block text-sm font-medium text-gray-700 mb-1">Business IDs (comma separated)</label>
                <Input
                    placeholder="e.g. 1, 2, 3"
                    value={formData.business_ids?.join(', ') || ''}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                        const val = e.target.value;
                        const ids = val.split(',').map(s => s.trim()).filter(s => s !== '').map(Number).filter(n => !isNaN(n));
                        handleChange('business_ids', ids);
                    }}
                />
                <p className="text-xs text-gray-500 mt-1">Enter Business IDs separated by commas.</p>
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
