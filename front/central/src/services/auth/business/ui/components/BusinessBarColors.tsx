'use client';

export interface PaletteColors {
    primary: string;
    secondary: string;
    tertiary: string;
    quaternary: string;
}

export type PaletteToken = keyof PaletteColors;

export const PALETTE_LABELS: Record<PaletteToken, string> = {
    primary: 'Primario',
    secondary: 'Secundario',
    tertiary: 'Terciario',
    quaternary: 'Cuaternario',
};

export function matchPaletteToken(palette: PaletteColors, color?: string): PaletteToken {
    if (!color) return 'primary';
    const found = (Object.keys(palette) as PaletteToken[]).find(
        key => (palette[key] || '').toLowerCase() === color.toLowerCase(),
    );
    return found || 'primary';
}

interface BusinessBarColorsProps {
    colors: PaletteColors;
    enabled: boolean;
    sidebarToken: PaletteToken;
    topbarToken: PaletteToken;
    onToggle: (enabled: boolean) => void;
    onTokenChange: (bar: 'sidebar' | 'topbar', token: PaletteToken) => void;
}

export function BusinessBarColors({
    colors,
    enabled,
    sidebarToken,
    topbarToken,
    onToggle,
    onTokenChange,
}: BusinessBarColorsProps) {
    const rows: { bar: 'sidebar' | 'topbar'; label: string; token: PaletteToken }[] = [
        { bar: 'sidebar', label: 'Barra izquierda', token: sidebarToken },
        { bar: 'topbar', label: 'Barra superior', token: topbarToken },
    ];

    return (
        <div className="rounded-xl border border-gray-200 dark:border-gray-600 p-2.5">
            <label className="flex cursor-pointer items-center justify-between gap-2">
                <span className="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                    Colorear barras de navegacion
                </span>
                <input
                    type="checkbox"
                    checked={enabled}
                    onChange={e => onToggle(e.target.checked)}
                    className="h-4 w-4 cursor-pointer"
                />
            </label>

            {enabled && (
                <div className="mt-2 space-y-2">
                    {rows.map(row => (
                        <div key={row.bar} className="flex items-center justify-between gap-2">
                            <span className="text-[10px] font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500">
                                {row.label}
                            </span>
                            <div className="flex items-center gap-1.5">
                                {(Object.keys(PALETTE_LABELS) as PaletteToken[]).map(token => (
                                    <button
                                        key={token}
                                        type="button"
                                        title={PALETTE_LABELS[token]}
                                        onClick={() => onTokenChange(row.bar, token)}
                                        className={`h-6 w-6 rounded-md border transition-transform hover:scale-110 ${
                                            row.token === token
                                                ? 'border-gray-900 dark:border-white ring-2 ring-offset-1 ring-gray-300 dark:ring-gray-500'
                                                : 'border-black/10'
                                        }`}
                                        style={{ backgroundColor: colors[token] }}
                                    />
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
