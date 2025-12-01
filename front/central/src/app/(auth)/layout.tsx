/**
 * Layout para páginas autenticadas
 * Incluye el sidebar de navegación
 */

'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { TokenStorage } from '@/shared/config';
import { Sidebar, Spinner } from '@/shared/ui';
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
  const [sidebarExpanded, setSidebarExpanded] = useState(false);
  const [showBusinessSelector, setShowBusinessSelector] = useState(false);

  // Páginas que NO deben tener sidebar (login)
  const isLoginPage = pathname === '/login';

  useEffect(() => {
    // Verificar autenticación (solo si no es login)
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

      // Si es super admin y no tiene business token, generarlo automáticamente
      if (isSuperAdmin && !businessToken) {
        console.log('🔑 Auth Layout - Generando business token para super admin');
        (async () => {
          try {
            const { generateBusinessTokenAction } = await import('@/services/auth/login/infra/actions');
            const result = await generateBusinessTokenAction({ business_id: 0 }, sessionToken);
            if (result.success && result.data) {
              TokenStorage.setBusinessToken(result.data.token);
              TokenStorage.removeUserPermissions(); // Limpiar permisos anteriores
              TokenStorage.setActiveBusiness(0);
              console.log('✅ Business token generado para super admin');
              window.location.reload();
              return;
            }
          } catch (err) {
            console.error('Error generando business token para super admin:', err);
          }
        })();
        return;
      }

      // Usuario business que NO es super admin: requiere business token
      if (isBusinessUser && !isSuperAdmin) {
        // Verificar si tiene negocios asignados
        if (!businessesData || businessesData.length === 0) {
          // No tiene negocios, redirigir al login con mensaje
          console.error('❌ Usuario business sin negocios asignados');
          TokenStorage.clearSession();
          router.push('/login?error=no_business');
          return;
        }

        // Tiene negocios pero no tiene business token: mostrar selector
        if (!businessToken) {
          setShowBusinessSelector(true);
          setLoading(false);
          return;
        }
      }

      setUser(userData);
    }

    setLoading(false);
  }, [router, isLoginPage, pathname]);

  const handleBusinessSelected = () => {
    setShowBusinessSelector(false);
    // Recargar la página para asegurar que todos los componentes usen el nuevo token
    window.location.reload();
  };

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
      <div className="min-h-screen flex items-center justify-center">
        <Spinner size="xl" color="primary" text="Cargando..." />
      </div>
    );
  }

  // Si es la página de login, renderizar sin sidebar
  if (isLoginPage) {
    return <>{children}</>;
  }

  // Páginas autenticadas con sidebar
  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Sidebar */}
      <Sidebar user={user} />

      {/* Contenido principal */}
      <main
        className="flex-1 transition-all duration-300"
        style={{
          marginLeft: sidebarExpanded ? '250px' : '80px'
        }}
        onMouseEnter={() => setSidebarExpanded(false)}
      >
        {children}
      </main>
    </div>
  );
}

