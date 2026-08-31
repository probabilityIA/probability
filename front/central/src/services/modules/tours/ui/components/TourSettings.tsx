'use client';

import { useState } from 'react';
import { ArrowPathIcon, CheckCircleIcon } from '@heroicons/react/24/outline';
import { useTour } from './TourProvider';

export function TourSettings() {
    const { allTours, progress, resetTour, resetAllTours } = useTour();
    const [ocupado, setOcupado] = useState<string | null>(null);

    if (allTours.length === 0) return null;

    const ejecutar = async (clave: string, accion: () => Promise<void>) => {
        setOcupado(clave);
        try {
            await accion();
        } finally {
            setOcupado(null);
        }
    };

    const completados = allTours.filter((tour) => {
        const estado = progress[tour.key]?.status;
        return estado === 'completed' || estado === 'skipped';
    }).length;

    return (
        <section className="bg-white dark:bg-white/5 border border-gray-100 dark:border-white/10 rounded-3xl p-6 shadow-xl shadow-gray-200/50 dark:shadow-none">
            <div className="flex items-start justify-between gap-4 mb-5">
                <div>
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white">Tutoriales</h3>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        {completados} de {allTours.length} vistos. Reinicia uno para volver a verlo al entrar al modulo.
                    </p>
                </div>
                <button
                    type="button"
                    onClick={() => ejecutar('__all__', resetAllTours)}
                    disabled={ocupado !== null}
                    className="shrink-0 inline-flex items-center gap-1.5 px-3 py-2 rounded-xl text-xs font-semibold text-gray-600 dark:text-gray-300 border border-gray-200 dark:border-white/10 hover:bg-gray-50 dark:hover:bg-white/5 disabled:opacity-50 transition-colors"
                >
                    <ArrowPathIcon className="w-4 h-4" />
                    Reiniciar todos
                </button>
            </div>

            <ul className="space-y-2">
                {allTours.map((tour) => {
                    const estado = progress[tour.key]?.status;
                    const visto = estado === 'completed' || estado === 'skipped';
                    return (
                        <li
                            key={tour.key}
                            className="flex items-center justify-between gap-3 p-3 rounded-2xl hover:bg-gray-50 dark:hover:bg-white/5 transition-colors"
                        >
                            <div className="min-w-0 flex items-center gap-2">
                                {visto && <CheckCircleIcon className="w-4 h-4 text-green-500 shrink-0" />}
                                <span className="text-sm font-medium text-gray-900 dark:text-white truncate">{tour.title}</span>
                            </div>
                            <button
                                type="button"
                                onClick={() => ejecutar(tour.key, () => resetTour(tour.key))}
                                disabled={ocupado !== null || !visto}
                                className="shrink-0 text-xs font-semibold text-indigo-600 dark:text-indigo-400 hover:underline disabled:opacity-40 disabled:no-underline disabled:cursor-not-allowed"
                            >
                                Reiniciar
                            </button>
                        </li>
                    );
                })}
            </ul>
        </section>
    );
}
