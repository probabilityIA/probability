import { cookies } from 'next/headers';

function getApiUrl() {
    return process.env.MONITORING_API_URL || 'http://localhost:3070';
}

export async function GET() {
    const cookieStore = await cookies();
    const token = cookieStore.get('monitoring_token')?.value;

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
