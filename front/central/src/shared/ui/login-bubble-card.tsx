'use client';

import { ReactNode } from 'react';

interface Props {
    children: ReactNode;
    isDark?: boolean;
    visible: boolean;
}

export function LoginBubbleCard({ children, isDark = false, visible }: Props) {
    return (
        <div
            className={`relative z-30 h-full w-full lg:w-[42%] lg:max-w-[560px] flex items-center justify-center px-6 sm:px-10 lg:rounded-r-[56px] ${
                isDark ? 'bg-[#0f0a1e]' : 'bg-white'
            }`}
            style={{
                boxShadow: visible ? '24px 0 80px rgba(0, 0, 0, 0.28)' : 'none',
                opacity: visible ? 1 : 0,
                transform: visible ? 'translateX(0)' : 'translateX(-100%)',
                transition:
                    'opacity 600ms ease-out, transform 1000ms cubic-bezier(0.16, 1, 0.3, 1), box-shadow 1000ms ease-out',
                pointerEvents: visible ? 'auto' : 'none',
            }}
        >
            <div
                className="w-full max-w-[360px]"
                style={{
                    opacity: visible ? 1 : 0,
                    transform: visible ? 'translateY(0)' : 'translateY(16px)',
                    transition: 'opacity 700ms ease-out 400ms, transform 700ms ease-out 400ms',
                }}
            >
                {children}
            </div>
        </div>
    );
}
