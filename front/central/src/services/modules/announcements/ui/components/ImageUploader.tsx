'use client';

import { useState, useRef } from 'react';
import { TrashIcon, ArrowUpIcon, ArrowDownIcon, PhotoIcon } from '@heroicons/react/24/outline';

interface ImageItem {
    id?: number;
    image_url?: string;
    file?: File;
    preview: string;
    sort_order: number;
}

interface ImageUploaderProps {
    images: ImageItem[];
    onChange: (images: ImageItem[]) => void;
    maxImages?: number;
}

export default function ImageUploader({ images, onChange, maxImages = 10 }: ImageUploaderProps) {
    const fileInputRef = useRef<HTMLInputElement>(null);
    const [dragOver, setDragOver] = useState(false);

    const handleFiles = (files: FileList | null) => {
        if (!files) return;
        const remaining = maxImages - images.length;
        const newFiles = Array.from(files).slice(0, remaining);

        const newImages: ImageItem[] = newFiles.map((file, i) => ({
            file,
            preview: URL.createObjectURL(file),
            sort_order: images.length + i,
        }));

        onChange([...images, ...newImages]);
    };

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        setDragOver(false);
        handleFiles(e.dataTransfer.files);
    };

    const handleRemove = (index: number) => {
        const item = images[index];
        if (item.file && item.preview) {
            URL.revokeObjectURL(item.preview);
        }
        const updated = images.filter((_, i) => i !== index).map((img, i) => ({ ...img, sort_order: i }));
        onChange(updated);
    };

    const handleMove = (index: number, direction: -1 | 1) => {
        const newIndex = index + direction;
        if (newIndex < 0 || newIndex >= images.length) return;
        const updated = [...images];
        [updated[index], updated[newIndex]] = [updated[newIndex], updated[index]];
        onChange(updated.map((img, i) => ({ ...img, sort_order: i })));
    };

    return (
        <div className="space-y-3">
            <div
                onClick={() => fileInputRef.current?.click()}
                onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
                onDragLeave={() => setDragOver(false)}
                onDrop={handleDrop}
                className={`
                    flex flex-col items-center justify-center p-6 border-2 border-dashed rounded-lg cursor-pointer transition-colors
                    ${dragOver
                        ? 'border-purple-500 bg-purple-50 dark:bg-purple-900/10'
                        : 'border-gray-300 dark:border-gray-600 hover:border-purple-400 hover:bg-gray-50 dark:hover:bg-gray-700/50'
                    }
                    ${images.length >= maxImages ? 'opacity-50 pointer-events-none' : ''}
                `}
            >
                <PhotoIcon className="w-8 h-8 text-gray-400 mb-2" />
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Arrastra imagenes o haz click para seleccionar
                </p>
                <p className="text-xs text-gray-400 mt-1">
                    {images.length} / {maxImages} imagenes
                </p>
                <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    multiple
                    onChange={(e) => handleFiles(e.target.files)}
                    className="hidden"
                />
            </div>

            {images.length > 0 && (
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
                    {images.map((img, index) => (
                        <div key={index} className="relative group rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700">
                            <img
                                src={img.preview || img.image_url}
                                alt={`Imagen ${index + 1}`}
                                className="w-full h-24 object-cover"
                            />
                            <div className="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-colors flex items-center justify-center gap-1 opacity-0 group-hover:opacity-100">
                                {index > 0 && (
                                    <button
                                        type="button"
                                        onClick={() => handleMove(index, -1)}
                                        className="p-1 bg-white/90 rounded text-gray-700 hover:bg-white"
                                    >
                                        <ArrowUpIcon className="w-3.5 h-3.5" />
                                    </button>
                                )}
                                {index < images.length - 1 && (
                                    <button
                                        type="button"
                                        onClick={() => handleMove(index, 1)}
                                        className="p-1 bg-white/90 rounded text-gray-700 hover:bg-white"
                                    >
                                        <ArrowDownIcon className="w-3.5 h-3.5" />
                                    </button>
                                )}
                                <button
                                    type="button"
                                    onClick={() => handleRemove(index)}
                                    className="p-1 bg-red-500/90 rounded text-white hover:bg-red-600"
                                >
                                    <TrashIcon className="w-3.5 h-3.5" />
                                </button>
                            </div>
                            <div className="absolute top-1 left-1 bg-black/50 text-white text-xs px-1.5 py-0.5 rounded">
                                {index + 1}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
