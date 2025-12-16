import React, { useEffect, useState, useRef } from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Select } from '@/shared/ui/select';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { AvatarUpload } from '@/shared/ui/avatar-upload';
import { User } from '../../domain/types';
import { useUserForm } from '../hooks/useUserForm';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';

// Opciones de scope - IDs deben coincidir con la base de datos
const SCOPE_OPTIONS = [
    { value: '2', label: 'Business - Usuario de negocio' },
    { value: '1', label: 'Platform - Super administrador' },
];

interface Business {
    id: number;
    name: string;
}

interface MultiSelectOption {
    value: number;
    label: string;
}

interface MultiSelectProps {
    options: MultiSelectOption[];
    selectedValues: number[];
    onChange: (values: number[]) => void;
    placeholder?: string;
    disabled?: boolean;
}

// Componente MultiSelect simple con checkboxes
const MultiSelect: React.FC<MultiSelectProps> = ({ options, selectedValues, onChange, placeholder, disabled }) => {
    const [isOpen, setIsOpen] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
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

    const handleToggle = (value: number) => {
        if (disabled) return;
        
        if (selectedValues.includes(value)) {
            onChange(selectedValues.filter(v => v !== value));
        } else {
            onChange([...selectedValues, value]);
        }
    };

    const selectedLabels = options
        .filter(opt => selectedValues.includes(opt.value))
        .map(opt => opt.label)
        .join(', ');

    return (
        <div className="relative" ref={dropdownRef}>
            <button
                type="button"
                onClick={() => !disabled && setIsOpen(!isOpen)}
                disabled={disabled}
                className={`w-full input text-left flex items-center justify-between ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
            >
                <span className={selectedLabels ? 'text-gray-900' : 'text-gray-400'}>
                    {selectedLabels || placeholder || 'Seleccionar...'}
                </span>
                <svg 
                    className={`w-5 h-5 text-gray-400 transition-transform ${isOpen ? 'transform rotate-180' : ''}`}
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
            </button>

            {isOpen && !disabled && (
                <div className="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto">
                    {options.length === 0 ? (
                        <div className="px-4 py-2 text-sm text-gray-500">No hay opciones disponibles</div>
                    ) : (
                        options.map((option) => (
                            <label
                                key={option.value}
                                className="flex items-center px-4 py-2 hover:bg-gray-50 cursor-pointer"
                            >
                                <input
                                    type="checkbox"
                                    checked={selectedValues.includes(option.value)}
                                    onChange={() => handleToggle(option.value)}
                                    className="mr-2 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                />
                                <span className="text-sm text-gray-900">{option.label}</span>
                            </label>
                        ))
                    )}
                </div>
            )}
        </div>
    );
};

interface UserFormProps {
    initialData?: User;
    onSuccess: () => void;
    onCancel: () => void;
}

export const UserForm: React.FC<UserFormProps> = ({ initialData, onSuccess, onCancel }) => {
    const {
        formData,
        loading,
        error,
        successMessage,
        handleChange,
        handleFileChange,
        submit,
        setError
    } = useUserForm(initialData, onSuccess);

    const [businesses, setBusinesses] = useState<Business[]>([]);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);
    const [currentUser, setCurrentUser] = useState<{ is_super_admin?: boolean; scope?: string; business_id?: number } | null>(null);
    const [selectedScope, setSelectedScope] = useState<string>('2'); // Default: business

    // Obtener información del usuario actual
    useEffect(() => {
        const userData = TokenStorage.getUser();
        const businessesData = TokenStorage.getBusinessesData();
        
        if (userData) {
            // Obtener el business_id del primer negocio si es usuario de negocio
            let businessId: number | undefined;
            if (businessesData && businessesData.length > 0) {
                businessId = businessesData[0].id;
            }
            
            setCurrentUser({
                is_super_admin: userData.is_super_admin || false,
                scope: userData.scope || '',
                business_id: businessId
            });
        }
    }, []);

    // Inicializar scope_id al crear usuario (default: business = 2)
    useEffect(() => {
        if (!initialData && currentUser?.is_super_admin) {
            handleChange('scope_id', 2); // Default: business
        }
    }, [initialData, currentUser]);

    // Manejar cambio de scope
    const handleScopeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const scopeId = parseInt(e.target.value);
        setSelectedScope(e.target.value);
        handleChange('scope_id', scopeId);
        
        // Si es platform (1), limpiar business_ids
        if (scopeId === 1) {
            handleChange('business_ids', []);
        }
    };

    // Cargar negocios disponibles
    useEffect(() => {
        const loadBusinesses = async () => {
            if (!currentUser) return;

            // Si es usuario de negocio (no super admin), no cargar negocios
            // Se asignará automáticamente su negocio
            if (currentUser.scope === 'business' && !currentUser.is_super_admin) {
                return;
            }

            setLoadingBusinesses(true);
            try {
                const response = await getBusinessesAction({ page: 1, per_page: 1000 });
                if (response.success && response.data) {
                    setBusinesses(response.data.map((b: any) => ({ id: b.id, name: b.name })));
                }
            } catch (err: any) {
                console.error('Error loading businesses:', err);
            } finally {
                setLoadingBusinesses(false);
            }
        };

        loadBusinesses();
    }, [currentUser]);

    // Asignar automáticamente el negocio del usuario si es usuario de negocio
    useEffect(() => {
        if (currentUser && currentUser.scope === 'business' && !currentUser.is_super_admin && currentUser.business_id) {
            // Si es usuario de negocio, asignar automáticamente su negocio
            if (!initialData) {
                // Solo al crear, no al editar
                handleChange('business_ids', [currentUser.business_id]);
                handleChange('scope_id', 2); // business scope
            }
        }
    }, [currentUser, initialData]);

    // Debug: Log para verificar que initialData tiene avatar_url
    useEffect(() => {
        if (initialData) {
            console.log('UserForm - initialData:', initialData);
            console.log('UserForm - avatar_url:', initialData.avatar_url);
        }
    }, [initialData]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        await submit();
    };

    if (successMessage) {
        return (
            <div className="space-y-4">
                <Alert type="success">{successMessage}</Alert>
                <div className="flex justify-end">
                    <Button onClick={onSuccess}>Close</Button>
                </div>
            </div>
        );
    }

    // Determinar si se debe mostrar el selector de negocios
    // Solo super admins pueden ver y seleccionar scope/negocios
    const isSuperAdmin = currentUser?.is_super_admin === true;
    const isBusinessUser = currentUser?.scope === 'business' && !currentUser?.is_super_admin;
    const isCreating = !initialData;
    
    // Mostrar selector de negocios si:
    // - Es super admin Y
    // - Es crear usuario (no editar) Y
    // - El scope seleccionado es "business" (2)
    const showBusinessSelector = isSuperAdmin && isCreating && selectedScope === '2';
    
    // Mostrar selector de scope solo para super admin al crear
    const showScopeSelector = isSuperAdmin && isCreating;

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <Input
                label="Nombre"
                value={formData.name || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('name', e.target.value)}
                required
            />
            <Input
                label="Email"
                type="email"
                value={formData.email || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('email', e.target.value)}
                required
            />
            <Input
                label="Teléfono"
                value={formData.phone || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('phone', e.target.value)}
            />

            <div className="flex flex-col items-center gap-4 py-4">
                <AvatarUpload
                    key={`avatar-${initialData?.id || 'new'}-${initialData?.avatar_url || 'none'}`}
                    currentAvatarUrl={initialData?.avatar_url || null}
                    onFileSelect={handleFileChange}
                    size="lg"
                />
            </div>

            <div className="flex items-center gap-2 mt-4">
                <input
                    type="checkbox"
                    checked={formData.is_active}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('is_active', e.target.checked)}
                    id="is_active"
                />
                <label htmlFor="is_active">Activo</label>
            </div>

            {/* Selector de Scope - Solo para super admins al CREAR */}
            {showScopeSelector && (
                <div className="mt-4">
                    <Select
                        label="Tipo de Usuario *"
                        value={selectedScope}
                        onChange={handleScopeChange}
                        options={SCOPE_OPTIONS}
                    />
                    <p className="text-xs text-gray-500 mt-1">
                        {selectedScope === '1' 
                            ? 'Super administrador con acceso completo a la plataforma.' 
                            : 'Usuario que pertenece a uno o más negocios.'}
                    </p>
                </div>
            )}

            {/* Mostrar scope actual en edición (no editable) */}
            {!isCreating && initialData?.scope_code && (
                <div className="mt-4 p-3 bg-gray-50 border border-gray-200 rounded-md">
                    <p className="text-sm text-gray-700">
                        <span className="font-medium">Tipo de Usuario:</span>{' '}
                        {initialData.scope_code === 'platform' ? 'Platform (Super Admin)' : 'Business (Usuario de Negocio)'}
                    </p>
                    <p className="text-xs text-gray-500 mt-1">
                        El tipo de usuario no puede ser modificado después de la creación.
                    </p>
                </div>
            )}

            {/* Selector de Negocios - Solo para super admins al crear con scope business */}
            {showBusinessSelector && (
                <div className="mt-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Negocios *
                    </label>
                    {loadingBusinesses ? (
                        <div className="flex items-center gap-2 py-2">
                            <Spinner size="sm" />
                            <span className="text-sm text-gray-500">Cargando negocios...</span>
                        </div>
                    ) : (
                        <MultiSelect
                            options={businesses.map(b => ({ value: b.id, label: b.name }))}
                            selectedValues={formData.business_ids || []}
                            onChange={(ids) => handleChange('business_ids', ids)}
                            placeholder="Seleccionar negocios..."
                        />
                    )}
                    <p className="text-xs text-gray-500 mt-1">
                        Seleccione uno o más negocios para asignar al usuario.
                    </p>
                </div>
            )}

            {/* Mensaje para usuarios platform (no necesitan negocios) */}
            {showScopeSelector && selectedScope === '1' && (
                <div className="mt-4 p-3 bg-purple-50 border border-purple-200 rounded-md">
                    <p className="text-xs text-purple-700">
                        Los usuarios Platform (Super Admin) no requieren asignación de negocios.
                    </p>
                </div>
            )}

            {/* Mensaje para usuarios de negocio (el campo no se muestra, se asigna automáticamente) */}
            {isBusinessUser && currentUser?.business_id && isCreating && (
                <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
                    <p className="text-xs text-blue-700">
                        El usuario será asignado automáticamente a su negocio.
                    </p>
                </div>
            )}

            <div className="flex justify-end gap-2 mt-6">
                <Button type="button" variant="secondary" onClick={onCancel}>Cancelar</Button>
                <Button type="submit" disabled={loading}>
                    {loading ? <Spinner size="sm" /> : (initialData ? 'Actualizar' : 'Crear')}
                </Button>
            </div>
        </form>
    );
};
