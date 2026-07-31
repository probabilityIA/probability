import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function proxy(request: NextRequest) {
    const response = NextResponse.next();

    const shop = request.nextUrl.searchParams.get('shop') ||
                 request.headers.get('x-shopify-shop-domain');

    let frameAncestors: string;

    if (shop && shop.endsWith('.myshopify.com')) {
        frameAncestors = `frame-ancestors https://${shop} https://admin.shopify.com`;
    } else if (shop) {
        frameAncestors = "frame-ancestors 'none'";
    } else {
        frameAncestors = 'frame-ancestors https://admin.shopify.com';
    }

    const csp = [
        frameAncestors,
        "default-src 'self'",
        "script-src 'self' 'unsafe-eval' 'unsafe-inline' https://cdn.shopify.com https://*.bold.co",
        "style-src 'self' 'unsafe-inline'",
        "img-src 'self' data: https: https://*.bold.co",
        "font-src 'self' data:",
        "connect-src 'self' http://localhost:3050 https://*.probabilityia.com.co wss://*.probabilityia.com.co https://cdn.shopify.com https://*.bold.co",
        "frame-src 'self' https://admin.shopify.com https://*.bold.co https://www.openstreetmap.org",
    ].join('; ');

    response.headers.set('Content-Security-Policy', csp);
    response.headers.set('Strict-Transport-Security', 'max-age=63072000; includeSubDomains; preload');
    response.headers.set('X-Content-Type-Options', 'nosniff');
    response.headers.set('Referrer-Policy', 'strict-origin-when-cross-origin');

    if (shop) {
        response.headers.set('X-Frame-Options', 'ALLOWALL');

        response.headers.set('Access-Control-Allow-Origin', `https://${shop}`);
        response.headers.set('Access-Control-Allow-Credentials', 'true');
        response.headers.set('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
        response.headers.set('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-Requested-With');
    }

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

export const config = {
    matcher: [
        '/((?!api|_next/static|_next/image|favicon.ico|.*\\.png$|.*\\.jpg$|.*\\.svg$).*)',
    ],
};
