'use client';

import { useState } from 'react';
import { TicketsManager } from '@/services/modules/tickets/ui';

export default function TicketsPage() {
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <TicketsManager selectedBusinessId={selectedBusinessId} onBusinessChange={setSelectedBusinessId} />
        </div>
    );
}
