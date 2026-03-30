'use client';

import { useState, useTransition } from 'react';
import { loginAction } from '@/services/monitoring/infra/actions';

export default function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isPending, startTransition] = useTransition();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        startTransition(async () => {
            const result = await loginAction(email, password);
            if (result?.error) {
                setError(result.error);
            }
        });
    };

    return (
        <div className="login-glow min-h-screen flex items-center justify-center px-4">
            <div className="w-full max-w-sm">
                {/* Logo / Title */}
                <div className="text-center mb-8">
                    <div
                        className="inline-flex items-center justify-center w-16 h-16 rounded-2xl mb-4"
                        style={{
                            background: 'linear-gradient(135deg, #00f0ff15, #00ff8815)',
                            border: '1px solid #00f0ff20',
                            boxShadow: '0 0 40px #00f0ff08, 0 0 80px #00f0ff04',
                        }}
                    >
                        <svg className="w-8 h-8" style={{ color: '#00f0ff' }} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M5.25 14.25h13.5m-13.5 0a3 3 0 01-3-3m3 3a3 3 0 100 6h13.5a3 3 0 100-6m-16.5-3a3 3 0 013-3h13.5a3 3 0 013 3m-19.5 0a4.5 4.5 0 01.9-2.7L5.737 5.1a3.375 3.375 0 012.7-1.35h7.126c1.062 0 2.062.5 2.7 1.35l2.587 3.45a4.5 4.5 0 01.9 2.7m0 0a3 3 0 01-3 3m0 3h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008zm-3 6h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008z" />
                        </svg>
                    </div>
                    <h1 className="text-2xl font-semibold tracking-tight" style={{ color: '#e4e4ef' }}>
                        Monitoring
                    </h1>
                    <p className="text-xs mt-1 uppercase tracking-widest" style={{ color: '#55556a' }}>
                        Probability Infrastructure
                    </p>
                </div>

                {/* Form */}
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div
                        className="rounded-xl p-6 space-y-4"
                        style={{
                            background: '#12121a',
                            border: '1px solid #1e1e2e',
                            boxShadow: '0 4px 24px rgba(0,0,0,0.3)',
                        }}
                    >
                        <div>
                            <label className="block text-[11px] font-medium mb-1.5 uppercase tracking-wider" style={{ color: '#55556a' }}>
                                Email
                            </label>
                            <input
                                type="email"
                                value={email}
                                onChange={e => setEmail(e.target.value)}
                                required
                                autoFocus
                                className="w-full px-3 py-2.5 rounded-lg text-sm outline-none transition-all font-mono"
                                style={{
                                    background: '#08080d',
                                    border: '1px solid #1e1e2e',
                                    color: '#e4e4ef',
                                }}
                                onFocus={e => {
                                    e.target.style.borderColor = '#00f0ff40';
                                    e.target.style.boxShadow = '0 0 12px #00f0ff08';
                                }}
                                onBlur={e => {
                                    e.target.style.borderColor = '#1e1e2e';
                                    e.target.style.boxShadow = 'none';
                                }}
                                placeholder="admin@example.com"
                            />
                        </div>

                        <div>
                            <label className="block text-[11px] font-medium mb-1.5 uppercase tracking-wider" style={{ color: '#55556a' }}>
                                Password
                            </label>
                            <input
                                type="password"
                                value={password}
                                onChange={e => setPassword(e.target.value)}
                                required
                                className="w-full px-3 py-2.5 rounded-lg text-sm outline-none transition-all font-mono"
                                style={{
                                    background: '#08080d',
                                    border: '1px solid #1e1e2e',
                                    color: '#e4e4ef',
                                }}
                                onFocus={e => {
                                    e.target.style.borderColor = '#00f0ff40';
                                    e.target.style.boxShadow = '0 0 12px #00f0ff08';
                                }}
                                onBlur={e => {
                                    e.target.style.borderColor = '#1e1e2e';
                                    e.target.style.boxShadow = 'none';
                                }}
                                placeholder="••••••••"
                            />
                        </div>
                    </div>

                    {error && (
                        <div
                            className="text-sm px-3 py-2 rounded-lg"
                            style={{
                                background: '#ff336610',
                                color: '#ff3366',
                                border: '1px solid #ff336625',
                                boxShadow: '0 0 12px #ff336608',
                            }}
                        >
                            {error}
                        </div>
                    )}

                    <button
                        type="submit"
                        disabled={isPending}
                        className="w-full py-2.5 rounded-lg text-sm font-semibold transition-all cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed uppercase tracking-wider"
                        style={{
                            background: 'linear-gradient(135deg, #00f0ff, #00c8ff)',
                            color: '#0a0a0f',
                            boxShadow: '0 0 20px #00f0ff15',
                        }}
                        onMouseEnter={e => {
                            if (!isPending) e.currentTarget.style.boxShadow = '0 0 30px #00f0ff25';
                        }}
                        onMouseLeave={e => {
                            e.currentTarget.style.boxShadow = '0 0 20px #00f0ff15';
                        }}
                    >
                        {isPending ? 'Signing in...' : 'Sign In'}
                    </button>
                </form>
            </div>
        </div>
    );
}
