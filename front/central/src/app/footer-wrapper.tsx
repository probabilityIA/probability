'use client';

import { usePathname } from 'next/navigation';
import { Footer } from "@/shared/ui";

export function FooterWrapper() {
    const pathname = usePathname();

    // No mostrar footer en la página de login
    if (pathname === '/login') {
        return null;
    }

    return <Footer />;
}
