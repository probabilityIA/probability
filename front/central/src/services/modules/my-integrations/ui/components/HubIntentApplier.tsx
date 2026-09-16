'use client';

import { useEffect } from 'react';
import { useSyncActivity, type SyncEnvironment } from '../sync-activity-context';

interface HubIntentApplierProps {
    environment: SyncEnvironment;
}

export function HubIntentApplier({ environment }: HubIntentApplierProps) {
    const { setEnvironment, setView } = useSyncActivity();

    useEffect(() => {
        setEnvironment(environment);
        if (environment === 'orders_compare') setView('informe');
    }, [environment, setEnvironment, setView]);

    return null;
}
