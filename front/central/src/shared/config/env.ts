/**
 * Configuración de variables de entorno
 * IMPORTANTE: Sin valores por defecto
 * Si falta alguna variable requerida, la aplicación lanzará un error
 */

/**
 * Obtiene una variable de entorno requerida (servidor)
 * Lanza error si no existe
 */
function getRequiredEnv(key: string): string {
  const value = process.env[key];
  if (!value) {
    throw new Error(
      `❌ Variable de entorno requerida no encontrada: ${key}\n` +
      `Por favor, configúrala en tu archivo .env.local`
    );
  }
  return value;
}

/**
 * Variables de entorno del proyecto
 */
export const env = {
  // API Backend (privada - solo servidor)
  // REQUERIDA para hacer peticiones al backend
  get API_BASE_URL(): string {
    return getRequiredEnv('API_BASE_URL');
  },
} as const;

/**
 * Variables de entorno públicas (cliente)
 * IMPORTANTE: NEXT_PUBLIC_* se inyectan en build time, no runtime
 */
export const envPublic = {
  // API Backend para cliente (solo SSE)
  get API_BASE_URL(): string {
    const value = process.env.NEXT_PUBLIC_API_BASE_URL;
    if (!value) {
      throw new Error(
        `❌ NEXT_PUBLIC_API_BASE_URL no está definida.\n` +
        `Agrégala a .env.local y ejecuta: rm -rf .next && pnpm run dev`
      );
    }
    return value;
  },
} as const;
