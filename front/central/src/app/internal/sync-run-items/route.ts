import { NextRequest } from 'next/server';
import { getAuthToken } from '@/shared/utils/server-auth';
import { env } from '@/shared/config/env';

const PASSTHROUGH = ['integration_id', 'kind', 'group', 'q', 'page', 'page_size', 'business_id'];

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

    const upstream = await fetch(`${env.API_BASE_URL}/integrations/sync-runs/items?${qs.toString()}`, {
        headers: { Authorization: `Bearer ${token}` },
        cache: 'no-store',
    });

    const body = await upstream.text();
    return new Response(body, {
        status: upstream.status,
        headers: { 'Content-Type': 'application/json', 'Cache-Control': 'no-store' },
    });
}
