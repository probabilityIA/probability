'use client';

import { useEffect, useRef, useState } from 'react';
import {
    ArrowPathIcon,
    ArrowRightIcon,
    HandThumbDownIcon,
    HandThumbUpIcon,
    XMarkIcon,
} from '@heroicons/react/24/outline';
import {
    HandThumbDownIcon as HandThumbDownSolid,
    HandThumbUpIcon as HandThumbUpSolid,
} from '@heroicons/react/24/solid';
import { AssistantAlertsList } from './AssistantAlertsList';
import { AssistantAvatar } from './AssistantAvatar';
import { AssistantDestinationCard } from './AssistantDestinationCard';
import { ASSISTANT_NAME, formatResetTime, isCurrentRoute, isHubDestination } from '../../app/use-cases';
import type { AssistantChat } from '../hooks/use-assistant-chat';
import type { AssistantAlerts } from '../hooks/use-assistant-alerts';
import type { AssistantDestination, AssistantTab } from '../../domain/types';

interface AssistantPanelProps {
    chat: AssistantChat;
    alerts: AssistantAlerts;
    tab: AssistantTab;
    onTabChange: (tab: AssistantTab) => void;
    suggestions: string[];
    onClose: () => void;
    onShowMe: (destination: AssistantDestination, messageId?: string) => void;
}

const tabButton = 'relative flex-1 py-2 text-sm font-semibold transition-colors focus-visible:outline-2 focus-visible:outline-offset-[-2px] focus-visible:outline-[#5B1BE6]';

const feedbackButton =
    'rounded-md p-1 transition-colors focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-[#5B1BE6]';

export function AssistantPanel({ chat, alerts, tab, onTabChange, suggestions, onClose, onShowMe }: AssistantPanelProps) {
    const [draft, setDraft] = useState('');
    const inputRef = useRef<HTMLInputElement>(null);
    const listRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (tab === 'chat') inputRef.current?.focus();
    }, [tab]);

    useEffect(() => {
        const list = listRef.current;
        if (list) list.scrollTo({ top: list.scrollHeight, behavior: 'smooth' });
    }, [chat.entries.length, chat.pending, tab]);

    useEffect(() => {
        const onKey = (event: KeyboardEvent) => {
            if (event.key === 'Escape') onClose();
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, [onClose]);

    const submit = (event: React.FormEvent) => {
        event.preventDefault();
        if (!draft.trim()) return;
        chat.send(draft);
        setDraft('');
    };

    const lastAssistantId = [...chat.entries].reverse().find((e) => e.kind === 'assistant')?.id;
    const showSuggestions = chat.entries.length === 1 && !chat.pending && suggestions.length > 0;
    const resetTime = formatResetTime(chat.blockedUntil);
    const placeholder = chat.blocked
        ? resetTime ? `Disponible a las ${resetTime}` : 'Disponible m\u00e1s tarde'
        : 'Pregunta d\u00f3nde est\u00e1 algo\u2026';

    return (
        <section
            role="dialog"
            aria-label={`Asistente ${ASSISTANT_NAME}`}
            className="fixed inset-x-3 bottom-3 top-16 z-40 flex flex-col overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl dark:border-gray-700 dark:bg-gray-800 sm:inset-auto sm:bottom-5 sm:right-5 sm:h-[min(640px,calc(100vh-6rem))] sm:w-[380px]"
        >
            <header className="flex items-center gap-3 border-b border-gray-200 px-4 py-3 dark:border-gray-700">
                <AssistantAvatar mood={chat.mood} size={36} />
                <div className="min-w-0 leading-tight">
                    <p className="text-base font-bold text-gray-900 dark:text-white">{ASSISTANT_NAME}</p>
                    <p className="truncate text-xs text-gray-500 dark:text-gray-400">{'Conoce los m\u00f3dulos de tu cuenta'}</p>
                </div>
                {tab === 'chat' && (
                    <button
                        type="button"
                        onClick={chat.reset}
                        title={'Nueva conversaci\u00f3n'}
                        aria-label={'Nueva conversaci\u00f3n'}
                        className="ml-auto rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-gray-700 dark:hover:text-gray-200"
                    >
                        <ArrowPathIcon className="h-5 w-5" />
                    </button>
                )}
                <button
                    type="button"
                    onClick={onClose}
                    title="Cerrar"
                    aria-label="Cerrar asistente"
                    className="rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-gray-700 dark:hover:text-gray-200"
                >
                    <XMarkIcon className="h-5 w-5" />
                </button>
            </header>

            <div role="tablist" className="flex border-b border-gray-200 dark:border-gray-700">
                <button
                    type="button"
                    role="tab"
                    aria-selected={tab === 'chat'}
                    onClick={() => onTabChange('chat')}
                    className={`${tabButton} ${tab === 'chat' ? 'text-[#5B1BE6] dark:text-violet-300' : 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'}`}
                >
                    {'Chat'}
                    {tab === 'chat' && <span className="absolute inset-x-6 bottom-0 h-0.5 rounded-full bg-[#5B1BE6] dark:bg-violet-300" />}
                </button>
                <button
                    type="button"
                    role="tab"
                    aria-selected={tab === 'alerts'}
                    onClick={() => onTabChange('alerts')}
                    className={`${tabButton} ${tab === 'alerts' ? 'text-[#5B1BE6] dark:text-violet-300' : 'text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'}`}
                >
                    <span className="inline-flex items-center gap-1.5">
                        {'Alertas'}
                        {alerts.unread > 0 && (
                            <span className="rounded-full bg-red-500 px-1.5 text-[10px] font-bold leading-4 text-white">{alerts.unread > 99 ? '99+' : alerts.unread}</span>
                        )}
                    </span>
                    {tab === 'alerts' && <span className="absolute inset-x-6 bottom-0 h-0.5 rounded-full bg-[#5B1BE6] dark:bg-violet-300" />}
                </button>
            </div>

            {tab === 'alerts' && (
                <AssistantAlertsList
                    alerts={alerts}
                    pathname={chat.pathname}
                    onGo={(destination) => chat.goTo(destination)}
                    onShowMe={(destination) => onShowMe(destination)}
                />
            )}

            {tab === 'chat' && (<>
            <div ref={listRef} className="flex-1 space-y-3 overflow-y-auto px-4 py-4" aria-live="polite">
                {chat.entries.map((entry) => {
                    if (entry.kind === 'user') {
                        return (
                            <div key={entry.id} className="flex justify-end">
                                <p className="max-w-[85%] whitespace-pre-wrap rounded-2xl rounded-br-md bg-[#5B1BE6] px-3 py-2 text-sm text-white dark:bg-[#7148FF]">
                                    {entry.text}
                                </p>
                            </div>
                        );
                    }
                    if (entry.kind === 'system') {
                        return (
                            <p key={entry.id} className="text-center font-mono text-[11px] text-gray-400 dark:text-gray-500">
                                {entry.text}
                            </p>
                        );
                    }
                    if (entry.kind === 'notice') {
                        const tone = entry.tone === 'warning'
                            ? 'bg-amber-50 text-amber-800 dark:bg-amber-500/10 dark:text-amber-300'
                            : 'bg-red-50 text-red-700 dark:bg-red-500/10 dark:text-red-300';
                        return (
                            <div key={entry.id} className={`rounded-xl px-3 py-2 text-sm ${tone}`}>
                                <p>{entry.text}</p>
                                {entry.retryable && (
                                    <button
                                        type="button"
                                        onClick={chat.retry}
                                        className="mt-2 rounded-lg bg-[#5B1BE6] px-3 py-1.5 text-xs font-semibold text-white hover:bg-[#4A12C9] dark:bg-[#7148FF]"
                                    >
                                        Reintentar
                                    </button>
                                )}
                            </div>
                        );
                    }
                    return (
                        <div key={entry.id} className="flex items-start gap-2">
                            <AssistantAvatar
                                size={24}
                                className="mt-0.5"
                                mood={entry.id === lastAssistantId && chat.mood === 'pointing' ? 'pointing' : 'idle'}
                            />
                            <div className="max-w-[85%] space-y-1.5">
                                <div className="space-y-2.5 rounded-2xl rounded-tl-md border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-800 dark:border-gray-700 dark:bg-gray-900/60 dark:text-gray-100">
                                    <p className="whitespace-pre-wrap">{entry.text}</p>
                                    {entry.destination && (
                                        <AssistantDestinationCard
                                            destination={entry.destination}
                                            isCurrent={!isHubDestination(entry.destination.key) && isCurrentRoute(chat.pathname, entry.destination.route)}
                                            onGo={(destination) => chat.goTo(destination, entry.messageId)}
                                            onShowMe={(destination) => onShowMe(destination, entry.messageId)}
                                        />
                                    )}
                                </div>
                                {entry.messageId && (
                                    <div className="flex items-center gap-0.5 pl-1">
                                        <span className="mr-1 text-[11px] text-gray-400 dark:text-gray-500">{'\u00bfTe sirvi\u00f3?'}</span>
                                        <button
                                            type="button"
                                            onClick={() => chat.rate(entry.id, 1)}
                                            aria-label={'Me sirvi\u00f3'}
                                            aria-pressed={entry.feedback === 1}
                                            title={'Me sirvi\u00f3'}
                                            className={`${feedbackButton} ${entry.feedback === 1 ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'}`}
                                        >
                                            {entry.feedback === 1 ? <HandThumbUpSolid className="h-4 w-4" /> : <HandThumbUpIcon className="h-4 w-4" />}
                                        </button>
                                        <button
                                            type="button"
                                            onClick={() => chat.rate(entry.id, -1)}
                                            aria-label={'No me sirvi\u00f3'}
                                            aria-pressed={entry.feedback === -1}
                                            title={'No me sirvi\u00f3'}
                                            className={`${feedbackButton} ${entry.feedback === -1 ? 'text-red-600 dark:text-red-400' : 'text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'}`}
                                        >
                                            {entry.feedback === -1 ? <HandThumbDownSolid className="h-4 w-4" /> : <HandThumbDownIcon className="h-4 w-4" />}
                                        </button>
                                    </div>
                                )}
                            </div>
                        </div>
                    );
                })}

                {chat.pending && (
                    <div className="flex items-start gap-2">
                        <AssistantAvatar size={24} mood="thinking" className="mt-0.5" />
                        <div className="flex gap-1 rounded-2xl rounded-tl-md border border-gray-200 bg-gray-50 px-3 py-3 dark:border-gray-700 dark:bg-gray-900/60">
                            <span className="h-1.5 w-1.5 animate-bounce rounded-full bg-gray-400 [animation-delay:-0.3s]" />
                            <span className="h-1.5 w-1.5 animate-bounce rounded-full bg-gray-400 [animation-delay:-0.15s]" />
                            <span className="h-1.5 w-1.5 animate-bounce rounded-full bg-gray-400" />
                        </div>
                    </div>
                )}
            </div>

            {showSuggestions && (
                <div className="flex flex-wrap gap-1.5 px-4 pb-2">
                    {suggestions.map((text) => (
                        <button
                            key={text}
                            type="button"
                            onClick={() => chat.send(text)}
                            className="rounded-full border border-gray-200 bg-white px-2.5 py-1 text-xs text-gray-600 transition-colors hover:border-[#5B1BE6] hover:text-[#3E0FA8] dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:border-[#8B63FF] dark:hover:text-violet-200"
                        >
                            {text}
                        </button>
                    ))}
                </div>
            )}

            <form onSubmit={submit} className="flex gap-2 border-t border-gray-200 px-3 py-2.5 dark:border-gray-700">
                <input
                    id="assistant-message"
                    ref={inputRef}
                    type="text"
                    value={draft}
                    maxLength={1000}
                    autoComplete="off"
                    disabled={chat.blocked}
                    onChange={(event) => setDraft(event.target.value)}
                    placeholder={placeholder}
                    aria-label={`Mensaje para ${ASSISTANT_NAME}`}
                    className="min-w-0 flex-1 rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-900 placeholder:text-gray-400 focus:border-[#5B1BE6] focus:outline-none focus:ring-2 focus:ring-[#5B1BE6]/20 disabled:opacity-60 dark:border-gray-600 dark:bg-gray-900 dark:text-white"
                />
                <button
                    type="submit"
                    disabled={chat.pending || chat.blocked || !draft.trim()}
                    aria-label="Enviar"
                    className="grid h-10 w-10 shrink-0 place-items-center rounded-xl bg-[#5B1BE6] text-white transition-colors hover:bg-[#4A12C9] disabled:opacity-50 dark:bg-[#7148FF]"
                >
                    <ArrowRightIcon className="h-4 w-4" />
                </button>
            </form>
            <p className="px-4 pb-2 text-center text-[11px] leading-snug text-gray-400 dark:text-gray-500">
                {`${ASSISTANT_NAME} puede equivocarse y nunca cambia de pantalla sin tu clic. Guardamos las conversaciones para mejorarlo.`}
            </p>
            </>)}
        </section>
    );
}
