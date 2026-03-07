'use client';

import { useRouter, useSearchParams } from 'next/navigation';

interface Props {
    currentPage: number;
    totalPages: number;
    total: number;
    basePath: string;
    label?: string;
}

export function StorefrontPagination({ currentPage, totalPages, total, basePath, label = 'registros' }: Props) {
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
                className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-800 text-gray-700 dark:text-gray-300"
            >
                Anterior
            </button>
            <span className="text-sm text-gray-600 dark:text-gray-400">
                Pagina {currentPage} de {totalPages} ({total} {label})
            </span>
            <button
                disabled={currentPage >= totalPages}
                onClick={() => handlePageChange(currentPage + 1)}
                className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-800 text-gray-700 dark:text-gray-300"
            >
                Siguiente
            </button>
        </div>
    );
}
