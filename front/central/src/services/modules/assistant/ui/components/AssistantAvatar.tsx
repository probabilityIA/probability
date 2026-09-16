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
                <defs>
                    <radialGradient id="via-sphere" cx="0.35" cy="0.3" r="0.9">
                        <stop offset="0" stopColor="#6D2BFF" />
                        <stop offset="0.55" stopColor="#4310C8" />
                        <stop offset="1" stopColor="#2A0A86" />
                    </radialGradient>
                    <linearGradient id="via-sheen" x1="0" y1="0" x2="1" y2="1">
                        <stop offset="0" stopColor="#B08CFF" stopOpacity="0.55" />
                        <stop offset="0.5" stopColor="#8B63FF" stopOpacity="0.12" />
                        <stop offset="1" stopColor="#5B1BE6" stopOpacity="0" />
                    </linearGradient>
                    <filter id="via-glow" x="-60%" y="-60%" width="220%" height="220%">
                        <feGaussianBlur stdDeviation="1.6" result="blur" />
                        <feMerge>
                            <feMergeNode in="blur" />
                            <feMergeNode in="SourceGraphic" />
                        </feMerge>
                    </filter>
                    <clipPath id="via-clip">
                        <circle cx="32" cy="32" r="32" />
                    </clipPath>
                </defs>
                <circle cx="32" cy="32" r="32" fill="url(#via-sphere)" />
                <g clipPath="url(#via-clip)">
                    <path d="M-6 30 C 14 8, 40 4, 70 22 L 70 -10 L -6 -10 Z" fill="url(#via-sheen)" />
                    <path d="M-10 46 C 20 64, 44 62, 74 38 L 74 80 L -10 80 Z" fill="#8B63FF" opacity="0.28" />
                </g>
                <g className="face">
                    <g className="eyes" filter="url(#via-glow)" fill="none" stroke="#F3ECFF" strokeWidth="4.2" strokeLinecap="round">
                        <path className="eye" d="M17 31 Q 22.5 24.5, 28 31" />
                        <path className="eye" d="M36 31 Q 41.5 24.5, 47 31" />
                    </g>
                    <g className="mouth" fill="none" stroke="#F3ECFF" strokeWidth="2.6" strokeLinecap="round" opacity="0">
                        <path d="M27 42 Q 32 46, 37 42" />
                    </g>
                </g>
            </svg>
            <style jsx>{`
                .avatar {
                    animation: float 4.5s ease-in-out infinite;
                }
                .face {
                    transform-box: view-box;
                    transform-origin: 32px 32px;
                    transition: transform 0.5s cubic-bezier(0.3, 1.4, 0.5, 1);
                }
                .eyes {
                    transform-box: view-box;
                    transform-origin: 32px 30px;
                    animation: lookaround 7s ease-in-out infinite;
                }
                .eye {
                    transform-box: fill-box;
                    transform-origin: center;
                    animation: blink 4.6s ease-in-out infinite;
                }
                .eye:last-child {
                    animation-delay: 0.05s;
                }
                .idle .face {
                    animation: tilt 7s ease-in-out infinite;
                }
                .pointing .face {
                    transform: translate(3px, -1px) rotate(-6deg);
                }
                .thinking .face {
                    animation: bob 1.1s ease-in-out infinite;
                }
                .thinking .eyes {
                    animation: look 1.6s ease-in-out infinite;
                }
                .thinking .eye {
                    animation: blink 2.2s ease-in-out infinite;
                }
                .happy {
                    animation: jump 0.32s cubic-bezier(0.3, 0, 0.4, 1) 5;
                }
                .happy .face {
                    animation: none;
                    transform: translateY(-0.5px);
                }
                .happy .eyes {
                    animation: none;
                }
                .happy .mouth {
                    opacity: 1;
                    transform-box: view-box;
                    transform-origin: 32px 43px;
                    transform: scaleX(1.3) scaleY(1.2);
                }
                .talking .mouth {
                    opacity: 1;
                    transform-box: view-box;
                    transform-origin: 32px 43px;
                    animation: talk 0.5s ease-in-out infinite;
                }
                @keyframes float {
                    0%, 100% { transform: translateY(0); }
                    50% { transform: translateY(-1.5px); }
                }
                @keyframes lookaround {
                    0%, 8% { transform: translateX(0); }
                    16%, 34% { transform: translateX(-3px); }
                    42%, 50% { transform: translateX(0); }
                    58%, 76% { transform: translateX(3px); }
                    84%, 100% { transform: translateX(0); }
                }
                @keyframes tilt {
                    0%, 8%, 42%, 50%, 84%, 100% { transform: rotate(0deg); }
                    16%, 34% { transform: rotate(-4deg); }
                    58%, 76% { transform: rotate(4deg); }
                }
                @keyframes blink {
                    0%, 90%, 100% { transform: scaleY(1); }
                    94% { transform: scaleY(0.1); }
                }
                @keyframes bob {
                    0%, 100% { transform: translateY(0); }
                    50% { transform: translateY(-2.5px); }
                }
                @keyframes look {
                    0%, 100% { transform: translateX(-2px); }
                    50% { transform: translateX(2px); }
                }
                @keyframes jump {
                    0%, 100% { transform: translateY(0) scale(1, 1); }
                    20% { transform: translateY(1px) scale(1.08, 0.9); }
                    55% { transform: translateY(-9px) scale(0.95, 1.08); }
                    85% { transform: translateY(0) scale(1.04, 0.94); }
                }
                @keyframes talk {
                    0%, 100% { transform: scaleY(0.6); }
                    50% { transform: scaleY(1.3); }
                }
                @media (prefers-reduced-motion: reduce) {
                    .avatar, .face, .eyes, .eye, .mouth { animation: none !important; transition: none; }
                }
            `}</style>
        </span>
    );
}
