export * from './http-logger';
export * from './apply-business-theme';
export * from './sound';

// Storage
export { TokenStorage } from './token-storage';
export { CookieStorage } from './cookie-storage';
export { SimpleCookieStorage } from './cookie-storage-simple';

// Server Auth - NO exportar aquí (usa next/headers, solo para Server Actions)
// Importar directamente: import { getAuthToken } from '@/shared/utils/server-auth';

// Tipos
export type {
    BusinessColors,
    BusinessData,
    UserData,
    ResourcePermission,
    UserPermissions
} from './cookie-storage';
