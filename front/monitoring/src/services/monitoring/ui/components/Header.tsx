'use client';

import { logoutAction } from '../../infra/actions';
import { useTransition } from 'react';

interface HeaderProps {
    userName?: string;
}

export function Header({ userName }: HeaderProps) {
    const [isPending, startTransition] = useTransition();

    return (
        <header
            className="flex items-center justify-between px-6 py-3"
            style={{ borderBottom: '1px solid #1e1e2e' }}
        >
            <div className="flex items-center gap-3">
                <div
                    className="w-8 h-8 rounded-lg flex items-center justify-center"
                    style={{ background: 'linear-gradient(135deg, #00f0ff20, #00ff8820)', border: '1px solid #00f0ff30' }}
                >
                    <svg className="w-4 h-4" style={{ color: '#00f0ff' }} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M5.25 14.25h13.5m-13.5 0a3 3 0 01-3-3m3 3a3 3 0 100 6h13.5a3 3 0 100-6m-16.5-3a3 3 0 013-3h13.5a3 3 0 013 3m-19.5 0a4.5 4.5 0 01.9-2.7L5.737 5.1a3.375 3.375 0 012.7-1.35h7.126c1.062 0 2.062.5 2.7 1.35l2.587 3.45a4.5 4.5 0 01.9 2.7m0 0a3 3 0 01-3 3m0 3h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008zm-3 6h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008z" />
                    </svg>
                </div>
                <div>
                    <h1 className="text-sm font-semibold" style={{ color: '#e4e4ef' }}>Monitoring</h1>
                    <p className="text-[10px]" style={{ color: '#55556a' }}>Probability Infrastructure</p>
                </div>
            </div>

            <div className="flex items-center gap-4">
                {userName && (
                    <span className="text-xs" style={{ color: '#8888a0' }}>{userName}</span>
                )}
                <button
                    onClick={() => startTransition(() => logoutAction())}
                    disabled={isPending}
                    className="text-xs px-3 py-1.5 rounded-md transition-colors cursor-pointer"
                    style={{ color: '#8888a0', border: '1px solid #1e1e2e' }}
                    onMouseEnter={e => {
                        e.currentTarget.style.borderColor = '#ff336630';
                        e.currentTarget.style.color = '#ff3366';
                    }}
                    onMouseLeave={e => {
                        e.currentTarget.style.borderColor = '#1e1e2e';
                        e.currentTarget.style.color = '#8888a0';
                    }}
                >
                    Logout
                </button>
            </div>
        </header>
    );
}
