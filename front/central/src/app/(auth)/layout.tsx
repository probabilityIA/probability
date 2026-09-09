
'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { TokenStorage } from '@/shared/config';
import { Spinner, ShopifyIframeDetector } from '@/shared/ui';
import { ToastProvider } from '@/shared/providers/toast-provider';
import { SidebarProvider } from '@/shared/contexts/sidebar-context';
import { PermissionsProvider } from '@/shared/contexts/permissions-context';
import { NavbarProvider } from '@/shared/contexts/navbar-context';
import { useShopifyAuth } from '@/providers/ShopifyAuthProvider';
import LayoutContent from './layout-content';

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const { isShopifyEmbedded, sessionToken: shopifySessionToken, isLoading: isShopifyLoading } = useShopifyAuth();
  const [user, setUser] = useState<{ userId: string; name: string; email: string; role: string; avatarUrl?: string; is_super_admin?: boolean; scope?: string } | null>(null);
  const [loading, setLoading] = useState(true);
  const [showBusinessSelector] = useState(false);

  const isLoginPage = pathname === '/login';
  const isPublicPage = isLoginPage || pathname === '/verify-email' || pathname === '/verify-demo' || pathname === '/forgot-password' || pathname === '/reset-password' || pathname === '/verify-code' || pathname === '/auth/google/callback' || pathname === '/registro-demo';

  useEffect(() => {
    if (isShopifyEmbedded && isShopifyLoading) {
      return;
    }

    if (!isPublicPage) {
      try {
        const userData = TokenStorage.getUser();

        if (!userData) {
          console.warn('⚠️ No user data, redirecting to login');
          const destino = typeof window !== 'undefined'
            ? `${window.location.pathname}${window.location.search}`
            : pathname;
          const esRaiz = !destino || destino === '/' || destino === '/home';
          router.push(esRaiz ? '/login' : `/login?next=${encodeURIComponent(destino)}`);
          return;
        }

        const isSuperAdmin = userData.is_super_admin || false;
        const scope = userData.scope || '';
        const businessesData = TokenStorage.getBusinessesData();
        const isBusinessUser = scope === 'business';

        if (isBusinessUser && !isSuperAdmin) {
          if (!businessesData || businessesData.length === 0) {
            console.error('❌ Usuario business sin negocios asignados');
            TokenStorage.clearSession();
            router.push('/login?error=no_business');
            return;
          }
        }

        setTimeout(() => {
          setUser(userData);
          setLoading(false);
        }, 0);
      } catch (error) {
        console.error('❌ Error checking authentication:', error);
        router.push('/login');
      }
    } else {
      setTimeout(() => setLoading(false), 0);
    }
  }, [router, isPublicPage, pathname, isShopifyEmbedded, isShopifyLoading, shopifySessionToken]);



  if (showBusinessSelector && !isPublicPage) {
    const businessesData = TokenStorage.getBusinessesData();
    if (businessesData && businessesData.length > 0) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white">
          <div className="text-center">
            <h2 className="text-xl font-bold mb-4">Seleccionar Negocio</h2>
            <p>El componente de selección de negocio está en migración.</p>
          </div>
        </div>
      );
    }
  }

  if ((loading || !user) && !isPublicPage) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <Spinner size="xl" color="primary" text={isShopifyEmbedded ? "Conectando con Shopify..." : "Cargando..."} />
          {isShopifyEmbedded && (
            <p className="mt-4 text-sm text-gray-600 dark:text-gray-300">
              🛍️ Inicializando integración de Shopify
            </p>
          )}
        </div>
      </div>
    );
  }

  if (isPublicPage) {
    return (
      <ShopifyIframeDetector>
        {children}
      </ShopifyIframeDetector>
    );
  }

  return (
    <ShopifyIframeDetector>
      <ToastProvider>
        <PermissionsProvider>
          <NavbarProvider>
            <SidebarProvider>
              <LayoutContent user={user}>
                {children}
              </LayoutContent>
            </SidebarProvider>
          </NavbarProvider>
        </PermissionsProvider>
      </ToastProvider>
    </ShopifyIframeDetector>
  );
}
