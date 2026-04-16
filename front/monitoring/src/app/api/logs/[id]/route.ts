import { cookies } from 'next/headers';
import { NextRequest } from 'next/server';

function getApiUrl() {
    return process.env.MONITORING_API_URL || 'http://localhost:3070';
}

export async function GET(
    request: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
    const { id } = await params;

    // Token from query param (preferred) or cookie
    let token = request.nextUrl.searchParams.get('token');
    if (!token) {
        const cookieStore = await cookies();
        token = cookieStore.get('monitoring_token')?.value || null;
    }

    if (!token) {
        return new Response('Unauthorized', { status: 401 });
    }

    const apiUrl = getApiUrl();

    try {
        const upstream = await fetch(`${apiUrl}/api/v1/containers/${id}/logs/stream`, {
            headers: { Authorization: `Bearer ${token}` },
            signal: request.signal,
        });

        if (!upstream.ok) {
            const body = await upstream.text();
            return new Response(body || upstream.statusText, { status: upstream.status });
        }

        return new Response(upstream.body, {
            headers: {
                'Content-Type': 'text/event-stream',
                'Cache-Control': 'no-cache',
                'Connection': 'keep-alive',
                'X-Accel-Buffering': 'no',
            },
        });
    } catch (err) {
        return new Response(`Stream error: ${err}`, { status: 502 });
    }
}
