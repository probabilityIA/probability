'use client';

import * as React from 'react';
import { cn } from '@/lib/utils';

// Contexto para compartir config entre Chart y sub-componentes
type ChartConfig = {
    [key: string]: {
        label?: string;
        color?: string;
        icon?: React.ComponentType;
    };
};

const ChartContext = React.createContext<{
    config: ChartConfig;
} | null>(null);

function useChart() {
    const context = React.useContext(ChartContext);
    if (!context) {
        throw new Error('useChart must be used within a <ChartContainer />');
    }
    return context;
}

const ChartContainer = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement> & {
        config: ChartConfig;
        children: React.ReactNode;
    }
>(({ className, config, children, ...props }, ref) => {
    return (
        <ChartContext.Provider value={{ config }}>
            <div
                ref={ref}
                className={cn(
                    'relative w-full',
                    '[&_.recharts-surface]:outline-none',
                    '[&_.recharts-sector]:outline-none',
                    '[&_.recharts-rectangle]:outline-none',
                    '[&_path]:outline-none', // Elimina outline en paths SVG genericos
                    '[&_rect]:outline-none', // Elimina outline en rects SVG genericos
                    '[&_*:focus]:outline-none', // Elimina outline en cualquier elemento enfocado
                    '[&_*:focus-visible]:outline-none',
                    className
                )}
                {...props}
            >
                {children}
            </div>
        </ChartContext.Provider>
    );
});
ChartContainer.displayName = 'ChartContainer';

const ChartTooltip = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => {
    return (
        <div
            ref={ref}
            className={cn(
                'rounded-lg border bg-white p-3 shadow-lg',
                'dark:bg-gray-800 dark:border-gray-700',
                className
            )}
            {...props}
        />
    );
});
ChartTooltip.displayName = 'ChartTooltip';

const ChartTooltipContent = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement> & {
        active?: boolean;
        payload?: any[];
        label?: string;
        hideLabel?: boolean;
        hideIndicator?: boolean;
        indicator?: 'line' | 'dot' | 'dashed';
        nameKey?: string;
        labelKey?: string;
    }
>(
    (
        {
            active,
            payload,
            className,
            label,
            hideLabel = false,
            hideIndicator = false,
            indicator = 'dot',
            nameKey,
            labelKey,
        },
        ref
    ) => {
        if (!active || !payload?.length) {
            return null;
        }

        return (
            <div
                ref={ref}
                className={cn(
                    'grid gap-1.5 rounded-lg border border-gray-200 bg-white px-3 py-2 text-xs shadow-lg',
                    'dark:border-gray-700 dark:bg-gray-800',
                    className
                )}
            >
                {!hideLabel && label && (
                    <div className="font-semibold text-gray-900 dark:text-gray-100">
                        {label}
                    </div>
                )}
                <div className="grid gap-1.5">
                    {payload.map((entry: any, index: number) => {
                        const key = `item-${index}`;
                        const itemConfig = entry.dataKey;
                        const indicatorColor = entry.color || entry.fill || '#000';

                        return (
                            <div
                                key={key}
                                className="flex items-center gap-2 text-sm"
                            >
                                {!hideIndicator && (
                                    <div
                                        className={cn(
                                            'h-2.5 w-2.5 shrink-0 rounded-full',
                                            indicator === 'line' && 'h-0.5 w-4 rounded-none',
                                            indicator === 'dashed' &&
                                            'h-0.5 w-4 rounded-none border-b-2 border-dashed'
                                        )}
                                        style={{
                                            backgroundColor:
                                                indicator === 'dashed'
                                                    ? 'transparent'
                                                    : indicatorColor,
                                            borderColor:
                                                indicator === 'dashed' ? indicatorColor : undefined,
                                        }}
                                    />
                                )}
                                <span className="text-gray-600 dark:text-gray-400">
                                    {entry.name || entry.dataKey}:
                                </span>
                                <span className="font-bold text-gray-900 dark:text-gray-100">
                                    {typeof entry.value === 'number'
                                        ? entry.value.toLocaleString()
                                        : entry.value}
                                </span>
                            </div>
                        );
                    })}
                </div>
            </div>
        );
    }
);
ChartTooltipContent.displayName = 'ChartTooltipContent';

const ChartLegend = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement> & {
        payload?: any[];
    }
>(({ className, payload, ...props }, ref) => {
    if (!payload?.length) {
        return null;
    }

    return (
        <div
            ref={ref}
            className={cn(
                'flex flex-wrap items-center justify-center gap-4 pt-4',
                className
            )}
            {...props}
        >
            {payload.map((entry: any, index: number) => (
                <div
                    key={`legend-${index}`}
                    className="flex items-center gap-2 text-sm"
                >
                    <div
                        className="h-3 w-3 rounded-full"
                        style={{
                            backgroundColor: entry.color || entry.fill,
                        }}
                    />
                    <span className="text-gray-700 dark:text-gray-300">
                        {entry.value}
                    </span>
                </div>
            ))}
        </div>
    );
});
ChartLegend.displayName = 'ChartLegend';

// Definir gradientes modernos para los gráficos
const CHART_GRADIENTS = {
    purple: {
        id: 'gradientPurple',
        colors: ['#8B5CF6', '#6366F1'],
        opacities: [0.9, 0.7],
    },
    blue: {
        id: 'gradientBlue',
        colors: ['#3B82F6', '#06B6D4'],
        opacities: [0.9, 0.7],
    },
    green: {
        id: 'gradientGreen',
        colors: ['#10B981', '#34D399'],
        opacities: [0.9, 0.7],
    },
    orange: {
        id: 'gradientOrange',
        colors: ['#F97316', '#FB923C'],
        opacities: [0.9, 0.7],
    },
    pink: {
        id: 'gradientPink',
        colors: ['#EC4899', '#F472B6'],
        opacities: [0.9, 0.7],
    },
    indigo: {
        id: 'gradientIndigo',
        colors: ['#6366F1', '#8B5CF6'],
        opacities: [0.9, 0.7],
    },
    teal: {
        id: 'gradientTeal',
        colors: ['#14B8A6', '#2DD4BF'],
        opacities: [0.9, 0.7],
    },
    amber: {
        id: 'gradientAmber',
        colors: ['#F59E0B', '#FBBF24'],
        opacities: [0.9, 0.7],
    },
};

const ChartGradientDefs = () => {
    return (
        <defs>
            {Object.values(CHART_GRADIENTS).map((gradient) => (
                <linearGradient
                    key={gradient.id}
                    id={gradient.id}
                    x1="0"
                    y1="0"
                    x2="0"
                    y2="1"
                >
                    <stop
                        offset="0%"
                        stopColor={gradient.colors[0]}
                        stopOpacity={gradient.opacities[0]}
                    />
                    <stop
                        offset="100%"
                        stopColor={gradient.colors[1]}
                        stopOpacity={gradient.opacities[1]}
                    />
                </linearGradient>
            ))}
        </defs>
    );
};

export {
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
    ChartLegend,
    ChartGradientDefs,
    CHART_GRADIENTS,
    useChart,
    ChartCustomGradientBar,
};

const ChartCustomGradientBar = (
    props: any // Recharts pasa muchas props inesperadas
) => {
    const { fill, x, y, width, height, index } = props;
    // Usamos index para asegurar unicidad si dataKey no viene o es común
    const uniqueId = `gradient-bar-${index}-${x}-${y}`;

    return (
        <>
            {/* Barra con gradiente de opacidad */}
            <rect
                x={x}
                y={y}
                width={width}
                height={height}
                stroke="none"
                fill={`url(#${uniqueId})`}
            />
            {/* Línea superior sólida */}
            <rect x={x} y={y} width={width} height={2} stroke="none" fill={fill} />
            <defs>
                <linearGradient
                    id={uniqueId}
                    x1="0"
                    y1="0"
                    x2="0"
                    y2="1"
                >
                    <stop offset="0%" stopColor={fill} stopOpacity={0.5} />
                    <stop offset="100%" stopColor={fill} stopOpacity={0} />
                </linearGradient>
            </defs>
        </>
    );
};
export type { ChartConfig };
