import { useState, useEffect } from 'react';
import { createPermissionAction, updatePermissionAction } from '../../infra/actions';
import { Permission, CreatePermissionDTO, UpdatePermissionDTO } from '../../domain/types';

export const usePermissionForm = (initialData?: Permission, onSuccess?: () => void) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [formData, setFormData] = useState<Partial<CreatePermissionDTO>>({
        name: '',
        code: '',
        description: '',
        resource_id: undefined,
        action_id: undefined,
        scope_id: undefined,
        business_type_id: undefined,
    });

    useEffect(() => {
        if (initialData) {
            setFormData({
                name: initialData.name,
                code: initialData.code,
                description: initialData.description,
                resource_id: initialData.resource_id,
                action_id: initialData.action_id,
                scope_id: initialData.scope_id,
                business_type_id: initialData.business_type_id,
            });
        }
    }, [initialData]);

    const handleChange = (field: keyof CreatePermissionDTO, value: string | number | boolean | null) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const submit = async () => {
        setLoading(true);
        setError(null);

        try {
            if (initialData) {
                await updatePermissionAction(initialData.id, formData as UpdatePermissionDTO);
            } else {
                await createPermissionAction(formData as CreatePermissionDTO);
            }
            if (onSuccess) onSuccess();
            return true;
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error saving permission';
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
