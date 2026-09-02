'use client';

import { useState, useRef, useEffect, ReactNode } from 'react';
import { FunnelIcon, XMarkIcon, ArrowDownTrayIcon } from '@heroicons/react/24/outline';
import { Button } from './button';
import { IconActionButton } from './icon-action-button';
import { DateRangePicker } from './date-range-picker';

export interface FilterOption {
    key: string;
    label: string;
    type: 'text' | 'select' | 'date-range' | 'boolean';
    placeholder?: string;
    options?: FilterSelectOption[];
    dot?: string;
    imageUrl?: string;
    avatarOptions?: boolean;
}

export interface FilterSelectOption {
    value: string;
    label: string;
    avatarUrl?: string;
}

export interface ActiveFilter {
    key: string;
    label: string;
    value: string | { start?: string; end?: string } | boolean;
    type: 'text' | 'select' | 'date-range' | 'boolean';
}

function OptionAvatar({ option, small = false }: { option: FilterSelectOption; small?: boolean }) {
    const size = small ? 'h-4 w-4 text-[9px]' : 'h-6 w-6 text-[10px]';
    if (option.avatarUrl) {
        return (
            <img
                src={option.avatarUrl}
                alt=""
                className={`${size} flex-shrink-0 rounded-full object-cover ring-1 ring-gray-200 dark:ring-gray-600`}
            />
        );
    }
    return (
        <span
            className={`${size} flex flex-shrink-0 items-center justify-center rounded-full bg-gray-200 font-semibold text-gray-600 dark:bg-gray-600 dark:text-gray-200`}
        >
            {option.label.trim().charAt(0).toUpperCase() || '?'}
        </span>
    );
}

interface DynamicFiltersProps {
    availableFilters: FilterOption[];
    activeFilters: ActiveFilter[];
    onAddFilter: (filterKey: string, value: any) => void;
    onRemoveFilter: (filterKey: string) => void;
    onCreate?: () => void;
    createButtonText?: string;
    createButtonIconOnly?: boolean;
    createButtonPosition?: 'left' | 'right';
    createButtonAriaLabel?: string;
    onTestGuide?: () => void;
    onDownload?: () => void;
    downloadButtonText?: string;
    downloadTooltip?: string;
    sortBy?: string;
    sortOrder?: 'asc' | 'desc';
    onSortChange?: (sortBy: string, sortOrder: 'asc' | 'desc') => void;
    sortOptions?: Array<{ value: string; label: string }>;
    className?: string;
    variant?: 'card' | 'bar';
    extraActions?: ReactNode;
    triggerClassName?: string;
}

export function DynamicFilters({
    availableFilters,
    activeFilters,
    onAddFilter,
    onRemoveFilter,
    onCreate,
    createButtonText = '+ Crear Orden',
    createButtonIconOnly = false,
    createButtonPosition = 'left',
    createButtonAriaLabel,
    onTestGuide,
    onDownload,
    downloadButtonText = '↓ Descargar \u00d3rdenes',
    downloadTooltip = 'Descargar \u00f3rdenes en Excel',
    sortBy = 'created_at',
    sortOrder = 'desc',
    onSortChange,
    sortOptions = [
        { value: 'created_at', label: 'Ordenar por fecha' },
        { value: 'updated_at', label: 'Ordenar por actualizaci\u00f3n' },
        { value: 'total_amount', label: 'Ordenar por monto' },
    ],
    className = '',
    variant = 'card',
    extraActions,
    triggerClassName = '',
}: DynamicFiltersProps) {
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [selectedFilterKey, setSelectedFilterKey] = useState<string | null>(null);
    const [tempValue, setTempValue] = useState<any>('');
    const dropdownRef = useRef<HTMLDivElement>(null);
    const buttonRef = useRef<HTMLButtonElement>(null);

    const isBar = variant === 'bar';

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

    const availableFilterOptions = availableFilters.filter((filter) => {
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
        if (start || end) {
            onAddFilter(selectedFilterKey, { start, end });
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
            return filter.value === true ? 'Si' : 'No';
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

    const primaryButtonClass = isBar
        ? 'btn btn-primary flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg hover:shadow-md transition-all'
        : 'btn btn-primary flex items-center gap-2 px-4 py-2 text-sm font-semibold rounded-lg hover:shadow-lg hover:scale-105 transition-all';


    const barSelectBase = 'px-2.5 py-1.5 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-700 text-xs';
    const cardSelectClass = 'px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-700 text-sm';
    const sortBySelectClass = isBar ? barSelectBase : cardSelectClass;
    const sortOrderSelectClass = isBar ? barSelectBase : cardSelectClass;
    const sortBySelectStyle = isBar ? { width: '11rem' } : undefined;
    const sortOrderSelectStyle = isBar ? { width: '8.5rem' } : undefined;

    const createButtonLabel = createButtonAriaLabel || 'Crear nuevo';

    const createButton = onCreate ? (
        <button
            style={{ background: 'var(--color-tertiary-500)' }}
            onClick={() => {
                onCreate();
                setIsDropdownOpen(false);
                setSelectedFilterKey(null);
                setTempValue('');
            }}
            className={
                isBar
                    ? 'flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold text-white rounded-lg hover:shadow-md transition-all'
                    : 'flex items-center gap-2 px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all'
            }
            title={createButtonIconOnly ? createButtonLabel : undefined}
            aria-label={createButtonIconOnly ? createButtonLabel : undefined}
        >
            {createButtonIconOnly ? (
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
            ) : (
                createButtonText
            )}
        </button>
    ) : null;

    const chips = activeFilters.map((filter) => {
        const filterDef = availableFilters.find((f) => f.key === filter.key);
        const chipOption = filterDef?.avatarOptions
            ? filterDef.options?.find((opt) => opt.value === filter.value.toString())
            : undefined;
        return (
        <div
            key={filter.key}
            className={
                isBar
                    ? 'chip-business inline-flex items-center gap-1 pl-2.5 pr-1 py-1 rounded-full text-xs font-medium whitespace-nowrap'
                    : 'chip-business inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm font-medium'
            }
        >
            {chipOption && <OptionAvatar option={chipOption} small />}
            <span>
                {getFilterLabel(filter.key)}: {getFilterDisplayValue(filter)}
            </span>
            <button
                onClick={() => onRemoveFilter(filter.key)}
                className="rounded-full p-0.5 transition-colors"
                aria-label={`Eliminar filtro ${filter.label}`}
            >
                <XMarkIcon className={isBar ? 'w-3.5 h-3.5' : 'w-4 h-4'} />
            </button>
        </div>
        );
    });

    return (
        <div
            className={
                isBar
                    ? `w-full ${className}`
                    : `p-4 sm:p-6 rounded-t-lg rounded-b-none shadow-sm border border-b-0 bg-primary-50 dark:bg-gray-900 border-primary-200 dark:border-gray-800 ${className}`
            }
        >
            <div
                className={
                    isBar
                        ? 'flex flex-wrap items-center gap-x-3 gap-y-2'
                        : 'flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4'
                }
            >
                <div className={isBar ? 'flex flex-wrap items-center gap-2' : 'flex-1 flex flex-wrap items-center gap-2'}>
                    <div className="relative">
                        <div className="flex items-center gap-2">
                            {isBar ? (
                                <IconActionButton
                                    ref={buttonRef}
                                    label={`A${'ñ'}adir filtro`}
                                    variant="primary"
                                    onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                                    icon={<FunnelIcon className="w-4 h-4" />}
                                    className={triggerClassName}
                                />
                            ) : (
                                <button
                                    ref={buttonRef}
                                    onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                                    className={primaryButtonClass}
                                >
                                    <FunnelIcon className="w-4 h-4" />
                                    A{'ñ'}adir Filtro
                                </button>
                            )}

                            {onTestGuide && (
                                <button
                                    onClick={() => {
                                        onTestGuide();
                                        setIsDropdownOpen(false);
                                    }}
                                    className={primaryButtonClass}
                                >
                                    Gu{'í'}a Env{'í'}o
                                </button>
                            )}

                            {createButtonPosition === 'left' && createButton}

                            {onDownload && (
                                isBar ? (
                                    <IconActionButton
                                        label={downloadTooltip}
                                        variant="tertiary"
                                        onClick={() => {
                                            onDownload();
                                            setIsDropdownOpen(false);
                                            setSelectedFilterKey(null);
                                            setTempValue('');
                                        }}
                                        icon={<ArrowDownTrayIcon className="w-4 h-4" />}
                                    />
                                ) : (
                                    <button
                                        onClick={() => {
                                            onDownload();
                                            setIsDropdownOpen(false);
                                            setSelectedFilterKey(null);
                                            setTempValue('');
                                        }}
                                        className="btn btn-tertiary flex items-center gap-2 px-4 py-2 text-sm font-semibold text-white rounded-lg"
                                        title={downloadTooltip}
                                    >
                                        {downloadButtonText}
                                    </button>
                                )
                            )}
                        </div>

                        {isDropdownOpen && (
                            <div
                                ref={dropdownRef}
                                className="absolute top-full left-0 mt-2 w-64 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg z-50 p-3"
                            >
                                {selectedFilterKey ? (
                                    (() => {
                                        const filter = availableFilters.find(
                                            (f) => f.key === selectedFilterKey
                                        );
                                        if (!filter) return null;

                                        if (filter.type === 'text') {
                                            return (
                                                <div className="space-y-3">
                                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                                                        {filter.label}
                                                    </label>
                                                    <input
                                                        type="text"
                                                        value={tempValue}
                                                        onChange={(e) => setTempValue(e.target.value)}
                                                        placeholder={filter.placeholder}
                                                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-700"
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

                                        if (filter.type === 'select' && filter.avatarOptions) {
                                            return (
                                                <div className="space-y-3">
                                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                                                        {filter.label}
                                                    </label>
                                                    <div className="max-h-64 space-y-0.5 overflow-y-auto pr-0.5">
                                                        {filter.options?.length ? (
                                                            filter.options.map((opt) => (
                                                                <button
                                                                    key={opt.value}
                                                                    onClick={() => {
                                                                        onAddFilter(selectedFilterKey, opt.value);
                                                                        setTempValue('');
                                                                        setSelectedFilterKey(null);
                                                                        setIsDropdownOpen(false);
                                                                    }}
                                                                    className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
                                                                >
                                                                    <OptionAvatar option={opt} />
                                                                    <span className="min-w-0 truncate">{opt.label}</span>
                                                                </button>
                                                            ))
                                                        ) : (
                                                            <p className="p-2 text-sm text-gray-500 dark:text-gray-400">Sin opciones</p>
                                                        )}
                                                    </div>
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

                                        if (filter.type === 'select') {
                                            return (
                                                <div className="space-y-3">
                                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                                                        {filter.label}
                                                    </label>
                                                    <select
                                                        value={tempValue}
                                                        onChange={(e) => setTempValue(e.target.value)}
                                                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-700"
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
                                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
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
                                    <div className="space-y-1">
                                        {availableFilterOptions.length === 0 ? (
                                            <p className="text-sm text-gray-500 dark:text-gray-400 p-2">
                                                Todos los filtros est{'á'}n aplicados
                                            </p>
                                        ) : (
                                            availableFilterOptions.map((filter) => (
                                                <button
                                                    key={filter.key}
                                                    onClick={() => handleFilterSelect(filter.key)}
                                                    className="w-full flex items-center gap-2 text-left px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition-colors"
                                                >
                                                    {filter.imageUrl ? (
                                                        <img
                                                            src={filter.imageUrl}
                                                            alt=""
                                                            className="h-5 w-5 flex-shrink-0 rounded border border-gray-200 bg-white object-contain p-0.5 dark:border-gray-600"
                                                        />
                                                    ) : filter.dot ? (
                                                        <span
                                                            className="h-2 w-2 flex-shrink-0 rounded-full"
                                                            style={{ backgroundColor: filter.dot }}
                                                        />
                                                    ) : null}
                                                    <span className="min-w-0">{filter.label}</span>
                                                </button>
                                            ))
                                        )}
                                    </div>
                                )}
                            </div>
                        )}
                    </div>

                    {extraActions}

                    {!isBar && chips}
                </div>

                {(isBar || onSortChange || (createButtonPosition === 'right' && createButton)) && (
                <div className={isBar ? 'flex items-center gap-2 ml-auto flex-shrink-0' : 'flex items-center gap-2 flex-shrink-0'}>
                    {isBar && chips.length > 0 && (
                        <div className="flex flex-wrap items-center justify-end gap-1.5 max-w-[38vw]">{chips}</div>
                    )}
                    {onSortChange && (
                        <>
                            {isBar && (
                                <span className="text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap">Ordenar</span>
                            )}
                            <select
                                value={sortBy}
                                onChange={(e) => onSortChange(e.target.value, sortOrder)}
                                className={sortBySelectClass}
                                style={sortBySelectStyle}
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
                                className={sortOrderSelectClass}
                                style={sortOrderSelectStyle}
                            >
                                <option value="desc">Descendente</option>
                                <option value="asc">Ascendente</option>
                            </select>
                        </>
                    )}
                    {createButtonPosition === 'right' && createButton}
                </div>
                )}
            </div>
        </div>
    );
}
