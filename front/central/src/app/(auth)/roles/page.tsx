'use client';

import { RoleList } from '@/services/auth/roles/ui';

export default function RolesPageRoute() {
    return (
        <div className="min-h-screen bg-gray-50">
            <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                <div className="px-4 py-6 sm:px-0">
                    <RoleList />
                </div>
            </div>
        </div>
    );
}
