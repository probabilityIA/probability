/**
 * Componente reutilizable de filtros
 * Diseño consistente para todas las tablas
 */

'use client';

import { useState, useEffect } from 'react';
import { 
  FunnelIcon, 
  XMarkIcon
} from '@heroicons/react/24/outline';
import { Button } from './button';
import { Input } from './input';
import { Badge } from './badge';

export interface FilterField {
  key: string;
  label: string;
  type: 'text' | 'number' | 'select' | 'boolean';
  placeholder?: string;
  icon?: React.ReactNode;
  options?: Array<{ value: string; label: string }>;
  min?: number;
  max?: number;
  advanced?: boolean;
}

export interface FiltersProps {
  fields: FilterField[];
  filters: Record<string, any>;
  onFiltersChange: (filters: Record<string, any>) => void;
  onClearFilters: () => void;
  showAdvanced?: boolean;
  className?: string;
}

export function Filters({ 
  fields, 
  filters, 
  onFiltersChange, 
  onClearFilters,
  showAdvanced: defaultShowAdvanced = false,
  className = ''
}: FiltersProps) {
  const [localFilters, setLocalFilters] = useState(filters);
  const [showAdvanced, setShowAdvanced] = useState(defaultShowAdvanced);

  // Sincronizar filtros locales con los props
  useEffect(() => {
    setLocalFilters(filters);
  }, [filters]);

  // Manejar cambios en filtros
  const handleFilterChange = (key: string, value: any) => {
    const newFilters = { ...localFilters, [key]: value };
    setLocalFilters(newFilters);
    onFiltersChange(newFilters);
  };

  // Limpiar filtros
  const handleClearFilters = () => {
    setLocalFilters({});
    onClearFilters();
  };

  // Contar filtros activos
  const activeFiltersCount = Object.values(localFilters).filter(value => 
    value !== undefined && value !== '' && value !== null
  ).length;

  // Separar campos básicos y avanzados
  const basicFields = fields.filter(f => !f.advanced);
  const advancedFields = fields.filter(f => f.advanced);

  const renderField = (field: FilterField, index: number) => {
    const value = localFilters[field.key];

    switch (field.type) {
      case 'text':
      case 'number':
        return (
          <div key={`${field.key}-${index}`}>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {field.icon}
              {field.label}
            </label>
            <Input
              type={field.type}
              value={value || ''}
              onChange={(e) => handleFilterChange(field.key, field.type === 'number' ? (e.target.value ? parseFloat(e.target.value) : undefined) : e.target.value)}
              placeholder={field.placeholder}
              className="input"
              min={field.min}
              max={field.max}
            />
          </div>
        );

      case 'select':
        return (
          <div key={`${field.key}-${index}`}>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {field.icon}
              {field.label}
            </label>
            <select
              value={value?.toString() || ''}
              onChange={(e) => handleFilterChange(field.key, e.target.value === '' ? undefined : e.target.value)}
              className="input"
            >
              {field.options?.map(option => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
        );

      case 'boolean':
        return (
          <div key={`${field.key}-${index}`}>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {field.icon}
              {field.label}
            </label>
            <select
              value={value !== undefined ? value.toString() : ''}
              onChange={(e) => handleFilterChange(field.key, e.target.value === '' ? undefined : e.target.value === 'true')}
              className="input"
            >
              <option value="">Todos</option>
              <option value="true">Sí</option>
              <option value="false">No</option>
            </select>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className={`card ${className}`}>
      <div className="space-y-4">
        {/* Filtros básicos */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {basicFields.map(renderField)}
        </div>

        {/* Filtros avanzados */}
        {showAdvanced && advancedFields.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-4 border-t border-gray-200">
            {advancedFields.map(renderField)}
          </div>
        )}

        {/* Controles */}
        <div className="flex items-center justify-between pt-4 border-t border-gray-200">
          <div className="flex items-center gap-3">
            {/* Botón de filtros avanzados */}
            {advancedFields.length > 0 && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowAdvanced(!showAdvanced)}
                className="btn-outline btn-sm"
              >
                <FunnelIcon className="w-4 h-4 mr-2" />
                {showAdvanced ? 'Ocultar' : 'Mostrar'} filtros avanzados
              </Button>
            )}

            {/* Contador de filtros activos */}
            {activeFiltersCount > 0 && (
              <Badge type="secondary" className="text-xs">
                {activeFiltersCount} filtro{activeFiltersCount !== 1 ? 's' : ''} activo{activeFiltersCount !== 1 ? 's' : ''}
              </Badge>
            )}
          </div>

          {/* Botón de limpiar filtros */}
          {activeFiltersCount > 0 && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleClearFilters}
              className="text-red-600 hover:text-red-700 hover:bg-red-50"
            >
              <XMarkIcon className="w-4 h-4 mr-2" />
              Limpiar filtros
            </Button>
          )}
        </div>

        {/* Filtros activos como badges */}
        {activeFiltersCount > 0 && (
          <div className="flex flex-wrap gap-2 pt-2">
            {Object.entries(localFilters).map(([key, value]) => {
              if (!value && value !== 0) return null;
              
              const field = fields.find(f => f.key === key);
              if (!field) return null;

              let displayValue = value;
              if (field.type === 'select' && field.options) {
                const option = field.options.find(opt => opt.value === value?.toString());
                displayValue = option?.label || value;
              }

              return (
                <Badge key={key} type="secondary" className="text-xs">
                  {field.label}: {displayValue}
                  <button
                    onClick={() => handleFilterChange(key, undefined)}
                    className="ml-2 hover:text-red-600"
                  >
                    <XMarkIcon className="w-3 h-3" />
                  </button>
                </Badge>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}

