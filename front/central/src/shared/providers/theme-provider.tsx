'use client';

import { useEffect } from 'react';
import { TokenStorage } from '@/shared/config';
import { updateAllColorScales } from '@/shared/utils/color-scales';
import { readableTextColor } from '@/shared/utils/apply-business-theme';

interface BusinessColors {
  primary: string;
  secondary: string;
  tertiary: string;
  quaternary: string;
  sidebar?: string;
  topbar?: string;
}

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    // Aplicar colores del negocio activo
    applyBusinessColors();

    // Escuchar cambios en localStorage (cuando cambie de negocio)
    const handleStorageChange = () => {
      applyBusinessColors();
    };

    window.addEventListener('storage', handleStorageChange);

    // También escuchar un evento custom para cambios locales
    window.addEventListener('businessChanged', handleStorageChange);

    return () => {
      window.removeEventListener('storage', handleStorageChange);
      window.removeEventListener('businessChanged', handleStorageChange);
    };
  }, []);

  return <>{children}</>;
}

/**
 * Aplica los colores del negocio activo a las CSS variables
 */
const DEFAULT_COLORS = {
  primary: '#0f172a',
  secondary: '#be185d',
  tertiary: '#06b6d4',
  quaternary: '#f59e0b',
};

function applyBusinessColors() {
  const colors = TokenStorage.getBusinessColors();
  const primaryColor = colors?.primary || DEFAULT_COLORS.primary;
  const secondaryColor = colors?.secondary || DEFAULT_COLORS.secondary;
  const tertiaryColor = colors?.tertiary || DEFAULT_COLORS.tertiary;
  const quaternaryColor = colors?.quaternary || DEFAULT_COLORS.quaternary;

  updateAllColorScales(primaryColor, secondaryColor, tertiaryColor, quaternaryColor);

  const root = document.documentElement;
  applyBarColor(root, 'sidebar', colors?.sidebar);
  applyBarColor(root, 'topbar', colors?.topbar);
}

function applyBarColor(root: HTMLElement, name: 'sidebar' | 'topbar', color?: string) {
  if (color) {
    root.style.setProperty(`--color-${name}`, color);
    root.style.setProperty(`--color-${name}-text`, readableTextColor(color));
    return;
  }

  root.style.removeProperty(`--color-${name}`);
  root.style.removeProperty(`--color-${name}-text`);
}

/**
 * Hook para cambiar los colores del tema programáticamente
 */
export function useTheme() {
  const setColors = (colors: BusinessColors) => {
    TokenStorage.setBusinessColors(colors);
    updateAllColorScales(colors.primary, colors.secondary, colors.tertiary, colors.quaternary);
    window.dispatchEvent(new Event('businessChanged'));
  };

  const getColors = (): BusinessColors | null => {
    const colors = TokenStorage.getBusinessColors();
    if (!colors) return null;
    return {
      primary: colors.primary || '',
      secondary: colors.secondary || '',
      tertiary: colors.tertiary || '',
      quaternary: colors.quaternary || '',
      sidebar: colors.sidebar || '',
      topbar: colors.topbar || '',
    };
  };

  return {
    setColors,
    getColors,
  };
}

