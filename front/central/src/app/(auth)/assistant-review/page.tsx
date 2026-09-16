'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { AssistantReviewDashboard } from '@/services/modules/assistant/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';

export default function AssistantReviewPage() {
    const router = useRouter();
    const { isSuperAdmin, isLoading } = usePermissions();

    useEffect(() => {
        if (isLoading || isSuperAdmin) return;
        router.replace('/home');
    }, [isLoading, isSuperAdmin, router]);

    if (!isLoading && !isSuperAdmin) return null;

    return (
        <div className="p-6">
            <AssistantReviewDashboard />
        </div>
    );
}
