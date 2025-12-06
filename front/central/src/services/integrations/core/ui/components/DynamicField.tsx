'use client';

import { useState } from 'react';
import { Input, Select, Alert } from '@/shared/ui';

export interface FieldSchema {
    // Tipo y validación básica
    type: 'string' | 'number' | 'boolean';
    required?: boolean;
    default?: any;

    // UI Labels
    label: string;
    placeholder?: string;
    description?: string;
    boldLabel?: boolean;  // Para hacer el label en negrilla

    // Help
    help_text?: string;
    help_link?: string;

    // Input
    input_type?: 'text' | 'password' | 'email' | 'url' | 'number';

    // Validaciones
    pattern?: string;  // Regex
    min?: number;
    max?: number;
    minLength?: number;
    maxLength?: number;
    enum?: Array<{ value: string; label: string }>;

    // Mensajes
    error_message?: string;

    // Ordenamiento
    order?: number;
}

interface DynamicFieldProps {
    name: string;
    schema: FieldSchema;
    value: any;
    onChange: (value: any) => void;
    error?: string;
}

export default function DynamicField({ name, schema, value, onChange, error }: DynamicFieldProps) {
    const [touched, setTouched] = useState(false);
    const [validationError, setValidationError] = useState<string | null>(null);
    const [showPassword, setShowPassword] = useState(false);
    
    // Log para debug
    if (schema.input_type === 'password') {
        console.log(`🔐 Campo password "${name}":`, {
            value: value,
            valueType: typeof value,
            valueLength: typeof value === 'string' ? value.length : 'N/A',
            isEmpty: !value || value === '',
            showPassword
        });
    }

    const handleChange = (newValue: any) => {
        setTouched(true);
        onChange(newValue);

        // Validate
        if (schema.required && !newValue) {
            setValidationError(schema.error_message || 'Este campo es requerido');
            return;
        }

        if (schema.pattern && newValue) {
            const regex = new RegExp(schema.pattern);
            if (!regex.test(newValue)) {
                setValidationError(schema.error_message || 'Formato inválido');
                return;
            }
        }

        if (schema.minLength && newValue && newValue.length < schema.minLength) {
            setValidationError(schema.error_message || `Mínimo ${schema.minLength} caracteres`);
            return;
        }

        if (schema.maxLength && newValue && newValue.length > schema.maxLength) {
            setValidationError(schema.error_message || `Máximo ${schema.maxLength} caracteres`);
            return;
        }

        if (schema.min !== undefined && newValue < schema.min) {
            setValidationError(schema.error_message || `Valor mínimo: ${schema.min}`);
            return;
        }

        if (schema.max !== undefined && newValue > schema.max) {
            setValidationError(schema.error_message || `Valor máximo: ${schema.max}`);
            return;
        }

        setValidationError(null);
    };

    const displayError = touched && (error || validationError);

    return (
        <div className="space-y-1.5">
            <label className={`block text-xs font-bold text-gray-700 flex items-center gap-1`}>
                {schema.label}
                {schema.required && <span className="text-red-500 ml-0.5">*</span>}
                {schema.help_text && (
                    <span className="group relative">
                        <span className="text-gray-400 cursor-help">ⓘ</span>
                        <span className="invisible group-hover:visible absolute left-0 top-6 w-64 bg-gray-900 text-white text-xs rounded py-1 px-2 z-10">
                            {schema.help_text}
                            {schema.help_link && (
                                <>
                                    {' '}
                                    <a
                                        href={schema.help_link}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="text-blue-300 hover:text-blue-200 underline"
                                    >
                                        Ver guía →
                                    </a>
                                </>
                            )}
                        </span>
                    </span>
                )}
            </label>

            {/* Text/Email/URL/Number inputs */}
            {(schema.type === 'string' || schema.type === 'number') && !schema.enum && (
                <div className="relative">
                    <input
                        type={schema.input_type === 'password' && showPassword ? 'text' : (schema.input_type || (schema.type === 'number' ? 'number' : 'text'))}
                        value={value ?? ''}
                        onChange={(e) => handleChange(schema.type === 'number' ? Number(e.target.value) : e.target.value)}
                        onBlur={() => setTouched(true)}
                        required={schema.required}
                        placeholder={schema.placeholder}
                        min={schema.min}
                        max={schema.max}
                        minLength={schema.minLength}
                        maxLength={schema.maxLength}
                        className="w-full px-3 py-2 text-sm text-gray-900 bg-gray-50 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent focus:bg-white"
                    />
                    {schema.input_type === 'password' && (
                        <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 p-1"
                            title={showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
                        >
                            {showPassword ? (
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                                </svg>
                            ) : (
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                </svg>
                            )}
                        </button>
                    )}
                </div>
            )}

            {/* Select/Enum */}
            {schema.enum && (
                <select
                    value={value || ''}
                    onChange={(e) => handleChange(e.target.value)}
                    onBlur={() => setTouched(true)}
                    required={schema.required}
                    className="w-full px-3 py-2 text-sm text-gray-900 bg-gray-50 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent focus:bg-white"
                >
                    <option value="">Selecciona una opción</option>
                    {schema.enum.map((option) => (
                        <option key={option.value} value={option.value}>
                            {option.label}
                        </option>
                    ))}
                </select>
            )}

            {/* Boolean/Checkbox */}
            {schema.type === 'boolean' && (
                <label className="flex items-center space-x-2 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={!!value}
                        onChange={(e) => handleChange(e.target.checked)}
                        className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
                    />
                    <span className="text-sm text-gray-700">{schema.description || schema.label}</span>
                </label>
            )}

            {/* Error message removed as per user request - relying on tooltips */}

            {error && (
                <p className="text-xs text-red-600">
                    {error}
                </p>
            )}
        </div>
    );
}
