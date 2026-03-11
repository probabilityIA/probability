'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

interface CatalogSearchProps {
    basePath: string;
    initialSearch: string;
}

export function CatalogSearch({ basePath, initialSearch }: CatalogSearchProps) {
    const [search, setSearch] = useState(initialSearch);
    const router = useRouter();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const params = new URLSearchParams();
        if (search) params.set('search', search);
        params.set('page', '1');
        router.push(`${basePath}?${params.toString()}`);
    };

    return (
        <form onSubmit={handleSubmit} className="flex gap-2">
            <input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Buscar productos..."
                className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <button
                type="submit"
                className="px-6 py-2 rounded-lg text-white font-medium transition-colors"
                style={{ backgroundColor: 'var(--brand-secondary)' }}
            >
                Buscar
            </button>
        </form>
    );
}
