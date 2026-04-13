'use client';

import { useState, useEffect } from 'react';
import { AnnouncementStats as StatsType, AnnouncementInfo } from '../../domain/types';
import { getStatsAction, getAnnouncementByIdAction } from '../../infra/actions';
import { Spinner, Alert } from '@/shared/ui';
import { EyeIcon, UserGroupIcon, CursorArrowRaysIcon, CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/outline';

interface AnnouncementStatsProps {
    announcementId: number;
}

export default function AnnouncementStats({ announcementId }: AnnouncementStatsProps) {
    const [stats, setStats] = useState<StatsType | null>(null);
    const [announcement, setAnnouncement] = useState<AnnouncementInfo | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetch = async () => {
            setLoading(true);
            try {
                const [s, a] = await Promise.all([
                    getStatsAction(announcementId),
                    getAnnouncementByIdAction(announcementId),
                ]);
                setStats(s);
                setAnnouncement(a);
            } catch (err: any) {
                setError(err.message || 'Error al cargar estadisticas');
            } finally {
                setLoading(false);
            }
        };
        fetch();
    }, [announcementId]);

    if (loading) {
        return (
            <div className="flex justify-center items-center p-12">
                <Spinner size="lg" />
            </div>
        );
    }

    if (error) {
        return <Alert type="error">{error}</Alert>;
    }

    if (!stats || !announcement) return null;

    const cards = [
        { label: 'Vistas totales', value: stats.total_views, icon: EyeIcon, color: 'bg-blue-500' },
        { label: 'Usuarios unicos', value: stats.unique_users, icon: UserGroupIcon, color: 'bg-purple-500' },
        { label: 'Clicks en links', value: stats.total_clicks, icon: CursorArrowRaysIcon, color: 'bg-green-500' },
        { label: 'Aceptaciones', value: stats.total_acceptances, icon: CheckCircleIcon, color: 'bg-emerald-500' },
        { label: 'Cerrados', value: stats.total_closed, icon: XCircleIcon, color: 'bg-red-500' },
    ];

    return (
        <div className="space-y-6">
            <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{announcement.title}</h2>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                    {announcement.category?.name} - {announcement.display_type === 'modal_image' ? 'Modal imagen' : announcement.display_type === 'modal_text' ? 'Modal texto' : 'Ticker'}
                </p>
            </div>

            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
                {cards.map((card) => (
                    <div key={card.label} className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
                        <div className="flex items-center gap-3 mb-2">
                            <div className={`${card.color} p-2 rounded-lg`}>
                                <card.icon className="w-4 h-4 text-white" />
                            </div>
                        </div>
                        <p className="text-2xl font-bold text-gray-900 dark:text-white">{card.value}</p>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">{card.label}</p>
                    </div>
                ))}
            </div>

            {announcement.links && announcement.links.length > 0 && (
                <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
                    <h3 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Links del anuncio</h3>
                    <div className="space-y-2">
                        {announcement.links.map(link => (
                            <div key={link.id} className="flex items-center justify-between text-sm">
                                <span className="text-gray-700 dark:text-gray-300">{link.label}</span>
                                <a
                                    href={link.url}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-purple-600 hover:text-purple-700 truncate ml-2 max-w-[200px]"
                                >
                                    {link.url}
                                </a>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
