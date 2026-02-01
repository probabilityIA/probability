import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
    const response = NextResponse.next();

    // Obtener el dominio de la tienda de Shopify desde query params o headers
    const shop = request.nextUrl.searchParams.get('shop') ||
                 request.headers.get('x-shopify-shop-domain');

    // Construir CSP dinámico basado en la tienda específica
    let frameAncestors: string;

    if (shop && shop.endsWith('.myshopify.com')) {
        // CRÍTICO: Frame-ancestors DEBE ser específico por tienda (NO wildcards)
        // Ref: https://shopify.dev/docs/apps/build/security/set-up-iframe-protection
        frameAncestors = `frame-ancestors https://${shop} https://admin.shopify.com`;
    } else if (shop) {
        // Shop inválido (no termina en .myshopify.com)
        frameAncestors = "frame-ancestors 'none'";
    } else {
        // Sin shop param, permitir solo admin.shopify.com
        frameAncestors = 'frame-ancestors https://admin.shopify.com';
    }

    // Construir CSP completo
    // IMPORTANTE: 'unsafe-inline' y 'unsafe-eval' son necesarios para Next.js en desarrollo
    // TODO: Ajustar para producción con nonces
    const csp = [
        frameAncestors,
        "default-src 'self'",
        "script-src 'self' 'unsafe-eval' 'unsafe-inline' https://cdn.shopify.com",
        "style-src 'self' 'unsafe-inline'",
        "img-src 'self' data: https:",
        "font-src 'self' data:",
        "connect-src 'self' https://*.probabilityia.com.co wss://*.probabilityia.com.co https://cdn.shopify.com",
        "frame-src 'self' https://admin.shopify.com",
    ].join('; ');

    // Security headers
    response.headers.set('Content-Security-Policy', csp);
    response.headers.set('Strict-Transport-Security', 'max-age=63072000; includeSubDomains; preload');
    response.headers.set('X-Content-Type-Options', 'nosniff');
    response.headers.set('Referrer-Policy', 'strict-origin-when-cross-origin');

    // Headers específicos para iframes de Shopify
    if (shop) {
        // Permitir que el iframe funcione en Shopify
        response.headers.set('X-Frame-Options', 'ALLOWALL'); // Deprecated pero algunos navegadores antiguos lo usan

        // Cross-Origin headers para permitir peticiones desde Shopify
        response.headers.set('Access-Control-Allow-Origin', `https://${shop}`);
        response.headers.set('Access-Control-Allow-Credentials', 'true');
        response.headers.set('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
        response.headers.set('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-Requested-With');
    }

    // CRÍTICO: Permitir cookies en contexto de terceros (necesario para iframes)
    // Esto hace que las cookies funcionen con SameSite=None; Secure
    response.cookies.set({
        name: '__session_test',
        value: 'test',
        sameSite: 'none',
        secure: true,
        httpOnly: false,
        path: '/',
    });

    return response;
}

// Aplicar el middleware a todas las rutas excepto static files y API routes
export const config = {
    matcher: [
        /*
         * Match all request paths except for the ones starting with:
         * - api (API routes)
         * - _next/static (static files)
         * - _next/image (image optimization files)
         * - favicon.ico (favicon file)
         * - public folder
         */
        '/((?!api|_next/static|_next/image|favicon.ico|.*\\.png$|.*\\.jpg$|.*\\.svg$).*)',
    ],
};
