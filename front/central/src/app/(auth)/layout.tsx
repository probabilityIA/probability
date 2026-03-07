/**
 * Layout para páginas autenticadas
 * Incluye el sidebar de navegación
 */

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
// import { BusinessSelector } from '@modules/auth/ui';

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

  // Páginas que NO deben tener sidebar (login)
  const isLoginPage = pathname === '/login';

  useEffect(() => {
    // Esperar a que Shopify Auth termine de cargar si estamos en iframe
    if (isShopifyEmbedded && isShopifyLoading) {
      return;
    }

    // Verificar autenticación (solo si no es login)
    if (!isLoginPage) {
      try {
        // ✅ NO verificar token (cookie HttpOnly se envía automáticamente)
        // Solo verificar que haya datos del usuario en sessionStorage
        const userData = TokenStorage.getUser();

        if (!userData) {
          console.warn('⚠️ No user data, redirecting to login');
          router.push('/login');
          setTimeout(() => setLoading(false), 0);
          return;
        }

        // Si el usuario es business y NO es super admin, debe tener business token
        const isSuperAdmin = userData.is_super_admin || false;
        const scope = userData.scope || '';
        const businessesData = TokenStorage.getBusinessesData();
        const isBusinessUser = scope === 'business';

        // Usuario business: validación básica
        if (isBusinessUser && !isSuperAdmin) {
          // Verificar si tiene negocios asignados
          if (!businessesData || businessesData.length === 0) {
            // No tiene negocios, redirigir al login con mensaje
            console.error('❌ Usuario business sin negocios asignados');
            TokenStorage.clearSession();
            router.push('/login?error=no_business');
            setTimeout(() => setLoading(false), 0);
            return;
          }
        }

        // Check if user is cliente_final (level 5) - redirect to storefront
        const permissions = TokenStorage.getPermissions();
        if (permissions?.role_name === 'cliente_final') {
          router.push('/storefront/catalogo');
          return;
        }

        // Todo OK, setear usuario
        setTimeout(() => {
          setUser(userData);
          setLoading(false);
        }, 0);
      } catch (error) {
        console.error('❌ Error checking authentication:', error);
        // En caso de error (ej: localStorage bloqueado en iframe), redirigir a login
        router.push('/login');
        setTimeout(() => setLoading(false), 0);
      }
    } else {
      setTimeout(() => setLoading(false), 0);
    }
  }, [router, isLoginPage, pathname, isShopifyEmbedded, isShopifyLoading, shopifySessionToken]);



  // Si debe mostrar el selector de negocios
  if (showBusinessSelector && !isLoginPage) {
    const businessesData = TokenStorage.getBusinessesData();
    if (businessesData && businessesData.length > 0) {
      // TODO: Migrar BusinessSelector a la nueva arquitectura
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white">
          <div className="text-center">
            <h2 className="text-xl font-bold mb-4">Seleccionar Negocio</h2>
            <p>El componente de selección de negocio está en migración.</p>
            {/*
            <BusinessSelector
              businesses={mappedBusinesses}
              isOpen={true}
              onClose={handleBusinessSelected}
              showSuperAdminButton={false}
              skipRedirect={true}
            />
            */}
          </div>
        </div>
      );
    }
  }

  if (loading && !isLoginPage) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <Spinner size="xl" color="primary" text={isShopifyEmbedded ? "Conectando con Shopify..." : "Cargando..."} />
          {isShopifyEmbedded && (
            <p className="mt-4 text-sm text-gray-600">
              🛍️ Inicializando integración de Shopify
            </p>
          )}
        </div>
      </div>
    );
  }

  // Si es la página de login, renderizar sin sidebar
  if (isLoginPage) {
    return (
      <ShopifyIframeDetector>
        {children}
      </ShopifyIframeDetector>
    );
  }

  // Páginas autenticadas con sidebar
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
