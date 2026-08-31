'use client';

import { TruckIcon } from '@heroicons/react/24/outline';
import { SlidersHorizontal } from 'lucide-react';

interface ShippingConfigNodeProps {
    onOpen: () => void;
}

export function ShippingConfigNode({ onOpen }: ShippingConfigNodeProps) {
    return (
        <div className="flex flex-col items-center gap-2">
            <span
                className="rounded-full px-3 py-1 text-xs font-semibold text-white"
                style={{ backgroundColor: 'var(--color-primary)' }}
            >
                Envios
            </span>
            <div
                className="flex w-full min-w-[260px] items-center gap-3 rounded-2xl border bg-white p-4 shadow-sm dark:bg-gray-800"
                style={{ borderColor: 'var(--color-primary)' }}
            >
                <span
                    className="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-full"
                    style={{ backgroundColor: 'var(--color-quaternary, rgba(0,0,0,0.05))' }}
                >
                    <TruckIcon className="h-5 w-5" style={{ color: 'var(--color-primary)' }} />
                </span>

                <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-semibold text-gray-900 dark:text-white">Configuración de envíos</p>
                    <p className="truncate text-[11px] leading-tight text-gray-500 dark:text-gray-400">Transportadoras, cajas y origen</p>
                </div>

                <button
                    onClick={onOpen}
                    title="Configurar envíos"
                    aria-label="Configurar envios"
                    className="flex-shrink-0 rounded-lg border border-indigo-200 bg-indigo-50 p-1.5 text-indigo-600 shadow-sm transition-all hover:scale-105 hover:bg-indigo-100 hover:shadow dark:border-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300 dark:hover:bg-indigo-900/60"
                >
                    <SlidersHorizontal className="h-4 w-4" />
                </button>
            </div>
        </div>
    );
}
