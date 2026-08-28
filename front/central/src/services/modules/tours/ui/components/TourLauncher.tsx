'use client';

import { useEffect } from 'react';
import { QuestionMarkCircleIcon } from '@heroicons/react/24/outline';
import { useTour } from './TourProvider';

const FALLBACK_TOUR = 'home';

export function TourLauncher({ className }: { className?: string }) {
    const { availableTour, allTours, startTour, isRunning, registerInlineLauncher } = useTour();

    useEffect(() => registerInlineLauncher(), [registerInlineLauncher]);

    if (isRunning) return null;

    const tour = availableTour ?? allTours.find((t) => t.key === FALLBACK_TOUR);
    if (!tour) return null;

    return (
        <button
            type="button"
            onClick={() => startTour(tour.key)}
            title={`Tutorial: ${tour.title}`}
            aria-label={`Abrir tutorial de ${tour.title}`}
            className={
                className ??
                'inline-flex h-9 shrink-0 items-center gap-1.5 rounded-lg border border-gray-300 bg-white px-2.5 text-xs font-semibold shadow-sm transition-all hover:shadow-md dark:border-gray-600 dark:bg-gray-700'
            }
            style={{ color: 'var(--color-primary, #7c3aed)' }}
        >
            <QuestionMarkCircleIcon className="h-5 w-5" />
            <span className="hidden lg:inline">Tutorial</span>
        </button>
    );
}

export function TourLauncherFloating() {
    const { availableTour, allTours, startTour, isRunning, hasInlineLauncher } = useTour();

    if (isRunning || hasInlineLauncher) return null;

    const tour = availableTour ?? allTours.find((t) => t.key === FALLBACK_TOUR);
    if (!tour) return null;

    return (
        <button
            type="button"
            onClick={() => startTour(tour.key)}
            title={`Tutorial: ${tour.title}`}
            aria-label={`Abrir tutorial de ${tour.title}`}
            className="fixed right-5 top-3 z-50 inline-flex h-9 items-center gap-1.5 rounded-full bg-white px-3 text-xs font-semibold shadow-lg ring-1 ring-gray-200 transition-all hover:shadow-xl dark:bg-gray-800 dark:ring-gray-700"
            style={{ color: 'var(--color-primary, #7c3aed)' }}
        >
            <QuestionMarkCircleIcon className="h-4 w-4" />
            Tutorial
        </button>
    );
}
