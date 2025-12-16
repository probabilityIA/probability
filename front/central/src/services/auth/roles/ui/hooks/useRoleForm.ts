import { useState, useEffect } from 'react';
import { createRoleAction, updateRoleAction } from '../../infra/actions';
import { Role, CreateRoleDTO, UpdateRoleDTO } from '../../domain/types';

export const useRoleForm = (initialData?: Role, onSuccess?: () => void) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [formData, setFormData] = useState<Partial<CreateRoleDTO>>({
        name: '',
        description: '',
        level: undefined,
        is_system: false,
        scope_id: undefined,
        business_type_id: 1, // Siempre es 1
    });

    useEffect(() => {
        if (initialData) {
            setFormData({
                name: initialData.name,
                description: initialData.description,
                level: initialData.level,
                is_system: initialData.is_system,
                scope_id: initialData.scope_id,
                business_type_id: initialData.business_type_id,
            });
        }
    }, [initialData]);

    const handleChange = (field: keyof CreateRoleDTO, value: string | number | boolean | null) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const submit = async () => {
        setLoading(true);
        setError(null);

        try {
            // Asegurar que business_type_id siempre sea 1
            const dataToSubmit = {
                ...formData,
                business_type_id: 1,
            };

            if (initialData) {
                await updateRoleAction(initialData.id, dataToSubmit as UpdateRoleDTO);
            } else {
                await createRoleAction(dataToSubmit as CreateRoleDTO);
            }
            if (onSuccess) onSuccess();
            return true;
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error saving role';
            setError(errorMessage);
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
        setError
    };
};
