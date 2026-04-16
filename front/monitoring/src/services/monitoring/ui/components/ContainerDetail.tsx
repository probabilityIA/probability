import type { Container } from '../../domain/types';

interface ContainerDetailProps {
    container: Container;
}

function formatDate(dateStr: string): string {
    if (!dateStr) return '--';
    return new Date(dateStr).toLocaleString();
}

function formatUptime(startedAt: string): string {
    if (!startedAt) return '--';
    const diff = Date.now() - new Date(startedAt).getTime();
    if (diff < 0) return '--';
    const days = Math.floor(diff / 86400000);
    const hours = Math.floor((diff % 86400000) / 3600000);
    const minutes = Math.floor((diff % 3600000) / 60000);
    if (days > 0) return `${days}d ${hours}h ${minutes}m`;
    return `${hours}h ${minutes}m`;
}

function getStateStyle(state: string) {
    switch (state) {
        case 'running': return { color: '#00ff88', bg: '#00ff8810' };
        case 'exited': return { color: '#ff3366', bg: '#ff336610' };
        case 'restarting': return { color: '#ffaa00', bg: '#ffaa0010' };
        default: return { color: '#8888a0', bg: '#8888a010' };
    }
}

export function ContainerDetail({ container }: ContainerDetailProps) {
    const stateStyle = getStateStyle(container.state);

    const rows = [
        { label: 'ID', value: container.id.slice(0, 12) },
        { label: 'Service', value: container.service },
        { label: 'Image', value: container.image },
        { label: 'Project', value: container.project },
        { label: 'Created', value: formatDate(container.created_at) },
        { label: 'Started', value: formatDate(container.started_at) },
        { label: 'Uptime', value: container.state === 'running' ? formatUptime(container.started_at) : '--' },
        { label: 'Ports', value: container.ports?.length > 0 ? container.ports.map(p => `${p.host_port}:${p.container_port}/${p.protocol}`).join(', ') : 'None' },
    ];

    return (
        <div className="rounded-xl p-4" style={{ background: '#12121a', border: '1px solid #1e1e2e' }}>
            {/* Header with state */}
            <div className="flex items-center gap-3 mb-4">
                <h2 className="text-base font-semibold" style={{ color: '#e4e4ef' }}>
                    {container.service || container.name}
                </h2>
                <span
                    className="text-[11px] px-2 py-0.5 rounded-md font-medium"
                    style={{ background: stateStyle.bg, color: stateStyle.color }}
                >
                    {container.state}
                </span>
                {container.health && container.health !== 'none' && container.health !== '' && (
                    <span
                        className="text-[11px] px-2 py-0.5 rounded-md"
                        style={{
                            background: container.health === 'healthy' ? '#00ff8815' : '#ff336615',
                            color: container.health === 'healthy' ? '#00ff88' : '#ff3366',
                        }}
                    >
                        {container.health}
                    </span>
                )}
            </div>

            {/* Info grid */}
            <div className="grid grid-cols-2 gap-x-6 gap-y-2">
                {rows.map(row => (
                    <div key={row.label} className="flex items-baseline gap-2 text-xs py-1">
                        <span className="shrink-0" style={{ color: '#55556a', minWidth: '60px' }}>{row.label}</span>
                        <span className="font-mono truncate" style={{ color: '#c8c8d8' }}>{row.value}</span>
                    </div>
                ))}
            </div>
        </div>
    );
}
