import { env } from '@/shared/config/env';

export async function GET() {
    const res = await fetch(`${env.API_BASE_URL}/woocommerce/plugin-download`, {
        cache: 'no-store',
    });

    if (!res.ok) {
        return new Response('No se pudo descargar el plugin', { status: 502 });
    }

    const buf = await res.arrayBuffer();

    return new Response(buf, {
        headers: {
            'Content-Type': 'application/zip',
            'Content-Disposition': 'attachment; filename="probability-shipping.zip"',
        },
    });
}
