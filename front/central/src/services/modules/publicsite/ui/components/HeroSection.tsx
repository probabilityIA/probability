'use client';

import Link from 'next/link';
import { useRef, useState } from 'react';
import { HeroBlock, HeroContent, HeroPosition, PublicBusiness } from '../../domain/types';

interface HeroSectionProps {
    content: HeroContent | null;
    business: PublicBusiness;
    slug: string;
    editable?: boolean;
    onContentPositionChange?: (position: HeroPosition) => void;
    onBlockPositionChange?: (index: number, position: HeroPosition) => void;
}

const DEFAULT_POSITION: HeroPosition = { x: 50, y: 50 };
const SNAP_THRESHOLD = 4;

function DraggableArea({
    editable,
    position,
    onChange,
    className,
    children,
}: {
    editable?: boolean;
    position: HeroPosition;
    onChange?: (position: HeroPosition) => void;
    className?: string;
    children: React.ReactNode;
}) {
    const containerRef = useRef<HTMLDivElement>(null);
    const [dragging, setDragging] = useState(false);
    const [guide, setGuide] = useState({ x: false, y: false });

    const handleMove = (clientX: number, clientY: number) => {
        const rect = containerRef.current?.getBoundingClientRect();
        if (!rect || rect.width === 0 || rect.height === 0) return;
        let x = ((clientX - rect.left) / rect.width) * 100;
        let y = ((clientY - rect.top) / rect.height) * 100;
        x = Math.min(96, Math.max(4, x));
        y = Math.min(96, Math.max(4, y));
        const snapX = Math.abs(x - 50) < SNAP_THRESHOLD;
        const snapY = Math.abs(y - 50) < SNAP_THRESHOLD;
        if (snapX) x = 50;
        if (snapY) y = 50;
        setGuide({ x: snapX, y: snapY });
        onChange?.({ x, y });
    };

    return (
        <div ref={containerRef} className={`relative ${className || ''}`}>
            {editable && dragging && guide.x && (
                <div className="absolute left-1/2 top-0 bottom-0 w-px bg-blue-400 -translate-x-1/2 pointer-events-none z-20" />
            )}
            {editable && dragging && guide.y && (
                <div className="absolute top-1/2 left-0 right-0 h-px bg-blue-400 -translate-y-1/2 pointer-events-none z-20" />
            )}
            <div
                className={editable ? 'absolute cursor-move select-none touch-none' : 'absolute'}
                style={{ left: `${position.x}%`, top: `${position.y}%`, transform: 'translate(-50%, -50%)' }}
                onPointerDown={editable ? (e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
                    setDragging(true);
                } : undefined}
                onPointerMove={editable ? (e) => {
                    if (!dragging) return;
                    e.stopPropagation();
                    handleMove(e.clientX, e.clientY);
                } : undefined}
                onPointerUp={editable ? (e) => {
                    e.stopPropagation();
                    setDragging(false);
                    setGuide({ x: false, y: false });
                    (e.currentTarget as HTMLElement).releasePointerCapture(e.pointerId);
                } : undefined}
            >
                {children}
            </div>
        </div>
    );
}

function HeroBlockCell({ block, slug, onDark }: { block: HeroBlock; slug: string; onDark: boolean }) {
    const textColor = onDark ? 'text-white' : 'text-gray-900 dark:text-white';
    const mutedColor = onDark ? 'text-white/80' : 'text-gray-600 dark:text-gray-300';

    if (block.type === 'image') {
        return block.image ? (
            <img src={block.image} alt="" className="w-40 h-28 object-cover rounded-lg" />
        ) : null;
    }

    if (block.type === 'button') {
        const href = block.button_link
            ? (block.button_link.startsWith('http') ? block.button_link : `/tienda/${slug}${block.button_link}`)
            : `/tienda/${slug}/productos`;
        return (
            <Link
                href={href}
                className="inline-block px-8 py-4 bg-white dark:bg-gray-800 rounded-lg font-bold text-lg transition-transform hover:scale-105 whitespace-nowrap"
                style={{ color: 'var(--brand-primary)' }}
            >
                {block.button_text || 'Ver más'}
            </Link>
        );
    }

    if (block.type === 'text') {
        return (
            <div className="flex flex-col items-center text-center gap-2 max-w-xs">
                {block.title && (
                    <h2 className={`text-2xl md:text-3xl font-bold ${textColor}`} style={{ fontSize: block.title_size ? `${block.title_size}px` : undefined }}>
                        {block.title}
                    </h2>
                )}
                {block.subtitle && (
                    <p className={`text-base ${mutedColor}`} style={{ fontSize: block.subtitle_size ? `${block.subtitle_size}px` : undefined }}>
                        {block.subtitle}
                    </p>
                )}
            </div>
        );
    }

    return null;
}

export function HeroSection({ content, business, slug, editable, onContentPositionChange, onBlockPositionChange }: HeroSectionProps) {
    const title = content?.title || `Bienvenido a ${business.name}`;
    const subtitle = content?.subtitle || business.description || 'Descubre nuestros productos';
    const ctaText = content?.cta_text || 'Ver Productos';
    const bgImage = content?.background_image;

    const rows = content?.layout_rows || 1;
    const columns = content?.layout_columns || 1;
    const blocks = content?.blocks;
    const isGrid = (rows > 1 || columns > 1) && !!blocks && blocks.length > 0;
    const columnWidths = content?.column_widths;
    const gridTemplateColumns = columnWidths && columnWidths.length === columns
        ? columnWidths.map(w => `${w}fr`).join(' ')
        : `repeat(${columns}, 1fr)`;

    return (
        <section
            className="relative py-16 px-4 min-h-[400px]"
            style={{
                backgroundColor: bgImage ? undefined : 'var(--brand-primary)',
                backgroundImage: bgImage ? `url(${bgImage})` : undefined,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
            }}
        >
            {bgImage && <div className="absolute inset-0 bg-black/50" />}

            {isGrid ? (
                <div
                    className="relative z-10 max-w-5xl mx-auto grid gap-6 h-full min-h-[360px]"
                    style={{ gridTemplateColumns, gridTemplateRows: `repeat(${rows}, minmax(120px, auto))` }}
                >
                    {blocks!.map((block, i) => (
                        <DraggableArea
                            key={i}
                            editable={editable}
                            position={block.position || DEFAULT_POSITION}
                            onChange={(p) => onBlockPositionChange?.(i, p)}
                            className="border border-dashed border-white/0 hover:border-white/30"
                        >
                            <HeroBlockCell block={block} slug={slug} onDark={!!bgImage} />
                        </DraggableArea>
                    ))}
                </div>
            ) : (
                <div className="relative z-10 h-full min-h-[360px]">
                    <DraggableArea
                        editable={editable}
                        position={content?.position || DEFAULT_POSITION}
                        onChange={(p) => onContentPositionChange?.(p)}
                        className="w-full h-full"
                    >
                        <div className="max-w-3xl flex flex-col items-center text-center px-4">
                            <h1
                                className="text-4xl md:text-5xl font-bold text-white mb-4"
                                style={{ fontSize: content?.title_size ? `${content.title_size}px` : undefined }}
                            >
                                {title}
                            </h1>
                            <p
                                className="text-xl text-white/80 mb-8"
                                style={{ fontSize: content?.subtitle_size ? `${content.subtitle_size}px` : undefined }}
                            >
                                {subtitle}
                            </p>
                            <Link
                                href={`/tienda/${slug}/productos`}
                                className="inline-block px-8 py-4 bg-white dark:bg-gray-800 rounded-lg font-bold text-lg transition-transform hover:scale-105"
                                style={{ color: 'var(--brand-primary)' }}
                            >
                                {ctaText}
                            </Link>
                        </div>
                    </DraggableArea>
                </div>
            )}
        </section>
    );
}
