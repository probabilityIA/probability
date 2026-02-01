'use client';

import { IntegrationCategory } from '../../domain/types';

interface CategoryTabsProps {
    categories: IntegrationCategory[];
    activeCategory: string | null;
    activeTab: 'integrations' | 'types';
    onSelectCategory: (categoryCode: string | null) => void;
    onSelectTypes: () => void;
    isSuperUser?: boolean;
}

export function CategoryTabs({
    categories,
    activeCategory,
    activeTab,
    onSelectCategory,
    onSelectTypes,
    isSuperUser = false
}: CategoryTabsProps) {
    // Filter and sort categories
    const sortedCategories = [...categories]
        .filter(c => c.is_visible && c.is_active)
        .sort((a, b) => a.display_order - b.display_order);

    return (
        <div className="border-b border-gray-200 mb-6">
            <nav className="-mb-px flex space-x-8">
                {/* Tabs por categoría - solo nombre */}
                {sortedCategories.map((category) => (
                    <button
                        key={category.code}
                        onClick={() => onSelectCategory(category.code)}
                        className={`
                            py-4 px-1 border-b-2 font-medium text-sm whitespace-nowrap
                            ${activeCategory === category.code && activeTab === 'integrations'
                                ? 'border-blue-500 text-blue-600'
                                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                        `}
                    >
                        {category.name}
                    </button>
                ))}

                {/* Tab "Tipos de Integración" - solo para super usuario */}
                {isSuperUser && (
                    <button
                        onClick={onSelectTypes}
                        className={`
                            py-4 px-1 border-b-2 font-medium text-sm whitespace-nowrap
                            ${activeTab === 'types'
                                ? 'border-blue-500 text-blue-600'
                                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                        `}
                    >
                        Tipos de Integración
                    </button>
                )}
            </nav>
        </div>
    );
}
