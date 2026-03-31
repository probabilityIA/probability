'use client';

import { useRef, useEffect, useState, useCallback } from 'react';
import type { Container } from '../../domain/types';
import Link from 'next/link';

interface ArchitectureViewProps {
    containers: Container[];
}

function norm(service: string): string {
    return service.replace(/^font-/, 'front-');
}

function find(containers: Container[], name: string): Container | undefined {
    return containers.find(c => norm(c.service) === name);
}

// ── Service Node ──
function ServiceNode({
    container,
    color,
    id,
}: {
    container: Container | undefined;
    color: string;
    id: string;
}) {
    if (!container) return null;
    const isRunning = container.state === 'running';
    const stateColor = isRunning ? '#00ff88' : '#ff3366';

    return (
        <Link
            href={`/dashboard/${container.id}`}
            id={id}
            className="card-hover block rounded-lg px-3 py-2 relative overflow-hidden group"
            style={{ background: '#12121a', border: `1px solid ${color}20` }}
            onMouseEnter={e => { e.currentTarget.style.borderColor = `${color}45`; e.currentTarget.style.boxShadow = `0 0 20px ${color}12`; }}
            onMouseLeave={e => { e.currentTarget.style.borderColor = `${color}20`; e.currentTarget.style.boxShadow = 'none'; }}
        >
            <div className="absolute top-0 left-0 right-0 h-px opacity-0 group-hover:opacity-100 transition-opacity"
                style={{ background: `linear-gradient(90deg, transparent, ${color}50, transparent)` }} />
            <div className="flex items-center justify-between gap-2">
                <div className="min-w-0">
                    <div className="text-[11px] font-medium truncate" style={{ color: '#e4e4ef' }}>
                        {norm(container.service)}
                    </div>
                    <div className="text-[8px] font-mono mt-0.5 truncate" style={{ color: '#55556a' }}>
                        {container.ports?.filter(p => p.host_port > 0).map(p => `${p.host_port}:${p.container_port}`).join(', ') || 'internal'}
                    </div>
                </div>
                <div className="flex items-center gap-1 shrink-0">
                    <div className={isRunning ? 'pulse-dot' : ''} style={{
                        width: 5, height: 5, borderRadius: '50%',
                        background: stateColor,
                        boxShadow: isRunning ? `0 0 6px ${stateColor}` : 'none',
                    }} />
                </div>
            </div>
        </Link>
    );
}

// ── Connection definition ──
interface Connection {
    from: string;
    to: string;
    label: string;
    color: string;
}

const CONNECTIONS: Connection[] = [
    { from: 'node-fc', to: 'node-bc', label: 'HTTP', color: '#00f0ff' },
    { from: 'node-ft', to: 'node-bt', label: 'HTTP', color: '#a855f7' },
    { from: 'node-mw', to: 'node-ma', label: 'HTTP', color: '#8888a0' },
    { from: 'node-bt', to: 'node-bc', label: 'testea APIs', color: '#ff6b6b' },
    { from: 'node-bc', to: 'node-rd', label: 'TCP', color: '#ffaa00' },
    { from: 'node-bc', to: 'node-rq', label: 'AMQP', color: '#ffaa00' },
];

// ── SVG Arrow drawing ──
function SvgConnections({ connections, containerRef }: { connections: Connection[]; containerRef: React.RefObject<HTMLDivElement | null> }) {
    const [lines, setLines] = useState<{ x1: number; y1: number; x2: number; y2: number; label: string; color: string }[]>([]);

    const calculate = useCallback(() => {
        if (!containerRef.current) return;
        const box = containerRef.current.getBoundingClientRect();
        const newLines: { x1: number; y1: number; x2: number; y2: number; label: string; color: string }[] = [];

        for (const conn of connections) {
            const fromEl = document.getElementById(conn.from);
            const toEl = document.getElementById(conn.to);
            if (!fromEl || !toEl) continue;

            const fromRect = fromEl.getBoundingClientRect();
            const toRect = toEl.getBoundingClientRect();

            // Calculate center-bottom of from, center-top of to (or sides for horizontal)
            const fromCx = fromRect.left + fromRect.width / 2 - box.left;
            const fromCy = fromRect.top + fromRect.height / 2 - box.top;
            const toCx = toRect.left + toRect.width / 2 - box.left;
            const toCy = toRect.top + toRect.height / 2 - box.top;

            // Determine best connection points
            let x1: number, y1: number, x2: number, y2: number;

            const dy = Math.abs(toCy - fromCy);
            const dx = Math.abs(toCx - fromCx);

            if (dy > dx * 0.4) {
                // Mostly vertical
                x1 = fromCx;
                y1 = fromCy > toCy ? fromRect.top - box.top : fromRect.bottom - box.top;
                x2 = toCx;
                y2 = toCy > fromCy ? toRect.top - box.top : toRect.bottom - box.top;
            } else {
                // Mostly horizontal
                x1 = fromCx > toCx ? fromRect.left - box.left : fromRect.right - box.left;
                y1 = fromCy;
                x2 = toCx > fromCx ? toRect.left - box.left : toRect.right - box.left;
                y2 = toCy;
            }

            newLines.push({ x1, y1, x2, y2, label: conn.label, color: conn.color });
        }

        setLines(newLines);
    }, [connections, containerRef]);

    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        setMounted(true);
    }, []);

    useEffect(() => {
        if (!mounted) return;
        // Delay to ensure DOM nodes are rendered
        const t1 = setTimeout(calculate, 100);
        const t2 = setTimeout(calculate, 500);
        window.addEventListener('resize', calculate);
        return () => { window.removeEventListener('resize', calculate); clearTimeout(t1); clearTimeout(t2); };
    }, [calculate, mounted]);

    if (!mounted || lines.length === 0) return null;

    return (
        <svg className="absolute inset-0 w-full h-full pointer-events-none" style={{ zIndex: 1 }}>
            <defs>
                {lines.map((_, i) => (
                    <marker key={`arrow-${i}`} id={`arrowhead-${i}`} markerWidth="6" markerHeight="4" refX="5" refY="2" orient="auto">
                        <polygon points="0 0, 6 2, 0 4" fill={lines[i].color} opacity="0.6" />
                    </marker>
                ))}
            </defs>

            {lines.map((line, i) => {
                const mx = (line.x1 + line.x2) / 2;
                const my = (line.y1 + line.y2) / 2;

                return (
                    <g key={i}>
                        {/* Line */}
                        <line
                            x1={line.x1} y1={line.y1} x2={line.x2} y2={line.y2}
                            stroke={line.color}
                            strokeWidth="1"
                            opacity="0.25"
                            markerEnd={`url(#arrowhead-${i})`}
                        />
                        {/* Animated dot */}
                        <circle r="2" fill={line.color} opacity="0.8">
                            <animate
                                attributeName="cx"
                                from={line.x1} to={line.x2}
                                dur="2s" repeatCount="indefinite"
                            />
                            <animate
                                attributeName="cy"
                                from={line.y1} to={line.y2}
                                dur="2s" repeatCount="indefinite"
                            />
                            <animate
                                attributeName="opacity"
                                values="0;0.9;0.9;0" dur="2s" repeatCount="indefinite"
                            />
                        </circle>
                        {/* Glow on dot */}
                        <circle r="4" fill={line.color} opacity="0">
                            <animate attributeName="cx" from={line.x1} to={line.x2} dur="2s" repeatCount="indefinite" />
                            <animate attributeName="cy" from={line.y1} to={line.y2} dur="2s" repeatCount="indefinite" />
                            <animate attributeName="opacity" values="0;0.3;0.3;0" dur="2s" repeatCount="indefinite" />
                        </circle>
                        {/* Label */}
                        <text x={mx} y={my - 6} textAnchor="middle" fill={line.color} opacity="0.4"
                            style={{ fontSize: '8px', fontFamily: 'monospace' }}>
                            {line.label}
                        </text>
                    </g>
                );
            })}
        </svg>
    );
}

// ── Main diagram ──
export function ArchitectureView({ containers }: ArchitectureViewProps) {
    const containerRef = useRef<HTMLDivElement>(null);

    const fc = find(containers, 'front-central');
    const fw = find(containers, 'front-website');
    const ft = find(containers, 'front-testing');
    const mw = find(containers, 'monitoring-web');
    const bc = find(containers, 'back-central');
    const bt = find(containers, 'back-testing');
    const ma = find(containers, 'monitoring-api');
    const nx = find(containers, 'nginx');
    const rd = find(containers, 'redis');
    const rq = find(containers, 'rabbitmq');

    return (
        <div ref={containerRef} className="relative" style={{ minHeight: 420 }}>
            {/* SVG overlay for arrows */}
            <SvgConnections connections={CONNECTIONS} containerRef={containerRef} />

            {/* Nodes grid - positioned to match the concept map */}
            <div className="relative" style={{ zIndex: 2 }}>

                {/* Row 1: Frontends */}
                <div className="grid grid-cols-4 gap-x-4 gap-y-0 mb-2">
                    <div className="col-span-4 flex items-center gap-2 mb-1">
                        <div className="w-1 h-3 rounded-full" style={{ background: '#00f0ff' }} />
                        <span className="text-[10px] font-semibold uppercase tracking-wider" style={{ color: '#00f0ff' }}>Frontend</span>
                    </div>
                </div>
                <div className="grid grid-cols-4 gap-3">
                    <ServiceNode container={mw} color="#8888a0" id="node-mw" />
                    <ServiceNode container={fc} color="#00f0ff" id="node-fc" />
                    <ServiceNode container={fw} color="#00ff88" id="node-fw" />
                    <ServiceNode container={ft} color="#a855f7" id="node-ft" />
                </div>

                {/* Spacer for arrows */}
                <div className="h-16" />

                {/* Row 2: Backends */}
                <div className="grid grid-cols-4 gap-x-4 gap-y-0 mb-2">
                    <div className="col-span-4 flex items-center gap-2 mb-1">
                        <div className="w-1 h-3 rounded-full" style={{ background: '#a855f7' }} />
                        <span className="text-[10px] font-semibold uppercase tracking-wider" style={{ color: '#a855f7' }}>Backend</span>
                    </div>
                </div>
                <div className="grid grid-cols-4 gap-3">
                    <ServiceNode container={ma} color="#8888a0" id="node-ma" />
                    <ServiceNode container={bc} color="#00f0ff" id="node-bc" />
                    <ServiceNode container={nx} color="#ff6b6b" id="node-nx" />
                    <ServiceNode container={bt} color="#a855f7" id="node-bt" />
                </div>

                {/* Spacer for arrows */}
                <div className="h-16" />

                {/* Row 3: Data */}
                <div className="grid grid-cols-4 gap-x-4 gap-y-0 mb-2">
                    <div className="col-span-4 flex items-center gap-2 mb-1">
                        <div className="w-1 h-3 rounded-full" style={{ background: '#ffaa00' }} />
                        <span className="text-[10px] font-semibold uppercase tracking-wider" style={{ color: '#ffaa00' }}>Data & Messaging</span>
                    </div>
                </div>
                <div className="grid grid-cols-4 gap-3">
                    <div /> {/* empty col1 */}
                    <ServiceNode container={rq} color="#ffaa00" id="node-rq" />
                    <ServiceNode container={rd} color="#ffaa00" id="node-rd" />
                    <div /> {/* empty col4 */}
                </div>
            </div>
        </div>
    );
}
