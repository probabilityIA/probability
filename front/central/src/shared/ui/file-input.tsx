'use client';

import { InputHTMLAttributes, ReactNode, useRef, useState } from 'react';
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
  const [fileName, setFileName] = useState<string | null>(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setFileName(file?.name || null);
    if (onChange) {
      onChange(file);
    }
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
        
        {fileName && (
          <span className="text-sm text-gray-600 truncate max-w-xs">
            {fileName}
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
