'use client';

import { BusinessList } from '@/services/auth/business/ui';

export default function BusinessesPageRoute() {
    return (
        <div className="w-full px-6 py-8">
            <BusinessList />
        </div>
    );
}
