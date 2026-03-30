import type { Container } from '../../domain/types';
import { ContainerCard } from './ContainerCard';

interface ContainerGridProps {
    containers: Container[];
}

export function ContainerGrid({ containers }: ContainerGridProps) {
    const running = containers.filter(c => c.state === 'running');
    const stopped = containers.filter(c => c.state !== 'running');

    return (
        <div className="space-y-5">
            {/* Summary bar */}
            <div
                className="flex items-center gap-5 text-xs px-4 py-2.5 rounded-lg"
                style={{ background: '#12121a', border: '1px solid #1e1e2e' }}
            >
                <div className="flex items-center gap-1.5">
                    <div className="w-2 h-2 rounded-full pulse-dot" style={{ background: '#00ff88', boxShadow: '0 0 6px #00ff88' }} />
                    <span style={{ color: '#00ff88' }}>{running.length}</span>
                    <span style={{ color: '#8888a0' }}>running</span>
                </div>
                {stopped.length > 0 && (
                    <div className="flex items-center gap-1.5">
                        <div className="w-2 h-2 rounded-full" style={{ background: '#ff3366' }} />
                        <span style={{ color: '#ff3366' }}>{stopped.length}</span>
                        <span style={{ color: '#8888a0' }}>stopped</span>
                    </div>
                )}
                <div className="ml-auto font-mono" style={{ color: '#55556a' }}>
                    {containers.length} containers
                </div>
            </div>

            {/* Grid */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                {containers.map(container => (
                    <ContainerCard key={container.id} container={container} />
                ))}
            </div>
        </div>
    );
}
