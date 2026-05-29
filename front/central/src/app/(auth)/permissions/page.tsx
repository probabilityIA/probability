'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { PermissionList } from '@/services/auth/permissions/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

export default function PermissionsPageRoute() {
    const router = useRouter();
    const { isSuperAdmin, hasPermission, isLoading } = usePermissions();

    useEffect(() => {
        if (isLoading || isSuperAdmin) return;
        if (!hasPermission('Permisos', 'Read')) router.replace('/users');
    }, [isLoading, isSuperAdmin, hasPermission, router]);

    if (!isLoading && !isSuperAdmin && !hasPermission('Permisos', 'Read')) return null;

    return (
        <div className="p-6 space-y-6">
            <PermissionList />
        </div>
    );
}
