'use client';

import { useState, useRef, useEffect } from 'react';
import { DayPicker, DateRange } from 'react-day-picker';
import { es } from 'date-fns/locale';
import { format } from 'date-fns';

interface DateRangePickerProps {
    startDate?: string;
    endDate?: string;
    onChange: (startDate: string | undefined, endDate: string | undefined) => void;
    placeholder?: string;
    className?: string;
    /** Mostrar inputs de hora junto a cada fecha */
    showTime?: boolean;
}

function parseDate(dateString: string | undefined): Date | undefined {
    if (!dateString) return undefined;
    // YYYY-MM-DD
    if (/^\d{4}-\d{2}-\d{2}$/.test(dateString)) {
        const [year, month, day] = dateString.split('-').map(Number);
        return new Date(year, month - 1, day);
    }
    // YYYY-MM-DDTHH:MM
    if (/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/.test(dateString)) {
        const [datePart, timePart] = dateString.split('T');
        const [year, month, day] = datePart.split('-').map(Number);
        const [hours, minutes] = timePart.split(':').map(Number);
        return new Date(year, month - 1, day, hours, minutes);
    }
    const parsed = new Date(dateString);
    return isNaN(parsed.getTime()) ? undefined : parsed;
}

function extractTime(dateString: string | undefined): string {
    if (!dateString) return '00:00';
    const match = dateString.match(/T(\d{2}:\d{2})/);
    return match ? match[1] : '00:00';
}

function formatDateOutput(date: Date, time: string): string {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}T${time}`;
}

function formatDateOnly(date: Date): string {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
}

export function DateRangePicker({
    startDate,
    endDate,
    onChange,
    placeholder = 'Seleccionar rango de fechas',
    className = '',
    showTime = false
}: DateRangePickerProps) {
    const [isOpen, setIsOpen] = useState(false);
    const [tempRange, setTempRange] = useState<DateRange | undefined>(() => {
        const from = parseDate(startDate);
        const to = parseDate(endDate);
        return (from || to) ? { from, to } : undefined;
    });
    const [startTime, setStartTime] = useState(() => extractTime(startDate));
    const [endTime, setEndTime] = useState(() => extractTime(endDate) || '23:59');
    const containerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (isOpen) {
            const from = parseDate(startDate);
            const to = parseDate(endDate);
            setTempRange((from || to) ? { from, to } : undefined);
            setStartTime(extractTime(startDate));
            setEndTime(extractTime(endDate) || '23:59');
        }
    }, [isOpen, startDate, endDate]);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        };

        if (isOpen) {
            document.addEventListener('mousedown', handleClickOutside);
        }

        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [isOpen]);

    const handleSelect = (range: DateRange | undefined) => {
        if (range?.from) {
            range.from = new Date(range.from.getFullYear(), range.from.getMonth(), range.from.getDate());
        }
        if (range?.to) {
            range.to = new Date(range.to.getFullYear(), range.to.getMonth(), range.to.getDate());
        }
        setTempRange(range);
    };

    const handleApply = () => {
        let fromString: string | undefined = undefined;
        let toString: string | undefined = undefined;

        if (tempRange?.from) {
            fromString = showTime
                ? formatDateOutput(tempRange.from, startTime)
                : formatDateOnly(tempRange.from);
        }

        if (tempRange?.to) {
            toString = showTime
                ? formatDateOutput(tempRange.to, endTime)
                : formatDateOnly(tempRange.to);
        }

        onChange(fromString, toString);
        setIsOpen(false);
    };

    const getDisplayText = () => {
        const from = parseDate(startDate);
        const to = parseDate(endDate);

        const fromTimeStr = showTime && startDate?.includes('T') ? ` ${extractTime(startDate)}` : '';
        const toTimeStr = showTime && endDate?.includes('T') ? ` ${extractTime(endDate)}` : '';

        if (from && to) {
            return `${format(from, 'dd/MM/yyyy', { locale: es })}${fromTimeStr} → ${format(to, 'dd/MM/yyyy', { locale: es })}${toTimeStr}`;
        } else if (from) {
            return `Desde: ${format(from, 'dd/MM/yyyy', { locale: es })}${fromTimeStr}`;
        } else if (to) {
            return `Hasta: ${format(to, 'dd/MM/yyyy', { locale: es })}${toTimeStr}`;
        }
        return '';
    };

    const clearDates = () => {
        setTempRange(undefined);
        setStartTime('00:00');
        setEndTime('23:59');
        onChange(undefined, undefined);
        setIsOpen(false);
    };

    return (
        <div ref={containerRef} className={`relative ${className}`}>
            <div className="relative">
                <input
                    type="text"
                    readOnly
                    value={getDisplayText()}
                    placeholder={placeholder}
                    onClick={() => setIsOpen(!isOpen)}
                    className="w-full px-3 py-2 pr-8 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder:text-gray-500 bg-white cursor-pointer text-sm"
                    style={{ textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }}
                />
                <div className="absolute right-2 top-1/2 transform -translate-y-1/2 pointer-events-none">
                    <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                </div>
            </div>

            {isOpen && (
                <div className="absolute z-50 mt-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-xl p-4 w-auto">
                    {/* Indicador de selección con horas */}
                    <div className="mb-3 px-2 py-2 bg-gray-50 dark:bg-gray-700 rounded-md">
                        <div className="flex items-center gap-2 text-sm dark:text-gray-200">
                            <div className="flex items-center gap-1">
                                <span className={`px-2 py-1 rounded text-xs font-medium ${tempRange?.from ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200' : 'text-gray-500 dark:text-gray-400'}`}>
                                    {tempRange?.from ? format(tempRange.from, 'dd/MM/yyyy', { locale: es }) : 'Inicio'}
                                </span>
                                {showTime && tempRange?.from && (
                                    <input
                                        type="time"
                                        value={startTime}
                                        onChange={(e) => setStartTime(e.target.value)}
                                        className="px-1 py-0.5 text-xs border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500 w-[70px]"
                                    />
                                )}
                            </div>
                            <span className="text-gray-400 dark:text-gray-600">→</span>
                            <div className="flex items-center gap-1">
                                <span className={`px-2 py-1 rounded text-xs font-medium ${tempRange?.to ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200' : 'text-gray-500 dark:text-gray-400'}`}>
                                    {tempRange?.to ? format(tempRange.to, 'dd/MM/yyyy', { locale: es }) : 'Fin'}
                                </span>
                                {showTime && tempRange?.to && (
                                    <input
                                        type="time"
                                        value={endTime}
                                        onChange={(e) => setEndTime(e.target.value)}
                                        className="px-1 py-0.5 text-xs border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500 w-[70px]"
                                    />
                                )}
                            </div>
                        </div>
                    </div>

                    <DayPicker
                        mode="range"
                        selected={tempRange}
                        onSelect={handleSelect}
                        locale={es}
                        numberOfMonths={1}
                        className="rounded-lg"
                        classNames={{
                            months: 'flex flex-col sm:flex-row space-y-4 sm:space-x-4 sm:space-y-0',
                            month: 'space-y-4',
                            caption: 'flex justify-center pt-1 relative items-center mb-2',
                            caption_label: 'text-base font-bold text-black dark:text-gray-100',
                            nav: 'space-x-1 flex items-center',
                            nav_button: 'h-8 w-8 bg-transparent p-0 text-black dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors',
                            nav_button_previous: 'absolute left-1',
                            nav_button_next: 'absolute right-1',
                            table: 'w-full border-collapse space-y-1',
                            head_row: 'flex mb-1',
                            head_cell: 'text-black dark:text-gray-200 rounded-md w-10 font-bold text-xs uppercase tracking-wider',
                            row: 'flex w-full mt-1',
                            cell: 'text-center text-sm p-0 relative',
                            day: 'h-10 w-10 p-0 font-medium text-black dark:text-gray-200 rounded-md transition-colors hover:bg-gray-100 dark:hover:bg-gray-700',
                            day_range_start: 'bg-blue-500 text-white hover:bg-blue-600 font-semibold rounded-l-md',
                            day_range_end: 'bg-blue-500 text-white hover:bg-blue-600 font-semibold rounded-r-md',
                            day_selected: 'bg-blue-500 text-white hover:bg-blue-600 font-semibold',
                            day_today: 'bg-gray-100 dark:bg-gray-700 text-black dark:text-gray-100 font-bold border-2 border-blue-500',
                            day_outside: 'text-gray-500 dark:text-gray-500 opacity-60',
                            day_disabled: 'text-gray-300 dark:text-gray-600 opacity-50 cursor-not-allowed',
                            day_range_middle: 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200 font-medium',
                            day_hidden: 'invisible',
                        }}
                    />

                    {/* Botones de acción */}
                    <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700 flex gap-2">
                        <button
                            onClick={clearDates}
                            className="flex-1 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-md transition-colors font-medium border border-gray-300 dark:border-gray-600"
                            type="button"
                        >
                            Limpiar
                        </button>
                        <button
                            onClick={handleApply}
                            className="flex-1 px-4 py-2 text-sm bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors font-medium shadow-sm"
                            type="button"
                        >
                            Aplicar
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}
