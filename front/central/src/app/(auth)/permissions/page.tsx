'use client';

import { PermissionList } from '@/services/auth/permissions/ui';

export default function PermissionsPageRoute() {
    return (
        <div className="w-full px-6 py-8">
            <PermissionList />
        </div>
    );
}
