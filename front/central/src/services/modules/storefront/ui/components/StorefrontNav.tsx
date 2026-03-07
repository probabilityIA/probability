'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { TokenStorage } from '@/shared/config';
import { ShoppingBagIcon, ClipboardDocumentListIcon, ArrowRightOnRectangleIcon } from '@heroicons/react/24/outline';

export function StorefrontNav() {
    const pathname = usePathname();
    const router = useRouter();

    const handleLogout = () => {
        TokenStorage.clearSession();
        router.push('/storefront/login');
    };

    const links = [
        { href: '/storefront/catalogo', label: 'Catalogo', icon: ShoppingBagIcon },
        { href: '/storefront/pedidos', label: 'Mis Pedidos', icon: ClipboardDocumentListIcon },
    ];

    return (
        <nav className="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700 px-4 py-3">
            <div className="max-w-7xl mx-auto flex items-center justify-between">
                <Link href="/storefront/catalogo" className="text-xl font-bold text-indigo-600 dark:text-indigo-400">
                    Tienda
                </Link>

                <div className="flex items-center gap-6">
                    {links.map((link) => {
                        const Icon = link.icon;
                        const isActive = pathname.startsWith(link.href);
                        return (
                            <Link
                                key={link.href}
                                href={link.href}
                                className={`flex items-center gap-2 text-sm font-medium transition-colors ${
                                    isActive
                                        ? 'text-indigo-600 dark:text-indigo-400'
                                        : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
                                }`}
                            >
                                <Icon className="w-5 h-5" />
                                <span className="hidden sm:inline">{link.label}</span>
                            </Link>
                        );
                    })}

                    <button
                        onClick={handleLogout}
                        className="flex items-center gap-2 text-sm font-medium text-gray-500 hover:text-red-600 transition-colors"
                    >
                        <ArrowRightOnRectangleIcon className="w-5 h-5" />
                        <span className="hidden sm:inline">Salir</span>
                    </button>
                </div>
            </div>
        </nav>
    );
}
