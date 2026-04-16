'use client';

import { PermissionList } from '@/services/auth/permissions/ui';

export default function PermissionsPageRoute() {
    return (
        <div className="p-6 space-y-6">
            <PermissionList />
        </div>
    );
}
