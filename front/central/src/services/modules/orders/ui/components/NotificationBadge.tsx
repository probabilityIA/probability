'use client';

import { NotificationCounter } from '../../domain/types';

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
            <svg viewBox="0 0 24 24" className="h-4 w-4 fill-current" aria-hidden>
                <path d="M12.04 2C6.58 2 2.13 6.45 2.13 11.91c0 1.75.46 3.45 1.32 4.95L2 22l5.25-1.38a9.9 9.9 0 0 0 4.79 1.22h.01c5.46 0 9.9-4.45 9.9-9.91 0-2.65-1.03-5.14-2.9-7.01A9.82 9.82 0 0 0 12.04 2Zm0 18.15h-.01a8.2 8.2 0 0 1-4.19-1.15l-.3-.18-3.12.82.83-3.04-.2-.31a8.18 8.18 0 0 1-1.26-4.38c0-4.54 3.7-8.23 8.25-8.23a8.2 8.2 0 0 1 8.24 8.24c0 4.54-3.7 8.23-8.24 8.23Zm4.52-6.16c-.25-.12-1.47-.72-1.69-.81-.23-.08-.39-.12-.56.13-.16.24-.64.8-.78.97-.15.16-.29.18-.54.06-.25-.13-1.05-.39-1.99-1.23-.74-.66-1.23-1.47-1.38-1.72-.14-.25-.01-.38.11-.51.11-.11.25-.29.37-.44.13-.15.17-.25.25-.42.08-.17.04-.31-.02-.44-.06-.12-.56-1.35-.76-1.85-.2-.48-.41-.42-.56-.43h-.48c-.17 0-.44.06-.67.31-.23.25-.87.86-.87 2.09s.9 2.43 1.02 2.6c.12.16 1.76 2.68 4.26 3.76.6.26 1.06.41 1.42.53.6.19 1.14.16 1.57.1.48-.07 1.47-.6 1.68-1.18.21-.58.21-1.08.14-1.18-.06-.11-.22-.17-.47-.29Z" />
            </svg>
            <span>
                {counter.sent}/{counter.expected}
            </span>
        </button>
    );
}
