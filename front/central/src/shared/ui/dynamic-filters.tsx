'use client';

import { useState, useRef, useEffect } from 'react';
import { FunnelIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { Button } from './button';
import { DateRangePicker } from './date-range-picker';

export interface FilterOption {
    key: string;
    label: string;
    type: 'text' | 'select' | 'date-range' | 'boolean';
    placeholder?: string;
    options?: Array<{ value: string; label: string }>;
}

export interface ActiveFilter {
    key: string;
    label: string;
    value: string | { start?: string; end?: string } | boolean;
    type: 'text' | 'select' | 'date-range' | 'boolean';
}

interface DynamicFiltersProps {
    availableFilters: FilterOption[];
    activeFilters: ActiveFilter[];
    onAddFilter: (filterKey: string, value: any) => void;
    onRemoveFilter: (filterKey: string) => void;
    sortBy?: string;
    sortOrder?: 'asc' | 'desc';
    onSortChange?: (sortBy: string, sortOrder: 'asc' | 'desc') => void;
    sortOptions?: Array<{ value: string; label: string }>;
    className?: string;
}

export function DynamicFilters({
    availableFilters,
    activeFilters,
    onAddFilter,
    onRemoveFilter,
    sortBy = 'created_at',
    sortOrder = 'desc',
    onSortChange,
    sortOptions = [
        { value: 'created_at', label: 'Ordenar por fecha' },
        { value: 'updated_at', label: 'Ordenar por actualización' },
        { value: 'total_amount', label: 'Ordenar por monto' },
    ],
    className = '',
}: DynamicFiltersProps) {
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [selectedFilterKey, setSelectedFilterKey] = useState<string | null>(null);
    const [tempValue, setTempValue] = useState<any>('');
    const dropdownRef = useRef<HTMLDivElement>(null);
    const buttonRef = useRef<HTMLButtonElement>(null);

    // Cerrar dropdown al hacer clic fuera
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (
                dropdownRef.current &&
                !dropdownRef.current.contains(event.target as Node) &&
                buttonRef.current &&
                !buttonRef.current.contains(event.target as Node)
            ) {
                setIsDropdownOpen(false);
                setSelectedFilterKey(null);
                setTempValue('');
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    // Obtener filtros disponibles (que no estén ya activos)
    const availableFilterOptions = availableFilters.filter((filter) => {
        // Para date-range, verificar si ya hay un filtro de fecha activo
        if (filter.type === 'date-range') {
            return !activeFilters.some((active) => active.type === 'date-range');
        }
        return !activeFilters.some((active) => active.key === filter.key);
    });

    const handleFilterSelect = (filterKey: string) => {
        const filter = availableFilters.find((f) => f.key === filterKey);
        if (!filter) return;

        setSelectedFilterKey(filterKey);
        setTempValue('');

        // Si es boolean, aplicar directamente
        if (filter.type === 'boolean') {
            onAddFilter(filterKey, true);
            setIsDropdownOpen(false);
            setSelectedFilterKey(null);
        }
    };

    const handleApplyFilter = () => {
        if (!selectedFilterKey || tempValue === '') return;

        const filter = availableFilters.find((f) => f.key === selectedFilterKey);
        if (!filter) return;

        onAddFilter(selectedFilterKey, tempValue);
        setTempValue('');
        setSelectedFilterKey(null);
        setIsDropdownOpen(false);
    };

    const handleDateRangeChange = (start: string | undefined, end: string | undefined) => {
        if (!selectedFilterKey) return;
        // Permitir aplicar incluso si solo hay una fecha (el usuario puede completar después)
        if (start || end) {
            onAddFilter(selectedFilterKey, { start, end });
            // Solo cerrar si ambas fechas están completas
            if (start && end) {
                setSelectedFilterKey(null);
                setIsDropdownOpen(false);
            }
        }
    };

    const getFilterDisplayValue = (filter: ActiveFilter): string => {
        if (filter.type === 'date-range' && typeof filter.value === 'object') {
            const range = filter.value as { start?: string; end?: string };
            if (range.start && range.end) {
                const startDate = new Date(range.start).toLocaleDateString('es-CO', {
                    month: 'short',
                    day: 'numeric',
                });
                const endDate = new Date(range.end).toLocaleDateString('es-CO', {
                    month: 'short',
                    day: 'numeric',
                });
                return `${startDate} - ${endDate}`;
            }
            return 'Rango de fechas';
        }
        if (filter.type === 'boolean') {
            return filter.value === true ? 'Sí' : 'No';
        }
        if (filter.type === 'select') {
            const filterOption = availableFilters.find((f) => f.key === filter.key);
            const option = filterOption?.options?.find(
                (opt) => opt.value === filter.value.toString()
            );
            return option?.label || filter.value.toString();
        }
        return filter.value.toString();
    };

    const getFilterLabel = (key: string): string => {
        const filter = availableFilters.find((f) => f.key === key);
        return filter?.label || key;
    };

    return (
        <div className={`bg-white p-4 sm:p-6 rounded-t-lg rounded-b-none shadow-sm border border-gray-200 border-b-0 ${className}`}>
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                {/* Botón Añadir Filtro y Chips */}
                <div className="flex-1 flex flex-wrap items-center gap-2">
                    <div className="relative">
                        <Button
                            ref={buttonRef}
                            variant="primary"
                            size="sm"
                            onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                            className="flex items-center gap-2"
                        >
                            <FunnelIcon className="w-4 h-4" />
                            Añadir Filtro
                        </Button>

                        {/* Dropdown de filtros */}
                        {isDropdownOpen && (
                            <div
                                ref={dropdownRef}
                                className="absolute top-full left-0 mt-2 w-64 bg-white border border-gray-200 rounded-lg shadow-lg z-50 p-3"
                            >
                                {selectedFilterKey ? (
                                    // Mostrar input para el filtro seleccionado
                                    (() => {
                                        const filter = availableFilters.find(
                                            (f) => f.key === selectedFilterKey
                                        );
                                        if (!filter) return null;

                                        if (filter.type === 'text') {
                                            return (
                                                <div className="space-y-3">
                                                    <label className="block text-sm font-medium text-gray-700">
                                                        {filter.label}
                                                    </label>
                                                    <input
                                                        type="text"
                                                        value={tempValue}
                                                        onChange={(e) => setTempValue(e.target.value)}
                                                        placeholder={filter.placeholder}
                                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900"
                                                        autoFocus
                                                        onKeyDown={(e) => {
                                                            if (e.key === 'Enter') {
                                                                handleApplyFilter();
                                                            }
                                                        }}
                                                    />
                                                    <div className="flex gap-2">
                                                        <Button
                                                            variant="primary"
                                                            size="sm"
                                                            onClick={handleApplyFilter}
                                                            disabled={!tempValue}
                                                            className="flex-1"
                                                        >
                                                            Aplicar
                                                        </Button>
                                                        <Button
                                                            variant="outline"
                                                            size="sm"
                                                            onClick={() => {
                                                                setSelectedFilterKey(null);
                                                                setTempValue('');
                                                            }}
                                                            className="flex-1"
                                                        >
                                                            Cancelar
                                                        </Button>
                                                    </div>
                                                </div>
                                            );
                                        }

                                        if (filter.type === 'select') {
                                            return (
                                                <div className="space-y-3">
                                                    <label className="block text-sm font-medium text-gray-700">
                                                        {filter.label}
                                                    </label>
                                                    <select
                                                        value={tempValue}
                                                        onChange={(e) => setTempValue(e.target.value)}
                                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                                                        autoFocus
                                                    >
                                                        <option value="">Seleccionar...</option>
                                                        {filter.options?.map((opt) => (
                                                            <option key={opt.value} value={opt.value}>
                                                                {opt.label}
                                                            </option>
                                                        ))}
                                                    </select>
                                                    <div className="flex gap-2">
                                                        <Button
                                                            variant="primary"
                                                            size="sm"
                                                            onClick={handleApplyFilter}
                                                            disabled={!tempValue}
                                                            className="flex-1"
                                                        >
                                                            Aplicar
                                                        </Button>
                                                        <Button
                                                            variant="outline"
                                                            size="sm"
                                                            onClick={() => {
                                                                setSelectedFilterKey(null);
                                                                setTempValue('');
                                                            }}
                                                            className="flex-1"
                                                        >
                                                            Cancelar
                                                        </Button>
                                                    </div>
                                                </div>
                                            );
                                        }

                                        if (filter.type === 'date-range') {
                                            return (
                                                <div className="space-y-3">
                                                    <label className="block text-sm font-medium text-gray-700">
                                                        {filter.label}
                                                    </label>
                                                    <DateRangePicker
                                                        startDate={
                                                            typeof tempValue === 'object' && tempValue?.start
                                                                ? tempValue.start
                                                                : undefined
                                                        }
                                                        endDate={
                                                            typeof tempValue === 'object' && tempValue?.end
                                                                ? tempValue.end
                                                                : undefined
                                                        }
                                                        onChange={(start, end) => {
                                                            setTempValue({ start, end });
                                                            if (start && end) {
                                                                handleDateRangeChange(start, end);
                                                            }
                                                        }}
                                                        placeholder="Seleccionar rango"
                                                    />
                                                    <Button
                                                        variant="outline"
                                                        size="sm"
                                                        onClick={() => {
                                                            setSelectedFilterKey(null);
                                                            setTempValue('');
                                                        }}
                                                        className="w-full"
                                                    >
                                                        Cancelar
                                                    </Button>
                                                </div>
                                            );
                                        }

                                        return null;
                                    })()
                                ) : (
                                    // Mostrar lista de filtros disponibles
                                    <div className="space-y-1">
                                        {availableFilterOptions.length === 0 ? (
                                            <p className="text-sm text-gray-500 p-2">
                                                Todos los filtros están aplicados
                                            </p>
                                        ) : (
                                            availableFilterOptions.map((filter) => (
                                                <button
                                                    key={filter.key}
                                                    onClick={() => handleFilterSelect(filter.key)}
                                                    className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded-md transition-colors"
                                                >
                                                    {filter.label}
                                                </button>
                                            ))
                                        )}
                                    </div>
                                )}
                            </div>
                        )}
                    </div>

                    {/* Chips de filtros activos */}
                    {activeFilters.map((filter) => (
                        <div
                            key={filter.key}
                            className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-blue-50 text-blue-800 rounded-full text-sm font-medium"
                        >
                            <span>
                                {getFilterLabel(filter.key)}: {getFilterDisplayValue(filter)}
                            </span>
                            <button
                                onClick={() => onRemoveFilter(filter.key)}
                                className="hover:bg-blue-100 rounded-full p-0.5 transition-colors"
                                aria-label={`Eliminar filtro ${filter.label}`}
                            >
                                <XMarkIcon className="w-4 h-4" />
                            </button>
                        </div>
                    ))}
                </div>

                {/* Selectores de ordenamiento */}
                {onSortChange && (
                    <div className="flex items-center gap-2 flex-shrink-0">
                        <select
                            value={sortBy}
                            onChange={(e) => onSortChange(e.target.value, sortOrder)}
                            className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white text-sm"
                        >
                            {sortOptions.map((opt) => (
                                <option key={opt.value} value={opt.value}>
                                    {opt.label}
                                </option>
                            ))}
                        </select>
                        <select
                            value={sortOrder}
                            onChange={(e) => onSortChange(sortBy, e.target.value as 'asc' | 'desc')}
                            className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white text-sm"
                        >
                            <option value="desc">Descendente</option>
                            <option value="asc">Ascendente</option>
                        </select>
                    </div>
                )}
            </div>
        </div>
    );
}
