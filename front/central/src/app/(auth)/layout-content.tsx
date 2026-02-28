'use client';

import React from 'react';
import { Sidebar, OrdersSubNavbar, InventorySubNavbar, IntegrationsSubNavbar } from '@/shared/ui';
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
  const {
    primaryExpanded,
    secondaryExpanded,
    requestCollapse,
    setHasSecondarySidebar,
    requestSecondaryCollapse
  } = useSidebar();

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
    // Solo en escritorio para evitar comportamientos extraños en móvil
    if (typeof window !== 'undefined' && window.innerWidth >= 768) {
      requestCollapse(false);
      requestSecondaryCollapse();
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Sidebar Principal */}
      <Sidebar user={user} />

      {/* Sidebar Secundario (IAM) eliminado: contenido integrado en Sidebar principal */}

      {/* Sidebar Secundario (Ordenes) ya está integrado en Sidebar principal */}

      {/* Contenido principal */}
      <main
        className="flex-1 transition-all duration-300 w-full overflow-x-hidden main-content flex flex-col"
        onMouseEnter={handleMainMouseEnter}
      >
        <OrdersSubNavbar />
        <InventorySubNavbar />
        <IntegrationsSubNavbar />
        <div className="w-full min-w-0 flex-1">
          {children}
        </div>
        <style jsx>{`
          .main-content {
            margin-left: 0;
          }
          @media (min-width: 768px) {
            .main-content {
              margin-left: ${totalSidebarWidth}px;
            }
          }
        `}</style>
      </main>
    </div>
  );
}

export default LayoutContent;
