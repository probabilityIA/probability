import { NextRequest } from 'next/server';
import { getAuthToken } from '@/shared/utils/server-auth';
import { env } from '@/shared/config/env';

export async function GET(request: NextRequest, { params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const businessId = new URL(request.url).searchParams.get('business_id') || '';

    const token = await getAuthToken();
    if (!token) {
        return new Response('Unauthorized', { status: 401 });
    }

    const qs = businessId ? `?business_id=${encodeURIComponent(businessId)}` : '';
    const upstream = await fetch(`${env.API_BASE_URL}/siigo/invoices/${id}/pdf${qs}`, {
        headers: { Authorization: `Bearer ${token}` },
        cache: 'no-store',
    });

    if (!upstream.ok) {
        const cuerpo = await upstream.text().catch(() => '');
        let mensaje = 'No se pudo obtener el PDF de la factura';
        try {
            const data = JSON.parse(cuerpo);
            mensaje = data?.error || data?.message || mensaje;
        } catch {
            if (cuerpo) mensaje = cuerpo;
        }
        return new Response(mensaje, { status: upstream.status });
    }

    const headers = new Headers();
    headers.set('Content-Type', upstream.headers.get('Content-Type') || 'application/pdf');
    const disposition = upstream.headers.get('Content-Disposition');
    if (disposition) headers.set('Content-Disposition', disposition);

    return new Response(upstream.body, { status: 200, headers });
}
