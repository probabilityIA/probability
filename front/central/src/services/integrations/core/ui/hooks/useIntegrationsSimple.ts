import { useState, useEffect } from 'react';
import { getIntegrationsSimpleAction } from '../../infra/actions';
import { IntegrationSimple } from '../../domain/types';

interface UseIntegrationsSimpleOptions {
    businessId?: number;
}

/**
 * Hook optimizado para obtener lista simple de integraciones (solo id, name, type, business_id, is_active)
 * Ideal para dropdowns, selectores y otros componentes que no necesitan todos los datos
 */
export const useIntegrationsSimple = (options?: UseIntegrationsSimpleOptions) => {
    const [integrations, setIntegrations] = useState<IntegrationSimple[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchIntegrations = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getIntegrationsSimpleAction(options?.businessId);
            if (response.success) {
                setIntegrations(response.data);
            } else {
                setError(response.message);
            }
        } catch (err: any) {
            setError(err.message || 'Error al obtener integraciones');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchIntegrations();
<<<<<<< HEAD
        // eslint-disable-next-line react-hooks/exhaustive-deps
=======
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
    }, [options?.businessId]);

    return {
        integrations,
        loading,
        error,
        refresh: fetchIntegrations,
    };
};
