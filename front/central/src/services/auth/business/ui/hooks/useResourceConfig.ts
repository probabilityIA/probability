import { useState, useEffect, useCallback } from 'react';
import {
    getBusinessConfiguredResourcesAction,
    activateResourceAction,
    deactivateResourceAction
} from '../../infra/actions';
import { BusinessConfiguredResources, ConfiguredResource } from '../../domain/types';

export const useResourceConfig = (businessId: number) => {
    const [config, setConfig] = useState<BusinessConfiguredResources | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [actionLoading, setActionLoading] = useState<number | null>(null);

    const fetchConfig = useCallback(async () => {
        setLoading(true);
        try {
            const res = await getBusinessConfiguredResourcesAction(businessId);
            setConfig(res.data);
        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => {
        if (businessId) fetchConfig();
    }, [businessId, fetchConfig]);

    const toggleResource = async (resource: ConfiguredResource) => {
        setActionLoading(resource.resource_id);
        try {
            if (resource.is_active) {
                await deactivateResourceAction(resource.resource_id, businessId);
            } else {
                await activateResourceAction(resource.resource_id, businessId);
            }
            await fetchConfig();
        } catch (err: any) {
            setError(err.message);
        } finally {
            setActionLoading(null);
        }
    };

    return {
        config,
        loading,
        error,
        actionLoading,
        toggleResource,
        setError
    };
};
