import type { TourDefinition, TourProgress } from '../domain/types';

export function findTourForRoute(
    definitions: TourDefinition[],
    pathname: string,
): TourDefinition | undefined {
    const candidatos = definitions.filter((tour) =>
        tour.routes.some((route) => pathname === route || pathname.startsWith(`${route}/`)),
    );
    if (candidatos.length === 0) return undefined;

    return candidatos.reduce((mejor, actual) => {
        const largoMejor = Math.max(...mejor.routes.map((r) => r.length));
        const largoActual = Math.max(...actual.routes.map((r) => r.length));
        return largoActual > largoMejor ? actual : mejor;
    });
}

export function shouldAutoStart(tour: TourDefinition, progress?: TourProgress): boolean {
    if (!tour.autoStart) return false;
    if (!progress) return true;
    if (progress.status === 'skipped') return false;
    if (progress.status === 'completed') return progress.version < tour.version;
    return true;
}

export function resolveVisibleSteps(tour: TourDefinition): TourDefinition {
    if (typeof document === 'undefined') return tour;

    const visibles = tour.steps.filter((step) => {
        if (!step.target) return true;
        if (step.route) return true;
        if (!step.optional) return true;
        return Boolean(document.querySelector(step.target));
    });

    if (visibles.length === tour.steps.length) return tour;
    return { ...tour, steps: visibles };
}

export function readLegacySeen(tour: TourDefinition): boolean {
    if (!tour.legacyStorageKey) return false;
    try {
        return window.localStorage.getItem(tour.legacyStorageKey) === 'true';
    } catch {
        return false;
    }
}

export function clearLegacySeen(tour: TourDefinition): void {
    if (!tour.legacyStorageKey) return;
    try {
        window.localStorage.removeItem(tour.legacyStorageKey);
    } catch {
        return;
    }
}
