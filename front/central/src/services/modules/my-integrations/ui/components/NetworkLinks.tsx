'use client';

import { useEffect, useState, type CSSProperties, type RefObject } from 'react';
import { useSyncActivity, type SyncNodeState } from '../sync-activity-context';

export interface NetworkTarget {
    key: string;
    el: HTMLElement;
    dir: 'in' | 'out';
    color: string;
    nodeId?: number;
    lane?: number;
    laneCount?: number;
}

interface NetworkLinksProps {
    container: RefObject<HTMLDivElement | null>;
    hub: RefObject<HTMLDivElement | null>;
    getTargets: () => NetworkTarget[];
    revision: number;
}

interface LinkPath {
    key: string;
    d: string;
    dRev: string;
    color: string;
    nodeId?: number;
}

const HOT_COLOR: Record<string, string> = {
    active: '#3b82f6',
    scan: '#8b5cf6',
};

const PULSES_PER_LINE = 3;

function pulseStyle(d: string, color: string, index: number, count: number, duration: number): CSSProperties {
    return {
        position: 'absolute',
        top: 0,
        left: 0,
        width: '9px',
        height: '9px',
        borderRadius: '50%',
        background: `radial-gradient(circle, #fff 20%, ${color} 60%)`,
        boxShadow: `0 0 12px 3px ${color}99`,
        offsetPath: `path('${d}')`,
        offsetRotate: '0deg',
        animation: `cyber-travel ${duration}s linear infinite`,
        animationDelay: `-${((index * duration) / count).toFixed(2)}s`,
        pointerEvents: 'none',
    };
}

export function NetworkLinks({ container, hub, getTargets, revision }: NetworkLinksProps) {
    const [links, setLinks] = useState<LinkPath[]>([]);
    const [size, setSize] = useState({ w: 0, h: 0 });
    const { nodes } = useSyncActivity();

    useEffect(() => {
        const compute = () => {
            const c = container.current;
            const h = hub.current;
            if (!c || !h) return;
            const cRect = c.getBoundingClientRect();
            const hRect = h.getBoundingClientRect();
            if (cRect.width === 0 || hRect.width === 0) return;

            const hx = hRect.left + hRect.width / 2 - cRect.left;
            const hy = hRect.top + hRect.height / 2 - cRect.top;
            const hr = hRect.width / 2 + 8;
            const f = (n: number) => n.toFixed(1);
            const next: LinkPath[] = [];

            const all = getTargets();
            const laneOffset = new Map<string, number>();
            const channelTargets = all.filter(t => t.nodeId !== undefined);
            const sideOf = (t: NetworkTarget) => {
                const b = t.el.getBoundingClientRect();
                const tcx = b.left - cRect.left + b.width / 2;
                return tcx < hx - 60 ? 'left' : tcx > hx + 60 ? 'right' : 'center';
            };
            const topOf = (t: NetworkTarget) => t.el.getBoundingClientRect().top;
            (['left', 'right'] as const).forEach(sideKey => {
                const group = channelTargets
                    .filter(t => sideOf(t) === sideKey)
                    .sort((a, b) => topOf(a) - topOf(b));
                const dir = sideKey === 'left' ? -1 : 1;
                group.forEach((t, i) => {
                    laneOffset.set(t.key, dir * (i + 1) * 9);
                });
            });
            const centered = channelTargets.filter(t => sideOf(t) === 'center');
            centered.forEach((t, i) => {
                laneOffset.set(t.key, (i - (centered.length - 1) / 2) * 9);
            });

            for (const t of all) {
                const r = t.el.getBoundingClientRect();
                const cx = r.left + r.width / 2 - cRect.left;
                const cy = r.top + r.height / 2 - cRect.top;
                const hw = r.width / 2;
                const hh = r.height / 2;

                if (t.nodeId !== undefined && t.lane !== undefined && t.laneCount !== undefined) {
                    const lanes = t.laneCount;
                    const laneX = hx + (laneOffset.get(t.key) ?? 0);
                    const hubTop = hy - hr;
                    const side = cx < hx - 60 ? 1 : cx > hx + 60 ? -1 : 0;

                    if (side !== 0) {
                        const ax = cx + side * hw;
                        const ay = cy;
                        const offset = laneX - hx;
                        const sy = hy - Math.sqrt(Math.max(0, hr * hr - offset * offset));
                        const rc = Math.min(16, Math.max(2, sy - ay - 4));
                        const dirX = ax > laneX ? 1 : -1;
                        const d = `M ${f(laneX)} ${f(sy)} L ${f(laneX)} ${f(ay + rc)} Q ${f(laneX)} ${f(ay)} ${f(laneX + dirX * rc)} ${f(ay)} L ${f(ax)} ${f(ay)}`;
                        const dRev = `M ${f(ax)} ${f(ay)} L ${f(laneX + dirX * rc)} ${f(ay)} Q ${f(laneX)} ${f(ay)} ${f(laneX)} ${f(ay + rc)} L ${f(laneX)} ${f(sy)}`;
                        next.push({ key: t.key, d, dRev, color: t.color, nodeId: t.nodeId });
                        continue;
                    }

                    const cardBottom = r.bottom - cRect.top;
                    const gap = Math.max(24, hubTop - cardBottom);
                    const busY = cardBottom + gap * (0.35 + (t.lane / Math.max(1, lanes)) * 0.4);
                    const ax = cx;
                    const rc = Math.min(14, Math.max(4, Math.abs(busY - cardBottom) / 2), Math.abs(ax - laneX) / 2 || 14);
                    const dirX = ax >= laneX ? 1 : -1;
                    const d = [
                        `M ${f(laneX)} ${f(hubTop)}`,
                        `L ${f(laneX)} ${f(busY + rc)}`,
                        `Q ${f(laneX)} ${f(busY)} ${f(laneX + dirX * rc)} ${f(busY)}`,
                        `L ${f(ax - dirX * rc)} ${f(busY)}`,
                        `Q ${f(ax)} ${f(busY)} ${f(ax)} ${f(busY - rc)}`,
                        `L ${f(ax)} ${f(cardBottom)}`,
                    ].join(' ');
                    const dRev = [
                        `M ${f(ax)} ${f(cardBottom)}`,
                        `L ${f(ax)} ${f(busY - rc)}`,
                        `Q ${f(ax)} ${f(busY)} ${f(ax - dirX * rc)} ${f(busY)}`,
                        `L ${f(laneX + dirX * rc)} ${f(busY)}`,
                        `Q ${f(laneX)} ${f(busY)} ${f(laneX)} ${f(busY + rc)}`,
                        `L ${f(laneX)} ${f(hubTop)}`,
                    ].join(' ');
                    next.push({ key: t.key, d, dRev, color: t.color, nodeId: t.nodeId });
                    continue;
                }

                const dxEdge = hx - cx || 1e-6;
                const dyEdge = hy - cy || 1e-6;
                const tt = Math.min(hw / Math.abs(dxEdge), hh / Math.abs(dyEdge));
                const ax = cx + dxEdge * tt;
                const ay = cy + dyEdge * tt;
                const dx = ax - hx;
                const dy = ay - hy;
                const dist = Math.hypot(dx, dy) || 1;
                const sx = hx + (dx / dist) * hr;
                const sy = hy + (dy / dist) * hr;
                const mx = (sx + ax) / 2 - dy * 0.13;
                const my = (sy + ay) / 2 + dx * 0.13;
                next.push({
                    key: t.key,
                    d: `M ${f(sx)} ${f(sy)} Q ${f(mx)} ${f(my)} ${f(ax)} ${f(ay)}`,
                    dRev: `M ${f(ax)} ${f(ay)} Q ${f(mx)} ${f(my)} ${f(sx)} ${f(sy)}`,
                    color: t.color,
                });
            }

            setSize({ w: cRect.width, h: cRect.height });
            setLinks(next);
        };

        compute();
        const c = container.current;
        if (!c) return;
        const ro = new ResizeObserver(() => compute());
        ro.observe(c);
        const raf = window.setTimeout(compute, 400);
        window.addEventListener('resize', compute);
        return () => {
            ro.disconnect();
            window.clearTimeout(raf);
            window.removeEventListener('resize', compute);
        };
    }, [container, hub, getTargets, revision]);

    if (links.length === 0 || size.w === 0) return null;

    const stateOf = (link: LinkPath): SyncNodeState =>
        (link.nodeId !== undefined ? nodes[link.nodeId] : undefined) || 'idle';

    return (
        <>
            <svg
                className="pointer-events-none absolute inset-0 z-20"
                width={size.w}
                height={size.h}
                viewBox={`0 0 ${size.w} ${size.h}`}
                fill="none"
                style={{ overflow: 'visible' }}
            >
                {links.map(link => {
                    const state = stateOf(link);
                    const hot = state === 'active' || state === 'scan';
                    return (
                        <g key={link.key}>
                            <path
                                d={link.d}
                                stroke={state === 'done' ? '#22c55e' : link.color}
                                strokeOpacity={state === 'done' ? 0.7 : state === 'queued' ? 0.18 : 0.3}
                                strokeWidth="1.5"
                                strokeDasharray="2 7"
                                strokeLinecap="round"
                            />
                            {hot && (
                                <path
                                    d={link.d}
                                    stroke={HOT_COLOR[state]}
                                    strokeWidth="2"
                                    strokeDasharray="7 11"
                                    strokeLinecap="round"
                                    opacity="0.85"
                                    style={{ animation: 'cyber-dashflow .5s linear infinite' }}
                                />
                            )}
                            {!hot && state !== 'done' && (
                                <circle r="2.5" fill={link.color} fillOpacity="0.85">
                                    <animateMotion dur="3s" repeatCount="indefinite" path={link.d} />
                                </circle>
                            )}
                        </g>
                    );
                })}
            </svg>

            <div className="pointer-events-none absolute inset-0 z-20">
                {links.flatMap(link => {
                    const state = stateOf(link);
                    if (state === 'active') {
                        return Array.from({ length: PULSES_PER_LINE }, (_, i) => (
                            <div
                                key={`${link.key}-out-${i}`}
                                style={pulseStyle(link.d, '#3b82f6', i, PULSES_PER_LINE, 1.4)}
                            />
                        ));
                    }
                    if (state === 'scan') {
                        const back = Array.from({ length: PULSES_PER_LINE }, (_, i) => (
                            <div
                                key={`${link.key}-in-${i}`}
                                style={pulseStyle(link.dRev, '#8b5cf6', i, PULSES_PER_LINE, 1.6)}
                            />
                        ));
                        return [
                            ...back,
                            <div key={`${link.key}-probe`} style={pulseStyle(link.d, '#22d3ee', 0, 1, 2.2)} />,
                        ];
                    }
                    return [];
                })}
            </div>
        </>
    );
}
