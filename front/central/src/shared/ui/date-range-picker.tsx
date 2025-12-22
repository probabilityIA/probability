'use client';

import { useState, useRef, useEffect } from 'react';
import { DayPicker, DateRange } from 'react-day-picker';
import { es } from 'date-fns/locale';
import { format, parse } from 'date-fns';
// Los estilos se aplican inline con styled-jsx

interface DateRangePickerProps {
    startDate?: string;
    endDate?: string;
    onChange: (startDate: string | undefined, endDate: string | undefined) => void;
    placeholder?: string;
    className?: string;
}

// Función helper para parsear fechas sin problemas de zona horaria
function parseDate(dateString: string | undefined): Date | undefined {
    if (!dateString) return undefined;
    // Si viene en formato YYYY-MM-DD, parsearlo directamente para evitar problemas de zona horaria
    if (/^\d{4}-\d{2}-\d{2}$/.test(dateString)) {
        const [year, month, day] = dateString.split('-').map(Number);
        return new Date(year, month - 1, day);
    }
    // Si viene en otro formato, intentar parsearlo
    const parsed = new Date(dateString);
    return isNaN(parsed.getTime()) ? undefined : parsed;
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
    const [tempRange, setTempRange] = useState<DateRange | undefined>(() => {
        const from = parseDate(startDate);
        const to = parseDate(endDate);
        return (from || to) ? { from, to } : undefined;
    });
    const containerRef = useRef<HTMLDivElement>(null);

    // Sincronizar el estado temporal con las props cuando se abre el calendario o cambian las props
    useEffect(() => {
        if (isOpen) {
            const from = parseDate(startDate);
            const to = parseDate(endDate);
            setTempRange((from || to) ? { from, to } : undefined);
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
        // Asegurarse de que las fechas estén en hora local (sin zona horaria)
        if (range?.from) {
            const year = range.from.getFullYear();
            const month = range.from.getMonth();
            const day = range.from.getDate();
            range.from = new Date(year, month, day);
        }
        if (range?.to) {
            const year = range.to.getFullYear();
            const month = range.to.getMonth();
            const day = range.to.getDate();
            range.to = new Date(year, month, day);
        }
        setTempRange(range);
        // NO cerrar el calendario, esperar a que el usuario haga clic en "Aplicar"
    };

    const handleApply = () => {
        // Aplicar los cambios solo cuando se hace clic en "Aplicar"
        // Usar getFullYear, getMonth, getDate para evitar problemas de zona horaria
        let fromString: string | undefined = undefined;
        let toString: string | undefined = undefined;
        
        if (tempRange?.from) {
            const year = tempRange.from.getFullYear();
            const month = String(tempRange.from.getMonth() + 1).padStart(2, '0');
            const day = String(tempRange.from.getDate()).padStart(2, '0');
            fromString = `${year}-${month}-${day}`;
        }
        
        if (tempRange?.to) {
            const year = tempRange.to.getFullYear();
            const month = String(tempRange.to.getMonth() + 1).padStart(2, '0');
            const day = String(tempRange.to.getDate()).padStart(2, '0');
            toString = `${year}-${month}-${day}`;
        }
        
        onChange(fromString, toString);
        setIsOpen(false);
    };

    const getDisplayText = () => {
        const from = parseDate(startDate);
        const to = parseDate(endDate);
        
        if (from && to) {
            // Formato compacto: "21/11/2025 → 12/12/2025"
            const fromStr = format(from, 'dd/MM/yyyy', { locale: es });
            const toStr = format(to, 'dd/MM/yyyy', { locale: es });
            return `${fromStr} → ${toStr}`;
        } else if (from) {
            const fromStr = format(from, 'dd/MM/yyyy', { locale: es });
            return `Desde: ${fromStr}`;
        } else if (to) {
            const toStr = format(to, 'dd/MM/yyyy', { locale: es });
            return `Hasta: ${toStr}`;
        }
        return '';
    };
    
    const getFullDisplayText = () => {
        const from = parseDate(startDate);
        const to = parseDate(endDate);
        
        if (from && to) {
            return `Rango: ${format(from, 'dd/MM/yyyy', { locale: es })} hasta ${format(to, 'dd/MM/yyyy', { locale: es })}`;
        } else if (from) {
            return `Fecha inicio: ${format(from, 'dd/MM/yyyy', { locale: es })} - Selecciona fecha fin`;
        } else if (to) {
            return `Fecha fin: ${format(to, 'dd/MM/yyyy', { locale: es })} - Selecciona fecha inicio`;
        }
        return placeholder;
    };

    const clearDates = () => {
        setTempRange(undefined);
        // No aplicar los cambios hasta hacer clic en "Aplicar"
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
                    title={getFullDisplayText()}
                    className="w-full px-3 py-2 pr-8 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder:text-gray-500 bg-white cursor-pointer text-sm"
                    style={{ 
                        textOverflow: 'ellipsis',
                        overflow: 'hidden',
                        whiteSpace: 'nowrap'
                    }}
                />
                {/* Icono de calendario */}
                <div className="absolute right-2 top-1/2 transform -translate-y-1/2 pointer-events-none">
                    <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                </div>
            </div>
            
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

