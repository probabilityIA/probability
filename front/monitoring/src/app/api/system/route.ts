import { cookies } from 'next/headers';
import { NextRequest } from 'next/server';

function getApiUrl() {
    return process.env.MONITORING_API_URL || 'http://localhost:3070';
}

export async function GET(request: NextRequest) {
    let token = request.nextUrl.searchParams.get('token');
    if (!token) {
        const cookieStore = await cookies();
        token = cookieStore.get('monitoring_token')?.value || null;
    }

    if (!token) {
        return Response.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const res = await fetch(`${getApiUrl()}/api/v1/system/stats`, {
        headers: { Authorization: `Bearer ${token}` },
        cache: 'no-store',
    });

    if (!res.ok) {
        return Response.json({ error: res.statusText }, { status: res.status });
    }

    return Response.json(await res.json());
}
