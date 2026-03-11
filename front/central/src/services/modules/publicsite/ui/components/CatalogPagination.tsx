'use client';

import { useRouter, useSearchParams } from 'next/navigation';

interface CatalogPaginationProps {
    currentPage: number;
    totalPages: number;
    total: number;
    basePath: string;
}

export function CatalogPagination({ currentPage, totalPages, total, basePath }: CatalogPaginationProps) {
    const router = useRouter();
    const searchParams = useSearchParams();

    const handlePageChange = (newPage: number) => {
        const params = new URLSearchParams(searchParams);
        params.set('page', newPage.toString());
        router.push(`${basePath}?${params.toString()}`);
    };

    return (
        <div className="flex items-center justify-center gap-4 mt-8">
            <button
                disabled={currentPage <= 1}
                onClick={() => handlePageChange(currentPage - 1)}
                className="px-4 py-2 border border-gray-300 rounded-lg disabled:opacity-50 hover:bg-gray-50 transition-colors"
            >
                Anterior
            </button>
            <span className="text-sm text-gray-600">
                Página {currentPage} de {totalPages} ({total} productos)
            </span>
            <button
                disabled={currentPage >= totalPages}
                onClick={() => handlePageChange(currentPage + 1)}
                className="px-4 py-2 border border-gray-300 rounded-lg disabled:opacity-50 hover:bg-gray-50 transition-colors"
            >
                Siguiente
            </button>
        </div>
    );
}
