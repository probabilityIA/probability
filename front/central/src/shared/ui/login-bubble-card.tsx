'use client';

import { ReactNode, useEffect, useState } from 'react';

const DELAY_MS = 900;

interface Props {
    children: ReactNode;
}

export function LoginBubbleCard({ children }: Props) {
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
            className="w-full max-w-[380px] px-2 py-2 text-gray-900"
            style={{
                opacity: visible ? 1 : 0,
                transform: visible ? 'translateY(0)' : 'translateY(18px)',
                transition: 'opacity 700ms ease-out, transform 800ms cubic-bezier(0.16, 1, 0.3, 1)',
            }}
        >
            {children}
        </div>
    );
}
