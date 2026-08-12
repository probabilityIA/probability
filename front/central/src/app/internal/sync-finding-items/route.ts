import { NextRequest } from 'next/server';
import { getAuthToken } from '@/shared/utils/server-auth';
import { env } from '@/shared/config/env';

const PASSTHROUGH = ['code', 'q', 'page', 'page_size', 'business_id', 'format'];

export async function GET(request: NextRequest) {
    const token = await getAuthToken();
    if (!token) {
        return Response.json({ success: false, error: 'Unauthorized' }, { status: 401 });
    }

    const incoming = new URL(request.url).searchParams;
    const qs = new URLSearchParams();
    for (const key of PASSTHROUGH) {
        const value = incoming.get(key);
        if (value) qs.set(key, value);
    }

    const upstream = await fetch(`${env.API_BASE_URL}/integrations/sync-runs/findings/items?${qs.toString()}`, {
        headers: { Authorization: `Bearer ${token}` },
        cache: 'no-store',
    });

    const headers = new Headers({ 'Cache-Control': 'no-store' });
    const contentType = upstream.headers.get('Content-Type');
    const disposition = upstream.headers.get('Content-Disposition');
    headers.set('Content-Type', contentType ?? 'application/json');
    if (disposition) headers.set('Content-Disposition', disposition);

    return new Response(await upstream.arrayBuffer(), { status: upstream.status, headers });
}
