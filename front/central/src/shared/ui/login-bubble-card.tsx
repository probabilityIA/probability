'use client';

import { ReactNode, useEffect, useState } from 'react';

const DELAY_MS = 2600;

interface Props {
    children: ReactNode;
    isDark?: boolean;
}

export function LoginBubbleCard({ children, isDark = false }: Props) {
    const [visible, setVisible] = useState(false);

    useEffect(() => {
        const reduce = typeof window !== 'undefined'
            && window.matchMedia?.('(prefers-reduced-motion: reduce)').matches;

        if (reduce) {
            setVisible(true);
            return;
        }

        const id = setTimeout(() => setVisible(true), DELAY_MS);
        return () => clearTimeout(id);
    }, []);

    return (
        <div
            className={`w-full max-w-md rounded-[32px] px-7 py-9 sm:px-10 ${
                isDark ? 'ring-1 ring-white/25' : 'ring-1 ring-white/60'
            }`}
            style={{
                backdropFilter: visible ? 'blur(22px) saturate(140%)' : 'blur(0px)',
                WebkitBackdropFilter: visible ? 'blur(22px) saturate(140%)' : 'blur(0px)',
                backgroundColor: visible
                    ? (isDark ? 'rgba(255, 255, 255, 0.28)' : 'rgba(255, 255, 255, 0.45)')
                    : 'rgba(255, 255, 255, 0)',
                boxShadow: visible
                    ? '0 30px 80px rgba(0, 0, 0, 0.22), inset 0 1px 0 rgba(255,255,255,0.25)'
                    : '0 0 0 rgba(0, 0, 0, 0)',
                opacity: visible ? 1 : 0,
                transform: visible ? 'scale(1) translateY(0)' : 'scale(0.82) translateY(28px)',
                transition:
                    'opacity 900ms ease-out, transform 1100ms cubic-bezier(0.22, 1.4, 0.36, 1), backdrop-filter 1200ms ease-out, box-shadow 1200ms ease-out, background-color 1200ms ease-out',
                pointerEvents: visible ? 'auto' : 'none',
            }}
        >
            {children}
        </div>
    );
}
