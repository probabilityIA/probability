'use client';

import { BusinessTypeList } from '@/services/auth/business/ui';

export default function BusinessTypesPageRoute() {
    return (
        <div className="w-full px-6 py-8">
            <BusinessTypeList />
        </div>
    );
}
