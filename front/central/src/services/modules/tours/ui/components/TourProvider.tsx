'use client';

import {
    createContext,
    useCallback,
    useContext,
    useEffect,
    useMemo,
    useRef,
    useState,
    type ReactNode,
} from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { TOUR_REGISTRY, TOUR_LIST } from '../../content';
import {
    clearLegacySeen,
    findTourForRoute,
    readLegacySeen,
    resolveVisibleSteps,
    shouldAutoStart,
} from '../../app/use-cases';
import {
    listTourProgressAction,
    resetAllToursAction,
    resetTourAction,
    saveTourProgressAction,
    skipAllToursAction,
} from '../../infra/actions';
import type { TourDefinition, TourProgress, TourStatus } from '../../domain/types';
import { TourRunner } from './TourRunner';

interface TourContextValue {
    availableTour?: TourDefinition;
    allTours: TourDefinition[];
    progress: Record<string, TourProgress>;
    isRunning: boolean;
    toursDisabled: boolean;
    hasInlineLauncher: boolean;
    registerInlineLauncher: () => () => void;
    startTour: (tourKey: string) => void;
    resetTour: (tourKey: string) => Promise<void>;
    resetAllTours: () => Promise<void>;
}

const TourContext = createContext<TourContextValue | null>(null);

function readBusinessId(): number | undefined {
    if (typeof window === 'undefined') return undefined;
    try {
        const stored = window.localStorage.getItem('selected_business_id');
        if (!stored) return undefined;
        const parsed = parseInt(stored, 10);
        return Number.isNaN(parsed) || parsed <= 0 ? undefined : parsed;
    } catch {
        return undefined;
    }
}

export function TourProvider({ children }: { children: ReactNode }) {
    const pathname = usePathname();
    const router = useRouter();
    const { hasPermission, isSuperAdmin, isLoading: permisosCargando } = usePermissions();

    const [progress, setProgress] = useState<Record<string, TourProgress>>({});
    const [cargado, setCargado] = useState(false);
    const [activeKey, setActiveKey] = useState<string | null>(null);
    const [activeTour, setActiveTour] = useState<TourDefinition | null>(null);
    const [stepIndex, setStepIndex] = useState(0);
    const [inlineLaunchers, setInlineLaunchers] = useState(0);
    const autoStartIntentado = useRef<Set<string>>(new Set());

    const businessId = useMemo(() => readBusinessId(), [pathname]);

    const cargarProgreso = useCallback(async () => {
        try {
            const items = await listTourProgressAction(businessId);
            const mapa: Record<string, TourProgress> = {};
            for (const item of items) mapa[item.tour_key] = item;
            setProgress(mapa);
        } catch {
            setProgress({});
        } finally {
            setCargado(true);
        }
    }, [businessId]);

    useEffect(() => {
        cargarProgreso();
    }, [cargarProgreso]);

    const registerInlineLauncher = useCallback(() => {
        setInlineLaunchers((n) => n + 1);
        return () => setInlineLaunchers((n) => Math.max(0, n - 1));
    }, []);

    const persistir = useCallback(
        async (tour: TourDefinition, status: TourStatus, indice: number) => {
            const optimista: TourProgress = {
                tour_key: tour.key,
                version: tour.version,
                status,
                step_index: indice,
            };
            setProgress((prev) => ({ ...prev, [tour.key]: optimista }));
            try {
                await saveTourProgressAction(
                    { tour_key: tour.key, version: tour.version, status, step_index: indice },
                    businessId,
                );
            } catch {
                return;
            }
        },
        [businessId],
    );

    const tourDisponible = useMemo(() => {
        const encontrado = findTourForRoute(TOUR_LIST, pathname ?? '');
        if (!encontrado) return undefined;
        if (encontrado.superAdminOnly && !isSuperAdmin) return undefined;
        if (encontrado.resource && !isSuperAdmin && !hasPermission(encontrado.resource, 'Read')) {
            return undefined;
        }
        return encontrado;
    }, [pathname, hasPermission, isSuperAdmin]);

    const toursVisibles = useMemo(
        () => TOUR_LIST.filter((tour) => !tour.superAdminOnly || isSuperAdmin),
        [isSuperAdmin],
    );

    const toursDisabled = useMemo(() => {
        if (!cargado) return false;
        const omitidos = toursVisibles.filter((tour) => progress[tour.key]?.status === 'skipped').length;
        return omitidos === toursVisibles.length;
    }, [cargado, progress, toursVisibles]);

    useEffect(() => {
        if (!cargado || permisosCargando || activeKey) return;
        if (!tourDisponible) return;
        if (autoStartIntentado.current.has(tourDisponible.key)) return;

        autoStartIntentado.current.add(tourDisponible.key);

        const guardado = progress[tourDisponible.key];
        if (!guardado && readLegacySeen(tourDisponible)) {
            clearLegacySeen(tourDisponible);
            persistir(tourDisponible, 'completed', 0);
            return;
        }

        if (!shouldAutoStart(tourDisponible, guardado)) return;

        const timer = setTimeout(() => {
            const resuelto = resolveVisibleSteps(tourDisponible, isSuperAdmin);
            if (resuelto.steps.length <= 1) return;
            setStepIndex(0);
            setActiveTour(resuelto);
            setActiveKey(tourDisponible.key);
        }, 600);

        return () => clearTimeout(timer);
    }, [cargado, permisosCargando, activeKey, tourDisponible, progress, persistir, isSuperAdmin]);

    const startTour = useCallback(
        (tourKey: string) => {
            const tour = TOUR_REGISTRY[tourKey];
            if (!tour) return;
            if (tour.superAdminOnly && !isSuperAdmin) return;
            autoStartIntentado.current.add(tourKey);
            setStepIndex(0);
            setActiveTour(resolveVisibleSteps(tour, isSuperAdmin));
            setActiveKey(tourKey);
        },
        [isSuperAdmin],
    );

    const resetTour = useCallback(
        async (tourKey: string) => {
            const tour = TOUR_REGISTRY[tourKey];
            if (tour) clearLegacySeen(tour);
            autoStartIntentado.current.delete(tourKey);
            setProgress((prev) => {
                const copia = { ...prev };
                delete copia[tourKey];
                return copia;
            });
            await resetTourAction(tourKey, businessId);
        },
        [businessId],
    );

    const resetAllTours = useCallback(async () => {
        for (const tour of toursVisibles) clearLegacySeen(tour);
        autoStartIntentado.current.clear();
        setProgress({});
        await resetAllToursAction(businessId);
    }, [businessId, toursVisibles]);

    const tourActivo = activeKey ? activeTour ?? TOUR_REGISTRY[activeKey] : undefined;

    const cerrar = useCallback(
        (status: TourStatus, indice: number) => {
            if (tourActivo) persistir(tourActivo, status, indice);
            setActiveKey(null);
            setActiveTour(null);
            setStepIndex(0);
        },
        [tourActivo, persistir],
    );

    const omitirTodos = useCallback(async () => {
        const omitidos: Record<string, TourProgress> = {};
        for (const tour of toursVisibles) {
            clearLegacySeen(tour);
            autoStartIntentado.current.add(tour.key);
            omitidos[tour.key] = {
                tour_key: tour.key,
                version: tour.version,
                status: 'skipped',
                step_index: 0,
            };
        }
        setProgress(omitidos);
        setActiveKey(null);
        setActiveTour(null);
        setStepIndex(0);
        try {
            await skipAllToursAction(
                toursVisibles.map((tour) => ({ tour_key: tour.key, version: tour.version })),
                businessId,
            );
        } catch {
            return;
        }
    }, [businessId, toursVisibles]);

    const value = useMemo<TourContextValue>(
        () => ({
            availableTour: tourDisponible,
            allTours: toursVisibles,
            progress,
            isRunning: Boolean(activeKey),
            toursDisabled,
            hasInlineLauncher: inlineLaunchers > 0,
            registerInlineLauncher,
            startTour,
            resetTour,
            resetAllTours,
        }),
        [
            tourDisponible,
            toursVisibles,
            progress,
            activeKey,
            toursDisabled,
            inlineLaunchers,
            registerInlineLauncher,
            startTour,
            resetTour,
            resetAllTours,
        ],
    );

    return (
        <TourContext.Provider value={value}>
            {children}
            {tourActivo && (
                <TourRunner
                    tour={tourActivo}
                    stepIndex={stepIndex}
                    onStepChange={setStepIndex}
                    onNavigate={(route) => router.push(route)}
                    onSkip={(indice) => cerrar('skipped', indice)}
                    onComplete={(indice) => cerrar('completed', indice)}
                    onSkipAll={omitirTodos}
                />
            )}
        </TourContext.Provider>
    );
}

export function useTour(): TourContextValue {
    const ctx = useContext(TourContext);
    if (!ctx) {
        return {
            availableTour: undefined,
            allTours: [],
            progress: {},
            isRunning: false,
            toursDisabled: false,
            hasInlineLauncher: false,
            registerInlineLauncher: () => () => undefined,
            startTour: () => undefined,
            resetTour: async () => undefined,
            resetAllTours: async () => undefined,
        };
    }
    return ctx;
}
