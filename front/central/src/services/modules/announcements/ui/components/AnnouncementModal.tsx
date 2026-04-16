'use client';

import { useState, useEffect, useCallback } from 'react';
import { XMarkIcon, ChevronLeftIcon, ChevronRightIcon, CheckIcon } from '@heroicons/react/24/outline';
import { AnnouncementInfo } from '../../domain/types';
import { getActiveAnnouncementsAction, registerViewAction } from '../../infra/actions';
import { TokenStorage } from '@/shared/config';

export default function AnnouncementModal() {
    const [announcements, setAnnouncements] = useState<AnnouncementInfo[]>([]);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [visible, setVisible] = useState(false);
    const [imageIndex, setImageIndex] = useState(0);

    useEffect(() => {
        const fetchActive = async () => {
            try {
                const userData = TokenStorage.getUser();
                if (!userData) return;

                const businessesData = TokenStorage.getBusinessesData();
                const businessId = businessesData?.[0]?.id;

                const all = await getActiveAnnouncementsAction(businessId);
                const modals = (all || []).filter(
                    a => a.display_type === 'modal_image' || a.display_type === 'modal_text'
                );
                if (modals.length > 0) {
                    setAnnouncements(modals);
                    setVisible(true);
                    registerViewAction(modals[0].id, { action: 'viewed' }).catch(() => {});
                }
            } catch {
            }
        };
        fetchActive();
    }, []);

    const current = announcements[currentIndex];

    const handleClose = useCallback(() => {
        if (current) {
            registerViewAction(current.id, { action: 'closed' }).catch(() => {});
        }
        if (currentIndex < announcements.length - 1) {
            const nextIdx = currentIndex + 1;
            setCurrentIndex(nextIdx);
            setImageIndex(0);
            registerViewAction(announcements[nextIdx].id, { action: 'viewed' }).catch(() => {});
        } else {
            setVisible(false);
        }
    }, [current, currentIndex, announcements]);

    const handleAccept = useCallback(() => {
        if (current) {
            registerViewAction(current.id, { action: 'accepted' }).catch(() => {});
        }
        if (currentIndex < announcements.length - 1) {
            const nextIdx = currentIndex + 1;
            setCurrentIndex(nextIdx);
            setImageIndex(0);
            registerViewAction(announcements[nextIdx].id, { action: 'viewed' }).catch(() => {});
        } else {
            setVisible(false);
        }
    }, [current, currentIndex, announcements]);

    const handleLinkClick = (linkId: number) => {
        if (current) {
            registerViewAction(current.id, { action: 'clicked_link', link_id: linkId }).catch(() => {});
        }
    };

    if (!visible || !current) return null;

    const hasImages = current.images && current.images.length > 0;
    const requiresAcceptance = current.frequency_type === 'requires_acceptance';

    return (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-hidden flex flex-col">
                <div className="flex items-center justify-between px-5 py-3 border-b border-gray-200 dark:border-gray-700">
                    <div className="flex items-center gap-2">
                        {current.category && (
                            <span
                                className="w-2.5 h-2.5 rounded-full"
                                style={{ backgroundColor: current.category.color }}
                            />
                        )}
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
                            {current.title}
                        </h2>
                    </div>
                    <div className="flex items-center gap-2">
                        {announcements.length > 1 && (
                            <span className="text-xs text-gray-400">
                                {currentIndex + 1} / {announcements.length}
                            </span>
                        )}
                        {!requiresAcceptance && (
                            <button
                                onClick={handleClose}
                                className="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 rounded transition-colors"
                            >
                                <XMarkIcon className="w-5 h-5" />
                            </button>
                        )}
                    </div>
                </div>

                <div className="flex-1 overflow-y-auto p-5">
                    {hasImages && (
                        <div className="relative mb-4">
                            <img
                                src={current.images[imageIndex]?.image_url}
                                alt={current.title}
                                className="w-full h-auto rounded-lg object-cover max-h-64"
                            />
                            {current.images.length > 1 && (
                                <div className="absolute inset-x-0 bottom-0 flex items-center justify-between px-2 py-1">
                                    <button
                                        onClick={() => setImageIndex(Math.max(0, imageIndex - 1))}
                                        disabled={imageIndex === 0}
                                        className="p-1 bg-black/40 text-white rounded-full disabled:opacity-30"
                                    >
                                        <ChevronLeftIcon className="w-4 h-4" />
                                    </button>
                                    <div className="flex gap-1">
                                        {current.images.map((_, i) => (
                                            <span
                                                key={i}
                                                className={`w-1.5 h-1.5 rounded-full ${i === imageIndex ? 'bg-white' : 'bg-white/50'}`}
                                            />
                                        ))}
                                    </div>
                                    <button
                                        onClick={() => setImageIndex(Math.min(current.images.length - 1, imageIndex + 1))}
                                        disabled={imageIndex === current.images.length - 1}
                                        className="p-1 bg-black/40 text-white rounded-full disabled:opacity-30"
                                    >
                                        <ChevronRightIcon className="w-4 h-4" />
                                    </button>
                                </div>
                            )}
                        </div>
                    )}

                    {current.message && (
                        <p className="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-wrap leading-relaxed">
                            {current.message}
                        </p>
                    )}

                    {current.links && current.links.length > 0 && (
                        <div className="mt-4 flex flex-wrap gap-2">
                            {current.links.map(link => (
                                <a
                                    key={link.id}
                                    href={link.url}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    onClick={() => handleLinkClick(link.id)}
                                    className="inline-flex items-center px-3 py-1.5 text-sm font-medium text-purple-600 bg-purple-50 hover:bg-purple-100 dark:text-purple-400 dark:bg-purple-900/20 dark:hover:bg-purple-900/40 rounded-lg transition-colors"
                                >
                                    {link.label}
                                </a>
                            ))}
                        </div>
                    )}
                </div>

                <div className="px-5 py-3 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-2">
                    {requiresAcceptance ? (
                        <button
                            onClick={handleAccept}
                            className="inline-flex items-center gap-1.5 px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white text-sm font-medium rounded-lg transition-colors"
                        >
                            <CheckIcon className="w-4 h-4" />
                            Aceptar
                        </button>
                    ) : (
                        <button
                            onClick={handleClose}
                            className="px-4 py-2 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200 text-sm font-medium rounded-lg transition-colors"
                        >
                            {currentIndex < announcements.length - 1 ? 'Siguiente' : 'Cerrar'}
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
}
