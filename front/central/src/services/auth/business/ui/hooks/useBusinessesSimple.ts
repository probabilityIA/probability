import { useState, useEffect } from 'react';
import { getBusinessesSimpleAction } from '../../infra/actions';
import { BusinessSimple } from '../../domain/types';

/**
 * Hook optimizado para obtener lista simple de businesses (solo id y name)
 * Ideal para dropdowns, selectores y otros componentes que no necesitan todos los datos
 */
export const useBusinessesSimple = () => {
    const [businesses, setBusinesses] = useState<BusinessSimple[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchBusinesses = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getBusinessesSimpleAction();
            if (response.success) {
                setBusinesses(response.data);
            } else {
                setError(response.message);
            }
        } catch (err: any) {
            setError(err.message || 'Error al obtener negocios');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchBusinesses();
    }, []);

    return {
        businesses,
        loading,
        error,
        refresh: fetchBusinesses,
    };
};
