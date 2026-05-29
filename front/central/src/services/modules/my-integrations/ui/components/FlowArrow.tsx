'use client';

const PIPE_GRADIENT_ID = 'flow-pipe-gradient';
const PULSE_GRADIENT_ID = 'flow-pulse-gradient';
const GLOW_FILTER_ID = 'flow-pulse-glow';

function FlowDefs() {
    return (
        <defs>
            <linearGradient id={PIPE_GRADIENT_ID} x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" stopColor="#6366F1" stopOpacity="0.55" />
                <stop offset="50%" stopColor="#A855F7" stopOpacity="0.45" />
                <stop offset="100%" stopColor="#EC4899" stopOpacity="0.55" />
            </linearGradient>
            <linearGradient id={PULSE_GRADIENT_ID} x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" stopColor="#22D3EE" stopOpacity="0" />
                <stop offset="50%" stopColor="#22D3EE" stopOpacity="1" />
                <stop offset="100%" stopColor="#F472B6" stopOpacity="0" />
            </linearGradient>
            <filter id={GLOW_FILTER_ID} x="-30%" y="-30%" width="160%" height="160%">
                <feGaussianBlur stdDeviation="3" result="blur" />
                <feMerge>
                    <feMergeNode in="blur" />
                    <feMergeNode in="SourceGraphic" />
                </feMerge>
            </filter>
        </defs>
    );
}

function curvePath(sourceX: number, sourceY: number, targetX: number, targetY: number) {
    const dx = targetX - sourceX;
    const dy = targetY - sourceY;
    const midY = sourceY + dy / 2;
    const sway = Math.min(Math.abs(dx) * 0.35 + 40, 90);
    const controlOffset = dx >= 0 ? sway : -sway;
    return `M${sourceX} ${sourceY} C${sourceX + controlOffset * 0.2} ${midY - 8}, ${targetX - controlOffset * 0.2} ${midY + 8}, ${targetX} ${targetY}`;
}

export function FlowConverge({ count }: { count: number }) {
    if (count <= 0) return null;

    const width = 600;
    const height = 80;
    const centerX = width / 2;
    const colWidth = width / count;
    const targetY = height - 4;

    return (
        <div className="flex items-center justify-center py-1 w-full">
            <svg
                width="100%"
                height={height}
                viewBox={`0 0 ${width} ${height}`}
                fill="none"
                className="max-w-full"
                preserveAspectRatio="xMidYMid meet"
            >
                <FlowDefs />
                {Array.from({ length: count }).map((_, i) => {
                    const sourceX = colWidth * i + colWidth / 2;
                    const d = curvePath(sourceX, 0, centerX, targetY);
                    return (
                        <g key={i}>
                            <path
                                d={d}
                                stroke="url(#flow-pipe-gradient)"
                                strokeWidth="8"
                                strokeLinecap="round"
                                strokeOpacity="0.18"
                                fill="none"
                            />
                            <path
                                d={d}
                                stroke="url(#flow-pipe-gradient)"
                                strokeWidth="3.5"
                                strokeLinecap="round"
                                strokeOpacity="0.6"
                                fill="none"
                            />
                            <path
                                d={d}
                                stroke="url(#flow-pulse-gradient)"
                                strokeWidth="3"
                                strokeLinecap="round"
                                strokeDasharray="22 200"
                                filter={`url(#${GLOW_FILTER_ID})`}
                                fill="none"
                            >
                                <animate
                                    attributeName="stroke-dashoffset"
                                    from="222"
                                    to="0"
                                    dur="2.4s"
                                    repeatCount="indefinite"
                                    begin={`${i * 0.35}s`}
                                />
                            </path>
                        </g>
                    );
                })}
            </svg>
        </div>
    );
}

export function FlowDiverge({ count }: { count: number }) {
    if (count <= 0) return null;

    const width = 600;
    const height = 80;
    const centerX = width / 2;
    const colWidth = width / count;
    const endY = height - 4;

    return (
        <div className="flex items-center justify-center py-1 w-full">
            <svg
                width="100%"
                height={height}
                viewBox={`0 0 ${width} ${height}`}
                fill="none"
                className="max-w-full"
                preserveAspectRatio="xMidYMid meet"
            >
                <FlowDefs />
                {Array.from({ length: count }).map((_, i) => {
                    const targetX = colWidth * i + colWidth / 2;
                    const d = curvePath(centerX, 0, targetX, endY);
                    return (
                        <g key={i}>
                            <path
                                d={d}
                                stroke="url(#flow-pipe-gradient)"
                                strokeWidth="8"
                                strokeLinecap="round"
                                strokeOpacity="0.18"
                                fill="none"
                            />
                            <path
                                d={d}
                                stroke="url(#flow-pipe-gradient)"
                                strokeWidth="3.5"
                                strokeLinecap="round"
                                strokeOpacity="0.6"
                                fill="none"
                            />
                            <path
                                d={d}
                                stroke="url(#flow-pulse-gradient)"
                                strokeWidth="3"
                                strokeLinecap="round"
                                strokeDasharray="22 200"
                                filter={`url(#${GLOW_FILTER_ID})`}
                                fill="none"
                            >
                                <animate
                                    attributeName="stroke-dashoffset"
                                    from="222"
                                    to="0"
                                    dur="2.4s"
                                    repeatCount="indefinite"
                                    begin={`${i * 0.28}s`}
                                />
                            </path>
                        </g>
                    );
                })}
            </svg>
        </div>
    );
}
