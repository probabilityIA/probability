'use client';

import { useState, useRef, useEffect } from 'react';
import { DayPicker, DateRange } from 'react-day-picker';
import { es } from 'date-fns/locale';
import { format } from 'date-fns';
// Los estilos se aplican inline con styled-jsx

interface DateRangePickerProps {
    startDate?: string;
    endDate?: string;
    onChange: (startDate: string | undefined, endDate: string | undefined) => void;
    placeholder?: string;
    className?: string;
}

export function DateRangePicker({ 
    startDate, 
    endDate, 
    onChange, 
    placeholder = 'Seleccionar rango de fechas',
    className = '' 
}: DateRangePickerProps) {
    const [isOpen, setIsOpen] = useState(false);
    // Estado temporal para la selección (no se aplica hasta hacer clic en "Aplicar")
    const [tempRange, setTempRange] = useState<DateRange | undefined>({
        from: startDate ? new Date(startDate) : undefined,
        to: endDate ? new Date(endDate) : undefined,
    });
    const containerRef = useRef<HTMLDivElement>(null);

    // Sincronizar el estado temporal con las props cuando se abre el calendario
    useEffect(() => {
        if (isOpen) {
            setTempRange(
                startDate || endDate
                    ? {
                          from: startDate ? new Date(startDate) : undefined,
                          to: endDate ? new Date(endDate) : undefined,
                      }
                    : undefined
            );
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
        // Solo actualizar el estado temporal, NO aplicar los cambios aún
        setTempRange(range);
        // NO cerrar el calendario, esperar a que el usuario haga clic en "Aplicar"
    };

    const handleApply = () => {
        // Aplicar los cambios solo cuando se hace clic en "Aplicar"
        const fromString = tempRange?.from ? format(tempRange.from, 'yyyy-MM-dd') : undefined;
        const toString = tempRange?.to ? format(tempRange.to, 'yyyy-MM-dd') : undefined;
        onChange(fromString, toString);
        setIsOpen(false);
    };

    const getDisplayText = () => {
        if (startDate && endDate) {
            const from = new Date(startDate);
            const to = new Date(endDate);
            return `${format(from, 'dd/MM/yyyy', { locale: es })} - ${format(to, 'dd/MM/yyyy', { locale: es })}`;
        } else if (startDate) {
            const from = new Date(startDate);
            return `Desde: ${format(from, 'dd/MM/yyyy', { locale: es })}`;
        } else if (endDate) {
            const to = new Date(endDate);
            return `Hasta: ${format(to, 'dd/MM/yyyy', { locale: es })}`;
        }
        return '';
    };

    const clearDates = () => {
        setTempRange(undefined);
        // No aplicar los cambios hasta hacer clic en "Aplicar"
    };

    return (
        <div ref={containerRef} className={`relative ${className}`}>
            <input
                type="text"
                readOnly
                value={getDisplayText()}
                placeholder={placeholder}
                onClick={() => setIsOpen(!isOpen)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder:text-gray-500 bg-white cursor-pointer"
            />
            
            {isOpen && (
                <div className="absolute z-50 mt-2 bg-white border border-gray-200 rounded-lg shadow-xl p-4 w-auto">
                    {/* Indicador de selección */}
                    <div className="mb-3 px-2 py-2 bg-gray-50 rounded-md">
                        <div className="flex items-center gap-2 text-sm">
                            <span className={`px-2 py-1 rounded text-xs font-medium ${tempRange?.from ? 'bg-blue-100 text-blue-700' : 'text-gray-500'}`}>
                                {tempRange?.from ? format(tempRange.from, 'dd/MM/yyyy', { locale: es }) : 'Seleccionar inicio'}
                            </span>
                            <span className="text-gray-400">→</span>
                            <span className={`px-2 py-1 rounded text-xs font-medium ${tempRange?.to ? 'bg-blue-100 text-blue-700' : 'text-gray-500'}`}>
                                {tempRange?.to ? format(tempRange.to, 'dd/MM/yyyy', { locale: es }) : 'Seleccionar fin'}
                            </span>
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
                            caption_label: 'text-base font-bold text-black',
                            nav: 'space-x-1 flex items-center',
                            nav_button: 'h-8 w-8 bg-transparent p-0 text-black hover:bg-gray-100 rounded transition-colors',
                            nav_button_previous: 'absolute left-1',
                            nav_button_next: 'absolute right-1',
                            table: 'w-full border-collapse space-y-1',
                            head_row: 'flex mb-1',
                            head_cell: 'text-black rounded-md w-10 font-bold text-xs uppercase tracking-wider',
                            row: 'flex w-full mt-1',
                            cell: 'text-center text-sm p-0 relative',
                            day: 'h-10 w-10 p-0 font-medium text-black rounded-md transition-colors hover:bg-gray-100',
                            day_range_start: 'bg-blue-500 text-white hover:bg-blue-600 font-semibold rounded-l-md',
                            day_range_end: 'bg-blue-500 text-white hover:bg-blue-600 font-semibold rounded-r-md',
                            day_selected: 'bg-blue-500 text-white hover:bg-blue-600 font-semibold',
                            day_today: 'bg-gray-100 text-black font-bold border-2 border-blue-500',
                            day_outside: 'text-gray-500 opacity-60',
                            day_disabled: 'text-gray-300 opacity-50 cursor-not-allowed',
                            day_range_middle: 'bg-blue-100 text-blue-700 font-medium',
                            day_hidden: 'invisible',
                        }}
                    />
                    
                    {/* Botones de acción */}
                    <div className="mt-4 pt-4 border-t border-gray-200 flex gap-2">
                        <button
                            onClick={clearDates}
                            className="flex-1 px-4 py-2 text-sm text-gray-700 hover:text-gray-900 hover:bg-gray-50 rounded-md transition-colors font-medium border border-gray-300"
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

