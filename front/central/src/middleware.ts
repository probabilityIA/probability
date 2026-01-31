import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
    const response = NextResponse.next();

    // Obtener el dominio de la tienda de Shopify desde query params
    const shop = request.nextUrl.searchParams.get('shop');

    // Construir CSP dinámico basado en la tienda específica
    let csp: string;

    if (shop && shop.endsWith('.myshopify.com')) {
        // CRÍTICO: Frame-ancestors DEBE ser específico por tienda (NO wildcards)
        // Ref: https://shopify.dev/docs/apps/build/security/set-up-iframe-protection
        csp = `frame-ancestors https://${shop} https://admin.shopify.com`;

        // X-Frame-Options para navegadores legacy
        response.headers.set('X-Frame-Options', `ALLOW-FROM https://${shop}`);
    } else if (shop) {
        // Shop inválido (no termina en .myshopify.com)
        csp = "frame-ancestors 'none'";
    } else {
        // Sin shop param, permitir solo admin.shopify.com
        csp = 'frame-ancestors https://admin.shopify.com';
    }

    // Security headers
    response.headers.set('Content-Security-Policy', csp);
    response.headers.set('Strict-Transport-Security', 'max-age=63072000; includeSubDomains; preload');
    response.headers.set('X-Content-Type-Options', 'nosniff');

    return response;
}

// Aplicar el middleware a todas las rutas
export const config = {
    matcher: [
        /*
         * Match all request paths except for the ones starting with:
         * - api (API routes)
         * - _next/static (static files)
         * - _next/image (image optimization files)
         * - favicon.ico (favicon file)
         */
        '/((?!api|_next/static|_next/image|favicon.ico).*)',
    ],
};
