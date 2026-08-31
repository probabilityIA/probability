'use client';

import { XMarkIcon } from '@heroicons/react/24/outline';

interface Props {
    title: string;
    body: string;
    stepNumber: number;
    totalSteps: number;
    style: React.CSSProperties;
    onNext: () => void;
    onPrev: () => void;
    onClose: () => void;
    onSkipAll: () => void;
}

export function TourCard({
    title,
    body,
    stepNumber,
    totalSteps,
    style,
    onNext,
    onPrev,
    onClose,
    onSkipAll,
}: Props) {
    const esPrimero = stepNumber === 1;
    const esUltimo = stepNumber === totalSteps;

    return (
        <div
            className="fixed z-[101] w-[370px] max-w-[calc(100vw-24px)] overflow-hidden rounded-xl bg-white shadow-2xl dark:bg-gray-800"
            style={style}
        >
            <div
                className="flex h-11 items-center justify-end px-3"
                style={{ background: 'var(--color-primary, #7c3aed)' }}
            >
                <button
                    type="button"
                    onClick={onClose}
                    aria-label="Cerrar tutorial"
                    className="rounded p-1 text-white/90 transition-colors hover:bg-white/20 hover:text-white"
                >
                    <XMarkIcon className="h-4 w-4" />
                </button>
            </div>

            <div className="px-5 pt-4 pb-3">
                <h3 className="mb-1.5 text-sm font-bold text-gray-900 dark:text-white">{title}</h3>
                <p className="whitespace-pre-line text-[13px] leading-relaxed text-gray-600 dark:text-gray-300">
                    {body}
                </p>
            </div>

            <div className="flex justify-center gap-1.5 pb-3">
                {Array.from({ length: totalSteps }).map((_, i) => (
                    <span
                        key={i}
                        className="h-1.5 rounded-full transition-all"
                        style={{
                            width: i === stepNumber - 1 ? 18 : 6,
                            background:
                                i === stepNumber - 1
                                    ? 'var(--color-primary, #7c3aed)'
                                    : 'rgb(209 213 219)',
                        }}
                    />
                ))}
            </div>

            <div className="flex items-center justify-between gap-2 px-5 pb-3">
                {esPrimero ? (
                    <span />
                ) : (
                    <button
                        type="button"
                        onClick={onPrev}
                        className="rounded-lg border px-3 py-1.5 text-xs font-semibold transition-colors hover:bg-gray-50 dark:hover:bg-white/5"
                        style={{
                            color: 'var(--color-primary, #7c3aed)',
                            borderColor: 'var(--color-primary, #7c3aed)',
                        }}
                    >
                        &larr; Anterior
                    </button>
                )}
                <button
                    type="button"
                    onClick={onNext}
                    className="rounded-lg px-4 py-1.5 text-xs font-semibold text-white transition-opacity hover:opacity-90"
                    style={{ background: 'var(--color-primary, #7c3aed)' }}
                >
                    {esUltimo ? 'Finalizar' : 'Siguiente →'}
                </button>
            </div>

            <div className="flex items-center justify-between gap-3 border-t border-gray-100 px-5 py-2.5 dark:border-gray-700">
                <button
                    type="button"
                    onClick={onClose}
                    className="text-[11px] font-medium text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-gray-200"
                >
                    Omitir este tutorial
                </button>
                <button
                    type="button"
                    onClick={onSkipAll}
                    className="text-[11px] font-medium text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-gray-200"
                >
                    No mostrar mas
                </button>
            </div>
        </div>
    );
}
