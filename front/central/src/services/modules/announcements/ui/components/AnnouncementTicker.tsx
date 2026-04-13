'use client';

import { useState, useEffect, useRef } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { AnnouncementInfo } from '../../domain/types';
import { getActiveAnnouncementsAction, registerViewAction } from '../../infra/actions';
import { TokenStorage } from '@/shared/config';

export default function AnnouncementTicker() {
    const [announcements, setAnnouncements] = useState<AnnouncementInfo[]>([]);
    const [dismissed, setDismissed] = useState(false);
    const viewedRef = useRef<Set<number>>(new Set());

    useEffect(() => {
        const fetchTicker = async () => {
            try {
                const userData = TokenStorage.getUser();
                if (!userData) return;

                const businessesData = TokenStorage.getBusinessesData();
                const businessId = businessesData?.[0]?.id;

                const all = await getActiveAnnouncementsAction(businessId);
                const tickers = (all || []).filter(a => a.display_type === 'ticker');
                setAnnouncements(tickers);

                tickers.forEach(t => {
                    if (!viewedRef.current.has(t.id)) {
                        viewedRef.current.add(t.id);
                        registerViewAction(t.id, { action: 'viewed' }).catch(() => {});
                    }
                });
            } catch {
            }
        };
        fetchTicker();
    }, []);

    if (dismissed || announcements.length === 0) return null;

    const tickerText = announcements.map(a => a.title + (a.message ? ` - ${a.message}` : '')).join('     ');

    return (
        <div className="relative bg-purple-600 text-white overflow-hidden h-8 flex items-center">
            <div className="ticker-track whitespace-nowrap">
                <span className="inline-block px-8 text-sm font-medium">{tickerText}</span>
                <span className="inline-block px-8 text-sm font-medium">{tickerText}</span>
            </div>
            <button
                onClick={() => setDismissed(true)}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-0.5 hover:bg-purple-700 rounded transition-colors z-10"
            >
                <XMarkIcon className="w-4 h-4" />
            </button>
            <style jsx>{`
                .ticker-track {
                    display: flex;
                    animation: ticker-scroll 30s linear infinite;
                }
                @keyframes ticker-scroll {
                    0% { transform: translateX(0); }
                    100% { transform: translateX(-50%); }
                }
            `}</style>
        </div>
    );
}
