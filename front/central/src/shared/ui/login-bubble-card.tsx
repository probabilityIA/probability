'use client';

import { ReactNode, useEffect, useState } from 'react';

const DELAY_MS = 900;

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
            className={`w-full max-w-[380px] rounded-2xl px-7 py-8 ${
                isDark ? 'bg-[#15102b]' : 'bg-white'
            }`}
            style={{
                boxShadow: '0 24px 70px rgba(0, 0, 0, 0.35)',
                opacity: visible ? 1 : 0,
                transform: visible ? 'translateY(0)' : 'translateY(18px)',
                transition: 'opacity 700ms ease-out, transform 800ms cubic-bezier(0.16, 1, 0.3, 1)',
            }}
        >
            {children}
        </div>
    );
}
