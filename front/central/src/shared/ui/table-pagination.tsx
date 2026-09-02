'use client';

import { Button } from './button';

export interface TablePaginationProps {
    currentPage: number;
    totalPages: number;
    totalItems: number;
    pageSize: number;
    onPageChange: (page: number) => void;
    onPageSizeChange?: (pageSize: number) => void;
    pageSizeOptions?: number[];
    maxVisiblePages?: number;
}

function buildPageList(currentPage: number, totalPages: number, maxVisible: number): (number | string)[] {
    const pages: (number | string)[] = [];

    if (totalPages <= maxVisible) {
        for (let i = 1; i <= totalPages; i++) pages.push(i);
        return pages;
    }

    const windowSize = maxVisible - 2;
    let start = Math.max(2, currentPage - Math.floor(windowSize / 2));
    let end = start + windowSize - 1;
    if (end >= totalPages) {
        end = totalPages - 1;
        start = Math.max(2, end - windowSize + 1);
    }

    pages.push(1);
    if (start > 2) pages.push('...');
    for (let i = start; i <= end; i++) pages.push(i);
    if (end < totalPages - 1) pages.push('...');
    pages.push(totalPages);

    return pages;
}

export function TablePagination({
    currentPage,
    totalPages,
    totalItems,
    pageSize,
    onPageChange,
    onPageSizeChange,
    pageSizeOptions = [10, 20, 50, 100],
    maxVisiblePages = 8,
}: TablePaginationProps) {
    if (totalPages <= 1 && totalItems === 0) return null;

    const safeTotalPages = Math.max(1, totalPages);
    const pages = buildPageList(currentPage, safeTotalPages, maxVisiblePages);
    const from = totalItems === 0 ? 0 : (currentPage - 1) * pageSize + 1;
    const to = Math.min(currentPage * pageSize, totalItems);

    const navButtonClass =
        'page-btn relative inline-flex items-center px-1.5 py-1 rounded-md border text-[11px] font-medium disabled:opacity-40 transition-all';

    const sizeSelect = (extraClass: string) => (
        onPageSizeChange ? (
            <div className="flex items-center gap-1">
                <label className={`text-white whitespace-nowrap ${extraClass}`}>Mostrar:</label>
                <select
                    value={pageSize}
                    onChange={(e) => onPageSizeChange(parseInt(e.target.value, 10))}
                    className={`px-1.5 py-1 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 dark:text-white bg-white dark:bg-gray-800 ${extraClass}`}
                >
                    {pageSizeOptions.map((option) => (
                        <option key={option} value={option}>{option}</option>
                    ))}
                </select>
            </div>
        ) : null
    );

    return (
        <div
            className="px-3 py-2 flex flex-col sm:flex-row items-center justify-between gap-2"
            style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}
        >
            <div className="flex-1 flex justify-between sm:hidden w-full">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onPageChange(Math.max(1, currentPage - 1))}
                    disabled={currentPage === 1}
                >
                    Anterior
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onPageChange(Math.min(safeTotalPages, currentPage + 1))}
                    disabled={currentPage === safeTotalPages}
                >
                    Siguiente
                </Button>
            </div>

            <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between w-full">
                <div className="flex items-center gap-3">
                    <p className="text-[11px] text-white">
                        Mostrando <span className="font-medium">{from}</span> a{' '}
                        <span className="font-medium">{to}</span> de{' '}
                        <span className="font-medium">{totalItems}</span> resultados
                    </p>
                    {sizeSelect('text-[11px]')}
                </div>
                <div className="flex items-center gap-1">
                    <nav className="relative z-0 inline-flex items-center gap-1 flex-wrap justify-end">
                        <button
                            onClick={() => onPageChange(1)}
                            disabled={currentPage === 1}
                            className={navButtonClass}
                            title={'Primera p\u00e1gina'}
                        >
                            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 19l-7-7 7-7m8 14l-7-7 7-7" /></svg>
                        </button>

                        <button
                            onClick={() => onPageChange(Math.max(1, currentPage - 1))}
                            disabled={currentPage === 1}
                            className={navButtonClass}
                            title="Anterior"
                        >
                            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" /></svg>
                        </button>

                        {pages.map((p, idx) =>
                            typeof p === 'string' ? (
                                <span key={`ellipsis-${idx}`} className="relative inline-flex items-center px-1 py-1 text-[11px] font-bold text-white/70">
                                    ...
                                </span>
                            ) : (
                                <button
                                    key={p}
                                    onClick={() => onPageChange(p)}
                                    className={`relative inline-flex items-center justify-center min-w-7 px-2 py-1 rounded-md border text-[11px] font-semibold transition-all ${p === currentPage ? 'page-btn-active z-10 shadow-md scale-105' : 'page-btn'}`}
                                    style={p === currentPage
                                        ? { backgroundColor: 'var(--color-secondary)', borderColor: 'var(--color-secondary)', color: 'white' }
                                        : undefined}
                                >
                                    {p}
                                </button>
                            )
                        )}

                        <button
                            onClick={() => onPageChange(Math.min(safeTotalPages, currentPage + 1))}
                            disabled={currentPage === safeTotalPages}
                            className={navButtonClass}
                            title="Siguiente"
                        >
                            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
                        </button>

                        <button
                            onClick={() => onPageChange(safeTotalPages)}
                            disabled={currentPage === safeTotalPages}
                            className={navButtonClass}
                            title={'\u00daltima p\u00e1gina'}
                        >
                            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 5l7 7-7 7M5 5l7 7-7 7" /></svg>
                        </button>
                    </nav>
                </div>
            </div>

            <div className="flex items-center justify-between w-full sm:hidden pt-2">
                {sizeSelect('text-xs')}
                <p className="text-xs text-white/80">
                    {'P\u00e1gina'} {currentPage} de {safeTotalPages}
                </p>
            </div>
        </div>
    );
}
