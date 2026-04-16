'use server';

import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { MonitoringApiRepository } from '../repository/api-repository';

export async function loginAction(email: string, password: string): Promise<{ error?: string }> {
    const repo = new MonitoringApiRepository();

    try {
        const result = await repo.login({ email, password });

        const cookieStore = await cookies();
        cookieStore.set('monitoring_token', result.token, {
            httpOnly: false,
            secure: false,
            sameSite: 'lax',
            path: '/',
            maxAge: 60 * 60 * 24, // 24 hours
        });
        cookieStore.set('monitoring_user', JSON.stringify({ name: result.name, email: result.email }), {
            httpOnly: false,
            secure: false,
            sameSite: 'lax',
            path: '/',
            maxAge: 60 * 60 * 24,
        });
    } catch (err) {
        return { error: err instanceof Error ? err.message : 'Login failed' };
    }

    redirect('/dashboard');
}

export async function logoutAction() {
    const cookieStore = await cookies();
    cookieStore.delete('monitoring_token');
    cookieStore.delete('monitoring_user');
    redirect('/login');
}

export async function getTokenAction(): Promise<string | null> {
    const cookieStore = await cookies();
    return cookieStore.get('monitoring_token')?.value ?? null;
}

export async function restartContainerAction(id: string): Promise<{ message?: string; error?: string }> {
    const cookieStore = await cookies();
    const token = cookieStore.get('monitoring_token')?.value;
    if (!token) return { error: 'Not authenticated' };

    try {
        const repo = new MonitoringApiRepository(token);
        const result = await repo.restartContainer(id);
        return { message: result.message };
    } catch (err) {
        return { error: err instanceof Error ? err.message : 'Failed to restart' };
    }
}

export async function stopContainerAction(id: string): Promise<{ message?: string; error?: string }> {
    const cookieStore = await cookies();
    const token = cookieStore.get('monitoring_token')?.value;
    if (!token) return { error: 'Not authenticated' };

    try {
        const repo = new MonitoringApiRepository(token);
        const result = await repo.stopContainer(id);
        return { message: result.message };
    } catch (err) {
        return { error: err instanceof Error ? err.message : 'Failed to stop' };
    }
}

export async function startContainerAction(id: string): Promise<{ message?: string; error?: string }> {
    const cookieStore = await cookies();
    const token = cookieStore.get('monitoring_token')?.value;
    if (!token) return { error: 'Not authenticated' };

    try {
        const repo = new MonitoringApiRepository(token);
        const result = await repo.startContainer(id);
        return { message: result.message };
    } catch (err) {
        return { error: err instanceof Error ? err.message : 'Failed to start' };
    }
}
