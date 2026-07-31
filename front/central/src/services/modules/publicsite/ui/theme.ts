import { ThemeContent } from '../domain/types';

export interface BrandFallback {
    primary: string;
    secondary: string;
    tertiary: string;
    quaternary: string;
}

export function hexToRgba(hex: string, alpha: number): string {
    const clean = hex.replace('#', '');
    const full = clean.length === 3 ? clean.split('').map(c => c + c).join('') : clean;
    const num = parseInt(full, 16);
    if (Number.isNaN(num) || full.length !== 6) return `rgba(255,255,255,${alpha})`;
    const r = (num >> 16) & 255;
    const g = (num >> 8) & 255;
    const b = num & 255;
    return `rgba(${r},${g},${b},${alpha})`;
}

export function buildThemeStyle(theme: ThemeContent | null | undefined, fallback: BrandFallback): React.CSSProperties {
    const primary = theme?.primary || fallback.primary;
    const secondary = theme?.secondary || fallback.secondary;
    const tertiary = theme?.tertiary || fallback.tertiary;
    const quaternary = theme?.quaternary || fallback.quaternary;

    const applyBackground = theme?.apply_background !== false;
    const background = applyBackground ? (theme?.background || '#ffffff') : '#ffffff';

    const applyNavbar = !!theme?.apply_navbar;
    const navBg = applyNavbar ? primary : '#ffffff';
    const navText = applyNavbar ? '#ffffff' : '#4b5563';
    const navTextStrong = applyNavbar ? '#ffffff' : '#111827';

    const style: Record<string, string> = {
        '--brand-primary': primary,
        '--brand-secondary': secondary,
        '--brand-tertiary': tertiary,
        '--brand-quaternary': quaternary,
        '--brand-background': background,
        '--brand-nav-bg': navBg,
        '--brand-nav-text': navText,
        '--brand-nav-text-strong': navTextStrong,
        backgroundColor: background,
    };

    if (theme?.background_image) {
        const overlay = Math.min(Math.max(theme.background_overlay ?? 80, 0), 98) / 100;
        const veil = hexToRgba(background, overlay);
        style.backgroundImage = `linear-gradient(${veil}, ${veil}), url(${theme.background_image})`;
        style.backgroundSize = 'cover';
        style.backgroundPosition = 'center';
        style.backgroundAttachment = 'fixed';
    }

    return style as React.CSSProperties;
}
