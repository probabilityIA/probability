'use client';

import { useState, useRef, useEffect } from 'react';

interface DatePickerProps {
    value?: string;
    onChange: (date: string) => void;
    placeholder?: string;
    className?: string;
}

export function DatePicker({ value, onChange, placeholder = 'Seleccionar fecha', className = '' }: DatePickerProps) {
    const [isOpen, setIsOpen] = useState(false);
    const [currentMonth, setCurrentMonth] = useState(new Date());
    const containerRef = useRef<HTMLDivElement>(null);

    const selectedDate = value ? new Date(value) : null;

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

    const getDaysInMonth = (date: Date) => {
        const year = date.getFullYear();
        const month = date.getMonth();
        const firstDay = new Date(year, month, 1);
        const lastDay = new Date(year, month + 1, 0);
        const daysInMonth = lastDay.getDate();
        const startingDayOfWeek = firstDay.getDay();

        const days: (number | null)[] = [];
        
        // Días vacíos al inicio
        for (let i = 0; i < startingDayOfWeek; i++) {
            days.push(null);
        }
        
        // Días del mes
        for (let i = 1; i <= daysInMonth; i++) {
            days.push(i);
        }

        return days;
    };

    const handleDateSelect = (day: number) => {
        const year = currentMonth.getFullYear();
        const month = currentMonth.getMonth();
        const date = new Date(year, month, day);
        const dateString = date.toISOString().split('T')[0];
        onChange(dateString);
        setIsOpen(false);
    };

    const handlePrevMonth = () => {
        setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1, 1));
    };

    const handleNextMonth = () => {
        setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1, 1));
    };

    const formatDisplayDate = (dateString: string) => {
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleDateString('es-ES', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric'
        });
    };

    const monthNames = [
        'Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
        'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'
    ];

    const weekDays = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb'];

    const days = getDaysInMonth(currentMonth);
    const today = new Date();
    const isToday = (day: number) => {
        return day === today.getDate() &&
               currentMonth.getMonth() === today.getMonth() &&
               currentMonth.getFullYear() === today.getFullYear();
    };

    const isSelected = (day: number) => {
        if (!selectedDate) return false;
        return day === selectedDate.getDate() &&
               currentMonth.getMonth() === selectedDate.getMonth() &&
               currentMonth.getFullYear() === selectedDate.getFullYear();
    };

    return (
        <div ref={containerRef} className={`relative ${className}`}>
            <input
                type="text"
                readOnly
                value={value ? formatDisplayDate(value) : ''}
                placeholder={placeholder}
                onClick={() => setIsOpen(!isOpen)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder:text-gray-500 bg-white cursor-pointer"
            />
            
            {isOpen && (
                <div className="absolute z-50 mt-2 bg-white border border-gray-200 rounded-lg shadow-lg p-4 w-72">
                    {/* Header del calendario */}
                    <div className="flex items-center justify-between mb-4">
                        <button
                            onClick={handlePrevMonth}
                            className="p-1 hover:bg-gray-100 rounded transition-colors"
                            type="button"
                        >
                            <svg className="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                            </svg>
                        </button>
                        <h3 className="text-sm font-semibold text-gray-900">
                            {monthNames[currentMonth.getMonth()]} {currentMonth.getFullYear()}
                        </h3>
                        <button
                            onClick={handleNextMonth}
                            className="p-1 hover:bg-gray-100 rounded transition-colors"
                            type="button"
                        >
                            <svg className="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                            </svg>
                        </button>
                    </div>

                    {/* Días de la semana */}
                    <div className="grid grid-cols-7 gap-1 mb-2">
                        {weekDays.map((day) => (
                            <div key={day} className="text-center text-xs font-medium text-gray-500 py-1">
                                {day}
                            </div>
                        ))}
                    </div>

                    {/* Días del mes */}
                    <div className="grid grid-cols-7 gap-1">
                        {days.map((day, index) => {
                            if (day === null) {
                                return <div key={`empty-${index}`} className="h-8" />;
                            }

                            const dayIsToday = isToday(day);
                            const dayIsSelected = isSelected(day);

                            return (
                                <button
                                    key={day}
                                    onClick={() => handleDateSelect(day)}
                                    className={`
                                        h-8 w-8 text-sm rounded transition-colors
                                        ${dayIsSelected
                                            ? 'bg-blue-500 text-white font-semibold'
                                            : dayIsToday
                                            ? 'bg-blue-100 text-blue-700 font-semibold'
                                            : 'text-gray-700 hover:bg-gray-100'
                                        }
                                    `}
                                    type="button"
                                >
                                    {day}
                                </button>
                            );
                        })}
                    </div>

                    {/* Botón para limpiar */}
                    {value && (
                        <div className="mt-3 pt-3 border-t border-gray-200">
                            <button
                                onClick={() => {
                                    onChange('');
                                    setIsOpen(false);
                                }}
                                className="w-full px-3 py-1.5 text-sm text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded transition-colors"
                                type="button"
                            >
                                Limpiar
                            </button>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}

