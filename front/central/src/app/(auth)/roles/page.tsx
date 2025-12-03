'use client';

import { RoleList } from '@/services/auth/roles/ui';

export default function RolesPageRoute() {
    return (
        <div className="w-full px-6 py-8">
            <RoleList />
        </div>
    );
}
