'use client';

import { useState, useTransition } from 'react';
import { restartContainerAction, stopContainerAction, startContainerAction } from '../../infra/actions';

interface ActionButtonsProps {
    containerId: string;
    state: string;
}

export function ActionButtons({ containerId, state }: ActionButtonsProps) {
    const [isPending, startTransition] = useTransition();
    const [message, setMessage] = useState<{ text: string; type: 'success' | 'error' } | null>(null);

    const handleAction = (action: 'restart' | 'stop' | 'start') => {
        setMessage(null);
        startTransition(async () => {
            const actions = { restart: restartContainerAction, stop: stopContainerAction, start: startContainerAction };
            const result = await actions[action](containerId);
            if (result.error) {
                setMessage({ text: result.error, type: 'error' });
            } else {
                setMessage({ text: result.message || `${action} successful`, type: 'success' });
            }
            setTimeout(() => setMessage(null), 4000);
        });
    };

    const isRunning = state === 'running';

    return (
        <div className="space-y-2">
            <div className="flex items-center gap-2">
                {isRunning ? (
                    <>
                        <button
                            onClick={() => handleAction('restart')}
                            disabled={isPending}
                            className="text-xs px-3 py-1.5 rounded-md transition-colors cursor-pointer disabled:opacity-50"
                            style={{ color: '#ffaa00', border: '1px solid #ffaa0030', background: '#ffaa0008' }}
                        >
                            {isPending ? 'Working...' : 'Restart'}
                        </button>
                        <button
                            onClick={() => handleAction('stop')}
                            disabled={isPending}
                            className="text-xs px-3 py-1.5 rounded-md transition-colors cursor-pointer disabled:opacity-50"
                            style={{ color: '#ff3366', border: '1px solid #ff336630', background: '#ff336608' }}
                        >
                            Stop
                        </button>
                    </>
                ) : (
                    <button
                        onClick={() => handleAction('start')}
                        disabled={isPending}
                        className="text-xs px-3 py-1.5 rounded-md transition-colors cursor-pointer disabled:opacity-50"
                        style={{ color: '#00ff88', border: '1px solid #00ff8830', background: '#00ff8808' }}
                    >
                        {isPending ? 'Starting...' : 'Start'}
                    </button>
                )}
            </div>

            {message && (
                <div
                    className="text-xs px-3 py-1.5 rounded-md"
                    style={{
                        background: message.type === 'success' ? '#00ff8810' : '#ff336610',
                        color: message.type === 'success' ? '#00ff88' : '#ff3366',
                        border: `1px solid ${message.type === 'success' ? '#00ff8830' : '#ff336630'}`,
                    }}
                >
                    {message.text}
                </div>
            )}
        </div>
    );
}
