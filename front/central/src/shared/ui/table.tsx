/**
 * Componente Table genérico reutilizable
 * Usa clases globales y colores dinámicos del negocio
 */

'use client';

import React, { ReactNode } from 'react';
import { DynamicFilters, FilterOption, ActiveFilter } from './dynamic-filters';

export interface TableColumn<T = Record<string, unknown>> {
  key: string;
  label: string;
  render?: (value: unknown, row: T, index: number) => ReactNode;
  width?: string;
  align?: 'left' | 'center' | 'right';
}

export interface PaginationProps {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  itemsPerPage: number;
  onPageChange: (page: number) => void;
  onItemsPerPageChange?: (itemsPerPage: number) => void;
  showItemsPerPageSelector?: boolean;
  itemsPerPageOptions?: number[];
}

export interface TableFiltersProps {
  availableFilters: FilterOption[];
  activeFilters: ActiveFilter[];
  onAddFilter: (filterKey: string, value: any) => void;
  onRemoveFilter: (filterKey: string) => void;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  onSortChange?: (sortBy: string, sortOrder: 'asc' | 'desc') => void;
  sortOptions?: Array<{ value: string; label: string }>;
}

interface TableProps<T = Record<string, unknown>> {
  columns: TableColumn<T>[];
  data: T[];
  keyExtractor?: (row: T, index: number) => string | number;
  emptyMessage?: string;
  loading?: boolean;
  onRowClick?: (row: T, index: number) => void;
  pagination?: PaginationProps;
  filters?: TableFiltersProps;
}

export function Table<T = Record<string, unknown>>({ 
  columns, 
  data, 
  keyExtractor = (_, i) => i,
  emptyMessage = 'No hay datos disponibles',
  loading = false,
  onRowClick,
  pagination,
  filters,
}: TableProps<T>) {
  const alignClass = {
    left: 'text-left',
    center: 'text-center',
    right: 'text-right',
  };

  const renderPagination = () => {
    if (!pagination) return null;

    const {
      currentPage,
      totalPages,
      totalItems,
      itemsPerPage,
      onPageChange,
      onItemsPerPageChange,
      showItemsPerPageSelector = true,
      itemsPerPageOptions = [10, 20, 50, 100],
    } = pagination;

    // Solo mostrar paginación si hay más de una página o hay items
    if (totalPages <= 1 && totalItems === 0) return null;

    return (
      <div className="pagination-alt border-t border-gray-200 rounded-t-none mt-0">
        <div className="flex items-center justify-between w-full flex-wrap gap-4">
          {/* Información y navegación - Centrado */}
          <div className="flex items-center gap-3 flex-1 justify-center">
            <button
              onClick={() => onPageChange(Math.max(1, currentPage - 1))}
              disabled={currentPage === 1 || loading}
              className="pagination-button"
            >
              ← Anterior
            </button>
            <span className="pagination-info">
              Página {currentPage} de {totalPages} ({totalItems} registros totales)
            </span>
            <button
              onClick={() => onPageChange(Math.min(totalPages, currentPage + 1))}
              disabled={currentPage === totalPages || loading}
              className="pagination-button"
            >
              Siguiente →
            </button>
          </div>
          {/* Selector de elementos por página */}
          {showItemsPerPageSelector && onItemsPerPageChange && (
            <div className="flex items-center gap-2">
              <label className="text-sm font-medium text-gray-700">Mostrar:</label>
              <select
                value={itemsPerPage}
                onChange={(e) => {
                  onItemsPerPageChange(parseInt(e.target.value));
                  onPageChange(1); // Reset a página 1 cuando cambia items per page
                }}
                className="input text-sm px-3 py-2"
                disabled={loading}
              >
                {itemsPerPageOptions.map((option) => (
                  <option key={option} value={option}>
                    {option}
                  </option>
                ))}
              </select>
              <span className="text-sm text-gray-700">por página</span>
            </div>
          )}
        </div>
      </div>
    );
  };

  return (
    <div className="w-full">
      {/* Filtros dinámicos */}
      {filters && (
        <div className="mb-0">
          <DynamicFilters
            availableFilters={filters.availableFilters}
            activeFilters={filters.activeFilters}
            onAddFilter={filters.onAddFilter}
            onRemoveFilter={filters.onRemoveFilter}
            sortBy={filters.sortBy || 'created_at'}
            sortOrder={filters.sortOrder || 'desc'}
            onSortChange={filters.onSortChange}
            sortOptions={filters.sortOptions || [
              { value: 'created_at', label: 'Ordenar por fecha' },
              { value: 'updated_at', label: 'Ordenar por actualización' },
              { value: 'total_amount', label: 'Ordenar por monto' },
            ]}
          />
        </div>
      )}

      {/* Tabla */}
      <div className={`card overflow-hidden p-0 w-full ${filters ? 'rounded-t-none border-t-0' : ''}`}>
        <div className="overflow-x-auto w-full">
          <table className="table w-full">
          {/* Header */}
          <thead>
            <tr>
              {columns.map((column) => (
                <th
                  key={column.key}
                  className={alignClass[column.align || 'left']}
                  style={{ width: column.width }}
                >
                  {column.label}
                </th>
              ))}
            </tr>
          </thead>

          {/* Body */}
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={columns.length} className="px-6 py-12 text-center text-gray-500">
                  <div className="flex justify-center items-center gap-3">
                    <div className="spinner"></div>
                    <span>Cargando...</span>
                  </div>
                </td>
              </tr>
            ) : data.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="px-6 py-12 text-center text-gray-500">
                  {emptyMessage}
                </td>
              </tr>
            ) : (
              data.map((row, rowIndex) => (
                <tr
                  key={keyExtractor(row, rowIndex)}
                  className={`transition-colors ${onRowClick ? 'cursor-pointer' : ''}`}
                  onClick={() => onRowClick?.(row, rowIndex)}
                >
                  {columns.map((column) => (
                    <td
                      key={column.key}
                      className={alignClass[column.align || 'left']}
                    >
                      {column.render 
                        ? column.render(row[column.key as keyof T], row, rowIndex)
                        : (() => {
                            const value = row[column.key as keyof T];
                            // Si es un elemento React (JSX), lo renderizamos directamente
                            if (React.isValidElement(value)) {
                              return value;
                            }
                            // Si es null/undefined, mostramos string vacío
                            if (value == null) {
                              return '';
                            }
                            // Para otros valores, los convertimos a string
                            return String(value);
                          })()}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
        {/* Paginación integrada */}
        {renderPagination()}
      </div>
    </div>
  );
}

// Re-exportar tipos de filtros para facilitar el uso
export type { FilterOption, ActiveFilter } from './dynamic-filters';

