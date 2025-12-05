/**
 * Sidebar de navegación
 * Componente compartido para todas las páginas autenticadas
 */

'use client';

import { useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import { TokenStorage } from '@/shared/config';
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
  const [sidebarExpanded, setSidebarExpanded] = useState(false);
  const [showUserModal, setShowUserModal] = useState(false);
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
          width: sidebarExpanded ? '250px' : '80px',
          backgroundColor: 'var(--color-primary)'
        }}
        onMouseEnter={() => setSidebarExpanded(true)}
        onMouseLeave={() => setSidebarExpanded(false)}
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
              {sidebarExpanded && (
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
                  {sidebarExpanded && (
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
                  {sidebarExpanded && (
                    <span className="text-sm font-medium transition-opacity duration-300">
                      Integraciones
                    </span>
                  )}
                </Link>
              </li>

              {/* Item Ordenes */}
              <li>
                <Link
                  href="/orders"
                  className={`
                    flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                    ${isActive('/orders') || pathname.startsWith('/orders')
                      ? 'bg-white/20 text-white shadow-lg scale-105'
                      : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
                    }
                  `}
                >
                  {/* Indicador activo (barra lateral) */}
                  {(isActive('/orders') || pathname.startsWith('/orders')) && (
                    <div
                      className="absolute left-0 w-1 h-8 rounded-r-full"
                      style={{ backgroundColor: 'var(--color-tertiary)' }}
                    />
                  )}

                  <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
                  </svg>
                  {sidebarExpanded && (
                    <span className="text-sm font-medium transition-opacity duration-300">
                      Ordenes
                    </span>
                  )}
                </Link>
              </li>

              {/* Item Order Status */}
              <li>
                <Link
                  href="/order-status"
                  className={`
                    flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                    ${isActive('/order-status') || pathname.startsWith('/order-status')
                      ? 'bg-white/20 text-white shadow-lg scale-105'
                      : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
                    }
                  `}
                >
                  {/* Indicador activo (barra lateral) */}
                  {(isActive('/order-status') || pathname.startsWith('/order-status')) && (
                    <div
                      className="absolute left-0 w-1 h-8 rounded-r-full"
                      style={{ backgroundColor: 'var(--color-tertiary)' }}
                    />
                  )}

                  <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  {sidebarExpanded && (
                    <span className="text-sm font-medium transition-opacity duration-300">
                      Order Status
                    </span>
                  )}
                </Link>
              </li>

              {/* Item Businesses */}
              {canAccessIAM && (
                <li>
                  <Link
                    href="/businesses"
                    className={`
                      flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                      ${isActive('/businesses') || pathname.startsWith('/businesses')
                        ? 'bg-white/20 text-white shadow-lg scale-105'
                        : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
                      }
                    `}
                  >
                    {(isActive('/businesses') || pathname.startsWith('/businesses')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                    </svg>
                    {sidebarExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">
                        Businesses
                      </span>
                    )}
                  </Link>
                </li>
              )}

              {/* Item Business Types */}
              {canAccessIAM && (
                <li>
                  <Link
                    href="/business-types"
                    className={`
                      flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                      ${isActive('/business-types') || pathname.startsWith('/business-types')
                        ? 'bg-white/20 text-white shadow-lg scale-105'
                        : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
                      }
                    `}
                  >
                    {(isActive('/business-types') || pathname.startsWith('/business-types')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                    {sidebarExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">
                        Business Types
                      </span>
                    )}
                  </Link>
                </li>
              )}

              {/* Item Users */}
              {canAccessIAM && (
                <li>
                  <Link
                    href="/users"
                    className={`
                      flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                      ${isActive('/users') || pathname.startsWith('/users')
                        ? 'bg-white/20 text-white shadow-lg scale-105'
                        : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
                      }
                    `}
                  >
                    {(isActive('/users') || pathname.startsWith('/users')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                    </svg>
                    {sidebarExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">
                        Users
                      </span>
                    )}
                  </Link>
                </li>
              )}

              {/* Item Roles */}
              {canAccessIAM && (
                <li>
                  <Link
                    href="/roles"
                    className={`
                      flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                      ${isActive('/roles') || pathname.startsWith('/roles')
                        ? 'bg-white/20 text-white shadow-lg scale-105'
                        : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
                      }
                    `}
                  >
                    {(isActive('/roles') || pathname.startsWith('/roles')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                    </svg>
                    {sidebarExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">
                        Roles
                      </span>
                    )}
                  </Link>
                </li>
              )}

              {/* Item Permissions */}
              {canAccessIAM && (
                <li>
                  <Link
                    href="/permissions"
                    className={`
                      flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                      ${isActive('/permissions') || pathname.startsWith('/permissions')
                        ? 'bg-white/20 text-white shadow-lg scale-105'
                        : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
                      }
                    `}
                  >
                    {(isActive('/permissions') || pathname.startsWith('/permissions')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                    </svg>
                    {sidebarExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">
                        Permissions
                      </span>
                    )}
                  </Link>
                </li>
              )}

              {/* Item Resources */}
              {canAccessIAM && (
                <li>
                  <Link
                    href="/resources"
                    className={`
                      flex items-center gap-3 p-3 rounded-lg transition-all duration-300
                      ${isActive('/resources') || pathname.startsWith('/resources')
                        ? 'bg-white/20 text-white shadow-lg scale-105'
                        : 'text-white/80 hover:bg-white/10 hover:text-white hover:scale-105'
                      }
                    `}
                  >
                    {(isActive('/resources') || pathname.startsWith('/resources')) && (
                      <div
                        className="absolute left-0 w-1 h-8 rounded-r-full"
                        style={{ backgroundColor: 'var(--color-tertiary)' }}
                      />
                    )}
                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
                    </svg>
                    {sidebarExpanded && (
                      <span className="text-sm font-medium transition-opacity duration-300">
                        Resources
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
              {sidebarExpanded && <span className="text-sm">Cerrar Sesión</span>}
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

