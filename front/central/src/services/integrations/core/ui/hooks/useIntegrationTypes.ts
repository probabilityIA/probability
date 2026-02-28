'use client';

import { useState, useEffect, useCallback } from 'react';
import {
    getIntegrationTypesAction,
    createIntegrationTypeAction,
    updateIntegrationTypeAction,
    deleteIntegrationTypeAction
} from '../../infra/actions';
import { IntegrationType, CreateIntegrationTypeDTO, UpdateIntegrationTypeDTO } from '../../domain/types';

export const useIntegrationTypes = (categoryId?: number) => {
    const [integrationTypes, setIntegrationTypes] = useState<IntegrationType[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchIntegrationTypes = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getIntegrationTypesAction(categoryId);
            if (response.success) {
                setIntegrationTypes(response.data);
            } else {
                setError(response.message);
            }
        } catch (err: any) {
            console.error('Error fetching integration types:', err);
            setError(err.message || 'Error fetching integration types');
        } finally {
            setLoading(false);
        }
    }, [categoryId]);

    const createIntegrationType = async (data: CreateIntegrationTypeDTO) => {
        try {
            const response = await createIntegrationTypeAction(data);
            if (response.success) {
                fetchIntegrationTypes();
                return true;
            } else {
                setError(response.message);
                return false;
            }
        } catch (err: any) {
            console.error('Error creating integration type:', err);
            setError(err.message || 'Error creating integration type');
            return false;
        }
    };

    const updateIntegrationType = async (id: number, data: UpdateIntegrationTypeDTO) => {
        try {
            const response = await updateIntegrationTypeAction(id, data);
            if (response.success) {
                fetchIntegrationTypes();
                return true;
            } else {
                setError(response.message);
                return false;
            }
        } catch (err: any) {
            console.error('Error updating integration type:', err);
            setError(err.message || 'Error updating integration type');
            return false;
        }
    };

    const deleteIntegrationType = async (id: number) => {
        try {
            const response = await deleteIntegrationTypeAction(id);
            if (response.success) {
                fetchIntegrationTypes();
                return true;
            } else {
                setError(response.message);
                return false;
            }
        } catch (err: any) {
            console.error('Error deleting integration type:', err);
            setError(err.message || 'Error deleting integration type');
            return false;
        }
    };

    useEffect(() => {
        fetchIntegrationTypes();
    }, [fetchIntegrationTypes]);

    return {
        integrationTypes,
        loading,
        error,
        setError,
        createIntegrationType,
        updateIntegrationType,
        deleteIntegrationType,
        refresh: fetchIntegrationTypes
    };
};
