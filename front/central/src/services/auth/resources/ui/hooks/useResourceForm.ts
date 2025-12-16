'use client';

import { useState, useEffect } from 'react';
import { createResourceAction, updateResourceAction } from '../../infra/actions';
import { Resource, CreateResourceDTO, UpdateResourceDTO } from '../../domain/types';

export const useResourceForm = (initialData?: Resource, onSuccess?: () => void) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [formData, setFormData] = useState<Partial<CreateResourceDTO>>({
        name: '',
        description: '',
        business_type_id: null,
    });

    useEffect(() => {
        if (initialData) {
            setFormData({
                name: initialData.name,
                description: initialData.description || '',
                business_type_id: initialData.business_type_id || null,
            });
        }
    }, [initialData]);

    const handleChange = (field: keyof CreateResourceDTO, value: string | number | null) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const submit = async () => {
        setLoading(true);
        setError(null);

        try {
            let response;
            if (initialData) {
                response = await updateResourceAction(initialData.id, formData as UpdateResourceDTO);
            } else {
                response = await createResourceAction(formData as CreateResourceDTO);
            }

            if (response.success) {
                if (onSuccess) onSuccess();
                return true;
            } else {
                setError(response.message || 'Error al guardar recurso');
                return false;
            }
        } catch (err: any) {
            setError(err.message || 'Error al guardar recurso');
            return false;
        } finally {
            setLoading(false);
        }
    };

    return {
        formData,
        loading,
        error,
        handleChange,
        submit,
        setError,
    };
};
