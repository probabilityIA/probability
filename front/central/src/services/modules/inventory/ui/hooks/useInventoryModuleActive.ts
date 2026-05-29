'use client';

import { useEffect, useState } from 'react';
import { getBusinessConfiguredResourcesAction } from '@/services/auth/business/infra/actions';

interface UseInventoryModuleActiveResult {
    isActive: boolean;
    loading: boolean;
}

export function useInventoryModuleActive(businessId?: number | null): UseInventoryModuleActiveResult {
    const [isActive, setIsActive] = useState(true);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!businessId) {
            setLoading(false);
            setIsActive(true);
            return;
        }

        let cancelled = false;
        setLoading(true);
        getBusinessConfiguredResourcesAction(businessId)
            .then(res => {
                if (cancelled) return;
                if (res?.success && res.data) {
                    const inventory = res.data.resources?.find(r => r.resource_name === 'Inventario');
                    setIsActive(inventory?.is_active === true);
                } else {
                    setIsActive(false);
                }
            })
            .catch(() => {
                if (!cancelled) setIsActive(false);
            })
            .finally(() => {
                if (!cancelled) setLoading(false);
            });

        return () => {
            cancelled = true;
        };
    }, [businessId]);

    return { isActive, loading };
}
