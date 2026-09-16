'use client';

import { useEffect, useState } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { AssistantAvatar } from './AssistantAvatar';
import { ASSISTANT_NAME, isWhatsAppAlert } from '../../app/use-cases';
import { WhatsAppMark } from './WhatsAppMark';
import type { AssistantAlert } from '../../domain/types';

interface AssistantAlertToastProps {
    alert: AssistantAlert;
    onOpen: () => void;
    onDismiss: () => void;
}

const TYPING_MS = 18;

export function AssistantAlertToast({ alert, onOpen, onDismiss }: AssistantAlertToastProps) {
    const [shown, setShown] = useState(0);
    const text = alert.body;
    const whatsapp = isWhatsAppAlert(alert.event_type);
    const frame = whatsapp
        ? 'border-emerald-200 bg-emerald-50 dark:border-emerald-500/40 dark:bg-emerald-950/60'
        : 'border-violet-200 bg-white dark:border-violet-500/40 dark:bg-gray-800';
    const label = whatsapp ? 'text-emerald-700 dark:text-emerald-300' : 'text-[#5B1BE6] dark:text-violet-300';
    const caret = whatsapp ? 'bg-emerald-600 dark:bg-emerald-300' : 'bg-[#5B1BE6] dark:bg-violet-300';

    useEffect(() => {
        setShown(0);
        let index = 0;
        const timer = window.setInterval(() => {
            index += 2;
            setShown(index);
            if (index >= text.length) window.clearInterval(timer);
        }, TYPING_MS);
        return () => window.clearInterval(timer);
    }, [text, alert.id]);

    return (
        <div
            role="status"
            className={`relative mb-2 w-[min(300px,calc(100vw-7rem))] rounded-2xl rounded-br-md border py-2.5 pl-3 pr-8 text-sm text-gray-800 shadow-xl dark:text-gray-100 ${frame}`}
        >
            <button type="button" onClick={onOpen} className="block w-full text-left">
                <p className="text-[11px] font-semibold uppercase tracking-wide text-[#5B1BE6] dark:text-violet-300">
                    {`${ASSISTANT_NAME} · ${alert.title}`}
                </p>
                <p className="mt-0.5 whitespace-pre-wrap leading-snug">
                    {text.slice(0, shown)}
                    {shown < text.length && <span className={`ml-0.5 inline-block h-3.5 w-1 animate-pulse rounded-sm align-middle ${caret}`} />}
                </p>
                <p className="mt-1 text-[11px] text-gray-400 dark:text-gray-500">{whatsapp ? 'Toca para responder' : 'Toca para ver las alertas'}</p>
            </button>
            <button
                type="button"
                onClick={onDismiss}
                aria-label="Ocultar alerta"
                className="absolute right-1.5 top-1.5 rounded p-0.5 text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
            >
                <XMarkIcon className="h-4 w-4" />
            </button>
        </div>
    );
}
