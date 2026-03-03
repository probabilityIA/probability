'use client';

import { useState, useEffect } from 'react';
import { Warehouse } from '../../../warehouses/domain/types';
import { getWarehousesAction } from '../../../warehouses/infra/actions';

interface UseWarehousesOptions {
    businessId: number;
}

interface UseWarehousesResult {
    warehouses: Warehouse[];
    loading: boolean;
}

export function useWarehouses({ businessId }: UseWarehousesOptions): UseWarehousesResult {
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (!businessId) return;

        setLoading(true);
        getWarehousesAction({
            business_id: businessId,
            is_active: true,
            page: 1,
            page_size: 100,
        })
            .then((res) => {
                setWarehouses(res.data || []);
            })
            .catch(() => {
                setWarehouses([]);
            })
            .finally(() => {
                setLoading(false);
            });
    }, [businessId]);

    return { warehouses, loading };
}
