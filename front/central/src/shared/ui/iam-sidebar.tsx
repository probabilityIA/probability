'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useSidebar } from '@/shared/contexts/sidebar-context';

export function IAMSidebar() {
    const pathname = usePathname();
    const { 
        primaryExpanded, 
        secondaryExpanded,
        requestSecondaryExpand,
        requestSecondaryCollapse
    } = useSidebar();
    const isActive = (path: string) => pathname === path || pathname.startsWith(path);
    
    // Calcular la posición izquierda basada en el estado del sidebar primario
    const leftPosition = primaryExpanded ? '250px' : '80px';
    
    // El sidebar secundario puede expandirse independientemente
    const isExpanded = secondaryExpanded || primaryExpanded;
    const width = isExpanded ? '240px' : '60px';

    const handleMouseEnter = () => {
        // Solo expandir el secundario, NO tocar el principal
        requestSecondaryExpand();
    };

    const handleMouseLeave = () => {
        // Solo colapsar el secundario
        requestSecondaryCollapse();
    };

    return (
        <aside
            className="fixed top-0 h-full bg-white border-r border-gray-200 z-20 overflow-y-auto transition-all duration-300 shadow-sm"
            style={{ 
                left: leftPosition,
                width: width
            }}
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
        >
            <div className="p-4">
                <div className="flex items-center gap-3 mb-6">
                    <div className="p-2 bg-gray-50 rounded-lg flex-shrink-0 border border-gray-200">
                        <svg className="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                        </svg>
                    </div>
                    {isExpanded && (
                        <h2 className="text-base font-bold text-gray-800 leading-tight whitespace-nowrap">
                            Gestión de<br />Identidad
                        </h2>
                    )}
                </div>

                <div className="space-y-6">
                    {/* ORGANIZACIÓN */}
                    <div>
                        {isExpanded && (
                            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">
                                ORGANIZACIÓN
                            </h3>
                        )}
                        <ul className="space-y-0.5">
                            <li>
                                <Link 
                                    href="/businesses" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/businesses') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                                    </svg>
                                    {isExpanded && <span>Empresas</span>}
                                </Link>
                            </li>
                            <li>
                                <Link 
                                    href="/business-types" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/business-types') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                                    </svg>
                                    {isExpanded && <span>Tipos de Empresa</span>}
                                </Link>
                            </li>
                        </ul>
                    </div>

                    {/* CONTROL DE ACCESO */}
                    <div>
                        {isExpanded && (
                            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">
                                CONTROL DE ACCESO
                            </h3>
                        )}
                        <ul className="space-y-0.5">
                            <li>
                                <Link 
                                    href="/users" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/users') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                                    </svg>
                                    {isExpanded && <span>Usuarios</span>}
                                </Link>
                            </li>
                            <li>
                                <Link 
                                    href="/roles" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/roles') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                                    </svg>
                                    {isExpanded && <span>Roles</span>}
                                </Link>
                            </li>
                            <li>
                                <Link 
                                    href="/permissions" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/permissions') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                                    </svg>
                                    {isExpanded && <span>Permisos</span>}
                                </Link>
                            </li>
                        </ul>
                    </div>

                    {/* SISTEMA */}
                    <div>
                        {isExpanded && (
                            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">
                                SISTEMA
                            </h3>
                        )}
                        <ul className="space-y-0.5">
                            <li>
                                <Link 
                                    href="/resources" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/resources') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
                                    </svg>
                                    {isExpanded && <span>Recursos</span>}
                                </Link>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </aside>
    );
}
