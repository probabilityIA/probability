/** Convierte un color hex (#RRGGBB) a { r, g, b } */
export function hexToRgb(hex: string): { r: number; g: number; b: number } | null {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result
        ? { r: parseInt(result[1], 16), g: parseInt(result[2], 16), b: parseInt(result[3], 16) }
        : null;
}

/**
 * Genera los estilos inline para un badge de estado con color dinámico.
 * Usa el color de la BD para fondo translúcido, borde y texto.
 */
export function getStatusBadgeStyle(color: string | undefined): React.CSSProperties {
    if (!color) return { backgroundColor: '#F3F4F6', color: '#374151' };

    const rgb = hexToRgb(color);
    if (!rgb) return { backgroundColor: '#F3F4F6', color: '#374151' };

    return {
        backgroundColor: `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, 0.12)`,
        color: color,
        border: `1px solid rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, 0.3)`,
    };
}
