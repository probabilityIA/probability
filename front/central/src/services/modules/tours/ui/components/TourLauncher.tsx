'use client';

import { useEffect } from 'react';
import { QuestionMarkCircleIcon } from '@heroicons/react/24/outline';
import { useTour } from './TourProvider';

export function TourLauncher({ className }: { className?: string }) {
    const { availableTour, startTour, isRunning, registerInlineLauncher } = useTour();

    useEffect(() => registerInlineLauncher(), [registerInlineLauncher]);

    if (!availableTour || isRunning) return null;

    return (
        <button
            type="button"
            onClick={() => startTour(availableTour.key)}
            title={`Tutorial: ${availableTour.title}`}
            aria-label={`Abrir tutorial de ${availableTour.title}`}
            className={
                className ??
                'inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-gray-300 bg-white text-gray-500 shadow-sm transition-all hover:text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-300 dark:hover:text-white'
            }
            style={{ color: 'var(--color-primary, #7c3aed)' }}
        >
            <QuestionMarkCircleIcon className="h-5 w-5" />
        </button>
    );
}

export function TourLauncherFloating() {
    const { availableTour, startTour, isRunning, hasInlineLauncher } = useTour();

    if (!availableTour || isRunning || hasInlineLauncher) return null;

    return (
        <button
            type="button"
            onClick={() => startTour(availableTour.key)}
            title={`Tutorial: ${availableTour.title}`}
            aria-label={`Abrir tutorial de ${availableTour.title}`}
            className="fixed right-5 top-3 z-50 inline-flex h-9 items-center gap-1.5 rounded-full bg-white px-3 text-xs font-semibold shadow-lg ring-1 ring-gray-200 transition-all hover:shadow-xl dark:bg-gray-800 dark:ring-gray-700"
            style={{ color: 'var(--color-primary, #7c3aed)' }}
        >
            <QuestionMarkCircleIcon className="h-4 w-4" />
            Tutorial
        </button>
    );
}
