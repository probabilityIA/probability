import { NextRequest } from 'next/server';
import { getAuthToken } from '@/shared/utils/server-auth';
import { env } from '@/shared/config/env';

export async function GET(request: NextRequest, { params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const url = new URL(request.url);
    const format = url.searchParams.get('format') || '';
    const download = url.searchParams.get('download') || '';

    const token = await getAuthToken();
    if (!token) {
        return new Response('Unauthorized', { status: 401 });
    }

    const qs = new URLSearchParams();
    if (format) qs.set('format', format);
    if (download === '1') qs.set('download', '1');
    const target = `${env.API_BASE_URL}/shipments/${id}/guide${qs.toString() ? '?' + qs.toString() : ''}`;

    const upstream = await fetch(target, {
        headers: { Authorization: `Bearer ${token}` },
    });

    if (!upstream.ok) {
        const text = await upstream.text().catch(() => '');
        return new Response(text || 'Error generando guia', { status: upstream.status });
    }

    const headers = new Headers();
    headers.set('Content-Type', upstream.headers.get('Content-Type') || 'application/pdf');
    const disposition = upstream.headers.get('Content-Disposition');
    if (disposition) headers.set('Content-Disposition', disposition);

    return new Response(upstream.body, { status: 200, headers });
}
