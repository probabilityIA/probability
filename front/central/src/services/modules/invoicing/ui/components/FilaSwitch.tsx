'use client';

import type { ReactNode } from 'react';

interface FilaSwitchProps {
  titulo: string;
  descripcion: ReactNode;
  checked: boolean;
  onToggle: (valor: boolean) => void;
  disabled?: boolean;
  children?: ReactNode;
}

export function FilaSwitch({ titulo, descripcion, checked, onToggle, disabled, children }: FilaSwitchProps) {
  return (
    <div className={`py-3 ${disabled ? 'opacity-50' : ''}`}>
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0">
          <span className="text-sm font-medium text-gray-900 dark:text-white">{titulo}</span>
          <p className="mt-0.5 text-xs text-gray-500 dark:text-gray-400">{descripcion}</p>
        </div>
        <button
          type="button"
          role="switch"
          aria-checked={checked}
          aria-label={titulo}
          onClick={() => onToggle(!checked)}
          disabled={disabled}
          className={`relative mt-0.5 inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors focus:outline-none disabled:cursor-not-allowed ${checked ? '' : 'bg-gray-200 dark:bg-gray-600'}`}
          style={{ backgroundColor: checked ? 'var(--color-primary)' : undefined }}
        >
          <span
            className={`inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform ${checked ? 'translate-x-6' : 'translate-x-1'}`}
          />
        </button>
      </div>
      {children}
    </div>
  );
}
