'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useSidebar } from '@/shared/contexts/sidebar-context';

export function OrdersSidebar() {
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
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
                        </svg>
                    </div>
                    {isExpanded && (
                        <h2 className="text-base font-bold text-gray-800 leading-tight whitespace-nowrap">
                            Gestión de<br />Ordenes
                        </h2>
                    )}
                </div>

                <div className="space-y-6">
                    {/* CATÁLOGO */}
                    <div>
                        {isExpanded && (
                            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">
                                CATÁLOGO
                            </h3>
                        )}
                        <ul className="space-y-0.5">
                            <li>
                                <Link 
                                    href="/products" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/products') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                                    </svg>
                                    {isExpanded && <span>Productos</span>}
                                </Link>
                            </li>
                        </ul>
                    </div>

                    {/* OPERACIONES */}
                    <div>
                        {isExpanded && (
                            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">
                                OPERACIONES
                            </h3>
                        )}
                        <ul className="space-y-0.5">
                            <li>
                                <Link 
                                    href="/orders" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/orders') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
                                    </svg>
                                    {isExpanded && <span>Ordenes</span>}
                                </Link>
                            </li>
                            <li>
                                <Link 
                                    href="/shipments" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/shipments') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8a1 1 0 011-1h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1m-6-1a1 1 0 001 1h1M5 17a2 2 0 104 0m-4 0a2 2 0 114 0m6 0a2 2 0 104 0m-4 0a2 2 0 114 0" />
                                    </svg>
                                    {isExpanded && <span>Envíos</span>}
                                </Link>
                            </li>
                        </ul>
                    </div>

                    {/* CONFIGURACIÓN */}
                    <div>
                        {isExpanded && (
                            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">
                                CONFIGURACIÓN
                            </h3>
                        )}
                        <ul className="space-y-0.5">
                            <li>
                                <Link 
                                    href="/order-status" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/order-status') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                    </svg>
                                    {isExpanded && <span>Estados de Orden</span>}
                                </Link>
                            </li>
                            <li>
                                <Link 
                                    href="/notification-config" 
                                    className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm font-medium transition-all ${
                                        isActive('/notification-config') 
                                            ? 'bg-blue-50 text-blue-700 border-l-2 border-blue-600' 
                                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                                    }`}
                                >
                                    <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                                    </svg>
                                    {isExpanded && <span>Notificaciones</span>}
                                </Link>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </aside>
    );
}
