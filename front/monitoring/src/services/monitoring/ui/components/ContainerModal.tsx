'use client';

import { useEffect } from 'react';
import type { Container } from '../../domain/types';
import { LogViewer } from './LogViewer';
import { ActionButtons } from './ActionButtons';
import { StatsBar } from './StatsBar';

interface ContainerModalProps {
    container: Container;
    onClose: () => void;
}

function formatUptime(startedAt: string): string {
    if (!startedAt) return '--';
    const diff = Date.now() - new Date(startedAt).getTime();
    if (diff < 0) return '--';
    const days = Math.floor(diff / 86400000);
    const hours = Math.floor((diff % 86400000) / 3600000);
    const minutes = Math.floor((diff % 3600000) / 60000);
    if (days > 0) return `${days}d ${hours}h ${minutes}m`;
    return `${hours}h ${minutes}m`;
}

export function ContainerModal({ container, onClose }: ContainerModalProps) {
    const isRunning = container.state === 'running';
    const stateColor = isRunning ? '#00ff88' : '#ff3366';
    const name = container.service?.replace(/^font-/, 'front-') || container.name;

    // Close on Escape
    useEffect(() => {
        const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose(); };
        window.addEventListener('keydown', handler);
        return () => window.removeEventListener('keydown', handler);
    }, [onClose]);

    // Prevent body scroll
    useEffect(() => {
        document.body.style.overflow = 'hidden';
        return () => { document.body.style.overflow = ''; };
    }, []);

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            {/* Backdrop */}
            <div
                className="absolute inset-0"
                style={{ background: 'rgba(5,5,10,0.85)', backdropFilter: 'blur(4px)' }}
                onClick={onClose}
            />

            {/* Modal */}
            <div
                className="relative w-full max-w-4xl max-h-[90vh] flex flex-col rounded-xl overflow-hidden"
                style={{ background: '#0e0e16', border: '1px solid #1e1e2e' }}
            >
                {/* Header */}
                <div className="flex items-center justify-between px-5 py-3 shrink-0"
                    style={{ borderBottom: '1px solid #1e1e2e' }}>
                    <div className="flex items-center gap-3">
                        <h2 className="text-sm font-semibold" style={{ color: '#e4e4ef' }}>{name}</h2>
                        <span className="text-[10px] px-2 py-0.5 rounded-md font-medium"
                            style={{ background: `${stateColor}15`, color: stateColor }}>
                            {container.state}
                        </span>
                        {isRunning && (
                            <span className="text-[10px] font-mono" style={{ color: '#55556a' }}>
                                uptime {formatUptime(container.started_at)}
                            </span>
                        )}
                    </div>
                    <div className="flex items-center gap-3">
                        <ActionButtons containerId={container.id} state={container.state} />
                        <button
                            onClick={onClose}
                            className="text-xs px-2 py-1 rounded-md cursor-pointer transition-colors"
                            style={{ color: '#8888a0', border: '1px solid #1e1e2e' }}
                            onMouseEnter={e => { e.currentTarget.style.borderColor = '#ff336630'; e.currentTarget.style.color = '#ff3366'; }}
                            onMouseLeave={e => { e.currentTarget.style.borderColor = '#1e1e2e'; e.currentTarget.style.color = '#8888a0'; }}
                        >
                            ESC
                        </button>
                    </div>
                </div>

                {/* Info bar */}
                <div className="flex items-center gap-4 px-5 py-2 text-[10px] shrink-0"
                    style={{ borderBottom: '1px solid #1e1e2e', color: '#55556a' }}>
                    <span className="font-mono">{container.image?.split(':')[0].split('/').pop()}</span>
                    {container.ports?.filter(p => p.host_port > 0).length > 0 && (
                        <span className="font-mono">
                            ports: {container.ports.filter(p => p.host_port > 0).map(p => `${p.host_port}:${p.container_port}`).join(', ')}
                        </span>
                    )}
                    <span className="font-mono">id: {container.id.slice(0, 12)}</span>
                </div>

                {/* Stats */}
                {isRunning && (
                    <div className="px-5 py-2 shrink-0" style={{ borderBottom: '1px solid #1e1e2e' }}>
                        <StatsBar containerId={container.id} state={container.state} />
                    </div>
                )}

                {/* Logs */}
                <div className="flex-1 min-h-0">
                    <LogViewer containerId={container.id} />
                </div>
            </div>
        </div>
    );
}
