'use client';

import { useParams } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import AnnouncementStats from '@/services/modules/announcements/ui/components/AnnouncementStats';

export default function AnnouncementStatsPage() {
    const params = useParams();
    const id = Number(params.id);

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="mb-4">
                <Link
                    href="/announcements"
                    className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
                >
                    <ArrowLeftIcon className="w-4 h-4" />
                    Volver a anuncios
                </Link>
            </div>
            <AnnouncementStats announcementId={id} />
        </div>
    );
}
