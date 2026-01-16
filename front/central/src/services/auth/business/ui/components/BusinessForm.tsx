import React, { useEffect, useState, useRef, useCallback } from 'react';
import { HexColorPicker, HexColorInput } from 'react-colorful';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { FileInput } from '@/shared/ui/file-input';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { Business } from '../../domain/types';
import { useBusinessForm } from '../hooks/useBusinessForm';

// Colores rápidos predeterminados
const PRESET_COLORS = [
    '#000000', '#374151', '#6B7280', '#FFFFFF',
    '#EF4444', '#F97316', '#EAB308', '#22C55E',
    '#06B6D4', '#3B82F6', '#8B5CF6', '#EC4899',
];

// Paletas de colores predefinidas (primario, secundario, terciario, cuaternario)
const COLOR_PALETTES = [
    {
        name: 'Corporativo',
        colors: { primary: '#1E3A5F', secondary: '#FFFFFF', tertiary: '#3B82F6', quaternary: '#E5E7EB' },
    },
    {
        name: 'Moderno',
        colors: { primary: '#111827', secondary: '#F9FAFB', tertiary: '#6366F1', quaternary: '#E0E7FF' },
    },
    {
        name: 'Natural',
        colors: { primary: '#166534', secondary: '#FFFFFF', tertiary: '#22C55E', quaternary: '#DCFCE7' },
    },
    {
        name: 'Elegante',
        colors: { primary: '#1F2937', secondary: '#F3F4F6', tertiary: '#9333EA', quaternary: '#F3E8FF' },
    },
    {
        name: 'Cálido',
        colors: { primary: '#92400E', secondary: '#FFFBEB', tertiary: '#F59E0B', quaternary: '#FEF3C7' },
    },
    {
        name: 'Energético',
        colors: { primary: '#DC2626', secondary: '#FFFFFF', tertiary: '#F97316', quaternary: '#FEE2E2' },
    },
    {
        name: 'Oceánico',
        colors: { primary: '#0E7490', secondary: '#ECFEFF', tertiary: '#06B6D4', quaternary: '#CFFAFE' },
    },
    {
        name: 'Minimalista',
        colors: { primary: '#000000', secondary: '#FFFFFF', tertiary: '#737373', quaternary: '#F5F5F5' },
    },
    {
        name: 'Rosado',
        colors: { primary: '#BE185D', secondary: '#FDF2F8', tertiary: '#EC4899', quaternary: '#FCE7F3' },
    },
    {
        name: 'Tech',
        colors: { primary: '#7C3AED', secondary: '#0F172A', tertiary: '#A78BFA', quaternary: '#1E293B' },
    },
];

interface ColorPickerProps {
    label: string;
    value: string;
    onChange: (color: string) => void;
}

const ColorPicker: React.FC<ColorPickerProps> = ({ label, value, onChange }) => {
    const [isOpen, setIsOpen] = useState(false);
    const popoverRef = useRef<HTMLDivElement>(null);

    // Cerrar al hacer clic fuera
    const handleClickOutside = useCallback((e: MouseEvent) => {
        if (popoverRef.current && !popoverRef.current.contains(e.target as Node)) {
            setIsOpen(false);
        }
    }, []);

    useEffect(() => {
        if (isOpen) {
            document.addEventListener('mousedown', handleClickOutside);
        }
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, [isOpen, handleClickOutside]);

    return (
        <div className="flex flex-col gap-1">
            <label className="text-sm font-medium text-gray-700">{label}</label>
            <div className="relative" ref={popoverRef}>
                <button
                    type="button"
                    onClick={() => setIsOpen(!isOpen)}
                    className="flex items-center gap-2 p-2 border rounded-lg hover:border-gray-400 transition-colors w-full bg-white"
                >
                    <div
                        className="w-8 h-8 rounded-md border border-gray-300 shadow-inner"
                        style={{ backgroundColor: value }}
                    />
                    <span className="text-sm text-gray-600 font-mono uppercase">{value}</span>
                </button>

                {isOpen && (
                    <div className="absolute z-50 mt-2 p-3 bg-white border rounded-xl shadow-xl">
                        {/* Color picker principal */}
                        <div className="color-picker-wrapper">
                            <HexColorPicker color={value} onChange={onChange} />
                        </div>

                        {/* Input hexadecimal */}
                        <div className="mt-3 flex items-center gap-2">
                            <span className="text-xs text-gray-500 font-medium">HEX:</span>
                            <HexColorInput
                                color={value}
                                onChange={onChange}
                                prefixed
                                className="flex-1 px-2 py-1 text-sm font-mono border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 uppercase"
                            />
                        </div>

                        {/* Colores rápidos */}
                        <div className="mt-3 pt-3 border-t">
                            <span className="text-xs text-gray-500 font-medium">Colores rápidos:</span>
                            <div className="grid grid-cols-6 gap-1.5 mt-2">
                                {PRESET_COLORS.map((color) => (
                                    <button
                                        key={color}
                                        type="button"
                                        onClick={() => onChange(color)}
                                        className={`w-6 h-6 rounded-md border-2 transition-all hover:scale-110 ${
                                            value.toLowerCase() === color.toLowerCase()
                                                ? 'border-blue-500 ring-2 ring-blue-200'
                                                : 'border-gray-200 hover:border-gray-300'
                                        }`}
                                        style={{ backgroundColor: color }}
                                        title={color}
                                    />
                                ))}
                            </div>
                        </div>
                    </div>
                )}
            </div>

            <style jsx>{`
                .color-picker-wrapper :global(.react-colorful) {
                    width: 200px;
                    height: 160px;
                }
                .color-picker-wrapper :global(.react-colorful__saturation) {
                    border-radius: 8px 8px 0 0;
                }
                .color-picker-wrapper :global(.react-colorful__hue) {
                    height: 14px;
                    border-radius: 0 0 8px 8px;
                }
                .color-picker-wrapper :global(.react-colorful__saturation-pointer),
                .color-picker-wrapper :global(.react-colorful__hue-pointer) {
                    width: 20px;
                    height: 20px;
                    border-width: 3px;
                }
            `}</style>
        </div>
    );
};

interface BusinessFormProps {
    initialData?: Business;
    onSuccess: () => void;
    onCancel: () => void;
}

export const BusinessForm: React.FC<BusinessFormProps> = ({ initialData, onSuccess, onCancel }) => {
    const {
        formData,
        loading,
        error,
        handleChange,
        submit,
        setError
    } = useBusinessForm(initialData, onSuccess);

    const [logoPreview, setLogoPreview] = useState<string | null>(null);
    const [navbarPreview, setNavbarPreview] = useState<string | null>(null);

    // Limpiar previews cuando se desmonta el componente
    useEffect(() => {
        return () => {
            if (logoPreview && logoPreview.startsWith('blob:')) {
                URL.revokeObjectURL(logoPreview);
            }
            if (navbarPreview && navbarPreview.startsWith('blob:')) {
                URL.revokeObjectURL(navbarPreview);
            }
        };
    }, [logoPreview, navbarPreview]);

    // Actualizar previews cuando cambian los archivos o hay datos iniciales
    useEffect(() => {
        if (formData.logo_file) {
            const url = URL.createObjectURL(formData.logo_file);
            // Usar setTimeout para evitar setState síncrono en efecto
            setTimeout(() => {
                setLogoPreview(url);
            }, 0);
            return () => URL.revokeObjectURL(url);
        } else if (initialData?.logo_url) {
            setTimeout(() => {
                setLogoPreview(initialData.logo_url);
            }, 0);
        } else {
            setTimeout(() => {
                setLogoPreview(null);
            }, 0);
        }
    }, [formData.logo_file, initialData?.logo_url]);

    useEffect(() => {
        if (formData.navbar_image_file) {
            const url = URL.createObjectURL(formData.navbar_image_file);
            // Usar setTimeout para evitar setState síncrono en efecto
            setTimeout(() => {
                setNavbarPreview(url);
            }, 0);
            return () => URL.revokeObjectURL(url);
        } else if (initialData?.navbar_image_url) {
            setTimeout(() => {
                setNavbarPreview(initialData.navbar_image_url);
            }, 0);
        } else {
            setTimeout(() => {
                setNavbarPreview(null);
            }, 0);
        }
    }, [formData.navbar_image_file, initialData?.navbar_image_url]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        await submit();
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            {/* Layout principal en dos columnas */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Columna izquierda: Información básica */}
                <div className="space-y-4">
                    <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide border-b pb-2">
                        Información del Negocio
                    </h3>

                    <Input
                        label="Nombre"
                        value={formData.name || ''}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('name', e.target.value)}
                        required
                    />

                    <Input
                        label="Descripción"
                        value={formData.description || ''}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('description', e.target.value)}
                    />

                    <Input
                        label="Dirección"
                        value={formData.address || ''}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('address', e.target.value)}
                    />

                    {/* Imágenes */}
                    <div className="grid grid-cols-2 gap-4 pt-2">
                        <div className="space-y-2">
                            <FileInput
                                label="Logo"
                                accept="image/*"
                                buttonText="Seleccionar"
                                onChange={(file: File | null) => handleChange('logo_file', file)}
                                helperText="JPG, PNG, WEBP"
                            />
                            {logoPreview && (
                                <img
                                    src={logoPreview}
                                    alt="Logo preview"
                                    className="w-full h-20 object-contain border rounded-lg bg-gray-50"
                                />
                            )}
                        </div>
                        <div className="space-y-2">
                            <FileInput
                                label="Navbar"
                                accept="image/*"
                                buttonText="Seleccionar"
                                onChange={(file: File | null) => handleChange('navbar_image_file', file)}
                                helperText="JPG, PNG, WEBP"
                            />
                            {navbarPreview && (
                                <img
                                    src={navbarPreview}
                                    alt="Navbar preview"
                                    className="w-full h-20 object-contain border rounded-lg bg-gray-50"
                                />
                            )}
                        </div>
                    </div>
                </div>

                {/* Columna derecha: Colores */}
                <div className="space-y-4">
                    <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide border-b pb-2">
                        Paleta de Colores
                    </h3>

                    {/* Paletas predefinidas */}
                    <div>
                        <span className="text-xs text-gray-500 font-medium">Combinaciones sugeridas:</span>
                        <div className="grid grid-cols-5 gap-2 mt-2">
                            {COLOR_PALETTES.map((palette) => (
                                <button
                                    key={palette.name}
                                    type="button"
                                    onClick={() => {
                                        handleChange('primary_color', palette.colors.primary);
                                        handleChange('secondary_color', palette.colors.secondary);
                                        handleChange('tertiary_color', palette.colors.tertiary);
                                        handleChange('quaternary_color', palette.colors.quaternary);
                                    }}
                                    className="group flex flex-col items-center p-2 border rounded-lg hover:border-blue-400 hover:bg-blue-50 transition-all"
                                    title={palette.name}
                                >
                                    <div className="flex w-full h-6 rounded overflow-hidden shadow-sm">
                                        <div className="flex-1" style={{ backgroundColor: palette.colors.primary }} />
                                        <div className="flex-1" style={{ backgroundColor: palette.colors.secondary }} />
                                        <div className="flex-1" style={{ backgroundColor: palette.colors.tertiary }} />
                                        <div className="flex-1" style={{ backgroundColor: palette.colors.quaternary }} />
                                    </div>
                                    <span className="text-[10px] text-gray-500 mt-1 group-hover:text-blue-600 truncate w-full text-center">
                                        {palette.name}
                                    </span>
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Personalización individual */}
                    <div className="pt-2">
                        <span className="text-xs text-gray-500 font-medium">Personalizar colores:</span>
                        <div className="grid grid-cols-2 gap-3 mt-2">
                            <ColorPicker
                                label="Primario"
                                value={formData.primary_color || '#000000'}
                                onChange={(color) => handleChange('primary_color', color)}
                            />
                            <ColorPicker
                                label="Secundario"
                                value={formData.secondary_color || '#ffffff'}
                                onChange={(color) => handleChange('secondary_color', color)}
                            />
                            <ColorPicker
                                label="Terciario"
                                value={formData.tertiary_color || '#cccccc'}
                                onChange={(color) => handleChange('tertiary_color', color)}
                            />
                            <ColorPicker
                                label="Cuaternario"
                                value={formData.quaternary_color || '#eeeeee'}
                                onChange={(color) => handleChange('quaternary_color', color)}
                            />
                        </div>
                    </div>

                    {/* Preview de colores */}
                    <div className="p-3 border rounded-lg bg-gray-50">
                        <span className="text-xs text-gray-500 font-medium">Tu paleta actual:</span>
                        <div className="flex gap-1 mt-2 h-8 rounded overflow-hidden shadow-sm">
                            <div
                                className="flex-1 flex items-center justify-center"
                                style={{ backgroundColor: formData.primary_color || '#000000' }}
                            >
                                <span className="text-[8px] font-medium" style={{ color: formData.secondary_color || '#ffffff' }}>P</span>
                            </div>
                            <div
                                className="flex-1 flex items-center justify-center"
                                style={{ backgroundColor: formData.secondary_color || '#ffffff' }}
                            >
                                <span className="text-[8px] font-medium" style={{ color: formData.primary_color || '#000000' }}>S</span>
                            </div>
                            <div
                                className="flex-1 flex items-center justify-center"
                                style={{ backgroundColor: formData.tertiary_color || '#cccccc' }}
                            >
                                <span className="text-[8px] font-medium" style={{ color: formData.primary_color || '#000000' }}>T</span>
                            </div>
                            <div
                                className="flex-1 flex items-center justify-center"
                                style={{ backgroundColor: formData.quaternary_color || '#eeeeee' }}
                            >
                                <span className="text-[8px] font-medium" style={{ color: formData.primary_color || '#000000' }}>C</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Botones de acción */}
            <div className="flex justify-end gap-2 pt-4 border-t">
                <Button type="button" variant="secondary" onClick={onCancel}>Cancelar</Button>
                <Button type="submit" disabled={loading}>
                    {loading ? <Spinner size="sm" /> : (initialData ? 'Actualizar' : 'Crear')}
                </Button>
            </div>
        </form>
    );
};
