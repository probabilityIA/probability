'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { BusinessList } from '@/services/auth/business/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

export default function BusinessesPageRoute() {
    const router = useRouter();
    const { isSuperAdmin, hasPermission, isLoading } = usePermissions();

    useEffect(() => {
        if (isLoading) return;
        if (isSuperAdmin) return;
        if (!hasPermission('Empresas', 'Read')) {
            router.replace('/users');
        }
    }, [isLoading, isSuperAdmin, hasPermission, router]);

    if (!isLoading && !isSuperAdmin && !hasPermission('Empresas', 'Read')) {
        return null;
    }

    return (
        <div className="p-6 space-y-6">
            <BusinessList />
        </div>
    );
}
