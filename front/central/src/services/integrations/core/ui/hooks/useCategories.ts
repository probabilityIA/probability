'use client';

import { useState, useEffect } from 'react';
import { getIntegrationCategoriesAction } from '../../infra/actions';
import { IntegrationCategory } from '../../domain/types';
import { getActionError } from '@/shared/utils/action-result';

export function useCategories() {
    const [categories, setCategories] = useState<IntegrationCategory[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        fetchCategories();
    }, []);

    const fetchCategories = async () => {
        setLoading(true);
        try {
            const response = await getIntegrationCategoriesAction();
            if (response.success && response.data) {
                setCategories(response.data);
                setError(null);
            } else {
                setError(response.message || 'Error al obtener categorías');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error desconocido'));
            console.error('Error fetching categories:', err);
        } finally {
            setLoading(false);
        }
    };

    return {
        categories,
        loading,
        error,
        refresh: fetchCategories
    };
}
