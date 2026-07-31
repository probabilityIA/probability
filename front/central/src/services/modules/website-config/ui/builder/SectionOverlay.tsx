'use client';

import { useState } from 'react';
import { EyeIcon, EyeSlashIcon, ChevronUpIcon, ChevronDownIcon, Bars3Icon, PencilSquareIcon } from '@heroicons/react/24/outline';

interface SectionOverlayProps {
    label: string;
    visible: boolean;
    selected: boolean;
    isFirst: boolean;
    isLast: boolean;
    hasEditor: boolean;
    onSelect: () => void;
    onToggleVisible: () => void;
    onMoveUp: () => void;
    onMoveDown: () => void;
    onDragStart: (e: React.DragEvent) => void;
    onDragOver: (e: React.DragEvent) => void;
    onDrop: (e: React.DragEvent) => void;
    children: React.ReactNode;
}

export function SectionOverlay({
    label, visible, selected, isFirst, isLast, hasEditor,
    onSelect, onToggleVisible, onMoveUp, onMoveDown,
    onDragStart, onDragOver, onDrop, children,
}: SectionOverlayProps) {
    const [hover, setHover] = useState(false);
    const active = hover || selected;

    return (
        <div
            className={`relative transition-shadow ${selected ? 'ring-2 ring-blue-500' : active ? 'ring-2 ring-blue-300' : ''}`}
            onMouseEnter={() => setHover(true)}
            onMouseLeave={() => setHover(false)}
            onClick={(e) => { e.stopPropagation(); onSelect(); }}
            onDragOver={onDragOver}
            onDrop={onDrop}
        >
            <div className={visible ? '' : 'opacity-40 grayscale'}>
                {children}
            </div>

            {active && (
                <div className="absolute top-2 left-2 z-20 flex items-center gap-1 bg-blue-600 text-white text-xs font-medium rounded-md px-2 py-1 shadow">
                    <span
                        draggable
                        onDragStart={onDragStart}
                        className="cursor-grab active:cursor-grabbing"
                        title="Arrastrar para reordenar"
                    >
                        <Bars3Icon className="w-4 h-4" />
                    </span>
                    <span>{label}</span>
                    {hasEditor && <PencilSquareIcon className="w-4 h-4 ml-1" />}
                </div>
            )}

            {active && (
                <div className="absolute top-2 right-2 z-20 flex items-center gap-1">
                    <button
                        type="button"
                        onClick={(e) => { e.stopPropagation(); onMoveUp(); }}
                        disabled={isFirst}
                        className="p-1.5 bg-white dark:bg-gray-700 rounded-md shadow border border-gray-200 dark:border-gray-600 text-gray-600 dark:text-gray-200 hover:bg-gray-50 disabled:opacity-30"
                        title="Subir seccion"
                    >
                        <ChevronUpIcon className="w-4 h-4" />
                    </button>
                    <button
                        type="button"
                        onClick={(e) => { e.stopPropagation(); onMoveDown(); }}
                        disabled={isLast}
                        className="p-1.5 bg-white dark:bg-gray-700 rounded-md shadow border border-gray-200 dark:border-gray-600 text-gray-600 dark:text-gray-200 hover:bg-gray-50 disabled:opacity-30"
                        title="Bajar seccion"
                    >
                        <ChevronDownIcon className="w-4 h-4" />
                    </button>
                    <button
                        type="button"
                        onClick={(e) => { e.stopPropagation(); onToggleVisible(); }}
                        className={`p-1.5 rounded-md shadow border ${visible ? 'bg-white dark:bg-gray-700 border-gray-200 dark:border-gray-600 text-gray-600 dark:text-gray-200 hover:bg-gray-50' : 'bg-amber-500 border-amber-500 text-white'}`}
                        title={visible ? 'Ocultar seccion' : 'Mostrar seccion'}
                    >
                        {visible ? <EyeIcon className="w-4 h-4" /> : <EyeSlashIcon className="w-4 h-4" />}
                    </button>
                </div>
            )}

            {!visible && active && (
                <div className="absolute inset-0 z-10 flex items-center justify-center pointer-events-none">
                    <span className="bg-amber-500 text-white text-xs font-semibold px-3 py-1 rounded-full shadow">Seccion oculta</span>
                </div>
            )}
        </div>
    );
}
