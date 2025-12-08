/**
 * Theme Provider
 * Inyecta los colores del negocio activo en las CSS variables
 * Los colores se cargan desde localStorage o se obtienen del backend
 */

'use client';

import { useEffect } from 'react';
import { TokenStorage } from '@/shared/config';

interface BusinessColors {
  primary: string;
  secondary: string;
  tertiary: string;
  quaternary: string;
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
function applyBusinessColors() {
  const colors = TokenStorage.getBusinessColors();

  if (colors) {
    document.documentElement.style.setProperty('--color-primary', colors.primary || '');
    document.documentElement.style.setProperty('--color-secondary', colors.secondary || '');
    document.documentElement.style.setProperty('--color-tertiary', colors.tertiary || '');
    document.documentElement.style.setProperty('--color-quaternary', colors.quaternary || '');
  }
}

/**
 * Hook para cambiar los colores del tema programáticamente
 */
export function useTheme() {
  const setColors = (colors: BusinessColors) => {
    // Guardar en localStorage
    TokenStorage.setBusinessColors(colors);

    // Aplicar inmediatamente
    document.documentElement.style.setProperty('--color-primary', colors.primary);
    document.documentElement.style.setProperty('--color-secondary', colors.secondary);
    document.documentElement.style.setProperty('--color-tertiary', colors.tertiary);
    document.documentElement.style.setProperty('--color-quaternary', colors.quaternary);

    // Disparar evento para otros componentes
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
    };
  };

  return {
    setColors,
    getColors,
  };
}

