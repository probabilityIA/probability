'use client';

import type { Container } from '../../domain/types';
import Link from 'next/link';

interface ContainerCardProps {
    container: Container;
}

function getStateConfig(state: string) {
    switch (state) {
        case 'running': return { dot: '#00ff88', bg: '#00ff8810', border: '#00ff8830', label: 'Running', glowClass: 'glow-running' };
        case 'exited': return { dot: '#ff3366', bg: '#ff336610', border: '#ff336630', label: 'Stopped', glowClass: 'glow-stopped' };
        case 'restarting': return { dot: '#ffaa00', bg: '#ffaa0010', border: '#ffaa0030', label: 'Restarting', glowClass: 'glow-restarting' };
        case 'paused': return { dot: '#a855f7', bg: '#a855f710', border: '#a855f730', label: 'Paused', glowClass: '' };
        default: return { dot: '#8888a0', bg: '#8888a010', border: '#8888a030', label: state, glowClass: '' };
    }
}

function getHealthBadge(health: string) {
    if (!health || health === 'none' || health === '') return null;
    const colors: Record<string, { bg: string; text: string }> = {
        healthy: { bg: '#00ff8815', text: '#00ff88' },
        unhealthy: { bg: '#ff336615', text: '#ff3366' },
        starting: { bg: '#ffaa0015', text: '#ffaa00' },
    };
    const c = colors[health] || { bg: '#8888a015', text: '#8888a0' };
    return (
        <span className="text-[10px] px-1.5 py-0.5 rounded" style={{ background: c.bg, color: c.text }}>
            {health}
        </span>
    );
}

function formatUptime(startedAt: string): string {
    if (!startedAt) return '--';
    const diff = Date.now() - new Date(startedAt).getTime();
    if (diff < 0) return '--';
    const hours = Math.floor(diff / 3600000);
    const minutes = Math.floor((diff % 3600000) / 60000);
    if (hours >= 24) {
        const days = Math.floor(hours / 24);
        return `${days}d ${hours % 24}h`;
    }
    return `${hours}h ${minutes}m`;
}

export function ContainerCard({ container }: ContainerCardProps) {
    const state = getStateConfig(container.state);
    const serviceName = container.service || container.name.replace(/^\//, '').replace(/_/g, ' ');

    return (
        <Link
            href={`/dashboard/${container.id}`}
            className="card-hover block rounded-xl p-4 relative overflow-hidden group"
            style={{
                background: '#12121a',
                border: '1px solid #1e1e2e',
            }}
            onMouseEnter={e => {
                e.currentTarget.style.borderColor = state.border;
                e.currentTarget.style.boxShadow = `0 0 24px ${state.dot}10, 0 0 48px ${state.dot}05`;
            }}
            onMouseLeave={e => {
                e.currentTarget.style.borderColor = '#1e1e2e';
                e.currentTarget.style.boxShadow = 'none';
            }}
        >
            {/* Subtle top accent line */}
            <div
                className="absolute top-0 left-0 right-0 h-px opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                style={{ background: `linear-gradient(90deg, transparent, ${state.dot}60, transparent)` }}
            />

            {/* Header */}
            <div className="flex items-start justify-between mb-3">
                <div className="min-w-0 flex-1">
                    <h3 className="text-sm font-medium truncate" style={{ color: '#e4e4ef' }}>
                        {serviceName}
                    </h3>
                    <p className="text-[11px] truncate mt-0.5 font-mono" style={{ color: '#55556a' }}>
                        {container.image.split(':')[0].split('/').pop()}
                    </p>
                </div>
                <div className="flex items-center gap-1.5 ml-2 shrink-0">
                    {getHealthBadge(container.health)}
                    <div className="flex items-center gap-1.5 px-2 py-1 rounded-md" style={{ background: state.bg }}>
                        <div
                            className={`w-1.5 h-1.5 rounded-full ${container.state === 'running' ? 'pulse-dot' : ''}`}
                            style={{
                                background: state.dot,
                                boxShadow: `0 0 6px ${state.dot}`,
                            }}
                        />
                        <span className="text-[10px] font-medium uppercase tracking-wider" style={{ color: state.dot }}>
                            {state.label}
                        </span>
                    </div>
                </div>
            </div>

            {/* Info row */}
            <div className="flex items-center gap-4 text-[11px]" style={{ color: '#8888a0' }}>
                {container.state === 'running' && (
                    <div className="flex items-center gap-1">
                        <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span className="font-mono">{formatUptime(container.started_at)}</span>
                    </div>
                )}
                {container.ports?.length > 0 && (
                    <div className="flex items-center gap-1">
                        <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m9.86-2.06a4.5 4.5 0 00-1.242-7.244l4.5-4.5a4.5 4.5 0 016.364 6.364l-1.757 1.757" />
                        </svg>
                        <span className="font-mono">{container.ports.map(p => `${p.host_port}:${p.container_port}`).join(', ')}</span>
                    </div>
                )}
                <div className="flex items-center gap-1 ml-auto opacity-0 group-hover:opacity-100 transition-opacity" style={{ color: '#00f0ff' }}>
                    <span>View</span>
                    <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
                    </svg>
                </div>
            </div>
        </Link>
    );
}
