'use client';

import type { Container } from '../../domain/types';
import Link from 'next/link';

interface ArchitectureViewProps {
    containers: Container[];
}

// Normalize font- → front-
function norm(service: string): string {
    return service.replace(/^font-/, 'front-');
}

function findContainer(containers: Container[], name: string): Container | undefined {
    return containers.find(c => norm(c.service) === name);
}

// Compact service node
function ServiceNode({ container, color }: { container: Container | undefined; color: string }) {
    if (!container) return <div className="w-full" />;

    const isRunning = container.state === 'running';
    const stateColor = isRunning ? '#00ff88' : '#ff3366';
    const name = norm(container.service);

    return (
        <Link
            href={`/dashboard/${container.id}`}
            className="card-hover block rounded-lg px-3 py-2.5 relative overflow-hidden group"
            style={{ background: '#12121a', border: `1px solid ${color}18` }}
            onMouseEnter={e => {
                e.currentTarget.style.borderColor = `${color}40`;
                e.currentTarget.style.boxShadow = `0 0 16px ${color}10`;
            }}
            onMouseLeave={e => {
                e.currentTarget.style.borderColor = `${color}18`;
                e.currentTarget.style.boxShadow = 'none';
            }}
        >
            <div className="absolute top-0 left-0 right-0 h-px opacity-0 group-hover:opacity-100 transition-opacity"
                style={{ background: `linear-gradient(90deg, transparent, ${color}50, transparent)` }} />
            <div className="flex items-center justify-between gap-2">
                <div className="min-w-0">
                    <div className="text-xs font-medium truncate" style={{ color: '#e4e4ef' }}>{name}</div>
                    <div className="text-[9px] font-mono mt-0.5 truncate" style={{ color: '#55556a' }}>
                        {container.ports?.filter(p => p.host_port > 0).map(p => `${p.host_port}:${p.container_port}`).join(', ') || 'internal'}
                    </div>
                </div>
                <div className="flex items-center gap-1 shrink-0">
                    <div
                        className={isRunning ? 'pulse-dot' : ''}
                        style={{
                            width: 5, height: 5, borderRadius: '50%',
                            background: stateColor,
                            boxShadow: isRunning ? `0 0 6px ${stateColor}` : 'none',
                        }}
                    />
                    <span className="text-[8px] uppercase tracking-wider font-medium" style={{ color: stateColor }}>
                        {isRunning ? 'UP' : 'DOWN'}
                    </span>
                </div>
            </div>
        </Link>
    );
}

// Vertical animated arrow between two layers
function VerticalArrow({ color, label }: { color: string; label: string }) {
    return (
        <div className="flex flex-col items-center py-0.5" style={{ minHeight: 32 }}>
            <div className="relative w-px flex-1 overflow-hidden" style={{ minHeight: 14, background: `${color}20` }}>
                <div
                    className="absolute w-px h-2 rounded-full"
                    style={{
                        background: color,
                        boxShadow: `0 0 4px ${color}`,
                        animation: 'flow-down 1.4s ease-in-out infinite',
                    }}
                />
            </div>
            <span className="text-[7px] font-mono py-0.5 whitespace-nowrap" style={{ color: `${color}50` }}>{label}</span>
            <div className="relative w-px flex-1 overflow-hidden" style={{ minHeight: 14, background: `${color}20` }}>
                <div
                    className="absolute w-px h-2 rounded-full"
                    style={{
                        background: color,
                        boxShadow: `0 0 4px ${color}`,
                        animation: 'flow-down 1.4s ease-in-out infinite',
                        animationDelay: '0.7s',
                    }}
                />
            </div>
            {/* Arrow tip */}
            <div style={{
                width: 0, height: 0,
                borderLeft: '3px solid transparent',
                borderRight: '3px solid transparent',
                borderTop: `4px solid ${color}40`,
            }} />
        </div>
    );
}

// Empty spacer matching arrow width
function ArrowSpacer() {
    return <div style={{ minHeight: 32 }} />;
}

export function ArchitectureView({ containers }: ArchitectureViewProps) {
    const fc = findContainer(containers, 'front-central');
    const fw = findContainer(containers, 'front-website');
    const ft = findContainer(containers, 'front-testing');
    const mw = findContainer(containers, 'monitoring-web');
    const bc = findContainer(containers, 'back-central');
    const bt = findContainer(containers, 'back-testing');
    const ma = findContainer(containers, 'monitoring-api');
    const nx = findContainer(containers, 'nginx');
    const rd = findContainer(containers, 'redis');
    const rq = findContainer(containers, 'rabbitmq');

    return (
        <div>
            {/* === FRONTEND ROW === */}
            <div className="flex items-center gap-2 mb-2">
                <div className="w-1 h-4 rounded-full" style={{ background: '#00f0ff' }} />
                <span className="text-[11px] font-semibold uppercase tracking-wider" style={{ color: '#00f0ff' }}>Frontend</span>
            </div>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-2">
                <ServiceNode container={fc} color="#00f0ff" />
                <ServiceNode container={fw} color="#00ff88" />
                <ServiceNode container={ft} color="#a855f7" />
                <ServiceNode container={mw} color="#8888a0" />
            </div>

            {/* === ARROWS: Frontend → Backend === */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-2">
                {/* front-central → back-central */}
                <VerticalArrow color="#00f0ff" label="HTTP :3050" />
                {/* front-website = standalone */}
                <ArrowSpacer />
                {/* front-testing → back-testing */}
                <VerticalArrow color="#a855f7" label="HTTP :9092" />
                {/* monitoring-web → monitoring-api */}
                <VerticalArrow color="#8888a0" label="HTTP :3070" />
            </div>

            {/* === BACKEND ROW === */}
            <div className="flex items-center gap-2 mb-2">
                <div className="w-1 h-4 rounded-full" style={{ background: '#a855f7' }} />
                <span className="text-[11px] font-semibold uppercase tracking-wider" style={{ color: '#a855f7' }}>Backend</span>
            </div>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-2">
                <ServiceNode container={bc} color="#a855f7" />
                <ServiceNode container={nx} color="#ff6b6b" />
                <ServiceNode container={bt} color="#a855f7" />
                <ServiceNode container={ma} color="#8888a0" />
            </div>

            {/* === ARROWS: Backend → Data (only from back-central) === */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-2">
                {/* back-central → redis + rabbitmq */}
                <div className="flex justify-center">
                    <div className="flex gap-6">
                        <VerticalArrow color="#ffaa00" label="TCP :6379" />
                        <VerticalArrow color="#ffaa00" label="AMQP :5672" />
                    </div>
                </div>
                <ArrowSpacer />
                <ArrowSpacer />
                <ArrowSpacer />
            </div>

            {/* === DATA ROW === */}
            <div className="flex items-center gap-2 mb-2">
                <div className="w-1 h-4 rounded-full" style={{ background: '#ffaa00' }} />
                <span className="text-[11px] font-semibold uppercase tracking-wider" style={{ color: '#ffaa00' }}>Data & Messaging</span>
            </div>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-2">
                <ServiceNode container={rd} color="#ffaa00" />
                <ServiceNode container={rq} color="#ffaa00" />
            </div>
        </div>
    );
}
