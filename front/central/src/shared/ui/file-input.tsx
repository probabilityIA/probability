'use client';

import { InputHTMLAttributes, ReactNode, useRef } from 'react';
import { Button } from './button';

interface FileInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'onChange'> {
  label?: string;
  error?: string;
  helperText?: string;
  accept?: string;
  buttonText?: string;
  icon?: ReactNode;
  onChange?: (file: File | null) => void;
}

export function FileInput({
  label,
  error,
  helperText,
  accept = '*/*',
  buttonText = 'Seleccionar archivo',
  icon,
  className = '',
  onChange,
  ...props
}: FileInputProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    if (onChange) {
      onChange(file);
    }
  };

  const getFileName = () => {
    const files = fileInputRef.current?.files;
    if (files && files.length > 0) {
      return files[0].name;
    }
    return null;
  };

  return (
    <div className="space-y-2">
      {label && (
        <label className="block text-sm font-medium text-gray-700">
          {label}
        </label>
      )}
      
      <div className="flex items-center gap-3">
        <input
          ref={fileInputRef}
          type="file"
          accept={accept}
          className="hidden"
          onChange={handleFileChange}
          {...props}
        />
        
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={handleButtonClick}
          className="flex items-center gap-2"
        >
          {icon}
          {buttonText}
        </Button>
        
        {getFileName() && (
          <span className="text-sm text-gray-600 truncate max-w-xs">
            {getFileName()}
          </span>
        )}
      </div>
      
      {helperText && (
        <p className="text-sm text-gray-500">{helperText}</p>
      )}
      
      {error && (
        <p className="text-sm text-red-500">{error}</p>
      )}
    </div>
  );
}
