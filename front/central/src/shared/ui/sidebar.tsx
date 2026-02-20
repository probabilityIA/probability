/**
 * Sidebar de navegación
 * Componente compartido para todas las páginas autenticadas
 */

'use client';

import { useState, useEffect, useMemo } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { TokenStorage } from '@/shared/config';
import { useSidebar } from '@/shared/contexts/sidebar-context';
import { UserProfileModal } from './user-profile-modal';
import { usePermissions } from '@/shared/contexts/permissions-context';

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
  const [ordersOpen, setOrdersOpen] = useState(false);
  const [invoicingOpen, setInvoicingOpen] = useState(false);
  const [iamOpen, setIamOpen] = useState(false);
  const { hasPermission, isSuperAdmin, isLoading, permissions } = usePermissions();

  const businessLogo = useMemo(() => {
    if (isSuperAdmin) return null;
    const businesses = TokenStorage.getBusinessesData();
    if (!businesses || !permissions?.business_id) return null;
    const active = businesses.find(b => b.id === permissions.business_id);
    return active?.logo_url || null;
  }, [isSuperAdmin, permissions?.business_id]);

  useEffect(() => {
    // When primary sidebar collapses, ensure submenus collapse too
    if (!primaryExpanded) {
      setOrdersOpen(false);
      setInvoicingOpen(false);
      setIamOpen(false);
    }
  }, [primaryExpanded]);

  useEffect(() => {
    // Cerrar sidebar móvil al cambiar de ruta
    setIsMobileOpen(false);
  }, [pathname, setIsMobileOpen]);

  // Determinar si hay sidebar secundario basado en la ruta actual
  const iamRoutes = ['/users', '/roles', '/permissions', '/businesses', '/resources'];
  const ordersRoutes = ['/products', '/orders', '/shipments', '/order-status', '/notification-config'];
  const invoicingRoutes = ['/invoicing'];
  const hasSecondarySidebar = iamRoutes.some(route => pathname.startsWith(route)) ||
    ordersRoutes.some(route => pathname.startsWith(route)) ||
    invoicingRoutes.some(route => pathname.startsWith(route));

  // Si está cargando, no hay permisos definidos, o resources es null/vacío, mostrar todo por defecto
  // Si está cargando, esperamos (no mostramos nada o mostramos skeleton si se implementara)
  // const permissionsNotLoaded = isLoading || !permissions || !permissions.resources || permissions.resources.length === 0;

  // Verificar permisos para cada módulo

  // Recursos: Solo para super admins (Plataforma)
  const canViewResources = isSuperAdmin;
  // Empresas: Visible para super admins y usuarios de negocio con permiso
  const canViewBusinesses = isSuperAdmin || hasPermission('Empresas', 'Read');

  // IAM Core: Visible para super admins Y administradores de negocio
  // Agregamos variantes de nombres de recursos para robustez
  const canViewUsers = isSuperAdmin || hasPermission('Usuarios', 'Read') || hasPermission('Users', 'Read') || hasPermission('Empleados', 'Read');
  const canViewRoles = isSuperAdmin || hasPermission('Roles', 'Read') || hasPermission('Roles y Permisos', 'Read');
  const canViewPermissions = isSuperAdmin || hasPermission('Permisos', 'Read') || hasPermission('Permissions', 'Read');

  // Orders Module
  const canViewProducts = isSuperAdmin || hasPermission('Productos', 'Read') || hasPermission('Products', 'Read');
  const canViewOrders = isSuperAdmin || hasPermission('Ordenes', 'Read') || hasPermission('Orders', 'Read');
  const canViewShipments = isSuperAdmin || hasPermission('Envios', 'Read') || hasPermission('Shipments', 'Read');

  // Configuración de Ordenes: Solo para super admins (Plataforma)
  const canViewOrderStatus = isSuperAdmin;
  const canViewNotifications = isSuperAdmin;

  // Integraciones: Visible para negocio (para crear integraciones)
  // Integraciones: Visible para negocio (para crear integraciones)
  const canViewIntegrations = isSuperAdmin || user?.role === 'Administrador' || hasPermission('Integraciones', 'Read') || hasPermission('Integrations', 'Read');

  // Facturación: Usa recurso único "Facturacion" de la BD (ID 10)
  const canViewInvoices = isSuperAdmin || hasPermission('Facturacion', 'Read');
  const canViewInvoicingProviders = isSuperAdmin || hasPermission('Facturacion', 'Read');
  const canViewInvoicingConfigs = isSuperAdmin || hasPermission('Facturacion', 'Read');

  // Verificar si tiene acceso a los módulos principales
  const canAccessIAM = canViewBusinesses || canViewUsers || canViewRoles || canViewPermissions || canViewResources;
  const canAccessOrders = canViewProducts || canViewOrders || canViewShipments || canViewOrderStatus || canViewNotifications;
  const canAccessInvoicing = canViewInvoices || canViewInvoicingProviders || canViewInvoicingConfigs;



  // Determinar la ruta de entrada para cada módulo (primera disponible)
  const getIAMEntryRoute = () => {
    if (canViewUsers) return '/users';
    if (canViewRoles) return '/roles';
    if (canViewPermissions) return '/permissions';
    if (canViewBusinesses) return '/businesses';
    if (canViewResources) return '/resources';
    return '/users';
  };

  const getOrdersEntryRoute = () => {
    if (canViewOrders) return '/orders';
    if (canViewProducts) return '/products';
    if (canViewShipments) return '/shipments';
    if (canViewOrderStatus) return '/order-status';
    if (canViewNotifications) return '/notification-config';
    return '/orders';
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

  // Helper para determinar si un link está activo
  const isActive = (path: string) => pathname === path;

  return (
    <>
      {/* Botón Burger - Fijo en la parte superior derecha para móvil */}
      <button
        onClick={() => {
          const newState = !isMobileOpen;
          setIsMobileOpen(newState);
          if (newState) {
            requestExpand();
          }
        }}
        className="fixed top-4 right-4 z-40 md:hidden p-3 bg-white rounded-xl shadow-lg border border-gray-100 text-gray-700 hover:bg-gray-50 transition-all active:scale-95"
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

      {/* Overlay para móvil */}
      {isMobileOpen && (
        <div
          className="fixed inset-0 bg-black/40 backdrop-blur-sm z-20 md:hidden"
        /* El overlay ya no cierra el menú al hacer clic, solo la burger lo hace */
        />
      )}

      {/* Sidebar - Menú lateral expandible */}
      <aside
        className={`
          fixed left-0 top-0 h-full transition-all duration-300 z-30 border-r border-gray-200 bg-white rounded-tr-lg rounded-br-lg
          ${isMobileOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'}
        `}
        style={{
          width: primaryExpanded ? '250px' : '80px',
        }}
        onMouseEnter={() => {
          // Solo expandir por hover si estamos en escritorio
          if (typeof window !== 'undefined' && window.innerWidth >= 768) {
            requestExpand();
          }
        }}
        onMouseLeave={() => {
          // Solo colapsar por hover si estamos en escritorio
          if (typeof window !== 'undefined' && window.innerWidth >= 768) {
            requestCollapse(hasSecondarySidebar);
          }
        }}
      >
        <div className="flex flex-col h-full">
          {/* Logo */}
          <div className="flex items-center justify-center py-4 transition-all duration-300">
            <div className={`relative transition-all duration-300 flex items-center justify-center ${primaryExpanded ? 'w-56 h-10' : 'w-8 h-8'}`}>
              {businessLogo ? (
                <img
                  src={businessLogo}
                  alt="Business Logo"
                  className={`object-contain transition-all duration-300 ${primaryExpanded ? 'max-w-full max-h-full' : 'w-8 h-8 rounded'}`}
                />
              ) : (
                <Image
                  src={primaryExpanded ? "/logo2recortado.png" : "/logo.ico"}
                  alt="Probability Logo"
                  fill
                  className="object-contain"
                  priority
                  sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                />
              )}
            </div>
          </div>
          <div className="mx-auto w-[85%] h-[1px] rounded-full bg-gradient-to-r from-transparent via-gray-200 to-transparent" />

          {/* Tarjeta de usuario arriba */}
          <Link
            href="/profile"
            className={`cursor-pointer hover:bg-gray-50 transition-colors rounded-xl mx-2 my-1 ${primaryExpanded ? 'p-4' : 'p-2 flex justify-center'} block`}
            title="Ver perfil completo"
          >
            <div className={`flex items-center ${primaryExpanded ? 'gap-3' : 'justify-center'}`}>
              {/* Avatar clickeable */}
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
                {/* Indicador de que es clickeable */}
                <div className="absolute inset-0 rounded-full bg-black/0 group-hover:bg-black/20 transition-all flex items-center justify-center">
                  <svg className="w-4 h-4 text-white opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </div>
              </div>

              {/* Nombre (solo visible cuando está expandido) */}
              {primaryExpanded && (
                <div className="text-gray-800 overflow-hidden">
                  <p className="font-semibold text-sm truncate">{user.name}</p>
                  <p className="text-xs text-gray-500 truncate">{user.email}</p>
                </div>
              )}
            </div>
          </Link>

          {/* Menú de navegación */}
          <nav className="flex-1 py-6 px-3">
            <ul className="space-y-2">
              {/* Item Home  visible */}
              <li>
                <Link
                  href="/home"
                  className={`
                    flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                    ${isActive('/home')
                      ? 'bg-gray-100 text-gray-900 shadow-sm scale-105'
                      : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 hover:scale-105'
                    }
                  `}
                >
                  {/* Indicador activo (barra lateral) */}
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

              {/* Item Integraciones - Solo si tiene permiso */}
              {canViewIntegrations && (
                <li>
                  <Link
                    href="/integrations"
                    className={`
                      flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                      ${isActive('/integrations') || pathname.startsWith('/integrations')
                        ? 'bg-gray-100 text-gray-900 shadow-sm scale-105'
                        : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 hover:scale-105'
                      }
                    `}
                  >
                    {/* Indicador activo (barra lateral) */}
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

              {/* Item Ordenes (Gestión de Ordenes) - Solo si tiene permiso */}
              {canAccessOrders && (
                <li>
                  <div>
                    <div className="flex items-center justify-between">
                      <button
                        type="button"
                        onClick={() => setOrdersOpen(v => {
                          const nv = !v;
                          if (nv) setIamOpen(false);
                          return nv;
                        })}
                        aria-expanded={ordersOpen}
                        aria-controls="orders-submenu"
                        className={`flex items-center gap-3 p-3 rounded-lg transition-all duration-300 text-left w-full
                          ${isActive('/orders') || isActive('/products') || isActive('/shipments') || isActive('/order-status') || isActive('/notification-config')
                            ? 'bg-gray-100 text-gray-900 shadow-sm scale-105'
                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 hover:scale-105'
                          }
                        `}
                      >
                        {(isActive('/orders') || isActive('/products') || isActive('/shipments') || isActive('/order-status') || isActive('/notification-config')) && (
                          <div
                            className="absolute left-0 w-1 h-8 rounded-r-full"
                            style={{ backgroundColor: 'var(--color-tertiary)' }}
                          />
                        )}

                        <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
                        </svg>
                        {primaryExpanded && (
                          <>
                            <span className="text-sm font-medium transition-opacity duration-300">Ordenes</span>
                            <svg
                              className={`w-4 h-4 transform transition-transform duration-150 ml-auto select-none ${ordersOpen ? '-rotate-90' : 'rotate-90'}`}
                              viewBox="0 0 20 20"
                              fill="none"
                              stroke="currentColor"
                            >
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 6l6 4-6 4" />
                            </svg>
                          </>
                        )}
                      </button>

                      {primaryExpanded && (
                        <Link
                          href={getOrdersEntryRoute()}
                          className="p-2 rounded-md text-gray-500 hover:bg-gray-100 hover:text-gray-700 transition-colors"
                          title="Ir a Ordenes"
                        >

                        </Link>
                      )}
                    </div>

                    {/* Submenu: mostrar solo cuando se haga click para expandir */}
                    {primaryExpanded && ordersOpen && (
                      <div id="orders-submenu" className="mt-2 pl-8 pr-2">
                        {/* CATÁLOGO */}
                        {canViewProducts && (
                          <div className="mb-3">
                            <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">CATÁLOGO</h4>
                            <ul className="space-y-1">
                              <li>
                                <Link
                                  href="/products"
                                  className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/products') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}
                                >
                                  <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                                  </svg>
                                  <span>Productos</span>
                                </Link>
                              </li>
                            </ul>
                          </div>
                        )}

                        {/* OPERACIONES */}
                        {(canViewOrders || canViewShipments) && (
                          <div className="mb-3">
                            <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">OPERACIONES</h4>
                            <ul className="space-y-1">
                              {canViewOrders && (
                                <li>
                                  <Link href="/orders" className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/orders') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}>
                                    <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
                                    </svg>
                                    <span>Ordenes</span>
                                  </Link>
                                </li>
                              )}
                              {canViewShipments && (
                                <>
                                  <li>
                                    <Link href="/shipments" className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/shipments') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}>
                                      <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8a1 1 0 011-1h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1m-6-1a1 1 0 001 1h1M5 17a2 2 0 104 0m-4 0a2 2 0 114 0m6 0a2 2 0 104 0m-4 0a2 2 0 114 0" />
                                      </svg>
                                      <span>Envíos</span>
                                    </Link>
                                  </li>
                                </>
                              )}
                            </ul>
                          </div>
                        )}

                        {/* CONFIGURACIÓN */}
                        {(canViewOrderStatus || canViewNotifications || canViewShipments) && (
                          <div>
                            <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">CONFIGURACIÓN</h4>
                            <ul className="space-y-1">
                              {canViewOrderStatus && (
                                <li>
                                  <Link href="/order-status" className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/order-status') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}>
                                    <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                    </svg>
                                    <span>Estados de Orden</span>
                                  </Link>
                                </li>
                              )}
                              {canViewNotifications && (
                                <li>
                                  <Link href="/notification-config" className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/notification-config') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}>
                                    <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                                    </svg>
                                    <span>Notificaciones</span>
                                  </Link>
                                </li>
                              )}
                              {canViewShipments && (
                                <li>
                                  <Link href="/shipments/origin-addresses" className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/shipments/origin-addresses') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}>
                                    <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                                    </svg>
                                    <span>Direcciones de Origen</span>
                                  </Link>
                                </li>
                              )}
                            </ul>
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                </li>
              )}

              {/* Item Facturación - Link directo sin submenu */}
              {canAccessInvoicing && (
                <li>
                  <Link
                    href="/invoicing/invoices"
                    className={`
                      flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                      ${pathname.startsWith('/invoicing')
                        ? 'bg-gray-100 text-gray-900 shadow-sm scale-105'
                        : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 hover:scale-105'
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
                      <span className="text-sm font-medium transition-opacity duration-300">Facturación</span>
                    )}
                  </Link>
                </li>
              )}

              {/* Item Billetera - Visible para todos */}
              <li>
                <Link
                  href="/wallet"
                  className={`
                    flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                    ${isActive('/wallet')
                      ? 'bg-gray-100 text-gray-900 shadow-sm scale-105'
                      : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 hover:scale-105'
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

              {/* Item IAM (Gestión de Identidad) - Solo si tiene permiso */}
              {canAccessIAM && (
                <li>
                  <div>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center justify-between">
                        <button
                          type="button"
                          onClick={() => setIamOpen(v => {
                            const nv = !v;
                            if (nv) setOrdersOpen(false);
                            return nv;
                          })}
                          aria-expanded={iamOpen}
                          aria-controls="iam-submenu"
                          className={`flex items-center gap-3 p-3 rounded-lg transition-all duration-300 text-left w-full
                            ${isActive('/users') || isActive('/roles') || isActive('/permissions') || isActive('/businesses') || isActive('/resources')
                              ? 'bg-gray-100 text-gray-900 shadow-sm scale-105'
                              : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 hover:scale-105'
                            }
                          `}
                        >
                          {(isActive('/users') || isActive('/roles') || isActive('/permissions') || isActive('/businesses') || isActive('/resources')) && (
                            <div
                              className="absolute left-0 w-1 h-8 rounded-r-full"
                              style={{ backgroundColor: 'var(--color-tertiary)' }}
                            />
                          )}

                          <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                          </svg>
                          {primaryExpanded && (
                            <>
                              <span className="text-sm font-medium transition-opacity duration-300">IAM</span>
                              <svg
                                className={`w-4 h-4 transform transition-transform duration-150 ml-auto select-none ${iamOpen ? '-rotate-90' : 'rotate-90'}`}
                                viewBox="0 0 20 20"
                                fill="none"
                                stroke="currentColor"
                              >
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 6l6 4-6 4" />
                              </svg>
                            </>
                          )}
                        </button>

                        {primaryExpanded && (
                          <Link
                            href={getIAMEntryRoute()}
                            className="p-2 rounded-md text-gray-500 hover:bg-gray-100 hover:text-gray-700 transition-colors"
                            title="Ir a IAM"
                          >

                          </Link>
                        )}
                      </div>
                    </div>

                    {/* Submenu IAM: mostrar solo cuando se haga click para expandir */}
                    {primaryExpanded && iamOpen && (
                      <div id="iam-submenu" className="mt-2 pl-8 pr-2">
                        {/* ORGANIZACIÓN */}
                        {canViewBusinesses && (
                          <div className="mb-3">
                            <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">ORGANIZACIÓN</h4>
                            <ul className="space-y-1">
                              <li>
                                <Link
                                  href="/businesses"
                                  className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/businesses') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}
                                >
                                  <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                                  </svg>
                                  <span>Empresas</span>
                                </Link>
                              </li>
                            </ul>
                          </div>
                        )}

                        {/* CONTROL DE ACCESO */}
                        {(canViewUsers || canViewRoles || canViewPermissions) && (
                          <div className="mb-3">
                            <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">CONTROL DE ACCESO</h4>
                            <ul className="space-y-1">
                              {canViewUsers && (
                                <li>
                                  <Link
                                    href="/users"
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/users') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}
                                  >
                                    <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                                    </svg>
                                    <span>Usuarios</span>
                                  </Link>
                                </li>
                              )}
                              {canViewRoles && (
                                <li>
                                  <Link
                                    href="/roles"
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/roles') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}
                                  >
                                    <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                                    </svg>
                                    <span>Roles</span>
                                  </Link>
                                </li>
                              )}
                              {canViewPermissions && (
                                <li>
                                  <Link
                                    href="/permissions"
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/permissions') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}
                                  >
                                    <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                                    </svg>
                                    <span>Permisos</span>
                                  </Link>
                                </li>
                              )}
                            </ul>
                          </div>
                        )}

                        {/* SISTEMA - Solo super admin */}
                        {canViewResources && (
                          <div>
                            <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">SISTEMA</h4>
                            <ul className="space-y-0.5">
                              <li>
                                <Link
                                  href="/resources"
                                  className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${isActive('/resources') ? 'bg-gray-100 text-gray-900' : 'text-gray-700 hover:bg-gray-50'}`}
                                >
                                  <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
                                  </svg>
                                  <span>Recursos</span>
                                </Link>
                              </li>
                            </ul>
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                </li>
              )}
            </ul>
          </nav>

          {/* Botón logout abajo */}
          <div className="mx-auto w-[85%] h-[1px] rounded-full bg-gradient-to-r from-transparent via-gray-200 to-transparent mb-2" />
          <div className="p-4 pt-2">
            <button
              onClick={handleLogout}
              className="w-full flex items-center gap-3 text-gray-700 hover:bg-gray-50 p-3 rounded-lg transition-colors"
            >
              <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              {primaryExpanded && <span className="text-sm">Cerrar Sesión</span>}
            </button>
          </div>
        </div>
      </aside >

      {/* Modal para cambiar foto de perfil */}
      < UserProfileModal
        isOpen={showUserModal}
        onClose={() => setShowUserModal(false)
        }
        user={user}
        onUpdate={() => {
          // Recargar la página para actualizar el avatar en el sidebar
          window.location.reload();
        }}
      />
    </>
  );
}
