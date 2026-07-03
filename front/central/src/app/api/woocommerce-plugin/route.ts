import { env } from '@/shared/config/env';

export const dynamic = 'force-dynamic';

export async function GET() {
    let res: Response;
    try {
        res = await fetch(`${env.API_BASE_URL}/woocommerce/plugin-download`, {
            cache: 'no-store',
        });
    } catch {
        return new Response('No se pudo descargar el plugin', { status: 502 });
    }

    if (!res.ok) {
        return new Response('No se pudo descargar el plugin', { status: 502 });
    }

    const buf = await res.arrayBuffer();

    return new Response(buf, {
        headers: {
            'Content-Type': 'application/zip',
            'Content-Disposition': 'attachment; filename="probability-shipping.zip"',
            'Content-Length': String(buf.byteLength),
            'Cache-Control': 'no-store',
        },
    });
}
