import { cookies } from 'next/headers';
import { NextRequest } from 'next/server';

const API_URL = process.env.MONITORING_API_URL || 'http://localhost:3070';

export async function GET(
    _request: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
    const { id } = await params;
    const cookieStore = await cookies();
    const token = cookieStore.get('monitoring_token')?.value;

    if (!token) {
        return Response.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const res = await fetch(`${API_URL}/api/v1/containers/${id}/stats`, {
        headers: { Authorization: `Bearer ${token}` },
        cache: 'no-store',
    });

    if (!res.ok) {
        return Response.json({ error: res.statusText }, { status: res.status });
    }

    const data = await res.json();
    return Response.json(data);
}
