'use client';

import Link from 'next/link';
import { useState } from 'react';
import { NavbarContent, NavbarLink, PublicBusiness } from '../../domain/types';
import { TiendaAccountButton } from '../session/TiendaAccountButton';

interface PublicNavProps {
    business: PublicBusiness;
    navbarContent?: NavbarContent | null;
}

const DEFAULT_LINKS: NavbarLink[] = [
    { label: 'Inicio', href: '' },
    { label: 'Productos', href: '/productos' },
    { label: 'Contacto', href: '/contacto' },
];

const ALIGNMENT_CLASS: Record<string, string> = {
    left: 'justify-start',
    center: 'justify-center',
    right: 'justify-end',
};

function resolveHref(slug: string, href: string): string {
    if (!href) return `/tienda/${slug}`;
    if (href.startsWith('http')) return href;
    return `/tienda/${slug}${href}`;
}

function defaultIcon(href: string): string {
    if (!href) return '⌂';
    if (href.includes('producto')) return '▤';
    if (href.includes('contacto')) return '✉';
    return '•';
}

function NavIcon({ link }: { link: NavbarLink }) {
    if (link.icon_type === 'image' && link.icon_image) {
        return <img src={link.icon_image} alt="" className="w-5 h-5 object-contain shrink-0" />;
    }
    return <span className="text-lg leading-none shrink-0 grayscale">{link.icon || defaultIcon(link.href)}</span>;
}

export function PublicNav({ business, navbarContent }: PublicNavProps) {
    const [mobileOpen, setMobileOpen] = useState(false);
    const slug = business.code;
    const position = navbarContent?.position || 'top';
    const alignment = navbarContent?.alignment || 'right';
    const links = navbarContent?.links && navbarContent.links.length > 0 ? navbarContent.links : DEFAULT_LINKS;

    const logo = business.logo_url ? (
        <img src={business.logo_url} alt={business.name} className="h-10 object-contain" />
    ) : (
        <span className="text-xl font-bold" style={{ color: 'var(--brand-nav-text-strong, var(--brand-primary))' }}>
            {business.name}
        </span>
    );

    if (position === 'left' || position === 'right') {
        const collapsible = !!navbarContent?.collapsible;
        return (
            <nav
                className={`group shrink-0 h-full flex flex-col gap-6 p-4 border-black/5 overflow-hidden transition-all duration-200 ${position === 'left' ? 'border-r' : 'border-l'} ${collapsible ? 'w-16 hover:w-56' : 'w-56'}`}
                style={{ backgroundColor: 'var(--brand-nav-bg, #ffffff)' }}
            >
                <Link href={`/tienda/${slug}`} className="shrink-0 flex items-center justify-center group-hover:justify-start">
                    {logo}
                </Link>
                <div className="flex flex-col gap-1">
                    {links.map((link, i) => (
                        <Link
                            key={i}
                            href={resolveHref(slug, link.href)}
                            className={`flex items-center gap-3 text-sm font-medium opacity-90 hover:opacity-100 py-2 rounded-lg hover:bg-black/5 whitespace-nowrap ${collapsible ? 'justify-center px-0 group-hover:px-3' : 'px-3'}`}
                            style={{ color: 'var(--brand-nav-text, #4b5563)' }}
                        >
                            {collapsible && <NavIcon link={link} />}
                            <span className={collapsible ? 'opacity-0 max-w-0 group-hover:opacity-100 group-hover:max-w-xs overflow-hidden transition-all' : ''}>
                                {link.label}
                            </span>
                        </Link>
                    ))}
                </div>
                <div className={`mt-auto ${collapsible ? 'flex justify-center group-hover:justify-start' : ''}`}>
                    <TiendaAccountButton slug={slug} collapsible={collapsible} />
                </div>
            </nav>
        );
    }

    return (
        <nav className="sticky top-0 z-50 shadow-sm border-b border-black/5" style={{ backgroundColor: 'var(--brand-nav-bg, #ffffff)' }}>
            <div className="max-w-7xl mx-auto px-4">
                <div className="flex items-center justify-between h-16">
                    <Link href={`/tienda/${slug}`} className="flex items-center gap-3 shrink-0">
                        {logo}
                    </Link>

                    <div className={`hidden md:flex items-center gap-6 flex-1 ${ALIGNMENT_CLASS[alignment]}`}>
                        {links.map((link, i) => (
                            <Link
                                key={i}
                                href={resolveHref(slug, link.href)}
                                className="text-sm font-medium opacity-90 hover:opacity-100"
                                style={{ color: 'var(--brand-nav-text, #4b5563)' }}
                            >
                                {link.label}
                            </Link>
                        ))}
                    </div>

                    <div className="hidden md:block shrink-0">
                        <TiendaAccountButton slug={slug} />
                    </div>

                    <button onClick={() => setMobileOpen(!mobileOpen)} className="md:hidden p-2" style={{ color: 'var(--brand-nav-text, #4b5563)' }}>
                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            {mobileOpen ? (
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            ) : (
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                            )}
                        </svg>
                    </button>
                </div>

                {mobileOpen && (
                    <div className="md:hidden py-4 border-t border-black/5 space-y-2">
                        {links.map((link, i) => (
                            <Link
                                key={i}
                                href={resolveHref(slug, link.href)}
                                className="block px-4 py-2 rounded-lg opacity-90 hover:opacity-100"
                                style={{ color: 'var(--brand-nav-text, #4b5563)' }}
                                onClick={() => setMobileOpen(false)}
                            >
                                {link.label}
                            </Link>
                        ))}
                        <div className="px-4 py-2">
                            <TiendaAccountButton slug={slug} />
                        </div>
                    </div>
                )}
            </div>
        </nav>
    );
}

export function AnnouncementBar({ navbarContent }: { navbarContent?: NavbarContent | null }) {
    if (!navbarContent?.announcement_text) return null;
    const animation = navbarContent.announcement_animation || 'none';
    const style = { backgroundColor: navbarContent.announcement_color || '#111827', color: '#ffffff' };

    if (animation === 'marquee') {
        return (
            <div className="announcement-marquee-track text-sm font-medium py-2" style={style}>
                <span className="announcement-marquee-content">{navbarContent.announcement_text}</span>
            </div>
        );
    }

    const animationClass = animation === 'fade' ? 'announcement-fade' : animation === 'slide' ? 'announcement-slide' : '';
    return (
        <div className={`text-center text-sm font-medium py-2 px-4 ${animationClass}`} style={style}>
            {navbarContent.announcement_text}
        </div>
    );
}
