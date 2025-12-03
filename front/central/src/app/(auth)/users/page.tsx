'use client';

import { UserList } from '@/services/auth/users/ui';

export default function UsersPageRoute() {
  return (
    <div className="w-full px-6 py-8">
      <UserList />
    </div>
  );
}
