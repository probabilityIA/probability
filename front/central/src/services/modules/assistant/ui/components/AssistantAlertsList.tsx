'use client';

import { useEffect } from 'react';
import { BellSlashIcon } from '@heroicons/react/24/outline';
import { AssistantAvatar } from './AssistantAvatar';
import { AssistantDestinationCard } from './AssistantDestinationCard';
import { isCurrentRoute, isWhatsAppAlert, relativeTime } from '../../app/use-cases';
import { WhatsAppMark } from './WhatsAppMark';
import type { AssistantAlerts } from '../hooks/use-assistant-alerts';
import type { AssistantAlert, AssistantDestination } from '../../domain/types';

interface AssistantAlertsListProps {
    alerts: AssistantAlerts;
    pathname: string;
    onGo: (destination: AssistantDestination) => void;
    onShowMe: (destination: AssistantDestination) => void;
}

const severityDot: Record<AssistantAlert['severity'], string> = {
    info: 'bg-sky-500',
    warning: 'bg-amber-500',
    critical: 'bg-red-500',
};

export function AssistantAlertsList({ alerts, pathname, onGo, onShowMe }: AssistantAlertsListProps) {
    const { alerts: items, loaded, loading, hasMore, load, loadMore, markSeen } = alerts;

    useEffect(() => {
        if (!loaded) void load(1);
    }, [loaded, load]);

    useEffect(() => {
        if (loaded) void markSeen();
    }, [loaded, markSeen]);

    if (loaded && items.length === 0) {
        return (
            <div className="flex flex-1 flex-col items-center justify-center gap-2 px-6 text-center">
                <BellSlashIcon className="h-8 w-8 text-gray-300 dark:text-gray-600" />
                <p className="text-sm font-semibold text-gray-700 dark:text-gray-200">{'Sin alertas por ahora'}</p>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                    {'Cuando pase algo que debas revisar (una cancelación, una guía rechazada, saldo bajo) te lo cuento aquí. Eliges qué avisos recibir en Notificaciones, canal Vía.'}
                </p>
            </div>
        );
    }

    return (
        <div className="flex-1 space-y-3 overflow-y-auto px-4 py-4">
            {items.map((alert) => {
                const whatsapp = isWhatsAppAlert(alert.event_type);
                const card = whatsapp
                    ? 'border-emerald-200 bg-emerald-50 dark:border-emerald-500/40 dark:bg-emerald-950/50'
                    : 'border-gray-200 bg-gray-50 dark:border-gray-700 dark:bg-gray-900/60';
                return (
                <article key={alert.id} className="flex items-start gap-2">
                    <AssistantAvatar size={24} className="mt-0.5" mood="idle" />
                    <div className="max-w-[88%] space-y-1.5">
                        <div className={`space-y-2 rounded-2xl rounded-tl-md border px-3 py-2 text-sm text-gray-800 dark:text-gray-100 ${card}`}>
                            <div className="flex items-center gap-1.5">
                                {whatsapp
                                    ? <WhatsAppMark className="h-4 w-4 shrink-0 text-emerald-600 dark:text-emerald-300" />
                                    : <span className={`h-2 w-2 shrink-0 rounded-full ${severityDot[alert.severity] || severityDot.info}`} />}
                                <p className={`font-semibold ${whatsapp ? 'text-emerald-800 dark:text-emerald-200' : ''}`}>{alert.title}</p>
                                {alert.unread && (
                                    <span className="rounded-full bg-[#5B1BE6] px-1.5 text-[10px] font-semibold text-white dark:bg-[#7148FF]">{'Nueva'}</span>
                                )}
                            </div>
                            <p className="whitespace-pre-wrap">{alert.body}</p>
                            {alert.destination && (
                                <AssistantDestinationCard
                                    destination={alert.destination}
                                    isCurrent={isCurrentRoute(pathname, alert.destination.route)}
                                    onGo={onGo}
                                    onShowMe={onShowMe}
                                />
                            )}
                        </div>
                        <p className="pl-1 text-[11px] text-gray-400 dark:text-gray-500">{relativeTime(alert.created_at)}</p>
                    </div>
                </article>
                );
            })}
            {hasMore && (
                <div className="flex justify-center pt-1">
                    <button
                        type="button"
                        disabled={loading}
                        onClick={loadMore}
                        className="rounded-full border border-gray-200 bg-white px-3 py-1 text-xs text-gray-600 hover:border-[#5B1BE6] hover:text-[#3E0FA8] disabled:opacity-50 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300"
                    >
                        {loading ? 'Cargando…' : 'Ver anteriores'}
                    </button>
                </div>
            )}
            {!loaded && loading && (
                <p className="text-center text-xs text-gray-400">{'Cargando alertas…'}</p>
            )}
        </div>
    );
}
