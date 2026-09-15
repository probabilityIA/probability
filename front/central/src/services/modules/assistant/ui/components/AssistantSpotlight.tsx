'use client';

import { useEffect, useMemo, useRef } from 'react';
import { useTargetRect } from '@/services/modules/tours/ui/hooks/use-target-rect';
import { AssistantAvatar } from './AssistantAvatar';
import { ASSISTANT_NAME } from '../../app/use-cases';
import type { AssistantHighlight } from '../../domain/types';

interface AssistantSpotlightProps {
    highlight: AssistantHighlight;
    onClose: () => void;
}

const PADDING = 6;
const CARD_WIDTH = 320;
const CARD_HEIGHT = 190;
const GAP = 14;
const WAIT_MS = 8000;

export function AssistantSpotlight({ highlight, onClose }: AssistantSpotlightProps) {
    const { rect, notFound } = useTargetRect(highlight.target, true, WAIT_MS);
    const scrolled = useRef(false);

    useEffect(() => {
        const onKey = (event: KeyboardEvent) => {
            if (event.key === 'Escape') onClose();
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, [onClose]);

    useEffect(() => {
        if (!rect) return;
        const element = document.querySelector(highlight.target) as HTMLElement | null;
        if (!element) return;
        if (!scrolled.current) {
            scrolled.current = true;
            element.scrollIntoView({ behavior: 'smooth', block: 'center', inline: 'center' });
        }
        element.addEventListener('click', onClose, { once: true });
        return () => element.removeEventListener('click', onClose);
    }, [rect, highlight.target, onClose]);

    const hole = useMemo(() => {
        if (!rect) return null;
        return {
            top: rect.top - PADDING,
            left: rect.left - PADDING,
            width: rect.width + PADDING * 2,
            height: rect.height + PADDING * 2,
        };
    }, [rect]);

    const cardStyle = useMemo((): React.CSSProperties => {
        if (!hole) return { top: '50%', left: '50%', transform: 'translate(-50%, -50%)' };

        const viewportW = window.innerWidth;
        const viewportH = window.innerHeight;
        let top = hole.top + hole.height + GAP;
        if (top + CARD_HEIGHT > viewportH - 12) top = hole.top - CARD_HEIGHT - GAP;
        let left = hole.left + hole.width - CARD_WIDTH;
        left = Math.max(12, Math.min(left, viewportW - CARD_WIDTH - 12));
        top = Math.max(12, Math.min(top, viewportH - CARD_HEIGHT - 12));
        return { top, left };
    }, [hole]);

    const searching = !rect && !notFound;
    const title = notFound ? 'No lo encontr\u00e9 en esta pantalla' : highlight.title;
    const body = notFound
        ? 'Puede que todav\u00eda no haya registros, que ya est\u00e9n procesados o que tu rol no tenga ese bot\u00f3n. Si eres super admin, primero elige un negocio.'
        : highlight.hint;

    return (
        <div className="pointer-events-none fixed inset-0 z-[90]" role="dialog" aria-label={`${ASSISTANT_NAME} te muestra d\u00f3nde`}>
            {hole ? (
                <div
                    className="absolute rounded-[10px] transition-all duration-300"
                    style={{
                        top: hole.top,
                        left: hole.left,
                        width: hole.width,
                        height: hole.height,
                        boxShadow: '0 0 0 9999px rgba(17, 24, 39, 0.5)',
                        border: '2px solid #5B1BE6',
                    }}
                >
                    <span className="absolute inset-0 animate-ping rounded-[10px] border-2 border-[#8B63FF] motion-reduce:animate-none" />
                </div>
            ) : (
                !searching && <div className="absolute inset-0 bg-gray-900/50" />
            )}

            {!searching && (
                <div
                    className="pointer-events-auto fixed w-[320px] max-w-[calc(100vw-24px)] rounded-2xl border border-gray-200 bg-white p-4 shadow-2xl dark:border-gray-700 dark:bg-gray-800"
                    style={cardStyle}
                >
                    <div className="flex items-start gap-3">
                        <AssistantAvatar size={32} mood={notFound ? 'idle' : 'pointing'} />
                        <div className="min-w-0">
                            <p className="text-sm font-bold text-gray-900 dark:text-white">{title}</p>
                            <p className="mt-1 text-[13px] leading-relaxed text-gray-600 dark:text-gray-300">{body}</p>
                        </div>
                    </div>
                    <div className="mt-3 flex justify-end">
                        <button
                            type="button"
                            onClick={onClose}
                            className="rounded-lg bg-[#5B1BE6] px-3 py-1.5 text-xs font-semibold text-white hover:bg-[#4A12C9] dark:bg-[#7148FF]"
                        >
                            Entendido
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}
