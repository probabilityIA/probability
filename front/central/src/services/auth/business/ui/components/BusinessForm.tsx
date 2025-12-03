import React from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Select } from '@/shared/ui/select';
import { FileInput } from '@/shared/ui/file-input';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { Business, BusinessType } from '../../domain/types';
import { useBusinessForm } from '../hooks/useBusinessForm';

interface BusinessFormProps {
    initialData?: Business;
    onSuccess: () => void;
    onCancel: () => void;
    businessTypes: BusinessType[];
}

export const BusinessForm: React.FC<BusinessFormProps> = ({ initialData, onSuccess, onCancel, businessTypes }) => {
    const {
        formData,
        loading,
        error,
        handleChange,
        submit,
        setError
    } = useBusinessForm(initialData, onSuccess);

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

            <Select
                label="Tipo de Negocio"
                value={String(formData.business_type_id || '')}
                onChange={(e: React.ChangeEvent<HTMLSelectElement>) => handleChange('business_type_id', Number(e.target.value))}
                options={businessTypes.map(t => ({ label: t.name, value: String(t.id) }))}
                required
            />

            <Input
                label="Descripción"
                value={formData.description || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('description', e.target.value)}
            />

            <Input
                label="Dirección"
                value={formData.address || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('address', e.target.value)}
            />

            <div className="grid grid-cols-2 gap-4">
                <div className="flex flex-col gap-1">
                    <label className="text-sm font-medium text-gray-700">Color Primario</label>
                    <div className="flex items-center gap-2">
                        <input
                            type="color"
                            value={formData.primary_color || '#000000'}
                            onChange={(e) => handleChange('primary_color', e.target.value)}
                            className="h-10 w-20 p-1 rounded border cursor-pointer"
                        />
                        <span className="text-sm text-gray-500">{formData.primary_color || '#000000'}</span>
                    </div>
                </div>
                <div className="flex flex-col gap-1">
                    <label className="text-sm font-medium text-gray-700">Color Secundario</label>
                    <div className="flex items-center gap-2">
                        <input
                            type="color"
                            value={formData.secondary_color || '#ffffff'}
                            onChange={(e) => handleChange('secondary_color', e.target.value)}
                            className="h-10 w-20 p-1 rounded border cursor-pointer"
                        />
                        <span className="text-sm text-gray-500">{formData.secondary_color || '#ffffff'}</span>
                    </div>
                </div>
            </div>

            {/* 
            <div className="grid grid-cols-2 gap-4 hidden">
                <FileInput
                    label="Logo"
                    onChange={(file: File | null) => handleChange('logo_file', file)}
                />
                <FileInput
                    label="Imagen de Navbar"
                    onChange={(file: File | null) => handleChange('navbar_image_file', file)}
                />
            </div>
            */}

            <div className="flex justify-end gap-2 mt-6">
                <Button type="button" variant="secondary" onClick={onCancel}>Cancelar</Button>
                <Button type="submit" disabled={loading}>
                    {loading ? <Spinner size="sm" /> : (initialData ? 'Actualizar' : 'Crear')}
                </Button>
            </div>
        </form>
    );
};
