'use client';

import type { Container, ServiceLayer } from '../../domain/types';
import { SERVICE_LAYERS, LAYER_CONFIG } from '../../domain/types';
import { ContainerCard } from './ContainerCard';

interface ArchitectureViewProps {
    containers: Container[];
}

function getLayer(container: Container): ServiceLayer {
    return SERVICE_LAYERS[container.service] || 'monitoring';
}

function LayerSection({
    layer,
    containers,
}: {
    layer: ServiceLayer;
    containers: Container[];
}) {
    const config = LAYER_CONFIG[layer];
    const running = containers.filter(c => c.state === 'running').length;

    return (
        <div className="relative">
            {/* Layer header */}
            <div className="flex items-center gap-3 mb-3">
                <div className="flex items-center gap-2">
                    <div className="w-1 h-5 rounded-full" style={{ background: config.color }} />
                    <span className="text-xs font-semibold uppercase tracking-wider" style={{ color: config.color }}>
                        {config.label}
                    </span>
                </div>
                <span className="text-[10px] font-mono" style={{ color: '#55556a' }}>
                    {running}/{containers.length} running
                </span>
            </div>

            {/* Cards */}
            <div
                className="rounded-xl p-3 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3"
                style={{ background: config.bg, border: `1px solid ${config.border}` }}
            >
                {containers.map(c => (
                    <ContainerCard key={c.id} container={c} accentColor={config.color} />
                ))}
            </div>
        </div>
    );
}

function ConnectionArrow({ fromColor, toColor, label }: { fromColor: string; toColor: string; label: string }) {
    return (
        <div className="flex items-center justify-center py-2">
            <div className="flex flex-col items-center gap-0.5">
                {/* Animated dots flowing down */}
                <div className="relative h-8 w-px overflow-hidden">
                    <div
                        className="absolute inset-0"
                        style={{
                            background: `linear-gradient(to bottom, ${fromColor}40, ${toColor}40)`,
                        }}
                    />
                    <div
                        className="absolute w-1 h-3 -left-px rounded-full animate-flow"
                        style={{
                            background: `linear-gradient(to bottom, ${fromColor}, ${toColor})`,
                            boxShadow: `0 0 6px ${fromColor}60`,
                        }}
                    />
                </div>
                <span className="text-[9px] uppercase tracking-widest" style={{ color: '#33334480' }}>
                    {label}
                </span>
                <div className="relative h-8 w-px overflow-hidden">
                    <div
                        className="absolute inset-0"
                        style={{
                            background: `linear-gradient(to bottom, ${fromColor}40, ${toColor}40)`,
                        }}
                    />
                    <div
                        className="absolute w-1 h-3 -left-px rounded-full animate-flow-delayed"
                        style={{
                            background: `linear-gradient(to bottom, ${fromColor}, ${toColor})`,
                            boxShadow: `0 0 6px ${toColor}60`,
                        }}
                    />
                </div>
            </div>
        </div>
    );
}

export function ArchitectureView({ containers }: ArchitectureViewProps) {
    const layers: Record<ServiceLayer, Container[]> = {
        frontend: [],
        backend: [],
        infra: [],
        monitoring: [],
    };

    for (const c of containers) {
        const layer = getLayer(c);
        layers[layer].push(c);
    }

    return (
        <div className="space-y-0">
            {/* Frontend Layer */}
            {layers.frontend.length > 0 && (
                <LayerSection layer="frontend" containers={layers.frontend} />
            )}

            {/* Connection: Frontend → Backend */}
            {layers.frontend.length > 0 && layers.backend.length > 0 && (
                <ConnectionArrow
                    fromColor={LAYER_CONFIG.frontend.color}
                    toColor={LAYER_CONFIG.backend.color}
                    label="HTTP / API"
                />
            )}

            {/* Backend Layer */}
            {layers.backend.length > 0 && (
                <LayerSection layer="backend" containers={layers.backend} />
            )}

            {/* Connection: Backend → Infra */}
            {layers.backend.length > 0 && layers.infra.length > 0 && (
                <ConnectionArrow
                    fromColor={LAYER_CONFIG.backend.color}
                    toColor={LAYER_CONFIG.infra.color}
                    label="TCP / AMQP"
                />
            )}

            {/* Infra Layer */}
            {layers.infra.length > 0 && (
                <LayerSection layer="infra" containers={layers.infra} />
            )}

            {/* Monitoring Layer (small, bottom) */}
            {layers.monitoring.length > 0 && (
                <div className="mt-6 pt-4" style={{ borderTop: '1px solid #1e1e2e' }}>
                    <LayerSection layer="monitoring" containers={layers.monitoring} />
                </div>
            )}
        </div>
    );
}
