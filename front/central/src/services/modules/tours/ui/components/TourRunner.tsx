'use client';

import { useCallback, useEffect } from 'react';
import type { TourDefinition } from '../../domain/types';
import { TourSpotlight } from './TourSpotlight';

interface Props {
    tour: TourDefinition;
    stepIndex: number;
    onStepChange: (index: number) => void;
    onNavigate: (route: string) => void;
    onSkip: (index: number) => void;
    onComplete: (index: number) => void;
    onSkipAll: () => void;
}

export function TourRunner({ tour, stepIndex, onStepChange, onNavigate, onSkip, onComplete, onSkipAll }: Props) {
    const step = tour.steps[stepIndex];

    const avanzar = useCallback(() => {
        if (stepIndex >= tour.steps.length - 1) {
            onComplete(stepIndex);
            return;
        }
        onStepChange(stepIndex + 1);
    }, [stepIndex, tour.steps.length, onComplete, onStepChange]);

    const retroceder = useCallback(() => {
        if (stepIndex > 0) onStepChange(stepIndex - 1);
    }, [stepIndex, onStepChange]);

    useEffect(() => {
        if (step?.route) onNavigate(step.route);
    }, [step, onNavigate]);

    useEffect(() => {
        const handleKey = (e: KeyboardEvent) => {
            if (e.key === 'Escape') onSkip(stepIndex);
            if (e.key === 'ArrowRight') avanzar();
            if (e.key === 'ArrowLeft') retroceder();
        };
        window.addEventListener('keydown', handleKey);
        return () => window.removeEventListener('keydown', handleKey);
    }, [avanzar, retroceder, onSkip, stepIndex]);

    if (!step) return null;

    return (
        <TourSpotlight
            key={step.id}
            step={step}
            stepNumber={stepIndex + 1}
            totalSteps={tour.steps.length}
            onNext={avanzar}
            onPrev={retroceder}
            onClose={() => onSkip(stepIndex)}
            onUnavailable={avanzar}
            onSkipAll={onSkipAll}
        />
    );
}
