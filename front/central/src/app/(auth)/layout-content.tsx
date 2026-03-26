'use client';

import React from 'react';
import { usePathname } from 'next/navigation';
import { Sidebar, OrdersSubNavbar, InventorySubNavbar, IntegrationsSubNavbar, NotificationsSubNavbar, InvoicingSubNavbar, DeliverySubNavbar, StorefrontSubNavbar, IAMSubNavbar } from '@/shared/ui';
import { useSidebar } from '@/shared/contexts/sidebar-context';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { InventoryBusinessProvider } from '@/shared/contexts/inventory-business-context';
import { NotificationBusinessProvider } from '@/shared/contexts/notification-business-context';
import { InvoicingBusinessProvider } from '@/shared/contexts/invoicing-business-context';
import { OrdersBusinessProvider } from '@/shared/contexts/orders-business-context';
import { DeliveryBusinessProvider } from '@/shared/contexts/delivery-business-context';
import { StorefrontBusinessProvider } from '@/shared/contexts/storefront-business-context';
import { StorefrontNav } from '@/services/modules/storefront/ui/components/StorefrontNav';
import { LinaChatbot } from '@/shared/ui/LinaChatbot';
import { SubscriptionGuard } from '@/shared/ui/SubscriptionGuard';

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
  const { permissions } = usePermissions();
  const isClienteFinal = permissions?.role_name === 'cliente_final';
  const isStorefrontPath = pathname.startsWith('/storefront');

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

  // Cliente final en storefront: layout simple sin sidebar
  if (isClienteFinal && isStorefrontPath) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <StorefrontNav />
        <main className="max-w-7xl mx-auto px-4 py-6">
          {children}
        </main>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Sidebar Principal */}
      <Sidebar user={user} />

      {/* Sidebar Secundario (IAM) eliminado: contenido integrado en Sidebar principal */}

      {/* Sidebar Secundario (Ordenes) ya está integrado en Sidebar principal */}

      {/* Contenido principal */}
      <main
        className="flex-1 transition-all duration-300 w-full overflow-x-hidden main-content flex flex-col"
        onMouseEnter={handleMainMouseEnter}
      >
        <OrdersBusinessProvider>
          <InventoryBusinessProvider>
            <NotificationBusinessProvider>
              <InvoicingBusinessProvider>
                <DeliveryBusinessProvider>
                  <StorefrontBusinessProvider>
                    <OrdersSubNavbar />
                    <InventorySubNavbar />
                    <IntegrationsSubNavbar />
                    <NotificationsSubNavbar />
                    <InvoicingSubNavbar />
                    <DeliverySubNavbar />
                    <StorefrontSubNavbar />
                    <IAMSubNavbar />
                    <div className="w-full min-w-0 flex-1">
                      <SubscriptionGuard>
                        {children}
                      </SubscriptionGuard>
                    </div>
                  </StorefrontBusinessProvider>
                </DeliveryBusinessProvider>
              </InvoicingBusinessProvider>
            </NotificationBusinessProvider>
          </InventoryBusinessProvider>
        </OrdersBusinessProvider>
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

      {/* Lina — Asistente Virtual (solo para roles business / super admin) */}
      <LinaChatbot
        userScope={user?.scope}
        isSuperAdmin={user?.is_super_admin}
      />
    </div>
  );
}

export default LayoutContent;
