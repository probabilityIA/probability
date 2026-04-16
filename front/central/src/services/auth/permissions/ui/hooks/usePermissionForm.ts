'use client';

import { useState, useEffect } from 'react';
import { createPermissionAction, updatePermissionAction } from '../../infra/actions';
import { Permission, CreatePermissionDTO, UpdatePermissionDTO } from '../../domain/types';
import { getResourcesAction } from '@/services/auth/resources/infra/actions';
import { getActionsAction } from '@/services/auth/actions/infra/actions';
import { Resource } from '@/services/auth/resources/domain/types';
import { Action } from '@/services/auth/actions/domain/types';

// Scopes predefinidos del sistema
export const SCOPE_OPTIONS = [
    { id: 1, name: 'Platform', code: 'platform' },
    { id: 2, name: 'Business', code: 'business' },
];

export const usePermissionForm = (initialData?: Permission, onSuccess?: () => void) => {
    const [loading, setLoading] = useState(false);
    const [loadingData, setLoadingData] = useState(true);
    const [error, setError] = useState<string | null>(null);
    
    // Datos para los dropdowns
    const [resources, setResources] = useState<Resource[]>([]);
    const [actions, setActions] = useState<Action[]>([]);

    const [formData, setFormData] = useState<Partial<CreatePermissionDTO>>({
        name: '',
        code: '',
        description: '',
        resource_id: undefined,
        action_id: undefined,
        scope_id: undefined,
        business_type_id: 1, // Siempre 1 por defecto
    });

    // Cargar recursos y acciones al montar
    useEffect(() => {
        const loadData = async () => {
            setLoadingData(true);
            try {
                const [resourcesRes, actionsRes] = await Promise.all([
                    getResourcesAction({ page_size: 100 }),
                    getActionsAction({ page_size: 100 })
                ]);

                if (resourcesRes.success && resourcesRes.data?.resources) {
                    setResources(resourcesRes.data.resources);
                }
                if (actionsRes.success && actionsRes.data?.actions) {
                    setActions(actionsRes.data.actions);
                }
            } catch (err) {
                console.error('Error loading form data:', err);
            } finally {
                setLoadingData(false);
            }
        };

        loadData();
    }, []);

    useEffect(() => {
        if (initialData) {
            setFormData({
                name: initialData.name,
                code: initialData.code,
                description: initialData.description,
                resource_id: initialData.resource_id,
                action_id: initialData.action_id,
                scope_id: initialData.scope_id,
                business_type_id: 1, // Siempre 1
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
            // Asegurar que business_type_id siempre sea 1
            const dataToSubmit = { ...formData, business_type_id: 1 };
            
            if (initialData) {
                await updatePermissionAction(initialData.id, dataToSubmit as UpdatePermissionDTO);
            } else {
                await createPermissionAction(dataToSubmit as CreatePermissionDTO);
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
        loadingData,
        error,
        resources,
        actions,
        handleChange,
        submit,
        setError
    };
};
