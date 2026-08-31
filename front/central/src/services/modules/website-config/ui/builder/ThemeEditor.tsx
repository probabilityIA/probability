'use client';

import { useState } from 'react';
import { ImageField } from './editors';

const HEX_RE = /^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$/;

type ThemeRecord = Record<string, any>;

interface ThemePreset {
    name: string;
    primary: string;
    secondary: string;
    tertiary: string;
    quaternary: string;
    background: string;
}

const PRESETS: ThemePreset[] = [
    { name: 'Clasico', primary: '#1f2937', secondary: '#3b82f6', tertiary: '#10b981', quaternary: '#fbbf24', background: '#ffffff' },
    { name: 'Oceano', primary: '#0c4a6e', secondary: '#0284c7', tertiary: '#14b8a6', quaternary: '#f59e0b', background: '#f0f9ff' },
    { name: 'Bosque', primary: '#14532d', secondary: '#16a34a', tertiary: '#84cc16', quaternary: '#eab308', background: '#f0fdf4' },
    { name: 'Vino', primary: '#4c0519', secondary: '#be123c', tertiary: '#fb7185', quaternary: '#f59e0b', background: '#fff1f2' },
    { name: 'Purpura', primary: '#3b0764', secondary: '#7c3aed', tertiary: '#c026d3', quaternary: '#fbbf24', background: '#faf5ff' },
    { name: 'Minimal', primary: '#111827', secondary: '#374151', tertiary: '#6b7280', quaternary: '#d1d5db', background: '#fafafa' },
];

const COLOR_ROLES: Array<{ key: string; label: string }> = [
    { key: 'primary', label: 'Principal (hero, fondos)' },
    { key: 'secondary', label: 'Botones y enlaces' },
    { key: 'tertiary', label: 'Acento' },
    { key: 'quaternary', label: 'Detalle' },
    { key: 'background', label: 'Fondo de la página' },
];

interface ThemeEditorProps {
    theme: ThemeRecord | null;
    fallback: { primary: string; secondary: string; tertiary: string; quaternary: string; background: string };
    onChange: (theme: ThemeRecord) => void;
    businessId?: number;
    onImageDeleted?: () => void;
}

export function ThemeEditor({ theme, fallback, onChange, businessId, onImageDeleted }: ThemeEditorProps) {
    const current = {
        primary: theme?.primary || fallback.primary,
        secondary: theme?.secondary || fallback.secondary,
        tertiary: theme?.tertiary || fallback.tertiary,
        quaternary: theme?.quaternary || fallback.quaternary,
        background: theme?.background || fallback.background,
    };

    const patch = (p: ThemeRecord) => onChange({ ...(theme || {}), ...current, ...p });

    const [hexDrafts, setHexDrafts] = useState<Record<string, string>>({});

    const handleHexChange = (key: string, value: string) => {
        setHexDrafts(prev => ({ ...prev, [key]: value }));
        if (HEX_RE.test(value)) patch({ [key]: value });
    };

    const handleHexBlur = (key: string) => {
        setHexDrafts(prev => {
            const { [key]: _omit, ...rest } = prev;
            return rest;
        });
    };

    const isActivePreset = (p: ThemePreset) =>
        p.primary === current.primary && p.secondary === current.secondary &&
        p.tertiary === current.tertiary && p.quaternary === current.quaternary &&
        p.background === current.background;

    return (
        <div className="space-y-4">
            <div>
                <p className="text-xs font-medium text-gray-600 dark:text-gray-300 mb-2">Paletas</p>
                <div className="grid grid-cols-3 gap-2">
                    {PRESETS.map(p => (
                        <button
                            key={p.name}
                            type="button"
                            onClick={() => patch({ primary: p.primary, secondary: p.secondary, tertiary: p.tertiary, quaternary: p.quaternary, background: p.background })}
                            className={`rounded-lg border p-2 text-left hover:border-blue-400 ${isActivePreset(p) ? 'border-blue-500 ring-1 ring-blue-500' : 'border-gray-200 dark:border-gray-600'}`}
                        >
                            <span className="flex gap-1 mb-1">
                                {[p.primary, p.secondary, p.tertiary, p.quaternary, p.background].map((c, i) => (
                                    <span key={i} className="w-4 h-4 rounded-full border border-black/10" style={{ backgroundColor: c }} />
                                ))}
                            </span>
                            <span className="text-xs text-gray-600 dark:text-gray-300">{p.name}</span>
                        </button>
                    ))}
                </div>
            </div>

            <div className="space-y-2">
                <p className="text-xs font-medium text-gray-600 dark:text-gray-300">Personalizado</p>
                {COLOR_ROLES.map(role => (
                    <label key={role.key} className="flex items-center justify-between gap-3 text-sm text-gray-700 dark:text-gray-200">
                        <span className="text-xs">{role.label}</span>
                        <span className="flex items-center gap-2">
                            <input
                                type="color"
                                value={(current as ThemeRecord)[role.key]}
                                onChange={(e) => patch({ [role.key]: e.target.value })}
                                className="w-8 h-8 rounded cursor-pointer border border-gray-300 dark:border-gray-600 bg-transparent p-0"
                            />
                            <input
                                type="text"
                                value={hexDrafts[role.key] ?? (current as ThemeRecord)[role.key]}
                                onChange={(e) => handleHexChange(role.key, e.target.value)}
                                onBlur={() => handleHexBlur(role.key)}
                                placeholder="#000000"
                                spellCheck={false}
                                className="text-xs font-mono text-gray-600 dark:text-gray-300 w-20 rounded border border-gray-300 dark:border-gray-600 bg-transparent px-1.5 py-1"
                            />
                        </span>
                    </label>
                ))}
            </div>

            <div className="space-y-2 border-t border-gray-200 dark:border-gray-700 pt-3">
                <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                    <input
                        type="checkbox"
                        checked={theme?.apply_background !== false}
                        onChange={(e) => patch({ apply_background: e.target.checked })}
                        className="rounded border-gray-300"
                    />
                    Aplicar tema al fondo de la pagina
                </label>
                <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                    <input
                        type="checkbox"
                        checked={!!theme?.apply_navbar}
                        onChange={(e) => patch({ apply_navbar: e.target.checked })}
                        className="rounded border-gray-300"
                    />
                    Aplicar tema a la barra superior
                </label>
            </div>

            <div className="border-t border-gray-200 dark:border-gray-700 pt-3 space-y-3">
                <ImageField
                    label="Imagen de fondo de la página"
                    value={theme?.background_image || ''}
                    onChange={(v) => patch({ background_image: v })}
                    businessId={businessId}
                    onDeleted={onImageDeleted}
                />
                {theme?.background_image && (
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">
                            Velo de color sobre la imagen: {theme?.background_overlay ?? 80}%
                        </label>
                        <input
                            type="range"
                            min={0}
                            max={95}
                            step={5}
                            value={theme?.background_overlay ?? 80}
                            onChange={(e) => patch({ background_overlay: Number(e.target.value) })}
                            className="w-full"
                        />
                        <p className="text-xs text-gray-400">0% muestra la imagen pura; más alto la funde con el color de fondo.</p>
                    </div>
                )}
            </div>

            {theme && Object.keys(theme).length > 0 && (
                <button
                    type="button"
                    onClick={() => onChange({})}
                    className="text-xs text-gray-400 hover:text-gray-600 underline"
                >
                    Volver a los colores de mi marca
                </button>
            )}
        </div>
    );
}
