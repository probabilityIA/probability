'use client';

import type { Container } from '../../domain/types';
import { ContainerCard } from './ContainerCard';

interface ArchitectureViewProps {
    containers: Container[];
}

type Layer = 'frontend' | 'backend' | 'data';

const LAYER_MAP: Record<string, Layer> = {
    'front-central': 'frontend',
    'font-central': 'frontend',
    'front-website': 'frontend',
    'font-website': 'frontend',
    'front-testing': 'frontend',
    'monitoring-web': 'frontend',
    'back-central': 'backend',
    'back-testing': 'backend',
    'monitoring-api': 'backend',
    'nginx': 'backend',
    'redis': 'data',
    'rabbitmq': 'data',
};

const LAYER_CONFIG: Record<Layer, { label: string; color: string }> = {
    frontend: { label: 'Frontend', color: '#00f0ff' },
    backend: { label: 'Backend', color: '#a855f7' },
    data: { label: 'Data & Messaging', color: '#ffaa00' },
};

// Specific connections between actual service pairs
// Each connection: [source service, target service, label]
const CONNECTIONS: [string, string, string][] = [
    ['front-central', 'back-central', 'HTTP :3050'],
    ['front-testing', 'back-testing', 'HTTP :9092'],
    ['monitoring-web', 'monitoring-api', 'HTTP :3070'],
    ['back-central', 'redis', 'TCP :6379'],
    ['back-central', 'rabbitmq', 'AMQP :5672'],
];

// Normalize service names (font- → front-)
function normalizeService(name: string): string {
    return name.replace(/^font-/, 'front-');
}

function LayerRow({ label, color, containers }: { label: string; color: string; containers: Container[] }) {
    return (
        <div>
            <div className="flex items-center gap-2 mb-2">
                <div className="w-1 h-4 rounded-full" style={{ background: color }} />
                <span className="text-[11px] font-semibold uppercase tracking-wider" style={{ color }}>
                    {label}
                </span>
                <span className="text-[10px] font-mono" style={{ color: '#444455' }}>
                    {containers.filter(c => c.state === 'running').length}/{containers.length}
                </span>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
                {containers.map(c => (
                    <ContainerCard key={c.id} container={c} accentColor={color} />
                ))}
            </div>
        </div>
    );
}

function ConnectionsBridge({
    fromContainers,
    toContainers,
    connections,
    fromColor,
    toColor,
}: {
    fromContainers: Container[];
    toContainers: Container[];
    connections: [string, string, string][];
    fromColor: string;
    toColor: string;
}) {
    // Filter to only connections that have both endpoints present
    const activeConnections = connections.filter(([from, to]) => {
        const hasFrom = fromContainers.some(c => normalizeService(c.service) === from);
        const hasTo = toContainers.some(c => normalizeService(c.service) === to);
        return hasFrom && hasTo;
    });

    if (activeConnections.length === 0) return <div className="h-3" />;

    return (
        <div className="flex items-center justify-center py-2">
            <div className="flex items-center gap-4 flex-wrap justify-center">
                {activeConnections.map(([from, to, label], i) => (
                    <div key={i} className="flex items-center gap-1.5">
                        {/* From label */}
                        <span className="text-[8px] font-mono" style={{ color: `${fromColor}70` }}>
                            {from.replace('front-', '').replace('back-', '')}
                        </span>

                        {/* Animated line */}
                        <div className="flex items-center">
                            <div className="h-px w-3" style={{ background: `${fromColor}40` }} />
                            <div className="relative w-10 h-3 overflow-hidden flex items-center">
                                <div className="absolute h-px w-full" style={{ background: `linear-gradient(90deg, ${fromColor}30, ${toColor}30)` }} />
                                <div
                                    className="absolute h-px w-3 rounded-full animate-flow"
                                    style={{
                                        background: `linear-gradient(90deg, ${fromColor}, ${toColor})`,
                                        boxShadow: `0 0 4px ${toColor}60`,
                                    }}
                                />
                            </div>
                            <div className="h-px w-3" style={{ background: `${toColor}40` }} />
                        </div>

                        {/* Label */}
                        <span className="text-[7px] font-mono px-1 py-0.5 rounded"
                            style={{ color: '#555570', background: '#ffffff04', border: '1px solid #ffffff06' }}>
                            {label}
                        </span>

                        {/* To label */}
                        <span className="text-[8px] font-mono" style={{ color: `${toColor}70` }}>
                            {to.replace('front-', '').replace('back-', '')}
                        </span>
                    </div>
                ))}
            </div>
        </div>
    );
}

export function ArchitectureView({ containers }: ArchitectureViewProps) {
    // Deduplicate by normalized service name
    const seen = new Set<string>();
    const uniqueContainers = containers.filter(c => {
        const norm = normalizeService(c.service);
        if (seen.has(norm)) return false;
        seen.add(norm);
        return true;
    });

    // Split into layers
    const layers: Record<Layer, Container[]> = { frontend: [], backend: [], data: [] };
    for (const c of uniqueContainers) {
        const layer = LAYER_MAP[normalizeService(c.service)] || LAYER_MAP[c.service];
        if (layer) layers[layer].push(c);
    }

    // Connections between layers
    const frontToBack = CONNECTIONS.filter(([from]) =>
        layers.frontend.some(c => normalizeService(c.service) === from)
    );
    const backToData = CONNECTIONS.filter(([from]) =>
        layers.backend.some(c => normalizeService(c.service) === from)
    );

    return (
        <div className="space-y-0">
            {/* Frontend row */}
            {layers.frontend.length > 0 && (
                <LayerRow
                    label={LAYER_CONFIG.frontend.label}
                    color={LAYER_CONFIG.frontend.color}
                    containers={layers.frontend}
                />
            )}

            {/* Connections: Frontend → Backend */}
            <ConnectionsBridge
                fromContainers={layers.frontend}
                toContainers={layers.backend}
                connections={frontToBack}
                fromColor={LAYER_CONFIG.frontend.color}
                toColor={LAYER_CONFIG.backend.color}
            />

            {/* Backend row */}
            {layers.backend.length > 0 && (
                <LayerRow
                    label={LAYER_CONFIG.backend.label}
                    color={LAYER_CONFIG.backend.color}
                    containers={layers.backend}
                />
            )}

            {/* Connections: Backend → Data */}
            <ConnectionsBridge
                fromContainers={layers.backend}
                toContainers={layers.data}
                connections={backToData}
                fromColor={LAYER_CONFIG.backend.color}
                toColor={LAYER_CONFIG.data.color}
            />

            {/* Data row */}
            {layers.data.length > 0 && (
                <LayerRow
                    label={LAYER_CONFIG.data.label}
                    color={LAYER_CONFIG.data.color}
                    containers={layers.data}
                />
            )}
        </div>
    );
}
