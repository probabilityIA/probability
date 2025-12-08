/**
 * Layout para p?ginas autenticadas
 * Incluye el sidebar de navegaci?n
 */

'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { TokenStorage } from '@/shared/config';
import { Spinner } from '@/shared/ui';
import { ToastProvider } from '@/shared/providers/toast-provider';
import { SidebarProvider } from '@/shared/contexts/sidebar-context';
import LayoutContent from './layout-content';
// import { BusinessSelector } from '@modules/auth/ui';

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const [user, setUser] = useState<{ userId: string; name: string; email: string; role: string; avatarUrl?: string; is_super_admin?: boolean; scope?: string } | null>(null);
  const [loading, setLoading] = useState(true);
  const [showBusinessSelector, setShowBusinessSelector] = useState(false);

  // P?ginas que NO deben tener sidebar (login)
  const isLoginPage = pathname === '/login';

  useEffect(() => {
    // Verificar autenticaci?n (solo si no es login)
    if (!isLoginPage) {
      const sessionToken = TokenStorage.getSessionToken();
      const businessToken = TokenStorage.getBusinessToken();
      const userData = TokenStorage.getUser();

      if (!sessionToken || !userData) {
        router.push('/login');
        return;
      }

      // Si el usuario es business y NO es super admin, debe tener business token
      const isSuperAdmin = userData.is_super_admin || false;
      const scope = userData.scope || '';
      const businessesData = TokenStorage.getBusinessesData();
      const isBusinessUser = scope === 'business';

      // Si es super admin, no necesitamos generar token de negocio adicional
      // El token de sesi?n ya tiene los permisos necesarios

      // Usuario business: validaci?n b?sica
      if (isBusinessUser && !isSuperAdmin) {
        // Verificar si tiene negocios asignados
        if (!businessesData || businessesData.length === 0) {
          // No tiene negocios, redirigir al login con mensaje
          console.error('? Usuario business sin negocios asignados');
          TokenStorage.clearSession();
          router.push('/login?error=no_business');
          return;
        }
      }

      setUser(userData);
    }

    setLoading(false);
  }, [router, isLoginPage, pathname]);



  // Si debe mostrar el selector de negocios
  if (showBusinessSelector && !isLoginPage) {
    const businessesData = TokenStorage.getBusinessesData();
    if (businessesData && businessesData.length > 0) {
      // TODO: Migrar BusinessSelector a la nueva arquitectura
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white">
          <div className="text-center">
            <h2 className="text-xl font-bold mb-4">Seleccionar Negocio</h2>
            <p>El componente de selecci?n de negocio est? en migraci?n.</p>
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
      <div className="min-h-screen flex items-center justify-center">
        <Spinner size="xl" color="primary" text="Cargando..." />
      </div>
    );
  }

  // Si es la p?gina de login, renderizar sin sidebar
  if (isLoginPage) {
    return <>{children}</>;
  }

  // P?ginas autenticadas con sidebar
  return (
    <ToastProvider>
      <SidebarProvider>
        <LayoutContent user={user}>
          {children}
        </LayoutContent>
      </SidebarProvider>
    </ToastProvider>
  );
}

