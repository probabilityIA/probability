'use client';

const GRADIENT_ID = 'flow-arrow-gradient';
const FILTER_ID = 'flow-arrow-glow';

function FlowDefs() {
    return (
        <defs>
            <linearGradient id={GRADIENT_ID} x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" stopColor="#8B5CF6" stopOpacity="0.15" />
                <stop offset="50%" stopColor="#8B5CF6" stopOpacity="1" />
                <stop offset="100%" stopColor="#8B5CF6" stopOpacity="0.15" />
            </linearGradient>
            <filter id={FILTER_ID} x="-20%" y="-20%" width="140%" height="140%">
                <feGaussianBlur stdDeviation="2.5" result="blur" />
                <feMerge>
                    <feMergeNode in="blur" />
                    <feMergeNode in="SourceGraphic" />
                </feMerge>
            </filter>
        </defs>
    );
}

export function FlowConverge({ count }: { count: number }) {
    if (count <= 0) return null;

    const width = 600;
    const height = 48;
    const centerX = width / 2;
    const colWidth = width / count;

    return (
        <div className="flex items-center justify-center py-1 w-full">
            <svg
                width="100%"
                height={height}
                viewBox={`0 0 ${width} ${height}`}
                fill="none"
                className="text-purple-400 dark:text-purple-500 max-w-full"
                preserveAspectRatio="xMidYMid meet"
            >
                <FlowDefs />
                {Array.from({ length: count }).map((_, i) => {
                    const sourceX = colWidth * i + colWidth / 2;
                    const pathId = `flow-converge-${i}`;
                    const d = `M${sourceX} 0 L${centerX} ${height - 10}`;
                    return (
                        <g key={i}>
                            <path
                                d={d}
                                stroke="currentColor"
                                strokeWidth="1.5"
                                strokeOpacity="0.25"
                                strokeLinecap="round"
                            />
                            <path
                                id={pathId}
                                d={d}
                                stroke={`url(#${GRADIENT_ID})`}
                                strokeWidth="2.5"
                                strokeLinecap="round"
                                strokeDasharray="14 86"
                                filter={`url(#${FILTER_ID})`}
                            >
                                <animate
                                    attributeName="stroke-dashoffset"
                                    from="100"
                                    to="0"
                                    dur="1.6s"
                                    repeatCount="indefinite"
                                    begin={`${i * 0.25}s`}
                                />
                            </path>
                        </g>
                    );
                })}
                <path
                    d={`M${centerX - 5} ${height - 18}L${centerX} ${height - 10}L${centerX + 5} ${height - 18}`}
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeOpacity="0.6"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    fill="none"
                />
            </svg>
        </div>
    );
}

export function FlowDiverge({ count }: { count: number }) {
    if (count <= 0) return null;

    const width = 600;
    const height = 48;
    const centerX = width / 2;
    const colWidth = width / count;

    return (
        <div className="flex items-center justify-center py-1 w-full">
            <svg
                width="100%"
                height={height}
                viewBox={`0 0 ${width} ${height}`}
                fill="none"
                className="text-purple-400 dark:text-purple-500 max-w-full"
                preserveAspectRatio="xMidYMid meet"
            >
                <FlowDefs />
                {Array.from({ length: count }).map((_, i) => {
                    const targetX = colWidth * i + colWidth / 2;
                    const endY = height - 10;
                    const angle = Math.atan2(endY, targetX - centerX);
                    const len = 7;
                    const spread = Math.PI / 6;
                    const lx = targetX - len * Math.cos(angle - spread);
                    const ly = endY - len * Math.sin(angle - spread);
                    const rx = targetX - len * Math.cos(angle + spread);
                    const ry = endY - len * Math.sin(angle + spread);
                    const d = `M${centerX} 0 L${targetX} ${endY}`;

                    return (
                        <g key={i}>
                            <path
                                d={d}
                                stroke="currentColor"
                                strokeWidth="1.5"
                                strokeOpacity="0.25"
                                strokeLinecap="round"
                            />
                            <path
                                d={d}
                                stroke={`url(#${GRADIENT_ID})`}
                                strokeWidth="2.5"
                                strokeLinecap="round"
                                strokeDasharray="14 86"
                                filter={`url(#${FILTER_ID})`}
                            >
                                <animate
                                    attributeName="stroke-dashoffset"
                                    from="100"
                                    to="0"
                                    dur="1.6s"
                                    repeatCount="indefinite"
                                    begin={`${i * 0.2}s`}
                                />
                            </path>
                            <path
                                d={`M${lx} ${ly}L${targetX} ${endY}L${rx} ${ry}`}
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeOpacity="0.6"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                            />
                        </g>
                    );
                })}
            </svg>
        </div>
    );
}
