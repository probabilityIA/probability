'use client';

import { IntegrationCategory } from '../../domain/types';

interface CategoryTabsProps {
    categories: IntegrationCategory[];
    activeCategory: string | null;
    onSelectCategory: (categoryCode: string | null) => void;
}

export function CategoryTabs({ categories, activeCategory, onSelectCategory }: CategoryTabsProps) {
    // Filter and sort categories
    const sortedCategories = [...categories]
        .filter(c => c.is_visible && c.is_active)
        .sort((a, b) => a.display_order - b.display_order);

    return (
        <div className="border-b border-gray-200 mb-6">
            <nav className="-mb-px flex space-x-8">
                {/* Tab "Todas" */}
                <button
                    onClick={() => onSelectCategory(null)}
                    className={`
                        py-4 px-1 border-b-2 font-medium text-sm whitespace-nowrap
                        ${activeCategory === null
                            ? 'border-blue-500 text-blue-600'
                            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                    `}
                >
                    Todas
                </button>

                {/* Tabs por categoría */}
                {sortedCategories.map((category) => (
                    <button
                        key={category.code}
                        onClick={() => onSelectCategory(category.code)}
                        className={`
                            py-4 px-1 border-b-2 font-medium text-sm whitespace-nowrap flex items-center gap-2
                            ${activeCategory === category.code
                                ? 'border-blue-500 text-blue-600'
                                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                        `}
                    >
                        {/* Icon opcional - se puede agregar después con heroicons */}
                        {category.icon && (
                            <span className="text-lg">{category.icon}</span>
                        )}
                        <span>{category.name}</span>
                    </button>
                ))}
            </nav>
        </div>
    );
}
