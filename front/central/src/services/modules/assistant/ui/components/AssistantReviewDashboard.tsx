'use client';

import { Fragment, useState } from 'react';
import { ArrowPathIcon, HandThumbDownIcon, HandThumbUpIcon } from '@heroicons/react/24/solid';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { AssistantAvatar } from './AssistantAvatar';
import { useAssistantReview } from '../hooks/use-assistant-review';
import { ASSISTANT_NAME, percent } from '../../app/use-cases';
import type { ReviewKind, ReviewMessage } from '../../domain/types';

const PERIODS = [
    { days: 7, label: '\u00daltimos 7 d\u00edas' },
    { days: 30, label: '\u00daltimos 30 d\u00edas' },
    { days: 90, label: '\u00daltimos 90 d\u00edas' },
    { days: 365, label: '\u00daltimo a\u00f1o' },
];

const KINDS: { key: ReviewKind; label: string }[] = [
    { key: 'all', label: 'Todas' },
    { key: 'no_destination', label: 'Sin destino' },
    { key: 'negative', label: 'No sirvi\u00f3' },
    { key: 'positive', label: 'Sirvi\u00f3' },
    { key: 'not_clicked', label: 'No sigui\u00f3 el destino' },
    { key: 'errors', label: 'Con error' },
];

const ERROR_LABELS: Record<string, string> = {
    rate_limited: 'L\u00edmite',
    assistant_unavailable: 'Sin respuesta',
};

const numberFormat = new Intl.NumberFormat('es-CO');
const dateFormat = new Intl.DateTimeFormat('es-CO', { day: 'numeric', month: 'short', hour: 'numeric', minute: '2-digit' });

function usd(value: number): string {
    return `USD ${value.toLocaleString('es-CO', { minimumFractionDigits: 2, maximumFractionDigits: 4 })}`;
}

function Kpi({ label, value, detail }: { label: string; value: string; detail?: string }) {
    return (
        <div className="rounded-xl border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
            <p className="text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-gray-400">{label}</p>
            <p className="mt-1 text-2xl font-bold tabular-nums text-gray-900 dark:text-white">{value}</p>
            {detail && <p className="mt-0.5 text-xs text-gray-500 dark:text-gray-400">{detail}</p>}
        </div>
    );
}

function Outcome({ message }: { message: ReviewMessage }) {
    if (message.error_code) {
        return (
            <span className="inline-flex rounded-full bg-red-50 px-2 py-0.5 text-xs font-medium text-red-700 dark:bg-red-500/10 dark:text-red-300">
                {ERROR_LABELS[message.error_code] || message.error_code}
            </span>
        );
    }
    if (!message.destination_key) {
        return (
            <span className="inline-flex rounded-full bg-amber-50 px-2 py-0.5 text-xs font-medium text-amber-800 dark:bg-amber-500/10 dark:text-amber-300">
                Sin destino
            </span>
        );
    }
    return (
        <span className="inline-flex flex-col">
            <code className="font-mono text-xs text-[#3E0FA8] dark:text-violet-200">{message.destination_key}</code>
            <span className="text-[11px] text-gray-500 dark:text-gray-400">{message.clicked_at ? 'Hizo clic' : 'Sin clic'}</span>
        </span>
    );
}

function Feedback({ value }: { value: number }) {
    if (value === 1) return <HandThumbUpIcon className="h-4 w-4 text-emerald-600 dark:text-emerald-400" aria-label={'Sirvi\u00f3'} />;
    if (value === -1) return <HandThumbDownIcon className="h-4 w-4 text-red-600 dark:text-red-400" aria-label={'No sirvi\u00f3'} />;
    return <span className="text-gray-300 dark:text-gray-600">\u2014</span>;
}

export function AssistantReviewDashboard() {
    const review = useAssistantReview();
    const { businesses } = useBusinessesSimple();
    const [expanded, setExpanded] = useState<string | null>(null);

    const summary = review.summary;
    const list = review.messages;
    const answered = summary ? summary.messages - summary.errors : 0;

    return (
        <div className="space-y-6">
            <div className="flex flex-wrap items-center gap-4">
                <AssistantAvatar size={44} />
                <div className="min-w-0 flex-1">
                    <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{`Conversaciones con ${ASSISTANT_NAME}`}</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        {'Qu\u00e9 le pregunta la gente al asistente y d\u00f3nde falla. Se guardan durante un a\u00f1o.'}
                    </p>
                </div>
                <button
                    type="button"
                    onClick={() => void review.reload()}
                    className="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-200 dark:hover:bg-gray-700"
                >
                    <ArrowPathIcon className={`h-4 w-4 ${review.loading ? 'animate-spin' : ''}`} />
                    Actualizar
                </button>
            </div>

            <div className="flex flex-wrap gap-3">
                <select
                    id="assistant-review-period"
                    value={review.days}
                    onChange={(e) => review.changeDays(Number(e.target.value))}
                    aria-label={'Per\u00edodo'}
                    className="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                >
                    {PERIODS.map((p) => (
                        <option key={p.days} value={p.days}>{p.label}</option>
                    ))}
                </select>
                <select
                    id="assistant-review-business"
                    value={review.businessId ?? ''}
                    onChange={(e) => review.changeBusiness(e.target.value ? Number(e.target.value) : undefined)}
                    aria-label="Negocio"
                    className="min-w-[200px] rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                >
                    <option value="">Todos los negocios</option>
                    {businesses.map((b) => (
                        <option key={b.id} value={b.id}>{b.name}</option>
                    ))}
                </select>
                <input
                    id="assistant-review-search"
                    type="search"
                    value={review.searchInput}
                    onChange={(e) => review.setSearchInput(e.target.value)}
                    placeholder="Buscar en preguntas y respuestas"
                    className="min-w-[240px] flex-1 rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                />
            </div>

            {review.error && (
                <div className="rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700 dark:bg-red-500/10 dark:text-red-300">{review.error}</div>
            )}

            <div className="grid grid-cols-2 gap-3 lg:grid-cols-3 xl:grid-cols-6">
                <Kpi
                    label="Mensajes"
                    value={numberFormat.format(summary?.messages ?? 0)}
                    detail={`${numberFormat.format(summary?.conversations ?? 0)} conversaciones \u00b7 ${numberFormat.format(summary?.users ?? 0)} usuarios`}
                />
                <Kpi
                    label="Sin destino"
                    value={numberFormat.format(summary?.no_destination ?? 0)}
                    detail={`${percent(summary?.no_destination ?? 0, answered)} de las respuestas`}
                />
                <Kpi
                    label="Siguieron el destino"
                    value={percent(summary?.clicked ?? 0, summary?.with_destination ?? 0)}
                    detail={`${numberFormat.format(summary?.clicked ?? 0)} de ${numberFormat.format(summary?.with_destination ?? 0)} con destino`}
                />
                <Kpi
                    label={'Calificaci\u00f3n'}
                    value={`${numberFormat.format(summary?.positive ?? 0)} / ${numberFormat.format(summary?.negative ?? 0)}`}
                    detail={'Sirvi\u00f3 / no sirvi\u00f3'}
                />
                <Kpi
                    label="Errores"
                    value={numberFormat.format(summary?.errors ?? 0)}
                    detail={'L\u00edmite o modelo sin respuesta'}
                />
                <Kpi
                    label="Costo estimado"
                    value={usd(summary?.estimated_cost_usd ?? 0)}
                    detail={`${numberFormat.format(summary?.input_tokens ?? 0)} tokens entrada \u00b7 ${numberFormat.format(summary?.output_tokens ?? 0)} salida`}
                />
            </div>

            {summary && summary.top_destinations.length > 0 && (
                <div className="flex flex-wrap items-center gap-2">
                    <span className="text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-gray-400">{'Destinos m\u00e1s propuestos'}</span>
                    {summary.top_destinations.map((d) => (
                        <span key={d.key} className="inline-flex items-center gap-1.5 rounded-full border border-gray-200 bg-white px-2.5 py-1 text-xs dark:border-gray-700 dark:bg-gray-800">
                            <code className="font-mono text-[#3E0FA8] dark:text-violet-200">{d.key}</code>
                            <span className="tabular-nums text-gray-500 dark:text-gray-400">{numberFormat.format(d.count)}</span>
                        </span>
                    ))}
                </div>
            )}

            <div className="flex flex-wrap gap-1 border-b border-gray-200 dark:border-gray-700">
                {KINDS.map((k) => (
                    <button
                        key={k.key}
                        type="button"
                        onClick={() => review.changeKind(k.key)}
                        className={`-mb-px border-b-2 px-3 py-2 text-sm font-medium transition-colors ${
                            review.kind === k.key
                                ? 'border-[#5B1BE6] text-[#3E0FA8] dark:border-[#8B63FF] dark:text-violet-200'
                                : 'border-transparent text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200'
                        }`}
                    >
                        {k.label}
                    </button>
                ))}
            </div>

            <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
                        <thead className="bg-gray-50 text-left text-xs font-medium uppercase tracking-wide text-gray-500 dark:bg-gray-900/40 dark:text-gray-400">
                            <tr>
                                <th className="px-4 py-3">Fecha</th>
                                <th className="px-4 py-3">{'Negocio \u00b7 usuario'}</th>
                                <th className="px-4 py-3">Pregunta</th>
                                <th className="px-4 py-3">Resultado</th>
                                <th className="px-4 py-3 text-center">{'Calif.'}</th>
                                <th className="px-4 py-3 text-right">Tokens</th>
                                <th className="px-4 py-3 text-right">Tiempo</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                            {review.loading && !list ? (
                                <tr><td colSpan={7} className="px-4 py-8 text-center text-gray-500">Cargando conversaciones...</td></tr>
                            ) : !list || list.data.length === 0 ? (
                                <tr><td colSpan={7} className="px-4 py-8 text-center text-gray-500">{'No hay mensajes con estos filtros.'}</td></tr>
                            ) : (
                                list.data.map((m) => (
                                    <Fragment key={m.id}>
                                        <tr
                                            onClick={() => setExpanded((prev) => (prev === m.id ? null : m.id))}
                                            className="cursor-pointer align-top hover:bg-gray-50 dark:hover:bg-gray-700/40"
                                        >
                                            <td className="whitespace-nowrap px-4 py-3 tabular-nums text-gray-500 dark:text-gray-400">{dateFormat.format(new Date(m.created_at))}</td>
                                            <td className="px-4 py-3">
                                                <p className="font-medium text-gray-900 dark:text-white">{m.business_name || 'Sin negocio'}</p>
                                                <p className="text-xs text-gray-500 dark:text-gray-400">{m.user_name || m.user_email}</p>
                                            </td>
                                            <td className="max-w-md px-4 py-3">
                                                <p className="line-clamp-2 text-gray-900 dark:text-gray-100">{m.question}</p>
                                                {m.pathname && <p className="font-mono text-[11px] text-gray-400">{m.pathname}</p>}
                                            </td>
                                            <td className="px-4 py-3"><Outcome message={m} /></td>
                                            <td className="px-4 py-3"><div className="flex justify-center"><Feedback value={m.feedback} /></div></td>
                                            <td className="whitespace-nowrap px-4 py-3 text-right tabular-nums text-gray-500 dark:text-gray-400">{`${numberFormat.format(m.input_tokens)} / ${numberFormat.format(m.output_tokens)}`}</td>
                                            <td className="whitespace-nowrap px-4 py-3 text-right tabular-nums text-gray-500 dark:text-gray-400">{`${(m.latency_ms / 1000).toFixed(1)} s`}</td>
                                        </tr>
                                        {expanded === m.id && (
                                            <tr className="bg-gray-50/70 dark:bg-gray-900/30">
                                                <td colSpan={7} className="px-4 py-4">
                                                    <div className="grid gap-4 lg:grid-cols-2">
                                                        <div>
                                                            <p className="mb-1 text-xs font-medium uppercase tracking-wide text-gray-500">Pregunta</p>
                                                            <p className="whitespace-pre-wrap text-gray-900 dark:text-gray-100">{m.question}</p>
                                                        </div>
                                                        <div>
                                                            <p className="mb-1 text-xs font-medium uppercase tracking-wide text-gray-500">{`Respuesta de ${ASSISTANT_NAME}`}</p>
                                                            <p className="whitespace-pre-wrap text-gray-900 dark:text-gray-100">{m.answer || '\u2014'}</p>
                                                            {m.destination_route && (
                                                                <p className="mt-2 font-mono text-xs text-[#3E0FA8] dark:text-violet-200">{m.destination_route}</p>
                                                            )}
                                                        </div>
                                                    </div>
                                                    <p className="mt-3 font-mono text-[11px] text-gray-400">{`${m.model || 'sin modelo'} \u00b7 conversaci\u00f3n ${m.conversation_id}`}</p>
                                                </td>
                                            </tr>
                                        )}
                                    </Fragment>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
                {list && list.total > 0 && (
                    <div className="flex items-center justify-between border-t border-gray-200 px-4 py-3 text-sm text-gray-600 dark:border-gray-700 dark:text-gray-300">
                        <span className="tabular-nums">{`P\u00e1gina ${list.page} de ${list.total_pages} \u00b7 ${numberFormat.format(list.total)} mensajes`}</span>
                        <div className="flex gap-2">
                            <button
                                type="button"
                                disabled={review.page <= 1}
                                onClick={() => review.setPage(review.page - 1)}
                                className="rounded-lg border border-gray-200 px-3 py-1.5 disabled:opacity-40 dark:border-gray-700"
                            >
                                Anterior
                            </button>
                            <button
                                type="button"
                                disabled={review.page >= list.total_pages}
                                onClick={() => review.setPage(review.page + 1)}
                                className="rounded-lg border border-gray-200 px-3 py-1.5 disabled:opacity-40 dark:border-gray-700"
                            >
                                Siguiente
                            </button>
                        </div>
                    </div>
                )}
            </div>
            <p className="text-xs text-gray-400">
                {'Costo estimado con USD 0,15 por mill\u00f3n de tokens de entrada y USD 1,20 por mill\u00f3n de salida (Qwen3 Next en Bedrock).'}
            </p>
        </div>
    );
}
