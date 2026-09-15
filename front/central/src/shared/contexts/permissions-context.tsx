'use client';

import React, { createContext, useContext, useState, useEffect, useCallback, useMemo, ReactNode } from 'react';
import { TokenStorage } from '../utils';
import type { UserPermissions } from '../utils';
import { getMyAccessAction } from '@/services/auth/access/infra/actions';
import type { Access, AccessNavItem } from '@/services/auth/access/domain/types';
import { AccessDeniedModal } from '../ui/access-denied-modal';
import {
    ACCESS_DENIED_EVENT,
    infoFromBody,
    isAccessDeniedMessage,
    isPublicAuthRequest,
    readAccessDeniedCookie,
    type AccessDeniedInfo,
} from '../utils/access-denied';

const COOKIE_POLL_MS = 1000;

const ALWAYS_ALLOWED_ROUTES = ['/home', '/profile', '/subscription'];

interface PermissionsContextType {
    access: Access | null;
    permissions: UserPermissions | null;
    isLoading: boolean;
    loadError: string | null;
    isSuperAdmin: boolean;
    roleCode: string;
    navigation: AccessNavItem[];
    can: (permission: string) => boolean;
    hasNav: (key: string) => boolean;
    canAccessRoute: (pathname: string) => boolean;
    reloadAccess: () => Promise<void>;
}

const PermissionsContext = createContext<PermissionsContextType | undefined>(undefined);

function toLegacyPermissions(access: Access): UserPermissions {
    return {
        is_super: access.is_super,
        business_id: access.business?.id ?? 0,
        business_name: access.business?.name ?? '',
        role_id: access.role.id,
        role_name: access.role.name,
        resources: [],
        subscription_status: access.subscription.status,
    };
}

function routeBase(route: string): string {
    const first = route.split('/').filter(Boolean)[0];
    return first ? `/${first}` : '/';
}

function matchesRoute(pathname: string, route: string): boolean {
    return pathname === route || pathname.startsWith(`${route}/`);
}

export const PermissionsProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [access, setAccess] = useState<Access | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [loadError, setLoadError] = useState<string | null>(null);

    const reloadAccess = useCallback(async () => {
        setIsLoading(true);
        const res = await getMyAccessAction();
        if (res.success && res.data) {
            setAccess(res.data);
            setLoadError(null);
            TokenStorage.setPermissions(toLegacyPermissions(res.data));
        } else {
            setAccess(null);
            setLoadError(res.error || 'No se pudo cargar el acceso');
            TokenStorage.removeUserPermissions();
        }
        setIsLoading(false);
    }, []);

    useEffect(() => {
        reloadAccess();
    }, [reloadAccess]);

    const [denied, setDenied] = useState<AccessDeniedInfo | null>(null);

    const handleDenied = useCallback(
        (info: AccessDeniedInfo) => {
            setDenied((current) => current ?? info);
            if (info.code !== 'unauthenticated') {
                reloadAccess();
            }
        },
        [reloadAccess],
    );

    useEffect(() => {
        const onDenied = (event: Event) => handleDenied((event as CustomEvent<AccessDeniedInfo>).detail);
        window.addEventListener(ACCESS_DENIED_EVENT, onDenied);

        const poll = window.setInterval(() => {
            const info = readAccessDeniedCookie();
            if (info) handleDenied(info);
        }, COOKIE_POLL_MS);

        const originalFetch = window.fetch;
        window.fetch = async (input: RequestInfo | URL, init?: RequestInit) => {
            const response = await originalFetch(input, init);
            if (response.status === 401 || response.status === 402 || response.status === 403) {
                const url = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
                if (url.includes('/api/v1/') && !isPublicAuthRequest(url)) {
                    const body = await response.clone().json().catch(() => ({}));
                    const info = infoFromBody(response.status, body);
                    if (info) handleDenied(info);
                }
            }
            return response;
        };

        const originalAlert = window.alert;
        window.alert = (message?: unknown) => {
            const denied = isAccessDeniedMessage(String(message ?? ''));
            if (denied) {
                handleDenied(denied);
                return;
            }
            originalAlert(message as string);
        };

        return () => {
            window.removeEventListener(ACCESS_DENIED_EVENT, onDenied);
            window.clearInterval(poll);
            window.fetch = originalFetch;
            window.alert = originalAlert;
        };
    }, [handleDenied]);

    const permissionSet = useMemo(() => new Set(access?.permissions ?? []), [access]);
    const navKeys = useMemo(() => new Set((access?.navigation ?? []).map((n) => n.key)), [access]);
    const isSuperAdmin = access?.is_super === true;

    const can = useCallback(
        (permission: string) => isSuperAdmin || permissionSet.has(permission),
        [isSuperAdmin, permissionSet],
    );

    const hasNav = useCallback((key: string) => navKeys.has(key), [navKeys]);

    const canAccessRoute = useCallback(
        (pathname: string) => {
            if (!access) return false;
            if (isSuperAdmin) return true;
            if (ALWAYS_ALLOWED_ROUTES.some((route) => matchesRoute(pathname, route))) return true;
            return access.navigation.some((item) => matchesRoute(pathname, routeBase(item.route)));
        },
        [access, isSuperAdmin],
    );

    const value = useMemo<PermissionsContextType>(
        () => ({
            access,
            permissions: access ? toLegacyPermissions(access) : null,
            isLoading,
            loadError,
            isSuperAdmin,
            roleCode: access?.role.code ?? '',
            navigation: access?.navigation ?? [],
            can,
            hasNav,
            canAccessRoute,
            reloadAccess,
        }),
        [access, isLoading, loadError, isSuperAdmin, can, hasNav, canAccessRoute, reloadAccess],
    );

    return (
        <PermissionsContext.Provider value={value}>
            {children}
            <AccessDeniedModal info={denied} onClose={() => setDenied(null)} />
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

export const useCan = (permission: string): boolean => {
    const { can, isLoading } = usePermissions();
    if (isLoading) return false;
    return can(permission);
};
