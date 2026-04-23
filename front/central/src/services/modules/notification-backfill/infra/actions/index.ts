'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';
import { BackfillApiRepository } from '../repository/api-repository';
import type { PreviewRequest, RunRequest } from '../../domain/types';

const getRepository = async () => {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || '';
    return new BackfillApiRepository(env.API_BASE_URL, token);
};

export async function listBackfillEventsAction() {
    try {
        const repo = await getRepository();
        const events = await repo.listEvents();
        return { success: true, data: events };
    } catch (error: any) {
        return { success: false, error: error.message };
    }
}

export async function previewBackfillAction(req: PreviewRequest) {
    try {
        const repo = await getRepository();
        const data = await repo.preview(req);
        return { success: true, data };
    } catch (error: any) {
        return { success: false, error: error.message };
    }
}

export async function runBackfillAction(req: RunRequest) {
    try {
        const repo = await getRepository();
        const data = await repo.run(req);
        return { success: true, data };
    } catch (error: any) {
        return { success: false, error: error.message };
    }
}

export async function getBackfillJobAction(jobId: string) {
    try {
        const repo = await getRepository();
        const data = await repo.getJob(jobId);
        return { success: true, data };
    } catch (error: any) {
        return { success: false, error: error.message };
    }
}
