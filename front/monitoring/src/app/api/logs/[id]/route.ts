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
    const cookieStore = await cookies();
    const token = cookieStore.get('monitoring_token')?.value;

    if (!token) {
        return new Response('Unauthorized', { status: 401 });
    }

    const upstream = await fetch(`${getApiUrl()}/api/v1/containers/${id}/logs/stream`, {
        headers: { Authorization: `Bearer ${token}` },
        signal: request.signal,
    });

    if (!upstream.ok) {
        return new Response(upstream.statusText, { status: upstream.status });
    }

    return new Response(upstream.body, {
        headers: {
            'Content-Type': 'text/event-stream',
            'Cache-Control': 'no-cache',
            'Connection': 'keep-alive',
        },
    });
}
