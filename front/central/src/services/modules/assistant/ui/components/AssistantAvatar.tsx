'use client';

import type { AvatarMood } from '../../domain/types';

interface AssistantAvatarProps {
    mood?: AvatarMood;
    size?: number;
    className?: string;
}

export function AssistantAvatar({ mood = 'idle', size = 40, className = '' }: AssistantAvatarProps) {
    return (
        <span className={`inline-flex shrink-0 ${className}`} style={{ width: size, height: size }}>
            <svg viewBox="0 0 64 64" width={size} height={size} aria-hidden="true" className={`avatar ${mood}`}>
                <circle cx="32" cy="32" r="32" className="fill-[#5B1BE6] dark:fill-[#7148FF]" />
                <rect className="pill" x="31" y="17" width="12" height="36" rx="6" fill="#ffffff" />
                <g className="head">
                    <circle cx="21" cy="21" r="8" fill="#ffffff" />
                    <g className="eyes">
                        <circle cx="18.6" cy="21" r="1.5" className="fill-[#5B1BE6] dark:fill-[#7148FF]" />
                        <circle cx="23.4" cy="21" r="1.5" className="fill-[#5B1BE6] dark:fill-[#7148FF]" />
                    </g>
                </g>
            </svg>
            <style jsx>{`
                .pill {
                    transform-box: view-box;
                    transform-origin: 37px 35px;
                    transform: rotate(28deg);
                    transition: transform 0.5s cubic-bezier(0.3, 1.4, 0.5, 1);
                }
                .head {
                    transition: transform 0.5s cubic-bezier(0.3, 1.4, 0.5, 1);
                }
                .pointing .pill {
                    transform: rotate(44deg);
                }
                .pointing .head {
                    transform: translate(4px, -1px);
                }
                .thinking .head {
                    animation: bob 1.1s ease-in-out infinite;
                }
                .thinking .eyes {
                    animation: look 2.2s ease-in-out infinite;
                }
                @keyframes bob {
                    0%, 100% { transform: translateY(0); }
                    50% { transform: translateY(-3px); }
                }
                @keyframes look {
                    0%, 100% { transform: translateX(0); }
                    50% { transform: translateX(1.5px); }
                }
                @media (prefers-reduced-motion: reduce) {
                    .pill, .head { transition: none; }
                    .thinking .head, .thinking .eyes { animation: none; }
                }
            `}</style>
        </span>
    );
}
