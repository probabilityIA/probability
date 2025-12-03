import { useState, useEffect, useCallback } from 'react';
import { getBusinessTypesAction, deleteBusinessTypeAction, createBusinessTypeAction, updateBusinessTypeAction } from '../../infra/actions';
import { BusinessType, CreateBusinessTypeDTO } from '../../domain/types';

export const useBusinessTypes = () => {
    const [types, setTypes] = useState<BusinessType[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchTypes = useCallback(async () => {
        setLoading(true);
        try {
            const res = await getBusinessTypesAction();
            setTypes(res.data);
        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    }, []);

    const deleteType = async (id: number) => {
        try {
            await deleteBusinessTypeAction(id);
            fetchTypes();
            return true;
        } catch (err: any) {
            setError(err.message);
            return false;
        }
    };

    useEffect(() => {
        fetchTypes();
    }, [fetchTypes]);

    return {
        types,
        loading,
        error,
        deleteType,
        refresh: fetchTypes,
        setError
    };
};

export const useBusinessTypeForm = (initialData?: BusinessType, onSuccess?: () => void) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [formData, setFormData] = useState<Partial<CreateBusinessTypeDTO>>({
        name: '',
        code: '',
        description: '',
        icon: '',
        is_active: true,
    });

    useEffect(() => {
        if (initialData) {
            setFormData({
                name: initialData.name,
                code: initialData.code,
                description: initialData.description,
                icon: initialData.icon,
                is_active: initialData.is_active,
            });
        }
    }, [initialData]);

    const submit = async () => {
        setLoading(true);
        try {
            if (initialData) {
                await updateBusinessTypeAction(initialData.id, formData);
            } else {
                await createBusinessTypeAction(formData as CreateBusinessTypeDTO);
            }
            if (onSuccess) onSuccess();
            return true;
        } catch (err: any) {
            setError(err.message);
            return false;
        } finally {
            setLoading(false);
        }
    };

    return {
        formData,
        setFormData,
        loading,
        error,
        submit,
        setError
    }
}
