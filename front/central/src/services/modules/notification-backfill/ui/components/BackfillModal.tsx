'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
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
    BusinessGroup,
    JobState,
    PreviewResponse,
} from '../../domain/types';

interface BackfillModalProps {
    isOpen: boolean;
    onClose: () => void;
}

const ALL_BUSINESSES = 'all';

export function BackfillModal({ isOpen, onClose }: BackfillModalProps) {
    const [events, setEvents] = useState<BackfillEvent[]>([]);
    const [eventCode, setEventCode] = useState<string>('');
    const [days, setDays] = useState<number>(4);
    const [businessScope, setBusinessScope] = useState<string>(ALL_BUSINESSES);
    const [preview, setPreview] = useState<PreviewResponse | null>(null);
    const [loadingPreview, setLoadingPreview] = useState(false);
    const [loadingRun, setLoadingRun] = useState(false);
    const [error, setError] = useState<string>('');
    const [jobId, setJobId] = useState<string>('');
    const [job, setJob] = useState<JobState | null>(null);
    const [expandedBusinesses, setExpandedBusinesses] = useState<Set<number>>(new Set());

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
        setBusinessScope(ALL_BUSINESSES);
        previewBackfillAction({ event_code: eventCode, days })
            .then((r) => {
                if (r.success && r.data) setPreview(r.data);
                else setError(r.error || 'Error en preview');
            })
            .finally(() => setLoadingPreview(false));
    }, [eventCode, days]);

    const selectedEvent = useMemo(
        () => events.find((e) => e.event_code === eventCode),
        [events, eventCode],
    );
    const isGuideEvent = eventCode === 'guia_envio_generada';

    const filteredBusinesses: BusinessGroup[] = useMemo(() => {
        if (!preview) return [];
        if (businessScope === ALL_BUSINESSES) return preview.businesses;
        const id = Number(businessScope);
        return preview.businesses.filter((b) => b.business_id === id);
    }, [preview, businessScope]);

    const scopeTotal = useMemo(
        () => filteredBusinesses.reduce((acc, b) => acc + b.count, 0),
        [filteredBusinesses],
    );

    const handleSSEMessage = useCallback(
        (event: MessageEvent) => {
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
                // ignore
            }
        },
        [jobId],
    );

    useSSE({
        eventTypes: jobId ? ['backfill.progress'] : undefined,
        onMessage: handleSSEMessage,
    });

    const toggleBusiness = (id: number) => {
        setExpandedBusinesses((prev) => {
            const next = new Set(prev);
            if (next.has(id)) next.delete(id);
            else next.add(id);
            return next;
        });
    };

    const handleRun = async () => {
        if (!eventCode || scopeTotal === 0) return;
        const bizID = businessScope === ALL_BUSINESSES ? undefined : Number(businessScope);
        setLoadingRun(true);
        setError('');
        const res = await runBackfillAction({ event_code: eventCode, days, business_id: bizID });
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
        setExpandedBusinesses(new Set());
        setBusinessScope(ALL_BUSINESSES);
        onClose();
    };

    const sent = job?.sent ?? 0;
    const failed = job?.failed ?? 0;
    const total = job?.total_eligible ?? scopeTotal;
    const progress = total > 0 ? Math.min(100, Math.round(((sent + failed) / total) * 100)) : 0;
    const running = job?.status === 'running';
    const completed = job?.status === 'completed' || job?.status === 'failed';

    return (
        <Modal isOpen={isOpen} onClose={resetAndClose} title="Envío masivo a faltantes" size="2xl">
            <div className="space-y-5">
                <div className="rounded-md border border-blue-200 dark:border-blue-900 bg-blue-50 dark:bg-blue-950 p-3 text-xs text-blue-800 dark:text-blue-200">
                    <div className="font-semibold mb-1">Reglas de elegibilidad</div>
                    <ul className="list-disc list-inside space-y-0.5">
                        <li>Solo se consideran negocios con la integración de <b>WhatsApp activa</b>.</li>
                        <li>Solo se consideran negocios que tengan el evento seleccionado <b>habilitado</b> en sus configuraciones de notificación.</li>
                        <li>Se excluyen órdenes a las que ya se les envió este mensaje (evita duplicados).</li>
                    </ul>
                </div>

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
                                setExpandedBusinesses(new Set());
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

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Negocio destino
                    </label>
                    <select
                        value={businessScope}
                        onChange={(e) => setBusinessScope(e.target.value)}
                        disabled={running || !preview || (preview.businesses?.length ?? 0) === 0}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-sm"
                    >
                        <option value={ALL_BUSINESSES}>
                            Todos ({preview?.total_eligible ?? 0})
                        </option>
                        {preview?.businesses.map((b) => (
                            <option key={b.business_id} value={String(b.business_id)}>
                                #{b.business_id} {b.business_name || 'Sin nombre'} ({b.count})
                            </option>
                        ))}
                    </select>
                </div>

                <div className="rounded-md border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900">
                    {loadingPreview ? (
                        <div className="flex items-center gap-2 text-sm p-4">
                            <Spinner /> Calculando elegibles…
                        </div>
                    ) : !preview ? (
                        <div className="text-sm text-gray-500 p-4">Selecciona un evento</div>
                    ) : filteredBusinesses.length === 0 ? (
                        <div className="text-sm text-gray-500 p-4">Sin elegibles para esta selección</div>
                    ) : (
                        <div className="divide-y divide-gray-200 dark:divide-gray-700">
                            <div className="p-3 flex items-center justify-between text-sm">
                                <span className="text-gray-600 dark:text-gray-400">Elegibles en selección</span>
                                <span className="text-2xl font-semibold text-gray-900 dark:text-white">
                                    {scopeTotal}
                                </span>
                            </div>
                            {filteredBusinesses.map((b) => {
                                const isExpanded = expandedBusinesses.has(b.business_id);
                                return (
                                    <div key={b.business_id}>
                                        <button
                                            type="button"
                                            onClick={() => toggleBusiness(b.business_id)}
                                            className="w-full flex items-center justify-between px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-800 text-left"
                                        >
                                            <div className="flex items-center gap-2 text-sm">
                                                <span className="text-gray-400">{isExpanded ? '▾' : '▸'}</span>
                                                <span className="font-mono text-xs text-gray-500">#{b.business_id}</span>
                                                <span className="font-medium">{b.business_name || 'Sin nombre'}</span>
                                            </div>
                                            <span className="text-sm font-semibold text-gray-700 dark:text-gray-300">
                                                {b.count} {b.count === 1 ? 'orden' : 'órdenes'}
                                            </span>
                                        </button>
                                        {isExpanded && (
                                            <div className="px-3 pb-3 pt-1 space-y-1">
                                                {b.orders.map((o) => (
                                                    <div
                                                        key={o.order_id}
                                                        className="flex items-center gap-3 text-xs bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md px-3 py-2"
                                                    >
                                                        <span className="font-mono font-medium text-gray-900 dark:text-white">
                                                            {o.order_number}
                                                        </span>
                                                        {isGuideEvent && o.tracking_number && (
                                                            <span className="text-gray-600 dark:text-gray-400">
                                                                Guía: <span className="font-mono">{o.tracking_number}</span>
                                                            </span>
                                                        )}
                                                        {isGuideEvent && (o.carrier || o.carrier_logo_url) && (
                                                            <span className="flex items-center gap-1 text-gray-600 dark:text-gray-400 ml-auto">
                                                                {o.carrier_logo_url && (
                                                                    <img
                                                                        src={o.carrier_logo_url}
                                                                        alt={o.carrier || 'carrier'}
                                                                        className="h-4 w-auto object-contain"
                                                                    />
                                                                )}
                                                                {o.carrier && <span>{o.carrier}</span>}
                                                            </span>
                                                        )}
                                                        {!isGuideEvent && (
                                                            <span className="text-gray-500 ml-auto">{o.status}</span>
                                                        )}
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                );
                            })}
                        </div>
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
                        {job.error_message && <div className="text-sm text-red-600">{job.error_message}</div>}
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
                            disabled={loadingRun || loadingPreview || running || scopeTotal === 0}
                        >
                            {loadingRun
                                ? 'Iniciando…'
                                : running
                                ? 'En progreso…'
                                : `Enviar ${scopeTotal} ${scopeTotal === 1 ? 'mensaje' : 'mensajes'}`}
                        </Button>
                    )}
                </div>
            </div>
        </Modal>
    );
}
