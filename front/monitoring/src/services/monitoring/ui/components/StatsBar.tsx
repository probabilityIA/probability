'use client';

import { useContainerStats } from '../hooks/useContainerStats';

interface StatsBarProps {
    containerId: string;
    state: string;
}

function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B';
    const units = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`;
}

function ProgressBar({ label, value, max, color }: { label: string; value: string; max: number; color: string }) {
    const pct = Math.min(max, 100);
    return (
        <div className="space-y-1">
            <div className="flex items-center justify-between text-[11px]">
                <span style={{ color: '#8888a0' }}>{label}</span>
                <span style={{ color }}>{value}</span>
            </div>
            <div className="h-1.5 rounded-full overflow-hidden" style={{ background: '#1e1e2e' }}>
                <div
                    className="h-full rounded-full transition-all duration-500 neon-bar"
                    style={{
                        width: `${pct}%`,
                        background: `linear-gradient(90deg, ${color}60, ${color})`,
                        boxShadow: `0 0 10px ${color}30, 0 0 20px ${color}10`,
                    }}
                />
            </div>
        </div>
    );
}

export function StatsBar({ containerId, state }: StatsBarProps) {
    const { stats, loading } = useContainerStats({
        containerId,
        enabled: state === 'running',
    });

    if (state !== 'running') {
        return (
            <div className="rounded-xl p-4 text-xs text-center" style={{ background: '#12121a', border: '1px solid #1e1e2e', color: '#55556a' }}>
                Container is not running
            </div>
        );
    }

    if (loading || !stats) {
        return (
            <div className="rounded-xl p-4" style={{ background: '#12121a', border: '1px solid #1e1e2e' }}>
                <div className="animate-pulse space-y-3">
                    <div className="h-4 rounded" style={{ background: '#1e1e2e', width: '60%' }} />
                    <div className="h-1.5 rounded-full" style={{ background: '#1e1e2e' }} />
                    <div className="h-4 rounded" style={{ background: '#1e1e2e', width: '50%' }} />
                    <div className="h-1.5 rounded-full" style={{ background: '#1e1e2e' }} />
                </div>
            </div>
        );
    }

    const cpuColor = stats.cpu_percent > 80 ? '#ff3366' : stats.cpu_percent > 50 ? '#ffaa00' : '#00f0ff';
    const memColor = stats.memory_percent > 80 ? '#ff3366' : stats.memory_percent > 50 ? '#ffaa00' : '#a855f7';

    return (
        <div className="rounded-xl p-4 space-y-3" style={{ background: '#12121a', border: '1px solid #1e1e2e' }}>
            <ProgressBar
                label="CPU"
                value={`${stats.cpu_percent.toFixed(1)}%`}
                max={stats.cpu_percent}
                color={cpuColor}
            />
            <ProgressBar
                label="Memory"
                value={`${formatBytes(stats.memory_usage)} / ${formatBytes(stats.memory_limit)}`}
                max={stats.memory_percent}
                color={memColor}
            />
            <div className="flex items-center justify-between text-[10px] pt-1" style={{ color: '#55556a' }}>
                <span>RX: {formatBytes(stats.network_rx)}</span>
                <span>TX: {formatBytes(stats.network_tx)}</span>
            </div>
        </div>
    );
}
