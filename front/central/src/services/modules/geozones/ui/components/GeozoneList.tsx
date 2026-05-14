'use client';

import { TrashIcon, MapPinIcon } from '@heroicons/react/24/outline';
import { Geozone } from '../../domain/types';
import TypeChip from './TypeChip';

interface GeozoneListProps {
    items: Geozone[];
    selectedId?: number | null;
    onSelect?: (g: Geozone) => void;
    onDelete?: (g: Geozone) => void;
    canDelete?: (g: Geozone) => boolean;
}

export default function GeozoneList({ items, selectedId, onSelect, onDelete, canDelete }: GeozoneListProps) {
    if (items.length === 0) {
        return (
            <div className="text-center py-10 px-4">
                <MapPinIcon className="w-10 h-10 mx-auto text-gray-300 dark:text-gray-600 mb-2" />
                <p className="text-sm text-gray-500 dark:text-gray-400">No hay geozonas con los filtros aplicados</p>
            </div>
        );
    }

    return (
        <ul className="divide-y divide-gray-200 dark:divide-gray-700">
            {items.map((g) => {
                const isSelected = selectedId === g.id;
                const isGlobal = g.business_id === 0;
                const allowDelete = onDelete && (canDelete ? canDelete(g) : !isGlobal);
                return (
                    <li
                        key={g.id}
                        onClick={() => onSelect?.(g)}
                        className={`px-4 py-3 cursor-pointer transition-all ${
                            isSelected
                                ? 'bg-purple-50 dark:bg-purple-900/20 border-l-4 border-purple-500'
                                : 'hover:bg-gray-50 dark:hover:bg-gray-700/50 border-l-4 border-transparent'
                        }`}
                    >
                        <div className="flex items-center justify-between gap-3">
                            <div className="min-w-0 flex-1">
                                <div className="flex items-center gap-2 mb-1 flex-wrap">
                                    <TypeChip type={g.type} />
                                    {isGlobal && (
                                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300 border border-blue-300 dark:border-blue-700">
                                            DANE
                                        </span>
                                    )}
                                    {g.code && (
                                        <span className="text-[10px] font-mono text-gray-500 dark:text-gray-400">#{g.code}</span>
                                    )}
                                </div>
                                <p className="text-sm font-medium text-gray-900 dark:text-white truncate">{g.name}</p>
                                <p className="text-[11px] text-gray-500 dark:text-gray-400">id: {g.id}{g.parent_id ? ` · padre: ${g.parent_id}` : ''}</p>
                            </div>
                            {allowDelete && (
                                <button
                                    onClick={(e) => { e.stopPropagation(); onDelete?.(g); }}
                                    className="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-md transition-colors"
                                    title="Eliminar"
                                >
                                    <TrashIcon className="w-4 h-4" />
                                </button>
                            )}
                        </div>
                    </li>
                );
            })}
        </ul>
    );
}
