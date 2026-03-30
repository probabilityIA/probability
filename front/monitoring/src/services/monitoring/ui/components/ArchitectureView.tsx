'use client';

import type { Container } from '../../domain/types';
import { ContainerCard } from './ContainerCard';

interface ArchitectureViewProps {
    containers: Container[];
}

// Real service topology
interface ServiceGroup {
    id: string;
    label: string;
    description: string;
    color: string;
    bg: string;
    border: string;
    serviceNames: string[];
    connections: { from: string; to: string; label: string; protocol: string }[];
}

const SERVICE_GROUPS: ServiceGroup[] = [
    {
        id: 'gateway',
        label: 'Gateway',
        description: 'Reverse proxy - routes traffic to frontends',
        color: '#ff6b6b',
        bg: '#ff6b6b06',
        border: '#ff6b6b18',
        serviceNames: ['nginx'],
        connections: [],
    },
    {
        id: 'central',
        label: 'Central App',
        description: 'Main platform - orders, products, invoicing',
        color: '#00f0ff',
        bg: '#00f0ff06',
        border: '#00f0ff18',
        serviceNames: ['font-central', 'front-central', 'back-central'],
        connections: [
            { from: 'font-central', to: 'back-central', label: 'HTTP API', protocol: ':3050' },
            { from: 'front-central', to: 'back-central', label: 'HTTP API', protocol: ':3050' },
        ],
    },
    {
        id: 'testing',
        label: 'Testing',
        description: 'Integration testing environment',
        color: '#a855f7',
        bg: '#a855f706',
        border: '#a855f718',
        serviceNames: ['front-testing', 'back-testing'],
        connections: [
            { from: 'front-testing', to: 'back-testing', label: 'HTTP API', protocol: ':9092' },
        ],
    },
    {
        id: 'website',
        label: 'Website',
        description: 'Landing page - static, no backend',
        color: '#00ff88',
        bg: '#00ff8806',
        border: '#00ff8818',
        serviceNames: ['font-website', 'front-website'],
        connections: [],
    },
    {
        id: 'datastore',
        label: 'Data & Messaging',
        description: 'Shared by back-central',
        color: '#ffaa00',
        bg: '#ffaa0006',
        border: '#ffaa0018',
        serviceNames: ['redis', 'rabbitmq'],
        connections: [],
    },
    {
        id: 'monitoring',
        label: 'Monitoring',
        description: 'This dashboard',
        color: '#8888a0',
        bg: '#8888a006',
        border: '#8888a018',
        serviceNames: ['monitoring-web', 'monitoring-api'],
        connections: [
            { from: 'monitoring-web', to: 'monitoring-api', label: 'HTTP API', protocol: ':3070' },
        ],
    },
];

// Inter-group connections
const GROUP_CONNECTIONS: { fromGroup: string; toGroup: string; label: string; fromColor: string; toColor: string }[] = [
    { fromGroup: 'gateway', toGroup: 'central', label: 'HTTPS → :8080', fromColor: '#ff6b6b', toColor: '#00f0ff' },
    { fromGroup: 'gateway', toGroup: 'website', label: 'HTTPS → :8081', fromColor: '#ff6b6b', toColor: '#00ff88' },
    { fromGroup: 'central', toGroup: 'datastore', label: 'TCP :6379 / AMQP :5672', fromColor: '#00f0ff', toColor: '#ffaa00' },
];

function findContainer(containers: Container[], names: string[]): Container | undefined {
    return containers.find(c => names.includes(c.service));
}

function ConnectionLine({ fromColor, toColor, label }: { fromColor: string; toColor: string; label: string }) {
    return (
        <div className="flex items-center justify-center py-1.5">
            <div className="flex items-center gap-2">
                <div className="h-px w-8" style={{ background: `linear-gradient(90deg, ${fromColor}50, ${fromColor}20)` }} />
                <div className="relative flex items-center">
                    {/* Animated dot */}
                    <div className="relative h-5 w-px overflow-hidden">
                        <div className="absolute inset-0" style={{ background: `linear-gradient(to bottom, ${fromColor}30, ${toColor}30)` }} />
                        <div
                            className="absolute w-1 h-2 -left-px rounded-full animate-flow"
                            style={{ background: toColor, boxShadow: `0 0 4px ${toColor}60` }}
                        />
                    </div>
                    <span className="text-[8px] font-mono px-2 tracking-wider" style={{ color: '#444460' }}>{label}</span>
                    <div className="relative h-5 w-px overflow-hidden">
                        <div className="absolute inset-0" style={{ background: `linear-gradient(to bottom, ${fromColor}30, ${toColor}30)` }} />
                        <div
                            className="absolute w-1 h-2 -left-px rounded-full animate-flow-delayed"
                            style={{ background: toColor, boxShadow: `0 0 4px ${toColor}60` }}
                        />
                    </div>
                </div>
                <div className="h-px w-8" style={{ background: `linear-gradient(90deg, ${toColor}20, ${toColor}50)` }} />
            </div>
        </div>
    );
}

function InternalConnection({ color, label, protocol }: { color: string; label: string; protocol: string }) {
    return (
        <div className="flex items-center justify-center">
            <div className="flex items-center gap-1">
                <div className="h-px w-4" style={{ background: `${color}30` }} />
                <div className="relative w-6 h-px overflow-hidden">
                    <div className="absolute inset-0" style={{ background: `${color}25` }} />
                    <div
                        className="absolute h-px w-2 animate-flow"
                        style={{ background: color, boxShadow: `0 0 3px ${color}` }}
                    />
                </div>
                <span className="text-[7px] font-mono px-1 whitespace-nowrap" style={{ color: `${color}60` }}>
                    {label} {protocol}
                </span>
                <div className="relative w-6 h-px overflow-hidden">
                    <div className="absolute inset-0" style={{ background: `${color}25` }} />
                    <div
                        className="absolute h-px w-2 animate-flow-delayed"
                        style={{ background: color, boxShadow: `0 0 3px ${color}` }}
                    />
                </div>
                <div className="h-px w-4" style={{ background: `${color}30` }} />
            </div>
        </div>
    );
}

function ServiceGroupSection({
    group,
    containers,
}: {
    group: ServiceGroup;
    containers: Container[];
}) {
    const groupContainers = containers.filter(c => group.serviceNames.includes(c.service));
    if (groupContainers.length === 0) return null;

    const running = groupContainers.filter(c => c.state === 'running').length;
    const hasInternalConnections = group.connections.length > 0;

    // Separate front and back for connected groups
    const frontContainers = groupContainers.filter(c =>
        c.service.startsWith('front-') || c.service.startsWith('font-') || c.service === 'monitoring-web'
    );
    const backContainers = groupContainers.filter(c =>
        c.service.startsWith('back-') || c.service === 'monitoring-api'
    );
    const otherContainers = groupContainers.filter(c =>
        !frontContainers.includes(c) && !backContainers.includes(c)
    );

    return (
        <div>
            {/* Group header */}
            <div className="flex items-center gap-3 mb-2">
                <div className="w-1 h-4 rounded-full" style={{ background: group.color }} />
                <span className="text-xs font-semibold uppercase tracking-wider" style={{ color: group.color }}>
                    {group.label}
                </span>
                <span className="text-[10px] font-mono" style={{ color: '#444455' }}>
                    {running}/{groupContainers.length}
                </span>
                <span className="text-[9px]" style={{ color: '#333344' }}>
                    {group.description}
                </span>
            </div>

            {/* Group container */}
            <div className="rounded-xl p-3" style={{ background: group.bg, border: `1px solid ${group.border}` }}>
                {hasInternalConnections && frontContainers.length > 0 && backContainers.length > 0 ? (
                    <div className="space-y-0">
                        {/* Front tier */}
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                            {frontContainers.map(c => (
                                <ContainerCard key={c.id} container={c} accentColor={group.color} />
                            ))}
                        </div>

                        {/* Internal connection arrows */}
                        {group.connections.slice(0, 1).map((conn, i) => (
                            <InternalConnection key={i} color={group.color} label={conn.label} protocol={conn.protocol} />
                        ))}

                        {/* Back tier */}
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                            {backContainers.map(c => (
                                <ContainerCard key={c.id} container={c} accentColor={group.color} />
                            ))}
                        </div>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                        {groupContainers.map(c => (
                            <ContainerCard key={c.id} container={c} accentColor={group.color} />
                        ))}
                    </div>
                )}

                {otherContainers.length > 0 && (
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 mt-3">
                        {otherContainers.map(c => (
                            <ContainerCard key={c.id} container={c} accentColor={group.color} />
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}

export function ArchitectureView({ containers }: ArchitectureViewProps) {
    const renderedGroups = SERVICE_GROUPS.filter(g =>
        containers.some(c => g.serviceNames.includes(c.service))
    );

    return (
        <div className="space-y-0">
            {renderedGroups.map((group, idx) => {
                // Find group connections TO this group
                const incomingConnection = GROUP_CONNECTIONS.find(gc => gc.toGroup === group.id);
                // Find group connections FROM prev group
                const prevGroup = idx > 0 ? renderedGroups[idx - 1] : null;
                const connectionBetween = prevGroup
                    ? GROUP_CONNECTIONS.find(gc => gc.fromGroup === prevGroup.id && gc.toGroup === group.id)
                    : null;

                return (
                    <div key={group.id}>
                        {/* Inter-group connection if exists */}
                        {connectionBetween && (
                            <ConnectionLine
                                fromColor={connectionBetween.fromColor}
                                toColor={connectionBetween.toColor}
                                label={connectionBetween.label}
                            />
                        )}

                        {/* Spacer if no connection but not first */}
                        {!connectionBetween && idx > 0 && <div className="h-4" />}

                        <ServiceGroupSection group={group} containers={containers} />
                    </div>
                );
            })}
        </div>
    );
}
