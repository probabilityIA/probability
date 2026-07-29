'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { MarketingLeadsList } from '@/services/modules/marketingleads/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

export default function MarketingLeadsPage() {
    const router = useRouter();
    const { isSuperAdmin, isLoading } = usePermissions();

    useEffect(() => {
        if (isLoading || isSuperAdmin) return;
        router.replace('/home');
    }, [isLoading, isSuperAdmin, router]);

    if (!isLoading && !isSuperAdmin) return null;

    return (
        <div className="p-6 space-y-6">
            <MarketingLeadsList />
        </div>
    );
}
