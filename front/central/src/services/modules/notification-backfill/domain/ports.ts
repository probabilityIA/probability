import type { BackfillEvent, JobState, PreviewRequest, PreviewResponse, RunRequest, RunResponse } from './types';

export interface IBackfillRepository {
    listEvents(): Promise<BackfillEvent[]>;
    preview(req: PreviewRequest): Promise<PreviewResponse>;
    run(req: RunRequest): Promise<RunResponse>;
    getJob(jobId: string): Promise<JobState>;
}
