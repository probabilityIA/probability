'use client';

import React from 'react';
import { usePathname } from 'next/navigation';
import { Sidebar, Spinner } from '@/shared/ui';
import { useSidebar } from '@/shared/contexts/sidebar-context';

interface LayoutContentProps {
  user: {
    userId: string;
    name: string;
    email: string;
    role: string;
    avatarUrl?: string;
    is_super_admin?: boolean;
    scope?: string;
  } | null;
  children: React.ReactNode;
}

function LayoutContent({ user, children }: LayoutContentProps) {
  const pathname = usePathname();
  const { 
    primaryExpanded, 
    secondaryExpanded, 
    requestExpand, 
    requestCollapse, 
    setHasSecondarySidebar,
    requestSecondaryCollapse
  } = useSidebar();

  // Rutas que pertenecen al módulo IAM
  const iamRoutes = ['/users', '/roles', '/permissions', '/businesses', '/resources'];
  const showIAMSidebar = iamRoutes.some(route => pathname.startsWith(route));

  // Rutas que pertenecen al módulo de Ordenes
  const ordersRoutes = ['/products', '/orders', '/shipments', '/order-status', '/notification-config'];
  const showOrdersSidebar = ordersRoutes.some(route => pathname.startsWith(route));

  // No usamos sidebars secundarios separados: todo está integrado en el `Sidebar` principal.
  const showSecondarySidebar = false;

  // Actualizar el contexto cuando cambia el estado del sidebar secundario
  React.useEffect(() => {
    setHasSecondarySidebar(showSecondarySidebar);
  }, [showSecondarySidebar, setHasSecondarySidebar]);

  // Calcular el marginLeft del contenido principal
  const primaryWidth = primaryExpanded ? 250 : 80;
  const secondaryWidth = showSecondarySidebar ? (secondaryExpanded ? 240 : 60) : 0;
  const totalSidebarWidth = primaryWidth + secondaryWidth;

  const handleMainMouseEnter = () => {
    // Cuando el cursor entra al contenido principal, cerrar ambos sidebars
    requestCollapse(false);
    requestSecondaryCollapse();
  };

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Sidebar Principal */}
      <Sidebar user={user} />

      {/* Sidebar Secundario (IAM) eliminado: contenido integrado en Sidebar principal */}

      {/* Sidebar Secundario (Ordenes) ya está integrado en Sidebar principal */}

      {/* Contenido principal */}
      <main
        className="flex-1 transition-all duration-300 w-full overflow-x-hidden"
        style={{
          marginLeft: `${totalSidebarWidth}px`
        }}
        onMouseEnter={handleMainMouseEnter}
      >
        <div className="w-full min-w-0">
          {children}
        </div>
      </main>
    </div>
  );
}

export default LayoutContent;

