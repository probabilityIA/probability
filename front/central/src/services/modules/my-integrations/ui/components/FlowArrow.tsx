'use client';

/**
 * Líneas convergentes: desde N puntos arriba convergen al centro abajo.
 * Conecta los canales de venta paralelos con el hub central.
 */
export function FlowConverge({ count }: { count: number }) {
    if (count <= 0) return null;

    const width = 600;
    const height = 36;
    const centerX = width / 2;
    const colWidth = width / count;

    return (
        <div className="flex items-center justify-center py-1 w-full">
            <svg
                width="100%"
                height={height}
                viewBox={`0 0 ${width} ${height}`}
                fill="none"
                className="text-gray-400 dark:text-gray-500 dark:text-gray-400 max-w-full"
                preserveAspectRatio="xMidYMid meet"
            >
                {Array.from({ length: count }).map((_, i) => {
                    const sourceX = colWidth * i + colWidth / 2;
                    return (
                        <path
                            key={i}
                            d={`M${sourceX} 0 L${centerX} ${height - 6}`}
                            stroke="currentColor"
                            strokeWidth="2"
                            strokeLinecap="round"
                        />
                    );
                })}
                {/* Punta de flecha en el centro */}
                <path
                    d={`M${centerX - 5} ${height - 14}L${centerX} ${height - 6}L${centerX + 5} ${height - 14}`}
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                />
            </svg>
        </div>
    );
}

/**
 * Líneas divergentes: desde un punto central se abren a N puntos abajo.
 * Conecta el hub central con cada servicio independiente.
 */
export function FlowDiverge({ count }: { count: number }) {
    if (count <= 0) return null;

    const width = 600;
    const height = 36;
    const centerX = width / 2;
    const colWidth = width / count;

    return (
        <div className="flex items-center justify-center py-1 w-full">
            <svg
                width="100%"
                height={height}
                viewBox={`0 0 ${width} ${height}`}
                fill="none"
                className="text-gray-400 dark:text-gray-500 dark:text-gray-400 max-w-full"
                preserveAspectRatio="xMidYMid meet"
            >
                {Array.from({ length: count }).map((_, i) => {
                    const targetX = colWidth * i + colWidth / 2;
                    const endY = height - 6;
                    const angle = Math.atan2(endY, targetX - centerX);
                    const len = 7;
                    const spread = Math.PI / 6;
                    const lx = targetX - len * Math.cos(angle - spread);
                    const ly = endY - len * Math.sin(angle - spread);
                    const rx = targetX - len * Math.cos(angle + spread);
                    const ry = endY - len * Math.sin(angle + spread);

                    return (
                        <g key={i}>
                            <path
                                d={`M${centerX} 0 L${targetX} ${endY}`}
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                            />
                            <path
                                d={`M${lx} ${ly}L${targetX} ${endY}L${rx} ${ry}`}
                                stroke="currentColor"
                                strokeWidth="2"
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
