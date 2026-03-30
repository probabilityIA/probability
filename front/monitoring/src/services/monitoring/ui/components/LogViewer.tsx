'use client';

import { useRef, useEffect, useState } from 'react';
import { useLogStream } from '../hooks/useLogStream';

interface LogViewerProps {
    containerId: string;
}

function colorizeLog(line: string): { color: string } {
    const lower = line.toLowerCase();
    if (lower.includes('error') || lower.includes('fatal') || lower.includes('panic')) return { color: '#ff3366' };
    if (lower.includes('warn')) return { color: '#ffaa00' };
    if (lower.includes('debug')) return { color: '#8888a0' };
    if (lower.includes('info')) return { color: '#00f0ff' };
    return { color: '#c8c8d8' };
}

export function LogViewer({ containerId }: LogViewerProps) {
    const { lines, connected, error, clear, reconnect } = useLogStream({ containerId });
    const scrollRef = useRef<HTMLDivElement>(null);
    const [followTail, setFollowTail] = useState(true);

    useEffect(() => {
        if (followTail && scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [lines, followTail]);

    const handleScroll = () => {
        if (!scrollRef.current) return;
        const { scrollTop, scrollHeight, clientHeight } = scrollRef.current;
        const isAtBottom = scrollHeight - scrollTop - clientHeight < 40;
        setFollowTail(isAtBottom);
    };

    return (
        <div className="flex flex-col h-full rounded-xl overflow-hidden crt-flicker" style={{ border: '1px solid #1e1e2e' }}>
            {/* Toolbar */}
            <div
                className="flex items-center justify-between px-4 py-2 text-xs shrink-0"
                style={{ background: '#0d0d14', borderBottom: '1px solid #1e1e2e' }}
            >
                <div className="flex items-center gap-3">
                    <div className="flex items-center gap-1.5">
                        <div
                            className={connected ? 'pulse-dot' : ''}
                            style={{
                                width: 6, height: 6, borderRadius: '50%',
                                background: connected ? '#00ff88' : error ? '#ff3366' : '#8888a0',
                                boxShadow: connected ? '0 0 8px #00ff88' : 'none',
                            }}
                        />
                        <span className="uppercase tracking-wider font-medium" style={{ color: connected ? '#00ff88' : '#8888a0' }}>
                            {connected ? 'Live' : error || 'Disconnected'}
                        </span>
                    </div>
                    <span className="font-mono" style={{ color: '#55556a' }}>{lines.length} lines</span>
                </div>

                <div className="flex items-center gap-2">
                    <button
                        onClick={() => setFollowTail(!followTail)}
                        className="px-2 py-1 rounded transition-all cursor-pointer"
                        style={{
                            background: followTail ? '#00f0ff15' : 'transparent',
                            color: followTail ? '#00f0ff' : '#8888a0',
                            border: `1px solid ${followTail ? '#00f0ff30' : '#1e1e2e'}`,
                            boxShadow: followTail ? '0 0 8px #00f0ff10' : 'none',
                        }}
                    >
                        Follow
                    </button>
                    <button
                        onClick={clear}
                        className="px-2 py-1 rounded transition-colors cursor-pointer"
                        style={{ color: '#8888a0', border: '1px solid #1e1e2e' }}
                    >
                        Clear
                    </button>
                    {!connected && (
                        <button
                            onClick={reconnect}
                            className="px-2 py-1 rounded transition-colors cursor-pointer"
                            style={{ color: '#ffaa00', border: '1px solid #ffaa0030', background: '#ffaa0008' }}
                        >
                            Reconnect
                        </button>
                    )}
                </div>
            </div>

            {/* Log content */}
            <div className="flex-1 relative overflow-hidden" style={{ background: '#06060a' }}>
                {/* Scanline overlay */}
                <div className="scanline-overlay" />

                <div
                    ref={scrollRef}
                    onScroll={handleScroll}
                    className="absolute inset-0 overflow-auto p-4 font-mono text-[11px] leading-[1.6]"
                >
                    {lines.length === 0 ? (
                        <div className="flex items-center justify-center h-full" style={{ color: '#55556a' }}>
                            <div className="text-center">
                                <div className="text-lg mb-1" style={{ color: '#1e1e2e' }}>{'>'}_</div>
                                <div>{connected ? 'Waiting for logs...' : 'Not connected'}</div>
                            </div>
                        </div>
                    ) : (
                        lines.map((line, i) => {
                            const { color } = colorizeLog(line);
                            return (
                                <div key={i} className="hover:bg-[#ffffff06] px-1 -mx-1 rounded whitespace-pre-wrap break-all">
                                    <span style={{ color: '#333344', userSelect: 'none' }}>
                                        {String(i + 1).padStart(4, ' ')}{' '}
                                    </span>
                                    <span style={{ color }}>{line}</span>
                                </div>
                            );
                        })
                    )}
                </div>
            </div>
        </div>
    );
}
