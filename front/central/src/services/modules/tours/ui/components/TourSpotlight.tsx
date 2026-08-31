'use client';

import { useEffect, useMemo } from 'react';
import type { TourStep } from '../../domain/types';
import { useTargetRect } from '../hooks/use-target-rect';
import { TourCard } from './TourCard';

interface Props {
    step: TourStep;
    stepNumber: number;
    totalSteps: number;
    onNext: () => void;
    onPrev: () => void;
    onClose: () => void;
    onUnavailable: () => void;
    onSkipAll: () => void;
}

const PADDING = 6;
const CARD_WIDTH = 370;
const CARD_HEIGHT = 210;
const GAP = 14;

export function TourSpotlight({ step, stepNumber, totalSteps, onNext, onPrev, onClose, onUnavailable, onSkipAll }: Props) {
    const anclado = Boolean(step.target);
    const { rect, notFound } = useTargetRect(step.target ?? null, anclado);

    useEffect(() => {
        if (notFound) onUnavailable();
    }, [notFound, onUnavailable]);

    useEffect(() => {
        if (!step.target) return;
        const el = document.querySelector(step.target) as HTMLElement | null;
        if (el) el.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }, [step.target]);

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
        if (!hole) {
            return { top: '50%', left: '50%', transform: 'translate(-50%, -50%)' };
        }

        const viewportW = typeof window !== 'undefined' ? window.innerWidth : 1280;
        const viewportH = typeof window !== 'undefined' ? window.innerHeight : 800;

        const espacioAbajo = viewportH - (hole.top + hole.height);
        const espacioDerecha = viewportW - (hole.left + hole.width);

        let placement = step.placement ?? 'bottom';
        if (placement === 'bottom' && espacioAbajo < CARD_HEIGHT + GAP) {
            placement = hole.top > CARD_HEIGHT + GAP ? 'top' : 'right';
        }
        if (placement === 'right' && espacioDerecha < CARD_WIDTH + GAP) placement = 'left';
        if (placement === 'left' && hole.left < CARD_WIDTH + GAP) placement = 'bottom';

        let top: number;
        let left: number;

        switch (placement) {
            case 'top':
                top = hole.top - CARD_HEIGHT - GAP;
                left = hole.left + hole.width / 2 - CARD_WIDTH / 2;
                break;
            case 'left':
                top = hole.top;
                left = hole.left - CARD_WIDTH - GAP;
                break;
            case 'right':
                top = hole.top;
                left = hole.left + hole.width + GAP;
                break;
            default:
                top = hole.top + hole.height + GAP;
                left = hole.left + hole.width / 2 - CARD_WIDTH / 2;
        }

        left = Math.max(12, Math.min(left, viewportW - CARD_WIDTH - 12));
        top = Math.max(12, Math.min(top, viewportH - CARD_HEIGHT - 12));

        return { top, left };
    }, [hole, step.placement]);

    return (
        <div className="fixed inset-0 z-[100]">
            {hole ? (
                <div
                    className="pointer-events-none absolute transition-all duration-200"
                    style={{
                        boxShadow: '0 0 0 9999px rgba(17, 24, 39, 0.55)',
                        top: hole.top,
                        left: hole.left,
                        width: hole.width,
                        height: hole.height,
                        borderRadius: 10,
                        border: '2px solid var(--color-primary, #7c3aed)',
                    }}
                />
            ) : (
                <div className="absolute inset-0" style={{ background: 'rgba(17, 24, 39, 0.55)' }} />
            )}

            <TourCard
                title={step.title}
                body={step.body}
                stepNumber={stepNumber}
                totalSteps={totalSteps}
                style={cardStyle}
                onNext={onNext}
                onPrev={onPrev}
                onClose={onClose}
                onSkipAll={onSkipAll}
            />
        </div>
    );
}
