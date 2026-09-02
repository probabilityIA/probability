'use client';

import React from 'react';
import { usePathname } from 'next/navigation';
import { TourProvider, TourLauncherFloating } from '@/services/modules/tours/ui';
import { Sidebar, OrdersSubNavbar, TicketsSubNavbar, InventorySubNavbar, IntegrationsSubNavbar, NotificationsSubNavbar, InvoicingSubNavbar, DeliverySubNavbar, StorefrontSubNavbar, IAMSubNavbar, WalletSubNavbar, AccountingSubNavbar } from '@/shared/ui';
import { SelectedBusinessProvider } from '@/shared/contexts/selected-business-context';
import { useSidebar } from '@/shared/contexts/sidebar-context';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { StorefrontNav } from '@/services/modules/storefront/ui/components/StorefrontNav';
import { SubscriptionGuard } from '@/shared/ui/SubscriptionGuard';
import AnnouncementModal from '@/services/modules/announcements/ui/components/AnnouncementModal';
import AnnouncementTicker from '@/services/modules/announcements/ui/components/AnnouncementTicker';
import LegalAcceptanceGate from '@/services/modules/legal/ui/components/LegalAcceptanceGate';

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

  const showSecondarySidebar = false;

  React.useEffect(() => {
    setHasSecondarySidebar(showSecondarySidebar);
  }, [showSecondarySidebar, setHasSecondarySidebar]);

  const primaryWidth = primaryExpanded ? 250 : 80;
  const secondaryWidth = showSecondarySidebar ? (secondaryExpanded ? 240 : 60) : 0;
  const totalSidebarWidth = primaryWidth + secondaryWidth;

  const handleMainMouseEnter = () => {
    if (typeof window !== 'undefined' && window.innerWidth >= 768) {
      requestCollapse(false);
      requestSecondaryCollapse();
    }
  };

  if (isClienteFinal && isStorefrontPath) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <StorefrontNav />
        <main className="max-w-7xl mx-auto px-4 py-6">
          {children}
        </main>
        <LegalAcceptanceGate />
      </div>
    );
  }

  return (
    <div className="flex min-h-screen bg-gray-50 dark:bg-gray-900">
      <Sidebar user={user} />

      <main
        className="flex-1 transition-all duration-300 w-full overflow-x-hidden main-content flex flex-col"
        onMouseEnter={handleMainMouseEnter}
      >
        <SelectedBusinessProvider>
          <TourProvider>
                      <AnnouncementTicker />
                      <OrdersSubNavbar />
                      <TicketsSubNavbar />
                      <InventorySubNavbar />
                      <IntegrationsSubNavbar />
                      <NotificationsSubNavbar />
                      <InvoicingSubNavbar />
                      <DeliverySubNavbar />
                      <StorefrontSubNavbar />
                      <IAMSubNavbar />
                      <WalletSubNavbar />
                      <AccountingSubNavbar />
                      <div className="w-full min-w-0 flex-1">
                        <SubscriptionGuard>
                          {children}
                        </SubscriptionGuard>
                      </div>
                      {!pathname.startsWith('/tickets') && <TourLauncherFloating />}
          </TourProvider>
        </SelectedBusinessProvider>
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

      <AnnouncementModal />
      <LegalAcceptanceGate />
    </div>
  );
}

export default LayoutContent;
