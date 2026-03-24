'use client';

import { useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

export function CatalogSearch({ initialSearch }: { initialSearch: string }) {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [search, setSearch] = useState(initialSearch);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const params = new URLSearchParams(searchParams);
        if (search) params.set('search', search);
        else params.delete('search');
        params.set('page', '1');
        router.push(`/storefront/catalogo?${params.toString()}`);
    };

    return (
        <form onSubmit={handleSubmit} className="mb-6 flex gap-2">
            <input
                type="text"
                value={search}
                onChange={e => setSearch(e.target.value)}
                placeholder="Buscar productos..."
                className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
            />
            <button type="submit" className="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
                Buscar
            </button>
        </form>
    );
}
