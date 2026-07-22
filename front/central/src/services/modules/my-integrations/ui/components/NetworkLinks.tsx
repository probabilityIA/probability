'use client';

import { useEffect, useState, type RefObject } from 'react';

export interface NetworkTarget {
    key: string;
    el: HTMLElement;
    dir: 'in' | 'out';
    color: string;
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
    color: string;
}

export function NetworkLinks({ container, hub, getTargets, revision }: NetworkLinksProps) {
    const [links, setLinks] = useState<LinkPath[]>([]);
    const [size, setSize] = useState({ w: 0, h: 0 });

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
            const next: LinkPath[] = [];
            for (const t of getTargets()) {
                const r = t.el.getBoundingClientRect();
                const above = r.top < hRect.top;
                const sx = r.left + r.width / 2 - cRect.left;
                const sy = above ? r.bottom - cRect.top : r.top - cRect.top;
                const dx = hx - sx;
                const dy = hy - sy;
                const dist = Math.sqrt(dx * dx + dy * dy) || 1;
                const ex = hx - (dx / dist) * hr;
                const ey = hy - (dy / dist) * hr;
                const bend = Math.max(Math.abs(ey - sy) * 0.45, 30);
                const c1y = above ? sy + bend : sy - bend;
                const c2y = above ? ey - bend : ey + bend;
                const d = t.dir === 'in'
                    ? `M ${sx} ${sy} C ${sx} ${c1y}, ${ex} ${c2y}, ${ex} ${ey}`
                    : `M ${ex} ${ey} C ${ex} ${c2y}, ${sx} ${c1y}, ${sx} ${sy}`;
                next.push({ key: t.key, d, color: t.color });
            }
            setSize({ w: cRect.width, h: cRect.height });
            setLinks(next);
        };

        compute();
        const c = container.current;
        if (!c) return;
        const ro = new ResizeObserver(() => compute());
        ro.observe(c);
        window.addEventListener('resize', compute);
        return () => {
            ro.disconnect();
            window.removeEventListener('resize', compute);
        };
    }, [container, hub, getTargets, revision]);

    if (links.length === 0 || size.w === 0) return null;

    return (
        <svg
            className="pointer-events-none absolute inset-0 z-0"
            width={size.w}
            height={size.h}
            viewBox={`0 0 ${size.w} ${size.h}`}
            fill="none"
        >
            {links.map(link => (
                <g key={link.key}>
                    <path
                        d={link.d}
                        stroke={link.color}
                        strokeOpacity="0.3"
                        strokeWidth="1.5"
                        strokeDasharray="4 8"
                        style={{ animation: 'cyber-dash 1.6s linear infinite' }}
                    />
                    <circle r="2.5" fill={link.color} fillOpacity="0.85">
                        <animateMotion dur="3s" repeatCount="indefinite" path={link.d} />
                    </circle>
                </g>
            ))}
        </svg>
    );
}
