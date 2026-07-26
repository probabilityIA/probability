'use client';

interface PreviewColors {
    primary: string;
    secondary: string;
    tertiary: string;
    quaternary: string;
}

interface PreviewBars {
    sidebar: string;
    topbar: string;
}

interface BusinessThemePreviewProps {
    colors: PreviewColors;
    bars?: PreviewBars | null;
    businessName?: string;
    logoUrl?: string | null;
}

function readableText(hex: string): string {
    const value = hex.replace('#', '');
    const full = value.length === 3 ? value.split('').map(c => c + c).join('') : value;
    if (full.length !== 6) return '#ffffff';

    const r = parseInt(full.slice(0, 2), 16) / 255;
    const g = parseInt(full.slice(2, 4), 16) / 255;
    const b = parseInt(full.slice(4, 6), 16) / 255;
    const channel = (c: number) => (c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4));
    const luminance = 0.2126 * channel(r) + 0.7152 * channel(g) + 0.0722 * channel(b);

    return luminance > 0.55 ? '#1f2937' : '#ffffff';
}

const ROWS = [
    { order: '145849', client: 'Eduar Perdomo', total: '223.353', status: 'Completada' },
    { order: '145836', client: 'Daniel Trujillo', total: '52.798', status: 'En proceso' },
    { order: '145831', client: 'Ricardo Lizarazo', total: '31.376', status: 'Pendiente' },
];

export function BusinessThemePreview({ colors, bars, businessName, logoUrl }: BusinessThemePreviewProps) {
    const onPrimary = readableText(colors.primary);
    const topbarStyle = bars ? { backgroundColor: bars.topbar, color: readableText(bars.topbar) } : undefined;
    const sidebarStyle = bars ? { backgroundColor: bars.sidebar, color: readableText(bars.sidebar) } : undefined;
    const onSecondary = readableText(colors.secondary);
    const onTertiary = readableText(colors.tertiary);
    const onQuaternary = readableText(colors.quaternary);

    return (
        <div>
            <span className="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                Vista previa
            </span>

            <div className="mt-2 flex overflow-hidden rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-900 shadow-sm">
                <div
                    className={`flex w-7 flex-shrink-0 flex-col items-center gap-1.5 py-2 border-r border-gray-200 dark:border-gray-700 ${
                        sidebarStyle ? '' : 'bg-white dark:bg-gray-800'
                    }`}
                    style={sidebarStyle}
                >
                    <span
                        className="flex h-4 w-4 items-center justify-center overflow-hidden rounded"
                        style={{ backgroundColor: colors.primary, color: onPrimary }}
                    >
                        {logoUrl ? (
                            <img src={logoUrl} alt="" className="max-h-full max-w-full object-contain" />
                        ) : (
                            <span className="text-[7px] font-bold">{(businessName || 'P').charAt(0)}</span>
                        )}
                    </span>
                    {[0, 1, 2, 3, 4].map(item => (
                        <span
                            key={item}
                            className="h-1.5 w-3 rounded-full"
                            style={{
                                backgroundColor: item === 1 ? colors.primary : 'currentColor',
                                opacity: item === 1 ? 1 : 0.25,
                            }}
                        />
                    ))}
                </div>

                <div className="min-w-0 flex-1">
                <div
                    className={`flex items-center justify-between gap-2 px-2 py-1.5 border-b border-gray-200 dark:border-gray-700 ${
                        topbarStyle ? '' : 'bg-white dark:bg-gray-800'
                    }`}
                    style={topbarStyle}
                >
                    <div className="flex items-center gap-1.5">
                        <span
                            className="rounded-md px-1.5 py-0.5 text-[8px] font-semibold"
                            style={{ backgroundColor: colors.primary, color: onPrimary }}
                        >
                            Ordenes
                        </span>
                        <span className={`rounded-md px-1.5 py-0.5 text-[8px] font-medium ${topbarStyle ? 'opacity-80' : 'text-gray-500 dark:text-gray-400'}`}>
                            Envios
                        </span>
                        <span className={`rounded-md px-1.5 py-0.5 text-[8px] font-medium ${topbarStyle ? 'opacity-80' : 'text-gray-500 dark:text-gray-400'}`}>
                            Cotizaciones
                        </span>
                    </div>
                    <div className="flex items-center gap-1">
                        {[colors.primary, colors.secondary, colors.tertiary, colors.quaternary].map((color, index) => (
                            <span
                                key={index}
                                className="h-4 w-4 rounded-md border border-black/10"
                                style={{ backgroundColor: color }}
                            />
                        ))}
                    </div>
                </div>

                <div className="p-2">
                    <div className="overflow-hidden rounded-lg">
                        <div
                            className="grid grid-cols-[1fr_1.6fr_1fr_1.1fr] gap-1 px-2 py-1 text-[8px] font-semibold uppercase tracking-wide"
                            style={{ backgroundColor: colors.primary, color: onPrimary }}
                        >
                            <span>Orden</span>
                            <span>Cliente</span>
                            <span>Total</span>
                            <span>Estado</span>
                        </div>

                        {ROWS.map((row, index) => (
                            <div
                                key={row.order}
                                className={`grid grid-cols-[1fr_1.6fr_1fr_1.1fr] items-center gap-1 px-2 py-1 text-[8px] ${
                                    index % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50 dark:bg-gray-700/40'
                                }`}
                            >
                                <span className="font-semibold text-gray-700 dark:text-gray-200">#{row.order}</span>
                                <span className="truncate text-gray-500 dark:text-gray-400">{row.client}</span>
                                <span className="font-medium text-gray-700 dark:text-gray-200">{row.total}</span>
                                <span
                                    className="w-fit rounded-full px-1.5 py-0.5 font-semibold"
                                    style={{
                                        backgroundColor: index === 0 ? colors.tertiary : colors.secondary,
                                        color: index === 0 ? onTertiary : onSecondary,
                                    }}
                                >
                                    {row.status}
                                </span>
                            </div>
                        ))}

                        <div
                            className="flex items-center justify-between px-2 py-1 text-[8px] font-medium"
                            style={{ backgroundColor: colors.primary, color: onPrimary }}
                        >
                            <span>3 de 429</span>
                            <span className="flex items-center gap-1">
                                <span
                                    className="rounded px-1 py-0.5"
                                    style={{ backgroundColor: colors.secondary, color: onSecondary }}
                                >
                                    1
                                </span>
                                <span className="px-1 py-0.5 opacity-70">2</span>
                                <span className="px-1 py-0.5 opacity-70">3</span>
                            </span>
                        </div>
                    </div>

                    <div className="mt-2 flex flex-wrap items-center gap-1.5">
                        <span
                            className="rounded-md px-2 py-1 text-[8px] font-semibold"
                            style={{ backgroundColor: colors.primary, color: onPrimary }}
                        >
                            Nueva orden
                        </span>
                        <span
                            className="rounded-md px-2 py-1 text-[8px] font-semibold"
                            style={{ backgroundColor: colors.secondary, color: onSecondary }}
                        >
                            Guardar
                        </span>
                        <span
                            className="rounded-md px-2 py-1 text-[8px] font-semibold"
                            style={{ backgroundColor: colors.tertiary, color: onTertiary }}
                        >
                            Cotizar
                        </span>
                        <span
                            className="rounded-md px-2 py-1 text-[8px] font-semibold"
                            style={{ backgroundColor: colors.quaternary, color: onQuaternary }}
                        >
                            Exportar
                        </span>
                    </div>
                </div>
                </div>
            </div>
        </div>
    );
}
