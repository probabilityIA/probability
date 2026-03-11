'use client';

import Link from 'next/link';
import { useState } from 'react';
import { PublicBusiness } from '../../domain/types';

interface PublicNavProps {
    business: PublicBusiness;
}

export function PublicNav({ business }: PublicNavProps) {
    const [mobileOpen, setMobileOpen] = useState(false);
    const slug = business.code;

    return (
        <nav className="sticky top-0 z-50 bg-white shadow-sm border-b border-gray-100">
            <div className="max-w-7xl mx-auto px-4">
                <div className="flex items-center justify-between h-16">
                    <Link href={`/tienda/${slug}`} className="flex items-center gap-3">
                        {business.logo_url ? (
                            <img src={business.logo_url} alt={business.name} className="h-10 object-contain" />
                        ) : (
                            <span className="text-xl font-bold" style={{ color: 'var(--brand-primary)' }}>
                                {business.name}
                            </span>
                        )}
                    </Link>

                    {/* Desktop nav */}
                    <div className="hidden md:flex items-center gap-6">
                        <Link href={`/tienda/${slug}`} className="text-gray-600 hover:text-gray-900 text-sm font-medium">
                            Inicio
                        </Link>
                        <Link href={`/tienda/${slug}/productos`} className="text-gray-600 hover:text-gray-900 text-sm font-medium">
                            Productos
                        </Link>
                        <Link href={`/tienda/${slug}/contacto`} className="text-gray-600 hover:text-gray-900 text-sm font-medium">
                            Contacto
                        </Link>
                        <Link
                            href={`/login?redirect=/storefront/catalogo&business_code=${slug}`}
                            className="px-4 py-2 rounded-lg text-white text-sm font-medium transition-colors"
                            style={{ backgroundColor: 'var(--brand-secondary)' }}
                        >
                            Hacer Pedido
                        </Link>
                    </div>

                    {/* Mobile burger */}
                    <button onClick={() => setMobileOpen(!mobileOpen)} className="md:hidden p-2">
                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            {mobileOpen ? (
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            ) : (
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                            )}
                        </svg>
                    </button>
                </div>

                {/* Mobile menu */}
                {mobileOpen && (
                    <div className="md:hidden py-4 border-t border-gray-100 space-y-2">
                        <Link href={`/tienda/${slug}`} className="block px-4 py-2 text-gray-600 hover:bg-gray-50 rounded-lg" onClick={() => setMobileOpen(false)}>
                            Inicio
                        </Link>
                        <Link href={`/tienda/${slug}/productos`} className="block px-4 py-2 text-gray-600 hover:bg-gray-50 rounded-lg" onClick={() => setMobileOpen(false)}>
                            Productos
                        </Link>
                        <Link href={`/tienda/${slug}/contacto`} className="block px-4 py-2 text-gray-600 hover:bg-gray-50 rounded-lg" onClick={() => setMobileOpen(false)}>
                            Contacto
                        </Link>
                        <Link
                            href={`/login?redirect=/storefront/catalogo&business_code=${slug}`}
                            className="block px-4 py-2 rounded-lg text-white text-center font-medium"
                            style={{ backgroundColor: 'var(--brand-secondary)' }}
                            onClick={() => setMobileOpen(false)}
                        >
                            Hacer Pedido
                        </Link>
                    </div>
                )}
            </div>
        </nav>
    );
}
