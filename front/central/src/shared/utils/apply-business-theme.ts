import { TokenStorage } from '@/shared/config';
import type { BusinessColors } from './cookie-storage';
import { onColor } from './color-scales';

interface Business {
  primary_color?: string;
  secondary_color?: string;
  tertiary_color?: string;
  quaternary_color?: string;
  sidebar_color?: string;
  topbar_color?: string;
  logo_url?: string;
  name: string;
}

export function readableTextColor(hex?: string): string {
  if (!hex) return '#4b5563';

  const value = hex.replace('#', '');
  const full = value.length === 3 ? value.split('').map(c => c + c).join('') : value;
  if (full.length !== 6) return '#4b5563';

  const r = parseInt(full.slice(0, 2), 16) / 255;
  const g = parseInt(full.slice(2, 4), 16) / 255;
  const b = parseInt(full.slice(4, 6), 16) / 255;
  const channel = (c: number) => (c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4));
  const luminance = 0.2126 * channel(r) + 0.7152 * channel(g) + 0.0722 * channel(b);

  return luminance > 0.55 ? '#1f2937' : '#f9fafb';
}

function applyBarVariables(root: HTMLElement, name: 'sidebar' | 'topbar', color?: string): void {
  if (color) {
    root.style.setProperty(`--color-${name}`, color);
    root.style.setProperty(`--color-${name}-text`, readableTextColor(color));
    root.classList.add(`${name}-tinted`);
    return;
  }

  root.style.removeProperty(`--color-${name}`);
  root.style.removeProperty(`--color-${name}-text`);
  root.classList.remove(`${name}-tinted`);
}

export function applyBusinessTheme(business: Business): void {
  const colors: BusinessColors = {
    primary: business.primary_color,
    secondary: business.secondary_color,
    tertiary: business.tertiary_color,
    quaternary: business.quaternary_color,
    sidebar: business.sidebar_color,
    topbar: business.topbar_color,
  };

  TokenStorage.setBusinessColors(colors);

  if (typeof window !== 'undefined') {
    const root = document.documentElement;
    root.style.setProperty('--color-primary', colors.primary || '#0f172a');
    root.style.setProperty('--color-secondary', colors.secondary || '#be185d');
    root.style.setProperty('--color-tertiary', colors.tertiary || '#06b6d4');
    root.style.setProperty('--color-quaternary', colors.quaternary || '#f59e0b');
    root.style.setProperty('--color-on-primary', onColor(colors.primary || '#0f172a'));
    root.style.setProperty('--color-on-secondary', onColor(colors.secondary || '#be185d'));
    root.style.setProperty('--color-on-tertiary', onColor(colors.tertiary || '#06b6d4'));
    root.style.setProperty('--color-on-quaternary', onColor(colors.quaternary || '#f59e0b'));
    applyBarVariables(root, 'sidebar', colors.sidebar);
    applyBarVariables(root, 'topbar', colors.topbar);

    if (business.logo_url) {
      localStorage.setItem('selected_business_logo', business.logo_url);
    } else {
      localStorage.removeItem('selected_business_logo');
    }
    localStorage.setItem('selected_business_name', business.name);

    window.dispatchEvent(new Event('businessChanged'));
  }
}

export function resetTheme(): void {
  const defaultColors: BusinessColors = {
    primary: '#0f172a',
    secondary: '#be185d',
    tertiary: '#06b6d4',
    quaternary: '#f59e0b',
  };

  TokenStorage.setBusinessColors(defaultColors);

  if (typeof window !== 'undefined') {
    const root = document.documentElement;
    root.style.setProperty('--color-primary', defaultColors.primary || '#0f172a');
    root.style.setProperty('--color-secondary', defaultColors.secondary || '#be185d');
    root.style.setProperty('--color-tertiary', defaultColors.tertiary || '#06b6d4');
    root.style.setProperty('--color-quaternary', defaultColors.quaternary || '#f59e0b');
    applyBarVariables(root, 'sidebar', undefined);
    applyBarVariables(root, 'topbar', undefined);

    localStorage.removeItem('selected_business_logo');
    localStorage.removeItem('selected_business_name');

    window.dispatchEvent(new Event('businessChanged'));
  }
}
