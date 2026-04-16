import { RoleList } from '@/services/auth/roles/ui';

export default function RolesPageRoute() {
    return (
        <div className="p-6 space-y-6">
            <RoleList />
        </div>
    );
}
