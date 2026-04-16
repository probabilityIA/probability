/**
 * API Configuration - Dinámica en runtime
 * Funciona igual en local y producción sin rebuild
 */

export function getApiUrl(): string {
  // En el navegador, detecta la URL base dinámicamente
  if (typeof window !== 'undefined') {
    const protocol = window.location.protocol; // http: o https:
    const host = window.location.host;         // localhost:4322 o www.probabilityia.com.co

    // En producción, usa ruta relativa (nginx hace proxy)
    // En local, construye URL completa al backend
    if (host.includes('localhost') || host.includes('127.0.0.1')) {
      // Local: apunta al backend en localhost:3050
      return 'http://localhost:3050/api/v1';
    }

    // Producción: usa ruta relativa (nginx proxy)
    return '/api/v1';
  }

  // Server-side (Astro SSR) - fallback
  return '/api/v1';
}

// Exportar como getter para que se evalúe en cada llamada (no en importación)
export const API_BASE_URL = getApiUrl();
