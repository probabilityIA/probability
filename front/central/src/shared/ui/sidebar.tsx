/**
 * Sidebar de navegación
 * Componente compartido para todas las páginas autenticadas
 */

'use client';

import { useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import { TokenStorage } from '@/shared/config';
import { useSidebar } from '@/shared/contexts/sidebar-context';
// import { UserInfoModal } from '@modules/auth';
// import { usePermissions } from '@modules/auth/ui/hooks';

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
  const { primaryExpanded, requestExpand, requestCollapse } = useSidebar();
  const [showUserModal, setShowUserModal] = useState(false);
  
  // Determinar si hay sidebar secundario basado en la ruta actual
  const iamRoutes = ['/users', '/roles', '/permissions', '/businesses', '/business-types', '/resources'];
  const ordersRoutes = ['/products', '/orders', '/shipments', '/order-status', '/notification-config'];
  const hasSecondarySidebar = iamRoutes.some(route => pathname.startsWith(route)) || 
                               ordersRoutes.some(route => pathname.startsWith(route));
  
  // const { hasResource, hasRouteAccess } = usePermissions();

  const handleLogout = () => {
    TokenStorage.clearSession();
    router.push('/login');
  };

  if (!user) return null;

  // Helper para determinar si un link está activo
  const isActive = (path: string) => pathname === path;

  // Verificar permisos para mostrar items
  // TODO: Migrar usePermissions a la nueva arquitectura

  const canAccessIAM = true; // hasResource('Usuarios') || hasResource('Roles') || hasResource('Permisos') || hasResource('Recursos') || hasResource('Tipos de Negocio') || hasResource('Negocios');

  return (
    <>
      {/* Sidebar - Menú lateral expandible */}
      <aside
        className="fixed left-0 top-0 h-full transition-all duration-300 z-30"
        style={{
          width: primaryExpanded ? '250px' : '80px',
          backgroundColor: 'var(--color-primary)'
        }}
        onMouseEnter={requestExpand}
        onMouseLeave={() => requestCollapse(hasSecondarySidebar)}
      >
        <div className="flex flex-col h-full">
          {/* Tarjeta de usuario arriba */}
          <div
            className="p-4 border-b border-white/10 cursor-pointer hover:bg-white/5 transition-colors"
            onClick={() => setShowUserModal(true)}
          >
            <div className="flex items-center gap-3">
              {/* Avatar */}
              {user.avatarUrl ? (
                <img
                  src={user.avatarUrl}
                  alt={user.name}
                  className="w-12 h-12 rounded-full object-cover flex-shrink-0 border-2 border-white/20"
                />
              ) : (
                <div
                  className="w-12 h-12 rounded-full flex items-center justify-center text-white text-lg font-bold flex-shrink-0"
                  style={{ backgroundColor: 'var(--color-secondary)' }}
                >
                  {user.name.charAt(0).toUpperCase()}
                </div>
              )}

              {/* Nombre (solo visible cuando está expandido) */}
              {primaryExpanded && (
                <div className="text-white overflow-hidden">
                  <p className="font-semibold text-sm truncate">{user.name}</p>
                  <p className="text-xs text-white/70 truncate">{user.email}</p>
                </div>
              )}
            </div>
          </div>

          {/* Menú de navegación */}
          <nav className="flex-1 py-6 px-3">
            <ul className="space-y-2">
              {/* Item Home */}
              <li>
                <Link
                  href="/home"
                  className={`
                    flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                    ${isActive('/home')
                      ? 'bg-white/20 text-white shadow-lg scale-105'
                      : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
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

              {/* Item Integraciones */}
              <li>
                <Link
                  href="/integrations"
                  className={`
                    flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                    ${isActive('/integrations') || pathname.startsWith('/integrations')
                      ? 'bg-white/20 text-white shadow-lg scale-105'
                      : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
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

              {/* Item Ordenes (Gestión de Ordenes) */}
              <li>
                <Link
                  href="/orders" // Default entry point for Orders
                  className={`
                    flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                    ${isActive('/orders') || isActive('/products') || isActive('/shipments') || isActive('/order-status') || isActive('/notification-config')
                      ? 'bg-white/20 text-white shadow-lg scale-105'
                      : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
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
                    <span className="text-sm font-medium transition-opacity duration-300">
                      Ordenes
                    </span>
                  )}
                </Link>
              </li>

              {/* Item IAM (Gestión de Identidad) */}
              {canAccessIAM && (
                <li>
                  <Link
                    href="/users" // Default entry point for IAM
                    className={`
                      flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                      ${isActive('/users') || isActive('/roles') || isActive('/permissions') || isActive('/businesses') || isActive('/business-types') || isActive('/resources')
                        ? 'bg-white/20 text-white shadow-lg scale-105'
                        : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
                      }
                    `}
                  >
                    {(isActive('/users') || isActive('/roles') || isActive('/permissions') || isActive('/businesses') || isActive('/business-types') || isActive('/resources')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                    </svg>
                    {primaryExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">
                        IAM
                      </span>
                    )}
                  </Link>
                </li>
              )}
            </ul>
          </nav>

          {/* Botón logout abajo */}
          <div className="p-4 border-t border-white/10">
            <button
              onClick={handleLogout}
              className="w-full flex items-center gap-3 text-white hover:bg-white/10 p-3 rounded-lg transition-colors"
            >
              <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              {primaryExpanded && <span className="text-sm">Cerrar Sesión</span>}
            </button>
          </div>
        </div>
      </aside>

      {/* Modal con información del usuario */}
      {/* <UserInfoModal
        isOpen={showUserModal}
        onClose={() => setShowUserModal(false)}
        onLogout={handleLogout}
        user={user}
      /> */}
    </>
  );
}

