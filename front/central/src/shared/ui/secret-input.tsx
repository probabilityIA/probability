'use client';

import { InputHTMLAttributes, useId, useState } from 'react';
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';

export interface SecretInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
    label?: string;
    helperText?: string;
    error?: string;
    containerClassName?: string;
}

export function SecretInput({
    label,
    helperText,
    error,
    containerClassName,
    className,
    id,
    name,
    placeholder,
    ...rest
}: SecretInputProps) {
    const [show, setShow] = useState(false);
    const reactId = useId();
    const fieldId = id || `secret-${reactId.replace(/:/g, '')}`;
    const fieldName = name || fieldId;

    const maskStyle = show
        ? undefined
        : ({ WebkitTextSecurity: 'disc', textSecurity: 'disc' } as React.CSSProperties);

    return (
        <div className={containerClassName}>
            {label && (
                <label
                    htmlFor={fieldId}
                    className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1.5"
                >
                    {label}
                </label>
            )}
            <div className="relative">
                <input
                    type="text"
                    autoComplete="off"
                    autoCorrect="off"
                    autoCapitalize="off"
                    spellCheck={false}
                    data-1p-ignore="true"
                    data-lpignore="true"
                    data-bwignore="true"
                    data-form-type="other"
                    {...rest}
                    id={fieldId}
                    name={fieldName}
                    placeholder={placeholder}
                    style={{ ...maskStyle, ...(rest.style || {}) }}
                    className={
                        'w-full rounded-md border bg-white dark:bg-gray-800 text-gray-900 dark:text-white font-mono text-sm px-3 py-2 pr-10 focus:outline-none focus:ring-2 focus:ring-blue-500 ' +
                        (error
                            ? 'border-red-400 focus:ring-red-500 '
                            : 'border-gray-300 dark:border-gray-600 ') +
                        (className || '')
                    }
                />
                <button
                    type="button"
                    onClick={() => setShow((s) => !s)}
                    className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
                    tabIndex={-1}
                    aria-label={show ? 'Ocultar' : 'Mostrar'}
                >
                    {show ? <EyeSlashIcon className="w-5 h-5" /> : <EyeIcon className="w-5 h-5" />}
                </button>
            </div>
            {error ? (
                <p className="mt-1 text-xs text-red-600 dark:text-red-400">{error}</p>
            ) : helperText ? (
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">{helperText}</p>
            ) : null}
        </div>
    );
}
