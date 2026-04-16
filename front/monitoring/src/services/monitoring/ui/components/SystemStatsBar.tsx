'use client';

import { useSystemStats } from '../hooks/useSystemStats';

function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B';
    const units = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`;
}

function MiniBar({ label, value, percent, color }: { label: string; value: string; percent: number; color: string }) {
    return (
        <div className="flex items-center gap-3 min-w-0">
            <span className="text-[10px] uppercase tracking-wider shrink-0 w-10" style={{ color: '#55556a' }}>{label}</span>
            <div className="flex-1 h-1.5 rounded-full overflow-hidden" style={{ background: '#1e1e2e', minWidth: 60 }}>
                <div
                    className="h-full rounded-full neon-bar transition-all duration-700"
                    style={{
                        width: `${Math.min(percent, 100)}%`,
                        background: `linear-gradient(90deg, ${color}60, ${color})`,
                        boxShadow: `0 0 8px ${color}30`,
                    }}
                />
            </div>
            <span className="text-[11px] font-mono shrink-0" style={{ color }}>{value}</span>
        </div>
    );
}

export function SystemStatsBar() {
    const stats = useSystemStats(5000);

    if (!stats) {
        return (
            <div className="flex items-center gap-6 px-5 py-3 rounded-lg animate-pulse"
                style={{ background: '#12121a', border: '1px solid #1e1e2e' }}>
                <div className="h-3 w-24 rounded" style={{ background: '#1e1e2e' }} />
                <div className="h-3 w-24 rounded" style={{ background: '#1e1e2e' }} />
                <div className="h-3 w-24 rounded" style={{ background: '#1e1e2e' }} />
            </div>
        );
    }

    const cpuColor = stats.cpu_percent > 80 ? '#ff3366' : stats.cpu_percent > 50 ? '#ffaa00' : '#00f0ff';
    const memColor = stats.memory_percent > 80 ? '#ff3366' : stats.memory_percent > 50 ? '#ffaa00' : '#a855f7';
    const diskColor = stats.disk_percent > 90 ? '#ff3366' : stats.disk_percent > 70 ? '#ffaa00' : '#00ff88';

    return (
        <div
            className="grid grid-cols-1 sm:grid-cols-3 gap-3 sm:gap-6 px-5 py-3 rounded-lg"
            style={{ background: '#12121a', border: '1px solid #1e1e2e' }}
        >
            <MiniBar label="CPU" value={`${stats.cpu_percent.toFixed(1)}% (${stats.cpu_cores} cores)`} percent={stats.cpu_percent} color={cpuColor} />
            <MiniBar label="RAM" value={`${formatBytes(stats.memory_used)} / ${formatBytes(stats.memory_total)}`} percent={stats.memory_percent} color={memColor} />
            <MiniBar label="Disk" value={`${formatBytes(stats.disk_used)} / ${formatBytes(stats.disk_total)}`} percent={stats.disk_percent} color={diskColor} />
        </div>
    );
}
