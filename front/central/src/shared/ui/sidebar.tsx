
'use client';

import { useState, useEffect, useMemo } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { TokenStorage } from '@/shared/config';
import { useSidebar } from '@/shared/contexts/sidebar-context';
import { UserProfileModal } from './user-profile-modal';
import { BusinessSettingsModal } from '@/services/auth/business/ui/components/BusinessSettingsModal';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { getMyModulesAction } from '@/services/modules/wallet/infra/subscription-actions';

interface SidebarProps {
  user: {
    userId: string;
    name: string;
    email: string;
    role: string;
    avatarUrl?: string;
  } | null;
}

export function Sidebar({ user }: SidebarProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { primaryExpanded, requestExpand, requestCollapse, isMobileOpen, setIsMobileOpen } = useSidebar();
  const [showUserModal, setShowUserModal] = useState(false);
  const [invoicingOpen, setInvoicingOpen] = useState(false);
  const { hasPermission, isSuperAdmin, isLoading, permissions } = usePermissions();
  const isDemo = permissions?.role_name === 'demo';

  const [accessibleModules, setAccessibleModules] = useState<string[] | null>(null);
  useEffect(() => {
    if (isSuperAdmin) return;
    getMyModulesAction().then((res) => {
      if (res.success && res.data) setAccessibleModules(res.data);
    });
  }, [isSuperAdmin]);
  const hasSubscriptionModule = (code: string) => isSuperAdmin || (accessibleModules?.includes(code) ?? false);

  const businessLogo = useMemo(() => {
    if (isSuperAdmin) return null;
    const businesses = TokenStorage.getBusinessesData();
    if (!businesses || !permissions?.business_id) return null;
    const active = businesses.find(b => b.id === permissions.business_id);
    return active?.logo_url || null;
  }, [isSuperAdmin, permissions]);

  const [superAdminBusinessLogo, setSuperAdminBusinessLogo] = useState<string | null>(null);
  const [superAdminBusinessId, setSuperAdminBusinessId] = useState<number | null>(null);
  const [showBusinessModal, setShowBusinessModal] = useState(false);

  useEffect(() => {
    if (!isSuperAdmin) return;

    const updateLogo = () => {
      setSuperAdminBusinessLogo(localStorage.getItem('selected_business_logo'));
      const stored = localStorage.getItem('selected_business_id');
      const parsed = stored ? parseInt(stored, 10) : NaN;
      setSuperAdminBusinessId(Number.isNaN(parsed) ? null : parsed);
    };

    updateLogo();
    window.addEventListener('businessChanged', updateLogo);
    window.addEventListener('storage', updateLogo);
    return () => {
      window.removeEventListener('businessChanged', updateLogo);
      window.removeEventListener('storage', updateLogo);
    };
  }, [isSuperAdmin]);

  const activeBusinessId = isSuperAdmin ? superAdminBusinessId : permissions?.business_id ?? null;

  useEffect(() => {
    if (!primaryExpanded) {
      setInvoicingOpen(false);
    }
  }, [primaryExpanded]);

  useEffect(() => {
    setIsMobileOpen(false);
  }, [pathname, setIsMobileOpen]);

  const iamRoutes = ['/users', '/roles', '/permissions', '/businesses', '/resources'];
  const ordersRoutes = ['/orders', '/shipments'];
  const inventoryRoutes = ['/products', '/warehouses', '/inventory'];
  const invoicingRoutes = ['/invoicing'];
  const deliveryRoutes = ['/delivery'];
  const hasSecondarySidebar = iamRoutes.some(route => pathname.startsWith(route)) ||
    ordersRoutes.some(route => pathname.startsWith(route)) ||
    inventoryRoutes.some(route => pathname.startsWith(route)) ||
    invoicingRoutes.some(route => pathname.startsWith(route)) ||
    deliveryRoutes.some(route => pathname.startsWith(route));

  const canViewResources = isSuperAdmin;
  const canViewBusinesses = isSuperAdmin || hasPermission('Empresas', 'Read');

  const canViewUsers = isSuperAdmin || hasPermission('Usuarios', 'Read') || hasPermission('Users', 'Read') || hasPermission('Empleados', 'Read');
  const canViewRoles = isSuperAdmin || hasPermission('Roles', 'Read') || hasPermission('Roles y Permisos', 'Read');
  const canViewPermissions = isSuperAdmin || hasPermission('Permisos', 'Read') || hasPermission('Permissions', 'Read');

  const canViewProducts = isSuperAdmin || hasPermission('Productos', 'Read') || hasPermission('Products', 'Read');
  const canViewOrders = isSuperAdmin || hasPermission('Ordenes', 'Read') || hasPermission('Orders', 'Read');
  const canViewShipments = isSuperAdmin || hasPermission('Envios', 'Read') || hasPermission('Shipments', 'Read');

  const canViewCustomers = isSuperAdmin || hasPermission('Clientes', 'Read') || hasPermission('Customers', 'Read');

  const canViewAnnouncements = isSuperAdmin;
  const canViewTickets = !isDemo && hasSubscriptionModule('tickets');

  const canViewWarehouses = isSuperAdmin || hasPermission('Bodegas', 'Read') || hasPermission('Warehouses', 'Read');
  const canViewInventory = isSuperAdmin
    || hasPermission('Inventario', 'Read') || hasPermission('Inventory', 'Read')
    || hasPermission('Inventario-Stock', 'Read')
    || hasPermission('Inventario-Movimientos', 'Read')
    || hasPermission('Inventario-Trazabilidad', 'Read')
    || hasPermission('Inventario-Kardex', 'Read')
    || hasPermission('Inventario-Operaciones', 'Read')
    || hasPermission('Inventario-Slotting', 'Read')
    || hasPermission('Inventario-Auditoria', 'Read')
    || hasPermission('Inventario-LPN', 'Read')
    || hasPermission('Inventario-Scan', 'Read')
    || hasPermission('Inventario-Sync-Logs', 'Read');

  const canViewNotifications = isSuperAdmin || hasPermission('Notificaciones', 'Read');

  const canViewWallet = isSuperAdmin || hasPermission('Billetera', 'Read');

  const canViewIntegrations = isSuperAdmin || user?.role === 'Administrador' || hasPermission('Integraciones', 'Read') || hasPermission('Integrations', 'Read');

  const canViewInvoices = isSuperAdmin || hasPermission('Facturacion', 'Read');
  const canViewInvoicingProviders = isSuperAdmin || hasPermission('Facturacion', 'Read');
  const canViewInvoicingConfigs = isSuperAdmin || hasPermission('Facturacion', 'Read');

  const canViewAccounting = isSuperAdmin;

  const canViewStorefront = (isSuperAdmin || hasPermission('Storefront', 'Read')) && hasSubscriptionModule('storefront');
  const canViewWebsiteConfig = isSuperAdmin || user?.role === 'Administrador';

  const canViewDelivery = (isSuperAdmin || hasPermission('Ultima Milla', 'Read') || hasPermission('Delivery', 'Read')) && hasSubscriptionModule('delivery');

  const canAccessIAM = canViewBusinesses || canViewUsers || canViewRoles || canViewPermissions || canViewResources;
  const canAccessOrders = canViewOrders || canViewShipments;
  const canAccessInventory = canViewProducts || canViewWarehouses || canViewInventory;
  const canAccessInvoicing = canViewInvoices || canViewInvoicingProviders || canViewInvoicingConfigs;

  const getIAMEntryRoute = () => {
    if (canViewBusinesses) return '/businesses';
    if (canViewUsers) return '/users';
    if (canViewRoles) return '/roles';
    if (canViewPermissions) return '/permissions';
    if (canViewResources) return '/resources';
    return '/businesses';
  };

  const getOrdersEntryRoute = () => {
    if (canViewOrders) return '/orders';
    if (canViewProducts) return '/products';
    if (canViewShipments) return '/shipments';
    if (canViewNotifications) return '/notification-config';
    return '/orders';
  };

  const getInventoryEntryRoute = () => {
    if (canViewProducts) return '/products';
    if (canViewWarehouses) return '/warehouses';
    if (canViewInventory) return '/inventory';
    return '/products';
  };

  const getInvoicingEntryRoute = () => {
    if (canViewInvoices) return '/invoicing/invoices';
    if (canViewInvoicingProviders) return '/invoicing/providers';
    if (canViewInvoicingConfigs) return '/invoicing/configs';
    return '/invoicing/invoices';
  };

  const handleLogout = () => {
    TokenStorage.clearSession();
    router.push('/login');
  };

  if (!user) return null;

  const isActive = (path: string) => pathname === path;

  return (
    <>
      <button
        onClick={() => {
          const newState = !isMobileOpen;
          setIsMobileOpen(newState);
          if (newState) {
            requestExpand();
          }
        }}
        className="fixed top-4 right-4 z-40 md:hidden p-3 bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-100 dark:border-gray-700 text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-all active:scale-95"
        aria-label="Toggle Menu"
      >
        {isMobileOpen ? (
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        ) : (
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        )}
      </button>

      {isMobileOpen && (
        <div
          className="fixed inset-0 bg-black/40 backdrop-blur-sm z-20 md:hidden"
        />
      )}

      <aside
        className={`
          fixed left-0 top-0 h-full transition-all duration-300 z-30 border-r border-gray-200 dark:border-gray-700 sidebar-surface rounded-tr-lg rounded-br-lg
          ${isMobileOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'}
        `}
        style={{
          width: primaryExpanded ? '200px' : '80px',
        }}
        onMouseEnter={() => {
          if (typeof window !== 'undefined' && window.innerWidth >= 768) {
            requestExpand();
          }
        }}
        onMouseLeave={() => {
          if (typeof window !== 'undefined' && window.innerWidth >= 768) {
            requestCollapse(hasSecondarySidebar);
          }
        }}
      >
        <div className="flex flex-col h-full">
          <div className="flex flex-col items-center py-4 gap-2">
            {businessLogo ? (
              <button
                type="button"
                onClick={() => setShowBusinessModal(true)}
                title="Informacion del negocio"
                className={`relative transition-all duration-300 flex items-center justify-center cursor-pointer hover:opacity-80 ${primaryExpanded ? 'w-56 h-10' : 'w-8 h-8'}`}
              >
                <img
                  src={businessLogo}
                  alt="Business Logo"
                  className={`object-contain transition-all duration-300 ${primaryExpanded ? 'max-w-full max-h-full' : 'w-8 h-8 rounded'}`}
                />
              </button>
            ) : (
              <a
                href="https://www.probabilityia.com.co/"
                target="_blank"
                rel="noopener noreferrer"
                className={`relative transition-all duration-300 flex items-center justify-center cursor-pointer hover:opacity-80 ${primaryExpanded ? 'w-56 h-10' : 'w-8 h-8'}`}
              >
                <Image
                  src={primaryExpanded ? "/logo2recortado.png" : "/logo.ico"}
                  alt="Probability Logo"
                  fill
                  className="object-contain"
                  priority
                  sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                />
              </a>
            )}
            {isSuperAdmin && superAdminBusinessLogo && (
              <button
                type="button"
                onClick={() => setShowBusinessModal(true)}
                title="Informacion del negocio"
                className={`transition-all duration-300 flex items-center justify-center cursor-pointer hover:opacity-80 ${primaryExpanded ? 'w-32 h-8' : 'w-7 h-7'}`}
              >
                <img
                  src={superAdminBusinessLogo}
                  alt="Negocio seleccionado"
                  className={`object-contain transition-all duration-300 ${primaryExpanded ? 'max-w-full max-h-full' : 'w-7 h-7 rounded'}`}
                />
              </button>
            )}
          </div>
          <div className="mx-auto w-[85%] h-[1px] rounded-full bg-gradient-to-r from-transparent via-gray-200 dark:via-gray-600 to-transparent" />

          <button
            type="button"
            onClick={() => setShowUserModal(true)}
            className={`cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors rounded-xl mx-2 my-1 ${primaryExpanded ? 'p-4' : 'p-2 flex justify-center'} block w-[calc(100%-1rem)] text-left`}
            title="Abrir perfil"
          >
            <div className={`flex items-center ${primaryExpanded ? 'gap-3' : 'justify-center'}`}>
              <div className="relative group">
                {user.avatarUrl ? (
                  <img
                    src={user.avatarUrl}
                    alt={user.name}
                    className="w-12 h-12 rounded-full object-cover flex-shrink-0 border-2 border-white/20 transition-all group-hover:border-white/40 group-hover:ring-2 group-hover:ring-white/20"
                  />
                ) : (
                  <div
                    className="w-12 h-12 rounded-full flex items-center justify-center text-white text-lg font-bold flex-shrink-0 transition-all group-hover:ring-2 group-hover:ring-white/20"
                    style={{ backgroundColor: 'var(--color-secondary)' }}
                  >
                    {user.name.charAt(0).toUpperCase()}
                  </div>
                )}
                <div className="absolute inset-0 rounded-full bg-black/0 group-hover:bg-black/20 transition-all flex items-center justify-center">
                  <svg className="w-4 h-4 text-white opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </div>
              </div>

              {primaryExpanded && (
                <div className="text-gray-800 dark:text-gray-100 dark:text-gray-200 overflow-hidden">
                  <p className="font-semibold text-sm truncate">{user.name}</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 truncate">{user.email}</p>
                </div>
              )}
            </div>
          </button>

          <nav className="flex-1 py-6 px-3 overflow-y-auto overflow-x-hidden">
            <ul className="space-y-2">
              <li>
                <Link
                  href="/home"
                  className={`
                    flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                    ${isActive('/home')
                      ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                      : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                    }
                  `}
                >
                  {isActive('/home') && (
                    <div
                      className="absolute left-0 w-1 h-8 rounded-r-full"
                      style={{ backgroundColor: 'var(--color-tertiary)' }}
                    />
                  )}

                  <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                  </svg>
                  {primaryExpanded && (
                    <span className="text-sm font-medium transition-opacity duration-300">
                      Inicio
                    </span>
                  )}
                </Link>
              </li>

              {canViewWallet && (
                <li>
                  <Link
                    href="/wallet"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${isActive('/wallet')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {isActive('/wallet') && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">
                        Billetera
                      </span>
                    )}
                  </Link>
                </li>
              )}

              {canViewIntegrations && (
                <li>
                  <Link
                    href="/integrations"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${isActive('/integrations') || pathname.startsWith('/integrations')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {(isActive('/integrations') || pathname.startsWith('/integrations')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}

                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">
                        Integraciones
                      </span>
                    )}
                  </Link>
                </li>
              )}

              {canAccessInventory && (
                <li>
                  <Link
                    href={getInventoryEntryRoute()}
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/products') || pathname.startsWith('/warehouses') || pathname.startsWith('/inventory')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {(pathname.startsWith('/products') || pathname.startsWith('/warehouses') || pathname.startsWith('/inventory')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">Inventario</span>
                    )}
                  </Link>
                </li>
              )}

              {canAccessOrders && (
                <li>
                  <Link
                    href="/orders"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/orders') || pathname.startsWith('/shipments')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {(pathname.startsWith('/orders') || pathname.startsWith('/shipments')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">{'\u00d3rdenes'}</span>
                    )}
                  </Link>
                </li>
              )}

              {canViewDelivery && (
                <li>
                  <Link
                    href="/delivery/routes"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/delivery')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {pathname.startsWith('/delivery') && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8a1 1 0 011-1h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1m-6-1a1 1 0 001 1h1M5 17a2 2 0 104 0m-4 0a2 2 0 114 0m6 0a2 2 0 104 0m-4 0a2 2 0 114 0" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">{'\u00daltima Milla'}</span>
                    )}
                  </Link>
                </li>
              )}

              {canViewNotifications && (
                <li>
                  <Link
                    href="/notification-config"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/notification-config')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {pathname.startsWith('/notification-config') && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <img src="https://cdn-icons-png.flaticon.com/512/1384/1384162.png" alt="Notificaciones" className="w-5 h-5 flex-shrink-0 object-contain" />
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">Notificaciones</span>
                    )}
                  </Link>
                </li>
              )}

              {canViewCustomers && (
                <li>
                  <Link
                    href="/customers"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/customers')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {pathname.startsWith('/customers') && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}

                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">Clientes</span>
                    )}
                  </Link>
                </li>
              )}

              {canViewTickets && (
                <li>
                  <Link
                    href="/tickets"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/tickets')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {pathname.startsWith('/tickets') && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 8h10M7 12h6m-6 4h10M5 4h14a2 2 0 012 2v12a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">Tickets</span>
                    )}
                  </Link>
                </li>
              )}

              {canViewAnnouncements && (
                <li>
                  <Link
                    href="/announcements"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/announcements')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {pathname.startsWith('/announcements') && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5.882V19.24a1.76 1.76 0 01-3.5 0V5.882m0 0C7.5 4.334 9.167 3 11 3s3.5 1.334 3.5 2.882M11 5.882h3.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">Anuncios</span>
                    )}
                  </Link>
                </li>
              )}

              {(canViewStorefront || canViewWebsiteConfig) && (
                <li>
                  <Link
                    href="/website-config"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/storefront') || pathname.startsWith('/website-config')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {(pathname.startsWith('/storefront') || pathname.startsWith('/website-config')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">Tienda</span>
                    )}
                  </Link>
                </li>
              )}

              {canAccessInvoicing && (
                <li>
                  <Link
                    href="/invoicing/invoices"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/invoicing')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {pathname.startsWith('/invoicing') && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}

                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 14l6-6m-5.5.5h.01m4.99 5h.01M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16l3.5-2 3.5 2 3.5-2 3.5 2zM10 8.5a.5.5 0 11-1 0 .5.5 0 011 0zm5 5a.5.5 0 11-1 0 .5.5 0 011 0z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">{'Facturaci\u00f3n'}</span>
                    )}
                  </Link>
                </li>
              )}

              {canViewAccounting && (
                <li>
                  <Link
                    href="/accounting"
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/accounting')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {pathname.startsWith('/accounting') && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 7h6m-6 4h6m-6 4h3m-7.5 6h15a1.5 1.5 0 001.5-1.5v-15A1.5 1.5 0 0019.5 3h-15A1.5 1.5 0 003 4.5v15A1.5 1.5 0 004.5 21z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">Contabilidad</span>
                    )}
                  </Link>
                </li>
              )}

              {!isDemo && (
              <li>
                <Link
                  href="/subscription"
                  className={`
                    flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                    ${pathname.startsWith('/subscription')
                      ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                      : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                    }
                  `}
                >
                  {pathname.startsWith('/subscription') && (
                    <div
                      className="absolute left-0 w-1 h-8 rounded-r-full"
                      style={{ backgroundColor: 'var(--color-tertiary)' }}
                    />
                  )}
                  <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
                  </svg>
                  {primaryExpanded && (
                    <span className="text-sm font-medium transition-opacity duration-300">
                      {'Suscripci\u00f3n'}
                    </span>
                  )}
                </Link>
              </li>
              )}

              {canAccessIAM && (
                <li>
                  <Link
                    href={getIAMEntryRoute()}
                    className={`
                      flex ${primaryExpanded ? 'items-center' : 'justify-center items-center'} gap-1 px-3 py-1.5 rounded-lg transition-colors duration-300 w-full
                      ${pathname.startsWith('/users') || pathname.startsWith('/roles') || pathname.startsWith('/permissions') || pathname.startsWith('/businesses') || pathname.startsWith('/resources')
                        ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white dark:text-gray-100 shadow-sm'
                        : 'text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:text-white dark:hover:text-gray-100'
                      }
                    `}
                  >
                    {(pathname.startsWith('/users') || pathname.startsWith('/roles') || pathname.startsWith('/permissions') || pathname.startsWith('/businesses') || pathname.startsWith('/resources')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">IAM</span>
                    )}
                  </Link>
                </li>
              )}
            </ul>
          </nav>

          {businessLogo && (
            <a
              href="https://www.probabilityia.com.co/"
              target="_blank"
              rel="noopener noreferrer"
              className={`flex items-center justify-center gap-1.5 py-2 opacity-50 hover:opacity-90 transition-opacity ${primaryExpanded ? '' : 'px-1'}`}
              title="Probability"
            >
              <Image src="/logo.ico" alt="Probability" width={16} height={16} className="object-contain flex-shrink-0" />
              {primaryExpanded && <span className="text-[10px] text-gray-500 dark:text-gray-400 whitespace-nowrap">Powered by Probability</span>}
            </a>
          )}
          <div className="mx-auto w-[85%] h-[1px] rounded-full bg-gradient-to-r from-transparent via-gray-200 dark:via-gray-600 to-transparent mb-2" />
          <div className="p-4 pt-2">
            <button
              onClick={handleLogout}
              className="w-full flex items-center gap-3 text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 p-3 rounded-lg transition-colors"
            >
              <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              {primaryExpanded && <span className="text-sm">{'Cerrar Sesi\u00f3n'}</span>}
            </button>
          </div>
        </div>
      </aside >

      < UserProfileModal
        isOpen={showUserModal}
        onClose={() => setShowUserModal(false)
        }
        user={user}
        onUpdate={() => {
          window.location.reload();
        }}
      />

      <BusinessSettingsModal
        isOpen={showBusinessModal}
        onClose={() => setShowBusinessModal(false)}
        businessId={activeBusinessId}
      />
    </>
  );
}
