'use client';

import { NotificationCounter } from '../../domain/types';
import { WhatsAppIcon } from '@/shared/ui';

interface Props {
    counter: NotificationCounter;
    onClick: () => void;
}

const STATE_STYLES: Record<string, string> = {
    none: 'bg-red-50 text-red-700 ring-red-200 hover:bg-red-100',
    partial: 'bg-orange-50 text-orange-700 ring-orange-200 hover:bg-orange-100',
    done: 'bg-green-50 text-green-700 ring-green-200 hover:bg-green-100',
};

const STATE_TITLES: Record<string, string> = {
    none: 'Sin notificaciones enviadas',
    partial: 'Faltan notificaciones por enviar',
    done: 'Todas las notificaciones enviadas',
};

export function NotificationBadge({ counter, onClick }: Props) {
    const style = STATE_STYLES[counter.state] || STATE_STYLES.none;
    const title = `${STATE_TITLES[counter.state] || ''} (${counter.sent}/${counter.expected})${
        counter.failed > 0 ? ` - ${counter.failed} fallida(s)` : ''
    }`;

    return (
        <button
            type="button"
            onClick={(e) => {
                e.stopPropagation();
                onClick();
            }}
            title={title}
            aria-label={title}
            className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-bold ring-1 transition-colors ${style}`}
        >
            <WhatsAppIcon className="h-4 w-4" />
            <span>
                {counter.sent}/{counter.expected}
            </span>
        </button>
    );
}
