'use client';

import { useEffect, useRef, useState } from 'react';
import { Printer, ChevronDown, Crop } from 'lucide-react';
import { useGuideFormats } from '../hooks/useGuideFormats';
import { GuideFormat } from '../../domain/types';

interface Props {
    shipmentId: number;
    carrier?: string | null;
    className?: string;
}

export function ProbabilityGuideButton({ shipmentId, carrier, className = '' }: Props) {
    const { formats, loading } = useGuideFormats(carrier || undefined);
    const [open, setOpen] = useState(false);
    const ref = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handler = (e: MouseEvent) => {
            if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
        };
        document.addEventListener('mousedown', handler);
        return () => document.removeEventListener('mousedown', handler);
    }, []);

    const carrierFormat = formats.find((f) => f.strategy !== 'rebuild');
    const PREFERRED_DEFAULT_CODE = 'probability-10x10';
    const probabilityFormats = formats
        .filter((f) => f.strategy === 'rebuild')
        .sort((a, b) => {
            if (a.code === PREFERRED_DEFAULT_CODE) return -1;
            if (b.code === PREFERRED_DEFAULT_CODE) return 1;
            return a.sort_order - b.sort_order;
        });

    const openGuide = (format?: GuideFormat) => {
        const code = format?.code || '';
        const url = `/internal/shipment-guide/${shipmentId}${code ? `?format=${encodeURIComponent(code)}` : ''}`;
        window.open(url, '_blank');
        setOpen(false);
    };

    if (loading && formats.length === 0) {
        return (
            <div className={`${className} flex gap-1.5`}>
                <button disabled className="flex-1 flex items-center justify-center gap-1 py-1.5 px-2 rounded-md bg-zinc-300 text-white text-[11px] font-semibold cursor-wait">
                    <Printer size={11} /> Cargando...
                </button>
            </div>
        );
    }

    return (
        <div className={`${className} flex gap-1.5`}>
            {carrierFormat && (
                <button
                    onClick={() => openGuide(carrierFormat)}
                    className="flex-1 flex items-center justify-center gap-1 py-1.5 px-2 rounded-md bg-amber-600 hover:bg-amber-700 text-white text-[11px] font-semibold transition-colors"
                    title={`Imprimir ${carrierFormat.label}`}
                >
                    <Crop size={11} />
                    Recortada
                </button>
            )}

            {probabilityFormats.length > 0 && (
                <div className="relative inline-flex flex-1" ref={ref}>
                    <button
                        onClick={() => openGuide(probabilityFormats[0])}
                        className="flex-1 flex items-center justify-center gap-1 py-1.5 px-2 rounded-l-md bg-emerald-600 hover:bg-emerald-700 text-white text-[11px] font-semibold transition-colors"
                        title={`Imprimir ${probabilityFormats[0].label}`}
                    >
                        <Printer size={11} />
                        Probability
                    </button>
                    <button
                        onClick={() => setOpen((v) => !v)}
                        className="flex items-center justify-center px-1.5 rounded-r-md bg-emerald-700 hover:bg-emerald-800 text-white border-l border-emerald-800"
                        aria-label="Elegir tamano"
                        title="Elegir tamano"
                    >
                        <ChevronDown size={12} />
                    </button>
                    {open && (
                        <div className="absolute right-0 top-full mt-1 z-50 min-w-[200px] bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 rounded-md shadow-lg overflow-hidden">
                            <div className="px-3 py-1.5 text-[10px] font-semibold text-zinc-500 uppercase border-b border-zinc-100 dark:border-zinc-800">
                                Tamano etiqueta Probability
                            </div>
                            {probabilityFormats.map((f) => (
                                <button
                                    key={f.code}
                                    onClick={() => openGuide(f)}
                                    className="block w-full text-left px-3 py-2 hover:bg-zinc-100 dark:hover:bg-zinc-800 text-[12px] text-zinc-700 dark:text-zinc-200"
                                >
                                    <div className="font-medium">{f.label}</div>
                                    <div className="text-[10px] text-zinc-500">
                                        {f.width_cm} x {f.height_cm} cm
                                        {f.adhesive ? ' (adhesiva)' : ''}
                                    </div>
                                </button>
                            ))}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}
