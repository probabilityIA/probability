'use client';

import { useCallback, useEffect, useState } from 'react';
import { Modal } from '@/shared/ui/modal';
import { Button } from '@/shared/ui/button';
import { Spinner } from '@/shared/ui/spinner';
import { useSSE } from '@/shared/hooks/use-sse';
import {
    getBackfillJobAction,
    listBackfillEventsAction,
    previewBackfillAction,
    runBackfillAction,
} from '../../infra/actions';
import type {
    BackfillEvent,
    BackfillProgressEvent,
    JobState,
    PreviewResponse,
} from '../../domain/types';

interface BackfillModalProps {
    isOpen: boolean;
    onClose: () => void;
}

export function BackfillModal({ isOpen, onClose }: BackfillModalProps) {
    const [events, setEvents] = useState<BackfillEvent[]>([]);
    const [eventCode, setEventCode] = useState<string>('');
    const [days, setDays] = useState<number>(4);
    const [preview, setPreview] = useState<PreviewResponse | null>(null);
    const [loadingPreview, setLoadingPreview] = useState(false);
    const [loadingRun, setLoadingRun] = useState(false);
    const [error, setError] = useState<string>('');
    const [jobId, setJobId] = useState<string>('');
    const [job, setJob] = useState<JobState | null>(null);

    useEffect(() => {
        if (!isOpen) return;
        setError('');
        listBackfillEventsAction().then((r) => {
            if (r.success && r.data) {
                setEvents(r.data);
                if (r.data.length > 0 && !eventCode) setEventCode(r.data[0].event_code);
            } else if (!r.success) {
                setError(r.error || 'Error cargando eventos');
            }
        });
    }, [isOpen]);

    useEffect(() => {
        if (!eventCode) return;
        setLoadingPreview(true);
        setError('');
        previewBackfillAction({ event_code: eventCode, days })
            .then((r) => {
                if (r.success && r.data) setPreview(r.data);
                else setError(r.error || 'Error en preview');
            })
            .finally(() => setLoadingPreview(false));
    }, [eventCode, days]);

    const handleSSEMessage = useCallback((event: MessageEvent) => {
        if (!jobId) return;
        try {
            const data = JSON.parse(event.data);
            if (data?.type !== 'backfill.progress') return;
            const payload: BackfillProgressEvent | undefined = data?.data;
            if (!payload || payload.job_id !== jobId) return;
            setJob((prev) => ({
                id: payload.job_id,
                event_code: payload.event_code,
                status: payload.status as JobState['status'],
                total_eligible: payload.total_eligible,
                sent: payload.sent,
                skipped: payload.skipped,
                failed: payload.failed,
                started_at: prev?.started_at ?? new Date().toISOString(),
                error_message: payload.error_message,
            }));
        } catch {
            // ignore malformed SSE payloads
        }
    }, [jobId]);

    useSSE({
        eventTypes: jobId ? ['backfill.progress'] : undefined,
        onMessage: handleSSEMessage,
    });

    const handleRun = async () => {
        if (!eventCode || !preview || preview.total_eligible === 0) return;
        setLoadingRun(true);
        setError('');
        const res = await runBackfillAction({ event_code: eventCode, days });
        setLoadingRun(false);
        if (!res.success || !res.data) {
            setError(res.error || 'No se pudo iniciar el masivo');
            return;
        }
        setJobId(res.data.job_id);
        const initial = await getBackfillJobAction(res.data.job_id);
        if (initial.success && initial.data) setJob(initial.data);
    };

    const resetAndClose = () => {
        setJobId('');
        setJob(null);
        setPreview(null);
        setError('');
        onClose();
    };

    const total = job?.total_eligible ?? preview?.total_eligible ?? 0;
    const sent = job?.sent ?? 0;
    const failed = job?.failed ?? 0;
    const progress = total > 0 ? Math.min(100, Math.round(((sent + failed) / total) * 100)) : 0;
    const running = job?.status === 'running';
    const completed = job?.status === 'completed' || job?.status === 'failed';

    return (
        <Modal isOpen={isOpen} onClose={resetAndClose} title="Envío masivo a faltantes" size="2xl">
            <div className="space-y-5">
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                            Evento
                        </label>
                        <select
                            value={eventCode}
                            onChange={(e) => {
                                setEventCode(e.target.value);
                                setJob(null);
                                setJobId('');
                            }}
                            disabled={running || events.length === 0}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-sm"
                        >
                            {events.length === 0 && <option value="">Cargando…</option>}
                            {events.map((ev) => (
                                <option key={ev.event_code} value={ev.event_code}>
                                    {ev.event_name} ({ev.channel})
                                </option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                            Días hacia atrás
                        </label>
                        <input
                            type="number"
                            min={1}
                            max={30}
                            value={days}
                            disabled={running}
                            onChange={(e) => setDays(Number(e.target.value) || 4)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-sm"
                        />
                    </div>
                </div>

                <div className="rounded-md border border-gray-200 dark:border-gray-700 p-4 bg-gray-50 dark:bg-gray-900">
                    {loadingPreview ? (
                        <div className="flex items-center gap-2 text-sm">
                            <Spinner /> Calculando elegibles…
                        </div>
                    ) : preview ? (
                        <div className="space-y-2 text-sm">
                            <div className="flex items-center justify-between">
                                <span className="text-gray-600 dark:text-gray-400">Elegibles</span>
                                <span className="text-2xl font-semibold text-gray-900 dark:text-white">
                                    {preview.total_eligible}
                                </span>
                            </div>
                            {Object.keys(preview.breakdown_by_business ?? {}).length > 0 && (
                                <div>
                                    <div className="text-xs text-gray-500 mb-1">Por negocio:</div>
                                    <div className="flex flex-wrap gap-2">
                                        {Object.entries(preview.breakdown_by_business).map(([bizId, n]) => (
                                            <span
                                                key={bizId}
                                                className="px-2 py-1 rounded-full text-xs bg-gray-200 dark:bg-gray-700"
                                            >
                                                Business {bizId}: {n}
                                            </span>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>
                    ) : (
                        <div className="text-sm text-gray-500">Selecciona un evento</div>
                    )}
                </div>

                {(running || completed) && job && (
                    <div className="space-y-2">
                        <div className="flex items-center justify-between text-sm">
                            <span className="text-gray-600 dark:text-gray-400">
                                {running ? 'Enviando…' : job.status === 'completed' ? 'Completado' : 'Fallido'}
                            </span>
                            <span className="font-medium">
                                {sent + failed} / {total} ({progress}%)
                            </span>
                        </div>
                        <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2 overflow-hidden">
                            <div
                                className={`h-2 transition-all duration-300 ${
                                    job.status === 'failed' ? 'bg-red-500' : 'bg-blue-600'
                                }`}
                                style={{ width: `${progress}%` }}
                            />
                        </div>
                        <div className="flex gap-4 text-xs text-gray-600 dark:text-gray-400">
                            <span>Enviados: {sent}</span>
                            <span>Fallidos: {failed}</span>
                            {job.skipped > 0 && <span>Omitidos: {job.skipped}</span>}
                        </div>
                        {job.error_message && (
                            <div className="text-sm text-red-600">{job.error_message}</div>
                        )}
                    </div>
                )}

                {error && <div className="text-sm text-red-600">{error}</div>}

                <div className="flex justify-end gap-3 pt-2">
                    <Button variant="secondary" onClick={resetAndClose} disabled={running}>
                        {completed ? 'Cerrar' : 'Cancelar'}
                    </Button>
                    {!completed && (
                        <Button
                            variant="primary"
                            onClick={handleRun}
                            disabled={
                                loadingRun ||
                                loadingPreview ||
                                running ||
                                !preview ||
                                preview.total_eligible === 0
                            }
                        >
                            {loadingRun ? 'Iniciando…' : running ? 'En progreso…' : 'Enviar masivo a faltantes'}
                        </Button>
                    )}
                </div>
            </div>
        </Modal>
    );
}
