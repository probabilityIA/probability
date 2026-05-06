import { cookies } from 'next/headers';
import { NextRequest } from 'next/server';
import { env } from '@/shared/config/env';

export async function GET(req: NextRequest) {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;

    const url = new URL(req.url);
    const params = url.searchParams.toString();
    const target = `${env.API_BASE_URL}/geozones/display${params ? `?${params}` : ''}`;

    const headers: Record<string, string> = { Accept: 'application/json' };
    if (token) headers['Authorization'] = `Bearer ${token}`;

    const ifNone = req.headers.get('if-none-match');
    if (ifNone) headers['If-None-Match'] = ifNone;

    const upstream = await fetch(target, { headers, cache: 'no-store' });

    const respHeaders = new Headers();
    for (const h of ['etag', 'cache-control', 'content-type', 'content-encoding', 'vary']) {
        const v = upstream.headers.get(h);
        if (v) respHeaders.set(h, v);
    }
    if (!respHeaders.has('content-type')) respHeaders.set('content-type', 'application/json');

    return new Response(upstream.body, {
        status: upstream.status,
        headers: respHeaders,
    });
}
