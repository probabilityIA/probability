'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { RoleList } from '@/services/auth/roles/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

export default function RolesPageRoute() {
    const router = useRouter();
    const { isSuperAdmin, hasPermission, isLoading } = usePermissions();

    useEffect(() => {
        if (isLoading || isSuperAdmin) return;
        if (!hasPermission('Roles', 'Read')) router.replace('/users');
    }, [isLoading, isSuperAdmin, hasPermission, router]);

    if (!isLoading && !isSuperAdmin && !hasPermission('Roles', 'Read')) return null;

    return (
        <div className="p-6 space-y-6">
            <RoleList />
        </div>
    );
}
