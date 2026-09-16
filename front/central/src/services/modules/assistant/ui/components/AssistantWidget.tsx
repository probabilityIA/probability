'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { AssistantAlertToast } from './AssistantAlertToast';
import { AssistantAvatar } from './AssistantAvatar';
import { AssistantPanel } from './AssistantPanel';
import { AssistantSpotlight } from './AssistantSpotlight';
import { useAssistantAlerts } from '../hooks/use-assistant-alerts';
import { useAssistantChat } from '../hooks/use-assistant-chat';
import { ASSISTANT_NAME, INTRO_DELAY_MS, pickSuggestions } from '../../app/use-cases';
import { getAssistantStateAction, markAssistantIntroSeenAction } from '../../infra/actions';
import type { AssistantDestination, AssistantHighlight, AssistantTab } from '../../domain/types';

export function AssistantWidget() {
    const { access, isLoading, roleCode, hasNav } = usePermissions();
    const chat = useAssistantChat();
    const { applyState, prepareHighlight } = chat;
    const [open, setOpen] = useState(false);
    const [tab, setTab] = useState<AssistantTab>('chat');
    const [showHint, setShowHint] = useState(false);
    const [highlight, setHighlight] = useState<AssistantHighlight | null>(null);
    const introPending = useRef(false);

    const enabled = !isLoading && !!access && roleCode !== 'cliente_final';
    const alerts = useAssistantAlerts(enabled);

    useEffect(() => {
        if (!enabled) return;
        let cancelled = false;
        let timer: number | null = null;
        getAssistantStateAction().then((result) => {
            if (cancelled || !result.success) return;
            applyState(result.data);
            if (!result.data.intro_seen) {
                introPending.current = true;
                timer = window.setTimeout(() => setShowHint(true), INTRO_DELAY_MS);
            }
        });
        return () => {
            cancelled = true;
            if (timer) window.clearTimeout(timer);
        };
    }, [enabled, applyState]);

    const markIntroSeen = useCallback(() => {
        setShowHint(false);
        if (!introPending.current) return;
        introPending.current = false;
        void markAssistantIntroSeenAction();
    }, []);

    const openPanel = (nextTab: AssistantTab = 'chat') => {
        markIntroSeen();
        setHighlight(null);
        setTab(nextTab);
        setOpen(true);
    };

    const openAlerts = () => {
        alerts.dismissToast();
        openPanel('alerts');
    };

    const closePanel = useCallback(() => setOpen(false), []);

    const showMe = useCallback(
        (destination: AssistantDestination, messageId?: string) => {
            if (!destination.highlight) return;
            prepareHighlight(destination, messageId);
            setOpen(false);
            setHighlight(destination.highlight);
        },
        [prepareHighlight],
    );

    const closeHighlight = useCallback(() => setHighlight(null), []);

    if (!enabled) return null;

    return (
        <>
            {highlight && <AssistantSpotlight key={highlight.target} highlight={highlight} onClose={closeHighlight} />}

            {open ? (
                <AssistantPanel
                    chat={chat}
                    alerts={alerts}
                    tab={tab}
                    onTabChange={setTab}
                    suggestions={pickSuggestions(hasNav)}
                    onClose={closePanel}
                    onShowMe={showMe}
                />
            ) : (
                <div className="fixed bottom-5 right-5 z-40 flex items-end gap-3">
                    {alerts.toast && <AssistantAlertToast alert={alerts.toast} onOpen={openAlerts} onDismiss={alerts.dismissToast} />}
                    {showHint && !alerts.toast && (
                        <div className="relative mb-2 flex max-w-[220px] items-start gap-2 rounded-2xl rounded-br-md border border-gray-200 bg-white py-2 pl-3 pr-2 text-sm text-gray-700 shadow-lg dark:border-gray-700 dark:bg-gray-800 dark:text-gray-200">
                            <button type="button" onClick={() => openPanel()} className="text-left">
                                {'\u00bfBuscas algo? Preg\u00fantame.'}
                            </button>
                            <button
                                type="button"
                                onClick={markIntroSeen}
                                aria-label="Ocultar aviso"
                                className="rounded p-0.5 text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
                            >
                                <XMarkIcon className="h-4 w-4" />
                            </button>
                        </div>
                    )}
                    <button
                        type="button"
                        onClick={() => (alerts.unread > 0 ? openAlerts() : openPanel())}
                        title={alerts.unread > 0 ? `${alerts.unread} alerta(s) sin ver` : `Asistente ${ASSISTANT_NAME}`}
                        aria-label={`Abrir asistente ${ASSISTANT_NAME}`}
                        className="relative rounded-full shadow-lg transition-transform hover:scale-105 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#5B1BE6]"
                    >
                        <AssistantAvatar size={56} mood={highlight ? 'pointing' : chat.pending ? 'thinking' : alerts.toast ? 'happy' : 'idle'} />
                        {alerts.unread > 0 ? (
                            <span className="absolute -right-1 -top-1 flex h-5 min-w-5 items-center justify-center rounded-full border-2 border-white bg-red-500 px-1 text-[10px] font-bold leading-none text-white dark:border-gray-900">
                                {alerts.unread > 99 ? '99+' : alerts.unread}
                            </span>
                        ) : (
                            <span className="absolute -right-0.5 -top-0.5 h-3.5 w-3.5 rounded-full border-2 border-white bg-emerald-500 dark:border-gray-900" />
                        )}
                    </button>
                </div>
            )}
        </>
    );
}
