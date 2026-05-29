'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { ResourceList } from '@/services/auth/resources/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

export default function ResourcesPage() {
    const router = useRouter();
    const { isSuperAdmin, hasPermission, isLoading } = usePermissions();

    useEffect(() => {
        if (isLoading || isSuperAdmin) return;
        if (!hasPermission('Recursos', 'Read')) router.replace('/users');
    }, [isLoading, isSuperAdmin, hasPermission, router]);

    if (!isLoading && !isSuperAdmin && !hasPermission('Recursos', 'Read')) return null;

    return (
        <div className="p-6 space-y-6">
            <ResourceList />
        </div>
    );
}
