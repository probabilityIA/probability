'use client';

import { useEffect, useRef, useState } from 'react';
import type { AssistantDestination } from '../../domain/types';

interface AssistantDestinationCardProps {
    destination: AssistantDestination;
    isCurrent: boolean;
    onGo: (destination: AssistantDestination) => void;
}

export function AssistantDestinationCard({ destination, isCurrent, onGo }: AssistantDestinationCardProps) {
    const [copied, setCopied] = useState(false);
    const timer = useRef<number | null>(null);

    useEffect(() => () => {
        if (timer.current) window.clearTimeout(timer.current);
    }, []);

    const copy = () => {
        const url = `${window.location.origin}${destination.route}`;
        navigator.clipboard
            ?.writeText(url)
            .then(() => {
                setCopied(true);
                timer.current = window.setTimeout(() => setCopied(false), 1500);
            })
            .catch(() => undefined);
    };

    return (
        <div className="rounded-xl border border-gray-200 bg-white p-3 dark:border-gray-700 dark:bg-gray-800">
            <code className="inline-block rounded-md bg-violet-50 px-1.5 py-0.5 font-mono text-[11px] text-[#3E0FA8] dark:bg-violet-500/15 dark:text-violet-200">
                {destination.route}
            </code>
            <p className="mt-1.5 text-sm font-semibold text-gray-900 dark:text-white">{destination.label}</p>
            {destination.description && (
                <p className="text-xs text-gray-500 dark:text-gray-400">{destination.description}</p>
            )}
            <div className="mt-2.5 flex flex-wrap gap-2">
                <button
                    type="button"
                    disabled={isCurrent}
                    onClick={() => onGo(destination)}
                    className="rounded-lg bg-[#5B1BE6] px-3 py-1.5 text-xs font-semibold text-white transition-colors hover:bg-[#4A12C9] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#5B1BE6] disabled:cursor-default disabled:opacity-60 dark:bg-[#7148FF] dark:hover:bg-[#8B63FF]"
                >
                    {isCurrent ? 'Ya est\u00e1s aqu\u00ed' : `Ir a ${destination.label}`}
                </button>
                <button
                    type="button"
                    onClick={copy}
                    className="rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-semibold text-gray-600 transition-colors hover:border-gray-400 hover:text-gray-900 dark:border-gray-600 dark:text-gray-300 dark:hover:text-white"
                >
                    {copied ? 'Copiado' : 'Copiar enlace'}
                </button>
            </div>
        </div>
    );
}
