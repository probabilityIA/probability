'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { TokenStorage } from '@/shared/config';
import { StorefrontNav } from '@/services/modules/storefront/ui/components/StorefrontNav';
import { Sidebar } from '@/shared/ui/sidebar';
import { StorefrontSubNavbar } from '@/shared/ui/storefront-subnavbar';
import { SidebarProvider, useSidebar } from '@/shared/contexts/sidebar-context';
import { PermissionsProvider } from '@/shared/contexts/permissions-context';
import { NavbarProvider } from '@/shared/contexts/navbar-context';
import { StorefrontBusinessProvider } from '@/shared/contexts/storefront-business-context';
import { ToastProvider } from '@/shared/providers/toast-provider';

const PUBLIC_PATHS = ['/storefront/login', '/storefront/registro'];

function AdminStorefrontContent({ user, children }: { user: any; children: React.ReactNode }) {
    const { primaryExpanded, requestCollapse, requestSecondaryCollapse } = useSidebar();
    const primaryWidth = primaryExpanded ? 250 : 80;

    const handleMainMouseEnter = () => {
        if (typeof window !== 'undefined' && window.innerWidth >= 768) {
            requestCollapse(false);
            requestSecondaryCollapse();
        }
    };

    return (
        <div className="flex min-h-screen bg-gray-50 dark:bg-gray-900">
            <Sidebar user={user} />
            <main
                className="flex-1 transition-all duration-300 w-full overflow-x-hidden main-content flex flex-col"
                onMouseEnter={handleMainMouseEnter}
            >
                <StorefrontBusinessProvider>
                    <StorefrontSubNavbar />
                    <div className="w-full min-w-0 flex-1">
                        {children}
                    </div>
                </StorefrontBusinessProvider>
                <style jsx>{`
                    .main-content {
                        margin-left: 0;
                    }
                    @media (min-width: 768px) {
                        .main-content {
                            margin-left: ${primaryWidth}px;
                        }
                    }
                `}</style>
            </main>
        </div>
    );
}

export default function StorefrontLayout({ children }: { children: React.ReactNode }) {
    const router = useRouter();
    const pathname = usePathname();
    const [loading, setLoading] = useState(true);
    const [authenticated, setAuthenticated] = useState(false);
    const [isAdmin, setIsAdmin] = useState(false);
    const [user, setUser] = useState<any>(null);

    const isPublicPage = PUBLIC_PATHS.includes(pathname);

    useEffect(() => {
        if (isPublicPage) {
            setLoading(false);
            return;
        }

        const userData = TokenStorage.getUser();
        if (!userData) {
            router.push('/storefront/login');
            setLoading(false);
            return;
        }

        const permissions = TokenStorage.getPermissions();
        const roleName = permissions?.role_name || '';
        const isClienteFinal = roleName === 'cliente_final';

        setUser(userData);
        setIsAdmin(!isClienteFinal);
        setAuthenticated(true);
        setLoading(false);
    }, [router, isPublicPage, pathname]);

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
                <div className="text-gray-500">Cargando...</div>
            </div>
        );
    }

    if (isPublicPage) {
        return <div className="min-h-screen bg-gray-50 dark:bg-gray-900">{children}</div>;
    }

    if (!authenticated) {
        return null;
    }

    // Admin / Super Admin: show full admin layout with sidebar + subnavbar
    if (isAdmin) {
        return (
            <ToastProvider>
                <PermissionsProvider>
                    <NavbarProvider>
                        <SidebarProvider>
                            <AdminStorefrontContent user={user}>
                                {children}
                            </AdminStorefrontContent>
                        </SidebarProvider>
                    </NavbarProvider>
                </PermissionsProvider>
            </ToastProvider>
        );
    }

    // Cliente final: simple storefront nav, no sidebar
    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
            <StorefrontNav />
            <main className="max-w-7xl mx-auto px-4 py-6">
                {children}
            </main>
        </div>
    );
}
