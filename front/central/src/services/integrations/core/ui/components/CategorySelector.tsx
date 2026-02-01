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

// Mapeo de códigos de categoría a iconos
const getCategoryIcon = (code: string) => {
    const icons: Record<string, ReactElement> = {
        'ecommerce': <ShoppingCartIcon className="w-10 h-10" />,
        'invoicing': <DocumentTextIcon className="w-10 h-10" />,
        'messaging': <ChatBubbleLeftRightIcon className="w-10 h-10" />,
        'payment': <BanknotesIcon className="w-10 h-10" />,
        'shipping': <TruckIcon className="w-10 h-10" />,
        'system': <CogIcon className="w-10 h-10" />,
    };
    return icons[code] || <CogIcon className="w-10 h-10" />;
};

export function CategorySelector({ categories, onSelect }: CategorySelectorProps) {
    const sortedCategories = [...categories]
        .filter(c => c.is_visible && c.is_active)
        .sort((a, b) => a.display_order - b.display_order);

    return (
        <div className="p-6">
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Seleccionar Categoría</h2>
            <p className="text-gray-600 mb-6">
                Elige el tipo de integración que deseas configurar
            </p>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {sortedCategories.map((category) => (
                    <button
                        key={category.code}
                        onClick={() => onSelect(category)}
                        className="p-6 border-2 border-gray-200 rounded-lg hover:border-blue-500 hover:shadow-md transition-all text-left group"
                        style={{
                            borderColor: category.color ? `${category.color}20` : undefined,
                        }}
                    >
                        <div className="flex items-start gap-4">
                            {/* Icon with category color */}
                            <div
                                className="flex-shrink-0 p-3 rounded-lg group-hover:scale-110 transition-transform"
                                style={{
                                    backgroundColor: category.color ? `${category.color}15` : '#f3f4f6',
                                    color: category.color || '#6b7280',
                                }}
                            >
                                {getCategoryIcon(category.code)}
                            </div>

                            <div className="flex-1">
                                {/* Category Name */}
                                <h3 className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 mb-1">
                                    {category.name}
                                </h3>

                                {/* Description */}
                                {category.description && (
                                    <p className="text-sm text-gray-600">
                                        {category.description}
                                    </p>
                                )}
                            </div>
                        </div>
                    </button>
                ))}
            </div>

            {sortedCategories.length === 0 && (
                <div className="text-center py-12 text-gray-500">
                    No hay categorías disponibles
                </div>
            )}
        </div>
    );
}
