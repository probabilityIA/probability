'use client';

import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import { TokenStorage } from '../utils';
import type { UserPermissions, ResourcePermission } from '../utils';

// Mapeo de recursos del backend a rutas del frontend
const RESOURCE_ROUTE_MAP: Record<string, string> = {
    'Usuarios': '/users',
    'Roles': '/roles',
    'Permisos': '/permissions',
    'Recursos': '/resources',
    'Empresas': '/businesses',
    'Ordenes': '/orders',
    'Productos': '/products',
    'Envios': '/shipments',
    'Integraciones': '/integrations',
    'Configuración de Notificaciones': '/notification-config',
    'Estado de Ordenes': '/order-status',
    'Facturacion': '/invoicing',
};

// Mapeo inverso: rutas a recursos
const ROUTE_RESOURCE_MAP: Record<string, string> = Object.entries(RESOURCE_ROUTE_MAP).reduce(
    (acc, [resource, route]) => ({ ...acc, [route]: resource }),
    {}
);

// Acciones estándar
export type ActionType = 'Create' | 'Read' | 'Update' | 'Delete' | 'List';

interface PermissionsContextType {
    permissions: UserPermissions | null;
    isLoading: boolean;
    isSuperAdmin: boolean;
    // Verificar si tiene permiso sobre un recurso y acción
    hasPermission: (resource: string, action: ActionType) => boolean;
    // Verificar si tiene acceso a una ruta
    hasRouteAccess: (route: string) => boolean;
    // Obtener las acciones permitidas para un recurso
    getResourceActions: (resource: string) => string[];
    // Recargar permisos
    reloadPermissions: () => void;
    // Establecer permisos (después del login)
    setUserPermissions: (permissions: UserPermissions) => void;
}

const PermissionsContext = createContext<PermissionsContextType | undefined>(undefined);

export const PermissionsProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [permissions, setPermissions] = useState<UserPermissions | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    // Cargar permisos del localStorage al iniciar
    useEffect(() => {
        const stored = TokenStorage.getPermissions();
        if (stored) {
            setPermissions(stored);
        }
        setIsLoading(false);
    }, []);

    const isSuperAdmin = permissions?.is_super === true;

    // Verificar permiso sobre recurso y acción
    const hasPermission = useCallback((resource: string, action: ActionType): boolean => {
        // Super admin tiene todos los permisos
        if (isSuperAdmin) return true;

        if (!permissions?.resources) return false;

        const resourcePermission = permissions.resources.find(
            (r) => r.resource.toLowerCase() === resource.toLowerCase() && r.active
        );

        if (!resourcePermission) return false;

        return resourcePermission.actions.some(
            (a) => a.toLowerCase() === action.toLowerCase()
        );
    }, [permissions, isSuperAdmin]);

    // Verificar acceso a una ruta
    const hasRouteAccess = useCallback((route: string): boolean => {
        // Super admin tiene acceso a todo
        if (isSuperAdmin) return true;

        if (!permissions?.resources) return false;

        // Normalizar la ruta (quitar parámetros y trailing slash)
        const normalizedRoute = '/' + route.split('/').filter(Boolean)[0];
        const resourceName = ROUTE_RESOURCE_MAP[normalizedRoute];

        if (!resourceName) {
            // Si no está mapeado, permitir acceso (rutas públicas o no controladas)
            return true;
        }

        const resourcePermission = permissions.resources.find(
            (r) => r.resource === resourceName && r.active
        );

        // Si tiene el recurso y al menos una acción, tiene acceso
        return resourcePermission ? resourcePermission.actions.length > 0 : false;
    }, [permissions, isSuperAdmin]);

    // Obtener acciones permitidas para un recurso
    const getResourceActions = useCallback((resource: string): string[] => {
        // Super admin tiene todas las acciones
        if (isSuperAdmin) return ['Create', 'Read', 'Update', 'Delete', 'List'];

        if (!permissions?.resources) return [];

        const resourcePermission = permissions.resources.find(
            (r) => r.resource.toLowerCase() === resource.toLowerCase() && r.active
        );

        return resourcePermission?.actions || [];
    }, [permissions, isSuperAdmin]);

    // Recargar permisos del localStorage
    const reloadPermissions = useCallback(() => {
        const stored = TokenStorage.getPermissions();
        setPermissions(stored);
    }, []);

    // Establecer permisos (después del login)
    const setUserPermissions = useCallback((newPermissions: UserPermissions) => {
        TokenStorage.setPermissions(newPermissions);
        setPermissions(newPermissions);
    }, []);

    return (
        <PermissionsContext.Provider
            value={{
                permissions,
                isLoading,
                isSuperAdmin,
                hasPermission,
                hasRouteAccess,
                getResourceActions,
                reloadPermissions,
                setUserPermissions,
            }}
        >
            {children}
        </PermissionsContext.Provider>
    );
};

export const usePermissions = (): PermissionsContextType => {
    const context = useContext(PermissionsContext);
    if (context === undefined) {
        throw new Error('usePermissions must be used within a PermissionsProvider');
    }
    return context;
};

// Hook para verificar un permiso específico
export const useHasPermission = (resource: string, action: ActionType): boolean => {
    const { hasPermission, isLoading } = usePermissions();
    if (isLoading) return false;
    return hasPermission(resource, action);
};

// Componente para renderizado condicional basado en permisos
export const PermissionGate: React.FC<{
    resource: string;
    action: ActionType;
    children: ReactNode;
    fallback?: ReactNode;
}> = ({ resource, action, children, fallback = null }) => {
    const { hasPermission, isLoading } = usePermissions();

    if (isLoading) return null;

    return hasPermission(resource, action) ? <>{children}</> : <>{fallback}</>;
};
