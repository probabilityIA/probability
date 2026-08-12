import { NextRequest } from 'next/server';
import { getAuthToken } from '@/shared/utils/server-auth';
import { env } from '@/shared/config/env';

export async function GET(request: NextRequest) {
    const token = await getAuthToken();
    if (!token) {
        return Response.json({ success: false, error: 'Unauthorized' }, { status: 401 });
    }

    const qs = new URLSearchParams();
    const businessId = new URL(request.url).searchParams.get('business_id');
    if (businessId) qs.set('business_id', businessId);

    const upstream = await fetch(`${env.API_BASE_URL}/integrations/sync-runs/findings?${qs.toString()}`, {
        headers: { Authorization: `Bearer ${token}` },
        cache: 'no-store',
    });

    const body = await upstream.text();
    return new Response(body, {
        status: upstream.status,
        headers: { 'Content-Type': 'application/json', 'Cache-Control': 'no-store' },
    });
}
