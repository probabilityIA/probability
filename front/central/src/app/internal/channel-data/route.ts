import { NextRequest } from 'next/server';
import { getAuthToken } from '@/shared/utils/server-auth';
import { env } from '@/shared/config/env';

const RUTAS: Record<string, { path: string; method: 'GET' | 'POST' }> = {
    summary: { path: '/integrations/sync-runs/data-summary', method: 'GET' },
    preview: { path: '/integrations/sync-runs/data-preview', method: 'POST' },
    apply: { path: '/products/channel-data/apply', method: 'POST' },
    undo: { path: '/products/channel-data/undo', method: 'POST' },
    batches: { path: '/products/channel-data/batches', method: 'GET' },
};

async function proxy(request: NextRequest, accion: string | null, body?: string) {
    const token = await getAuthToken();
    if (!token) {
        return Response.json({ success: false, error: 'Unauthorized' }, { status: 401 });
    }

    const ruta = accion ? RUTAS[accion] : undefined;
    if (!ruta) {
        return Response.json({ success: false, error: 'accion invalida' }, { status: 400 });
    }

    const businessId = new URL(request.url).searchParams.get('business_id');
    const qs = businessId ? `?business_id=${encodeURIComponent(businessId)}` : '';

    const upstream = await fetch(`${env.API_BASE_URL}${ruta.path}${qs}`, {
        method: ruta.method,
        headers: {
            Authorization: `Bearer ${token}`,
            ...(body ? { 'Content-Type': 'application/json' } : {}),
        },
        body,
        cache: 'no-store',
    });

    return new Response(await upstream.text(), {
        status: upstream.status,
        headers: { 'Content-Type': 'application/json', 'Cache-Control': 'no-store' },
    });
}

export async function GET(request: NextRequest) {
    return proxy(request, new URL(request.url).searchParams.get('accion'));
}

export async function POST(request: NextRequest) {
    return proxy(request, new URL(request.url).searchParams.get('accion'), await request.text());
}
