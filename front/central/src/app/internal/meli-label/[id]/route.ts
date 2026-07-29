import { NextRequest } from 'next/server';
import { getAuthToken } from '@/shared/utils/server-auth';
import { env } from '@/shared/config/env';

export async function GET(request: NextRequest, { params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const url = new URL(request.url);
    const businessId = url.searchParams.get('business_id') || '';
    const download = url.searchParams.get('download') || '';
    const responseType = url.searchParams.get('response_type') || '';

    const token = await getAuthToken();
    if (!token) {
        return new Response('Unauthorized', { status: 401 });
    }

    const qs = new URLSearchParams();
    if (businessId) qs.set('business_id', businessId);
    if (download === '1') qs.set('download', '1');
    if (responseType) qs.set('response_type', responseType);
    const target = `${env.API_BASE_URL}/integrations/meli/shipments/${id}/label${qs.toString() ? '?' + qs.toString() : ''}`;

    const upstream = await fetch(target, {
        headers: { Authorization: `Bearer ${token}` },
    });

    if (!upstream.ok) {
        const text = await upstream.text().catch(() => '');
        return new Response(text || 'Error obteniendo la etiqueta de MercadoLibre', { status: upstream.status });
    }

    const headers = new Headers();
    headers.set('Content-Type', upstream.headers.get('Content-Type') || 'application/pdf');
    const disposition = upstream.headers.get('Content-Disposition');
    if (disposition) headers.set('Content-Disposition', disposition);

    return new Response(upstream.body, { status: 200, headers });
}
