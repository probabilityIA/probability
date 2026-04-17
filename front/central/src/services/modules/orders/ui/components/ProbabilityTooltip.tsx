'use client';

import { useState, useRef, useEffect, ReactNode } from 'react';
import { createPortal } from 'react-dom';
import { Fingerprint } from 'lucide-react';

interface ProbabilityTooltipProps {
    children: ReactNode;
}

export function ProbabilityTooltip({ children }: ProbabilityTooltipProps) {
    const [open, setOpen] = useState(false);
    const [position, setPosition] = useState<{ top: number; left: number; placement: 'top' | 'bottom' }>({ top: 0, left: 0, placement: 'top' });
    const triggerRef = useRef<HTMLDivElement>(null);
    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        setMounted(true);
    }, []);

    useEffect(() => {
        if (!open || !triggerRef.current) return;

        const rect = triggerRef.current.getBoundingClientRect();
        const viewportHeight = window.innerHeight;
        const tooltipHeight = 240;
        const tooltipWidth = 260;

        const spaceAbove = rect.top;
        const spaceBelow = viewportHeight - rect.bottom;
        const placement: 'top' | 'bottom' = spaceAbove > tooltipHeight || spaceAbove > spaceBelow ? 'top' : 'bottom';

        const top = placement === 'top'
            ? rect.top - 8
            : rect.bottom + 8;

        const left = Math.min(
            Math.max(rect.left + rect.width / 2 - tooltipWidth / 2, 8),
            window.innerWidth - tooltipWidth - 8
        );

        setPosition({ top, left, placement });
    }, [open]);

    return (
        <div
            ref={triggerRef}
            className="relative inline-flex"
            onMouseEnter={() => setOpen(true)}
            onMouseLeave={() => setOpen(false)}
        >
            <div className="flex items-center justify-center w-6 h-6 rounded-full bg-indigo-100 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-300 cursor-pointer hover:bg-indigo-200 dark:hover:bg-indigo-800 transition-colors">
                <Fingerprint size={14} strokeWidth={2} />
            </div>
            {mounted && open && createPortal(
                <div
                    className="fixed z-[9999] pointer-events-none"
                    style={{
                        top: position.top,
                        left: position.left,
                        transform: position.placement === 'top' ? 'translateY(-100%)' : 'translateY(0)',
                    }}
                >
                    <div className="bg-gray-900 text-white text-xs rounded-lg py-3 px-4 shadow-2xl w-[260px] border border-gray-700">
                        {children}
                    </div>
                </div>,
                document.body
            )}
        </div>
    );
}
