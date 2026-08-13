'use client';

import { ChevronLeft, ChevronRight } from 'lucide-react';
import { ACCENT, ACCENT_BORDER, ACCENT_SOFT, CARD_BORDER } from '../panel-theme';

export const TAMANOS_PAGINA = [10, 25, 50, 100];

interface PanelPagerProps {
    page: number;
    totalPages: number;
    total: number;
    shown: number;
    noun: string;
    pageSize: number;
    onPage: (page: number) => void;
    onPageSize?: (pageSize: number) => void;
}

function paginas(page: number, totalPages: number): (number | 'gap')[] {
    if (totalPages <= 7) {
        return Array.from({ length: totalPages }, (_, i) => i + 1);
    }
    const out: (number | 'gap')[] = [1];
    const desde = Math.max(2, page - 1);
    const hasta = Math.min(totalPages - 1, page + 1);
    if (desde > 2) out.push('gap');
    for (let i = desde; i <= hasta; i++) out.push(i);
    if (hasta < totalPages - 1) out.push('gap');
    out.push(totalPages);
    return out;
}

function Total({ total, noun, detalle }: { total: number; noun: string; detalle?: string }) {
    return (
        <span
            className="inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1 text-[11.5px] font-semibold"
            style={{ borderColor: ACCENT_BORDER, backgroundColor: ACCENT_SOFT, color: ACCENT }}
        >
            <span className="text-[13px] font-black">{total.toLocaleString('es-CO')}</span>
            {noun}
            {detalle && <span className="font-normal opacity-70">{detalle}</span>}
        </span>
    );
}

export function PanelPager({ page, totalPages, total, shown, noun, pageSize, onPage, onPageSize }: PanelPagerProps) {
    const selector = onPageSize && (
        <label className="inline-flex items-center gap-1 text-[11px] text-gray-400">
            ver
            <select
                value={pageSize}
                onChange={event => onPageSize(Number(event.target.value))}
                className="rounded-lg border bg-white px-1.5 py-0.5 text-[11px] font-semibold text-gray-600 outline-none dark:bg-gray-800 dark:text-gray-200"
                style={{ borderColor: CARD_BORDER }}
            >
                {TAMANOS_PAGINA.map(tamano => (
                    <option key={tamano} value={tamano}>{tamano}</option>
                ))}
            </select>
        </label>
    );

    if (totalPages <= 1) {
        return (
            <div className="flex flex-wrap items-center justify-between gap-2">
                <Total total={total} noun={noun} />
                {selector}
            </div>
        );
    }

    const desde = total === 0 ? 0 : (page - 1) * pageSize + 1;
    const hasta = (page - 1) * pageSize + shown;
    const fmt = (n: number) => n.toLocaleString('es-CO');

    return (
        <div className="flex flex-wrap items-center justify-between gap-2">
            <div className="flex flex-wrap items-center gap-2">
                <Total total={total} noun={noun} detalle={`viendo ${fmt(desde)}–${fmt(hasta)}`} />
                {selector}
            </div>

            <div className="flex items-center gap-1">
                <button
                    onClick={() => onPage(page - 1)}
                    disabled={page <= 1}
                    aria-label="Pagina anterior"
                    className="flex h-7 w-7 items-center justify-center rounded-lg border bg-white text-gray-500 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40 dark:bg-gray-800 dark:text-gray-300"
                    style={{ borderColor: CARD_BORDER }}
                >
                    <ChevronLeft size={13} />
                </button>

                {paginas(page, totalPages).map((p, i) =>
                    p === 'gap' ? (
                        <span key={`gap-${i}`} className="px-1 text-[11px] text-gray-400">
                            ...
                        </span>
                    ) : (
                        <button
                            key={p}
                            onClick={() => onPage(p)}
                            aria-current={p === page ? 'page' : undefined}
                            className="h-7 min-w-7 rounded-lg border px-2 text-[11.5px] font-semibold transition-colors"
                            style={
                                p === page
                                    ? { backgroundColor: ACCENT, borderColor: ACCENT, color: 'var(--color-on-primary)' }
                                    : { borderColor: ACCENT_BORDER, color: ACCENT, backgroundColor: 'white' }
                            }
                        >
                            {p}
                        </button>
                    ),
                )}

                <button
                    onClick={() => onPage(page + 1)}
                    disabled={page >= totalPages}
                    aria-label="Pagina siguiente"
                    className="flex h-7 w-7 items-center justify-center rounded-lg border bg-white text-gray-500 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40 dark:bg-gray-800 dark:text-gray-300"
                    style={{ borderColor: CARD_BORDER }}
                >
                    <ChevronRight size={13} />
                </button>
            </div>
        </div>
    );
}
