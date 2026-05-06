'use client';

import { GeozoneType } from '../../domain/types';

const STYLES: Record<GeozoneType, { bg: string; text: string; border: string; label: string; emoji: string }> = {
    country:        { bg: 'bg-sky-100 dark:bg-sky-900/40',         text: 'text-sky-700 dark:text-sky-300',         border: 'border-sky-300 dark:border-sky-700',         label: 'Pais',          emoji: '🌎' },
    state:          { bg: 'bg-violet-100 dark:bg-violet-900/40',   text: 'text-violet-700 dark:text-violet-300',   border: 'border-violet-300 dark:border-violet-700',   label: 'Departamento',  emoji: '🗺️' },
    city:           { bg: 'bg-emerald-100 dark:bg-emerald-900/40', text: 'text-emerald-700 dark:text-emerald-300', border: 'border-emerald-300 dark:border-emerald-700', label: 'Municipio',     emoji: '🏙️' },
    admin_district: { bg: 'bg-indigo-100 dark:bg-indigo-900/40',   text: 'text-indigo-700 dark:text-indigo-300',   border: 'border-indigo-300 dark:border-indigo-700',   label: 'Localidad',     emoji: '🏛️' },
    locality:       { bg: 'bg-amber-100 dark:bg-amber-900/40',     text: 'text-amber-700 dark:text-amber-300',     border: 'border-amber-300 dark:border-amber-700',     label: 'Corregimiento', emoji: '🌾' },
    neighborhood:   { bg: 'bg-rose-100 dark:bg-rose-900/40',       text: 'text-rose-700 dark:text-rose-300',       border: 'border-rose-300 dark:border-rose-700',       label: 'UPZ',           emoji: '🏘️' },
    barrio:         { bg: 'bg-red-100 dark:bg-red-900/40',         text: 'text-red-700 dark:text-red-300',         border: 'border-red-300 dark:border-red-700',         label: 'Barrio',        emoji: '🏠' },
    custom:         { bg: 'bg-pink-100 dark:bg-pink-900/40',       text: 'text-pink-700 dark:text-pink-300',       border: 'border-pink-300 dark:border-pink-700',       label: 'Personalizada', emoji: '✨' },
};

export function getTypeStyle(t: GeozoneType) { return STYLES[t]; }

export default function TypeChip({ type, count }: { type: GeozoneType; count?: number }) {
    const s = STYLES[type];
    return (
        <span className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium border ${s.bg} ${s.text} ${s.border}`}>
            <span>{s.emoji}</span>
            <span>{s.label}</span>
            {count !== undefined && (
                <span className="ml-1 px-1.5 py-0.5 rounded-full bg-white/70 dark:bg-black/30 text-[10px] font-bold">{count}</span>
            )}
        </span>
    );
}
