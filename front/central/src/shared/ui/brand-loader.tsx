'use client';

import Image from 'next/image';

interface BrandLoaderProps {
    title?: string;
    subtitle?: string;
}

export function BrandLoader({ title = 'Procesando...', subtitle }: BrandLoaderProps) {
    return (
        <div className="flex flex-col items-center justify-center gap-4 py-6">
            <div className="relative w-24 h-24 flex items-center justify-center">
                <div
                    className="absolute inset-0 rounded-full animate-spin"
                    style={{
                        borderWidth: '3px',
                        borderStyle: 'solid',
                        borderColor: 'rgba(0,0,0,0.08)',
                        borderTopColor: 'var(--color-primary)',
                    }}
                />
                <Image
                    src="/isotipo.png"
                    alt="Probability"
                    width={44}
                    height={44}
                    className="object-contain animate-pulse"
                    priority
                />
            </div>
            <div className="text-center">
                <p className="text-base font-semibold" style={{ color: 'var(--color-primary)' }}>{title}</p>
                {subtitle && (
                    <p className="mt-1 text-sm text-gray-600 dark:text-gray-300 max-w-xs">{subtitle}</p>
                )}
            </div>
        </div>
    );
}

export function BrandLoaderOverlay({ title, subtitle }: BrandLoaderProps) {
    return (
        <div className="absolute inset-0 z-20 flex items-center justify-center bg-white/85 dark:bg-gray-800/85 backdrop-blur-sm rounded-2xl">
            <BrandLoader title={title} subtitle={subtitle} />
        </div>
    );
}
