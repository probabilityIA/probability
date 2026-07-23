'use client';

import { ReactElement } from 'react';
import {
    ShoppingCartIcon,
    DocumentTextIcon,
    ChatBubbleLeftRightIcon,
    CogIcon,
    BanknotesIcon,
    TruckIcon,
} from '@heroicons/react/24/outline';
import { IntegrationCategory } from '../../domain/types';

interface CategorySelectorProps {
    categories: IntegrationCategory[];
    onSelect: (category: IntegrationCategory) => void;
}

const getCategoryIcon = (code: string) => {
    const icons: Record<string, ReactElement> = {
        'ecommerce': <ShoppingCartIcon className="w-6 h-6" />,
        'invoicing': <DocumentTextIcon className="w-6 h-6" />,
        'messaging': <ChatBubbleLeftRightIcon className="w-6 h-6" />,
        'payment': <BanknotesIcon className="w-6 h-6" />,
        'shipping': <TruckIcon className="w-6 h-6" />,
        'system': <CogIcon className="w-6 h-6" />,
    };
    return icons[code] || <CogIcon className="w-6 h-6" />;
};

export function CategorySelector({ categories, onSelect }: CategorySelectorProps) {
    const sortedCategories = [...categories]
        .filter(c => c.is_visible && c.is_active)
        .sort((a, b) => a.display_order - b.display_order);

    return (
        <div className="w-[34rem] max-w-[88vw] p-2">
            <h2 className="text-base font-bold text-gray-900 dark:text-white">Seleccionar Categoria</h2>
            <p className="mb-4 text-sm text-gray-500 dark:text-gray-400">
                Elige el tipo de integracion que deseas configurar
            </p>

            <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                {sortedCategories.map((category) => (
                    <button
                        key={category.code}
                        onClick={() => onSelect(category)}
                        className="group rounded-lg border-2 border-gray-200 p-3.5 text-left transition-all hover:border-blue-500 hover:shadow-md dark:border-gray-700"
                        style={{
                            borderColor: category.color ? `${category.color}20` : undefined,
                        }}
                    >
                        <div className="flex items-center gap-3">
                            <div
                                className="flex-shrink-0 rounded-lg p-2 transition-transform group-hover:scale-110"
                                style={{
                                    backgroundColor: category.color ? `${category.color}15` : '#f3f4f6',
                                    color: category.color || '#6b7280',
                                }}
                            >
                                {getCategoryIcon(category.code)}
                            </div>

                            <div className="min-w-0 flex-1">
                                <h3 className="text-sm font-semibold text-gray-900 group-hover:text-blue-600 dark:text-white">
                                    {category.name}
                                </h3>
                                {category.description && (
                                    <p className="truncate text-xs text-gray-500 dark:text-gray-400">
                                        {category.description}
                                    </p>
                                )}
                            </div>
                        </div>
                    </button>
                ))}
            </div>

            {sortedCategories.length === 0 && (
                <div className="py-8 text-center text-sm text-gray-500 dark:text-gray-400">
                    No hay categorias disponibles
                </div>
            )}
        </div>
    );
}
