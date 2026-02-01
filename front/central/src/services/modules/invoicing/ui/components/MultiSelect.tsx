'use client';

import { useState } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';

interface MultiSelectProps {
  options: Array<{ value: string | number; label: string }>;
  value: Array<string | number>;
  onChange: (value: Array<string | number>) => void;
  placeholder?: string;
  disabled?: boolean;
}

/**
 * Componente MultiSelect reutilizable
 */
export function MultiSelect({
  options,
  value,
  onChange,
  placeholder = 'Seleccione opciones...',
  disabled = false,
}: MultiSelectProps) {
  const [isOpen, setIsOpen] = useState(false);

  const handleToggle = (optionValue: string | number) => {
    if (value.includes(optionValue)) {
      onChange(value.filter((v) => v !== optionValue));
    } else {
      onChange([...value, optionValue]);
    }
  };

  const handleRemove = (optionValue: string | number) => {
    onChange(value.filter((v) => v !== optionValue));
  };

  const selectedOptions = options.filter((opt) => value.includes(opt.value));

  return (
    <div className="relative">
      {/* Selected tags */}
      {selectedOptions.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-2">
          {selectedOptions.map((option) => (
            <span
              key={option.value}
              className="inline-flex items-center gap-1 px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded"
            >
              {option.label}
              <button
                type="button"
                onClick={() => handleRemove(option.value)}
                disabled={disabled}
                className="hover:text-blue-600 disabled:opacity-50"
              >
                <XMarkIcon className="w-3 h-3" />
              </button>
            </span>
          ))}
        </div>
      )}

      {/* Dropdown toggle */}
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        disabled={disabled}
        className="w-full px-3 py-2 text-left border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {selectedOptions.length === 0 ? (
          <span className="text-gray-400">{placeholder}</span>
        ) : (
          <span className="text-gray-700">
            {selectedOptions.length} seleccionado(s)
          </span>
        )}
      </button>

      {/* Dropdown menu */}
      {isOpen && !disabled && (
        <>
          {/* Backdrop */}
          <div
            className="fixed inset-0 z-10"
            onClick={() => setIsOpen(false)}
          />

          {/* Options */}
          <div className="absolute z-20 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto">
            {options.map((option) => {
              const isSelected = value.includes(option.value);

              return (
                <div
                  key={option.value}
                  onClick={() => handleToggle(option.value)}
                  className={`px-3 py-2 cursor-pointer hover:bg-gray-100 ${
                    isSelected ? 'bg-blue-50 text-blue-800' : 'text-gray-700'
                  }`}
                >
                  <div className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      checked={isSelected}
                      onChange={() => {}}
                      className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <span>{option.label}</span>
                  </div>
                </div>
              );
            })}
          </div>
        </>
      )}
    </div>
  );
}
